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
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/identity"

	radix "github.com/armon/go-radix"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/source"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/publisher"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/langs/i18n"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/tpl"
	"github.com/gohugoio/hugo/tpl/tplimpl"
)

// HugoSites represents the sites to build. Each site represents a language.
type HugoSites struct {
	Sites []*Site

	multilingual *Multilingual

	// Multihost is set if multilingual and baseURL set on the language level.
	multihost bool

	// If this is running in the dev server.
	running bool

	// Serializes rebuilds when server is running.
	runningMu sync.Mutex

	// Render output formats for all sites.
	renderFormats output.Formats

	*deps.Deps

	gitInfo *gitInfo

	// As loaded from the /data dirs
	data map[string]interface{}

	content *pageMaps

	// Keeps track of bundle directories and symlinks to enable partial rebuilding.
	ContentChanges *contentChangeMap

	init *hugoSitesInit

	workers    *para.Workers
	numWorkers int

	*fatalErrorHandler
	*testCounters
}

// Only used in tests.
type testCounters struct {
	contentRenderCounter uint64
}

func (h *testCounters) IncrContentRender() {
	if h == nil {
		return
	}
	atomic.AddUint64(&h.contentRenderCounter, 1)
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

	// Loads the Git info for all the pages if enabled.
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

func (h *HugoSites) Data() map[string]interface{} {
	if _, err := h.init.data.Do(); err != nil {
		h.SendError(errors.Wrap(err, "failed to load data"))
		return nil
	}
	return h.data
}

func (h *HugoSites) gitInfoForPage(p page.Page) (*gitmap.GitInfo, error) {
	if _, err := h.init.gitInfo.Do(); err != nil {
		return nil, err
	}

	if h.gitInfo == nil {
		return nil, nil
	}

	return h.gitInfo.forPage(p), nil
}

func (h *HugoSites) siteInfos() page.Sites {
	infos := make(page.Sites, len(h.Sites))
	for i, site := range h.Sites {
		infos[i] = site.Info
	}
	return infos
}

func (h *HugoSites) pickOneAndLogTheRest(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	var i int

	for j, err := range errors {
		// If this is in server mode, we want to return an error to the client
		// with a file context, if possible.
		if herrors.UnwrapErrorWithFileContext(err) != nil {
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

		h.Log.ERROR.Println(err)
	}

	return errors[i]
}

func (h *HugoSites) IsMultihost() bool {
	return h != nil && h.multihost
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
	return int(h.Log.ErrorCounter.Count())
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

	h.content.walkBundles(func(b *contentNode) bool {
		if b.p == nil || b.fi == nil {
			return false
		}

		if b.fi.Meta().Filename() == filename {
			p = b.p
			return true
		}

		return false
	})

	return p
}

// NewHugoSites creates a new collection of sites given the input sites, building
// a language configuration based on those.
func newHugoSites(cfg deps.DepsCfg, sites ...*Site) (*HugoSites, error) {

	if cfg.Language != nil {
		return nil, errors.New("Cannot provide Language in Cfg when sites are provided")
	}

	langConfig, err := newMultiLingualFromSites(cfg.Cfg, sites...)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create language config")
	}

	var contentChangeTracker *contentChangeMap

	numWorkers := config.GetNumWorkerMultiplier()
	if numWorkers > len(sites) {
		numWorkers = len(sites)
	}
	var workers *para.Workers
	if numWorkers > 1 {
		workers = para.New(numWorkers)
	}

	h := &HugoSites{
		running:      cfg.Running,
		multilingual: langConfig,
		multihost:    cfg.Cfg.GetBool("multihost"),
		Sites:        sites,
		workers:      workers,
		numWorkers:   numWorkers,
		init: &hugoSitesInit{
			data:         lazy.New(),
			layouts:      lazy.New(),
			gitInfo:      lazy.New(),
			translations: lazy.New(),
		},
	}

	h.fatalErrorHandler = &fatalErrorHandler{
		h:     h,
		donec: make(chan bool),
	}

	h.init.data.Add(func() (interface{}, error) {
		err := h.loadData(h.PathSpec.BaseFs.Data.Dirs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load data")
		}
		return nil, nil
	})

	h.init.layouts.Add(func() (interface{}, error) {
		for _, s := range h.Sites {
			if err := s.Tmpl().(tpl.TemplateManager).MarkReady(); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})

	h.init.translations.Add(func() (interface{}, error) {
		if len(h.Sites) > 1 {
			allTranslations := pagesToTranslationsMap(h.Sites)
			assignTranslationsToPages(allTranslations, h.Sites)
		}

		return nil, nil
	})

	h.init.gitInfo.Add(func() (interface{}, error) {
		err := h.loadGitInfo()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load Git info")
		}
		return nil, nil
	})

	for _, s := range sites {
		s.h = h
	}

	if err := applyDeps(cfg, sites...); err != nil {
		return nil, errors.Wrap(err, "add site dependencies")
	}

	h.Deps = sites[0].Deps

	// Only needed in server mode.
	// TODO(bep) clean up the running vs watching terms
	if cfg.Running {
		contentChangeTracker = &contentChangeMap{
			pathSpec:      h.PathSpec,
			symContent:    make(map[string]map[string]bool),
			leafBundles:   radix.New(),
			branchBundles: make(map[string]bool),
		}
		h.ContentChanges = contentChangeTracker
	}

	return h, nil
}

func (h *HugoSites) loadGitInfo() error {
	if h.Cfg.GetBool("enableGitInfo") {
		gi, err := newGitInfo(h.Cfg)
		if err != nil {
			h.Log.ERROR.Println("Failed to read Git log:", err)
		} else {
			h.gitInfo = gi
		}
	}
	return nil
}

func applyDeps(cfg deps.DepsCfg, sites ...*Site) error {
	if cfg.TemplateProvider == nil {
		cfg.TemplateProvider = tplimpl.DefaultTemplateProvider
	}

	if cfg.TranslationProvider == nil {
		cfg.TranslationProvider = i18n.NewTranslationProvider()
	}

	var (
		d   *deps.Deps
		err error
	)

	for _, s := range sites {
		if s.Deps != nil {
			continue
		}

		onCreated := func(d *deps.Deps) error {
			s.Deps = d

			// Set up the main publishing chain.
			pub, err := publisher.NewDestinationPublisher(d.PathSpec.BaseFs.PublishFs, s.outputFormatsConfig, s.mediaTypesConfig, cfg.Cfg)

			if err != nil {
				return err
			}
			s.publisher = pub

			if err := s.initializeSiteInfo(); err != nil {
				return err
			}

			d.Site = s.Info

			siteConfig, err := loadSiteConfig(s.language)
			if err != nil {
				return errors.Wrap(err, "load site config")
			}
			s.siteConfigConfig = siteConfig

			pm := &pageMap{
				contentMap: newContentMap(contentMapConfig{
					lang:                 s.Lang(),
					taxonomyConfig:       s.siteCfg.taxonomiesConfig.Values(),
					taxonomyDisabled:     !s.isEnabled(page.KindTaxonomy),
					taxonomyTermDisabled: !s.isEnabled(page.KindTaxonomyTerm),
					pageDisabled:         !s.isEnabled(page.KindPage),
				}),
				s: s,
			}

			s.PageCollections = newPageCollections(pm)

			s.siteRefLinker, err = newSiteRefLinker(s.language, s)
			return err
		}

		cfg.Language = s.language
		cfg.MediaTypes = s.mediaTypesConfig
		cfg.OutputFormats = s.outputFormatsConfig

		if d == nil {
			cfg.WithTemplate = s.withSiteTemplates(cfg.WithTemplate)

			var err error
			d, err = deps.New(cfg)
			if err != nil {
				return errors.Wrap(err, "create deps")
			}

			d.OutputFormatsConfig = s.outputFormatsConfig

			if err := onCreated(d); err != nil {
				return errors.Wrap(err, "on created")
			}

			if err = d.LoadResources(); err != nil {
				return errors.Wrap(err, "load resources")
			}

		} else {
			d, err = d.ForLanguage(cfg, onCreated)
			if err != nil {
				return err
			}
			d.OutputFormatsConfig = s.outputFormatsConfig
		}
	}

	return nil
}

// NewHugoSites creates HugoSites from the given config.
func NewHugoSites(cfg deps.DepsCfg) (*HugoSites, error) {
	sites, err := createSitesFromConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "from config")
	}
	return newHugoSites(cfg, sites...)
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

func createSitesFromConfig(cfg deps.DepsCfg) ([]*Site, error) {

	var (
		sites []*Site
	)

	languages := getLanguages(cfg.Cfg)

	for _, lang := range languages {
		if lang.Disabled {
			continue
		}
		var s *Site
		var err error
		cfg.Language = lang
		s, err = newSite(cfg)

		if err != nil {
			return nil, err
		}

		sites = append(sites, s)
	}

	return sites, nil
}

// Reset resets the sites and template caches etc., making it ready for a full rebuild.
func (h *HugoSites) reset(config *BuildCfg) {
	if config.ResetState {
		for i, s := range h.Sites {
			h.Sites[i] = s.reset()
			if r, ok := s.Fs.Destination.(hugofs.Reseter); ok {
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
		s.Deps.DistinctErrorLog = helpers.NewDistinctLogger(h.Log.ERROR)
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

func (h *HugoSites) createSitesFromConfig(cfg config.Provider) error {
	oldLangs, _ := h.Cfg.Get("languagesSorted").(langs.Languages)

	if err := loadLanguageSettings(h.Cfg, oldLangs); err != nil {
		return err
	}

	depsCfg := deps.DepsCfg{Fs: h.Fs, Cfg: cfg}

	sites, err := createSitesFromConfig(depsCfg)

	if err != nil {
		return err
	}

	langConfig, err := newMultiLingualFromSites(depsCfg.Cfg, sites...)

	if err != nil {
		return err
	}

	h.Sites = sites

	for _, s := range sites {
		s.h = h
	}

	if err := applyDeps(depsCfg, sites...); err != nil {
		return err
	}

	h.Deps = sites[0].Deps

	h.multilingual = langConfig
	h.multihost = h.Deps.Cfg.GetBool("multihost")

	return nil
}

func (h *HugoSites) toSiteInfos() []*SiteInfo {
	infos := make([]*SiteInfo, len(h.Sites))
	for i, s := range h.Sites {
		infos[i] = s.Info
	}
	return infos
}

// BuildCfg holds build options used to, as an example, skip the render step.
type BuildCfg struct {
	// Reset site state before build. Use to force full rebuilds.
	ResetState bool
	// If set, we re-create the sites from the given configuration before a build.
	// This is needed if new languages are added.
	NewConfig config.Provider
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

	testCounters *testCounters
}

// shouldRender is used in the Fast Render Mode to determine if we need to re-render
// a Page: If it is recently visited (the home pages will always be in this set) or changed.
// Note that a page does not have to have a content page / file.
// For regular builds, this will allways return true.
// TODO(bep) rename/work this.
func (cfg *BuildCfg) shouldRender(p *pageState) bool {
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

func (h *HugoSites) renderCrossSitesArtifacts() error {

	if !h.multilingual.enabled() || h.IsMultihost() {
		return nil
	}

	sitemapEnabled := false
	for _, s := range h.Sites {
		if s.isEnabled(kindSitemap) {
			sitemapEnabled = true
			break
		}
	}

	if !sitemapEnabled {
		return nil
	}

	s := h.Sites[0]

	templ := s.lookupLayouts("sitemapindex.xml", "_default/sitemapindex.xml", "_internal/_default/sitemapindex.xml")

	return s.renderAndWriteXML(&s.PathSpec.ProcessingStats.Sitemaps, "sitemapindex",
		s.siteCfg.sitemap.Filename, h.toSiteInfos(), templ)
}

func (h *HugoSites) removePageByFilename(filename string) {
	h.content.withMaps(func(m *pageMap) error {
		m.deleteBundleMatching(func(b *contentNode) bool {
			if b.p == nil {
				return false
			}

			if b.fi == nil {
				return false
			}

			return b.fi.Meta().Filename() == filename
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
	spec := source.NewSourceSpec(h.PathSpec, nil)

	h.data = make(map[string]interface{})
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
	var current map[string]interface{}

	f, err := r.FileInfo().Meta().Open()
	if err != nil {
		return errors.Wrapf(err, "data: failed to open %q:", r.LogicalName())
	}
	defer f.Close()

	// Crawl in data tree to insert data
	current = h.data
	keyParts := strings.Split(r.Dir(), helpers.FilePathSeparator)

	for _, key := range keyParts {
		if key != "" {
			if _, ok := current[key]; !ok {
				current[key] = make(map[string]interface{})
			}
			current = current[key].(map[string]interface{})
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
	case map[string]interface{}:

		switch higherPrecedentData.(type) {
		case nil:
			current[r.BaseFileName()] = data
		case map[string]interface{}:
			// merge maps: insert entries from data for keys that
			// don't already exist in higherPrecedentData
			higherPrecedentMap := higherPrecedentData.(map[string]interface{})
			for key, value := range data.(map[string]interface{}) {
				if _, exists := higherPrecedentMap[key]; exists {
					// this warning could happen if
					// 1. A theme uses the same key; the main data folder wins
					// 2. A sub folder uses the same key: the sub folder wins
					// TODO(bep) figure out a way to detect 2) above and make that a WARN
					h.Log.INFO.Printf("Data for key '%s' in path '%s' is overridden by higher precedence data already in the data tree", key, r.Path())
				} else {
					higherPrecedentMap[key] = value
				}
			}
		default:
			// can't merge: higherPrecedentData is not a map
			h.Log.WARN.Printf("The %T data from '%s' overridden by "+
				"higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
		}

	case []interface{}:
		if higherPrecedentData == nil {
			current[r.BaseFileName()] = data
		} else {
			// we don't merge array data
			h.Log.WARN.Printf("The %T data from '%s' overridden by "+
				"higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
		}

	default:
		h.Log.ERROR.Printf("unexpected data type %T in file %s", data, r.LogicalName())
	}

	return nil
}

func (h *HugoSites) errWithFileContext(err error, f source.File) error {
	fim, ok := f.FileInfo().(hugofs.FileMetaInfo)
	if !ok {
		return err
	}

	realFilename := fim.Meta().Filename()

	err, _ = herrors.WithFileContextForFile(
		err,
		realFilename,
		realFilename,
		h.SourceSpec.Fs.Source,
		herrors.SimpleLineMatcher)

	return err
}

func (h *HugoSites) readData(f source.File) (interface{}, error) {
	file, err := f.FileInfo().Meta().Open()
	if err != nil {
		return nil, errors.Wrap(err, "readData: failed to open data file")
	}
	defer file.Close()
	content := helpers.ReaderToBytes(file)

	format := metadecoders.FormatFromString(f.Extension())
	return metadecoders.Default.Unmarshal(content, format)
}

func (h *HugoSites) findPagesByKindIn(kind string, inPages page.Pages) page.Pages {
	return h.Sites[0].findPagesByKindIn(kind, inPages)
}

func (h *HugoSites) resetPageState() {
	h.content.walkBundles(func(n *contentNode) bool {
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
	h.content.walkBundles(func(n *contentNode) bool {
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
			for id := range idset {
				if idm, ok := s.info.(identity.Manager); ok && idm.Search(id) != nil {
					for _, po := range p.pageOutputs {
						if po.cp != nil {
							po.cp.Reset()
						}
					}
					return false
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
	if !meta.IsSymlink() {
		return
	}
	m.symContentMu.Lock()

	from, to := meta.Filename(), meta.OriginalFilename()
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
