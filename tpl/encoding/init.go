// Copyright 2020 The Hugo Authors. All rights reserved.
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

package encoding

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "encoding"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Base64Decode,
			[]string{"base64Decode"},
			[][2]string{
				{`{{ "SGVsbG8gd29ybGQ=" | base64Decode }}`, `Hello world`},
				{`{{ 42 | base64Encode | base64Decode }}`, `42`},
			},
		)

		ns.AddMethodMapping(ctx.Base64Encode,
			[]string{"base64Encode"},
			[][2]string{
				{`{{ "Hello world" | base64Encode }}`, `SGVsbG8gd29ybGQ=`},
			},
		)

		ns.AddMethodMapping(ctx.Jsonify,
			[]string{"jsonify"},
			[][2]string{
				{`{{ (slice "A" "B" "C") | jsonify }}`, `["A","B","C"]`},
				{`{{ (slice "A" "B" "C") | jsonify (dict "indent" "  ") }}`, "[\n  \"A\",\n  \"B\",\n  \"C\"\n]"},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
