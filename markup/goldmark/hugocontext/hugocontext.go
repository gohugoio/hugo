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

package hugocontext

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"io"

	"github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/common/constants"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

// New returns a goldmark v2 extension (parser + HTML renderer) for the Hugo
// context markers used by .RenderShortcodes.
func New(logger loggers.Logger, unsafe bool) *hugoContextExtension {
	return &hugoContextExtension{logger: logger, unsafe: unsafe}
}

// Wrap wraps the given byte slice in a Hugo context that used to determine the correct Page
// in .RenderShortcodes.
func Wrap(b []byte, pid uint64) string {
	buf := bufferpool.GetBuffer()
	defer bufferpool.PutBuffer(buf)
	buf.Write(hugoCtxPrefix)
	buf.WriteString(" pid=")
	buf.WriteString(strconv.FormatUint(pid, 10))
	buf.Write(hugoCtxEndDelim)
	buf.WriteByte('\n')
	buf.Write(b)
	// To make sure that we're able to parse it, make sure it ends with a newline.
	if len(b) > 0 && b[len(b)-1] != '\n' {
		buf.WriteByte('\n')
	}
	buf.Write(hugoCtxPrefix)
	buf.Write(hugoCtxClosingDelim)
	buf.WriteByte('\n')
	return buf.String()
}

var kindHugoContext = ast.NewNodeKind("HugoContext")

// HugoContext is a node that represents a Hugo context.
type HugoContext struct {
	ast.BaseInline

	Closing bool

	// Internal page ID. Not persisted.
	Pid uint64
}

// Dump implements Node.Dump.
// GOLDMARK-V2: Dump now returns *ast.NodeDump.
func (n *HugoContext) Dump(source []byte) *ast.NodeDump {
	return ast.NewNodeDump(n, map[string]any{
		"Pid": fmt.Sprintf("%v", n.Pid),
	})
}

func (n *HugoContext) parseAttrs(attrBytes []byte) {
	keyPairs := bytes.SplitSeq(attrBytes, []byte(" "))
	for keyPair := range keyPairs {
		kv := bytes.Split(keyPair, []byte("="))
		if len(kv) != 2 {
			continue
		}
		key := string(kv[0])
		val := string(kv[1])
		switch key {
		case "pid":
			pid, _ := strconv.ParseUint(val, 10, 64)
			n.Pid = pid
		}
	}
}

func (h *HugoContext) Kind() ast.NodeKind {
	return kindHugoContext
}

var (
	hugoCtxPrefix       = []byte("{{__hugo_ctx")
	hugoCtxEndDelim     = []byte("}}")
	hugoCtxClosingDelim = []byte("/}}")
	hugoCtxRe           = regexp.MustCompile(`{{__hugo_ctx( pid=\d+)?/?}}\n?`)
	hugoCtxIndentedRe   = regexp.MustCompile(`(?m)^[ \t]+({{__hugo_ctx[^\n]*}})`)
)

// DedentMarkers removes leading whitespace from Hugo context marker lines
// to prevent them from being treated as indented code blocks by Goldmark.
func DedentMarkers(b []byte) []byte {
	if !bytes.Contains(b, hugoCtxPrefix) {
		return b
	}
	return hugoCtxIndentedRe.ReplaceAll(b, []byte("$1"))
}

// Strip strips any Hugo context markers from b.
func Strip(b []byte) []byte {
	if !bytes.Contains(b, hugoCtxPrefix) {
		return b
	}
	return hugoCtxRe.ReplaceAll(b, nil)
}

var _ parser.InlineParser = (*hugoContextParser)(nil)

type hugoContextParser struct{}

func (a *hugoContextParser) Trigger() []byte {
	return []byte{'{'}
}

func (s *hugoContextParser) Parse(parent ast.Node, reader text.Reader, pc parser.Context) ast.Node {
	line, _ := reader.PeekLine()
	if !bytes.HasPrefix(line, hugoCtxPrefix) {
		return nil
	}
	end := bytes.Index(line, hugoCtxEndDelim)
	if end == -1 {
		return nil
	}

	reader.Advance(end + len(hugoCtxEndDelim) + 1) // +1 for the newline

	if line[end-1] == '/' {
		return &HugoContext{Closing: true}
	}

	attrBytes := line[len(hugoCtxPrefix)+1 : end]
	h := &HugoContext{}
	h.parseAttrs(attrBytes)
	return h
}

type hugoContextRenderer struct {
	logger loggers.Logger
	unsafe bool
	writer html.Writer
}

func (r *hugoContextRenderer) rendererOptions() []html.Option {
	return []html.Option{
		html.WithNodeRenderer(kindHugoContext, html.NodeRendererFunc(r.handleHugoContext)),
		html.WithNodeRenderer(ast.KindRawHTML, html.NodeRendererFunc(r.renderRawHTML)),
		html.WithNodeRenderer(ast.KindHTMLBlock, html.NodeRendererFunc(r.renderHTMLBlock)),
	}
}

func (r *hugoContextRenderer) stripHugoCtx(b []byte) ([]byte, bool) {
	if !bytes.Contains(b, hugoCtxPrefix) {
		return b, false
	}
	return hugoCtxRe.ReplaceAll(b, nil), true
}

func (r *hugoContextRenderer) logRawHTMLEmittedWarn(w io.Writer) {
	r.logger.Warnidf(constants.WarnGoldmarkRawHTML, "Raw HTML omitted while rendering %q; see https://gohugo.io/getting-started/configuration-markup/#rendererunsafe", r.getPage(w))
}

func (r *hugoContextRenderer) getPage(w io.Writer) any {
	var p any
	ctx, ok := w.(*render.Context)
	if ok {
		p, _ = render.GetPageAndPageInner(ctx)
	}
	return p
}

func (r *hugoContextRenderer) isHTMLComment(b []byte) bool {
	return len(b) > 4 && b[0] == '<' && b[1] == '!' && b[2] == '-' && b[3] == '-'
}

// HTML rendering based on Goldmark implementation.
// GOLDMARK-V2: HTMLBlock.Lines()/ClosureLine/HasClosure are gone; the whole
// block content (including any closing line) is now in HTMLBlock.Value.
func (r *hugoContextRenderer) renderHTMLBlock(
	w io.Writer, source []byte, node ast.Node, entering bool, _ renderer.Context,
) (ast.WalkStatus, error) {
	bw := w.(*render.Context)
	n := node.(*ast.HTMLBlock)

	if entering {
		v := n.Value.Bytes(source)
		if r.unsafe {
			var stripped bool
			v, stripped = r.stripHugoCtx(v)
			if stripped {
				r.logger.Warnidf(constants.WarnRenderShortcodesInHTML, ".RenderShortcodes detected inside HTML block in %q; this may not be what you intended, see https://gohugo.io/methods/page/rendershortcodes/#limitations", r.getPage(w))
			}
			r.writer.WriteHTML(bw, v)
		} else {
			if !r.isHTMLComment(v) {
				r.logRawHTMLEmittedWarn(w)
				_, _ = bw.WriteString("<!-- raw HTML omitted -->\n")
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *hugoContextRenderer) renderRawHTML(
	w io.Writer, source []byte, node ast.Node, entering bool, _ renderer.Context,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}
	bw := w.(*render.Context)
	n := node.(*ast.RawHTML)
	// GOLDMARK-V2: RawHTML.Segments -> RawHTML.Value (text.MultilineValue).
	v := n.Value.Bytes(source)
	if r.unsafe {
		_, _ = bw.Write(v)
		return ast.WalkSkipChildren, nil
	}
	if !r.isHTMLComment(v) {
		r.logRawHTMLEmittedWarn(w)
		_, _ = bw.WriteString("<!-- raw HTML omitted -->")
	}
	return ast.WalkSkipChildren, nil
}

func (r *hugoContextRenderer) handleHugoContext(w io.Writer, source []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	hctx := node.(*HugoContext)
	ctx, ok := w.(*render.Context)
	if !ok {
		return ast.WalkContinue, nil
	}
	if hctx.Closing {
		_ = ctx.PopPid()
	} else {
		ctx.PushPid(hctx.Pid)
	}
	return ast.WalkContinue, nil
}

type hugoContextTransformer struct{}

var _ parser.ASTTransformer = (*hugoContextTransformer)(nil)

func (a *hugoContextTransformer) Transform(n *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkContinue
		if !entering || n.Kind() != kindHugoContext {
			return s, nil
		}

		if p, ok := n.Parent().(*ast.Paragraph); ok {
			if p.ChildCount() == 1 {
				// Avoid empty paragraphs.
				// GOLDMARK-V2: ReplaceChild/RemoveChild no longer take the receiver.
				p.Parent().ReplaceChild(p, n)
			} else {
				if t, ok := n.PreviousSibling().(*ast.Text); ok {
					// Remove the newline produced by the Hugo context markers.
					if t.SoftLineBreak() {
						// GOLDMARK-V2: Text.Segment -> Text.Value.
						if len(t.Value.Bytes(reader.Source())) == 0 {
							p.RemoveChild(t)
						} else {
							t.SetSoftLineBreak(false)
						}
					}
				}
			}
		}

		return s, nil
	})
}

type hugoContextExtension struct {
	logger loggers.Logger
	unsafe bool
}

func (a *hugoContextExtension) ParserOptions() []parser.Option {
	return []parser.Option{
		parser.WithInlineParsers(
			util.Prioritized[parser.InlineParser](&hugoContextParser{}, 50),
		),
		parser.WithASTTransformers(util.Prioritized[parser.ASTTransformer](&hugoContextTransformer{}, 10)),
	}
}

func (a *hugoContextExtension) RendererOptions() []html.Option {
	r := &hugoContextRenderer{
		logger: a.logger,
		unsafe: a.unsafe,
		writer: html.DefaultWriter,
	}
	return r.rendererOptions()
}
