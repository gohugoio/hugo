// Copyright 2019 The Hugo Authors. All rights reserved.
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

package goldmark

import (
	"bytes"

	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var (
	tocResultKey = parser.NewContextKey()
	tocEnableKey = parser.NewContextKey()
)

type tocTransformer struct {
	r renderer.Renderer
}

func (t *tocTransformer) Transform(n *ast.Document, reader text.Reader, pc parser.Context) {
	if b, ok := pc.Get(tocEnableKey).(bool); !ok || !b {
		return
	}

	var (
		toc         tableofcontents.Root
		header      tableofcontents.Header
		level       int
		row         = -1
		inHeading   bool
		headingText bytes.Buffer
	)

	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		if n.Kind() == ast.KindHeading {
			if inHeading && !entering {
				header.Text = headingText.String()
				headingText.Reset()
				toc.AddAt(header, row, level-1)
				header = tableofcontents.Header{}
				inHeading = false
				return s, nil
			}

			inHeading = true
		}

		if !(inHeading && entering) {
			return s, nil
		}

		switch n.Kind() {
		case ast.KindHeading:
			heading := n.(*ast.Heading)
			level = heading.Level

			if level == 1 || row == -1 {
				row++
			}

			id, found := heading.AttributeString("id")
			if found {
				header.ID = string(id.([]byte))
			}
		case
			ast.KindCodeSpan,
			ast.KindLink,
			ast.KindImage,
			ast.KindEmphasis:
			err := t.r.Render(&headingText, reader.Source(), n)
			if err != nil {
				return s, err
			}

			return ast.WalkSkipChildren, nil
		case
			ast.KindAutoLink,
			ast.KindRawHTML,
			ast.KindText,
			ast.KindString:
			err := t.r.Render(&headingText, reader.Source(), n)
			if err != nil {
				return s, err
			}
		}

		return s, nil
	})

	pc.Set(tocResultKey, toc)
}

type tocExtension struct {
	options []renderer.Option
}

func newTocExtension(options []renderer.Option) goldmark.Extender {
	return &tocExtension{
		options: options,
	}
}

func (e *tocExtension) Extend(m goldmark.Markdown) {
	r := goldmark.DefaultRenderer()
	r.AddOptions(e.options...)
	m.Parser().AddOptions(parser.WithASTTransformers(util.Prioritized(&tocTransformer{
		r: r,
	}, 10)))
}
