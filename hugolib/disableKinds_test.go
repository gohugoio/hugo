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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/kinds"
)

func TestDisableKinds(t *testing.T) {
	filesForDisabledKind := func(disableKind string) string {
		return fmt.Sprintf(`
-- hugo.toml --
baseURL = "http://example.com/blog"
enableRobotsTXT = true
ignoreErrors = ["error-disable-taxonomy"]
disableKinds = ["%s"]
-- layouts/_default/single.html --
single
-- content/sect/page.md --
---
title: Page
categories: ["mycat"]
tags: ["mytag"]
---
-- content/sect/no-list.md --
---
title: No List
build:
  list: false
---
-- content/sect/no-render.md --
---
title: No List
build:
  render: false
---
-- content/sect/no-render-link.md --
---
title: No Render Link
aliases: ["/link-alias"]
build:
  render: link
---
-- content/sect/no-publishresources/index.md --
---
title: No Publish Resources
build:
  publishResources: false
---
-- content/sect/headlessbundle/index.md --
---
title: Headless
headless: true
---
-- content/headless-local/_index.md --
---
title: Headless Local Lists
cascade:
    build:
        render: false
        list: local
        publishResources: false
---
-- content/headless-local/headless-local-page.md --
---
title: Headless Local Page
---
-- content/headless-local/sub/_index.md --
---
title: Headless Local Lists Sub
---
-- content/headless-local/sub/headless-local-sub-page.md --
---
title: Headless Local Sub Page
---
-- content/sect/headlessbundle/data.json --
DATA
-- content/sect/no-publishresources/data.json --
DATA
`, disableKind)
	}
	t.Run("Disable "+kinds.KindPage, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindPage)
		b := Test(t, files)
		s := b.H.Sites[0]
		b.AssertFileExists("public/sect/page/index.html", false)
		b.AssertFileExists("public/categories/mycat/index.html", false)
		b.Assert(len(s.Taxonomies()["categories"]), qt.Equals, 0)
	})

	t.Run("Disable "+kinds.KindTerm, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindTerm)
		b := Test(t, files)
		s := b.H.Sites[0]
		b.AssertFileExists("public/categories/index.html", false)
		b.AssertFileExists("public/categories/mycat/index.html", false)
		b.Assert(len(s.Taxonomies()["categories"]), qt.Equals, 0)
	})

	t.Run("Disable "+kinds.KindTaxonomy, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindTaxonomy)
		b := Test(t, files)
		s := b.H.Sites[0]
		b.AssertFileExists("public/categories/mycat/index.html", false)
		b.AssertFileExists("public/categories/index.html", false)
		b.Assert(len(s.Taxonomies()["categories"]), qt.Equals, 1)
	})

	t.Run("Disable "+kinds.KindHome, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindHome)
		b := Test(t, files)
		b.AssertFileExists("public/index.html", false)
	})

	t.Run("Disable "+kinds.KindSection, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindSection)
		b := Test(t, files)
		b.AssertFileExists("public/sect/index.html", false)
		b.AssertFileContent("public/sitemap.xml", "sitemap")
		b.AssertFileContent("public/index.xml", "rss")
	})

	t.Run("Disable "+kinds.KindRSS, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindRSS)
		b := Test(t, files)
		b.AssertFileExists("public/index.xml", false)
	})

	t.Run("Disable "+kinds.KindSitemap, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindSitemap)
		b := Test(t, files)
		b.AssertFileExists("public/sitemap.xml", false)
	})

	t.Run("Disable "+kinds.KindStatus404, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindStatus404)
		b := Test(t, files)
		b.AssertFileExists("public/404.html", false)
	})

	t.Run("Disable "+kinds.KindRobotsTXT, func(t *testing.T) {
		files := filesForDisabledKind(kinds.KindRobotsTXT)
		b := Test(t, files)
		b.AssertFileExists("public/robots.txt", false)
	})

	t.Run("Headless bundle", func(t *testing.T) {
		files := filesForDisabledKind("")
		b := Test(t, files)
		b.AssertFileExists("public/sect/headlessbundle/index.html", false)
		b.AssertFileExists("public/sect/headlessbundle/data.json", true)
	})

	t.Run("Build config, no list", func(t *testing.T) {
		files := filesForDisabledKind("")
		b := Test(t, files)
		b.AssertFileExists("public/sect/no-list/index.html", true)
	})

	t.Run("Build config, local list", func(t *testing.T) {
		files := filesForDisabledKind("")
		b := Test(t, files)
		// Assert that the pages are not rendered to disk, as list:local implies.
		b.AssertFileExists("public/headless-local/index.html", false)
		b.AssertFileExists("public/headless-local/headless-local-page/index.html", false)
		b.AssertFileExists("public/headless-local/sub/index.html", false)
		b.AssertFileExists("public/headless-local/sub/headless-local-sub-page/index.html", false)
	})

	t.Run("Build config, no render", func(t *testing.T) {
		files := filesForDisabledKind("")
		b := Test(t, files)
		b.AssertFileExists("public/sect/no-render/index.html", false)
	})

	t.Run("Build config, no render link", func(t *testing.T) {
		files := filesForDisabledKind("")
		b := Test(t, files)
		b.AssertFileExists("public/sect/no-render/index.html", false)
		b.AssertFileContent("public/link-alias/index.html", "refresh")
	})

	t.Run("Build config, no publish resources", func(t *testing.T) {
		files := filesForDisabledKind("")
		b := Test(t, files)
		b.AssertFileExists("public/sect/no-publishresources/index.html", true)
		b.AssertFileExists("public/sect/no-publishresources/data.json", false)
	})
}

// https://github.com/gohugoio/hugo/issues/6897#issuecomment-587947078
func TestDisableRSSWithRSSInCustomOutputs(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["term", "taxonomy", "RSS"]
[outputs]
home = [ "HTML", "RSS" ]
-- layouts/index.html --
Home
`
	b := Test(t, files)

	// The config above is a little conflicting, but it exists in the real world.
	// In Hugo 0.65 we consolidated the code paths and made RSS a pure output format,
	// but we should make sure to not break existing sites.
	b.AssertFileExists("public/index.xml", false)
}

func TestBundleNoPublishResources(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "http://example.com"
-- layouts/index.html --
{{ $bundle := site.GetPage "section/bundle-false" }}
{{ $data1 := $bundle.Resources.GetMatch "data1*" }}
Data1: {{ $data1.RelPermalink }}
-- content/section/bundle-false/index.md --
---
title: BundleFalse
build:
  publishResources: false
---
-- content/section/bundle-false/data1.json --
Some data1
-- content/section/bundle-false/data2.json --
Some data2
-- content/section/bundle-true/index.md --
---
title: BundleTrue
---
-- content/section/bundle-true/data3.json --
Some data 3
`
	b := Test(t, files)
	b.AssertFileContent("public/index.html", `Data1: /section/bundle-false/data1.json`)
	b.AssertFileContent("public/section/bundle-false/data1.json", `Some data1`)
	b.AssertFileExists("public/section/bundle-false/data2.json", false)
	b.AssertFileContent("public/section/bundle-true/data3.json", `Some data 3`)
}

func TestNoRenderAndNoPublishResources(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "http://example.com"
-- layouts/index.html --
{{ $page := site.GetPage "sect/no-render" }}
{{ $sect := site.GetPage "sect-no-render" }}

Page: {{ $page.Title }}|RelPermalink: {{ $page.RelPermalink }}|Outputs: {{ len $page.OutputFormats }}
Section: {{ $sect.Title }}|RelPermalink: {{ $sect.RelPermalink }}|Outputs: {{ len $sect.OutputFormats }}
-- content/sect-no-render/_index.md --
---
title: MySection
build:
    render: false
    publishResources: false
---
-- content/sect/no-render.md --
---
title: MyPage
build:
    render: false
    publishResources: false
---
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", `
Page: MyPage|RelPermalink: |Outputs: 0
Section: MySection|RelPermalink: |Outputs: 0
`)

	b.AssertFileExists("public/sect/no-render/index.html", false)
	b.AssertFileExists("public/sect-no-render/index.html", false)
}

func TestDisableOneOfThreeLanguages(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
title = "English"
[languages.nn]
weight = 2
title = "Nynorsk"
disabled = true
[languages.nb]
weight = 3
title = "Bokmål"
-- content/p1.nn.md --
---
title: "Page 1 nn"
---
-- content/p1.nb.md --
---
title: "Page 1 nb"
---
-- content/p1.en.md --
---
title: "Page 1 en"
---
-- content/p2.nn.md --
---
title: "Page 2 nn"
---
-- layouts/_default/single.html --
{{ .Title }}
`
	b := Test(t, files)

	b.Assert(len(b.H.Sites), qt.Equals, 2)
	b.AssertFileContent("public/en/p1/index.html", "Page 1 en")
	b.AssertFileContent("public/nb/p1/index.html", "Page 1 nb")

	b.AssertFileExists("public/en/p2/index.html", false)
	b.AssertFileExists("public/nn/p1/index.html", false)
	b.AssertFileExists("public/nn/p2/index.html", false)
}
