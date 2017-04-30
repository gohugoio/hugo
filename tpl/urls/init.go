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

package urls

import (
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl/internal"
)

const name = "urls"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		examples := [][2]string{
			{`{{ "index.html" | absLangURL }}`, `http://mysite.com/hugo/en/index.html`},
			{`{{ "http://gohugo.io/" | absURL }}`, `http://gohugo.io/`},
			{`{{ "mystyle.css" | absURL }}`, `http://mysite.com/hugo/mystyle.css`},
			{`{{ 42 | absURL }}`, `http://mysite.com/hugo/42`},
			{`{{ "index.html" | relLangURL }}`, `/hugo/en/index.html`},
			{`{{ "http://gohugo.io/" | relURL }}`, `http://gohugo.io/`},
			{`{{ "mystyle.css" | relURL }}`, `/hugo/mystyle.css`},
			{`{{ mul 2 21 | relURL }}`, `/hugo/42`},
		}

		return &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func() interface{} { return ctx },
			Aliases: map[string]interface{}{
				"absURL":     ctx.AbsURL,
				"absLangURL": ctx.AbsLangURL,
				"ref":        ctx.Ref,
				"relURL":     ctx.RelURL,
				"relLangURL": ctx.RelLangURL,
				"relref":     ctx.RelRef,
			},
			Examples: examples,
		}

	}

	internal.AddTemplateFuncsNamespace(f)
}
