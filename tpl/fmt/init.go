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

package fmt

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "fmt"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Print,
			[]string{"print"},
			[][2]string{
				{`{{ print "works!" }}`, `works!`},
			},
		)

		ns.AddMethodMapping(ctx.Println,
			[]string{"println"},
			[][2]string{
				{`{{ println "works!" }}`, "works!\n"},
			},
		)

		ns.AddMethodMapping(ctx.Printf,
			[]string{"printf"},
			[][2]string{
				{`{{ printf "%s!" "works" }}`, `works!`},
			},
		)

		ns.AddMethodMapping(ctx.Errorf,
			[]string{"errorf"},
			[][2]string{
				{`{{ errorf "%s." "failed" }}`, ``},
			},
		)

		ns.AddMethodMapping(ctx.Warnf,
			[]string{"warnf"},
			[][2]string{
				{`{{ warnf "%s." "warning" }}`, ``},
			},
		)
		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
