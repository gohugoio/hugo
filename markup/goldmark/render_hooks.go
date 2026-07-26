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
	"io"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/goldmark/images"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/internal/attributes"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/util"
)

func newLinkRenderer(cfg goldmark_config.Config) *hookedRenderer {
	r := &hookedRenderer{
		linkifyProtocol: []byte(cfg.Extensions.LinkifyProtocol),
		unsafe:          cfg.Renderer.Unsafe,
		xhtml:           cfg.Renderer.XHTML,
		hardWraps:       cfg.Renderer.HardWraps,
		writer:          html.DefaultWriter,
	}
	return r
}

func newLinks(cfg goldmark_config.Config) *links {
	return &links{cfg: cfg}
}

type linkContext struct {
	hooks.BaseContext
	destination string
	title       string
	text        hstring.HTML
	plainText   string
	*attributes.AttributesHolder
}

func (ctx linkContext) Destination() string {
	return ctx.destination
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
	isBlock bool
}

func (ctx imageLinkContext) IsBlock() bool {
	return ctx.isBlock
}

type headingContext struct {
	hooks.BaseContext
	level     int
	anchor    string
	text      hstring.HTML
	plainText string
	*attributes.AttributesHolder
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
	unsafe          bool
	xhtml           bool
	hardWraps       bool
	writer          html.Writer
}

func (r *hookedRenderer) rendererOptions() []html.Option {
	return []html.Option{
		html.WithNodeRenderer(ast.KindLink, html.NodeRendererFunc(r.renderLink)),
		html.WithNodeRenderer(ast.KindAutoLink, html.NodeRendererFunc(r.renderAutoLink)),
		html.WithNodeRenderer(ast.KindImage, html.NodeRendererFunc(r.renderImage)),
		html.WithNodeRenderer(ast.KindHeading, html.NodeRendererFunc(r.renderHeading)),
	}
}

func (r *hookedRenderer) renderImage(w io.Writer, source []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
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
	// GOLDMARK-V2: attribute values are text.MultilineValue now, so the internal
	// signal attributes set by the images transform are strings we decode here.
	if b, ok := n.Attribute(images.AttrIsBlock); ok && string(b.Bytes(nil)) == "true" {
		isBlock = true
	}
	if o, ok := n.Attribute(images.AttrOrdinal); ok {
		ordinal, _ = strconv.Atoi(string(o.Bytes(nil)))
	}

	// We use the attributes to signal from the parser whether the image is in
	// a block context or not.
	// We may find a better way to do that, but for now, we'll need to remove any
	// internal attributes before rendering.
	attrs := r.filterInternalAttributes(n.Attributes())

	err := lr.RenderLink(
		ctx.RenderContext().Ctx,
		w,
		imageLinkContext{
			linkContext: linkContext{
				BaseContext:      render.NewBaseContext(ctx, lr, node, source, ordinal),
				destination:      string(n.Destination.Bytes(source)),
				title:            string(n.Title.Bytes(source)),
				text:             hstring.HTML(text),
				plainText:        render.TextPlain(n, source),
				AttributesHolder: attributes.New(attrs, attributes.AttributesOwnerGeneral),
			},
			isBlock: isBlock,
		},
	)

	return ast.WalkContinue, err
}

func (r *hookedRenderer) filterInternalAttributes(attrs []ast.Attribute) []ast.Attribute {
	n := 0
	for _, x := range attrs {
		if !strings.HasPrefix(x.Name, internalAttrPrefix) {
			attrs[n] = x
			n++
		}
	}
	return attrs[:n]
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark
func (r *hookedRenderer) renderImageDefault(w io.Writer, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	bw := w.(*render.Context)
	n := node.(*ast.Image)
	_, _ = bw.WriteString("<img src=\"")
	dest := util.URLEscape(n.Destination.Bytes(source), true)
	if r.unsafe || !html.IsDangerousURL(dest) {
		_, _ = bw.Write(util.EscapeHTML(dest))
	}
	_, _ = bw.WriteString(`" alt="`)
	r.renderTexts(bw, source, n)
	_ = bw.WriteByte('"')
	if title := n.Title.Bytes(source); len(title) > 0 {
		_, _ = bw.WriteString(` title="`)
		r.writer.WriteText(bw, title)
		_ = bw.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(bw, n, html.ImageAttributeFilter)
	}
	if r.xhtml {
		_, _ = bw.WriteString(" />")
	} else {
		_, _ = bw.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

func (r *hookedRenderer) renderLink(w io.Writer, source []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
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
	ordinal := ctx.GetAndIncrementOrdinal(node.Kind())

	err := lr.RenderLink(
		ctx.RenderContext().Ctx,
		w,
		linkContext{
			BaseContext:      render.NewBaseContext(ctx, lr, node, source, ordinal),
			destination:      string(n.Destination.Bytes(source)),
			title:            string(n.Title.Bytes(source)),
			text:             hstring.HTML(text),
			plainText:        render.TextPlain(n, source),
			AttributesHolder: attributes.Empty,
		},
	)

	return ast.WalkContinue, err
}

// Borrowed from Goldmark's HTML renderer.
// GOLDMARK-V2: ast.String was removed; textual content is all *ast.Text now.
func (r *hookedRenderer) renderTexts(w io.Writer, source []byte, n ast.Node) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			_, _ = r.renderText(w, source, t, true)
		} else {
			r.renderTexts(w, source, c)
		}
	}
}

// Borrowed from Goldmark's HTML renderer.
func (r *hookedRenderer) renderText(w io.Writer, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	bw := w.(*render.Context)
	n := node.(*ast.Text)
	// GOLDMARK-V2: Text.Segment -> Text.Value; Writer.Write/RawWrite ->
	// WriteText/RawWriteText.
	value := n.Value.Bytes(source)
	r.writer.WriteText(bw, value)
	{
		if n.HardLineBreak() || (n.SoftLineBreak() && r.hardWraps) {
			if r.xhtml {
				_, _ = bw.WriteString("<br />\n")
			} else {
				_, _ = bw.WriteString("<br>\n")
			}
		} else if n.SoftLineBreak() {
			// TODO(bep) we use these methods a fallback to default rendering when no image/link hooks are defined.
			// I don't think the below is relevant in these situations, but if so, we need to create a PR
			// upstream to export softLineBreak.
			/*if r.EastAsianLineBreaks != html.EastAsianLineBreaksNone && len(value) != 0 {
				sibling := node.NextSibling()
				if sibling != nil && sibling.Kind() == ast.KindText {
					if siblingText := sibling.(*ast.Text).Value(source); len(siblingText) != 0 {
						thisLastRune := util.ToRune(value, len(value)-1)
						siblingFirstRune, _ := utf8.DecodeRune(siblingText)
						if r.EastAsianLineBreaks.softLineBreak(thisLastRune, siblingFirstRune) {
							_ = w.WriteByte('\n')
						}
					}
				}
			} else {
				_ = bw.WriteByte('\n')
			}*/
			_ = bw.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

// Fall back to the default Goldmark render funcs. Method below borrowed from:
// https://github.com/yuin/goldmark/blob/b611cd333a492416b56aa8d94b04a67bf0096ab2/renderer/html/html.go#L404
func (r *hookedRenderer) renderLinkDefault(w io.Writer, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	bw := w.(*render.Context)
	n := node.(*ast.Link)
	if entering {
		_, _ = bw.WriteString("<a href=\"")
		dest := util.URLEscape(n.Destination.Bytes(source), true)
		if r.unsafe || !html.IsDangerousURL(dest) {
			_, _ = bw.Write(util.EscapeHTML(dest))
		}
		_ = bw.WriteByte('"')
		if title := n.Title.Bytes(source); len(title) > 0 {
			_, _ = bw.WriteString(` title="`)
			r.writer.WriteText(bw, title)
			_ = bw.WriteByte('"')
		}
		_ = bw.WriteByte('>')
	} else {
		_, _ = bw.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) renderAutoLink(w io.Writer, source []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
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

	// GOLDMARK-V2: AutoLink no longer exposes AutoLinkType/Protocol/URL()/Label().
	// Destination already includes the "mailto:" prefix for emails, so the
	// previous email handling and the linkifyProtocol rewrite are gone (the
	// latter has no v2 equivalent).
	url := string(n.Destination.Bytes(source))
	label := string(n.Label.Bytes(source))

	ordinal := ctx.GetAndIncrementOrdinal(n.Kind())

	err := lr.RenderLink(
		ctx.RenderContext().Ctx,
		w,
		linkContext{
			BaseContext:      render.NewBaseContext(ctx, lr, node, source, ordinal),
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
func (r *hookedRenderer) renderAutoLinkDefault(w io.Writer, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	bw := w.(*render.Context)
	n := node.(*ast.AutoLink)
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = bw.WriteString(`<a href="`)
	// GOLDMARK-V2: Destination already includes any mailto: prefix.
	url := util.URLEscape(n.Destination.Bytes(source), false)
	label := n.Label.Bytes(source)
	if r.unsafe || !html.IsDangerousURL(url) {
		_, _ = bw.Write(util.EscapeHTML(url))
	}
	if n.Attributes() != nil {
		_ = bw.WriteByte('"')
		html.RenderAttributes(bw, n, html.LinkAttributeFilter)
		_ = bw.WriteByte('>')
	} else {
		_, _ = bw.WriteString(`">`)
	}
	_, _ = bw.Write(util.EscapeHTML(label))
	_, _ = bw.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

func (r *hookedRenderer) renderHeading(w io.Writer, source []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
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

	var anchor []byte
	if anchori, ok := n.Attribute("id"); ok {
		anchor = anchori.Bytes(source)
	}
	ordinal := ctx.GetAndIncrementOrdinal(n.Kind())

	err := hr.RenderHeading(
		ctx.RenderContext().Ctx,
		w,
		headingContext{
			BaseContext:      render.NewBaseContext(ctx, hr, node, source, ordinal),
			level:            n.Level,
			anchor:           string(anchor),
			text:             hstring.HTML(text),
			plainText:        render.TextPlain(n, source),
			AttributesHolder: attributes.New(n.Attributes(), attributes.AttributesOwnerGeneral),
		},
	)

	return ast.WalkContinue, err
}

func (r *hookedRenderer) renderHeadingDefault(w io.Writer, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	bw := w.(*render.Context)
	n := node.(*ast.Heading)
	if entering {
		_, _ = bw.WriteString("<h")
		_ = bw.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			attributes.RenderASTAttributes(bw, node.Attributes()...)
		}
		_ = bw.WriteByte('>')
	} else {
		_, _ = bw.WriteString("</h")
		_ = bw.WriteByte("0123456"[n.Level])
		_, _ = bw.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

type links struct {
	cfg goldmark_config.Config
}

// RendererOptions returns the goldmark v2 HTML renderer options for links,
// autolinks, images and headings.
func (e *links) RendererOptions() []html.Option {
	return newLinkRenderer(e.cfg).rendererOptions()
}
