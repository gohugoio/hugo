// Copyright 2021 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/rungroup"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/hugofs"
)

func newPagesCollector(
	ctx context.Context,
	h *HugoSites,
	sp *source.SourceSpec,
	logger loggers.Logger,
	infoLogger logg.LevelLogger,
	m *pageMap,
	buildConfig *BuildCfg,
	ids []pathChange,
) *pagesCollector {
	return &pagesCollector{
		ctx:         ctx,
		h:           h,
		fs:          sp.BaseFs.Content.Fs,
		m:           m,
		sp:          sp,
		logger:      logger,
		infoLogger:  infoLogger,
		buildConfig: buildConfig,
		ids:         ids,
		seenDirs:    make(map[string]bool),
	}
}

type pagesCollector struct {
	ctx        context.Context
	h          *HugoSites
	sp         *source.SourceSpec
	logger     loggers.Logger
	infoLogger logg.LevelLogger

	m *pageMap

	fs afero.Fs

	buildConfig *BuildCfg

	// List of paths that have changed. Used in partial builds.
	ids      []pathChange
	seenDirs map[string]bool

	g rungroup.Group[hugofs.FileMetaInfo]
}

// Collect collects content by walking the file system and storing
// it in the content tree.
// It may be restricted by filenames set on the collector (partial build).
func (c *pagesCollector) Collect() (collectErr error) {
	var (
		numWorkers             = c.h.numWorkers
		numFilesProcessedTotal atomic.Uint64
		numPagesProcessedTotal atomic.Uint64
		numResourcesProcessed  atomic.Uint64
		numFilesProcessedLast  uint64
		fileBatchTimer         = time.Now()
		fileBatchTimerMu       sync.Mutex
	)

	l := c.infoLogger.WithField("substep", "collect")

	logFilesProcessed := func(force bool) {
		fileBatchTimerMu.Lock()
		if force || time.Since(fileBatchTimer) > 3*time.Second {
			numFilesProcessedBatch := numFilesProcessedTotal.Load() - numFilesProcessedLast
			numFilesProcessedLast = numFilesProcessedTotal.Load()
			loggers.TimeTrackf(l, fileBatchTimer,
				logg.Fields{
					logg.Field{Name: "files", Value: numFilesProcessedBatch},
					logg.Field{Name: "files_total", Value: numFilesProcessedTotal.Load()},
					logg.Field{Name: "pages_total", Value: numPagesProcessedTotal.Load()},
					logg.Field{Name: "resources_total", Value: numResourcesProcessed.Load()},
				},
				"",
			)
			fileBatchTimer = time.Now()
		}
		fileBatchTimerMu.Unlock()
	}

	defer func() {
		logFilesProcessed(true)
	}()

	c.g = rungroup.Run[hugofs.FileMetaInfo](c.ctx, rungroup.Config[hugofs.FileMetaInfo]{
		NumWorkers: numWorkers,
		Handle: func(ctx context.Context, fi hugofs.FileMetaInfo) error {
			numPages, numResources, err := c.m.AddFi(fi, c.buildConfig)
			if err != nil {
				return hugofs.AddFileInfoToError(err, fi, c.fs)
			}
			numFilesProcessedTotal.Add(1)
			numPagesProcessedTotal.Add(numPages)
			numResourcesProcessed.Add(numResources)
			if numFilesProcessedTotal.Load()%1000 == 0 {
				logFilesProcessed(false)
			}
			return nil
		},
	})

	if c.ids == nil {
		// Collect everything.
		collectErr = c.collectDir(nil, false, nil)
	} else {
		for _, s := range c.h.Sites {
			s.pageMap.cfg.isRebuild = true
		}

		for _, id := range c.ids {
			if id.p.IsLeafBundle() {
				collectErr = c.collectDir(
					id.p,
					false,
					func(fim hugofs.FileMetaInfo) bool {
						return true
					},
				)
			} else if id.p.IsBranchBundle() {
				collectErr = c.collectDir(
					id.p,
					false,
					func(fim hugofs.FileMetaInfo) bool {
						if fim.IsDir() {
							return id.isStructuralChange()
						}
						fimp := fim.Meta().PathInfo
						if fimp == nil {
							return false
						}

						return strings.HasPrefix(fimp.Path(), paths.AddTrailingSlash(id.p.Dir()))
					},
				)
			} else {
				// We always start from a directory.
				collectErr = c.collectDir(id.p, id.isDir, func(fim hugofs.FileMetaInfo) bool {
					if id.isStructuralChange() {
						if id.isDir && fim.Meta().PathInfo.IsLeafBundle() {
							return strings.HasPrefix(fim.Meta().PathInfo.Path(), paths.AddTrailingSlash(id.p.Path()))
						}

						return id.p.Dir() == fim.Meta().PathInfo.Dir()
					}

					if fim.Meta().PathInfo.IsLeafBundle() && id.p.BundleType() == paths.PathTypeContentSingle {
						return id.p.Dir() == fim.Meta().PathInfo.Dir()
					}

					return id.p.Path() == fim.Meta().PathInfo.Path()
				})
			}

			if collectErr != nil {
				break
			}
		}

	}

	werr := c.g.Wait()
	if collectErr == nil {
		collectErr = werr
	}

	return
}

func (c *pagesCollector) collectDir(dirPath *paths.Path, isDir bool, inFilter func(fim hugofs.FileMetaInfo) bool) error {
	var dpath string
	if dirPath != nil {
		if isDir {
			dpath = filepath.FromSlash(dirPath.Unnormalized().Path())
		} else {
			dpath = filepath.FromSlash(dirPath.Unnormalized().Dir())
		}
	}

	if c.seenDirs[dpath] {
		return nil
	}
	c.seenDirs[dpath] = true

	root, err := c.fs.Stat(dpath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	rootm := root.(hugofs.FileMetaInfo)

	if err := c.collectDirDir(dpath, rootm, inFilter); err != nil {
		return err
	}

	return nil
}

func (c *pagesCollector) collectDirDir(path string, root hugofs.FileMetaInfo, inFilter func(fim hugofs.FileMetaInfo) bool) error {
	filter := func(fim hugofs.FileMetaInfo) bool {
		if inFilter != nil {
			return inFilter(fim)
		}
		return true
	}

	preHook := func(dir hugofs.FileMetaInfo, path string, readdir []hugofs.FileMetaInfo) ([]hugofs.FileMetaInfo, error) {
		filtered := readdir[:0]
		for _, fi := range readdir {
			if filter(fi) {
				filtered = append(filtered, fi)
			}
		}
		readdir = filtered
		if len(readdir) == 0 {
			return nil, nil
		}

		n := 0
		for _, fi := range readdir {
			if fi.Meta().PathInfo.IsContentData() {
				// _content.json
				// These are not part of any bundle, so just add them directly and remove them from the readdir slice.
				if err := c.g.Enqueue(fi); err != nil {
					return nil, err
				}
			} else {
				readdir[n] = fi
				n++
			}
		}
		readdir = readdir[:n]

		// Pick the first regular file.
		var first hugofs.FileMetaInfo
		for _, fi := range readdir {
			if fi.IsDir() {
				continue
			}
			first = fi
			break
		}

		if first == nil {
			// Only dirs, keep walking.
			return readdir, nil
		}

		// Any bundle file will always be first.
		firstPi := first.Meta().PathInfo

		if firstPi == nil {
			panic(fmt.Sprintf("collectDirDir: no path info for %q", first.Meta().Filename))
		}

		if firstPi.IsLeafBundle() {
			if err := c.handleBundleLeaf(dir, first, path, readdir); err != nil {
				return nil, err
			}
			return nil, filepath.SkipDir
		}

		seen := map[hstrings.Tuple]bool{}
		for _, fi := range readdir {
			if fi.IsDir() {
				continue
			}

			pi := fi.Meta().PathInfo
			meta := fi.Meta()

			// Filter out duplicate page or resource.
			// These would eventually have been filtered out as duplicates when
			// inserting them into the document store,
			// but doing it here will preserve a consistent ordering.
			baseLang := hstrings.Tuple{First: pi.Base(), Second: meta.Lang}
			if seen[baseLang] {
				continue
			}
			seen[baseLang] = true

			if pi == nil {
				panic(fmt.Sprintf("no path info for %q", meta.Filename))
			}

			if meta.Lang == "" {
				panic("lang not set")
			}

			if err := c.g.Enqueue(fi); err != nil {
				return nil, err
			}
		}

		// Keep walking.
		return readdir, nil
	}

	var postHook hugofs.WalkHook

	wfn := func(path string, fi hugofs.FileMetaInfo) error {
		return nil
	}

	w := hugofs.NewWalkway(
		hugofs.WalkwayConfig{
			Logger:     c.logger,
			Root:       path,
			Info:       root,
			Fs:         c.fs,
			IgnoreFile: c.h.SourceSpec.IgnoreFile,
			PathParser: c.h.Conf.PathParser(),
			HookPre:    preHook,
			HookPost:   postHook,
			WalkFn:     wfn,
		})

	return w.Walk()
}

func (c *pagesCollector) handleBundleLeaf(dir, bundle hugofs.FileMetaInfo, inPath string, readdir []hugofs.FileMetaInfo) error {
	bundlePi := bundle.Meta().PathInfo
	seen := map[hstrings.Tuple]bool{}

	walk := func(path string, info hugofs.FileMetaInfo) error {
		if info.IsDir() {
			return nil
		}

		pi := info.Meta().PathInfo

		if info != bundle {
			// Everything inside a leaf bundle is a Resource,
			// even the content pages.
			// Note that we do allow index.md as page resources, but not in the bundle root.
			if !pi.IsLeafBundle() || pi.Dir() != bundlePi.Dir() {
				paths.ModifyPathBundleTypeResource(pi)
			}
		}

		// Filter out duplicate page or resource.
		// These would eventually have been filtered out as duplicates when
		// inserting them into the document store,
		// but doing it here will preserve a consistent ordering.
		baseLang := hstrings.Tuple{First: pi.Base(), Second: info.Meta().Lang}
		if seen[baseLang] {
			return nil
		}
		seen[baseLang] = true

		return c.g.Enqueue(info)
	}

	// Start a new walker from the given path.
	w := hugofs.NewWalkway(
		hugofs.WalkwayConfig{
			Root:       inPath,
			Fs:         c.fs,
			Logger:     c.logger,
			Info:       dir,
			DirEntries: readdir,
			IgnoreFile: c.h.SourceSpec.IgnoreFile,
			PathParser: c.h.Conf.PathParser(),
			WalkFn:     walk,
		})

	return w.Walk()
}
