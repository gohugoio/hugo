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

package crypto

import (
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl/internal"
)

const name = "crypto"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		examples := [][2]string{
			{`{{ md5 "Hello world, gophers!" }}`, `b3029f756f98f79e7f1b7f1d1f0dd53b`},
			{`{{ sha1 "Hello world, gophers!" }}`, `c8b5b0e33d408246e30f53e32b8f7627a7a649d4`},
			{`{{ sha256 "Hello world, gophers!" }}`, `6ec43b78da9669f50e4e422575c54bf87536954ccd58280219c393f2ce352b46`},
		}

		return &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func() interface{} { return ctx },
			Aliases: map[string]interface{}{
				"md5":    ctx.MD5,
				"sha1":   ctx.SHA1,
				"sha256": ctx.SHA256,
			},
			Examples: examples,
		}

	}

	internal.AddTemplateFuncsNamespace(f)
}
