// Copyright 2024 The Hugo Authors. All rights reserved.
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

package bit

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "bit"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.And,
			[]string{"band"},
			[][2]string{
				{"{{ bit.And 0b0011 0b0110 }}", "2"},
			},
		)

		ns.AddMethodMapping(ctx.Clear,
			[]string{"bandn", "bclear"},
			[][2]string{
				{"{{ bit.Clear 0xFF 0x0F }}", "240"},
			},
		)
		
		ns.AddMethodMapping(ctx.Extract,
			nil,
			[][2]string{
				{"{{ bit.Extract 0x170F 4 8 }}", "7"},
			},
		)

		ns.AddMethodMapping(ctx.LeadingZeros,
			[]string{"clz"},
			[][2]string{
				{"{{ bit.LeadingZeros 0b11 }}", "62"},
			},
		)

		ns.AddMethodMapping(ctx.Not,
			[]string{"bnot"},
			[][2]string{
				{"{{ bit.Not 0xFF }}", "-256"},
			},
		)

		ns.AddMethodMapping(ctx.OnesCount,
			[]string{"popcnt"},
			[][2]string{
				{"{{ bit.OnesCount 0b1101011 }}", "5"},
			},
		)

		ns.AddMethodMapping(ctx.Or,
			[]string{"bor"},
			[][2]string{
				{"{{ bit.Or 0b0011 0b0110 }}", "7"},
			},
		)

		ns.AddMethodMapping(ctx.ShiftLeft,
			[]string{"lsl"},
			[][2]string{
				{"{{ bit.ShiftLeft 0b1001011 2 }}", "300"},
			},
		)
		
		ns.AddMethodMapping(ctx.ShiftRight,
			[]string{"asr"},
			[][2]string{
				{"{{ bit.ShiftRight 0b1001011 2 }}", "18"},
			},
		)

		ns.AddMethodMapping(ctx.TrailingZeros,
			[]string{"ctz"},
			[][2]string{
				{"{{ bit.TrailingZeros 0b100000 }}", "5"},
			},
		)

		ns.AddMethodMapping(ctx.Xnor,
			[]string{"bxnor"},
			[][2]string{
				{"{{ bit.Xnor 0b0011 0b0110 }}", "-6"},
			},
		)

		ns.AddMethodMapping(ctx.Xor,
			[]string{"bxor"},
			[][2]string{
				{"{{ bit.Xor 0b0011 0b0110 }}", "5"},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
