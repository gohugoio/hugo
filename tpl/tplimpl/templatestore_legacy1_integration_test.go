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

package tplimpl_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Old as in before Hugo v0.146.0.
func TestLayoutsOldSetup(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
title = "Title in English"
weight = 1
[languages.nn]
title = "Tittel på nynorsk"
weight = 2
-- layouts/index.html --
Home.
{{ template "_internal/twitter_cards.html" . }}
-- layouts/_default/single.html --
Single.
-- layouts/_default/single.nn.html --
Single NN.
-- layouts/_default/list.html --
List HTML.
-- layouts/docs/list-baseof.html --
Docs Baseof List HTML.
{{ block "main" . }}Docs Baseof List HTML main block.{{ end }}
-- layouts/docs/list.section.html --
{{ define "main" }}
Docs List HTML.
{{ end }}
-- layouts/_default/list.json --
List JSON.
-- layouts/_default/list.rss.xml --
List RSS.
-- layouts/_default/list.nn.rss.xml --
List NN RSS.
-- layouts/_default/baseof.html --
Base.
-- layouts/partials/mypartial.html --
Partial.
-- layouts/shortcodes/myshortcode.html --
Shortcode.
-- content/docs/p1.md --
---
title: "P1"
---

	`

	b := hugolib.Test(t, files)

	//	b.DebugPrint("", tplimpl.CategoryBaseof)

	b.AssertFileContent("public/en/docs/index.html", "Docs Baseof List HTML.\n\nDocs List HTML.")
}

func TestLayoutsOldSetupBaseofPrefix(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/_default/layout1-baseof.html --
Baseof layout1. {{ block "main" . }}{{ end }}
-- layouts/_default/layout2-baseof.html --
Baseof layout2. {{ block "main" . }}{{ end }}
-- layouts/_default/layout1.html --
{{ define "main" }}Layout1. {{ .Title }}{{ end }}
-- layouts/_default/layout2.html --
{{ define "main" }}Layout2. {{ .Title }}{{ end }}
-- content/p1.md --
---
title: "P1"
layout: "layout1"
---
-- content/p2.md --
---
title: "P2"
layout: "layout2"
---
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Baseof layout1. Layout1. P1")
	b.AssertFileContent("public/p2/index.html", "Baseof layout2. Layout2. P2")
}

func TestLayoutsOldSetupTaxonomyAndTerm(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[taxonomies]
cat = 'cats'
dog = 'dogs'
# Templates for term taxonomy, old setup.
-- layouts/dogs/terms.html --
Dogs Terms. Most specific taxonomy template.
-- layouts/taxonomy/terms.html --
Taxonomy Terms. Down the list.
# Templates for term term, old setup.
-- layouts/dogs/term.html --
Dogs Term. Most specific term template.
-- layouts/term/term.html --
Term Term. Down the list.
-- layouts/dogs/max/list.html --
max: {{ .Title }}
-- layouts/_default/list.html --
Default list.
-- layouts/_default/single.html --
Default single.
-- content/p1.md --
---
title: "P1"
dogs: ["luna", "daisy", "max"]
---

`
	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("! WARN")

	b.AssertFileContent("public/dogs/index.html", "Dogs Terms. Most specific taxonomy template.")
	b.AssertFileContent("public/dogs/luna/index.html", "Dogs Term. Most specific term template.")
	b.AssertFileContent("public/dogs/max/index.html", "max: Max") // layouts/dogs/max/list.html wins over layouts/term/term.html because of distance.
}

func TestLayoutsOldSetupCustomRSS(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "page"]
[outputs]
home = ["rss"]
-- layouts/_default/list.rss.xml --
List RSS.
`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.xml", "List RSS.")
}

func TestLayoutWithLanguagesLegacy(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "sitemap", "rss"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[languages.sv]
weight = 3
-- layouts/_default/single.en.html --
layouts/_default/single.html
-- layouts/_default/single.nn.html --
layouts/_default/single.nn.html
-- layouts/_default/single.sv.html --
layouts/_default/single.sv.html
-- content/p1.md --
---
title: "P1"
---
-- content/p1.nn.md --
---
title: "P1 NN"
---
-- content/p1.sv.md --
---
title: "P1 SV"
---
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/p1/index.html", "layouts/_default/single.html")
	b.AssertFileContent("public/nn/p1/index.html", "layouts/_default/single.nn.html")
	b.AssertFileContent("public/sv/p1/index.html", "layouts/_default/single.sv.html")
}

func TestLayoutWithLanguagesLegacyMounts(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "sitemap", "rss"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[languages.sv]
weight = 3


[[module.mounts]]
source = 'layouts/en'
target = 'layouts'
lang = 'en'

[[module.mounts]]
source = 'layouts/nn'
target = 'layouts'
lang = 'nn'
[[module.mounts]]
source = 'layouts/sv'
target = 'layouts'
lang = 'sv'

-- layouts/en/_default/single.html --
layouts/en/_default/single.html
-- layouts/nn/_default/single.html --
layouts/nn/_default/single.html
-- layouts/sv/_default/single.html --
layouts/sv/_default/single.html
-- content/p1.md --
---
title: "P1"
---
-- content/p1.nn.md --
---
title: "P1 NN"
---
-- content/p1.sv.md --
---
title: "P1 SV"
---
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/p1/index.html", "layouts/en/_default/single.html")
	b.AssertFileContent("public/nn/p1/index.html", "layouts/nn/_default/single.html")
	b.AssertFileContent("public/sv/p1/index.html", "layouts/sv/_default/single.html")
}

func TestLegacyPartialIssue13599(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/partials/mypartial.html --
Mypartial.
-- layouts/index.html --
mypartial:   {{ template "partials/mypartial.html" . }}

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Mypartial.")
}
