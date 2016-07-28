// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/viper"
)

// File represents a source content file.
// All paths are relative from the source directory base
type File struct {
	relpath     string // Original relative path, e.g. content/foo.txt
	logicalName string // foo.txt
	baseName    string // `post` for `post.md`, also `post.en` for `post.en.md`
	Contents    io.Reader
	section     string // The first directory
	dir         string // The relative directory Path (minus file name)
	ext         string // Just the ext (eg txt)
	uniqueID    string // MD5 of the filename

	translationBaseName string // `post` for `post.es.md` (if `Multilingual` is enabled.)
	lang                string // The language code if `Multilingual` is enabled
}

// UniqueID is the MD5 hash of the filename and is for most practical applications,
// Hugo content files being one of them, considered to be unique.
func (f *File) UniqueID() string {
	return f.uniqueID
}

// String returns the file's content as a string.
func (f *File) String() string {
	return helpers.ReaderToString(f.Contents)
}

// Bytes returns the file's content as a byte slice.
func (f *File) Bytes() []byte {
	return helpers.ReaderToBytes(f.Contents)
}

// BaseFileName Filename without extension.
func (f *File) BaseFileName() string {
	return f.baseName
}

// Filename with no extension, not even the optional language extension part.
func (f *File) TranslationBaseName() string {
	return f.translationBaseName
}

// Lang for this page, if `Multilingual` is enabled on your site.
func (f *File) Lang() string {
	return f.lang
}

// Section is first directory below the content root.
func (f *File) Section() string {
	return f.section
}

// LogicalName is filename and extension of the file.
func (f *File) LogicalName() string {
	return f.logicalName
}

// SetDir sets the relative directory where this file lives.
// TODO(bep) Get rid of this.
func (f *File) SetDir(dir string) {
	f.dir = dir
}

// Dir gets the name of the directory that contains this file.
// The directory is relative to the content root.
func (f *File) Dir() string {
	return f.dir
}

// Extension gets the file extension, i.e "myblogpost.md" will return "md".
func (f *File) Extension() string {
	return f.ext
}

// Ext is an alias for Extension.
func (f *File) Ext() string {
	return f.Extension()
}

// Path gets the relative path including file name and extension from
// the base of the source directory.
func (f *File) Path() string {
	return f.relpath
}

// NewFileWithContents creates a new File pointer with the given relative path and
// content. The language defaults to "en".
func NewFileWithContents(relpath string, content io.Reader) *File {
	file := NewFile(relpath)
	file.Contents = content
	file.lang = "en"
	return file
}

// NewFile creates a new File pointer with the given relative path.
func NewFile(relpath string) *File {
	f := &File{
		relpath: relpath,
	}

	f.dir, f.logicalName = filepath.Split(f.relpath)
	f.ext = strings.TrimPrefix(filepath.Ext(f.LogicalName()), ".")
	f.baseName = helpers.Filename(f.LogicalName())

	f.lang = strings.TrimPrefix(filepath.Ext(f.baseName), ".")
	if f.lang == "" {
		f.lang = viper.GetString("DefaultContentLanguage")
		if f.lang == "" {
			// TODO(bep) ml
			f.lang = "en"
		}
	}
	f.translationBaseName = helpers.Filename(f.baseName)

	f.section = helpers.GuessSection(f.Dir())
	f.uniqueID = helpers.Md5String(f.LogicalName())

	return f
}

// NewFileFromAbs creates a new File pointer with the given full file path path and
// content.
func NewFileFromAbs(base, fullpath string, content io.Reader) (f *File, err error) {
	var name string
	if name, err = helpers.GetRelativePath(fullpath, base); err != nil {
		return nil, err
	}

	return NewFileWithContents(name, content), nil
}
