package hugolib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/hugo/source"
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

	// finally by title
	setSortVals([3]time.Time{d3, d3, d3}, [3]string{"b", "c", "a"}, [3]int{1, 1, 1}, p)
	p.Sort()
	assert.Equal(t, "a", p[0].Title)
	assert.Equal(t, "b", p[1].Title)
	assert.Equal(t, "c", p[2].Title)
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
		pages[i].Weight = weights[i]
		pages[i].Title = titles[i]
	}

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
