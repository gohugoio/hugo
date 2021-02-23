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

package hugio

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/spf13/afero"
)

// CopyFile copies a file.
func CopyFile(fs afero.Fs, from, to string) error {
	sf, err := fs.Open(from)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := fs.Create(to)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}
	si, err := fs.Stat(from)
	if err != nil {
		err = fs.Chmod(to, si.Mode())

		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDir copies a directory.
func CopyDir(fs afero.Fs, from, to string, shouldCopy func(filename string) bool) error {
	fi, err := os.Stat(from)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.Errorf("%q is not a directory", from)
	}

	err = fs.MkdirAll(to, 0777) // before umask
	if err != nil {
		return err
	}

	entries, _ := ioutil.ReadDir(from)
	for _, entry := range entries {
		fromFilename := filepath.Join(from, entry.Name())
		toFilename := filepath.Join(to, entry.Name())
		if entry.IsDir() {
			if shouldCopy != nil && !shouldCopy(fromFilename) {
				continue
			}
			if err := CopyDir(fs, fromFilename, toFilename, shouldCopy); err != nil {
				return err
			}
		} else {
			if err := CopyFile(fs, fromFilename, toFilename); err != nil {
				return err
			}
		}

	}

	return nil
}
