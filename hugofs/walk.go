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
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/media"

	"github.com/spf13/afero"
)

type (
	WalkFunc func(path string, info FileMetaInfo) error
	WalkHook func(dir FileMetaInfo, path string, readdir []FileMetaInfo) ([]FileMetaInfo, error)
)

type Walkway struct {
	logger loggers.Logger

	// Prevent a walkway to be walked more than once.
	walked bool

	// Config from client.
	cfg WalkwayConfig
}

type WalkwayConfig struct {
	// The filesystem to walk.
	Fs afero.Fs

	// The root to start from in Fs.
	Root string

	// The logger to use.
	Logger     loggers.Logger
	PathParser *paths.PathParser

	// One or both of these may be pre-set.
	Info       FileMetaInfo               // The start info.
	DirEntries []FileMetaInfo             // The start info's dir entries.
	IgnoreFile func(filename string) bool // Optional

	// Will be called in order.
	HookPre  WalkHook // Optional.
	WalkFn   WalkFunc
	HookPost WalkHook // Optional.

	// Some optional flags.
	FailOnNotExist bool // If set, return an error if a directory is not found.
	SortDirEntries bool // If set, sort the dir entries by Name before calling the WalkFn, default is ReaDir order.
}

func NewWalkway(cfg WalkwayConfig) *Walkway {
	if cfg.Fs == nil {
		panic("fs must be set")
	}

	if cfg.PathParser == nil {
		cfg.PathParser = media.DefaultPathParser
	}

	logger := cfg.Logger
	if logger == nil {
		logger = loggers.NewDefault()
	}

	return &Walkway{
		cfg:    cfg,
		logger: logger,
	}
}

func (w *Walkway) Walk() error {
	if w.walked {
		panic("this walkway is already walked")
	}
	w.walked = true

	if w.cfg.Fs == NoOpFs {
		return nil
	}

	return w.walk(w.cfg.Root, w.cfg.Info, w.cfg.DirEntries)
}

// checkErr returns true if the error is handled.
func (w *Walkway) checkErr(filename string, err error) bool {
	if herrors.IsNotExist(err) && !w.cfg.FailOnNotExist {
		// The file may be removed in process.
		// This may be a ERROR situation, but it is not possible
		// to determine as a general case.
		w.logger.Warnf("File %q not found, skipping.", filename)
		return true
	}

	return false
}

// walk recursively descends path, calling walkFn.
func (w *Walkway) walk(path string, info FileMetaInfo, dirEntries []FileMetaInfo) error {
	pathRel := strings.TrimPrefix(path, w.cfg.Root)

	if info == nil {
		var err error
		fi, err := w.cfg.Fs.Stat(path)
		if err != nil {
			if path == w.cfg.Root && herrors.IsNotExist(err) {
				return nil
			}
			if w.checkErr(path, err) {
				return nil
			}
			return fmt.Errorf("walk: stat: %s", err)
		}
		info = fi.(FileMetaInfo)
	}

	err := w.cfg.WalkFn(path, info)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	if dirEntries == nil {
		f, err := w.cfg.Fs.Open(path)
		if err != nil {
			if w.checkErr(path, err) {
				return nil
			}
			return fmt.Errorf("walk: open: path: %q filename: %q: %s", path, info.Meta().Filename, err)
		}
		fis, err := f.(fs.ReadDirFile).ReadDir(-1)

		f.Close()
		if err != nil {
			if w.checkErr(path, err) {
				return nil
			}
			return fmt.Errorf("walk: Readdir: %w", err)
		}

		dirEntries = DirEntriesToFileMetaInfos(fis)
		for _, fi := range dirEntries {
			if fi.Meta().PathInfo == nil {
				fi.Meta().PathInfo = w.cfg.PathParser.Parse("", filepath.Join(pathRel, fi.Name()))
			}
		}

		if w.cfg.SortDirEntries {
			sort.Slice(dirEntries, func(i, j int) bool {
				return dirEntries[i].Name() < dirEntries[j].Name()
			})
		}

	}

	if w.cfg.IgnoreFile != nil {
		n := 0
		for _, fi := range dirEntries {
			if !w.cfg.IgnoreFile(fi.Meta().Filename) {
				dirEntries[n] = fi
				n++
			}
		}
		dirEntries = dirEntries[:n]
	}

	if w.cfg.HookPre != nil {
		var err error
		dirEntries, err = w.cfg.HookPre(info, path, dirEntries)
		if err != nil {
			if err == filepath.SkipDir {
				return nil
			}
			return err
		}
	}

	for _, fim := range dirEntries {
		nextPath := filepath.Join(path, fim.Name())
		err := w.walk(nextPath, fim, nil)
		if err != nil {
			if !fim.IsDir() || err != filepath.SkipDir {
				return err
			}
		}
	}

	if w.cfg.HookPost != nil {
		var err error
		dirEntries, err = w.cfg.HookPost(info, path, dirEntries)
		if err != nil {
			if err == filepath.SkipDir {
				return nil
			}
			return err
		}
	}
	return nil
}
