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
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl/internal"
)

const name = "safe"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		examples := [][2]string{
			{`{{ "Bat&Man" | safeCSS | safeCSS }}`, `Bat&amp;Man`},
			{`{{ "Bat&Man" | safeHTML | safeHTML }}`, `Bat&Man`},
			{`{{ "Bat&Man" | safeHTML }}`, `Bat&Man`},
			{`{{ "(1*2)" | safeJS | safeJS }}`, `(1*2)`},
			{`{{ "http://gohugo.io" | safeURL | safeURL }}`, `http://gohugo.io`},
		}

		return &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func() interface{} { return ctx },
			Aliases: map[string]interface{}{
				"safeCSS":      ctx.CSS,
				"safeHTML":     ctx.HTML,
				"safeHTMLAttr": ctx.HTMLAttr,
				"safeJS":       ctx.JS,
				"safeJSStr":    ctx.JSStr,
				"safeURL":      ctx.URL,
				"sanitizeURL":  ctx.SanitizeURL,
				"sanitizeurl":  ctx.SanitizeURL,
			},
			Examples: examples,
		}

	}

	internal.AddTemplateFuncsNamespace(f)
}
