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
	"strconv"

	"github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func New() goldmark.Extender {
	return &hugoContextExtension{}
}

// Wrap wraps the given byte slice in a Hugo context that used to determine the correct Page
// in .RenderShortcodes.
func Wrap(b []byte, pid uint64) string {
	buf := bufferpool.GetBuffer()
	defer bufferpool.PutBuffer(buf)
	buf.Write(prefix)
	buf.WriteString(" pid=")
	buf.WriteString(strconv.FormatUint(pid, 10))
	buf.Write(endDelim)
	buf.WriteByte('\n')
	buf.Write(b)
	buf.Write(prefix)
	buf.Write(closingDelimAndNewline)
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
	prefix                 = []byte("{{__hugo_ctx")
	endDelim               = []byte("}}")
	closingDelimAndNewline = []byte("/}}\n")
)

var _ parser.InlineParser = (*hugoContextParser)(nil)

type hugoContextParser struct{}

func (s *hugoContextParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, prefix) {
		return nil
	}
	end := bytes.Index(line, endDelim)
	if end == -1 {
		return nil
	}

	block.Advance(end + len(endDelim) + 1) // +1 for the newline

	if line[end-1] == '/' {
		return &HugoContext{Closing: true}
	}

	attrBytes := line[len(prefix)+1 : end]
	h := &HugoContext{}
	h.parseAttrs(attrBytes)
	return h
}

func (a *hugoContextParser) Trigger() []byte {
	return []byte{'{'}
}

type hugoContextRenderer struct{}

func (r *hugoContextRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindHugoContext, r.handleHugoContext)
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

type hugoContextExtension struct{}

func (a *hugoContextExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&hugoContextParser{}, 50),
		),
	)

	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&hugoContextRenderer{}, 50),
		),
	)
}
