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
	"github.com/spf13/afero"
)

// hasSymlinkParent reports whether any directory between name and base is a symlink.
// Lstat only refrains from following the last element of a path, so callers checking
// name itself need this to keep a symlinked parent from escaping base.
// base itself is not checked; an empty base walks to the root of fs.
func hasSymlinkParent(fs afero.Fs, base, name string) (bool, error) {
	for name != base {
		parent := filepath.Dir(name)
		if parent == name || parent == base || parent == "." {
			return false, nil
		}
		fi, err := LstatIfPossible(fs, parent)
		if err != nil {
			if herrors.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return true, nil
		}
		name = parent
	}
	return false, nil
}

// NewDropSymlinksFs returns an afero.Fs wrapper that treats symlinks as non-existing files.
func NewDropSymlinksFs(base afero.Fs) *DropSymlinksFs {
	return &DropSymlinksFs{base}
}

// DropSymlinksFs is an afero.Fs wrapper that treats symlinks as non-existing files.
type DropSymlinksFs struct {
	afero.Fs
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
	symlinkParent, err := hasSymlinkParent(fs.Fs, "", name)
	if err != nil {
		return nil, err
	}
	if symlinkParent {
		return nil, os.ErrNotExist
	}
	return fi, nil
}
