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

package page_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
)

func TestPermalinks(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
[taxonomies]
tag = "tags"
[permalinks.page]
withpageslug = '/pageslug/:slug/'
withallbutlastsection = '/:sections[:last]/:slug/'
withallbutlastsectionslug = '/:sectionslugs[:last]/:slug/'
withsectionslug = '/sectionslug/:sectionslug/:slug/'
withsectionslugs = '/sectionslugs/:sectionslugs/:slug/'
[permalinks.section]
withfilefilename = '/sectionwithfilefilename/:contentbasename/'
withfilefiletitle = '/sectionwithfilefiletitle/:title/'
withfileslug = '/sectionwithfileslug/:slug/'
nofileslug = '/sectionnofileslug/:slug/'
nofilefilename = '/sectionnofilefilename/:contentbasename/'
nofiletitle1 = '/sectionnofiletitle1/:title/'
nofiletitle2 = '/sectionnofiletitle2/:sections[:last]/'
[permalinks.term]
tags = '/tagsslug/tag/:slug/'
[permalinks.taxonomy]
tags = '/tagsslug/:slug/'
-- content/withpageslug/p1.md --
---
slug: "p1slugvalue"
tags: ["mytag"]
---
-- content/withfilefilename/_index.md --
-- content/withfileslug/_index.md --
---
slug: "withfileslugvalue"
---
-- content/nofileslug/p1.md --
-- content/nofilefilename/p1.md --
-- content/nofiletitle1/p1.md --
-- content/nofiletitle2/asdf/p1.md --
-- content/withallbutlastsection/subsection/p1.md --
-- content/withallbutlastsectionslug/_index.md --
---
slug: "root-section-slug"
---
-- content/withallbutlastsectionslug/subsection/_index.md --
---
slug: "sub-section-slug"
---
-- content/withallbutlastsectionslug/subsection/p1.md --
---
slug: "page-slug"
---
-- content/withsectionslug/_index.md --
---
slug: "section-root-slug"
---
-- content/withsectionslug/subsection/_index.md --
-- content/withsectionslug/subsection/p1.md --
---
slug: "page1-slug"
---
-- content/withsectionslugs/_index.md --
---
slug: "sections-root-slug"
---
-- content/withsectionslugs/level1/_index.md --
---
slug: "level1-slug"
---
-- content/withsectionslugs/level1/p1.md --
---
slug: "page1-slug"
---
-- content/tags/_index.md --
---
slug: "tagsslug"
---
-- content/tags/mytag/_index.md --
---
slug: "mytagslug"
---


`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	t.Log(b.LogString())
	// No .File.TranslationBaseName on zero object etc. warnings.
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/pageslug/p1slugvalue/index.html", "Single|page|/pageslug/p1slugvalue/|")
	b.AssertFileContent("public/sectionslug/section-root-slug/page1-slug/index.html", "Single|page|/sectionslug/section-root-slug/page1-slug/|")
	b.AssertFileContent("public/sectionslugs/sections-root-slug/level1-slug/page1-slug/index.html", "Single|page|/sectionslugs/sections-root-slug/level1-slug/page1-slug/|")

	b.AssertFileContent("public/sectionwithfilefilename/withfilefilename/index.html", "List|section|/sectionwithfilefilename/withfilefilename/|")

	b.AssertFileContent("public/sectionwithfileslug/withfileslugvalue/index.html", "List|section|/sectionwithfileslug/withfileslugvalue/|")

	b.AssertFileContent("public/sectionnofilefilename/nofilefilename/index.html", "List|section|/sectionnofilefilename/nofilefilename/|")
	b.AssertFileContent("public/sectionnofileslug/nofileslugs/index.html", "List|section|/sectionnofileslug/nofileslugs/|")
	b.AssertFileContent("public/sectionnofiletitle1/nofiletitle1s/index.html", "List|section|/sectionnofiletitle1/nofiletitle1s/|")
	b.AssertFileContent("public/sectionnofiletitle2/index.html", "List|section|/sectionnofiletitle2/|")

	b.AssertFileContent("public/tagsslug/tag/mytagslug/index.html", "List|term|/tagsslug/tag/mytagslug/|")
	b.AssertFileContent("public/tagsslug/tagsslug/index.html", "List|taxonomy|/tagsslug/tagsslug/|")

	permalinksConf := b.H.Configs.Base.Permalinks
	// 5 page + 7 section + 1 taxonomy + 1 term = 14 rules.
	b.Assert(len(permalinksConf), qt.Equals, 14)
}

func TestPermalinksOldSetup(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
[permalinks]
withpageslug = '/pageslug/:slug/'
-- content/withpageslug/p1.md --
---
slug: "p1slugvalue"
---




`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	t.Log(b.LogString())
	// No .File.TranslationBaseName on zero object etc. warnings.
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/pageslug/p1slugvalue/index.html", "Single|page|/pageslug/p1slugvalue/|")

	permalinksConf := b.H.Configs.Base.Permalinks
	// Old flat format: 1 entry creates 2 rules (page + term).
	b.Assert(len(permalinksConf), qt.Equals, 2)
}

func TestPermalinksNestedSections(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[permalinks.page]
books = '/libros/:sections[1:]/:contentbasename'

[permalinks.section]
books = '/libros/:sections[1:]'
-- content/books/_index.md --
---
title: Books
---
-- content/books/fiction/_index.md --
---
title: Fiction
---
-- content/books/fiction/2023/_index.md --
---
title: 2023
---
-- content/books/fiction/2023/book1/index.md --
---
title: Book 1
---
-- layouts/single.html --
Single.
-- layouts/list.html --
List.
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	t.Log(b.LogString())
	// No .File.TranslationBaseName on zero object etc. warnings.
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)

	b.AssertFileContent("public/libros/index.html", "List.")
	b.AssertFileContent("public/libros/fiction/index.html", "List.")
	b.AssertFileContent("public/libros/fiction/2023/book1/index.html", "Single.")
}

func TestPermalinksNestedSectionsWithSlugs(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[permalinks.page]
books = '/libros/:sectionslugs[1:]/:slug'

[permalinks.section]
books = '/libros/:sectionslugs[1:]'
-- content/books/_index.md --
---
title: Books
---
-- content/books/fiction/_index.md --
---
title: Fiction
slug: fictionslug
---
-- content/books/fiction/2023/_index.md --
---
title: 2023
---
-- content/books/fiction/2023/book1/index.md --
---
title: Book One
---
-- layouts/single.html --
Single.
-- layouts/list.html --
List.
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	t.Log(b.LogString())
	// No .File.TranslationBaseName on zero object etc. warnings.
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)

	b.AssertFileContent("public/libros/index.html", "List.")
	b.AssertFileContent("public/libros/fictionslug/index.html", "List.")
	b.AssertFileContent("public/libros/fictionslug/2023/book-one/index.html", "Single.")
}

func TestPermalinksUrlCascade(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
-- content/cooking/delicious-recipes/_index.md --
---
url: /delicious-recipe/
cascade:
  url: /delicious-recipe/:slug/
---
-- content/cooking/delicious-recipes/example1.md --
---
title: Recipe 1
---
-- content/cooking/delicious-recipes/example2.md --
---
title: Recipe 2
slug: custom-recipe-2
---
`
	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	t.Log(b.LogString())
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/delicious-recipe/index.html", "List|section|/delicious-recipe/")
	b.AssertFileContent("public/delicious-recipe/recipe-1/index.html", "Single|page|/delicious-recipe/recipe-1/")
	b.AssertFileContent("public/delicious-recipe/custom-recipe-2/index.html", "Single|page|/delicious-recipe/custom-recipe-2/")
}

// Issue 12948.
// Issue 12954.
func TestPermalinksWithEscapedColons(t *testing.T) {
	t.Parallel()

	if htesting.IsWindows() {
		t.Skip("Windows does not support colons in paths")
	}

	files := `
-- hugo.toml --
disableKinds = ['home','rss','sitemap','taxonomy','term']
[permalinks.page]
s2 = "/c\\:d/:slug/"
-- content/s1/_index.md --
---
title: s1
url: "/a\\:b/:slug/"
---
-- content/s1/p1.md --
---
title: p1
url: "/a\\:b/:slug/"
---
-- content/s2/p2.md --
---
title: p2
---
-- layouts/single.html --
{{ .Title }}
-- layouts/list.html --
{{ .Title }}
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/a:b/p1/index.html", true)
	b.AssertFileExists("public/a:b/s1/index.html", true)

	// The above URLs come from the URL front matter field where everything is allowed.
	// We strip colons from paths constructed by Hugo (they are not supported on Windows).
	b.AssertFileExists("public/cd/p2/index.html", true)
}

func TestPermalinksContentbasenameContentAdapter(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[permalinks]
[permalinks.page]
a = "/:slugorcontentbasename/"
b = "/:sections/:contentbasename/"
-- content/_content.gotmpl --
{{ $.AddPage  (dict "kind" "page" "path" "a/b/contentbasename1" "title" "My A Page No Slug")  }}
{{ $.AddPage (dict "kind" "page" "path" "a/b/contentbasename2" "slug" "myslug"  "title" "My A Page With Slug")  }}
 {{ $.AddPage  (dict "kind" "section" "path" "b/c" "title" "My B Section")  }}
{{ $.AddPage  (dict "kind" "page" "path" "b/c/contentbasename3" "title" "My B Page No Slug")  }}
-- layouts/single.html --
{{ .Title }}|{{ .RelPermalink }}|{{ .Kind }}|
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/contentbasename1/index.html", "My A Page No Slug|/contentbasename1/|page|")
	b.AssertFileContent("public/myslug/index.html", "My A Page With Slug|/myslug/|page|")
}

func TestPermalinksContentbasenameWithAndWithoutFile(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[permalinks.section]
a = "/mya/:contentbasename/"
[permalinks.page]
a = "/myapage/:contentbasename/"
[permalinks.term]
categories = "/myc/:slugorcontentbasename/"
-- content/b/c/_index.md --
---
title: "C section"
---
-- content/a/b/index.md --
---
title: "My Title"
categories: ["c1", "c2"]
---
-- content/categories/c2/_index.md --
---
title: "C2"
slug: "c2slug"
---
-- layouts/single.html --
{{ .Title }}|{{ .RelPermalink }}|{{ .Kind }}|
-- layouts/list.html --
{{ .Title }}|{{ .RelPermalink }}|{{ .Kind }}|
`
	b := hugolib.Test(t, files)

	// Sections.
	b.AssertFileContent("public/mya/a/index.html", "As|/mya/a/|section|")

	// Pages.
	b.AssertFileContent("public/myapage/b/index.html", "My Title|/myapage/b/|page|")

	// Taxonomies.
	b.AssertFileContent("public/myc/c1/index.html", "C1|/myc/c1/|term|")
	b.AssertFileContent("public/myc/c2slug/index.html", "C2|/myc/c2slug/|term|")
}

func TestIssue13755(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
disablePathToLower = false
[permalinks.page]
s1 = "/:contentbasename"
-- content/s1/aBc.md --
---
title: aBc
---
-- layouts/all.html --
{{ .Title }}
`

	b := hugolib.Test(t, files)
	b.AssertFileExists("public/abc/index.html", true)

	files = strings.ReplaceAll(files, "disablePathToLower = false", "disablePathToLower = true")

	b = hugolib.Test(t, files)
	b.AssertFileExists("public/aBc/index.html", true)
}

func TestPermalinksTaxonomyAndPageConsistencyIssue14325(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
disableKinds = ['rss','sitemap']
%s
-- content/s1/p1.md --
---
title: p1
date: 2026-04-02
tags: ['tag-a']
categories: ['category-a']
---
%s
-- layouts/all.html --
{{ .Title }}|
`

	tests := []struct {
		explicitTermContent bool
		kindSpecificConfig  bool
	}{
		{true, true},
		{false, true},
		{true, false},
		{false, false},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("explicitTermContent=%t/kindSpecificConfig=%t", tt.explicitTermContent, tt.kindSpecificConfig)

		t.Run(name, func(t *testing.T) {
			var permalinkConfig, termContent string

			if tt.kindSpecificConfig {
				permalinkConfig = `
[permalinks.page]
s1 = '/:year/:month/:day/:contentbasename/'
[permalinks.term]
categories = '/:contentbasename/'
tags = '/tags/:contentbasename/'`
			} else {
				permalinkConfig = `
[permalinks]
s1 = '/:year/:month/:day/:contentbasename/'
categories = '/:contentbasename/'
tags = '/tags/:contentbasename/'`
			}

			if tt.explicitTermContent {
				termContent = `
-- content/categories/category-a/_index.md --
---
title: Category A (set in front matter)
---
-- content/tags/tag-a/_index.md --
---
title: Tag A (set in front matter)
---`
			}

			f := fmt.Sprintf(filesTemplate, permalinkConfig, termContent)
			b := hugolib.Test(t, f)

			b.AssertFileExists("public/2026/04/02/p1/index.html", true)
			b.AssertFileExists("public/categories/index.html", true)
			b.AssertFileExists("public/category-a/index.html", true)
			b.AssertFileExists("public/index.html", true)
			b.AssertFileExists("public/s1/index.html", true)
			b.AssertFileExists("public/tags/index.html", true)
			b.AssertFileExists("public/tags/tag-a/index.html", true)

			b.AssertFileExists("public/categories/category-a/index.html", false)
			b.AssertFileExists("public/s1/p1/index.html", false)
		})
	}
}

func TestPermalinksNewSliceFormat(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.yaml --
permalinks:
  - target:
      kind: page
      path: "/books/**"
    pattern: /books/:year/:slug/
  - target:
      kind: section
      path: "/{books,books/**}"
    pattern: /libros/:sections[1:]
  - target:
      kind: page
    pattern: /other/:slug/
-- content/books/_index.md --
---
title: Books
---
-- content/books/fiction/_index.md --
---
title: Fiction
---
-- content/books/fiction/book1.md --
---
title: Book One
date: 2023-06-15
slug: book-one
---
-- content/other/p1.md --
---
title: Other Page
slug: other-page
---
-- content/unmatched/p2.md --
---
title: Unmatched Page
slug: unmatched-page
---
`

	b := hugolib.Test(t, files)

	// Page in /books section gets the books-specific pattern.
	b.AssertFileContent("public/books/2023/book-one/index.html", "Single|page|/books/2023/book-one/|")
	// Section page for books/fiction gets the section pattern.
	b.AssertFileContent("public/libros/fiction/index.html", "List|section|/libros/fiction/|")
	// Page outside /books matches the default page pattern.
	b.AssertFileContent("public/other/other-page/index.html", "Single|page|/other/other-page/|")
	// Page in unmatched section also matches the default page pattern.
	b.AssertFileContent("public/other/unmatched-page/index.html", "Single|page|/other/unmatched-page/|")
}

func TestPermalinksNewSliceFormatEnvironment(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
disableKinds = ['home','rss','sitemap','taxonomy','term']
[[permalinks]]
pattern = "/testing/:slug/"
[permalinks.target]
path = "/books/**"
environment = "test"
[[permalinks]]
pattern = "/prod/:slug/"
[permalinks.target]
path = "/books/**"
environment = "production"
-- content/books/fiction/book1.md --
---
title: Book One
date: 2023-06-15
slug: book-one
---

`

	b := hugolib.Test(t, files)

	b.AssertPublishDir("prod", "! testing")
}

func TestPermalinksNewSliceFormatSitesMatrix(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|{{ .Language.Lang }}|
-- hugo.yaml --
defaultContentLanguage: en
defaultContentLanguageInSubdir: true
disableKinds: ['home','rss','section','sitemap','taxonomy','term']
languages:
  en:
    weight: 1
  de:
    weight: 2
permalinks:
  - target:
      kind: page
      sites:
        matrix:
          languages: ["en"]
    pattern: /en-posts/:slug/
  - target:
      kind: page
      sites:
        matrix:
          languages: ["de"]
    pattern: /de-posts/:slug/
  - target:
      kind: page
    pattern: /other/:slug/
-- content/p1.en.md --
---
title: Hello
slug: hello
---
-- content/p1.de.md --
---
title: Hallo
slug: hallo
---
`

	b := hugolib.Test(t, files)

	// English page matches the en-specific rule.
	b.AssertFileContent("public/en/en-posts/hello/index.html", "Single|page|/en/en-posts/hello/|en|")
	// German page matches the de-specific rule.
	b.AssertFileContent("public/de/de-posts/hallo/index.html", "Single|page|/de/de-posts/hallo/|de|")
}

// We only apply permalink patterns to kinds of type hone, page, section, taxonomy and term.
func TestPermalinksNewSliceFormatAsteriskKind(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
enableRobotsTXT = true
[[permalinks]]
pattern = "/mylink/:slug/"
[permalinks.target]
kind = "*"
-- content/mysection/p1.md --
---
title: My Page
slug: my-page
---
-- layouts/all.html --
{{ .Title }}|{{ .RelPermalink }}|{{ .Kind }}|
-- layouts/404.html --
404.
`

	b := hugolib.Test(t, files)
	b.AssertPublishDir(`
404.html
index.html
index.xml
mylink/categories/index.html 
mylink/my-page/index.html
mylink/mysections/index.html
mylink/tags/index.html 
sitemap.xml
`)
}

func TestPermalinksHomeKind(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/home.html --
Home|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.yaml --
disableKinds: ['rss','sitemap','taxonomy','term']
permalinks:
  - target:
      kind: home
    pattern: /welcome/
  - target:
      kind: section
    pattern: /s/:slug/
-- content/mysection/_index.md --
---
title: My Section
slug: my-section
---
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())
	t.Log(b.LogString())
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/welcome/index.html", "Home|home|/welcome/|")
	b.AssertFileContent("public/s/my-section/index.html", "List|section|/s/my-section/|")
}

func TestPermalinksHomeKindLegacyMap(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/home.html --
Home|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term']
[permalinks.home]
"/" = '/welcome/'
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())
	t.Log(b.LogString())
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/welcome/index.html", "Home|home|/welcome/|")
}

func TestPermalinksNewSliceFormatRootSection(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.yaml --
permalinks:
  - target:
      kind: page
      path: "/*"
    pattern: /root/:slug/
  - target:
      kind: page
    pattern: /deep/:slug/
-- content/p1.md --
---
title: Root Page
slug: root-page
---
-- content/sub/p2.md --
---
title: Sub Page
slug: sub-page
---
`

	b := hugolib.Test(t, files)

	// Root section page matches /* (non-recursive).
	b.AssertFileContent("public/root/root-page/index.html", "Single|page|/root/root-page/|")
	// Nested page matches the default rule.
	b.AssertFileContent("public/deep/sub-page/index.html", "Single|page|/deep/sub-page/|")
}
