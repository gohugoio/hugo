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

var _ causer = (*fileError)(nil)

// FileError represents an error when handling a file: Parsing a config file,
// execute a template etc.
type FileError interface {
	error

	// LineNumber gets the error location, starting at line 1.
	LineNumber() int

	ColumnNumber() int

	// A string identifying the type of file, e.g. JSON, TOML, markdown etc.
	Type() string
}

var _ FileError = (*fileError)(nil)

type fileError struct {
	lineNumber   int
	columnNumber int
	fileType     string

	cause error
}

func (e *fileError) LineNumber() int {
	return e.lineNumber
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
func NewFileError(fileType string, lineNumber, columnNumber int, err error) FileError {
	return &fileError{cause: err, fileType: fileType, lineNumber: lineNumber, columnNumber: columnNumber}
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

// ToFileError will try to convert the given error to an error supporting
// the FileError interface.
// If will fall back to returning the original error if a line number cannot be extracted.
func ToFileError(fileType string, err error) error {
	return ToFileErrorWithOffset(fileType, err, 0)
}

// ToFileErrorWithOffset will try to convert the given error to an error supporting
// the FileError interface. It will take any line number offset given into account.
// If will fall back to returning the original error if a line number cannot be extracted.
func ToFileErrorWithOffset(fileType string, err error, offset int) error {
	for _, handle := range lineNumberExtractors {

		lno, col := handle(err)
		if lno > 0 {
			return NewFileError(fileType, lno+offset, col, err)
		}
	}
	// Fall back to the original.
	return err
}
