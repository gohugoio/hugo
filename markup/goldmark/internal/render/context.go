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
	positions []int
	ContextData
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
