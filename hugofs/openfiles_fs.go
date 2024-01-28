// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"io/fs"
	"os"
	"sync"

	"github.com/spf13/afero"
)

var _ FilesystemUnwrapper = (*OpenFilesFs)(nil)

// OpenFilesFs is a wrapper around afero.Fs that keeps track of open files.
type OpenFilesFs struct {
	afero.Fs

	mu        sync.Mutex
	openFiles map[string]int
}

func (fs *OpenFilesFs) UnwrapFilesystem() afero.Fs {
	return fs.Fs
}

func (fs *OpenFilesFs) Create(name string) (afero.File, error) {
	f, err := fs.Fs.Create(name)
	if err != nil {
		return nil, err
	}
	return fs.trackAndWrapFile(f), nil
}

func (fs *OpenFilesFs) Open(name string) (afero.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	return fs.trackAndWrapFile(f), nil
}

func (fs *OpenFilesFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return fs.trackAndWrapFile(f), nil
}

func (fs *OpenFilesFs) trackAndWrapFile(f afero.File) afero.File {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if fs.openFiles == nil {
		fs.openFiles = make(map[string]int)
	}

	fs.openFiles[f.Name()]++

	return &openFilesFsFile{fs: fs, File: f}
}

type openFilesFsFile struct {
	fs *OpenFilesFs
	afero.File
}

func (f *openFilesFsFile) ReadDir(count int) ([]fs.DirEntry, error) {
	return f.File.(fs.ReadDirFile).ReadDir(count)
}

func (f *openFilesFsFile) Close() (err error) {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()

	err = f.File.Close()

	if f.fs.openFiles == nil {
		return
	}

	name := f.Name()

	f.fs.openFiles[name]--

	if f.fs.openFiles[name] <= 0 {
		delete(f.fs.openFiles, name)
	}

	return
}

func (fs *OpenFilesFs) OpenFiles() map[string]int {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return fs.openFiles
}
