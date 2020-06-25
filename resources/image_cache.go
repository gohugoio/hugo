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
	"image"
	"io"
	"path/filepath"

	"github.com/gohugoio/hugo/cache/memcache"

	"github.com/gohugoio/hugo/resources/images"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/helpers"
)

type imageCache struct {
	pathSpec *helpers.PathSpec

	fileCache *filecache.Cache
	mCache    memcache.Getter
}

func (c *imageCache) getOrCreate(
	parent *imageResource, conf images.ImageConfig,
	createImage func() (*imageResource, image.Image, error)) (*resourceAdapter, error) {
	relTarget := parent.relTargetPathFromConfig(conf)
	memKey := parent.relTargetPathForRel(relTarget.path(), false, false, false)
	memKey = memcache.CleanKey(memKey)

	v, err := c.mCache.GetOrCreate(memKey, func() memcache.Entry {
		// For the file cache we want to generate and store it once if possible.
		fileKeyPath := relTarget
		if fi := parent.root.getFileInfo(); fi != nil {
			fileKeyPath.dir = filepath.ToSlash(filepath.Dir(fi.Meta().Path()))
		}
		fileKey := fileKeyPath.path()

		var img *imageResource

		// These funcs are protected by a named lock.
		// read clones the parent to its new name and copies
		// the content to the destinations.
		read := func(info filecache.ItemInfo, r io.ReadSeeker) error {
			img = parent.clone(nil)
			rp := img.getResourcePaths()
			rp.relTargetDirFile.file = relTarget.file
			img.setSourceFilename(info.Name)

			if err := img.InitConfig(r); err != nil {
				return err
			}

			r.Seek(0, 0)

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

		// create creates the image and encodes it to the cache (w).
		create := func(info filecache.ItemInfo, w io.WriteCloser) (err error) {
			defer w.Close()

			var conv image.Image
			img, conv, err = createImage()
			if err != nil {
				return
			}
			rp := img.getResourcePaths()
			rp.relTargetDirFile.file = relTarget.file
			img.setSourceFilename(info.Name)

			return img.EncodeTo(conf, conv, w)
		}

		// Now look in the file cache.

		// The definition of this counter is not that we have processed that amount
		// (e.g. resized etc.), it can be fetched from file cache,
		//  but the count of processed image variations for this site.
		c.pathSpec.ProcessingStats.Incr(&c.pathSpec.ProcessingStats.ProcessedImages)

		_, err := c.fileCache.ReadOrCreate(fileKey, read, create)
		if err != nil {
			return memcache.Entry{Err: err}
		}

		// The file is now stored in this cache.
		img.setSourceFs(c.fileCache.Fs)

		imgAdapter := newResourceAdapter(parent.getSpec(), true, img)

		return memcache.Entry{Value: imgAdapter, ClearWhen: memcache.ClearOnChange}
	})

	if err != nil {
		return nil, err
	}
	return v.(*resourceAdapter), nil

}

func newImageCache(fileCache *filecache.Cache, memCache *memcache.Cache, ps *helpers.PathSpec) *imageCache {
	return &imageCache{fileCache: fileCache, mCache: memCache.GetOrCreatePartition("images", memcache.ClearOnChange), pathSpec: ps}
}
