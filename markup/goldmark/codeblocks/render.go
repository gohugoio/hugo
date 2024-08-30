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

package codeblocks

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	htext "github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
	"github.com/gohugoio/hugo/markup/highlight/chromalexers"
	"github.com/gohugoio/hugo/markup/internal/attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type (
	codeBlocksExtension struct{}
	htmlRenderer        struct{}
)

func New() goldmark.Extender {
	return &codeBlocksExtension{}
}

func (e *codeBlocksExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newHTMLRenderer(), 100),
	))
}

func newHTMLRenderer() renderer.NodeRenderer {
	r := &htmlRenderer{}
	return r
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
}

func (r *htmlRenderer) renderCodeBlock(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)

	if entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.FencedCodeBlock)

	lang := getLang(n, src)
	renderer := ctx.RenderContext().GetRenderer(hooks.CodeBlockRendererType, lang)
	if renderer == nil {
		return ast.WalkStop, fmt.Errorf("no code renderer found for %q", lang)
	}

	ordinal := ctx.GetAndIncrementOrdinal(ast.KindFencedCodeBlock)

	var buff bytes.Buffer

	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		buff.Write(line.Value(src))
	}

	s := htext.Chomp(buff.String())

	var info []byte
	if n.Info != nil {
		info = n.Info.Segment.Value(src)
	}

	attrtp := attributes.AttributesOwnerCodeBlockCustom
	if isd, ok := renderer.(hooks.IsDefaultCodeBlockRendererProvider); (ok && isd.IsDefaultCodeBlockRenderer()) || chromalexers.Get(lang) != nil {
		// We say that this is a Chroma code block if it's the default code block renderer
		// or if the language is supported by Chroma.
		attrtp = attributes.AttributesOwnerCodeBlockChroma
	}

	attrs, attrStr, err := getAttributes(n, info)
	if err != nil {
		return ast.WalkStop, &herrors.TextSegmentError{Err: err, Segment: attrStr}
	}

	cbctx := &codeBlockContext{
		BaseContext:      render.NewBaseContext(ctx, renderer, node, src, func() []byte { return []byte(s) }, ordinal),
		lang:             lang,
		code:             s,
		AttributesHolder: attributes.New(attrs, attrtp),
	}

	cr := renderer.(hooks.CodeBlockRenderer)

	err = cr.RenderCodeblock(
		ctx.RenderContext().Ctx,
		w,
		cbctx,
	)
	if err != nil {
		return ast.WalkContinue, herrors.NewFileErrorFromPos(err, cbctx.Position())
	}

	return ast.WalkContinue, nil
}

type codeBlockContext struct {
	hooks.BaseContext
	lang string
	code string

	*attributes.AttributesHolder
}

func (c *codeBlockContext) Type() string {
	return c.lang
}

func (c *codeBlockContext) Inner() string {
	return c.code
}

func getLang(node *ast.FencedCodeBlock, src []byte) string {
	langWithAttributes := string(node.Language(src))
	lang, _, _ := strings.Cut(langWithAttributes, "{")
	return lang
}

func getAttributes(node *ast.FencedCodeBlock, infostr []byte) ([]ast.Attribute, string, error) {
	if node.Attributes() != nil {
		return node.Attributes(), "", nil
	}
	if infostr != nil {
		attrStartIdx := -1
		attrEndIdx := -1

		for idx, char := range infostr {
			if attrEndIdx == -1 && char == '{' {
				attrStartIdx = idx
			}
			if attrStartIdx != -1 && char == '}' {
				attrEndIdx = idx
				break
			}
		}

		if attrStartIdx != -1 && attrEndIdx != -1 {
			n := ast.NewTextBlock() // dummy node for storing attributes
			attrStr := infostr[attrStartIdx : attrEndIdx+1]
			if attrs, hasAttr := parser.ParseAttributes(text.NewReader(attrStr)); hasAttr {
				for _, attr := range attrs {
					n.SetAttribute(attr.Name, attr.Value)
				}
				return n.Attributes(), "", nil
			} else {
				return nil, string(attrStr), errors.New("failed to parse Markdown attributes; you may need to quote the values")
			}
		}
	}
	return nil, "", nil
}
