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
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/spf13/afero"
)

var (
	_ afero.Fs      = (*SliceFs)(nil)
	_ afero.Lstater = (*SliceFs)(nil)
	_ afero.File    = (*sliceDir)(nil)
)

func NewSliceFs(dirs ...FileMetaInfo) (afero.Fs, error) {
	if len(dirs) == 0 {
		return NoOpFs, nil
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			return nil, errors.New("this fs supports directories only")
		}
	}

	fs := &SliceFs{
		dirs: dirs,
	}

	return fs, nil

}

// SliceFs is an ordered composite filesystem.
type SliceFs struct {
	dirs []FileMetaInfo
}

func (fs *SliceFs) Chmod(n string, m os.FileMode) error {
	return syscall.EPERM
}

func (fs *SliceFs) Chtimes(n string, a, m time.Time) error {
	return syscall.EPERM
}

func (fs *SliceFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	fi, _, err := fs.pickFirst(name)

	if err != nil {
		return nil, false, err
	}

	if fi.IsDir() {
		return decorateFileInfo(fi, fs, fs.getOpener(name), "", "", nil), false, nil
	}

	return nil, false, errors.Errorf("lstat: files not supported: %q", name)

}

func (fs *SliceFs) Mkdir(n string, p os.FileMode) error {
	return syscall.EPERM
}

func (fs *SliceFs) MkdirAll(n string, p os.FileMode) error {
	return syscall.EPERM
}

func (fs *SliceFs) Name() string {
	return "SliceFs"
}

func (fs *SliceFs) Open(name string) (afero.File, error) {
	fi, idx, err := fs.pickFirst(name)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		panic("currently only dirs in here")
	}

	return &sliceDir{
		lfs:     fs,
		idx:     idx,
		dirname: name,
	}, nil

}

func (fs *SliceFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	panic("not implemented")
}

func (fs *SliceFs) ReadDir(name string) ([]os.FileInfo, error) {
	panic("not implemented")
}

func (fs *SliceFs) Remove(n string) error {
	return syscall.EPERM
}

func (fs *SliceFs) RemoveAll(p string) error {
	return syscall.EPERM
}

func (fs *SliceFs) Rename(o, n string) error {
	return syscall.EPERM
}

func (fs *SliceFs) Stat(name string) (os.FileInfo, error) {
	fi, _, err := fs.LstatIfPossible(name)
	return fi, err
}

func (fs *SliceFs) Create(n string) (afero.File, error) {
	return nil, syscall.EPERM
}

func (fs *SliceFs) getOpener(name string) func() (afero.File, error) {
	return func() (afero.File, error) {
		return fs.Open(name)
	}
}

func (fs *SliceFs) pickFirst(name string) (os.FileInfo, int, error) {
	for i, mfs := range fs.dirs {
		meta := mfs.Meta()
		fs := meta.Fs()
		fi, _, err := lstatIfPossible(fs, name)
		if err == nil {
			// Gotta match!
			return fi, i, nil
		}

		if !os.IsNotExist(err) {
			// Real error
			return nil, -1, err
		}
	}

	// Not found
	return nil, -1, os.ErrNotExist
}

func (fs *SliceFs) readDirs(name string, startIdx, count int) ([]os.FileInfo, error) {
	collect := func(lfs FileMeta) ([]os.FileInfo, error) {
		d, err := lfs.Fs().Open(name)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
			return nil, nil
		} else {
			defer d.Close()
			dirs, err := d.Readdir(-1)
			if err != nil {
				return nil, err
			}
			return dirs, nil
		}
	}

	var dirs []os.FileInfo

	for i := startIdx; i < len(fs.dirs); i++ {
		mfs := fs.dirs[i]

		fis, err := collect(mfs.Meta())
		if err != nil {
			return nil, err
		}

		dirs = append(dirs, fis...)

	}

	seen := make(map[string]bool)
	var duplicates []int
	for i, fi := range dirs {
		if !fi.IsDir() {
			continue
		}

		if seen[fi.Name()] {
			duplicates = append(duplicates, i)
		} else {
			// Make sure it's opened by this filesystem.
			dirs[i] = decorateFileInfo(fi, fs, fs.getOpener(fi.(FileMetaInfo).Meta().Filename()), "", "", nil)
			seen[fi.Name()] = true
		}
	}

	// Remove duplicate directories, keep first.
	if len(duplicates) > 0 {
		for i := len(duplicates) - 1; i >= 0; i-- {
			idx := duplicates[i]
			dirs = append(dirs[:idx], dirs[idx+1:]...)
		}
	}

	if count > 0 && len(dirs) >= count {
		return dirs[:count], nil
	}

	return dirs, nil

}

type sliceDir struct {
	lfs     *SliceFs
	idx     int
	dirname string
}

func (f *sliceDir) Close() error {
	return nil
}

func (f *sliceDir) Name() string {
	return f.dirname
}

func (f *sliceDir) Read(p []byte) (n int, err error) {
	panic("not implemented")
}

func (f *sliceDir) ReadAt(p []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (f *sliceDir) Readdir(count int) ([]os.FileInfo, error) {
	return f.lfs.readDirs(f.dirname, f.idx, count)
}

func (f *sliceDir) Readdirnames(count int) ([]string, error) {
	dirsi, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, len(dirsi))
	for i, d := range dirsi {
		dirs[i] = d.Name()
	}
	return dirs, nil
}

func (f *sliceDir) Seek(offset int64, whence int) (int64, error) {
	panic("not implemented")
}

func (f *sliceDir) Stat() (os.FileInfo, error) {
	panic("not implemented")
}

func (f *sliceDir) Sync() error {
	panic("not implemented")
}

func (f *sliceDir) Truncate(size int64) error {
	panic("not implemented")
}

func (f *sliceDir) Write(p []byte) (n int, err error) {
	panic("not implemented")
}

func (f *sliceDir) WriteAt(p []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (f *sliceDir) WriteString(s string) (ret int, err error) {
	panic("not implemented")
}
