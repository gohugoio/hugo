// Copyright 2017 The Hugo Authors. All rights reserved.
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

package safe

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "safe"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.CSS,
			[]string{"safeCSS"},
			[][2]string{
				{`{{ "Bat&Man" | safeCSS | safeCSS }}`, `Bat&amp;Man`},
			},
		)

		ns.AddMethodMapping(ctx.HTML,
			[]string{"safeHTML"},
			[][2]string{
				{`{{ "Bat&Man" | safeHTML | safeHTML }}`, `Bat&Man`},
				{`{{ "Bat&Man" | safeHTML }}`, `Bat&Man`},
			},
		)

		ns.AddMethodMapping(ctx.HTMLAttr,
			[]string{"safeHTMLAttr"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.JS,
			[]string{"safeJS"},
			[][2]string{
				{`{{ "(1*2)" | safeJS | safeJS }}`, `(1*2)`},
			},
		)

		ns.AddMethodMapping(ctx.JSStr,
			[]string{"safeJSStr"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.URL,
			[]string{"safeURL"},
			[][2]string{
				{`{{ "http://gohugo.io" | safeURL | safeURL }}`, `http://gohugo.io`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
