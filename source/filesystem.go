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
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gohugoio/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
	"golang.org/x/text/unicode/norm"
)

type Filesystem struct {
	files     []ReadableFile
	filesInit sync.Once

	Base string

	SourceSpec
}

type Input interface {
	Files() []ReadableFile
}

func (sp SourceSpec) NewFilesystem(base string) *Filesystem {
	return &Filesystem{SourceSpec: sp, Base: base}
}

func (f *Filesystem) Files() []ReadableFile {
	f.filesInit.Do(func() {
		f.captureFiles()
	})
	return f.files
}

// add populates a file in the Filesystem.files
func (f *Filesystem) add(name string, fi os.FileInfo) (err error) {
	var file ReadableFile

	if runtime.GOOS == "darwin" {
		// When a file system is HFS+, its filepath is in NFD form.
		name = norm.NFC.String(name)
	}

	file = f.SourceSpec.NewFileInfo(f.Base, name, false, fi)
	f.files = append(f.files, file)

	return err
}

func (f *Filesystem) captureFiles() {
	walker := func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		b, err := f.shouldRead(filePath, fi)
		if err != nil {
			return err
		}
		if b {
			f.add(filePath, fi)
		}
		return err
	}

	if f.SourceFs == nil {
		panic("Must have a fs")
	}
	err := helpers.SymbolicWalk(f.SourceFs, f.Base, walker)

	if err != nil {
		jww.ERROR.Println(err)
	}

}

func (f *Filesystem) shouldRead(filename string, fi os.FileInfo) (bool, error) {
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, err := filepath.EvalSymlinks(filename)
		if err != nil {
			jww.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", filename, err)
			return false, nil
		}
		linkfi, err := f.SourceFs.Stat(link)
		if err != nil {
			jww.ERROR.Printf("Cannot stat '%s', error was: %s", link, err)
			return false, nil
		}

		if !linkfi.Mode().IsRegular() {
			jww.ERROR.Printf("Symbolic links for directories not supported, skipping '%s'", filename)
		}
		return false, nil
	}

	ignore := f.SourceSpec.IgnoreFile(filename)

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
