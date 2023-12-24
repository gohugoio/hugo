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

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources/images"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/helpers"
)

// ImageCache is a cache for image resources. The backing caches are shared between all sites.
type ImageCache struct {
	pathSpec *helpers.PathSpec

	fcache *filecache.Cache
	mcache *dynacache.Partition[string, *resourceAdapter]
}

func (c *ImageCache) getOrCreate(
	parent *imageResource, conf images.ImageConfig,
	createImage func() (*imageResource, image.Image, error),
) (*resourceAdapter, error) {
	relTarget := parent.relTargetPathFromConfig(conf)
	relTargetPath := relTarget.TargetPath()
	memKey := dynacache.CleanKey(relTargetPath)

	v, err := c.mcache.GetOrCreate(memKey, func(key string) (*resourceAdapter, error) {
		var img *imageResource

		// These funcs are protected by a named lock.
		// read clones the parent to its new name and copies
		// the content to the destinations.
		read := func(info filecache.ItemInfo, r io.ReadSeeker) error {
			img = parent.clone(nil)
			targetPath := img.getResourcePaths()
			targetPath.File = relTarget.File
			img.setTargetPath(targetPath)
			img.setOpenSource(func() (hugio.ReadSeekCloser, error) {
				return c.fcache.Fs.Open(info.Name)
			})
			img.setSourceFilenameIsHash(true)
			img.setMediaType(conf.TargetFormat.MediaType())

			if err := img.InitConfig(r); err != nil {
				return err
			}

			return nil
		}

		// create creates the image and encodes it to the cache (w).
		create := func(info filecache.ItemInfo, w io.WriteCloser) (err error) {
			defer w.Close()

			var conv image.Image
			img, conv, err = createImage()
			if err != nil {
				return
			}
			targetPath := img.getResourcePaths()
			targetPath.File = relTarget.File
			img.setTargetPath(targetPath)
			img.setOpenSource(func() (hugio.ReadSeekCloser, error) {
				return c.fcache.Fs.Open(info.Name)
			})
			return img.EncodeTo(conf, conv, w)
		}

		// Now look in the file cache.

		// The definition of this counter is not that we have processed that amount
		// (e.g. resized etc.), it can be fetched from file cache,
		//  but the count of processed image variations for this site.
		c.pathSpec.ProcessingStats.Incr(&c.pathSpec.ProcessingStats.ProcessedImages)

		_, err := c.fcache.ReadOrCreate(relTargetPath, read, create)
		if err != nil {
			return nil, err
		}

		imgAdapter := newResourceAdapter(parent.getSpec(), true, img)

		return imgAdapter, nil
	})

	return v, err
}

func newImageCache(fileCache *filecache.Cache, memCache *dynacache.Cache, ps *helpers.PathSpec) *ImageCache {
	return &ImageCache{
		fcache: fileCache,
		mcache: dynacache.GetOrCreatePartition[string, *resourceAdapter](
			memCache,
			"/imgs",
			dynacache.OptionsPartition{ClearWhen: dynacache.ClearOnChange, Weight: 70},
		),
		pathSpec: ps,
	}
}
