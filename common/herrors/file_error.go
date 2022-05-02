// Copyright 2022 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable lfmtaw or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package herrors

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/text"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
	"github.com/tdewolff/parse/v2"

	"errors"
)

// FileError represents an error when handling a file: Parsing a config file,
// execute a template etc.
type FileError interface {
	error

	// ErroContext holds some context information about the error.
	ErrorContext() *ErrorContext

	text.Positioner

	// UpdatePosition updates the position of the error.
	UpdatePosition(pos text.Position) FileError

	// UpdateContent updates the error with a new ErrorContext from the content of the file.
	UpdateContent(r io.Reader, linematcher LineMatcherFn) FileError
}

// Unwrapper can unwrap errors created with fmt.Errorf.
type Unwrapper interface {
	Unwrap() error
}

var (
	_ FileError = (*fileError)(nil)
	_ Unwrapper = (*fileError)(nil)
)

func (fe *fileError) UpdatePosition(pos text.Position) FileError {
	oldFilename := fe.Position().Filename
	if pos.Filename != "" && fe.fileType == "" {
		_, fe.fileType = paths.FileAndExtNoDelimiter(filepath.Clean(pos.Filename))
	}
	if pos.Filename == "" {
		pos.Filename = oldFilename
	}
	fe.position = pos
	return fe
}

func (fe *fileError) UpdateContent(r io.Reader, linematcher LineMatcherFn) FileError {
	if linematcher == nil {
		linematcher = SimpleLineMatcher
	}

	var (
		contentPos   text.Position
		posle        = fe.position
		errorContext *ErrorContext
	)

	if posle.LineNumber <= 1 && posle.Offset > 0 {
		// Try to locate the line number from the content if offset is set.
		errorContext, contentPos = locateError(r, fe, func(m LineMatcher) bool {
			if posle.Offset >= m.Offset && posle.Offset < m.Offset+len(m.Line) {
				lno := posle.LineNumber - m.Position.LineNumber + m.LineNumber
				m.Position = text.Position{LineNumber: lno}
				return linematcher(m)
			}
			return false
		})
	} else {
		errorContext, contentPos = locateError(r, fe, linematcher)
	}

	if errorContext.ChromaLexer == "" {
		if fe.fileType != "" {
			errorContext.ChromaLexer = chromaLexerFromType(fe.fileType)
		} else {
			errorContext.ChromaLexer = chromaLexerFromFilename(fe.Position().Filename)
		}
	}

	fe.errorContext = errorContext

	if contentPos.LineNumber > 0 {
		fe.position.LineNumber = contentPos.LineNumber
	}

	return fe

}

type fileError struct {
	position     text.Position
	errorContext *ErrorContext

	fileType string

	cause error
}

type fileErrorWithErrorContext struct {
	*fileError
}

func (e *fileError) ErrorContext() *ErrorContext {
	return e.errorContext
}

// Position returns the text position of this error.
func (e fileError) Position() text.Position {
	return e.position
}

func (e *fileError) Error() string {
	return fmt.Sprintf("%s: %s", e.position, e.cause)
}

func (e *fileError) Unwrap() error {
	return e.cause
}

// NewFileError creates a new FileError that wraps err.
// The value for name should identify the file, the best
// being the full filename to the file on disk.
func NewFileError(name string, err error) FileError {
	if err == nil {
		panic("err is nil")
	}

	// Filetype is used to determine the Chroma lexer to use.
	fileType, pos := extractFileTypePos(err)
	pos.Filename = name
	if fileType == "" {
		_, fileType = paths.FileAndExtNoDelimiter(filepath.Clean(name))
	}

	if pos.LineNumber < 0 {
		panic(fmt.Sprintf("invalid line number: %d", pos.LineNumber))
	}

	return &fileError{cause: err, fileType: fileType, position: pos}

}

// NewFileErrorFromFile is a convenience method to create a new FileError from a file.
func NewFileErrorFromFile(err error, filename, realFilename string, fs afero.Fs, linematcher LineMatcherFn) FileError {
	if err == nil {
		panic("err is nil")
	}
	if linematcher == nil {
		linematcher = SimpleLineMatcher
	}
	f, err2 := fs.Open(filename)
	if err2 != nil {
		return NewFileError(realFilename, err)
	}
	defer f.Close()
	return NewFileError(realFilename, err).UpdateContent(f, linematcher)
}

// Cause returns the underlying error or itself if it does not implement Unwrap.
func Cause(err error) error {
	if u := errors.Unwrap(err); u != nil {
		return u
	}
	return err
}

func extractFileTypePos(err error) (string, text.Position) {
	err = Cause(err)
	var fileType string

	// Fall back to line/col 1:1 if we cannot find any better information.
	pos := text.Position{
		Offset:       -1,
		LineNumber:   1,
		ColumnNumber: 1,
	}

	// JSON errors.
	offset, typ := extractOffsetAndType(err)
	if fileType == "" {
		fileType = typ
	}

	if offset >= 0 {
		pos.Offset = offset
	}

	// The error type from the minifier contains line number and column number.
	if line, col := exctractLineNumberAndColumnNumber(err); line >= 0 {
		pos.LineNumber = line
		pos.ColumnNumber = col
		return fileType, pos
	}

	// Look in the error message for the line number.
	for _, handle := range lineNumberExtractors {
		lno, col := handle(err)
		if lno > 0 {
			pos.ColumnNumber = col
			pos.LineNumber = lno
			break
		}
	}

	return fileType, pos
}

// UnwrapFileError tries to unwrap a FileError from err.
// It returns nil if this is not possible.
func UnwrapFileError(err error) FileError {
	for err != nil {
		switch v := err.(type) {
		case FileError:
			return v
		default:
			err = errors.Unwrap(err)
		}
	}
	return nil
}

// UnwrapFileErrorsWithErrorContext tries to unwrap all FileError in err that has an ErrorContext.
func UnwrapFileErrorsWithErrorContext(err error) []FileError {
	var errs []FileError
	for err != nil {
		if v, ok := err.(FileError); ok && v.ErrorContext() != nil {
			errs = append(errs, v)
		}
		err = errors.Unwrap(err)
	}
	return errs
}

func extractOffsetAndType(e error) (int, string) {
	switch v := e.(type) {
	case *json.UnmarshalTypeError:
		return int(v.Offset), "json"
	case *json.SyntaxError:
		return int(v.Offset), "json"
	default:
		return -1, ""
	}
}

func exctractLineNumberAndColumnNumber(e error) (int, int) {
	switch v := e.(type) {
	case *parse.Error:
		return v.Line, v.Column
	case *toml.DecodeError:
		return v.Position()

	}

	return -1, -1
}
