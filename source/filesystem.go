// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
)

type Input interface {
	Files() []*File
}

type Filesystem struct {
	files      []*File
	Base       string
	AvoidPaths []string
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

	file, err = NewFileFromAbs(f.Base, name, reader)

	if err == nil {
		f.files = append(f.files, file)
	}
	return err
}

func (f *Filesystem) getRelativePath(name string) (final string, err error) {
	return helpers.GetRelativePath(name, f.Base)
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
			rd, err := NewLazyFileReader(filePath)
			if err != nil {
				return err
			}
			f.add(filePath, rd)
		}
		return err
	}

	filepath.Walk(f.Base, walker)
}

func (f *Filesystem) shouldRead(filePath string, fi os.FileInfo) (bool, error) {
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		link, err := filepath.EvalSymlinks(filePath)
		if err != nil {
			jww.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", filePath, err)
			return false, nil
		}
		linkfi, err := os.Stat(link)
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
		if f.avoid(filePath) || isNonProcessablePath(filePath) {
			return false, filepath.SkipDir
		}
		return false, nil
	}

	if isNonProcessablePath(filePath) {
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

func isNonProcessablePath(filePath string) bool {
	base := filepath.Base(filePath)
	if base[0] == '.' {
		return true
	}

	if base[0] == '#' {
		return true
	}

	if base[len(base)-1] == '~' {
		return true
	}

	ignoreFiles := viper.GetStringSlice("IgnoreFiles")
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
