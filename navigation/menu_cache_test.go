// Copyright 2021 The Hugo Authors. All rights reserved.
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

package navigation

import (
	"sync"
	"sync/atomic"
	"testing"

	qt "github.com/frankban/quicktest"
)

func createSortTestMenu(num int) Menu {
	menu := make(Menu, num)
	for i := 0; i < num; i++ {
		m := &MenuEntry{}
		menu[i] = m
	}
	return menu
}

func TestMenuCache(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	c1 := newMenuCache()

	changeFirst := func(m Menu) {
		m[0].MenuConfig.Title = "changed"
	}

	var o1 uint64
	var o2 uint64

	var wg sync.WaitGroup

	var l1 sync.Mutex
	var l2 sync.Mutex

	var testMenuSets []Menu

	for i := 0; i < 50; i++ {
		testMenuSets = append(testMenuSets, createSortTestMenu(i+1))
	}

	for j := 0; j < 100; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k, menu := range testMenuSets {
				l1.Lock()
				m, ca := c1.get("k1", nil, menu)
				c.Assert(ca, qt.Equals, !atomic.CompareAndSwapUint64(&o1, uint64(k), uint64(k+1)))
				l1.Unlock()
				m2, c2 := c1.get("k1", nil, m)
				c.Assert(c2, qt.Equals, true)
				c.Assert(menuEqual(m, m2), qt.Equals, true)
				c.Assert(menuEqual(m, menu), qt.Equals, true)
				c.Assert(m, qt.Not(qt.IsNil))

				l2.Lock()
				m3, c3 := c1.get("k2", changeFirst, menu)
				c.Assert(c3, qt.Equals, !atomic.CompareAndSwapUint64(&o2, uint64(k), uint64(k+1)))
				l2.Unlock()
				c.Assert(m3, qt.Not(qt.IsNil))
				c.Assert("changed", qt.Equals, m3[0].Title)
			}
		}()
	}
	wg.Wait()
}
