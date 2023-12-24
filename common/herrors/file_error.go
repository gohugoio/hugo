// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"errors"
	"fmt"
	"io"
	"path/filepath"

	godartsassv1 "github.com/bep/godartsass"

	"github.com/bep/godartsass/v2"
	"github.com/bep/golibsass/libsass/libsasserrors"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/text"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
	"github.com/tdewolff/parse/v2"
)

// FileError represents an error when handling a file: Parsing a config file,
// execute a template etc.
type FileError interface {
	error

	// ErrorContext holds some context information about the error.
	ErrorContext() *ErrorContext

	text.Positioner

	// UpdatePosition updates the position of the error.
	UpdatePosition(pos text.Position) FileError

	// UpdateContent updates the error with a new ErrorContext from the content of the file.
	UpdateContent(r io.Reader, linematcher LineMatcherFn) FileError

	// SetFilename sets the filename of the error.
	SetFilename(filename string) FileError
}

// Unwrapper can unwrap errors created with fmt.Errorf.
type Unwrapper interface {
	Unwrap() error
}

var (
	_ FileError = (*fileError)(nil)
	_ Unwrapper = (*fileError)(nil)
)

func (fe *fileError) SetFilename(filename string) FileError {
	fe.position.Filename = filename
	return fe
}

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
		posle = fe.position
		ectx  *ErrorContext
	)

	if posle.LineNumber <= 1 && posle.Offset > 0 {
		// Try to locate the line number from the content if offset is set.
		ectx = locateError(r, fe, func(m LineMatcher) int {
			if posle.Offset >= m.Offset && posle.Offset < m.Offset+len(m.Line) {
				lno := posle.LineNumber - m.Position.LineNumber + m.LineNumber
				m.Position = text.Position{LineNumber: lno}
				return linematcher(m)
			}
			return -1
		})
	} else {
		ectx = locateError(r, fe, linematcher)
	}

	if ectx.ChromaLexer == "" {
		if fe.fileType != "" {
			ectx.ChromaLexer = chromaLexerFromType(fe.fileType)
		} else {
			ectx.ChromaLexer = chromaLexerFromFilename(fe.Position().Filename)
		}
	}

	fe.errorContext = ectx

	if ectx.Position.LineNumber > 0 {
		fe.position.LineNumber = ectx.Position.LineNumber
	}

	if ectx.Position.ColumnNumber > 0 {
		fe.position.ColumnNumber = ectx.Position.ColumnNumber
	}

	return fe
}

type fileError struct {
	position     text.Position
	errorContext *ErrorContext

	fileType string

	cause error
}

func (e *fileError) ErrorContext() *ErrorContext {
	return e.errorContext
}

// Position returns the text position of this error.
func (e fileError) Position() text.Position {
	return e.position
}

func (e *fileError) Error() string {
	return fmt.Sprintf("%s: %s", e.position, e.causeString())
}

func (e *fileError) causeString() string {
	if e.cause == nil {
		return ""
	}
	switch v := e.cause.(type) {
	// Avoid repeating the file info in the error message.
	case godartsass.SassError:
		return v.Message
	case godartsassv1.SassError:
		return v.Message
	case libsasserrors.Error:
		return v.Message
	default:
		return v.Error()
	}
}

func (e *fileError) Unwrap() error {
	return e.cause
}

// NewFileError creates a new FileError that wraps err.
// It will try to extract the filename and line number from err.
func NewFileError(err error) FileError {
	// Filetype is used to determine the Chroma lexer to use.
	fileType, pos := extractFileTypePos(err)
	return &fileError{cause: err, fileType: fileType, position: pos}
}

// NewFileErrorFromName creates a new FileError that wraps err.
// The value for name should identify the file, the best
// being the full filename to the file on disk.
func NewFileErrorFromName(err error, name string) FileError {
	// Filetype is used to determine the Chroma lexer to use.
	fileType, pos := extractFileTypePos(err)
	pos.Filename = name
	if fileType == "" {
		_, fileType = paths.FileAndExtNoDelimiter(filepath.Clean(name))
	}

	return &fileError{cause: err, fileType: fileType, position: pos}
}

// NewFileErrorFromPos will use the filename and line number from pos to create a new FileError, wrapping err.
func NewFileErrorFromPos(err error, pos text.Position) FileError {
	// Filetype is used to determine the Chroma lexer to use.
	fileType, _ := extractFileTypePos(err)
	if fileType == "" {
		_, fileType = paths.FileAndExtNoDelimiter(filepath.Clean(pos.Filename))
	}
	return &fileError{cause: err, fileType: fileType, position: pos}
}

func NewFileErrorFromFileInErr(err error, fs afero.Fs, linematcher LineMatcherFn) FileError {
	fe := NewFileError(err)
	pos := fe.Position()
	if pos.Filename == "" {
		return fe
	}

	f, realFilename, err2 := openFile(pos.Filename, fs)
	if err2 != nil {
		return fe
	}

	pos.Filename = realFilename
	defer f.Close()
	return fe.UpdateContent(f, linematcher)
}

func NewFileErrorFromFileInPos(err error, pos text.Position, fs afero.Fs, linematcher LineMatcherFn) FileError {
	if err == nil {
		panic("err is nil")
	}
	f, realFilename, err2 := openFile(pos.Filename, fs)
	if err2 != nil {
		return NewFileErrorFromPos(err, pos)
	}
	pos.Filename = realFilename
	defer f.Close()
	return NewFileErrorFromPos(err, pos).UpdateContent(f, linematcher)
}

// NewFileErrorFromFile is a convenience method to create a new FileError from a file.
func NewFileErrorFromFile(err error, filename string, fs afero.Fs, linematcher LineMatcherFn) FileError {
	if err == nil {
		panic("err is nil")
	}
	f, realFilename, err2 := openFile(filename, fs)
	if err2 != nil {
		return NewFileErrorFromName(err, realFilename)
	}
	defer f.Close()
	return NewFileErrorFromName(err, realFilename).UpdateContent(f, linematcher)
}

func openFile(filename string, fs afero.Fs) (afero.File, string, error) {
	realFilename := filename

	// We want the most specific filename possible in the error message.
	fi, err2 := fs.Stat(filename)
	if err2 == nil {
		if s, ok := fi.(interface {
			Filename() string
		}); ok {
			realFilename = s.Filename()
		}
	}

	f, err2 := fs.Open(filename)
	if err2 != nil {
		return nil, realFilename, err2
	}

	return f, realFilename, nil
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

	// LibSass, DartSass
	if pos := extractPosition(err); pos.LineNumber > 0 || pos.Offset > 0 {
		_, fileType = paths.FileAndExtNoDelimiter(pos.Filename)
		return fileType, pos
	}

	// Default to line 1 col 1 if we don't find any better.
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
	if line, col := extractLineNumberAndColumnNumber(err); line >= 0 {
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

	if fileType == "" && pos.Filename != "" {
		_, fileType = paths.FileAndExtNoDelimiter(pos.Filename)
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

// UnwrapFileErrors tries to unwrap all FileError.
func UnwrapFileErrors(err error) []FileError {
	var errs []FileError
	for err != nil {
		if v, ok := err.(FileError); ok {
			errs = append(errs, v)
		}
		err = errors.Unwrap(err)
	}
	return errs
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

func extractLineNumberAndColumnNumber(e error) (int, int) {
	switch v := e.(type) {
	case *parse.Error:
		return v.Line, v.Column
	case *toml.DecodeError:
		return v.Position()

	}

	return -1, -1
}

func extractPosition(e error) (pos text.Position) {
	switch v := e.(type) {
	case godartsass.SassError:
		span := v.Span
		start := span.Start
		filename, _ := paths.UrlToFilename(span.Url)
		pos.Filename = filename
		pos.Offset = start.Offset
		pos.ColumnNumber = start.Column
	case godartsassv1.SassError:
		span := v.Span
		start := span.Start
		filename, _ := paths.UrlToFilename(span.Url)
		pos.Filename = filename
		pos.Offset = start.Offset
		pos.ColumnNumber = start.Column
	case libsasserrors.Error:
		pos.Filename = v.File
		pos.LineNumber = v.Line
		pos.ColumnNumber = v.Column
	}
	return
}

// TextSegmentError is an error with a text segment attached.
type TextSegmentError struct {
	Segment string
	Err     error
}

func (e TextSegmentError) Unwrap() error {
	return e.Err
}

func (e TextSegmentError) Error() string {
	return e.Err.Error()
}
