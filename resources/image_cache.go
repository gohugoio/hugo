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

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/images"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/helpers"
)

type imageCache struct {
	pathSpec *helpers.PathSpec

	fileCache *filecache.Cache

	mu    sync.RWMutex
	store map[string]baseResource
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
	c.store = make(map[string]baseResource)
}

func (c *imageCache) getOrCreate(
	parent *imageResource, conf images.ImageConfig,
	createImage func() (baseResource, interface{}, error)) (baseResource, error) {

	relTarget := parent.relTargetPathFromConfig(conf)
	key := parent.relTargetPathForRel(relTarget.path(), false, false, false)

	var img baseResource
	var found bool

	// First check the in-memory store, then the disk.
	c.mu.RLock()
	img, found = c.store[key]
	c.mu.RUnlock()

	if found {
		return img, nil
	}

	// These funcs are protected by a named lock.
	// read clones the parent to its new name and copies
	// the content to the destinations.
	read := func(info filecache.ItemInfo, r io.Reader) error {
		if conf.Action == "trace" {
			// trace produces a SVG
			img = parent.baseResource.Clone().(baseResource)
			img.setMediaType(media.SVGType)

		} else {
			img = parent.clone(nil)
		}

		rp := img.getResourcePaths()
		rp.relTargetDirFile.file = relTarget.file
		img.setSourceFilename(info.Name)

		w, err := img.openDestinationsForWriting()
		if err != nil {
			return err
		}

		if w == nil {
			// Nothing to write.
			return nil
		}

		defer w.Close()
		_, err = io.Copy(w, r)
		return err
	}

	// create creates the image and encodes it to w (cache) and to its destinations.
	create := func(info filecache.ItemInfo, w io.WriteCloser) (err error) {
		var conv interface{}
		img, conv, err = createImage()
		if err != nil {
			w.Close()
			return
		}

		rp := img.getResourcePaths()
		rp.relTargetDirFile.file = relTarget.file
		img.setSourceFilename(info.Name)

		destinations, err := img.openDestinationsForWriting()
		if err != nil {
			w.Close()
			return err
		}

		if destinations != nil {
			w = hugio.NewMultiWriteCloser(w, destinations)
		}
		defer w.Close()

		switch v := conv.(type) {
		case string:
			_, err := fmt.Fprint(w, v)
			return err
		case image.Image:
			return img.(*imageResource).EncodeTo(conf, v, w)
		default:
			panic(fmt.Sprintf("unknown type %T", conv))
		}
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
	img.setSourceFs(c.fileCache.Fs)

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
	return &imageCache{fileCache: fileCache, pathSpec: ps, store: make(map[string]baseResource)}
}
