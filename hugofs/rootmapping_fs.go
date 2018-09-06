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
	"path/filepath"
	"strings"
	"time"

	radix "github.com/hashicorp/go-immutable-radix"
	"github.com/spf13/afero"
)

var filepathSeparator = string(filepath.Separator)

// A RootMappingFs maps several roots into one. Note that the root of this filesystem
// is directories only, and they will be returned in Readdir and Readdirnames
// in the order given.
type RootMappingFs struct {
	afero.Fs
	rootMapToReal *radix.Node
	virtualRoots  []string
}

type rootMappingFile struct {
	afero.File
	fs   *RootMappingFs
	name string
}

type rootMappingFileInfo struct {
	name string
}

func (fi *rootMappingFileInfo) Name() string {
	return fi.name
}

func (fi *rootMappingFileInfo) Size() int64 {
	panic("not implemented")
}

func (fi *rootMappingFileInfo) Mode() os.FileMode {
	return os.ModeDir
}

func (fi *rootMappingFileInfo) ModTime() time.Time {
	panic("not implemented")
}

func (fi *rootMappingFileInfo) IsDir() bool {
	return true
}

func (fi *rootMappingFileInfo) Sys() interface{} {
	return nil
}

func newRootMappingDirFileInfo(name string) *rootMappingFileInfo {
	return &rootMappingFileInfo{name: name}
}

// NewRootMappingFs creates a new RootMappingFs on top of the provided with
// a list of from, to string pairs of root mappings.
// Note that 'from' represents a virtual root that maps to the actual filename in 'to'.
func NewRootMappingFs(fs afero.Fs, fromTo ...string) (*RootMappingFs, error) {
	rootMapToReal := radix.New().Txn()
	var virtualRoots []string

	for i := 0; i < len(fromTo); i += 2 {
		vr := filepath.Clean(fromTo[i])
		rr := filepath.Clean(fromTo[i+1])

		// We need to preserve the original order for Readdir
		virtualRoots = append(virtualRoots, vr)

		rootMapToReal.Insert([]byte(vr), rr)
	}

	return &RootMappingFs{Fs: fs,
		virtualRoots:  virtualRoots,
		rootMapToReal: rootMapToReal.Commit().Root()}, nil
}

// Stat returns the os.FileInfo structure describing a given file.  If there is
// an error, it will be of type *os.PathError.
func (fs *RootMappingFs) Stat(name string) (os.FileInfo, error) {
	if fs.isRoot(name) {
		return newRootMappingDirFileInfo(name), nil
	}
	realName := fs.realName(name)
	return fs.Fs.Stat(realName)
}

func (fs *RootMappingFs) isRoot(name string) bool {
	return name == "" || name == filepathSeparator

}

// Open opens the named file for reading.
func (fs *RootMappingFs) Open(name string) (afero.File, error) {
	if fs.isRoot(name) {
		return &rootMappingFile{name: name, fs: fs}, nil
	}
	realName := fs.realName(name)
	f, err := fs.Fs.Open(realName)
	if err != nil {
		return nil, err
	}
	return &rootMappingFile{File: f, name: name, fs: fs}, nil
}

// LstatIfPossible returns the os.FileInfo structure describing a given file.
// It attempts to use Lstat if supported or defers to the os.  In addition to
// the FileInfo, a boolean is returned telling whether Lstat was called.
func (fs *RootMappingFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	if fs.isRoot(name) {
		return newRootMappingDirFileInfo(name), false, nil
	}
	name = fs.realName(name)
	if ls, ok := fs.Fs.(afero.Lstater); ok {
		return ls.LstatIfPossible(name)
	}
	fi, err := fs.Stat(name)
	return fi, false, err
}

func (fs *RootMappingFs) realName(name string) string {
	key, val, found := fs.rootMapToReal.LongestPrefix([]byte(filepath.Clean(name)))
	if !found {
		return name
	}
	keystr := string(key)

	return filepath.Join(val.(string), strings.TrimPrefix(name, keystr))
}

func (f *rootMappingFile) Readdir(count int) ([]os.FileInfo, error) {
	if f.File == nil {
		dirsn := make([]os.FileInfo, 0)
		for i := 0; i < len(f.fs.virtualRoots); i++ {
			if count != -1 && i >= count {
				break
			}
			dirsn = append(dirsn, newRootMappingDirFileInfo(f.fs.virtualRoots[i]))
		}
		return dirsn, nil
	}
	return f.File.Readdir(count)

}

func (f *rootMappingFile) Readdirnames(count int) ([]string, error) {
	dirs, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}
	dirss := make([]string, len(dirs))
	for i, d := range dirs {
		dirss[i] = d.Name()
	}
	return dirss, nil
}

func (f *rootMappingFile) Name() string {
	return f.name
}

func (f *rootMappingFile) Close() error {
	if f.File == nil {
		return nil
	}
	return f.File.Close()
}
