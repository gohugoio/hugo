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
	"testing"
	"time"

	"github.com/spf13/hugo/source"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSort(t *testing.T) {

	d1 := time.Now()
	d2 := d1.Add(-1 * time.Hour)
	d3 := d1.Add(-2 * time.Hour)

	p := createSortTestPages(3)

	// first by weight
	setSortVals([3]time.Time{d1, d2, d3}, [3]string{"b", "a", "c"}, [3]int{3, 2, 1}, p)
	p.Sort()

	assert.Equal(t, 1, p[0].Weight)

	// next by date
	setSortVals([3]time.Time{d3, d1, d2}, [3]string{"a", "b", "c"}, [3]int{1, 1, 1}, p)
	p.Sort()
	assert.Equal(t, d1, p[0].Date)

	// finally by link title
	setSortVals([3]time.Time{d3, d3, d3}, [3]string{"b", "c", "a"}, [3]int{1, 1, 1}, p)
	p.Sort()
	assert.Equal(t, "al", p[0].LinkTitle())
	assert.Equal(t, "bl", p[1].LinkTitle())
	assert.Equal(t, "cl", p[2].LinkTitle())
}

func TestSortByN(t *testing.T) {

	d1 := time.Now()
	d2 := d1.Add(-2 * time.Hour)
	d3 := d1.Add(-10 * time.Hour)

	p := createSortTestPages(3)

	for i, this := range []struct {
		sortFunc   func(p Pages) Pages
		assertFunc func(p Pages) bool
	}{
		{(Pages).ByWeight, func(p Pages) bool { return p[0].Weight == 1 }},
		{(Pages).ByTitle, func(p Pages) bool { return p[0].Title == "ab" }},
		{(Pages).ByLinkTitle, func(p Pages) bool { return p[0].LinkTitle() == "abl" }},
		{(Pages).ByDate, func(p Pages) bool { return p[0].Date == d3 }},
		{(Pages).ByPublishDate, func(p Pages) bool { return p[0].PublishDate == d3 }},
		{(Pages).ByLastmod, func(p Pages) bool { return p[1].Lastmod == d2 }},
		{(Pages).ByLength, func(p Pages) bool { return p[0].Content == "b_content" }},
	} {
		setSortVals([3]time.Time{d1, d2, d3}, [3]string{"b", "ab", "cde"}, [3]int{3, 2, 1}, p)

		sorted := this.sortFunc(p)
		if !this.assertFunc(sorted) {
			t.Errorf("[%d] sort error", i)
		}
	}

}

func TestLimit(t *testing.T) {
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
	p1 := createSortTestPages(10)
	assert.Equal(t, 0, p1[0].FuzzyWordCount)
	assert.Equal(t, 9, p1[9].FuzzyWordCount)
	p2 := p1.Reverse()
	assert.Equal(t, 9, p2[0].FuzzyWordCount)
	assert.Equal(t, 0, p2[9].FuzzyWordCount)
	// cached
	assert.True(t, probablyEqualPages(p2, p1.Reverse()))
}

func BenchmarkSortByWeightAndReverse(b *testing.B) {

	p := createSortTestPages(300)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p = p.ByWeight().Reverse()
	}
}

func setSortVals(dates [3]time.Time, titles [3]string, weights [3]int, pages Pages) {
	for i := range dates {
		pages[i].Date = dates[i]
		pages[i].Lastmod = dates[i]
		pages[i].Weight = weights[i]
		pages[i].Title = titles[i]
		// make sure we compare apples and ... apples ...
		pages[len(dates)-1-i].linkTitle = pages[i].Title + "l"
		pages[len(dates)-1-i].PublishDate = dates[i]
		pages[len(dates)-1-i].Content = template.HTML(titles[i] + "_content")
	}
	lastLastMod := pages[2].Lastmod
	pages[2].Lastmod = pages[1].Lastmod
	pages[1].Lastmod = lastLastMod
}

func createSortTestPages(num int) Pages {
	pages := make(Pages, num)

	for i := 0; i < num; i++ {
		pages[i] = &Page{
			Node: Node{
				URLPath: URLPath{
					Section: "z",
					URL:     fmt.Sprintf("http://base/x/y/p%d.html", i),
				},
				Site: &SiteInfo{
					BaseURL: "http://base/",
				},
			},
			Source: Source{File: *source.NewFile(filepath.FromSlash(fmt.Sprintf("/x/y/p%d.md", i)))},
		}
		w := 5
		if i%2 == 0 {
			w = 10
		}
		pages[i].FuzzyWordCount = i
		pages[i].Weight = w
		pages[i].Description = "initial"
	}

	return pages
}
