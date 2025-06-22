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

package hugolib

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFrontMatterParamsInItsOwnSection(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
-- content/_index.md --
+++
title = "Home"
[[cascade]]
background = 'yosemite.jpg'
[cascade.params]
a = "home-a"
b = "home-b"
[cascade._target]
kind = 'page'
+++
-- content/p1.md --
---
title: "P1"
summary: "frontmatter.summary"
params:
   a: "p1-a"
   summary: "params.summary"
---	
-- layouts/_default/single.html --
Params: {{ range $k, $v := .Params }}{{ $k }}: {{ $v }}|{{ end }}$
Summary: {{ .Summary }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Params: a: p1-a|b: home-b|background: yosemite.jpg|draft: false|iscjklanguage: false|summary: params.summary|title: P1|$",
		"Summary: frontmatter.summary|",
	)
}

func TestFrontMatterParamsPath(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term"]

-- content/p1.md --
---
title: "P1"
date: 2019-08-07
path: "/a/b/c"
slug: "s1"
---
-- content/mysection/_index.md --
---
title: "My Section"
date: 2022-08-07
path: "/a/b"
---
-- layouts/index.html --
RegularPages: {{ range site.RegularPages }}{{ .Path }}|{{ .RelPermalink }}|{{ .Title }}|{{ .Date.Format "2006-02-01" }}| Slug: {{ .Params.slug }}|{{ end }}$
Sections: {{ range site.Sections }}{{ .Path }}|{{ .RelPermalink }}|{{ .Title }}|{{ .Date.Format "2006-02-01" }}| Slug: {{ .Params.slug }}|{{ end }}$
{{ $ab := site.GetPage "a/b" }}
a/b pages: {{ range $ab.RegularPages }}{{ .Path }}|{{ .RelPermalink }}|{{ end }}$
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html",
		"RegularPages: /a/b/c|/a/b/s1/|P1|2019-07-08| Slug: s1|$",
		"Sections: /a|/a/|As",
		"a/b pages: /a/b/c|/a/b/s1/|$",
	)
}

func TestFrontMatterParamsLangNoCascade(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/_index.md --
+++
[[cascade]]
background = 'yosemite.jpg'
lang = 'nn'
+++

`

	b, err := TestE(t, files)
	b.Assert(err, qt.IsNotNil)
}

// Issue 11970.
func TestFrontMatterBuildIsHugoKeyword(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
-- content/p1.md --
---
title: "P1"
build: "foo"
---
-- layouts/_default/single.html --
Params: {{ range $k, $v := .Params }}{{ $k }}: {{ $v }}|{{ end }}$
`
	b, err := TestE(t, files)

	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "We renamed the _build keyword")
}
