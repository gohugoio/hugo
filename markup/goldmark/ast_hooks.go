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
	"github.com/gohugoio/hugo/markup/converter"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func newASTExtension(cfg converter.ProviderConfig) goldmark.Extender {
	return &astExtension{cfg: cfg}
}

type astExtension struct {
	cfg converter.ProviderConfig
}

func (e *astExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(util.Prioritized(&astTransformer{cfg: e.cfg}, 10)))
}

type astHook interface {
	Done()
	Visit(n ast.Node, entering bool) (ast.WalkStatus, error)
}

type astTransformer struct {
	cfg converter.ProviderConfig
}

func (t *astTransformer) Transform(n *ast.Document, reader text.Reader, pc parser.Context) {
	var hooks []astHook

	if b, ok := pc.Get(renderContextKey).(converter.RenderContext); ok && b.RenderTOC {
		hooks = append(hooks, newTocAstHook(reader, pc))
	}

	if len(hooks) == 0 {
		return
	}

	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)

		for i, hook := range hooks {
			t, err := hook.Visit(n, entering)
			if err != nil {
				return t, err
			}
			if i == 0 || t > s {
				s = t
			}
		}

		return s, nil
	})

	for _, hook := range hooks {
		hook.Done()
	}
}
