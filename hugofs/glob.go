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

package hugofs

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs/hglob"
	"github.com/spf13/afero"
)

// Glob walks the fs and passes all matches to the handle func.
// The handle func can return true to signal a stop.
func Glob(fs afero.Fs, pattern string, handle func(fi FileMetaInfo) (bool, error)) error {
	pattern = hglob.NormalizePathNoLower(pattern)
	if pattern == "" {
		return nil
	}
	root := hglob.ResolveRootDir(pattern)
	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}
	pattern = strings.ToLower(pattern)

	g, err := hglob.GetGlob(pattern)
	if err != nil {
		return err
	}

	hasSuperAsterisk := strings.Contains(pattern, "**")
	levels := strings.Count(pattern, "/")

	// Signals that we're done.
	done := errors.New("done")

	wfn := func(ctx context.Context, p string, info FileMetaInfo) error {
		p = hglob.NormalizePath(p)
		if info.IsDir() {
			if !hasSuperAsterisk {
				// Avoid walking to the bottom if we can avoid it.
				if p != "" && strings.Count(p, "/") >= levels {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if g.Match(p) {
			d, err := handle(info)
			if err != nil {
				return err
			}
			if d {
				return done
			}
		}

		return nil
	}

	w := NewWalkway(
		WalkwayConfig{
			Root:           root,
			Fs:             fs,
			WalkFn:         wfn,
			FailOnNotExist: true,
		})

	err = w.Walk()

	if err != done {
		return err
	}

	return nil
}
