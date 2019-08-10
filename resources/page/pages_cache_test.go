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
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPageCache(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	c1 := newPageCache()

	changeFirst := func(p Pages) {
		p[0].(*testPage).description = "changed"
	}

	var o1 uint64
	var o2 uint64

	var wg sync.WaitGroup

	var l1 sync.Mutex
	var l2 sync.Mutex

	var testPageSets []Pages

	for i := 0; i < 50; i++ {
		testPageSets = append(testPageSets, createSortTestPages(i+1))
	}

	for j := 0; j < 100; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k, pages := range testPageSets {
				l1.Lock()
				p, ca := c1.get("k1", nil, pages)
				c.Assert(ca, qt.Equals, !atomic.CompareAndSwapUint64(&o1, uint64(k), uint64(k+1)))
				l1.Unlock()
				p2, c2 := c1.get("k1", nil, p)
				c.Assert(c2, qt.Equals, true)
				c.Assert(pagesEqual(p, p2), qt.Equals, true)
				c.Assert(pagesEqual(p, pages), qt.Equals, true)
				c.Assert(p, qt.Not(qt.IsNil))

				l2.Lock()
				p3, c3 := c1.get("k2", changeFirst, pages)
				c.Assert(c3, qt.Equals, !atomic.CompareAndSwapUint64(&o2, uint64(k), uint64(k+1)))
				l2.Unlock()
				c.Assert(p3, qt.Not(qt.IsNil))
				c.Assert("changed", qt.Equals, p3[0].(*testPage).description)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkPageCache(b *testing.B) {
	cache := newPageCache()
	pages := make(Pages, 30)
	for i := 0; i < 30; i++ {
		pages[i] = &testPage{title: "p" + strconv.Itoa(i)}
	}
	key := "key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.getP(key, nil, pages)
	}
}
