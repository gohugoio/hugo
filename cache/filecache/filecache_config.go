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

// Package filecache provides a file based cache for Hugo.
package filecache

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
)

const (
	resourcesGenDir = ":resourceDir/_gen"
	cacheDirProject = ":cacheDir/:project"
)

var defaultCacheConfig = FileCacheConfig{
	MaxAge: -1, // Never expire
	Dir:    cacheDirProject,
}

const (
	CacheKeyGetJSON     = "getjson"
	CacheKeyGetCSV      = "getcsv"
	CacheKeyImages      = "images"
	CacheKeyAssets      = "assets"
	CacheKeyModules     = "modules"
	CacheKeyGetResource = "getresource"
	CacheKeyMisc        = "misc"
)

type Configs map[string]FileCacheConfig

// For internal use.
func (c Configs) CacheDirModules() string {
	return c[CacheKeyModules].DirCompiled
}

var defaultCacheConfigs = Configs{
	CacheKeyModules: {
		MaxAge: -1,
		Dir:    ":cacheDir/modules",
	},
	CacheKeyGetJSON: defaultCacheConfig,
	CacheKeyGetCSV:  defaultCacheConfig,
	CacheKeyImages: {
		MaxAge: -1,
		Dir:    resourcesGenDir,
	},
	CacheKeyAssets: {
		MaxAge: -1,
		Dir:    resourcesGenDir,
	},
	CacheKeyGetResource: {
		MaxAge: -1, // Never expire
		Dir:    cacheDirProject,
	},
	CacheKeyMisc: {
		MaxAge: -1,
		Dir:    cacheDirProject,
	},
}

type FileCacheConfig struct {
	// Max age of cache entries in this cache. Any items older than this will
	// be removed and not returned from the cache.
	// A negative value means forever, 0 means cache is disabled.
	// Hugo is lenient with what types it accepts here, but we recommend using
	// a duration string, a sequence of  decimal numbers, each with optional fraction and a unit suffix,
	// such as "300ms", "1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	MaxAge time.Duration

	// The directory where files are stored.
	Dir         string
	DirCompiled string `json:"-"`

	// Will resources/_gen will get its own composite filesystem that
	// also checks any theme.
	IsResourceDir bool `json:"-"`
}

// GetJSONCache gets the file cache for getJSON.
func (f Caches) GetJSONCache() *Cache {
	return f[CacheKeyGetJSON]
}

// GetCSVCache gets the file cache for getCSV.
func (f Caches) GetCSVCache() *Cache {
	return f[CacheKeyGetCSV]
}

// ImageCache gets the file cache for processed images.
func (f Caches) ImageCache() *Cache {
	return f[CacheKeyImages]
}

// ModulesCache gets the file cache for Hugo Modules.
func (f Caches) ModulesCache() *Cache {
	return f[CacheKeyModules]
}

// AssetsCache gets the file cache for assets (processed resources, SCSS etc.).
func (f Caches) AssetsCache() *Cache {
	return f[CacheKeyAssets]
}

// MiscCache gets the file cache for miscellaneous stuff.
func (f Caches) MiscCache() *Cache {
	return f[CacheKeyMisc]
}

// GetResourceCache gets the file cache for remote resources.
func (f Caches) GetResourceCache() *Cache {
	return f[CacheKeyGetResource]
}

func DecodeConfig(fs afero.Fs, bcfg config.BaseConfig, m map[string]any) (Configs, error) {
	c := make(Configs)
	valid := make(map[string]bool)
	// Add defaults
	for k, v := range defaultCacheConfigs {
		c[k] = v
		valid[k] = true
	}

	_, isOsFs := fs.(*afero.OsFs)

	for k, v := range m {
		if _, ok := v.(maps.Params); !ok {
			continue
		}
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
			return nil, fmt.Errorf("failed to decode filecache config: %w", err)
		}

		if cc.Dir == "" {
			return c, errors.New("must provide cache Dir")
		}

		name := strings.ToLower(k)
		if !valid[name] {
			return nil, fmt.Errorf("%q is not a valid cache name", name)
		}

		c[name] = cc
	}

	for k, v := range c {
		dir := filepath.ToSlash(filepath.Clean(v.Dir))
		hadSlash := strings.HasPrefix(dir, "/")
		parts := strings.Split(dir, "/")

		for i, part := range parts {
			if strings.HasPrefix(part, ":") {
				resolved, isResource, err := resolveDirPlaceholder(fs, bcfg, part)
				if err != nil {
					return c, err
				}
				if isResource {
					v.IsResourceDir = true
				}
				parts[i] = resolved
			}
		}

		dir = path.Join(parts...)
		if hadSlash {
			dir = "/" + dir
		}
		v.DirCompiled = filepath.Clean(filepath.FromSlash(dir))

		if !v.IsResourceDir {
			if isOsFs && !filepath.IsAbs(v.DirCompiled) {
				return c, fmt.Errorf("%q must resolve to an absolute directory", v.DirCompiled)
			}

			// Avoid cache in root, e.g. / (Unix) or c:\ (Windows)
			if len(strings.TrimPrefix(v.DirCompiled, filepath.VolumeName(v.DirCompiled))) == 1 {
				return c, fmt.Errorf("%q is a root folder and not allowed as cache dir", v.DirCompiled)
			}
		}

		if !strings.HasPrefix(v.DirCompiled, "_gen") {
			// We do cache eviction (file removes) and since the user can set
			// his/hers own cache directory, we really want to make sure
			// we do not delete any files that do not belong to this cache.
			// We do add the cache name as the root, but this is an extra safe
			// guard. We skip the files inside /resources/_gen/ because
			// that would be breaking.
			v.DirCompiled = filepath.Join(v.DirCompiled, FilecacheRootDirname, k)
		} else {
			v.DirCompiled = filepath.Join(v.DirCompiled, k)
		}

		c[k] = v
	}

	return c, nil
}

// Resolves :resourceDir => /myproject/resources etc., :cacheDir => ...
func resolveDirPlaceholder(fs afero.Fs, bcfg config.BaseConfig, placeholder string) (cacheDir string, isResource bool, err error) {
	switch strings.ToLower(placeholder) {
	case ":resourcedir":
		return "", true, nil
	case ":cachedir":
		return bcfg.CacheDir, false, nil
	case ":project":
		return filepath.Base(bcfg.WorkingDir), false, nil
	}

	return "", false, fmt.Errorf("%q is not a valid placeholder (valid values are :cacheDir or :resourceDir)", placeholder)
}
