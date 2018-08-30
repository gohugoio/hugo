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

package os

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "os"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Getenv,
			[]string{"getenv"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.ReadDir,
			[]string{"readDir"},
			[][2]string{
				{`{{ range (readDir "files") }}{{ .Name }}{{ end }}`, "README.txt"},
			},
		)

		ns.AddMethodMapping(ctx.ReadFile,
			[]string{"readFile"},
			[][2]string{
				{`{{ readFile "files/README.txt" }}`, `Hugo Rocks!`},
			},
		)

		ns.AddMethodMapping(ctx.FileExists,
			[]string{"fileExists"},
			[][2]string{
				{`{{ fileExists "foo.txt" }}`, `false`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
