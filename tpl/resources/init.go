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

package resources

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "resources"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx, err := New(d)
		if err != nil {
			// TODO(bep) no panic.
			panic(err)
		}

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.ByType,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Concat,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Copy,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.ExecuteAsTemplate,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Fingerprint,
			[]string{"fingerprint"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.FromString,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Get,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Match,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.GetRemote,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Match,
			nil,
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Minify,
			[]string{"minify"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.PostProcess,
			nil,
			[][2]string{},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
