// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/helpers"

	"sort"

	"github.com/spf13/cast"
)

var spc = newPageCache()

/*
 * Implementation of a custom sorter for Pages
 */

// A pageSorter implements the sort interface for Pages
type pageSorter struct {
	pages Pages
	by    pageBy
}

// pageBy is a closure used in the Sort.Less method.
type pageBy func(p1, p2 *Page) bool

// Sort stable sorts the pages given the receiver's sort order.
func (by pageBy) Sort(pages Pages) {
	ps := &pageSorter{
		pages: pages,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ps)
}

// defaultPageSort is the default sort for pages in Hugo:
// Order by Weight, Date, LinkTitle and then full file path.
var defaultPageSort = func(p1, p2 *Page) bool {
	if p1.Weight == p2.Weight {
		if p1.Date.Unix() == p2.Date.Unix() {
			if p1.LinkTitle() == p2.LinkTitle() {
				return (p1.FullFilePath() < p2.FullFilePath())
			}
			return (p1.LinkTitle() < p2.LinkTitle())
		}
		return p1.Date.Unix() > p2.Date.Unix()
	}

	if p2.Weight == 0 {
		return true
	}

	if p1.Weight == 0 {
		return false
	}

	return p1.Weight < p2.Weight
}

var languagePageSort = func(p1, p2 *Page) bool {
	if p1.Language().Weight == p2.Language().Weight {
		if p1.Date.Unix() == p2.Date.Unix() {
			if p1.LinkTitle() == p2.LinkTitle() {
				return (p1.FullFilePath() < p2.FullFilePath())
			}
			return (p1.LinkTitle() < p2.LinkTitle())
		}
		return p1.Date.Unix() > p2.Date.Unix()
	}

	if p2.Language().Weight == 0 {
		return true
	}

	if p1.Language().Weight == 0 {
		return false
	}

	return p1.Language().Weight < p2.Language().Weight
}

func (ps *pageSorter) Len() int      { return len(ps.pages) }
func (ps *pageSorter) Swap(i, j int) { ps.pages[i], ps.pages[j] = ps.pages[j], ps.pages[i] }

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ps *pageSorter) Less(i, j int) bool { return ps.by(ps.pages[i], ps.pages[j]) }

// Sort sorts the pages by the default sort order defined:
// Order by Weight, Date, LinkTitle and then full file path.
func (p Pages) Sort() {
	// Remove in Hugo 0.51
	helpers.Deprecated("Pages", "Sort", "Use .ByWeight", false)
	p.sort()
}

func (p Pages) sort() {
	pageBy(defaultPageSort).Sort(p)
}

// Limit limits the number of pages returned to n.
func (p Pages) Limit(n int) Pages {
	if len(p) > n {
		return p[0:n]
	}
	return p
}

// ByWeight sorts the Pages by weight and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByWeight() Pages {
	const key = "pageSort.ByWeight"
	pages, _ := spc.get(key, pageBy(defaultPageSort).Sort, p)
	return pages
}

// ByTitle sorts the Pages by title and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByTitle() Pages {

	const key = "pageSort.ByTitle"

	title := func(p1, p2 *Page) bool {
		return p1.title < p2.title
	}

	pages, _ := spc.get(key, pageBy(title).Sort, p)
	return pages
}

// ByLinkTitle sorts the Pages by link title and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByLinkTitle() Pages {

	const key = "pageSort.ByLinkTitle"

	linkTitle := func(p1, p2 *Page) bool {
		return p1.linkTitle < p2.linkTitle
	}

	pages, _ := spc.get(key, pageBy(linkTitle).Sort, p)

	return pages
}

// ByDate sorts the Pages by date and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByDate() Pages {

	const key = "pageSort.ByDate"

	date := func(p1, p2 *Page) bool {
		return p1.Date.Unix() < p2.Date.Unix()
	}

	pages, _ := spc.get(key, pageBy(date).Sort, p)

	return pages
}

// ByPublishDate sorts the Pages by publish date and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByPublishDate() Pages {

	const key = "pageSort.ByPublishDate"

	pubDate := func(p1, p2 *Page) bool {
		return p1.PublishDate.Unix() < p2.PublishDate.Unix()
	}

	pages, _ := spc.get(key, pageBy(pubDate).Sort, p)

	return pages
}

// ByExpiryDate sorts the Pages by publish date and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByExpiryDate() Pages {

	const key = "pageSort.ByExpiryDate"

	expDate := func(p1, p2 *Page) bool {
		return p1.ExpiryDate.Unix() < p2.ExpiryDate.Unix()
	}

	pages, _ := spc.get(key, pageBy(expDate).Sort, p)

	return pages
}

// ByLastmod sorts the Pages by the last modification date and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByLastmod() Pages {

	const key = "pageSort.ByLastmod"

	date := func(p1, p2 *Page) bool {
		return p1.Lastmod.Unix() < p2.Lastmod.Unix()
	}

	pages, _ := spc.get(key, pageBy(date).Sort, p)

	return pages
}

// ByLength sorts the Pages by length and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByLength() Pages {

	const key = "pageSort.ByLength"

	length := func(p1, p2 *Page) bool {
		return len(p1.content()) < len(p2.content())
	}

	pages, _ := spc.get(key, pageBy(length).Sort, p)

	return pages
}

// ByLanguage sorts the Pages by the language's Weight.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByLanguage() Pages {

	const key = "pageSort.ByLanguage"

	pages, _ := spc.get(key, pageBy(languagePageSort).Sort, p)

	return pages
}

// Reverse reverses the order in Pages and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) Reverse() Pages {
	const key = "pageSort.Reverse"

	reverseFunc := func(pages Pages) {
		for i, j := 0, len(pages)-1; i < j; i, j = i+1, j-1 {
			pages[i], pages[j] = pages[j], pages[i]
		}
	}

	pages, _ := spc.get(key, reverseFunc, p)

	return pages
}

// ByParam sorts the pages according to the given page Params key.
//
// Adjacent invocations on the same receiver with the same paramsKey will return a cached result.
//
// This may safely be executed  in parallel.

func (p Pages) ByParam(paramsKey interface{}) Pages {
	paramsKeyStr := cast.ToString(paramsKey)
	key := "pageSort.ByParam." + paramsKeyStr

	paramsKeyComparator := func(p1, p2 *Page) bool {
		v1, _ := p1.Param(paramsKeyStr)
		v2, _ := p2.Param(paramsKeyStr)
		s1 := cast.ToString(v1)
		s2 := cast.ToString(v2)

		// Sort nils last.
		if s1 == "" {
			return false
		} else if s2 == "" {
			return true
		}

		return s1 < s2
	}

	pages, _ := spc.get(key, pageBy(paramsKeyComparator).Sort, p)

	return pages
}
