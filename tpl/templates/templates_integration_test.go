// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
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

func TestTemplatesCurrent(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/baseof.html --
baseof: {{ block "main" . }}{{ end }}
-- layouts/all.html --
{{ define "main" }}
all.current: {{ templates.Current.Name }}
all.current.filename: {{ templates.Current.Filename }}
all.base: {{ with templates.Current.Base }}{{ .Name }}{{ end }}|
all.parent: {{ with .Parent }}Name: {{ .Name }}{{ end }}|
{{ partial "p1.html" . }}
{{ end }}
-- layouts/_partials/p1.html --
p1.current: {{ with templates.Current }}Name: {{ .Name }}|{{ with .Parent }}Parent.Name: {{ .Name }}{{ end }}{{ end }}|
p1.current.Ancestors: {{ with templates.Current }}{{ range .Ancestors }}{{ .Name }}|{{ end }}{{ end }}
{{ partial "p2.html" . }}
-- layouts/_partials/p2.html --
p2.current: {{ with templates.Current }}Name: {{ .Name }}|{{ with .Parent }}Parent.Name: {{ .Name }}{{ end }}{{ end }}|
p2.current.Ancestors: {{ with templates.Current }}{{ range .Ancestors }}{{ .Name }}|{{ end }}{{ end }}
p3.current.Ancestors.Reverse: {{ with templates.Current }}{{ range .Ancestors.Reverse }}{{ .Name }}|{{ end }}{{ end }}

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"all.current: all.html",
		filepath.FromSlash("all.current.filename: /layouts/all.html"),
		"all.base: baseof.html",
		"all.parent: |",
		"p1.current: Name: _partials/p1.html|Parent.Name: all.html|",
		"p1.current.Ancestors: all.html|",
		"p2.current.Ancestors: _partials/p1.html|all.html",
	)
}

func TestBaseOfIssue13583(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/_index.md --
---
title: "Home"
outputs: ["html", "amp"]
---
title: "Home"
-- layouts/baseof.html --
layouts/baseof.html
{{ block "main" . }}{{ end }}
-- layouts/baseof.amp.html --
layouts/baseof.amp.html
{{ block "main" . }}{{ end }}
-- layouts/home.html --
{{ define "main" }}
Home.
{{ end }}

`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "layouts/baseof.html")
	b.AssertFileContent("public/amp/index.html", "layouts/baseof.amp.html")
}

func TestAllVsAmp(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/_index.md --
---
title: "Home"
outputs: ["html", "amp"]
---
title: "Home"
-- layouts/all.html --
All.

`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "All.")
	b.AssertFileContent("public/amp/index.html", "All.")
}

// Issue #13584.
func TestLegacySectionSection(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/mysection/_index.md --
-- layouts/section/section.html --
layouts/section/section.html

`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/mysection/index.html", "layouts/section/section.html")
}

func TestErrorMessageParseError(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/home.html --
Line 1.
Line 2. {{ foo }} <- this func does not exist.
Line 3.
`

	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, filepath.FromSlash(`"/layouts/home.html:2:1": parse of template failed: template: home.html:2: function "foo" not defined`))
}

func TestErrorMessageExecuteError(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/home.html --
Line 1.
Line 2. {{ .Foo }} <- this method does not exist.
Line 3.
`

	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, filepath.FromSlash(` "/layouts/home.html:2:11": execute of template failed`))
}

func TestPartialReturnPanicIssue13600(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/home.html --
Partial: {{ partial "p1.html" . }}
-- layouts/_partials/p1.html --
P1.
{{ return ( delimit . ", " ) | string }}
`

	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "wrong number of args for string: want 1 got 0")
}

func TestPartialWithoutSuffixIssue13601(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/home.html --
P1: {{ partial "p1" . }}
P2: {{ partial "p2" . }}
-- layouts/_partials/p1 --
P1.
-- layouts/_partials/p2 --
P2.
{{ return "foo bar" }}

`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "P1: P1.\nP2: foo bar")
}

func TestTemplateExistsCaseIssue13684(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/home.html --
P1: {{ templates.Exists "_partials/MyPartial.html" }}|P1: {{ templates.Exists "_partials/mypartial.html" }}|
-- layouts/_partials/MyPartial.html --
MyPartial.

`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "P1: true|P1: true|")
}
