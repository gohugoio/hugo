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

package hugofs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/spf13/afero"
)

var (
	_ afero.Fs      = (*FilterFs)(nil)
	_ afero.Lstater = (*FilterFs)(nil)
	_ afero.File    = (*filterDir)(nil)
)

func NewLanguageFs(langs map[string]int, fs afero.Fs) (afero.Fs, error) {

	applyMeta := func(fs *FilterFs, name string, fis []os.FileInfo) {

		for i, fi := range fis {
			if fi.IsDir() {
				filename := filepath.Join(name, fi.Name())
				fis[i] = decorateFileInfo(fi, fs, fs.getOpener(filename), "", "", nil)
				continue
			}

			meta := fi.(FileMetaInfo).Meta()
			lang := meta.Lang()

			fileLang, translationBaseName, translationBaseNameWithExt := langInfoFrom(langs, fi.Name())
			weight := 0

			if fileLang != "" {
				weight = 1
				if fileLang == lang {
					// Give priority to myfile.sv.txt inside the sv filesystem.
					weight++
				}
				lang = fileLang
			}

			fim := NewFileMetaInfo(fi, FileMeta{
				metaKeyLang:                       lang,
				metaKeyWeight:                     weight,
				metaKeyOrdinal:                    langs[lang],
				metaKeyTranslationBaseName:        translationBaseName,
				metaKeyTranslationBaseNameWithExt: translationBaseNameWithExt,
				metaKeyClassifier:                 files.ClassifyContentFile(fi.Name(), meta.GetOpener()),
			})

			fis[i] = fim
		}
	}

	all := func(fis []os.FileInfo) {
		// Maps translation base name to a list of language codes.
		translations := make(map[string][]string)
		trackTranslation := func(meta FileMeta) {
			name := meta.TranslationBaseNameWithExt()
			translations[name] = append(translations[name], meta.Lang())
		}
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			meta := fi.(FileMetaInfo).Meta()

			trackTranslation(meta)

		}

		for _, fi := range fis {
			fim := fi.(FileMetaInfo)
			langs := translations[fim.Meta().TranslationBaseNameWithExt()]
			if len(langs) > 0 {
				fim.Meta()["translations"] = sortAndremoveStringDuplicates(langs)
			}
		}
	}

	return &FilterFs{
		fs:             fs,
		applyPerSource: applyMeta,
		applyAll:       all,
	}, nil

}

func NewFilterFs(fs afero.Fs) (afero.Fs, error) {

	applyMeta := func(fs *FilterFs, name string, fis []os.FileInfo) {
		for i, fi := range fis {
			if fi.IsDir() {
				fis[i] = decorateFileInfo(fi, fs, fs.getOpener(fi.(FileMetaInfo).Meta().Filename()), "", "", nil)
			}
		}
	}

	ffs := &FilterFs{
		fs:             fs,
		applyPerSource: applyMeta,
	}

	return ffs, nil

}

// FilterFs is an ordered composite filesystem.
type FilterFs struct {
	fs afero.Fs

	applyPerSource func(fs *FilterFs, name string, fis []os.FileInfo)
	applyAll       func(fis []os.FileInfo)
}

func (fs *FilterFs) Chmod(n string, m os.FileMode) error {
	return syscall.EPERM
}

func (fs *FilterFs) Chtimes(n string, a, m time.Time) error {
	return syscall.EPERM
}

func (fs *FilterFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	fi, b, err := lstatIfPossible(fs.fs, name)

	if err != nil {
		return nil, false, err
	}

	if fi.IsDir() {
		return decorateFileInfo(fi, fs, fs.getOpener(name), "", "", nil), false, nil
	}

	parent := filepath.Dir(name)
	fs.applyFilters(parent, -1, fi)

	return fi, b, nil

}

func (fs *FilterFs) Mkdir(n string, p os.FileMode) error {
	return syscall.EPERM
}

func (fs *FilterFs) MkdirAll(n string, p os.FileMode) error {
	return syscall.EPERM
}

func (fs *FilterFs) Name() string {
	return "WeightedFileSystem"
}

func (fs *FilterFs) Open(name string) (afero.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	return &filterDir{
		File: f,
		ffs:  fs,
	}, nil

}

func (fs *FilterFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return fs.fs.Open(name)
}

func (fs *FilterFs) ReadDir(name string) ([]os.FileInfo, error) {
	panic("not implemented")
}

func (fs *FilterFs) Remove(n string) error {
	return syscall.EPERM
}

func (fs *FilterFs) RemoveAll(p string) error {
	return syscall.EPERM
}

func (fs *FilterFs) Rename(o, n string) error {
	return syscall.EPERM
}

func (fs *FilterFs) Stat(name string) (os.FileInfo, error) {
	fi, _, err := fs.LstatIfPossible(name)
	return fi, err
}

func (fs *FilterFs) Create(n string) (afero.File, error) {
	return nil, syscall.EPERM
}

func (fs *FilterFs) getOpener(name string) func() (afero.File, error) {
	return func() (afero.File, error) {
		return fs.Open(name)
	}
}

func (fs *FilterFs) applyFilters(name string, count int, fis ...os.FileInfo) ([]os.FileInfo, error) {
	if fs.applyPerSource != nil {
		fs.applyPerSource(fs, name, fis)
	}

	seen := make(map[string]bool)
	var duplicates []int
	for i, dir := range fis {
		if !dir.IsDir() {
			continue
		}
		if seen[dir.Name()] {
			duplicates = append(duplicates, i)
		} else {
			seen[dir.Name()] = true
		}
	}

	// Remove duplicate directories, keep first.
	if len(duplicates) > 0 {
		for i := len(duplicates) - 1; i >= 0; i-- {
			idx := duplicates[i]
			fis = append(fis[:idx], fis[idx+1:]...)
		}
	}

	if fs.applyAll != nil {
		fs.applyAll(fis)
	}

	if count > 0 && len(fis) >= count {
		return fis[:count], nil
	}

	return fis, nil

}

type filterDir struct {
	afero.File
	ffs *FilterFs
}

func (f *filterDir) Readdir(count int) ([]os.FileInfo, error) {
	fis, err := f.File.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return f.ffs.applyFilters(f.Name(), count, fis...)
}

func (f *filterDir) Readdirnames(count int) ([]string, error) {
	dirsi, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, len(dirsi))
	for i, d := range dirsi {
		dirs[i] = d.Name()
	}
	return dirs, nil
}

// Try to extract the language from the given filename.
// Any valid language identificator in the name will win over the
// language set on the file system, e.g. "mypost.en.md".
func langInfoFrom(languages map[string]int, name string) (string, string, string) {
	var lang string

	baseName := filepath.Base(name)
	ext := filepath.Ext(baseName)
	translationBaseName := baseName

	if ext != "" {
		translationBaseName = strings.TrimSuffix(translationBaseName, ext)
	}

	fileLangExt := filepath.Ext(translationBaseName)
	fileLang := strings.TrimPrefix(fileLangExt, ".")

	if _, found := languages[fileLang]; found {
		lang = fileLang
		translationBaseName = strings.TrimSuffix(translationBaseName, fileLangExt)
	}

	translationBaseNameWithExt := translationBaseName

	if ext != "" {
		translationBaseNameWithExt += ext
	}

	return lang, translationBaseName, translationBaseNameWithExt

}

func printFs(fs afero.Fs, path string, w io.Writer) {
	if fs == nil {
		return
	}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		fmt.Println("p:::", path)
		return nil
	})
}

func sortAndremoveStringDuplicates(s []string) []string {
	ss := sort.StringSlice(s)
	ss.Sort()
	i := 0
	for j := 1; j < len(s); j++ {
		if !ss.Less(i, j) {
			continue
		}
		i++
		s[i] = s[j]
	}

	return s[:i+1]
}
