// Copyright 2024 The Hugo Authors. All rights reserved.
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

package templates_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestExists(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
index.html: {{ templates.Exists "index.html" }}
post/single.html: {{ templates.Exists "post/single.html" }}
partials/foo.html: {{ templates.Exists "partials/foo.html" }}
partials/doesnotexist.html: {{ templates.Exists "partials/doesnotexist.html" }}
-- layouts/post/single.html --
-- layouts/partials/foo.html --
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
index.html: true
post/single.html: true
partials/foo.html: true
partials/doesnotexist.html: false  
`)
}

func TestExistsWithBaseOf(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/baseof.html --
{{ block "main" . }}{{ end }}
-- layouts/index.html --
{{ define "main" }}
index.html: {{ templates.Exists "index.html" }}
post/single.html: {{ templates.Exists "post/single.html" }}
post/doesnotexist.html: {{ templates.Exists "post/doesnotexist.html" }}
{{ end }}
-- layouts/post/single.html --
{{ define "main" }}MAIN{{ end }}


`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
index.html: true
post/single.html: true
post/doesnotexist.html: false

`)
}

// See #10774
func TestPageFunctionExists(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
Home: {{ page.IsHome }}

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
Home: true

`)
}

func TestTry(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
Home.
{{ $g :=  try ("hello = \"Hello Hugo\"" | transform.Unmarshal)   }}
{{ with $g.Err }}
Err1: {{ . }}
{{ else }}
Value1: {{ $g.Value.hello | safeHTML }}|
{{ end }}
{{ $g :=  try ("hello != \"Hello Hugo\"" | transform.Unmarshal)   }}
{{ with $g.Err }}
Err2: {{ . | safeHTML }}
{{ else }}
Value2: {{ $g.Value.hello | safeHTML }}|
{{ end }}
Try upper: {{ (try ("hello" | upper)).Value }}
Try printf: {{ (try (printf "hello %s" "world")).Value }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"Value1: Hello Hugo|",
		"Err2: template: index.html:",
		"Try upper: HELLO",
		"Try printf: hello world",
	)
}
