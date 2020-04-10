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

// Package herrors contains common Hugo errors and error related utilities.
package herrors

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/text"

	"github.com/spf13/afero"
)

// LineMatcher contains the elements used to match an error to a line
type LineMatcher struct {
	Position text.Position
	Error    error

	LineNumber int
	Offset     int
	Line       string
}

// LineMatcherFn is used to match a line with an error.
type LineMatcherFn func(m LineMatcher) bool

// SimpleLineMatcher simply matches by line number.
var SimpleLineMatcher = func(m LineMatcher) bool {
	return m.Position.LineNumber == m.LineNumber
}

var _ text.Positioner = ErrorContext{}

// ErrorContext contains contextual information about an error. This will
// typically be the lines surrounding some problem in a file.
type ErrorContext struct {

	// If a match will contain the matched line and up to 2 lines before and after.
	// Will be empty if no match.
	Lines []string

	// The position of the error in the Lines above. 0 based.
	LinesPos int

	position text.Position

	// The lexer to use for syntax highlighting.
	// https://gohugo.io/content-management/syntax-highlighting/#list-of-chroma-highlighting-languages
	ChromaLexer string
}

// Position returns the text position of this error.
func (e ErrorContext) Position() text.Position {
	return e.position
}

var _ causer = (*ErrorWithFileContext)(nil)

// ErrorWithFileContext is an error with some additional file context related
// to that error.
type ErrorWithFileContext struct {
	cause error
	ErrorContext
}

func (e *ErrorWithFileContext) Error() string {
	pos := e.Position()
	if pos.IsValid() {
		return pos.String() + ": " + e.cause.Error()
	}
	return e.cause.Error()
}

func (e *ErrorWithFileContext) Cause() error {
	return e.cause
}

// WithFileContextForFile will try to add a file context with lines matching the given matcher.
// If no match could be found, the original error is returned with false as the second return value.
func WithFileContextForFile(e error, realFilename, filename string, fs afero.Fs, matcher LineMatcherFn) (error, bool) {
	f, err := fs.Open(filename)
	if err != nil {
		return e, false
	}
	defer f.Close()
	return WithFileContext(e, realFilename, f, matcher)
}

// WithFileContextForFile will try to add a file context with lines matching the given matcher.
// If no match could be found, the original error is returned with false as the second return value.
func WithFileContext(e error, realFilename string, r io.Reader, matcher LineMatcherFn) (error, bool) {
	if e == nil {
		panic("error missing")
	}
	le := UnwrapFileError(e)

	if le == nil {
		var ok bool
		if le, ok = ToFileError("", e).(FileError); !ok {
			return e, false
		}
	}

	var errCtx ErrorContext

	posle := le.Position()

	if posle.Offset != -1 {
		errCtx = locateError(r, le, func(m LineMatcher) bool {
			if posle.Offset >= m.Offset && posle.Offset < m.Offset+len(m.Line) {
				lno := posle.LineNumber - m.Position.LineNumber + m.LineNumber
				m.Position = text.Position{LineNumber: lno}
			}
			return matcher(m)
		})
	} else {
		errCtx = locateError(r, le, matcher)
	}

	pos := &errCtx.position

	if pos.LineNumber == -1 {
		return e, false
	}

	pos.Filename = realFilename

	if le.Type() != "" {
		errCtx.ChromaLexer = chromaLexerFromType(le.Type())
	} else {
		errCtx.ChromaLexer = chromaLexerFromFilename(realFilename)
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
	switch fileType {
	case "html", "htm":
		return "go-html-template"
	}
	return fileType
}

func extNoDelimiter(filename string) string {
	return strings.TrimPrefix(filepath.Ext(filename), ".")
}

func chromaLexerFromFilename(filename string) string {
	if strings.Contains(filename, "layouts") {
		return "go-html-template"
	}

	ext := extNoDelimiter(filename)
	return chromaLexerFromType(ext)
}

func locateErrorInString(src string, matcher LineMatcherFn) ErrorContext {
	return locateError(strings.NewReader(src), &fileError{}, matcher)
}

func locateError(r io.Reader, le FileError, matches LineMatcherFn) ErrorContext {
	if le == nil {
		panic("must provide an error")
	}

	errCtx := ErrorContext{position: text.Position{LineNumber: -1, ColumnNumber: 1, Offset: -1}, LinesPos: -1}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errCtx
	}

	pos := &errCtx.position
	lepos := le.Position()

	lines := strings.Split(string(b), "\n")

	if lepos.ColumnNumber >= 0 {
		pos.ColumnNumber = lepos.ColumnNumber
	}

	lineNo := 0
	posBytes := 0

	for li, line := range lines {
		lineNo = li + 1
		m := LineMatcher{
			Position:   le.Position(),
			Error:      le,
			LineNumber: lineNo,
			Offset:     posBytes,
			Line:       line,
		}
		if errCtx.LinesPos == -1 && matches(m) {
			pos.LineNumber = lineNo
			break
		}

		posBytes += len(line)
	}

	if pos.LineNumber != -1 {
		low := pos.LineNumber - 3
		if low < 0 {
			low = 0
		}

		if pos.LineNumber > 2 {
			errCtx.LinesPos = 2
		} else {
			errCtx.LinesPos = pos.LineNumber - 1
		}

		high := pos.LineNumber + 2
		if high > len(lines) {
			high = len(lines)
		}

		errCtx.Lines = lines[low:high]

	}

	return errCtx
}
