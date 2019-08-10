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

package page

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/output"
)

func TestSplitPages(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	pages := createTestPages(21)
	chunks := splitPages(pages, 5)
	c.Assert(len(chunks), qt.Equals, 5)

	for i := 0; i < 4; i++ {
		c.Assert(chunks[i].Len(), qt.Equals, 5)
	}

	lastChunk := chunks[4]
	c.Assert(lastChunk.Len(), qt.Equals, 1)

}

func TestSplitPageGroups(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	pages := createTestPages(21)
	groups, _ := pages.GroupBy("Weight", "desc")
	chunks := splitPageGroups(groups, 5)
	c.Assert(len(chunks), qt.Equals, 5)

	firstChunk := chunks[0]

	// alternate weight 5 and 10
	if groups, ok := firstChunk.(PagesGroup); ok {
		c.Assert(groups.Len(), qt.Equals, 5)
		for _, pg := range groups {
			// first group 10 in weight
			c.Assert(pg.Key, qt.Equals, 10)
			for _, p := range pg.Pages {
				c.Assert(p.FuzzyWordCount()%2 == 0, qt.Equals, true) // magic test
			}
		}
	} else {
		t.Fatal("Excepted PageGroup")
	}

	lastChunk := chunks[4]

	if groups, ok := lastChunk.(PagesGroup); ok {
		c.Assert(groups.Len(), qt.Equals, 1)
		for _, pg := range groups {
			// last should have 5 in weight
			c.Assert(pg.Key, qt.Equals, 5)
			for _, p := range pg.Pages {
				c.Assert(p.FuzzyWordCount()%2 != 0, qt.Equals, true) // magic test
			}
		}
	} else {
		t.Fatal("Excepted PageGroup")
	}

}

func TestPager(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	pages := createTestPages(21)
	groups, _ := pages.GroupBy("Weight", "desc")

	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	_, err := newPaginatorFromPages(pages, -1, urlFactory)
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = newPaginatorFromPageGroups(groups, -1, urlFactory)
	c.Assert(err, qt.Not(qt.IsNil))

	pag, err := newPaginatorFromPages(pages, 5, urlFactory)
	c.Assert(err, qt.IsNil)
	doTestPages(t, pag)
	first := pag.Pagers()[0].First()
	c.Assert(first.String(), qt.Equals, "Pager 1")
	c.Assert(first.Pages(), qt.Not(qt.HasLen), 0)
	c.Assert(first.PageGroups(), qt.HasLen, 0)

	pag, err = newPaginatorFromPageGroups(groups, 5, urlFactory)
	c.Assert(err, qt.IsNil)
	doTestPages(t, pag)
	first = pag.Pagers()[0].First()
	c.Assert(first.PageGroups(), qt.Not(qt.HasLen), 0)
	c.Assert(first.Pages(), qt.HasLen, 0)

}

func doTestPages(t *testing.T, paginator *Paginator) {
	c := qt.New(t)
	paginatorPages := paginator.Pagers()

	c.Assert(len(paginatorPages), qt.Equals, 5)
	c.Assert(paginator.TotalNumberOfElements(), qt.Equals, 21)
	c.Assert(paginator.PageSize(), qt.Equals, 5)
	c.Assert(paginator.TotalPages(), qt.Equals, 5)

	first := paginatorPages[0]
	c.Assert(first.URL(), qt.Equals, template.HTML("page/1/"))
	c.Assert(first.First(), qt.Equals, first)
	c.Assert(first.HasNext(), qt.Equals, true)
	c.Assert(first.Next(), qt.Equals, paginatorPages[1])
	c.Assert(first.HasPrev(), qt.Equals, false)
	c.Assert(first.Prev(), qt.IsNil)
	c.Assert(first.NumberOfElements(), qt.Equals, 5)
	c.Assert(first.PageNumber(), qt.Equals, 1)

	third := paginatorPages[2]
	c.Assert(third.HasNext(), qt.Equals, true)
	c.Assert(third.HasPrev(), qt.Equals, true)
	c.Assert(third.Prev(), qt.Equals, paginatorPages[1])

	last := paginatorPages[4]
	c.Assert(last.URL(), qt.Equals, template.HTML("page/5/"))
	c.Assert(last.Last(), qt.Equals, last)
	c.Assert(last.HasNext(), qt.Equals, false)
	c.Assert(last.Next(), qt.IsNil)
	c.Assert(last.HasPrev(), qt.Equals, true)
	c.Assert(last.NumberOfElements(), qt.Equals, 1)
	c.Assert(last.PageNumber(), qt.Equals, 5)
}

func TestPagerNoPages(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	pages := createTestPages(0)
	groups, _ := pages.GroupBy("Weight", "desc")

	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	paginator, _ := newPaginatorFromPages(pages, 5, urlFactory)
	doTestPagerNoPages(t, paginator)

	first := paginator.Pagers()[0].First()
	c.Assert(first.PageGroups(), qt.HasLen, 0)
	c.Assert(first.Pages(), qt.HasLen, 0)

	paginator, _ = newPaginatorFromPageGroups(groups, 5, urlFactory)
	doTestPagerNoPages(t, paginator)

	first = paginator.Pagers()[0].First()
	c.Assert(first.PageGroups(), qt.HasLen, 0)
	c.Assert(first.Pages(), qt.HasLen, 0)

}

func doTestPagerNoPages(t *testing.T, paginator *Paginator) {
	paginatorPages := paginator.Pagers()
	c := qt.New(t)
	c.Assert(len(paginatorPages), qt.Equals, 1)
	c.Assert(paginator.TotalNumberOfElements(), qt.Equals, 0)
	c.Assert(paginator.PageSize(), qt.Equals, 5)
	c.Assert(paginator.TotalPages(), qt.Equals, 0)

	// pageOne should be nothing but the first
	pageOne := paginatorPages[0]
	c.Assert(pageOne.First(), qt.Not(qt.IsNil))
	c.Assert(pageOne.HasNext(), qt.Equals, false)
	c.Assert(pageOne.HasPrev(), qt.Equals, false)
	c.Assert(pageOne.Next(), qt.IsNil)
	c.Assert(len(pageOne.Pagers()), qt.Equals, 1)
	c.Assert(pageOne.Pages().Len(), qt.Equals, 0)
	c.Assert(pageOne.NumberOfElements(), qt.Equals, 0)
	c.Assert(pageOne.TotalNumberOfElements(), qt.Equals, 0)
	c.Assert(pageOne.TotalPages(), qt.Equals, 0)
	c.Assert(pageOne.PageNumber(), qt.Equals, 1)
	c.Assert(pageOne.PageSize(), qt.Equals, 5)

}

func TestPaginationURLFactory(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	cfg := viper.New()
	cfg.Set("paginatePath", "zoo")

	for _, uglyURLs := range []bool{false, true} {
		c.Run(fmt.Sprintf("uglyURLs=%t", uglyURLs), func(c *qt.C) {

			tests := []struct {
				name         string
				d            TargetPathDescriptor
				baseURL      string
				page         int
				expected     string
				expectedUgly string
			}{
				{"HTML home page 32",
					TargetPathDescriptor{Kind: KindHome, Type: output.HTMLFormat}, "http://example.com/", 32, "/zoo/32/", "/zoo/32.html"},
				{"JSON home page 42",
					TargetPathDescriptor{Kind: KindHome, Type: output.JSONFormat}, "http://example.com/", 42, "/zoo/42/index.json", "/zoo/42.json"},
			}

			for _, test := range tests {
				d := test.d
				cfg.Set("baseURL", test.baseURL)
				cfg.Set("uglyURLs", uglyURLs)
				d.UglyURLs = uglyURLs

				pathSpec := newTestPathSpecFor(cfg)
				d.PathSpec = pathSpec

				factory := newPaginationURLFactory(d)

				got := factory(test.page)

				if uglyURLs {
					c.Assert(got, qt.Equals, test.expectedUgly)
				} else {
					c.Assert(got, qt.Equals, test.expected)
				}

			}
		})

	}
}

func TestProbablyEqualPageLists(t *testing.T) {
	t.Parallel()
	fivePages := createTestPages(5)
	zeroPages := createTestPages(0)
	zeroPagesByWeight, _ := createTestPages(0).GroupBy("Weight", "asc")
	fivePagesByWeight, _ := createTestPages(5).GroupBy("Weight", "asc")
	ninePagesByWeight, _ := createTestPages(9).GroupBy("Weight", "asc")

	for i, this := range []struct {
		v1     interface{}
		v2     interface{}
		expect bool
	}{
		{nil, nil, true},
		{"a", "b", true},
		{"a", fivePages, false},
		{fivePages, "a", false},
		{fivePages, createTestPages(2), false},
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

func TestPaginationPage(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	fivePages := createTestPages(7)
	fivePagesFuzzyWordCount, _ := createTestPages(7).GroupBy("FuzzyWordCount", "asc")

	p1, _ := newPaginatorFromPages(fivePages, 2, urlFactory)
	p2, _ := newPaginatorFromPageGroups(fivePagesFuzzyWordCount, 2, urlFactory)

	f1 := p1.pagers[0].First()
	f2 := p2.pagers[0].First()

	page11, _ := f1.page(1)
	page1Nil, _ := f1.page(3)

	page21, _ := f2.page(1)
	page2Nil, _ := f2.page(3)

	c.Assert(page11.FuzzyWordCount(), qt.Equals, 3)
	c.Assert(page1Nil, qt.IsNil)

	c.Assert(page21, qt.Not(qt.IsNil))
	c.Assert(page21.FuzzyWordCount(), qt.Equals, 3)
	c.Assert(page2Nil, qt.IsNil)
}
