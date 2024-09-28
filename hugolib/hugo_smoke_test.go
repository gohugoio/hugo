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
	"math/rand"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
)

// The most basic build test.
func TestHello(t *testing.T) {
	files := `
-- hugo.toml --
title = "Hello"
baseURL="https://example.org"
disableKinds = ["term", "taxonomy", "section", "page"]
-- content/p1.md --
---
title: Page
---
-- layouts/index.html --
Home: {{ .Title }}
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			// LogLevel:    logg.LevelTrace,
		},
	).Build()

	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/index.html", `Hello`)
}

func TestSmokeOutputFormats(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
defaultContentLanguage = "en"
disableKinds = ["term", "taxonomy", "robotsTXT", "sitemap"]
[outputs]
home = ["html",  "rss"]
section = ["html", "rss"]
page = ["html"]
-- content/p1.md --
---
title: Page
---
Page.

-- layouts/_default/list.html --
List: {{ .Title }}|{{ .RelPermalink}}|{{ range .OutputFormats }}{{ .Name }}: {{ .RelPermalink }}|{{ end }}$
-- layouts/_default/list.xml --
List xml: {{ .Title }}|{{ .RelPermalink}}|{{ range .OutputFormats }}{{ .Name }}: {{ .RelPermalink }}|{{ end }}$
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .RelPermalink}}|{{ range .OutputFormats }}{{ .Name }}: {{ .RelPermalink }}|{{ end }}$

`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", `List: |/|html: /|rss: /index.xml|$`)
	b.AssertFileContent("public/index.xml", `List xml: |/|html: /|rss: /index.xml|$`)
	b.AssertFileContent("public/p1/index.html", `Single: Page|/p1/|html: /p1/|$`)
	b.AssertFileExists("public/p1/index.xml", false)
}

func TestSmoke(t *testing.T) {
	t.Parallel()

	// Basic test cases.
	// OK translations
	// OK page collections
	// OK next, prev in section
	// OK GetPage
	// OK Pagination
	// OK RenderString with shortcode
	// OK cascade
	// OK site last mod, section last mod.
	// OK main sections
	// OK taxonomies
	// OK GetTerms
	// OK Resource page
	// OK Resource txt

	const files = `
-- hugo.toml --
baseURL = "https://example.com"
title = "Smoke Site"
rssLimit = 3
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
enableRobotsTXT = true

[pagination]
pagerSize = 1

[taxonomies]
category = 'categories'
tag = 'tags'

[languages]
[languages.en]
weight = 1
title = "In English"
[languages.no]
weight = 2
title = "På norsk"

[params]
hugo = "Rules!"

[outputs]
  home = ["html", "json", "rss"]
-- layouts/index.html --
Home: {{ .Lang}}|{{ .Kind }}|{{ .RelPermalink }}|{{ .Title }}|{{ .Content }}|Len Resources: {{ len .Resources }}|HTML
Resources: {{ range .Resources }}{{ .ResourceType }}|{{ .RelPermalink }}|{{ .MediaType }} - {{ end }}|
Site last mod: {{ site.Lastmod.Format "2006-02-01" }}|
Home last mod: {{ .Lastmod.Format "2006-02-01" }}|
Len Translations: {{ len .Translations }}|
Len home.RegularPagesRecursive: {{ len .RegularPagesRecursive }}|
RegularPagesRecursive: {{ range .RegularPagesRecursive }}{{ .RelPermalink }}|{{ end }}@
Len site.RegularPages: {{ len site.RegularPages }}|
Len site.Pages: {{ len site.Pages }}|
Len site.AllPages: {{ len site.AllPages }}|
GetPage: {{ with .Site.GetPage "posts/p1" }}{{ .RelPermalink }}|{{ .Title }}{{ end }}|
RenderString with shortcode: {{ .RenderString "{{% hello %}}" }}|
Paginate: {{ .Paginator.PageNumber }}/{{ .Paginator.TotalPages }}|
-- layouts/index.json --
Home:{{ .Lang}}|{{ .Kind }}|{{ .RelPermalink }}|{{ .Title }}|{{ .Content }}|Len Resources: {{ len .Resources }}|JSON
-- layouts/_default/list.html --
List: {{ .Lang}}|{{ .Kind }}|{{ .RelPermalink }}|{{ .Title }}|{{ .Content }}|Len Resources: {{ len .Resources }}|
Resources: {{ range .Resources }}{{ .ResourceType }}|{{ .RelPermalink }}|{{ .MediaType }} - {{ end }}
Pages Length: {{ len .Pages }}
RegularPages Length: {{ len .RegularPages }}
RegularPagesRecursive Length: {{ len .RegularPagesRecursive }}
List last mod: {{ .Lastmod.Format "2006-02-01" }}
Background: {{ .Params.background }}|
Kind: {{ .Kind }}
Type: {{ .Type }}
Paginate: {{ .Paginator.PageNumber }}/{{ .Paginator.TotalPages }}|
-- layouts/_default/single.html --
Single: {{ .Lang}}|{{ .Kind }}|{{ .RelPermalink }}|{{ .Title }}|{{ .Content }}|Len Resources: {{ len .Resources }}|Background: {{ .Params.background }}|
Resources: {{ range .Resources }}{{ .ResourceType }}|{{ .RelPermalink }}|{{ .MediaType }}|{{ .Params }} - {{ end }}
{{ $textResource := .Resources.GetMatch "**.txt" }}
{{ with $textResource }}
Icon: {{ .Params.icon }}|
{{ $textResourceFingerprinted :=  . | fingerprint }}
Icon fingerprinted: {{ with $textResourceFingerprinted }}{{ .Params.icon }}|{{ .RelPermalink }}{{ end }}|
{{ end }}
NextInSection: {{ with .NextInSection }}{{ .RelPermalink }}|{{ .Title }}{{ end }}|
PrevInSection: {{ with .PrevInSection }}{{ .RelPermalink }}|{{ .Title }}{{ end }}|
GetTerms: {{ range .GetTerms "tags" }}name: {{ .Name }}, title: {{ .Title }}|{{ end }}
-- layouts/shortcodes/hello.html --
Hello.
-- content/_index.md --
---
title: Home in English
---
Home Content.
-- content/_index.no.md --
---
title: Hjem
cascade:
  - _target:
      kind: page
      path: /posts/**
    background: post.jpg
  - _target:
      kind: term
    background: term.jpg
---
Hjem Innhold.
-- content/posts/f1.txt --
posts f1 text.
-- content/posts/sub/f1.txt --
posts sub f1 text.
-- content/posts/p1/index.md --
+++
title = "Post 1"
lastMod = "2001-01-01"
tags = ["tag1"]
[[resources]]
src = '**'
[resources.params]
icon = 'enicon'
+++
Content 1.
-- content/posts/p1/index.no.md --
+++
title = "Post 1 no"
lastMod = "2002-02-02"
tags = ["tag1", "tag2"]
[[resources]]
src = '**'
[resources.params]
icon = 'noicon'
+++
Content 1 no.
-- content/posts/_index.md --
---
title: Posts
---
-- content/posts/p1/f1.txt --
posts p1 f1 text.
-- content/posts/p1/sub/ps1.md --
---
title: Post Sub 1
---
Content Sub 1.
-- content/posts/p2.md --
---
title: Post 2
tags: ["tag1", "tag3"]
---
Content 2.
-- content/posts/p2.no.md --
---
title: Post 2 No
---
Content 2 No.
-- content/tags/_index.md --
---
title: Tags
---
Content Tags.
-- content/tags/tag1/_index.md --
---
title: Tag 1
---
Content Tag 1.


`

	b := NewIntegrationTestBuilder(IntegrationTestConfig{
		T:           t,
		TxtarString: files,
		NeedsOsFS:   true,
		// Verbose:     true,
		// LogLevel:    logg.LevelTrace,
	}).Build()

	b.AssertFileContent("public/en/index.html",
		"Home: en|home|/en/|Home in English|<p>Home Content.</p>\n|HTML",
		"Site last mod: 2001-01-01",
		"Home last mod: 2001-01-01",
		"Translations: 1|",
		"Len home.RegularPagesRecursive: 2|",
		"Len site.RegularPages: 2|",
		"Len site.Pages: 8|",
		"Len site.AllPages: 16|",
		"GetPage: /en/posts/p1/|Post 1|",
		"RenderString with shortcode: Hello.|",
		"Paginate: 1/2|",
	)
	b.AssertFileContent("public/en/page/2/index.html", "Paginate: 2/2|")

	b.AssertFileContent("public/no/index.html",
		"Home: no|home|/no/|Hjem|<p>Hjem Innhold.</p>\n|HTML",
		"Site last mod: 2002-02-02",
		"Home last mod: 2002-02-02",
		"Translations: 1",
		"GetPage: /no/posts/p1/|Post 1 no|",
	)

	b.AssertFileContent("public/en/index.json", "Home:en|home|/en/|Home in English|<p>Home Content.</p>\n|JSON")
	b.AssertFileContent("public/no/index.json", "Home:no|home|/no/|Hjem|<p>Hjem Innhold.</p>\n|JSON")

	b.AssertFileContent("public/en/posts/p1/index.html",
		"Single: en|page|/en/posts/p1/|Post 1|<p>Content 1.</p>\n|Len Resources: 2|",
		"Resources: text|/en/posts/p1/f1.txt|text/plain|map[icon:enicon] - page||application/octet-stream|map[draft:false iscjklanguage:false title:Post Sub 1] -",
		"Icon: enicon",
		"Icon fingerprinted: enicon|/en/posts/p1/f1.e5746577af5cbfc4f34c558051b7955a9a5a795a84f1c6ab0609cb3473a924cb.txt|",
		"NextInSection: |\nPrevInSection: /en/posts/p2/|Post 2|",
		"GetTerms: name: tag1, title: Tag 1|",
	)

	b.AssertFileContent("public/no/posts/p1/index.html",
		"Resources: 1",
		"Resources: text|/en/posts/p1/f1.txt|text/plain|map[icon:noicon] -",
		"Icon: noicon",
		"Icon fingerprinted: noicon|/en/posts/p1/f1.e5746577af5cbfc4f34c558051b7955a9a5a795a84f1c6ab0609cb3473a924cb.txt|",
		"Background: post.jpg",
		"NextInSection: |\nPrevInSection: /no/posts/p2/|Post 2 No|",
	)

	b.AssertFileContent("public/en/posts/index.html",
		"List: en|section|/en/posts/|Posts||Len Resources: 2|",
		"Resources: text|/en/posts/f1.txt|text/plain - text|/en/posts/sub/f1.txt|text/plain -",
		"List last mod: 2001-01-01",
	)

	b.AssertFileContent("public/no/posts/index.html",
		"List last mod: 2002-02-02",
	)

	b.AssertFileContent("public/en/posts/p2/index.html", "Single: en|page|/en/posts/p2/|Post 2|<p>Content 2.</p>\n|",
		"|Len Resources: 0",
		"GetTerms: name: tag1, title: Tag 1|name: tag3, title: Tag3|",
	)
	b.AssertFileContent("public/no/posts/p2/index.html", "Single: no|page|/no/posts/p2/|Post 2 No|<p>Content 2 No.</p>\n|")

	b.AssertFileContent("public/no/categories/index.html",
		"Kind: taxonomy",
		"Type: categories",
	)
	b.AssertFileContent("public/no/tags/index.html",
		"Kind: taxonomy",
		"Type: tags",
	)

	b.AssertFileContent("public/no/tags/tag1/index.html",
		"Background: term.jpg",
		"Kind: term",
		"Type: tags",
		"Paginate: 1/1|",
	)

	b.AssertFileContent("public/en/tags/tag1/index.html",
		"Kind: term",
		"Type: tags",
		"Paginate: 1/2|",
	)
}

// Basic tests that verifies that the different file systems work as expected.
func TestSmokeFilesystems(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
title = "In English"
[languages.nn]
title = "På nynorsk"
[module]
[[module.mounts]]
source = "i18n"
target = "i18n"
[[module.mounts]]
source = "data"
target = "data"
[[module.mounts]]
source = "content/en"
target = "content"
lang = "en"
[[module.mounts]]
source = "content/nn"
target = "content"
lang = "nn"
[[module.imports]]
path = "mytheme"
-- layouts/index.html --
i18n s1: {{ i18n "s1" }}|
i18n s2: {{ i18n "s2" }}|
data s1: {{ site.Data.d1.s1 }}|
data s2: {{ site.Data.d1.s2 }}|
title: {{ .Title }}|
-- themes/mytheme/hugo.toml --
[[module.mounts]]
source = "i18n"
target = "i18n"
[[module.mounts]]
source = "data"
target = "data"
# i18n files both project and theme.
-- i18n/en.toml --
[s1]
other = 's1project'
-- i18n/nn.toml --
[s1]
other = 's1prosjekt'
-- themes/mytheme/i18n/en.toml --
[s1]
other = 's1theme'
[s2]
other = 's2theme'
# data files both project and theme.
-- data/d1.yaml --
s1: s1project
-- themes/mytheme/data/d1.yaml --
s1: s1theme
s2: s2theme
# Content
-- content/en/_index.md --
---
title: "Home"
---
-- content/nn/_index.md --
---
title: "Heim"
---

`
	b := Test(t, files)

	b.AssertFileContent("public/en/index.html",
		"i18n s1: s1project", "i18n s2: s2theme",
		"data s1: s1project", "data s2: s2theme",
		"title: Home",
	)

	b.AssertFileContent("public/nn/index.html",
		"i18n s1: s1prosjekt", "i18n s2: s2theme",
		"data s1: s1project", "data s2: s2theme",
		"title: Heim",
	)
}

// https://github.com/golang/go/issues/30286
func TestDataRace(t *testing.T) {
	const page = `
---
title: "The Page"
outputs: ["HTML", "JSON"]
---

The content.


	`

	b := newTestSitesBuilder(t).WithSimpleConfigFile()
	for i := 1; i <= 50; i++ {
		b.WithContent(fmt.Sprintf("blog/page%d.md", i), page)
	}

	b.WithContent("_index.md", `
---
title: "The Home"
outputs: ["HTML", "JSON", "CSV", "RSS"]
---

The content.


`)

	commonTemplate := `{{ .Data.Pages }}`

	b.WithTemplatesAdded("_default/single.html", "HTML Single: "+commonTemplate)
	b.WithTemplatesAdded("_default/list.html", "HTML List: "+commonTemplate)

	b.CreateSites().Build(BuildCfg{})
}

// This is just a test to verify that BenchmarkBaseline is working as intended.
func TestBenchmarkBaseline(t *testing.T) {
	cfg := IntegrationTestConfig{
		T:           t,
		TxtarString: benchmarkBaselineFiles(true),
	}
	b := NewIntegrationTestBuilder(cfg).Build()

	b.Assert(len(b.H.Sites), qt.Equals, 4)
	b.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 161)
	b.Assert(len(b.H.Sites[0].Pages()), qt.Equals, 197)
	b.Assert(len(b.H.Sites[2].RegularPages()), qt.Equals, 158)
	b.Assert(len(b.H.Sites[2].Pages()), qt.Equals, 194)
}

func BenchmarkBaseline(b *testing.B) {
	cfg := IntegrationTestConfig{
		T:           b,
		TxtarString: benchmarkBaselineFiles(false),
	}
	builders := make([]*IntegrationTestBuilder, b.N)

	for i := range builders {
		builders[i] = NewIntegrationTestBuilder(cfg)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builders[i].Build()
	}
}

func benchmarkBaselineFiles(leafBundles bool) string {
	rnd := rand.New(rand.NewSource(32))

	files := `
-- config.toml --
baseURL = "https://example.com"
defaultContentLanguage = 'en'

[module]
[[module.mounts]]
source = 'content/en'
target = 'content/en'
lang = 'en'
[[module.mounts]]
source = 'content/nn'
target = 'content/nn'
lang = 'nn'
[[module.mounts]]
source = 'content/no'
target = 'content/no'
lang = 'no'
[[module.mounts]]
source = 'content/sv'
target = 'content/sv'
lang = 'sv'
[[module.mounts]]
source = 'layouts'
target = 'layouts'

[languages]
[languages.en]
title = "English"
weight = 1
[languages.nn]
title = "Nynorsk"
weight = 2
[languages.no]
title = "Norsk"
weight = 3
[languages.sv]
title = "Svenska"
weight = 4
-- layouts/_default/list.html --
{{ .Title }}
{{ .Content }}
-- layouts/_default/single.html --
{{ .Title }}
{{ .Content }}
-- layouts/shortcodes/myshort.html --
{{ .Inner }}
`

	contentTemplate := `
---
title: "Page %d"
date: "2018-01-01"
weight: %d
---

## Heading 1

Duis nisi reprehenderit nisi cupidatat cillum aliquip ea id eu esse commodo et.

## Heading 2

Aliqua labore enim et sint anim amet excepteur ea dolore.

{{< myshort >}}
Hello, World!
{{< /myshort >}}

Aliqua labore enim et sint anim amet excepteur ea dolore.


`

	for _, lang := range []string{"en", "nn", "no", "sv"} {
		files += fmt.Sprintf("\n-- content/%s/_index.md --\n"+contentTemplate, lang, 1, 1, 1)
		for i, root := range []string{"", "foo", "bar", "baz"} {
			for j, section := range []string{"posts", "posts/funny", "posts/science", "posts/politics", "posts/world", "posts/technology", "posts/world/news", "posts/world/news/europe"} {
				n := i + j + 1
				files += fmt.Sprintf("\n-- content/%s/%s/%s/_index.md --\n"+contentTemplate, lang, root, section, n, n, n)
				for k := 1; k < rnd.Intn(30)+1; k++ {
					n := n + k
					ns := fmt.Sprintf("%d", n)
					if leafBundles {
						ns = fmt.Sprintf("%d/index", n)
					}
					file := fmt.Sprintf("\n-- content/%s/%s/%s/p%s.md --\n"+contentTemplate, lang, root, section, ns, n, n)
					files += file
				}
			}
		}
	}

	return files
}
