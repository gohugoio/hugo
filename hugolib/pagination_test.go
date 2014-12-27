package hugolib

import (
	"fmt"
	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
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

func TestPaginator(t *testing.T) {

	pages := createTestPages(21)
	urlFactory := func(page int) string {
		return fmt.Sprintf("page/%d/", page)
	}

	paginator := newPaginator(pages, 5, urlFactory)
	paginatorPages := paginator.Pagers()

	assert.Equal(t, 5, len(paginatorPages))
	assert.Equal(t, 21, paginator.TotalNumberOfElements())
	assert.Equal(t, 5, paginator.PageSize())
	assert.Equal(t, 5, paginator.TotalPages())

	first := paginatorPages[0]
	assert.Equal(t, "page/1/", first.Url())
	assert.Equal(t, first, first.First())
	assert.Equal(t, true, first.HasNext())
	assert.Equal(t, false, first.HasPrev())
	assert.Equal(t, 5, first.NumberOfElements())
	assert.Equal(t, 1, first.PageNumber())

	third := paginatorPages[2]
	assert.Equal(t, true, third.HasNext())
	assert.Equal(t, true, third.HasPrev())

	last := paginatorPages[4]
	assert.Equal(t, "page/5/", last.Url())
	assert.Equal(t, last, last.Last())
	assert.Equal(t, false, last.HasNext())
	assert.Equal(t, true, last.HasPrev())
	assert.Equal(t, 1, last.NumberOfElements())
	assert.Equal(t, 5, last.PageNumber())

}

func TestPaginationUrlFactory(t *testing.T) {
	viper.Set("PaginatePath", "zoo")
	unicode := newPaginationUrlFactory("новости проекта")
	fooBar := newPaginationUrlFactory("foo", "bar")

	assert.Equal(t, "/%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D0%B8-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0/", unicode(1))
	assert.Equal(t, "/foo/bar/", fooBar(1))
	assert.Equal(t, "/%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D0%B8-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0/zoo/4/", unicode(4))
	assert.Equal(t, "/foo/bar/zoo/12345/", fooBar(12345))

}

func createTestPages(num int) Pages {
	pages := make(Pages, num)

	for i := 0; i < num; i++ {
		pages[i] = &Page{
			Node: Node{
				UrlPath: UrlPath{
					Section: "z",
					Url:     fmt.Sprintf("http://base/x/y/p%d.html", num),
				},
				Site: &SiteInfo{
					BaseUrl: "http://base/",
				},
			},
			Source: Source{File: *source.NewFile(filepath.FromSlash(fmt.Sprintf("/x/y/p%d.md", num)))},
		}
	}

	return pages
}
