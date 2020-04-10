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
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/afero"
)

func decorateDirs(fs afero.Fs, meta FileMeta) afero.Fs {
	ffs := &baseFileDecoratorFs{Fs: fs}

	decorator := func(fi os.FileInfo, name string) (os.FileInfo, error) {
		if !fi.IsDir() {
			// Leave regular files as they are.
			return fi, nil
		}

		return decorateFileInfo(fi, fs, nil, "", "", meta), nil
	}

	ffs.decorate = decorator

	return ffs

}

func decoratePath(fs afero.Fs, createPath func(name string) string) afero.Fs {

	ffs := &baseFileDecoratorFs{Fs: fs}

	decorator := func(fi os.FileInfo, name string) (os.FileInfo, error) {
		path := createPath(name)

		return decorateFileInfo(fi, fs, nil, "", path, nil), nil
	}

	ffs.decorate = decorator

	return ffs

}

// DecorateBasePathFs adds Path info to files and directories in the
// provided BasePathFs, using the base as base.
func DecorateBasePathFs(base *afero.BasePathFs) afero.Fs {
	basePath, _ := base.RealPath("")
	if !strings.HasSuffix(basePath, filepathSeparator) {
		basePath += filepathSeparator
	}

	ffs := &baseFileDecoratorFs{Fs: base}

	decorator := func(fi os.FileInfo, name string) (os.FileInfo, error) {
		path := strings.TrimPrefix(name, basePath)

		return decorateFileInfo(fi, base, nil, "", path, nil), nil
	}

	ffs.decorate = decorator

	return ffs
}

// NewBaseFileDecorator decorates the given Fs to provide the real filename
// and an Opener func.
func NewBaseFileDecorator(fs afero.Fs, callbacks ...func(fi FileMetaInfo)) afero.Fs {

	ffs := &baseFileDecoratorFs{Fs: fs}

	decorator := func(fi os.FileInfo, filename string) (os.FileInfo, error) {
		// Store away the original in case it's a symlink.
		meta := FileMeta{metaKeyName: fi.Name()}
		if fi.IsDir() {
			meta[metaKeyJoinStat] = func(name string) (FileMetaInfo, error) {
				joinedFilename := filepath.Join(filename, name)
				fi, _, err := lstatIfPossible(fs, joinedFilename)
				if err != nil {
					return nil, err
				}

				fi, err = ffs.decorate(fi, joinedFilename)
				if err != nil {
					return nil, err
				}

				return fi.(FileMetaInfo), nil
			}
		}

		isSymlink := isSymlink(fi)
		if isSymlink {
			meta[metaKeyOriginalFilename] = filename
			var link string
			var err error
			link, fi, err = evalSymlinks(fs, filename)
			if err != nil {
				return nil, err
			}
			filename = link
			meta[metaKeyIsSymlink] = true
		}

		opener := func() (afero.File, error) {
			return ffs.open(filename)
		}

		fim := decorateFileInfo(fi, ffs, opener, filename, "", meta)

		for _, cb := range callbacks {
			cb(fim)
		}

		return fim, nil

	}

	ffs.decorate = decorator
	return ffs
}

func evalSymlinks(fs afero.Fs, filename string) (string, os.FileInfo, error) {
	link, err := filepath.EvalSymlinks(filename)
	if err != nil {
		return "", nil, err
	}

	fi, err := fs.Stat(link)
	if err != nil {
		return "", nil, err
	}

	return link, fi, nil
}

type baseFileDecoratorFs struct {
	afero.Fs
	decorate func(fi os.FileInfo, filename string) (os.FileInfo, error)
}

func (fs *baseFileDecoratorFs) Stat(name string) (os.FileInfo, error) {
	fi, err := fs.Fs.Stat(name)
	if err != nil {
		return nil, err
	}

	return fs.decorate(fi, name)

}

func (fs *baseFileDecoratorFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	var (
		fi  os.FileInfo
		err error
		ok  bool
	)

	if lstater, isLstater := fs.Fs.(afero.Lstater); isLstater {
		fi, ok, err = lstater.LstatIfPossible(name)
	} else {
		fi, err = fs.Fs.Stat(name)
	}

	if err != nil {
		return nil, false, err
	}

	fi, err = fs.decorate(fi, name)

	return fi, ok, err
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

func (l *baseFileDecoratorFile) Readdir(c int) (ofi []os.FileInfo, err error) {
	dirnames, err := l.File.Readdirnames(c)
	if err != nil {
		return nil, err
	}

	fisp := make([]os.FileInfo, 0, len(dirnames))

	for _, dirname := range dirnames {
		filename := dirname

		if l.Name() != "" && l.Name() != filepathSeparator {
			filename = filepath.Join(l.Name(), dirname)
		}

		// We need to resolve any symlink info.
		fi, _, err := lstatIfPossible(l.fs.Fs, filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		fi, err = l.fs.decorate(fi, filename)
		if err != nil {
			return nil, errors.Wrap(err, "decorate")
		}
		fisp = append(fisp, fi)
	}

	return fisp, err
}
