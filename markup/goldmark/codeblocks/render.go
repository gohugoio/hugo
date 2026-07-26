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
	"github.com/gohugoio/hugo/markup/internal/attributes"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"

	"io"
)

type (
	codeBlocksExtension struct{}
	htmlRenderer        struct{}
)

// New returns a goldmark v2 HTML renderer extension for fenced code blocks.
func New() html.Extension {
	return &codeBlocksExtension{}
}

func (e *codeBlocksExtension) RendererOptions(*html.Config) []html.Option {
	r := &htmlRenderer{}
	// GOLDMARK-V2: FencedCodeBlock was merged into CodeBlock; we register
	// KindCodeBlock and dispatch on CodeBlockKind below.
	return []html.Option{
		html.WithNodeRenderer(ast.KindCodeBlock, html.NodeRendererFunc(r.renderCodeBlock)),
	}
}

func (r *htmlRenderer) renderCodeBlock(w io.Writer, src []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)

	n := node.(*ast.CodeBlock)

	// GOLDMARK-V2: Only fenced blocks get Hugo's highlighting; indented code
	// blocks fall back to a default rendering.
	if n.CodeBlockKind != ast.CodeBlockKindFenced {
		return renderIndentedCodeBlockDefault(ctx, src, n, entering)
	}

	if entering {
		return ast.WalkContinue, nil
	}

	lang := getLang(n, src)
	renderer := ctx.RenderContext().GetRenderer(hooks.CodeBlockRendererType, lang)
	if renderer == nil {
		return ast.WalkStop, fmt.Errorf("no code renderer found for %q", lang)
	}

	ordinal := ctx.GetAndIncrementOrdinal(ast.KindCodeBlock)

	var buff bytes.Buffer
	buff.Write(n.Value.Bytes(src))

	s := htext.Chomp(buff.String())

	info := n.Info.Bytes(src)

	attrs, attrStr, err := getAttributes(n, info)
	if err != nil {
		return ast.WalkStop, &herrors.TextSegmentError{Err: err, Segment: attrStr}
	}

	cbctx := codeBlockContext{
		BaseContext:      render.NewBaseContext(ctx, renderer, node, src, ordinal),
		lang:             lang,
		code:             s,
		AttributesHolder: attributes.New(attrs, attributes.AttributesOwnerCodeBlockChroma),
	}

	cr := renderer.(hooks.CodeBlockRenderer)

	err = cr.RenderCodeblock(
		ctx.RenderContext().Ctx,
		ctx,
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

func (c codeBlockContext) Type() string {
	return c.lang
}

func (c codeBlockContext) Inner() string {
	return c.code
}

// renderIndentedCodeBlockDefault renders an indented code block the same way
// goldmark's default HTML renderer does.
// GOLDMARK-V2: needed because fenced and indented code blocks now share
// KindCodeBlock and we override the whole kind.
func renderIndentedCodeBlockDefault(ctx *render.Context, src []byte, n *ast.CodeBlock, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = ctx.WriteString("<pre><code>")
		html.DefaultWriter.WriteText(ctx, n.Value.Bytes(src))
	} else {
		_, _ = ctx.WriteString("</code></pre>\n")
	}
	return ast.WalkContinue, nil
}

func getLang(node *ast.CodeBlock, src []byte) string {
	langValue, _ := node.Language(src)
	langWithAttributes := string(langValue.Bytes(src))
	lang, _, _ := strings.Cut(langWithAttributes, "{")
	return lang
}

func getAttributes(node *ast.CodeBlock, infostr []byte) ([]ast.Attribute, string, error) {
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
			attrStr := infostr[attrStartIdx : attrEndIdx+1]
			// GOLDMARK-V2: parser.ParseAttributes now returns []ast.Attribute
			// directly (values are text.MultilineValue), so we no longer need a
			// dummy node to collect them.
			if attrs, hasAttr := parser.ParseAttributes(text.NewReader(attrStr)); hasAttr {
				return attrs, "", nil
			} else {
				return nil, string(attrStr), errors.New("failed to parse Markdown attributes; you may need to quote the values")
			}
		}
	}
	return nil, "", nil
}
