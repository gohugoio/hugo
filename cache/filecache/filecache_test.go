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

package filecache

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
)

func TestFileCache(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	tempWorkingDir, err := ioutil.TempDir("", "hugo_filecache_test_work")
	assert.NoError(err)
	defer os.Remove(tempWorkingDir)

	tempCacheDir, err := ioutil.TempDir("", "hugo_filecache_test_cache")
	assert.NoError(err)
	defer os.Remove(tempCacheDir)

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

		cfg, err := config.FromConfigString(configStr, "toml")
		assert.NoError(err)

		fs := hugofs.NewFrom(osfs, cfg)
		p, err := helpers.NewPathSpec(fs, cfg)
		assert.NoError(err)

		caches, err := NewCaches(p)
		assert.NoError(err)

		cache := caches.Get("GetJSON")
		assert.NotNil(cache)
		assert.Equal("10h0m0s", cache.maxAge.String())

		bfs, ok := cache.Fs.(*afero.BasePathFs)
		assert.True(ok)
		filename, err := bfs.RealPath("key")
		assert.NoError(err)
		if test.cacheDir != "" {
			assert.Equal(filepath.Join(test.cacheDir, "c/"+filecacheRootDirname+"/getjson/key"), filename)
		} else {
			// Temp dir.
			assert.Regexp(regexp.MustCompile(".*hugo_cache.*"+filecacheRootDirname+".*key"), filename)
		}

		cache = caches.Get("Images")
		assert.NotNil(cache)
		assert.Equal(time.Duration(-1), cache.maxAge)
		bfs, ok = cache.Fs.(*afero.BasePathFs)
		assert.True(ok)
		filename, _ = bfs.RealPath("key")
		assert.Equal(filepath.FromSlash("_gen/images/key"), filename)

		rf := func(s string) func() (io.ReadCloser, error) {
			return func() (io.ReadCloser, error) {
				return struct {
					io.ReadSeeker
					io.Closer
				}{
					strings.NewReader(s),
					ioutil.NopCloser(nil),
				}, nil
			}
		}

		bf := func() ([]byte, error) {
			return []byte("bcd"), nil
		}

		for _, c := range []*Cache{caches.ImageCache(), caches.AssetsCache(), caches.GetJSONCache(), caches.GetCSVCache()} {
			for i := 0; i < 2; i++ {
				info, r, err := c.GetOrCreate("a", rf("abc"))
				assert.NoError(err)
				assert.NotNil(r)
				assert.Equal("a", info.Name)
				b, _ := ioutil.ReadAll(r)
				r.Close()
				assert.Equal("abc", string(b))

				info, b, err = c.GetOrCreateBytes("b", bf)
				assert.NoError(err)
				assert.NotNil(r)
				assert.Equal("b", info.Name)
				assert.Equal("bcd", string(b))

				_, b, err = c.GetOrCreateBytes("a", bf)
				assert.NoError(err)
				assert.Equal("abc", string(b))

				_, r, err = c.GetOrCreate("a", rf("bcd"))
				assert.NoError(err)
				b, _ = ioutil.ReadAll(r)
				r.Close()
				assert.Equal("abc", string(b))
			}
		}

		assert.NotNil(caches.Get("getJSON"))

		info, w, err := caches.ImageCache().WriteCloser("mykey")
		assert.NoError(err)
		assert.Equal("mykey", info.Name)
		io.WriteString(w, "Hugo is great!")
		w.Close()
		assert.Equal("Hugo is great!", caches.ImageCache().getString("mykey"))

		info, r, err := caches.ImageCache().Get("mykey")
		assert.NoError(err)
		assert.NotNil(r)
		assert.Equal("mykey", info.Name)
		b, _ := ioutil.ReadAll(r)
		r.Close()
		assert.Equal("Hugo is great!", string(b))

		info, b, err = caches.ImageCache().GetBytes("mykey")
		assert.NoError(err)
		assert.Equal("mykey", info.Name)
		assert.Equal("Hugo is great!", string(b))

	}

}

func TestFileCacheConcurrent(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

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

	cfg, err := config.FromConfigString(configStr, "toml")
	assert.NoError(err)
	fs := hugofs.NewMem(cfg)
	p, err := helpers.NewPathSpec(fs, cfg)
	assert.NoError(err)

	caches, err := NewCaches(p)
	assert.NoError(err)

	const cacheName = "getjson"

	filenameData := func(i int) (string, string) {
		data := fmt.Sprintf("data: %d", i)
		filename := fmt.Sprintf("file%d", i)
		return filename, data
	}

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				c := caches.Get(cacheName)
				assert.NotNil(c)
				filename, data := filenameData(i)
				_, r, err := c.GetOrCreate(filename, func() (io.ReadCloser, error) {
					return hugio.ToReadCloser(strings.NewReader(data)), nil
				})
				assert.NoError(err)
				b, _ := ioutil.ReadAll(r)
				r.Close()
				assert.Equal(data, string(b))
				// Trigger some expiration.
				time.Sleep(50 * time.Millisecond)
			}
		}(i)

	}
	wg.Wait()
}

func TestCleanID(t *testing.T) {
	assert := require.New(t)
	assert.Equal(filepath.FromSlash("a/b/c.txt"), cleanID(filepath.FromSlash("/a/b//c.txt")))
	assert.Equal(filepath.FromSlash("a/b/c.txt"), cleanID(filepath.FromSlash("a/b//c.txt")))
}
