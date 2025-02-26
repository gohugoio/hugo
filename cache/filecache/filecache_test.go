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

package filecache_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestFileCache(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	tempWorkingDir := t.TempDir()
	tempCacheDir := t.TempDir()

	osfs := afero.NewOsFs()

	for _, test := range []struct {
		cacheDir   string
		workingDir string
	}{
		// Run with same dirs twice to make sure that works.
		{tempCacheDir, tempWorkingDir},
		{tempCacheDir, tempWorkingDir},
	} {

		configStr := `
workingDir = "WORKING_DIR"
resourceDir = "resources"
cacheDir = "CACHEDIR"
contentDir = "content"
dataDir = "data"
i18nDir = "i18n"
layoutDir = "layouts"
assetDir = "assets"
archeTypedir = "archetypes"

[caches]
[caches.getJSON]
maxAge = "10h"
dir = ":cacheDir/c"

`

		winPathSep := "\\\\"

		replacer := strings.NewReplacer("CACHEDIR", test.cacheDir, "WORKING_DIR", test.workingDir)

		configStr = replacer.Replace(configStr)
		configStr = strings.Replace(configStr, "\\", winPathSep, -1)

		p := newPathsSpec(t, osfs, configStr)

		caches, err := filecache.NewCaches(p)
		c.Assert(err, qt.IsNil)

		cache := caches.Get("GetJSON")
		c.Assert(cache, qt.Not(qt.IsNil))

		cache = caches.Get("Images")
		c.Assert(cache, qt.Not(qt.IsNil))

		rf := func(s string) func() (io.ReadCloser, error) {
			return func() (io.ReadCloser, error) {
				return struct {
					io.ReadSeeker
					io.Closer
				}{
					strings.NewReader(s),
					io.NopCloser(nil),
				}, nil
			}
		}

		bf := func() ([]byte, error) {
			return []byte("bcd"), nil
		}

		for _, ca := range []*filecache.Cache{caches.ImageCache(), caches.AssetsCache(), caches.GetJSONCache(), caches.GetCSVCache()} {
			for range 2 {
				info, r, err := ca.GetOrCreate("a", rf("abc"))
				c.Assert(err, qt.IsNil)
				c.Assert(r, qt.Not(qt.IsNil))
				c.Assert(info.Name, qt.Equals, "a")
				b, _ := io.ReadAll(r)
				r.Close()
				c.Assert(string(b), qt.Equals, "abc")

				info, b, err = ca.GetOrCreateBytes("b", bf)
				c.Assert(err, qt.IsNil)
				c.Assert(r, qt.Not(qt.IsNil))
				c.Assert(info.Name, qt.Equals, "b")
				c.Assert(string(b), qt.Equals, "bcd")

				_, b, err = ca.GetOrCreateBytes("a", bf)
				c.Assert(err, qt.IsNil)
				c.Assert(string(b), qt.Equals, "abc")

				_, r, err = ca.GetOrCreate("a", rf("bcd"))
				c.Assert(err, qt.IsNil)
				b, _ = io.ReadAll(r)
				r.Close()
				c.Assert(string(b), qt.Equals, "abc")
			}
		}

		c.Assert(caches.Get("getJSON"), qt.Not(qt.IsNil))

		info, w, err := caches.ImageCache().WriteCloser("mykey")
		c.Assert(err, qt.IsNil)
		c.Assert(info.Name, qt.Equals, "mykey")
		io.WriteString(w, "Hugo is great!")
		w.Close()
		c.Assert(caches.ImageCache().GetString("mykey"), qt.Equals, "Hugo is great!")

		info, r, err := caches.ImageCache().Get("mykey")
		c.Assert(err, qt.IsNil)
		c.Assert(r, qt.Not(qt.IsNil))
		c.Assert(info.Name, qt.Equals, "mykey")
		b, _ := io.ReadAll(r)
		r.Close()
		c.Assert(string(b), qt.Equals, "Hugo is great!")

		info, b, err = caches.ImageCache().GetBytes("mykey")
		c.Assert(err, qt.IsNil)
		c.Assert(info.Name, qt.Equals, "mykey")
		c.Assert(string(b), qt.Equals, "Hugo is great!")

	}
}

func TestFileCacheConcurrent(t *testing.T) {
	t.Parallel()

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
[caches.getjson]
maxAge = "1s"
dir = "/cache/c"

`

	p := newPathsSpec(t, afero.NewMemMapFs(), configStr)

	caches, err := filecache.NewCaches(p)
	c.Assert(err, qt.IsNil)

	const cacheName = "getjson"

	filenameData := func(i int) (string, string) {
		data := fmt.Sprintf("data: %d", i)
		filename := fmt.Sprintf("file%d", i)
		return filename, data
	}

	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for range 20 {
				ca := caches.Get(cacheName)
				c.Assert(ca, qt.Not(qt.IsNil))
				filename, data := filenameData(i)
				_, r, err := ca.GetOrCreate(filename, func() (io.ReadCloser, error) {
					return hugio.ToReadCloser(strings.NewReader(data)), nil
				})
				c.Assert(err, qt.IsNil)
				b, _ := io.ReadAll(r)
				r.Close()
				c.Assert(string(b), qt.Equals, data)
				// Trigger some expiration.
				time.Sleep(50 * time.Millisecond)
			}
		}(i)

	}
	wg.Wait()
}

func TestFileCacheReadOrCreateErrorInRead(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var result string

	rf := func(failLevel int) func(info filecache.ItemInfo, r io.ReadSeeker) error {
		return func(info filecache.ItemInfo, r io.ReadSeeker) error {
			if failLevel > 0 {
				if failLevel > 1 {
					return filecache.ErrFatal
				}
				return errors.New("fail")
			}

			b, _ := io.ReadAll(r)
			result = string(b)

			return nil
		}
	}

	bf := func(s string) func(info filecache.ItemInfo, w io.WriteCloser) error {
		return func(info filecache.ItemInfo, w io.WriteCloser) error {
			defer w.Close()
			result = s
			_, err := w.Write([]byte(s))
			return err
		}
	}

	cache := filecache.NewCache(afero.NewMemMapFs(), 100*time.Hour, "")

	const id = "a32"

	_, err := cache.ReadOrCreate(id, rf(0), bf("v1"))
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Equals, "v1")
	_, err = cache.ReadOrCreate(id, rf(0), bf("v2"))
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Equals, "v1")
	_, err = cache.ReadOrCreate(id, rf(1), bf("v3"))
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.Equals, "v3")
	_, err = cache.ReadOrCreate(id, rf(2), bf("v3"))
	c.Assert(err, qt.Equals, filecache.ErrFatal)
}

func newPathsSpec(t *testing.T, fs afero.Fs, configStr string) *helpers.PathSpec {
	c := qt.New(t)
	cfg, err := config.FromConfigString(configStr, "toml")
	c.Assert(err, qt.IsNil)
	acfg := testconfig.GetTestConfig(fs, cfg)
	p, err := helpers.NewPathSpec(hugofs.NewFrom(fs, acfg.BaseConfig()), acfg, nil)
	c.Assert(err, qt.IsNil)
	return p
}
