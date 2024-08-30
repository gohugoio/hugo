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

	"github.com/gohugoio/hugo-goldmark-extensions/passthrough"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/internal/attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func New(cfg goldmark_config.Passthrough) goldmark.Extender {
	if !cfg.Enable {
		return nil
	}
	return &passthroughExtension{cfg: cfg}
}

type (
	passthroughExtension struct {
		cfg goldmark_config.Passthrough
	}
	htmlRenderer struct{}
)

func (e *passthroughExtension) Extend(m goldmark.Markdown) {
	configuredInlines := e.cfg.Delimiters.Inline
	configuredBlocks := e.cfg.Delimiters.Block

	inlineDelimiters := make([]passthrough.Delimiters, len(configuredInlines))
	blockDelimiters := make([]passthrough.Delimiters, len(configuredBlocks))

	for i, d := range configuredInlines {
		inlineDelimiters[i] = passthrough.Delimiters{
			Open:  d[0],
			Close: d[1],
		}
	}

	for i, d := range configuredBlocks {
		blockDelimiters[i] = passthrough.Delimiters{
			Open:  d[0],
			Close: d[1],
		}
	}

	pse := passthrough.New(
		passthrough.Config{
			InlineDelimiters: inlineDelimiters,
			BlockDelimiters:  blockDelimiters,
		},
	)

	pse.Extend(m)

	// Set up render hooks if configured.
	// Upstream passthrough inline = 101
	// Upstream passthrough block = 99
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newHTMLRenderer(), 90),
	))
}

func newHTMLRenderer() renderer.NodeRenderer {
	r := &htmlRenderer{}
	return r
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(passthrough.KindPassthroughBlock, r.renderPassthroughBlock)
	reg.Register(passthrough.KindPassthroughInline, r.renderPassthroughBlock)
}

func (r *htmlRenderer) renderPassthroughBlock(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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
		s = string(nn.Text(src))
		typ = "inline"
		delims = nn.Delimiters
	case (*passthrough.PassthroughBlock):
		l := nn.Lines().Len()
		var buff bytes.Buffer
		for i := 0; i < l; i++ {
			line := nn.Lines().At(i)
			buff.Write(line.Value(src))
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

	pctx := &passthroughContext{
		BaseContext:      render.NewBaseContext(ctx, renderer, node, src, nil, ordinal),
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

func (p *passthroughContext) Type() string {
	return p.typ
}

func (p *passthroughContext) Inner() string {
	return p.inner
}
