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

package hugolib

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/resources"

	"github.com/fsnotify/fsnotify"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/page"
)

// HugoSites represents the sites to build. Each site represents a language.
type HugoSites struct {
	Sites []*Site

	Configs *allconfig.Configs

	hugoInfo hugo.HugoInfo

	// Render output formats for all sites.
	renderFormats output.Formats

	// The currently rendered Site.
	currentSite *Site

	*deps.Deps

	gitInfo       *gitInfo
	codeownerInfo *codeownerInfo

	// As loaded from the /data dirs
	data map[string]any

	// Cache for page listings.
	cachePages *dynacache.Partition[string, page.Pages]
	// Cache for content sources.
	cacheContentSource *dynacache.Partition[string, *resources.StaleValue[[]byte]]

	// Before Hugo 0.122.0 we managed all translations in a map using a translationKey
	// that could be overridden in front matter.
	// Now the different page dimensions (e.g. language) are built-in to the page trees above.
	// But we sill need to support the overridden translationKey, but that should
	// be relatively rare and low volume.
	translationKeyPages *maps.SliceCache[page.Page]

	pageTrees *pageTrees

	postRenderInit sync.Once

	// File change events with filename stored in this map will be skipped.
	skipRebuildForFilenamesMu sync.Mutex
	skipRebuildForFilenames   map[string]bool

	init *hugoSitesInit

	workersSite     *para.Workers
	numWorkersSites int
	numWorkers      int

	*fatalErrorHandler
	*buildCounters
	// Tracks invocations of the Build method.
	buildCounter atomic.Uint64
}

// ShouldSkipFileChangeEvent allows skipping filesystem event early before
// the build is started.
func (h *HugoSites) ShouldSkipFileChangeEvent(ev fsnotify.Event) bool {
	h.skipRebuildForFilenamesMu.Lock()
	defer h.skipRebuildForFilenamesMu.Unlock()
	return h.skipRebuildForFilenames[ev.Name]
}

func (h *HugoSites) isRebuild() bool {
	return h.buildCounter.Load() > 0
}

func (h *HugoSites) resolveSite(lang string) *Site {
	if lang == "" {
		lang = h.Conf.DefaultContentLanguage()
	}

	for _, s := range h.Sites {
		if s.Lang() == lang {
			return s
		}
	}

	return nil
}

// Only used in tests.
type buildCounters struct {
	contentRenderCounter atomic.Uint64
	pageRenderCounter    atomic.Uint64
}

func (c *buildCounters) loggFields() logg.Fields {
	return logg.Fields{
		{Name: "pages", Value: c.pageRenderCounter.Load()},
		{Name: "content", Value: c.contentRenderCounter.Load()},
	}
}

type fatalErrorHandler struct {
	mu sync.Mutex

	h *HugoSites

	err error

	done  bool
	donec chan bool // will be closed when done
}

// FatalError error is used in some rare situations where it does not make sense to
// continue processing, to abort as soon as possible and log the error.
func (f *fatalErrorHandler) FatalError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.done {
		f.done = true
		close(f.donec)
	}
	f.err = err
}

func (f *fatalErrorHandler) getErr() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.err
}

func (f *fatalErrorHandler) Done() <-chan bool {
	return f.donec
}

type hugoSitesInit struct {
	// Loads the data from all of the /data folders.
	data *lazy.Init

	// Performs late initialization (before render) of the templates.
	layouts *lazy.Init

	// Loads the Git info and CODEOWNERS for all the pages if enabled.
	gitInfo *lazy.Init
}

func (h *HugoSites) Data() map[string]any {
	if _, err := h.init.data.Do(context.Background()); err != nil {
		h.SendError(fmt.Errorf("failed to load data: %w", err))
		return nil
	}
	return h.data
}

// Pages returns all pages for all sites.
func (h *HugoSites) Pages() page.Pages {
	key := "pages"
	v, err := h.cachePages.GetOrCreate(key, func(string) (page.Pages, error) {
		var pages page.Pages
		for _, s := range h.Sites {
			pages = append(pages, s.Pages()...)
		}
		page.SortByDefault(pages)
		return pages, nil
	})
	if err != nil {
		panic(err)
	}
	return v
}

// Pages returns all regularpages for all sites.
func (h *HugoSites) RegularPages() page.Pages {
	key := "regular-pages"
	v, err := h.cachePages.GetOrCreate(key, func(string) (page.Pages, error) {
		var pages page.Pages
		for _, s := range h.Sites {
			pages = append(pages, s.RegularPages()...)
		}
		page.SortByDefault(pages)

		return pages, nil
	})
	if err != nil {
		panic(err)
	}
	return v
}

func (h *HugoSites) gitInfoForPage(p page.Page) (source.GitInfo, error) {
	if _, err := h.init.gitInfo.Do(context.Background()); err != nil {
		return source.GitInfo{}, err
	}

	if h.gitInfo == nil {
		return source.GitInfo{}, nil
	}

	return h.gitInfo.forPage(p), nil
}

func (h *HugoSites) codeownersForPage(p page.Page) ([]string, error) {
	if _, err := h.init.gitInfo.Do(context.Background()); err != nil {
		return nil, err
	}

	if h.codeownerInfo == nil {
		return nil, nil
	}

	return h.codeownerInfo.forPage(p), nil
}

func (h *HugoSites) pickOneAndLogTheRest(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	var i int

	for j, err := range errors {
		// If this is in server mode, we want to return an error to the client
		// with a file context, if possible.
		if herrors.UnwrapFileError(err) != nil {
			i = j
			break
		}
	}

	// Log the rest, but add a threshold to avoid flooding the log.
	const errLogThreshold = 5

	for j, err := range errors {
		if j == i || err == nil {
			continue
		}

		if j >= errLogThreshold {
			break
		}

		h.Log.Errorln(err)
	}

	return errors[i]
}

func (h *HugoSites) isMultilingual() bool {
	return len(h.Sites) > 1
}

// TODO(bep) consolidate
func (h *HugoSites) LanguageSet() map[string]int {
	set := make(map[string]int)
	for i, s := range h.Sites {
		set[s.language.Lang] = i
	}
	return set
}

func (h *HugoSites) NumLogErrors() int {
	if h == nil {
		return 0
	}
	return h.Log.LoggCount(logg.LevelError)
}

func (h *HugoSites) PrintProcessingStats(w io.Writer) {
	stats := make([]*helpers.ProcessingStats, len(h.Sites))
	for i := 0; i < len(h.Sites); i++ {
		stats[i] = h.Sites[i].PathSpec.ProcessingStats
	}
	helpers.ProcessingStatsTable(w, stats...)
}

// GetContentPage finds a Page with content given the absolute filename.
// Returns nil if none found.
func (h *HugoSites) GetContentPage(filename string) page.Page {
	var p page.Page

	h.withPage(func(s string, p2 *pageState) bool {
		if p2.File() == nil {
			return false
		}

		if p2.File().FileInfo().Meta().Filename == filename {
			p = p2
			return true
		}

		for _, r := range p2.Resources().ByType(pageResourceType) {
			p3 := r.(page.Page)
			if p3.File() != nil && p3.File().FileInfo().Meta().Filename == filename {
				p = p3
				return true
			}
		}

		return false
	})

	return p
}

func (h *HugoSites) loadGitInfo() error {
	if h.Configs.Base.EnableGitInfo {
		gi, err := newGitInfo(h.Deps)
		if err != nil {
			h.Log.Errorln("Failed to read Git log:", err)
		} else {
			h.gitInfo = gi
		}

		co, err := newCodeOwners(h.Configs.LoadingInfo.BaseConfig.WorkingDir)
		if err != nil {
			h.Log.Errorln("Failed to read CODEOWNERS:", err)
		} else {
			h.codeownerInfo = co
		}
	}
	return nil
}

// Reset resets the sites and template caches etc., making it ready for a full rebuild.
func (h *HugoSites) reset(config *BuildCfg) {
	h.fatalErrorHandler = &fatalErrorHandler{
		h:     h,
		donec: make(chan bool),
	}
}

// resetLogs resets the log counters etc. Used to do a new build on the same sites.
func (h *HugoSites) resetLogs() {
	h.Log.Reset()
	for _, s := range h.Sites {
		s.Deps.Log.Reset()
	}
}

func (h *HugoSites) withSite(fn func(s *Site) error) error {
	for _, s := range h.Sites {
		if err := fn(s); err != nil {
			return err
		}
	}
	return nil
}

func (h *HugoSites) withPage(fn func(s string, p *pageState) bool) {
	h.withSite(func(s *Site) error {
		w := &doctree.NodeShiftTreeWalker[contentNodeI]{
			Tree:     s.pageMap.treePages,
			LockType: doctree.LockTypeRead,
			Handle: func(s string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
				return fn(s, n.(*pageState)), nil
			},
		}
		return w.Walk(context.Background())
	})
}

// BuildCfg holds build options used to, as an example, skip the render step.
type BuildCfg struct {
	// Skip rendering. Useful for testing.
	SkipRender bool

	// Use this to indicate what changed (for rebuilds).
	WhatChanged *WhatChanged

	// This is a partial re-render of some selected pages.
	PartialReRender bool

	// Set in server mode when the last build failed for some reason.
	ErrRecovery bool

	// Recently visited URLs. This is used for partial re-rendering.
	RecentlyVisited *types.EvictingStringQueue

	// Can be set to build only with a sub set of the content source.
	ContentInclusionFilter *glob.FilenameFilter

	// Set when the buildlock is already acquired (e.g. the archetype content builder).
	NoBuildLock bool

	testCounters *buildCounters
}

// shouldRender returns whether this output format should be rendered or not.
func (cfg *BuildCfg) shouldRender(p *pageState) bool {
	if p.skipRender() {
		return false
	}

	if !p.renderOnce {
		return true
	}

	// The render state is incremented on render and reset when a related change is detected.
	// Note that this is set per output format.
	shouldRender := p.renderState == 0

	if !shouldRender {
		return false
	}

	fastRenderMode := p.s.Conf.FastRenderMode()

	if !fastRenderMode || p.s.h.buildCounter.Load() == 0 {
		return shouldRender
	}

	if !p.render {
		// Not be to rendered for this output format.
		return false
	}

	if p.outputFormat().IsHTML {
		// This is fast render mode and the output format is HTML,
		// rerender if this page is one of the recently visited.
		return cfg.RecentlyVisited.Contains(p.RelPermalink())
	}

	// In fast render mode, we want to avoid re-rendering the sitemaps etc. and
	// other big listings whenever we e.g. change a content file,
	// but we want partial renders of the recently visited pages to also include
	// alternative formats of the same HTML page (e.g. RSS, JSON).
	for _, po := range p.pageOutputs {
		if po.render && po.f.IsHTML && cfg.RecentlyVisited.Contains(po.RelPermalink()) {
			return true
		}
	}

	return false
}

func (s *Site) preparePagesForRender(isRenderingSite bool, idx int) error {
	var err error

	initPage := func(p *pageState) error {
		if err = p.shiftToOutputFormat(isRenderingSite, idx); err != nil {
			return err
		}
		return nil
	}

	return s.pageMap.forEeachPageIncludingBundledPages(nil,
		func(p *pageState) (bool, error) {
			return false, initPage(p)
		},
	)
}

func (h *HugoSites) loadData() error {
	h.data = make(map[string]any)
	w := hugofs.NewWalkway(
		hugofs.WalkwayConfig{
			Fs:         h.PathSpec.BaseFs.Data.Fs,
			IgnoreFile: h.SourceSpec.IgnoreFile,
			PathParser: h.Conf.PathParser(),
			WalkFn: func(path string, fi hugofs.FileMetaInfo) error {
				if fi.IsDir() {
					return nil
				}
				pi := fi.Meta().PathInfo
				if pi == nil {
					panic("no path info")
				}
				return h.handleDataFile(source.NewFileInfo(fi))
			},
		})

	if err := w.Walk(); err != nil {
		return err
	}
	return nil
}

func (h *HugoSites) handleDataFile(r *source.File) error {
	var current map[string]any

	f, err := r.FileInfo().Meta().Open()
	if err != nil {
		return fmt.Errorf("data: failed to open %q: %w", r.LogicalName(), err)
	}
	defer f.Close()

	// Crawl in data tree to insert data
	current = h.data
	dataPath := r.FileInfo().Meta().PathInfo.Unnormalized().Dir()[1:]
	keyParts := strings.Split(dataPath, "/")

	for _, key := range keyParts {
		if key != "" {
			if _, ok := current[key]; !ok {
				current[key] = make(map[string]any)
			}
			current = current[key].(map[string]any)
		}
	}

	data, err := h.readData(r)
	if err != nil {
		return h.errWithFileContext(err, r)
	}

	if data == nil {
		return nil
	}

	// filepath.Walk walks the files in lexical order, '/' comes before '.'
	higherPrecedentData := current[r.BaseFileName()]

	switch data.(type) {
	case nil:
	case map[string]any:

		switch higherPrecedentData.(type) {
		case nil:
			current[r.BaseFileName()] = data
		case map[string]any:
			// merge maps: insert entries from data for keys that
			// don't already exist in higherPrecedentData
			higherPrecedentMap := higherPrecedentData.(map[string]any)
			for key, value := range data.(map[string]any) {
				if _, exists := higherPrecedentMap[key]; exists {
					// this warning could happen if
					// 1. A theme uses the same key; the main data folder wins
					// 2. A sub folder uses the same key: the sub folder wins
					// TODO(bep) figure out a way to detect 2) above and make that a WARN
					h.Log.Infof("Data for key '%s' in path '%s' is overridden by higher precedence data already in the data tree", key, r.Path())
				} else {
					higherPrecedentMap[key] = value
				}
			}
		default:
			// can't merge: higherPrecedentData is not a map
			h.Log.Warnf("The %T data from '%s' overridden by "+
				"higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
		}

	case []any:
		if higherPrecedentData == nil {
			current[r.BaseFileName()] = data
		} else {
			// we don't merge array data
			h.Log.Warnf("The %T data from '%s' overridden by "+
				"higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
		}

	default:
		h.Log.Errorf("unexpected data type %T in file %s", data, r.LogicalName())
	}

	return nil
}

func (h *HugoSites) errWithFileContext(err error, f *source.File) error {
	realFilename := f.FileInfo().Meta().Filename
	return herrors.NewFileErrorFromFile(err, realFilename, h.Fs.Source, nil)
}

func (h *HugoSites) readData(f *source.File) (any, error) {
	file, err := f.FileInfo().Meta().Open()
	if err != nil {
		return nil, fmt.Errorf("readData: failed to open data file: %w", err)
	}
	defer file.Close()
	content := helpers.ReaderToBytes(file)

	format := metadecoders.FormatFromString(f.Ext())
	return metadecoders.Default.Unmarshal(content, format)
}
