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

package langs_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestLanguagesContentSimple(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 2
[languages.nn]
weight = 1
-- content/_index.md --
---
title: "Home"
---

Welcome to the home page.
-- content/_index.nn.md --
---
title: "Heim"
---
Welkomen heim!
-- layouts/all.html --
title: {{ .Title }}|
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/index.html", `title: Home|`)
	b.AssertFileContent("public/nn/index.html", `title: Heim|`)
}

func TestLanguagesContentSharedResource(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "sitemap", "404"]
baseURL = "https://example.org/"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 2
[languages.nn]
weight = 1
-- content/mytext.txt --
This is a shared text file.
-- content/_index.md --
---
title: "Home"
---

Welcome to the home page.
-- content/_index.nn.md --
---
title: "Heim"
---
Welkomen heim!
-- layouts/home.html --
{{ $text := .Resources.Get "mytext.txt" }}
title: {{ .Title }}|text: {{ with $text }}{{ .Content | safeHTML }}{{ end }}|{{  .Resources | len}}

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/nn/index.html", `title: Heim|text: This is a shared text file.|1`)
	b.AssertFileContent("public/en/index.html", `title: Home|text: This is a shared text file.|1`)
}

// Issue 14031
func TestDraftTermIssue14031(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap']
capitalizeListTitles = false
pluralizeListTitles = false
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages.en]
weight = 1
[languages.fr]
weight = 2
[taxonomies]
tag = 'tags'
-- content/p1.en.md --
---
title: P1 (en)
tags: [a,b]
---
-- content/p1.fr.md --
---
title: P1 (fr)
tags: [a,b]
---
-- content/tags/a/_index.en.md --
---
title: a (en)
---
-- content/tags/a/_index.fr.md --
---
title: a (fr)
---
-- content/tags/b/_index.en.md --
---
title: b (en)
---
-- content/tags/b/_index.fr.md --
---
title: b (fr)
draft: false
---
-- layouts/list.html --
{{- .Title -}}|
{{- range .Pages -}}
TITLE: {{ .Title }} RELPERMALINK: {{ .RelPermalink }}|
{{- end -}}
-- layouts/page.html --
{{- .Title -}}|
{{- range .GetTerms "tags" -}}
TITLE: {{ .Title }} RELPERMALINK: {{ .RelPermalink }}|
{{- end -}}
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/en/tags/a/index.html", true)
	b.AssertFileExists("public/fr/tags/a/index.html", true)
	b.AssertFileExists("public/en/tags/b/index.html", true)
	b.AssertFileExists("public/fr/tags/b/index.html", true)

	b.AssertFileContentEquals("public/en/p1/index.html",
		"P1 (en)|TITLE: a (en) RELPERMALINK: /en/tags/a/|TITLE: b (en) RELPERMALINK: /en/tags/b/|",
	)
	b.AssertFileContentEquals("public/fr/p1/index.html",
		"P1 (fr)|TITLE: a (fr) RELPERMALINK: /fr/tags/a/|TITLE: b (fr) RELPERMALINK: /fr/tags/b/|",
	)
	b.AssertFileContentEquals("public/en/tags/index.html",
		"tags|TITLE: a (en) RELPERMALINK: /en/tags/a/|TITLE: b (en) RELPERMALINK: /en/tags/b/|",
	)
	b.AssertFileContentEquals("public/fr/tags/index.html",
		"tags|TITLE: a (fr) RELPERMALINK: /fr/tags/a/|TITLE: b (fr) RELPERMALINK: /fr/tags/b/|",
	)
	b.AssertFileContentEquals("public/en/tags/a/index.html",
		"a (en)|TITLE: P1 (en) RELPERMALINK: /en/p1/|",
	)
	b.AssertFileContentEquals("public/fr/tags/a/index.html",
		"a (fr)|TITLE: P1 (fr) RELPERMALINK: /fr/p1/|",
	)
	b.AssertFileContentEquals("public/en/tags/b/index.html",
		"b (en)|TITLE: P1 (en) RELPERMALINK: /en/p1/|",
	)
	b.AssertFileContentEquals("public/fr/tags/b/index.html",
		"b (fr)|TITLE: P1 (fr) RELPERMALINK: /fr/p1/|",
	)

	// Set draft to true on content/tags/b/_index.fr.md.
	files = strings.ReplaceAll(files, "draft: false", "draft: true")

	b = hugolib.Test(t, files)

	b.AssertFileExists("public/en/tags/a/index.html", true)
	b.AssertFileExists("public/fr/tags/a/index.html", true)
	b.AssertFileExists("public/en/tags/b/index.html", true)
	b.AssertFileExists("public/fr/tags/b/index.html", false)

	// The assertion below fails.
	// Got: P1 (en)|TITLE: a (en) RELPERMALINK: /en/tags/a/|
	b.AssertFileContentEquals("public/en/p1/index.html",
		"P1 (en)|TITLE: a (en) RELPERMALINK: /en/tags/a/|TITLE: b (en) RELPERMALINK: /en/tags/b/|",
	)
	b.AssertFileContentEquals("public/fr/p1/index.html",
		"P1 (fr)|TITLE: a (fr) RELPERMALINK: /fr/tags/a/|",
	)
	b.AssertFileContentEquals("public/en/tags/index.html",
		"tags|TITLE: a (en) RELPERMALINK: /en/tags/a/|TITLE: b (en) RELPERMALINK: /en/tags/b/|",
	)
	b.AssertFileContentEquals("public/fr/tags/index.html",
		"tags|TITLE: a (fr) RELPERMALINK: /fr/tags/a/|",
	)
	b.AssertFileContentEquals("public/en/tags/a/index.html",
		"a (en)|TITLE: P1 (en) RELPERMALINK: /en/p1/|",
	)
	b.AssertFileContentEquals("public/fr/tags/a/index.html",
		"a (fr)|TITLE: P1 (fr) RELPERMALINK: /fr/p1/|",
	)
	// The assertion below fails.
	// Got: b (en)|
	b.AssertFileContentEquals("public/en/tags/b/index.html",
		"b (en)|TITLE: P1 (en) RELPERMALINK: /en/p1/|",
	)
}

func TestLocaleAsLocalizationKeyIssue9109(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true

[languages.de]
weight = 1

[languages.en]
locale = 'en-US'
weight = 2

[languages.fr]
locale = 'bogus'
weight = 3

[languages.nn]
weight = 4

[languages.xx]
locale = 'bogus'
weight = 5

[languages.zh-cn]
locale = 'zh-Hans'
weight = 6
-- layouts/home.html --
Time: {{ "2026-04-01T20:11:31+08:00" | time.Format ":date_long" }}|
FormatAccounting: {{ 512.5032 | lang.FormatAccounting 2 "USD" }}|
FormatCurrency: {{ 512.5032 | lang.FormatCurrency 2 "USD" }}|
`

	b := hugolib.Test(t, files)

	// de: no locale => localize with lang.
	b.AssertFileContent("public/de/index.html",
		`Time: 1. April 2026|`,
		"FormatAccounting: 512,50\u00a0$|",
		"FormatCurrency: 512,50\u00a0$|",
	)

	// en: valid locale => localize with locale.
	b.AssertFileContent("public/en/index.html",
		`Time: April 1, 2026|`,
		`FormatAccounting: $512.50|`,
		`FormatCurrency: $512.50|`,
	)

	// fr: invalid locale => localize with lang.
	b.AssertFileContent("public/fr/index.html",
		`Time: 1 avril 2026|`,
		"FormatAccounting: 512,50\u00a0$US|",
		"FormatCurrency: 512,50\u00a0$US|",
	)

	// nn: no locale => localize with lang.
	b.AssertFileContent("public/nn/index.html",
		`Time: 1. april 2026|`,
		"FormatAccounting: 512,50\u00a0USD|",
		"FormatCurrency: 512,50\u00a0USD|",
	)

	// xx: invalid locale + invalid lang => localize with defaultContentLanguage.
	b.AssertFileContent("public/xx/index.html",
		`Time: 1. April 2026|`,
		"FormatAccounting: 512,50\u00a0$|",
		"FormatCurrency: 512,50\u00a0$|",
	)

	// zh-cn: invalid lang, but valid locale => localize with locale.
	b.AssertFileContent("public/zh-cn/index.html",
		`Time: 2026年4月1日|`,
		`FormatAccounting: US$512.50|`,
		`FormatCurrency: US$512.50|`,
	)
}

func TestMultipleLanguageVariants7982(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','sitemap','taxonomy','term']
defaultContentLanguageInSubdir = true
[languages.en]
weight = 1
[languages.de]
weight = 2
[languages.de-de]
weight = 3
-- i18n/en.toml --
file = 'en'
-- i18n/de.toml --
file = 'de'
-- i18n/de-de.toml --
file = 'de-de'
-- layouts/index.html --
language: {{ site.Language.Name }} file: {{ T "file" }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/index.html", "language: en file: en|")
	b.AssertFileContent("public/de/index.html", "language: de file: de|")
	b.AssertFileContent("public/de-de/index.html", "language: de-de file: de-de|")
}

func TestDefaultContentLanguageFallback14243(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
defaultContentLanguage = 'es'
defaultContentLanguageInSubdir = true

[languages.es]
locale = 'es-AR'
weight = 1

[languages.pt]
locale = 'pt-BR'
weight = 2
-- layouts/home.html --
{{ T "foo"}}|
-- i18n/es-ar.toml --
foo = 'foo es-ar'
-- i18n/es.toml --
foo = 'foo es'
-- i18n/pt-br.toml --
foo = 'foo pt-br'
-- i18n/pt.toml --
foo = 'foo pt'
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/pt/index.html", "foo pt-br|")

	files = strings.ReplaceAll(files, "i18n/pt-br.toml", "unused-a.toml")

	b = hugolib.Test(t, files)
	b.AssertFileContent("public/pt/index.html", "foo pt|")

	files = strings.ReplaceAll(files, "i18n/pt.toml", "unused-b.toml")

	b = hugolib.Test(t, files)
	b.AssertFileContent("public/pt/index.html", "foo es-ar|")

	files = strings.ReplaceAll(files, "i18n/es-ar.toml", "unused-c.toml")

	b = hugolib.Test(t, files)
	b.AssertFileContent("public/pt/index.html", "foo es|")

	files = strings.ReplaceAll(files, "i18n/es.toml", "unused-d.toml")

	b = hugolib.Test(t, files)
	b.AssertFileContent("public/pt/index.html", "|")
}
