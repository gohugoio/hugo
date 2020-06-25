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
	"strconv"

	"github.com/gohugoio/hugo/markup/converter/hooks"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
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

type footnoteLinkContext struct {
	page  interface{}
	index int
}

func (ctx footnoteLinkContext) Page() interface{} {
	return ctx.page
}

func (ctx footnoteLinkContext) Index() int {
	return ctx.index
}

type footnoteContext struct {
	page      interface{}
	ref       string
	index     int
	text      string
	plainText string
}

func (ctx footnoteContext) Page() interface{} {
	return ctx.page
}

func (ctx footnoteContext) Ref() string {
	return ctx.ref
}

func (ctx footnoteContext) Index() int {
	return ctx.index
}

func (ctx footnoteContext) Text() string {
	return ctx.text
}

func (ctx footnoteContext) PlainText() string {
	return ctx.plainText
}

type footnotesContext struct {
	page      interface{}
	footnotes []hooks.FootnoteContext
}

func (ctx footnotesContext) Page() interface{} {
	return ctx.page
}

func (ctx footnotesContext) Footnotes() []hooks.FootnoteContext {
	return ctx.footnotes
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
	reg.Register(extast.KindFootnoteLink, r.renderFootnoteLink)
	reg.Register(extast.KindFootnoteBackLink, r.renderFootnoteBackLink)
	reg.Register(extast.KindFootnote, r.renderFootnote)
	reg.Register(extast.KindFootnoteList, r.renderFootnoteList)
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
		ctx.pushPosition(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	pos := ctx.popPosition()
	text := ctx.Buffer.Bytes()[pos:]
	ctx.Buffer.Truncate(pos)

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

	ctx.AddIdentity(h.ImageRenderer.GetIdentity())

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
		ctx.pushPosition(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	pos := ctx.popPosition()
	text := ctx.Buffer.Bytes()[pos:]
	ctx.Buffer.Truncate(pos)

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

	ctx.AddIdentity(h.LinkRenderer.GetIdentity())

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
		ctx.pushPosition(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	pos := ctx.popPosition()
	text := ctx.Buffer.Bytes()[pos:]
	ctx.Buffer.Truncate(pos)
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

	ctx.AddIdentity(h.HeadingRenderer.GetIdentity())

	return ast.WalkContinue, err
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/64d4e16bf4548242453a87665afd78954f1aae3e/extension/footnote.go#L242
func (r *hookedRenderer) renderDefaultFootnoteLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*extast.FootnoteLink)
		is := strconv.Itoa(n.Index)
		_, _ = w.WriteString(`<sup id="fnref:`)
		_, _ = w.WriteString(is)
		_, _ = w.WriteString(`"><a href="#fn:`)
		_, _ = w.WriteString(is)
		_, _ = w.WriteString(`" class="footnote-ref" role="doc-noteref">`)
		_, _ = w.WriteString(is)
		_, _ = w.WriteString(`</a></sup>`)
	}
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) renderFootnoteLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*extast.FootnoteLink)
	var h *hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h != nil && h.FootnoteLinkRenderer != nil
	}

	if !ok {
		return r.renderDefaultFootnoteLink(w, source, node, entering)
	}

	if entering {
		ctx.pushPosition(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	pos := ctx.popPosition()
	ctx.Buffer.Truncate(pos)

	err := h.FootnoteLinkRenderer.RenderFootnoteLink(
		w,
		footnoteLinkContext{
			page:  ctx.DocumentContext().Document,
			index: n.Index,
		},
	)

	ctx.AddIdentity(h.FootnoteLinkRenderer.GetIdentity())

	return ast.WalkContinue, err
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/64d4e16bf4548242453a87665afd78954f1aae3e/extension/footnote.go#L257
func (r *hookedRenderer) renderDefaultFootnoteBackLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*extast.FootnoteBackLink)
		is := strconv.Itoa(n.Index)
		_, _ = w.WriteString(` <a href="#fnref:`)
		_, _ = w.WriteString(is)
		_, _ = w.WriteString(`" class="footnote-backref" role="doc-backlink">`)
		_, _ = w.WriteString("&#x21a9;&#xfe0e;")
		_, _ = w.WriteString(`</a>`)
	}
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) renderFootnoteBackLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var h *hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h != nil && h.FootnotesRenderer != nil
	}

	if !ok {
		return r.renderDefaultFootnoteBackLink(w, source, node, entering)
	}

	return ast.WalkContinue, nil
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/64d4e16bf4548242453a87665afd78954f1aae3e/extension/footnote.go#L270
func (r *hookedRenderer) renderDefaultFootnote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*extast.Footnote)
	is := strconv.Itoa(n.Index)
	if entering {
		_, _ = w.WriteString(`<li id="fn:`)
		_, _ = w.WriteString(is)
		_, _ = w.WriteString(`" role="doc-endnote"`)
		if node.Attributes() != nil {
			html.RenderAttributes(w, node, html.ListItemAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) renderFootnote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*extast.Footnote)
	var h *hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h != nil && h.FootnotesRenderer != nil
	}

	if !ok {
		return r.renderDefaultFootnote(w, source, node, entering)
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.pushPosition(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	pos := ctx.popPosition()
	text := ctx.Buffer.Bytes()[pos:]
	ctx.Buffer.Truncate(pos)

	ctx.pushFootnote(footnoteContext{
		page:      ctx.DocumentContext().Document,
		ref:       string(n.Ref),
		index:     n.Index,
		text:      string(text),
		plainText: string(n.Text(source)),
	})

	return ast.WalkContinue, nil
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/64d4e16bf4548242453a87665afd78954f1aae3e/extension/footnote.go#L287
func (r *hookedRenderer) renderDefaultFootnoteList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	tag := "section"
	if r.Config.XHTML {
		tag = "div"
	}
	if entering {
		_, _ = w.WriteString("<")
		_, _ = w.WriteString(tag)
		_, _ = w.WriteString(` class="footnotes" role="doc-endnotes"`)
		if node.Attributes() != nil {
			html.RenderAttributes(w, node, html.GlobalAttributeFilter)
		}
		_ = w.WriteByte('>')
		if r.Config.XHTML {
			_, _ = w.WriteString("\n<hr />\n")
		} else {
			_, _ = w.WriteString("\n<hr>\n")
		}
		_, _ = w.WriteString("<ol>\n")
	} else {
		_, _ = w.WriteString("</ol>\n")
		_, _ = w.WriteString("</")
		_, _ = w.WriteString(tag)
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) renderFootnoteList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var h *hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h != nil && h.FootnotesRenderer != nil
	}

	if !ok {
		return r.renderDefaultFootnoteList(w, source, node, entering)
	}

	if entering {
		return ast.WalkContinue, nil
	}

	err := h.FootnotesRenderer.RenderFootnotes(
		w,
		footnotesContext{
			page:      ctx.DocumentContext().Document,
			footnotes: ctx.flushFootnotes(),
		},
	)

	ctx.AddIdentity(h.FootnotesRenderer.GetIdentity())

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
