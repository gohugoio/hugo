// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"sync"

	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/helpers"
)

// fileInfo implements the File interface.
var (
	_ File = (*FileInfo)(nil)
)

// File represents a source file.
// This is a temporary construct until we resolve page.Page conflicts.
// TODO(bep) remove this construct once we have resolved page deprecations
type File interface {
	fileOverlap
	FileWithoutOverlap
}

// Temporary to solve duplicate/deprecated names in page.Page
type fileOverlap interface {
	// Path gets the relative path including file name and extension.
	// The directory is relative to the content root.
	Path() string

	// Section is first directory below the content root.
	// For page bundles in root, the Section will be empty.
	Section() string

	// Lang is the language code for this page. It will be the
	// same as the site's language code.
	Lang() string

	IsZero() bool
}

type FileWithoutOverlap interface {

	// Filename gets the full path and filename to the file.
	Filename() string

	// Dir gets the name of the directory that contains this file.
	// The directory is relative to the content root.
	Dir() string

	// Extension gets the file extension, i.e "myblogpost.md" will return "md".
	Extension() string

	// Ext is an alias for Extension.
	Ext() string // Hmm... Deprecate Extension

	// LogicalName is filename and extension of the file.
	LogicalName() string

	// BaseFileName is a filename without extension.
	BaseFileName() string

	// TranslationBaseName is a filename with no extension,
	// not even the optional language extension part.
	TranslationBaseName() string

	// ContentBaseName is a either TranslationBaseName or name of containing folder
	// if file is a leaf bundle.
	ContentBaseName() string

	// UniqueID is the MD5 hash of the file's path and is for most practical applications,
	// Hugo content files being one of them, considered to be unique.
	UniqueID() string

	FileInfo() hugofs.FileMetaInfo
}

// FileInfo describes a source file.
type FileInfo struct {

	// Absolute filename to the file on disk.
	filename string

	sp *SourceSpec

	fim hugofs.FileMetaInfo

	uniqueID string

	lazyInit sync.Once
}

func (fi *FileInfo) pathToDir(s string) string {
	return filepath.FromSlash(s[1:] + "/")
}

func (fi *FileInfo) p() paths.Path {
	return fi.fim.Meta().PathInfo
}

// Filename returns a file's absolute path and filename on disk.
func (fi *FileInfo) Filename() string { return fi.fim.Meta().Filename }

// Path gets the relative path including file name and extension.  The directory
// is relative to the content root.
func (fi *FileInfo) Path() string { return filepath.Join(fi.p().Dir()[1:], fi.p().Name()) }

// Dir gets the name of the directory that contains this file.  The directory is
// relative to the content root.
func (fi *FileInfo) Dir() string {
	return fi.pathToDir(fi.p().Dir())
}

// Extension is an alias to Ext().
func (fi *FileInfo) Extension() string {
	helpers.Deprecated(".File.Extension()", ".File.Ext()", false)
	return fi.Ext()
}

// Ext returns a file's extension without the leading period (ie. "md").
func (fi *FileInfo) Ext() string { return fi.p().Ext() }

// Lang returns a file's language (ie. "sv").
func (fi *FileInfo) Lang() string { return fi.p().Identifier(1) }

// LogicalName returns a file's name and extension (ie. "page.sv.md").
func (fi *FileInfo) LogicalName() string {
	return fi.p().Name()
}

// BaseFileName returns a file's name without extension (ie. "page.sv").
func (fi *FileInfo) BaseFileName() string {
	return fi.p().NameNoExt()
}

// TranslationBaseName returns a file's translation base name without the
// language segment (ie. "page").
func (fi *FileInfo) TranslationBaseName() string { return fi.p().NameNoIdentifier() }

// ContentBaseName is a either TranslationBaseName or name of containing folder
// if file is a bundle.
func (fi *FileInfo) ContentBaseName() string {
	if fi.p().IsBundle() {
		return fi.p().Container()
	}
	return fi.p().NameNoIdentifier()
}

// Section returns a file's section.
func (fi *FileInfo) Section() string {
	return fi.p().Section()
}

// UniqueID returns a file's unique, MD5 hash identifier.
func (fi *FileInfo) UniqueID() string {
	fi.init()
	return fi.uniqueID
}

// FileInfo returns a file's underlying os.FileInfo.
func (fi *FileInfo) FileInfo() hugofs.FileMetaInfo { return fi.fim }

func (fi *FileInfo) String() string { return fi.BaseFileName() }

// Open implements ReadableFile.
func (fi *FileInfo) Open() (hugio.ReadSeekCloser, error) {
	f, err := fi.fim.Meta().Open()

	return f, err
}

func (fi *FileInfo) IsZero() bool {
	return fi == nil
}

// We create a lot of these FileInfo objects, but there are parts of it used only
// in some cases that is slightly expensive to construct.
func (fi *FileInfo) init() {
	fi.lazyInit.Do(func() {
		fi.uniqueID = helpers.MD5String(filepath.ToSlash(fi.Path()))
	})
}

// NewTestFile creates a partially filled File used in unit tests.
// TODO(bep) improve this package
func NewTestFile(filename string) *FileInfo {
	return &FileInfo{
		filename: filename,
	}
}

func (sp *SourceSpec) NewFileInfoFrom(path, filename string) (*FileInfo, error) {
	meta := &hugofs.FileMeta{
		Filename: filename,
		Path:     path,
	}

	return sp.NewFileInfo(hugofs.NewFileMetaInfo(nil, meta))
}

func (sp *SourceSpec) NewFileInfo(fi hugofs.FileMetaInfo) (*FileInfo, error) {
	m := fi.Meta()

	f := &FileInfo{
		sp:       sp,
		filename: m.Filename,
		fim:      fi,
	}

	return f, nil
}
