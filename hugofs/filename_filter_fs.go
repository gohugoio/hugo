// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"strings"
	"syscall"
	"time"

	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/spf13/afero"
)

func newFilenameFilterFs(fs afero.Fs, base string, filter *glob.FilenameFilter) afero.Fs {
	return &filenameFilterFs{
		fs:     fs,
		base:   base,
		filter: filter,
	}
}

// filenameFilterFs is a filesystem that filters by filename.
type filenameFilterFs struct {
	base string
	fs   afero.Fs

	filter *glob.FilenameFilter
}

func (fs *filenameFilterFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	fi, b, err := fs.fs.(afero.Lstater).LstatIfPossible(name)
	if err != nil {
		return nil, false, err
	}
	if !fs.filter.Match(name, fi.IsDir()) {
		return nil, false, os.ErrNotExist
	}
	return fi, b, nil
}

func (fs *filenameFilterFs) Open(name string) (afero.File, error) {
	fi, err := fs.fs.Stat(name)
	if err != nil {
		return nil, err
	}

	if !fs.filter.Match(name, fi.IsDir()) {
		return nil, os.ErrNotExist
	}

	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return f, nil
	}

	return &filenameFilterDir{
		File:   f,
		base:   fs.base,
		filter: fs.filter,
	}, nil
}

func (fs *filenameFilterFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return fs.Open(name)
}

func (fs *filenameFilterFs) Stat(name string) (os.FileInfo, error) {
	fi, _, err := fs.LstatIfPossible(name)
	return fi, err
}

func (fs *filenameFilterFs) getOpener(name string) func() (afero.File, error) {
	return func() (afero.File, error) {
		return fs.Open(name)
	}
}

type filenameFilterDir struct {
	afero.File
	base   string
	filter *glob.FilenameFilter
}

func (f *filenameFilterDir) Readdir(count int) ([]os.FileInfo, error) {
	fis, err := f.File.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var result []os.FileInfo
	for _, fi := range fis {
		fim := fi.(FileMetaInfo)
		if f.filter.Match(strings.TrimPrefix(fim.Meta().Filename, f.base), fim.IsDir()) {
			result = append(result, fi)
		}
	}

	return result, nil
}

func (f *filenameFilterDir) Readdirnames(count int) ([]string, error) {
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

func (fs *filenameFilterFs) Chmod(n string, m os.FileMode) error {
	return syscall.EPERM
}

func (fs *filenameFilterFs) Chtimes(n string, a, m time.Time) error {
	return syscall.EPERM
}

func (fs *filenameFilterFs) Chown(n string, uid, gid int) error {
	return syscall.EPERM
}

func (fs *filenameFilterFs) ReadDir(name string) ([]os.FileInfo, error) {
	panic("not implemented")
}

func (fs *filenameFilterFs) Remove(n string) error {
	return syscall.EPERM
}

func (fs *filenameFilterFs) RemoveAll(p string) error {
	return syscall.EPERM
}

func (fs *filenameFilterFs) Rename(o, n string) error {
	return syscall.EPERM
}
func (fs *filenameFilterFs) Create(n string) (afero.File, error) {
	return nil, syscall.EPERM
}
func (fs *filenameFilterFs) Name() string {
	return "FinameFilterFS"
}

func (fs *filenameFilterFs) Mkdir(n string, p os.FileMode) error {
	return syscall.EPERM
}

func (fs *filenameFilterFs) MkdirAll(n string, p os.FileMode) error {
	return syscall.EPERM
}
