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
	"path"
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
	root := resolveRootDirCase(fs, hglob.ResolveRootDir(pattern))
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

// resolveRootDirCase resolves the casing of each path element in root (which may be lowercased)
// by looking it up lowercased in its parent's directory listing.
// We need this to walk from the correct root on case-sensitive file systems,
// and to preserve the original casing in the walked paths.
// If no match is found (e.g. a virtual mount dir not listed in its parent), the element is used as is.
func resolveRootDirCase(fs afero.Fs, root string) string {
	resolved := "/"
	for part := range strings.SplitSeq(root, "/") {
		if part == "" {
			continue
		}
		resolved = path.Join(resolved, findDirEntryName(fs, resolved, part))
	}
	return resolved
}

func findDirEntryName(fs afero.Fs, dir, name string) string {
	f, err := fs.Open(dir)
	if err != nil {
		return name
	}
	defer f.Close()
	fis, _, err := ReadDirWithContext(context.Background(), f, -1)
	if err != nil {
		return name
	}
	nameLower := strings.ToLower(name)
	var candidate string
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		if fi.Name() == name {
			return name
		}
		// Use the same lowercasing as the glob matcher (not strings.EqualFold, which differs for e.g. "İ").
		if candidate == "" && strings.ToLower(fi.Name()) == nameLower {
			candidate = fi.Name()
		}
	}
	if candidate != "" {
		return candidate
	}
	return name
}
