// Copyright 2022 The Hugo Authors. All rights reserved.
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
	"fmt"

	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/goldmark/internal/render"
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
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{}, 100),
		),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newHTMLRenderer(), 100),
	))
}

func newHTMLRenderer() renderer.NodeRenderer {
	r := &htmlRenderer{}
	return r
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCodeBlock, r.renderCodeBlock)
}

func (r *htmlRenderer) renderCodeBlock(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)

	if entering {
		return ast.WalkContinue, nil
	}

	n := node.(*codeBlock)
	lang := string(n.b.Language(src))
	ordinal := n.ordinal

	var buff bytes.Buffer

	l := n.b.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.b.Lines().At(i)
		buff.Write(line.Value(src))
	}
	text := buff.String()

	var info []byte
	if n.b.Info != nil {
		info = n.b.Info.Segment.Value(src)
	}
	attrs := getAttributes(n.b, info)

	v := ctx.RenderContext().GetRenderer(hooks.CodeBlockRendererType, lang)
	if v == nil {
		return ast.WalkStop, fmt.Errorf("no code renderer found for %q", lang)
	}

	cr := v.(hooks.CodeBlockRenderer)

	err := cr.RenderCodeblock(
		w,
		codeBlockContext{
			page:             ctx.DocumentContext().Document,
			lang:             lang,
			code:             text,
			ordinal:          ordinal,
			AttributesHolder: attributes.New(attrs, attributes.AttributesOwnerCodeBlock),
		},
	)

	ctx.AddIdentity(cr)

	return ast.WalkContinue, err
}

type codeBlockContext struct {
	page    interface{}
	lang    string
	code    string
	ordinal int
	*attributes.AttributesHolder
}

func (c codeBlockContext) Page() interface{} {
	return c.page
}

func (c codeBlockContext) Lang() string {
	return c.lang
}

func (c codeBlockContext) Code() string {
	return c.code
}

func (c codeBlockContext) Ordinal() int {
	return c.ordinal
}

func getAttributes(node *ast.FencedCodeBlock, infostr []byte) []ast.Attribute {
	if node.Attributes() != nil {
		return node.Attributes()
	}
	if infostr != nil {
		attrStartIdx := -1

		for idx, char := range infostr {
			if char == '{' {
				attrStartIdx = idx
				break
			}
		}

		if attrStartIdx > 0 {
			n := ast.NewTextBlock() // dummy node for storing attributes
			attrStr := infostr[attrStartIdx:]
			if attrs, hasAttr := parser.ParseAttributes(text.NewReader(attrStr)); hasAttr {
				for _, attr := range attrs {
					n.SetAttribute(attr.Name, attr.Value)
				}
				return n.Attributes()
			}
		}
	}
	return nil
}
