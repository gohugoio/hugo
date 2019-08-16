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

// Package blackfriday converts Markdown to HTML using Blackfriday v1.
package blackfriday

import (
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
	"github.com/russross/blackfriday"
)

// Provider is the package entry point.
var Provider converter.NewProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	defaultBlackFriday, err := internal.NewBlackfriday(cfg)
	if err != nil {
		return nil, err
	}

	defaultExtensions := getMarkdownExtensions(defaultBlackFriday)

	pygmentsCodeFences := cfg.Cfg.GetBool("pygmentsCodeFences")
	pygmentsCodeFencesGuessSyntax := cfg.Cfg.GetBool("pygmentsCodeFencesGuessSyntax")
	pygmentsOptions := cfg.Cfg.GetString("pygmentsOptions")

	var n converter.NewConverter = func(ctx converter.DocumentContext) (converter.Converter, error) {
		b := defaultBlackFriday
		extensions := defaultExtensions

		if ctx.ConfigOverrides != nil {
			var err error
			b, err = internal.UpdateBlackFriday(b, ctx.ConfigOverrides)
			if err != nil {
				return nil, err
			}
			extensions = getMarkdownExtensions(b)
		}

		return &blackfridayConverter{
			ctx:        ctx,
			bf:         b,
			extensions: extensions,
			cfg:        cfg,

			pygmentsCodeFences:            pygmentsCodeFences,
			pygmentsCodeFencesGuessSyntax: pygmentsCodeFencesGuessSyntax,
			pygmentsOptions:               pygmentsOptions,
		}, nil
	}

	return n, nil

}

type blackfridayConverter struct {
	ctx        converter.DocumentContext
	bf         *internal.BlackFriday
	extensions int

	pygmentsCodeFences            bool
	pygmentsCodeFencesGuessSyntax bool
	pygmentsOptions               string

	cfg converter.ProviderConfig
}

func (c *blackfridayConverter) AnchorSuffix() string {
	if c.bf.PlainIDAnchors {
		return ""
	}
	return ":" + c.ctx.DocumentID
}

func (c *blackfridayConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	r := c.getHTMLRenderer(ctx.RenderTOC)

	return converter.Bytes(blackfriday.Markdown(ctx.Src, r, c.extensions)), nil

}

func (c *blackfridayConverter) getHTMLRenderer(renderTOC bool) blackfriday.Renderer {
	flags := getFlags(renderTOC, c.bf)

	documentID := c.ctx.DocumentID

	renderParameters := blackfriday.HtmlRendererParameters{
		FootnoteAnchorPrefix:       c.bf.FootnoteAnchorPrefix,
		FootnoteReturnLinkContents: c.bf.FootnoteReturnLinkContents,
	}

	if documentID != "" && !c.bf.PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = documentID + ":" + renderParameters.FootnoteAnchorPrefix
		renderParameters.HeaderIDSuffix = ":" + documentID
	}

	return &hugoHTMLRenderer{
		c:        c,
		Renderer: blackfriday.HtmlRendererWithParameters(flags, "", "", renderParameters),
	}
}

func getFlags(renderTOC bool, cfg *internal.BlackFriday) int {

	var flags int

	if renderTOC {
		flags = blackfriday.HTML_TOC
	}

	flags |= blackfriday.HTML_USE_XHTML
	flags |= blackfriday.HTML_FOOTNOTE_RETURN_LINKS

	if cfg.Smartypants {
		flags |= blackfriday.HTML_USE_SMARTYPANTS
	}

	if cfg.SmartypantsQuotesNBSP {
		flags |= blackfriday.HTML_SMARTYPANTS_QUOTES_NBSP
	}

	if cfg.AngledQuotes {
		flags |= blackfriday.HTML_SMARTYPANTS_ANGLED_QUOTES
	}

	if cfg.Fractions {
		flags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	}

	if cfg.HrefTargetBlank {
		flags |= blackfriday.HTML_HREF_TARGET_BLANK
	}

	if cfg.NofollowLinks {
		flags |= blackfriday.HTML_NOFOLLOW_LINKS
	}

	if cfg.NoreferrerLinks {
		flags |= blackfriday.HTML_NOREFERRER_LINKS
	}

	if cfg.SmartDashes {
		flags |= blackfriday.HTML_SMARTYPANTS_DASHES
	}

	if cfg.LatexDashes {
		flags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	}

	if cfg.SkipHTML {
		flags |= blackfriday.HTML_SKIP_HTML
	}

	return flags
}

func getMarkdownExtensions(cfg *internal.BlackFriday) int {
	// Default Blackfriday common extensions
	commonExtensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	// Extra Blackfriday extensions that Hugo enables by default
	flags := commonExtensions |
		blackfriday.EXTENSION_AUTO_HEADER_IDS |
		blackfriday.EXTENSION_FOOTNOTES

	for _, extension := range cfg.Extensions {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	for _, extension := range cfg.ExtensionsMask {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags &= ^flag
		}
	}
	return flags
}

var blackfridayExtensionMap = map[string]int{
	"noIntraEmphasis":        blackfriday.EXTENSION_NO_INTRA_EMPHASIS,
	"tables":                 blackfriday.EXTENSION_TABLES,
	"fencedCode":             blackfriday.EXTENSION_FENCED_CODE,
	"autolink":               blackfriday.EXTENSION_AUTOLINK,
	"strikethrough":          blackfriday.EXTENSION_STRIKETHROUGH,
	"laxHtmlBlocks":          blackfriday.EXTENSION_LAX_HTML_BLOCKS,
	"spaceHeaders":           blackfriday.EXTENSION_SPACE_HEADERS,
	"hardLineBreak":          blackfriday.EXTENSION_HARD_LINE_BREAK,
	"tabSizeEight":           blackfriday.EXTENSION_TAB_SIZE_EIGHT,
	"footnotes":              blackfriday.EXTENSION_FOOTNOTES,
	"noEmptyLineBeforeBlock": blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK,
	"headerIds":              blackfriday.EXTENSION_HEADER_IDS,
	"titleblock":             blackfriday.EXTENSION_TITLEBLOCK,
	"autoHeaderIds":          blackfriday.EXTENSION_AUTO_HEADER_IDS,
	"backslashLineBreak":     blackfriday.EXTENSION_BACKSLASH_LINE_BREAK,
	"definitionLists":        blackfriday.EXTENSION_DEFINITION_LISTS,
	"joinLines":              blackfriday.EXTENSION_JOIN_LINES,
}

var (
	_ converter.DocumentInfo = (*blackfridayConverter)(nil)
)
