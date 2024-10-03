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
	"reflect"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/deps"
)

func TestTaxonomiesCountOrder(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	taxonomies := make(map[string]string)
	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	cfg, fs := newTestCfg()

	cfg.Set("titleCaseStyle", "none")
	cfg.Set("taxonomies", taxonomies)
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	const pageContent = `---
tags: ['a', 'B', 'c']
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

	writeSource(t, fs, filepath.Join("content", "page.md"), pageContent)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

	st := make([]string, 0)
	for _, t := range s.Taxonomies()["tags"].ByCount() {
		st = append(st, t.Page().Title()+":"+t.Name)
	}

	expect := []string{"a:a", "B:b", "c:c"}

	if !reflect.DeepEqual(st, expect) {
		t.Fatalf("ordered taxonomies mismatch, expected\n%v\ngot\n%q", expect, st)
	}
}

func TestTaxonomiesWithAndWithoutContentFile(t *testing.T) {
	for _, uglyURLs := range []bool{false, true} {
		uglyURLs := uglyURLs
		t.Run(fmt.Sprintf("uglyURLs=%t", uglyURLs), func(t *testing.T) {
			t.Parallel()
			doTestTaxonomiesWithAndWithoutContentFile(t, uglyURLs)
		})
	}
}

func doTestTaxonomiesWithAndWithoutContentFile(t *testing.T, uglyURLs bool) {
	siteConfig := `
baseURL = "http://example.com/blog"
titleCaseStyle = "firstupper"
uglyURLs = %t
defaultContentLanguage = "en"
[pagination]
pagerSize = 1
[Taxonomies]
tag = "tags"
category = "categories"
other = "others"
empty = "empties"
permalinked = "permalinkeds"
[permalinks]
permalinkeds = "/perma/:slug/"
`

	pageTemplate := `---
title: "%s"
tags:
%s
categories:
%s
others:
%s
permalinkeds:
%s
---
# Doc
`

	siteConfig = fmt.Sprintf(siteConfig, uglyURLs)

	b := newTestSitesBuilder(t).WithConfigFile("toml", siteConfig)

	b.WithContent(
		"p1.md", fmt.Sprintf(pageTemplate, "t1/c1", "- Tag1", "- cAt1", "- o1", "- Pl1"),
		"p2.md", fmt.Sprintf(pageTemplate, "t2/c1", "- tag2", "- cAt1", "- o1", "- Pl1"),
		"p3.md", fmt.Sprintf(pageTemplate, "t2/c12", "- tag2", "- cat2", "- o1", "- Pl1"),
		"p4.md", fmt.Sprintf(pageTemplate, "Hello World", "", "", "- \"Hello Hugo world\"", "- Pl1"),
		"categories/_index.md", newTestPage("Category Terms", "2017-01-01", 10),
		"tags/Tag1/_index.md", newTestPage("Tag1 List", "2017-01-01", 10),
		// https://github.com/gohugoio/hugo/issues/5847
		"/tags/not-used/_index.md", newTestPage("Unused Tag List", "2018-01-01", 10),
	)

	b.Build(BuildCfg{})

	// So what we have now is:
	// 1. categories with terms content page, but no content page for the only c1 category
	// 2. tags with no terms content page, but content page for one of 2 tags (tag1)
	// 3. the "others" taxonomy with no content pages.
	// 4. the "permalinkeds" taxonomy with permalinks configuration.

	pathFunc := func(s string) string {
		if uglyURLs {
			return strings.Replace(s, "/index.html", ".html", 1)
		}
		return s
	}

	// 1.
	b.AssertFileContent(pathFunc("public/categories/cat1/index.html"), "List", "CAt1")
	b.AssertFileContent(pathFunc("public/categories/index.html"), "Taxonomy Term Page", "Category Terms")

	// 2.
	b.AssertFileContent(pathFunc("public/tags/tag2/index.html"), "List", "tag2")
	b.AssertFileContent(pathFunc("public/tags/tag1/index.html"), "List", "Tag1")
	b.AssertFileContent(pathFunc("public/tags/index.html"), "Taxonomy Term Page", "Tags")

	// 3.
	b.AssertFileContent(pathFunc("public/others/o1/index.html"), "List", "o1")
	b.AssertFileContent(pathFunc("public/others/index.html"), "Taxonomy Term Page", "Others")

	// 4.
	b.AssertFileContent(pathFunc("public/perma/pl1/index.html"), "List", "Pl1")

	// This looks kind of funky, but the taxonomy terms do not have a permalinks definition,
	// for good reasons.
	b.AssertFileContent(pathFunc("public/permalinkeds/index.html"), "Taxonomy Term Page", "Permalinkeds")

	s := b.H.Sites[0]

	// Make sure that each kinds.KindTaxonomyTerm page has an appropriate number
	// of kinds.KindTaxonomy pages in its Pages slice.
	taxonomyTermPageCounts := map[string]int{
		"tags":         3,
		"categories":   2,
		"others":       2,
		"empties":      0,
		"permalinkeds": 1,
	}

	for taxonomy, count := range taxonomyTermPageCounts {
		msg := qt.Commentf(taxonomy)
		term := s.getPageOldVersion(kinds.KindTaxonomy, taxonomy)
		b.Assert(term, qt.Not(qt.IsNil), msg)
		b.Assert(len(term.Pages()), qt.Equals, count, msg)

		for _, p := range term.Pages() {
			b.Assert(p.Kind(), qt.Equals, kinds.KindTerm)
		}
	}

	cat1 := s.getPageOldVersion(kinds.KindTerm, "categories", "cat1")
	b.Assert(cat1, qt.Not(qt.IsNil))
	if uglyURLs {
		b.Assert(cat1.RelPermalink(), qt.Equals, "/blog/categories/cat1.html")
	} else {
		b.Assert(cat1.RelPermalink(), qt.Equals, "/blog/categories/cat1/")
	}

	pl1 := s.getPageOldVersion(kinds.KindTerm, "permalinkeds", "pl1")
	permalinkeds := s.getPageOldVersion(kinds.KindTaxonomy, "permalinkeds")
	b.Assert(pl1, qt.Not(qt.IsNil))
	b.Assert(permalinkeds, qt.Not(qt.IsNil))
	if uglyURLs {
		b.Assert(pl1.RelPermalink(), qt.Equals, "/blog/perma/pl1.html")
		b.Assert(permalinkeds.RelPermalink(), qt.Equals, "/blog/permalinkeds.html")
	} else {
		b.Assert(pl1.RelPermalink(), qt.Equals, "/blog/perma/pl1/")
		b.Assert(permalinkeds.RelPermalink(), qt.Equals, "/blog/permalinkeds/")
	}

	helloWorld := s.getPageOldVersion(kinds.KindTerm, "others", "hello-hugo-world")
	b.Assert(helloWorld, qt.Not(qt.IsNil))
	b.Assert(helloWorld.Title(), qt.Equals, "Hello Hugo world")

	// Issue #2977
	b.AssertFileContent(pathFunc("public/empties/index.html"), "Taxonomy Term Page", "Empties")
}

// https://github.com/gohugoio/hugo/issues/5513
// https://github.com/gohugoio/hugo/issues/5571
func TestTaxonomiesPathSeparation(t *testing.T) {
	t.Parallel()

	config := `
baseURL = "https://example.com"
titleCaseStyle = "none"
[taxonomies]
"news/tag" = "news/tags"
"news/category" = "news/categories"
"t1/t2/t3" = "t1/t2/t3s"
"s1/s2/s3" = "s1/s2/s3s"
`

	pageContent := `
+++
title = "foo"
"news/categories" = ["a", "b", "c", "d/e", "f/g/h"]
"t1/t2/t3s" = ["t4/t5", "t4/t5/t6"]
+++
Content.
`

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", config)
	b.WithContent("page.md", pageContent)
	b.WithContent("news/categories/b/_index.md", `
---
title: "This is B"
---
`)

	b.WithContent("news/categories/f/g/h/_index.md", `
---
title: "This is H"
---
`)

	b.WithContent("t1/t2/t3s/t4/t5/_index.md", `
---
title: "This is T5"
---
`)

	b.WithContent("s1/s2/s3s/_index.md", `
---
title: "This is S3s"
---
`)

	b.CreateSites().Build(BuildCfg{})

	s := b.H.Sites[0]

	filterbyKind := func(kind string) page.Pages {
		var pages page.Pages
		for _, p := range s.Pages() {
			if p.Kind() == kind {
				pages = append(pages, p)
			}
		}
		return pages
	}

	ta := filterbyKind(kinds.KindTerm)
	te := filterbyKind(kinds.KindTaxonomy)

	b.Assert(len(te), qt.Equals, 4)
	b.Assert(len(ta), qt.Equals, 7)

	b.AssertFileContent("public/news/categories/a/index.html", "Taxonomy List Page 1|a|Hello|https://example.com/news/categories/a/|")
	b.AssertFileContent("public/news/categories/b/index.html", "Taxonomy List Page 1|This is B|Hello|https://example.com/news/categories/b/|")
	b.AssertFileContent("public/news/categories/d/e/index.html", "Taxonomy List Page 1|d/e|Hello|https://example.com/news/categories/d/e/|")
	b.AssertFileContent("public/news/categories/f/g/h/index.html", "Taxonomy List Page 1|This is H|Hello|https://example.com/news/categories/f/g/h/|")
	b.AssertFileContent("public/t1/t2/t3s/t4/t5/index.html", "Taxonomy List Page 1|This is T5|Hello|https://example.com/t1/t2/t3s/t4/t5/|")
	b.AssertFileContent("public/t1/t2/t3s/t4/t5/t6/index.html", "Taxonomy List Page 1|t4/t5/t6|Hello|https://example.com/t1/t2/t3s/t4/t5/t6/|")

	b.AssertFileContent("public/news/categories/index.html", "Taxonomy Term Page 1|categories|Hello|https://example.com/news/categories/|")
	b.AssertFileContent("public/t1/t2/t3s/index.html", "Taxonomy Term Page 1|t3s|Hello|https://example.com/t1/t2/t3s/|")
	b.AssertFileContent("public/s1/s2/s3s/index.html", "Taxonomy Term Page 1|This is S3s|Hello|https://example.com/s1/s2/s3s/|")
}

// https://github.com/gohugoio/hugo/issues/5719
func TestTaxonomiesNextGenLoops(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	b.WithTemplatesAdded("index.html", `
<h1>Tags</h1>
<ul>
    {{ range .Site.Taxonomies.tags }}
            <li><a href="{{ .Page.Permalink }}">{{ .Page.Title }}</a> {{ .Count }}</li>
    {{ end }}
</ul>

`)

	b.WithTemplatesAdded("_default/terms.html", `
<h1>Terms</h1>
<ul>
    {{ range .Data.Terms.Alphabetical }}
            <li><a href="{{ .Page.Permalink }}">{{ .Page.Title }}</a> {{ .Count }}</li>
    {{ end }}
</ul>
`)

	for i := 0; i < 10; i++ {
		b.WithContent(fmt.Sprintf("page%d.md", i+1), `
---
Title: "Taxonomy!"
tags: ["Hugo Rocks!", "Rocks I say!" ]
categories: ["This is Cool", "And new" ]
---

Content.

		`)
	}

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `<li><a href="http://example.com/tags/hugo-rocks/">Hugo Rocks!</a> 10</li>`)
	b.AssertFileContent("public/categories/index.html", `<li><a href="http://example.com/categories/this-is-cool/">This Is Cool</a> 10</li>`)
	b.AssertFileContent("public/tags/index.html", `<li><a href="http://example.com/tags/rocks-i-say/">Rocks I Say!</a> 10</li>`)
}

// Issue 6213
func TestTaxonomiesNotForDrafts(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithContent("draft.md", `---
title: "Draft"
draft: true
categories: ["drafts"]
---

`,
		"regular.md", `---
title: "Not Draft"
categories: ["regular"]
---

`)

	b.Build(BuildCfg{})
	s := b.H.Sites[0]

	b.Assert(b.CheckExists("public/categories/regular/index.html"), qt.Equals, true)
	b.Assert(b.CheckExists("public/categories/drafts/index.html"), qt.Equals, false)

	reg, _ := s.getPage(nil, "categories/regular")
	dra, _ := s.getPage(nil, "categories/draft")
	b.Assert(reg, qt.Not(qt.IsNil))
	b.Assert(dra, qt.IsNil)
}

func TestTaxonomiesIndexDraft(t *testing.T) {
	t.Parallel()
	b := newTestSitesBuilder(t)
	b.WithContent(
		"categories/_index.md", `---
title: "The Categories"
draft: true
---

Content.

`,
		"page.md", `---
title: "The Page"
categories: ["cool"]
---

Content.

`,
	)

	b.WithTemplates("index.html", `
{{ range .Site.Pages }}
{{ .RelPermalink }}|{{ .Title }}|{{ .WordCount }}|{{ .Content }}|
{{ end }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContentFn("public/index.html", func(s string) bool {
		return !strings.Contains(s, "/categories/|")
	})
}

// https://github.com/gohugoio/hugo/issues/6927
func TestTaxonomiesHomeDraft(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithContent(
		"_index.md", `---
title: "Home"
draft: true
---

Content.

`,
		"posts/_index.md", `---
title: "Posts"
draft: true
---

Content.

`,
		"posts/page.md", `---
title: "The Page"
categories: ["cool"]
---

Content.

`,
	)

	b.WithTemplates("index.html", `
NO HOME FOR YOU
`)

	b.Build(BuildCfg{})

	b.Assert(b.CheckExists("public/index.html"), qt.Equals, false)
	b.Assert(b.CheckExists("public/categories/index.html"), qt.Equals, false)
	b.Assert(b.CheckExists("public/posts/index.html"), qt.Equals, false)
}

// https://github.com/gohugoio/hugo/issues/6173
func TestTaxonomiesWithBundledResources(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithTemplates("_default/list.html", `
List {{ .Title }}:
{{ range .Resources }}
Resource: {{ .RelPermalink }}|{{ .MediaType }}
{{ end }}
	`)

	b.WithContent("p1.md", `---
title: Page
categories: ["funny"]
---
	`,
		"categories/_index.md", "---\ntitle: Categories Page\n---",
		"categories/data.json", "Category data",
		"categories/funny/_index.md", "---\ntitle: Funny Category\n---",
		"categories/funny/funnydata.json", "Category funny data",
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/categories/index.html", `Resource: /categories/data.json|application/json`)
	b.AssertFileContent("public/categories/funny/index.html", `Resource: /categories/funny/funnydata.json|application/json`)
}

func TestTaxonomiesRemoveOne(t *testing.T) {
	files := `
-- hugo.toml --
disableLiveReload = true
-- layouts/index.html --
{{ $cats := .Site.Taxonomies.categories.cats }}
{{ if $cats }}
Len cats: {{ len $cats }}
{{ range $cats }}
	Cats:|{{ .Page.RelPermalink }}|
{{ end }}
{{ end }}
{{ $funny := .Site.Taxonomies.categories.funny }}
{{ if $funny }}
Len funny: {{ len $funny }}
{{ range $funny }}
	Funny:|{{ .Page.RelPermalink }}|
{{ end }}
{{ end }}
-- content/p1.md --
---
title: Page
categories: ["funny", "cats"]
---
-- content/p2.md --
---
title: Page2
categories: ["funny", "cats"]
---

`
	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", `
Len cats: 2
Len funny: 2
Cats:|/p1/|
Cats:|/p2/|
Funny:|/p1/|
Funny:|/p2/|`)

	// Remove one category from one of the pages.
	b.EditFiles("content/p1.md", `---
title: Page
categories: ["funny"]
---
	`)

	b.Build()

	b.AssertFileContent("public/index.html", `
Len cats: 1
Len funny: 2
Cats:|/p2/|
Funny:|/p1/|
Funny:|/p2/|`)
}

// https://github.com/gohugoio/hugo/issues/6590
func TestTaxonomiesListPages(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithTemplates("_default/list.html", `

{{ template "print-taxo" "categories.cats" }}
{{ template "print-taxo" "categories.funny" }}

{{ define "print-taxo" }}
{{ $node := index site.Taxonomies (split $ ".") }}
{{ if $node }}
Len {{ $ }}: {{ len $node }}
{{ range $node }}
    {{ $ }}:|{{ .Page.RelPermalink }}|
{{ end }}
{{ else }}
{{ $ }} not found.
{{ end }}
{{ end }}
	`)

	b.WithContent("_index.md", `---
title: Home
categories: ["funny", "cats"]
---
	`, "blog/p1.md", `---
title: Page1
categories: ["funny"]
---
	`, "blog/_index.md", `---
title: Blog Section
categories: ["cats"]
---
	`,
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `

Len categories.cats: 2
categories.cats:|/blog/|
categories.cats:|/|

Len categories.funny: 2
categories.funny:|/|
categories.funny:|/blog/p1/|
`)
}

func TestTaxonomiesPageCollections(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithContent(
		"_index.md", `---
title: "Home Sweet Home"
categories: [ "dogs", "gorillas"]
---
`,
		"section/_index.md", `---
title: "Section"
categories: [ "cats", "dogs", "birds"]
---
`,
		"section/p1.md", `---
title: "Page1"
categories: ["funny", "cats"]
---
`, "section/p2.md", `---
title: "Page2"
categories: ["funny"]
---
`)

	b.WithTemplatesAdded("index.html", `
{{ $home := site.Home }}
{{ $section := site.GetPage "section" }}
{{ $categories := site.GetPage "categories" }}
{{ $funny := site.GetPage "categories/funny" }}
{{ $cats := site.GetPage "categories/cats" }}
{{ $p1 := site.GetPage "section/p1" }}

Categories Pages: {{ range $categories.Pages}}{{.RelPermalink }}|{{ end }}:END
Funny Pages: {{ range $funny.Pages}}{{.RelPermalink }}|{{ end }}:END
Cats Pages: {{ range $cats.Pages}}{{.RelPermalink }}|{{ end }}:END
P1 Terms: {{ range $p1.GetTerms "categories" }}{{.RelPermalink }}|{{ end }}:END
Section Terms: {{ range $section.GetTerms "categories" }}{{.RelPermalink }}|{{ end }}:END
Home Terms: {{ range $home.GetTerms "categories" }}{{.RelPermalink }}|{{ end }}:END
Category Paginator {{ range $categories.Paginator.Pages }}{{ .RelPermalink }}|{{ end }}:END
Cats Paginator {{ range $cats.Paginator.Pages }}{{ .RelPermalink }}|{{ end }}:END

`)
	b.WithTemplatesAdded("404.html", `
404 Terms: {{ range .GetTerms "categories" }}{{.RelPermalink }}|{{ end }}:END
	`)
	b.Build(BuildCfg{})

	cat := b.GetPage("categories")
	funny := b.GetPage("categories/funny")

	b.Assert(cat, qt.Not(qt.IsNil))
	b.Assert(funny, qt.Not(qt.IsNil))

	b.Assert(cat.Parent().IsHome(), qt.Equals, true)
	b.Assert(funny.Kind(), qt.Equals, "term")
	b.Assert(funny.Parent(), qt.Equals, cat)

	b.AssertFileContent("public/index.html", `
Categories Pages: /categories/birds/|/categories/cats/|/categories/dogs/|/categories/funny/|/categories/gorillas/|:END
Funny Pages: /section/p1/|/section/p2/|:END
Cats Pages: /section/p1/|/section/|:END
P1 Terms: /categories/funny/|/categories/cats/|:END
Section Terms: /categories/cats/|/categories/dogs/|/categories/birds/|:END
Home Terms: /categories/dogs/|/categories/gorillas/|:END
Cats Paginator /section/p1/|/section/|:END
Category Paginator /categories/birds/|/categories/cats/|/categories/dogs/|/categories/funny/|/categories/gorillas/|:END`,
	)
	b.AssertFileContent("public/404.html", "\n404 Terms: :END\n\t")
	b.AssertFileContent("public/categories/funny/index.xml", `<link>http://example.com/section/p1/</link>`)
	b.AssertFileContent("public/categories/index.xml", `<link>http://example.com/categories/funny/</link>`)
}

func TestTaxonomiesDirectoryOverlaps(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).WithContent(
		"abc/_index.md", "---\ntitle: \"abc\"\nabcdefgs: [abc]\n---",
		"abc/p1.md", "---\ntitle: \"abc-p\"\n---",
		"abcdefgh/_index.md", "---\ntitle: \"abcdefgh\"\n---",
		"abcdefgh/p1.md", "---\ntitle: \"abcdefgh-p\"\n---",
		"abcdefghijk/index.md", "---\ntitle: \"abcdefghijk\"\n---",
	)

	b.WithConfigFile("toml", `
baseURL = "https://example.org"
titleCaseStyle = "none"

[taxonomies]
  abcdef = "abcdefs"
  abcdefg = "abcdefgs"
  abcdefghi = "abcdefghis"
`)

	b.WithTemplatesAdded("index.html", `
{{ range site.Pages }}Page: {{ template "print-page" . }}
{{ end }}
{{ $abc := site.GetPage "abcdefgs/abc" }}
{{ $abcdefgs := site.GetPage "abcdefgs" }}
abc: {{ template "print-page" $abc }}|IsAncestor: {{ $abc.IsAncestor $abcdefgs }}|IsDescendant: {{ $abc.IsDescendant $abcdefgs }}
abcdefgs: {{ template "print-page" $abcdefgs }}|IsAncestor: {{ $abcdefgs.IsAncestor $abc }}|IsDescendant: {{ $abcdefgs.IsDescendant $abc }}

{{ define "print-page" }}{{ .RelPermalink }}|{{ .Title }}|{{.Kind }}|Parent: {{ with .Parent }}{{ .RelPermalink }}{{ end }}|CurrentSection: {{ .CurrentSection.RelPermalink}}|FirstSection: {{ .FirstSection.RelPermalink }}{{ end }}

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
    Page: /||home|Parent: |CurrentSection: /|
    Page: /abc/|abc|section|Parent: /|CurrentSection: /abc/|
    Page: /abc/p1/|abc-p|page|Parent: /abc/|CurrentSection: /abc/|
    Page: /abcdefgh/|abcdefgh|section|Parent: /|CurrentSection: /abcdefgh/|
    Page: /abcdefgh/p1/|abcdefgh-p|page|Parent: /abcdefgh/|CurrentSection: /abcdefgh/|
    Page: /abcdefghijk/|abcdefghijk|page|Parent: /|CurrentSection: /|
    Page: /abcdefghis/|abcdefghis|taxonomy|Parent: /|CurrentSection: /abcdefghis/|
    Page: /abcdefgs/|abcdefgs|taxonomy|Parent: /|CurrentSection: /abcdefgs/|
    Page: /abcdefs/|abcdefs|taxonomy|Parent: /|CurrentSection: /abcdefs/|
    abc: /abcdefgs/abc/|abc|term|Parent: /abcdefgs/|CurrentSection: /abcdefgs/abc/|
    abcdefgs: /abcdefgs/|abcdefgs|taxonomy|Parent: /|CurrentSection: /abcdefgs/|
    abc: /abcdefgs/abc/|abc|term|Parent: /abcdefgs/|CurrentSection: /abcdefgs/abc/|FirstSection: /abcdefgs/|IsAncestor: false|IsDescendant: true
    abcdefgs: /abcdefgs/|abcdefgs|taxonomy|Parent: /|CurrentSection: /abcdefgs/|FirstSection: /abcdefgs/|IsAncestor: true|IsDescendant: false
`)
}

func TestTaxonomiesWeightSort(t *testing.T) {
	files := `
-- layouts/index.html --
{{ $a := site.GetPage "tags/a"}}
:{{ range $a.Pages }}{{ .RelPermalink }}|{{ end }}:
-- content/p1.md --
---
title: P1
weight: 100
tags: ['a']
tags_weight: 20
---
-- content/p3.md --
---
title: P2
weight: 200
tags: ['a']
tags_weight: 30
---
-- content/p2.md --
---
title: P3
weight: 50
tags: ['a']
tags_weight: 40
---
	`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", `:/p1/|/p3/|/p2/|:`)
}

func TestTaxonomiesEmptyTagsString(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[taxonomies]
tag = 'tags'
-- content/p1.md --
+++
title = "P1"
tags = ''
+++
-- layouts/_default/single.html --
Single.

`
	Test(t, files)
}

func TestTaxonomiesSpaceInName(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[taxonomies]
authors = 'book authors'
-- content/p1.md --
---
title: Good Omens
book authors:
  - Neil Gaiman
  - Terry Pratchett
---
-- layouts/index.html --
{{- $taxonomy := "book authors" }}
Len Book Authors: {{ len (index .Site.Taxonomies $taxonomy) }}
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Len Book Authors: 2")
}

func TestTaxonomiesListTermsHome(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
[taxonomies]
tag = "tags"
-- content/_index.md --
---
title: "Home"
tags: ["a", "b", "c", "hello world"]
---
-- content/tags/a/_index.md --
---
title: "A"
---
-- content/tags/b/_index.md --
---
title: "B"
---
-- content/tags/c/_index.md --
---
title: "C"
---
-- content/tags/d/_index.md --
---
title: "D"
---
-- content/tags/hello-world/_index.md --
---
title: "Hello World!"
---
-- layouts/home.html --
Terms: {{ range site.Taxonomies.tags }}{{ .Page.Title }}: {{ .Count }}|{{ end }}$
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Terms: A: 1|B: 1|C: 1|Hello World!: 1|$")
}

func TestTaxonomiesTermTitleAndTerm(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
[taxonomies]
tag = "tags"
-- content/_index.md --
---
title: "Home"
tags: ["hellO world"]
---
-- layouts/_default/term.html --
{{ .Title }}|{{ .Kind }}|{{ .Data.Singular }}|{{ .Data.Plural }}|{{ .Page.Data.Term }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/tags/hello-world/index.html", "HellO World|term|tag|tags|hellO world|")
}

func TestTermDraft(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/_default/list.html --
|{{ .Title }}|
-- content/p1.md --
---
title: p1
tags: [a]
---
-- content/tags/a/_index.md --
---
title: tag-a-title-override
draft: true
---
  `

	b := Test(t, files)

	b.AssertFileExists("public/tags/a/index.html", false)
}

func TestTermBuildNeverRenderNorList(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/index.html --
|{{ len site.Taxonomies.tags }}|
-- content/p1.md --
---
title: p1
tags: [a]
---
-- content/tags/a/_index.md --
---
title: tag-a-title-override
build:
  render: never
  list: never
---

  `

	b := Test(t, files)

	b.AssertFileExists("public/tags/a/index.html", false)
	b.AssertFileContent("public/index.html", "|0|")
}

func TestTaxonomiesTermLookup(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
[taxonomies]
tag = "tags"
-- content/_index.md --
---
title: "Home"
tags: ["a", "b"]
---
-- layouts/taxonomy/tag.html --
Tag: {{ .Title }}|
-- content/tags/a/_index.md --
---
title: tag-a-title-override
---
`

	b := Test(t, files)

	b.AssertFileContent("public/tags/a/index.html", "Tag: tag-a-title-override|")
}

func TestTaxonomyLookupIssue12193(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap']
[taxonomies]
author = 'authors'
-- layouts/_default/list.html --
{{ .Title }}|
-- layouts/_default/author.terms.html --
layouts/_default/author.terms.html
-- content/authors/_index.md --
---
title: Authors Page
---
`

	b := Test(t, files)

	b.AssertFileExists("public/index.html", true)
	b.AssertFileExists("public/authors/index.html", true)
	b.AssertFileContent("public/authors/index.html", "layouts/_default/author.terms.html") // failing test
}

func TestTaxonomyNestedEmptySectionsIssue12188(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap']
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages.en]
weight = 1
[languages.ja]
weight = 2
[taxonomies]
's1/category' = 's1/category'
-- layouts/_default/single.html --
{{ .Title }}|
-- layouts/_default/list.html --
{{ .Title }}|
-- content/s1/p1.en.md --
---
title: p1
---
`

	b := Test(t, files)

	b.AssertFileExists("public/en/s1/index.html", true)
	b.AssertFileExists("public/en/s1/p1/index.html", true)
	b.AssertFileExists("public/en/s1/category/index.html", true)

	b.AssertFileExists("public/ja/s1/index.html", false) // failing test
	b.AssertFileExists("public/ja/s1/category/index.html", true)
}

func BenchmarkTaxonomiesGetTerms(b *testing.B) {
	createBuilders := func(b *testing.B, numPages int) []*IntegrationTestBuilder {
		files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["RSS", "sitemap", "section"]
[taxononomies]
tag = "tags"
-- layouts/_default/list.html --
List.
-- layouts/_default/single.html --
GetTerms.tags: {{ range .GetTerms "tags" }}{{ .Title }}|{{ end }}
-- content/_index.md --
`

		tagsVariants := []string{
			"tags: ['a']",
			"tags: ['a', 'b']",
			"tags: ['a', 'b', 'c']",
			"tags: ['a', 'b', 'c', 'd']",
			"tags: ['a', 'b',  'd', 'e']",
			"tags: ['a', 'b', 'c', 'd', 'e']",
			"tags: ['a', 'd']",
			"tags: ['a',  'f']",
		}

		for i := 1; i < numPages; i++ {
			tags := tagsVariants[i%len(tagsVariants)]
			files += fmt.Sprintf("\n-- content/posts/p%d.md --\n---\n%s\n---", i+1, tags)
		}
		cfg := IntegrationTestConfig{
			T:           b,
			TxtarString: files,
		}
		builders := make([]*IntegrationTestBuilder, b.N)

		for i := range builders {
			builders[i] = NewIntegrationTestBuilder(cfg)
		}

		b.ResetTimer()

		return builders
	}

	for _, numPages := range []int{100, 1000, 10000, 20000} {
		b.Run(fmt.Sprintf("pages_%d", numPages), func(b *testing.B) {
			builders := createBuilders(b, numPages)
			for i := 0; i < b.N; i++ {
				builders[i].Build()
			}
		})
	}
}
