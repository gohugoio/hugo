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
	"html/template"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/output"
	"github.com/stretchr/testify/require"
)

func TestSplitPages(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)

	pages := createTestPages(s, 21)
	chunks := splitPages(pages, 5)
	require.Equal(t, 5, len(chunks))

	for i := 0; i < 4; i++ {
		require.Equal(t, 5, chunks[i].Len())
	}

	lastChunk := chunks[4]
	require.Equal(t, 1, lastChunk.Len())

}

func TestSplitPageGroups(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	pages := createTestPages(s, 21)
	groups, _ := pages.GroupBy("Weight", "desc")
	chunks := splitPageGroups(groups, 5)
	require.Equal(t, 5, len(chunks))

	firstChunk := chunks[0]

	// alternate weight 5 and 10
	if groups, ok := firstChunk.(PagesGroup); ok {
		require.Equal(t, 5, groups.Len())
		for _, pg := range groups {
			// first group 10 in weight
			require.Equal(t, 10, pg.Key)
			for _, p := range pg.Pages {
				require.True(t, p.fuzzyWordCount%2 == 0) // magic test
			}
		}
	} else {
		t.Fatal("Excepted PageGroup")
	}

	lastChunk := chunks[4]

	if groups, ok := lastChunk.(PagesGroup); ok {
		require.Equal(t, 1, groups.Len())
		for _, pg := range groups {
			// last should have 5 in weight
			require.Equal(t, 5, pg.Key)
			for _, p := range pg.Pages {
				require.True(t, p.fuzzyWordCount%2 != 0) // magic test
			}
		}
	} else {
		t.Fatal("Excepted PageGroup")
	}

}

func TestPager(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	pages := createTestPages(s, 21)
	groups, _ := pages.GroupBy("Weight", "desc")

	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	_, err := newPaginatorFromPages(pages, -1, urlFactory)
	require.NotNil(t, err)

	_, err = newPaginatorFromPageGroups(groups, -1, urlFactory)
	require.NotNil(t, err)

	pag, err := newPaginatorFromPages(pages, 5, urlFactory)
	require.Nil(t, err)
	doTestPages(t, pag)
	first := pag.Pagers()[0].First()
	require.Equal(t, "Pager 1", first.String())
	require.NotEmpty(t, first.Pages())
	require.Empty(t, first.PageGroups())

	pag, err = newPaginatorFromPageGroups(groups, 5, urlFactory)
	require.Nil(t, err)
	doTestPages(t, pag)
	first = pag.Pagers()[0].First()
	require.NotEmpty(t, first.PageGroups())
	require.Empty(t, first.Pages())

}

func doTestPages(t *testing.T, paginator *paginator) {

	paginatorPages := paginator.Pagers()

	require.Equal(t, 5, len(paginatorPages))
	require.Equal(t, 21, paginator.TotalNumberOfElements())
	require.Equal(t, 5, paginator.PageSize())
	require.Equal(t, 5, paginator.TotalPages())

	first := paginatorPages[0]
	require.Equal(t, template.HTML("page/1/"), first.URL())
	require.Equal(t, first, first.First())
	require.True(t, first.HasNext())
	require.Equal(t, paginatorPages[1], first.Next())
	require.False(t, first.HasPrev())
	require.Nil(t, first.Prev())
	require.Equal(t, 5, first.NumberOfElements())
	require.Equal(t, 1, first.PageNumber())

	third := paginatorPages[2]
	require.True(t, third.HasNext())
	require.True(t, third.HasPrev())
	require.Equal(t, paginatorPages[1], third.Prev())

	last := paginatorPages[4]
	require.Equal(t, template.HTML("page/5/"), last.URL())
	require.Equal(t, last, last.Last())
	require.False(t, last.HasNext())
	require.Nil(t, last.Next())
	require.True(t, last.HasPrev())
	require.Equal(t, 1, last.NumberOfElements())
	require.Equal(t, 5, last.PageNumber())
}

func TestPagerNoPages(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	pages := createTestPages(s, 0)
	groups, _ := pages.GroupBy("Weight", "desc")

	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	paginator, _ := newPaginatorFromPages(pages, 5, urlFactory)
	doTestPagerNoPages(t, paginator)

	first := paginator.Pagers()[0].First()
	require.Empty(t, first.PageGroups())
	require.Empty(t, first.Pages())

	paginator, _ = newPaginatorFromPageGroups(groups, 5, urlFactory)
	doTestPagerNoPages(t, paginator)

	first = paginator.Pagers()[0].First()
	require.Empty(t, first.PageGroups())
	require.Empty(t, first.Pages())

}

func doTestPagerNoPages(t *testing.T, paginator *paginator) {
	paginatorPages := paginator.Pagers()

	require.Equal(t, 1, len(paginatorPages))
	require.Equal(t, 0, paginator.TotalNumberOfElements())
	require.Equal(t, 5, paginator.PageSize())
	require.Equal(t, 0, paginator.TotalPages())

	// pageOne should be nothing but the first
	pageOne := paginatorPages[0]
	require.NotNil(t, pageOne.First())
	require.False(t, pageOne.HasNext())
	require.False(t, pageOne.HasPrev())
	require.Nil(t, pageOne.Next())
	require.Equal(t, 1, len(pageOne.Pagers()))
	require.Equal(t, 0, pageOne.Pages().Len())
	require.Equal(t, 0, pageOne.NumberOfElements())
	require.Equal(t, 0, pageOne.TotalNumberOfElements())
	require.Equal(t, 0, pageOne.TotalPages())
	require.Equal(t, 1, pageOne.PageNumber())
	require.Equal(t, 5, pageOne.PageSize())

}

func TestPaginationURLFactory(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()
	cfg.Set("paginatePath", "zoo")

	for _, uglyURLs := range []bool{false, true} {
		for _, canonifyURLs := range []bool{false, true} {
			t.Run(fmt.Sprintf("uglyURLs=%t,canonifyURLs=%t", uglyURLs, canonifyURLs), func(t *testing.T) {

				tests := []struct {
					name     string
					d        targetPathDescriptor
					baseURL  string
					page     int
					expected string
				}{
					{"HTML home page 32",
						targetPathDescriptor{Kind: KindHome, Type: output.HTMLFormat}, "http://example.com/", 32, "/zoo/32/"},
					{"JSON home page 42",
						targetPathDescriptor{Kind: KindHome, Type: output.JSONFormat}, "http://example.com/", 42, "/zoo/42/"},
					// Issue #1252
					{"BaseURL with sub path",
						targetPathDescriptor{Kind: KindHome, Type: output.HTMLFormat}, "http://example.com/sub/", 999, "/sub/zoo/999/"},
				}

				for _, test := range tests {
					d := test.d
					cfg.Set("baseURL", test.baseURL)
					cfg.Set("canonifyURLs", canonifyURLs)
					cfg.Set("uglyURLs", uglyURLs)
					d.UglyURLs = uglyURLs

					expected := test.expected

					if canonifyURLs {
						expected = strings.Replace(expected, "/sub", "", 1)
					}

					if uglyURLs {
						expected = expected[:len(expected)-1] + "." + test.d.Type.MediaType.Suffix
					}

					pathSpec := newTestPathSpec(fs, cfg)
					d.PathSpec = pathSpec

					factory := newPaginationURLFactory(d)

					got := factory(test.page)

					require.Equal(t, expected, got)

				}
			})
		}
	}
}

func TestPaginator(t *testing.T) {
	t.Parallel()
	for _, useViper := range []bool{false, true} {
		doTestPaginator(t, useViper)
	}
}

func doTestPaginator(t *testing.T, useViper bool) {

	cfg, fs := newTestCfg()

	pagerSize := 5
	if useViper {
		cfg.Set("paginate", pagerSize)
	} else {
		cfg.Set("paginate", -1)
	}

	s, err := NewSiteForCfg(deps.DepsCfg{Cfg: cfg, Fs: fs})
	require.NoError(t, err)

	pages := createTestPages(s, 12)
	n1, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)
	n2, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)
	n1.Data["Pages"] = pages

	var paginator1 *Pager

	if useViper {
		paginator1, err = n1.Paginator()
	} else {
		paginator1, err = n1.Paginator(pagerSize)
	}

	require.Nil(t, err)
	require.NotNil(t, paginator1)
	require.Equal(t, 3, paginator1.TotalPages())
	require.Equal(t, 12, paginator1.TotalNumberOfElements())

	n2.paginator = paginator1.Next()
	paginator2, err := n2.Paginator()
	require.Nil(t, err)
	require.Equal(t, paginator2, paginator1.Next())

	n1.Data["Pages"] = createTestPages(s, 1)
	samePaginator, _ := n1.Paginator()
	require.Equal(t, paginator1, samePaginator)

	pp, _ := s.NewPage("test")
	p, _ := newPageOutput(pp, false, output.HTMLFormat)

	_, err = p.Paginator()
	require.NotNil(t, err)
}

func TestPaginatorWithNegativePaginate(t *testing.T) {
	t.Parallel()
	s := newTestSite(t, "paginate", -1)
	n1, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)
	_, err := n1.Paginator()
	require.Error(t, err)
}

func TestPaginate(t *testing.T) {
	t.Parallel()
	for _, useViper := range []bool{false, true} {
		doTestPaginate(t, useViper)
	}
}

func TestPaginatorURL(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	cfg.Set("paginate", 2)
	cfg.Set("paginatePath", "testing")

	for i := 0; i < 10; i++ {
		// Issue #2177, do not double encode URLs
		writeSource(t, fs, filepath.Join("content", "阅读", fmt.Sprintf("page%d.md", (i+1))),
			fmt.Sprintf(`---
title: Page%d
---
Conten%d
`, (i+1), i+1))

	}
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), "<html><body>{{.Content}}</body></html>")
	writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"),
		`
<html><body>
Count: {{ .Paginator.TotalNumberOfElements }}
Pages: {{ .Paginator.TotalPages }}
{{ range .Paginator.Pagers -}}
 {{ .PageNumber }}: {{ .URL }} 
{{ end }}
</body></html>`)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	th := testHelper{s.Cfg, s.Fs, t}

	th.assertFileContent(filepath.Join("public", "阅读", "testing", "2", "index.html"), "2: /%E9%98%85%E8%AF%BB/testing/2/")

}

func doTestPaginate(t *testing.T, useViper bool) {
	pagerSize := 5

	var (
		s   *Site
		err error
	)

	if useViper {
		s = newTestSite(t, "paginate", pagerSize)
	} else {
		s = newTestSite(t, "paginate", -1)
	}

	pages := createTestPages(s, 6)
	n1, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)
	n2, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)

	var paginator1, paginator2 *Pager

	if useViper {
		paginator1, err = n1.Paginate(pages)
	} else {
		paginator1, err = n1.Paginate(pages, pagerSize)
	}

	require.Nil(t, err)
	require.NotNil(t, paginator1)
	require.Equal(t, 2, paginator1.TotalPages())
	require.Equal(t, 6, paginator1.TotalNumberOfElements())

	n2.paginator = paginator1.Next()
	if useViper {
		paginator2, err = n2.Paginate(pages)
	} else {
		paginator2, err = n2.Paginate(pages, pagerSize)
	}
	require.Nil(t, err)
	require.Equal(t, paginator2, paginator1.Next())

	pp, err := s.NewPage("test")
	p, _ := newPageOutput(pp, false, output.HTMLFormat)

	_, err = p.Paginate(pages)
	require.NotNil(t, err)
}

func TestInvalidOptions(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	n1, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)

	_, err := n1.Paginate(createTestPages(s, 1), 1, 2)
	require.NotNil(t, err)
	_, err = n1.Paginator(1, 2)
	require.NotNil(t, err)
	_, err = n1.Paginator(-1)
	require.NotNil(t, err)
}

func TestPaginateWithNegativePaginate(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()
	cfg.Set("paginate", -1)

	s, err := NewSiteForCfg(deps.DepsCfg{Cfg: cfg, Fs: fs})
	require.NoError(t, err)

	n, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)

	_, err = n.Paginate(createTestPages(s, 2))
	require.NotNil(t, err)
}

func TestPaginatePages(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)

	groups, _ := createTestPages(s, 31).GroupBy("Weight", "desc")
	pd := targetPathDescriptor{Kind: KindHome, Type: output.HTMLFormat, PathSpec: s.PathSpec, Addends: "t"}

	for i, seq := range []interface{}{createTestPages(s, 11), groups, WeightedPages{}, PageGroup{}, &Pages{}} {
		v, err := paginatePages(pd, seq, 11)
		require.NotNil(t, v, "Val %d", i)
		require.Nil(t, err, "Err %d", i)
	}
	_, err := paginatePages(pd, Site{}, 11)
	require.NotNil(t, err)

}

// Issue #993
func TestPaginatorFollowedByPaginateShouldFail(t *testing.T) {
	t.Parallel()
	s := newTestSite(t, "paginate", 10)
	n1, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)
	n2, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)

	_, err := n1.Paginator()
	require.Nil(t, err)
	_, err = n1.Paginate(createTestPages(s, 2))
	require.NotNil(t, err)

	_, err = n2.Paginate(createTestPages(s, 2))
	require.Nil(t, err)

}

func TestPaginateFollowedByDifferentPaginateShouldFail(t *testing.T) {
	t.Parallel()
	s := newTestSite(t, "paginate", 10)

	n1, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)
	n2, _ := newPageOutput(s.newHomePage(), false, output.HTMLFormat)

	p1 := createTestPages(s, 2)
	p2 := createTestPages(s, 10)

	_, err := n1.Paginate(p1)
	require.Nil(t, err)

	_, err = n1.Paginate(p1)
	require.Nil(t, err)

	_, err = n1.Paginate(p2)
	require.NotNil(t, err)

	_, err = n2.Paginate(p2)
	require.Nil(t, err)
}

func TestProbablyEqualPageLists(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	fivePages := createTestPages(s, 5)
	zeroPages := createTestPages(s, 0)
	zeroPagesByWeight, _ := createTestPages(s, 0).GroupBy("Weight", "asc")
	fivePagesByWeight, _ := createTestPages(s, 5).GroupBy("Weight", "asc")
	ninePagesByWeight, _ := createTestPages(s, 9).GroupBy("Weight", "asc")

	for i, this := range []struct {
		v1     interface{}
		v2     interface{}
		expect bool
	}{
		{nil, nil, true},
		{"a", "b", true},
		{"a", fivePages, false},
		{fivePages, "a", false},
		{fivePages, createTestPages(s, 2), false},
		{fivePages, fivePages, true},
		{zeroPages, zeroPages, true},
		{fivePagesByWeight, fivePagesByWeight, true},
		{zeroPagesByWeight, fivePagesByWeight, false},
		{zeroPagesByWeight, zeroPagesByWeight, true},
		{fivePagesByWeight, fivePages, false},
		{fivePagesByWeight, ninePagesByWeight, false},
	} {
		result := probablyEqualPageLists(this.v1, this.v2)

		if result != this.expect {
			t.Errorf("[%d] got %t but expected %t", i, result, this.expect)

		}
	}
}

func TestPage(t *testing.T) {
	t.Parallel()
	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	s := newTestSite(t)

	fivePages := createTestPages(s, 7)
	fivePagesFuzzyWordCount, _ := createTestPages(s, 7).GroupBy("FuzzyWordCount", "asc")

	p1, _ := newPaginatorFromPages(fivePages, 2, urlFactory)
	p2, _ := newPaginatorFromPageGroups(fivePagesFuzzyWordCount, 2, urlFactory)

	f1 := p1.pagers[0].First()
	f2 := p2.pagers[0].First()

	page11, _ := f1.page(1)
	page1Nil, _ := f1.page(3)

	page21, _ := f2.page(1)
	page2Nil, _ := f2.page(3)

	require.Equal(t, 3, page11.fuzzyWordCount)
	require.Nil(t, page1Nil)

	require.Equal(t, 3, page21.fuzzyWordCount)
	require.Nil(t, page2Nil)
}

func createTestPages(s *Site, num int) Pages {
	pages := make(Pages, num)

	for i := 0; i < num; i++ {
		p := s.newPage(filepath.FromSlash(fmt.Sprintf("/x/y/z/p%d.md", i)))
		w := 5
		if i%2 == 0 {
			w = 10
		}
		p.fuzzyWordCount = i + 2
		p.Weight = w
		pages[i] = p

	}

	return pages
}
