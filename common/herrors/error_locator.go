// Copyright 2022 The Hugo Authors. All rights reserved.
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

// ErrorContext contains contextual information about an error. This will
// typically be the lines surrounding some problem in a file.
type ErrorContext struct {

	// If a match will contain the matched line and up to 2 lines before and after.
	// Will be empty if no match.
	Lines []string

	// The position of the error in the Lines above. 0 based.
	LinesPos int

	// The lexer to use for syntax highlighting.
	// https://gohugo.io/content-management/syntax-highlighting/#list-of-chroma-highlighting-languages
	ChromaLexer string
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

func locateErrorInString(src string, matcher LineMatcherFn) (*ErrorContext, text.Position) {
	return locateError(strings.NewReader(src), &fileError{}, matcher)
}

func locateError(r io.Reader, le FileError, matches LineMatcherFn) (*ErrorContext, text.Position) {
	if le == nil {
		panic("must provide an error")
	}

	errCtx := &ErrorContext{LinesPos: -1}
	pos := text.Position{LineNumber: -1, ColumnNumber: 1, Offset: -1}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errCtx, pos
	}

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

	return errCtx, pos
}
