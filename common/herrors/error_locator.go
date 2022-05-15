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
// It returns the column number or 0 if the line was found, but column could not be determinde. Returns -1 if no line match.
type LineMatcherFn func(m LineMatcher) int

// SimpleLineMatcher simply matches by line number.
var SimpleLineMatcher = func(m LineMatcher) int {
	if m.Position.LineNumber == m.LineNumber {
		// We found the line, but don't know the column.
		return 0
	}
	return -1
}

// NopLineMatcher is a matcher that always returns 1.
// This will effectively give line 1, column 1.
var NopLineMatcher = func(m LineMatcher) int {
	return 1
}

// OffsetMatcher is a line matcher that matches by offset.
var OffsetMatcher = func(m LineMatcher) int {
	if m.Offset+len(m.Line) >= m.Position.Offset {
		// We found the line, but return 0 to signal that we want to determine
		// the column from the error.
		return 0
	}
	return -1
}

// ErrorContext contains contextual information about an error. This will
// typically be the lines surrounding some problem in a file.
type ErrorContext struct {

	// If a match will contain the matched line and up to 2 lines before and after.
	// Will be empty if no match.
	Lines []string

	// The position of the error in the Lines above. 0 based.
	LinesPos int

	// The position of the content in the file. Note that this may be different from the error's position set
	// in FileError.
	Position text.Position

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

func locateErrorInString(src string, matcher LineMatcherFn) *ErrorContext {
	return locateError(strings.NewReader(src), &fileError{}, matcher)
}

func locateError(r io.Reader, le FileError, matches LineMatcherFn) *ErrorContext {
	if le == nil {
		panic("must provide an error")
	}

	ectx := &ErrorContext{LinesPos: -1, Position: text.Position{Offset: -1}}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return ectx
	}

	lines := strings.Split(string(b), "\n")

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
		v := matches(m)
		if ectx.LinesPos == -1 && v != -1 {
			ectx.Position.LineNumber = lineNo
			ectx.Position.ColumnNumber = v
			break
		}

		posBytes += len(line)
	}

	if ectx.Position.LineNumber > 0 {
		low := ectx.Position.LineNumber - 3
		if low < 0 {
			low = 0
		}

		if ectx.Position.LineNumber > 2 {
			ectx.LinesPos = 2
		} else {
			ectx.LinesPos = ectx.Position.LineNumber - 1
		}

		high := ectx.Position.LineNumber + 2
		if high > len(lines) {
			high = len(lines)
		}

		ectx.Lines = lines[low:high]

	}

	return ectx
}
