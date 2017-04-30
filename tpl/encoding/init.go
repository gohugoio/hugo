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

package encoding

import (
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl/internal"
)

const name = "encoding"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		examples := [][2]string{
			{`{{ (slice "A" "B" "C") | jsonify }}`, `["A","B","C"]`},
			{`{{ "SGVsbG8gd29ybGQ=" | base64Decode }}`, `Hello world`},
			{`{{ 42 | base64Encode | base64Decode }}`, `42`},
			{`{{ "Hello world" | base64Encode }}`, `SGVsbG8gd29ybGQ=`},
		}

		return &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func() interface{} { return ctx },
			Aliases: map[string]interface{}{
				"base64Decode": ctx.Base64Decode,
				"base64Encode": ctx.Base64Encode,
				"jsonify":      ctx.Jsonify,
			},
			Examples: examples,
		}

	}

	internal.AddTemplateFuncsNamespace(f)
}
