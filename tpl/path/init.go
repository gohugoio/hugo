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

package path

import (
	"fmt"
	"path/filepath"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "path"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Split,
			nil,
			[][2]string{
				{`{{ "/my/path/filename.txt" | path.Split }}`, `/my/path/|filename.txt`},
				{fmt.Sprintf(`{{ %q | path.Split }}`, filepath.FromSlash("/my/path/filename.txt")), `/my/path/|filename.txt`},
			},
		)

		ns.AddMethodMapping(ctx.Join,
			nil,
			[][2]string{
				{fmt.Sprintf(`{{ slice %q "filename.txt" | path.Join  }}`, "my"+helpers.FilePathSeparator+"path"), `my/path/filename.txt`},
				{`{{  path.Join "my" "path" "filename.txt" }}`, `my/path/filename.txt`},
			},
		)

		return ns

	}
	internal.AddTemplateFuncsNamespace(f)
}
