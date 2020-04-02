// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package herrors

import (
	"encoding/json"

	"github.com/gohugoio/hugo/common/text"

	"github.com/pkg/errors"
)

var (
	_ causer = (*fileError)(nil)
)

// FileError represents an error when handling a file: Parsing a config file,
// execute a template etc.
type FileError interface {
	error

	text.Positioner

	// A string identifying the type of file, e.g. JSON, TOML, markdown etc.
	Type() string
}

var _ FileError = (*fileError)(nil)

type fileError struct {
	position text.Position

	fileType string

	cause error
}

// Position returns the text position of this error.
func (e fileError) Position() text.Position {
	return e.position
}

func (e *fileError) Type() string {
	return e.fileType
}

func (e *fileError) Error() string {
	if e.cause == nil {
		return ""
	}
	return e.cause.Error()
}

func (f *fileError) Cause() error {
	return f.cause
}

// NewFileError creates a new FileError.
func NewFileError(fileType string, offset, lineNumber, columnNumber int, err error) FileError {
	pos := text.Position{Offset: offset, LineNumber: lineNumber, ColumnNumber: columnNumber}
	return &fileError{cause: err, fileType: fileType, position: pos}
}

// UnwrapFileError tries to unwrap a FileError from err.
// It returns nil if this is not possible.
func UnwrapFileError(err error) FileError {
	for err != nil {
		switch v := err.(type) {
		case FileError:
			return v
		case causer:
			err = v.Cause()
		default:
			return nil
		}
	}
	return nil
}

// ToFileErrorWithOffset will return a new FileError with a line number
// with the given offset from the original.
func ToFileErrorWithOffset(fe FileError, offset int) FileError {
	pos := fe.Position()
	return ToFileErrorWithLineNumber(fe, pos.LineNumber+offset)
}

// ToFileErrorWithOffset will return a new FileError with the given line number.
func ToFileErrorWithLineNumber(fe FileError, lineNumber int) FileError {
	pos := fe.Position()
	pos.LineNumber = lineNumber
	return &fileError{cause: fe, fileType: fe.Type(), position: pos}
}

// ToFileError will convert the given error to an error supporting
// the FileError interface.
func ToFileError(fileType string, err error) FileError {
	for _, handle := range lineNumberExtractors {
		lno, col := handle(err)
		offset, typ := extractOffsetAndType(err)
		if fileType == "" {
			fileType = typ
		}

		if lno > 0 || offset != -1 {
			return NewFileError(fileType, offset, lno, col, err)
		}
	}
	// Fall back to the pointing to line number 1.
	return NewFileError(fileType, -1, 1, 1, err)
}

func extractOffsetAndType(e error) (int, string) {
	e = errors.Cause(e)
	switch v := e.(type) {
	case *json.UnmarshalTypeError:
		return int(v.Offset), "json"
	case *json.SyntaxError:
		return int(v.Offset), "json"
	default:
		return -1, ""
	}
}
