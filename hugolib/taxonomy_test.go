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

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/hugofs"
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

// Issue #2992
func TestTaxonomiesWithAndWithoutContentFile(t *testing.T) {
	t.Parallel()

	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 1
defaultContentLanguage = "en"

[Taxonomies]
tag = "tags"
category = "categories"
other = "others"
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

	mf := afero.NewMemMapFs()

	writeToFs(t, mf, "config.toml", siteConfig)

	cfg, err := LoadConfig(mf, "", "config.toml")
	require.NoError(t, err)

	fs := hugofs.NewFrom(mf, cfg)
	th := testHelper{cfg, fs, t}

	writeSource(t, fs, "layouts/_default/single.html", "Single|{{ .Title }}|{{ .Content }}")
	writeSource(t, fs, "layouts/_default/list.html", "List|{{ .Title }}|{{ .Content }}")
	writeSource(t, fs, "layouts/_default/terms.html", "Terms List|{{ .Title }}|{{ .Content }}")

	writeSource(t, fs, "content/p1.md", fmt.Sprintf(pageTemplate, "t1/c1", "- tag1", "- cat1", "- o1"))
	writeSource(t, fs, "content/p2.md", fmt.Sprintf(pageTemplate, "t2/c1", "- tag2", "- cat1", "- o1"))

	writeNewContentFile(t, fs, "Category Terms", "2017-01-01", "content/categories/_index.md", 10)
	writeNewContentFile(t, fs, "Tag1 List", "2017-01-01", "content/tags/tag1/_index.md", 10)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)
	require.Len(t, h.Sites, 1)

	err = h.Build(BuildCfg{})

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

}
