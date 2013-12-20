// Copyright © 2013 Steve Francia <spf@spf13.com>.
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
	"github.com/spf13/hugo/helpers"
	"sort"
)

/*
 *  An index list is a list of all indexes and their values
 *  EG. List['tags'] => TagIndex (from above)
 */
type IndexList map[string]Index

/*
 *  An index is a map of keywords to a list of pages.
 *  For example
 *    TagIndex['technology'] = WeightedPages
 *    TagIndex['golang']  =  WeightedPages2
 */
type Index map[string]WeightedPages

/*
 *  A list of Pages with their corresponding (and relative) weight
 *  [{Weight: 30, Page: *1}, {Weight: 40, Page: *2}]
 */
type WeightedPages []WeightedPage
type WeightedPage struct {
	Weight int
	Page   *Page
}

/*
 * This is another representation of an Index using an array rather than a map.
 * Important because you can't order a map.
 */
type OrderedIndex []OrderedIndexEntry

/*
 * Similar to an element of an Index, but with the key embedded (as name)
 * Eg:  {Name: Technology, WeightedPages: Indexedpages}
 */
type OrderedIndexEntry struct {
	Name          string
	WeightedPages WeightedPages
}

// KeyPrep... Indexes should be case insensitive. Can make it easily conditional later.
func kp(in string) string {
	return helpers.Urlize(in)
}

func (i Index) Get(key string) WeightedPages { return i[kp(key)] }
func (i Index) Count(key string) int         { return len(i[kp(key)]) }
func (i Index) Add(key string, w WeightedPage) {
	key = kp(key)
	i[key] = append(i[key], w)
}

// Returns an ordered index with a non defined order
func (i Index) IndexArray() OrderedIndex {
	ies := make([]OrderedIndexEntry, len(i))
	count := 0
	for k, v := range i {
		ies[count] = OrderedIndexEntry{Name: k, WeightedPages: v}
		count++
	}
	return ies
}

// Returns an ordered index sorted by key name
func (i Index) Alphabetical() OrderedIndex {
	name := func(i1, i2 *OrderedIndexEntry) bool {
		return i1.Name < i2.Name
	}

	ia := i.IndexArray()
	OIby(name).Sort(ia)
	return ia
}

// Returns an ordered index sorted by # of pages per key
func (i Index) ByCount() OrderedIndex {
	count := func(i1, i2 *OrderedIndexEntry) bool {
		return len(i1.WeightedPages) > len(i2.WeightedPages)
	}

	ia := i.IndexArray()
	OIby(count).Sort(ia)
	return ia
}

// Helper to move the page access up a level
func (ie OrderedIndexEntry) Pages() []*Page {
	return ie.WeightedPages.Pages()
}

func (ie OrderedIndexEntry) Count() int {
	return len(ie.WeightedPages)
}

/*
 * Implementation of a custom sorter for OrderedIndexes
 */

// A type to implement the sort interface for IndexEntries.
type orderedIndexSorter struct {
	index OrderedIndex
	by    OIby
}

// Closure used in the Sort.Less method.
type OIby func(i1, i2 *OrderedIndexEntry) bool

func (by OIby) Sort(index OrderedIndex) {
	ps := &orderedIndexSorter{
		index: index,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// Len is part of sort.Interface.
func (s *orderedIndexSorter) Len() int {
	return len(s.index)
}

// Swap is part of sort.Interface.
func (s *orderedIndexSorter) Swap(i, j int) {
	s.index[i], s.index[j] = s.index[j], s.index[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *orderedIndexSorter) Less(i, j int) bool {
	return s.by(&s.index[i], &s.index[j])
}

func (wp WeightedPages) Pages() Pages {
	pages := make(Pages, len(wp))
	for i := range wp {
		pages[i] = wp[i].Page
	}
	return pages
}

func (p WeightedPages) Len() int      { return len(p) }
func (p WeightedPages) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p WeightedPages) Sort()         { sort.Sort(p) }
func (p WeightedPages) Count() int    { return len(p) }
func (p WeightedPages) Less(i, j int) bool {
	if p[i].Weight == p[j].Weight {
		return p[i].Page.Date.Unix() > p[j].Page.Date.Unix()
	} else {
		return p[i].Weight < p[j].Weight
	}
}

// TODO mimic PagesSorter for WeightedPages
