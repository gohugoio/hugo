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
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type lexerTest struct {
	name  string
	input string
	items []Item
}

func nti(tp ItemType, val string) Item {
	return Item{tp, 0, []byte(val), false}
}

var (
	tstJSON                = `{ "a": { "b": "\"Hugo\"}" } }`
	tstFrontMatterTOML     = nti(TypeFrontMatterTOML, "foo = \"bar\"\n")
	tstFrontMatterYAML     = nti(TypeFrontMatterYAML, "foo: \"bar\"\n")
	tstFrontMatterYAMLCRLF = nti(TypeFrontMatterYAML, "foo: \"bar\"\r\n")
	tstFrontMatterJSON     = nti(TypeFrontMatterJSON, tstJSON+"\r\n")
	tstSomeText            = nti(tText, "\nSome text.\n")
	tstSummaryDivider      = nti(TypeLeadSummaryDivider, "<!--more-->\n")
	tstNewline             = nti(tText, "\n")

	tstORG = `
#+TITLE: T1
#+AUTHOR: A1
#+DESCRIPTION: D1
`
	tstFrontMatterORG = nti(TypeFrontMatterORG, tstORG)
)

var crLfReplacer = strings.NewReplacer("\r", "#", "\n", "$")

// TODO(bep) a way to toggle ORG mode vs the rest.
var frontMatterTests = []lexerTest{
	{"empty", "", []Item{tstEOF}},
	{"Byte order mark", "\ufeff\nSome text.\n", []Item{nti(TypeIgnore, "\ufeff"), tstSomeText, tstEOF}},
	{"HTML Document", `  <html>  `, []Item{nti(tError, "plain HTML documents not supported")}},
	{"HTML Document with shortcode", `<html>{{< sc1 >}}</html>`, []Item{nti(tError, "plain HTML documents not supported")}},
	{"No front matter", "\nSome text.\n", []Item{tstSomeText, tstEOF}},
	{"YAML front matter", "---\nfoo: \"bar\"\n---\n\nSome text.\n", []Item{tstFrontMatterYAML, tstSomeText, tstEOF}},
	{"YAML empty front matter", "---\n---\n\nSome text.\n", []Item{nti(TypeFrontMatterYAML, ""), tstSomeText, tstEOF}},
	{"YAML commented out front matter", "<!--\n---\nfoo: \"bar\"\n---\n-->\nSome text.\n", []Item{nti(TypeIgnore, "<!--\n"), tstFrontMatterYAML, nti(TypeIgnore, "-->"), tstSomeText, tstEOF}},
	{"YAML commented out front matter, no end", "<!--\n---\nfoo: \"bar\"\n---\nSome text.\n", []Item{nti(TypeIgnore, "<!--\n"), tstFrontMatterYAML, nti(tError, "starting HTML comment with no end")}},
	// Note that we keep all bytes as they are, but we need to handle CRLF
	{"YAML front matter CRLF", "---\r\nfoo: \"bar\"\r\n---\n\nSome text.\n", []Item{tstFrontMatterYAMLCRLF, tstSomeText, tstEOF}},
	{"TOML front matter", "+++\nfoo = \"bar\"\n+++\n\nSome text.\n", []Item{tstFrontMatterTOML, tstSomeText, tstEOF}},
	{"JSON front matter", tstJSON + "\r\n\nSome text.\n", []Item{tstFrontMatterJSON, tstSomeText, tstEOF}},
	{"ORG front matter", tstORG + "\nSome text.\n", []Item{tstFrontMatterORG, tstSomeText, tstEOF}},
	{"Summary divider ORG", tstORG + "\nSome text.\n# more\nSome text.\n", []Item{tstFrontMatterORG, tstSomeText, nti(TypeLeadSummaryDivider, "# more\n"), nti(tText, "Some text.\n"), tstEOF}},
	{"Summary divider", "+++\nfoo = \"bar\"\n+++\n\nSome text.\n<!--more-->\nSome text.\n", []Item{tstFrontMatterTOML, tstSomeText, tstSummaryDivider, nti(tText, "Some text.\n"), tstEOF}},
	{"Summary divider same line", "+++\nfoo = \"bar\"\n+++\n\nSome text.<!--more-->Some text.\n", []Item{tstFrontMatterTOML, nti(tText, "\nSome text."), nti(TypeLeadSummaryDivider, "<!--more-->"), nti(tText, "Some text.\n"), tstEOF}},
	// https://github.com/gohugoio/hugo/issues/5402
	{"Summary and shortcode, no space", "+++\nfoo = \"bar\"\n+++\n\nSome text.\n<!--more-->{{< sc1 >}}\nSome text.\n", []Item{tstFrontMatterTOML, tstSomeText, nti(TypeLeadSummaryDivider, "<!--more-->"), tstLeftNoMD, tstSC1, tstRightNoMD, tstSomeText, tstEOF}},
	// https://github.com/gohugoio/hugo/issues/5464
	{"Summary and shortcode only", "+++\nfoo = \"bar\"\n+++\n{{< sc1 >}}\n<!--more-->\n{{< sc2 >}}", []Item{tstFrontMatterTOML, tstLeftNoMD, tstSC1, tstRightNoMD, tstNewline, tstSummaryDivider, tstLeftNoMD, tstSC2, tstRightNoMD, tstEOF}},
}

func TestFrontMatter(t *testing.T) {
	t.Parallel()
	for i, test := range frontMatterTests {
		items := collect([]byte(test.input), false, lexIntroSection)
		if !equal(items, test.items) {
			got := crLfReplacer.Replace(fmt.Sprint(items))
			expected := crLfReplacer.Replace(fmt.Sprint(test.items))
			t.Errorf("[%d] %s: got\n\t%v\nexpected\n\t%v", i, test.name, got, expected)
		}
	}
}

func collectWithConfig(input []byte, skipFrontMatter bool, stateStart stateFunc, cfg Config) (items []Item) {
	l := newPageLexer(input, stateStart, cfg)
	l.run()
	t := l.newIterator()

	for {
		item := t.Next()
		items = append(items, item)
		if item.Type == tEOF || item.Type == tError {
			break
		}
	}
	return
}

func collect(input []byte, skipFrontMatter bool, stateStart stateFunc) (items []Item) {
	var cfg Config

	return collectWithConfig(input, skipFrontMatter, stateStart, cfg)

}

// no positional checking, for now ...
func equal(i1, i2 []Item) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].Type != i2[k].Type {
			return false
		}

		if !reflect.DeepEqual(i1[k].Val, i2[k].Val) {
			return false
		}
	}
	return true
}
