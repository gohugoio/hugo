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

package secretbox

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "secretbox"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Open,
			nil,
			[][2]string{
				{
					`{{ secretbox.Open "KEY" (hexDecode "444f4e27542053454e442041204e4f4e434500000000000036f13ebda2c8537738c4958c367744f3c1b949c9872bb59e5f53706fd6ad") }}`,
					`Secret Message`,
				},
			},
		)

		ns.AddMethodMapping(ctx.Seal,
			nil,
			[][2]string{
				{
					`{{ secretbox.Seal "KEY" "Secret Message" "DON'T SEND A NONCE" | hexEncode }}`,
					`444f4e27542053454e442041204e4f4e434500000000000036f13ebda2c8537738c4958c367744f3c1b949c9872bb59e5f53706fd6ad`,
				},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
