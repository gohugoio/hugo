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
	"fmt"
	"path/filepath"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/alecthomas/chroma/styles"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/highlight"
	"github.com/gohugoio/hugo/markup/markup_config"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provide{}

type provide struct {
}

func (p provide) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	md := newMarkdown(cfg.MarkupConfig)
	return converter.NewProvider("goldmark", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &goldmarkConverter{
			ctx: ctx,
			cfg: cfg,
			md:  md,
		}, nil
	}), nil
}

type goldmarkConverter struct {
	md  goldmark.Markdown
	ctx converter.DocumentContext
	cfg converter.ProviderConfig
}

func newMarkdown(mcfg markup_config.Config) goldmark.Markdown {
	cfg := mcfg.Goldmark

	var (
		extensions = []goldmark.Extender{
			newTocExtension(),
		}
		rendererOptions []renderer.Option
		parserOptions   []parser.Option
	)

	if cfg.Renderer.HardWraps {
		rendererOptions = append(rendererOptions, html.WithHardWraps())
	}

	if cfg.Renderer.XHTML {
		rendererOptions = append(rendererOptions, html.WithXHTML())
	}

	if cfg.Renderer.Unsafe {
		rendererOptions = append(rendererOptions, html.WithUnsafe())
	}

	if mcfg.Highlight.CodeFences {
		extensions = append(extensions, newHighlighting(mcfg.Highlight))
	}

	if cfg.Extensions.Table {
		extensions = append(extensions, extension.Table)
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

	if cfg.Extensions.Typographer {
		extensions = append(extensions, extension.Typographer)
	}

	if cfg.Extensions.DefinitionList {
		extensions = append(extensions, extension.DefinitionList)
	}

	if cfg.Extensions.Footnote {
		extensions = append(extensions, extension.Footnote)
	}

	if cfg.Parser.AutoHeadingID {
		parserOptions = append(parserOptions, parser.WithAutoHeadingID())
	}

	if cfg.Parser.Attribute {
		parserOptions = append(parserOptions, parser.WithAttribute())
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

type converterResult struct {
	converter.Result
	toc tableofcontents.Root
}

func (c converterResult) TableOfContents() tableofcontents.Root {
	return c.toc
}

func (c *goldmarkConverter) Convert(ctx converter.RenderContext) (result converter.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			dir := afero.GetTempDir(hugofs.Os, "hugo_bugs")
			name := fmt.Sprintf("goldmark_%s.txt", c.ctx.DocumentID)
			filename := filepath.Join(dir, name)
			afero.WriteFile(hugofs.Os, filename, ctx.Src, 07555)
			fmt.Print(string(debug.Stack()))
			err = errors.Errorf("[BUG] goldmark: %s: create an issue on GitHub attaching the file in: %s", r, filename)
		}
	}()

	buf := &bytes.Buffer{}
	result = buf
	pctx := parser.NewContext()
	pctx.Set(tocEnableKey, ctx.RenderTOC)

	reader := text.NewReader(ctx.Src)

	doc := c.md.Parser().Parse(
		reader,
		parser.WithContext(pctx),
	)

	if err := c.md.Renderer().Render(buf, ctx.Src, doc); err != nil {
		return nil, err
	}

	if toc, ok := pctx.Get(tocResultKey).(tableofcontents.Root); ok {
		return converterResult{
			Result: buf,
			toc:    toc,
		}, nil
	}

	return buf, nil
}

func newHighlighting(cfg highlight.Config) goldmark.Extender {
	style := styles.Get(cfg.Style)
	if style == nil {
		style = styles.Fallback
	}

	e := hl.NewHighlighting(
		hl.WithStyle(cfg.Style),
		hl.WithCodeBlockOptions(highlight.GetCodeBlockOptions()),
		hl.WithFormatOptions(
			cfg.ToHTMLOptions()...,
		),

		hl.WithWrapperRenderer(func(w util.BufWriter, ctx hl.CodeBlockContext, entering bool) {
			l, hasLang := ctx.Language()
			var language string
			if hasLang {
				language = string(l)
			}

			if entering {
				if !ctx.Highlighted() {
					w.WriteString(`<pre>`)
					highlight.WriteCodeTag(w, language)
					return
				}
				w.WriteString(`<div class="highlight">`)
				return
			}

			if !ctx.Highlighted() {
				w.WriteString(`</code></pre>`)
				return
			}

			w.WriteString("</div>")

		}),
	)

	return e
}
