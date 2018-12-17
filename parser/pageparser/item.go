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
)

type Item struct {
	Type ItemType
	Pos  int
	Val  []byte
}

type Items []Item

func (i Item) ValStr() string {
	return string(i.Val)
}

func (i Item) IsText() bool {
	return i.Type == tText
}

func (i Item) IsNonWhitespace() bool {
	return len(bytes.TrimSpace(i.Val)) > 0
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

func (i Item) String() string {
	switch {
	case i.Type == tEOF:
		return "EOF"
	case i.Type == tError:
		return string(i.Val)
	case i.Type > tKeywordMarker:
		return fmt.Sprintf("<%s>", i.Val)
	case len(i.Val) > 50:
		return fmt.Sprintf("%v:%.20q...", i.Type, i.Val)
	}
	return fmt.Sprintf("%v:[%s]", i.Type, i.Val)
}

type ItemType int

const (
	tError ItemType = iota
	tEOF

	// page items
	TypeHTMLStart          // document starting with < as first non-whitespace
	TypeLeadSummaryDivider // <!--more-->,  # more
	TypeFrontMatterYAML
	TypeFrontMatterTOML
	TypeFrontMatterJSON
	TypeFrontMatterORG
	TypeEmoji
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

	tText // plain text

	// preserved for later - keywords come after this
	tKeywordMarker
)
