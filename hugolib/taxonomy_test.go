// Copyright 2015 The Hugo Authors. All rights reserved.
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

	"github.com/stretchr/testify/require"

	"github.com/gohugoio/hugo/deps"
)

func TestByCountOrderOfTaxonomies(t *testing.T) {
	t.Parallel()
	taxonomies := make(map[string]string)

	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	cfg, fs := newTestCfg()

	cfg.Set("taxonomies", taxonomies)

	writeSource(t, fs, filepath.Join("content", "page.md"), pageYamlWithTaxonomiesA)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	st := make([]string, 0)
	for _, t := range s.Taxonomies["tags"].ByCount() {
		st = append(st, t.Name)
	}

	if !reflect.DeepEqual(st, []string{"a", "b", "c"}) {
		t.Fatalf("ordered taxonomies do not match [a, b, c].  Got: %s", st)
	}
}

//
func TestTaxonomiesWithAndWithoutContentFile(t *testing.T) {
	for _, uglyURLs := range []bool{false, true} {
		for _, preserveTaxonomyNames := range []bool{false, true} {
			t.Run(fmt.Sprintf("uglyURLs=%t,preserveTaxonomyNames=%t", uglyURLs, preserveTaxonomyNames), func(t *testing.T) {
				doTestTaxonomiesWithAndWithoutContentFile(t, preserveTaxonomyNames, uglyURLs)
			})
		}
	}
}

func doTestTaxonomiesWithAndWithoutContentFile(t *testing.T, preserveTaxonomyNames, uglyURLs bool) {
	t.Parallel()

	siteConfig := `
baseURL = "http://example.com/blog"
preserveTaxonomyNames = %t
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

	siteConfig = fmt.Sprintf(siteConfig, preserveTaxonomyNames, uglyURLs)

	th, h := newTestSitesFromConfigWithDefaultTemplates(t, siteConfig)
	require.Len(t, h.Sites, 1)

	fs := th.Fs

	if preserveTaxonomyNames {
		writeSource(t, fs, "content/p1.md", fmt.Sprintf(pageTemplate, "t1/c1", "- tag1", "- cat1", "- o1", "- pl1"))
	} else {
		// Check lower-casing of tags
		writeSource(t, fs, "content/p1.md", fmt.Sprintf(pageTemplate, "t1/c1", "- Tag1", "- cAt1", "- o1", "- pl1"))

	}
	writeSource(t, fs, "content/p2.md", fmt.Sprintf(pageTemplate, "t2/c1", "- tag2", "- cat1", "- o1", "- pl1"))
	writeSource(t, fs, "content/p3.md", fmt.Sprintf(pageTemplate, "t2/c12", "- tag2", "- cat2", "- o1", "- pl1"))
	writeSource(t, fs, "content/p4.md", fmt.Sprintf(pageTemplate, "Hello World", "", "", "- \"Hello Hugo world\"", "- pl1"))

	writeNewContentFile(t, fs.Source, "Category Terms", "2017-01-01", "content/categories/_index.md", 10)
	writeNewContentFile(t, fs.Source, "Tag1 List", "2017-01-01", "content/tags/Tag1/_index.md", 10)

	err := h.Build(BuildCfg{})

	require.NoError(t, err)

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
	if preserveTaxonomyNames {
		th.assertFileContent(pathFunc("public/categories/cat1/index.html"), "List", "cat1")
	} else {
		th.assertFileContent(pathFunc("public/categories/cat1/index.html"), "List", "Cat1")
	}

	th.assertFileContent(pathFunc("public/categories/index.html"), "Terms List", "Category Terms")

	// 2.
	if preserveTaxonomyNames {
		th.assertFileContent(pathFunc("public/tags/tag2/index.html"), "List", "tag2")
	} else {
		th.assertFileContent(pathFunc("public/tags/tag2/index.html"), "List", "Tag2")
	}
	th.assertFileContent(pathFunc("public/tags/tag1/index.html"), "List", "Tag1")
	th.assertFileContent(pathFunc("public/tags/index.html"), "Terms List", "Tags")

	// 3.
	if preserveTaxonomyNames {
		th.assertFileContent(pathFunc("public/others/o1/index.html"), "List", "o1")
	} else {
		th.assertFileContent(pathFunc("public/others/o1/index.html"), "List", "O1")
	}
	th.assertFileContent(pathFunc("public/others/index.html"), "Terms List", "Others")

	// 4.
	if preserveTaxonomyNames {
		th.assertFileContent(pathFunc("public/perma/pl1/index.html"), "List", "pl1")
	} else {
		th.assertFileContent(pathFunc("public/perma/pl1/index.html"), "List", "Pl1")
	}
	// This looks kind of funky, but the taxonomy terms do not have a permalinks definition,
	// for good reasons.
	th.assertFileContent(pathFunc("public/permalinkeds/index.html"), "Terms List", "Permalinkeds")

	s := h.Sites[0]

	// Make sure that each KindTaxonomyTerm page has an appropriate number
	// of KindTaxonomy pages in its Pages slice.
	taxonomyTermPageCounts := map[string]int{
		"tags":         2,
		"categories":   2,
		"others":       2,
		"empties":      0,
		"permalinkeds": 1,
	}

	for taxonomy, count := range taxonomyTermPageCounts {
		term := s.getPage(KindTaxonomyTerm, taxonomy)
		require.NotNil(t, term)
		require.Len(t, term.Pages, count)

		for _, page := range term.Pages {
			require.Equal(t, KindTaxonomy, page.Kind)
		}
	}

	cat1 := s.getPage(KindTaxonomy, "categories", "cat1")
	require.NotNil(t, cat1)
	if uglyURLs {
		require.Equal(t, "/blog/categories/cat1.html", cat1.RelPermalink())
	} else {
		require.Equal(t, "/blog/categories/cat1/", cat1.RelPermalink())
	}

	pl1 := s.getPage(KindTaxonomy, "permalinkeds", "pl1")
	permalinkeds := s.getPage(KindTaxonomyTerm, "permalinkeds")
	require.NotNil(t, pl1)
	require.NotNil(t, permalinkeds)
	if uglyURLs {
		require.Equal(t, "/blog/perma/pl1.html", pl1.RelPermalink())
		require.Equal(t, "/blog/permalinkeds.html", permalinkeds.RelPermalink())
	} else {
		require.Equal(t, "/blog/perma/pl1/", pl1.RelPermalink())
		require.Equal(t, "/blog/permalinkeds/", permalinkeds.RelPermalink())
	}

	// Issue #3070 preserveTaxonomyNames
	if preserveTaxonomyNames {
		helloWorld := s.getPage(KindTaxonomy, "others", "Hello Hugo world")
		require.NotNil(t, helloWorld)
		require.Equal(t, "Hello Hugo world", helloWorld.title)
	} else {
		helloWorld := s.getPage(KindTaxonomy, "others", "hello-hugo-world")
		require.NotNil(t, helloWorld)
		require.Equal(t, "Hello Hugo World", helloWorld.title)
	}

	// Issue #2977
	th.assertFileContent(pathFunc("public/empties/index.html"), "Terms List", "Empties")

}

// https://github.com/gohugoio/hugo/issues/5513
// https://github.com/gohugoio/hugo/issues/5571
func TestTaxonomiesPathSeparation(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

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

	ta := s.findPagesByKind(KindTaxonomy)
	te := s.findPagesByKind(KindTaxonomyTerm)

	assert.Equal(4, len(te))
	assert.Equal(7, len(ta))

	b.AssertFileContent("public/news/categories/a/index.html", "Taxonomy List Page 1|A|Hello|https://example.com/news/categories/a/|")
	b.AssertFileContent("public/news/categories/b/index.html", "Taxonomy List Page 1|This is B|Hello|https://example.com/news/categories/b/|")
	b.AssertFileContent("public/news/categories/d/e/index.html", "Taxonomy List Page 1|D/E|Hello|https://example.com/news/categories/d/e/|")
	b.AssertFileContent("public/news/categories/f/g/h/index.html", "Taxonomy List Page 1|This is H|Hello|https://example.com/news/categories/f/g/h/|")
	b.AssertFileContent("public/t1/t2/t3s/t4/t5/index.html", "Taxonomy List Page 1|This is T5|Hello|https://example.com/t1/t2/t3s/t4/t5/|")
	b.AssertFileContent("public/t1/t2/t3s/t4/t5/t6/index.html", "Taxonomy List Page 1|T4/T5/T6|Hello|https://example.com/t1/t2/t3s/t4/t5/t6/|")

	b.AssertFileContent("public/news/categories/index.html", "Taxonomy Term Page 1|News/Categories|Hello|https://example.com/news/categories/|")
	b.AssertFileContent("public/t1/t2/t3s/index.html", "Taxonomy Term Page 1|T1/T2/T3s|Hello|https://example.com/t1/t2/t3s/|")
	b.AssertFileContent("public/s1/s2/s3s/index.html", "Taxonomy Term Page 1|This is S3s|Hello|https://example.com/s1/s2/s3s/|")

}
