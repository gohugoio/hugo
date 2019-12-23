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

// Package mmark converts Markdown to HTML using MMark v1.
package mmark

import (
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/blackfriday/blackfriday_config"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/miekg/mmark"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	defaultBlackFriday := cfg.MarkupConfig.BlackFriday
	defaultExtensions := getMmarkExtensions(defaultBlackFriday)

	return converter.NewProvider("mmark", func(ctx converter.DocumentContext) (converter.Converter, error) {
		b := defaultBlackFriday
		extensions := defaultExtensions

		if ctx.ConfigOverrides != nil {
			var err error
			b, err = blackfriday_config.UpdateConfig(b, ctx.ConfigOverrides)
			if err != nil {
				return nil, err
			}
			extensions = getMmarkExtensions(b)
		}

		return &mmarkConverter{
			ctx:        ctx,
			b:          b,
			extensions: extensions,
			cfg:        cfg,
		}, nil
	}), nil

}

type mmarkConverter struct {
	ctx        converter.DocumentContext
	extensions int
	b          blackfriday_config.Config
	cfg        converter.ProviderConfig
}

func (c *mmarkConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	r := getHTMLRenderer(c.ctx, c.b, c.cfg)
	return mmark.Parse(ctx.Src, r, c.extensions), nil
}

func (c *mmarkConverter) Supports(feature identity.Identity) bool {
	return false
}

func getHTMLRenderer(
	ctx converter.DocumentContext,
	cfg blackfriday_config.Config,
	pcfg converter.ProviderConfig) mmark.Renderer {

	var (
		flags      int
		documentID string
	)

	documentID = ctx.DocumentID

	renderParameters := mmark.HtmlRendererParameters{
		FootnoteAnchorPrefix:       cfg.FootnoteAnchorPrefix,
		FootnoteReturnLinkContents: cfg.FootnoteReturnLinkContents,
	}

	if documentID != "" && !cfg.PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = documentID + ":" + renderParameters.FootnoteAnchorPrefix
	}

	htmlFlags := flags
	htmlFlags |= mmark.HTML_FOOTNOTE_RETURN_LINKS

	return &mmarkRenderer{
		BlackfridayConfig: cfg,
		Config:            pcfg,
		Renderer:          mmark.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters),
	}

}

func getMmarkExtensions(cfg blackfriday_config.Config) int {
	flags := 0
	flags |= mmark.EXTENSION_TABLES
	flags |= mmark.EXTENSION_FENCED_CODE
	flags |= mmark.EXTENSION_AUTOLINK
	flags |= mmark.EXTENSION_SPACE_HEADERS
	flags |= mmark.EXTENSION_CITATION
	flags |= mmark.EXTENSION_TITLEBLOCK_TOML
	flags |= mmark.EXTENSION_HEADER_IDS
	flags |= mmark.EXTENSION_AUTO_HEADER_IDS
	flags |= mmark.EXTENSION_UNIQUE_HEADER_IDS
	flags |= mmark.EXTENSION_FOOTNOTES
	flags |= mmark.EXTENSION_SHORT_REF
	flags |= mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	flags |= mmark.EXTENSION_INCLUDE

	for _, extension := range cfg.Extensions {
		if flag, ok := mmarkExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	return flags
}

var mmarkExtensionMap = map[string]int{
	"tables":                 mmark.EXTENSION_TABLES,
	"fencedCode":             mmark.EXTENSION_FENCED_CODE,
	"autolink":               mmark.EXTENSION_AUTOLINK,
	"laxHtmlBlocks":          mmark.EXTENSION_LAX_HTML_BLOCKS,
	"spaceHeaders":           mmark.EXTENSION_SPACE_HEADERS,
	"hardLineBreak":          mmark.EXTENSION_HARD_LINE_BREAK,
	"footnotes":              mmark.EXTENSION_FOOTNOTES,
	"noEmptyLineBeforeBlock": mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK,
	"headerIds":              mmark.EXTENSION_HEADER_IDS,
	"autoHeaderIds":          mmark.EXTENSION_AUTO_HEADER_IDS,
}
