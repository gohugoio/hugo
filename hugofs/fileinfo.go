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
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/hugofs/files"
	"golang.org/x/text/unicode/norm"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/htime"

	"github.com/spf13/afero"
)

func NewFileMeta() *FileMeta {
	return &FileMeta{}
}

// PathFile returns the relative file path for the file source.
func (f *FileMeta) PathFile() string {
	if f.BaseDir == "" {
		return ""
	}
	return strings.TrimPrefix(strings.TrimPrefix(f.Filename, f.BaseDir), filepathSeparator)
}

type FileMeta struct {
	Name             string
	Filename         string
	Path             string
	PathWalk         string
	OriginalFilename string
	BaseDir          string

	SourceRoot string
	MountRoot  string
	Module     string

	Weight     int
	IsOrdered  bool
	IsSymlink  bool
	IsRootFile bool
	IsProject  bool
	Watch      bool

	Classifier files.ContentClass

	SkipDir bool

	Lang                       string
	TranslationBaseName        string
	TranslationBaseNameWithExt string
	Translations               []string

	Fs           afero.Fs
	OpenFunc     func() (afero.File, error)
	JoinStatFunc func(name string) (FileMetaInfo, error)

	// Include only files or directories that match.
	InclusionFilter *glob.FilenameFilter
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

func (f *FileMeta) JoinStat(name string) (FileMetaInfo, error) {
	if f.JoinStatFunc == nil {
		return nil, os.ErrNotExist
	}
	return f.JoinStatFunc(name)
}

type FileMetaInfo interface {
	os.FileInfo
	// Meta is for internal use.
	Meta() *FileMeta
}

type fileInfoMeta struct {
	os.FileInfo

	m *FileMeta
}

type filenameProvider interface {
	Filename() string
}

var _ filenameProvider = (*fileInfoMeta)(nil)

// Filename returns the full filename.
func (fi *fileInfoMeta) Filename() string {
	return fi.m.Filename
}

// Name returns the file's name. Note that we follow symlinks,
// if supported by the file system, and the Name given here will be the
// name of the symlink, which is what Hugo needs in all situations.
func (fi *fileInfoMeta) Name() string {
	if name := fi.m.Name; name != "" {
		return name
	}
	return fi.FileInfo.Name()
}

func (fi *fileInfoMeta) Meta() *FileMeta {
	return fi.m
}

func NewFileMetaInfo(fi os.FileInfo, m *FileMeta) FileMetaInfo {
	if m == nil {
		panic("FileMeta must be set")
	}
	if fim, ok := fi.(FileMetaInfo); ok {
		m.Merge(fim.Meta())
	}
	return &fileInfoMeta{FileInfo: fi, m: m}
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
	m.IsOrdered = false

	return NewFileMetaInfo(
		&dirNameOnlyFileInfo{name: base, modTime: htime.Now()},
		m,
	)
}

func decorateFileInfo(
	fi os.FileInfo,
	fs afero.Fs, opener func() (afero.File, error),
	filename, filepath string, inMeta *FileMeta,
) FileMetaInfo {
	var meta *FileMeta
	var fim FileMetaInfo

	filepath = strings.TrimPrefix(filepath, filepathSeparator)

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
	if fs != nil {
		meta.Fs = fs
	}
	nfilepath := normalizeFilename(filepath)
	nfilename := normalizeFilename(filename)
	if nfilepath != "" {
		meta.Path = nfilepath
	}
	if nfilename != "" {
		meta.Filename = nfilename
	}

	meta.Merge(inMeta)

	return fim
}

func isSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink == os.ModeSymlink
}

func fileInfosToFileMetaInfos(fis []os.FileInfo) []FileMetaInfo {
	fims := make([]FileMetaInfo, len(fis))
	for i, v := range fis {
		fims[i] = v.(FileMetaInfo)
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

func fileInfosToNames(fis []os.FileInfo) []string {
	names := make([]string, len(fis))
	for i, d := range fis {
		names[i] = d.Name()
	}
	return names
}

func sortFileInfos(fis []os.FileInfo) {
	sort.Slice(fis, func(i, j int) bool {
		fimi, fimj := fis[i].(FileMetaInfo), fis[j].(FileMetaInfo)
		return fimi.Meta().Filename < fimj.Meta().Filename
	})
}
