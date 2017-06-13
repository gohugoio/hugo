// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"

	"fmt"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/stretchr/testify/require"
)

func TestDisableKindsNoneDisabled(t *testing.T) {
	t.Parallel()
	doTestDisableKinds(t)
}

func TestDisableKindsSomeDisabled(t *testing.T) {
	t.Parallel()
	doTestDisableKinds(t, KindSection, kind404)
}

func TestDisableKindsOneDisabled(t *testing.T) {
	t.Parallel()
	for _, kind := range allKinds {
		if kind == KindPage {
			// Turning off regular page generation have some side-effects
			// not handled by the assertions below (no sections), so
			// skip that for now.
			continue
		}
		doTestDisableKinds(t, kind)
	}
}

func TestDisableKindsAllDisabled(t *testing.T) {
	t.Parallel()
	doTestDisableKinds(t, allKinds...)
}

func doTestDisableKinds(t *testing.T, disabled ...string) {
	siteConfigTemplate := `
baseURL = "http://example.com/blog"
enableRobotsTXT = true
disableKinds = %s

paginate = 1
defaultContentLanguage = "en"

[Taxonomies]
tag = "tags"
category = "categories"
`

	pageTemplate := `---
title: "%s"
tags:
%s
categories:
- Hugo
---
# Doc
`

	mf := afero.NewMemMapFs()

	disabledStr := "[]"

	if len(disabled) > 0 {
		disabledStr = strings.Replace(fmt.Sprintf("%#v", disabled), "[]string{", "[", -1)
		disabledStr = strings.Replace(disabledStr, "}", "]", -1)
	}

	siteConfig := fmt.Sprintf(siteConfigTemplate, disabledStr)
	writeToFs(t, mf, "config.toml", siteConfig)

	cfg, err := LoadConfig(mf, "", "config.toml")
	require.NoError(t, err)

	fs := hugofs.NewFrom(mf, cfg)
	th := testHelper{cfg, fs, t}

	writeSource(t, fs, "layouts/index.html", "Home|{{ .Title }}|{{ .Content }}")
	writeSource(t, fs, "layouts/_default/single.html", "Single|{{ .Title }}|{{ .Content }}")
	writeSource(t, fs, "layouts/_default/list.html", "List|{{ .Title }}|{{ .Content }}")
	writeSource(t, fs, "layouts/_default/terms.html", "Terms List|{{ .Title }}|{{ .Content }}")
	writeSource(t, fs, "layouts/404.html", "Page Not Found")

	writeSource(t, fs, "content/sect/p1.md", fmt.Sprintf(pageTemplate, "P1", "- tag1"))

	writeNewContentFile(t, fs, "Category Terms", "2017-01-01", "content/categories/_index.md", 10)
	writeNewContentFile(t, fs, "Tag1 List", "2017-01-01", "content/tags/tag1/_index.md", 10)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)
	require.Len(t, h.Sites, 1)

	err = h.Build(BuildCfg{})

	require.NoError(t, err)

	assertDisabledKinds(th, h.Sites[0], disabled...)

}

func assertDisabledKinds(th testHelper, s *Site, disabled ...string) {
	assertDisabledKind(th,
		func(isDisabled bool) bool {
			if isDisabled {
				return len(s.RegularPages) == 0
			}
			return len(s.RegularPages) > 0
		}, disabled, KindPage, "public/sect/p1/index.html", "Single|P1")
	assertDisabledKind(th,
		func(isDisabled bool) bool {
			p := s.getPage(KindHome)
			if isDisabled {
				return p == nil
			}
			return p != nil
		}, disabled, KindHome, "public/index.html", "Home")
	assertDisabledKind(th,
		func(isDisabled bool) bool {
			p := s.getPage(KindSection, "sect")
			if isDisabled {
				return p == nil
			}
			return p != nil
		}, disabled, KindSection, "public/sect/index.html", "Sects")
	assertDisabledKind(th,
		func(isDisabled bool) bool {
			p := s.getPage(KindTaxonomy, "tags", "tag1")

			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, KindTaxonomy, "public/tags/tag1/index.html", "Tag1")
	assertDisabledKind(th,
		func(isDisabled bool) bool {
			p := s.getPage(KindTaxonomyTerm, "tags")
			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, KindTaxonomyTerm, "public/tags/index.html", "Tags")
	assertDisabledKind(th,
		func(isDisabled bool) bool {
			p := s.getPage(KindTaxonomyTerm, "categories")

			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, KindTaxonomyTerm, "public/categories/index.html", "Category Terms")
	assertDisabledKind(th,
		func(isDisabled bool) bool {
			p := s.getPage(KindTaxonomy, "categories", "hugo")
			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, KindTaxonomy, "public/categories/hugo/index.html", "Hugo")
	// The below have no page in any collection.
	assertDisabledKind(th, func(isDisabled bool) bool { return true }, disabled, kindRSS, "public/index.xml", "<link>")
	assertDisabledKind(th, func(isDisabled bool) bool { return true }, disabled, kindSitemap, "public/sitemap.xml", "sitemap")
	assertDisabledKind(th, func(isDisabled bool) bool { return true }, disabled, kindRobotsTXT, "public/robots.txt", "User-agent")
	assertDisabledKind(th, func(isDisabled bool) bool { return true }, disabled, kind404, "public/404.html", "Page Not Found")
}

func assertDisabledKind(th testHelper, kindAssert func(bool) bool, disabled []string, kind, path, matcher string) {
	isDisabled := stringSliceContains(kind, disabled...)
	require.True(th.T, kindAssert(isDisabled), fmt.Sprintf("%s: %t", kind, isDisabled))

	if kind == kindRSS && !isDisabled {
		// If the home page is also disabled, there is not RSS to look for.
		if stringSliceContains(KindHome, disabled...) {
			isDisabled = true
		}
	}

	if isDisabled {
		// Path should not exist
		fileExists, err := helpers.Exists(path, th.Fs.Destination)
		require.False(th.T, fileExists)
		require.NoError(th.T, err)

	} else {
		th.assertFileContent(path, matcher)
	}
}

func stringSliceContains(k string, values ...string) bool {
	for _, v := range values {
		if k == v {
			return true
		}
	}
	return false
}
