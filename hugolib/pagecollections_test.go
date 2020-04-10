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
	"math/rand"
	"path"
	"path/filepath"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/deps"
)

const pageCollectionsPageTemplate = `---
title: "%s"
categories:
- Hugo
---
# Doc
`

func BenchmarkGetPage(b *testing.B) {
	var (
		cfg, fs = newTestCfg()
		r       = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			writeSource(b, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), "CONTENT")
		}
	}

	s := buildSingleSite(b, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	pagePaths := make([]string, b.N)

	for i := 0; i < b.N; i++ {
		pagePaths[i] = fmt.Sprintf("sect%d", r.Intn(10))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		home, _ := s.getPageNew(nil, "/")
		if home == nil {
			b.Fatal("Home is nil")
		}

		p, _ := s.getPageNew(nil, pagePaths[i])
		if p == nil {
			b.Fatal("Section is nil")
		}

	}
}

func createGetPageRegularBenchmarkSite(t testing.TB) *Site {

	var (
		c       = qt.New(t)
		cfg, fs = newTestCfg()
	)

	pc := func(title string) string {
		return fmt.Sprintf(pageCollectionsPageTemplate, title)
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			content := pc(fmt.Sprintf("Title%d_%d", i, j))
			writeSource(c, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), content)
		}
	}

	return buildSingleSite(c, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

}

func TestBenchmarkGetPageRegular(t *testing.T) {
	c := qt.New(t)
	s := createGetPageRegularBenchmarkSite(t)

	for i := 0; i < 10; i++ {
		pp := path.Join("/", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", i))
		page, _ := s.getPageNew(nil, pp)
		c.Assert(page, qt.Not(qt.IsNil), qt.Commentf(pp))
	}
}

func BenchmarkGetPageRegular(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.Run("From root", func(b *testing.B) {
		s := createGetPageRegularBenchmarkSite(b)
		c := qt.New(b)

		pagePaths := make([]string, b.N)

		for i := 0; i < b.N; i++ {
			pagePaths[i] = path.Join(fmt.Sprintf("/sect%d", r.Intn(10)), fmt.Sprintf("page%d.md", r.Intn(100)))
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			page, _ := s.getPageNew(nil, pagePaths[i])
			c.Assert(page, qt.Not(qt.IsNil))
		}
	})

	b.Run("Page relative", func(b *testing.B) {
		s := createGetPageRegularBenchmarkSite(b)
		c := qt.New(b)
		allPages := s.RegularPages()

		pagePaths := make([]string, b.N)
		pages := make([]page.Page, b.N)

		for i := 0; i < b.N; i++ {
			pagePaths[i] = fmt.Sprintf("page%d.md", r.Intn(100))
			pages[i] = allPages[r.Intn(len(allPages)/3)]
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			page, _ := s.getPageNew(pages[i], pagePaths[i])
			c.Assert(page, qt.Not(qt.IsNil))
		}
	})

}

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

	var (
		cfg, fs = newTestCfg()
		c       = qt.New(t)
	)

	pc := func(title string) string {
		return fmt.Sprintf(pageCollectionsPageTemplate, title)
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			content := pc(fmt.Sprintf("Title%d_%d", i, j))
			writeSource(t, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), content)
		}
	}

	content := pc("home page")
	writeSource(t, fs, filepath.Join("content", "_index.md"), content)

	content = pc("about page")
	writeSource(t, fs, filepath.Join("content", "about.md"), content)

	content = pc("section 3")
	writeSource(t, fs, filepath.Join("content", "sect3", "_index.md"), content)

	writeSource(t, fs, filepath.Join("content", "sect3", "unique.md"), pc("UniqueBase"))
	writeSource(t, fs, filepath.Join("content", "sect3", "Unique2.md"), pc("UniqueBase2"))

	content = pc("another sect7")
	writeSource(t, fs, filepath.Join("content", "sect3", "sect7", "_index.md"), content)

	content = pc("deep page")
	writeSource(t, fs, filepath.Join("content", "sect3", "subsect", "deep.md"), content)

	// Bundle variants
	writeSource(t, fs, filepath.Join("content", "sect3", "b1", "index.md"), pc("b1 bundle"))
	writeSource(t, fs, filepath.Join("content", "sect3", "index", "index.md"), pc("index bundle"))

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	sec3, err := s.getPageNew(nil, "/sect3")
	c.Assert(err, qt.IsNil)
	c.Assert(sec3, qt.Not(qt.IsNil))

	tests := []getPageTest{
		// legacy content root relative paths
		{"Root relative, no slash, home", page.KindHome, nil, []string{""}, "home page"},
		{"Root relative, no slash, root page", page.KindPage, nil, []string{"about.md", "ABOUT.md"}, "about page"},
		{"Root relative, no slash, section", page.KindSection, nil, []string{"sect3"}, "section 3"},
		{"Root relative, no slash, section page", page.KindPage, nil, []string{"sect3/page1.md"}, "Title3_1"},
		{"Root relative, no slash, sub setion", page.KindSection, nil, []string{"sect3/sect7"}, "another sect7"},
		{"Root relative, no slash, nested page", page.KindPage, nil, []string{"sect3/subsect/deep.md"}, "deep page"},
		{"Root relative, no slash, OS slashes", page.KindPage, nil, []string{filepath.FromSlash("sect5/page3.md")}, "Title5_3"},

		{"Short ref, unique", page.KindPage, nil, []string{"unique.md", "unique"}, "UniqueBase"},
		{"Short ref, unique, upper case", page.KindPage, nil, []string{"Unique2.md", "unique2.md", "unique2"}, "UniqueBase2"},
		{"Short ref, ambiguous", "Ambiguous", nil, []string{"page1.md"}, ""},

		// ISSUE: This is an ambiguous ref, but because we have to support the legacy
		// content root relative paths without a leading slash, the lookup
		// returns /sect7. This undermines ambiguity detection, but we have no choice.
		//{"Ambiguous", nil, []string{"sect7"}, ""},
		{"Section, ambigous", page.KindSection, nil, []string{"sect7"}, "Sect7s"},

		{"Absolute, home", page.KindHome, nil, []string{"/", ""}, "home page"},
		{"Absolute, page", page.KindPage, nil, []string{"/about.md", "/about"}, "about page"},
		{"Absolute, sect", page.KindSection, nil, []string{"/sect3"}, "section 3"},
		{"Absolute, page in subsection", page.KindPage, nil, []string{"/sect3/page1.md", "/Sect3/Page1.md"}, "Title3_1"},
		{"Absolute, section, subsection with same name", page.KindSection, nil, []string{"/sect3/sect7"}, "another sect7"},
		{"Absolute, page, deep", page.KindPage, nil, []string{"/sect3/subsect/deep.md"}, "deep page"},
		{"Absolute, page, OS slashes", page.KindPage, nil, []string{filepath.FromSlash("/sect5/page3.md")}, "Title5_3"}, //test OS-specific path
		{"Absolute, unique", page.KindPage, nil, []string{"/sect3/unique.md"}, "UniqueBase"},
		{"Absolute, unique, case", page.KindPage, nil, []string{"/sect3/Unique2.md", "/sect3/unique2.md", "/sect3/unique2", "/sect3/Unique2"}, "UniqueBase2"},
		//next test depends on this page existing
		// {"NoPage", nil, []string{"/unique.md"}, ""},  // ISSUE #4969: this is resolving to /sect3/unique.md
		{"Absolute, missing page", "NoPage", nil, []string{"/missing-page.md"}, ""},
		{"Absolute, missing section", "NoPage", nil, []string{"/missing-section"}, ""},

		// relative paths
		{"Dot relative, home", page.KindHome, sec3, []string{".."}, "home page"},
		{"Dot relative, home, slash", page.KindHome, sec3, []string{"../"}, "home page"},
		{"Dot relative about", page.KindPage, sec3, []string{"../about.md"}, "about page"},
		{"Dot", page.KindSection, sec3, []string{"."}, "section 3"},
		{"Dot slash", page.KindSection, sec3, []string{"./"}, "section 3"},
		{"Page relative, no dot", page.KindPage, sec3, []string{"page1.md"}, "Title3_1"},
		{"Page relative, dot", page.KindPage, sec3, []string{"./page1.md"}, "Title3_1"},
		{"Up and down another section", page.KindPage, sec3, []string{"../sect4/page2.md"}, "Title4_2"},
		{"Rel sect7", page.KindSection, sec3, []string{"sect7"}, "another sect7"},
		{"Rel sect7 dot", page.KindSection, sec3, []string{"./sect7"}, "another sect7"},
		{"Dot deep", page.KindPage, sec3, []string{"./subsect/deep.md"}, "deep page"},
		{"Dot dot inner", page.KindPage, sec3, []string{"./subsect/../../sect7/page9.md"}, "Title7_9"},
		{"Dot OS slash", page.KindPage, sec3, []string{filepath.FromSlash("../sect5/page3.md")}, "Title5_3"}, //test OS-specific path
		{"Dot unique", page.KindPage, sec3, []string{"./unique.md"}, "UniqueBase"},
		{"Dot sect", "NoPage", sec3, []string{"./sect2"}, ""},
		//{"NoPage", sec3, []string{"sect2"}, ""}, // ISSUE: /sect3 page relative query is resolving to /sect2

		{"Abs, ignore context, home", page.KindHome, sec3, []string{"/"}, "home page"},
		{"Abs, ignore context, about", page.KindPage, sec3, []string{"/about.md"}, "about page"},
		{"Abs, ignore context, page in section", page.KindPage, sec3, []string{"/sect4/page2.md"}, "Title4_2"},
		{"Abs, ignore context, page subsect deep", page.KindPage, sec3, []string{"/sect3/subsect/deep.md"}, "deep page"}, //next test depends on this page existing
		{"Abs, ignore context, page deep", "NoPage", sec3, []string{"/subsect/deep.md"}, ""},

		// Taxonomies
		{"Taxonomy term", page.KindTaxonomyTerm, nil, []string{"categories"}, "Categories"},
		{"Taxonomy", page.KindTaxonomy, nil, []string{"categories/hugo", "categories/Hugo"}, "Hugo"},

		// Bundle variants
		{"Bundle regular", page.KindPage, nil, []string{"sect3/b1", "sect3/b1/index.md", "sect3/b1/index.en.md"}, "b1 bundle"},
		{"Bundle index name", page.KindPage, nil, []string{"sect3/index/index.md", "sect3/index"}, "index bundle"},
	}

	for _, test := range tests {
		c.Run(test.name, func(c *qt.C) {
			errorMsg := fmt.Sprintf("Test case %v %v -> %s", test.context, test.pathVariants, test.expectedTitle)

			// test legacy public Site.GetPage (which does not support page context relative queries)
			if test.context == nil {
				for _, ref := range test.pathVariants {
					args := append([]string{test.kind}, ref)
					page, err := s.Info.GetPage(args...)
					test.check(page, err, errorMsg, c)
				}
			}

			// test new internal Site.getPageNew
			for _, ref := range test.pathVariants {
				page2, err := s.getPageNew(test.context, ref)
				test.check(page2, err, errorMsg, c)
			}

		})
	}

}

// https://github.com/gohugoio/hugo/issues/6034
func TestGetPageRelative(t *testing.T) {
	b := newTestSitesBuilder(t)
	for i, section := range []string{"what", "where", "who"} {
		isDraft := i == 2
		b.WithContent(
			section+"/_index.md", fmt.Sprintf("---title: %s\n---", section),
			section+"/members.md", fmt.Sprintf("---title: members %s\ndraft: %t\n---", section, isDraft),
		)
	}

	b.WithTemplates("_default/list.html", `
{{ with .GetPage "members.md" }}
    Members: {{ .Title }}
{{ else }}
NOT FOUND
{{ end }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/what/index.html", `Members: members what`)
	b.AssertFileContent("public/where/index.html", `Members: members where`)
	b.AssertFileContent("public/who/index.html", `NOT FOUND`)

}

// https://github.com/gohugoio/hugo/issues/7016
func TestGetPageMultilingual(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithConfigFile("yaml", `
baseURL: "http://example.org/"
languageCode: "en-us"
defaultContentLanguage: ru
title: "My New Hugo Site"
uglyurls: true

languages:
  ru: {}
  en: {}
`)

	b.WithContent(
		"docs/1.md", "\n---title: p1\n---",
		"news/1.md", "\n---title: p1\n---",
		"news/1.en.md", "\n---title: p1en\n---",
		"news/about/1.md", "\n---title: about1\n---",
		"news/about/1.en.md", "\n---title: about1en\n---",
	)

	b.WithTemplates("index.html", `
{{ with site.GetPage "docs/1" }}
    Docs p1: {{ .Title }}
{{ else }}
NOT FOUND
{{ end }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `Docs p1: p1`)
	b.AssertFileContent("public/en/index.html", `NOT FOUND`)

}

func TestShouldDoSimpleLookup(t *testing.T) {
	c := qt.New(t)

	c.Assert(shouldDoSimpleLookup("foo.md"), qt.Equals, true)
	c.Assert(shouldDoSimpleLookup("/foo.md"), qt.Equals, true)
	c.Assert(shouldDoSimpleLookup("./foo.md"), qt.Equals, false)
	c.Assert(shouldDoSimpleLookup("docs/foo.md"), qt.Equals, false)

}

func TestRegularPagesRecursive(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithConfigFile("yaml", `
baseURL: "http://example.org/"
title: "My New Hugo Site"

`)

	b.WithContent(
		"docs/1.md", "\n---title: docs1\n---",
		"docs/sect1/_index.md", "\n---title: docs_sect1\n---",
		"docs/sect1/ps1.md", "\n---title: docs_sect1_ps1\n---",
		"docs/sect1/ps2.md", "\n---title: docs_sect1_ps2\n---",
		"docs/sect1/sect1_s2/_index.md", "\n---title: docs_sect1_s2\n---",
		"docs/sect1/sect1_s2/ps2_1.md", "\n---title: docs_sect1_s2_1\n---",
		"docs/sect2/_index.md", "\n---title: docs_sect2\n---",
		"docs/sect2/ps1.md", "\n---title: docs_sect2_ps1\n---",
		"docs/sect2/ps2.md", "\n---title: docs_sect2_ps2\n---",
		"news/1.md", "\n---title: news1\n---",
	)

	b.WithTemplates("index.html", `
{{ $sect1 := site.GetPage "sect1" }}

Sect1 RegularPagesRecursive: {{ range $sect1.RegularPagesRecursive }}{{ .Kind }}:{{ .RelPermalink}}|{{ end }}|End.

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
Sect1 RegularPagesRecursive: page:/docs/sect1/ps1/|page:/docs/sect1/ps2/|page:/docs/sect1/sect1_s2/ps2_1/||End.


`)

}
