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

package time

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "time"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name: name,
			Context: func(args ...interface{}) interface{} {
				// Handle overlapping "time" namespace and func.
				//
				// If no args are passed to `time`, assume namespace usage and
				// return namespace context.
				//
				// If args are passed, call AsTime().

				if len(args) == 0 {
					return ctx
				}

				t, err := ctx.AsTime(args[0])
				if err != nil {
					return err
				}
				return t
			},
		}

		ns.AddMethodMapping(ctx.Format,
			[]string{"dateFormat"},
			[][2]string{
				{`dateFormat: {{ dateFormat "Monday, Jan 2, 2006" "2015-01-21" }}`, `dateFormat: Wednesday, Jan 21, 2015`},
			},
		)

		ns.AddMethodMapping(ctx.Now,
			[]string{"now"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.AsTime,
			nil,
			[][2]string{
				{`{{ (time "2015-01-21").Year }}`, `2015`},
			},
		)

		ns.AddMethodMapping(ctx.Duration,
			[]string{"duration"},
			[][2]string{
				{`{{ mul 60 60 | duration "second" }}`, `1h0m0s`},
			},
		)

		ns.AddMethodMapping(ctx.ParseDuration,
			nil,
			[][2]string{
				{`{{ "1h12m10s" | time.ParseDuration }}`, `1h12m10s`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
