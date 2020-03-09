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

package hugolib

import (
	"fmt"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/parser/pageparser"
)

var (
	internalSummaryDividerBase      = "HUGOMORE42"
	internalSummaryDividerBaseBytes = []byte(internalSummaryDividerBase)
	internalSummaryDividerPre       = []byte("\n\n" + internalSummaryDividerBase + "\n\n")
)

// The content related items on a Page.
type pageContent struct {
	selfLayout string
	truncated  bool

	cmap *pageContentMap

	shortcodeState *shortcodeHandler

	source rawPageContent
}

// returns the content to be processed by Blackfriday or similar.
func (p pageContent) contentToRender(renderedShortcodes map[string]string) []byte {
	source := p.source.parsed.Input()

	c := make([]byte, 0, len(source)+(len(source)/10))

	for _, it := range p.cmap.items {
		switch v := it.(type) {
		case pageparser.Item:
			c = append(c, source[v.Pos:v.Pos+len(v.Val)]...)
		case pageContentReplacement:
			c = append(c, v.val...)
		case *shortcode:
			if !v.insertPlaceholder() {
				// Insert the rendered shortcode.
				renderedShortcode, found := renderedShortcodes[v.placeholder]
				if !found {
					// This should never happen.
					panic(fmt.Sprintf("rendered shortcode %q not found", v.placeholder))
				}

				c = append(c, []byte(renderedShortcode)...)

			} else {
				// Insert the placeholder so we can insert the content after
				// markdown processing.
				c = append(c, []byte(v.placeholder)...)

			}
		default:
			panic(fmt.Sprintf("unknown item type %T", it))
		}
	}

	return c
}

func (p pageContent) selfLayoutForOutput(f output.Format) string {
	if p.selfLayout == "" {
		return ""
	}
	return p.selfLayout + f.Name
}

type rawPageContent struct {
	hasSummaryDivider bool

	// The AST of the parsed page. Contains information about:
	// shortcodes, front matter, summary indicators.
	parsed pageparser.Result

	// Returns the position in bytes after any front matter.
	posMainContent int

	// These are set if we're able to determine this from the source.
	posSummaryEnd int
	posBodyStart  int
}

type pageContentReplacement struct {
	val []byte

	source pageparser.Item
}

type pageContentMap struct {

	// If not, we can skip any pre-rendering of shortcodes.
	hasMarkdownShortcode bool

	// Indicates whether we must do placeholder replacements.
	hasNonMarkdownShortcode bool

	//  *shortcode, pageContentReplacement or pageparser.Item
	items []interface{}
}

func (p *pageContentMap) AddBytes(item pageparser.Item) {
	p.items = append(p.items, item)
}

func (p *pageContentMap) AddReplacement(val []byte, source pageparser.Item) {
	p.items = append(p.items, pageContentReplacement{val: val, source: source})
}

func (p *pageContentMap) AddShortcode(s *shortcode) {
	p.items = append(p.items, s)
	if s.insertPlaceholder() {
		p.hasNonMarkdownShortcode = true
	} else {
		p.hasMarkdownShortcode = true
	}
}
