// Copyright 2019 The Hugo Authors. All rights reserved.
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

package page

import (
	"sort"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/compare"
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
type pageBy func(p1, p2 Page) bool

// Sort stable sorts the pages given the receiver's sort order.
func (by pageBy) Sort(pages Pages) {
	ps := &pageSorter{
		pages: pages,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ps)
}

var (

	// DefaultPageSort is the default sort func for pages in Hugo:
	// Order by Weight, Date, LinkTitle and then full file path.
	DefaultPageSort = func(p1, p2 Page) bool {
		if p1.Weight() == p2.Weight() {
			if p1.Date().Unix() == p2.Date().Unix() {
				c := compare.Strings(p1.LinkTitle(), p2.LinkTitle())
				if c == 0 {
					if p1.File().IsZero() || p2.File().IsZero() {
						return p1.File().IsZero()
					}
					return compare.LessStrings(p1.File().Filename(), p2.File().Filename())
				}
				return c < 0
			}
			return p1.Date().Unix() > p2.Date().Unix()
		}

		if p2.Weight() == 0 {
			return true
		}

		if p1.Weight() == 0 {
			return false
		}

		return p1.Weight() < p2.Weight()
	}

	lessPageLanguage = func(p1, p2 Page) bool {

		if p1.Language().Weight == p2.Language().Weight {
			if p1.Date().Unix() == p2.Date().Unix() {
				c := compare.Strings(p1.LinkTitle(), p2.LinkTitle())
				if c == 0 {
					if !p1.File().IsZero() && !p2.File().IsZero() {
						return compare.LessStrings(p1.File().Filename(), p2.File().Filename())
					}
				}
				return c < 0
			}
			return p1.Date().Unix() > p2.Date().Unix()
		}

		if p2.Language().Weight == 0 {
			return true
		}

		if p1.Language().Weight == 0 {
			return false
		}

		return p1.Language().Weight < p2.Language().Weight
	}

	lessPageTitle = func(p1, p2 Page) bool {
		return compare.LessStrings(p1.Title(), p2.Title())
	}

	lessPageLinkTitle = func(p1, p2 Page) bool {
		return compare.LessStrings(p1.LinkTitle(), p2.LinkTitle())
	}

	lessPageDate = func(p1, p2 Page) bool {
		return p1.Date().Unix() < p2.Date().Unix()
	}

	lessPagePubDate = func(p1, p2 Page) bool {
		return p1.PublishDate().Unix() < p2.PublishDate().Unix()
	}
)

func (ps *pageSorter) Len() int      { return len(ps.pages) }
func (ps *pageSorter) Swap(i, j int) { ps.pages[i], ps.pages[j] = ps.pages[j], ps.pages[i] }

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (ps *pageSorter) Less(i, j int) bool { return ps.by(ps.pages[i], ps.pages[j]) }

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
	pages, _ := spc.get(key, pageBy(DefaultPageSort).Sort, p)
	return pages
}

// SortByDefault sorts pages by the default sort.
func SortByDefault(pages Pages) {
	pageBy(DefaultPageSort).Sort(pages)
}

// ByTitle sorts the Pages by title and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByTitle() Pages {

	const key = "pageSort.ByTitle"

	pages, _ := spc.get(key, pageBy(lessPageTitle).Sort, p)
	return pages
}

// ByLinkTitle sorts the Pages by link title and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByLinkTitle() Pages {

	const key = "pageSort.ByLinkTitle"

	pages, _ := spc.get(key, pageBy(lessPageLinkTitle).Sort, p)

	return pages
}

// ByDate sorts the Pages by date and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByDate() Pages {

	const key = "pageSort.ByDate"

	pages, _ := spc.get(key, pageBy(lessPageDate).Sort, p)

	return pages
}

// ByPublishDate sorts the Pages by publish date and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByPublishDate() Pages {

	const key = "pageSort.ByPublishDate"

	pages, _ := spc.get(key, pageBy(lessPagePubDate).Sort, p)

	return pages
}

// ByExpiryDate sorts the Pages by publish date and returns a copy.
//
// Adjacent invocations on the same receiver will return a cached result.
//
// This may safely be executed  in parallel.
func (p Pages) ByExpiryDate() Pages {

	const key = "pageSort.ByExpiryDate"

	expDate := func(p1, p2 Page) bool {
		return p1.ExpiryDate().Unix() < p2.ExpiryDate().Unix()
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

	date := func(p1, p2 Page) bool {
		return p1.Lastmod().Unix() < p2.Lastmod().Unix()
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

	length := func(p1, p2 Page) bool {

		p1l, ok1 := p1.(resource.LengthProvider)
		p2l, ok2 := p2.(resource.LengthProvider)

		if !ok1 {
			return true
		}

		if !ok2 {
			return false
		}

		return p1l.Len() < p2l.Len()
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

	pages, _ := spc.get(key, pageBy(lessPageLanguage).Sort, p)

	return pages
}

// SortByLanguage sorts the pages by language.
func SortByLanguage(pages Pages) {
	pageBy(lessPageLanguage).Sort(pages)
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

	paramsKeyComparator := func(p1, p2 Page) bool {
		v1, _ := p1.Param(paramsKeyStr)
		v2, _ := p2.Param(paramsKeyStr)

		if v1 == nil {
			return false
		}

		if v2 == nil {
			return true
		}

		isNumeric := func(v interface{}) bool {
			switch v.(type) {
			case uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
				return true
			default:
				return false
			}
		}

		if isNumeric(v1) && isNumeric(v2) {
			return cast.ToFloat64(v1) < cast.ToFloat64(v2)
		}

		s1 := cast.ToString(v1)
		s2 := cast.ToString(v2)

		return compare.LessStrings(s1, s2)

	}

	pages, _ := spc.get(key, pageBy(paramsKeyComparator).Sort, p)

	return pages
}
