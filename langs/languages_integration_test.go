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
