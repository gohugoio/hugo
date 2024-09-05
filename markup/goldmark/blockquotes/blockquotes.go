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

package blockquotes

import (
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/internal/attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type (
	blockquotesExtension struct{}
	htmlRenderer         struct{}
)

func New() goldmark.Extender {
	return &blockquotesExtension{}
}

func (e *blockquotesExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newHTMLRenderer(), 100),
	))
}

func newHTMLRenderer() renderer.NodeRenderer {
	r := &htmlRenderer{}
	return r
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
}

const (
	typeRegular = "regular"
	typeAlert   = "alert"
)

func (r *htmlRenderer) renderBlockquote(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)

	n := node.(*ast.Blockquote)

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.PushPos(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	text := ctx.PopRenderedString()

	ordinal := ctx.GetAndIncrementOrdinal(ast.KindBlockquote)

	typ := typeRegular
	alert := resolveBlockQuoteAlert(string(text))
	if alert.typ != "" {
		typ = typeAlert
	}

	renderer := ctx.RenderContext().GetRenderer(hooks.BlockquoteRendererType, typ)
	if renderer == nil {
		return r.renderBlockquoteDefault(w, n, text)
	}

	if typ == typeAlert {
		// Trim preamble: <p>[!NOTE]<br>\n but preserve leading paragraph.
		// We could possibly complicate this by moving this to the parser, but
		// keep it simple for now.
		text = "<p>" + text[strings.Index(text, "\n")+1:]
	}

	bqctx := &blockquoteContext{
		BaseContext:      render.NewBaseContext(ctx, renderer, n, src, nil, ordinal),
		typ:              typ,
		alert:            alert,
		text:             hstring.HTML(text),
		AttributesHolder: attributes.New(n.Attributes(), attributes.AttributesOwnerGeneral),
	}

	cr := renderer.(hooks.BlockquoteRenderer)

	err := cr.RenderBlockquote(
		ctx.RenderContext().Ctx,
		w,
		bqctx,
	)
	if err != nil {
		return ast.WalkContinue, herrors.NewFileErrorFromPos(err, bqctx.Position())
	}

	return ast.WalkContinue, nil
}

// Code borrowed from goldmark's html renderer.
func (r *htmlRenderer) renderBlockquoteDefault(
	w util.BufWriter, n ast.Node, text string,
) (ast.WalkStatus, error) {
	if n.Attributes() != nil {
		_, _ = w.WriteString("<blockquote")
		html.RenderAttributes(w, n, html.BlockquoteAttributeFilter)
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("<blockquote>\n")
	}

	_, _ = w.WriteString(text)

	_, _ = w.WriteString("</blockquote>\n")
	return ast.WalkContinue, nil
}

type blockquoteContext struct {
	hooks.BaseContext
	text  hstring.HTML
	typ   string
	alert blockQuoteAlert
	*attributes.AttributesHolder
}

func (c *blockquoteContext) Type() string {
	return c.typ
}

func (c *blockquoteContext) AlertType() string {
	return c.alert.typ
}

func (c *blockquoteContext) AlertTitle() hstring.HTML {
	return hstring.HTML(c.alert.title)
}

func (c *blockquoteContext) AlertSign() string {
	return c.alert.sign
}

func (c *blockquoteContext) Text() hstring.HTML {
	return c.text
}

var blockQuoteAlertRe = regexp.MustCompile(`^<p>\[!([a-zA-Z]+)\](-|\+)?[^\S\r\n]?([^\n]*)\n?`)

func resolveBlockQuoteAlert(s string) blockQuoteAlert {
	m := blockQuoteAlertRe.FindStringSubmatch(s)
	if len(m) == 4 {
		title := strings.TrimSpace(m[3])
		title = strings.TrimRight(title, "</p>")
		return blockQuoteAlert{
			typ:   strings.ToLower(m[1]),
			sign:  m[2],
			title: title,
		}
	}

	return blockQuoteAlert{}
}

// Blockquote alert syntax was introduced by GitHub, but is also used
// by Obsidian which also support some extended attributes: More types, alert titles and a +/- sign for folding.
type blockQuoteAlert struct {
	typ   string
	sign  string
	title string
}
