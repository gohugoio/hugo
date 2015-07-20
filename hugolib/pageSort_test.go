package hugolib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/source"
)

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
