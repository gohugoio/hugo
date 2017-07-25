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

package lang

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "lang"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Translate,
			[]string{"i18n", "T"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.NumFmt,
			nil,
			[][2]string{
				{`{{ lang.NumFmt 2 12345.6789 }}`, `12,345.68`},
				{`{{ lang.NumFmt 2 12345.6789 "- , ." }}`, `12.345,68`},
				{`{{ lang.NumFmt 6 -12345.6789 "- ." }}`, `-12345.678900`},
				{`{{ lang.NumFmt 0 -12345.6789 "- . ," }}`, `-12,346`},
				{`{{ -98765.4321 | lang.NumFmt 2 }}`, `-98,765.43`},
			},
		)
		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
