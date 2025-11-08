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

package hugolib

import (
	"fmt"
	"strings"
	"testing"
)

// https://github.com/gohugoio/hugo/issues/4895
func TestTemplateBOM(t *testing.T) {
	t.Parallel()

	bom := "\ufeff"

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/_default/baseof.html --
` + bom + `
		Base: {{ block "main" . }}base main{{ end }}
-- layouts/_default/single.html --
` + bom + `{{ define "main" }}Hi!?{{ end }}
-- content/page.md --
---
title: "Page"
---

Page Content
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/page/index.html", "Base: Hi!?")
}

func TestTemplateManyBaseTemplates(t *testing.T) {
	t.Parallel()

	numPages := 100

	var b strings.Builder
	b.WriteString("-- hugo.toml --\n")
	b.WriteString("baseURL = \"http://example.com/\"\n")

	for i := range numPages {
		id := i + 1
		b.WriteString(fmt.Sprintf("-- content/page%d.md --\n", id))
		b.WriteString(fmt.Sprintf(`---
title: "Page %d"
layout: "layout%d"
---

Content.
`, id, id))

		b.WriteString(fmt.Sprintf("-- layouts/_default/layout%d.html --\n", id))
		b.WriteString(fmt.Sprintf(`
{{ define "main" }}%d{{ end }}
`, id))

		b.WriteString(fmt.Sprintf("-- layouts/_default/layout%d-baseof.html --\n", id))
		b.WriteString(fmt.Sprintf(`
Base %d: {{ block "main" . }}FOO{{ end }}
`, id))
	}

	files := b.String()

	builder := Test(t, files)

	for i := range numPages {
		id := i + 1
		builder.AssertFileContent(fmt.Sprintf("public/page%d/index.html", id), fmt.Sprintf(`Base %d: %d`, id, id))
	}
}

// https://github.com/gohugoio/hugo/issues/6790
func TestTemplateNoBasePlease(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/_default/list.html --
{{ define "main" }}
  Bonjour
{{ end }}

{{ printf "list" }}
-- layouts/_default/single.html --
{{ printf "single" }}
{{ define "main" }}
  Bonjour
{{ end }}
-- content/blog/p1.md --
---
title: The Page
---
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/blog/p1/index.html", `single`)
	b.AssertFileContent("public/blog/index.html", `list`)
}

// https://github.com/gohugoio/hugo/issues/6816
func TestTemplateBaseWithComment(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/baseof.html --
Base: {{ block "main" . }}{{ end }}
-- layouts/index.html --
{{/*  A comment */}}
{{ define "main" }}
  Bonjour
{{ end }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `Base:
Bonjour`)
}

func TestTemplateLookupSite(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/_default/single.html --
Single: {{ .Title }}
-- layouts/_default/list.html --
List: {{ .Title }}
-- content/_index.md --
---
title: Home Sweet Home
---
-- content/p1.md --
---
title: P1
---
`
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
				BuildCfg:    BuildCfg{},
			},
		).Build()

		b.AssertFileContent("public/index.html", `List: Home Sweet Home`)
		b.AssertFileContent("public/p1/index.html", `Single: P1`)
	})
}

func TestTemplateLookupSitBaseOf(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/blog"
disablePathToLower = true
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true

[languages]
[languages.en]
weight = 10
[languages.fr]
weight = 20

-- layouts/index.html --
{{ define "main" }}Main Home En{{ end }}
-- layouts/index.fr.html --
{{ define "main" }}Main Home Fr{{ end }}
-- layouts/baseof.html --
Baseof en: {{ block "main" . }}main block{{ end }}
-- layouts/baseof.fr.html --
Baseof fr: {{ block "main" . }}main block{{ end }}
-- layouts/mysection/baseof.html --
Baseof mysection: {{ block "main" .  }}mysection block{{ end }}
-- layouts/_default/single.html --
{{ define "main" }}Main Default Single{{ end }}
-- layouts/_default/list.html --
{{ define "main" }}Main Default List{{ end }}
-- content/mysection/p1.md --
---
title: My Page
---
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/en/index.html", `Baseof en: Main Home En`)
	b.AssertFileContent("public/fr/index.html", `Baseof fr: Main Home Fr`)
	b.AssertFileContent("public/en/mysection/index.html", `Baseof mysection: Main Default List`)
	b.AssertFileContent("public/en/mysection/p1/index.html", `Baseof mysection: Main Default Single`)
}

func TestTemplateFuncs(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/blog"
disablePathToLower = true
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true

[languages]
[languages.en]
weight = 10
[languages.fr]
weight = 20

-- layouts/index.html --
Site: {{ site.Language.Lang }} / {{ .Site.Language.Lang }} / {{ site.BaseURL }}
Sites: {{ site.Sites.Default.Home.Language.Lang }}
Hugo: {{ hugo.Generator }}
-- layouts/index.fr.html --
Site: {{ site.Language.Lang }} / {{ .Site.Language.Lang }} / {{ site.BaseURL }}
Sites: {{ site.Sites.Default.Home.Language.Lang }}
Hugo: {{ hugo.Generator }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

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
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/partials/add42.tpl --
{{ $v := add . 42 }}
{{ return $v }}
-- layouts/partials/dollarContext.tpl --
{{ $v := add $ 42 }}
{{ return $v }}
-- layouts/partials/dict.tpl --
{{ $v := add $.adder 42 }}
{{ return $v }}
-- layouts/partials/complex.tpl --
{{ return add . 42 }}
-- layouts/partials/hello.tpl --
{{ $v := printf "hello %s" . }}
{{ return $v }}
-- layouts/index.html --
Test Partials With Return Values:

add42: 50: {{ partial "add42.tpl" 8 }}
hello world: {{ partial "hello.tpl" "world" }}
dollarContext: 60: {{ partial "dollarContext.tpl" 18 }}
adder: 70: {{ partial "dict.tpl" (dict "adder" 28) }}
complex: 80: {{ partial "complex.tpl" 38 }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `
add42: 50: 50
hello world: hello world
dollarContext: 60: 60
adder: 70: 70
complex: 80: 80
`)
}

// Issue 7528
func TestPartialWithZeroedArgs(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/index.html --
X{{ partial "retval" dict }}X
X{{ partial "retval" slice }}X
X{{ partial "retval" "" }}X
X{{ partial "retval" false }}X
X{{ partial "retval" 0 }}X
{{ define "partials/retval" }}
  {{ return 123 }}
{{ end }}
-- content/p.md --
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

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
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/index.html --
{{ $key1 := (dict "a" "av" ) }}
{{ $key2 := (dict "a" "av2" ) }}
Partial cached1: {{ partialCached "p1" "input1" $key1 }}
Partial cached2: {{ partialCached "p1" "input2" $key1 }}
Partial cached3: {{ partialCached "p1" "input3" $key2 }}
-- layouts/partials/p1.html --
partial: {{ . }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `
 Partial cached1: partial: input1
 Partial cached2: partial: input1
 Partial cached3: partial: input3
`)
}

// https://github.com/gohugoio/hugo/issues/6615
func TestTemplateTruth(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/index.html --
{{ $p := index site.RegularPages 0 }}
{{ $zero := $p.ExpiryDate }}
{{ $notZero := time.Now }}

if: Zero: {{ if $zero }}FAIL{{ else }}OK{{ end }}
if: Not Zero: {{ if $notZero }}OK{{ else }}Fail{{ end }}
not: Zero: {{ if not $zero }}OK{{ else }}FAIL{{ end }}
not: Not Zero: {{ if not $notZero }}FAIL{{ else }}OK{{ end }}

with: Zero {{ with $zero }}FAIL{{ else }}OK{{ end }}
-- content/p1.md --
---
title: p1
---
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `
if: Zero: OK
if: Not Zero: OK
not: Zero: OK
not: Not Zero: OK
with: Zero OK
`)
}

func TestTemplateGoIssues(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- layouts/index.html --
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
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `
<script type="application/ld+json">{"@type":"WebPage","headline":"a \u0026 b"}</script>
Population in Norway is 5 MILLIONS

`)
}

func TestPartialInline(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/p1.md --
-- layouts/index.html --
{{ $p1 := partial "p1" . }}
{{ $p2 := partial "p2" . }}

P1: {{ $p1 }}
P2: {{ $p2 }}

{{ define "partials/p1" }}Inline: p1{{ end }}

{{ define "partials/p2" }}
{{ $value := 32 }}
{{ return $value }}
{{ end }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html",
		`
P1: Inline: p1
P2: 32`,
	)
}

func TestPartialInlineBase(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/p1.md --
-- layouts/baseof.html --
{{ $p3 := partial "p3" . }}P3: {{ $p3 }}
{{ block "main" . }}{{ end }}{{ define "partials/p3" }}Inline: p3{{ end }}
-- layouts/index.html --
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
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

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
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/p1.md --
---
title: P
---
Content
-- layouts/_default/baseof.html --
::Header Start:{{ block "header" . }}{{ end }}:Header End:
::{{ block "main" . }}Main{{ end }}::
-- layouts/index.html --
{{ define "header" }}
Home Header
{{ end }}
{{ define "main" }}
This is home main
{{ end }}
-- layouts/_default/single.html --
{{ define "main" }}
This is single main
{{ end }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

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
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/p1.md --
-- layouts/index.html --
{{ $b := slice " a " "     b "   "       c" }}
{{ $a := apply $b "strings.Trim" "." " " }}
a: {{ $a }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

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
