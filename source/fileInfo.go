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
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/hugofs/files"

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

	// Extension is an alias to Ext().
	// Deprecated: Use Ext instead.
	Extension() string

	// Ext gets the file extension, i.e "myblogpost.md" will return "md".
	Ext() string

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

	fi hugofs.FileMetaInfo

	// Derived from filename
	ext  string // Extension without any "."
	lang string

	name string

	dir                 string
	relDir              string
	relPath             string
	baseName            string
	translationBaseName string
	contentBaseName     string
	section             string
	classifier          files.ContentClass

	uniqueID string

	lazyInit sync.Once
}

// Filename returns a file's absolute path and filename on disk.
func (fi *FileInfo) Filename() string { return fi.filename }

// Path gets the relative path including file name and extension.  The directory
// is relative to the content root.
func (fi *FileInfo) Path() string { return fi.relPath }

// Dir gets the name of the directory that contains this file.  The directory is
// relative to the content root.
func (fi *FileInfo) Dir() string { return fi.relDir }

// Extension is an alias to Ext().
func (fi *FileInfo) Extension() string {
	helpers.Deprecated(".File.Extension", "Use .File.Ext instead. ", false)
	return fi.Ext()
}

// Ext returns a file's extension without the leading period (ie. "md").
func (fi *FileInfo) Ext() string { return fi.ext }

// Lang returns a file's language (ie. "sv").
func (fi *FileInfo) Lang() string { return fi.lang }

// LogicalName returns a file's name and extension (ie. "page.sv.md").
func (fi *FileInfo) LogicalName() string { return fi.name }

// BaseFileName returns a file's name without extension (ie. "page.sv").
func (fi *FileInfo) BaseFileName() string { return fi.baseName }

// TranslationBaseName returns a file's translation base name without the
// language segment (ie. "page").
func (fi *FileInfo) TranslationBaseName() string { return fi.translationBaseName }

// ContentBaseName is a either TranslationBaseName or name of containing folder
// if file is a leaf bundle.
func (fi *FileInfo) ContentBaseName() string {
	fi.init()
	return fi.contentBaseName
}

// Section returns a file's section.
func (fi *FileInfo) Section() string {
	fi.init()
	return fi.section
}

// UniqueID returns a file's unique, MD5 hash identifier.
func (fi *FileInfo) UniqueID() string {
	fi.init()
	return fi.uniqueID
}

// FileInfo returns a file's underlying os.FileInfo.
func (fi *FileInfo) FileInfo() hugofs.FileMetaInfo { return fi.fi }

func (fi *FileInfo) String() string { return fi.BaseFileName() }

// Open implements ReadableFile.
func (fi *FileInfo) Open() (hugio.ReadSeekCloser, error) {
	f, err := fi.fi.Meta().Open()

	return f, err
}

func (fi *FileInfo) IsZero() bool {
	return fi == nil
}

// We create a lot of these FileInfo objects, but there are parts of it used only
// in some cases that is slightly expensive to construct.
func (fi *FileInfo) init() {
	fi.lazyInit.Do(func() {
		relDir := strings.Trim(fi.relDir, helpers.FilePathSeparator)
		parts := strings.Split(relDir, helpers.FilePathSeparator)
		var section string
		if (fi.classifier != files.ContentClassLeaf && len(parts) == 1) || len(parts) > 1 {
			section = parts[0]
		}
		fi.section = section

		if fi.classifier.IsBundle() && len(parts) > 0 {
			fi.contentBaseName = parts[len(parts)-1]
		} else {
			fi.contentBaseName = fi.translationBaseName
		}

		fi.uniqueID = helpers.MD5String(filepath.ToSlash(fi.relPath))
	})
}

// NewTestFile creates a partially filled File used in unit tests.
// TODO(bep) improve this package
func NewTestFile(filename string) *FileInfo {
	base := filepath.Base(filepath.Dir(filename))
	return &FileInfo{
		filename:            filename,
		translationBaseName: base,
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

	filename := m.Filename
	relPath := m.Path

	if relPath == "" {
		return nil, fmt.Errorf("no Path provided by %v (%T)", m, m.Fs)
	}

	if filename == "" {
		return nil, fmt.Errorf("no Filename provided by %v (%T)", m, m.Fs)
	}

	relDir := filepath.Dir(relPath)
	if relDir == "." {
		relDir = ""
	}
	if !strings.HasSuffix(relDir, helpers.FilePathSeparator) {
		relDir = relDir + helpers.FilePathSeparator
	}

	lang := m.Lang
	translationBaseName := m.TranslationBaseName

	dir, name := filepath.Split(relPath)
	if !strings.HasSuffix(dir, helpers.FilePathSeparator) {
		dir = dir + helpers.FilePathSeparator
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	baseName := paths.Filename(name)

	if translationBaseName == "" {
		// This is usually provided by the filesystem. But this FileInfo is also
		// created in a standalone context when doing "hugo new". This is
		// an approximate implementation, which is "good enough" in that case.
		fileLangExt := filepath.Ext(baseName)
		translationBaseName = strings.TrimSuffix(baseName, fileLangExt)
	}

	f := &FileInfo{
		sp:                  sp,
		filename:            filename,
		fi:                  fi,
		lang:                lang,
		ext:                 ext,
		dir:                 dir,
		relDir:              relDir,  // Dir()
		relPath:             relPath, // Path()
		name:                name,
		baseName:            baseName, // BaseFileName()
		translationBaseName: translationBaseName,
		classifier:          m.Classifier,
	}

	return f, nil
}

func NewGitInfo(info gitmap.GitInfo) GitInfo {
	return GitInfo(info)
}

// GitInfo provides information about a version controled source file.
type GitInfo struct {
	// Commit hash.
	Hash string `json:"hash"`
	// Abbreviated commit hash.
	AbbreviatedHash string `json:"abbreviatedHash"`
	// The commit message's subject/title line.
	Subject string `json:"subject"`
	// The author name, respecting .mailmap.
	AuthorName string `json:"authorName"`
	// The author email address, respecting .mailmap.
	AuthorEmail string `json:"authorEmail"`
	// The author date.
	AuthorDate time.Time `json:"authorDate"`
	// The commit date.
	CommitDate time.Time `json:"commitDate"`
}

// IsZero returns true if the GitInfo is empty,
// meaning it will also be falsy in the Go templates.
func (g GitInfo) IsZero() bool {
	return g.Hash == ""
}
