// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"errors"
	"fmt"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/viper"
	"html/template"
	"math"
	"path"
)

type pager struct {
	number int
	*paginator
}

type pagers []*pager

var paginatorEmptyPages Pages

type paginator struct {
	paginatedPages []Pages
	pagers
	paginationUrlFactory
	total int
	size  int
}

type paginationUrlFactory func(int) string

// PageNumber returns the current page's number in the pager sequence.
func (p *pager) PageNumber() int {
	return p.number
}

// Url returns the url to the current page.
func (p *pager) Url() template.HTML {
	return template.HTML(p.paginationUrlFactory(p.PageNumber()))
}

// Pages returns the elements on this page.
func (p *pager) Pages() Pages {
	if len(p.paginatedPages) == 0 {
		return paginatorEmptyPages
	}
	return p.paginatedPages[p.PageNumber()-1]
}

// NumberOfElements gets the number of elements on this page.
func (p *pager) NumberOfElements() int {
	return len(p.Pages())
}

// HasPrev tests whether there are page(s) before the current.
func (p *pager) HasPrev() bool {
	return p.PageNumber() > 1
}

// Prev returns the pager for the previous page.
func (p *pager) Prev() *pager {
	if !p.HasPrev() {
		return nil
	}
	return p.pagers[p.PageNumber()-2]
}

// HasNext tests whether there are page(s) after the current.
func (p *pager) HasNext() bool {
	return p.PageNumber() < len(p.paginatedPages)
}

// Next returns the pager for the next page.
func (p *pager) Next() *pager {
	if !p.HasNext() {
		return nil
	}
	return p.pagers[p.PageNumber()]
}

// First returns the pager for the first page.
func (p *pager) First() *pager {
	return p.pagers[0]
}

// Last returns the pager for the last page.
func (p *pager) Last() *pager {
	return p.pagers[len(p.pagers)-1]
}

// Pagers returns a list of pagers that can be used to build a pagination menu.
func (p *paginator) Pagers() pagers {
	return p.pagers
}

// PageSize returns the size of each paginator page.
func (p *paginator) PageSize() int {
	return p.size
}

// TotalPages returns the number of pages in the paginator.
func (p *paginator) TotalPages() int {
	return len(p.paginatedPages)
}

// TotalNumberOfElements returns the number of elements on all pages in this paginator.
func (p *paginator) TotalNumberOfElements() int {
	return p.total
}

func splitPages(pages Pages, size int) []Pages {
	var split []Pages
	for low, j := 0, len(pages); low < j; low += size {
		high := int(math.Min(float64(low+size), float64(len(pages))))
		split = append(split, pages[low:high])
	}

	return split
}

// Paginator gets this Node's paginator if it's already created.
// If it's not, one will be created with all pages in Data["Pages"].
func (n *Node) Paginator() (*pager, error) {

	var initError error

	n.paginatorInit.Do(func() {
		if n.paginator != nil {
			return
		}

		pagers, err := paginatePages(n.Data["Pages"], n.Url)

		if err != nil {
			initError = err
		}

		if len(pagers) > 0 {
			// the rest of the nodes will be created later
			n.paginator = pagers[0]
			n.Site.addToPaginationPageCount(uint64(n.paginator.TotalPages()))
		}

	})

	if initError != nil {
		return nil, initError
	}

	return n.paginator, nil
}

// Paginator on Page isn't supported, calling this yields an error.
func (p *Page) Paginator() (*pager, error) {
	return nil, errors.New("Paginators not supported for content pages.")
}

// Paginate on Page isn't supported, calling this yields an error.
func (p *Page) Paginate(seq interface{}) (*pager, error) {
	return nil, errors.New("Paginators not supported for content pages.")
}

// Paginate gets this Node's paginator if it's already created.
// If it's not, one will be created with the qiven sequence.
// Note that repeated calls will return the same result, even if the sequence is different.
func (n *Node) Paginate(seq interface{}) (*pager, error) {

	var initError error

	n.paginatorInit.Do(func() {
		if n.paginator != nil {
			return
		}
		pagers, err := paginatePages(seq, n.Url)

		if err != nil {
			initError = err
		}

		if len(pagers) > 0 {
			// the rest of the nodes will be created later
			n.paginator = pagers[0]
			n.Site.addToPaginationPageCount(uint64(n.paginator.TotalPages()))
		}

	})

	if initError != nil {
		return nil, initError
	}

	return n.paginator, nil
}

func paginatePages(seq interface{}, section string) (pagers, error) {
	paginateSize := viper.GetInt("paginate")

	if paginateSize <= 0 {
		return nil, errors.New("'paginate' configuration setting must be positive to paginate")
	}
	var pages Pages
	switch seq.(type) {
	case Pages:
		pages = seq.(Pages)
	case *Pages:
		pages = *(seq.(*Pages))
	case WeightedPages:
		pages = (seq.(WeightedPages)).Pages()
	case PageGroup:
		pages = (seq.(PageGroup)).Pages
	default:
		return nil, errors.New(fmt.Sprintf("unsupported type in paginate, got %T", seq))
	}

	urlFactory := newPaginationUrlFactory(section)
	paginator, _ := newPaginator(pages, paginateSize, urlFactory)
	pagers := paginator.Pagers()

	return pagers, nil
}

func newPaginator(pages Pages, size int, urlFactory paginationUrlFactory) (*paginator, error) {

	if size <= 0 {
		return nil, errors.New("Paginator size must be positive")
	}

	split := splitPages(pages, size)

	p := &paginator{total: len(pages), paginatedPages: split, size: size, paginationUrlFactory: urlFactory}

	var ps pagers

	if len(split) > 0 {
		ps = make(pagers, len(split))
		for i := range p.paginatedPages {
			ps[i] = &pager{number: (i + 1), paginator: p}
		}
	} else {
		ps = make(pagers, 1)
		ps[0] = &pager{number: 1, paginator: p}
	}

	p.pagers = ps

	return p, nil
}

func newPaginationUrlFactory(pathElements ...string) paginationUrlFactory {
	paginatePath := viper.GetString("paginatePath")

	return func(page int) string {
		var rel string
		if page == 1 {
			rel = fmt.Sprintf("/%s/", path.Join(pathElements...))
		} else {
			rel = fmt.Sprintf("/%s/%s/%d/", path.Join(pathElements...), paginatePath, page)
		}

		return helpers.UrlizeAndPrep(rel)
	}
}
