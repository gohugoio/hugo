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
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestContentMapSite(t *testing.T) {
	b := newTestSitesBuilder(t)

	pageTempl := `
---
title: "Page %d"
date: "2019-06-%02d"
lastMod: "2019-06-%02d"
categories: [%q]
---

Page content.
`
	createPage := func(i int) string {
		return fmt.Sprintf(pageTempl, i, i, i+1, "funny")
	}

	createPageInCategory := func(i int, category string) string {
		return fmt.Sprintf(pageTempl, i, i, i+1, category)
	}

	draftTemplate := `---
title: "Draft"
draft: true
---

`

	b.WithContent("_index.md", `
---
title: "Hugo Home"
cascade:
    description: "Common Description"
---

Home Content.
`)

	b.WithContent("blog/page1.md", createPage(1))
	b.WithContent("blog/page2.md", createPage(2))
	b.WithContent("blog/page3.md", createPage(3))
	b.WithContent("blog/bundle/index.md", createPage(4))
	b.WithContent("blog/bundle/data.json", "data")
	b.WithContent("blog/bundle/page.md", createPage(5))
	b.WithContent("blog/subsection/_index.md", createPage(6))
	b.WithContent("blog/subsection/subdata.json", "data")
	b.WithContent("blog/subsection/page4.md", createPage(7))
	b.WithContent("blog/subsection/page5.md", createPage(8))
	b.WithContent("blog/subsection/draft/index.md", draftTemplate)
	b.WithContent("blog/subsection/draft/data.json", "data")
	b.WithContent("blog/draftsection/_index.md", draftTemplate)
	b.WithContent("blog/draftsection/page/index.md", createPage(9))
	b.WithContent("blog/draftsection/page/folder/data.json", "data")
	b.WithContent("blog/draftsection/sub/_index.md", createPage(10))
	b.WithContent("blog/draftsection/sub/page.md", createPage(11))
	b.WithContent("docs/page6.md", createPage(12))
	b.WithContent("tags/_index.md", createPageInCategory(13, "sad"))
	b.WithContent("overlap/_index.md", createPageInCategory(14, "sad"))
	b.WithContent("overlap2/_index.md", createPage(15))

	b.WithTemplatesAdded("layouts/index.html", `
Num Regular: {{ len .Site.RegularPages }}|{{ range .Site.RegularPages }}{{ .RelPermalink }}|{{ end }}$
Main Sections: {{ .Site.Params.mainSections }}
Pag Num Pages: {{ len .Paginator.Pages }}
{{ $home := .Site.Home }}
{{ $blog := .Site.GetPage "blog" }}
{{ $categories := .Site.GetPage "categories" }}
{{ $funny := .Site.GetPage "categories/funny" }}
{{ $blogSub := .Site.GetPage "blog/subsection" }}
{{ $page := .Site.GetPage "blog/page1" }}
{{ $page2 := .Site.GetPage "blog/page2" }}
{{ $page4 := .Site.GetPage "blog/subsection/page4" }}
{{ $bundle := .Site.GetPage "blog/bundle" }}
{{ $overlap1 := .Site.GetPage "overlap" }}
{{ $overlap2 := .Site.GetPage "overlap2" }}

Home: {{ template "print-page" $home }}
Blog Section: {{ template "print-page" $blog }}
Blog Sub Section: {{ template "print-page" $blogSub }}
Page: {{ template "print-page" $page }}
Bundle: {{ template "print-page" $bundle }}
IsDescendant: true: {{ $page.IsDescendant $blog }} true: {{ $blogSub.IsDescendant $blog }} true: {{ $bundle.IsDescendant $blog }} true: {{ $page4.IsDescendant $blog }} true: {{ $blog.IsDescendant $home }} false: {{ $blog.IsDescendant $blog }} false: {{ $home.IsDescendant $blog }}
IsAncestor: true: {{ $blog.IsAncestor $page }} true: {{ $home.IsAncestor $blog }} true: {{ $blog.IsAncestor $blogSub }} true: {{ $blog.IsAncestor $bundle }} true: {{ $blog.IsAncestor $page4 }} true: {{ $home.IsAncestor $page }} false: {{ $blog.IsAncestor $blog }} false: {{ $page.IsAncestor $blog }} false: {{ $blog.IsAncestor $home }}  false: {{ $blogSub.IsAncestor $blog }}
IsDescendant overlap1: false: {{ $overlap1.IsDescendant $overlap2 }}
IsDescendant overlap2: false: {{ $overlap2.IsDescendant $overlap1 }}
IsAncestor overlap1: false: {{ $overlap1.IsAncestor $overlap2 }}
IsAncestor overlap2: false: {{ $overlap2.IsAncestor $overlap1 }}
FirstSection: {{ $blogSub.FirstSection.RelPermalink }} {{ $blog.FirstSection.RelPermalink }} {{ $home.FirstSection.RelPermalink }} {{ $page.FirstSection.RelPermalink }}
InSection: true: {{ $page.InSection $blog }} false: {{ $page.InSection $blogSub }}
Next: {{ $page2.Next.RelPermalink }}
NextInSection: {{ $page2.NextInSection.RelPermalink }}
Pages: {{ range $blog.Pages }}{{ .RelPermalink }}|{{ end }}
Sections: {{ range $home.Sections }}{{ .RelPermalink }}|{{ end }}:END
Categories: {{ range .Site.Taxonomies.categories }}{{ .Page.RelPermalink }}; {{ .Page.Title }}; {{ .Count }}|{{ end }}:END
Category Terms:  {{ $categories.Kind}}: {{ range $categories.Data.Terms.Alphabetical }}{{ .Page.RelPermalink }}; {{ .Page.Title }}; {{ .Count }}|{{ end }}:END
Category Funny:  {{ $funny.Kind}}; {{ $funny.Data.Term }}: {{ range $funny.Pages }}{{ .RelPermalink }}|{{ end }}:END
Pag Num Pages: {{ len .Paginator.Pages }}
Pag Blog Num Pages: {{ len $blog.Paginator.Pages }}
Blog Num RegularPages: {{ len $blog.RegularPages }}|{{ range $blog.RegularPages }}P: {{ .RelPermalink }}|{{ end }}
Blog Num Pages: {{ len $blog.Pages }}

Draft1: {{ if (.Site.GetPage "blog/subsection/draft") }}FOUND{{ end }}|
Draft2: {{ if (.Site.GetPage "blog/draftsection") }}FOUND{{ end }}|
Draft3: {{ if (.Site.GetPage "blog/draftsection/page") }}FOUND{{ end }}|
Draft4: {{ if (.Site.GetPage "blog/draftsection/sub") }}FOUND{{ end }}|
Draft5: {{ if (.Site.GetPage "blog/draftsection/sub/page") }}FOUND{{ end }}|

{{ define "print-page" }}{{ .Title }}|{{ .RelPermalink }}|{{ .Date.Format "2006-01-02" }}|Current Section: {{ with .CurrentSection }}{{ .Path }}{{ else }}NIL{{ end }}|Resources: {{ range .Resources }}{{ .ResourceType }}: {{ .RelPermalink }}|{{ end }}{{ end }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",

		`
	 Num Regular: 9
        Main Sections: [blog]
        Pag Num Pages: 9

      Home: Hugo Home|/|2019-06-15|Current Section: /|Resources:
        Blog Section: Blogs|/blog/|2019-06-11|Current Section: /blog|Resources:
        Blog Sub Section: Page 6|/blog/subsection/|2019-06-06|Current Section: /blog/subsection|Resources: application: /blog/subsection/subdata.json|
        Page: Page 1|/blog/page1/|2019-06-01|Current Section: /blog|Resources:
        Bundle: Page 4|/blog/bundle/|2019-06-04|Current Section: /blog|Resources: application: /blog/bundle/data.json|page: |
        IsDescendant: true: true true: true true: true true: true true: true false: false false: false
        IsAncestor: true: true true: true true: true true: true true: true true: true false: false false: false false: false  false: false
        IsDescendant overlap1: false: false
        IsDescendant overlap2: false: false
        IsAncestor overlap1: false: false
        IsAncestor overlap2: false: false
        FirstSection: /blog/ /blog/ / /blog/
        InSection: true: true false: false
        Next: /blog/page3/
        NextInSection: /blog/page3/
		Pages: /blog/subsection/|/blog/bundle/|/blog/page3/|/blog/page2/|/blog/page1/|
		Sections: /overlap2/|/overlap/|/docs/|/blog/|:END
		Categories: /categories/funny/; Funny; 12|/categories/sad/; Sad; 2|:END
        Category Terms:  taxonomy: /categories/funny/; Funny; 12|/categories/sad/; Sad; 2|:END
		Category Funny:  term; funny: /overlap2/|/docs/page6/|/blog/draftsection/sub/page/|/blog/draftsection/sub/|/blog/draftsection/page/|/blog/subsection/page5/|/blog/subsection/page4/|/blog/subsection/|/blog/bundle/|/blog/page3/|/blog/page2/|/blog/page1/|:END
 		Pag Num Pages: 9
        Pag Blog Num Pages: 4
        Blog Num RegularPages: 4
        Blog Num Pages: 5

        Draft1: |
        Draft2: FOUND|
        Draft3: FOUND|
        Draft4: FOUND|
        Draft5: FOUND|

`)
}

func TestIntegrationTestTemplate(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	files := `
-- config.toml --
title = "Integration Test"
disableKinds=["page", "section", "taxonomy", "term", "sitemap", "robotsTXT", "RSS"]
-- layouts/index.html --
Home: {{ .Title }}|

	`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           c,
			TxtarString: files,
		}).Build()

	b.AssertFileContent("public/index.html", "Home: Integration Test|")
}

// Issue #11840
func TestBundleResourceLanguageBestMatch(t *testing.T) {
	files := `
-- hugo.toml --
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.fr]
weight = 2
[languages.de]
weight = 3
-- layouts/index.html --
{{ $bundle := site.GetPage "bundle" }}
{{ $r := $bundle.Resources.GetMatch "*.txt" }}
{{ .Language.Lang }}: {{ $r.RelPermalink }}|{{ $r.Content }}
-- content/bundle/index.fr.md --
---
title: "Bundle Fr"
---
-- content/bundle/index.en.md --
---
title: "Bundle En"
---
-- content/bundle/index.de.md --
---
title: "Bundle De"
---
-- content/bundle/data.fr.txt --
Data fr
-- content/bundle/data.en.txt --
Data en

`
	b := Test(t, files)

	b.AssertFileContent("public/fr/index.html", "fr: /fr/bundle/data.fr.txt|Data fr")
	b.AssertFileContent("public/en/index.html", "en: /en/bundle/data.en.txt|Data en")
	b.AssertFileContent("public/de/index.html", "de: /fr/bundle/data.fr.txt|Data fr")
}

func TestBundleMultipleContentPageWithSamePath(t *testing.T) {
	files := `
-- hugo.toml --
-- content/bundle/index.md --
---
title: "Bundle md"
foo: md
---
-- content/bundle/index.html --
---
title: "Bundle html"
foo: html
---
-- content/bundle/data.txt --
Data.
-- content/p1.md --
---
title: "P1 md"
foo: md
---
-- content/p1.html --
---
title: "P1 html"
foo: html
---
-- layouts/index.html --
{{ $bundle := site.GetPage "bundle" }}
Bundle: {{ $bundle.Title }}|{{ $bundle.Params.foo }}|{{ $bundle.File.Filename }}|
{{ $p1 := site.GetPage "p1" }}
P1: {{ $p1.Title }}|{{ $p1.Params.foo }}|{{ $p1.File.Filename }}|
`

	b := Test(t, files)

	// There's multiple content files sharing the same logical path and language.
	// This is a little arbitrary, but we have to pick one and prefer the Markdown version.
	b.AssertFileContent("public/index.html",
		filepath.FromSlash("Bundle: Bundle md|md|/content/bundle/index.md|"),
		filepath.FromSlash("P1: P1 md|md|/content/p1.md|"),
	)
}

// Issue #11944
func TestBundleResourcesGetWithSpacesInFilename(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
-- content/bundle/index.md --
-- content/bundle/data with Spaces.txt --
Data.
-- layouts/index.html --
{{ $bundle := site.GetPage "bundle" }}
{{ $r := $bundle.Resources.Get "data with Spaces.txt" }}
R: {{ with $r }}{{ .Content }}{{ end }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "R: Data.")
}

// Issue #11946.
func TestBundleResourcesGetDuplicateSortOrder(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
-- content/bundle/index.md --
-- content/bundle/data-1.txt --
data-1.txt
-- content/bundle/data 1.txt --
data 1.txt
-- content/bundle/Data 1.txt --
Data 1.txt
-- content/bundle/Data-1.txt --
Data-1.txt
-- layouts/index.html --
{{ $bundle := site.GetPage "bundle" }}
{{ $r := $bundle.Resources.Get "data-1.txt" }}
R: {{ with $r }}{{ .Content }}{{ end }}|Len: {{ len $bundle.Resources }}|$

`

	for i := 0; i < 3; i++ {
		b := Test(t, files)
		b.AssertFileContent("public/index.html", "R: Data 1.txt|", "Len: 1|")
	}
}

func TestBundleResourcesNoPublishedIssue12198(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','sitemap','taxonomy','term']
-- content/s1/p1.md --
---
title: p1
---
-- content/s1/foo.txt --
foo.txt
-- content/s1/p1.txt --
p1.txt
-- content/s1/p1-foo.txt --
p1-foo.txt
-- layouts/_default/list.html --
{{.Title }}|
-- layouts/_default/single.html --
{{.Title }}|
	`

	b := Test(t, files)
	b.Build()

	b.AssertFileExists("public/s1/index.html", true)
	b.AssertFileExists("public/s1/foo.txt", true)
	b.AssertFileExists("public/s1/p1.txt", true)     // failing test
	b.AssertFileExists("public/s1/p1-foo.txt", true) // failing test
	b.AssertFileExists("public/s1/p1/index.html", true)
}

func TestSitemapOverrideFilename(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'https://example.org/'
disableKinds = ['page','rss','section','taxonomy','term']
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true
[languages.de]
[languages.en]
[sitemap]
filename = 'foo.xml'
-- layouts/index.html --
irrelevant
`

	b := Test(t, files)

	b.AssertFileExists("public/de/foo.xml", true)
	b.AssertFileExists("public/en/foo.xml", true)
	b.AssertFileContent("public/foo.xml",
		"<loc>https://example.org/de/foo.xml</loc>",
		"<loc>https://example.org/en/foo.xml</loc>",
	)

	files = strings.ReplaceAll(files, "filename = 'foo.xml'", "")

	b = Test(t, files)

	b.AssertFileExists("public/de/sitemap.xml", true)
	b.AssertFileExists("public/en/sitemap.xml", true)
	b.AssertFileContent("public/sitemap.xml",
		"<loc>https://example.org/de/sitemap.xml</loc>",
		"<loc>https://example.org/en/sitemap.xml</loc>",
	)
}
