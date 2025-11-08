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
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDraftAndFutureRender(t *testing.T) {
	t.Parallel()

	basefiles := `
-- content/sect/doc1.md --
---
title: doc1
draft: true
publishdate: "2414-05-29"
---
# doc1
*some content*
-- content/sect/doc2.md --
---
title: doc2
draft: true
publishdate: "2012-05-29"
---
# doc2
*some content*
-- content/sect/doc3.md --
---
title: doc3
draft: false
publishdate: "2414-05-29"
---
# doc3
*some content*
-- content/sect/doc4.md --
---
title: doc4
draft: false
publishdate: "2012-05-29"
---
# doc4
*some content*
`

	t.Run("defaults", func(t *testing.T) {
		t.Parallel()
		files := `
-- hugo.toml --
baseURL = "http://auth/bub"
` + basefiles

		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()
		b.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 1)
	})

	t.Run("buildDrafts", func(t *testing.T) {
		t.Parallel()
		files := `
-- hugo.toml --
baseURL = "http://auth/bub"
buildDrafts = true
` + basefiles
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()
		b.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 2)
	})

	t.Run("buildFuture", func(t *testing.T) {
		t.Parallel()
		files := `
-- hugo.toml --
baseURL = "http://auth/bub"
buildFuture = true
` + basefiles
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()
		b.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 2)
	})

	t.Run("buildDrafts and buildFuture", func(t *testing.T) {
		t.Parallel()
		files := `
-- hugo.toml --
baseURL = "http://auth/bub"
buildDrafts = true
buildFuture = true
` + basefiles
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()
		b.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 4)
	})
}

func TestFutureExpirationRender(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://auth/bub"
-- content/sect/doc3.md --
---
title: doc1
expirydate: "2400-05-29"
---
# doc1
*some content*
-- content/sect/doc4.md --
---
title: doc2
expirydate: "2000-05-29"
---
# doc2
*some content*
`
	b := Test(t, files)

	b.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 1)
	b.Assert(b.H.Sites[0].RegularPages()[0].Title(), qt.Equals, "doc1")
}

func TestLastChange(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/sect/doc1.md --
---
title: doc1
weight: 1
date: 2014-05-29
---
# doc1
*some content*
-- content/sect/doc2.md --
---
title: doc2
weight: 2
date: 2015-05-29
---
# doc2
*some content*
-- content/sect/doc3.md --
---
title: doc3
weight: 3
date: 2017-05-29
---
# doc3
*some content*
-- content/sect/doc4.md --
---
title: doc4
weight: 4
date: 2016-05-29
---
# doc4
*some content*
-- content/sect/doc5.md --
---
title: doc5
weight: 3
---
# doc5
*some content*
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{SkipRender: true},
		},
	).Build()

	b.Assert(b.H.Sites[0].Lastmod().IsZero(), qt.Equals, false)
	b.Assert(b.H.Sites[0].Lastmod().Year(), qt.Equals, 2017)
}

// Issue #_index
func TestPageWithUnderScoreIndexInFilename(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/sect/my_index_file.md --
---
title: doc1
weight: 1
date: 2014-05-29
---
# doc1
*some content*
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{SkipRender: true},
		},
	).Build()

	b.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 1)
}

// Issue #939
// Issue #1923
func TestShouldAlwaysHaveUglyURLs(t *testing.T) {
	t.Parallel()

	basefiles := `
-- layouts/index.html --
Home Sweet {{ if.IsHome  }}Home{{ end }}.
-- layouts/_default/single.html --
{{.Content}}{{ if.IsHome  }}This is not home!{{ end }}
-- layouts/404.html --
Page Not Found.{{ if.IsHome  }}This is not home!{{ end }}
-- layouts/rss.xml --
<root>RSS</root>
-- layouts/sitemap.xml --
<root>SITEMAP</root>
-- content/sect/doc1.md --
---
markup: markdown
---
# title
some *content*
-- content/sect/doc2.md --
---
url: /ugly.html
markup: markdown
---
# title
doc2 *content*
`

	t.Run("uglyURLs=true", func(t *testing.T) {
		t.Parallel()
		files := `
-- hugo.toml --
baseURL = "http://auth/bub"
uglyURLs = true
` + basefiles
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()

		b.AssertFileContent("public/index.html", "Home Sweet Home.")
		b.AssertFileContent("public/sect/doc1.html", "<h1 id=\"title\">title</h1>\n<p>some <em>content</em></p>\n")
		b.AssertFileContent("public/404.html", "Page Not Found.")
		b.AssertFileContent("public/index.xml", "<root>RSS</root>")
		b.AssertFileContent("public/sitemap.xml", "<root>SITEMAP</root>")
		b.AssertFileContent("public/ugly.html", "<h1 id=\"title\">title</h1>\n<p>doc2 <em>content</em></p>\n")

		for _, p := range b.H.Sites[0].RegularPages() {
			b.Assert(p.IsHome(), qt.Equals, false)
		}
	})

	t.Run("uglyURLs=false", func(t *testing.T) {
		t.Parallel()
		files := `
-- hugo.toml --
baseURL = "http://auth/bub"
uglyURLs = false
` + basefiles
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()

		b.AssertFileContent("public/index.html", "Home Sweet Home.")
		b.AssertFileContent("public/sect/doc1/index.html", "<h1 id=\"title\">title</h1>\n<p>some <em>content</em></p>\n")
		b.AssertFileContent("public/404.html", "Page Not Found.")
		b.AssertFileContent("public/index.xml", "<root>RSS</root>")
		b.AssertFileContent("public/sitemap.xml", "<root>SITEMAP</root>")
		b.AssertFileContent("public/ugly.html", "<h1 id=\"title\">title</h1>\n<p>doc2 <em>content</em></p>\n")

		for _, p := range b.H.Sites[0].RegularPages() {
			b.Assert(p.IsHome(), qt.Equals, false)
		}
	})
}

func TestMainSectionsMoveToSite(t *testing.T) {
	t.Run("defined in params", func(t *testing.T) {
		t.Parallel()

		files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
[params]
mainSections=["a", "b"]
-- content/mysect/page1.md --
-- layouts/index.html --
{{/* Behaviour before Hugo 0.112.0. */}}
MainSections Params: {{ site.Params.mainSections }}|
MainSections Site method: {{ site.MainSections }}|


	`

		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
MainSections Params: [a b]|
MainSections Site method: [a b]|
	`)
	})

	t.Run("defined in top level config", func(t *testing.T) {
		t.Parallel()

		files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
mainSections=["a", "b"]
[params]
[params.sub]
mainSections=["c", "d"]
-- content/mysect/page1.md --
-- layouts/index.html --
{{/* Behaviour before Hugo 0.112.0. */}}
MainSections Params: {{ site.Params.mainSections }}|
MainSections Param sub: {{ site.Params.sub.mainSections }}|
MainSections Site method: {{ site.MainSections }}|


`

		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
MainSections Params: [a b]|
MainSections Param sub: [c d]|
MainSections Site method: [a b]|
`)
	})

	t.Run("guessed from pages", func(t *testing.T) {
		t.Parallel()

		files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
-- content/mysect/page1.md --
-- layouts/index.html --
MainSections Params: {{ site.Params.mainSections }}|
MainSections Site method: {{ site.MainSections }}|


	`

		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
MainSections Params: [mysect]|
MainSections Site method: [mysect]|
	`)
	})
}

func TestRelRefWithTrailingSlash(t *testing.T) {
	files := `
-- hugo.toml --
-- content/docs/5.3/examples/_index.md --
---
title: "Examples"
---
-- content/_index.md --
---
title: "Home"
---

Examples: {{< relref "/docs/5.3/examples/" >}}
-- layouts/home.html --
Content: {{ .Content }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Examples: /docs/5.3/examples/")
}

// https://github.com/gohugoio/hugo/issues/6952
func TestRefIssues(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com"
-- content/post/b1/index.md --
---
title: pb1
---
Ref: {{< ref "b2" >}}
-- content/post/b2/index.md --
---
title: pb2
---
-- content/post/nested-a/content-a.md --
---
title: ca
---
{{< ref "content-b" >}}
-- content/post/nested-b/content-b.md --
---
title: ca
---
-- layouts/index.html --
Home
-- layouts/_default/single.html --
Content: {{ .Content }}
`

	b := Test(t, files)

	b.AssertFileContent("public/post/b1/index.html", `Content: <p>Ref: http://example.com/post/b2/</p>`)
	b.AssertFileContent("public/post/nested-a/content-a/index.html", `Content: http://example.com/post/nested-b/content-b/`)
}

func TestClassCollector(t *testing.T) {
	for _, minify := range []bool{false, true} {
		t.Run(fmt.Sprintf("minify-%t", minify), func(t *testing.T) {
			statsFilename := "hugo_stats.json"
			defer os.Remove(statsFilename)

			files := fmt.Sprintf(`
-- hugo.toml --
minify = %t

[build]
  writeStats = true
-- layouts/index.html --

<div id="el1" class="a b c">Foo</div>

Some text.

<div class="c d e [&>p]:text-red-600" id="el2">Foo</div>

<span class=z>FOO</span>

 <a class="text-base hover:text-gradient inline-block px-3 pb-1 rounded lowercase" href="{{ .RelPermalink }}">{{ .Title }}</a>

-- content/p1.md --
`, minify)

			b := Test(t, files, TestOptOsFs())

			b.AssertFileContent("hugo_stats.json", `
 {
          "htmlElements": {
            "tags": [
              "a",
              "div",
              "span"
            ],
            "classes": [
              "a",
              "b",
              "c",
              "d",
              "e",
              "hover:text-gradient",
			  "[&>p]:text-red-600",
              "inline-block",
              "lowercase",
              "pb-1",
              "px-3",
              "rounded",
              "text-base",
              "z"
            ],
            "ids": [
              "el1",
              "el2"
            ]
          }
        }
`)
		})
	}
}

func TestClassCollectorConfigWriteStats(t *testing.T) {
	r := func(writeStatsConfig string) *IntegrationTestBuilder {
		files := `
-- hugo.toml --
` + writeStatsConfig + `
-- layouts/_default/list.html --
<div id="myid" class="myclass">Foo</div>

`
		b := Test(t, files, TestOptOsFs())
		return b
	}

	// Legacy config.
	var b *IntegrationTestBuilder // Declare 'b' once
	b = r(`
[build]
writeStats = true
`)

	b.AssertFileContent("hugo_stats.json", "myclass", "div", "myid")

	b = r(`
[build]
writeStats = false
	`)

	b.AssertFileExists("public/hugo_stats.json", false)

	b = r(`
[build.buildStats]
enable = true
`)

	b.AssertFileContent("hugo_stats.json", "myclass", "div", "myid")

	b = r(`
[build.buildStats]
enable = true
disableids = true
`)

	b.AssertFileContent("hugo_stats.json", "myclass", "div", "! myid")

	b = r(`
[build.buildStats]
enable = true
disableclasses = true
`)

	b.AssertFileContent("hugo_stats.json", "! myclass", "div", "myid")

	b = r(`
[build.buildStats]
enable = true
disabletags = true
	`)

	b.AssertFileContent("hugo_stats.json", "myclass", "! div", "myid")

	b = r(`
[build.buildStats]
enable = true
disabletags = true
disableclasses = true
	`)

	b.AssertFileContent("hugo_stats.json", "! myclass", "! div", "myid")

	b = r(`
[build.buildStats]
enable = false
	`)
	b.AssertFileExists("public/hugo_stats.json", false)
}
