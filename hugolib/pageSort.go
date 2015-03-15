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

func (p Pages) ByWeight() Pages {
	PageBy(DefaultPageSort).Sort(p)
	return p
}

func (p Pages) ByTitle() Pages {
	title := func(p1, p2 *Page) bool {
		return p1.Title < p2.Title
	}

	PageBy(title).Sort(p)
	return p
}

func (p Pages) ByLinkTitle() Pages {
	linkTitle := func(p1, p2 *Page) bool {
		return p1.linkTitle < p2.linkTitle
	}

	PageBy(linkTitle).Sort(p)
	return p
}

func (p Pages) ByDate() Pages {
	date := func(p1, p2 *Page) bool {
		return p1.Date.Unix() < p2.Date.Unix()
	}

	PageBy(date).Sort(p)
	return p
}

func (p Pages) ByPublishDate() Pages {
	pubDate := func(p1, p2 *Page) bool {
		return p1.PublishDate.Unix() < p2.PublishDate.Unix()
	}

	PageBy(pubDate).Sort(p)
	return p
}

func (p Pages) ByLength() Pages {
	length := func(p1, p2 *Page) bool {
		return len(p1.Content) < len(p2.Content)
	}

	PageBy(length).Sort(p)
	return p
}

func (p Pages) Reverse() Pages {
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}

	return p
}
