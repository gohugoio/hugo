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

func BenchmarkGetPageRegular(b *testing.B) {
	var (
		c       = qt.New(b)
		cfg, fs = newTestCfg()
		r       = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			content := fmt.Sprintf(pageCollectionsPageTemplate, fmt.Sprintf("Title%d_%d", i, j))
			writeSource(b, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), content)
		}
	}

	s := buildSingleSite(b, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	pagePaths := make([]string, b.N)

	for i := 0; i < b.N; i++ {
		pagePaths[i] = path.Join(fmt.Sprintf("sect%d", r.Intn(10)), fmt.Sprintf("page%d.md", r.Intn(100)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		page, _ := s.getPageNew(nil, pagePaths[i])
		c.Assert(page, qt.Not(qt.IsNil))
	}
}

type testCase struct {
	kind          string
	context       page.Page
	path          []string
	expectedTitle string
}

func (t *testCase) check(p page.Page, err error, errorMsg string, c *qt.C) {
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

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			content := fmt.Sprintf(pageCollectionsPageTemplate, fmt.Sprintf("Title%d_%d", i, j))
			writeSource(t, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), content)
		}
	}

	content := fmt.Sprintf(pageCollectionsPageTemplate, "home page")
	writeSource(t, fs, filepath.Join("content", "_index.md"), content)

	content = fmt.Sprintf(pageCollectionsPageTemplate, "about page")
	writeSource(t, fs, filepath.Join("content", "about.md"), content)

	content = fmt.Sprintf(pageCollectionsPageTemplate, "section 3")
	writeSource(t, fs, filepath.Join("content", "sect3", "_index.md"), content)

	content = fmt.Sprintf(pageCollectionsPageTemplate, "UniqueBase")
	writeSource(t, fs, filepath.Join("content", "sect3", "unique.md"), content)

	content = fmt.Sprintf(pageCollectionsPageTemplate, "another sect7")
	writeSource(t, fs, filepath.Join("content", "sect3", "sect7", "_index.md"), content)

	content = fmt.Sprintf(pageCollectionsPageTemplate, "deep page")
	writeSource(t, fs, filepath.Join("content", "sect3", "subsect", "deep.md"), content)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	sec3, err := s.getPageNew(nil, "/sect3")
	c.Assert(err, qt.IsNil)
	c.Assert(sec3, qt.Not(qt.IsNil))

	tests := []testCase{
		// legacy content root relative paths
		{page.KindHome, nil, []string{}, "home page"},
		{page.KindPage, nil, []string{"about.md"}, "about page"},
		{page.KindSection, nil, []string{"sect3"}, "section 3"},
		{page.KindPage, nil, []string{"sect3/page1.md"}, "Title3_1"},
		{page.KindPage, nil, []string{"sect4/page2.md"}, "Title4_2"},
		{page.KindSection, nil, []string{"sect3/sect7"}, "another sect7"},
		{page.KindPage, nil, []string{"sect3/subsect/deep.md"}, "deep page"},
		{page.KindPage, nil, []string{filepath.FromSlash("sect5/page3.md")}, "Title5_3"}, //test OS-specific path

		// shorthand refs (potentially ambiguous)
		{page.KindPage, nil, []string{"unique.md"}, "UniqueBase"},
		{"Ambiguous", nil, []string{"page1.md"}, ""},

		// ISSUE: This is an ambiguous ref, but because we have to support the legacy
		// content root relative paths without a leading slash, the lookup
		// returns /sect7. This undermines ambiguity detection, but we have no choice.
		//{"Ambiguous", nil, []string{"sect7"}, ""},
		{page.KindSection, nil, []string{"sect7"}, "Sect7s"},

		// absolute paths
		{page.KindHome, nil, []string{"/"}, "home page"},
		{page.KindPage, nil, []string{"/about.md"}, "about page"},
		{page.KindSection, nil, []string{"/sect3"}, "section 3"},
		{page.KindPage, nil, []string{"/sect3/page1.md"}, "Title3_1"},
		{page.KindPage, nil, []string{"/sect4/page2.md"}, "Title4_2"},
		{page.KindSection, nil, []string{"/sect3/sect7"}, "another sect7"},
		{page.KindPage, nil, []string{"/sect3/subsect/deep.md"}, "deep page"},
		{page.KindPage, nil, []string{filepath.FromSlash("/sect5/page3.md")}, "Title5_3"}, //test OS-specific path
		{page.KindPage, nil, []string{"/sect3/unique.md"}, "UniqueBase"},                  //next test depends on this page existing
		// {"NoPage", nil, []string{"/unique.md"}, ""},  // ISSUE #4969: this is resolving to /sect3/unique.md
		{"NoPage", nil, []string{"/missing-page.md"}, ""},
		{"NoPage", nil, []string{"/missing-section"}, ""},

		// relative paths
		{page.KindHome, sec3, []string{".."}, "home page"},
		{page.KindHome, sec3, []string{"../"}, "home page"},
		{page.KindPage, sec3, []string{"../about.md"}, "about page"},
		{page.KindSection, sec3, []string{"."}, "section 3"},
		{page.KindSection, sec3, []string{"./"}, "section 3"},
		{page.KindPage, sec3, []string{"page1.md"}, "Title3_1"},
		{page.KindPage, sec3, []string{"./page1.md"}, "Title3_1"},
		{page.KindPage, sec3, []string{"../sect4/page2.md"}, "Title4_2"},
		{page.KindSection, sec3, []string{"sect7"}, "another sect7"},
		{page.KindSection, sec3, []string{"./sect7"}, "another sect7"},
		{page.KindPage, sec3, []string{"./subsect/deep.md"}, "deep page"},
		{page.KindPage, sec3, []string{"./subsect/../../sect7/page9.md"}, "Title7_9"},
		{page.KindPage, sec3, []string{filepath.FromSlash("../sect5/page3.md")}, "Title5_3"}, //test OS-specific path
		{page.KindPage, sec3, []string{"./unique.md"}, "UniqueBase"},
		{"NoPage", sec3, []string{"./sect2"}, ""},
		//{"NoPage", sec3, []string{"sect2"}, ""}, // ISSUE: /sect3 page relative query is resolving to /sect2

		// absolute paths ignore context
		{page.KindHome, sec3, []string{"/"}, "home page"},
		{page.KindPage, sec3, []string{"/about.md"}, "about page"},
		{page.KindPage, sec3, []string{"/sect4/page2.md"}, "Title4_2"},
		{page.KindPage, sec3, []string{"/sect3/subsect/deep.md"}, "deep page"}, //next test depends on this page existing
		{"NoPage", sec3, []string{"/subsect/deep.md"}, ""},
	}

	for _, test := range tests {
		errorMsg := fmt.Sprintf("Test case %s %v -> %s", test.context, test.path, test.expectedTitle)

		// test legacy public Site.GetPage (which does not support page context relative queries)
		if test.context == nil {
			args := append([]string{test.kind}, test.path...)
			page, err := s.Info.GetPage(args...)
			test.check(page, err, errorMsg, c)
		}

		// test new internal Site.getPageNew
		var ref string
		if len(test.path) == 1 {
			ref = filepath.ToSlash(test.path[0])
		} else {
			ref = path.Join(test.path...)
		}
		page2, err := s.getPageNew(test.context, ref)
		test.check(page2, err, errorMsg, c)
	}

}
