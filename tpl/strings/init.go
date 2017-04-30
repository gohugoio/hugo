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

package strings

import (
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl/internal"
)

const name = "strings"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		examples := [][2]string{
			{`{{chomp "<p>Blockhead</p>\n" }}`, `<p>Blockhead</p>`},
			{
				`{{ findRE "[G|g]o" "Hugo is a static side generator written in Go." "1" }}`,
				`[go]`},
			{`{{ hasPrefix "Hugo" "Hu" }}`, `true`},
			{`{{ hasPrefix "Hugo" "Fu" }}`, `false`},
			{`{{lower "BatMan"}}`, `batman`},
			{
				`{{ replace "Batman and Robin" "Robin" "Catwoman" }}`,
				`Batman and Catwoman`},
			{
				`{{ "http://gohugo.io/docs" | replaceRE "^https?://([^/]+).*" "$1" }}`,
				`gohugo.io`},
			{`{{slicestr "BatMan" 0 3}}`, `Bat`},
			{`{{slicestr "BatMan" 3}}`, `Man`},
			{`{{substr "BatMan" 0 -3}}`, `Bat`},
			{`{{substr "BatMan" 3 3}}`, `Man`},
			{`{{title "Bat man"}}`, `Bat Man`},
			{`{{ trim "++Batman--" "+-" }}`, `Batman`},
			{`{{ "this is a very long text" | truncate 10 " ..." }}`, `this is a ...`},
			{`{{ "With [Markdown](/markdown) inside." | markdownify | truncate 14 }}`, `With <a href="/markdown">Markdown …</a>`},
			{`{{upper "BatMan"}}`, `BATMAN`},
		}

		return &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func() interface{} { return ctx },
			Aliases: map[string]interface{}{
				"chomp":      ctx.Chomp,
				"countrunes": ctx.CountRunes,
				"countwords": ctx.CountWords,
				"findRE":     ctx.FindRE,
				"hasPrefix":  ctx.HasPrefix,
				"lower":      ctx.ToLower,
				"replace":    ctx.Replace,
				"replaceRE":  ctx.ReplaceRE,
				"slicestr":   ctx.SliceString,
				"split":      ctx.Split,
				"substr":     ctx.Substr,
				"title":      ctx.Title,
				"trim":       ctx.Trim,
				"truncate":   ctx.Truncate,
				"upper":      ctx.ToUpper,
			},
			Examples: examples,
		}

	}

	internal.AddTemplateFuncsNamespace(f)
}
