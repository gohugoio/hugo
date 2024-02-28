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

// Package hugofs provides the file systems used by Hugo.
package hugofs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/gohugoio/hugo/hugofs/glob"

	"golang.org/x/text/unicode/norm"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/spf13/afero"
)

func NewFileMeta() *FileMeta {
	return &FileMeta{}
}

type FileMeta struct {
	PathInfo *paths.Path
	Name     string
	Filename string

	BaseDir       string
	SourceRoot    string
	Module        string
	ModuleOrdinal int
	Component     string

	Weight    int
	IsProject bool
	Watch     bool

	// The lang associated with this file. This may be
	// either the language set in the filename or
	// the language defined in the source mount configuration.
	Lang string
	// The language index for the above lang. This is the index
	// in the sorted list of languages/sites.
	LangIndex int

	OpenFunc     func() (afero.File, error)
	JoinStatFunc func(name string) (FileMetaInfo, error)

	// Include only files or directories that match.
	InclusionFilter *glob.FilenameFilter

	// Rename the name part of the file (not the directory).
	// Returns the new name and a boolean indicating if the file
	// should be included.
	Rename func(name string, toFrom bool) (string, bool)
}

func (m *FileMeta) Copy() *FileMeta {
	if m == nil {
		return NewFileMeta()
	}
	c := *m
	return &c
}

func (m *FileMeta) Merge(from *FileMeta) {
	if m == nil || from == nil {
		return
	}
	dstv := reflect.Indirect(reflect.ValueOf(m))
	srcv := reflect.Indirect(reflect.ValueOf(from))

	for i := 0; i < dstv.NumField(); i++ {
		v := dstv.Field(i)
		if !v.CanSet() {
			continue
		}
		if !hreflect.IsTruthfulValue(v) {
			v.Set(srcv.Field(i))
		}
	}

	if m.InclusionFilter == nil {
		m.InclusionFilter = from.InclusionFilter
	}
}

func (f *FileMeta) Open() (afero.File, error) {
	if f.OpenFunc == nil {
		return nil, errors.New("OpenFunc not set")
	}
	return f.OpenFunc()
}

func (f *FileMeta) ReadAll() ([]byte, error) {
	file, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func (f *FileMeta) JoinStat(name string) (FileMetaInfo, error) {
	if f.JoinStatFunc == nil {
		return nil, os.ErrNotExist
	}
	return f.JoinStatFunc(name)
}

type FileMetaInfo interface {
	fs.DirEntry
	MetaProvider

	// This is a real hybrid as it also implements the fs.FileInfo interface.
	FileInfoOptionals
}

type MetaProvider interface {
	Meta() *FileMeta
}

type FileInfoOptionals interface {
	Size() int64
	Mode() fs.FileMode
	ModTime() time.Time
	Sys() any
}

type FileNameIsDir interface {
	Name() string
	IsDir() bool
}

type FileInfoProvider interface {
	FileInfo() FileMetaInfo
}

// DirOnlyOps is a subset of the afero.File interface covering
// the methods needed for directory operations.
type DirOnlyOps interface {
	io.Closer
	Name() string
	Readdir(count int) ([]os.FileInfo, error)
	Readdirnames(n int) ([]string, error)
	Stat() (os.FileInfo, error)
}

type dirEntryMeta struct {
	fs.DirEntry
	m *FileMeta

	fi     fs.FileInfo
	fiInit sync.Once
}

func (fi *dirEntryMeta) Meta() *FileMeta {
	return fi.m
}

// Filename returns the full filename.
func (fi *dirEntryMeta) Filename() string {
	return fi.m.Filename
}

func (fi *dirEntryMeta) fileInfo() fs.FileInfo {
	var err error
	fi.fiInit.Do(func() {
		fi.fi, err = fi.DirEntry.Info()
	})
	if err != nil {
		panic(err)
	}
	return fi.fi
}

func (fi *dirEntryMeta) Size() int64 {
	return fi.fileInfo().Size()
}

func (fi *dirEntryMeta) Mode() fs.FileMode {
	return fi.fileInfo().Mode()
}

func (fi *dirEntryMeta) ModTime() time.Time {
	return fi.fileInfo().ModTime()
}

func (fi *dirEntryMeta) Sys() any {
	return fi.fileInfo().Sys()
}

// Name returns the file's name.
func (fi *dirEntryMeta) Name() string {
	if name := fi.m.Name; name != "" {
		return name
	}
	return fi.DirEntry.Name()
}

// dirEntry is an adapter from os.FileInfo to fs.DirEntry
type dirEntry struct {
	fs.FileInfo
}

var _ fs.DirEntry = dirEntry{}

func (d dirEntry) Type() fs.FileMode { return d.FileInfo.Mode().Type() }

func (d dirEntry) Info() (fs.FileInfo, error) { return d.FileInfo, nil }

func NewFileMetaInfo(fi FileNameIsDir, m *FileMeta) FileMetaInfo {
	if m == nil {
		panic("FileMeta must be set")
	}
	if fim, ok := fi.(MetaProvider); ok {
		m.Merge(fim.Meta())
	}
	switch v := fi.(type) {
	case fs.DirEntry:
		return &dirEntryMeta{DirEntry: v, m: m}
	case fs.FileInfo:
		return &dirEntryMeta{DirEntry: dirEntry{v}, m: m}
	case nil:
		return &dirEntryMeta{DirEntry: dirEntry{}, m: m}
	default:
		panic(fmt.Sprintf("Unsupported type: %T", fi))
	}
}

type dirNameOnlyFileInfo struct {
	name    string
	modTime time.Time
}

func (fi *dirNameOnlyFileInfo) Name() string {
	return fi.name
}

func (fi *dirNameOnlyFileInfo) Size() int64 {
	panic("not implemented")
}

func (fi *dirNameOnlyFileInfo) Mode() os.FileMode {
	return os.ModeDir
}

func (fi *dirNameOnlyFileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *dirNameOnlyFileInfo) IsDir() bool {
	return true
}

func (fi *dirNameOnlyFileInfo) Sys() any {
	return nil
}

func newDirNameOnlyFileInfo(name string, meta *FileMeta, fileOpener func() (afero.File, error)) FileMetaInfo {
	name = normalizeFilename(name)
	_, base := filepath.Split(name)

	m := meta.Copy()
	if m.Filename == "" {
		m.Filename = name
	}
	m.OpenFunc = fileOpener

	return NewFileMetaInfo(
		&dirNameOnlyFileInfo{name: base, modTime: htime.Now()},
		m,
	)
}

func decorateFileInfo(fi FileNameIsDir, opener func() (afero.File, error), filename string, inMeta *FileMeta) FileMetaInfo {
	var meta *FileMeta
	var fim FileMetaInfo

	var ok bool
	if fim, ok = fi.(FileMetaInfo); ok {
		meta = fim.Meta()
	} else {
		meta = NewFileMeta()
		fim = NewFileMetaInfo(fi, meta)
	}

	if opener != nil {
		meta.OpenFunc = opener
	}

	nfilename := normalizeFilename(filename)
	if nfilename != "" {
		meta.Filename = nfilename
	}

	meta.Merge(inMeta)

	return fim
}

func DirEntriesToFileMetaInfos(fis []fs.DirEntry) []FileMetaInfo {
	fims := make([]FileMetaInfo, len(fis))
	for i, v := range fis {
		fim := v.(FileMetaInfo)
		fims[i] = fim
	}
	return fims
}

func normalizeFilename(filename string) string {
	if filename == "" {
		return ""
	}
	if runtime.GOOS == "darwin" {
		// When a file system is HFS+, its filepath is in NFD form.
		return norm.NFC.String(filename)
	}
	return filename
}

func sortDirEntries(fis []fs.DirEntry) {
	sort.Slice(fis, func(i, j int) bool {
		fimi, fimj := fis[i].(FileMetaInfo), fis[j].(FileMetaInfo)
		return fimi.Meta().Filename < fimj.Meta().Filename
	})
}

// AddFileInfoToError adds file info to the given error.
func AddFileInfoToError(err error, fi FileMetaInfo, fs afero.Fs) error {
	if err == nil {
		return nil
	}

	meta := fi.Meta()
	filename := meta.Filename

	// Check if it's already added.
	for _, ferr := range herrors.UnwrapFileErrors(err) {
		pos := ferr.Position()
		errfilename := pos.Filename
		if errfilename == "" {
			pos.Filename = filename
			ferr.UpdatePosition(pos)
		}

		if errfilename == "" || errfilename == filename {
			if filename != "" && ferr.ErrorContext() == nil {
				f, ioerr := fs.Open(filename)
				if ioerr != nil {
					return err
				}
				defer f.Close()
				ferr.UpdateContent(f, nil)
			}
			return err
		}
	}

	lineMatcher := herrors.NopLineMatcher

	if textSegmentErr, ok := err.(*herrors.TextSegmentError); ok {
		lineMatcher = herrors.ContainsMatcher(textSegmentErr.Segment)
	}

	return herrors.NewFileErrorFromFile(err, filename, fs, lineMatcher)
}
