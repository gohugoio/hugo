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
	"github.com/gohugoio/hugo-goldmark-extensions/extras/v2"

	"github.com/gohugoio/hugo/markup/goldmark/blockquotes"
	"github.com/gohugoio/hugo/markup/goldmark/codeblocks"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/hugocontext"
	"github.com/gohugoio/hugo/markup/goldmark/images"
	"github.com/gohugoio/hugo/markup/goldmark/internal/extensions/attributes"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/goldmark/passthrough"
	"github.com/gohugoio/hugo/markup/goldmark/tables"
	emoji "github.com/yuin/goldmark-emoji/v2"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"
)

const (
	// Don't change this. This pattern is lso used in the image render hooks.
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
				return sanitizeAnchorNameString(s, cfg.MarkupConfig().Goldmark.Parser.AutoIDType)
			},
		}, nil
	}), nil
}

var _ converter.AnchorNameSanitizer = (*goldmarkConverter)(nil)

// markdownHandler holds the goldmark v2 parser and renderers.
// GOLDMARK-V2: v2 has no unified goldmark.Markdown; the parser and renderer are
// built and used separately.
type markdownHandler struct {
	parser parser.Parser
	// renderer is the main HTML renderer (with Hugo's render hooks).
	renderer html.Renderer
	// tocRenderer renders heading inline content for the table of contents.
	tocRenderer html.Renderer
}

type goldmarkConverter struct {
	md  *markdownHandler
	ctx converter.DocumentContext
	cfg converter.ProviderConfig

	sanitizeAnchorName func(s string) string
}

// rendererExtension adapts a slice of html.Option into an html.Extension so
// that Hugo's node renderers are registered after (and thus override) the
// CommonMark renderers.
type rendererExtension []html.Option

func (e rendererExtension) RendererOptions(*html.Config) []html.Option {
	return e
}

func (c *goldmarkConverter) SanitizeAnchorName(s string) string {
	return c.sanitizeAnchorName(s)
}

func newMarkdown(pcfg converter.ProviderConfig) *markdownHandler {
	mcfg := pcfg.MarkupConfig()
	cfg := mcfg.Goldmark

	// GOLDMARK-V2: build the parser and renderer separately. Hugo's own node
	// renderers are wrapped as html.Extensions (via rendererExtension) so that
	// they are registered after CommonMark and override it.
	var (
		parserOptions   []parser.Option
		rendererOptions []html.Option
		parserExts      []parser.Extension
		rendererExts    []html.Extension
	)

	// Base HTML renderer options.
	if cfg.Renderer.HardWraps {
		rendererOptions = append(rendererOptions, html.WithHardWraps())
	}
	if cfg.Renderer.XHTML {
		rendererOptions = append(rendererOptions, html.WithXHTML())
	}
	if cfg.Renderer.Unsafe {
		rendererOptions = append(rendererOptions, html.WithUnsafe())
	}

	// Hugo context (.RenderShortcodes) — parser + renderer.
	hc := hugocontext.New(pcfg.Logger, cfg.Renderer.Unsafe)
	parserOptions = append(parserOptions, hc.ParserOptions()...)
	rendererExts = append(rendererExts, rendererExtension(hc.RendererOptions()))

	// Link/image/heading render hooks.
	rendererExts = append(rendererExts, rendererExtension(newLinks(cfg).RendererOptions()))

	// Blockquotes.
	rendererExts = append(rendererExts, blockquotes.New())

	// Images.
	parserExts = append(parserExts, images.New(cfg.Parser.WrapStandAloneImageWithinParagraph))

	parserExts = append(parserExts, extras.NewParser(
		extras.Config{
			Delete:      extras.DeleteConfig{Enable: cfg.Extensions.Extras.Delete.Enable},
			Insert:      extras.InsertConfig{Enable: cfg.Extensions.Extras.Insert.Enable},
			Mark:        extras.MarkConfig{Enable: cfg.Extensions.Extras.Mark.Enable},
			Subscript:   extras.SubscriptConfig{Enable: cfg.Extensions.Extras.Subscript.Enable},
			Superscript: extras.SuperscriptConfig{Enable: cfg.Extensions.Extras.Superscript.Enable},
		},
	))

	if mcfg.Highlight.CodeFences {
		rendererExts = append(rendererExts, codeblocks.New())
	}

	if cfg.Extensions.Table {
		parserExts = append(parserExts, extension.NewTableParser())
		rendererExts = append(rendererExts, extension.NewTableHTMLRenderer())
		rendererExts = append(rendererExts, tables.New())
	}

	if cfg.Extensions.Strikethrough {
		parserExts = append(parserExts, extension.NewStrikethroughParser())
		rendererExts = append(rendererExts, extension.NewStrikethroughHTMLRenderer())
	}

	if cfg.Extensions.Linkify {
		parserExts = append(parserExts, extension.NewLinkifyParser())
	}

	if cfg.Extensions.TaskList {
		parserExts = append(parserExts, extension.NewTaskCheckBoxParser())
		rendererExts = append(rendererExts, extension.NewTaskListItemHTMLRenderer())
	}

	if !cfg.Extensions.Typographer.Disable {
		parserExts = append(parserExts, extension.NewTypographerParser(
			extension.WithTypographicSubstitutions(toTypographicPunctuationMap(cfg.Extensions.Typographer)),
		))
	}

	if cfg.Extensions.DefinitionList {
		parserExts = append(parserExts, extension.NewDefinitionListParser())
		rendererExts = append(rendererExts, extension.NewDefinitionListHTMLRenderer())
	}

	if cfg.Extensions.Footnote.Enable {
		parserExts = append(parserExts, extension.NewFootnoteParser())
		// GOLDMARK-V2: footnote options moved from the parser to the HTML
		// renderer, and the option names were shortened.
		opts := []extension.FootnoteHTMLRendererOption{
			extension.WithBacklinkHTML(cfg.Extensions.Footnote.BacklinkHTML),
		}
		if cfg.Extensions.Footnote.EnableAutoIDPrefix {
			opts = append(opts,
				extension.WithIDPrefixFunction(func(n ast.Node) []byte {
					documentID := n.OwnerDocument().Metadata()["documentID"].(string)
					return []byte("h" + documentID)
				}))
		}
		rendererExts = append(rendererExts, extension.NewFootnoteHTMLRenderer(opts...))
	}

	if cfg.Extensions.CJK.Enable {
		// GOLDMARK-V2: extension.NewCJK is gone; CJK is now expressed through
		// parser/renderer options.
		if cfg.Extensions.CJK.EastAsianLineBreaks {
			if cfg.Extensions.CJK.EastAsianLineBreaksStyle == "css3draft" {
				rendererOptions = append(rendererOptions, html.WithEastAsianLineBreaks(html.EastAsianLineBreaksCSS3Draft))
			} else {
				rendererOptions = append(rendererOptions, html.WithEastAsianLineBreaks(html.EastAsianLineBreaksSimple))
			}
		}
		if cfg.Extensions.CJK.EscapedSpace {
			parserOptions = append(parserOptions, parser.WithEscapedSpace())
			rendererOptions = append(rendererOptions, html.WithEscapedSpace())
		}
	}

	if cfg.Extensions.Passthrough.Enable {
		if pe, re := passthrough.New(cfg.Extensions.Passthrough); pe != nil {
			parserExts = append(parserExts, pe)
			rendererExts = append(rendererExts, re)
		}
	}

	if pcfg.Conf.EnableEmoji() {
		parserExts = append(parserExts, emoji.Parser)
	}

	if cfg.Parser.Attribute.Title {
		parserOptions = append(parserOptions, parser.WithAttribute())
	}

	if cfg.Parser.Attribute.Block || cfg.Parser.AutoHeadingID || cfg.Parser.AutoDefinitionTermID {
		parserExts = append(parserExts, attributes.New(cfg.Parser))
	}

	// Stateless ID generator; per-document uniqueness is handled by goldmark.
	parserOptions = append(parserOptions,
		parser.WithIDGenerator(newIDFactory(cfg.Parser.AutoIDType)))

	parserOptions = append(parserOptions, parser.WithExtensions(parserExts...))
	rendererOptions = append(rendererOptions, html.WithExtensions(rendererExts...))

	// The TOC renderer renders heading inline content to plain-ish HTML. It uses
	// the default CommonMark renderers plus strikethrough (when enabled).
	var tocRendererOptions []html.Option
	tocRendererOptions = append(tocRendererOptions, rendererOptions...)
	if cfg.Extensions.Strikethrough {
		tocRendererOptions = append(tocRendererOptions,
			html.WithExtensions(extension.NewStrikethroughHTMLRenderer()))
	}

	return &markdownHandler{
		parser:      parser.New(parserOptions...),
		renderer:    html.New(rendererOptions...),
		tocRenderer: html.New(tocRendererOptions...),
	}
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
	// GOLDMARK-V2: Parser.Parse takes only the source now; there is no way to
	// pass a custom parser.Context, so the ID generator is configured on the
	// parser and the table of contents is built by walking the returned AST.
	doc := c.md.parser.Parse(ctx.Src)
	doc.OwnerDocument().AddMeta("documentID", c.ctx.DocumentID)

	var toc *tableofcontents.Fragments
	if ctx.RenderTOC {
		toc = buildTableOfContents(doc, ctx, c.ctx, c.md.tocRenderer)
	}

	return parserResult{
		doc: doc,
		toc: toc,
	}, nil
}

func (c *goldmarkConverter) Render(ctx converter.RenderContext, doc any) (converter.ResultRender, error) {
	n := doc.(ast.Node)

	w := render.NewContext(ctx, c.ctx)

	// TODO1 goldmrk v2. Would be great to have a context.Context or something.
	if err := c.md.renderer.Render(w, ctx.Src, n); err != nil {
		return nil, err
	}

	return renderResult{
		ResultRender: w.Buffer,
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
