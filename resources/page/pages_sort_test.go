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
	"testing"
	"time"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSort(t *testing.T) {
	t.Parallel()
	d1 := time.Now()
	d2 := d1.Add(-1 * time.Hour)
	d3 := d1.Add(-2 * time.Hour)
	d4 := d1.Add(-3 * time.Hour)

	p := createSortTestPages(4)

	// first by weight
	setSortVals([4]time.Time{d1, d2, d3, d4}, [4]string{"b", "a", "c", "d"}, [4]int{4, 3, 2, 1}, p)
	SortByDefault(p)

	assert.Equal(t, 1, p[0].Weight())

	// Consider zero weight, issue #2673
	setSortVals([4]time.Time{d1, d2, d3, d4}, [4]string{"b", "a", "d", "c"}, [4]int{0, 0, 0, 1}, p)
	SortByDefault(p)

	assert.Equal(t, 1, p[0].Weight())

	// next by date
	setSortVals([4]time.Time{d3, d4, d1, d2}, [4]string{"a", "b", "c", "d"}, [4]int{1, 1, 1, 1}, p)
	SortByDefault(p)
	assert.Equal(t, d1, p[0].Date())

	// finally by link title
	setSortVals([4]time.Time{d3, d3, d3, d3}, [4]string{"b", "c", "a", "d"}, [4]int{1, 1, 1, 1}, p)
	SortByDefault(p)
	assert.Equal(t, "al", p[0].LinkTitle())
	assert.Equal(t, "bl", p[1].LinkTitle())
	assert.Equal(t, "cl", p[2].LinkTitle())
}

// https://github.com/gohugoio/hugo/issues/4953
func TestSortByLinkTitle(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	pages := createSortTestPages(6)

	for i, p := range pages {
		pp := p.(*testPage)
		if i < 5 {
			pp.title = fmt.Sprintf("title%d", i)
		}

		if i > 2 {
			pp.linkTitle = fmt.Sprintf("linkTitle%d", i)
		}

	}

	pages.shuffle()

	bylt := pages.ByLinkTitle()

	for i, p := range bylt {
		msg := fmt.Sprintf("test: %d", i)
		if i < 3 {
			assert.Equal(fmt.Sprintf("linkTitle%d", i+3), p.LinkTitle(), msg)
		} else {
			assert.Equal(fmt.Sprintf("title%d", i-3), p.LinkTitle(), msg)
		}
	}
}

func TestSortByN(t *testing.T) {
	t.Parallel()
	d1 := time.Now()
	d2 := d1.Add(-2 * time.Hour)
	d3 := d1.Add(-10 * time.Hour)
	d4 := d1.Add(-20 * time.Hour)

	p := createSortTestPages(4)

	for i, this := range []struct {
		sortFunc   func(p Pages) Pages
		assertFunc func(p Pages) bool
	}{
		{(Pages).ByWeight, func(p Pages) bool { return p[0].Weight() == 1 }},
		{(Pages).ByTitle, func(p Pages) bool { return p[0].Title() == "ab" }},
		{(Pages).ByLinkTitle, func(p Pages) bool { return p[0].LinkTitle() == "abl" }},
		{(Pages).ByDate, func(p Pages) bool { return p[0].Date() == d4 }},
		{(Pages).ByPublishDate, func(p Pages) bool { return p[0].PublishDate() == d4 }},
		{(Pages).ByExpiryDate, func(p Pages) bool { return p[0].ExpiryDate() == d4 }},
		{(Pages).ByLastmod, func(p Pages) bool { return p[1].Lastmod() == d3 }},
		{(Pages).ByLength, func(p Pages) bool { return p[0].(resource.LengthProvider).Len() == len(p[0].(*testPage).content) }},
	} {
		setSortVals([4]time.Time{d1, d2, d3, d4}, [4]string{"b", "ab", "cde", "fg"}, [4]int{0, 3, 2, 1}, p)

		sorted := this.sortFunc(p)
		if !this.assertFunc(sorted) {
			t.Errorf("[%d] sort error", i)
		}
	}

}

func TestLimit(t *testing.T) {
	t.Parallel()
	p := createSortTestPages(10)
	firstFive := p.Limit(5)
	assert.Equal(t, 5, len(firstFive))
	for i := 0; i < 5; i++ {
		assert.Equal(t, p[i], firstFive[i])
	}
	assert.Equal(t, p, p.Limit(10))
	assert.Equal(t, p, p.Limit(11))
}

func TestPageSortReverse(t *testing.T) {
	t.Parallel()
	p1 := createSortTestPages(10)
	assert.Equal(t, 0, p1[0].(*testPage).fuzzyWordCount)
	assert.Equal(t, 9, p1[9].(*testPage).fuzzyWordCount)
	p2 := p1.Reverse()
	assert.Equal(t, 9, p2[0].(*testPage).fuzzyWordCount)
	assert.Equal(t, 0, p2[9].(*testPage).fuzzyWordCount)
	// cached
	assert.True(t, pagesEqual(p2, p1.Reverse()))
}

func TestPageSortByParam(t *testing.T) {
	t.Parallel()
	var k interface{} = "arbitrarily.nested"

	unsorted := createSortTestPages(10)
	delete(unsorted[9].Params(), "arbitrarily")

	firstSetValue, _ := unsorted[0].Param(k)
	secondSetValue, _ := unsorted[1].Param(k)
	lastSetValue, _ := unsorted[8].Param(k)
	unsetValue, _ := unsorted[9].Param(k)

	assert.Equal(t, "xyz100", firstSetValue)
	assert.Equal(t, "xyz99", secondSetValue)
	assert.Equal(t, "xyz92", lastSetValue)
	assert.Equal(t, nil, unsetValue)

	sorted := unsorted.ByParam("arbitrarily.nested")
	firstSetSortedValue, _ := sorted[0].Param(k)
	secondSetSortedValue, _ := sorted[1].Param(k)
	lastSetSortedValue, _ := sorted[8].Param(k)
	unsetSortedValue, _ := sorted[9].Param(k)

	assert.Equal(t, firstSetValue, firstSetSortedValue)
	assert.Equal(t, secondSetValue, lastSetSortedValue)
	assert.Equal(t, lastSetValue, secondSetSortedValue)
	assert.Equal(t, unsetValue, unsetSortedValue)
}

func TestPageSortByParamNumeric(t *testing.T) {
	t.Parallel()
	var k interface{} = "arbitrarily.nested"

	n := 10
	unsorted := createSortTestPages(n)
	for i := 0; i < n; i++ {
		v := 100 - i
		if i%2 == 0 {
			v = 100.0 - i
		}

		unsorted[i].(*testPage).params = map[string]interface{}{
			"arbitrarily": map[string]interface{}{
				"nested": v,
			},
		}
	}
	delete(unsorted[9].Params(), "arbitrarily")

	firstSetValue, _ := unsorted[0].Param(k)
	secondSetValue, _ := unsorted[1].Param(k)
	lastSetValue, _ := unsorted[8].Param(k)
	unsetValue, _ := unsorted[9].Param(k)

	assert.Equal(t, 100, firstSetValue)
	assert.Equal(t, 99, secondSetValue)
	assert.Equal(t, 92, lastSetValue)
	assert.Equal(t, nil, unsetValue)

	sorted := unsorted.ByParam("arbitrarily.nested")
	firstSetSortedValue, _ := sorted[0].Param(k)
	secondSetSortedValue, _ := sorted[1].Param(k)
	lastSetSortedValue, _ := sorted[8].Param(k)
	unsetSortedValue, _ := sorted[9].Param(k)

	assert.Equal(t, 92, firstSetSortedValue)
	assert.Equal(t, 93, secondSetSortedValue)
	assert.Equal(t, 100, lastSetSortedValue)
	assert.Equal(t, unsetValue, unsetSortedValue)
}

func BenchmarkSortByWeightAndReverse(b *testing.B) {
	p := createSortTestPages(300)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p = p.ByWeight().Reverse()
	}
}

func setSortVals(dates [4]time.Time, titles [4]string, weights [4]int, pages Pages) {
	for i := range dates {
		this := pages[i].(*testPage)
		other := pages[len(dates)-1-i].(*testPage)

		this.date = dates[i]
		this.lastMod = dates[i]
		this.weight = weights[i]
		this.title = titles[i]
		// make sure we compare apples and ... apples ...
		other.linkTitle = this.Title() + "l"
		other.pubDate = dates[i]
		other.expiryDate = dates[i]
		other.content = titles[i] + "_content"
	}
	lastLastMod := pages[2].Lastmod()
	pages[2].(*testPage).lastMod = pages[1].Lastmod()
	pages[1].(*testPage).lastMod = lastLastMod

	for _, p := range pages {
		p.(*testPage).content = ""
	}

}

func createSortTestPages(num int) Pages {
	pages := make(Pages, num)

	for i := 0; i < num; i++ {
		p := newTestPage()
		p.path = fmt.Sprintf("/x/y/p%d.md", i)
		p.params = map[string]interface{}{
			"arbitrarily": map[string]interface{}{
				"nested": ("xyz" + fmt.Sprintf("%v", 100-i)),
			},
		}

		w := 5

		if i%2 == 0 {
			w = 10
		}
		p.fuzzyWordCount = i
		p.weight = w
		p.description = "initial"

		pages[i] = p
	}

	return pages
}
