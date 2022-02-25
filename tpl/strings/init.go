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
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "strings"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		sc := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(ctx context.Context, args ...interface{}) (interface{}, error) { return sc, nil },
		}

		ns.AddMethodMapping(sc.Chomp,
			[]string{"chomp"},
			[][2]string{
				{`{{chomp "<p>Blockhead</p>\n" | safeHTML }}`, `<p>Blockhead</p>`},
			},
		)

		ns.AddMethodMapping(sc.CountRunes,
			[]string{"countrunes"},
			[][2]string{},
		)

		ns.AddMethodMapping(sc.RuneCount,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(sc.CountWords,
			[]string{"countwords"},
			[][2]string{},
		)

		ns.AddMethodMapping(sc.Count,
			nil,
			[][2]string{
				{`{{"aabab" | strings.Count "a" }}`, `3`},
			},
		)

		ns.AddMethodMapping(sc.Contains,
			nil,
			[][2]string{
				{`{{ strings.Contains "abc" "b" }}`, `true`},
				{`{{ strings.Contains "abc" "d" }}`, `false`},
			},
		)

		ns.AddMethodMapping(sc.ContainsAny,
			nil,
			[][2]string{
				{`{{ strings.ContainsAny "abc" "bcd" }}`, `true`},
				{`{{ strings.ContainsAny "abc" "def" }}`, `false`},
			},
		)

		ns.AddMethodMapping(sc.FindRE,
			[]string{"findRE"},
			[][2]string{
				{
					`{{ findRE "[G|g]o" "Hugo is a static side generator written in Go." "1" }}`,
					`[go]`,
				},
			},
		)

		ns.AddMethodMapping(sc.HasPrefix,
			[]string{"hasPrefix"},
			[][2]string{
				{`{{ hasPrefix "Hugo" "Hu" }}`, `true`},
				{`{{ hasPrefix "Hugo" "Fu" }}`, `false`},
			},
		)

		ns.AddMethodMapping(sc.ToLower,
			[]string{"lower"},
			[][2]string{
				{`{{lower "BatMan"}}`, `batman`},
			},
		)

		ns.AddMethodMapping(sc.Replace,
			[]string{"replace"},
			[][2]string{
				{
					`{{ replace "Batman and Robin" "Robin" "Catwoman" }}`,
					`Batman and Catwoman`,
				},
				{
					`{{ replace "aabbaabb" "a" "z" 2 }}`,
					`zzbbaabb`,
				},
			},
		)

		ns.AddMethodMapping(sc.ReplaceRE,
			[]string{"replaceRE"},
			[][2]string{
				{
					`{{ replaceRE "a+b" "X" "aabbaabbab" }}`,
					`XbXbX`,
				},
				{
					`{{ replaceRE "a+b" "X" "aabbaabbab" 1 }}`,
					`Xbaabbab`,
				},
			},
		)

		ns.AddMethodMapping(sc.SliceString,
			[]string{"slicestr"},
			[][2]string{
				{`{{slicestr "BatMan" 0 3}}`, `Bat`},
				{`{{slicestr "BatMan" 3}}`, `Man`},
			},
		)

		ns.AddMethodMapping(sc.Split,
			[]string{"split"},
			[][2]string{},
		)

		ns.AddMethodMapping(sc.Substr,
			[]string{"substr"},
			[][2]string{
				{`{{substr "BatMan" 0 -3}}`, `Bat`},
				{`{{substr "BatMan" 3 3}}`, `Man`},
			},
		)

		ns.AddMethodMapping(sc.Trim,
			[]string{"trim"},
			[][2]string{
				{`{{ trim "++Batman--" "+-" }}`, `Batman`},
			},
		)

		ns.AddMethodMapping(sc.TrimLeft,
			nil,
			[][2]string{
				{`{{ "aabbaa" | strings.TrimLeft "a" }}`, `bbaa`},
			},
		)

		ns.AddMethodMapping(sc.TrimPrefix,
			nil,
			[][2]string{
				{`{{ "aabbaa" | strings.TrimPrefix "a" }}`, `abbaa`},
				{`{{ "aabbaa" | strings.TrimPrefix "aa" }}`, `bbaa`},
			},
		)

		ns.AddMethodMapping(sc.TrimRight,
			nil,
			[][2]string{
				{`{{ "aabbaa" | strings.TrimRight "a" }}`, `aabb`},
			},
		)

		ns.AddMethodMapping(sc.TrimSuffix,
			nil,
			[][2]string{
				{`{{ "aabbaa" | strings.TrimSuffix "a" }}`, `aabba`},
				{`{{ "aabbaa" | strings.TrimSuffix "aa" }}`, `aabb`},
			},
		)

		ns.AddMethodMapping(sc.Title,
			[]string{"title"},
			[][2]string{
				{`{{title "Bat man"}}`, `Bat Man`},
				{`{{title "somewhere over the rainbow"}}`, `Somewhere Over the Rainbow`},
			},
		)

		ns.AddMethodMapping(sc.FirstUpper,
			nil,
			[][2]string{
				{`{{ "hugo rocks!" | strings.FirstUpper }}`, `Hugo rocks!`},
			},
		)

		ns.AddMethodMapping(sc.Truncate,
			[]string{"truncate"},
			[][2]string{
				{`{{ "this is a very long text" | truncate 10 " ..." }}`, `this is a ...`},
				{`{{ "With [Markdown](/markdown) inside." | markdownify | truncate 14 }}`, `With <a href="/markdown">Markdown …</a>`},
			},
		)

		ns.AddMethodMapping(sc.Repeat,
			nil,
			[][2]string{
				{`{{ "yo" | strings.Repeat 4 }}`, `yoyoyoyo`},
			},
		)

		ns.AddMethodMapping(sc.ToUpper,
			[]string{"upper"},
			[][2]string{
				{`{{upper "BatMan"}}`, `BATMAN`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
