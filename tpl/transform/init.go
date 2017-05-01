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

package transform

import (
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl/internal"
)

const name = "transform"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		examples := [][2]string{
			{`{{ "I :heart: Hugo" | emojify }}`, `I ❤️ Hugo`},
			{`{{ .Title | markdownify}}`, `<strong>BatMan</strong>`},
			{`{{ plainify  "Hello <strong>world</strong>, gophers!" }}`, `Hello world, gophers!`},
			{
				`htmlEscape 1: {{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>" | safeHTML}}`,
				`htmlEscape 1: Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;`},
			{
				`htmlEscape 2: {{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>"}}`,
				`htmlEscape 2: Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;`},
			{
				`htmlUnescape 1: {{htmlUnescape "Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;" | safeHTML}}`,
				`htmlUnescape 1: Cathal Garvey & The Sunshine Band <cathal@foo.bar>`},
			{
				`htmlUnescape 2: {{"Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;" | htmlUnescape | htmlUnescape | safeHTML}}`,
				`htmlUnescape 2: Cathal Garvey & The Sunshine Band <cathal@foo.bar>`},
			{
				`htmlUnescape 3: {{"Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;" | htmlUnescape | htmlUnescape }}`,
				`htmlUnescape 3: Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;`},
			{
				`htmlUnescape 4: {{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>" | htmlUnescape | safeHTML }}`,
				`htmlUnescape 4: Cathal Garvey & The Sunshine Band <cathal@foo.bar>`},
			{
				`htmlUnescape 5: {{ htmlUnescape "Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;" | htmlEscape | safeHTML }}`,
				`htmlUnescape 5: Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;`},
		}

		return &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func() interface{} { return ctx },
			Aliases: map[string]interface{}{
				"emojify":      ctx.Emojify,
				"highlight":    ctx.Highlight,
				"htmlEscape":   ctx.HTMLEscape,
				"htmlUnescape": ctx.HTMLUnescape,
				"markdownify":  ctx.Markdownify,
				"plainify":     ctx.Plainify,
			},
			Examples: examples,
		}

	}

	internal.AddTemplateFuncsNamespace(f)
}
