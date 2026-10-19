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
	"fmt"
	"path/filepath"
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
[cascade.target]
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
-- layouts/single.html --
Params: {{ range $k, $v := .Params }}{{ $k }}: {{ $v }}|{{ end }}$
Summary: {{ .Summary }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Params: a: p1-a|b: home-b|background: yosemite.jpg|draft: false|iscjklanguage: false|summary: params.summary|title: P1|$",
		"Summary: frontmatter.summary|",
	)
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
-- layouts/single.html --
Params: {{ range $k, $v := .Params }}{{ $k }}: {{ $v }}|{{ end }}$
`
	b, err := TestE(t, files)

	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "We renamed the _build keyword")
}

// Issue 15411.
func TestParamsReservedFieldsRedirect(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
-- content/_index.md --
---
title: Home
cascade:
  - params:
      title: cascade-title
    target:
      path: "{/p3}"
---
-- content/p1.md --
---
title: p1
summary: A **bold** word.
magic: 41
params:
  magic: 42
---
-- content/p2.md --
---
title: p2
params:
  title: params-title
---
-- content/p3.md --
---
---
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "p4" "title" "p4" "params" (dict "magic" 42)) }}
{{ .AddPage (dict "path" "p5" "title" "p5" "params" (dict "title" "params-title-p5")) }}
-- layouts/page.html --
{{ partial "params.html" . }}
-- layouts/home.html --
Where: {{ range where site.RegularPages "Params.title" "p1" }}{{ .Path }}|{{ end }}
WhereMagic: {{ range where site.RegularPages "Params.magic" 42 }}{{ .Path }}|{{ end }}
ByParam: {{ range site.RegularPages.ByParam "title" }}{{ .Title }}|{{ end }}
GroupByParam: {{ range site.RegularPages.GroupByParam "title" }}{{ .Key }}|{{ end }}
-- layouts/_partials/params.html --
Title: {{ .Params.title }}|{{ $.Params.Title }}|{{ (.).Params.TITLE }}|
Summary: {{ .Params.summary }}|{{ .RawSummary }}|
Draft: {{ .Params.draft }}|
Magic: {{ .Params.magic }}|
`

	b := Test(t, files, TestOptInfo())

	b.AssertFileContent("public/p1/index.html",
		"Title: p1|p1|p1|",
		"Summary: A **bold** word.|A **bold** word.|",
		"Draft: false|",
		"Magic: 42|",
	)
	b.AssertFileContent("public/p2/index.html", "Title: params-title|params-title|params-title|")
	b.AssertFileContent("public/p3/index.html", "Title: cascade-title|cascade-title|cascade-title|")
	b.AssertFileContent("public/p4/index.html", "Title: |||", "Magic: 42|")
	b.AssertFileContent("public/p5/index.html", "Title: params-title-p5|params-title-p5|params-title-p5|")

	b.AssertFileContent("public/index.html",
		"Where: /p1|",
		"WhereMagic: /p1|/p4|",
		"ByParam: cascade-title|p1|p2|p5|p4|",
		"GroupByParam: cascade-title|p1|params-title|params-title-p5|",
	)

	b.AssertLogContains(
		fmt.Sprintf("Using .Params.title in template %q to get the title front matter field was deprecated in Hugo 0.169.0", filepath.FromSlash("/layouts/_partials/params.html")),
		"Use .Title instead.",
		"Use .RawSummary instead.",
		"Use .Draft instead.",
		`Using where with "Params.title" to get the title front matter field was deprecated in Hugo 0.169.0 and will be removed in a future release. Use where with "Title" instead.`,
		`Using ByParam "title" to get the title front matter field was deprecated in Hugo 0.169.0 and will be removed in a future release. Use sort with "Title" instead.`,
		`Using GroupByParam "title" to get the title front matter field was deprecated in Hugo 0.169.0 and will be removed in a future release. Use GroupBy "Title" instead.`,
		"! magic",
	)
}

// See issue 15411.
func TestPageParamReservedFrontMatterField(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[params]
description = "site"
-- content/p1.md --
---
title: p1
description: fm
---
-- content/p2.md --
---
title: p2
---
-- layouts/page.html --
{{ .Param "description" }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "fm|")
	b.AssertFileContent("public/p2/index.html", "site|")
}

// See issue 15411.
func TestPageParamReservedFrontMatterFieldDeprecation(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/p1.md --
---
title: p1
description: fm
---
-- layouts/page.html --
{{ .Param "description" }}|
`
	b := Test(t, files, TestOptInfo())

	b.AssertFileContent("public/p1/index.html", "fm|")
	b.AssertLogContains("to get the description front matter field was deprecated")
}

// See issue 15411.
func TestParamsReservedDateFieldsRedirect(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
timeZone = "America/Los_Angeles"
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
[frontmatter]
date = [":filename", "date"]
-- content/p1.md --
---
title: p1
date: 2001-12-15T02:59:43.1
lastmod: 2002-01-01
publishDate: 2001-12-01
expiryDate: 2099-01-01
---
-- content/2020-01-01-p2.md --
---
title: p2
date: 2001-12-15
---
-- content/p3.md --
---
title: p3
date: 2001-12-15
params:
  date: 1999-01-01
---
-- layouts/page.html --
Date: {{ .Params.date }}|{{ .Param "date" }}|
Lastmod: {{ with .Params.lastmod }}{{ .Format "2006-01-02" }}{{ end }}|
PublishDate: {{ with .Params.publishdate }}{{ .Format "2006-01-02" }}{{ end }}|
ExpiryDate: {{ with .Params.expirydate }}{{ .Format "2006-01-02" }}{{ end }}|
-- layouts/home.html --
GroupByParamDate: {{ range site.RegularPages.GroupByParamDate "date" "2006" }}{{ .Key }}|{{ end }}
`

	b := Test(t, files, TestOptInfo())

	b.AssertFileContent("public/p1/index.html",
		"Date: 2001-12-15 02:59:43.1 -0800 PST|2001-12-15 02:59:43.1 -0800 PST|",
		"Lastmod: 2002-01-01|",
		"PublishDate: 2001-12-01|",
		"ExpiryDate: 2099-01-01|",
	)
	// The date is taken from the filename.
	b.AssertFileContent("public/p2/index.html", "Date: 2020-01-01 00:00:00 -0800 PST|2020-01-01 00:00:00 -0800 PST|")
	// params.date wins, and is not parsed as a date.
	b.AssertFileContent("public/p3/index.html", "Date: 1999-01-01|1999-01-01|")

	b.AssertLogContains(
		"Use .Date instead.",
		"Use .Lastmod instead.",
		"Use .PublishDate instead.",
		"Use .ExpiryDate instead.",
		`Using GroupByParamDate "date" to get the date front matter field was deprecated in Hugo 0.169.0 and will be removed in a future release. Use GroupByDate instead.`,
	)
}
