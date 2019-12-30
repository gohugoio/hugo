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
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/pkg/errors"

	radix "github.com/armon/go-radix"
	"github.com/spf13/afero"
)

var filepathSeparator = string(filepath.Separator)

// NewRootMappingFs creates a new RootMappingFs on top of the provided with
// of root mappings with some optional metadata about the root.
// Note that From represents a virtual root that maps to the actual filename in To.
func NewRootMappingFs(fs afero.Fs, rms ...RootMapping) (*RootMappingFs, error) {
	rootMapToReal := radix.New()

	for i, rm := range rms {
		(&rm).clean()

		fromBase := files.ResolveComponentFolder(rm.From)
		if fromBase == "" {
			panic("unrecognised component folder in" + rm.From)
		}

		if len(rm.To) < 2 {
			panic(fmt.Sprintf("invalid root mapping; from/to: %s/%s", rm.From, rm.To))
		}

		fi, err := fs.Stat(rm.To)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		// Extract "blog" from "content/blog"
		rm.path = strings.TrimPrefix(strings.TrimPrefix(rm.From, fromBase), filepathSeparator)
		if rm.Meta != nil {
			rm.Meta[metaKeyIsFileMount] = !fi.IsDir()
			rm.Meta[metaKeyMountRoot] = rm.path
			if rm.ToBasedir != "" {
				pathFile := strings.TrimPrefix(strings.TrimPrefix(rm.To, rm.ToBasedir), filepathSeparator)
				rm.Meta[metaKeyPathFile] = pathFile
			}
		}

		meta := copyFileMeta(rm.Meta)

		if !fi.IsDir() {
			_, name := filepath.Split(rm.From)
			meta[metaKeyName] = name
		}

		rm.fi = NewFileMetaInfo(fi, meta)

		key := rm.rootKey()
		var mappings []RootMapping
		v, found := rootMapToReal.Get(key)
		if found {
			// There may be more than one language pointing to the same root.
			mappings = v.([]RootMapping)
		}
		mappings = append(mappings, rm)
		rootMapToReal.Insert(key, mappings)

		rms[i] = rm
	}

	rfs := &RootMappingFs{Fs: fs,
		virtualRoots:  rms,
		rootMapToReal: rootMapToReal}

	return rfs, nil
}

// NewRootMappingFsFromFromTo is a convenicence variant of NewRootMappingFs taking
// From and To as string pairs.
func NewRootMappingFsFromFromTo(fs afero.Fs, fromTo ...string) (*RootMappingFs, error) {
	rms := make([]RootMapping, len(fromTo)/2)
	for i, j := 0, 0; j < len(fromTo); i, j = i+1, j+2 {
		rms[i] = RootMapping{
			From: fromTo[j],
			To:   fromTo[j+1],
		}
	}

	return NewRootMappingFs(fs, rms...)
}

type RootMapping struct {
	From      string   // The virtual mount.
	To        string   // The source directory or file.
	ToBasedir string   // The base of To. May be empty if an absolute path was provided.
	Meta      FileMeta // File metadata (lang etc.)

	fi   FileMetaInfo
	path string // The virtual mount point, e.g. "blog".

}

func (rm *RootMapping) clean() {
	rm.From = strings.Trim(filepath.Clean(rm.From), filepathSeparator)
	rm.To = filepath.Clean(rm.To)
}

func (r RootMapping) filename(name string) string {
	if name == "" {
		return r.To
	}
	return filepath.Join(r.To, strings.TrimPrefix(name, r.From))
}

func (r RootMapping) rootKey() string {
	return r.From
}

// A RootMappingFs maps several roots into one. Note that the root of this filesystem
// is directories only, and they will be returned in Readdir and Readdirnames
// in the order given.
type RootMappingFs struct {
	afero.Fs
	rootMapToReal *radix.Tree
	virtualRoots  []RootMapping
	filter        func(r RootMapping) bool
}

func (fs *RootMappingFs) Dirs(base string) ([]FileMetaInfo, error) {
	roots := fs.getRootsWithPrefix(base)

	if roots == nil {
		return nil, nil
	}

	fss := make([]FileMetaInfo, len(roots))
	for i, r := range roots {
		bfs := afero.NewBasePathFs(fs.Fs, r.To)
		bfs = decoratePath(bfs, func(name string) string {
			p := strings.TrimPrefix(name, r.To)
			if r.path != "" {
				// Make sure it's mounted to a any sub path, e.g. blog
				p = filepath.Join(r.path, p)
			}
			p = strings.TrimLeft(p, filepathSeparator)
			return p
		})
		fs := decorateDirs(bfs, r.Meta)
		fi, err := fs.Stat("")
		if err != nil {
			return nil, errors.Wrap(err, "RootMappingFs.Dirs")
		}

		if !fi.IsDir() {
			mergeFileMeta(r.Meta, fi.(FileMetaInfo).Meta())
		}

		fss[i] = fi.(FileMetaInfo)
	}

	return fss, nil
}

// LstatIfPossible returns the os.FileInfo structure describing a given file.
func (fs *RootMappingFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	fis, _, b, err := fs.doLstat(name, false)
	if err != nil {
		return nil, b, err
	}
	return fis[0], b, nil
}

func (fs *RootMappingFs) virtualDirOpener(name string, isRoot bool) func() (afero.File, error) {
	return func() (afero.File, error) { return &rootMappingFile{name: name, isRoot: isRoot, fs: fs}, nil }
}

func (fs *RootMappingFs) doLstat(name string, allowMultiple bool) ([]FileMetaInfo, []FileMetaInfo, bool, error) {
	if fs.isRoot(name) {
		return []FileMetaInfo{newDirNameOnlyFileInfo(name, true, fs.virtualDirOpener(name, true))}, nil, false, nil
	}

	roots := fs.getRoots(name)
	rootsWithPrefix := fs.getRootsWithPrefix(name)
	hasRootMappingsBelow := len(rootsWithPrefix) != 0

	if len(roots) == 0 {
		if hasRootMappingsBelow {
			// No exact matches, but we have root mappings below name,
			// let's make it look like a directory.
			return []FileMetaInfo{newDirNameOnlyFileInfo(name, true, fs.virtualDirOpener(name, false))}, nil, false, nil
		}

		return nil, nil, false, os.ErrNotExist
	}

	// We may have a mapping for both static and static/subdir.
	// These will not show in any Readdir so append them
	// manually.
	rootsInDir := fs.filterRootsBelow(rootsWithPrefix, name)

	var (
		fis  []FileMetaInfo
		dirs []FileMetaInfo
		b    bool
		root RootMapping
		err  error
	)

	for _, root = range roots {
		var fi os.FileInfo
		fi, b, err = fs.statRoot(root, name)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, nil, false, err
		}
		fim := fi.(FileMetaInfo)

		fis = append(fis, fim)
	}

	for _, root = range rootsInDir {

		fi, _, err := fs.statRoot(root, "")
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, nil, false, err
		}
		fim := fi.(FileMetaInfo)
		dirs = append(dirs, fim)
	}

	if len(fis) == 0 && len(dirs) == 0 {
		return nil, nil, false, os.ErrNotExist
	}

	if allowMultiple || len(fis) == 1 {
		return fis, dirs, b, nil
	}

	if len(fis) == 0 {
		return nil, nil, false, os.ErrNotExist
	}

	// Open it in this composite filesystem.
	opener := func() (afero.File, error) {
		return fs.Open(name)
	}

	return []FileMetaInfo{decorateFileInfo(fis[0], fs, opener, "", "", root.Meta)}, nil, b, nil

}

// Open opens the namedrootMappingFile file for reading.
func (fs *RootMappingFs) Open(name string) (afero.File, error) {
	if fs.isRoot(name) {
		return &rootMappingFile{name: name, fs: fs, isRoot: true}, nil
	}

	fis, dirs, _, err := fs.doLstat(name, true)
	if err != nil {
		return nil, err
	}

	if len(fis) == 1 {
		fi := fis[0]
		meta := fi.(FileMetaInfo).Meta()
		f, err := meta.Open()
		if err != nil {
			return nil, err
		}

		f = &rootMappingFile{File: f, fs: fs, name: name, meta: meta}

		if len(dirs) > 0 {
			return &readDirDirsAppender{File: f, dirs: dirs}, nil
		}

		return f, nil
	}

	f, err := fs.newUnionFile(fis...)
	if err != nil {
		return nil, err
	}

	if len(dirs) > 0 {
		return &readDirDirsAppender{File: f, dirs: dirs}, nil
	}

	return f, nil

}

// Stat returns the os.FileInfo structure describing a given file.  If there is
// an error, it will be of type *os.PathError.
func (fs *RootMappingFs) Stat(name string) (os.FileInfo, error) {
	fi, _, err := fs.LstatIfPossible(name)
	return fi, err

}

// Filter creates a copy of this filesystem with the applied filter.
func (fs RootMappingFs) Filter(f func(m RootMapping) bool) *RootMappingFs {
	fs.filter = f
	return &fs
}

func (fs *RootMappingFs) isRoot(name string) bool {
	return name == "" || name == filepathSeparator

}

func (fs *RootMappingFs) getRoots(name string) []RootMapping {
	name = filepath.Clean(name)
	_, v, found := fs.rootMapToReal.LongestPrefix(name)
	if !found {
		return nil
	}

	rm := v.([]RootMapping)

	return fs.applyFilterToRoots(rm)
}

func (fs *RootMappingFs) applyFilterToRoots(rm []RootMapping) []RootMapping {
	if fs.filter == nil {
		return rm
	}

	var filtered []RootMapping
	for _, m := range rm {
		if fs.filter(m) {
			filtered = append(filtered, m)
		}
	}

	return filtered
}

func (fs *RootMappingFs) getRootsWithPrefix(prefix string) []RootMapping {
	if fs.isRoot(prefix) {
		return fs.virtualRoots
	}
	prefix = filepath.Clean(prefix)
	var roots []RootMapping

	fs.rootMapToReal.WalkPrefix(prefix, func(b string, v interface{}) bool {
		roots = append(roots, v.([]RootMapping)...)
		return false
	})

	return fs.applyFilterToRoots(roots)
}

// Filter out the mappings inside the name directory.
func (fs *RootMappingFs) filterRootsBelow(roots []RootMapping, name string) []RootMapping {
	if len(roots) == 0 {
		return nil
	}

	sepCount := strings.Count(name, filepathSeparator)
	var filtered []RootMapping
	for _, x := range roots {
		if name == x.From {
			continue
		}

		if strings.Count(x.From, filepathSeparator)-sepCount != 1 {
			continue
		}

		filtered = append(filtered, x)

	}
	return filtered
}

func (fs *RootMappingFs) newUnionFile(fis ...FileMetaInfo) (afero.File, error) {
	meta := fis[0].Meta()
	f, err := meta.Open()
	if err != nil {
		return nil, err
	}
	rf := &rootMappingFile{File: f, fs: fs, name: meta.Name(), meta: meta}
	if len(fis) == 1 {
		return rf, err
	}

	next, err := fs.newUnionFile(fis[1:]...)
	if err != nil {
		return nil, err
	}

	uf := &afero.UnionFile{Base: rf, Layer: next}

	uf.Merger = func(lofi, bofi []os.FileInfo) ([]os.FileInfo, error) {
		// Ignore duplicate directory entries
		seen := make(map[string]bool)
		var result []os.FileInfo

		for _, fis := range [][]os.FileInfo{bofi, lofi} {
			for _, fi := range fis {

				if fi.IsDir() && seen[fi.Name()] {
					continue
				}

				if fi.IsDir() {
					seen[fi.Name()] = true
				}

				result = append(result, fi)
			}
		}

		return result, nil
	}

	return uf, nil

}

func (fs *RootMappingFs) statRoot(root RootMapping, name string) (os.FileInfo, bool, error) {
	filename := root.filename(name)

	var b bool
	var fi os.FileInfo
	var err error

	if ls, ok := fs.Fs.(afero.Lstater); ok {
		fi, b, err = ls.LstatIfPossible(filename)
		if err != nil {
			return nil, b, err
		}

	} else {
		fi, err = fs.Fs.Stat(filename)
		if err != nil {
			return nil, b, err
		}
	}

	// Opens the real directory/file.
	opener := func() (afero.File, error) {
		return fs.Fs.Open(filename)
	}

	if fi.IsDir() {
		if name == "" {
			name = root.From
		}
		_, name = filepath.Split(name)
		fi = newDirNameOnlyFileInfo(name, false, opener)
	}

	return decorateFileInfo(fi, fs.Fs, opener, "", "", root.Meta), b, nil

}

type rootMappingFile struct {
	afero.File
	fs     *RootMappingFs
	name   string
	meta   FileMeta
	isRoot bool
}

type readDirDirsAppender struct {
	afero.File
	dirs []FileMetaInfo
}

func (f *readDirDirsAppender) Readdir(count int) ([]os.FileInfo, error) {
	fis, err := f.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	for _, dir := range f.dirs {
		fis = append(fis, dir)
	}
	return fis, nil

}

func (f *readDirDirsAppender) Readdirnames(count int) ([]string, error) {
	fis, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}
	return fileInfosToNames(fis), nil
}

func (f *rootMappingFile) Close() error {
	if f.File == nil {
		return nil
	}
	return f.File.Close()
}

func (f *rootMappingFile) Name() string {
	return f.name
}

func (f *rootMappingFile) Readdir(count int) ([]os.FileInfo, error) {
	if f.File == nil {
		filesn := make([]os.FileInfo, 0)
		roots := f.fs.getRootsWithPrefix(f.name)
		seen := make(map[string]bool) // Do not return duplicate directories

		j := 0
		for _, rm := range roots {
			if count != -1 && j >= count {
				break
			}

			if !rm.fi.IsDir() {
				// A single file mount
				filesn = append(filesn, rm.fi)
				continue
			}

			from := rm.From
			name := from
			if !f.isRoot {
				_, name = filepath.Split(from)
			}

			if seen[name] {
				continue
			}
			seen[name] = true

			opener := func() (afero.File, error) {
				return f.fs.Open(from)
			}

			j++

			fi := newDirNameOnlyFileInfo(name, false, opener)

			if rm.Meta != nil {
				mergeFileMeta(rm.Meta, fi.Meta())
			}

			filesn = append(filesn, fi)
		}
		return filesn, nil
	}

	if f.File == nil {
		panic(fmt.Sprintf("no File for %q", f.name))
	}

	fis, err := f.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	for i, fi := range fis {
		fis[i] = decorateFileInfo(fi, f.fs, nil, "", "", f.meta)
	}

	return fis, nil
}

func (f *rootMappingFile) Readdirnames(count int) ([]string, error) {
	dirs, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}
	return fileInfosToNames(dirs), nil
}
