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

package render

import (
	"bytes"
	"math/bits"
	"sync"

	htext "github.com/gohugoio/hugo/common/text"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/yuin/goldmark/ast"
)

type BufWriter struct {
	*bytes.Buffer
}

const maxInt = 1<<(bits.UintSize-1) - 1

func (b *BufWriter) Available() int {
	return maxInt
}

func (b *BufWriter) Buffered() int {
	return b.Len()
}

func (b *BufWriter) Flush() error {
	return nil
}

type Context struct {
	*BufWriter
	ContextData
	positions []int
	pids      []uint64
	ordinals  map[ast.NodeKind]int
	values    map[ast.NodeKind][]any
}

func (ctx *Context) GetAndIncrementOrdinal(kind ast.NodeKind) int {
	if ctx.ordinals == nil {
		ctx.ordinals = make(map[ast.NodeKind]int)
	}
	i := ctx.ordinals[kind]
	ctx.ordinals[kind]++
	return i
}

func (ctx *Context) PushPos(n int) {
	ctx.positions = append(ctx.positions, n)
}

func (ctx *Context) PopPos() int {
	i := len(ctx.positions) - 1
	p := ctx.positions[i]
	ctx.positions = ctx.positions[:i]
	return p
}

func (ctx *Context) PopRenderedString() string {
	pos := ctx.PopPos()
	text := string(ctx.Bytes()[pos:])
	ctx.Truncate(pos)
	return text
}

// PushPid pushes a new page ID to the stack.
func (ctx *Context) PushPid(pid uint64) {
	ctx.pids = append(ctx.pids, pid)
}

// PeekPid returns the current page ID without removing it from the stack.
func (ctx *Context) PeekPid() uint64 {
	if len(ctx.pids) == 0 {
		return 0
	}
	return ctx.pids[len(ctx.pids)-1]
}

// PopPid pops the last page ID from the stack.
func (ctx *Context) PopPid() uint64 {
	if len(ctx.pids) == 0 {
		return 0
	}
	i := len(ctx.pids) - 1
	p := ctx.pids[i]
	ctx.pids = ctx.pids[:i]
	return p
}

func (ctx *Context) PushValue(k ast.NodeKind, v any) {
	if ctx.values == nil {
		ctx.values = make(map[ast.NodeKind][]any)
	}
	ctx.values[k] = append(ctx.values[k], v)
}

func (ctx *Context) PopValue(k ast.NodeKind) any {
	if ctx.values == nil {
		return nil
	}
	v := ctx.values[k]
	if len(v) == 0 {
		return nil
	}
	i := len(v) - 1
	r := v[i]
	ctx.values[k] = v[:i]
	return r
}

func (ctx *Context) PeekValue(k ast.NodeKind) any {
	if ctx.values == nil {
		return nil
	}
	v := ctx.values[k]
	if len(v) == 0 {
		return nil
	}
	return v[len(v)-1]
}

type ContextData interface {
	RenderContext() converter.RenderContext
	DocumentContext() converter.DocumentContext
}

type RenderContextDataHolder struct {
	Rctx converter.RenderContext
	Dctx converter.DocumentContext
}

func (ctx *RenderContextDataHolder) RenderContext() converter.RenderContext {
	return ctx.Rctx
}

func (ctx *RenderContextDataHolder) DocumentContext() converter.DocumentContext {
	return ctx.Dctx
}

// extractSourceSample returns a sample of the source for the given node.
// Note that this is not a copy of the source, but a slice of it,
// so it assumes that the source is not mutated.
func extractSourceSample(n ast.Node, src []byte) []byte {
	var sample []byte

	// Extract a source sample to use for position information.
	if nn := n.FirstChild(); nn != nil {
		var start, stop int
		for i := 0; i < nn.Lines().Len() && i < 2; i++ {
			line := nn.Lines().At(i)
			if i == 0 {
				start = line.Start
			}
			stop = line.Stop
		}
		// We do not mutate the source, so this is safe.
		sample = src[start:stop]
	}
	return sample
}

// GetPageAndPageInner returns the current page and the inner page for the given context.
func GetPageAndPageInner(rctx *Context) (any, any) {
	p := rctx.DocumentContext().Document
	pid := rctx.PeekPid()
	if pid > 0 {
		if lookup := rctx.DocumentContext().DocumentLookup; lookup != nil {
			if v := rctx.DocumentContext().DocumentLookup(pid); v != nil {
				return p, v
			}
		}
	}
	return p, p
}

// NewBaseContext creates a new BaseContext.
func NewBaseContext(rctx *Context, renderer any, n ast.Node, src []byte, getSourceSample func() []byte, ordinal int) hooks.BaseContext {
	if getSourceSample == nil {
		getSourceSample = func() []byte {
			return extractSourceSample(n, src)
		}
	}
	page, pageInner := GetPageAndPageInner(rctx)
	b := &hookBase{
		page:      page,
		pageInner: pageInner,

		getSourceSample: getSourceSample,
		ordinal:         ordinal,
	}

	b.createPos = func() htext.Position {
		if resolver, ok := renderer.(hooks.ElementPositionResolver); ok {
			return resolver.ResolvePosition(b)
		}

		return htext.Position{
			Filename:     rctx.DocumentContext().Filename,
			LineNumber:   1,
			ColumnNumber: 1,
		}
	}

	return b
}

var _ hooks.PositionerSourceTargetProvider = (*hookBase)(nil)

type hookBase struct {
	page      any
	pageInner any
	ordinal   int

	// This is only used in error situations and is expensive to create,
	// so delay creation until needed.
	pos             htext.Position
	posInit         sync.Once
	createPos       func() htext.Position
	getSourceSample func() []byte
}

func (c *hookBase) Page() any {
	return c.page
}

func (c *hookBase) PageInner() any {
	return c.pageInner
}

func (c *hookBase) Ordinal() int {
	return c.ordinal
}

func (c *hookBase) Position() htext.Position {
	c.posInit.Do(func() {
		c.pos = c.createPos()
	})
	return c.pos
}

// For internal use.
func (c *hookBase) PositionerSourceTarget() []byte {
	return c.getSourceSample()
}
