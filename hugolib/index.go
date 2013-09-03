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

type IndexCount struct {
	Name  string
	Count int
}

type Index map[string]Pages
type IndexList map[string]Index

type OrderedIndex []IndexCount
type OrderedIndexList map[string]OrderedIndex

// KeyPrep... Indexes should be case insensitive. Can make it easily conditional later.
func kp(in string) string {
	return template.Urlize(in)
}

func (i Index) Get(key string) Pages { return i[kp(key)] }
func (i Index) Count(key string) int { return len(i[kp(key)]) }
func (i Index) Add(key string, p *Page) {
	key = kp(key)
	i[key] = append(i[key], p)
}

func (l IndexList) BuildOrderedIndexList() OrderedIndexList {
	oil := make(OrderedIndexList, len(l))
	for idx_name, index := range l {
		i := 0
		oi := make(OrderedIndex, len(index))
		for name, pages := range index {
			oi[i] = IndexCount{name, len(pages)}
			i++
		}
		sort.Sort(oi)
		oil[idx_name] = oi
	}
	return oil
}

func (idx OrderedIndex) Len() int           { return len(idx) }
func (idx OrderedIndex) Less(i, j int) bool { return idx[i].Count > idx[j].Count }
func (idx OrderedIndex) Swap(i, j int)      { idx[i], idx[j] = idx[j], idx[i] }
