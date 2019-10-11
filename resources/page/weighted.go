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
	"fmt"
	"sort"

	"github.com/gohugoio/hugo/common/collections"
)

var (
	_ collections.Slicer = WeightedPage{}
)

// WeightedPages is a list of Pages with their corresponding (and relative) weight
// [{Weight: 30, Page: *1}, {Weight: 40, Page: *2}]
type WeightedPages []WeightedPage

// Page will return the Page (of Kind taxonomyList) that represents this set
// of pages. This method will panic if p is empty, as that should never happen.
func (p WeightedPages) Page() Page {
	if len(p) == 0 {
		panic("WeightedPages is empty")
	}

	first := p[0]

	// TODO(bep) fix tests
	if first.owner == nil {
		return nil
	}

	return first.owner
}

// A WeightedPage is a Page with a weight.
type WeightedPage struct {
	Weight int
	Page

	// Reference to the owning Page. This avoids having to do
	// manual .Site.GetPage lookups. It is implemented in this roundabout way
	// because we cannot add additional state to the WeightedPages slice
	// without breaking lots of templates in the wild.
	owner Page
}

func NewWeightedPage(weight int, p Page, owner Page) WeightedPage {
	return WeightedPage{Weight: weight, Page: p, owner: owner}
}

func (w WeightedPage) String() string {
	return fmt.Sprintf("WeightedPage(%d,%q)", w.Weight, w.Page.Title())
}

// Slice is not meant to be used externally. It's a bridge function
// for the template functions. See collections.Slice.
func (p WeightedPage) Slice(in interface{}) (interface{}, error) {
	switch items := in.(type) {
	case WeightedPages:
		return items, nil
	case []interface{}:
		weighted := make(WeightedPages, len(items))
		for i, v := range items {
			g, ok := v.(WeightedPage)
			if !ok {
				return nil, fmt.Errorf("type %T is not a WeightedPage", v)
			}
			weighted[i] = g
		}
		return weighted, nil
	default:
		return nil, fmt.Errorf("invalid slice type %T", items)
	}
}

// Pages returns the Pages in this weighted page set.
func (wp WeightedPages) Pages() Pages {
	pages := make(Pages, len(wp))
	for i := range wp {
		pages[i] = wp[i].Page
	}
	return pages
}

// Next returns the next Page relative to the given Page in
// this weighted page set.
func (wp WeightedPages) Next(cur Page) Page {
	for x, c := range wp {
		if c.Page.Eq(cur) {
			if x == 0 {
				return nil
			}
			return wp[x-1].Page
		}
	}
	return nil
}

// Prev returns the previous Page relative to the given Page in
// this weighted page set.
func (wp WeightedPages) Prev(cur Page) Page {
	for x, c := range wp {
		if c.Page.Eq(cur) {
			if x < len(wp)-1 {
				return wp[x+1].Page
			}
			return nil
		}
	}
	return nil
}

func (wp WeightedPages) Len() int      { return len(wp) }
func (wp WeightedPages) Swap(i, j int) { wp[i], wp[j] = wp[j], wp[i] }

// Sort stable sorts this weighted page set.
func (wp WeightedPages) Sort() { sort.Stable(wp) }

// Count returns the number of pages in this weighted page set.
func (wp WeightedPages) Count() int { return len(wp) }

func (wp WeightedPages) Less(i, j int) bool {
	if wp[i].Weight == wp[j].Weight {
		return DefaultPageSort(wp[i].Page, wp[j].Page)
	}
	return wp[i].Weight < wp[j].Weight
}
