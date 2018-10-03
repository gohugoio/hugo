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

// Package errors contains common Hugo errors and error related utilities.
package herrors

import (
	"bufio"
	"io"
	"strings"

	"github.com/spf13/afero"
)

// LineMatcher is used to match a line with an error.
type LineMatcher func(le FileError, lineNumber int, line string) bool

// SimpleLineMatcher matches if the current line number matches the line number
// in the error.
var SimpleLineMatcher = func(le FileError, lineNumber int, line string) bool {
	return le.LineNumber() == lineNumber
}

// ErrorContext contains contextual information about an error. This will
// typically be the lines surrounding some problem in a file.
type ErrorContext struct {

	// If a match will contain the matched line and up to 2 lines before and after.
	// Will be empty if no match.
	Lines []string

	// The position of the error in the Lines above. 0 based.
	Pos int

	// The linenumber in the source file from where the Lines start. Starting at 1.
	LineNumber int

	// The lexer to use for syntax highlighting.
	// https://gohugo.io/content-management/syntax-highlighting/#list-of-chroma-highlighting-languages
	ChromaLexer string
}

var _ causer = (*ErrorWithFileContext)(nil)

// ErrorWithFileContext is an error with some additional file context related
// to that error.
type ErrorWithFileContext struct {
	cause error
	ErrorContext
}

func (e *ErrorWithFileContext) Error() string {
	return e.cause.Error()
}

func (e *ErrorWithFileContext) Cause() error {
	return e.cause
}

// WithFileContextForFile will try to add a file context with lines matching the given matcher.
// If no match could be found, the original error is returned with false as the second return value.
func WithFileContextForFile(e error, filename string, fs afero.Fs, chromaLexer string, matcher LineMatcher) (error, bool) {
	f, err := fs.Open(filename)
	if err != nil {
		return e, false
	}
	defer f.Close()
	return WithFileContext(e, f, chromaLexer, matcher)
}

// WithFileContextForFile will try to add a file context with lines matching the given matcher.
// If no match could be found, the original error is returned with false as the second return value.
func WithFileContext(e error, r io.Reader, chromaLexer string, matcher LineMatcher) (error, bool) {
	if e == nil {
		panic("error missing")
	}
	le := UnwrapFileError(e)
	if le == nil {
		var ok bool
		if le, ok = ToFileError("bash", e).(FileError); !ok {
			return e, false
		}
	}

	errCtx := locateError(r, le, matcher)

	if errCtx.LineNumber == -1 {
		return e, false
	}

	if chromaLexer != "" {
		errCtx.ChromaLexer = chromaLexer
	} else {
		errCtx.ChromaLexer = chromaLexerFromType(le.Type())
	}

	return &ErrorWithFileContext{cause: e, ErrorContext: errCtx}, true
}

// UnwrapErrorWithFileContext tries to unwrap an ErrorWithFileContext from err.
// It returns nil if this is not possible.
func UnwrapErrorWithFileContext(err error) *ErrorWithFileContext {
	for err != nil {
		switch v := err.(type) {
		case *ErrorWithFileContext:
			return v
		case causer:
			err = v.Cause()
		default:
			return nil
		}
	}
	return nil
}

func chromaLexerFromType(fileType string) string {
	return fileType
}

func locateErrorInString(le FileError, src string, matcher LineMatcher) ErrorContext {
	return locateError(strings.NewReader(src), nil, matcher)
}

func locateError(r io.Reader, le FileError, matches LineMatcher) ErrorContext {
	var errCtx ErrorContext
	s := bufio.NewScanner(r)

	lineNo := 0

	var buff [6]string
	i := 0
	errCtx.Pos = -1

	for s.Scan() {
		lineNo++
		txt := s.Text()
		buff[i] = txt

		if errCtx.Pos != -1 && i >= 5 {
			break
		}

		if errCtx.Pos == -1 && matches(le, lineNo, txt) {
			errCtx.Pos = i
			errCtx.LineNumber = lineNo - i
		}

		if errCtx.Pos == -1 && i == 2 {
			// Shift left
			buff[0], buff[1] = buff[i-1], buff[i]
		} else {
			i++
		}
	}

	// Go's template parser will typically report "unexpected EOF" errors on the
	// empty last line that is supressed by the scanner.
	// Do an explicit check for that.
	if errCtx.Pos == -1 {
		lineNo++
		if matches(le, lineNo, "") {
			buff[i] = ""
			errCtx.Pos = i
			errCtx.LineNumber = lineNo - 1

			i++
		}
	}

	if errCtx.Pos != -1 {
		low := errCtx.Pos - 2
		if low < 0 {
			low = 0
		}
		high := i
		errCtx.Lines = buff[low:high]

	} else {
		errCtx.Pos = -1
		errCtx.LineNumber = -1
	}

	return errCtx
}
