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

package filecache

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/httpcache"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/helpers"

	"github.com/BurntSushi/locker"
	"github.com/spf13/afero"
)

// ErrFatal can be used to signal an unrecoverable error.
var ErrFatal = errors.New("fatal filecache error")

const (
	FilecacheRootDirname = "filecache"
)

// Cache caches a set of files in a directory. This is usually a file on
// disk, but since this is backed by an Afero file system, it can be anything.
type Cache struct {
	Fs afero.Fs

	// Max age for items in this cache. Negative duration means forever,
	// 0 is effectively turning this cache off.
	maxAge time.Duration

	// When set, we just remove this entire root directory on expiration.
	pruneAllRootDir string

	nlocker *lockTracker

	initOnce sync.Once
	initErr  error
}

type lockTracker struct {
	seenMu sync.RWMutex
	seen   map[string]struct{}

	*locker.Locker
}

// Lock tracks the ids in use. We use this information to do garbage collection
// after a Hugo build.
func (l *lockTracker) Lock(id string) {
	l.seenMu.RLock()
	if _, seen := l.seen[id]; !seen {
		l.seenMu.RUnlock()
		l.seenMu.Lock()
		l.seen[id] = struct{}{}
		l.seenMu.Unlock()
	} else {
		l.seenMu.RUnlock()
	}

	l.Locker.Lock(id)
}

// ItemInfo contains info about a cached file.
type ItemInfo struct {
	// This is the file's name relative to the cache's filesystem.
	Name string
}

// NewCache creates a new file cache with the given filesystem and max age.
func NewCache(fs afero.Fs, maxAge time.Duration, pruneAllRootDir string) *Cache {
	return &Cache{
		Fs:              fs,
		nlocker:         &lockTracker{Locker: locker.NewLocker(), seen: make(map[string]struct{})},
		maxAge:          maxAge,
		pruneAllRootDir: pruneAllRootDir,
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

func (c *Cache) init() error {
	c.initOnce.Do(func() {
		// Create the base dir if it does not exist.
		if err := c.Fs.MkdirAll("", 0o777); err != nil && !os.IsExist(err) {
			c.initErr = err
		}
	})
	return c.initErr
}

// WriteCloser returns a transactional writer into the cache.
// It's important that it's closed when done.
func (c *Cache) WriteCloser(id string) (ItemInfo, io.WriteCloser, error) {
	if err := c.init(); err != nil {
		return ItemInfo{}, nil, err
	}

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
	read func(info ItemInfo, r io.ReadSeeker) error,
	create func(info ItemInfo, w io.WriteCloser) error,
) (info ItemInfo, err error) {
	if err := c.init(); err != nil {
		return ItemInfo{}, err
	}

	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info = ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		err = read(info, r)
		defer r.Close()
		if err == nil || err == ErrFatal {
			// See https://github.com/gohugoio/hugo/issues/6401
			// To recover from file corruption we handle read errors
			// as the cache item was not found.
			// Any file permission issue will also fail in the next step.
			return
		}
	}

	f, err := helpers.OpenFileForWriting(c.Fs, id)
	if err != nil {
		return
	}

	err = create(info, f)

	return
}

// NamedLock locks the given id. The lock is released when the returned function is called.
func (c *Cache) NamedLock(id string) func() {
	id = cleanID(id)
	c.nlocker.Lock(id)
	return func() {
		c.nlocker.Unlock(id)
	}
}

// GetOrCreate tries to get the file with the given id from cache. If not found or expired, create will
// be invoked and the result cached.
// This method is protected by a named lock using the given id as identifier.
func (c *Cache) GetOrCreate(id string, create func() (io.ReadCloser, error)) (ItemInfo, io.ReadCloser, error) {
	if err := c.init(); err != nil {
		return ItemInfo{}, nil, err
	}
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info := ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		return info, r, nil
	}

	var (
		r   io.ReadCloser
		err error
	)

	r, err = create()
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
		c.writeReader(id, io.TeeReader(r, &buff))
}

func (c *Cache) writeReader(id string, r io.Reader) error {
	dir := filepath.Dir(id)
	if dir != "" {
		_ = c.Fs.MkdirAll(dir, 0o777)
	}
	f, err := c.Fs.Create(id)
	if err != nil {
		return err
	}
	defer f.Close()

	_, _ = io.Copy(f, r)

	return nil
}

// GetOrCreateBytes is the same as GetOrCreate, but produces a byte slice.
func (c *Cache) GetOrCreateBytes(id string, create func() ([]byte, error)) (ItemInfo, []byte, error) {
	if err := c.init(); err != nil {
		return ItemInfo{}, nil, err
	}
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info := ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		defer r.Close()
		b, err := io.ReadAll(r)
		return info, b, err
	}

	var (
		b   []byte
		err error
	)

	b, err = create()
	if err != nil {
		return info, nil, err
	}

	if c.maxAge == 0 {
		return info, b, nil
	}

	if err := c.writeReader(id, bytes.NewReader(b)); err != nil {
		return info, nil, err
	}

	return info, b, nil
}

// GetBytes gets the file content with the given id from the cache, nil if none found.
func (c *Cache) GetBytes(id string) (ItemInfo, []byte, error) {
	if err := c.init(); err != nil {
		return ItemInfo{}, nil, err
	}
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	info := ItemInfo{Name: id}

	if r := c.getOrRemove(id); r != nil {
		defer r.Close()
		b, err := io.ReadAll(r)
		return info, b, err
	}

	return info, nil, nil
}

// Get gets the file with the given id from the cache, nil if none found.
func (c *Cache) Get(id string) (ItemInfo, io.ReadCloser, error) {
	if err := c.init(); err != nil {
		return ItemInfo{}, nil, err
	}
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

	if removed, err := c.removeIfExpired(id); err != nil || removed {
		return nil
	}

	f, err := c.Fs.Open(id)
	if err != nil {
		return nil
	}

	return f
}

func (c *Cache) getBytesAndRemoveIfExpired(id string) ([]byte, bool) {
	if c.maxAge == 0 {
		// No caching.
		return nil, false
	}

	f, err := c.Fs.Open(id)
	if err != nil {
		return nil, false
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, false
	}

	removed, err := c.removeIfExpired(id)
	if err != nil {
		return nil, false
	}

	return b, removed
}

func (c *Cache) removeIfExpired(id string) (bool, error) {
	if c.maxAge <= 0 {
		return false, nil
	}

	fi, err := c.Fs.Stat(id)
	if err != nil {
		return false, err
	}

	if c.isExpired(fi.ModTime()) {
		c.Fs.Remove(id)
		return true, nil
	}

	return false, nil
}

func (c *Cache) isExpired(modTime time.Time) bool {
	if c.maxAge < 0 {
		return false
	}

	// Note the use of time.Since here.
	// We cannot use Hugo's global Clock for this.
	return c.maxAge == 0 || time.Since(modTime) > c.maxAge
}

// For testing
func (c *Cache) GetString(id string) string {
	id = cleanID(id)

	c.nlocker.Lock(id)
	defer c.nlocker.Unlock(id)

	f, err := c.Fs.Open(id)
	if err != nil {
		return ""
	}
	defer f.Close()

	b, _ := io.ReadAll(f)
	return string(b)
}

// Caches is a named set of caches.
type Caches map[string]*Cache

// Get gets a named cache, nil if none found.
func (f Caches) Get(name string) *Cache {
	return f[strings.ToLower(name)]
}

// NewCaches creates a new set of file caches from the given
// configuration.
func NewCaches(p *helpers.PathSpec) (Caches, error) {
	dcfg := p.Cfg.GetConfigSection("caches").(Configs)
	fs := p.Fs.Source

	m := make(Caches)
	for k, v := range dcfg {
		var cfs afero.Fs

		if v.IsResourceDir {
			cfs = p.BaseFs.ResourcesCache
		} else {
			cfs = fs
		}

		if cfs == nil {
			panic("nil fs")
		}

		baseDir := v.DirCompiled

		bfs := hugofs.NewBasePathFs(cfs, baseDir)

		var pruneAllRootDir string
		if k == CacheKeyModules {
			pruneAllRootDir = "pkg"
		}

		m[k] = NewCache(bfs, v.MaxAge, pruneAllRootDir)
	}

	return m, nil
}

func cleanID(name string) string {
	return strings.TrimPrefix(filepath.Clean(name), helpers.FilePathSeparator)
}

// AsHTTPCache returns an httpcache.Cache implementation for this file cache.
// Note that none of the methods are protected by named locks, so you need to make sure
// to do that in your own code.
func (c *Cache) AsHTTPCache() httpcache.Cache {
	return &httpCache{c: c}
}

type httpCache struct {
	c *Cache
}

func (h *httpCache) Get(id string) (resp []byte, ok bool) {
	id = cleanID(id)
	b, removed := h.c.getBytesAndRemoveIfExpired(id)

	return b, !removed
}

func (h *httpCache) Set(id string, resp []byte) {
	if h.c.maxAge == 0 {
		return
	}

	id = cleanID(id)

	if err := h.c.writeReader(id, bytes.NewReader(resp)); err != nil {
		panic(err)
	}
}

func (h *httpCache) Delete(key string) {
	h.c.Fs.Remove(key)
}
