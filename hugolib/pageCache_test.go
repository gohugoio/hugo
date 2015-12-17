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
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
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

	for j := 0; j < 50; j++ {
		testPageSets = append(testPageSets, createSortTestPages(j+1))
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j, pages := range testPageSets {
				l1.Lock()
				l2.Lock()
				msg := fmt.Sprintf("Go %d %d %d %d", i, j, o1, o2)
				p, c := c1.get("k1", pages, nil)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o1, uint64(j), uint64(j+1)), c, "c1: "+msg)
				l1.Unlock()
				p2, c2 := c1.get("k1", p, nil)
				assert.True(t, c2)
				assert.True(t, probablyEqualPages(p, p2))
				assert.True(t, probablyEqualPages(p, pages))
				assert.NotNil(t, p, msg)

				p3, c3 := c1.get("k2", pages, changeFirst)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o2, uint64(j), uint64(j+1)), c3, "c3: "+msg)
				l2.Unlock()
				assert.NotNil(t, p3, msg)
				assert.Equal(t, p3[0].Description, "changed", msg)
			}
		}(i)
	}

	wg.Wait()

}
