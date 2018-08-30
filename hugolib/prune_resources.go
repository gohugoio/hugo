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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gohugoio/hugo/helpers"

	"github.com/spf13/afero"
)

// GC requires a build first.
func (h *HugoSites) GC() (int, error) {
	s := h.Sites[0]
	fs := h.PathSpec.BaseFs.Resources.Fs

	imageCacheDir := s.ResourceSpec.GenImagePath
	if len(imageCacheDir) < 10 {
		panic("invalid image cache")
	}
	assetsCacheDir := s.ResourceSpec.GenAssetsPath
	if len(assetsCacheDir) < 10 {
		panic("invalid assets cache")
	}

	isImageInUse := func(filename string) bool {
		key := strings.TrimPrefix(filename, imageCacheDir)
		for _, site := range h.Sites {
			if site.ResourceSpec.IsInImageCache(key) {
				return true
			}
		}

		return false
	}

	isAssetInUse := func(filename string) bool {
		key := strings.TrimPrefix(filename, assetsCacheDir)
		// These assets are stored in tuplets with an added extension to the key.
		key = strings.TrimSuffix(key, helpers.Ext(key))
		for _, site := range h.Sites {
			if site.ResourceSpec.ResourceCache.Contains(key) {
				return true
			}
		}

		return false
	}

	walker := func(dirname string, inUse func(filename string) bool) (int, error) {
		counter := 0
		err := afero.Walk(fs, dirname, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}

			if !strings.HasPrefix(path, dirname) {
				return fmt.Errorf("Invalid state, walk outside of resource dir: %q", path)
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

	imageCounter, err1 := walker(imageCacheDir, isImageInUse)
	assetsCounter, err2 := walker(assetsCacheDir, isAssetInUse)
	totalCount := imageCounter + assetsCounter

	if err1 != nil {
		return totalCount, err1
	}

	return totalCount, err2

}
