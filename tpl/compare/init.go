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

package compare

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "compare"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		language := d.Conf.Language()
		if language == nil {
			panic("language must be set")
		}

		ctx := New(langs.GetLocation(language), false)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Default,
			[]string{"default"},
			[][2]string{
				{`{{ "Hugo Rocks!" | default "Hugo Rules!" }}`, `Hugo Rocks!`},
				{`{{ "" | default "Hugo Rules!" }}`, `Hugo Rules!`},
			},
		)

		ns.AddMethodMapping(ctx.Eq,
			[]string{"eq"},
			[][2]string{
				{`{{ if eq .Section "blog" }}current-section{{ end }}`, `current-section`},
			},
		)

		ns.AddMethodMapping(ctx.Ge,
			[]string{"ge"},
			[][2]string{
				{`{{ if ge hugo.Version "0.80" }}Reasonable new Hugo version!{{ end }}`, `Reasonable new Hugo version!`},
			},
		)

		ns.AddMethodMapping(ctx.Gt,
			[]string{"gt"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Le,
			[]string{"le"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Lt,
			[]string{"lt"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Ne,
			[]string{"ne"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Conditional,
			[]string{"cond"},
			[][2]string{
				{`{{ cond (eq (add 2 2) 4) "2+2 is 4" "what?" | safeHTML }}`, `2+2 is 4`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
