// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"sync"
)

type pageCacheEntry struct {
	in  []Pages
	out Pages
}

func (entry pageCacheEntry) matches(pageLists []Pages) bool {
	if len(entry.in) != len(pageLists) {
		return false
	}
	for i, p := range pageLists {
		if !fastEqualPages(p, entry.in[i]) {
			return false
		}
	}

	return true
}

type pageCache struct {
	sync.RWMutex
	m map[string][]pageCacheEntry
}

func newPageCache() *pageCache {
	return &pageCache{m: make(map[string][]pageCacheEntry)}
}

// get/getP gets a Pages slice from the cache matching the given key and
// all the provided Pages slices.
// If none found in cache, a copy of the first slice is created.
//
// If an apply func is provided, that func is applied to the newly created copy.
//
// The getP variant' apply func takes a pointer to Pages.
//
// The cache and the execution of the apply func is protected by a RWMutex.
func (c *pageCache) get(key string, apply func(p Pages), pageLists ...Pages) (Pages, bool) {
	return c.getP(key, func(p *Pages) {
		if apply != nil {
			apply(*p)
		}
	}, pageLists...)
}

func (c *pageCache) getP(key string, apply func(p *Pages), pageLists ...Pages) (Pages, bool) {
	c.RLock()
	if cached, ok := c.m[key]; ok {
		for _, entry := range cached {
			if entry.matches(pageLists) {
				c.RUnlock()
				return entry.out, true
			}
		}
	}
	c.RUnlock()

	c.Lock()
	defer c.Unlock()

	// double-check
	if cached, ok := c.m[key]; ok {
		for _, entry := range cached {
			if entry.matches(pageLists) {
				return entry.out, true
			}
		}
	}

	p := pageLists[0]
	pagesCopy := append(Pages(nil), p...)

	if apply != nil {
		apply(&pagesCopy)
	}

	entry := pageCacheEntry{in: pageLists, out: pagesCopy}
	if v, ok := c.m[key]; ok {
		c.m[key] = append(v, entry)
	} else {
		c.m[key] = []pageCacheEntry{entry}
	}

	return pagesCopy, false

}

// "fast" as in: we do not compare every element for big slices, but that is
// good enough for our use cases.
// TODO(bep) there is a similar method in pagination.go. DRY.
func fastEqualPages(p1, p2 Pages) bool {
	if p1 == nil && p2 == nil {
		return true
	}

	if p1 == nil || p2 == nil {
		return false
	}

	if p1.Len() != p2.Len() {
		return false
	}

	if p1.Len() == 0 {
		return true
	}

	step := 1

	if len(p1) >= 50 {
		step = len(p1) / 10
	}

	for i := 0; i < len(p1); i += step {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}
