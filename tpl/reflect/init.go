// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Package reflect provides template functions for run-time object reflection.
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

		ns.AddMethodMapping(ctx.IsMap,
			nil,
			[][2]string{
				{`{{ if reflect.IsMap (dict "a" 1) }}Map{{ end }}`, `Map`},
			},
		)

		ns.AddMethodMapping(ctx.IsSlice,
			nil,
			[][2]string{
				{`{{ if reflect.IsSlice (slice 1 2 3) }}Slice{{ end }}`, `Slice`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
