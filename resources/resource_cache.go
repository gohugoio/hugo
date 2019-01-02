// Copyright 2019 The Hugo Authors. All rights reserved.
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

package resources

import (
	"encoding/json"
	"io"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/cache/filecache"

	"github.com/BurntSushi/locker"
)

const (
	CACHE_CLEAR_ALL = "clear_all"
	CACHE_OTHER     = "other"
)

type ResourceCache struct {
	rs *Spec

	sync.RWMutex
	cache map[string]resource.Resource

	fileCache *filecache.Cache

	// Provides named resource locks.
	nlocker *locker.Locker
}

// ResourceKeyPartition returns a partition name
// to  allow for more fine grained cache flushes.
// It will return the file extension without the leading ".". If no
// extension, it will return "other".
func ResourceKeyPartition(filename string) string {
	ext := strings.TrimPrefix(path.Ext(filepath.ToSlash(filename)), ".")
	if ext == "" {
		ext = CACHE_OTHER
	}
	return ext
}

func newResourceCache(rs *Spec) *ResourceCache {
	return &ResourceCache{
		rs:        rs,
		fileCache: rs.FileCaches.AssetsCache(),
		cache:     make(map[string]resource.Resource),
		nlocker:   locker.NewLocker(),
	}
}

func (c *ResourceCache) clear() {
	c.Lock()
	defer c.Unlock()

	c.cache = make(map[string]resource.Resource)
	c.nlocker = locker.NewLocker()
}

func (c *ResourceCache) Contains(key string) bool {
	key = c.cleanKey(filepath.ToSlash(key))
	_, found := c.get(key)
	return found
}

func (c *ResourceCache) cleanKey(key string) string {
	return strings.TrimPrefix(path.Clean(key), "/")
}

func (c *ResourceCache) get(key string) (resource.Resource, bool) {
	c.RLock()
	defer c.RUnlock()
	r, found := c.cache[key]
	return r, found
}

func (c *ResourceCache) GetOrCreate(partition, key string, f func() (resource.Resource, error)) (resource.Resource, error) {
	key = c.cleanKey(path.Join(partition, key))
	// First check in-memory cache.
	r, found := c.get(key)
	if found {
		return r, nil
	}
	// This is a potentially long running operation, so get a named lock.
	c.nlocker.Lock(key)

	// Double check in-memory cache.
	r, found = c.get(key)
	if found {
		c.nlocker.Unlock(key)
		return r, nil
	}

	defer c.nlocker.Unlock(key)

	r, err := f()
	if err != nil {
		return nil, err
	}

	c.set(key, r)

	return r, nil

}

func (c *ResourceCache) getFilenames(key string) (string, string) {
	filenameMeta := key + ".json"
	filenameContent := key + ".content"

	return filenameMeta, filenameContent
}

func (c *ResourceCache) getFromFile(key string) (filecache.ItemInfo, io.ReadCloser, transformedResourceMetadata, bool) {
	c.RLock()
	defer c.RUnlock()

	var meta transformedResourceMetadata
	filenameMeta, filenameContent := c.getFilenames(key)

	_, jsonContent, _ := c.fileCache.GetBytes(filenameMeta)
	if jsonContent == nil {
		return filecache.ItemInfo{}, nil, meta, false
	}

	if err := json.Unmarshal(jsonContent, &meta); err != nil {
		return filecache.ItemInfo{}, nil, meta, false
	}

	fi, rc, _ := c.fileCache.Get(filenameContent)

	return fi, rc, meta, rc != nil

}

// writeMeta writes the metadata to file and returns a writer for the content part.
func (c *ResourceCache) writeMeta(key string, meta transformedResourceMetadata) (filecache.ItemInfo, io.WriteCloser, error) {
	filenameMeta, filenameContent := c.getFilenames(key)
	raw, err := json.Marshal(meta)
	if err != nil {
		return filecache.ItemInfo{}, nil, err
	}

	_, fm, err := c.fileCache.WriteCloser(filenameMeta)
	if err != nil {
		return filecache.ItemInfo{}, nil, err
	}
	defer fm.Close()

	if _, err := fm.Write(raw); err != nil {
		return filecache.ItemInfo{}, nil, err
	}

	fi, fc, err := c.fileCache.WriteCloser(filenameContent)

	return fi, fc, err

}

func (c *ResourceCache) set(key string, r resource.Resource) {
	c.Lock()
	defer c.Unlock()
	c.cache[key] = r
}

func (c *ResourceCache) DeletePartitions(partitions ...string) {
	partitionsSet := map[string]bool{
		// Always clear out the resources not matching the partition.
		"other": true,
	}
	for _, p := range partitions {
		partitionsSet[p] = true
	}

	if partitionsSet[CACHE_CLEAR_ALL] {
		c.clear()
		return
	}

	c.Lock()
	defer c.Unlock()

	for k := range c.cache {
		clear := false
		partIdx := strings.Index(k, "/")
		if partIdx == -1 {
			clear = true
		} else {
			partition := k[:partIdx]
			if partitionsSet[partition] {
				clear = true
			}
		}

		if clear {
			delete(c.cache, k)
		}
	}

}
