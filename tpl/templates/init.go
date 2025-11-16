// Copyright 2018 The Hugo Authors. All rights reserved.
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

package templates

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "templates"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Exists,
			nil,
			[][2]string{
				{`{{ if (templates.Exists "partials/header.html") }}Yes!{{ end }}`, `Yes!`},
				{`{{ if not (templates.Exists "partials/doesnotexist.html") }}No!{{ end }}`, `No!`},
			},
		)

		ns.AddMethodMapping(ctx.Defer,
			nil, // No aliases to keep the AST parsing simple.
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.DoDefer,
			[]string{"doDefer"},
			[][2]string{},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
