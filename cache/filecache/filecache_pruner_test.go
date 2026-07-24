// Copyright 2018 The Hugo Authors. All rights reserved.
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

package filecache_test

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/htesting"
	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

// A cache entry created before Hugo started lowercasing content paths in v0.123
// (e.g. _gen/images/MyBundle) is on a case-insensitive filesystem the same file as
// the lowercased cache key used today, and must not be pruned.
// See issue 15101.
func TestPruneCacheEntryWithOtherCase(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	dir := t.TempDir()
	if isCaseInsensitive, err := htesting.IsCaseInsensitiveFs(dir); err != nil {
		t.Fatal(err)
	} else if !isCaseInsensitive {
		t.Skip("skip test on case-sensitive filesystem")
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), dir)
	newCache := func() *filecache.Cache {
		return filecache.NewCache(fs, filecache.FileCacheConfig{Dir: "cache", MaxAge: -1})
	}

	c.Assert(newCache().SetBytes("MyBundle/i1", []byte("abc")), qt.IsNil)

	cache := newCache()
	_, b, err := cache.GetOrCreateBytes("mybundle/i1", func() ([]byte, error) {
		return []byte("def"), nil
	})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "abc")

	count, err := cache.Prune(false)
	c.Assert(err, qt.IsNil)
	c.Assert(count, qt.Equals, 0)
	c.Assert(cache.GetString("MyBundle/i1"), qt.Equals, "abc")
}

// On a case-sensitive filesystem the entries above are distinct files,
// and the one not used in this build should be pruned.
func TestPruneCacheEntryWithOtherCaseCaseSensitiveFs(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	fs := afero.NewMemMapFs()
	cache := filecache.NewCache(fs, filecache.FileCacheConfig{Dir: "cache", MaxAge: -1})

	c.Assert(cache.SetBytes("MyBundle/i1", []byte("abc")), qt.IsNil)

	cache = filecache.NewCache(fs, filecache.FileCacheConfig{Dir: "cache", MaxAge: -1})
	_, b, err := cache.GetOrCreateBytes("mybundle/i1", func() ([]byte, error) {
		return []byte("def"), nil
	})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "def")

	count, err := cache.Prune(false)
	c.Assert(err, qt.IsNil)
	c.Assert(count, qt.Equals, 1)
	c.Assert(cache.GetString("MyBundle/i1"), qt.Equals, "")
	c.Assert(cache.GetString("mybundle/i1"), qt.Equals, "def")
}

func TestPrune(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		c := qt.New(t)

		configStr := `
resourceDir = "myresources"
contentDir = "content"
dataDir = "data"
i18nDir = "i18n"
layoutDir = "layouts"
assetDir = "assets"
archeTypedir = "archetypes"

[caches]
[caches.misc]
maxAge = "200ms"
dir = "/cache/c"
[caches.assets]
maxAge = "200ms"
dir = ":resourceDir/_gen"
[caches.images]
maxAge = "200ms"
dir = ":resourceDir/_gen"
`

		for _, name := range []string{filecache.CacheKeyAssets, filecache.CacheKeyImages} {
			msg := qt.Commentf("cache: %s", name)
			fs := afero.NewMemMapFs()
			p := newPathsSpec(t, fs, configStr)
			fileCachConfig := p.Cfg.GetConfigSection("caches").(filecache.Configs)
			caches, err := filecache.NewCaches(fileCachConfig, fs)
			c.Assert(err, qt.IsNil)
			caches.SetResourceFs(fs)
			cache := caches[name]
			for i := range 10 {
				id := fmt.Sprintf("i%d", i)
				cache.GetOrCreateBytes(id, func() ([]byte, error) {
					return []byte("abc"), nil
				})
				if i == 4 {
					// This will expire the first 5
					time.Sleep(201 * time.Millisecond)
				}
			}

			count, err := caches.Prune()
			c.Assert(err, qt.IsNil)
			c.Assert(count, qt.Equals, 5, msg)

			for i := range 10 {
				id := fmt.Sprintf("i%d", i)
				v := cache.GetString(id)
				if i < 5 {
					c.Assert(v, qt.Equals, "")
				} else {
					c.Assert(v, qt.Equals, "abc")
				}
			}

			caches, err = filecache.NewCaches(fileCachConfig, fs)
			c.Assert(err, qt.IsNil)
			caches.SetResourceFs(fs)
			cache = caches[name]
			// Touch one and then prune.
			cache.GetOrCreateBytes("i5", func() ([]byte, error) {
				return []byte("abc"), nil
			})

			count, err = caches.Prune()
			c.Assert(err, qt.IsNil)
			c.Assert(count, qt.Equals, 4)

			// Now only the i5 should be left.
			for i := range 10 {
				id := fmt.Sprintf("i%d", i)
				v := cache.GetString(id)
				if i != 5 {
					c.Assert(v, qt.Equals, "")
				} else {
					c.Assert(v, qt.Equals, "abc")
				}
			}

		}
	})
}
