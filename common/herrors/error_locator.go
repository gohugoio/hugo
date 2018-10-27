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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gohugoio/hugo/common/terminal"
	"github.com/gohugoio/hugo/helpers"

	"github.com/spf13/afero"
)

var fileErrorFormatFunc func(e ErrorContext) string

func createFileLogFormatter(formatStr string) func(e ErrorContext) string {

	if formatStr == "" {
		formatStr = "\":file::line::col\""
	}

	var identifiers = []string{":file", ":line", ":col"}
	var identifiersFound []string

	for i := range formatStr {
		for _, id := range identifiers {
			if strings.HasPrefix(formatStr[i:], id) {
				identifiersFound = append(identifiersFound, id)
			}
		}
	}

	replacer := strings.NewReplacer(":file", "%s", ":line", "%d", ":col", "%d")
	format := replacer.Replace(formatStr)

	f := func(e ErrorContext) string {
		args := make([]interface{}, len(identifiersFound))
		for i, id := range identifiersFound {
			switch id {
			case ":file":
				args[i] = e.Filename
			case ":line":
				args[i] = e.LineNumber
			case ":col":
				args[i] = e.ColumnNumber
			}
		}

		msg := fmt.Sprintf(format, args...)

		if terminal.IsTerminal(os.Stdout) {
			return terminal.Notice(msg)
		}

		return msg
	}

	return f
}

func init() {
	fileErrorFormatFunc = createFileLogFormatter(os.Getenv("HUGO_FILE_LOG_FORMAT"))
}

// LineMatcher contains the elements used to match an error to a line
type LineMatcher struct {
	FileError  FileError
	LineNumber int
	Offset     int
	Line       string
}

// LineMatcherFn is used to match a line with an error.
type LineMatcherFn func(m LineMatcher) bool

// SimpleLineMatcher simply matches by line number.
var SimpleLineMatcher = func(m LineMatcher) bool {
	return m.FileError.LineNumber() == m.LineNumber
}

// ErrorContext contains contextual information about an error. This will
// typically be the lines surrounding some problem in a file.
type ErrorContext struct {
	// The source filename.
	Filename string

	// If a match will contain the matched line and up to 2 lines before and after.
	// Will be empty if no match.
	Lines []string

	// The position of the error in the Lines above. 0 based.
	Pos int

	// The linenumber in the source file from where the Lines start. Starting at 1.
	LineNumber int

	// The column number in the source file. Starting at 1.
	ColumnNumber int

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
	return fileErrorFormatFunc(e.ErrorContext) + ": " + e.cause.Error()
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

	if le.Offset() != -1 {
		errCtx = locateError(r, le, func(m LineMatcher) bool {
			if le.Offset() >= m.Offset && le.Offset() < m.Offset+len(m.Line) {
				fe := m.FileError
				m.FileError = ToFileErrorWithOffset(fe, -fe.LineNumber()+m.LineNumber)
			}
			return matcher(m)
		})

	} else {
		errCtx = locateError(r, le, matcher)
	}

	if errCtx.LineNumber == -1 {
		return e, false
	}

	errCtx.Filename = realFilename

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

func chromaLexerFromFilename(filename string) string {
	if strings.Contains(filename, "layouts") {
		return "go-html-template"
	}

	ext := helpers.ExtNoDelimiter(filename)
	return chromaLexerFromType(ext)
}

func locateErrorInString(src string, matcher LineMatcherFn) ErrorContext {
	return locateError(strings.NewReader(src), &fileError{}, matcher)
}

func locateError(r io.Reader, le FileError, matches LineMatcherFn) ErrorContext {
	if le == nil {
		panic("must provide an error")
	}

	errCtx := ErrorContext{LineNumber: -1, ColumnNumber: 1, Pos: -1}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errCtx
	}

	lines := strings.Split(string(b), "\n")

	if le != nil && le.ColumnNumber() >= 0 {
		errCtx.ColumnNumber = le.ColumnNumber()
	}

	lineNo := 0
	posBytes := 0

	for li, line := range lines {
		lineNo = li + 1
		m := LineMatcher{
			FileError:  le,
			LineNumber: lineNo,
			Offset:     posBytes,
			Line:       line,
		}
		if errCtx.Pos == -1 && matches(m) {
			errCtx.LineNumber = lineNo
			break
		}

		posBytes += len(line)
	}

	if errCtx.LineNumber != -1 {
		low := errCtx.LineNumber - 3
		if low < 0 {
			low = 0
		}

		if errCtx.LineNumber > 2 {
			errCtx.Pos = 2
		} else {
			errCtx.Pos = errCtx.LineNumber - 1
		}

		high := errCtx.LineNumber + 2
		if high > len(lines) {
			high = len(lines)
		}

		errCtx.Lines = lines[low:high]

	}

	return errCtx
}
