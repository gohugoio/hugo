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

var (
	filepathSeparator = string(filepath.Separator)
)

// NewRootMappingFs creates a new RootMappingFs on top of the provided with
// root mappings with some optional metadata about the root.
// Note that From represents a virtual root that maps to the actual filename in To.
func NewRootMappingFs(fs afero.Fs, rms ...RootMapping) (*RootMappingFs, error) {
	rootMapToReal := radix.New()
	var virtualRoots []RootMapping

	for _, rm := range rms {
		(&rm).clean()

		fromBase := files.ResolveComponentFolder(rm.From)

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
		if rm.Meta == nil {
			rm.Meta = make(FileMeta)
		}

		rm.Meta[metaKeySourceRoot] = rm.To
		rm.Meta[metaKeyBaseDir] = rm.ToBasedir
		rm.Meta[metaKeyMountRoot] = rm.path
		rm.Meta[metaKeyModule] = rm.Module

		meta := copyFileMeta(rm.Meta)

		if !fi.IsDir() {
			_, name := filepath.Split(rm.From)
			meta[metaKeyName] = name
		}

		rm.fi = NewFileMetaInfo(fi, meta)

		key := filepathSeparator + rm.From
		var mappings []RootMapping
		v, found := rootMapToReal.Get(key)
		if found {
			// There may be more than one language pointing to the same root.
			mappings = v.([]RootMapping)
		}
		mappings = append(mappings, rm)
		rootMapToReal.Insert(key, mappings)

		virtualRoots = append(virtualRoots, rm)
	}

	rootMapToReal.Insert(filepathSeparator, virtualRoots)

	rfs := &RootMappingFs{
		Fs:            fs,
		rootMapToReal: rootMapToReal,
	}

	return rfs, nil
}

func newRootMappingFsFromFromTo(
	baseDir string,
	fs afero.Fs,
	fromTo ...string,
) (*RootMappingFs, error) {

	rms := make([]RootMapping, len(fromTo)/2)
	for i, j := 0, 0; j < len(fromTo); i, j = i+1, j+2 {
		rms[i] = RootMapping{
			From:      fromTo[j],
			To:        fromTo[j+1],
			ToBasedir: baseDir,
		}
	}

	return NewRootMappingFs(fs, rms...)
}

// RootMapping describes a virtual file or directory mount.
type RootMapping struct {
	From      string   // The virtual mount.
	To        string   // The source directory or file.
	ToBasedir string   // The base of To. May be empty if an absolute path was provided.
	Module    string   // The module path/ID.
	Meta      FileMeta // File metadata (lang etc.)

	fi   FileMetaInfo
	path string // The virtual mount point, e.g. "blog".

}

type keyRootMappings struct {
	key   string
	roots []RootMapping
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

// A RootMappingFs maps several roots into one. Note that the root of this filesystem
// is directories only, and they will be returned in Readdir and Readdirnames
// in the order given.
type RootMappingFs struct {
	afero.Fs
	rootMapToReal *radix.Tree
}

func (fs *RootMappingFs) Dirs(base string) ([]FileMetaInfo, error) {
	base = filepathSeparator + fs.cleanName(base)
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

// Filter creates a copy of this filesystem with only mappings matching a filter.
func (fs RootMappingFs) Filter(f func(m RootMapping) bool) *RootMappingFs {
	rootMapToReal := radix.New()
	fs.rootMapToReal.Walk(func(b string, v interface{}) bool {
		rms := v.([]RootMapping)
		var nrms []RootMapping
		for _, rm := range rms {
			if f(rm) {
				nrms = append(nrms, rm)
			}
		}
		if len(nrms) != 0 {
			rootMapToReal.Insert(b, nrms)
		}
		return false
	})

	fs.rootMapToReal = rootMapToReal

	return &fs
}

// LstatIfPossible returns the os.FileInfo structure describing a given file.
func (fs *RootMappingFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	fis, err := fs.doLstat(name)
	if err != nil {
		return nil, false, err
	}
	return fis[0], false, nil
}

// Open opens the named file for reading.
func (fs *RootMappingFs) Open(name string) (afero.File, error) {
	fis, err := fs.doLstat(name)

	if err != nil {
		return nil, err
	}

	return fs.newUnionFile(fis...)
}

// Stat returns the os.FileInfo structure describing a given file.  If there is
// an error, it will be of type *os.PathError.
func (fs *RootMappingFs) Stat(name string) (os.FileInfo, error) {
	fi, _, err := fs.LstatIfPossible(name)
	return fi, err

}

func (fs *RootMappingFs) hasPrefix(prefix string) bool {
	hasPrefix := false
	fs.rootMapToReal.WalkPrefix(prefix, func(b string, v interface{}) bool {
		hasPrefix = true
		return true
	})

	return hasPrefix
}

func (fs *RootMappingFs) getRoot(key string) []RootMapping {
	v, found := fs.rootMapToReal.Get(key)
	if !found {
		return nil
	}

	return v.([]RootMapping)
}

func (fs *RootMappingFs) getRoots(key string) (string, []RootMapping) {
	s, v, found := fs.rootMapToReal.LongestPrefix(key)
	if !found || (s == filepathSeparator && key != filepathSeparator) {
		return "", nil
	}
	return s, v.([]RootMapping)

}

func (fs *RootMappingFs) debug() {
	fmt.Println("debug():")
	fs.rootMapToReal.Walk(func(s string, v interface{}) bool {
		fmt.Println("Key", s)
		return false
	})

}

func (fs *RootMappingFs) getRootsWithPrefix(prefix string) []RootMapping {
	var roots []RootMapping
	fs.rootMapToReal.WalkPrefix(prefix, func(b string, v interface{}) bool {
		roots = append(roots, v.([]RootMapping)...)
		return false
	})

	return roots
}

func (fs *RootMappingFs) getAncestors(prefix string) []keyRootMappings {
	var roots []keyRootMappings
	fs.rootMapToReal.WalkPath(prefix, func(s string, v interface{}) bool {
		if strings.HasPrefix(prefix, s+filepathSeparator) {
			roots = append(roots, keyRootMappings{
				key:   s,
				roots: v.([]RootMapping),
			})
		}
		return false
	})

	return roots
}

func (fs *RootMappingFs) newUnionFile(fis ...FileMetaInfo) (afero.File, error) {
	meta := fis[0].Meta()
	f, err := meta.Open()
	if err != nil {
		return nil, err
	}
	if len(fis) == 1 {
		return f, nil
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

func (fs *RootMappingFs) cleanName(name string) string {
	return strings.Trim(filepath.Clean(name), filepathSeparator)
}

func (fs *RootMappingFs) collectDirEntries(prefix string) ([]os.FileInfo, error) {
	prefix = filepathSeparator + fs.cleanName(prefix)

	var fis []os.FileInfo

	seen := make(map[string]bool) // Prevent duplicate directories
	level := strings.Count(prefix, filepathSeparator)

	collectDir := func(rm RootMapping, fi FileMetaInfo) error {
		f, err := fi.Meta().Open()
		if err != nil {
			return err
		}
		direntries, err := f.Readdir(-1)
		if err != nil {
			f.Close()
			return err
		}

		for _, fi := range direntries {
			meta := fi.(FileMetaInfo).Meta()
			mergeFileMeta(rm.Meta, meta)
			if fi.IsDir() {
				name := fi.Name()
				if seen[name] {
					continue
				}
				seen[name] = true
				opener := func() (afero.File, error) {
					return fs.Open(filepath.Join(rm.From, name))
				}
				fi = newDirNameOnlyFileInfo(name, meta, opener)
			}

			fis = append(fis, fi)
		}

		f.Close()

		return nil
	}

	// First add any real files/directories.
	rms := fs.getRoot(prefix)
	for _, rm := range rms {
		if err := collectDir(rm, rm.fi); err != nil {
			return nil, err
		}
	}

	// Next add any file mounts inside the given directory.
	prefixInside := prefix + filepathSeparator
	fs.rootMapToReal.WalkPrefix(prefixInside, func(s string, v interface{}) bool {

		if (strings.Count(s, filepathSeparator) - level) != 1 {
			// This directory is not part of the current, but we
			// need to include the first name part to make it
			// navigable.
			path := strings.TrimPrefix(s, prefixInside)
			parts := strings.Split(path, filepathSeparator)
			name := parts[0]

			if seen[name] {
				return false
			}
			seen[name] = true
			opener := func() (afero.File, error) {
				return fs.Open(path)
			}

			fi := newDirNameOnlyFileInfo(name, nil, opener)
			fis = append(fis, fi)

			return false
		}

		rms := v.([]RootMapping)
		for _, rm := range rms {
			if !rm.fi.IsDir() {
				// A single file mount
				fis = append(fis, rm.fi)
				continue
			}
			name := filepath.Base(rm.From)
			if seen[name] {
				continue
			}
			seen[name] = true

			opener := func() (afero.File, error) {
				return fs.Open(rm.From)
			}

			fi := newDirNameOnlyFileInfo(name, rm.Meta, opener)

			fis = append(fis, fi)

		}

		return false
	})

	// Finally add any ancestor dirs with files in this directory.
	ancestors := fs.getAncestors(prefix)
	for _, root := range ancestors {
		subdir := strings.TrimPrefix(prefix, root.key)
		for _, rm := range root.roots {
			if rm.fi.IsDir() {
				fi, err := rm.fi.Meta().JoinStat(subdir)
				if err == nil {
					if err := collectDir(rm, fi); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return fis, nil
}

func (fs *RootMappingFs) doLstat(name string) ([]FileMetaInfo, error) {
	name = fs.cleanName(name)
	key := filepathSeparator + name

	roots := fs.getRoot(key)

	if roots == nil {
		if fs.hasPrefix(key) {
			// We have directories mounted below this.
			// Make it look like a directory.
			return []FileMetaInfo{newDirNameOnlyFileInfo(name, nil, fs.virtualDirOpener(name))}, nil
		}

		// Find any real files or directories with this key.
		_, roots := fs.getRoots(key)
		if roots == nil {
			return nil, &os.PathError{Op: "LStat", Path: name, Err: os.ErrNotExist}
		}

		var err error
		var fis []FileMetaInfo

		for _, rm := range roots {
			var fi FileMetaInfo
			fi, _, err = fs.statRoot(rm, name)
			if err == nil {
				fis = append(fis, fi)
			}
		}

		if fis != nil {
			return fis, nil
		}

		if err == nil {
			err = &os.PathError{Op: "LStat", Path: name, Err: err}
		}

		return nil, err
	}

	fileCount := 0
	for _, root := range roots {
		if !root.fi.IsDir() {
			fileCount++
		}
		if fileCount > 1 {
			break
		}
	}

	if fileCount == 0 {
		// Dir only.
		return []FileMetaInfo{newDirNameOnlyFileInfo(name, roots[0].Meta, fs.virtualDirOpener(name))}, nil
	}

	if fileCount > 1 {
		// Not supported by this filesystem.
		return nil, errors.Errorf("found multiple files with name %q, use .Readdir or the source filesystem directly", name)

	}

	return []FileMetaInfo{roots[0].fi}, nil

}

func (fs *RootMappingFs) statRoot(root RootMapping, name string) (FileMetaInfo, bool, error) {
	filename := root.filename(name)

	fi, b, err := lstatIfPossible(fs.Fs, filename)
	if err != nil {
		return nil, b, err
	}

	var opener func() (afero.File, error)
	if fi.IsDir() {
		// Make sure metadata gets applied in Readdir.
		opener = fs.realDirOpener(filename, root.Meta)
	} else {
		// Opens the real file directly.
		opener = func() (afero.File, error) {
			return fs.Fs.Open(filename)
		}
	}

	return decorateFileInfo(fi, fs.Fs, opener, "", "", root.Meta), b, nil

}

func (fs *RootMappingFs) virtualDirOpener(name string) func() (afero.File, error) {
	return func() (afero.File, error) { return &rootMappingFile{name: name, fs: fs}, nil }
}

func (fs *RootMappingFs) realDirOpener(name string, meta FileMeta) func() (afero.File, error) {
	return func() (afero.File, error) {
		f, err := fs.Fs.Open(name)
		if err != nil {
			return nil, err
		}
		return &rootMappingFile{name: name, meta: meta, fs: fs, File: f}, nil
	}
}

type rootMappingFile struct {
	afero.File
	fs   *RootMappingFs
	name string
	meta FileMeta
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
	if f.File != nil {
		fis, err := f.File.Readdir(count)
		if err != nil {
			return nil, err
		}

		for i, fi := range fis {
			fis[i] = decorateFileInfo(fi, f.fs, nil, "", "", f.meta)
		}
		return fis, nil
	}
	return f.fs.collectDirEntries(f.name)
}

func (f *rootMappingFile) Readdirnames(count int) ([]string, error) {
	dirs, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}
	return fileInfosToNames(dirs), nil
}
