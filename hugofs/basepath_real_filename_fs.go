// Copyright 2018 The Hugo Authors. All rights reserved.
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

	"github.com/spf13/afero"
)

// RealFilenameInfo is a thin wrapper around os.FileInfo adding the real filename.
type RealFilenameInfo interface {
	os.FileInfo

	// This is the real filename to the file in the underlying filesystem.
	RealFilename() string
}

type realFilenameInfo struct {
	os.FileInfo
	realFilename string
}

func (f *realFilenameInfo) RealFilename() string {
	return f.realFilename
}

// NewBasePathRealFilenameFs returns a new BasePathRealFilenameFs instance
// using base.
func NewBasePathRealFilenameFs(base *afero.BasePathFs) *BasePathRealFilenameFs {
	return &BasePathRealFilenameFs{BasePathFs: base}
}

// BasePathRealFilenameFs is a thin wrapper around afero.BasePathFs that
// provides the real filename in Stat and LstatIfPossible.
type BasePathRealFilenameFs struct {
	*afero.BasePathFs
}

// Stat returns the os.FileInfo structure describing a given file.  If there is
// an error, it will be of type *os.PathError.
func (b *BasePathRealFilenameFs) Stat(name string) (os.FileInfo, error) {
	fi, err := b.BasePathFs.Stat(name)
	if err != nil {
		return nil, err
	}

	if _, ok := fi.(RealFilenameInfo); ok {
		return fi, nil
	}

	filename, err := b.RealPath(name)
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}

	return &realFilenameInfo{FileInfo: fi, realFilename: filename}, nil
}

// LstatIfPossible returns the os.FileInfo structure describing a given file.
// It attempts to use Lstat if supported or defers to the os.  In addition to
// the FileInfo, a boolean is returned telling whether Lstat was called.
func (b *BasePathRealFilenameFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {

	fi, ok, err := b.BasePathFs.LstatIfPossible(name)
	if err != nil {
		return nil, false, err
	}

	if _, ok := fi.(RealFilenameInfo); ok {
		return fi, ok, nil
	}

	filename, err := b.RealPath(name)
	if err != nil {
		return nil, false, &os.PathError{Op: "lstat", Path: name, Err: err}
	}

	return &realFilenameInfo{FileInfo: fi, realFilename: filename}, ok, nil
}
