// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestItemValTyped(t *testing.T) {
	c := qt.New(t)

	source := []byte("3.14")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, float64(3.14))
	source = []byte(".14")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, float64(0.14))
	source = []byte("314")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, 314)
	source = []byte("314")
	c.Assert(Item{low: 0, high: len(source), isString: true}.ValTyped(source), qt.Equals, "314")
	source = []byte("314x")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "314x")
	source = []byte("314 ")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "314 ")
	source = []byte("true")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, true)
	source = []byte("false")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, false)
	source = []byte("falsex")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "falsex")
	source = []byte("xfalse")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "xfalse")
	source = []byte("truex")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "truex")
	source = []byte("xtrue")
	c.Assert(Item{low: 0, high: len(source)}.ValTyped(source), qt.Equals, "xtrue")
}

func TestItemBoolMethods(t *testing.T) {
	c := qt.New(t)

	source := []byte("  shortcode ")
	tests := []struct {
		name   string
		item   Item
		source []byte
		want   bool
		call   func(Item, []byte) bool
	}{
		{
			name: "IsText true",
			item: Item{Type: tText},
			call: func(i Item, _ []byte) bool { return i.IsText() },
			want: true,
		},
		{
			name: "IsIndentation false",
			item: Item{Type: tText},
			call: func(i Item, _ []byte) bool { return i.IsIndentation() },
			want: false,
		},
		{
			name: "IsShortcodeName",
			item: Item{Type: tScName},
			call: func(i Item, _ []byte) bool { return i.IsShortcodeName() },
			want: true,
		},
		{
			name: "IsNonWhitespace true",
			item: Item{
				Type: tText,
				low:  2,
				high: 11,
			},
			source: source,
			call:   func(i Item, src []byte) bool { return i.IsNonWhitespace(src) },
			want:   true,
		},
		{
			name: "IsShortcodeParam false",
			item: Item{Type: tScParamVal},
			call: func(i Item, _ []byte) bool { return i.IsShortcodeParam() },
			want: false,
		},
		{
			name: "IsInlineShortcodeName",
			item: Item{Type: tScNameInline},
			call: func(i Item, _ []byte) bool { return i.IsInlineShortcodeName() },
			want: true,
		},
		{
			name: "IsLeftShortcodeDelim tLeftDelimScWithMarkup",
			item: Item{Type: tLeftDelimScWithMarkup},
			call: func(i Item, _ []byte) bool { return i.IsLeftShortcodeDelim() },
			want: true,
		},
		{
			name: "IsLeftShortcodeDelim tLeftDelimScNoMarkup",
			item: Item{Type: tLeftDelimScNoMarkup},
			call: func(i Item, _ []byte) bool { return i.IsLeftShortcodeDelim() },
			want: true,
		},
		{
			name: "IsRightShortcodeDelim tRightDelimScWithMarkup",
			item: Item{Type: tRightDelimScWithMarkup},
			call: func(i Item, _ []byte) bool { return i.IsRightShortcodeDelim() },
			want: true,
		},
		{
			name: "IsRightShortcodeDelim tRightDelimScNoMarkup",
			item: Item{Type: tRightDelimScNoMarkup},
			call: func(i Item, _ []byte) bool { return i.IsRightShortcodeDelim() },
			want: true,
		},
		{
			name: "IsShortcodeClose",
			item: Item{Type: tScClose},
			call: func(i Item, _ []byte) bool { return i.IsShortcodeClose() },
			want: true,
		},
		{
			name: "IsShortcodeParamVal",
			item: Item{Type: tScParamVal},
			call: func(i Item, _ []byte) bool { return i.IsShortcodeParamVal() },
			want: true,
		},
		{
			name: "IsShortcodeMarkupDelimiter tLeftDelimScWithMarkup",
			item: Item{Type: tLeftDelimScWithMarkup},
			call: func(i Item, _ []byte) bool { return i.IsShortcodeMarkupDelimiter() },
			want: true,
		},
		{
			name: "IsShortcodeMarkupDelimiter tRightDelimScWithMarkup",
			item: Item{Type: tRightDelimScWithMarkup},
			call: func(i Item, _ []byte) bool { return i.IsShortcodeMarkupDelimiter() },
			want: true,
		},
		{
			name: "IsFrontMatter TypeFrontMatterYAML",
			item: Item{Type: TypeFrontMatterYAML},
			call: func(i Item, _ []byte) bool { return i.IsFrontMatter() },
			want: true,
		},
		{
			name: "IsFrontMatter TypeFrontMatterTOML",
			item: Item{Type: TypeFrontMatterTOML},
			call: func(i Item, _ []byte) bool { return i.IsFrontMatter() },
			want: true,
		},
		{
			name: "IsFrontMatter TypeFrontMatterJSON",
			item: Item{Type: TypeFrontMatterJSON},
			call: func(i Item, _ []byte) bool { return i.IsFrontMatter() },
			want: true,
		},
		{
			name: "IsFrontMatter TypeFrontMatterORG",
			item: Item{Type: TypeFrontMatterORG},
			call: func(i Item, _ []byte) bool { return i.IsFrontMatter() },
			want: true,
		},
		{
			name: "IsDone tError",
			item: Item{Type: tError},
			call: func(i Item, _ []byte) bool { return i.IsDone() },
			want: true,
		},
		{
			name: "IsDone tEOF",
			item: Item{Type: tEOF},
			call: func(i Item, _ []byte) bool { return i.IsDone() },
			want: true,
		},
		{
			name: "IsEOF",
			item: Item{Type: tEOF},
			call: func(i Item, _ []byte) bool { return i.IsEOF() },
			want: true,
		},
		{
			name: "IsError",
			item: Item{Type: tError},
			call: func(i Item, _ []byte) bool { return i.IsError() },
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.call(tt.item, tt.source)
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}

func TestItem_ToString(t *testing.T) {
	c := qt.New(t)

	source := []byte("src")
	long := make([]byte, 100)
	for i := range long {
		long[i] = byte(i)
	}

	tests := []struct {
		name   string
		item   Item
		source []byte
		want   string
		call   func(Item, []byte) string
	}{
		{
			name: "EOF",
			item: Item{Type: tEOF},
			call: func(i Item, _ []byte) string { return i.ToString(source) },
			want: "EOF",
		},
		{
			name: "Error",
			item: Item{Type: tError},
			call: func(i Item, _ []byte) string { return i.ToString(source) },
			want: "",
		},
		{
			name: "Indentation",
			item: Item{Type: tIndentation},
			call: func(i Item, _ []byte) string { return i.ToString(source) },
			want: "tIndentation:[]",
		},
		{
			name: "Long",
			item: Item{Type: tKeywordMarker + 1, low: 0, high: 100},
			call: func(i Item, _ []byte) string { return i.ToString(long) },
			want: "<" + string(long) + ">",
		},
		{
			name: "Empty",
			item: Item{Type: tKeywordMarker + 1},
			call: func(i Item, _ []byte) string { return i.ToString([]byte("")) },
			want: "<>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.call(tt.item, tt.source)
			c.Assert(got, qt.Equals, tt.want)
		})
	}
}
