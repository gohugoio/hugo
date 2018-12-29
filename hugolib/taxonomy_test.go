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

	expect := []string{"a", "b", "c", "x/y"}
	if !reflect.DeepEqual(st, expect) {
		t.Fatalf("ordered taxonomies do not match %v.  Got: %s", expect, st)
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
subcats = "subcats"

[permalinks]
permalinkeds = "/perma/:slug/"
subcats = "/subcats/:slug/"
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
subcats:
%s
---
# Doc
`

	siteConfig = fmt.Sprintf(siteConfig, preserveTaxonomyNames, uglyURLs)

	th, h := newTestSitesFromConfigWithDefaultTemplates(t, siteConfig)
	require.Len(t, h.Sites, 1)

	fs := th.Fs

	writeSource(t, fs, "content/p1.md", fmt.Sprintf(pageTemplate, "t1/c1", "- Tag1", "- cat1\n- \"cAt/dOg\"", "- o1", "- pl1", ""))
	writeSource(t, fs, "content/p2.md", fmt.Sprintf(pageTemplate, "t2/c1", "- tag2", "- cat1", "- o1", "- pl1", ""))
	writeSource(t, fs, "content/p3.md", fmt.Sprintf(pageTemplate, "t2/c12", "- tag2", "- cat2", "- o1", "- pl1", ""))
	writeSource(t, fs, "content/p4.md", fmt.Sprintf(pageTemplate, "Hello World", "", "", "- \"Hello Hugo world\"", "- pl1", ""))
	writeSource(t, fs, "content/p5.md", fmt.Sprintf(pageTemplate, "Sub/categories", "", "", "", "", "- \"sc0/sp1\""))

	writeNewContentFile(t, fs.Source, "Category Terms", "2017-01-01", "content/categories/_index.md", 10)
	writeNewContentFile(t, fs.Source, "Tag1 List", "2017-01-01", "content/tags/Tag1/_index.md", 10)

	require.NoError(t, h.Build(BuildCfg{}))

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
		th.assertFileContent(pathFunc("public/categories/cat-dog/index.html"), "List", "cAt/dOg")
	} else {
		th.assertFileContent(pathFunc("public/categories/cat1/index.html"), "List", "Cat1")
		th.assertFileContent(pathFunc("public/categories/cat-dog/index.html"), "List", "Cat/Dog")
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
		"categories":   3,
		"others":       2,
		"empties":      0,
		"permalinkeds": 1,
		"subcats":      1,
	}

	for taxonomy, count := range taxonomyTermPageCounts {
		term := s.getPage(KindTaxonomyTerm, taxonomy)
		require.NotNil(t, term)
		require.Len(t, term.Pages, count, taxonomy)

		for _, page := range term.Pages {
			require.Equal(t, KindTaxonomy, page.Kind)
		}
	}

	fixTerm := func(s string) string {
		if preserveTaxonomyNames {
			return s
		}
		return strings.ToLower(s)
	}

	fixURL := func(s string) string {
		if uglyURLs {
			return strings.TrimRight(s, "/") + ".html"
		}
		return s
	}

	cat1 := s.getPage(KindTaxonomy, "categories", "cat1")
	require.NotNil(t, cat1)
	require.Equal(t, fixURL("/blog/categories/cat1/"), cat1.RelPermalink())

	catdog := s.getPage(KindTaxonomy, "categories", fixTerm("cAt/dOg"))
	require.NotNil(t, catdog)
	require.Equal(t, fixURL("/blog/categories/cat-dog/"), catdog.RelPermalink())

	pl1 := s.getPage(KindTaxonomy, "permalinkeds", "pl1")
	require.NotNil(t, pl1)
	require.Equal(t, fixURL("/blog/perma/pl1/"), pl1.RelPermalink())

	permalinkeds := s.getPage(KindTaxonomyTerm, "permalinkeds")
	require.NotNil(t, permalinkeds)
	require.Equal(t, fixURL("/blog/permalinkeds/"), permalinkeds.RelPermalink())

	// Issue #5223
	sp1 := s.getPage(KindTaxonomy, "subcats", "sc0/sp1")
	require.NotNil(t, sp1)
	require.Equal(t, fixURL("/blog/subcats/sc0/sp1/"), sp1.RelPermalink())

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
