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
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/pkg/errors"

	"github.com/spf13/afero"
)

type (
	WalkFunc func(path string, info FileMetaInfo, err error) error
	WalkHook func(dir FileMetaInfo, path string, readdir []FileMetaInfo) ([]FileMetaInfo, error)
)

type Walkway struct {
	fs       afero.Fs
	root     string
	basePath string

	logger *loggers.Logger

	// May be pre-set
	fi         FileMetaInfo
	dirEntries []FileMetaInfo

	walkFn WalkFunc
	walked bool

	// We may traverse symbolic links and bite ourself.
	seen map[string]bool

	// Optional hooks
	hookPre  WalkHook
	hookPost WalkHook
}

type WalkwayConfig struct {
	Fs       afero.Fs
	Root     string
	BasePath string

	Logger *loggers.Logger

	// One or both of these may be pre-set.
	Info       FileMetaInfo
	DirEntries []FileMetaInfo

	WalkFn   WalkFunc
	HookPre  WalkHook
	HookPost WalkHook
}

func NewWalkway(cfg WalkwayConfig) *Walkway {
	var fs afero.Fs
	if cfg.Info != nil {
		fs = cfg.Info.Meta().Fs()
	} else {
		fs = cfg.Fs
	}

	basePath := cfg.BasePath
	if basePath != "" && !strings.HasSuffix(basePath, filepathSeparator) {
		basePath += filepathSeparator
	}

	logger := cfg.Logger
	if logger == nil {
		logger = loggers.NewWarningLogger()
	}

	return &Walkway{
		fs:         fs,
		root:       cfg.Root,
		basePath:   basePath,
		fi:         cfg.Info,
		dirEntries: cfg.DirEntries,
		walkFn:     cfg.WalkFn,
		hookPre:    cfg.HookPre,
		hookPost:   cfg.HookPost,
		logger:     logger,
		seen:       make(map[string]bool)}
}

func (w *Walkway) Walk() error {
	if w.walked {
		panic("this walkway is already walked")
	}
	w.walked = true

	if w.fs == NoOpFs {
		return nil
	}

	var fi FileMetaInfo
	if w.fi != nil {
		fi = w.fi
	} else {
		info, _, err := lstatIfPossible(w.fs, w.root)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}

			if w.checkErr(w.root, err) {
				return nil
			}
			return w.walkFn(w.root, nil, errors.Wrapf(err, "walk: %q", w.root))
		}
		fi = info.(FileMetaInfo)
	}

	if !fi.IsDir() {
		return w.walkFn(w.root, nil, errors.New("file to walk must be a directory"))
	}

	return w.walk(w.root, fi, w.dirEntries, w.walkFn)

}

// if the filesystem supports it, use Lstat, else use fs.Stat
func lstatIfPossible(fs afero.Fs, path string) (os.FileInfo, bool, error) {
	if lfs, ok := fs.(afero.Lstater); ok {
		fi, b, err := lfs.LstatIfPossible(path)
		return fi, b, err
	}
	fi, err := fs.Stat(path)
	return fi, false, err
}

// checkErr returns true if the error is handled.
func (w *Walkway) checkErr(filename string, err error) bool {
	if err == ErrPermissionSymlink {
		logUnsupportedSymlink(filename, w.logger)
		return true
	}

	if os.IsNotExist(err) {
		// The file may be removed in process.
		// This may be a ERROR situation, but it is not possible
		// to determine as a general case.
		w.logger.WARN.Printf("File %q not found, skipping.", filename)
		return true
	}

	return false
}

func logUnsupportedSymlink(filename string, logger *loggers.Logger) {
	logger.WARN.Printf("Unsupported symlink found in %q, skipping.", filename)
}

// walk recursively descends path, calling walkFn.
// It follow symlinks if supported by the filesystem, but only the same path once.
func (w *Walkway) walk(path string, info FileMetaInfo, dirEntries []FileMetaInfo, walkFn WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return nil
	}

	meta := info.Meta()
	filename := meta.Filename()

	if dirEntries == nil {
		f, err := w.fs.Open(path)
		if err != nil {
			if w.checkErr(path, err) {
				return nil
			}
			return walkFn(path, info, errors.Wrapf(err, "walk: open %q (%q)", path, w.root))
		}

		fis, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			if w.checkErr(filename, err) {
				return nil
			}
			return walkFn(path, info, errors.Wrap(err, "walk: Readdir"))
		}

		dirEntries = fileInfosToFileMetaInfos(fis)

		if !meta.IsOrdered() {
			sort.Slice(dirEntries, func(i, j int) bool {
				fii := dirEntries[i]
				fij := dirEntries[j]

				fim, fjm := fii.Meta(), fij.Meta()

				// Pull bundle headers to the top.
				ficlass, fjclass := fim.Classifier(), fjm.Classifier()
				if ficlass != fjclass {
					return ficlass < fjclass
				}

				// With multiple content dirs with different languages,
				// there can be duplicate files, and a weight will be added
				// to the closest one.
				fiw, fjw := fim.Weight(), fjm.Weight()
				if fiw != fjw {
					return fiw > fjw
				}

				// Explicit order set.
				fio, fjo := fim.Ordinal(), fjm.Ordinal()
				if fio != fjo {
					return fio < fjo
				}

				// When we walk into a symlink, we keep the reference to
				// the original name.
				fin, fjn := fim.Name(), fjm.Name()
				if fin != "" && fjn != "" {
					return fin < fjn
				}

				return fii.Name() < fij.Name()
			})
		}
	}

	// First add some metadata to the dir entries
	for _, fi := range dirEntries {
		fim := fi.(FileMetaInfo)

		meta := fim.Meta()

		// Note that we use the original Name even if it's a symlink.
		name := meta.Name()
		if name == "" {
			name = fim.Name()
		}

		if name == "" {
			panic(fmt.Sprintf("[%s] no name set in %v", path, meta))
		}
		pathn := filepath.Join(path, name)

		pathMeta := pathn
		if w.basePath != "" {
			pathMeta = strings.TrimPrefix(pathn, w.basePath)
		}

		meta[metaKeyPath] = normalizeFilename(pathMeta)
		meta[metaKeyPathWalk] = pathn

		if fim.IsDir() && w.isSeen(meta.Filename()) {
			// Prevent infinite recursion
			// Possible cyclic reference
			meta[metaKeySkipDir] = true
		}
	}

	if w.hookPre != nil {
		dirEntries, err = w.hookPre(info, path, dirEntries)
		if err != nil {
			if err == filepath.SkipDir {
				return nil
			}
			return err
		}
	}

	for _, fi := range dirEntries {
		fim := fi.(FileMetaInfo)
		meta := fim.Meta()

		if meta.SkipDir() {
			continue
		}

		err := w.walk(meta.GetString(metaKeyPathWalk), fim, nil, walkFn)
		if err != nil {
			if !fi.IsDir() || err != filepath.SkipDir {
				return err
			}
		}
	}

	if w.hookPost != nil {
		dirEntries, err = w.hookPost(info, path, dirEntries)
		if err != nil {
			if err == filepath.SkipDir {
				return nil
			}
			return err
		}
	}
	return nil
}

func (w *Walkway) isSeen(filename string) bool {
	if filename == "" {
		return false
	}

	if w.seen[filename] {
		return true
	}

	w.seen[filename] = true
	return false
}
