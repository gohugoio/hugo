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
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/config"

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
[caches.getJSON]
maxAge = "10m"
dir = "/path/to/c1"
[caches.getCSV]
maxAge = "11h"
dir = "/path/to/c2"
[caches.images]
dir = "/path/to/c3"
[caches.getResource]
dir = "/path/to/c4"
`

	cfg, err := config.FromConfigString(configStr, "toml")
	c.Assert(err, qt.IsNil)
	fs := afero.NewMemMapFs()
	decoded, err := DecodeConfig(fs, cfg)
	c.Assert(err, qt.IsNil)

	c.Assert(len(decoded), qt.Equals, 6)

	c2 := decoded["getcsv"]
	c.Assert(c2.MaxAge.String(), qt.Equals, "11h0m0s")
	c.Assert(c2.Dir, qt.Equals, filepath.FromSlash("/path/to/c2/filecache/getcsv"))

	c3 := decoded["images"]
	c.Assert(c3.MaxAge, qt.Equals, time.Duration(-1))
	c.Assert(c3.Dir, qt.Equals, filepath.FromSlash("/path/to/c3/filecache/images"))

	c4 := decoded["getresource"]
	c.Assert(c4.MaxAge, qt.Equals, time.Duration(-1))
	c.Assert(c4.Dir, qt.Equals, filepath.FromSlash("/path/to/c4/filecache/getresource"))
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
[caches.getJSON]
maxAge = 1234
dir = "/path/to/c1"
[caches.getCSV]
maxAge = 3456
dir = "/path/to/c2"
[caches.images]
dir = "/path/to/c3"
[caches.getResource]
dir = "/path/to/c4"
`

	cfg, err := config.FromConfigString(configStr, "toml")
	c.Assert(err, qt.IsNil)
	fs := afero.NewMemMapFs()
	decoded, err := DecodeConfig(fs, cfg)
	c.Assert(err, qt.IsNil)

	c.Assert(len(decoded), qt.Equals, 6)

	for _, v := range decoded {
		c.Assert(v.MaxAge, qt.Equals, time.Duration(0))
	}
}

func TestDecodeConfigDefault(t *testing.T) {
	c := qt.New(t)
	cfg := newTestConfig()

	if runtime.GOOS == "windows" {
		cfg.Set("resourceDir", "c:\\cache\\resources")
		cfg.Set("cacheDir", "c:\\cache\\thecache")

	} else {
		cfg.Set("resourceDir", "/cache/resources")
		cfg.Set("cacheDir", "/cache/thecache")
	}

	fs := afero.NewMemMapFs()

	decoded, err := DecodeConfig(fs, cfg)

	c.Assert(err, qt.IsNil)

	c.Assert(len(decoded), qt.Equals, 6)

	imgConfig := decoded[cacheKeyImages]
	jsonConfig := decoded[cacheKeyGetJSON]

	if runtime.GOOS == "windows" {
		c.Assert(imgConfig.Dir, qt.Equals, filepath.FromSlash("_gen/images"))
	} else {
		c.Assert(imgConfig.Dir, qt.Equals, "_gen/images")
		c.Assert(jsonConfig.Dir, qt.Equals, "/cache/thecache/hugoproject/filecache/getjson")
	}

	c.Assert(imgConfig.isResourceDir, qt.Equals, true)
	c.Assert(jsonConfig.isResourceDir, qt.Equals, false)
}

func TestDecodeConfigInvalidDir(t *testing.T) {
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
[caches.getJSON]
maxAge = "10m"
dir = "/"

`
	if runtime.GOOS == "windows" {
		configStr = strings.Replace(configStr, "/", "c:\\\\", 1)
	}

	cfg, err := config.FromConfigString(configStr, "toml")
	c.Assert(err, qt.IsNil)
	fs := afero.NewMemMapFs()

	_, err = DecodeConfig(fs, cfg)
	c.Assert(err, qt.Not(qt.IsNil))
}

func newTestConfig() config.Provider {
	cfg := config.NewWithTestDefaults()
	cfg.Set("workingDir", filepath.FromSlash("/my/cool/hugoproject"))
	cfg.Set("contentDir", "content")
	cfg.Set("dataDir", "data")
	cfg.Set("resourceDir", "resources")
	cfg.Set("i18nDir", "i18n")
	cfg.Set("layoutDir", "layouts")
	cfg.Set("archetypeDir", "archetypes")
	cfg.Set("assetDir", "assets")

	return cfg
}
