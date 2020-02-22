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

package para

import (
	"context"
	"runtime"

	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestPara(t *testing.T) {
	if runtime.NumCPU() < 4 {
		t.Skipf("skip para test, CPU count is %d", runtime.NumCPU())
	}

	c := qt.New(t)

	c.Run("Order", func(c *qt.C) {
		n := 500
		ints := make([]int, n)
		for i := 0; i < n; i++ {
			ints[i] = i
		}

		p := New(4)
		r, _ := p.Start(context.Background())

		var result []int
		var mu sync.Mutex
		for i := 0; i < n; i++ {
			i := i
			r.Run(func() error {
				mu.Lock()
				defer mu.Unlock()
				result = append(result, i)
				return nil
			})
		}

		c.Assert(r.Wait(), qt.IsNil)
		c.Assert(result, qt.HasLen, len(ints))
		c.Assert(sort.IntsAreSorted(result), qt.Equals, false, qt.Commentf("Para does not seem to be parallel"))
		sort.Ints(result)
		c.Assert(result, qt.DeepEquals, ints)

	})

	c.Run("Time", func(c *qt.C) {
		const n = 100

		p := New(5)
		r, _ := p.Start(context.Background())

		start := time.Now()

		var counter int64

		for i := 0; i < n; i++ {
			r.Run(func() error {
				atomic.AddInt64(&counter, 1)
				time.Sleep(1 * time.Millisecond)
				return nil
			})
		}

		c.Assert(r.Wait(), qt.IsNil)
		c.Assert(counter, qt.Equals, int64(n))
		c.Assert(time.Since(start) < n/2*time.Millisecond, qt.Equals, true)

	})

}
