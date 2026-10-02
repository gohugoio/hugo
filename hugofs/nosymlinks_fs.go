// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/spf13/afero"
)

// symlinkChecker caches Lstat results for directories.
// We assume that symlinks are not created or removed while Hugo is running.
type symlinkChecker struct {
	fs    afero.Fs
	cache *hmaps.Cache[string, bool]
}

func newSymlinkChecker(fs afero.Fs) *symlinkChecker {
	return &symlinkChecker{fs: fs, cache: hmaps.NewCacheWithOptions[string, bool](hmaps.CacheOptions{Size: 10000})}
}

// isSymlink reports whether name is a symlink. A missing name is not a symlink.
func (c *symlinkChecker) isSymlink(name string) (bool, error) {
	return c.cache.GetOrCreate(name, func() (bool, error) {
		fi, err := LstatIfPossible(c.fs, name)
		if err != nil {
			if herrors.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		return fi.Mode()&os.ModeSymlink != 0, nil
	})
}

// hasSymlinkParent reports whether dir or any of its parents below base is a symlink.
// base itself is not checked; an empty base walks to the root of fs.
func (c *symlinkChecker) hasSymlinkParent(base, dir string) (bool, error) {
	for dir != base && dir != "." {
		parent := filepath.Dir(dir)
		if parent == dir {
			return false, nil
		}
		symlink, err := c.isSymlink(dir)
		if err != nil || symlink {
			return symlink, err
		}
		dir = parent
	}
	return false, nil
}

// isSymlinkOrHasSymlinkParent reports whether name is a symlink or has a symlinked parent below base.
// Parents are only checked when base is set; absolute mount sources may legitimately
// live below symlinked directories (e.g. /tmp on macOS).
func (c *symlinkChecker) isSymlinkOrHasSymlinkParent(base, name string) (bool, error) {
	symlink, err := c.isSymlink(name)
	if err != nil || symlink {
		return symlink, err
	}
	if base == "" {
		return false, nil
	}
	return c.hasSymlinkParent(filepath.Clean(base), filepath.Dir(name))
}

// NewDropSymlinksFs returns an afero.Fs wrapper that treats symlinks as non-existing files.
func NewDropSymlinksFs(base afero.Fs) *DropSymlinksFs {
	return &DropSymlinksFs{Fs: base, symlinks: newSymlinkChecker(base)}
}

// DropSymlinksFs is an afero.Fs wrapper that treats symlinks as non-existing files.
type DropSymlinksFs struct {
	afero.Fs
	symlinks *symlinkChecker
}

func (fs *DropSymlinksFs) Open(name string) (afero.File, error) {
	if _, err := fs.Stat(name); err != nil {
		return nil, err
	}
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs *DropSymlinksFs) Stat(name string) (os.FileInfo, error) {
	fi, err := LstatIfPossible(fs.Fs, name)
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return nil, os.ErrNotExist
	}
	symlinkParent, err := fs.symlinks.hasSymlinkParent("", filepath.Dir(name))
	if err != nil {
		return nil, err
	}
	if symlinkParent {
		return nil, os.ErrNotExist
	}
	return fi, nil
}
