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
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/gohugoio/hugo/hugofs/files"
	"golang.org/x/text/unicode/norm"

	"github.com/pkg/errors"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/hreflect"

	"github.com/spf13/afero"
)

const (
	metaKeyFilename = "filename"

	metaKeyBaseDir                    = "baseDir" // Abs base directory of source file.
	metaKeyMountRoot                  = "mountRoot"
	metaKeyModule                     = "module"
	metaKeyOriginalFilename           = "originalFilename"
	metaKeyName                       = "name"
	metaKeyPath                       = "path"
	metaKeyPathWalk                   = "pathWalk"
	metaKeyLang                       = "lang"
	metaKeyWeight                     = "weight"
	metaKeyOrdinal                    = "ordinal"
	metaKeyFs                         = "fs"
	metaKeyOpener                     = "opener"
	metaKeyIsOrdered                  = "isOrdered"
	metaKeyIsSymlink                  = "isSymlink"
	metaKeyJoinStat                   = "joinStat"
	metaKeySkipDir                    = "skipDir"
	metaKeyClassifier                 = "classifier"
	metaKeyTranslationBaseName        = "translationBaseName"
	metaKeyTranslationBaseNameWithExt = "translationBaseNameWithExt"
	metaKeyTranslations               = "translations"
	metaKeyDecoraterPath              = "decoratorPath"
)

type FileMeta map[string]interface{}

func (f FileMeta) GetInt(key string) int {
	return cast.ToInt(f[key])
}

func (f FileMeta) GetString(key string) string {
	return cast.ToString(f[key])
}

func (f FileMeta) GetBool(key string) bool {
	return cast.ToBool(f[key])
}

func (f FileMeta) Filename() string {
	return f.stringV(metaKeyFilename)
}

func (f FileMeta) OriginalFilename() string {
	return f.stringV(metaKeyOriginalFilename)
}

func (f FileMeta) SkipDir() bool {
	return f.GetBool(metaKeySkipDir)
}
func (f FileMeta) TranslationBaseName() string {
	return f.stringV(metaKeyTranslationBaseName)
}

func (f FileMeta) TranslationBaseNameWithExt() string {
	return f.stringV(metaKeyTranslationBaseNameWithExt)
}

func (f FileMeta) Translations() []string {
	return cast.ToStringSlice(f[metaKeyTranslations])
}

func (f FileMeta) Name() string {
	return f.stringV(metaKeyName)
}

func (f FileMeta) Classifier() files.ContentClass {
	c, found := f[metaKeyClassifier]
	if found {
		return c.(files.ContentClass)
	}

	return files.ContentClassFile // For sorting
}

func (f FileMeta) Lang() string {
	return f.stringV(metaKeyLang)
}

// Path returns the relative file path to where this file is mounted.
func (f FileMeta) Path() string {
	return f.stringV(metaKeyPath)
}

// PathFile returns the relative file path for the file source.
func (f FileMeta) PathFile() string {
	base := f.stringV(metaKeyBaseDir)
	if base == "" {
		return ""
	}
	return strings.TrimPrefix(strings.TrimPrefix(f.Filename(), base), filepathSeparator)
}

func (f FileMeta) MountRoot() string {
	return f.stringV(metaKeyMountRoot)
}

func (f FileMeta) Module() string {
	return f.stringV(metaKeyModule)
}

func (f FileMeta) Weight() int {
	return f.GetInt(metaKeyWeight)
}

func (f FileMeta) Ordinal() int {
	return f.GetInt(metaKeyOrdinal)
}

func (f FileMeta) IsOrdered() bool {
	return f.GetBool(metaKeyIsOrdered)
}

// IsSymlink returns whether this comes from a symlinked file or directory.
func (f FileMeta) IsSymlink() bool {
	return f.GetBool(metaKeyIsSymlink)
}

func (f FileMeta) Watch() bool {
	if v, found := f["watch"]; found {
		return v.(bool)
	}
	return false
}

func (f FileMeta) Fs() afero.Fs {
	if v, found := f[metaKeyFs]; found {
		return v.(afero.Fs)
	}
	return nil
}

func (f FileMeta) GetOpener() func() (afero.File, error) {
	o, found := f[metaKeyOpener]
	if !found {
		return nil
	}
	return o.(func() (afero.File, error))
}

func (f FileMeta) Open() (afero.File, error) {
	v, found := f[metaKeyOpener]
	if !found {
		return nil, errors.New("file opener not found")
	}
	return v.(func() (afero.File, error))()
}

func (f FileMeta) JoinStat(name string) (FileMetaInfo, error) {
	v, found := f[metaKeyJoinStat]
	if !found {
		return nil, os.ErrNotExist
	}
	return v.(func(name string) (FileMetaInfo, error))(name)
}

func (f FileMeta) stringV(key string) string {
	if v, found := f[key]; found {
		return v.(string)
	}
	return ""
}

func (f FileMeta) setIfNotZero(key string, val interface{}) {
	if !hreflect.IsTruthful(val) {
		return
	}
	f[key] = val
}

type FileMetaInfo interface {
	os.FileInfo
	Meta() FileMeta
}

type fileInfoMeta struct {
	os.FileInfo

	m FileMeta
}

// Name returns the file's name. Note that we follow symlinks,
// if supported by the file system, and the Name given here will be the
// name of the symlink, which is what Hugo needs in all situations.
func (fi *fileInfoMeta) Name() string {
	if name := fi.m.Name(); name != "" {
		return name
	}
	return fi.FileInfo.Name()
}

func (fi *fileInfoMeta) Meta() FileMeta {
	return fi.m
}

func NewFileMetaInfo(fi os.FileInfo, m FileMeta) FileMetaInfo {

	if fim, ok := fi.(FileMetaInfo); ok {
		mergeFileMeta(fim.Meta(), m)
	}
	return &fileInfoMeta{FileInfo: fi, m: m}
}

func copyFileMeta(m FileMeta) FileMeta {
	c := make(FileMeta)
	for k, v := range m {
		c[k] = v
	}
	return c
}

// Merge metadata, last entry wins.
func mergeFileMeta(from, to FileMeta) {
	if from == nil {
		return
	}
	for k, v := range from {
		if _, found := to[k]; !found {
			to[k] = v
		}
	}
}

type dirNameOnlyFileInfo struct {
	name string
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
	return time.Time{}
}

func (fi *dirNameOnlyFileInfo) IsDir() bool {
	return true
}

func (fi *dirNameOnlyFileInfo) Sys() interface{} {
	return nil
}

func newDirNameOnlyFileInfo(name string, meta FileMeta, fileOpener func() (afero.File, error)) FileMetaInfo {
	name = normalizeFilename(name)
	_, base := filepath.Split(name)

	m := copyFileMeta(meta)
	if _, found := m[metaKeyFilename]; !found {
		m.setIfNotZero(metaKeyFilename, name)
	}
	m[metaKeyOpener] = fileOpener
	m[metaKeyIsOrdered] = false

	return NewFileMetaInfo(
		&dirNameOnlyFileInfo{name: base},
		m,
	)
}

func decorateFileInfo(
	fi os.FileInfo,
	fs afero.Fs, opener func() (afero.File, error),
	filename, filepath string, inMeta FileMeta) FileMetaInfo {

	var meta FileMeta
	var fim FileMetaInfo

	filepath = strings.TrimPrefix(filepath, filepathSeparator)

	var ok bool
	if fim, ok = fi.(FileMetaInfo); ok {
		meta = fim.Meta()
	} else {
		meta = make(FileMeta)
		fim = NewFileMetaInfo(fi, meta)
	}

	meta.setIfNotZero(metaKeyOpener, opener)
	meta.setIfNotZero(metaKeyFs, fs)
	meta.setIfNotZero(metaKeyPath, normalizeFilename(filepath))
	meta.setIfNotZero(metaKeyFilename, normalizeFilename(filename))

	mergeFileMeta(inMeta, meta)

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

func fromSlash(filenames []string) []string {
	for i, name := range filenames {
		filenames[i] = filepath.FromSlash(name)
	}
	return filenames
}

func sortFileInfos(fis []os.FileInfo) {
	sort.Slice(fis, func(i, j int) bool {
		fimi, fimj := fis[i].(FileMetaInfo), fis[j].(FileMetaInfo)
		return fimi.Meta().Filename() < fimj.Meta().Filename()

	})
}
