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
	"context"
	iofs "io/fs"
	"os"
	"path"
	"runtime"
	"sort"

	"github.com/bep/helpers/contexthelpers"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
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

var (
	_ FilesystemUnwrapper   = (*componentFs)(nil)
	_ ReadDirWithContextDir = (*componentFsDir)(nil)
)

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

type (
	contextKey uint8
)

const (
	contextKeyIsInLeafBundle contextKey = iota
)

var componentFsContext = struct {
	IsInLeafBundle contexthelpers.ContextDispatcher[bool]
}{
	IsInLeafBundle: contexthelpers.NewContextDispatcher[bool](contextKeyIsInLeafBundle),
}

func (f *componentFsDir) ReadDirWithContext(ctx context.Context, count int) ([]iofs.DirEntry, context.Context, error) {
	fis, err := f.DirOnlyOps.(iofs.ReadDirFile).ReadDir(-1)
	if err != nil {
		return nil, ctx, err
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
				return nil, ctx, err
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
	n = 0

	sort.Slice(fis, func(i, j int) bool {
		fimi, fimj := fis[i].(FileMetaInfo), fis[j].(FileMetaInfo)
		if fimi.IsDir() != fimj.IsDir() {
			return fimi.IsDir()
		}
		fimim, fimjm := fimi.Meta(), fimj.Meta()

		bi, bj := fimim.PathInfo.Base(), fimjm.PathInfo.Base()
		if bi == bj {
			matrixi, matrixj := fimim.SitesMatrix, fimjm.SitesMatrix
			l1, l2 := matrixi.LenVectors(), matrixj.LenVectors()
			if l1 != l2 {
				// Pull the ones with the least number of sites defined to the top.
				return l1 < l2
			}
		}

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

	if f.fs.opts.Component == files.ComponentFolderContent {
		isInLeafBundle := componentFsContext.IsInLeafBundle.Get(ctx)
		var isCurrentLeafBundle bool
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}

			pi := fi.(FileMetaInfo).Meta().PathInfo

			if pi.IsLeafBundle() {
				isCurrentLeafBundle = true
				break
			}
		}

		if !isInLeafBundle && isCurrentLeafBundle {
			ctx = componentFsContext.IsInLeafBundle.Set(ctx, true)
		}

		if isInLeafBundle || isCurrentLeafBundle {
			for _, fi := range fis {
				if fi.IsDir() {
					continue
				}
				pi := fi.(FileMetaInfo).Meta().PathInfo

				// Everything below a leaf bundle is a resource.
				isResource := isInLeafBundle && pi.Type() > paths.TypeFile
				// Every sibling of a leaf bundle is a resource.
				isResource = isResource || (isCurrentLeafBundle && !pi.IsLeafBundle())

				if isResource {
					paths.ModifyPathBundleTypeResource(pi)
				}
			}
		}

	}

	type typeBase struct {
		Type paths.Type
		Base string
	}

	variants := make(map[typeBase][]sitesmatrix.VectorProvider)

	for _, fi := range fis {

		if !fi.IsDir() {
			meta := fi.(FileMetaInfo).Meta()

			pi := meta.PathInfo

			if pi.Component() == files.ComponentFolderLayouts || pi.Component() == files.ComponentFolderContent {

				var base string
				switch pi.Component() {
				case files.ComponentFolderContent:
					base = pi.Base() + pi.Custom()
				default:
					base = pi.PathNoLang()
				}

				baseName := typeBase{pi.Type(), base}

				// There may be multiple languge/version/role combinations for the same file.
				// The most important come early.
				matrixes, found := variants[baseName]

				if found {
					complement := meta.SitesMatrix.Complement(matrixes...)
					if complement == nil || complement.LenVectors() == 0 {
						continue
					}
					matrixes = append(matrixes, meta.SitesMatrix)
					meta.SitesMatrix = complement

					variants[baseName] = matrixes

				} else {
					matrixes = []sitesmatrix.VectorProvider{meta.SitesMatrix}
					variants[baseName] = matrixes

				}
			}
		}

		fis[n] = fi
		n++

	}
	fis = fis[:n]

	return fis, ctx, nil
}

// ReadDir reads count entries from this virtual directory and
// sorts the entries according to the component filesystem rules.
func (f *componentFsDir) ReadDir(count int) ([]iofs.DirEntry, error) {
	v, _, err := ReadDirWithContext(context.Background(), f.DirOnlyOps, count)
	return v, err
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
	pi := fs.opts.Cfg.PathParser().Parse(fs.opts.Component, name)
	if pi.Disabled() {
		return fim, false
	}

	meta.PathInfo = pi
	if !fim.IsDir() {
		if fileLang := meta.PathInfo.Lang(); fileLang != "" {
			if idx, ok := fs.opts.Cfg.PathParser().LanguageIndex[fileLang]; ok {
				// A valid lang set in filename.
				// Give priority to myfile.sv.txt inside the sv filesystem.
				meta.Weight++
				meta.SitesMatrix = meta.SitesMatrix.WithLanguageIndices(idx)
				if idx > 0 {
					// Not the default language, add some weight.
					meta.SitesMatrix = sitesmatrix.NewWeightedVectorStore(meta.SitesMatrix, 10)
				}

			}
		}
		switch meta.Component {
		case files.ComponentFolderLayouts:
			// Eg. index.fr.html when French isn't defined,
			// we want e.g. index.html to be used instead.
			if len(pi.IdentifiersUnknown()) > 0 {
				meta.Weight--
			}
		}

	}

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

	Cfg config.AllProvider
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
