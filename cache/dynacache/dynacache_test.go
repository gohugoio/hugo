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
	"path/filepath"
	"testing"

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

func TestCalculateMaxSizePerPartition(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Assert(calculateMaxSizePerPartition(1000, 500, 5), qt.Equals, 200)
	c.Assert(calculateMaxSizePerPartition(1000, 250, 5), qt.Equals, 400)
	c.Assert(func() { calculateMaxSizePerPartition(1000, 250, 0) }, qt.PanicMatches, ".*must be > 0.*")
	c.Assert(func() { calculateMaxSizePerPartition(1000, 0, 1) }, qt.PanicMatches, ".*must be > 0.*")
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

	cache.ClearOnRebuild()

	// Stale items are always cleared.
	c.Assert(cache.Keys(predicateAll), qt.HasLen, 2)

	cache = newTestCache(t)
	cache.ClearOnRebuild(identity.StringIdentity("changed"))

	c.Assert(cache.Keys(nil), qt.HasLen, 1)

	cache = newTestCache(t)

	cache.ClearMatching(nil, func(k, v any) bool {
		return k.(string) == "clearOnRebuild"
	})

	c.Assert(cache.Keys(predicateAll), qt.HasLen, 3)

	cache.adjustCurrentMaxSize()
}

func TestAdjustCurrentMaxSize(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	cache := newTestCache(t)
	alloc := cache.stats.memstatsCurrent.Alloc
	cache.adjustCurrentMaxSize()
	c.Assert(cache.stats.memstatsCurrent.Alloc, qt.Not(qt.Equals), alloc)
}
