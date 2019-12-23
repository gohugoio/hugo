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

package goldmark

import (
	"github.com/gohugoio/hugo/markup/converter/hooks"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

var _ renderer.SetOptioner = (*linkRenderer)(nil)

func newLinkRenderer() renderer.NodeRenderer {
	r := &linkRenderer{
		Config: html.Config{
			Writer: html.DefaultWriter,
		},
	}
	return r
}

func newLinks() goldmark.Extender {
	return &links{}
}

type linkContext struct {
	page        interface{}
	destination string
	title       string
	text        string
	plainText   string
}

func (ctx linkContext) Destination() string {
	return ctx.destination
}

func (ctx linkContext) Resolved() bool {
	return false
}

func (ctx linkContext) Page() interface{} {
	return ctx.page
}

func (ctx linkContext) Text() string {
	return ctx.text
}

func (ctx linkContext) PlainText() string {
	return ctx.plainText
}

func (ctx linkContext) Title() string {
	return ctx.title
}

type linkRenderer struct {
	html.Config
}

func (r *linkRenderer) SetOption(name renderer.OptionName, value interface{}) {
	r.Config.SetOption(name, value)
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *linkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *linkRenderer) renderDefaultImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString("<img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	}
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(n.Text(source))
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		_ = w.WriteByte('"')
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *linkRenderer) renderDefaultLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString("<a href=\"")
		if r.Unsafe || !html.IsDangerousURL(n.Destination) {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
		_ = w.WriteByte('"')
		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			r.Writer.Write(w, n.Title)
			_ = w.WriteByte('"')
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r *linkRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	var h *hooks.Render

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h != nil && h.ImageRenderer != nil
	}

	if !ok {
		return r.renderDefaultImage(w, source, node, entering)
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.pos = ctx.Buffer.Len()
		return ast.WalkContinue, nil
	}

	text := ctx.Buffer.Bytes()[ctx.pos:]
	ctx.Buffer.Truncate(ctx.pos)

	err := h.ImageRenderer.Render(
		w,
		linkContext{
			page:        ctx.DocumentContext().Document,
			destination: string(n.Destination),
			title:       string(n.Title),
			text:        string(text),
			plainText:   string(n.Text(source)),
		},
	)

	ctx.AddIdentity(h.ImageRenderer.GetIdentity())

	return ast.WalkContinue, err

}

func (r *linkRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	var h *hooks.Render

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h != nil && h.LinkRenderer != nil
	}

	if !ok {
		return r.renderDefaultLink(w, source, node, entering)
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.pos = ctx.Buffer.Len()
		return ast.WalkContinue, nil
	}

	text := ctx.Buffer.Bytes()[ctx.pos:]
	ctx.Buffer.Truncate(ctx.pos)

	err := h.LinkRenderer.Render(
		w,
		linkContext{
			page:        ctx.DocumentContext().Document,
			destination: string(n.Destination),
			title:       string(n.Title),
			text:        string(text),
			plainText:   string(n.Text(source)),
		},
	)

	ctx.AddIdentity(h.LinkRenderer.GetIdentity())

	return ast.WalkContinue, err

}

type links struct {
}

// Extend implements goldmark.Extender.
func (e *links) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newLinkRenderer(), 100),
	))
}
