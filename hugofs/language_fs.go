// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

const hugoFsMarker = "__hugofs"

var (
	_ LanguageAnnouncer = (*LanguageFileInfo)(nil)
	_ FilePather        = (*LanguageFileInfo)(nil)
	_ afero.Lstater     = (*LanguageFs)(nil)
)

// LanguageAnnouncer is aware of its language.
type LanguageAnnouncer interface {
	Lang() string
	TranslationBaseName() string
}

// FilePather is aware of its file's location.
type FilePather interface {
	// Filename gets the full path and filename to the file.
	Filename() string

	// Path gets the content relative path including file name and extension.
	// The directory is relative to the content root where "content" is a broad term.
	Path() string

	// RealName is FileInfo.Name in its original form.
	RealName() string

	BaseDir() string
}

// LanguageDirsMerger implements the afero.DirsMerger interface, which is used
// to merge two directories.
var LanguageDirsMerger = func(lofi, bofi []os.FileInfo) ([]os.FileInfo, error) {
	m := make(map[string]*LanguageFileInfo)

	for _, fi := range lofi {
		fil, ok := fi.(*LanguageFileInfo)
		if !ok {
			return nil, fmt.Errorf("received %T, expected *LanguageFileInfo", fi)
		}
		m[fil.virtualName] = fil
	}

	for _, fi := range bofi {
		fil, ok := fi.(*LanguageFileInfo)
		if !ok {
			return nil, fmt.Errorf("received %T, expected *LanguageFileInfo", fi)
		}
		existing, found := m[fil.virtualName]

		if !found || existing.weight < fil.weight {
			m[fil.virtualName] = fil
		}
	}

	merged := make([]os.FileInfo, len(m))
	i := 0
	for _, v := range m {
		merged[i] = v
		i++
	}

	return merged, nil
}

// LanguageFileInfo is a super-set of os.FileInfo with additional information
// about the file in relation to its Hugo language.
type LanguageFileInfo struct {
	os.FileInfo
	lang                string
	baseDir             string
	realFilename        string
	relFilename         string
	name                string
	realName            string
	virtualName         string
	translationBaseName string

	// We add some weight to the files in their own language's content directory.
	weight int
}

// Filename returns a file's real filename including the base (ie.
// "/my/base/sect/page.md").
func (fi *LanguageFileInfo) Filename() string {
	return fi.realFilename
}

// Path returns a file's filename relative to the base (ie. "sect/page.md").
func (fi *LanguageFileInfo) Path() string {
	return fi.relFilename
}

// RealName returns a file's real base name (ie. "page.md").
func (fi *LanguageFileInfo) RealName() string {
	return fi.realName
}

// BaseDir returns a file's base directory (ie. "/my/base").
func (fi *LanguageFileInfo) BaseDir() string {
	return fi.baseDir
}

// Lang returns a file's language (ie. "sv").
func (fi *LanguageFileInfo) Lang() string {
	return fi.lang
}

// TranslationBaseName returns the base filename without any extension or language
// identifiers (ie. "page").
func (fi *LanguageFileInfo) TranslationBaseName() string {
	return fi.translationBaseName
}

// Name is the name of the file within this filesystem without any path info.
// It will be marked with language information so we can identify it as ours
// (ie. "__hugofs_sv_page.md").
func (fi *LanguageFileInfo) Name() string {
	return fi.name
}

type languageFile struct {
	afero.File
	fs *LanguageFs
}

// Readdir creates FileInfo entries by calling Lstat if possible.
func (l *languageFile) Readdir(c int) (ofi []os.FileInfo, err error) {
	names, err := l.File.Readdirnames(c)
	if err != nil {
		return nil, err
	}

	fis := make([]os.FileInfo, len(names))

	for i, name := range names {
		fi, _, err := l.fs.LstatIfPossible(filepath.Join(l.Name(), name))

		if err != nil {
			return nil, err
		}
		fis[i] = fi
	}

	return fis, err
}

// LanguageFs represents a language filesystem.
type LanguageFs struct {
	// This Fs is usually created with a BasePathFs
	basePath   string
	lang       string
	nameMarker string
	languages  map[string]bool
	afero.Fs
}

// NewLanguageFs creates a new language filesystem.
func NewLanguageFs(lang string, languages map[string]bool, fs afero.Fs) *LanguageFs {
	if lang == "" {
		panic("no lang set for the language fs")
	}
	var basePath string

	if bfs, ok := fs.(*afero.BasePathFs); ok {
		basePath, _ = bfs.RealPath("")
	}

	marker := hugoFsMarker + "_" + lang + "_"

	return &LanguageFs{lang: lang, languages: languages, basePath: basePath, Fs: fs, nameMarker: marker}
}

// Lang returns a language filesystem's language (ie. "sv").
func (fs *LanguageFs) Lang() string {
	return fs.lang
}

// Stat returns the os.FileInfo of a given file.
func (fs *LanguageFs) Stat(name string) (os.FileInfo, error) {
	name, err := fs.realName(name)
	if err != nil {
		return nil, err
	}

	fi, err := fs.Fs.Stat(name)
	if err != nil {
		return nil, err
	}

	return fs.newLanguageFileInfo(name, fi)
}

// Open opens the named file for reading.
func (fs *LanguageFs) Open(name string) (afero.File, error) {
	name, err := fs.realName(name)
	if err != nil {
		return nil, err
	}
	f, err := fs.Fs.Open(name)

	if err != nil {
		return nil, err
	}
	return &languageFile{File: f, fs: fs}, nil
}

// LstatIfPossible returns the os.FileInfo structure describing a given file.
// It attempts to use Lstat if supported or defers to the os.  In addition to
// the FileInfo, a boolean is returned telling whether Lstat was called.
func (fs *LanguageFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	name, err := fs.realName(name)
	if err != nil {
		return nil, false, err
	}

	var fi os.FileInfo
	var b bool

	if lif, ok := fs.Fs.(afero.Lstater); ok {
		fi, b, err = lif.LstatIfPossible(name)
	} else {
		fi, err = fs.Fs.Stat(name)
	}

	if err != nil {
		return nil, b, err
	}

	lfi, err := fs.newLanguageFileInfo(name, fi)

	return lfi, b, err
}

func (fs *LanguageFs) realPath(name string) (string, error) {
	if baseFs, ok := fs.Fs.(*afero.BasePathFs); ok {
		return baseFs.RealPath(name)
	}
	return name, nil
}

func (fs *LanguageFs) realName(name string) (string, error) {
	if strings.Contains(name, hugoFsMarker) {
		if !strings.Contains(name, fs.nameMarker) {
			return "", os.ErrNotExist
		}
		return strings.Replace(name, fs.nameMarker, "", 1), nil
	}

	if fs.basePath == "" {
		return name, nil
	}

	return strings.TrimPrefix(name, fs.basePath), nil
}

func (fs *LanguageFs) newLanguageFileInfo(filename string, fi os.FileInfo) (*LanguageFileInfo, error) {
	filename = filepath.Clean(filename)
	_, name := filepath.Split(filename)

	realName := name
	virtualName := name

	realPath, err := fs.realPath(filename)
	if err != nil {
		return nil, err
	}

	lang := fs.Lang()

	baseNameNoExt := ""

	if !fi.IsDir() {

		// Try to extract the language from the file name.
		// Any valid language identificator in the name will win over the
		// language set on the file system, e.g. "mypost.en.md".
		baseName := filepath.Base(name)
		ext := filepath.Ext(baseName)
		baseNameNoExt = baseName

		if ext != "" {
			baseNameNoExt = strings.TrimSuffix(baseNameNoExt, ext)
		}

		fileLangExt := filepath.Ext(baseNameNoExt)
		fileLang := strings.TrimPrefix(fileLangExt, ".")

		if fs.languages[fileLang] {
			lang = fileLang
			baseNameNoExt = strings.TrimSuffix(baseNameNoExt, fileLangExt)
		}

		// This connects the filename to the filesystem, not the language.
		virtualName = baseNameNoExt + "." + lang + ext

		name = fs.nameMarker + name
	}

	weight := 1
	// If this file's language belongs in this directory, add some weight to it
	// to make it more important.
	if lang == fs.Lang() {
		weight = 2
	}

	if fi.IsDir() {
		// For directories we always want to start from the union view.
		realPath = strings.TrimPrefix(realPath, fs.basePath)
	}

	return &LanguageFileInfo{
		lang:                lang,
		weight:              weight,
		realFilename:        realPath,
		realName:            realName,
		relFilename:         strings.TrimPrefix(strings.TrimPrefix(realPath, fs.basePath), string(os.PathSeparator)),
		name:                name,
		virtualName:         virtualName,
		translationBaseName: baseNameNoExt,
		baseDir:             fs.basePath,
		FileInfo:            fi}, nil
}
