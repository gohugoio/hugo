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
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "transform"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Emojify,
			[]string{"emojify"},
			[][2]string{
				{`{{ "I :heart: Hugo" | emojify }}`, `I ❤ Hugo`},
			},
		)

		ns.AddMethodMapping(ctx.Highlight,
			[]string{"highlight"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.HTMLEscape,
			[]string{"htmlEscape"},
			[][2]string{
				{
					`{{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>" | safeHTML}}`,
					`Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;`},
				{
					`{{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>"}}`,
					`Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;`},
				{
					`{{ htmlEscape "Cathal Garvey & The Sunshine Band <cathal@foo.bar>" | htmlUnescape | safeHTML }}`,
					`Cathal Garvey & The Sunshine Band <cathal@foo.bar>`},
			},
		)

		ns.AddMethodMapping(ctx.HTMLUnescape,
			[]string{"htmlUnescape"},
			[][2]string{
				{
					`{{ htmlUnescape "Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;" | safeHTML}}`,
					`Cathal Garvey & The Sunshine Band <cathal@foo.bar>`},
				{
					`{{"Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;" | htmlUnescape | htmlUnescape | safeHTML}}`,
					`Cathal Garvey & The Sunshine Band <cathal@foo.bar>`},
				{
					`{{"Cathal Garvey &amp;amp; The Sunshine Band &amp;lt;cathal@foo.bar&amp;gt;" | htmlUnescape | htmlUnescape }}`,
					`Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;`},
				{
					`{{ htmlUnescape "Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;" | htmlEscape | safeHTML }}`,
					`Cathal Garvey &amp; The Sunshine Band &lt;cathal@foo.bar&gt;`},
			},
		)

		ns.AddMethodMapping(ctx.Markdownify,
			[]string{"markdownify"},
			[][2]string{
				{`{{ .Title | markdownify}}`, `<strong>BatMan</strong>`},
			},
		)

		ns.AddMethodMapping(ctx.Plainify,
			[]string{"plainify"},
			[][2]string{
				{`{{ plainify  "Hello <strong>world</strong>, gophers!" }}`, `Hello world, gophers!`},
			},
		)

		ns.AddMethodMapping(ctx.Remarshal,
			nil,
			[][2]string{
				{`{{ "title = \"Hello World\"" | transform.Remarshal "json" | safeHTML }}`, "{\n   \"title\": \"Hello World\"\n}\n"},
			},
		)

		ns.AddMethodMapping(ctx.Unmarshal,
			[]string{"unmarshal"},
			[][2]string{
				{`{{ "hello = \"Hello World\"" | transform.Unmarshal }}`, "map[hello:Hello World]"},
				{`{{ "hello = \"Hello World\"" | resources.FromString "data/greetings.toml" | transform.Unmarshal }}`, "map[hello:Hello World]"},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
