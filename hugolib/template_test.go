// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"testing"

	qt "github.com/frankban/quicktest"
)

// https://github.com/gohugoio/hugo/issues/4895
func TestTemplateBOM(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()
	bom := "\ufeff"

	b.WithTemplatesAdded(
		"_default/baseof.html", bom+`
		Base: {{ block "main" . }}base main{{ end }}`,
		"_default/single.html", bom+`{{ define "main" }}Hi!?{{ end }}`)

	b.WithContent("page.md", `---
title: "Page"
---

Page Content
`)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/page/index.html", "Base: Hi!?")
}

func TestTemplateManyBaseTemplates(t *testing.T) {
	t.Parallel()
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	numPages := 100 // To get some parallelism

	pageTemplate := `---
title: "Page %d"
layout: "layout%d"
---

Content.
`

	singleTemplate := `
{{ define "main" }}%d{{ end }}
`
	baseTemplate := `
Base %d: {{ block "main" . }}FOO{{ end }}
`

	for i := range numPages {
		id := i + 1
		b.WithContent(fmt.Sprintf("page%d.md", id), fmt.Sprintf(pageTemplate, id, id))
		b.WithTemplates(fmt.Sprintf("_default/layout%d.html", id), fmt.Sprintf(singleTemplate, id))
		b.WithTemplates(fmt.Sprintf("_default/layout%d-baseof.html", id), fmt.Sprintf(baseTemplate, id))
	}

	b.Build(BuildCfg{})
	for i := range numPages {
		id := i + 1
		b.AssertFileContent(fmt.Sprintf("public/page%d/index.html", id), fmt.Sprintf(`Base %d: %d`, id, id))
	}
}

// https://github.com/gohugoio/hugo/issues/6790
func TestTemplateNoBasePlease(t *testing.T) {
	t.Parallel()
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	b.WithTemplates("_default/list.html", `
{{ define "main" }}
  Bonjour
{{ end }}

{{ printf "list" }}


	`)

	b.WithTemplates(
		"_default/single.html", `
{{ printf "single" }}
{{ define "main" }}
  Bonjour
{{ end }}


`)

	b.WithContent("blog/p1.md", `---
title: The Page
---
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/blog/p1/index.html", `single`)
	b.AssertFileContent("public/blog/index.html", `list`)
}

// https://github.com/gohugoio/hugo/issues/6816
func TestTemplateBaseWithComment(t *testing.T) {
	t.Parallel()
	b := newTestSitesBuilder(t).WithSimpleConfigFile()
	b.WithTemplatesAdded(
		"baseof.html", `Base: {{ block "main" . }}{{ end }}`,
		"index.html", `
	{{/*  A comment */}}
	{{ define "main" }}
	  Bonjour
	{{ end }}


	`)

	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html", `Base:
Bonjour`)
}

func TestTemplateLookupSite(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		b := newTestSitesBuilder(t).WithSimpleConfigFile()
		b.WithTemplates(
			"_default/single.html", `Single: {{ .Title }}`,
			"_default/list.html", `List: {{ .Title }}`,
		)

		createContent := func(title string) string {
			return fmt.Sprintf(`---
title: %s
---`, title)
		}

		b.WithContent(
			"_index.md", createContent("Home Sweet Home"),
			"p1.md", createContent("P1"))

		b.CreateSites().Build(BuildCfg{})
		b.AssertFileContent("public/index.html", `List: Home Sweet Home`)
		b.AssertFileContent("public/p1/index.html", `Single: P1`)
	})

	{
	}
}

func TestTemplateLookupSitBaseOf(t *testing.T) {
	t.Parallel()
	b := newTestSitesBuilder(t).WithDefaultMultiSiteConfig()

	b.WithTemplatesAdded(
		"index.html", `{{ define "main" }}Main Home En{{ end }}`,
		"index.fr.html", `{{ define "main" }}Main Home Fr{{ end }}`,
		"baseof.html", `Baseof en: {{ block "main" . }}main block{{ end }}`,
		"baseof.fr.html", `Baseof fr: {{ block "main" . }}main block{{ end }}`,
		"mysection/baseof.html", `Baseof mysection: {{ block "main" .  }}mysection block{{ end }}`,
		"_default/single.html", `{{ define "main" }}Main Default Single{{ end }}`,
		"_default/list.html", `{{ define "main" }}Main Default List{{ end }}`,
	)

	b.WithContent("mysection/p1.md", `---
title: My Page
---

`)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/en/index.html", `Baseof en: Main Home En`)
	b.AssertFileContent("public/fr/index.html", `Baseof fr: Main Home Fr`)
	b.AssertFileContent("public/en/mysection/index.html", `Baseof mysection: Main Default List`)
	b.AssertFileContent("public/en/mysection/p1/index.html", `Baseof mysection: Main Default Single`)
}

func TestTemplateFuncs(t *testing.T) {
	b := newTestSitesBuilder(t).WithDefaultMultiSiteConfig()

	homeTpl := `Site: {{ site.Language.Lang }} / {{ .Site.Language.Lang }} / {{ site.BaseURL }}
Sites: {{ site.Sites.Default.Home.Language.Lang }}
Hugo: {{ hugo.Generator }}
`

	b.WithTemplatesAdded(
		"index.html", homeTpl,
		"index.fr.html", homeTpl,
	)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/en/index.html",
		"Site: en / en / http://example.com/blog",
		"Sites: en",
		"Hugo: <meta name=\"generator\" content=\"Hugo")
	b.AssertFileContent("public/fr/index.html",
		"Site: fr / fr / http://example.com/blog",
		"Sites: en",
		"Hugo: <meta name=\"generator\" content=\"Hugo",
	)
}

func TestPartialWithReturn(t *testing.T) {
	c := qt.New(t)

	newBuilder := func(t testing.TB) *sitesBuilder {
		b := newTestSitesBuilder(t).WithSimpleConfigFile()
		b.WithTemplatesAdded(
			"partials/add42.tpl", `
		{{ $v := add . 42 }}
		{{ return $v }}
		`,
			"partials/dollarContext.tpl", `
{{ $v := add $ 42 }}
{{ return $v }}
`,
			"partials/dict.tpl", `
{{ $v := add $.adder 42 }}
{{ return $v }}
`,
			"partials/complex.tpl", `
{{ return add . 42 }}
`, "partials/hello.tpl", `
		{{ $v := printf "hello %s" . }}
		{{ return $v }}
		`,
		)

		return b
	}

	c.Run("Return", func(c *qt.C) {
		for range 2 {
			b := newBuilder(c)

			b.WithTemplatesAdded(
				"index.html", `
Test Partials With Return Values:

add42: 50: {{ partial "add42.tpl" 8 }}
hello world: {{ partial "hello.tpl" "world" }}
dollarContext: 60: {{ partial "dollarContext.tpl" 18 }}
adder: 70: {{ partial "dict.tpl" (dict "adder" 28) }}
complex: 80: {{ partial "complex.tpl" 38 }}
`,
			)

			b.CreateSites().Build(BuildCfg{})

			b.AssertFileContent("public/index.html", `
add42: 50: 50
hello world: hello world
dollarContext: 60: 60
adder: 70: 70
complex: 80: 80
`,
			)
		}
	})
}

// Issue 7528
func TestPartialWithZeroedArgs(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithTemplatesAdded("index.html",
		`
X{{ partial "retval" dict }}X
X{{ partial "retval" slice }}X
X{{ partial "retval" "" }}X
X{{ partial "retval" false }}X
X{{ partial "retval" 0 }}X
{{ define "partials/retval" }}
  {{ return 123 }}
{{ end }}`)

	b.WithContentAdded("p.md", ``)
	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html",
		`
X123X
X123X
X123X
X123X
X123X
`)
}

func TestPartialCached(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithTemplatesAdded(
		"index.html", `
{{ $key1 := (dict "a" "av" ) }}
{{ $key2 := (dict "a" "av2" ) }}
Partial cached1: {{ partialCached "p1" "input1" $key1 }}
Partial cached2: {{ partialCached "p1" "input2" $key1 }}
Partial cached3: {{ partialCached "p1" "input3" $key2 }}
`,

		"partials/p1.html", `partial: {{ . }}`,
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
 Partial cached1: partial: input1
 Partial cached2: partial: input1
 Partial cached3: partial: input3
`)
}

// https://github.com/gohugoio/hugo/issues/6615
func TestTemplateTruth(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithTemplatesAdded("index.html", `
{{ $p := index site.RegularPages 0 }}
{{ $zero := $p.ExpiryDate }}
{{ $notZero := time.Now }}

if: Zero: {{ if $zero }}FAIL{{ else }}OK{{ end }}
if: Not Zero: {{ if $notZero }}OK{{ else }}Fail{{ end }}
not: Zero: {{ if not $zero }}OK{{ else }}FAIL{{ end }}
not: Not Zero: {{ if not $notZero }}FAIL{{ else }}OK{{ end }}

with: Zero {{ with $zero }}FAIL{{ else }}OK{{ end }}

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
if: Zero: OK
if: Not Zero: OK
not: Zero: OK
not: Not Zero: OK
with: Zero OK
`)
}

func TestTemplateGoIssues(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithTemplatesAdded(
		"index.html", `
{{ $title := "a & b" }}
<script type="application/ld+json">{"@type":"WebPage","headline":"{{$title}}"}</script>

{{/* Action/commands newlines, from Go 1.16, see https://github.com/golang/go/issues/29770 */}}
{{ $norway := dict
	"country" "Norway"
	"population" "5 millions"
	"language" "Norwegian"
	"language_code" "nb"
	"weather" "freezing cold"
	"capitol" "Oslo"
	"largest_city" "Oslo"
	"currency"  "Norwegian krone"
	"dialing_code" "+47"
}}

Population in Norway is {{
	  $norway.population
	| lower
	| upper
}}

`,
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
<script type="application/ld+json">{"@type":"WebPage","headline":"a \u0026 b"}</script>
Population in Norway is 5 MILLIONS

`)
}

func TestPartialInline(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithContent("p1.md", "")

	b.WithTemplates(
		"index.html", `

{{ $p1 := partial "p1" . }}
{{ $p2 := partial "p2" . }}

P1: {{ $p1 }}
P2: {{ $p2 }}

{{ define "partials/p1" }}Inline: p1{{ end }}

{{ define "partials/p2" }}
{{ $value := 32 }}
{{ return $value }}
{{ end }}


`,
	)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		`
P1: Inline: p1
P2: 32`,
	)
}

func TestPartialInlineBase(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithContent("p1.md", "")

	b.WithTemplates(
		"baseof.html", `{{ $p3 := partial "p3" . }}P3: {{ $p3 }}
{{ block "main" . }}{{ end }}{{ define "partials/p3" }}Inline: p3{{ end }}`,
		"index.html", `
{{ define "main" }}

{{ $p1 := partial "p1" . }}
{{ $p2 := partial "p2" . }}

P1: {{ $p1 }}
P2: {{ $p2 }}

{{ end }}


{{ define "partials/p1" }}Inline: p1{{ end }}

{{ define "partials/p2" }}
{{ $value := 32 }}
{{ return $value }}
{{ end }}


`,
	)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		`
P1: Inline: p1
P2: 32
P3: Inline: p3
`,
	)
}

// https://github.com/gohugoio/hugo/issues/7478
func TestBaseWithAndWithoutDefine(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithContent("p1.md", "---\ntitle: P\n---\nContent")

	b.WithTemplates(
		"_default/baseof.html", `
::Header Start:{{ block "header" . }}{{ end }}:Header End:
::{{ block "main" . }}Main{{ end }}::
`, "index.html", `
{{ define "header" }}
Home Header
{{ end }}
{{ define "main" }}
This is home main
{{ end }}
`,

		"_default/single.html", `
{{ define "main" }}
This is single main
{{ end }}
`,
	)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
Home Header
This is home main
`,
	)

	b.AssertFileContent("public/p1/index.html", `
 ::Header Start::Header End:
This is single main
`,
	)
}

// Issue 9393.
func TestApplyWithNamespace(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithTemplates(
		"index.html", `
{{ $b := slice " a " "     b "   "       c" }}
{{ $a := apply $b "strings.Trim" "." " " }}
a: {{ $a }}
`,
	).WithContent("p1.md", "")

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `a: [a b c]`)
}

// Legacy behavior for internal templates.
func TestOverrideInternalTemplate(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/index.html --
{{ template "_internal/google_analytics_async.html" . }}
-- layouts/_internal/google_analytics_async.html --
Overridden.
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Overridden.")
}
