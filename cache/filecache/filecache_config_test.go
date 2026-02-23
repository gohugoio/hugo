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
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"

	qt "github.com/frankban/quicktest"
)

func TestDecodeConfig(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	configStr := `
resourceDir = "myresources"
contentDir = "content"
dataDir = "data"
i18nDir = "i18n"
layoutDir = "layouts"
assetDir = "assets"
archetypeDir = "archetypes"

[caches]
[caches.misc]
maxAge = "10m"
dir = "/path/to/c1"
[caches.images]
dir = "/path/to/c3"
[caches.getResource]
dir = "/path/to/c4"
`

	cfg, err := config.FromConfigString(configStr, "toml")
	c.Assert(err, qt.IsNil)
	fs := afero.NewMemMapFs()
	decoded := testconfig.GetTestConfigs(fs, cfg).Base.Caches
	c.Assert(len(decoded), qt.Equals, 7)

	c2 := decoded["misc"]
	c.Assert(c2.MaxAge.String(), qt.Equals, "10m0s")
	c.Assert(c2.DirCompiled, qt.Equals, filepath.FromSlash("/path/to/c1/filecache/misc"))

	c3 := decoded["images"]
	c.Assert(c3.MaxAge, qt.Equals, time.Duration(-1))
	c.Assert(c3.DirCompiled, qt.Equals, filepath.FromSlash("/path/to/c3/filecache/images"))

	c4 := decoded["getresource"]
	c.Assert(c4.MaxAge, qt.Equals, time.Duration(-1))
	c.Assert(c4.DirCompiled, qt.Equals, filepath.FromSlash("/path/to/c4/filecache/getresource"))
}

func TestDecodeConfigIgnoreCache(t *testing.T) {
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

ignoreCache = true
[caches]
[caches.misc]
maxAge = 1234
dir = "/path/to/c1"
[caches.images]
dir = "/path/to/c3"
[caches.getResource]
dir = "/path/to/c4"
`

	cfg, err := config.FromConfigString(configStr, "toml")
	c.Assert(err, qt.IsNil)
	fs := afero.NewMemMapFs()
	decoded := testconfig.GetTestConfigs(fs, cfg).Base.Caches
	c.Assert(len(decoded), qt.Equals, 7)

	for _, v := range decoded {
		c.Assert(v.MaxAge, qt.Equals, time.Duration(0))
	}
}

func TestDecodeConfigDefault(t *testing.T) {
	c := qt.New(t)
	cfg := config.New()

	if runtime.GOOS == "windows" {
		cfg.Set("resourceDir", "c:\\cache\\resources")
		cfg.Set("cacheDir", "c:\\cache\\thecache")

	} else {
		cfg.Set("resourceDir", "/cache/resources")
		cfg.Set("cacheDir", "/cache/thecache")
	}
	cfg.Set("workingDir", filepath.FromSlash("/my/cool/hugoproject"))

	fs := afero.NewMemMapFs()
	decoded := testconfig.GetTestConfigs(fs, cfg).Base.Caches
	c.Assert(len(decoded), qt.Equals, 7)

	imgConfig := decoded[filecache.CacheKeyImages]
	miscConfig := decoded[filecache.CacheKeyMisc]

	if runtime.GOOS == "windows" {
		c.Assert(imgConfig.DirCompiled, qt.Equals, filepath.FromSlash("_gen/images"))
	} else {
		c.Assert(imgConfig.DirCompiled, qt.Equals, "_gen/images")
		c.Assert(miscConfig.DirCompiled, qt.Equals, "/cache/thecache/hugoproject/filecache/misc")
	}

	c.Assert(imgConfig.IsResourceDir, qt.Equals, true)
	c.Assert(miscConfig.IsResourceDir, qt.Equals, false)
}

func TestFileCacheConfigMarshalJSON(t *testing.T) {
	c := qt.New(t)

	cfg := config.New()
	cfg.Set("cacheDir", "/cache")
	cfg.Set("workingDir", "/my/project")

	fs := afero.NewMemMapFs()
	decoded := testconfig.GetTestConfigs(fs, cfg).Base.Caches

	moduleQueriesConfig := decoded[filecache.CacheKeyModuleQueries]
	c.Assert(moduleQueriesConfig.MaxAge, qt.Equals, 24*time.Hour)

	// Also verify the new moduleGitInfo cache.
	moduleGitInfoConfig := decoded[filecache.CacheKeyModuleGitInfo]
	c.Assert(moduleGitInfoConfig.MaxAge, qt.Equals, 24*time.Hour)

	b, err := json.Marshal(moduleQueriesConfig)
	c.Assert(err, qt.IsNil)

	c.Assert(string(b), qt.Contains, `"maxAge":"24h"`)
	c.Assert(string(b), qt.Not(qt.Contains), "86400000000000")
	c.Assert(string(b), qt.Not(qt.Contains), "8.64e")

	moduleQueriesConfig.MaxAge = -1
	b, err = json.Marshal(moduleQueriesConfig)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Contains, `"maxAge":-1`)
}
