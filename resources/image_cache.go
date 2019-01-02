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
	"fmt"
	"image"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/helpers"
)

type imageCache struct {
	pathSpec *helpers.PathSpec

	fileCache *filecache.Cache

	mu    sync.RWMutex
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
	parent *Image, conf imageConfig, createImage func() (*Image, image.Image, error)) (*Image, error) {

	relTarget := parent.relTargetPathFromConfig(conf)
	key := parent.relTargetPathForRel(relTarget.path(), false, false, false)

	// First check the in-memory store, then the disk.
	c.mu.RLock()
	img, found := c.store[key]
	c.mu.RUnlock()

	if found {
		return img, nil
	}

	// These funcs are protected by a named lock.
	// read clones the parent to its new name and copies
	// the content to the destinations.
	read := func(info filecache.ItemInfo, r io.Reader) error {
		img = parent.clone()
		img.relTargetDirFile.file = relTarget.file
		img.sourceFilename = info.Name

		w, err := img.openDestinationsForWriting()
		if err != nil {
			return err
		}

		defer w.Close()
		_, err = io.Copy(w, r)
		return err
	}

	// create creates the image and encodes it to w (cache) and to its destinations.
	create := func(info filecache.ItemInfo, w io.WriteCloser) (err error) {
		var conv image.Image
		img, conv, err = createImage()
		if err != nil {
			w.Close()
			return
		}
		img.relTargetDirFile.file = relTarget.file
		img.sourceFilename = info.Name

		destinations, err := img.openDestinationsForWriting()
		if err != nil {
			w.Close()
			return err
		}

		mw := hugio.NewMultiWriteCloser(w, destinations)
		defer mw.Close()

		return img.encodeTo(conf, conv, mw)
	}

	// Now look in the file cache.

	// The definition of this counter is not that we have processed that amount
	// (e.g. resized etc.), it can be fetched from file cache,
	//  but the count of processed image variations for this site.
	c.pathSpec.ProcessingStats.Incr(&c.pathSpec.ProcessingStats.ProcessedImages)

	_, err := c.fileCache.ReadOrCreate(key, read, create)
	if err != nil {
		return nil, err
	}

	// The file is now stored in this cache.
	img.overriddenSourceFs = c.fileCache.Fs

	c.mu.Lock()
	if img2, found := c.store[key]; found {
		c.mu.Unlock()
		return img2, nil
	}
	c.store[key] = img
	c.mu.Unlock()

	return img, nil

}

func newImageCache(fileCache *filecache.Cache, ps *helpers.PathSpec) *imageCache {
	return &imageCache{fileCache: fileCache, pathSpec: ps, store: make(map[string]*Image)}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
