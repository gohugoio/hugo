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

package debug

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "debug"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Dump,
			nil,
			[][2]string{
				{`{{ $m := newScratch }}
{{ $m.Set "Hugo" "Rocks!" }}
{{ $m.Values | debug.Dump | safeHTML }}`, "{\n  \"Hugo\": \"Rocks!\"\n}"},
			},
		)

		ns.AddMethodMapping(ctx.List,
			nil,
			[][2]string{
				{`{{ $s := dict "name" "Hugo" "type" "SSG" }}
{{ debug.List $s }}`, `[name type]`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
