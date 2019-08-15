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

	"github.com/gohugoio/hugo/resources/page"

	"reflect"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/deps"
)

func TestTaxonomiesCountOrder(t *testing.T) {
	t.Parallel()
	taxonomies := make(map[string]string)

	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	cfg, fs := newTestCfg()

	cfg.Set("taxonomies", taxonomies)

	const pageContent = `---
tags: ['a', 'B', 'c']
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

	writeSource(t, fs, filepath.Join("content", "page.md"), pageContent)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	st := make([]string, 0)
	for _, t := range s.Taxonomies["tags"].ByCount() {
		st = append(st, t.Page().Title()+":"+t.Name)
	}

	expect := []string{"a:a", "B:b", "c:c"}

	if !reflect.DeepEqual(st, expect) {
		t.Fatalf("ordered taxonomies mismatch, expected\n%v\ngot\n%q", expect, st)
	}
}

//
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
uglyURLs = %t
paginate = 1
defaultContentLanguage = "en"
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
	b.AssertFileContent(pathFunc("public/categories/cat1/index.html"), "List", "cAt1")
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

	// Make sure that each page.KindTaxonomyTerm page has an appropriate number
	// of page.KindTaxonomy pages in its Pages slice.
	taxonomyTermPageCounts := map[string]int{
		"tags":         3,
		"categories":   2,
		"others":       2,
		"empties":      0,
		"permalinkeds": 1,
	}

	for taxonomy, count := range taxonomyTermPageCounts {
		term := s.getPage(page.KindTaxonomyTerm, taxonomy)
		b.Assert(term, qt.Not(qt.IsNil))
		b.Assert(len(term.Pages()), qt.Equals, count, qt.Commentf(taxonomy))

		for _, p := range term.Pages() {
			b.Assert(p.Kind(), qt.Equals, page.KindTaxonomy)
		}
	}

	cat1 := s.getPage(page.KindTaxonomy, "categories", "cat1")
	b.Assert(cat1, qt.Not(qt.IsNil))
	if uglyURLs {
		b.Assert(cat1.RelPermalink(), qt.Equals, "/blog/categories/cat1.html")
	} else {
		b.Assert(cat1.RelPermalink(), qt.Equals, "/blog/categories/cat1/")
	}

	pl1 := s.getPage(page.KindTaxonomy, "permalinkeds", "pl1")
	permalinkeds := s.getPage(page.KindTaxonomyTerm, "permalinkeds")
	b.Assert(pl1, qt.Not(qt.IsNil))
	b.Assert(permalinkeds, qt.Not(qt.IsNil))
	if uglyURLs {
		b.Assert(pl1.RelPermalink(), qt.Equals, "/blog/perma/pl1.html")
		b.Assert(permalinkeds.RelPermalink(), qt.Equals, "/blog/permalinkeds.html")
	} else {
		b.Assert(pl1.RelPermalink(), qt.Equals, "/blog/perma/pl1/")
		b.Assert(permalinkeds.RelPermalink(), qt.Equals, "/blog/permalinkeds/")
	}

	helloWorld := s.getPage(page.KindTaxonomy, "others", "hello-hugo-world")
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

	ta := s.findPagesByKind(page.KindTaxonomy)
	te := s.findPagesByKind(page.KindTaxonomyTerm)

	b.Assert(len(te), qt.Equals, 4)
	b.Assert(len(ta), qt.Equals, 7)

	b.AssertFileContent("public/news/categories/a/index.html", "Taxonomy List Page 1|a|Hello|https://example.com/news/categories/a/|")
	b.AssertFileContent("public/news/categories/b/index.html", "Taxonomy List Page 1|This is B|Hello|https://example.com/news/categories/b/|")
	b.AssertFileContent("public/news/categories/d/e/index.html", "Taxonomy List Page 1|d/e|Hello|https://example.com/news/categories/d/e/|")
	b.AssertFileContent("public/news/categories/f/g/h/index.html", "Taxonomy List Page 1|This is H|Hello|https://example.com/news/categories/f/g/h/|")
	b.AssertFileContent("public/t1/t2/t3s/t4/t5/index.html", "Taxonomy List Page 1|This is T5|Hello|https://example.com/t1/t2/t3s/t4/t5/|")
	b.AssertFileContent("public/t1/t2/t3s/t4/t5/t6/index.html", "Taxonomy List Page 1|t4/t5/t6|Hello|https://example.com/t1/t2/t3s/t4/t5/t6/|")

	b.AssertFileContent("public/news/categories/index.html", "Taxonomy Term Page 1|News/Categories|Hello|https://example.com/news/categories/|")
	b.AssertFileContent("public/t1/t2/t3s/index.html", "Taxonomy Term Page 1|T1/T2/T3s|Hello|https://example.com/t1/t2/t3s/|")
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
	b.AssertFileContent("public/categories/index.html", `<li><a href="http://example.com/categories/this-is-cool/">This is Cool</a> 10</li>`)
	b.AssertFileContent("public/tags/index.html", `<li><a href="http://example.com/tags/rocks-i-say/">Rocks I say!</a> 10</li>`)

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

	reg, _ := s.getPageNew(nil, "categories/regular")
	dra, _ := s.getPageNew(nil, "categories/draft")
	b.Assert(reg, qt.Not(qt.IsNil))
	b.Assert(dra, qt.IsNil)

}

// See https://github.com/gohugoio/hugo/issues/6222
// We need to revisit this once we figure out what to do with the
// draft etc _index pages, but for now we need to avoid the crash.
func TestTaxonomiesIndexDraft(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithContent(
		"categories/_index.md", `---
title: "The Categories"
draft: true
---

This is the invisible content.

`)

	b.WithTemplates("index.html", `
{{ range .Site.Pages }}
{{ .RelPermalink }}|{{ .Title }}|{{ .WordCount }}|{{ .Content }}|
{{ end }}
`)

	b.Build(BuildCfg{})

	// We publish the index page, but the content will be empty.
	b.AssertFileContent("public/index.html", " /categories/|The Categories|0||")

}
