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
	"sync"

	"github.com/gohugoio/hugo/common/herrors"
	htext "github.com/gohugoio/hugo/common/text"
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

	pos := ctx.PopPos()
	text := ctx.Buffer.Bytes()[pos:]
	ctx.Buffer.Truncate(pos)

	ordinal := ctx.GetAndIncrementOrdinal(ast.KindBlockquote)

	texts := string(text)
	typ := typeRegular
	alertType := resolveGitHubAlert(texts)
	if alertType != "" {
		typ = typeAlert
	}

	renderer := ctx.RenderContext().GetRenderer(hooks.BlockquoteRendererType, typ)
	if renderer == nil {
		return r.renderBlockquoteDefault(w, n, texts)
	}

	if typ == typeAlert {
		// Trim preamble: <p>[!NOTE]<br>\n but preserve leading paragraph.
		// We could possibly complicate this by moving this to the parser, but
		// keep it simple for now.
		texts = "<p>" + texts[strings.Index(texts, "\n")+1:]
	}

	var sourceRef []byte

	// Extract a source sample to use for position information.
	if nn := n.FirstChild(); nn != nil {
		var start, stop int
		for i := 0; i < nn.Lines().Len() && i < 2; i++ {
			line := nn.Lines().At(i)
			if i == 0 {
				start = line.Start
			}
			stop = line.Stop
		}
		// We do not mutate the source, so this is safe.
		sourceRef = src[start:stop]
	}

	bqctx := &blockquoteContext{
		page:             ctx.DocumentContext().Document,
		pageInner:        r.getPageInner(ctx),
		typ:              typ,
		alertType:        alertType,
		text:             hstring.RenderedString(texts),
		sourceRef:        sourceRef,
		ordinal:          ordinal,
		AttributesHolder: attributes.New(n.Attributes(), attributes.AttributesOwnerGeneral),
	}

	bqctx.createPos = func() htext.Position {
		if resolver, ok := renderer.(hooks.ElementPositionResolver); ok {
			return resolver.ResolvePosition(bqctx)
		}

		return htext.Position{
			Filename:     ctx.DocumentContext().Filename,
			LineNumber:   1,
			ColumnNumber: 1,
		}
	}

	cr := renderer.(hooks.BlockquoteRenderer)

	err := cr.RenderBlockquote(
		ctx.RenderContext().Ctx,
		w,
		bqctx,
	)
	if err != nil {
		return ast.WalkContinue, herrors.NewFileErrorFromPos(err, bqctx.createPos())
	}

	return ast.WalkContinue, nil
}

func (r *htmlRenderer) getPageInner(rctx *render.Context) any {
	pid := rctx.PeekPid()
	if pid > 0 {
		if lookup := rctx.DocumentContext().DocumentLookup; lookup != nil {
			if v := rctx.DocumentContext().DocumentLookup(pid); v != nil {
				return v
			}
		}
	}
	return rctx.DocumentContext().Document
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
	page      any
	pageInner any
	text      hstring.RenderedString
	typ       string
	sourceRef []byte
	alertType string
	ordinal   int

	// This is only used in error situations and is expensive to create,
	// so delay creation until needed.
	pos       htext.Position
	posInit   sync.Once
	createPos func() htext.Position

	*attributes.AttributesHolder
}

func (c *blockquoteContext) Type() string {
	return c.typ
}

func (c *blockquoteContext) AlertType() string {
	return c.alertType
}

func (c *blockquoteContext) Page() any {
	return c.page
}

func (c *blockquoteContext) PageInner() any {
	return c.pageInner
}

func (c *blockquoteContext) Text() hstring.RenderedString {
	return c.text
}

func (c *blockquoteContext) Ordinal() int {
	return c.ordinal
}

func (c *blockquoteContext) Position() htext.Position {
	c.posInit.Do(func() {
		c.pos = c.createPos()
	})
	return c.pos
}

func (c *blockquoteContext) PositionerSourceTarget() []byte {
	return c.sourceRef
}

var _ hooks.PositionerSourceTargetProvider = (*blockquoteContext)(nil)

// https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax#alerts
// Five types:
// [!NOTE], [!TIP], [!WARNING], [!IMPORTANT], [!CAUTION]
var gitHubAlertRe = regexp.MustCompile(`^<p>\[!(NOTE|TIP|WARNING|IMPORTANT|CAUTION)\]`)

// resolveGitHubAlert returns one of note, tip, warning, important or caution.
// An empty string if no match.
func resolveGitHubAlert(s string) string {
	m := gitHubAlertRe.FindStringSubmatch(s)
	if len(m) == 2 {
		return strings.ToLower(m[1])
	}
	return ""
}
