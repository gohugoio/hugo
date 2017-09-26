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
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
	"golang.org/x/text/unicode/norm"
)

type Input interface {
	Files() []*File
}

type Filesystem struct {
	files      []*File
	Base       string
	AvoidPaths []string

	SourceSpec
}

func (sp SourceSpec) NewFilesystem(base string, avoidPaths ...string) *Filesystem {
	return &Filesystem{SourceSpec: sp, Base: base, AvoidPaths: avoidPaths}
}

func (f *Filesystem) FilesByExts(exts ...string) []*File {
	var newFiles []*File

	if len(exts) == 0 {
		return f.Files()
	}

	for _, x := range f.Files() {
		for _, e := range exts {
			if x.Ext() == strings.TrimPrefix(e, ".") {
				newFiles = append(newFiles, x)
			}
		}
	}
	return newFiles
}

func (f *Filesystem) Files() []*File {
	if len(f.files) < 1 {
		f.captureFiles()
	}
	return f.files
}

// add populates a file in the Filesystem.files
func (f *Filesystem) add(name string, reader io.Reader) (err error) {
	var file *File

	if runtime.GOOS == "darwin" {
		// When a file system is HFS+, its filepath is in NFD form.
		name = norm.NFC.String(name)
	}

	file, err = f.SourceSpec.NewFileFromAbs(f.Base, name, reader)

	if err == nil {
		f.files = append(f.files, file)
	}
	return err
}

func (f *Filesystem) captureFiles() {
	walker := func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		b, err := f.ShouldRead(filePath, fi)
		if err != nil {
			return err
		}
		if b {
			rd, err := NewLazyFileReader(f.Fs.Source, filePath)
			if err != nil {
				return err
			}
			f.add(filePath, rd)
		}
		return err
	}

	if f.Fs == nil {
		panic("Must have a fs")
	}
	err := helpers.SymbolicWalk(f.Fs.Source, f.Base, walker)

	if err != nil {
		jww.ERROR.Println(err)
		if err == helpers.ErrWalkRootTooShort {
			panic("The root path is too short. If this is a test, make sure to init the content paths.")
		}
	}

}

func (f *Filesystem) ShouldRead(filePath string, fi os.FileInfo) (bool, error) {
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, err := filepath.EvalSymlinks(filePath)
		if err != nil {
			jww.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", filePath, err)
			return false, nil
		}
		linkfi, err := f.Fs.Source.Stat(link)
		if err != nil {
			jww.ERROR.Printf("Cannot stat '%s', error was: %s", link, err)
			return false, nil
		}
		if !linkfi.Mode().IsRegular() {
			jww.ERROR.Printf("Symbolic links for directories not supported, skipping '%s'", filePath)
		}
		return false, nil
	}

	if fi.IsDir() {
		if f.avoid(filePath) || f.isNonProcessablePath(filePath) {
			return false, filepath.SkipDir
		}
		return false, nil
	}

	if f.isNonProcessablePath(filePath) {
		return false, nil
	}
	return true, nil
}

func (f *Filesystem) avoid(filePath string) bool {
	for _, avoid := range f.AvoidPaths {
		if avoid == filePath {
			return true
		}
	}
	return false
}

func (sp SourceSpec) isNonProcessablePath(filePath string) bool {
	base := filepath.Base(filePath)
	if strings.HasPrefix(base, ".") ||
		strings.HasPrefix(base, "#") ||
		strings.HasSuffix(base, "~") {
		return true
	}
	ignoreFiles := cast.ToStringSlice(sp.Cfg.Get("ignoreFiles"))
	if len(ignoreFiles) > 0 {
		for _, ignorePattern := range ignoreFiles {
			match, err := regexp.MatchString(ignorePattern, filePath)
			if err != nil {
				helpers.DistinctErrorLog.Printf("Invalid regexp '%s' in ignoreFiles: %s", ignorePattern, err)
				return false
			} else if match {
				return true
			}
		}
	}
	return false
}
