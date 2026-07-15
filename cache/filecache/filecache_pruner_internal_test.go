// Copyright 2026 The Hugo Authors. All rights reserved.
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

package filecache

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestPruneSeenFileCaseSensitiveFilesystem(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	cache := newPruneSeenTestCache(false)
	writePruneSeenTestFile(c, cache, "MyBundle/test.jpeg")
	markPruneSeenTestFile(cache, "mybundle/test.jpeg")

	count, err := cache.Prune(false)
	c.Assert(err, qt.IsNil)
	c.Assert(count, qt.Equals, 1)
	c.Assert(cache.GetString(filepath.FromSlash("MyBundle/test.jpeg")), qt.Equals, "")
}

// See issue 15101.
func TestPruneSeenFileCaseInsensitiveFilesystem(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	cache := newPruneSeenTestCache(true)
	writePruneSeenTestFile(c, cache, "MyBundle/test.jpeg")
	markPruneSeenTestFile(cache, "mybundle/test.jpeg")

	count, err := cache.Prune(false)
	c.Assert(err, qt.IsNil)
	c.Assert(count, qt.Equals, 0)
	c.Assert(cache.GetString(filepath.FromSlash("MyBundle/test.jpeg")), qt.Equals, "abc")
}

func newPruneSeenTestCache(caseInsensitiveFilesystem bool) *Cache {
	cache := NewCache(afero.NewMemMapFs(), FileCacheConfig{
		MaxAge: -1,
		Dir:    "cache/c",
	})
	cache.caseInsensitiveFilesystemInit.Do(func() {
		cache.caseInsensitiveFilesystem = caseInsensitiveFilesystem
	})
	return cache
}

func writePruneSeenTestFile(c *qt.C, cache *Cache, name string) {
	filename := filepath.FromSlash(name)
	c.Assert(cache.Fs.MkdirAll(filepath.Dir(filename), 0o777), qt.IsNil)
	c.Assert(afero.WriteFile(cache.Fs, filename, []byte("abc"), 0o666), qt.IsNil)
}

func markPruneSeenTestFile(cache *Cache, name string) {
	id := cleanID(filepath.FromSlash(name))
	cache.entryLocker.Lock(id)
	cache.entryLocker.Unlock(id)
}
