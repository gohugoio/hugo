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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spf13/hugo/deps"
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

func TestTaxonomiesWithAndWithoutContentFile(t *testing.T) {
	for _, preserveTaxonomyNames := range []bool{false, true} {
		t.Run(fmt.Sprintf("preserveTaxonomyNames %t", preserveTaxonomyNames), func(t *testing.T) {
			doTestTaxonomiesWithAndWithoutContentFile(t, preserveTaxonomyNames)
		})

	}
}

func doTestTaxonomiesWithAndWithoutContentFile(t *testing.T, preserveTaxonomyNames bool) {
	t.Parallel()

	siteConfig := `
baseURL = "http://example.com/blog"
preserveTaxonomyNames = %t

paginate = 1
defaultContentLanguage = "en"

[Taxonomies]
tag = "tags"
category = "categories"
other = "others"
empty = "empties"
`

	pageTemplate := `---
title: "%s"
tags:
%s
categories:
%s
others:
%s
---
# Doc
`

	siteConfig = fmt.Sprintf(siteConfig, preserveTaxonomyNames)

	th, h := newTestSitesFromConfigWithDefaultTemplates(t, siteConfig)
	require.Len(t, h.Sites, 1)

	fs := th.Fs

	writeSource(t, fs, "content/p1.md", fmt.Sprintf(pageTemplate, "t1/c1", "- tag1", "- cat1", "- o1"))
	writeSource(t, fs, "content/p2.md", fmt.Sprintf(pageTemplate, "t2/c1", "- tag2", "- cat1", "- o1"))
	writeSource(t, fs, "content/p3.md", fmt.Sprintf(pageTemplate, "t2/c12", "- tag2", "- cat2", "- o1"))
	writeSource(t, fs, "content/p4.md", fmt.Sprintf(pageTemplate, "Hello World", "", "", "- \"Hello Hugo world\""))

	writeNewContentFile(t, fs, "Category Terms", "2017-01-01", "content/categories/_index.md", 10)
	writeNewContentFile(t, fs, "Tag1 List", "2017-01-01", "content/tags/tag1/_index.md", 10)

	err := h.Build(BuildCfg{})

	require.NoError(t, err)

	// So what we have now is:
	// 1. categories with terms content page, but no content page for the only c1 category
	// 2. tags with no terms content page, but content page for one of 2 tags (tag1)
	// 3. the "others" taxonomy with no content pages.

	// 1.
	th.assertFileContent("public/categories/cat1/index.html", "List", "Cat1")
	th.assertFileContent("public/categories/index.html", "Terms List", "Category Terms")

	// 2.
	th.assertFileContent("public/tags/tag2/index.html", "List", "Tag2")
	th.assertFileContent("public/tags/tag1/index.html", "List", "Tag1")
	th.assertFileContent("public/tags/index.html", "Terms List", "Tags")

	// 3.
	th.assertFileContent("public/others/o1/index.html", "List", "O1")
	th.assertFileContent("public/others/index.html", "Terms List", "Others")

	s := h.Sites[0]

	// Make sure that each KindTaxonomyTerm page has an appropriate number
	// of KindTaxonomy pages in its Pages slice.
	taxonomyTermPageCounts := map[string]int{
		"tags":       2,
		"categories": 2,
		"others":     2,
		"empties":    0,
	}

	for taxonomy, count := range taxonomyTermPageCounts {
		term := s.getPage(KindTaxonomyTerm, taxonomy)
		require.NotNil(t, term)
		require.Len(t, term.Pages, count)

		for _, page := range term.Pages {
			require.Equal(t, KindTaxonomy, page.Kind)
		}
	}

	// Issue #3070 preserveTaxonomyNames
	if preserveTaxonomyNames {
		helloWorld := s.getPage(KindTaxonomy, "others", "Hello Hugo world")
		require.NotNil(t, helloWorld)
		require.Equal(t, "Hello Hugo world", helloWorld.Title)
	} else {
		helloWorld := s.getPage(KindTaxonomy, "others", "hello-hugo-world")
		require.NotNil(t, helloWorld)
		require.Equal(t, "Hello Hugo World", helloWorld.Title)
	}

	// Issue #2977
	th.assertFileContent("public/empties/index.html", "Terms List", "Empties")

}
