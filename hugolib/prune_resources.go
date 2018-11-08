// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"io"
	"os"
	"strings"

	"github.com/gohugoio/hugo/helpers"

	"github.com/spf13/afero"
)

// GC requires a build first.
func (h *HugoSites) GC() (int, error) {
	s := h.Sites[0]
	assetsCacheFs := h.Deps.FileCaches.AssetsCache().Fs
	imageCacheFs := h.Deps.FileCaches.ImageCache().Fs

	isImageInUse := func(name string) bool {
		for _, site := range h.Sites {
			if site.ResourceSpec.IsInImageCache(name) {
				return true
			}
		}

		return false
	}

	isAssetInUse := func(name string) bool {
		// These assets are stored in tuplets with an added extension to the key.
		key := strings.TrimSuffix(name, helpers.Ext(name))
		for _, site := range h.Sites {
			if site.ResourceSpec.ResourceCache.Contains(key) {
				return true
			}
		}

		return false
	}

	walker := func(fs afero.Fs, dirname string, inUse func(filename string) bool) (int, error) {
		counter := 0
		err := afero.Walk(fs, dirname, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}

			if info.IsDir() {
				f, err := fs.Open(path)
				if err != nil {
					return nil
				}
				defer f.Close()
				_, err = f.Readdirnames(1)
				if err == io.EOF {
					// Empty dir.
					s.Fs.Source.Remove(path)
				}

				return nil
			}

			inUse := inUse(path)
			if !inUse {
				err := fs.Remove(path)
				if err != nil && !os.IsNotExist(err) {
					s.Log.ERROR.Printf("Failed to remove %q: %s", path, err)
				} else {
					counter++
				}
			}
			return nil
		})

		return counter, err
	}

	imageCounter, err1 := walker(imageCacheFs, "", isImageInUse)
	assetsCounter, err2 := walker(assetsCacheFs, "", isAssetInUse)
	totalCount := imageCounter + assetsCounter

	if err1 != nil {
		return totalCount, err1
	}

	return totalCount, err2

}
