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
	typ itemType
	pos pos
	Val string
}

func (i Item) IsText() bool {
	return i.typ == tText
}

func (i Item) IsShortcodeName() bool {
	return i.typ == tScName
}

func (i Item) IsLeftShortcodeDelim() bool {
	return i.typ == tLeftDelimScWithMarkup || i.typ == tLeftDelimScNoMarkup
}

func (i Item) IsRightShortcodeDelim() bool {
	return i.typ == tRightDelimScWithMarkup || i.typ == tRightDelimScNoMarkup
}

func (i Item) IsShortcodeClose() bool {
	return i.typ == tScClose
}

func (i Item) IsShortcodeParam() bool {
	return i.typ == tScParam
}

func (i Item) IsShortcodeParamVal() bool {
	return i.typ == tScParamVal
}

func (i Item) IsShortcodeMarkupDelimiter() bool {
	return i.typ == tLeftDelimScWithMarkup || i.typ == tRightDelimScWithMarkup
}

func (i Item) IsDone() bool {
	return i.typ == tError || i.typ == tEOF
}

func (i Item) IsEOF() bool {
	return i.typ == tEOF
}

func (i Item) IsError() bool {
	return i.typ == tError
}

func (i Item) String() string {
	switch {
	case i.typ == tEOF:
		return "EOF"
	case i.typ == tError:
		return i.Val
	case i.typ > tKeywordMarker:
		return fmt.Sprintf("<%s>", i.Val)
	case len(i.Val) > 20:
		return fmt.Sprintf("%.20q...", i.Val)
	}
	return fmt.Sprintf("[%s]", i.Val)
}

type itemType int

const (
	tError itemType = iota
	tEOF

	// shortcode items
	tLeftDelimScNoMarkup
	tRightDelimScNoMarkup
	tLeftDelimScWithMarkup
	tRightDelimScWithMarkup
	tScClose
	tScName
	tScParam
	tScParamVal

	//itemIdentifier
	tText // plain text, used for everything outside the shortcodes

	// preserved for later - keywords come after this
	tKeywordMarker
)
