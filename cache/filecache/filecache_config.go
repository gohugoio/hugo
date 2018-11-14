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
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib/paths"

	"github.com/bep/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	cachesConfigKey = "caches"

	resourcesGenDir = ":resourceDir/_gen"
)

var defaultCacheConfig = cacheConfig{
	MaxAge: -1, // Never expire
	Dir:    ":cacheDir",
}

const (
	cacheKeyGetJSON = "getjson"
	cacheKeyGetCSV  = "getcsv"
	cacheKeyImages  = "images"
	cacheKeyAssets  = "assets"
)

var defaultCacheConfigs = map[string]cacheConfig{
	cacheKeyGetJSON: defaultCacheConfig,
	cacheKeyGetCSV:  defaultCacheConfig,
	cacheKeyImages: cacheConfig{
		MaxAge: -1,
		Dir:    resourcesGenDir,
	},
	cacheKeyAssets: cacheConfig{
		MaxAge: -1,
		Dir:    resourcesGenDir,
	},
}

type cachesConfig map[string]cacheConfig

type cacheConfig struct {
	// Max age of cache entries in this cache. Any items older than this will
	// be removed and not returned from the cache.
	// a negative value means forever, 0 means cache is disabled.
	MaxAge time.Duration

	// The directory where files are stored.
	Dir string
}

// GetJSONCache gets the file cache for getJSON.
func (f Caches) GetJSONCache() *Cache {
	return f[cacheKeyGetJSON]
}

// GetCSVCache gets the file cache for getCSV.
func (f Caches) GetCSVCache() *Cache {
	return f[cacheKeyGetCSV]
}

// ImageCache gets the file cache for processed images.
func (f Caches) ImageCache() *Cache {
	return f[cacheKeyImages]
}

// AssetsCache gets the file cache for assets (processed resources, SCSS etc.).
func (f Caches) AssetsCache() *Cache {
	return f[cacheKeyAssets]
}

func decodeConfig(p *paths.Paths) (cachesConfig, error) {
	c := make(cachesConfig)
	valid := make(map[string]bool)
	// Add defaults
	for k, v := range defaultCacheConfigs {
		c[k] = v
		valid[k] = true
	}

	cfg := p.Cfg

	m := cfg.GetStringMap(cachesConfigKey)

	_, isOsFs := p.Fs.Source.(*afero.OsFs)

	for k, v := range m {
		cc := defaultCacheConfig

		dc := &mapstructure.DecoderConfig{
			Result:           &cc,
			DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
			WeaklyTypedInput: true,
		}

		decoder, err := mapstructure.NewDecoder(dc)
		if err != nil {
			return c, err
		}

		if err := decoder.Decode(v); err != nil {
			return nil, err
		}

		if cc.Dir == "" {
			return c, errors.New("must provide cache Dir")
		}

		name := strings.ToLower(k)
		if !valid[name] {
			return nil, errors.Errorf("%q is not a valid cache name", name)
		}

		c[name] = cc
	}

	// This is a very old flag in Hugo, but we need to respect it.
	disabled := cfg.GetBool("ignoreCache")

	for k, v := range c {
		v.Dir = filepath.Clean(v.Dir)
		dir := filepath.ToSlash(v.Dir)
		parts := strings.Split(dir, "/")
		first := parts[0]

		if strings.HasPrefix(first, ":") {
			resolved, err := resolveDirPlaceholder(p, first)
			if err != nil {
				return c, err
			}
			resolved = filepath.ToSlash(resolved)

			v.Dir = filepath.FromSlash(path.Join((append([]string{resolved}, parts[1:]...))...))

		} else if isOsFs && !path.IsAbs(dir) {
			return c, errors.Errorf("%q must either start with a placeholder (e.g. :cacheDir, :resourceDir) or be absolute", v.Dir)
		}

		if len(v.Dir) < 5 {
			return c, errors.Errorf("%q is not a valid cache dir", v.Dir)
		}

		if disabled {
			v.MaxAge = 0
		}

		c[k] = v
	}

	return c, nil
}

// Resolves :resourceDir => /myproject/resources etc., :cacheDir => ...
func resolveDirPlaceholder(p *paths.Paths, placeholder string) (string, error) {
	switch strings.ToLower(placeholder) {
	case ":resourcedir":
		return p.AbsResourcesDir, nil
	case ":cachedir":
		return helpers.GetCacheDir(p.Fs.Source, p.Cfg)
	}

	return "", errors.Errorf("%q is not a valid placeholder (valid values are :cacheDir or :resourceDir)", placeholder)
}
