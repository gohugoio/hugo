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
	"sort"

	"github.com/spf13/hugo/helpers"
)

/*
 *  An taxonomy list is a list of all taxonomies and their values
 *  EG. List['tags'] => TagTaxonomy (from above)
 */
type TaxonomyList map[string]Taxonomy

/*
 *  An taxonomy is a map of keywords to a list of pages.
 *  For example
 *    TagTaxonomy['technology'] = WeightedPages
 *    TagTaxonomy['go']  =  WeightedPages2
 */
type Taxonomy map[string]WeightedPages

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
 * This is another representation of an Taxonomy using an array rather than a map.
 * Important because you can't order a map.
 */
type OrderedTaxonomy []OrderedTaxonomyEntry

/*
 * Similar to an element of an Taxonomy, but with the key embedded (as name)
 * Eg:  {Name: Technology, WeightedPages: Taxonomyedpages}
 */
type OrderedTaxonomyEntry struct {
	Name          string
	WeightedPages WeightedPages
}

// KeyPrep... Taxonomies should be case insensitive. Can make it easily conditional later.
func kp(in string) string {
	return helpers.MakePathToLower(in)
}

func (i Taxonomy) Get(key string) WeightedPages { return i[kp(key)] }
func (i Taxonomy) Count(key string) int         { return len(i[kp(key)]) }
func (i Taxonomy) Add(key string, w WeightedPage) {
	key = kp(key)
	i[key] = append(i[key], w)
}

// Returns an ordered taxonomy with a non defined order
func (i Taxonomy) TaxonomyArray() OrderedTaxonomy {
	ies := make([]OrderedTaxonomyEntry, len(i))
	count := 0
	for k, v := range i {
		ies[count] = OrderedTaxonomyEntry{Name: k, WeightedPages: v}
		count++
	}
	return ies
}

// Returns an ordered taxonomy sorted by key name
func (i Taxonomy) Alphabetical() OrderedTaxonomy {
	name := func(i1, i2 *OrderedTaxonomyEntry) bool {
		return i1.Name < i2.Name
	}

	ia := i.TaxonomyArray()
	OIby(name).Sort(ia)
	return ia
}

// Returns an ordered taxonomy sorted by # of pages per key
func (i Taxonomy) ByCount() OrderedTaxonomy {
	count := func(i1, i2 *OrderedTaxonomyEntry) bool {
		return len(i1.WeightedPages) > len(i2.WeightedPages)
	}

	ia := i.TaxonomyArray()
	OIby(count).Sort(ia)
	return ia
}

// Helper to move the page access up a level
func (ie OrderedTaxonomyEntry) Pages() Pages {
	return ie.WeightedPages.Pages()
}

func (ie OrderedTaxonomyEntry) Count() int {
	return len(ie.WeightedPages)
}

func (ie OrderedTaxonomyEntry) Term() string {
	return ie.Name
}

/*
 * Implementation of a custom sorter for OrderedTaxonomies
 */

// A type to implement the sort interface for TaxonomyEntries.
type orderedTaxonomySorter struct {
	taxonomy OrderedTaxonomy
	by       OIby
}

// Closure used in the Sort.Less method.
type OIby func(i1, i2 *OrderedTaxonomyEntry) bool

func (by OIby) Sort(taxonomy OrderedTaxonomy) {
	ps := &orderedTaxonomySorter{
		taxonomy: taxonomy,
		by:       by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ps)
}

// Len is part of sort.Interface.
func (s *orderedTaxonomySorter) Len() int {
	return len(s.taxonomy)
}

// Swap is part of sort.Interface.
func (s *orderedTaxonomySorter) Swap(i, j int) {
	s.taxonomy[i], s.taxonomy[j] = s.taxonomy[j], s.taxonomy[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *orderedTaxonomySorter) Less(i, j int) bool {
	return s.by(&s.taxonomy[i], &s.taxonomy[j])
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
func (p WeightedPages) Sort()         { sort.Stable(p) }
func (p WeightedPages) Count() int    { return len(p) }
func (p WeightedPages) Less(i, j int) bool {
	if p[i].Weight == p[j].Weight {
		if p[i].Page.Date.Equal(p[j].Page.Date) {
			return p[i].Page.Title < p[j].Page.Title
		}
		return p[i].Page.Date.After(p[i].Page.Date)
	}
	return p[i].Weight < p[j].Weight
}

// TODO mimic PagesSorter for WeightedPages
