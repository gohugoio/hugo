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

package pageparser

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/yuin/goldmark/util"
)

type lowHigh struct {
	Low  int
	High int
}

type Item struct {
	Type ItemType
	Err  error

	// The common case is a single segment.
	low  int
	high int

	// This is the uncommon case.
	segments []lowHigh

	// Used for validation.
	firstByte byte

	isString bool
}

type Items []Item

func (i Item) Pos() int {
	if len(i.segments) > 0 {
		return i.segments[0].Low
	}
	return i.low
}

func (i Item) Val(source []byte) []byte {
	if len(i.segments) == 0 {
		return source[i.low:i.high]
	}

	if len(i.segments) == 1 {
		return source[i.segments[0].Low:i.segments[0].High]
	}

	var b bytes.Buffer
	for _, s := range i.segments {
		b.Write(source[s.Low:s.High])
	}
	return b.Bytes()
}

func (i Item) ValStr(source []byte) string {
	return string(i.Val(source))
}

func (i Item) ValTyped(source []byte) any {
	str := i.ValStr(source)
	if i.isString {
		// A quoted value that is a string even if it looks like a number etc.
		return str
	}

	if boolRe.MatchString(str) {
		return str == "true"
	}

	if intRe.MatchString(str) {
		num, err := strconv.Atoi(str)
		if err != nil {
			return str
		}
		return num
	}

	if floatRe.MatchString(str) {
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return str
		}
		return num
	}

	return str
}

func (i Item) IsText() bool {
	return i.Type == tText || i.Type == tIndentation
}

func (i Item) IsIndentation() bool {
	return i.Type == tIndentation
}

func (i Item) IsNonWhitespace(source []byte) bool {
	return len(bytes.TrimSpace(i.Val(source))) > 0
}

func (i Item) IsShortcodeName() bool {
	return i.Type == tScName
}

func (i Item) IsInlineShortcodeName() bool {
	return i.Type == tScNameInline
}

func (i Item) IsLeftShortcodeDelim() bool {
	return i.Type == tLeftDelimScWithMarkup || i.Type == tLeftDelimScNoMarkup
}

func (i Item) IsRightShortcodeDelim() bool {
	return i.Type == tRightDelimScWithMarkup || i.Type == tRightDelimScNoMarkup
}

func (i Item) IsShortcodeClose() bool {
	return i.Type == tScClose
}

func (i Item) IsShortcodeParam() bool {
	return i.Type == tScParam
}

func (i Item) IsShortcodeParamVal() bool {
	return i.Type == tScParamVal
}

func (i Item) IsShortcodeMarkupDelimiter() bool {
	return i.Type == tLeftDelimScWithMarkup || i.Type == tRightDelimScWithMarkup
}

func (i Item) IsFrontMatter() bool {
	return i.Type >= TypeFrontMatterYAML && i.Type <= TypeFrontMatterORG
}

func (i Item) IsDone() bool {
	return i.Type == tError || i.Type == tEOF
}

func (i Item) IsEOF() bool {
	return i.Type == tEOF
}

func (i Item) IsError() bool {
	return i.Type == tError
}

func (i Item) ToString(source []byte) string {
	val := i.Val(source)
	switch {
	case i.Type == tEOF:
		return "EOF"
	case i.Type == tError:
		return string(val)
	case i.Type == tIndentation:
		return fmt.Sprintf("%s:[%s]", i.Type, util.VisualizeSpaces(val))
	case i.Type > tKeywordMarker:
		return fmt.Sprintf("<%s>", val)
	case len(val) > 50:
		return fmt.Sprintf("%v:%.20q...", i.Type, val)
	}
	return fmt.Sprintf("%v:[%s]", i.Type, val)
}

type ItemType int

const (
	tError ItemType = iota
	tEOF

	// page items
	TypeLeadSummaryDivider // <!--more-->,  # more
	TypeFrontMatterYAML
	TypeFrontMatterTOML
	TypeFrontMatterJSON
	TypeFrontMatterORG
	TypeIgnore // // The BOM Unicode byte order marker and possibly others

	// shortcode items
	tLeftDelimScNoMarkup
	tRightDelimScNoMarkup
	tLeftDelimScWithMarkup
	tRightDelimScWithMarkup
	tScClose
	tScName
	tScNameInline
	tScParam
	tScParamVal

	tIndentation

	tText // plain text

	// preserved for later - keywords come after this
	tKeywordMarker
)

var (
	boolRe  = regexp.MustCompile(`^(true|false)$`)
	intRe   = regexp.MustCompile(`^[-+]?\d+$`)
	floatRe = regexp.MustCompile(`^[-+]?\d*\.\d+$`)
)
