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
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "math"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Add,
			[]string{"add"},
			[][2]string{
				{"{{add 1 2}}", "3"},
			},
		)

		ns.AddMethodMapping(ctx.Ceil,
			nil,
			[][2]string{
				{"{{math.Ceil 2.1}}", "3"},
			},
		)

		ns.AddMethodMapping(ctx.Div,
			[]string{"div"},
			[][2]string{
				{"{{div 6 3}}", "2"},
			},
		)

		ns.AddMethodMapping(ctx.Floor,
			nil,
			[][2]string{
				{"{{math.Floor 1.9}}", "1"},
			},
		)

		ns.AddMethodMapping(ctx.Log,
			nil,
			[][2]string{
				{"{{math.Log 1}}", "0"},
			},
		)

		ns.AddMethodMapping(ctx.Sqrt,
			nil,
			[][2]string{
				{"{{math.Sqrt 81}}", "9"},
			},
		)

		ns.AddMethodMapping(ctx.Mod,
			[]string{"mod"},
			[][2]string{
				{"{{mod 15 3}}", "0"},
			},
		)

		ns.AddMethodMapping(ctx.ModBool,
			[]string{"modBool"},
			[][2]string{
				{"{{modBool 15 3}}", "true"},
			},
		)

		ns.AddMethodMapping(ctx.Mul,
			[]string{"mul"},
			[][2]string{
				{"{{mul 2 3}}", "6"},
			},
		)

		ns.AddMethodMapping(ctx.Round,
			nil,
			[][2]string{
				{"{{math.Round 1.5}}", "2"},
			},
		)

		ns.AddMethodMapping(ctx.Sub,
			[]string{"sub"},
			[][2]string{
				{"{{sub 3 2}}", "1"},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
