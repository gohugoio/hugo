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
// limitatio	ns under the License.

package herrors

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var _ causer = (*fileError)(nil)

// FileError represents an error when handling a file: Parsing a config file,
// execute a template etc.
type FileError interface {
	error

	// Offset gets the error location offset in bytes, starting at 0.
	// It will return -1 if not provided.
	Offset() int

	// LineNumber gets the error location, starting at line 1.
	LineNumber() int

	// Column number gets the column location, starting at 1.
	ColumnNumber() int

	// A string identifying the type of file, e.g. JSON, TOML, markdown etc.
	Type() string
}

var _ FileError = (*fileError)(nil)

type fileError struct {
	offset       int
	lineNumber   int
	columnNumber int
	fileType     string

	cause error
}

type fileErrorWithLineOffset struct {
	FileError
	offset int
}

func (e *fileErrorWithLineOffset) LineNumber() int {
	return e.FileError.LineNumber() + e.offset
}

func (e *fileError) LineNumber() int {
	return e.lineNumber
}

func (e *fileError) Offset() int {
	return e.offset
}

func (e *fileError) ColumnNumber() int {
	return e.columnNumber
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
	return &fileError{cause: err, fileType: fileType, offset: offset, lineNumber: lineNumber, columnNumber: columnNumber}
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
	return &fileErrorWithLineOffset{FileError: fe, offset: offset}
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
