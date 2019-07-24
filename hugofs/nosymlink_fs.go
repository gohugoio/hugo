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

package hugofs

import (
	"errors"
	"os"

	"github.com/spf13/afero"
)

var ErrPermissionSymlink = errors.New("symlinks not allowed in this filesystem")

func NewNoSymlinkFs(fs afero.Fs) afero.Fs {
	return &noSymlinkFs{Fs: fs}
}

// noSymlinkFs is a filesystem that prevents symlinking.
type noSymlinkFs struct {
	afero.Fs
}

func (fs *noSymlinkFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	return fs.stat(name)
}

func (fs *noSymlinkFs) Stat(name string) (os.FileInfo, error) {
	fi, _, err := fs.stat(name)
	return fi, err
}

func (fs *noSymlinkFs) stat(name string) (os.FileInfo, bool, error) {
	var (
		fi       os.FileInfo
		wasLstat bool
		err      error
	)

	if lstater, ok := fs.Fs.(afero.Lstater); ok {
		fi, wasLstat, err = lstater.LstatIfPossible(name)
	} else {
		fi, err = fs.Fs.Stat(name)
	}

	var metaIsSymlink bool

	if fim, ok := fi.(FileMetaInfo); ok {
		metaIsSymlink = fim.Meta().IsSymlink()
	}

	if metaIsSymlink || isSymlink(fi) {
		return nil, wasLstat, ErrPermissionSymlink
	}

	return fi, wasLstat, err
}

func (fs *noSymlinkFs) Open(name string) (afero.File, error) {
	if _, _, err := fs.stat(name); err != nil {
		return nil, err
	}
	return fs.Fs.Open(name)
}

func (fs *noSymlinkFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if _, _, err := fs.stat(name); err != nil {
		return nil, err
	}
	return fs.Fs.OpenFile(name, flag, perm)
}
