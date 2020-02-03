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
	"unicode"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/blackfriday/blackfriday_config"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/russross/blackfriday"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	defaultExtensions := getMarkdownExtensions(cfg.MarkupConfig.BlackFriday)

	return converter.NewProvider("blackfriday", func(ctx converter.DocumentContext) (converter.Converter, error) {
		b := cfg.MarkupConfig.BlackFriday
		extensions := defaultExtensions

		if ctx.ConfigOverrides != nil {
			var err error
			b, err = blackfriday_config.UpdateConfig(b, ctx.ConfigOverrides)
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
		}, nil
	}), nil

}

type blackfridayConverter struct {
	ctx        converter.DocumentContext
	bf         blackfriday_config.Config
	extensions int
	cfg        converter.ProviderConfig
}

func (c *blackfridayConverter) SanitizeAnchorName(s string) string {
	return SanitizedAnchorName(s)
}

// SanitizedAnchorName is how Blackfriday sanitizes anchor names.
// Implementation borrowed from https://github.com/russross/blackfriday/blob/a477dd1646916742841ed20379f941cfa6c5bb6f/block.go#L1464
func SanitizedAnchorName(text string) string {
	var anchorName []rune
	futureDash := false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
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

func (c *blackfridayConverter) Supports(feature identity.Identity) bool {
	return false
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

func getFlags(renderTOC bool, cfg blackfriday_config.Config) int {

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

func getMarkdownExtensions(cfg blackfriday_config.Config) int {
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
	_ converter.DocumentInfo        = (*blackfridayConverter)(nil)
	_ converter.AnchorNameSanitizer = (*blackfridayConverter)(nil)
)
