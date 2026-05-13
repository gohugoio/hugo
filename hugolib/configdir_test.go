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

import "testing"

func TestConfigDir(t *testing.T) {
	t.Parallel()

	files := `
-- config/_default/params.toml --
a = "acp1"
d = "dcp1"
-- config/_default/config.toml --
[params]
a = "ac1"
b = "bc1"

-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]
ignoreErrors = ["error-missing-instagram-accesstoken"]
[params]
a = "a1"
b = "b1"
c = "c1"
-- layouts/home.html --
Params: {{ site.Params}}
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", `
Params: map[a:acp1 b:bc1 c:c1 d:dcp1]


`)
}

// TOML/YAML can't represent a top-level array, so a basename-matched wrapper
// key in the config file unwraps to the slice-typed root key.
func TestConfigDirCascadeSliceIssue12899(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "home", "section"]
-- config/_default/cascade.yaml --
cascade:
  - target:
      path: /books/**
    params:
      color: red
  - target:
      path: /films/**
    params:
      color: blue
-- content/books/b1.md --
---
title: B1
---
-- content/films/f1.md --
---
title: F1
---
-- layouts/page.html --
{{ .Title }}|color:{{ .Params.color }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/books/b1/index.html", "B1|color:red|")
	b.AssertFileContent("public/films/f1/index.html", "F1|color:blue|")
}

func TestConfigDirPermalinksSliceIssue12899(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "home", "section"]
-- config/_default/permalinks.yaml --
permalinks:
  - target:
      path: /books/**
    pattern: /shelf/:slug/
-- content/books/b1.md --
---
title: B1
slug: novel
---
-- layouts/page.html --
{{ .Title }}|{{ .RelPermalink }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/shelf/novel/index.html", "B1|/shelf/novel/|")
}

func TestConfigDirCascadeEnvironmentOverrideIssue12899(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "home", "section"]
-- config/_default/cascade.yaml --
cascade:
  - target:
      path: /**
    params:
      color: default
-- config/production/cascade.yaml --
cascade:
  - target:
      path: /**
    params:
      color: production
-- content/p1.md --
---
title: P1
---
-- layouts/page.html --
{{ .Title }}|color:{{ .Params.color }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "P1|color:production|")
}

// The basename-match unwrap is type-agnostic — also unwraps maps.
func TestConfigDirRepeatRootMapIssue14882(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "section"]
-- config/_default/params.yaml --
params:
  a: aval
  b: bval
-- layouts/home.html --
a:{{ site.Params.a }}|b:{{ site.Params.b }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "a:aval|b:bval|")
}

// The unwrap fires only when the basename-matched key is the sole top-level
// key in the file. A mixed map must be left alone.
func TestConfigDirUnwrapOnlySoleKeyIssue12899(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "home", "section"]
-- config/_default/params.yaml --
params:
  nested: yes
other: top
-- content/p1.md --
---
title: P1
---
-- layouts/page.html --
{{ .Title }}|params.params.nested:{{ site.Params.params.nested }}|params.other:{{ site.Params.other }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "P1|params.params.nested:yes|params.other:top|")
}
