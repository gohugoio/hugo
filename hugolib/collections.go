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
	"fmt"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/common/collections"
)

var (
	_ collections.Grouper         = (*Page)(nil)
	_ collections.Slicer          = (*Page)(nil)
	_ collections.Slicer          = PageGroup{}
	_ collections.Slicer          = WeightedPage{}
	_ resource.ResourcesConverter = Pages{}
)

// collections.Slicer implementations below. We keep these bridge implementations
// here as it makes it easier to get an idea of "type coverage". These
// implementations have no value on their own.

// Slice is not meant to be used externally. It's a bridge function
// for the template functions. See collections.Slice.
func (p *Page) Slice(items interface{}) (interface{}, error) {
	return toPages(items)
}

// Slice is not meant to be used externally. It's a bridge function
// for the template functions. See collections.Slice.
func (p PageGroup) Slice(in interface{}) (interface{}, error) {
	switch items := in.(type) {
	case PageGroup:
		return items, nil
	case []interface{}:
		groups := make(PagesGroup, len(items))
		for i, v := range items {
			g, ok := v.(PageGroup)
			if !ok {
				return nil, fmt.Errorf("type %T is not a PageGroup", v)
			}
			groups[i] = g
		}
		return groups, nil
	default:
		return nil, fmt.Errorf("invalid slice type %T", items)
	}
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

// collections.Grouper  implementations below

// Group creates a PageGroup from a key and a Pages object
// This method is not meant for external use. It got its non-typed arguments to satisfy
// a very generic interface in the tpl package.
func (p *Page) Group(key interface{}, in interface{}) (interface{}, error) {
	pages, err := toPages(in)
	if err != nil {
		return nil, err
	}
	return PageGroup{Key: key, Pages: pages}, nil
}

// ToResources wraps resource.ResourcesConverter
func (pages Pages) ToResources() resource.Resources {
	r := make(resource.Resources, len(pages))
	for i, p := range pages {
		r[i] = p
	}
	return r
}
