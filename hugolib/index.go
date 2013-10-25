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
	"github.com/spf13/hugo/template"
	"sort"
)

type WeightedIndexEntry struct {
	Weight int
	Page   *Page
}

type IndexedPages []WeightedIndexEntry

func (p IndexedPages) Len() int      { return len(p) }
func (p IndexedPages) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p IndexedPages) Sort()         { sort.Sort(p) }
func (p IndexedPages) Less(i, j int) bool {
	if p[i].Weight == p[j].Weight {
		return p[i].Page.Date.Unix() > p[j].Page.Date.Unix()
	} else {
		return p[i].Weight > p[j].Weight
	}
}

func (ip IndexedPages) Pages() Pages {
	pages := make(Pages, len(ip))
	for i, _ := range ip {
		pages[i] = ip[i].Page
	}
	return pages
}

type Index map[string]IndexedPages
type IndexList map[string]Index

// KeyPrep... Indexes should be case insensitive. Can make it easily conditional later.
func kp(in string) string {
	return template.Urlize(in)
}

func (i Index) Get(key string) IndexedPages { return i[kp(key)] }
func (i Index) Count(key string) int        { return len(i[kp(key)]) }
func (i Index) Add(key string, w WeightedIndexEntry) {
	key = kp(key)
	i[key] = append(i[key], w)
}

func (i Index) IndexArray() []IndexEntry {
	ies := make([]IndexEntry, len(i))
	count := 0
	for k, v := range i {
		ies[count] = IndexEntry{Name: k, Pages: v}
		count++
	}
	return ies
}

func (i Index) Alphabetical() []IndexEntry {
	name := func(i1, i2 *IndexEntry) bool {
		return i1.Name < i2.Name
	}

	ia := i.IndexArray()
	By(name).Sort(ia)
	return ia
}

func (i Index) ByCount() []IndexEntry {
	count := func(i1, i2 *IndexEntry) bool {
		return len(i1.Pages) < len(i2.Pages)
	}

	ia := i.IndexArray()
	By(count).Sort(ia)
	return ia
}

type IndexEntry struct {
	Name  string
	Pages IndexedPages
}

type By func(i1, i2 *IndexEntry) bool

func (by By) Sort(indexEntrys []IndexEntry) {
	ps := &indexEntrySorter{
		indexEntrys: indexEntrys,
		by:          by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type indexEntrySorter struct {
	indexEntrys []IndexEntry
	by          func(p1, p2 *IndexEntry) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *indexEntrySorter) Len() int {
	return len(s.indexEntrys)
}

// Swap is part of sort.Interface.
func (s *indexEntrySorter) Swap(i, j int) {
	s.indexEntrys[i], s.indexEntrys[j] = s.indexEntrys[j], s.indexEntrys[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *indexEntrySorter) Less(i, j int) bool {
	return s.by(&s.indexEntrys[i], &s.indexEntrys[j])
}
