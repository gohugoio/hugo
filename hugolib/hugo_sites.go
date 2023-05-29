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

package hugolib

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/fsnotify/fsnotify"

	"github.com/gohugoio/hugo/identity"

	radix "github.com/armon/go-radix"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/tpl"
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

	contentInit sync.Once
	content     *pageMaps

	// Keeps track of bundle directories and symlinks to enable partial rebuilding.
	ContentChanges *contentChangeMap

	// File change events with filename stored in this map will be skipped.
	skipRebuildForFilenamesMu sync.Mutex
	skipRebuildForFilenames   map[string]bool

	init *hugoSitesInit

	workers    *para.Workers
	numWorkers int

	*fatalErrorHandler
	*testCounters
}

// ShouldSkipFileChangeEvent allows skipping filesystem event early before
// the build is started.
func (h *HugoSites) ShouldSkipFileChangeEvent(ev fsnotify.Event) bool {
	h.skipRebuildForFilenamesMu.Lock()
	defer h.skipRebuildForFilenamesMu.Unlock()
	return h.skipRebuildForFilenames[ev.Name]
}

func (h *HugoSites) getContentMaps() *pageMaps {
	h.contentInit.Do(func() {
		h.content = newPageMaps(h)
	})
	return h.content
}

// Only used in tests.
type testCounters struct {
	contentRenderCounter uint64
	pageRenderCounter    uint64
}

func (h *testCounters) IncrContentRender() {
	if h == nil {
		return
	}
	atomic.AddUint64(&h.contentRenderCounter, 1)
}

func (h *testCounters) IncrPageRender() {
	if h == nil {
		return
	}
	atomic.AddUint64(&h.pageRenderCounter, 1)
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

	// Maps page translations.
	translations *lazy.Init
}

func (h *hugoSitesInit) Reset() {
	h.data.Reset()
	h.layouts.Reset()
	h.gitInfo.Reset()
	h.translations.Reset()
}

func (h *HugoSites) Data() map[string]any {
	if _, err := h.init.data.Do(context.Background()); err != nil {
		h.SendError(fmt.Errorf("failed to load data: %w", err))
		return nil
	}
	return h.data
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

func (h *HugoSites) isMultiLingual() bool {
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
	return int(h.Log.LogCounters().ErrorCounter.Count())
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

	h.getContentMaps().walkBundles(func(b *contentNode) bool {
		if b.p == nil || b.fi == nil {
			return false
		}

		if b.fi.Meta().Filename == filename {
			p = b.p
			return true
		}

		return false
	})

	return p
}

func (h *HugoSites) loadGitInfo() error {
	if h.Configs.Base.EnableGitInfo {
		gi, err := newGitInfo(h.Conf)
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

func (s *Site) withSiteTemplates(withTemplates ...func(templ tpl.TemplateManager) error) func(templ tpl.TemplateManager) error {
	return func(templ tpl.TemplateManager) error {
		for _, wt := range withTemplates {
			if wt == nil {
				continue
			}
			if err := wt(templ); err != nil {
				return err
			}
		}

		return nil
	}
}

// Reset resets the sites and template caches etc., making it ready for a full rebuild.
func (h *HugoSites) reset(config *BuildCfg) {
	if config.ResetState {
		for _, s := range h.Sites {
			if r, ok := s.Fs.PublishDir.(hugofs.Reseter); ok {
				r.Reset()
			}
		}
	}

	h.fatalErrorHandler = &fatalErrorHandler{
		h:     h,
		donec: make(chan bool),
	}

	h.init.Reset()
}

// resetLogs resets the log counters etc. Used to do a new build on the same sites.
func (h *HugoSites) resetLogs() {
	h.Log.Reset()
	loggers.GlobalErrorCounter.Reset()
	for _, s := range h.Sites {
		s.Deps.Log.Reset()
		s.Deps.LogDistinct.Reset()
	}
}

func (h *HugoSites) withSite(fn func(s *Site) error) error {
	if h.workers == nil {
		for _, s := range h.Sites {
			if err := fn(s); err != nil {
				return err
			}
		}
		return nil
	}

	g, _ := h.workers.Start(context.Background())
	for _, s := range h.Sites {
		s := s
		g.Run(func() error {
			return fn(s)
		})
	}
	return g.Wait()
}

// BuildCfg holds build options used to, as an example, skip the render step.
type BuildCfg struct {
	// Reset site state before build. Use to force full rebuilds.
	ResetState bool
	// Skip rendering. Useful for testing.
	SkipRender bool
	// Use this to indicate what changed (for rebuilds).
	whatChanged *whatChanged

	// This is a partial re-render of some selected pages. This means
	// we should skip most of the processing.
	PartialReRender bool

	// Set in server mode when the last build failed for some reason.
	ErrRecovery bool

	// Recently visited URLs. This is used for partial re-rendering.
	RecentlyVisited map[string]bool

	// Can be set to build only with a sub set of the content source.
	ContentInclusionFilter *glob.FilenameFilter

	// Set when the buildlock is already acquired (e.g. the archetype content builder).
	NoBuildLock bool

	testCounters *testCounters
}

// shouldRender is used in the Fast Render Mode to determine if we need to re-render
// a Page: If it is recently visited (the home pages will always be in this set) or changed.
// Note that a page does not have to have a content page / file.
// For regular builds, this will always return true.
// TODO(bep) rename/work this.
func (cfg *BuildCfg) shouldRender(p *pageState) bool {
	if p == nil {
		return false
	}

	if p.forceRender {
		return true
	}

	if len(cfg.RecentlyVisited) == 0 {
		return true
	}

	if cfg.RecentlyVisited[p.RelPermalink()] {
		return true
	}

	if cfg.whatChanged != nil && !p.File().IsZero() {
		return cfg.whatChanged.files[p.File().Filename()]
	}

	return false
}

func (h *HugoSites) renderCrossSitesSitemap() error {
	if !h.isMultiLingual() || h.Conf.IsMultihost() {
		return nil
	}

	sitemapEnabled := false
	for _, s := range h.Sites {
		if s.conf.IsKindEnabled(kindSitemap) {
			sitemapEnabled = true
			break
		}
	}

	if !sitemapEnabled {
		return nil
	}

	s := h.Sites[0]
	// We don't have any page context to pass in here.
	ctx := context.Background()

	templ := s.lookupLayouts("sitemapindex.xml", "_default/sitemapindex.xml", "_internal/_default/sitemapindex.xml")
	return s.renderAndWriteXML(ctx, &s.PathSpec.ProcessingStats.Sitemaps, "sitemapindex",
		s.conf.Sitemap.Filename, h.Sites, templ)
}

func (h *HugoSites) renderCrossSitesRobotsTXT() error {
	if h.Configs.IsMultihost {
		return nil
	}
	if !h.Configs.Base.EnableRobotsTXT {
		return nil
	}

	s := h.Sites[0]

	p, err := newPageStandalone(&pageMeta{
		s:    s,
		kind: kindRobotsTXT,
		urlPaths: pagemeta.URLPath{
			URL: "robots.txt",
		},
	},
		output.RobotsTxtFormat)
	if err != nil {
		return err
	}

	if !p.render {
		return nil
	}

	templ := s.lookupLayouts("robots.txt", "_default/robots.txt", "_internal/_default/robots.txt")

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "Robots Txt", "robots.txt", p, templ)
}

func (h *HugoSites) removePageByFilename(filename string) {
	h.getContentMaps().withMaps(func(m *pageMap) error {
		m.deleteBundleMatching(func(b *contentNode) bool {
			if b.p == nil {
				return false
			}

			if b.fi == nil {
				return false
			}

			return b.fi.Meta().Filename == filename
		})
		return nil
	})
}

func (h *HugoSites) createPageCollections() error {
	allPages := newLazyPagesFactory(func() page.Pages {
		var pages page.Pages
		for _, s := range h.Sites {
			pages = append(pages, s.Pages()...)
		}

		page.SortByDefault(pages)

		return pages
	})

	allRegularPages := newLazyPagesFactory(func() page.Pages {
		return h.findPagesByKindIn(page.KindPage, allPages.get())
	})

	for _, s := range h.Sites {
		s.PageCollections.allPages = allPages
		s.PageCollections.allRegularPages = allRegularPages
	}

	return nil
}

func (s *Site) preparePagesForRender(isRenderingSite bool, idx int) error {
	var err error
	s.pageMap.withEveryBundlePage(func(p *pageState) bool {
		if err = p.initOutputFormat(isRenderingSite, idx); err != nil {
			return true
		}
		return false
	})
	return nil
}

// Pages returns all pages for all sites.
func (h *HugoSites) Pages() page.Pages {
	return h.Sites[0].AllPages()
}

func (h *HugoSites) loadData(fis []hugofs.FileMetaInfo) (err error) {
	spec := source.NewSourceSpec(h.PathSpec, nil, nil)

	h.data = make(map[string]any)
	for _, fi := range fis {
		fileSystem := spec.NewFilesystemFromFileMetaInfo(fi)
		files, err := fileSystem.Files()
		if err != nil {
			return err
		}
		for _, r := range files {
			if err := h.handleDataFile(r); err != nil {
				return err
			}
		}
	}

	return
}

func (h *HugoSites) handleDataFile(r source.File) error {
	var current map[string]any

	f, err := r.FileInfo().Meta().Open()
	if err != nil {
		return fmt.Errorf("data: failed to open %q: %w", r.LogicalName(), err)
	}
	defer f.Close()

	// Crawl in data tree to insert data
	current = h.data
	keyParts := strings.Split(r.Dir(), helpers.FilePathSeparator)

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

func (h *HugoSites) errWithFileContext(err error, f source.File) error {
	fim, ok := f.FileInfo().(hugofs.FileMetaInfo)
	if !ok {
		return err
	}
	realFilename := fim.Meta().Filename

	return herrors.NewFileErrorFromFile(err, realFilename, h.SourceSpec.Fs.Source, nil)

}

func (h *HugoSites) readData(f source.File) (any, error) {
	file, err := f.FileInfo().Meta().Open()
	if err != nil {
		return nil, fmt.Errorf("readData: failed to open data file: %w", err)
	}
	defer file.Close()
	content := helpers.ReaderToBytes(file)

	format := metadecoders.FormatFromString(f.Ext())
	return metadecoders.Default.Unmarshal(content, format)
}

func (h *HugoSites) findPagesByKindIn(kind string, inPages page.Pages) page.Pages {
	return h.Sites[0].findPagesByKindIn(kind, inPages)
}

func (h *HugoSites) resetPageState() {
	h.getContentMaps().walkBundles(func(n *contentNode) bool {
		if n.p == nil {
			return false
		}
		p := n.p
		for _, po := range p.pageOutputs {
			if po.cp == nil {
				continue
			}
			po.cp.Reset()
		}

		return false
	})
}

func (h *HugoSites) resetPageStateFromEvents(idset identity.Identities) {
	h.getContentMaps().walkBundles(func(n *contentNode) bool {
		if n.p == nil {
			return false
		}
		p := n.p
	OUTPUTS:
		for _, po := range p.pageOutputs {
			if po.cp == nil {
				continue
			}
			for id := range idset {
				if po.cp.dependencyTracker.Search(id) != nil {
					po.cp.Reset()
					continue OUTPUTS
				}
			}
		}

		if p.shortcodeState == nil {
			return false
		}

		for _, s := range p.shortcodeState.shortcodes {
			for _, templ := range s.templs {
				sid := templ.(identity.Manager)
				for id := range idset {
					if sid.Search(id) != nil {
						for _, po := range p.pageOutputs {
							if po.cp != nil {
								po.cp.Reset()
							}
						}
						return false
					}
				}
			}
		}
		return false
	})
}

// Used in partial reloading to determine if the change is in a bundle.
type contentChangeMap struct {
	mu sync.RWMutex

	// Holds directories with leaf bundles.
	leafBundles *radix.Tree

	// Holds directories with branch bundles.
	branchBundles map[string]bool

	pathSpec *helpers.PathSpec

	// Hugo supports symlinked content (both directories and files). This
	// can lead to situations where the same file can be referenced from several
	// locations in /content -- which is really cool, but also means we have to
	// go an extra mile to handle changes.
	// This map is only used in watch mode.
	// It maps either file to files or the real dir to a set of content directories
	// where it is in use.
	symContentMu sync.Mutex
	symContent   map[string]map[string]bool
}

func (m *contentChangeMap) add(dirname string, tp bundleDirType) {
	m.mu.Lock()
	if !strings.HasSuffix(dirname, helpers.FilePathSeparator) {
		dirname += helpers.FilePathSeparator
	}
	switch tp {
	case bundleBranch:
		m.branchBundles[dirname] = true
	case bundleLeaf:
		m.leafBundles.Insert(dirname, true)
	default:
		m.mu.Unlock()
		panic("invalid bundle type")
	}
	m.mu.Unlock()
}

func (m *contentChangeMap) resolveAndRemove(filename string) (string, bundleDirType) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Bundles share resources, so we need to start from the virtual root.
	relFilename := m.pathSpec.RelContentDir(filename)
	dir, name := filepath.Split(relFilename)
	if !strings.HasSuffix(dir, helpers.FilePathSeparator) {
		dir += helpers.FilePathSeparator
	}

	if _, found := m.branchBundles[dir]; found {
		delete(m.branchBundles, dir)
		return dir, bundleBranch
	}

	if key, _, found := m.leafBundles.LongestPrefix(dir); found {
		m.leafBundles.Delete(key)
		dir = string(key)
		return dir, bundleLeaf
	}

	fileTp, isContent := classifyBundledFile(name)
	if isContent && fileTp != bundleNot {
		// A new bundle.
		return dir, fileTp
	}

	return dir, bundleNot
}

func (m *contentChangeMap) addSymbolicLinkMapping(fim hugofs.FileMetaInfo) {
	meta := fim.Meta()
	if !meta.IsSymlink {
		return
	}
	m.symContentMu.Lock()

	from, to := meta.Filename, meta.OriginalFilename
	if fim.IsDir() {
		if !strings.HasSuffix(from, helpers.FilePathSeparator) {
			from += helpers.FilePathSeparator
		}
	}

	mm, found := m.symContent[from]

	if !found {
		mm = make(map[string]bool)
		m.symContent[from] = mm
	}
	mm[to] = true
	m.symContentMu.Unlock()
}

func (m *contentChangeMap) GetSymbolicLinkMappings(dir string) []string {
	mm, found := m.symContent[dir]
	if !found {
		return nil
	}
	dirs := make([]string, len(mm))
	i := 0
	for dir := range mm {
		dirs[i] = dir
		i++
	}

	sort.Strings(dirs)

	return dirs
}
