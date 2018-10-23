// Copyright 2017 The Hugo Authors. All rights reserved.
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

	"github.com/gohugoio/hugo/deps"
	"github.com/stretchr/testify/require"
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
		require.NotNil(b, page)
	}
}

type testCase struct {
	kind          string
	context       *Page
	path          []string
	expectedTitle string
}

func (t *testCase) check(p *Page, err error, errorMsg string, assert *require.Assertions) {
	switch t.kind {
	case "Ambiguous":
		assert.Error(err)
		assert.Nil(p, errorMsg)
	case "NoPage":
		assert.NoError(err)
		assert.Nil(p, errorMsg)
	default:
		assert.NoError(err, errorMsg)
		assert.NotNil(p, errorMsg)
		assert.Equal(t.kind, p.Kind, errorMsg)
		assert.Equal(t.expectedTitle, p.title, errorMsg)
	}
}

func TestGetPage(t *testing.T) {

	var (
		assert  = require.New(t)
		cfg, fs = newTestCfg()
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
	assert.NoError(err, "error getting Page for /sec3")
	assert.NotNil(sec3, "failed to get Page for /sec3")

	tests := []testCase{
		// legacy content root relative paths
		{KindHome, nil, []string{}, "home page"},
		{KindPage, nil, []string{"about.md"}, "about page"},
		{KindSection, nil, []string{"sect3"}, "section 3"},
		{KindPage, nil, []string{"sect3/page1.md"}, "Title3_1"},
		{KindPage, nil, []string{"sect4/page2.md"}, "Title4_2"},
		{KindSection, nil, []string{"sect3/sect7"}, "another sect7"},
		{KindPage, nil, []string{"sect3/subsect/deep.md"}, "deep page"},
		{KindPage, nil, []string{filepath.FromSlash("sect5/page3.md")}, "Title5_3"}, //test OS-specific path

		// shorthand refs (potentially ambiguous)
		{KindPage, nil, []string{"unique.md"}, "UniqueBase"},
		{"Ambiguous", nil, []string{"page1.md"}, ""},

		// ISSUE: This is an ambiguous ref, but because we have to support the legacy
		// content root relative paths without a leading slash, the lookup
		// returns /sect7. This undermines ambiguity detection, but we have no choice.
		//{"Ambiguous", nil, []string{"sect7"}, ""},
		{KindSection, nil, []string{"sect7"}, "Sect7s"},

		// absolute paths
		{KindHome, nil, []string{"/"}, "home page"},
		{KindPage, nil, []string{"/about.md"}, "about page"},
		{KindSection, nil, []string{"/sect3"}, "section 3"},
		{KindPage, nil, []string{"/sect3/page1.md"}, "Title3_1"},
		{KindPage, nil, []string{"/sect4/page2.md"}, "Title4_2"},
		{KindSection, nil, []string{"/sect3/sect7"}, "another sect7"},
		{KindPage, nil, []string{"/sect3/subsect/deep.md"}, "deep page"},
		{KindPage, nil, []string{filepath.FromSlash("/sect5/page3.md")}, "Title5_3"}, //test OS-specific path
		{KindPage, nil, []string{"/sect3/unique.md"}, "UniqueBase"},                  //next test depends on this page existing
		// {"NoPage", nil, []string{"/unique.md"}, ""},  // ISSUE #4969: this is resolving to /sect3/unique.md
		{"NoPage", nil, []string{"/missing-page.md"}, ""},
		{"NoPage", nil, []string{"/missing-section"}, ""},

		// relative paths
		{KindHome, sec3, []string{".."}, "home page"},
		{KindHome, sec3, []string{"../"}, "home page"},
		{KindPage, sec3, []string{"../about.md"}, "about page"},
		{KindSection, sec3, []string{"."}, "section 3"},
		{KindSection, sec3, []string{"./"}, "section 3"},
		{KindPage, sec3, []string{"page1.md"}, "Title3_1"},
		{KindPage, sec3, []string{"./page1.md"}, "Title3_1"},
		{KindPage, sec3, []string{"../sect4/page2.md"}, "Title4_2"},
		{KindSection, sec3, []string{"sect7"}, "another sect7"},
		{KindSection, sec3, []string{"./sect7"}, "another sect7"},
		{KindPage, sec3, []string{"./subsect/deep.md"}, "deep page"},
		{KindPage, sec3, []string{"./subsect/../../sect7/page9.md"}, "Title7_9"},
		{KindPage, sec3, []string{filepath.FromSlash("../sect5/page3.md")}, "Title5_3"}, //test OS-specific path
		{KindPage, sec3, []string{"./unique.md"}, "UniqueBase"},
		{"NoPage", sec3, []string{"./sect2"}, ""},
		//{"NoPage", sec3, []string{"sect2"}, ""}, // ISSUE: /sect3 page relative query is resolving to /sect2

		// absolute paths ignore context
		{KindHome, sec3, []string{"/"}, "home page"},
		{KindPage, sec3, []string{"/about.md"}, "about page"},
		{KindPage, sec3, []string{"/sect4/page2.md"}, "Title4_2"},
		{KindPage, sec3, []string{"/sect3/subsect/deep.md"}, "deep page"}, //next test depends on this page existing
		{"NoPage", sec3, []string{"/subsect/deep.md"}, ""},
	}

	for _, test := range tests {
		errorMsg := fmt.Sprintf("Test case %s %v -> %s", test.context, test.path, test.expectedTitle)

		// test legacy public Site.GetPage (which does not support page context relative queries)
		if test.context == nil {
			args := append([]string{test.kind}, test.path...)
			page, err := s.Info.GetPage(args...)
			test.check(page, err, errorMsg, assert)
		}

		// test new internal Site.getPageNew
		var ref string
		if len(test.path) == 1 {
			ref = filepath.ToSlash(test.path[0])
		} else {
			ref = path.Join(test.path...)
		}
		page2, err := s.getPageNew(test.context, ref)
		test.check(page2, err, errorMsg, assert)
	}

}
