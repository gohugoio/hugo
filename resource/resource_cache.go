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

package resource

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/afero"

	"github.com/BurntSushi/locker"
)

const (
	CACHE_CLEAR_ALL = "clear_all"
	CACHE_OTHER     = "other"
)

type ResourceCache struct {
	rs *Spec

	cache map[string]Resource
	sync.RWMutex

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
		rs:      rs,
		cache:   make(map[string]Resource),
		nlocker: locker.NewLocker(),
	}
}

func (c *ResourceCache) clear() {
	c.Lock()
	defer c.Unlock()

	c.cache = make(map[string]Resource)
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

func (c *ResourceCache) get(key string) (Resource, bool) {
	c.RLock()
	defer c.RUnlock()
	r, found := c.cache[key]
	return r, found
}

func (c *ResourceCache) GetOrCreate(partition, key string, f func() (Resource, error)) (Resource, error) {
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
	filenameBase := filepath.Join(c.rs.GenAssetsPath, key)
	filenameMeta := filenameBase + ".json"
	filenameContent := filenameBase + ".content"

	return filenameMeta, filenameContent
}

func (c *ResourceCache) getFromFile(key string) (afero.File, transformedResourceMetadata, bool) {
	c.RLock()
	defer c.RUnlock()

	var meta transformedResourceMetadata
	filenameMeta, filenameContent := c.getFilenames(key)
	fMeta, err := c.rs.Resources.Fs.Open(filenameMeta)
	if err != nil {
		return nil, meta, false
	}
	defer fMeta.Close()

	jsonContent, err := ioutil.ReadAll(fMeta)
	if err != nil {
		return nil, meta, false
	}

	if err := json.Unmarshal(jsonContent, &meta); err != nil {
		return nil, meta, false
	}

	fContent, err := c.rs.Resources.Fs.Open(filenameContent)
	if err != nil {
		return nil, meta, false
	}

	return fContent, meta, true
}

// writeMeta writes the metadata to file and returns a writer for the content part.
func (c *ResourceCache) writeMeta(key string, meta transformedResourceMetadata) (afero.File, error) {
	filenameMeta, filenameContent := c.getFilenames(key)
	raw, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	fm, err := c.openResourceFileForWriting(filenameMeta)
	if err != nil {
		return nil, err
	}

	if _, err := fm.Write(raw); err != nil {
		return nil, err
	}

	return c.openResourceFileForWriting(filenameContent)

}

func (c *ResourceCache) openResourceFileForWriting(filename string) (afero.File, error) {
	return openFileForWriting(c.rs.Resources.Fs, filename)
}

// openFileForWriting opens or creates the given file. If the target directory
// does not exist, it gets created.
func openFileForWriting(fs afero.Fs, filename string) (afero.File, error) {
	filename = filepath.Clean(filename)
	// Create will truncate if file already exists.
	f, err := fs.Create(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = fs.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			return nil, err
		}
		f, err = fs.Create(filename)
	}

	return f, err
}

func (c *ResourceCache) set(key string, r Resource) {
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
