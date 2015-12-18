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
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPageCache(t *testing.T) {
	c1 := newPageCache()

	changeFirst := func(p Pages) {
		p[0].Description = "changed"
	}

	var o1 uint64 = 0
	var o2 uint64 = 0

	var wg sync.WaitGroup

	var l1 sync.Mutex
	var l2 sync.Mutex

	var testPageSets []Pages

	var i, j int

	for j = 0; j < 50; j++ {
		testPageSets = append(testPageSets, createSortTestPages(j+1))
	}

	for i = 0; i < 100; i++ {
		wg.Add(1)
		go func(i1, i2 int) {
			defer wg.Done()
			for j, pages := range testPageSets {
				l1.Lock()
				p, c := c1.get("k1", pages, nil)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o1, uint64(j), uint64(j+1)), c)
				l1.Unlock()
				p2, c2 := c1.get("k1", p, nil)
				assert.True(t, c2)
				assert.True(t, probablyEqualPages(p, p2))
				assert.True(t, probablyEqualPages(p, pages))
				assert.NotNil(t, p)

				l2.Lock()
				p3, c3 := c1.get("k2", pages, changeFirst)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o2, uint64(j), uint64(j+1)), c3)
				l2.Unlock()
				assert.NotNil(t, p3)
				assert.Equal(t, p3[0].Description, "changed")
			}
		}(i, j)
	}

	wg.Wait()

}
