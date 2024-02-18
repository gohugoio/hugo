// Copyright 2024 The Hugo Authors. All rights reserved.
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
	iofs "io/fs"
	"os"
	"path"
	"runtime"
	"sort"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/spf13/afero"
	"golang.org/x/text/unicode/norm"
)

// NewComponentFs creates a new component filesystem.
func NewComponentFs(opts ComponentFsOptions) *componentFs {
	if opts.Component == "" {
		panic("ComponentFsOptions.PathParser.Component must be set")
	}
	if opts.Fs == nil {
		panic("ComponentFsOptions.Fs must be set")
	}
	bfs := NewBasePathFs(opts.Fs, opts.Component)
	return &componentFs{Fs: bfs, opts: opts}
}

var _ FilesystemUnwrapper = (*componentFs)(nil)

// componentFs is a filesystem that holds one of the Hugo components, e.g. content, layouts etc.
type componentFs struct {
	afero.Fs

	opts ComponentFsOptions
}

func (fs *componentFs) UnwrapFilesystem() afero.Fs {
	return fs.Fs
}

type componentFsDir struct {
	*noOpRegularFileOps
	DirOnlyOps
	name string // the name passed to Open
	fs   *componentFs
}

// ReadDir reads count entries from this virtual directory and
// sorts the entries according to the component filesystem rules.
func (f *componentFsDir) ReadDir(count int) ([]iofs.DirEntry, error) {
	fis, err := f.DirOnlyOps.(iofs.ReadDirFile).ReadDir(-1)
	if err != nil {
		return nil, err
	}

	// Filter out any symlinks.
	n := 0
	for _, fi := range fis {
		// IsDir will always be false for symlinks.
		keep := fi.IsDir()
		if !keep {
			// This is unfortunate, but is the only way to determine if it is a symlink.
			info, err := fi.Info()
			if err != nil {
				if herrors.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			if info.Mode()&os.ModeSymlink == 0 {
				keep = true
			}
		}
		if keep {
			fis[n] = fi
			n++
		}
	}

	fis = fis[:n]

	n = 0
	for _, fi := range fis {
		s := path.Join(f.name, fi.Name())
		if _, ok := f.fs.applyMeta(fi, s); ok {
			fis[n] = fi
			n++
		}
	}
	fis = fis[:n]

	sort.Slice(fis, func(i, j int) bool {
		fimi, fimj := fis[i].(FileMetaInfo), fis[j].(FileMetaInfo)
		if fimi.IsDir() != fimj.IsDir() {
			return fimi.IsDir()
		}
		fimim, fimjm := fimi.Meta(), fimj.Meta()

		if fimim.ModuleOrdinal != fimjm.ModuleOrdinal {
			switch f.fs.opts.Component {
			case files.ComponentFolderI18n:
				// The way the language files gets loaded means that
				// we need to provide the least important files first (e.g. the theme files).
				return fimim.ModuleOrdinal > fimjm.ModuleOrdinal
			default:
				return fimim.ModuleOrdinal < fimjm.ModuleOrdinal
			}
		}

		pii, pij := fimim.PathInfo, fimjm.PathInfo
		if pii != nil {
			basei, basej := pii.Base(), pij.Base()
			exti, extj := pii.Ext(), pij.Ext()
			if f.fs.opts.Component == files.ComponentFolderContent {
				// Pull bundles to the top.
				if pii.IsBundle() != pij.IsBundle() {
					return pii.IsBundle()
				}
			}

			if exti != extj {
				// This pulls .md above .html.
				return exti > extj
			}

			if basei != basej {
				return basei < basej
			}
		}

		if fimim.Weight != fimjm.Weight {
			return fimim.Weight > fimjm.Weight
		}

		return fimi.Name() < fimj.Name()
	})

	return fis, nil
}

func (f *componentFsDir) Stat() (iofs.FileInfo, error) {
	fi, err := f.DirOnlyOps.Stat()
	if err != nil {
		return nil, err
	}
	fim, _ := f.fs.applyMeta(fi, f.name)
	return fim, nil
}

func (fs *componentFs) Stat(name string) (os.FileInfo, error) {
	fi, err := fs.Fs.Stat(name)
	if err != nil {
		return nil, err
	}
	fim, _ := fs.applyMeta(fi, name)
	return fim, nil
}

func (fs *componentFs) applyMeta(fi FileNameIsDir, name string) (FileMetaInfo, bool) {
	if runtime.GOOS == "darwin" {
		name = norm.NFC.String(name)
	}
	fim := fi.(FileMetaInfo)
	meta := fim.Meta()
	pi := fs.opts.PathParser.Parse(fs.opts.Component, name)
	if pi.Disabled() {
		return fim, false
	}
	if meta.Lang != "" {
		if isLangDisabled := fs.opts.PathParser.IsLangDisabled; isLangDisabled != nil && isLangDisabled(meta.Lang) {
			return fim, false
		}
	}
	meta.PathInfo = pi
	if !fim.IsDir() {
		if fileLang := meta.PathInfo.Lang(); fileLang != "" {
			// A valid lang set in filename.
			// Give priority to myfile.sv.txt inside the sv filesystem.
			meta.Weight++
			meta.Lang = fileLang
		}
	}

	if meta.Lang == "" {
		meta.Lang = fs.opts.DefaultContentLanguage
	}

	langIdx, found := fs.opts.PathParser.LanguageIndex[meta.Lang]
	if !found {
		panic("no language found for " + meta.Lang)
	}
	meta.LangIndex = langIdx

	if fi.IsDir() {
		meta.OpenFunc = func() (afero.File, error) {
			return fs.Open(name)
		}
	}

	return fim, true
}

func (f *componentFsDir) Readdir(count int) ([]os.FileInfo, error) {
	panic("not supported: Use ReadDir")
}

func (f *componentFsDir) Readdirnames(count int) ([]string, error) {
	dirsi, err := f.DirOnlyOps.(iofs.ReadDirFile).ReadDir(count)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, len(dirsi))
	for i, d := range dirsi {
		dirs[i] = d.Name()
	}
	return dirs, nil
}

type ComponentFsOptions struct {
	// The filesystem where one or more components are mounted.
	Fs afero.Fs

	// The component name, e.g. "content", "layouts" etc.
	Component string

	DefaultContentLanguage string

	// The parser used to parse paths provided by this filesystem.
	PathParser *paths.PathParser
}

func (fs *componentFs) Open(name string) (afero.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		if err != errIsDir {
			f.Close()
			return nil, err
		}
	} else if !fi.IsDir() {
		return f, nil
	}

	return &componentFsDir{
		DirOnlyOps: f,
		name:       name,
		fs:         fs,
	}, nil
}

func (fs *componentFs) ReadDir(name string) ([]os.FileInfo, error) {
	panic("not implemented")
}
