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
	"errors"
	"fmt"
	iofs "io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/bep/overlayfs"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugofs/glob"

	radix "github.com/armon/go-radix"
	"github.com/spf13/afero"
)

var filepathSeparator = string(filepath.Separator)

var _ ReverseLookupProvder = (*RootMappingFs)(nil)

// NewRootMappingFs creates a new RootMappingFs on top of the provided with
// root mappings with some optional metadata about the root.
// Note that From represents a virtual root that maps to the actual filename in To.
func NewRootMappingFs(fs afero.Fs, rms ...RootMapping) (*RootMappingFs, error) {
	rootMapToReal := radix.New()
	realMapToRoot := radix.New()
	id := fmt.Sprintf("rfs-%d", rootMappingFsCounter.Add(1))

	addMapping := func(key string, rm RootMapping, to *radix.Tree) {
		var mappings []RootMapping
		v, found := to.Get(key)
		if found {
			// There may be more than one language pointing to the same root.
			mappings = v.([]RootMapping)
		}
		mappings = append(mappings, rm)
		to.Insert(key, mappings)
	}

	for _, rm := range rms {
		(&rm).clean()

		rm.FromBase = files.ResolveComponentFolder(rm.From)

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

		if rm.Meta == nil {
			rm.Meta = NewFileMeta()
		}

		if rm.FromBase == "" {
			panic(" rm.FromBase is empty")
		}

		rm.Meta.Component = rm.FromBase
		rm.Meta.Module = rm.Module
		rm.Meta.ModuleOrdinal = rm.ModuleOrdinal
		rm.Meta.IsProject = rm.IsProject
		rm.Meta.BaseDir = rm.ToBase

		if !fi.IsDir() {
			// We do allow single file mounts.
			// However, the file system logic will be much simpler with just directories.
			// So, convert this mount into a directory mount with a renamer,
			// which will tell the caller if name should be included.
			dirFrom, nameFrom := filepath.Split(rm.From)
			dirTo, nameTo := filepath.Split(rm.To)
			dirFrom, dirTo = strings.TrimSuffix(dirFrom, filepathSeparator), strings.TrimSuffix(dirTo, filepathSeparator)
			rm.From = dirFrom
			singleFileMeta := rm.Meta.Copy()
			singleFileMeta.Name = nameFrom
			rm.fiSingleFile = NewFileMetaInfo(fi, singleFileMeta)
			rm.To = dirTo

			rm.Meta.Rename = func(name string, toFrom bool) (string, bool) {
				if toFrom {
					if name == nameTo {
						return nameFrom, true
					}
					return "", false
				}

				if name == nameFrom {
					return nameTo, true
				}

				return "", false
			}
			nameToFilename := filepathSeparator + nameTo

			rm.Meta.InclusionFilter = rm.Meta.InclusionFilter.Append(glob.NewFilenameFilterForInclusionFunc(
				func(filename string) bool {
					return nameToFilename == filename
				},
			))

			// Refresh the FileInfo object.
			fi, err = fs.Stat(rm.To)
			if err != nil {
				if herrors.IsNotExist(err) {
					continue
				}
				return nil, err
			}
		}

		// Extract "blog" from "content/blog"
		rm.path = strings.TrimPrefix(strings.TrimPrefix(rm.From, rm.FromBase), filepathSeparator)
		rm.Meta.SourceRoot = fi.(MetaProvider).Meta().Filename

		meta := rm.Meta.Copy()

		if !fi.IsDir() {
			_, name := filepath.Split(rm.From)
			meta.Name = name
		}

		rm.fi = NewFileMetaInfo(fi, meta)

		addMapping(filepathSeparator+rm.From, rm, rootMapToReal)
		rev := rm.To
		if !strings.HasPrefix(rev, filepathSeparator) {
			rev = filepathSeparator + rev
		}

		addMapping(rev, rm, realMapToRoot)

	}

	rfs := &RootMappingFs{
		id:            id,
		Fs:            fs,
		rootMapToReal: rootMapToReal,
		realMapToRoot: realMapToRoot,
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
			From:   fromTo[j],
			To:     fromTo[j+1],
			ToBase: baseDir,
		}
	}

	return NewRootMappingFs(fs, rms...)
}

// RootMapping describes a virtual file or directory mount.
type RootMapping struct {
	From          string    // The virtual mount.
	FromBase      string    // The base directory of the virtual mount.
	To            string    // The source directory or file.
	ToBase        string    // The base of To. May be empty if an absolute path was provided.
	Module        string    // The module path/ID.
	ModuleOrdinal int       // The module ordinal starting with 0 which is the project.
	IsProject     bool      // Whether this is a mount in the main project.
	Meta          *FileMeta // File metadata (lang etc.)

	fi           FileMetaInfo
	fiSingleFile FileMetaInfo // Also set when this mounts represents a single file with a rename func.
	path         string       // The virtual mount point, e.g. "blog".
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

func (r RootMapping) trimFrom(name string) string {
	if name == "" {
		return ""
	}
	return strings.TrimPrefix(name, r.From)
}

var _ FilesystemUnwrapper = (*RootMappingFs)(nil)

// A RootMappingFs maps several roots into one. Note that the root of this filesystem
// is directories only, and they will be returned in Readdir and Readdirnames
// in the order given.
type RootMappingFs struct {
	id string
	afero.Fs
	rootMapToReal *radix.Tree
	realMapToRoot *radix.Tree
}

var rootMappingFsCounter atomic.Int32

func (fs *RootMappingFs) Mounts(base string) ([]FileMetaInfo, error) {
	base = filepathSeparator + fs.cleanName(base)
	roots := fs.getRootsWithPrefix(base)

	if roots == nil {
		return nil, nil
	}

	fss := make([]FileMetaInfo, len(roots))
	for i, r := range roots {
		if r.fiSingleFile != nil {
			// A single file mount.
			fss[i] = r.fiSingleFile
			continue
		}
		bfs := NewBasePathFs(fs.Fs, r.To)
		fs := bfs
		if r.Meta.InclusionFilter != nil {
			fs = newFilenameFilterFs(fs, r.To, r.Meta.InclusionFilter)
		}
		fs = decorateDirs(fs, r.Meta)
		fi, err := fs.Stat("")
		if err != nil {
			return nil, fmt.Errorf("RootMappingFs.Dirs: %w", err)
		}
		fss[i] = fi.(FileMetaInfo)
	}

	return fss, nil
}

func (fs *RootMappingFs) Key() string {
	return fs.id
}

func (fs *RootMappingFs) UnwrapFilesystem() afero.Fs {
	return fs.Fs
}

// Filter creates a copy of this filesystem with only mappings matching a filter.
func (fs RootMappingFs) Filter(f func(m RootMapping) bool) *RootMappingFs {
	rootMapToReal := radix.New()
	fs.rootMapToReal.Walk(func(b string, v any) bool {
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

// Open opens the named file for reading.
func (fs *RootMappingFs) Open(name string) (afero.File, error) {
	fis, err := fs.doStat(name)
	if err != nil {
		return nil, err
	}

	return fs.newUnionFile(fis...)
}

// Stat returns the os.FileInfo structure describing a given file.  If there is
// an error, it will be of type *os.PathError.
func (fs *RootMappingFs) Stat(name string) (os.FileInfo, error) {
	fis, err := fs.doStat(name)
	if err != nil {
		return nil, err
	}
	return fis[0], nil
}

type ComponentPath struct {
	Component string
	Path      string
	Lang      string
	Watch     bool
}

func (c ComponentPath) ComponentPathJoined() string {
	return path.Join(c.Component, c.Path)
}

type ReverseLookupProvder interface {
	ReverseLookup(filename string) ([]ComponentPath, error)
	ReverseLookupComponent(component, filename string) ([]ComponentPath, error)
}

// func (fs *RootMappingFs) ReverseStat(filename string) ([]FileMetaInfo, error)
func (fs *RootMappingFs) ReverseLookup(filename string) ([]ComponentPath, error) {
	return fs.ReverseLookupComponent("", filename)
}

func (fs *RootMappingFs) ReverseLookupComponent(component, filename string) ([]ComponentPath, error) {
	filename = fs.cleanName(filename)
	key := filepathSeparator + filename

	s, roots := fs.getRootsReverse(key)

	if len(roots) == 0 {
		return nil, nil
	}

	var cps []ComponentPath

	base := strings.TrimPrefix(key, s)
	dir, name := filepath.Split(base)

	for _, first := range roots {
		if component != "" && first.FromBase != component {
			continue
		}

		var filename string
		if first.Meta.Rename != nil {
			// Single file mount.
			if newname, ok := first.Meta.Rename(name, true); ok {
				filename = filepathSeparator + filepath.Join(first.path, dir, newname)
			} else {
				continue
			}
		} else {
			// Now we know that this file _could_ be in this fs.
			filename = filepathSeparator + filepath.Join(first.path, dir, name)
		}

		cps = append(cps, ComponentPath{
			Component: first.FromBase,
			Path:      paths.ToSlashTrimLeading(filename),
			Lang:      first.Meta.Lang,
			Watch:     first.Meta.Watch,
		})
	}

	return cps, nil
}

func (fs *RootMappingFs) hasPrefix(prefix string) bool {
	hasPrefix := false
	fs.rootMapToReal.WalkPrefix(prefix, func(b string, v any) bool {
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
	tree := fs.rootMapToReal
	levels := strings.Count(key, filepathSeparator)
	seen := make(map[RootMapping]bool)

	var roots []RootMapping
	var s string

	for {
		var found bool
		ss, vv, found := tree.LongestPrefix(key)

		if !found || (levels < 2 && ss == key) {
			break
		}

		for _, rm := range vv.([]RootMapping) {
			if !seen[rm] {
				seen[rm] = true
				roots = append(roots, rm)
			}
		}
		s = ss

		// We may have more than one root for this key, so walk up.
		oldKey := key
		key = filepath.Dir(key)
		if key == oldKey {
			break
		}
	}

	return s, roots
}

func (fs *RootMappingFs) getRootsReverse(key string) (string, []RootMapping) {
	tree := fs.realMapToRoot
	s, v, found := tree.LongestPrefix(key)
	if !found {
		return "", nil
	}
	return s, v.([]RootMapping)
}

func (fs *RootMappingFs) getRootsWithPrefix(prefix string) []RootMapping {
	var roots []RootMapping
	fs.rootMapToReal.WalkPrefix(prefix, func(b string, v any) bool {
		roots = append(roots, v.([]RootMapping)...)
		return false
	})

	return roots
}

func (fs *RootMappingFs) getAncestors(prefix string) []keyRootMappings {
	var roots []keyRootMappings
	fs.rootMapToReal.WalkPath(prefix, func(s string, v any) bool {
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
	if len(fis) == 1 {
		return fis[0].Meta().Open()
	}

	if !fis[0].IsDir() {
		// Pick the last file mount.
		return fis[len(fis)-1].Meta().Open()
	}

	openers := make([]func() (afero.File, error), len(fis))
	for i := len(fis) - 1; i >= 0; i-- {
		fi := fis[i]
		openers[i] = func() (afero.File, error) {
			meta := fi.Meta()
			f, err := meta.Open()
			if err != nil {
				return nil, err
			}
			return &rootMappingDir{DirOnlyOps: f, fs: fs, name: meta.Name, meta: meta}, nil
		}
	}

	merge := func(lofi, bofi []iofs.DirEntry) []iofs.DirEntry {
		// Ignore duplicate directory entries
		for _, fi1 := range bofi {
			var found bool
			for _, fi2 := range lofi {
				if !fi2.IsDir() {
					continue
				}
				if fi1.Name() == fi2.Name() {
					found = true
					break
				}
			}
			if !found {
				lofi = append(lofi, fi1)
			}
		}

		return lofi
	}

	info := func() (os.FileInfo, error) {
		return fis[0], nil
	}

	return overlayfs.OpenDir(merge, info, openers...)
}

func (fs *RootMappingFs) cleanName(name string) string {
	name = strings.Trim(filepath.Clean(name), filepathSeparator)
	if name == "." {
		name = ""
	}
	return name
}

func (rfs *RootMappingFs) collectDirEntries(prefix string) ([]iofs.DirEntry, error) {
	prefix = filepathSeparator + rfs.cleanName(prefix)

	var fis []iofs.DirEntry

	seen := make(map[string]bool) // Prevent duplicate directories
	level := strings.Count(prefix, filepathSeparator)

	collectDir := func(rm RootMapping, fi FileMetaInfo) error {
		f, err := fi.Meta().Open()
		if err != nil {
			return err
		}
		direntries, err := f.(iofs.ReadDirFile).ReadDir(-1)
		if err != nil {
			f.Close()
			return err
		}

		for _, fi := range direntries {

			meta := fi.(FileMetaInfo).Meta()
			meta.Merge(rm.Meta)

			if !rm.Meta.InclusionFilter.Match(strings.TrimPrefix(meta.Filename, meta.SourceRoot), fi.IsDir()) {
				continue
			}

			if fi.IsDir() {
				name := fi.Name()
				if seen[name] {
					continue
				}
				seen[name] = true
				opener := func() (afero.File, error) {
					return rfs.Open(filepath.Join(rm.From, name))
				}
				fi = newDirNameOnlyFileInfo(name, meta, opener)
			} else if rm.Meta.Rename != nil {
				n, ok := rm.Meta.Rename(fi.Name(), true)
				if !ok {
					continue
				}
				fi.(MetaProvider).Meta().Name = n
			}
			fis = append(fis, fi)
		}

		f.Close()

		return nil
	}

	// First add any real files/directories.
	rms := rfs.getRoot(prefix)
	for _, rm := range rms {
		if err := collectDir(rm, rm.fi); err != nil {
			return nil, err
		}
	}

	// Next add any file mounts inside the given directory.
	prefixInside := prefix + filepathSeparator
	rfs.rootMapToReal.WalkPrefix(prefixInside, func(s string, v any) bool {
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
				return rfs.Open(path)
			}

			fi := newDirNameOnlyFileInfo(name, nil, opener)
			fis = append(fis, fi)

			return false
		}

		rms := v.([]RootMapping)
		for _, rm := range rms {
			name := filepath.Base(rm.From)
			if seen[name] {
				continue
			}
			seen[name] = true
			opener := func() (afero.File, error) {
				return rfs.Open(rm.From)
			}
			fi := newDirNameOnlyFileInfo(name, rm.Meta, opener)
			fis = append(fis, fi)
		}

		return false
	})

	// Finally add any ancestor dirs with files in this directory.
	ancestors := rfs.getAncestors(prefix)
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

func (fs *RootMappingFs) doStat(name string) ([]FileMetaInfo, error) {
	fis, err := fs.doDoStat(name)
	if err != nil {
		return nil, err
	}
	// Sanity check. Check that all is either file or directories.
	var isDir, isFile bool
	for _, fi := range fis {
		if fi.IsDir() {
			isDir = true
		} else {
			isFile = true
		}
	}
	if isDir && isFile {
		// For now.
		return nil, os.ErrNotExist
	}

	return fis, nil
}

func (fs *RootMappingFs) doDoStat(name string) ([]FileMetaInfo, error) {
	name = fs.cleanName(name)
	key := filepathSeparator + name

	roots := fs.getRoot(key)

	if roots == nil {
		if fs.hasPrefix(key) {
			// We have directories mounted below this.
			// Make it look like a directory.
			return []FileMetaInfo{newDirNameOnlyFileInfo(name, nil, fs.virtualDirOpener(name))}, nil
		}

		// Find any real directories with this key.
		_, roots := fs.getRoots(key)
		if roots == nil {
			return nil, &os.PathError{Op: "LStat", Path: name, Err: os.ErrNotExist}
		}

		var err error
		var fis []FileMetaInfo

		for _, rm := range roots {
			var fi FileMetaInfo
			fi, err = fs.statRoot(rm, name)
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

	return []FileMetaInfo{newDirNameOnlyFileInfo(name, roots[0].Meta, fs.virtualDirOpener(name))}, nil
}

func (fs *RootMappingFs) statRoot(root RootMapping, filename string) (FileMetaInfo, error) {
	dir, name := filepath.Split(filename)
	if root.Meta.Rename != nil {
		n, ok := root.Meta.Rename(name, false)
		if !ok {
			return nil, os.ErrNotExist
		}
		filename = filepath.Join(dir, n)
	}

	if !root.Meta.InclusionFilter.Match(root.trimFrom(filename), true) {
		return nil, os.ErrNotExist
	}

	filename = root.filename(filename)
	fi, err := fs.Fs.Stat(filename)
	if err != nil {
		return nil, err
	}

	var opener func() (afero.File, error)
	if !fi.IsDir() {
		// Open the file directly.
		// Opens the real file directly.
		opener = func() (afero.File, error) {
			return fs.Fs.Open(filename)
		}
	} else if root.Meta.Rename != nil {
		// A single file mount where we have mounted the containing directory.
		n, ok := root.Meta.Rename(fi.Name(), true)
		if !ok {
			return nil, os.ErrNotExist
		}
		meta := fi.(MetaProvider).Meta()
		meta.Name = n
		// Opens the real file directly.
		opener = func() (afero.File, error) {
			return fs.Fs.Open(filename)
		}
	} else {
		// Make sure metadata gets applied in ReadDir.
		opener = fs.realDirOpener(filename, root.Meta)
	}

	fim := decorateFileInfo(fi, opener, "", root.Meta)

	return fim, nil
}

func (fs *RootMappingFs) virtualDirOpener(name string) func() (afero.File, error) {
	return func() (afero.File, error) { return &rootMappingDir{name: name, fs: fs}, nil }
}

func (fs *RootMappingFs) realDirOpener(name string, meta *FileMeta) func() (afero.File, error) {
	return func() (afero.File, error) {
		f, err := fs.Fs.Open(name)
		if err != nil {
			return nil, err
		}
		return &rootMappingDir{name: name, meta: meta, fs: fs, DirOnlyOps: f}, nil
	}
}

var _ iofs.ReadDirFile = (*rootMappingDir)(nil)

type rootMappingDir struct {
	*noOpRegularFileOps
	DirOnlyOps
	fs   *RootMappingFs
	name string
	meta *FileMeta
}

func (f *rootMappingDir) Close() error {
	if f.DirOnlyOps == nil {
		return nil
	}
	return f.DirOnlyOps.Close()
}

func (f *rootMappingDir) Name() string {
	return f.name
}

func (f *rootMappingDir) ReadDir(count int) ([]iofs.DirEntry, error) {
	if f.DirOnlyOps != nil {
		fis, err := f.DirOnlyOps.(iofs.ReadDirFile).ReadDir(count)
		if err != nil {
			return nil, err
		}

		var result []iofs.DirEntry
		for _, fi := range fis {
			fim := decorateFileInfo(fi, nil, "", f.meta)
			meta := fim.Meta()
			if f.meta.InclusionFilter.Match(strings.TrimPrefix(meta.Filename, meta.SourceRoot), fim.IsDir()) {
				result = append(result, fim)
			}
		}
		return result, nil
	}

	return f.fs.collectDirEntries(f.name)
}

// Sentinel error to signal that a file is a directory.
var errIsDir = errors.New("isDir")

func (f *rootMappingDir) Stat() (iofs.FileInfo, error) {
	return nil, errIsDir
}

func (f *rootMappingDir) Readdir(count int) ([]os.FileInfo, error) {
	panic("not supported: use ReadDir")
}

// Note that Readdirnames preserves the order of the underlying filesystem(s),
// which is usually directory order.
func (f *rootMappingDir) Readdirnames(count int) ([]string, error) {
	dirs, err := f.ReadDir(count)
	if err != nil {
		return nil, err
	}
	return dirEntriesToNames(dirs), nil
}

func dirEntriesToNames(fis []iofs.DirEntry) []string {
	names := make([]string, len(fis))
	for i, d := range fis {
		names[i] = d.Name()
	}
	return names
}
