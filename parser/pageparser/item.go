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

import "fmt"

type Item struct {
	Typ ItemType
	pos pos
	Val []byte
}

type Items []Item

func (i Item) ValStr() string {
	return string(i.Val)
}

func (i Item) IsText() bool {
	return i.Typ == tText
}

func (i Item) IsShortcodeName() bool {
	return i.Typ == tScName
}

func (i Item) IsLeftShortcodeDelim() bool {
	return i.Typ == tLeftDelimScWithMarkup || i.Typ == tLeftDelimScNoMarkup
}

func (i Item) IsRightShortcodeDelim() bool {
	return i.Typ == tRightDelimScWithMarkup || i.Typ == tRightDelimScNoMarkup
}

func (i Item) IsShortcodeClose() bool {
	return i.Typ == tScClose
}

func (i Item) IsShortcodeParam() bool {
	return i.Typ == tScParam
}

func (i Item) IsShortcodeParamVal() bool {
	return i.Typ == tScParamVal
}

func (i Item) IsShortcodeMarkupDelimiter() bool {
	return i.Typ == tLeftDelimScWithMarkup || i.Typ == tRightDelimScWithMarkup
}

func (i Item) IsFrontMatter() bool {
	return i.Typ >= TypeFrontMatterYAML && i.Typ <= TypeFrontMatterORG
}

func (i Item) IsDone() bool {
	return i.Typ == tError || i.Typ == tEOF
}

func (i Item) IsEOF() bool {
	return i.Typ == tEOF
}

func (i Item) IsError() bool {
	return i.Typ == tError
}

func (i Item) String() string {
	switch {
	case i.Typ == tEOF:
		return "EOF"
	case i.Typ == tError:
		return string(i.Val)
	case i.Typ > tKeywordMarker:
		return fmt.Sprintf("<%s>", i.Val)
	case len(i.Val) > 50:
		return fmt.Sprintf("%v:%.20q...", i.Typ, i.Val)
	}
	return fmt.Sprintf("%v:[%s]", i.Typ, i.Val)
}

type ItemType int

const (
	tError ItemType = iota
	tEOF

	// page items
	TypeHTMLDocument       // document starting with < as first non-whitespace
	TypeHTMLComment        // We ignore leading comments
	TypeLeadSummaryDivider // <!--more-->
	TypeSummaryDividerOrg  // # more
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
	tScParam
	tScParamVal

	tText // plain text

	// preserved for later - keywords come after this
	tKeywordMarker
)
