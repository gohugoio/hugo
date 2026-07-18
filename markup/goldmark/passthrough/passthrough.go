// Copyright 2024 The Hugo Authors. All rights reserved.
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

package passthrough

import (
	"bytes"
	"io"

	"github.com/gohugoio/hugo-goldmark-extensions/passthrough/v2"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/internal/attributes"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
)

// New returns the goldmark v2 parser and HTML renderer extensions for the
// passthrough feature, or (nil, nil) if it is disabled. The parser comes from
// the upstream module; the renderer is Hugo's own so the passthrough render
// hooks apply.
func New(cfg goldmark_config.Passthrough) (parser.Extension, html.Extension) {
	if !cfg.Enable {
		return nil, nil
	}

	inlineDelimiters := make([]passthrough.Delimiters, len(cfg.Delimiters.Inline))
	blockDelimiters := make([]passthrough.Delimiters, len(cfg.Delimiters.Block))

	for i, d := range cfg.Delimiters.Inline {
		inlineDelimiters[i] = passthrough.Delimiters{Open: d[0], Close: d[1]}
	}
	for i, d := range cfg.Delimiters.Block {
		blockDelimiters[i] = passthrough.Delimiters{Open: d[0], Close: d[1]}
	}

	pcfg := passthrough.Config{
		InlineDelimiters: inlineDelimiters,
		BlockDelimiters:  blockDelimiters,
	}

	return passthrough.NewParser(pcfg), &htmlRenderer{}
}

type htmlRenderer struct{}

func (r *htmlRenderer) RendererOptions(*html.Config) []html.Option {
	return []html.Option{
		html.WithNodeRenderer(passthrough.KindPassthroughBlock, html.NodeRendererFunc(r.renderPassthroughBlock)),
		html.WithNodeRenderer(passthrough.KindPassthroughInline, html.NodeRendererFunc(r.renderPassthroughBlock)),
	}
}

func (r *htmlRenderer) renderPassthroughBlock(w io.Writer, src []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)

	if entering {
		return ast.WalkContinue, nil
	}

	var (
		s      string
		typ    string
		delims *passthrough.Delimiters
	)

	switch nn := node.(type) {
	case *passthrough.PassthroughInline:
		// GOLDMARK-V2: PassthroughInline exposes a Segment field (Text() is gone).
		s = string(nn.Segment.Bytes(src))
		typ = "inline"
		delims = nn.Delimiters
	case *passthrough.PassthroughBlock:
		// GOLDMARK-V2: block content is read via BlockNode.Source() now.
		var buff bytes.Buffer
		for _, line := range nn.Source() {
			buff.Write(line.Bytes(src))
		}
		s = buff.String()
		typ = "block"
		delims = nn.Delimiters
	}

	renderer := ctx.RenderContext().GetRenderer(hooks.PassthroughRendererType, typ)
	if renderer == nil {
		// Write the raw content if no renderer is found.
		ctx.WriteString(s)
		return ast.WalkContinue, nil
	}

	// Inline and block passthroughs share the same ordinal counter.
	ordinal := ctx.GetAndIncrementOrdinal(passthrough.KindPassthroughBlock)

	// Trim the delimiters.
	s = s[len(delims.Open) : len(s)-len(delims.Close)]

	pctx := passthroughContext{
		BaseContext:      render.NewBaseContext(ctx, renderer, node, src, ordinal),
		inner:            s,
		typ:              typ,
		AttributesHolder: attributes.New(node.Attributes(), attributes.AttributesOwnerGeneral),
	}

	pr := renderer.(hooks.PassthroughRenderer)

	if err := pr.RenderPassthrough(ctx.RenderContext().Ctx, w, pctx); err != nil {
		return ast.WalkStop, err
	}

	return ast.WalkContinue, nil
}

type passthroughContext struct {
	hooks.BaseContext

	typ   string // inner or block
	inner string

	*attributes.AttributesHolder
}

func (p passthroughContext) Type() string {
	return p.typ
}

func (p passthroughContext) Inner() string {
	return p.inner
}
