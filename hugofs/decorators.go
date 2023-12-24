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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var _ FilesystemUnwrapper = (*baseFileDecoratorFs)(nil)

func decorateDirs(fs afero.Fs, meta *FileMeta) afero.Fs {
	ffs := &baseFileDecoratorFs{Fs: fs}

	decorator := func(fi FileNameIsDir, name string) (FileNameIsDir, error) {
		if !fi.IsDir() {
			// Leave regular files as they are.
			return fi, nil
		}

		return decorateFileInfo(fi, nil, "", meta), nil
	}

	ffs.decorate = decorator

	return ffs
}

// NewBaseFileDecorator decorates the given Fs to provide the real filename
// and an Opener func.
func NewBaseFileDecorator(fs afero.Fs, callbacks ...func(fi FileMetaInfo)) afero.Fs {
	ffs := &baseFileDecoratorFs{Fs: fs}

	decorator := func(fi FileNameIsDir, filename string) (FileNameIsDir, error) {
		// Store away the original in case it's a symlink.
		meta := NewFileMeta()
		meta.Name = fi.Name()

		if fi.IsDir() {
			meta.JoinStatFunc = func(name string) (FileMetaInfo, error) {
				joinedFilename := filepath.Join(filename, name)
				fi, err := fs.Stat(joinedFilename)
				if err != nil {
					return nil, err
				}
				fim, err := ffs.decorate(fi, joinedFilename)
				if err != nil {
					return nil, err
				}

				return fim.(FileMetaInfo), nil
			}
		}

		opener := func() (afero.File, error) {
			return ffs.open(filename)
		}

		fim := decorateFileInfo(fi, opener, filename, meta)

		for _, cb := range callbacks {
			cb(fim)
		}

		return fim, nil
	}

	ffs.decorate = decorator
	return ffs
}

type baseFileDecoratorFs struct {
	afero.Fs
	decorate func(fi FileNameIsDir, name string) (FileNameIsDir, error)
}

func (fs *baseFileDecoratorFs) UnwrapFilesystem() afero.Fs {
	return fs.Fs
}

func (fs *baseFileDecoratorFs) Stat(name string) (os.FileInfo, error) {
	fi, err := fs.Fs.Stat(name)
	if err != nil {
		return nil, err
	}

	fim, err := fs.decorate(fi, name)
	if err != nil {
		return nil, err
	}
	return fim.(os.FileInfo), nil
}

func (fs *baseFileDecoratorFs) Open(name string) (afero.File, error) {
	return fs.open(name)
}

func (fs *baseFileDecoratorFs) open(name string) (afero.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	return &baseFileDecoratorFile{File: f, fs: fs}, nil
}

type baseFileDecoratorFile struct {
	afero.File
	fs *baseFileDecoratorFs
}

func (l *baseFileDecoratorFile) ReadDir(n int) ([]fs.DirEntry, error) {
	fis, err := l.File.(fs.ReadDirFile).ReadDir(-1)
	if err != nil {
		return nil, err
	}

	fisp := make([]fs.DirEntry, len(fis))

	for i, fi := range fis {
		filename := fi.Name()
		if l.Name() != "" {
			filename = filepath.Join(l.Name(), fi.Name())
		}

		fid, err := l.fs.decorate(fi, filename)
		if err != nil {
			return nil, fmt.Errorf("decorate: %w", err)
		}

		fisp[i] = fid.(fs.DirEntry)

	}

	return fisp, err
}

func (l *baseFileDecoratorFile) Readdir(c int) (ofi []os.FileInfo, err error) {
	panic("not supported: Use ReadDir")
}
