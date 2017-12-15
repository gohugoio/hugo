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

package reflect

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "reflect"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.KindIs,
			nil,
			[][2]string{
				{`{{ reflect.KindIs "map" (dict "one" 1) }}`, "true"},
			},
		)

		ns.AddMethodMapping(ctx.KindOf,
			nil,
			[][2]string{
				{`{{ reflect.KindOf (dict "one" 1) }}`, "map"},
			},
		)

		ns.AddMethodMapping(ctx.TypeIs,
			nil,
			[][2]string{
				{`{{ reflect.TypeIs "map[string]interface {}" (dict "one" 1) }}`, "true"},
			},
		)

		ns.AddMethodMapping(ctx.TypeIsLike,
			nil,
			[][2]string{
				{`{{ reflect.TypeIsLike "map[string]interface {}" (dict "one" 1) }}`, "true"},
			},
		)

		ns.AddMethodMapping(ctx.TypeOf,
			nil,
			[][2]string{
				{`{{ reflect.TypeOf (dict "one" 1) }}`, "map[string]interface {}"},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
