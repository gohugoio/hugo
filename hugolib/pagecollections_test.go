// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

type getPageTest struct {
	name          string
	kind          string
	context       page.Page
	pathVariants  []string
	expectedTitle string
}

func (t *getPageTest) check(p page.Page, err error, errorMsg string, c *qt.C) {
	c.Helper()
	errorComment := qt.Commentf(errorMsg)
	switch t.kind {
	case "Ambiguous":
		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(p, qt.IsNil, errorComment)
	case "NoPage":
		c.Assert(err, qt.IsNil)
		c.Assert(p, qt.IsNil, errorComment)
	default:
		c.Assert(err, qt.IsNil, errorComment)
		c.Assert(p, qt.Not(qt.IsNil), errorComment)
		c.Assert(p.Kind(), qt.Equals, t.kind, errorComment)
		c.Assert(p.Title(), qt.Equals, t.expectedTitle, errorComment)
	}
}

func TestGetPage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/_default/single.html --
{{ .Title }}
-- layouts/_default/list.html --
{{ .Title }}

-- content/_index.md --
---
title: "home page"
categories:
- Hugo
---
# Doc

-- content/about.md --
---
title: "about page"
categories:
- Hugo
---
# Doc

-- content/sect0/page1.md --
---
title: "Title0_1"
categories:
- Hugo
---
# Doc

-- content/sect1/page1.md --
---
title: "Title1_1"
categories:
- Hugo
---
# Doc

-- content/sect2/_index.md --
---
title: "Section 2"
categories:
- Hugo
---
# Doc

-- content/sect3/_index.md --
---
title: "section 3"
categories:
- Hugo
---
# Doc

-- content/sect3/page1.md --
---
title: "Title3_1"
categories:
- Hugo
---
# Doc

-- content/sect3/unique.md --
---
title: "UniqueBase"
categories:
- Hugo
---
# Doc

-- content/sect3/Unique2.md --
---
title: "UniqueBase2"
categories:
- Hugo
---
# Doc

-- content/sect3/sect7/_index.md --
---
title: "another sect7"
categories:
- Hugo
---
# Doc

-- content/sect3/subsect/deep.md --
---
title: "deep page"
categories:
- Hugo
---
# Doc

-- content/sect3/b1/index.md --
---
title: "b1 bundle"
categories:
- Hugo
---
# Doc

-- content/sect3/index/index.md --
---
title: "index bundle"
categories:
- Hugo
---
# Doc

-- content/sect4/page2.md --
---
title: "Title4_2"
categories:
- Hugo
---
# Doc

-- content/sect5/page3.md --
---
title: "Title5_3"
categories:
- Hugo
---
# Doc

-- content/sect7/somesubpage.md --
---
title: "Some Sub Page in Sect7"
categories:
- Hugo
---
# Doc

-- content/sect7/page9.md --
---
title: "Title7_9"
categories:
- Hugo
---
# Doc

-- content/section_bundle_overlap/_index.md --
---
title: "index overlap section"
categories:
- Hugo
---
# Doc

-- content/section_bundle_overlap_bundle/index.md --
---
title: "index overlap bundle"
categories:
- Hugo
---
# Doc
`

	b := Test(t, files)

	s := b.H.Sites[0]

	sec3, err := s.getPage(nil, "/sect3")
	b.Assert(err, qt.IsNil)
	b.Assert(sec3, qt.Not(qt.IsNil))

	tests := []getPageTest{
		// legacy content root relative paths
		{"Root relative, no slash, home", kinds.KindHome, nil, []string{""}, "home page"},
		{"Root relative, no slash, root page", kinds.KindPage, nil, []string{"about.md", "ABOUT.md"}, "about page"},
		{"Root relative, no slash, section", kinds.KindSection, nil, []string{"sect3"}, "section 3"},
		{"Root relative, no slash, section page", kinds.KindPage, nil, []string{"sect3/page1.md"}, "Title3_1"},
		{"Root relative, no slash, sub section", kinds.KindSection, nil, []string{"sect3/sect7"}, "another sect7"},
		{"Root relative, no slash, nested page", kinds.KindPage, nil, []string{"sect3/subsect/deep.md"}, "deep page"},
		{"Root relative, no slash, OS slashes", kinds.KindPage, nil, []string{filepath.FromSlash("sect5/page3.md")}, "Title5_3"},

		{"Short ref, unique", kinds.KindPage, nil, []string{"unique.md", "unique"}, "UniqueBase"},
		{"Short ref, unique, upper case", kinds.KindPage, nil, []string{"Unique2.md", "unique2.md", "unique2"}, "UniqueBase2"},
		{"Short ref, ambiguous", "Ambiguous", nil, []string{"page1.md"}, ""},

		// ISSUE: This is an ambiguous ref, but because we have to support the legacy
		// content root relative paths without a leading slash, the lookup
		// returns /sect7. This undermines ambiguity detection, but we have no choice.
		//{"Ambiguous", nil, []string{"sect7"}, ""},
		{"Section, ambiguous", kinds.KindSection, nil, []string{"sect7"}, "Sect7s"},

		{"Absolute, home", kinds.KindHome, nil, []string{"/", ""}, "home page"},
		{"Absolute, page", kinds.KindPage, nil, []string{"/about.md", "/about"}, "about page"},
		{"Absolute, sect", kinds.KindSection, nil, []string{"/sect3"}, "section 3"},
		{"Absolute, page in subsection", kinds.KindPage, nil, []string{"/sect3/page1.md", "/Sect3/Page1.md"}, "Title3_1"},
		{"Absolute, section, subsection with same name", kinds.KindSection, nil, []string{"/sect3/sect7"}, "another sect7"},
		{"Absolute, page, deep", kinds.KindPage, nil, []string{"/sect3/subsect/deep.md"}, "deep page"},
		{"Absolute, page, OS slashes", kinds.KindPage, nil, []string{filepath.FromSlash("/sect5/page3.md")}, "Title5_3"}, // test OS-specific path
		{"Absolute, unique", kinds.KindPage, nil, []string{"/sect3/unique.md"}, "UniqueBase"},
		{"Absolute, unique, case", kinds.KindPage, nil, []string{"/sect3/Unique2.md", "/sect3/unique2.md", "/sect3/unique2", "/sect3/Unique2"}, "UniqueBase2"},
		// next test depends on this page existing
		// {"NoPage", nil, []string{"/unique.md"}, ""},  // ISSUE #4969: this is resolving to /sect3/unique.md
		{"Absolute, missing page", "NoPage", nil, []string{"/missing-page.md"}, ""},
		{"Absolute, missing section", "NoPage", nil, []string{"/missing-section"}, ""},

		// relative paths
		{"Dot relative, home", kinds.KindHome, sec3, []string{".."}, "home page"},
		{"Dot relative, home, slash", kinds.KindHome, sec3, []string{"../"}, "home page"},
		{"Dot relative about", kinds.KindPage, sec3, []string{"../about.md"}, "about page"},
		{"Dot", kinds.KindSection, sec3, []string{"."}, "section 3"},
		{"Dot slash", kinds.KindSection, sec3, []string{"./"}, "section 3"},
		{"Page relative, no dot", kinds.KindPage, sec3, []string{"page1.md"}, "Title3_1"},
		{"Page relative, dot", kinds.KindPage, sec3, []string{"./page1.md"}, "Title3_1"},
		{"Up and down another section", kinds.KindPage, sec3, []string{"../sect4/page2.md"}, "Title4_2"},
		{"Rel sect7", kinds.KindSection, sec3, []string{"sect7"}, "another sect7"},
		{"Rel sect7 dot", kinds.KindSection, sec3, []string{"./sect7"}, "another sect7"},
		{"Dot deep", kinds.KindPage, sec3, []string{"./subsect/deep.md"}, "deep page"},
		{"Dot dot inner", kinds.KindPage, sec3, []string{"./subsect/../../sect7/page9.md"}, "Title7_9"},
		{"Dot OS slash", kinds.KindPage, sec3, []string{filepath.FromSlash("../sect5/page3.md")}, "Title5_3"}, // test OS-specific path
		{"Dot unique", kinds.KindPage, sec3, []string{"./unique.md"}, "UniqueBase"},
		{"Dot sect", "NoPage", sec3, []string{"./sect2"}, ""},
		//{"NoPage", sec3, []string{"sect2"}, ""}, // ISSUE: /sect3 page relative query is resolving to /sect2

		{"Abs, ignore context, home", kinds.KindHome, sec3, []string{"/"}, "home page"},
		{"Abs, ignore context, about", kinds.KindPage, sec3, []string{"/about.md"}, "about page"},
		{"Abs, ignore context, page in section", kinds.KindPage, sec3, []string{"/sect4/page2.md"}, "Title4_2"},
		{"Abs, ignore context, page subsect deep", kinds.KindPage, sec3, []string{"/sect3/subsect/deep.md"}, "deep page"}, // next test depends on this page existing
		{"Abs, ignore context, page deep", "NoPage", sec3, []string{"/subsect/deep.md"}, ""},

		// Taxonomies
		{"Taxonomy term", kinds.KindTaxonomy, nil, []string{"categories"}, "Categories"},
		{"Taxonomy", kinds.KindTerm, nil, []string{"categories/hugo", "categories/Hugo"}, "Hugo"},

		// Bundle variants
		{"Bundle regular", kinds.KindPage, nil, []string{"sect3/b1", "sect3/b1/index.md", "sect3/b1/index.en.md"}, "b1 bundle"},
		{"Bundle index name", kinds.KindPage, nil, []string{"sect3/index/index.md", "sect3/index"}, "index bundle"},

		// https://github.com/gohugoio/hugo/issues/7301
		{"Section and bundle overlap", kinds.KindPage, nil, []string{"section_bundle_overlap_bundle"}, "index overlap bundle"},
	}

	for _, test := range tests {
		b.Run(test.name, func(c *qt.C) {
			errorMsg := fmt.Sprintf("Test case %v %v -> %s", test.context, test.pathVariants, test.expectedTitle)

			// test legacy public Site.GetPage (which does not support page context relative queries)
			if test.context == nil {
				for _, ref := range test.pathVariants {
					args := append([]string{test.kind}, ref)
					page, err := s.GetPage(args...)
					test.check(page, err, errorMsg, c)
				}
			}

			// test new internal Site.getPage
			for _, ref := range test.pathVariants {
				page2, err := s.getPage(test.context, ref)
				test.check(page2, err, errorMsg, c)
			}
		})
	}
}

// #11664
func TestGetPageIndexIndex(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]	
-- content/mysect/index/index.md --
---
title: "Mysect Index"
---
-- layouts/index.html --
GetPage 1: {{ with site.GetPage "mysect/index/index.md" }}{{ .Title }}|{{ .RelPermalink }}|{{ .Path }}{{ end }}|
GetPage 2: {{ with site.GetPage "mysect/index" }}{{ .Title }}|{{ .RelPermalink }}|{{ .Path }}{{ end }}|
`

	b := Test(t, files)
	b.AssertFileContent("public/index.html",
		"GetPage 1: Mysect Index|/mysect/index/|/mysect/index|",
		"GetPage 2: Mysect Index|/mysect/index/|/mysect/index|",
	)
}

// https://github.com/gohugoio/hugo/issues/6034
func TestGetPageRelative(t *testing.T) {
	files := `
-- hugo.toml --
-- content/what/_index.md --
---title: what
---
-- content/what/members.md --
---title: members what
draft: false
---
-- content/where/_index.md --
---title: where
---
-- content/where/members.md --
---title: members where
draft: false
---
-- content/who/_index.md --
---title: who
---
-- content/who/members.md --
---title: members who
draft: true
---
-- layouts/_default/list.html --
{{ with .GetPage "members.md" }}
    Members: {{ .Title }}
{{ else }}
NOT FOUND
{{ end }}
`

	b := Test(t, files)

	b.AssertFileContent("public/what/index.html", `Members: members what`)
	b.AssertFileContent("public/where/index.html", `Members: members where`)
	b.AssertFileContent("public/who/index.html", `NOT FOUND`)
}

func TestGetPageIssue11883(t *testing.T) {
	files := `
-- hugo.toml --
-- p1/index.md --
---
title: p1
---
-- p1/p1.xyz --
xyz.
-- layouts/index.html --
Home. {{ with .Page.GetPage "p1.xyz" }}{{ else }}OK 1{{ end }} {{ with .Site.GetPage "p1.xyz" }}{{ else }}OK 2{{ end }}
`

	b := Test(t, files)
	b.AssertFileContent("public/index.html", "Home. OK 1 OK 2")
}

func TestGetPageIssue12120(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
-- content/s1/p1/index.md --
---
title: p1
layout: p1
---
-- content/s1/p2.md --
---
title: p2
layout: p2
---
-- layouts/_default/p1.html --
{{ (.GetPage "p2.md").Title }}|
-- layouts/_default/p2.html --
{{ (.GetPage "p1").Title }}|
`

	b := Test(t, files)
	b.AssertFileContent("public/s1/p1/index.html", "p2") // failing test
	b.AssertFileContent("public/s1/p2/index.html", "p1")
}

func TestGetPageNewsVsTagsNewsIssue12638(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','section','sitemap']
[taxonomies]
  tag = "tags"
-- content/p1.md --
---
title: p1
tags: [news]
---
-- layouts/index.html --
/tags/news: {{ with .Site.GetPage "/tags/news" }}{{ .Title }}{{ end }}|
news: {{ with .Site.GetPage "news" }}{{ .Title }}{{ end }}|
/news: {{ with .Site.GetPage "/news" }}{{ .Title }}{{ end }}|

`

	b := Test(t, files)

	b.AssertFileContent("public/index.html",
		"/tags/news: News|",
		"news: News|",
		"/news: |",
	)
}

func TestGetPageBundleToRegular(t *testing.T) {
	files := `
-- hugo.toml --
-- content/s1/p1/index.md --
---
title: p1
---
-- content/s1/p2.md --
---
title: p2
---
-- layouts/_default/single.html --
{{ with .GetPage "p2" }}
  OK: {{ .LinkTitle }}
{{ else }}
   Unable to get p2.
{{ end }}
`

	b := Test(t, files)
	b.AssertFileContent("public/s1/p1/index.html", "OK: p2")
	b.AssertFileContent("public/s1/p2/index.html", "OK: p2")
}

func TestPageGetPageVariations(t *testing.T) {
	files := `
-- hugo.toml --
-- content/s1/_index.md --
---
title: s1 section
---
-- content/s1/p1/index.md --
---
title: p1
---
-- content/s1/p2.md --
---
title: p2
---
-- content/s2/p3/index.md --
---
title: p3
---
-- content/p2.md --
---
title: p2_root
---
-- layouts/index.html --
/s1: {{ with .GetPage "/s1" }}{{ .Title }}{{ end }}|
/s1/: {{ with .GetPage "/s1/" }}{{ .Title }}{{ end }}|
/s1/p2.md: {{ with .GetPage "/s1/p2.md" }}{{ .Title }}{{ end }}|
/s1/p2: {{ with .GetPage "/s1/p2" }}{{ .Title }}{{ end }}|
/s1/p1/index.md: {{ with .GetPage "/s1/p1/index.md" }}{{ .Title }}{{ end }}|
/s1/p1: {{ with .GetPage "/s1/p1" }}{{ .Title }}{{ end }}|
-- layouts/_default/single.html --
../p2: {{ with .GetPage "../p2" }}{{ .Title }}{{ end }}|
../p2.md: {{ with .GetPage "../p2.md" }}{{ .Title }}{{ end }}|
p1/index.md: {{ with .GetPage "p1/index.md" }}{{ .Title }}{{ end }}|
../s2/p3/index.md: {{ with .GetPage "../s2/p3/index.md" }}{{ .Title }}{{ end }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", `
/s1: s1 section|
/s1/: s1 section|
/s1/p2.md: p2|
/s1/p2: p2|
/s1/p1/index.md: p1|
/s1/p1: p1|
`)

	b.AssertFileContent("public/s1/p1/index.html", `
../p2: p2_root|
../p2.md: p2|

`)

	b.AssertFileContent("public/s1/p2/index.html", `
../p2: p2_root|	 
../p2.md: p2_root|
p1/index.md: p1|
../s2/p3/index.md: p3|

`)
}

func TestPageGetPageMountsReverseLookup(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
[module]
[[module.mounts]]
source = "layouts"
target = "layouts"
[[module.mounts]]
source = "README.md"
target = "content/_index.md"
[[module.mounts]]
source = "blog"
target = "content/posts"
[[module.mounts]]
source = "docs"
target = "content/mydocs"
-- README.md --
---
title: README
---
-- blog/b1.md --
---
title: b1
---
{{< ref "../docs/d1.md" >}}
-- blog/b2/index.md --
---
title: b2
---
{{< ref "../../docs/d1.md" >}}
-- docs/d1.md --
---
title: d1
---
-- layouts/shortcodes/ref.html --
{{ $ref := .Get 0 }}
.Page.GetPage({{ $ref }}).Title: {{ with .Page.GetPage $ref }}{{ .Title }}{{ end }}|
-- layouts/index.html --
Home.
/blog/b1.md: {{ with .GetPage "/blog/b1.md" }}{{ .Title }}{{ end }}|
/blog/b2/index.md: {{ with .GetPage "/blog/b2/index.md" }}{{ .Title }}{{ end }}|
/docs/d1.md: {{ with .GetPage "/docs/d1.md" }}{{ .Title }}{{ end }}|
/README.md: {{ with .GetPage "/README.md" }}{{ .Title }}{{ end }}|
-- layouts/_default/single.html --
Single.
/README.md: {{ with .GetPage "/README.md" }}{{ .Title }}{{ end }}|
{{ .Content }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html",
		`
/blog/b1.md: b1|
/blog/b2/index.md: b2|
/docs/d1.md: d1|
/README.md: README
`,
	)

	b.AssertFileContent("public/mydocs/d1/index.html", `README.md: README|`)

	b.AssertFileContent("public/posts/b1/index.html", `.Page.GetPage(../docs/d1.md).Title: d1|`)
	b.AssertFileContent("public/posts/b2/index.html", `.Page.GetPage(../../docs/d1.md).Title: d1|`)
}

// https://github.com/gohugoio/hugo/issues/7016
func TestGetPageMultilingual(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.org/"
languageCode = "en-us"
defaultContentLanguage = "ru"
title = "My New Hugo Site"
uglyurls = true

[languages]
[languages.ru]
[languages.en]
-- content/docs/1.md --
---
title: p1
---
-- content/news/1.md --
---
title: p1
---
-- content/news/1.en.md --
---
title: p1en
---
-- content/news/about/1.md --
---
title: about1
---
-- content/news/about/1.en.md --
---
title: about1en
---
-- layouts/index.html --
{{ with site.GetPage "docs/1" }}
    Docs p1: {{ .Title }}
{{ else }}
NOT FOUND
{{ end }}
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `Docs p1: p1`)
	b.AssertFileContent("public/en/index.html", `NOT FOUND`)
}

func TestRegularPagesRecursive(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.org/"
title = "My New Hugo Site"
-- content/docs/1.md --
---
title: docs1
---
-- content/docs/sect1/_index.md --
---
title: docs_sect1
---
-- content/docs/sect1/ps1.md --
---
title: docs_sect1_ps1
---
-- content/docs/sect1/ps2.md --
---
title: docs_sect1_ps2
---
-- content/docs/sect1/sect1_s2/_index.md --
---
title: docs_sect1_s2
---
-- content/docs/sect1/sect1_s2/ps2_1.md --
---
title: docs_sect1_s2_1
---
-- content/docs/sect2/_index.md --
---
title: docs_sect2
---
-- content/docs/sect2/ps1.md --
---
title: docs_sect2_ps1
---
-- content/docs/sect2/ps2.md --
---
title: docs_sect2_ps2
---
-- content/news/1.md --
---
title: news1
---
-- layouts/index.html --
{{ $sect1 := site.GetPage "sect1" }}

Sect1 RegularPagesRecursive: {{ range $sect1.RegularPagesRecursive }}{{ .Kind }}:{{ .RelPermalink}}|{{ end }}|End.
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `
Sect1 RegularPagesRecursive: page:/docs/sect1/ps1/|page:/docs/sect1/ps2/|page:/docs/sect1/sect1_s2/ps2_1/||End.


`)
}

func TestRegularPagesRecursiveHome(t *testing.T) {
	files := `
-- hugo.toml --
-- content/p1.md --
-- content/post/p2.md --
-- layouts/index.html --
RegularPagesRecursive: {{ range .RegularPagesRecursive }}{{ .Kind }}:{{ .RelPermalink}}|{{ end }}|End.
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		}).Build()

	b.AssertFileContent("public/index.html", `RegularPagesRecursive: page:/p1/|page:/post/p2/||End.`)
}

// Issue #12169.
func TestPagesSimilarSectionNames(t *testing.T) {
	files := `
-- hugo.toml --
-- content/draftsection/_index.md --
---
draft: true
---
-- content/draftsection/sub/_index.md --got
-- content/draftsection/sub/d1.md --
-- content/s1/_index.md --
-- content/s1/p1.md --
-- content/s1-foo/_index.md --
-- content/s1-foo/p2.md --
-- content/s1-foo/s2/_index.md --
-- content/s1-foo/s2/p3.md --
-- content/s1-foo/s2-foo/_index.md --
-- content/s1-foo/s2-foo/p4.md --
-- layouts/_default/list.html --
{{ .RelPermalink }}: Pages: {{ range .Pages }}{{ .RelPermalink }}|{{ end }}$

`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "/: Pages: /s1-foo/|/s1/|$")
	b.AssertFileContent("public/s1/index.html", "/s1/: Pages: /s1/p1/|$")
	b.AssertFileContent("public/s1-foo/index.html", "/s1-foo/: Pages: /s1-foo/p2/|/s1-foo/s2-foo/|/s1-foo/s2/|$")
	b.AssertFileContent("public/s1-foo/s2/index.html", "/s1-foo/s2/: Pages: /s1-foo/s2/p3/|$")
}

func TestGetPageContentAdapterBaseIssue12561(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
Test A: {{ (site.GetPage "/s1/p1").Title }}
Test B: {{ (site.GetPage "p1").Title }}
Test C: {{ (site.GetPage "/s2/p2").Title }}
Test D: {{ (site.GetPage "p2").Title }}
-- layouts/_default/single.html --
{{ .Title }}
-- content/s1/p1.md --
---
title: p1
---
-- content/s2/_content.gotmpl --
{{ .AddPage (dict "path" "p2" "title" "p2") }}
`

	b := Test(t, files)

	b.AssertFileExists("public/s1/p1/index.html", true)
	b.AssertFileExists("public/s2/p2/index.html", true)
	b.AssertFileContent("public/index.html",
		"Test A: p1",
		"Test B: p1",
		"Test C: p2",
		"Test D: p2", // fails
	)
}
