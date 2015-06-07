package hugolib

import (
	"fmt"
	"html/template"
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSplitPages(t *testing.T) {

	pages := createTestPages(21)
	chunks := splitPages(pages, 5)
	assert.Equal(t, 5, len(chunks))

	for i := 0; i < 4; i++ {
		assert.Equal(t, 5, len(chunks[i]))
	}

	lastChunk := chunks[4]
	assert.Equal(t, 1, len(lastChunk))

}

func TestPager(t *testing.T) {

	pages := createTestPages(21)
	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	_, err := newPaginator(pages, -1, urlFactory)
	assert.NotNil(t, err)

	paginator, _ := newPaginator(pages, 5, urlFactory)
	paginatorPages := paginator.Pagers()

	assert.Equal(t, 5, len(paginatorPages))
	assert.Equal(t, 21, paginator.TotalNumberOfElements())
	assert.Equal(t, 5, paginator.PageSize())
	assert.Equal(t, 5, paginator.TotalPages())

	first := paginatorPages[0]
	assert.Equal(t, template.HTML("page/1/"), first.URL())
	assert.Equal(t, first.URL(), first.Url())
	assert.Equal(t, first, first.First())
	assert.True(t, first.HasNext())
	assert.Equal(t, paginatorPages[1], first.Next())
	assert.False(t, first.HasPrev())
	assert.Nil(t, first.Prev())
	assert.Equal(t, 5, first.NumberOfElements())
	assert.Equal(t, 1, first.PageNumber())

	third := paginatorPages[2]
	assert.True(t, third.HasNext())
	assert.True(t, third.HasPrev())
	assert.Equal(t, paginatorPages[1], third.Prev())

	last := paginatorPages[4]
	assert.Equal(t, template.HTML("page/5/"), last.URL())
	assert.Equal(t, last, last.Last())
	assert.False(t, last.HasNext())
	assert.Nil(t, last.Next())
	assert.True(t, last.HasPrev())
	assert.Equal(t, 1, last.NumberOfElements())
	assert.Equal(t, 5, last.PageNumber())
}

func TestPagerNoPages(t *testing.T) {
	pages := createTestPages(0)
	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	paginator, _ := newPaginator(pages, 5, urlFactory)
	paginatorPages := paginator.Pagers()

	assert.Equal(t, 1, len(paginatorPages))
	assert.Equal(t, 0, paginator.TotalNumberOfElements())
	assert.Equal(t, 5, paginator.PageSize())
	assert.Equal(t, 0, paginator.TotalPages())

	// pageOne should be nothing but the first
	pageOne := paginatorPages[0]
	assert.NotNil(t, pageOne.First())
	assert.False(t, pageOne.HasNext())
	assert.False(t, pageOne.HasPrev())
	assert.Nil(t, pageOne.Next())
	assert.Equal(t, 1, len(pageOne.Pagers()))
	assert.Equal(t, 0, len(pageOne.Pages()))
	assert.Equal(t, 0, pageOne.NumberOfElements())
	assert.Equal(t, 0, pageOne.TotalNumberOfElements())
	assert.Equal(t, 0, pageOne.TotalPages())
	assert.Equal(t, 1, pageOne.PageNumber())
	assert.Equal(t, 5, pageOne.PageSize())

}

func TestPaginationURLFactory(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("PaginatePath", "zoo")
	unicode := newPaginationURLFactory("новости проекта")
	fooBar := newPaginationURLFactory("foo", "bar")

	assert.Equal(t, "/%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D0%B8-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0/", unicode(1))
	assert.Equal(t, "/foo/bar/", fooBar(1))
	assert.Equal(t, "/%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D0%B8-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0/zoo/4/", unicode(4))
	assert.Equal(t, "/foo/bar/zoo/12345/", fooBar(12345))

}

func TestPaginator(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	for _, useViper := range []bool{false, true} {
		doTestPaginator(t, useViper)
	}
}

func doTestPaginator(t *testing.T, useViper bool) {
	viper.Reset()
	defer viper.Reset()

	pagerSize := 5
	if useViper {
		viper.Set("paginate", pagerSize)
	} else {
		viper.Set("paginate", -1)
	}
	pages := createTestPages(12)
	s := &Site{}
	n1 := s.newHomeNode()
	n2 := s.newHomeNode()
	n1.Data["Pages"] = pages

	var paginator1 *Pager
	var err error

	if useViper {
		paginator1, err = n1.Paginator()
	} else {
		paginator1, err = n1.Paginator(pagerSize)
	}

	assert.Nil(t, err)
	assert.NotNil(t, paginator1)
	assert.Equal(t, 3, paginator1.TotalPages())
	assert.Equal(t, 12, paginator1.TotalNumberOfElements())

	n2.paginator = paginator1.Next()
	paginator2, err := n2.Paginator()
	assert.Nil(t, err)
	assert.Equal(t, paginator2, paginator1.Next())

	n1.Data["Pages"] = createTestPages(1)
	samePaginator, _ := n1.Paginator()
	assert.Equal(t, paginator1, samePaginator)

	p, _ := NewPage("test")
	_, err = p.Paginator()
	assert.NotNil(t, err)
}

func TestPaginatorWithNegativePaginate(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("paginate", -1)
	s := &Site{}
	_, err := s.newHomeNode().Paginator()
	assert.NotNil(t, err)
}

func TestPaginate(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	for _, useViper := range []bool{false, true} {
		doTestPaginate(t, useViper)
	}
}

func doTestPaginate(t *testing.T, useViper bool) {
	pagerSize := 5
	if useViper {
		viper.Set("paginate", pagerSize)
	} else {
		viper.Set("paginate", -1)
	}

	pages := createTestPages(6)
	s := &Site{}
	n1 := s.newHomeNode()
	n2 := s.newHomeNode()

	var paginator1, paginator2 *Pager
	var err error

	if useViper {
		paginator1, err = n1.Paginate(pages)
	} else {
		paginator1, err = n1.Paginate(pages, pagerSize)
	}

	assert.Nil(t, err)
	assert.NotNil(t, paginator1)
	assert.Equal(t, 2, paginator1.TotalPages())
	assert.Equal(t, 6, paginator1.TotalNumberOfElements())

	n2.paginator = paginator1.Next()
	if useViper {
		paginator2, err = n2.Paginate(pages)
	} else {
		paginator2, err = n2.Paginate(pages, pagerSize)
	}
	assert.Nil(t, err)
	assert.Equal(t, paginator2, paginator1.Next())

	p, _ := NewPage("test")
	_, err = p.Paginate(pages)
	assert.NotNil(t, err)
}

func TestInvalidOptions(t *testing.T) {
	s := &Site{}
	n1 := s.newHomeNode()
	_, err := n1.Paginate(createTestPages(1), 1, 2)
	assert.NotNil(t, err)
	_, err = n1.Paginator(1, 2)
	assert.NotNil(t, err)
	_, err = n1.Paginator(-1)
	assert.NotNil(t, err)
}

func TestPaginateWithNegativePaginate(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("paginate", -1)
	s := &Site{}
	_, err := s.newHomeNode().Paginate(createTestPages(2))
	assert.NotNil(t, err)
}

func TestPaginatePages(t *testing.T) {
	for i, seq := range []interface{}{createTestPages(11), WeightedPages{}, PageGroup{}, &Pages{}} {
		v, err := paginatePages(seq, 11, "t")
		assert.NotNil(t, v, "Val %d", i)
		assert.Nil(t, err, "Err %d", i)
	}
	_, err := paginatePages(Site{}, 11, "t")
	assert.NotNil(t, err)

}

// Issue #993
func TestPaginatorFollowedByPaginateShouldFail(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("paginate", 10)
	s := &Site{}
	n1 := s.newHomeNode()
	n2 := s.newHomeNode()

	_, err := n1.Paginator()
	assert.Nil(t, err)
	_, err = n1.Paginate(createTestPages(2))
	assert.NotNil(t, err)

	_, err = n2.Paginate(createTestPages(2))
	assert.Nil(t, err)

}

func TestPaginateFollowedByDifferentPaginateShouldFail(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("paginate", 10)
	s := &Site{}
	n1 := s.newHomeNode()
	n2 := s.newHomeNode()

	p1 := createTestPages(2)
	p2 := createTestPages(10)

	_, err := n1.Paginate(p1)
	assert.Nil(t, err)

	_, err = n1.Paginate(p1)
	assert.Nil(t, err)

	_, err = n1.Paginate(p2)
	assert.NotNil(t, err)

	_, err = n2.Paginate(p2)
	assert.Nil(t, err)
}

func TestProbablyEqualPageLists(t *testing.T) {
	fivePages := createTestPages(5)
	zeroPages := createTestPages(0)
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
	} {
		result := probablyEqualPageLists(this.v1, this.v2)

		if result != this.expect {
			t.Errorf("[%d] got %t but expected %t", i, result, this.expect)

		}
	}
}

func createTestPages(num int) Pages {
	pages := make(Pages, num)

	for i := 0; i < num; i++ {
		pages[i] = &Page{
			Node: Node{
				URLPath: URLPath{
					Section: "z",
					URL:     fmt.Sprintf("http://base/x/y/p%d.html", num),
				},
				Site: &SiteInfo{
					BaseURL: "http://base/",
				},
			},
			Source: Source{File: *source.NewFile(filepath.FromSlash(fmt.Sprintf("/x/y/p%d.md", num)))},
		}
	}

	return pages
}
