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

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/hugofs/glob"

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

	// Either resource.Resource or resource.Resources.
	cache map[string]interface{}

	fileCache *filecache.Cache

	// Provides named resource locks.
	nlocker *locker.Locker
}

// ResourceCacheKey converts the filename into the format used in the resource
// cache.
func ResourceCacheKey(filename string) string {
	filename = filepath.ToSlash(filename)
	return path.Join(resourceKeyPartition(filename), filename)
}

func resourceKeyPartition(filename string) string {
	ext := strings.TrimPrefix(path.Ext(filepath.ToSlash(filename)), ".")
	if ext == "" {
		ext = CACHE_OTHER
	}
	return ext
}

// Commonly used aliases and directory names used for some types.
var extAliasKeywords = map[string][]string{
	"sass": []string{"scss"},
	"scss": []string{"sass"},
}

// ResourceKeyPartitions resolves a ordered slice of partitions that is
// used to do resource cache invalidations.
//
// We use the first directory path element and the extension, so:
//     a/b.json => "a", "json"
//     b.json => "json"
//
// For some of the extensions we will also map to closely related types,
// e.g. "scss" will also return "sass".
//
func ResourceKeyPartitions(filename string) []string {
	var partitions []string
	filename = glob.NormalizePath(filename)
	dir, name := path.Split(filename)
	ext := strings.TrimPrefix(path.Ext(filepath.ToSlash(name)), ".")

	if dir != "" {
		partitions = append(partitions, strings.Split(dir, "/")[0])
	}

	if ext != "" {
		partitions = append(partitions, ext)
	}

	if aliases, found := extAliasKeywords[ext]; found {
		partitions = append(partitions, aliases...)
	}

	if len(partitions) == 0 {
		partitions = []string{CACHE_OTHER}
	}

	return helpers.UniqueStringsSorted(partitions)
}

// ResourceKeyContainsAny returns whether the key is a member of any of the
// given partitions.
//
// This is used for resource cache invalidation.
func ResourceKeyContainsAny(key string, partitions []string) bool {
	parts := strings.Split(key, "/")
	for _, p1 := range partitions {
		for _, p2 := range parts {
			if p1 == p2 {
				return true
			}
		}
	}
	return false
}

func newResourceCache(rs *Spec) *ResourceCache {
	return &ResourceCache{
		rs:        rs,
		fileCache: rs.FileCaches.AssetsCache(),
		cache:     make(map[string]interface{}),
		nlocker:   locker.NewLocker(),
	}
}

func (c *ResourceCache) clear() {
	c.Lock()
	defer c.Unlock()

	c.cache = make(map[string]interface{})
	c.nlocker = locker.NewLocker()
}

func (c *ResourceCache) Contains(key string) bool {
	key = c.cleanKey(filepath.ToSlash(key))
	_, found := c.get(key)
	return found
}

func (c *ResourceCache) cleanKey(key string) string {
	return strings.TrimPrefix(path.Clean(strings.ToLower(key)), "/")
}

func (c *ResourceCache) get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	r, found := c.cache[key]
	return r, found
}

func (c *ResourceCache) GetOrCreate(key string, f func() (resource.Resource, error)) (resource.Resource, error) {
	r, err := c.getOrCreate(key, func() (interface{}, error) { return f() })
	if r == nil || err != nil {
		return nil, err
	}
	return r.(resource.Resource), nil
}

func (c *ResourceCache) GetOrCreateResources(key string, f func() (resource.Resources, error)) (resource.Resources, error) {
	r, err := c.getOrCreate(key, func() (interface{}, error) { return f() })
	if r == nil || err != nil {
		return nil, err
	}
	return r.(resource.Resources), nil
}

func (c *ResourceCache) getOrCreate(key string, f func() (interface{}, error)) (interface{}, error) {
	key = c.cleanKey(key)
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

func (c *ResourceCache) set(key string, r interface{}) {
	c.Lock()
	defer c.Unlock()
	c.cache[key] = r
}

func (c *ResourceCache) DeletePartitions(partitions ...string) {
	partitionsSet := map[string]bool{
		// Always clear out the resources not matching any partition.
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
		for p := range partitionsSet {
			if strings.Contains(k, p) {
				// There will be some false positive, but that's fine.
				clear = true
				break
			}
		}

		if clear {
			delete(c.cache, k)
		}
	}

}
