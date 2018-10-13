// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/helpers"
)

// fileInfo implements the File interface.
var (
	_ File         = (*FileInfo)(nil)
	_ ReadableFile = (*FileInfo)(nil)
)

// File represents a source file.
type File interface {

	// Filename gets the full path and filename to the file.
	Filename() string

	// Path gets the relative path including file name and extension.
	// The directory is relative to the content root.
	Path() string

	// Dir gets the name of the directory that contains this file.
	// The directory is relative to the content root.
	Dir() string

	// Extension gets the file extension, i.e "myblogpost.md" will return "md".
	Extension() string
	// Ext is an alias for Extension.
	Ext() string // Hmm... Deprecate Extension

	// Lang for this page, if `Multilingual` is enabled on your site.
	Lang() string

	// LogicalName is filename and extension of the file.
	LogicalName() string

	// Section is first directory below the content root.
	// For page bundles in root, the Section will be empty.
	Section() string

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

	FileInfo() os.FileInfo

	String() string
}

// A ReadableFile is a File that is readable.
type ReadableFile interface {
	File
	Open() (io.ReadCloser, error)
}

// FileInfo describes a source file.
type FileInfo struct {

	// Absolute filename to the file on disk.
	filename string

	sp *SourceSpec

	fi os.FileInfo

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
	isLeafBundle        bool

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
func (fi *FileInfo) Extension() string { return fi.Ext() }

// Ext returns a file's extension without the leading period (ie. "md").
func (fi *FileInfo) Ext() string { return fi.ext }

// Lang returns a file's language (ie. "sv").
func (fi *FileInfo) Lang() string { return fi.lang }

// LogicalName returns a file's name and extension (ie. "page.sv.md").
func (fi *FileInfo) LogicalName() string { return fi.name }

// BaseFileName returns a file's name without extension (ie. "page.sv").
func (fi *FileInfo) BaseFileName() string { return fi.baseName }

// TranslationBaseName returns a file's translation base name without the
// language segement (ie. "page").
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
func (fi *FileInfo) FileInfo() os.FileInfo { return fi.fi }

func (fi *FileInfo) String() string { return fi.BaseFileName() }

// Open implements ReadableFile.
func (fi *FileInfo) Open() (io.ReadCloser, error) {
	f, err := fi.sp.SourceFs.Open(fi.Filename())
	return f, err
}

// We create a lot of these FileInfo objects, but there are parts of it used only
// in some cases that is slightly expensive to construct.
func (fi *FileInfo) init() {
	fi.lazyInit.Do(func() {
		relDir := strings.Trim(fi.relDir, helpers.FilePathSeparator)
		parts := strings.Split(relDir, helpers.FilePathSeparator)
		var section string
		if (!fi.isLeafBundle && len(parts) == 1) || len(parts) > 1 {
			section = parts[0]
		}
		fi.section = section

		if fi.isLeafBundle && len(parts) > 0 {
			fi.contentBaseName = parts[len(parts)-1]
		} else {
			fi.contentBaseName = fi.translationBaseName
		}

		fi.uniqueID = helpers.MD5String(filepath.ToSlash(fi.relPath))
	})
}

// NewFileInfo returns a new FileInfo structure.
func (sp *SourceSpec) NewFileInfo(baseDir, filename string, isLeafBundle bool, fi os.FileInfo) *FileInfo {

	var lang, translationBaseName, relPath string

	if fp, ok := fi.(hugofs.FilePather); ok {
		filename = fp.Filename()
		baseDir = fp.BaseDir()
		relPath = fp.Path()
	}

	if fl, ok := fi.(hugofs.LanguageAnnouncer); ok {
		lang = fl.Lang()
		translationBaseName = fl.TranslationBaseName()
	}

	dir, name := filepath.Split(filename)
	if !strings.HasSuffix(dir, helpers.FilePathSeparator) {
		dir = dir + helpers.FilePathSeparator
	}

	baseDir = strings.TrimSuffix(baseDir, helpers.FilePathSeparator)

	relDir := ""
	if dir != baseDir {
		relDir = strings.TrimPrefix(dir, baseDir)
	}

	relDir = strings.TrimPrefix(relDir, helpers.FilePathSeparator)

	if relPath == "" {
		relPath = filepath.Join(relDir, name)
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	baseName := helpers.Filename(name)

	if translationBaseName == "" {
		// This is usyally provided by the filesystem. But this FileInfo is also
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
		relDir:              relDir,
		relPath:             relPath,
		name:                name,
		baseName:            baseName,
		translationBaseName: translationBaseName,
		isLeafBundle:        isLeafBundle,
	}

	return f

}

func printFs(fs afero.Fs, path string, w io.Writer) {
	if fs == nil {
		return
	}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {

		if info != nil && !info.IsDir() {

			s := path
			if lang, ok := info.(hugofs.LanguageAnnouncer); ok {
				s = s + "\t" + lang.Lang()
			}
			if fp, ok := info.(hugofs.FilePather); ok {
				s = s + "\t" + fp.Filename()
			}
			fmt.Fprintln(w, "    ", s)
		}
		return nil
	})
}
