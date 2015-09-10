// Copyright © 2014 Steve Francia <spf@spf13.com>.
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
	"sort"
)

var spc = newPageCache()

/*
 * Implementation of a custom sorter for Pages
 */

// A PageSorter implements the sort interface for Pages
type PageSorter struct {
	pages Pages
	by    PageBy
}

// PageBy is a closure used in the Sort.Less method.
type PageBy func(p1, p2 *Page) bool

func (by PageBy) Sort(pages Pages) {
	ps := &PageSorter{
		pages: pages,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ps)
}

var DefaultPageSort = func(p1, p2 *Page) bool {
	if p1.Weight == p2.Weight {
		if p1.Date.Unix() == p2.Date.Unix() {
			return (p1.LinkTitle() < p2.LinkTitle())
		}
		return p1.Date.Unix() > p2.Date.Unix()
	}
	return p1.Weight < p2.Weight
}

func (ps *PageSorter) Len() int      { return len(ps.pages) }
func (ps *PageSorter) Swap(i, j int) { ps.pages[i], ps.pages[j] = ps.pages[j], ps.pages[i] }

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ps *PageSorter) Less(i, j int) bool { return ps.by(ps.pages[i], ps.pages[j]) }

func (p Pages) Sort() {
	PageBy(DefaultPageSort).Sort(p)
}

func (p Pages) Limit(n int) Pages {
	if len(p) < n {
		return p[0:n]
	}
	return p
}

// ByWeight sorts the Pages by weight and returns a copy.
//
// Adjacent invocactions on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByWeight() Pages {
	key := "pageSort.ByWeight"
	pages, _ := spc.get(key, p, PageBy(DefaultPageSort).Sort)
	return pages
}

// ByTitle sorts the Pages by title and returns a copy.
//
// Adjacent invocactions on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByTitle() Pages {

	key := "pageSort.ByTitle"

	title := func(p1, p2 *Page) bool {
		return p1.Title < p2.Title
	}

	pages, _ := spc.get(key, p, PageBy(title).Sort)
	return pages
}

// ByLinkTitle sorts the Pages by link title and returns a copy.
//
// Adjacent invocactions on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByLinkTitle() Pages {

	key := "pageSort.ByLinkTitle"

	linkTitle := func(p1, p2 *Page) bool {
		return p1.linkTitle < p2.linkTitle
	}

	pages, _ := spc.get(key, p, PageBy(linkTitle).Sort)

	return pages
}

// ByDate sorts the Pages by date and returns a copy.
//
// Adjacent invocactions on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByDate() Pages {

	key := "pageSort.ByDate"

	date := func(p1, p2 *Page) bool {
		return p1.Date.Unix() < p2.Date.Unix()
	}

	pages, _ := spc.get(key, p, PageBy(date).Sort)

	return pages
}

// ByPublishDate sorts the Pages by publish date and returns a copy.
//
// Adjacent invocactions on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByPublishDate() Pages {

	key := "pageSort.ByPublishDate"

	pubDate := func(p1, p2 *Page) bool {
		return p1.PublishDate.Unix() < p2.PublishDate.Unix()
	}

	pages, _ := spc.get(key, p, PageBy(pubDate).Sort)

	return pages
}

// ByLength sorts the Pages by length and returns a copy.
//
// Adjacent invocactions on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByLength() Pages {

	key := "pageSort.ByLength"

	length := func(p1, p2 *Page) bool {
		return len(p1.Content) < len(p2.Content)
	}

	pages, _ := spc.get(key, p, PageBy(length).Sort)

	return pages
}

// Reverse reverses the order in Pages and returns a copy.
//
// Adjacent invocactions on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) Reverse() Pages {
	key := "pageSort.Reverse"

	reverseFunc := func(pages Pages) {
		for i, j := 0, len(pages)-1; i < j; i, j = i+1, j-1 {
			pages[i], pages[j] = pages[j], pages[i]
		}
	}

	pages, _ := spc.get(key, p, reverseFunc)

	return pages
}
