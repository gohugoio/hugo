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

package tables

import (
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/internal/attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type (
	ext          struct{}
	htmlRenderer struct{}
)

func New() goldmark.Extender {
	return &ext{}
}

func (e *ext) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newHTMLRenderer(), 100),
	))
}

func newHTMLRenderer() renderer.NodeRenderer {
	r := &htmlRenderer{}
	return r
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gast.KindTable, r.renderTable)
	reg.Register(gast.KindTableHeader, r.renderHeaderOrRow)
	reg.Register(gast.KindTableRow, r.renderHeaderOrRow)
	reg.Register(gast.KindTableCell, r.renderCell)
}

func (r *htmlRenderer) renderTable(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)
	if entering {
		// This will be modified below.
		table := &hooks.Table{}
		ctx.PushValue(gast.KindTable, table)
		return ast.WalkContinue, nil
	}

	v := ctx.PopValue(gast.KindTable)
	if v == nil {
		panic("table not found")
	}

	table := v.(*hooks.Table)

	renderer := ctx.RenderContext().GetRenderer(hooks.TableRendererType, nil)
	if renderer == nil {
		panic("table hook renderer not found")
	}

	ordinal := ctx.GetAndIncrementOrdinal(gast.KindTable)

	tctx := &tableContext{
		BaseContext:      render.NewBaseContext(ctx, renderer, n, source, nil, ordinal),
		AttributesHolder: attributes.New(n.Attributes(), attributes.AttributesOwnerGeneral),
		tHead:            table.THead,
		tBody:            table.TBody,
	}

	cr := renderer.(hooks.TableRenderer)

	err := cr.RenderTable(
		ctx.RenderContext().Ctx,
		w,
		tctx,
	)
	if err != nil {
		return ast.WalkContinue, herrors.NewFileErrorFromPos(err, tctx.Position())
	}

	return ast.WalkContinue, nil
}

func (r *htmlRenderer) peekTable(ctx *render.Context) *hooks.Table {
	v := ctx.PeekValue(gast.KindTable)
	if v == nil {
		panic("table not found")
	}
	return v.(*hooks.Table)
}

func (r *htmlRenderer) renderCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.PushPos(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	n := node.(*gast.TableCell)

	text := ctx.PopRenderedString()

	table := r.peekTable(ctx)

	var alignment string
	switch n.Alignment {
	case gast.AlignLeft:
		alignment = "left"
	case gast.AlignRight:
		alignment = "right"
	case gast.AlignCenter:
		alignment = "center"
	default:
		alignment = ""
	}

	cell := hooks.TableCell{Text: hstring.HTML(text), Alignment: alignment}

	if node.Parent().Kind() == gast.KindTableHeader {
		table.THead[len(table.THead)-1] = append(table.THead[len(table.THead)-1], cell)
	} else {
		table.TBody[len(table.TBody)-1] = append(table.TBody[len(table.TBody)-1], cell)
	}

	return ast.WalkContinue, nil
}

func (r *htmlRenderer) renderHeaderOrRow(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)
	table := r.peekTable(ctx)
	if entering {
		if n.Kind() == gast.KindTableHeader {
			table.THead = append(table.THead, hooks.TableRow{})
		} else {
			table.TBody = append(table.TBody, hooks.TableRow{})
		}
		return ast.WalkContinue, nil
	}

	return ast.WalkContinue, nil
}

type tableContext struct {
	hooks.BaseContext
	*attributes.AttributesHolder

	tHead []hooks.TableRow
	tBody []hooks.TableRow
}

func (c *tableContext) THead() []hooks.TableRow {
	return c.tHead
}

func (c *tableContext) TBody() []hooks.TableRow {
	return c.tBody
}
