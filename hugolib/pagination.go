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
	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/viper"
	"html/template"
	"math"
	"path"
	"reflect"
	"strings"
)

type Pager struct {
	number int
	*paginator
}

type paginatedElement interface {
	Len() int
}

func (p Pages) Len() int {
	return len(p)
}

func (psg PagesGroup) Len() int {
	l := 0
	for _, pg := range psg {
		l += len(pg.Pages)
	}
	return l
}

type pagers []*Pager

var paginatorEmptyPages Pages
var paginatorEmptyPageGroups PagesGroup

type paginator struct {
	paginatedElements []paginatedElement
	pagers
	paginationURLFactory
	total   int
	size    int
	source  interface{}
	options []interface{}
}

type paginationURLFactory func(int) string

// PageNumber returns the current page's number in the pager sequence.
func (p *Pager) PageNumber() int {
	return p.number
}

// URL returns the URL to the current page.
func (p *Pager) URL() template.HTML {
	return template.HTML(p.paginationURLFactory(p.PageNumber()))
}

// Pages returns the Pages on this page.
// Note: If this return a non-empty result, then PageGroups() will return empty.
func (p *Pager) Pages() Pages {
	if len(p.paginatedElements) == 0 {
		return paginatorEmptyPages
	}

	if pages, ok := p.element().(Pages); ok {
		return pages
	}

	return paginatorEmptyPages
}

// PageGroups return Page groups for this page.
// Note: If this return non-empty result, then Pages() will return empty.
func (p *Pager) PageGroups() PagesGroup {
	if len(p.paginatedElements) == 0 {
		return paginatorEmptyPageGroups
	}

	if groups, ok := p.element().(PagesGroup); ok {
		return groups
	}

	return paginatorEmptyPageGroups
}

func (p *Pager) element() paginatedElement {
	if len(p.paginatedElements) == 0 {
		return paginatorEmptyPages
	}
	return p.paginatedElements[p.PageNumber()-1]
}

// page returns the Page with the given index
func (p *Pager) page(index int) (*Page, error) {

	if pages, ok := p.element().(Pages); ok {
		if pages != nil && len(pages) > index {
			return pages[index], nil
		}
		return nil, nil
	}

	// must be PagesGroup
	// this construction looks clumsy, but ...
	// ... it is the difference between 99.5% and 100% test coverage :-)
	groups := p.element().(PagesGroup)

	i := 0
	for _, v := range groups {
		for _, page := range v.Pages {
			if i == index {
				return page, nil
			}
			i++
		}
	}
	return nil, nil
}

// NumberOfElements gets the number of elements on this page.
func (p *Pager) NumberOfElements() int {
	return p.element().Len()
}

// HasPrev tests whether there are page(s) before the current.
func (p *Pager) HasPrev() bool {
	return p.PageNumber() > 1
}

// Prev returns the pager for the previous page.
func (p *Pager) Prev() *Pager {
	if !p.HasPrev() {
		return nil
	}
	return p.pagers[p.PageNumber()-2]
}

// HasNext tests whether there are page(s) after the current.
func (p *Pager) HasNext() bool {
	return p.PageNumber() < len(p.paginatedElements)
}

// Next returns the pager for the next page.
func (p *Pager) Next() *Pager {
	if !p.HasNext() {
		return nil
	}
	return p.pagers[p.PageNumber()]
}

// First returns the pager for the first page.
func (p *Pager) First() *Pager {
	return p.pagers[0]
}

// Last returns the pager for the last page.
func (p *Pager) Last() *Pager {
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
	return len(p.paginatedElements)
}

// TotalNumberOfElements returns the number of elements on all pages in this paginator.
func (p *paginator) TotalNumberOfElements() int {
	return p.total
}

func splitPages(pages Pages, size int) []paginatedElement {
	var split []paginatedElement
	for low, j := 0, len(pages); low < j; low += size {
		high := int(math.Min(float64(low+size), float64(len(pages))))
		split = append(split, pages[low:high])
	}

	return split
}

func splitPageGroups(pageGroups PagesGroup, size int) []paginatedElement {
	var split []paginatedElement

	type keyPage struct {
		key  interface{}
		page *Page
	}

	var flattened []keyPage

	for _, g := range pageGroups {
		for _, p := range g.Pages {
			flattened = append(flattened, keyPage{g.Key, p})
		}
	}

	numPages := len(flattened)

	for low, j := 0, numPages; low < j; low += size {
		high := int(math.Min(float64(low+size), float64(numPages)))
		var pg PagesGroup
		var key interface{} = nil
		var groupIndex = -1
		for k := low; k < high; k++ {
			kp := flattened[k]
			if key == nil || key != kp.key {
				key = kp.key
				pg = append(pg, PageGroup{Key: key})
				groupIndex++
			}
			pg[groupIndex].Pages = append(pg[groupIndex].Pages, kp.page)
		}
		split = append(split, pg)
	}

	return split
}

// Paginator gets this Node's paginator if it's already created.
// If it's not, one will be created with all pages in Data["Pages"].
func (n *Node) Paginator(options ...interface{}) (*Pager, error) {

	pagerSize, err := resolvePagerSize(options...)

	if err != nil {
		return nil, err
	}

	var initError error

	n.paginatorInit.Do(func() {
		if n.paginator != nil {
			return
		}

		pagers, err := paginatePages(n.Data["Pages"], pagerSize, n.URL)

		if err != nil {
			initError = err
		}

		if len(pagers) > 0 {
			// the rest of the nodes will be created later
			n.paginator = pagers[0]
			n.paginator.source = "paginator"
			n.paginator.options = options
			n.Site.addToPaginationPageCount(uint64(n.paginator.TotalPages()))
		}

	})

	if initError != nil {
		return nil, initError
	}

	return n.paginator, nil
}

// Paginator on Page isn't supported, calling this yields an error.
func (p *Page) Paginator(options ...interface{}) (*Pager, error) {
	return nil, errors.New("Paginators not supported for content pages.")
}

// Paginate on Page isn't supported, calling this yields an error.
func (p *Page) Paginate(seq interface{}, options ...interface{}) (*Pager, error) {
	return nil, errors.New("Paginators not supported for content pages.")
}

// Paginate gets this Node's paginator if it's already created.
// If it's not, one will be created with the qiven sequence.
// Note that repeated calls will return the same result, even if the sequence is different.
func (n *Node) Paginate(seq interface{}, options ...interface{}) (*Pager, error) {

	pagerSize, err := resolvePagerSize(options...)

	if err != nil {
		return nil, err
	}

	var initError error

	n.paginatorInit.Do(func() {
		if n.paginator != nil {
			return
		}
		pagers, err := paginatePages(seq, pagerSize, n.URL)

		if err != nil {
			initError = err
		}

		if len(pagers) > 0 {
			// the rest of the nodes will be created later
			n.paginator = pagers[0]
			n.paginator.source = seq
			n.paginator.options = options
			n.Site.addToPaginationPageCount(uint64(n.paginator.TotalPages()))
		}

	})

	if initError != nil {
		return nil, initError
	}

	if n.paginator.source == "paginator" {
		return nil, errors.New("a Paginator was previously built for this Node without filters; look for earlier .Paginator usage")
	}

	if !reflect.DeepEqual(options, n.paginator.options) || !probablyEqualPageLists(n.paginator.source, seq) {
		return nil, errors.New("invoked multiple times with different arguments")
	}

	return n.paginator, nil
}

func resolvePagerSize(options ...interface{}) (int, error) {
	if len(options) == 0 {
		return viper.GetInt("paginate"), nil
	}

	if len(options) > 1 {
		return -1, errors.New("too many arguments, 'pager size' is currently the only option")
	}

	pas, err := cast.ToIntE(options[0])

	if err != nil || pas <= 0 {
		return -1, errors.New(("'pager size' must be a positive integer"))
	}

	return pas, nil
}

func paginatePages(seq interface{}, pagerSize int, section string) (pagers, error) {

	if pagerSize <= 0 {
		return nil, errors.New("'paginate' configuration setting must be positive to paginate")
	}

	section = strings.TrimSuffix(section, ".html")
	urlFactory := newPaginationURLFactory(section)

	var paginator *paginator

	if groups, ok := seq.(PagesGroup); ok {
		paginator, _ = newPaginatorFromPageGroups(groups, pagerSize, urlFactory)
	} else {
		pages, err := toPages(seq)
		if err != nil {
			return nil, err
		}
		paginator, _ = newPaginatorFromPages(pages, pagerSize, urlFactory)
	}

	pagers := paginator.Pagers()

	return pagers, nil
}

func toPages(seq interface{}) (Pages, error) {
	switch seq.(type) {
	case Pages:
		return seq.(Pages), nil
	case *Pages:
		return *(seq.(*Pages)), nil
	case WeightedPages:
		return (seq.(WeightedPages)).Pages(), nil
	case PageGroup:
		return (seq.(PageGroup)).Pages, nil
	default:
		return nil, fmt.Errorf("unsupported type in paginate, got %T", seq)
	}
}

// probablyEqual checks page lists for probable equality.
// It may return false positives.
// The motivation behind this is to avoid potential costly reflect.DeepEqual
// when "probably" is good enough.
func probablyEqualPageLists(a1 interface{}, a2 interface{}) bool {

	if a1 == nil || a2 == nil {
		return a1 == a2
	}

	t1 := reflect.TypeOf(a1)
	t2 := reflect.TypeOf(a2)

	if t1 != t2 {
		return false
	}

	if g1, ok := a1.(PagesGroup); ok {
		g2 := a2.(PagesGroup)
		if len(g1) != len(g2) {
			return false
		}
		if len(g1) == 0 {
			return true
		}
		if g1.Len() != g2.Len() {
			return false
		}

		return g1[0].Pages[0] == g2[0].Pages[0]
	}

	p1, err1 := toPages(a1)
	p2, err2 := toPages(a2)

	// probably the same wrong type
	if err1 != nil && err2 != nil {
		return true
	}

	if len(p1) != len(p2) {
		return false
	}

	if len(p1) == 0 {
		return true
	}

	return p1[0] == p2[0]
}

func newPaginatorFromPages(pages Pages, size int, urlFactory paginationURLFactory) (*paginator, error) {

	if size <= 0 {
		return nil, errors.New("Paginator size must be positive")
	}

	split := splitPages(pages, size)

	return newPaginator(split, len(pages), size, urlFactory)
}

func newPaginatorFromPageGroups(pageGroups PagesGroup, size int, urlFactory paginationURLFactory) (*paginator, error) {

	if size <= 0 {
		return nil, errors.New("Paginator size must be positive")
	}

	split := splitPageGroups(pageGroups, size)

	return newPaginator(split, pageGroups.Len(), size, urlFactory)
}

func newPaginator(elements []paginatedElement, total, size int, urlFactory paginationURLFactory) (*paginator, error) {
	p := &paginator{total: total, paginatedElements: elements, size: size, paginationURLFactory: urlFactory}

	var ps pagers

	if len(elements) > 0 {
		ps = make(pagers, len(elements))
		for i := range p.paginatedElements {
			ps[i] = &Pager{number: (i + 1), paginator: p}
		}
	} else {
		ps = make(pagers, 1)
		ps[0] = &Pager{number: 1, paginator: p}
	}

	p.pagers = ps

	return p, nil
}

func newPaginationURLFactory(pathElements ...string) paginationURLFactory {
	paginatePath := viper.GetString("paginatePath")

	return func(page int) string {
		var rel string
		if page == 1 {
			rel = fmt.Sprintf("/%s/", path.Join(pathElements...))
		} else {
			rel = fmt.Sprintf("/%s/%s/%d/", path.Join(pathElements...), paginatePath, page)
		}

		return helpers.URLizeAndPrep(rel)
	}
}
