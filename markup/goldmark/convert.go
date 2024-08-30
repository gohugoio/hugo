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

// Package goldmark converts Markdown to HTML using Goldmark.
package goldmark

import (
	"bytes"

	"github.com/gohugoio/hugo-goldmark-extensions/extras"
	"github.com/gohugoio/hugo/markup/goldmark/blockquotes"
	"github.com/gohugoio/hugo/markup/goldmark/codeblocks"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/hugocontext"
	"github.com/gohugoio/hugo/markup/goldmark/images"
	"github.com/gohugoio/hugo/markup/goldmark/internal/extensions/attributes"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/goldmark/passthrough"
	"github.com/gohugoio/hugo/markup/goldmark/tables"
	"github.com/yuin/goldmark/util"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"
)

const (
	internalAttrPrefix = "_h__"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provide{}

type provide struct{}

func (p provide) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	md := newMarkdown(cfg)

	return converter.NewProvider("goldmark", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &goldmarkConverter{
			ctx: ctx,
			cfg: cfg,
			md:  md,
			sanitizeAnchorName: func(s string) string {
				return sanitizeAnchorNameString(s, cfg.MarkupConfig().Goldmark.Parser.AutoHeadingIDType)
			},
		}, nil
	}), nil
}

var _ converter.AnchorNameSanitizer = (*goldmarkConverter)(nil)

type goldmarkConverter struct {
	md  goldmark.Markdown
	ctx converter.DocumentContext
	cfg converter.ProviderConfig

	sanitizeAnchorName func(s string) string
}

func (c *goldmarkConverter) SanitizeAnchorName(s string) string {
	return c.sanitizeAnchorName(s)
}

func newMarkdown(pcfg converter.ProviderConfig) goldmark.Markdown {
	mcfg := pcfg.MarkupConfig()
	cfg := mcfg.Goldmark
	var rendererOptions []renderer.Option

	if cfg.Renderer.HardWraps {
		rendererOptions = append(rendererOptions, html.WithHardWraps())
	}

	if cfg.Renderer.XHTML {
		rendererOptions = append(rendererOptions, html.WithXHTML())
	}

	if cfg.Renderer.Unsafe {
		rendererOptions = append(rendererOptions, html.WithUnsafe())
	}

	tocRendererOptions := make([]renderer.Option, len(rendererOptions))
	if rendererOptions != nil {
		copy(tocRendererOptions, rendererOptions)
	}
	tocRendererOptions = append(tocRendererOptions,
		renderer.WithNodeRenderers(util.Prioritized(extension.NewStrikethroughHTMLRenderer(), 500)),
		renderer.WithNodeRenderers(util.Prioritized(emoji.NewHTMLRenderer(), 200)))
	var (
		extensions = []goldmark.Extender{
			hugocontext.New(),
			newLinks(cfg),
			newTocExtension(tocRendererOptions),
			blockquotes.New(),
		}
		parserOptions []parser.Option
	)

	extensions = append(extensions, images.New(cfg.Parser.WrapStandAloneImageWithinParagraph))

	extensions = append(extensions, extras.New(
		extras.Config{
			Delete:      extras.DeleteConfig{Enable: cfg.Extensions.Extras.Delete.Enable},
			Insert:      extras.InsertConfig{Enable: cfg.Extensions.Extras.Insert.Enable},
			Mark:        extras.MarkConfig{Enable: cfg.Extensions.Extras.Mark.Enable},
			Subscript:   extras.SubscriptConfig{Enable: cfg.Extensions.Extras.Subscript.Enable},
			Superscript: extras.SuperscriptConfig{Enable: cfg.Extensions.Extras.Superscript.Enable},
		},
	))

	if mcfg.Highlight.CodeFences {
		extensions = append(extensions, codeblocks.New())
	}

	if cfg.Extensions.Table {
		extensions = append(extensions, extension.Table)
		extensions = append(extensions, tables.New())
	}

	if cfg.Extensions.Strikethrough {
		extensions = append(extensions, extension.Strikethrough)
	}

	if cfg.Extensions.Linkify {
		extensions = append(extensions, extension.Linkify)
	}

	if cfg.Extensions.TaskList {
		extensions = append(extensions, extension.TaskList)
	}

	if !cfg.Extensions.Typographer.Disable {
		t := extension.NewTypographer(
			extension.WithTypographicSubstitutions(toTypographicPunctuationMap(cfg.Extensions.Typographer)),
		)
		extensions = append(extensions, t)
	}

	if cfg.Extensions.DefinitionList {
		extensions = append(extensions, extension.DefinitionList)
	}

	if cfg.Extensions.Footnote {
		extensions = append(extensions, extension.Footnote)
	}

	if cfg.Extensions.CJK.Enable {
		opts := []extension.CJKOption{}
		if cfg.Extensions.CJK.EastAsianLineBreaks {
			if cfg.Extensions.CJK.EastAsianLineBreaksStyle == "css3draft" {
				opts = append(opts, extension.WithEastAsianLineBreaks(extension.EastAsianLineBreaksCSS3Draft))
			} else {
				opts = append(opts, extension.WithEastAsianLineBreaks())
			}
		}

		if cfg.Extensions.CJK.EscapedSpace {
			opts = append(opts, extension.WithEscapedSpace())
		}
		c := extension.NewCJK(opts...)
		extensions = append(extensions, c)
	}

	if cfg.Extensions.Passthrough.Enable {
		extensions = append(extensions, passthrough.New(cfg.Extensions.Passthrough))
	}

	if pcfg.Conf.EnableEmoji() {
		extensions = append(extensions, emoji.Emoji)
	}

	if cfg.Parser.AutoHeadingID {
		parserOptions = append(parserOptions, parser.WithAutoHeadingID())
	}

	if cfg.Parser.Attribute.Title {
		parserOptions = append(parserOptions, parser.WithAttribute())
	}

	if cfg.Parser.Attribute.Block {
		extensions = append(extensions, attributes.New())
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extensions...,
		),
		goldmark.WithParserOptions(
			parserOptions...,
		),
		goldmark.WithRendererOptions(
			rendererOptions...,
		),
	)

	return md
}

type parserResult struct {
	doc any
	toc *tableofcontents.Fragments
}

func (p parserResult) Doc() any {
	return p.doc
}

func (p parserResult) TableOfContents() *tableofcontents.Fragments {
	return p.toc
}

type renderResult struct {
	converter.ResultRender
}

type converterResult struct {
	converter.ResultRender
	tableOfContentsProvider
}

type tableOfContentsProvider interface {
	TableOfContents() *tableofcontents.Fragments
}

func (c *goldmarkConverter) Parse(ctx converter.RenderContext) (converter.ResultParse, error) {
	pctx := c.newParserContext(ctx)
	reader := text.NewReader(ctx.Src)

	doc := c.md.Parser().Parse(
		reader,
		parser.WithContext(pctx),
	)

	return parserResult{
		doc: doc,
		toc: pctx.TableOfContents(),
	}, nil
}

func (c *goldmarkConverter) Render(ctx converter.RenderContext, doc any) (converter.ResultRender, error) {
	n := doc.(ast.Node)
	buf := &render.BufWriter{Buffer: &bytes.Buffer{}}

	rcx := &render.RenderContextDataHolder{
		Rctx: ctx,
		Dctx: c.ctx,
	}

	w := &render.Context{
		BufWriter:   buf,
		ContextData: rcx,
	}

	if err := c.md.Renderer().Render(w, ctx.Src, n); err != nil {
		return nil, err
	}

	return renderResult{
		ResultRender: buf,
	}, nil
}

func (c *goldmarkConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	parseResult, err := c.Parse(ctx)
	if err != nil {
		return nil, err
	}
	renderResult, err := c.Render(ctx, parseResult.Doc())
	if err != nil {
		return nil, err
	}
	return converterResult{
		ResultRender:            renderResult,
		tableOfContentsProvider: parseResult,
	}, nil
}

func (c *goldmarkConverter) newParserContext(rctx converter.RenderContext) *parserContext {
	ctx := parser.NewContext(parser.WithIDs(newIDFactory(c.cfg.MarkupConfig().Goldmark.Parser.AutoHeadingIDType)))
	ctx.Set(tocEnableKey, rctx.RenderTOC)
	return &parserContext{
		Context: ctx,
	}
}

type parserContext struct {
	parser.Context
}

func (p *parserContext) TableOfContents() *tableofcontents.Fragments {
	if v := p.Get(tocResultKey); v != nil {
		return v.(*tableofcontents.Fragments)
	}
	return nil
}

// Note: It's tempting to put this in the config package, but that doesn't work.
// TODO(bep) create upstream issue.
func toTypographicPunctuationMap(t goldmark_config.Typographer) map[extension.TypographicPunctuation][]byte {
	return map[extension.TypographicPunctuation][]byte{
		extension.LeftSingleQuote:  []byte(t.LeftSingleQuote),
		extension.RightSingleQuote: []byte(t.RightSingleQuote),
		extension.LeftDoubleQuote:  []byte(t.LeftDoubleQuote),
		extension.RightDoubleQuote: []byte(t.RightDoubleQuote),
		extension.EnDash:           []byte(t.EnDash),
		extension.EmDash:           []byte(t.EmDash),
		extension.Ellipsis:         []byte(t.Ellipsis),
		extension.LeftAngleQuote:   []byte(t.LeftAngleQuote),
		extension.RightAngleQuote:  []byte(t.RightAngleQuote),
		extension.Apostrophe:       []byte(t.Apostrophe),
	}
}
