// Copyright 2020 The Hugo Authors. All rights reserved.
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

package memcache

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestCache(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	cache := New(Config{})

	counter := 0
	create := func() Entry {
		counter++
		return Entry{Value: counter, ClearWhen: ClearOnChange}
	}

	a := cache.GetOrCreatePartition("a", ClearNever)

	for i := 0; i < 5; i++ {
		v1, err := a.GetOrCreate(context.TODO(), "a1", create)
		c.Assert(err, qt.IsNil)
		c.Assert(v1, qt.Equals, 1)
		v2, err := a.GetOrCreate(context.TODO(), "a2", create)
		c.Assert(err, qt.IsNil)
		c.Assert(v2, qt.Equals, 2)
	}

	cache.Clear()

	v3, err := a.GetOrCreate(context.TODO(), "a2", create)
	c.Assert(err, qt.IsNil)
	c.Assert(v3, qt.Equals, 3)
}

func TestCacheConcurrent(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	var wg sync.WaitGroup

	cache := New(Config{})

	create := func(i int) func() Entry {
		return func() Entry {
			return Entry{Value: i, ClearWhen: ClearOnChange}
		}
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				id := fmt.Sprintf("id%d", j)
				v, err := cache.getOrCreate("a", id, create(j))
				c.Assert(err, qt.IsNil)
				c.Assert(v, qt.Equals, j)
			}
		}()
	}
	wg.Wait()
}

func TestCacheMemStats(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	cache := New(Config{
		ItemsToPrune:  10,
		CheckInterval: 500 * time.Millisecond,
	})

	s := cache.stats

	c.Assert(s.memstatsStart.Alloc > 0, qt.Equals, true)
	c.Assert(s.memstatsCurrent.Alloc, qt.Equals, uint64(0))
	c.Assert(s.availableMemory > 0, qt.Equals, true)
	c.Assert(s.numItems, qt.Equals, uint64(0))

	counter := 0
	create := func() Entry {
		counter++
		return Entry{Value: counter, ClearWhen: ClearNever}
	}

	for i := 1; i <= 20; i++ {
		_, err := cache.getOrCreate("a", fmt.Sprintf("b%d", i), create)
		c.Assert(err, qt.IsNil)
	}

	c.Assert(s.getNumItems(), qt.Equals, uint64(20))
	cache.cache.SetMaxSize(10)
	time.Sleep(time.Millisecond * 1200)
	c.Assert(int(s.getNumItems()), qt.Equals, 10)
}

func TestSplitBasePathAndExt(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	tests := []struct {
		path string
		a    string
		b    string
	}{
		{"a/b.json", "a", "json"},
		{"a/b/c/d.json", "a", "json"},
	}
	for i, this := range tests {
		msg := qt.Commentf("test %d", i)
		a, b := splitBasePathAndExt(this.path)

		c.Assert(a, qt.Equals, this.a, msg)
		c.Assert(b, qt.Equals, this.b, msg)
	}
}

func TestCleanKey(t *testing.T) {
	c := qt.New(t)

	c.Assert(CleanKey(filepath.FromSlash("a/b/c.js")), qt.Equals, "a/b/c.js")
	c.Assert(CleanKey("a//b////c.js"), qt.Equals, "a/b/c.js")
	c.Assert(CleanKey("a.js"), qt.Equals, "_root/a.js")
	c.Assert(CleanKey("b/a"), qt.Equals, "b/a.unkn")
}

func TestKeyValid(t *testing.T) {
	c := qt.New(t)

	c.Assert(keyValid("a/b.j"), qt.Equals, true)
	c.Assert(keyValid("a/b."), qt.Equals, false)
	c.Assert(keyValid("a/b"), qt.Equals, false)
	c.Assert(keyValid("/a/b.txt"), qt.Equals, false)
	c.Assert(keyValid("a\\b.js"), qt.Equals, false)
}

func TestInsertKeyPathElement(t *testing.T) {
	c := qt.New(t)

	c.Assert(InsertKeyPathElements("a/b.j", "en"), qt.Equals, "a/en/b.j")
	c.Assert(InsertKeyPathElements("a/b.j", "en", "foo"), qt.Equals, "a/en/foo/b.j")
	c.Assert(InsertKeyPathElements("a/b.j", "", "foo"), qt.Equals, "a/foo/b.j")
}

func TestShouldEvict(t *testing.T) {
	// TODO1 remove?
	// c := qt.New(t)

	// fmt.Println("=>", CleanKey("kkk"))
	// c.Assert(shouldEvict("key", Entry{}, ClearNever, identity.NewPathIdentity(files.ComponentFolderAssets, "a/b/c.js")), qt.Equals, true)
}
