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
	"math/rand"

	"github.com/gohugoio/hugo/compare"

	"github.com/gohugoio/hugo/resources/resource"
)

var (
	_ resource.ResourcesConverter = Pages{}
	_ compare.ProbablyEqer        = Pages{}
)

// Pages is a slice of pages. This is the most common list type in Hugo.
type Pages []Page

func (ps Pages) String() string {
	return fmt.Sprintf("Pages(%d)", len(ps))
}

// Used in tests.
func (ps Pages) shuffle() {
	for i := range ps {
		j := rand.Intn(i + 1)
		ps[i], ps[j] = ps[j], ps[i]
	}
}

// ToResources wraps resource.ResourcesConverter
func (pages Pages) ToResources() resource.Resources {
	r := make(resource.Resources, len(pages))
	for i, p := range pages {
		r[i] = p
	}
	return r
}

// ToPages tries to convert seq into Pages.
func ToPages(seq interface{}) (Pages, error) {
	if seq == nil {
		return Pages{}, nil
	}

	switch v := seq.(type) {
	case Pages:
		return v, nil
	case *Pages:
		return *(v), nil
	case WeightedPages:
		return v.Pages(), nil
	case PageGroup:
		return v.Pages, nil
	case []Page:
		pages := make(Pages, len(v))
		for i, vv := range v {
			pages[i] = vv
		}
		return pages, nil
	case []interface{}:
		pages := make(Pages, len(v))
		success := true
		for i, vv := range v {
			p, ok := vv.(Page)
			if !ok {
				success = false
				break
			}
			pages[i] = p
		}
		if success {
			return pages, nil
		}
	}

	return nil, fmt.Errorf("cannot convert type %T to Pages", seq)
}

func (p Pages) Group(key interface{}, in interface{}) (interface{}, error) {
	pages, err := ToPages(in)
	if err != nil {
		return nil, err
	}
	return PageGroup{Key: key, Pages: pages}, nil
}

// Len returns the number of pages in the list.
func (p Pages) Len() int {
	return len(p)
}

// ProbablyEq wraps comare.ProbablyEqer
func (pages Pages) ProbablyEq(other interface{}) bool {
	otherPages, ok := other.(Pages)
	if !ok {
		return false
	}

	if len(pages) != len(otherPages) {
		return false
	}

	step := 1

	for i := 0; i < len(pages); i += step {
		if !pages[i].Eq(otherPages[i]) {
			return false
		}

		if i > 50 {
			// This is most likely the same.
			step = 50
		}
	}

	return true
}

func (ps Pages) removeFirstIfFound(p Page) Pages {
	ii := -1
	for i, pp := range ps {
		if p.Eq(pp) {
			ii = i
			break
		}
	}

	if ii != -1 {
		ps = append(ps[:ii], ps[ii+1:]...)
	}
	return ps
}

// PagesFactory somehow creates some Pages.
// We do a lot of lazy Pages initialization in Hugo, so we need a type.
type PagesFactory func() Pages
