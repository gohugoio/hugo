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
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Prune removes expired and unused items from this cache.
// The last one requires a full build so the cache usage can be tracked.
// Note that we operate directly on the filesystem here, so this is not
// thread safe.
func (c Caches) Prune() (int, error) {
	counter := 0
	for k, cache := range c {
		err := afero.Walk(cache.Fs, "", func(name string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}

			name = cleanID(name)

			if info.IsDir() {
				f, err := cache.Fs.Open(name)
				if err != nil {
					// This cache dir may not exist.
					return nil
				}
				defer f.Close()
				_, err = f.Readdirnames(1)
				if err == io.EOF {
					// Empty dir.
					return cache.Fs.Remove(name)
				}

				return nil
			}

			shouldRemove := cache.isExpired(info.ModTime())

			if !shouldRemove && len(cache.nlocker.seen) > 0 {
				// Remove it if it's not been touched/used in the last build.
				_, seen := cache.nlocker.seen[name]
				shouldRemove = !seen
			}

			if shouldRemove {
				err := cache.Fs.Remove(name)
				if err == nil {
					counter++
				}
				return err
			}

			return nil
		})

		if err != nil {
			return counter, errors.Wrapf(err, "failed to prune cache %q", k)
		}

	}

	return counter, nil
}
