// Copyright 2016 The Hugo Authors. All rights reserved.
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

package source

import (
	"path/filepath"
	"sync"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/hugofs"
)

// Filesystem represents a source filesystem.
type Filesystem struct {
	files        []File
	filesInit    sync.Once
	filesInitErr error

	Base string

	fi hugofs.FileMetaInfo

	SourceSpec
}

// NewFilesystem returns a new filesytem for a given source spec.
func (sp SourceSpec) NewFilesystem(base string) *Filesystem {
	return &Filesystem{SourceSpec: sp, Base: base}
}

func (sp SourceSpec) NewFilesystemFromFileMetaInfo(fi hugofs.FileMetaInfo) *Filesystem {
	return &Filesystem{SourceSpec: sp, fi: fi}
}

// Files returns a slice of readable files.
func (f *Filesystem) Files() ([]File, error) {
	f.filesInit.Do(func() {
		err := f.captureFiles()
		if err != nil {
			f.filesInitErr = errors.Wrap(err, "capture files")
		}
	})
	return f.files, f.filesInitErr
}

// add populates a file in the Filesystem.files
func (f *Filesystem) add(name string, fi hugofs.FileMetaInfo) (err error) {
	var file File

	file, err = f.SourceSpec.NewFileInfo(fi)
	if err != nil {
		return err
	}

	f.files = append(f.files, file)

	return err
}

func (f *Filesystem) captureFiles() error {
	walker := func(path string, fi hugofs.FileMetaInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		meta := fi.Meta()
		filename := meta.Filename()

		b, err := f.shouldRead(filename, fi)
		if err != nil {
			return err
		}

		if b {
			err = f.add(filename, fi)
		}

		return err
	}

	w := hugofs.NewWalkway(hugofs.WalkwayConfig{
		Fs:     f.SourceFs,
		Info:   f.fi,
		Root:   f.Base,
		WalkFn: walker,
	})

	return w.Walk()

}

func (f *Filesystem) shouldRead(filename string, fi hugofs.FileMetaInfo) (bool, error) {

	ignore := f.SourceSpec.IgnoreFile(fi.Meta().Filename())

	if fi.IsDir() {
		if ignore {
			return false, filepath.SkipDir
		}
		return false, nil
	}

	if ignore {
		return false, nil
	}

	return true, nil
}
