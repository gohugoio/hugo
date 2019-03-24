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

	"github.com/gohugoio/hugo/output"
	"github.com/stretchr/testify/require"
)

func TestSplitPages(t *testing.T) {
	t.Parallel()

	pages := createTestPages(21)
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
	pages := createTestPages(21)
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
				require.True(t, p.FuzzyWordCount()%2 == 0) // magic test
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
				require.True(t, p.FuzzyWordCount()%2 != 0) // magic test
			}
		}
	} else {
		t.Fatal("Excepted PageGroup")
	}

}

func TestPager(t *testing.T) {
	t.Parallel()
	pages := createTestPages(21)
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

func doTestPages(t *testing.T, paginator *Paginator) {

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
	pages := createTestPages(0)
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

func doTestPagerNoPages(t *testing.T, paginator *Paginator) {
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
	cfg := viper.New()
	cfg.Set("paginatePath", "zoo")

	for _, uglyURLs := range []bool{false, true} {
		t.Run(fmt.Sprintf("uglyURLs=%t", uglyURLs), func(t *testing.T) {

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
					require.Equal(t, test.expectedUgly, got)
				} else {
					require.Equal(t, test.expected, got)
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

	require.Equal(t, 3, page11.FuzzyWordCount())
	require.Nil(t, page1Nil)

	require.NotNil(t, page21)
	require.Equal(t, 3, page21.FuzzyWordCount())
	require.Nil(t, page2Nil)
}
