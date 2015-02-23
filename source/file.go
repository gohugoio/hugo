// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"github.com/spf13/hugo/helpers"
	"io"
	"path/filepath"
	"strings"
)

type File struct {
	relpath     string // Original Full Path eg. /Users/Home/Hugo/foo.txt
	logicalName string // foo.txt
	Contents    io.Reader
	section     string // The first directory
	dir         string // The full directory Path (minus file name)
	ext         string // Just the ext (eg txt)
	uniqueID    string // MD5 of the filename
}

func (f *File) UniqueID() string {
	if f.uniqueID == "" {
		f.uniqueID = helpers.Md5String(f.LogicalName())
	}
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
	return helpers.Filename(f.LogicalName())
}

func (f *File) Section() string {
	if f.section != "" {
		return f.section
	}
	f.section = helpers.GuessSection(f.Dir())
	return f.section
}

func (f *File) LogicalName() string {
	if f.logicalName != "" {
		return f.logicalName
	}
	_, f.logicalName = filepath.Split(f.relpath)
	return f.logicalName
}

func (f *File) SetDir(dir string) {
	f.dir = dir
}

func (f *File) Dir() string {
	if f.dir != "" {
		return f.dir
	}
	f.dir, _ = filepath.Split(f.relpath)
	return f.dir
}

func (f *File) Extension() string {
	if f.ext != "" {
		return f.ext
	}
	f.ext = strings.TrimPrefix(filepath.Ext(f.LogicalName()), ".")
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
	return &File{
		relpath: relpath,
	}
}

func NewFileFromAbs(base, fullpath string, content io.Reader) (f *File, err error) {
	var name string
	if name, err = helpers.GetRelativePath(fullpath, base); err != nil {
		return nil, err
	}

	return NewFileWithContents(name, content), nil
}
