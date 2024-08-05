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

	"github.com/gohugoio/hugo/markup/converter"
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
