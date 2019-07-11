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

package collections

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "collections"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.After,
			[]string{"after"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Apply,
			[]string{"apply"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Complement,
			[]string{"complement"},
			[][2]string{
				{`{{ slice "a" "b" "c" "d" "e" "f" | complement (slice "b" "c") (slice "d" "e")  }}`, `[a f]`},
			},
		)

		ns.AddMethodMapping(ctx.SymDiff,
			[]string{"symdiff"},
			[][2]string{
				{`{{ slice 1 2 3 | symdiff (slice 3 4) }}`, `[1 2 4]`},
			},
		)

		ns.AddMethodMapping(ctx.Delimit,
			[]string{"delimit"},
			[][2]string{
				{`{{ delimit (slice "A" "B" "C") ", " " and " }}`, `A, B and C`},
			},
		)

		ns.AddMethodMapping(ctx.Dictionary,
			[]string{"dict"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.EchoParam,
			[]string{"echoParam"},
			[][2]string{
				{`{{ echoParam .Params "langCode" }}`, `en`},
			},
		)

		ns.AddMethodMapping(ctx.First,
			[]string{"first"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.KeyVals,
			[]string{"keyVals"},
			[][2]string{
				{`{{ keyVals "key" "a" "b" }}`, `key: [a b]`},
			},
		)

		ns.AddMethodMapping(ctx.In,
			[]string{"in"},
			[][2]string{
				{`{{ if in "this string contains a substring" "substring" }}Substring found!{{ end }}`, `Substring found!`},
			},
		)

		ns.AddMethodMapping(ctx.Index,
			[]string{"index"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Intersect,
			[]string{"intersect"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.IsSet,
			[]string{"isSet", "isset"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Last,
			[]string{"last"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Querify,
			[]string{"querify"},
			[][2]string{
				{
					`{{ (querify "foo" 1 "bar" 2 "baz" "with spaces" "qux" "this&that=those") | safeHTML }}`,
					`bar=2&baz=with+spaces&foo=1&qux=this%26that%3Dthose`},
				{
					`<a href="https://www.google.com?{{ (querify "q" "test" "page" 3) | safeURL }}">Search</a>`,
					`<a href="https://www.google.com?page=3&amp;q=test">Search</a>`},
			},
		)

		ns.AddMethodMapping(ctx.Shuffle,
			[]string{"shuffle"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Slice,
			[]string{"slice"},
			[][2]string{
				{`{{ slice "B" "C" "A" | sort }}`, `[A B C]`},
			},
		)

		ns.AddMethodMapping(ctx.Sort,
			[]string{"sort"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Union,
			[]string{"union"},
			[][2]string{
				{`{{ union (slice 1 2 3) (slice 3 4 5) }}`, `[1 2 3 4 5]`},
			},
		)

		ns.AddMethodMapping(ctx.Where,
			[]string{"where"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Append,
			[]string{"append"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Group,
			[]string{"group"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Seq,
			[]string{"seq"},
			[][2]string{
				{`{{ seq 3 }}`, `[1 2 3]`},
			},
		)

		ns.AddMethodMapping(ctx.NewScratch,
			[]string{"newScratch"},
			[][2]string{
				{`{{ $scratch := newScratch }}{{ $scratch.Add "b" 2 }}{{ $scratch.Add "b" 2 }}{{ $scratch.Get "b" }}`, `4`},
			},
		)

		ns.AddMethodMapping(ctx.Uniq,
			[]string{"uniq"},
			[][2]string{
				{`{{ slice 1 2 3 2 | uniq }}`, `[1 2 3]`},
			},
		)

		ns.AddMethodMapping(ctx.Merge,
			[]string{"merge"},
			[][2]string{
				{`{{ dict "title" "Hugo Rocks!" | collections.Merge (dict "title" "Default Title" "description" "Yes, Hugo Rocks!") | sort }}`,
					`[Yes, Hugo Rocks! Hugo Rocks!]`},
				{`{{  merge (dict "title" "Default Title" "description" "Yes, Hugo Rocks!") (dict "title" "Hugo Rocks!") | sort }}`,
					`[Yes, Hugo Rocks! Hugo Rocks!]`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
