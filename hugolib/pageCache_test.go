package hugolib

import (
	"fmt"
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

	for j := 0; j < 50; j++ {
		testPageSets = append(testPageSets, createSortTestPages(j+1))
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j, pages := range testPageSets {
				msg := fmt.Sprintf("Go %d %d %d %d", i, j, o1, o2)
				l1.Lock()
				p, c := c1.get("k1", pages, nil)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o1, uint64(j), uint64(j+1)), c, "c1: "+msg)
				l1.Unlock()
				p2, c2 := c1.get("k1", p, nil)
				assert.True(t, c2)
				assert.True(t, probablyEqualPages(p, p2))
				assert.True(t, probablyEqualPages(p, pages))
				assert.NotNil(t, p, msg)

				l2.Lock()
				p3, c3 := c1.get("k2", pages, changeFirst)
				assert.Equal(t, !atomic.CompareAndSwapUint64(&o2, uint64(j), uint64(j+1)), c3, "c3: "+msg)
				l2.Unlock()
				assert.NotNil(t, p3, msg)
				assert.Equal(t, p3[0].Description, "changed", msg)
			}
		}()
	}

	wg.Wait()

}
