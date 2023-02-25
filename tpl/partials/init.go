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

package partials

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const namespaceName = "partials"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    namespaceName,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Include,
			[]string{"partial"},
			[][2]string{
				{`{{ partial "header.html" . }}`, `<title>Hugo Rocks!</title>`},
			},
		)

		// TODO(bep) we need the return to be a valid identifier, but
		// should consider another way of adding it.
		ns.AddMethodMapping(func() string { return "" },
			[]string{"return"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.IncludeCached,
			[]string{"partialCached"},
			[][2]string{},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
