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

package inflect

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "inflect"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Humanize,
			[]string{"humanize"},
			[][2]string{
				{`{{ humanize "my-first-post" }}`, `My first post`},
				{`{{ humanize "myCamelPost" }}`, `My camel post`},
				{`{{ humanize "52" }}`, `52nd`},
				{`{{ humanize 103 }}`, `103rd`},
			},
		)

		ns.AddMethodMapping(ctx.Pluralize,
			[]string{"pluralize"},
			[][2]string{
				{`{{ "cat" | pluralize }}`, `cats`},
			},
		)

		ns.AddMethodMapping(ctx.Singularize,
			[]string{"singularize"},
			[][2]string{
				{`{{ "cats" | singularize }}`, `cat`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
