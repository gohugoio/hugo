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
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageCache(t *testing.T) {
	t.Parallel()
	c1 := newPageCache()

	changeFirst := func(p Pages) {
		p[0].Description = "changed"
	}

	var o1 uint64
	var o2 uint64

	var wg sync.WaitGroup

	var l1 sync.Mutex
	var l2 sync.Mutex

	var testPageSets []Pages

	s := newTestSite(t)

	for i := 0; i < 50; i++ {
		testPageSets = append(testPageSets, createSortTestPages(s, i+1))
	}

	for j := 0; j < 100; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k, pages := range testPageSets {
				l1.Lock()
				p, c := c1.get("k1", pages, nil)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o1, uint64(k), uint64(k+1)), c)
				l1.Unlock()
				p2, c2 := c1.get("k1", p, nil)
				assert.True(t, c2)
				assert.True(t, fastEqualPages(p, p2))
				assert.True(t, fastEqualPages(p, pages))
				assert.NotNil(t, p)

				l2.Lock()
				p3, c3 := c1.get("k2", pages, changeFirst)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o2, uint64(k), uint64(k+1)), c3)
				l2.Unlock()
				assert.NotNil(t, p3)
				assert.Equal(t, p3[0].Description, "changed")
			}
		}()
	}
	wg.Wait()
}
