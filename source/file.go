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

type File struct {
	relpath     string // Original Full Path eg. /Users/Home/Hugo/foo.txt
	logicalName string // foo.txt
	baseName    string // `post` for `post.md`, also `post.en` for `post.en.md`
	Contents    io.Reader
	section     string // The first directory
	dir         string // The full directory Path (minus file name)
	ext         string // Just the ext (eg txt)
	uniqueID    string // MD5 of the filename

	translationBaseName string // `post` for `post.es.md` (if `Multilingual` is enabled.)
	lang                string // The language code if `Multilingual` is enabled
}

func (f *File) UniqueID() string {
	return f.uniqueID
}

func (f *File) String() string {
	return helpers.ReaderToString(f.Contents)
}

func (f *File) Bytes() []byte {
	return helpers.ReaderToBytes(f.Contents)
}

// Filename without extension
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

func (f *File) Section() string {
	return f.section
}

func (f *File) LogicalName() string {
	return f.logicalName
}

func (f *File) SetDir(dir string) {
	f.dir = dir
}

func (f *File) Dir() string {
	return f.dir
}

func (f *File) Extension() string {
	return f.ext
}

func (f *File) Ext() string {
	return f.Extension()
}

func (f *File) Path() string {
	return f.relpath
}

func NewFileWithContents(relpath string, content io.Reader) *File {
	file := NewFile(relpath)
	file.Contents = content
	return file
}

func NewFile(relpath string) *File {
	f := &File{
		relpath: relpath,
	}

	f.dir, _ = filepath.Split(f.relpath)
	_, f.logicalName = filepath.Split(f.relpath)
	f.ext = strings.TrimPrefix(filepath.Ext(f.LogicalName()), ".")
	f.baseName = helpers.Filename(f.LogicalName())
	if viper.GetBool("Multilingual") {
		f.lang = strings.TrimPrefix(filepath.Ext(f.baseName), ".")
		if f.lang == "" {
			f.lang = viper.GetString("DefaultContentLang")
		}
		f.translationBaseName = helpers.Filename(f.baseName)
	} else {
		f.translationBaseName = f.baseName
	}

	f.section = helpers.GuessSection(f.Dir())
	f.uniqueID = helpers.Md5String(f.LogicalName())

	return f
}

func NewFileFromAbs(base, fullpath string, content io.Reader) (f *File, err error) {
	var name string
	if name, err = helpers.GetRelativePath(fullpath, base); err != nil {
		return nil, err
	}

	return NewFileWithContents(name, content), nil
}
