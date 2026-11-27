// Copyright 2024 The Hugo Authors. All rights reserved.
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

package dynacache

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	_ resource.StaleInfo = (*testItem)(nil)
	_ identity.Identity  = (*testItem)(nil)
)

type testItem struct {
	name         string
	staleVersion uint32
}

func (t testItem) StaleVersion() uint32 {
	return t.staleVersion
}

func (t testItem) IdentifierBase() string {
	return t.name
}

func TestCache(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	cache := New(Options{
		Log: loggers.NewDefault(),
	})

	c.Cleanup(func() {
		cache.Stop()
	})

	opts := OptionsPartition{Weight: 30}

	c.Assert(cache, qt.Not(qt.IsNil))

	p1 := GetOrCreatePartition[string, testItem](cache, "/aaaa/bbbb", opts)
	c.Assert(p1, qt.Not(qt.IsNil))

	p2 := GetOrCreatePartition[string, testItem](cache, "/aaaa/bbbb", opts)

	c.Assert(func() { GetOrCreatePartition[string, testItem](cache, "foo bar", opts) }, qt.PanicMatches, ".*invalid partition name.*")
	c.Assert(func() { GetOrCreatePartition[string, testItem](cache, "/aaaa/cccc", OptionsPartition{Weight: 1234}) }, qt.PanicMatches, ".*invalid Weight.*")

	c.Assert(p2, qt.Equals, p1)

	p3 := GetOrCreatePartition[string, testItem](cache, "/aaaa/cccc", opts)
	c.Assert(p3, qt.Not(qt.IsNil))
	c.Assert(p3, qt.Not(qt.Equals), p1)

	c.Assert(func() { New(Options{}) }, qt.PanicMatches, ".*nil Log.*")
}

func TestCleanKey(t *testing.T) {
	c := qt.New(t)

	c.Assert(CleanKey("a/b/c"), qt.Equals, "/a/b/c")
	c.Assert(CleanKey("/a/b/c"), qt.Equals, "/a/b/c")
	c.Assert(CleanKey("a/b/c/"), qt.Equals, "/a/b/c")
	c.Assert(CleanKey(filepath.FromSlash("/a/b/c/")), qt.Equals, "/a/b/c")
}

func newTestCache(t *testing.T) *Cache {
	cache := New(
		Options{
			Log: loggers.NewDefault(),
		},
	)

	p1 := GetOrCreatePartition[string, testItem](cache, "/aaaa/bbbb", OptionsPartition{Weight: 30, ClearWhen: ClearOnRebuild})
	p2 := GetOrCreatePartition[string, testItem](cache, "/aaaa/cccc", OptionsPartition{Weight: 30, ClearWhen: ClearOnChange})

	p1.GetOrCreate("clearOnRebuild", func(string) (testItem, error) {
		return testItem{}, nil
	})

	p2.GetOrCreate("clearBecauseStale", func(string) (testItem, error) {
		return testItem{
			staleVersion: 32,
		}, nil
	})

	p2.GetOrCreate("clearBecauseIdentityChanged", func(string) (testItem, error) {
		return testItem{
			name: "changed",
		}, nil
	})

	p2.GetOrCreate("clearNever", func(string) (testItem, error) {
		return testItem{
			staleVersion: 0,
		}, nil
	})

	t.Cleanup(func() {
		cache.Stop()
	})

	return cache
}

func TestClear(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	predicateAll := func(string) bool {
		return true
	}

	cache := newTestCache(t)

	c.Assert(cache.Keys(predicateAll), qt.HasLen, 4)

	cache.ClearOnRebuild(nil)

	// Stale items are always cleared.
	c.Assert(cache.Keys(predicateAll), qt.HasLen, 2)

	cache = newTestCache(t)
	cache.ClearOnRebuild(nil, identity.StringIdentity("changed"))

	c.Assert(cache.Keys(nil), qt.HasLen, 1)

	cache = newTestCache(t)

	cache.ClearMatching(nil, func(k, v any) bool {
		return k.(string) == "clearOnRebuild"
	})

	c.Assert(cache.Keys(predicateAll), qt.HasLen, 3)

	cache.checkMemory()
}

func TestPanicInCreate(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	cache := newTestCache(t)

	p1 := GetOrCreatePartition[string, testItem](cache, "/aaaa/bbbb", OptionsPartition{Weight: 30, ClearWhen: ClearOnRebuild})

	willPanic := func(i int) func() {
		return func() {
			p1.GetOrCreate(fmt.Sprintf("panic-%d", i), func(key string) (testItem, error) {
				panic(errors.New(key))
			})
		}
	}

	// GetOrCreateWitTimeout needs to recover from panics in the create func.
	willErr := func(i int) error {
		_, err := p1.GetOrCreateWitTimeout(fmt.Sprintf("error-%d", i), 10*time.Second, func(key string) (testItem, error) {
			return testItem{}, errors.New(key)
		})
		return err
	}

	for i := range 3 {
		for range 3 {
			c.Assert(willPanic(i), qt.PanicMatches, fmt.Sprintf("panic-%d", i))
			c.Assert(willErr(i), qt.ErrorMatches, fmt.Sprintf("error-%d", i))
		}
	}

	// Test the same keys again without the panic.
	for i := range 3 {
		for range 3 {
			v, err := p1.GetOrCreate(fmt.Sprintf("panic-%d", i), func(key string) (testItem, error) {
				return testItem{
					name: key,
				}, nil
			})
			c.Assert(err, qt.IsNil)
			c.Assert(v.name, qt.Equals, fmt.Sprintf("panic-%d", i))

			v, err = p1.GetOrCreateWitTimeout(fmt.Sprintf("error-%d", i), 10*time.Second, func(key string) (testItem, error) {
				return testItem{
					name: key,
				}, nil
			})
			c.Assert(err, qt.IsNil)
			c.Assert(v.name, qt.Equals, fmt.Sprintf("error-%d", i))
		}
	}
}

func TestAdjustMaxSize(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	const limit uint64 = 1000

	cache := New(Options{
		Log:           loggers.NewDefault(),
		CheckInterval: time.Hour,
	})
	t.Cleanup(func() {
		cache.Stop()
	})
	cache.availableMemory = limit

	pLow := GetOrCreatePartition[string, int](cache, "/aaaa/low", OptionsPartition{Weight: 30})
	pHigh := GetOrCreatePartition[string, int](cache, "/aaaa/high", OptionsPartition{Weight: 70})
	pEmpty := GetOrCreatePartition[string, int](cache, "/aaaa/empty", OptionsPartition{Weight: 30})

	c.Assert(pLow.getMaxSize(), qt.Equals, 30*maxEntriesPerWeight)
	c.Assert(pHigh.getMaxSize(), qt.Equals, 70*maxEntriesPerWeight)
	c.Assert(pEmpty.getMaxSize(), qt.Equals, 30*maxEntriesPerWeight)

	for i := range 100 {
		key := fmt.Sprintf("k%d", i)
		_, err := pLow.GetOrCreate(key, func(string) (int, error) { return i, nil })
		c.Assert(err, qt.IsNil)
		_, err = pHigh.GetOrCreate(key, func(string) (int, error) { return i, nil })
		c.Assert(err, qt.IsNil)
	}

	// Over limit: shrinks immediately from actual occupancy (100), weighted by Weight.
	// Weight 30 -> retain 0.40 -> 40; Weight 70 -> retain 0.60 -> 60; empty untouched.
	cache.adjustMaxSize(runtime.MemStats{Alloc: limit + 1})
	c.Assert(pLow.getMaxSize(), qt.Equals, 40)
	c.Assert(pLow.c.Len(), qt.Equals, 40)
	c.Assert(pHigh.getMaxSize(), qt.Equals, 60)
	c.Assert(pHigh.c.Len(), qt.Equals, 60)
	c.Assert(pEmpty.getMaxSize(), qt.Equals, 30*maxEntriesPerWeight)

	// Deadband [80%, 100%] of limit: holds steady without shrinking or growing.
	cache.adjustMaxSize(runtime.MemStats{Alloc: limit * 9 / 10})
	c.Assert(pLow.getMaxSize(), qt.Equals, 40)
	c.Assert(pHigh.getMaxSize(), qt.Equals, 60)

	// Below 80% of limit: grows smoothly by 1.5x.
	cache.adjustMaxSize(runtime.MemStats{Alloc: limit / 2})
	c.Assert(pLow.getMaxSize(), qt.Equals, 60)
	c.Assert(pHigh.getMaxSize(), qt.Equals, 90)

	// Repeated over-limit ticks shrink all the way down to minMaxSize.
	for range 20 {
		cache.adjustMaxSize(runtime.MemStats{Alloc: limit + 1})
	}
	c.Assert(pLow.getMaxSize(), qt.Equals, minMaxSize)
	c.Assert(pLow.c.Len(), qt.Equals, minMaxSize)
	c.Assert(pHigh.getMaxSize(), qt.Equals, minMaxSize)
	c.Assert(pHigh.c.Len(), qt.Equals, minMaxSize)

	// Repeated healthy ticks recover back to defaultMaxSize without overshooting.
	for range 100 {
		cache.adjustMaxSize(runtime.MemStats{Alloc: limit / 2})
	}
	c.Assert(pLow.getMaxSize(), qt.Equals, 30*maxEntriesPerWeight)
	c.Assert(pHigh.getMaxSize(), qt.Equals, 70*maxEntriesPerWeight)
}
