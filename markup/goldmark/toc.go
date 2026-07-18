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

package goldmark

import (
	"regexp"
	"strings"

	passthrough "github.com/gohugoio/hugo-goldmark-extensions/passthrough/v2"
	"github.com/microcosm-cc/bluemonday"
	strikethroughAst "github.com/yuin/goldmark/v2/extension/ast"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/renderer/html"
)

// TODO1 goldmrk v2. This doesn't look great.
// buildTableOfContents walks the parsed document and builds the table of
// contents.
//
// GOLDMARK-V2: In v1 this was a parser AST transformer that stored its result
// in the parser.Context. Since v2's Parser.Parse(source) creates its own
// internal context that callers can neither seed nor read back, we now walk the
// returned AST directly after parsing. The heading anchor ids are read from the
// heading `id` attributes (already assigned and de-duplicated during parsing).
func buildTableOfContents(n ast.Node, rc converter.RenderContext, dc converter.DocumentContext, r html.Renderer) *tableofcontents.Fragments {
	var (
		toc         tableofcontents.Builder
		tocHeading  = &tableofcontents.Heading{}
		level       int
		row         = -1
		inHeading   bool
		headingText = render.NewContext(rc, dc)
		identifiers []string
	)

	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		if n.Kind() == ast.KindHeading {
			if inHeading && !entering {
				tocHeading.Title = sanitizeTOCHeadingTitle(headingText.String())
				headingText.Reset()
				toc.AddAt(tocHeading, row, level-1)
				tocHeading = &tableofcontents.Heading{}
				inHeading = false
				return s, nil
			}

			inHeading = true
		}

		if !(inHeading && entering) {
			return s, nil
		}

		switch n.Kind() {
		case ast.KindHeading:
			heading := n.(*ast.Heading)
			level = heading.Level

			if level == 1 || row == -1 {
				row++
			}

			if id, found := heading.Attribute("id"); found {
				idStr := string(id.Bytes(rc.Src))
				tocHeading.ID = idStr
				tocHeading.Level = level
				identifiers = append(identifiers, idStr)
			}
		case
			ast.KindCodeSpan,
			ast.KindLink,
			ast.KindImage,
			ast.KindEmphasis,
			strikethroughAst.KindStrikethrough:
			// GOLDMARK-V2: emoji (goldmark-emoji) and the hugo-goldmark-extensions
			// (extras: delete/insert/mark/subscript/superscript) node kinds are
			// omitted here because those modules have no v2 release yet.
			err := r.Render(headingText, rc.Src, n)
			if err != nil {
				return s, err
			}

			return ast.WalkSkipChildren, nil
		case
			ast.KindAutoLink,
			ast.KindRawHTML,
			ast.KindText,
			passthrough.KindPassthroughInline,
			passthrough.KindPassthroughBlock:
			// GOLDMARK-V2: ast.KindString was removed (string content is now Text).
			err := r.Render(headingText, rc.Src, n)
			if err != nil {
				return s, err
			}
		}

		return s, nil
	})

	if len(identifiers) > 0 {
		toc.SetIdentifiers(identifiers)
	}

	return toc.Build()
}

var tocSanitizerPolicy = newTOCSanitizerPolicy()

// newTOCSanitizerPolicy returns a bluemonday policy for sanitizing TOC heading
// titles against an allowlist of inline HTML elements and attributes,
// specifically excluding anchor elements to prevent links within TOC heading
// titles.
func newTOCSanitizerPolicy() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements(
		"abbr", "b", "bdi", "bdo", "br", "cite", "code", "data", "del", "dfn",
		"em", "i", "ins", "kbd", "mark", "q", "rp", "rt", "ruby", "s", "samp",
		"small", "span", "strong", "sub", "sup", "time", "u", "var", "wbr",
	)
	p.AllowStandardAttributes()
	p.AllowStyling()
	p.AllowImages()
	p.AllowAttrs("cite").OnElements("del", "ins", "q")
	p.AllowAttrs("datetime").OnElements("del", "ins", "time")
	p.AllowAttrs("value").OnElements("data")
	return p
}

var whiteSpaceRe = regexp.MustCompile(`\s+`)

// sanitizeTOCHeadingTitle sanitizes s for use as a TOC heading title.
func sanitizeTOCHeadingTitle(s string) string {
	if strings.IndexByte(s, '<') == -1 {
		return s
	}

	// Sanitize the string.
	ss := tocSanitizerPolicy.Sanitize(s)

	// Remove extraneous whitespace.
	return whiteSpaceRe.ReplaceAllString(strings.TrimSpace(ss), " ")
}
