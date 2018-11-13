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
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/hugolib/paths"

	"github.com/pkg/errors"

	"github.com/BurntSushi/locker"
	"github.com/bep/mapstructure"
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
	// Maxe age of ache entries in this cache. Any items older than this will
	// be removed and not returned from the cache.
	// -1 means forever, 0 means cache is disabled.
	MaxAge int

	// The directory where files are stored.
	Dir string
}

// Cache caches a set of files in a directory. This is usually a file on
// disk, but since this is backed by an Afero file system, it can be anything.
type Cache struct {
	Fs afero.Fs

	// Max age in seconds.
	maxAge int

	nlocker *locker.Locker
}

// ItemInfo contains info about a cached file.
type ItemInfo struct {
	// This is the file's name relative to the cache's filesystem.
	Name string
}

// NewCache creates a new file cache with the given filesystem and max age.
func NewCache(fs afero.Fs, maxAge int) *Cache {
	return &Cache{
		Fs:      fs,
		nlocker: locker.NewLocker(),
		maxAge:  maxAge,
	}
}

// lockedFile is a file with a lock that is released on Close.
type lockedFile struct {
	afero.File
	unlock func()
}

func (l *lockedFile) Close() error {
	defer l.unlock()
	return l.File.Close()
}

// WriteCloser returns a transactional writer into the cache.
// It's important that it's closed when done.
func (c *Cache) WriteCloser(id string) (ItemInfo, io.WriteCloser, error) {
	id = cleanID(id)
	c.nlocker.Lock(id)

	info := ItemInfo{Name: id}

	f, err := helpers.OpenFileForWriting(c.Fs, id)
	if err != nil {
		c.nlocker.Unlock(id)
		return info, nil, err
	}

	return info, &lockedFile{
		File:   f,
		unlock: func() { c.nlocker.Unlock(id) },
	}, nil
}

// ReadOrCreate tries to lookup the file in cache.
// If found, it is passed to read and then closed.
// If not found a new file is created and passed to create, which should close
// it when done.
func (c *Cache) ReadOrCreate(id string,
	read func(info ItemInfo, r io.Reader) error,
	create func(info ItemInfo, w io.WriteCloser) error) (info ItemInfo, err error) {
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info = ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		err = read(info, r)
		defer r.Close()
		return
	}

	f, err := helpers.OpenFileForWriting(c.Fs, id)
	if err != nil {
		return
	}

	err = create(info, f)

	return

}

// GetOrCreate tries to get the file with the given id from cache. If not found or expired, create will
// be invoked and the result cached.
// This method is protected by a named lock using the given id as identifier.
func (c *Cache) GetOrCreate(id string, create func() (io.ReadCloser, error)) (ItemInfo, io.ReadCloser, error) {
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info := ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		return info, r, nil
	}

	r, err := create()
	if err != nil {
		return info, nil, err
	}

	if c.maxAge == 0 {
		// No caching.
		return info, hugio.ToReadCloser(r), nil
	}

	var buff bytes.Buffer
	return info,
		hugio.ToReadCloser(&buff),
		afero.WriteReader(c.Fs, id, io.TeeReader(r, &buff))
}

// GetOrCreateBytes is the same as GetOrCreate, but produces a byte slice.
func (c *Cache) GetOrCreateBytes(id string, create func() ([]byte, error)) (ItemInfo, []byte, error) {
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info := ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		defer r.Close()
		b, err := ioutil.ReadAll(r)
		return info, b, err
	}

	b, err := create()
	if err != nil {
		return info, nil, err
	}

	if c.maxAge == 0 {
		return info, b, nil
	}

	if err := afero.WriteReader(c.Fs, id, bytes.NewReader(b)); err != nil {
		return info, nil, err
	}
	return info, b, nil

}

// GetBytes gets the file content with the given id from the cahce, nil if none found.
func (c *Cache) GetBytes(id string) (ItemInfo, []byte, error) {
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info := ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		defer r.Close()
		b, err := ioutil.ReadAll(r)
		return info, b, err
	}

	return info, nil, nil
}

// Get gets the file with the given id from the cahce, nil if none found.
func (c *Cache) Get(id string) (ItemInfo, io.ReadCloser, error) {
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info := ItemInfo{Name: id}

	r := c.getOrRemove(id)

	return info, r, nil
}

// getOrRemove gets the file with the given id. If it's expired, it will
// be removed.
func (c *Cache) getOrRemove(id string) hugio.ReadSeekCloser {
	if c.maxAge == 0 {
		// No caching.
		return nil
	}

	if c.maxAge > 0 {
		fi, err := c.Fs.Stat(id)
		if err != nil {
			return nil
		}

		expiry := time.Now().Add(-time.Duration(c.maxAge) * time.Second)
		expired := fi.ModTime().Before(expiry)
		if expired {
			c.Fs.Remove(id)
			return nil
		}
	}

	f, err := c.Fs.Open(id)

	if err != nil {
		return nil
	}

	return f
}

// For testing
func (c *Cache) getString(id string) string {
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	if r := c.getOrRemove(id); r != nil {
		defer r.Close()
		b, _ := ioutil.ReadAll(r)
		return string(b)
	}

	return ""

}

// Caches is a named set of caches.
type Caches map[string]*Cache

// Get gets a named cache, nil if none found.
func (f Caches) Get(name string) *Cache {
	return f[strings.ToLower(name)]
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

// NewCachesFromPaths creates a new set of file caches from the given
// configuration.
func NewCachesFromPaths(p *paths.Paths) (Caches, error) {
	dcfg, err := decodeConfig(p)
	if err != nil {
		return nil, err
	}

	fs := p.Fs.Source

	m := make(Caches)
	for k, v := range dcfg {
		baseDir := filepath.Join(v.Dir, k)
		if err = fs.MkdirAll(baseDir, 0777); err != nil {
			return nil, err
		}
		bfs := afero.NewBasePathFs(fs, baseDir)
		m[k] = NewCache(bfs, v.MaxAge)
	}

	return m, nil
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

		if err := mapstructure.WeakDecode(v, &cc); err != nil {
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

func cleanID(name string) string {
	return filepath.Clean(name)
}
