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

var _ renderer.SetOptioner = (*hookedRenderer)(nil)

func newLinkRenderer() renderer.NodeRenderer {
	r := &hookedRenderer{
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

type headingContext struct {
	page      interface{}
	level     int
	anchor    string
	text      string
	plainText string
}

func (ctx headingContext) Page() interface{} {
	return ctx.page
}

func (ctx headingContext) Level() int {
	return ctx.level
}

func (ctx headingContext) Anchor() string {
	return ctx.anchor
}

func (ctx headingContext) Text() string {
	return ctx.text
}

func (ctx headingContext) PlainText() string {
	return ctx.plainText
}

type hookedRenderer struct {
	html.Config
}

func (r *hookedRenderer) SetOption(name renderer.OptionName, value interface{}) {
	r.Config.SetOption(name, value)
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *hookedRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindHeading, r.renderHeading)
}

// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *hookedRenderer) RenderAttributes(w util.BufWriter, node ast.Node) {

	for _, attr := range node.Attributes() {
		_, _ = w.WriteString(" ")
		_, _ = w.Write(attr.Name)
		_, _ = w.WriteString(`="`)
		_, _ = w.Write(util.EscapeHTML(attr.Value.([]byte)))
		_ = w.WriteByte('"')
	}
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *hookedRenderer) renderDefaultImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *hookedRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	var h *hooks.Renderers

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

	err := h.ImageRenderer.RenderLink(
		w,
		linkContext{
			page:        ctx.DocumentContext().Document,
			destination: string(n.Destination),
			title:       string(n.Title),
			text:        string(text),
			plainText:   string(n.Text(source)),
		},
	)

	ctx.AddIdentity(h.ImageRenderer)

	return ast.WalkContinue, err

}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *hookedRenderer) renderDefaultLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *hookedRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	var h *hooks.Renderers

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

	err := h.LinkRenderer.RenderLink(
		w,
		linkContext{
			page:        ctx.DocumentContext().Document,
			destination: string(n.Destination),
			title:       string(n.Title),
			text:        string(text),
			plainText:   string(n.Text(source)),
		},
	)

	// TODO(bep) I have a working branch that fixes these rather confusing identity types,
	// but for now it's important that it's not .GetIdentity() that's added here,
	// to make sure we search the entire chain on changes.
	ctx.AddIdentity(h.LinkRenderer)

	return ast.WalkContinue, err
}

func (r *hookedRenderer) renderDefaultHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString("<h")
		_ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			r.RenderAttributes(w, node)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</h")
		_ = w.WriteByte("0123456"[n.Level])
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	var h *hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h != nil && h.HeadingRenderer != nil
	}

	if !ok {
		return r.renderDefaultHeading(w, source, node, entering)
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.pos = ctx.Buffer.Len()
		return ast.WalkContinue, nil
	}

	text := ctx.Buffer.Bytes()[ctx.pos:]
	ctx.Buffer.Truncate(ctx.pos)
	// All ast.Heading nodes are guaranteed to have an attribute called "id"
	// that is an array of bytes that encode a valid string.
	anchori, _ := n.AttributeString("id")
	anchor := anchori.([]byte)

	err := h.HeadingRenderer.RenderHeading(
		w,
		headingContext{
			page:      ctx.DocumentContext().Document,
			level:     n.Level,
			anchor:    string(anchor),
			text:      string(text),
			plainText: string(n.Text(source)),
		},
	)

	ctx.AddIdentity(h.HeadingRenderer)

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
