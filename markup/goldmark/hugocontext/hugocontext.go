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
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
}

func (r *hugoContextRenderer) stripHugoCtx(b []byte) ([]byte, bool) {
	if !bytes.Contains(b, hugoCtxPrefix) {
		return b, false
	}
	return hugoCtxRe.ReplaceAll(b, nil), true
}

func (r *hugoContextRenderer) renderHTMLBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	n := node.(*ast.HTMLBlock)
	if entering {
		if r.Unsafe {
			l := n.Lines().Len()
			for i := 0; i < l; i++ {
				line := n.Lines().At(i)
				linev := line.Value(source)
				var stripped bool
				linev, stripped = r.stripHugoCtx(linev)
				if stripped {
					var p any
					ctx, ok := w.(*render.Context)
					if ok {
						p, _ = render.GetPageAndPageInner(ctx)
					}
					r.logger.Warnidf(constants.WarnRenderShortcodesInHTML, ".RenderShortcodes detected inside HTML block in %q; this may not be what you intended, see https://gohugo.io/methods/page/rendershortcodes/#limitations", p)
				}

				r.Writer.SecureWrite(w, linev)
			}
		} else {
			_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
		}
	} else {
		if n.HasClosure() {
			if r.Unsafe {
				closure := n.ClosureLine
				r.Writer.SecureWrite(w, closure.Value(source))
			} else {
				_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
			}
		}
	}
	return ast.WalkContinue, nil
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

type hugoContextExtension struct {
	logger loggers.Logger
}

func (a *hugoContextExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&hugoContextParser{}, 50),
		),
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
