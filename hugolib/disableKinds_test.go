// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/resources/page"
)

func TestDisable(t *testing.T) {
	c := qt.New(t)

	newSitesBuilder := func(c *qt.C, disableKind string) *sitesBuilder {
		config := fmt.Sprintf(`
baseURL = "http://example.com/blog"
enableRobotsTXT = true
ignoreErrors = ["error-disable-taxonomy"]
disableKinds = [%q]
`, disableKind)

		b := newTestSitesBuilder(c)
		b.WithTemplatesAdded("_default/single.html", `single`)
		b.WithConfigFile("toml", config).WithContent("sect/page.md", `
---
title: Page
categories: ["mycat"]
tags: ["mytag"]
---

`, "sect/no-list.md", `
---
title: No List
_build:
  list: false
---

`, "sect/no-render.md", `
---
title: No List
_build:
  render: false
---
`,
			"sect/no-render-link.md", `
---
title: No Render Link
aliases: ["/link-alias"]
_build:
  render: link
---
`,
			"sect/no-publishresources/index.md", `
---
title: No Publish Resources
_build:
  publishResources: false
---

`, "sect/headlessbundle/index.md", `
---
title: Headless
headless: true
---


`, "headless-local/_index.md", `
---
title: Headless Local Lists
cascade:
    _build:
        render: false
        list: local
        publishResources: false
---

`, "headless-local/headless-local-page.md", "---\ntitle: Headless Local Page\n---",
			"headless-local/sub/_index.md", `
---
title: Headless Local Lists Sub
---

`, "headless-local/sub/headless-local-sub-page.md", "---\ntitle: Headless Local Sub Page\n---",
		)

		b.WithSourceFile("content/sect/headlessbundle/data.json", "DATA")
		b.WithSourceFile("content/sect/no-publishresources/data.json", "DATA")

		return b
	}

	getPage := func(b *sitesBuilder, ref string) page.Page {
		b.Helper()
		p, err := b.H.Sites[0].getPage(nil, ref)
		b.Assert(err, qt.IsNil)
		return p
	}

	getPageInSitePages := func(b *sitesBuilder, ref string) page.Page {
		b.Helper()
		for _, pages := range []page.Pages{b.H.Sites[0].Pages(), b.H.Sites[0].RegularPages()} {
			for _, p := range pages {
				if ref == p.Path() {
					return p
				}
			}
		}
		return nil
	}

	getPageInPagePages := func(p page.Page, ref string, pageCollections ...page.Pages) page.Page {
		if len(pageCollections) == 0 {
			pageCollections = []page.Pages{p.Pages(), p.RegularPages(), p.RegularPagesRecursive(), p.Sections()}
		}
		for _, pages := range pageCollections {
			for _, p := range pages {
				if ref == p.Path() {
					return p
				}
			}
		}
		return nil
	}

	disableKind := kinds.KindPage
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		s := b.H.Sites[0]
		b.Assert(getPage(b, "/sect/page.md"), qt.IsNil)
		b.Assert(b.CheckExists("public/sect/page/index.html"), qt.Equals, false)
		b.Assert(getPageInSitePages(b, "/sect/page.md"), qt.IsNil)
		b.Assert(getPageInPagePages(getPage(b, "/"), "/sect/page.md"), qt.IsNil)

		// Also check the side effects
		b.Assert(b.CheckExists("public/categories/mycat/index.html"), qt.Equals, false)
		b.Assert(len(s.Taxonomies()["categories"]), qt.Equals, 0)
	})

	disableKind = kinds.KindTerm
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		s := b.H.Sites[0]
		b.Assert(b.CheckExists("public/categories/index.html"), qt.Equals, true)
		b.Assert(b.CheckExists("public/categories/mycat/index.html"), qt.Equals, false)
		b.Assert(len(s.Taxonomies()["categories"]), qt.Equals, 0)
		b.Assert(getPage(b, "/categories"), qt.Not(qt.IsNil))
		b.Assert(getPage(b, "/categories/mycat"), qt.IsNil)
	})

	disableKind = kinds.KindTaxonomy
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		s := b.H.Sites[0]
		b.Assert(b.CheckExists("public/categories/mycat/index.html"), qt.Equals, true)
		b.Assert(b.CheckExists("public/categories/index.html"), qt.Equals, false)
		b.Assert(len(s.Taxonomies()["categories"]), qt.Equals, 1)
		b.Assert(getPage(b, "/categories/mycat"), qt.Not(qt.IsNil))
		categories := getPage(b, "/categories")
		b.Assert(categories, qt.Not(qt.IsNil))
		b.Assert(categories.RelPermalink(), qt.Equals, "")
		b.Assert(getPageInSitePages(b, "/categories"), qt.IsNil)
		b.Assert(getPageInPagePages(getPage(b, "/"), "/categories"), qt.IsNil)
	})

	disableKind = kinds.KindHome
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/index.html"), qt.Equals, false)
		home := getPage(b, "/")
		b.Assert(home, qt.Not(qt.IsNil))
		b.Assert(home.RelPermalink(), qt.Equals, "")
		b.Assert(getPageInSitePages(b, "/"), qt.IsNil)
		b.Assert(getPageInPagePages(home, "/"), qt.IsNil)
		b.Assert(getPage(b, "/sect/page.md"), qt.Not(qt.IsNil))
	})

	disableKind = kinds.KindSection
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/sect/index.html"), qt.Equals, false)
		sect := getPage(b, "/sect")
		b.Assert(sect, qt.Not(qt.IsNil))
		b.Assert(sect.RelPermalink(), qt.Equals, "")
		b.Assert(getPageInSitePages(b, "/sect"), qt.IsNil)
		home := getPage(b, "/")
		b.Assert(getPageInPagePages(home, "/sect"), qt.IsNil)
		b.Assert(home.OutputFormats(), qt.HasLen, 2)
		page := getPage(b, "/sect/page.md")
		b.Assert(page, qt.Not(qt.IsNil))
		b.Assert(page.CurrentSection(), qt.Equals, sect)
		b.Assert(getPageInPagePages(sect, "/sect/page"), qt.Not(qt.IsNil))
		b.AssertFileContent("public/sitemap.xml", "sitemap")
		b.AssertFileContent("public/index.xml", "rss")
	})

	disableKind = kinds.KindRSS
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/index.xml"), qt.Equals, false)
		home := getPage(b, "/")
		b.Assert(home.OutputFormats(), qt.HasLen, 1)
	})

	disableKind = kinds.KindSitemap
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/sitemap.xml"), qt.Equals, false)
	})

	disableKind = kinds.KindStatus404
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/404.html"), qt.Equals, false)
	})

	disableKind = kinds.KindRobotsTXT
	c.Run("Disable "+disableKind, func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.WithTemplatesAdded("robots.txt", "myrobots")
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/robots.txt"), qt.Equals, false)
	})

	c.Run("Headless bundle", func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/sect/headlessbundle/index.html"), qt.Equals, false)
		b.Assert(b.CheckExists("public/sect/headlessbundle/data.json"), qt.Equals, true)
		bundle := getPage(b, "/sect/headlessbundle/index.md")
		b.Assert(bundle, qt.Not(qt.IsNil))
		b.Assert(bundle.RelPermalink(), qt.Equals, "")
		resource := bundle.Resources()[0]
		b.Assert(resource.RelPermalink(), qt.Equals, "/blog/sect/headlessbundle/data.json")
		b.Assert(bundle.OutputFormats(), qt.HasLen, 0)
		b.Assert(bundle.AlternativeOutputFormats(), qt.HasLen, 0)
	})

	c.Run("Build config, no list", func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		ref := "/sect/no-list.md"
		b.Assert(b.CheckExists("public/sect/no-list/index.html"), qt.Equals, true)
		p := getPage(b, ref)
		b.Assert(p, qt.Not(qt.IsNil))
		b.Assert(p.RelPermalink(), qt.Equals, "/blog/sect/no-list/")
		b.Assert(getPageInSitePages(b, ref), qt.IsNil)
		sect := getPage(b, "/sect")
		b.Assert(getPageInPagePages(sect, ref), qt.IsNil)
	})

	c.Run("Build config, local list", func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		ref := "/headless-local"
		sect := getPage(b, ref)
		b.Assert(sect, qt.Not(qt.IsNil))
		b.Assert(getPageInSitePages(b, ref), qt.IsNil)

		b.Assert(getPageInSitePages(b, "/headless-local"), qt.IsNil)
		b.Assert(getPageInSitePages(b, "/headless-local/headless-local-page"), qt.IsNil)

		localPageRef := ref + "/headless-local-page"

		b.Assert(getPageInPagePages(sect, localPageRef, sect.RegularPages()), qt.Not(qt.IsNil))
		b.Assert(getPageInPagePages(sect, localPageRef, sect.RegularPagesRecursive()), qt.Not(qt.IsNil))
		b.Assert(getPageInPagePages(sect, localPageRef, sect.Pages()), qt.Not(qt.IsNil))

		ref = "/headless-local/sub"

		sect = getPage(b, ref)
		b.Assert(sect, qt.Not(qt.IsNil))

		localPageRef = ref + "/headless-local-sub-page"
		b.Assert(getPageInPagePages(sect, localPageRef), qt.Not(qt.IsNil))
	})

	c.Run("Build config, no render", func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		ref := "/sect/no-render"
		b.Assert(b.CheckExists("public/sect/no-render/index.html"), qt.Equals, false)
		p := getPage(b, ref)
		b.Assert(p, qt.Not(qt.IsNil))
		b.Assert(p.RelPermalink(), qt.Equals, "")
		b.Assert(p.OutputFormats(), qt.HasLen, 0)
		b.Assert(getPageInSitePages(b, ref), qt.Not(qt.IsNil))
		sect := getPage(b, "/sect")
		b.Assert(getPageInPagePages(sect, ref), qt.Not(qt.IsNil))
	})

	c.Run("Build config, no render link", func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		ref := "/sect/no-render-link"
		b.Assert(b.CheckExists("public/sect/no-render/index.html"), qt.Equals, false)
		p := getPage(b, ref)
		b.Assert(p, qt.Not(qt.IsNil))
		b.Assert(p.RelPermalink(), qt.Equals, "/blog/sect/no-render-link/")
		b.Assert(p.OutputFormats(), qt.HasLen, 1)
		b.Assert(getPageInSitePages(b, ref), qt.Not(qt.IsNil))
		sect := getPage(b, "/sect")
		b.Assert(getPageInPagePages(sect, ref), qt.Not(qt.IsNil))

		// https://github.com/gohugoio/hugo/issues/7832
		// It should still render any aliases.
		b.AssertFileContent("public/link-alias/index.html", "refresh")
	})

	c.Run("Build config, no publish resources", func(c *qt.C) {
		b := newSitesBuilder(c, disableKind)
		b.Build(BuildCfg{})
		b.Assert(b.CheckExists("public/sect/no-publishresources/index.html"), qt.Equals, true)
		b.Assert(b.CheckExists("public/sect/no-publishresources/data.json"), qt.Equals, false)
		bundle := getPage(b, "/sect/no-publishresources/index.md")
		b.Assert(bundle, qt.Not(qt.IsNil))
		b.Assert(bundle.RelPermalink(), qt.Equals, "/blog/sect/no-publishresources/")
		b.Assert(bundle.Resources(), qt.HasLen, 1)
		resource := bundle.Resources()[0]
		b.Assert(resource.RelPermalink(), qt.Equals, "/blog/sect/no-publishresources/data.json")
	})
}

// https://github.com/gohugoio/hugo/issues/6897#issuecomment-587947078
func TestDisableRSSWithRSSInCustomOutputs(t *testing.T) {
	b := newTestSitesBuilder(t).WithConfigFile("toml", `
disableKinds = ["term", "taxonomy", "RSS"]
[outputs]
home = [ "HTML", "RSS" ]
`).Build(BuildCfg{})

	// The config above is a little conflicting, but it exists in the real world.
	// In Hugo 0.65 we consolidated the code paths and made RSS a pure output format,
	// but we should make sure to not break existing sites.
	b.Assert(b.CheckExists("public/index.xml"), qt.Equals, false)
}

func TestBundleNoPublishResources(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithTemplates("index.html", `
{{ $bundle := site.GetPage "section/bundle-false" }}
{{ $data1 := $bundle.Resources.GetMatch "data1*" }}
Data1: {{ $data1.RelPermalink }}

`)

	b.WithContent("section/bundle-false/index.md", `---\ntitle: BundleFalse
_build:
  publishResources: false
---`,
		"section/bundle-false/data1.json", "Some data1",
		"section/bundle-false/data2.json", "Some data2",
	)

	b.WithContent("section/bundle-true/index.md", `---\ntitle: BundleTrue
---`,
		"section/bundle-true/data3.json", "Some data 3",
	)

	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html", `Data1: /section/bundle-false/data1.json`)
	b.AssertFileContent("public/section/bundle-false/data1.json", `Some data1`)
	b.Assert(b.CheckExists("public/section/bundle-false/data2.json"), qt.Equals, false)
	b.AssertFileContent("public/section/bundle-true/data3.json", `Some data 3`)
}

func TestNoRenderAndNoPublishResources(t *testing.T) {
	noRenderPage := `
---
title: %s
_build:
    render: false
    publishResources: false
---
`
	b := newTestSitesBuilder(t)
	b.WithTemplatesAdded("index.html", `
{{ $page := site.GetPage "sect/no-render" }}
{{ $sect := site.GetPage "sect-no-render" }}

Page: {{ $page.Title }}|RelPermalink: {{ $page.RelPermalink }}|Outputs: {{ len $page.OutputFormats }}
Section: {{ $sect.Title }}|RelPermalink: {{ $sect.RelPermalink }}|Outputs: {{ len $sect.OutputFormats }}


`)
	b.WithContent("sect-no-render/_index.md", fmt.Sprintf(noRenderPage, "MySection"))
	b.WithContent("sect/no-render.md", fmt.Sprintf(noRenderPage, "MyPage"))

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
Page: MyPage|RelPermalink: |Outputs: 0
Section: MySection|RelPermalink: |Outputs: 0
`)

	b.Assert(b.CheckExists("public/sect/no-render/index.html"), qt.Equals, false)
	b.Assert(b.CheckExists("public/sect-no-render/index.html"), qt.Equals, false)
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
