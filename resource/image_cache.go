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
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/helpers"
)

type imageCache struct {
	cacheDir string
	pathSpec *helpers.PathSpec
	mu       sync.RWMutex

	store map[string]*Image
}

func (c *imageCache) isInCache(key string) bool {
	c.mu.RLock()
	_, found := c.store[c.normalizeKey(key)]
	c.mu.RUnlock()
	return found
}

func (c *imageCache) deleteByPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prefix = c.normalizeKey(prefix)
	for k := range c.store {
		if strings.HasPrefix(k, prefix) {
			delete(c.store, k)
		}
	}
}

func (c *imageCache) normalizeKey(key string) string {
	// It is a path with Unix style slashes and it always starts with a leading slash.
	key = filepath.ToSlash(key)
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	return key
}

func (c *imageCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]*Image)
}

func (c *imageCache) getOrCreate(
	parent *Image, conf imageConfig, create func(resourceCacheFilename string) (*Image, error)) (*Image, error) {

	relTarget := parent.relTargetPathFromConfig(conf)
	key := parent.relTargetPathForRel(relTarget.path(), false, false)

	// First check the in-memory store, then the disk.
	c.mu.RLock()
	img, found := c.store[key]
	c.mu.RUnlock()

	if found {
		return img, nil
	}

	// Now look in the file cache.
	// Multiple Go routines can invoke same operation on the same image, so
	// we need to make sure this is serialized per source image.
	parent.createMu.Lock()
	defer parent.createMu.Unlock()

	cacheFilename := filepath.Join(c.cacheDir, key)

	// The definition of this counter is not that we have processed that amount
	// (e.g. resized etc.), it can be fetched from file cache,
	//  but the count of processed image variations for this site.
	c.pathSpec.ProcessingStats.Incr(&c.pathSpec.ProcessingStats.ProcessedImages)

	exists, err := helpers.Exists(cacheFilename, c.pathSpec.BaseFs.Resources.Fs)
	if err != nil {
		return nil, err
	}

	if exists {
		img = parent.clone()
	} else {
		img, err = create(cacheFilename)
		if err != nil {
			return nil, err
		}
	}
	img.relTargetDirFile.file = relTarget.file
	img.sourceFilename = cacheFilename
	// We have to look in the resources file system for this.
	img.overriddenSourceFs = img.spec.BaseFs.Resources.Fs

	c.mu.Lock()
	if img2, found := c.store[key]; found {
		c.mu.Unlock()
		return img2, nil
	}

	c.store[key] = img

	c.mu.Unlock()

	if !exists {
		// File already written to destination
		return img, nil
	}

	return img, img.copyToDestination(cacheFilename)

}

func newImageCache(ps *helpers.PathSpec, cacheDir string) *imageCache {
	return &imageCache{pathSpec: ps, store: make(map[string]*Image), cacheDir: cacheDir}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
