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
	"strings"
	"testing"

	"fmt"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/helpers"
)

func TestDisableKindsNoneDisabled(t *testing.T) {
	t.Parallel()
	doTestDisableKinds(t)
}

func TestDisableKindsSomeDisabled(t *testing.T) {
	t.Parallel()
	doTestDisableKinds(t, page.KindSection, kind404)
}

func TestDisableKindsOneDisabled(t *testing.T) {
	t.Parallel()
	for _, kind := range allKinds {
		if kind == page.KindPage {
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

	disabledStr := "[]"

	if len(disabled) > 0 {
		disabledStr = strings.Replace(fmt.Sprintf("%#v", disabled), "[]string{", "[", -1)
		disabledStr = strings.Replace(disabledStr, "}", "]", -1)
	}

	siteConfig := fmt.Sprintf(siteConfigTemplate, disabledStr)

	b := newTestSitesBuilder(t).WithConfigFile("toml", siteConfig)

	b.WithTemplates(
		"index.html", "Home|{{ .Title }}|{{ .Content }}",
		"_default/single.html", "Single|{{ .Title }}|{{ .Content }}",
		"_default/list.html", "List|{{ .Title }}|{{ .Content }}",
		"_default/terms.html", "Terms List|{{ .Title }}|{{ .Content }}",
		"layouts/404.html", "Page Not Found",
	)

	b.WithContent(
		"sect/p1.md", fmt.Sprintf(pageTemplate, "P1", "- tag1"),
		"categories/_index.md", newTestPage("Category Terms", "2017-01-01", 10),
		"tags/tag1/_index.md", newTestPage("Tag1 List", "2017-01-01", 10),
	)

	b.Build(BuildCfg{})
	h := b.H

	assertDisabledKinds(b, h.Sites[0], disabled...)

}

func assertDisabledKinds(b *sitesBuilder, s *Site, disabled ...string) {
	assertDisabledKind(b,
		func(isDisabled bool) bool {
			if isDisabled {
				return len(s.RegularPages()) == 0
			}
			return len(s.RegularPages()) > 0
		}, disabled, page.KindPage, "public/sect/p1/index.html", "Single|P1")
	assertDisabledKind(b,
		func(isDisabled bool) bool {
			p := s.getPage(page.KindHome)
			if isDisabled {
				return p == nil
			}
			return p != nil
		}, disabled, page.KindHome, "public/index.html", "Home")
	assertDisabledKind(b,
		func(isDisabled bool) bool {
			p := s.getPage(page.KindSection, "sect")
			if isDisabled {
				return p == nil
			}
			return p != nil
		}, disabled, page.KindSection, "public/sect/index.html", "Sects")
	assertDisabledKind(b,
		func(isDisabled bool) bool {
			p := s.getPage(page.KindTaxonomy, "tags", "tag1")

			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, page.KindTaxonomy, "public/tags/tag1/index.html", "Tag1")
	assertDisabledKind(b,
		func(isDisabled bool) bool {
			p := s.getPage(page.KindTaxonomyTerm, "tags")
			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, page.KindTaxonomyTerm, "public/tags/index.html", "Tags")
	assertDisabledKind(b,
		func(isDisabled bool) bool {
			p := s.getPage(page.KindTaxonomyTerm, "categories")

			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, page.KindTaxonomyTerm, "public/categories/index.html", "Category Terms")
	assertDisabledKind(b,
		func(isDisabled bool) bool {
			p := s.getPage(page.KindTaxonomy, "categories", "hugo")
			if isDisabled {
				return p == nil
			}
			return p != nil

		}, disabled, page.KindTaxonomy, "public/categories/hugo/index.html", "Hugo")
	// The below have no page in any collection.
	assertDisabledKind(b, func(isDisabled bool) bool { return true }, disabled, kindRSS, "public/index.xml", "<link>")
	assertDisabledKind(b, func(isDisabled bool) bool { return true }, disabled, kindSitemap, "public/sitemap.xml", "sitemap")
	assertDisabledKind(b, func(isDisabled bool) bool { return true }, disabled, kindRobotsTXT, "public/robots.txt", "User-agent")
	assertDisabledKind(b, func(isDisabled bool) bool { return true }, disabled, kind404, "public/404.html", "Page Not Found")
}

func assertDisabledKind(b *sitesBuilder, kindAssert func(bool) bool, disabled []string, kind, path, matcher string) {
	isDisabled := stringSliceContains(kind, disabled...)
	b.Assert(kindAssert(isDisabled), qt.Equals, true)

	if kind == kindRSS && !isDisabled {
		// If the home page is also disabled, there is not RSS to look for.
		if stringSliceContains(page.KindHome, disabled...) {
			isDisabled = true
		}
	}

	if isDisabled {
		// Path should not exist
		fileExists, err := helpers.Exists(path, b.Fs.Destination)
		b.Assert(err, qt.IsNil)
		b.Assert(fileExists, qt.Equals, false)

	} else {
		b.AssertFileContent(path, matcher)
	}
}
