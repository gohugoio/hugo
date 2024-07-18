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

package math

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "math"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Abs,
			nil,
			[][2]string{
				{"{{ math.Abs -2.1 }}", "2.1"},
			},
		)

		ns.AddMethodMapping(ctx.Acos,
			nil,
			[][2]string{
				{"{{ math.Acos 1 }}", "0"},
			},
		)

		ns.AddMethodMapping(ctx.Add,
			[]string{"add"},
			[][2]string{
				{"{{ add 1 2 }}", "3"},
			},
		)

		ns.AddMethodMapping(ctx.Asin,
			nil,
			[][2]string{
				{"{{ math.Asin 1 }}", "1.5707963267948966"},
			},
		)

		ns.AddMethodMapping(ctx.Atan,
			nil,
			[][2]string{
				{"{{ math.Atan 1 }}", "0.7853981633974483"},
			},
		)

		ns.AddMethodMapping(ctx.Atan2,
			nil,
			[][2]string{
				{"{{ math.Atan2 1 2 }}", "0.4636476090008061"},
			},
		)

		ns.AddMethodMapping(ctx.Ceil,
			nil,
			[][2]string{
				{"{{ math.Ceil 2.1 }}", "3"},
			},
		)

		ns.AddMethodMapping(ctx.Cos,
			nil,
			[][2]string{
				{"{{ math.Cos 1 }}", "0.5403023058681398"},
			},
		)

		ns.AddMethodMapping(ctx.Div,
			[]string{"div"},
			[][2]string{
				{"{{ div 6 3 }}", "2"},
			},
		)

		ns.AddMethodMapping(ctx.Floor,
			nil,
			[][2]string{
				{"{{ math.Floor 1.9 }}", "1"},
			},
		)

		ns.AddMethodMapping(ctx.Log,
			nil,
			[][2]string{
				{"{{ math.Log 1 }}", "0"},
			},
		)

		ns.AddMethodMapping(ctx.Max,
			nil,
			[][2]string{
				{"{{ math.Max 1 2 }}", "2"},
			},
		)

		ns.AddMethodMapping(ctx.Min,
			nil,
			[][2]string{
				{"{{ math.Min 1 2 }}", "1"},
			},
		)

		ns.AddMethodMapping(ctx.Mod,
			[]string{"mod"},
			[][2]string{
				{"{{ mod 15 3 }}", "0"},
			},
		)

		ns.AddMethodMapping(ctx.ModBool,
			[]string{"modBool"},
			[][2]string{
				{"{{ modBool 15 3 }}", "true"},
			},
		)

		ns.AddMethodMapping(ctx.Mul,
			[]string{"mul"},
			[][2]string{
				{"{{ mul 2 3 }}", "6"},
			},
		)

		ns.AddMethodMapping(ctx.Pi,
			nil,
			[][2]string{
				{"{{ math.Pi }}", "3.141592653589793"},
			},
		)

		ns.AddMethodMapping(ctx.Pow,
			[]string{"pow"},
			[][2]string{
				{"{{ math.Pow 2 3 }}", "8"},
			},
		)

		ns.AddMethodMapping(ctx.Rand,
			nil,
			[][2]string{
				{"{{ math.Rand }}", "0.6312770459590062"},
			},
		)

		ns.AddMethodMapping(ctx.Round,
			nil,
			[][2]string{
				{"{{ math.Round 1.5 }}", "2"},
			},
		)

		ns.AddMethodMapping(ctx.Sin,
			nil,
			[][2]string{
				{"{{ math.Sin 1 }}", "0.8414709848078965"},
			},
		)

		ns.AddMethodMapping(ctx.Sqrt,
			nil,
			[][2]string{
				{"{{ math.Sqrt 81 }}", "9"},
			},
		)

		ns.AddMethodMapping(ctx.Sub,
			[]string{"sub"},
			[][2]string{
				{"{{ sub 3 2 }}", "1"},
			},
		)

		ns.AddMethodMapping(ctx.Tan,
			nil,
			[][2]string{
				{"{{ math.Tan 1 }}", "1.557407724654902"},
			},
		)

		ns.AddMethodMapping(ctx.ToDegrees,
			nil,
			[][2]string{
				{"{{ math.ToDegrees 1.5707963267948966 }}", "90"},
			},
		)

		ns.AddMethodMapping(ctx.ToRadians,
			nil,
			[][2]string{
				{"{{ math.ToRadians 90 }}", "1.5707963267948966"},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
