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
	"path/filepath"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/spf13/afero"
)

var (
	ErrPermissionSymlink = errors.New("symlinks not allowed in this filesystem")
)

// NewNoSymlinkFs creates a new filesystem that prevents symlinks.
func NewNoSymlinkFs(fs afero.Fs, logger *loggers.Logger, allowFiles bool) afero.Fs {
	return &noSymlinkFs{Fs: fs, logger: logger, allowFiles: allowFiles}
}

// noSymlinkFs is a filesystem that prevents symlinking.
type noSymlinkFs struct {
	allowFiles bool // block dirs only
	logger     *loggers.Logger
	afero.Fs
}

type noSymlinkFile struct {
	fs *noSymlinkFs
	afero.File
}

func (f *noSymlinkFile) Readdir(count int) ([]os.FileInfo, error) {
	fis, err := f.File.Readdir(count)

	filtered := fis[:0]
	for _, x := range fis {
		filename := filepath.Join(f.Name(), x.Name())
		if _, err := f.fs.checkSymlinkStatus(filename, x); err != nil {
			// Log a warning and drop the file from the list
			logUnsupportedSymlink(filename, f.fs.logger)
		} else {
			filtered = append(filtered, x)
		}
	}

	return filtered, err
}

func (f *noSymlinkFile) Readdirnames(count int) ([]string, error) {
	dirs, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}
	return fileInfosToNames(dirs), nil
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

	if err != nil {
		return nil, false, err
	}

	fi, err = fs.checkSymlinkStatus(name, fi)

	return fi, wasLstat, err
}

func (fs *noSymlinkFs) checkSymlinkStatus(name string, fi os.FileInfo) (os.FileInfo, error) {
	var metaIsSymlink bool

	if fim, ok := fi.(FileMetaInfo); ok {
		meta := fim.Meta()
		metaIsSymlink = meta.IsSymlink()
	}

	if metaIsSymlink {
		if fs.allowFiles && !fi.IsDir() {
			return fi, nil
		}
		return nil, ErrPermissionSymlink
	}

	// Also support non-decorated filesystems, e.g. the Os fs.
	if isSymlink(fi) {
		// Need to determine if this is a directory or not.
		_, sfi, err := evalSymlinks(fs.Fs, name)
		if err != nil {
			return nil, err
		}
		if fs.allowFiles && !sfi.IsDir() {
			// Return the original FileInfo to get the expected Name.
			return fi, nil
		}
		return nil, ErrPermissionSymlink
	}

	return fi, nil
}

func (fs *noSymlinkFs) Open(name string) (afero.File, error) {
	if _, _, err := fs.stat(name); err != nil {
		return nil, err
	}
	return fs.wrapFile(fs.Fs.Open(name))
}

func (fs *noSymlinkFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if _, _, err := fs.stat(name); err != nil {
		return nil, err
	}
	return fs.wrapFile(fs.Fs.OpenFile(name, flag, perm))
}

func (fs *noSymlinkFs) wrapFile(f afero.File, err error) (afero.File, error) {
	if err != nil {
		return nil, err
	}

	return &noSymlinkFile{File: f, fs: fs}, nil
}
