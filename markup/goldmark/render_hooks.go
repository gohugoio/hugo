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
	"sync"

	"github.com/spf13/cast"

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

type attributesHolder struct {
	// What we get from Goldmark.
	astAttributes []ast.Attribute

	// What we send to the the render hooks.
	attributesInit sync.Once
	attributes     map[string]string
}

func (a *attributesHolder) Attributes() map[string]string {
	a.attributesInit.Do(func() {
		a.attributes = make(map[string]string)
		for _, attr := range a.astAttributes {
			a.attributes[string(attr.Name)] = string(util.EscapeHTML(attr.Value.([]byte)))
		}
	})
	return a.attributes
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
	*attributesHolder
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
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindHeading, r.renderHeading)
}

func (r *hookedRenderer) renderAttributesForNode(w util.BufWriter, node ast.Node) {
	renderAttributes(w, false, node.Attributes()...)
}

var (

	// Attributes with special meaning that does not make sense to render in HTML.
	attributeExcludes = map[string]bool{
		"hl_lines":    true,
		"hl_style":    true,
		"linenos":     true,
		"linenostart": true,
	}
)

func renderAttributes(w util.BufWriter, skipClass bool, attributes ...ast.Attribute) {
	for _, attr := range attributes {
		if skipClass && bytes.Equal(attr.Name, []byte("class")) {
			continue
		}

		a := strings.ToLower(string(attr.Name))
		if attributeExcludes[a] || strings.HasPrefix(a, "on") {
			continue
		}

		_, _ = w.WriteString(" ")
		_, _ = w.Write(attr.Name)
		_, _ = w.WriteString(`="`)

		switch v := attr.Value.(type) {
		case []byte:
			_, _ = w.Write(util.EscapeHTML(v))
		default:
			w.WriteString(cast.ToString(v))
		}

		_ = w.WriteByte('"')
	}
}

func (r *hookedRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	var h hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h.ImageRenderer != nil
	}

	if !ok {
		return r.renderImageDefault(w, source, node, entering)
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

func (r *hookedRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	var h hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h.LinkRenderer != nil
	}

	if !ok {
		return r.renderLinkDefault(w, source, node, entering)
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
	var h hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h.LinkRenderer != nil
	}

	if !ok {
		return r.renderAutoLinkDefault(w, source, node, entering)
	}

	url := string(n.URL(source))
	label := string(n.Label(source))
	if n.AutoLinkType == ast.AutoLinkEmail && !strings.HasPrefix(strings.ToLower(url), "mailto:") {
		url = "mailto:" + url
	}

	err := h.LinkRenderer.RenderLink(
		w,
		linkContext{
			page:        ctx.DocumentContext().Document,
			destination: url,
			text:        label,
			plainText:   label,
		},
	)

	// TODO(bep) I have a working branch that fixes these rather confusing identity types,
	// but for now it's important that it's not .GetIdentity() that's added here,
	// to make sure we search the entire chain on changes.
	ctx.AddIdentity(h.LinkRenderer)

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
	url := n.URL(source)
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

func (r *hookedRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	var h hooks.Renderers

	ctx, ok := w.(*renderContext)
	if ok {
		h = ctx.RenderContext().RenderHooks
		ok = h.HeadingRenderer != nil
	}

	if !ok {
		return r.renderHeadingDefault(w, source, node, entering)
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
			page:             ctx.DocumentContext().Document,
			level:            n.Level,
			anchor:           string(anchor),
			text:             string(text),
			plainText:        string(n.Text(source)),
			attributesHolder: &attributesHolder{astAttributes: n.Attributes()},
		},
	)

	ctx.AddIdentity(h.HeadingRenderer)

	return ast.WalkContinue, err
}

func (r *hookedRenderer) renderHeadingDefault(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString("<h")
		_ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			r.renderAttributesForNode(w, node)
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
}

// Extend implements goldmark.Extender.
func (e *links) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newLinkRenderer(), 100),
	))
}
