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
	"bytes"
	"strings"

	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/images"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/internal/attributes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

var _ renderer.SetOptioner = (*hookedRenderer)(nil)

func newLinkRenderer(cfg goldmark_config.Config) renderer.NodeRenderer {
	r := &hookedRenderer{
		linkifyProtocol: []byte(cfg.Extensions.LinkifyProtocol),
		Config: html.Config{
			Writer: html.DefaultWriter,
		},
	}
	return r
}

func newLinks(cfg goldmark_config.Config) goldmark.Extender {
	return &links{cfg: cfg}
}

type linkContext struct {
	page        any
	pageInner   any
	destination string
	title       string
	text        hstring.HTML
	plainText   string
	*attributes.AttributesHolder
}

func (ctx linkContext) Destination() string {
	return ctx.destination
}

func (ctx linkContext) Page() any {
	return ctx.page
}

func (ctx linkContext) PageInner() any {
	return ctx.pageInner
}

func (ctx linkContext) Text() hstring.HTML {
	return ctx.text
}

func (ctx linkContext) PlainText() string {
	return ctx.plainText
}

func (ctx linkContext) Title() string {
	return ctx.title
}

type imageLinkContext struct {
	linkContext
	ordinal int
	isBlock bool
}

func (ctx imageLinkContext) IsBlock() bool {
	return ctx.isBlock
}

func (ctx imageLinkContext) Ordinal() int {
	return ctx.ordinal
}

type headingContext struct {
	page      any
	pageInner any
	level     int
	anchor    string
	text      hstring.HTML
	plainText string
	*attributes.AttributesHolder
}

func (ctx headingContext) Page() any {
	return ctx.page
}

func (ctx headingContext) PageInner() any {
	return ctx.pageInner
}

func (ctx headingContext) Level() int {
	return ctx.level
}

func (ctx headingContext) Anchor() string {
	return ctx.anchor
}

func (ctx headingContext) Text() hstring.HTML {
	return ctx.text
}

func (ctx headingContext) PlainText() string {
	return ctx.plainText
}

type hookedRenderer struct {
	linkifyProtocol []byte
	html.Config
}

func (r *hookedRenderer) SetOption(name renderer.OptionName, value any) {
	r.Config.SetOption(name, value)
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *hookedRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindHeading, r.renderHeading)
}

func (r *hookedRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	var lr hooks.LinkRenderer

	ctx, ok := w.(*render.Context)
	if ok {
		h := ctx.RenderContext().GetRenderer(hooks.ImageRendererType, nil)
		ok = h != nil
		if ok {
			lr = h.(hooks.LinkRenderer)
		}
	}

	if !ok {
		return r.renderImageDefault(w, source, node, entering)
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.PushPos(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	text := ctx.PopRenderedString()

	var (
		isBlock bool
		ordinal int
	)
	if b, ok := n.AttributeString(images.AttrIsBlock); ok && b.(bool) {
		isBlock = true
	}
	if n, ok := n.AttributeString(images.AttrOrdinal); ok {
		ordinal = n.(int)
	}

	// We use the attributes to signal from the parser whether the image is in
	// a block context or not.
	// We may find a better way to do that, but for now, we'll need to remove any
	// internal attributes before rendering.
	attrs := r.filterInternalAttributes(n.Attributes())

	page, pageInner := render.GetPageAndPageInner(ctx)

	err := lr.RenderLink(
		ctx.RenderContext().Ctx,
		w,
		imageLinkContext{
			linkContext: linkContext{
				page:             page,
				pageInner:        pageInner,
				destination:      string(n.Destination),
				title:            string(n.Title),
				text:             hstring.HTML(text),
				plainText:        string(n.Text(source)),
				AttributesHolder: attributes.New(attrs, attributes.AttributesOwnerGeneral),
			},
			ordinal: ordinal,
			isBlock: isBlock,
		},
	)

	return ast.WalkContinue, err
}

func (r *hookedRenderer) filterInternalAttributes(attrs []ast.Attribute) []ast.Attribute {
	n := 0
	for _, x := range attrs {
		if !bytes.HasPrefix(x.Name, []byte(internalAttrPrefix)) {
			attrs[n] = x
			n++
		}
	}
	return attrs[:n]
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *hookedRenderer) renderImageDefault(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString("<img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	}
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(nodeToHTMLText(n, source))
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		attrs := r.filterInternalAttributes(n.Attributes())
		attributes.RenderASTAttributes(w, attrs...)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

func (r *hookedRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	var lr hooks.LinkRenderer

	ctx, ok := w.(*render.Context)
	if ok {
		h := ctx.RenderContext().GetRenderer(hooks.LinkRendererType, nil)
		ok = h != nil
		if ok {
			lr = h.(hooks.LinkRenderer)
		}
	}

	if !ok {
		return r.renderLinkDefault(w, source, node, entering)
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.PushPos(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	text := ctx.PopRenderedString()

	page, pageInner := render.GetPageAndPageInner(ctx)

	err := lr.RenderLink(
		ctx.RenderContext().Ctx,
		w,
		linkContext{
			page:             page,
			pageInner:        pageInner,
			destination:      string(n.Destination),
			title:            string(n.Title),
			text:             hstring.HTML(text),
			plainText:        string(n.Text(source)),
			AttributesHolder: attributes.Empty,
		},
	)

	return ast.WalkContinue, err
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *hookedRenderer) renderLinkDefault(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *hookedRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.AutoLink)
	var lr hooks.LinkRenderer

	ctx, ok := w.(*render.Context)
	if ok {
		h := ctx.RenderContext().GetRenderer(hooks.LinkRendererType, nil)
		ok = h != nil
		if ok {
			lr = h.(hooks.LinkRenderer)
		}
	}

	if !ok {
		return r.renderAutoLinkDefault(w, source, node, entering)
	}

	url := string(r.autoLinkURL(n, source))
	label := string(n.Label(source))
	if n.AutoLinkType == ast.AutoLinkEmail && !strings.HasPrefix(strings.ToLower(url), "mailto:") {
		url = "mailto:" + url
	}

	page, pageInner := render.GetPageAndPageInner(ctx)

	err := lr.RenderLink(
		ctx.RenderContext().Ctx,
		w,
		linkContext{
			page:             page,
			pageInner:        pageInner,
			destination:      url,
			text:             hstring.HTML(label),
			plainText:        label,
			AttributesHolder: attributes.Empty,
		},
	)

	return ast.WalkContinue, err
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/5588d92a56fe1642791cf4aa8e9eae8227cfeecd/renderer/html/html.go#L439
func (r *hookedRenderer) renderAutoLinkDefault(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = w.WriteString(`<a href="`)
	url := r.autoLinkURL(n, source)
	label := n.Label(source)
	if n.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		_, _ = w.WriteString("mailto:")
	}
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	if n.Attributes() != nil {
		_ = w.WriteByte('"')
		html.RenderAttributes(w, n, html.LinkAttributeFilter)
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString(`">`)
	}
	_, _ = w.Write(util.EscapeHTML(label))
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) autoLinkURL(n *ast.AutoLink, source []byte) []byte {
	url := n.URL(source)
	if len(n.Protocol) > 0 && !bytes.Equal(n.Protocol, r.linkifyProtocol) {
		// The CommonMark spec says "http" is the correct protocol for links,
		// but this doesn't make much sense (the fact that they should care about the rendered output).
		// Note that n.Protocol is not set if protocol is provided by user.
		url = append(r.linkifyProtocol, url[len(n.Protocol):]...)
	}
	return url
}

func (r *hookedRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	var hr hooks.HeadingRenderer

	ctx, ok := w.(*render.Context)
	if ok {
		h := ctx.RenderContext().GetRenderer(hooks.HeadingRendererType, nil)
		ok = h != nil
		if ok {
			hr = h.(hooks.HeadingRenderer)
		}
	}

	if !ok {
		return r.renderHeadingDefault(w, source, node, entering)
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.PushPos(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	text := ctx.PopRenderedString()

	// All ast.Heading nodes are guaranteed to have an attribute called "id"
	// that is an array of bytes that encode a valid string.
	anchori, _ := n.AttributeString("id")
	anchor := anchori.([]byte)

	page, pageInner := render.GetPageAndPageInner(ctx)

	err := hr.RenderHeading(
		ctx.RenderContext().Ctx,
		w,
		headingContext{
			page:             page,
			pageInner:        pageInner,
			level:            n.Level,
			anchor:           string(anchor),
			text:             hstring.HTML(text),
			plainText:        string(n.Text(source)),
			AttributesHolder: attributes.New(n.Attributes(), attributes.AttributesOwnerGeneral),
		},
	)

	return ast.WalkContinue, err
}

func (r *hookedRenderer) renderHeadingDefault(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString("<h")
		_ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			attributes.RenderASTAttributes(w, node.Attributes()...)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</h")
		_ = w.WriteByte("0123456"[n.Level])
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

type links struct {
	cfg goldmark_config.Config
}

// Extend implements goldmark.Extender.
func (e *links) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newLinkRenderer(e.cfg), 100),
	))
}

// Borrowed from Goldmark.
func nodeToHTMLText(n ast.Node, source []byte) []byte {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if s, ok := c.(*ast.String); ok && s.IsCode() {
			buf.Write(s.Text(source))
		} else if !c.HasChildren() {
			buf.Write(util.EscapeHTML(c.Text(source)))
			if t, ok := c.(*ast.Text); ok && t.SoftLineBreak() {
				buf.WriteByte('\n')
			}
		} else {
			buf.Write(nodeToHTMLText(c, source))
		}
	}
	return buf.Bytes()
}
