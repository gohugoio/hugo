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
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "lang"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d, langs.GetTranslator(d.Conf.Language()))

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Translate,
			[]string{"i18n", "T"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.FormatNumber,
			nil,
			[][2]string{
				{`{{ 512.5032 | lang.FormatNumber 2 }}`, `512.50`},
			},
		)

		ns.AddMethodMapping(ctx.FormatPercent,
			nil,
			[][2]string{
				{`{{ 512.5032 | lang.FormatPercent 2 }}`, `512.50%`},
			},
		)

		ns.AddMethodMapping(ctx.FormatCurrency,
			nil,
			[][2]string{
				{`{{ 512.5032 | lang.FormatCurrency 2 "USD" }}`, `$512.50`},
			},
		)

		ns.AddMethodMapping(ctx.FormatAccounting,
			nil,
			[][2]string{
				{`{{ 512.5032 | lang.FormatAccounting 2 "NOK" }}`, `NOK512.50`},
			},
		)

		ns.AddMethodMapping(ctx.FormatNumberCustom,
			nil,
			[][2]string{
				{`{{ lang.FormatNumberCustom 2 12345.6789 }}`, `12,345.68`},
				{`{{ lang.FormatNumberCustom 2 12345.6789 "- , ." }}`, `12.345,68`},
				{`{{ lang.FormatNumberCustom 6 -12345.6789 "- ." }}`, `-12345.678900`},
				{`{{ lang.FormatNumberCustom 0 -12345.6789 "- . ," }}`, `-12,346`},
				{`{{ lang.FormatNumberCustom 0 -12345.6789 "-|.| " "|" }}`, `-12 346`},
				{`{{ -98765.4321 | lang.FormatNumberCustom 2 }}`, `-98,765.43`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
