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

	"github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/common/constants"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func New(logger loggers.Logger) goldmark.Extender {
	return &hugoContextExtension{logger: logger}
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
func (n *HugoContext) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Pid"] = fmt.Sprintf("%v", n.Pid)
	ast.DumpHelper(n, source, level, m, nil)
}

func (n *HugoContext) parseAttrs(attrBytes []byte) {
	keyPairs := bytes.Split(attrBytes, []byte(" "))
	for _, keyPair := range keyPairs {
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
)

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
	html.Config
}

func (r *hugoContextRenderer) SetOption(name renderer.OptionName, value any) {
	r.Config.SetOption(name, value)
}

func (r *hugoContextRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindHugoContext, r.handleHugoContext)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
}

func (r *hugoContextRenderer) stripHugoCtx(b []byte) ([]byte, bool) {
	if !bytes.Contains(b, hugoCtxPrefix) {
		return b, false
	}
	return hugoCtxRe.ReplaceAll(b, nil), true
}

func (r *hugoContextRenderer) logRawHTMLEmittedWarn(w util.BufWriter) {
	r.logger.Warnidf(constants.WarnGoldmarkRawHTML, "Raw HTML omitted while rendering %q; see https://gohugo.io/getting-started/configuration-markup/#rendererunsafe", r.getPage(w))
}

func (r *hugoContextRenderer) getPage(w util.BufWriter) any {
	var p any
	ctx, ok := w.(*render.Context)
	if ok {
		p, _ = render.GetPageAndPageInner(ctx)
	}
	return p
}

// HTML rendering based on Goldmark implementation.
func (r *hugoContextRenderer) renderHTMLBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	n := node.(*ast.HTMLBlock)
	isHTMLComment := func(b []byte) bool {
		return len(b) > 4 && b[0] == '<' && b[1] == '!' && b[2] == '-' && b[3] == '-'
	}
	if entering {
		if r.Unsafe {
			l := n.Lines().Len()
			for i := 0; i < l; i++ {
				line := n.Lines().At(i)
				linev := line.Value(source)
				var stripped bool
				linev, stripped = r.stripHugoCtx(linev)
				if stripped {
					r.logger.Warnidf(constants.WarnRenderShortcodesInHTML, ".RenderShortcodes detected inside HTML block in %q; this may not be what you intended, see https://gohugo.io/methods/page/rendershortcodes/#limitations", r.getPage(w))
				}
				r.Writer.SecureWrite(w, linev)
			}
		} else {
			l := n.Lines().At(0)
			v := l.Value(source)
			if !isHTMLComment(v) {
				r.logRawHTMLEmittedWarn(w)
				_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
			}
		}
	} else {
		if n.HasClosure() {
			if r.Unsafe {
				closure := n.ClosureLine
				r.Writer.SecureWrite(w, closure.Value(source))
			} else {
				l := n.Lines().At(0)
				v := l.Value(source)
				if !isHTMLComment(v) {
					_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
				}
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *hugoContextRenderer) renderRawHTML(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}
	if r.Unsafe {
		n := node.(*ast.RawHTML)
		l := n.Segments.Len()
		for i := 0; i < l; i++ {
			segment := n.Segments.At(i)
			_, _ = w.Write(segment.Value(source))
		}
		return ast.WalkSkipChildren, nil
	}
	r.logRawHTMLEmittedWarn(w)
	_, _ = w.WriteString("<!-- raw HTML omitted -->")
	return ast.WalkSkipChildren, nil
}

func (r *hugoContextRenderer) handleHugoContext(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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
				p.Parent().ReplaceChild(p.Parent(), p, n)
			} else {
				if t, ok := n.PreviousSibling().(*ast.Text); ok {
					// Remove the newline produced by the Hugo context markers.
					if t.SoftLineBreak() {
						if t.Segment.Len() == 0 {
							p.RemoveChild(p, t)
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
}

func (a *hugoContextExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&hugoContextParser{}, 50),
		),
		parser.WithASTTransformers(util.Prioritized(&hugoContextTransformer{}, 10)),
	)

	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&hugoContextRenderer{
				logger: a.logger,
				Config: html.Config{
					Writer: html.DefaultWriter,
				},
			}, 50),
		),
	)
}
