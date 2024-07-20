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

package filecache

import (
	"fmt"
	"io"
	"os"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/afero"
)

// Prune removes expired and unused items from this cache.
// The last one requires a full build so the cache usage can be tracked.
// Note that we operate directly on the filesystem here, so this is not
// thread safe.
func (c Caches) Prune() (int, error) {
	counter := 0
	for k, cache := range c {
		count, err := cache.Prune(false)

		counter += count

		if err != nil {
			if herrors.IsNotExist(err) {
				continue
			}
			return counter, fmt.Errorf("failed to prune cache %q: %w", k, err)
		}

	}

	return counter, nil
}

// Prune removes expired and unused items from this cache.
// If force is set, everything will be removed not considering expiry time.
func (c *Cache) Prune(force bool) (int, error) {
	if c.pruneAllRootDir != "" {
		return c.pruneRootDir(force)
	}
	if err := c.init(); err != nil {
		return 0, err
	}

	counter := 0

	err := afero.Walk(c.Fs, "", func(name string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		name = cleanID(name)

		if info.IsDir() {
			f, err := c.Fs.Open(name)
			if err != nil {
				// This cache dir may not exist.
				return nil
			}
			_, err = f.Readdirnames(1)
			f.Close()
			if err == io.EOF {
				// Empty dir.
				if name == "." {
					// e.g. /_gen/images -- keep it even if empty.
					err = nil
				} else {
					err = c.Fs.Remove(name)
				}
			}

			if err != nil && !herrors.IsNotExist(err) {
				return err
			}

			return nil
		}

		shouldRemove := force || c.isExpired(info.ModTime())

		if !shouldRemove && len(c.nlocker.seen) > 0 {
			// Remove it if it's not been touched/used in the last build.
			_, seen := c.nlocker.seen[name]
			shouldRemove = !seen
		}

		if shouldRemove {
			err := c.Fs.Remove(name)
			if err == nil {
				counter++
			}

			if err != nil && !herrors.IsNotExist(err) {
				return err
			}

		}

		return nil
	})

	return counter, err
}

func (c *Cache) pruneRootDir(force bool) (int, error) {
	if err := c.init(); err != nil {
		return 0, err
	}
	info, err := c.Fs.Stat(c.pruneAllRootDir)
	if err != nil {
		if herrors.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	if !force && !c.isExpired(info.ModTime()) {
		return 0, nil
	}

	return hugofs.MakeReadableAndRemoveAllModulePkgDir(c.Fs, c.pruneAllRootDir)
}
