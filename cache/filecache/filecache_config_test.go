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
	"testing"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/paths"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestDecodeConfig(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	configStr := `
[caches]
[caches.getJSON]
maxAge = 1234
dir = "/path/to/c1"
[caches.getCSV]
maxAge = 3456
dir = "/path/to/c2"
[caches.images]
dir = "/path/to/c3"

`

	cfg, err := config.FromConfigString(configStr, "toml")
	assert.NoError(err)
	fs := hugofs.NewMem(cfg)
	p, err := paths.New(fs, cfg)
	assert.NoError(err)

	decoded, err := decodeConfig(p)
	assert.NoError(err)

	assert.Equal(4, len(decoded))

	c2 := decoded["getcsv"]
	assert.Equal(3456, c2.MaxAge)
	assert.Equal(filepath.FromSlash("/path/to/c2"), c2.Dir)

	c3 := decoded["images"]
	assert.Equal(-1, c3.MaxAge)
	assert.Equal(filepath.FromSlash("/path/to/c3"), c3.Dir)

}

func TestDecodeConfigIgnoreCache(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	configStr := `
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

`

	cfg, err := config.FromConfigString(configStr, "toml")
	assert.NoError(err)
	fs := hugofs.NewMem(cfg)
	p, err := paths.New(fs, cfg)
	assert.NoError(err)

	decoded, err := decodeConfig(p)
	assert.NoError(err)

	assert.Equal(4, len(decoded))

	for _, v := range decoded {
		assert.Equal(0, v.MaxAge)
	}

}

func TestDecodeConfigDefault(t *testing.T) {
	assert := require.New(t)
	cfg := viper.New()
	if runtime.GOOS == "windows" {
		cfg.Set("resourceDir", "c:\\cache\\resources")
		cfg.Set("cacheDir", "c:\\cache\\thecache")

	} else {
		cfg.Set("resourceDir", "/cache/resources")
		cfg.Set("cacheDir", "/cache/thecache")
	}

	fs := hugofs.NewMem(cfg)
	p, err := paths.New(fs, cfg)
	assert.NoError(err)

	decoded, err := decodeConfig(p)

	assert.NoError(err)

	assert.Equal(4, len(decoded))

	if runtime.GOOS == "windows" {
		assert.Equal("c:\\cache\\resources\\_gen", decoded[cacheKeyImages].Dir)
	} else {
		assert.Equal("/cache/resources/_gen", decoded[cacheKeyImages].Dir)
	}
}
