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
	"errors"
	"fmt"
	"html/template"
	"os"
	"sort"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/langs/i18n"
	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/publisher"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/page/siteidentities"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/tpl"
	"github.com/gohugoio/hugo/tpl/tplimpl"
)

var _ page.Site = (*Site)(nil)

type Site struct {
	conf      *allconfig.Config
	language  *langs.Language
	languagei int
	pageMap   *pageMap

	// The owning container.
	h *HugoSites

	*deps.Deps

	// Page navigation.
	*pageFinder
	taxonomies page.TaxonomyList
	menus      navigation.Menus

	// Shortcut to the home page. Note that this may be nil if
	// home page, for some odd reason, is disabled.
	home *pageState

	// The last modification date of this site.
	lastmod time.Time

	relatedDocsHandler *page.RelatedDocsHandler
	siteRefLinker
	publisher          publisher.Publisher
	frontmatterHandler pagemeta.FrontMatterHandler

	// We render each site for all the relevant output formats in serial with
	// this rendering context pointing to the current one.
	rc *siteRenderingContext

	// The output formats that we need to render this site in. This slice
	// will be fixed once set.
	// This will be the union of Site.Pages' outputFormats.
	// This slice will be sorted.
	renderFormats output.Formats

	// Lazily loaded site dependencies
	init *siteInit
}

func (s *Site) Debug() {
	fmt.Println("Debugging site", s.Lang(), "=>")
	// fmt.Println(s.pageMap.testDump())
}

// NewHugoSites creates HugoSites from the given config.
func NewHugoSites(cfg deps.DepsCfg) (*HugoSites, error) {
	conf := cfg.Configs.GetFirstLanguageConfig()

	var logger loggers.Logger
	if cfg.TestLogger != nil {
		logger = cfg.TestLogger
	} else {
		var logHookLast func(e *logg.Entry) error
		if cfg.Configs.Base.PanicOnWarning {
			logHookLast = loggers.PanicOnWarningHook
		}
		if cfg.LogOut == nil {
			cfg.LogOut = os.Stdout
		}
		if cfg.LogLevel == 0 {
			cfg.LogLevel = logg.LevelWarn
		}

		logOpts := loggers.Options{
			Level:              cfg.LogLevel,
			DistinctLevel:      logg.LevelWarn, // This will drop duplicate log warning and errors.
			HandlerPost:        logHookLast,
			Stdout:             cfg.LogOut,
			Stderr:             cfg.LogOut,
			StoreErrors:        conf.Running(),
			SuppressStatements: conf.IgnoredErrors(),
		}
		logger = loggers.New(logOpts)
	}

	memCache := dynacache.New(dynacache.Options{Running: conf.Running(), Log: logger})

	firstSiteDeps := &deps.Deps{
		Fs:                  cfg.Fs,
		Log:                 logger,
		Conf:                conf,
		MemCache:            memCache,
		TemplateProvider:    tplimpl.DefaultTemplateProvider,
		TranslationProvider: i18n.NewTranslationProvider(),
	}

	if err := firstSiteDeps.Init(); err != nil {
		return nil, err
	}

	confm := cfg.Configs
	var sites []*Site

	ns := &contentNodeShifter{
		numLanguages: len(confm.Languages),
	}

	treeConfig := doctree.Config[contentNodeI]{
		Shifter: ns,
	}

	pageTrees := &pageTrees{
		treePages: doctree.New(
			treeConfig,
		),
		treeResources: doctree.New(
			treeConfig,
		),
		treeTaxonomyEntries: doctree.NewTreeShiftTree[*weightedContentNode](doctree.DimensionLanguage.Index(), len(confm.Languages)),
	}

	pageTrees.treePagesResources = doctree.WalkableTrees[contentNodeI]{
		pageTrees.treePages,
		pageTrees.treeResources,
	}

	pageTrees.resourceTrees = doctree.MutableTrees{
		pageTrees.treeResources,
	}

	for i, confp := range confm.ConfigLangs() {
		language := confp.Language()
		if language.Disabled {
			continue
		}
		k := language.Lang
		conf := confm.LanguageConfigMap[k]
		frontmatterHandler, err := pagemeta.NewFrontmatterHandler(firstSiteDeps.Log, conf.Frontmatter)
		if err != nil {
			return nil, err
		}

		langs.SetParams(language, conf.Params)

		s := &Site{
			conf:               conf,
			language:           language,
			languagei:          i,
			frontmatterHandler: frontmatterHandler,
		}

		if i == 0 {
			firstSiteDeps.Site = s
			s.Deps = firstSiteDeps
		} else {
			d, err := firstSiteDeps.Clone(s, confp)
			if err != nil {
				return nil, err
			}
			s.Deps = d
		}

		s.pageMap = newPageMap(i, s, memCache, pageTrees)

		s.pageFinder = newPageFinder(s.pageMap)
		s.siteRefLinker, err = newSiteRefLinker(s)
		if err != nil {
			return nil, err
		}
		// Set up the main publishing chain.
		pub, err := publisher.NewDestinationPublisher(
			firstSiteDeps.ResourceSpec,
			s.conf.OutputFormats.Config,
			s.conf.MediaTypes.Config,
		)
		if err != nil {
			return nil, err
		}

		s.publisher = pub
		s.relatedDocsHandler = page.NewRelatedDocsHandler(s.conf.Related)
		// Site deps end.

		s.prepareInits()
		sites = append(sites, s)
	}

	if len(sites) == 0 {
		return nil, errors.New("no sites to build")
	}

	// Pull the default content language to the top, then sort the sites by language weight (if set) or lang.
	defaultContentLanguage := confm.Base.DefaultContentLanguage
	sort.Slice(sites, func(i, j int) bool {
		li := sites[i].language
		lj := sites[j].language
		if li.Lang == defaultContentLanguage {
			return true
		}

		if lj.Lang == defaultContentLanguage {
			return false
		}

		if li.Weight != lj.Weight {
			return li.Weight < lj.Weight
		}
		return li.Lang < lj.Lang
	})

	h, err := newHugoSites(cfg, firstSiteDeps, pageTrees, sites)
	if err == nil && h == nil {
		panic("hugo: newHugoSitesNew returned nil error and nil HugoSites")
	}

	return h, err
}

func newHugoSites(cfg deps.DepsCfg, d *deps.Deps, pageTrees *pageTrees, sites []*Site) (*HugoSites, error) {
	numWorkers := config.GetNumWorkerMultiplier()
	numWorkersSite := numWorkers
	if numWorkersSite > len(sites) {
		numWorkersSite = len(sites)
	}
	workersSite := para.New(numWorkersSite)

	h := &HugoSites{
		Sites:           sites,
		Deps:            sites[0].Deps,
		Configs:         cfg.Configs,
		workersSite:     workersSite,
		numWorkersSites: numWorkers,
		numWorkers:      numWorkers,
		pageTrees:       pageTrees,
		cachePages: dynacache.GetOrCreatePartition[string,
			page.Pages](d.MemCache, "/pags/all",
			dynacache.OptionsPartition{Weight: 10, ClearWhen: dynacache.ClearOnRebuild},
		),
		cacheContentSource:      dynacache.GetOrCreatePartition[string, *resources.StaleValue[[]byte]](d.MemCache, "/cont/src", dynacache.OptionsPartition{Weight: 70, ClearWhen: dynacache.ClearOnChange}),
		translationKeyPages:     maps.NewSliceCache[page.Page](),
		currentSite:             sites[0],
		skipRebuildForFilenames: make(map[string]bool),
		init: &hugoSitesInit{
			data:    lazy.New(),
			layouts: lazy.New(),
			gitInfo: lazy.New(),
		},
	}

	// Assemble dependencies to be used in hugo.Deps.
	var dependencies []*hugo.Dependency
	var depFromMod func(m modules.Module) *hugo.Dependency
	depFromMod = func(m modules.Module) *hugo.Dependency {
		dep := &hugo.Dependency{
			Path:    m.Path(),
			Version: m.Version(),
			Time:    m.Time(),
			Vendor:  m.Vendor(),
		}

		// These are pointers, but this all came from JSON so there's no recursive navigation,
		// so just create new values.
		if m.Replace() != nil {
			dep.Replace = depFromMod(m.Replace())
		}
		if m.Owner() != nil {
			dep.Owner = depFromMod(m.Owner())
		}
		return dep
	}
	for _, m := range d.Paths.AllModules() {
		dependencies = append(dependencies, depFromMod(m))
	}

	h.hugoInfo = hugo.NewInfo(h.Configs.GetFirstLanguageConfig(), dependencies)

	var prototype *deps.Deps
	for i, s := range sites {
		s.h = h
		if err := s.Deps.Compile(prototype); err != nil {
			return nil, err
		}
		if i == 0 {
			prototype = s.Deps
		}
	}

	h.fatalErrorHandler = &fatalErrorHandler{
		h:     h,
		donec: make(chan bool),
	}

	h.init.data.Add(func(context.Context) (any, error) {
		err := h.loadData()
		if err != nil {
			return nil, fmt.Errorf("failed to load data: %w", err)
		}
		return nil, nil
	})

	h.init.layouts.Add(func(context.Context) (any, error) {
		for _, s := range h.Sites {
			if err := s.Tmpl().(tpl.TemplateManager).MarkReady(); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})

	h.init.gitInfo.Add(func(context.Context) (any, error) {
		err := h.loadGitInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to load Git info: %w", err)
		}
		return nil, nil
	})

	return h, nil
}

// Returns true if we're running in a server.
// Deprecated: use hugo.IsServer instead
func (s *Site) IsServer() bool {
	hugo.Deprecate(".Site.IsServer", "Use hugo.IsServer instead.", "v0.120.0")
	return s.conf.Internal.Running
}

// Returns the server port.
func (s *Site) ServerPort() int {
	return s.conf.C.BaseURL.Port()
}

// Returns the configured title for this Site.
func (s *Site) Title() string {
	return s.conf.Title
}

func (s *Site) Copyright() string {
	return s.conf.Copyright
}

func (s *Site) RSSLink() template.URL {
	hugo.Deprecate("Site.RSSLink", "Use the Output Format's Permalink method instead, e.g. .OutputFormats.Get \"RSS\".Permalink", "v0.114.0")
	rssOutputFormat := s.home.OutputFormats().Get("rss")
	return template.URL(rssOutputFormat.Permalink())
}

func (s *Site) Config() page.SiteConfig {
	return page.SiteConfig{
		Privacy:  s.conf.Privacy,
		Services: s.conf.Services,
	}
}

func (s *Site) LanguageCode() string {
	return s.Language().LanguageCode()
}

// Returns all Sites for all languages.
func (s *Site) Sites() page.Sites {
	sites := make(page.Sites, len(s.h.Sites))
	for i, s := range s.h.Sites {
		sites[i] = s.Site()
	}
	return sites
}

// Returns Site currently rendering.
func (s *Site) Current() page.Site {
	return s.h.currentSite
}

// MainSections returns the list of main sections.
func (s *Site) MainSections() []string {
	return s.conf.C.MainSections
}

// Returns a struct with some information about the build.
func (s *Site) Hugo() hugo.HugoInfo {
	if s.h == nil || s.h.hugoInfo.Environment == "" {
		panic("site: hugo: hugoInfo not initialized")
	}
	return s.h.hugoInfo
}

// Returns the BaseURL for this Site.
func (s *Site) BaseURL() string {
	return s.conf.C.BaseURL.WithPath
}

// Returns the last modification date of the content.
// Deprecated: Use .Lastmod instead.
func (s *Site) LastChange() time.Time {
	return s.lastmod
}

// Returns the last modification date of the content.
func (s *Site) Lastmod() time.Time {
	return s.lastmod
}

// Returns the Params configured for this site.
func (s *Site) Params() maps.Params {
	return s.conf.Params
}

func (s *Site) Author() map[string]any {
	return s.conf.Author
}

func (s *Site) Authors() page.AuthorList {
	return page.AuthorList{}
}

func (s *Site) Social() map[string]string {
	return s.conf.Social
}

// Deprecated: Use .Site.Config.Services.Disqus.Shortname instead
func (s *Site) DisqusShortname() string {
	hugo.Deprecate(".Site.DisqusShortname", "Use .Site.Config.Services.Disqus.Shortname instead.", "v0.120.0")
	return s.Config().Services.Disqus.Shortname
}

// Deprecated: Use .Site.Config.Services.GoogleAnalytics.ID instead
func (s *Site) GoogleAnalytics() string {
	hugo.Deprecate(".Site.GoogleAnalytics", "Use .Site.Config.Services.GoogleAnalytics.ID instead.", "v0.120.0")
	return s.Config().Services.GoogleAnalytics.ID
}

func (s *Site) Param(key any) (any, error) {
	return resource.Param(s, nil, key)
}

// Returns a map of all the data inside /data.
func (s *Site) Data() map[string]any {
	return s.s.h.Data()
}

func (s *Site) BuildDrafts() bool {
	return s.conf.BuildDrafts
}

func (s *Site) IsMultiLingual() bool {
	return s.h.isMultiLingual()
}

func (s *Site) LanguagePrefix() string {
	prefix := s.GetLanguagePrefix()
	if prefix == "" {
		return ""
	}
	return "/" + prefix
}

func (s *Site) Site() page.Site {
	return page.WrapSite(s)
}

func (s *Site) ForEeachIdentityByName(name string, f func(identity.Identity) bool) {
	if id, found := siteidentities.FromString(name); found {
		if f(id) {
			return
		}
	}
}

// Pages returns all pages.
// This is for the current language only.
func (s *Site) Pages() page.Pages {
	return s.pageMap.getPagesInSection(
		pageMapQueryPagesInSection{
			pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
				Path:    "",
				KeyPart: "global",
				Include: pagePredicates.ShouldListGlobal,
			},
			Recursive:   true,
			IncludeSelf: true,
		},
	)
}

// RegularPages returns all the regular pages.
// This is for the current language only.
func (s *Site) RegularPages() page.Pages {
	return s.pageMap.getPagesInSection(
		pageMapQueryPagesInSection{
			pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
				Path:    "",
				KeyPart: "global",
				Include: pagePredicates.ShouldListGlobal.And(pagePredicates.KindPage),
			},
			Recursive: true,
		},
	)
}

// AllPages returns all pages for all sites.
func (s *Site) AllPages() page.Pages {
	return s.h.Pages()
}

// AllRegularPages returns all regular pages for all sites.
func (s *Site) AllRegularPages() page.Pages {
	return s.h.RegularPages()
}
