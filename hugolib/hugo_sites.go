// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"errors"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/publisher"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/i18n"
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

	*deps.Deps

	// Keeps track of bundle directories and symlinks to enable partial rebuilding.
	ContentChanges *contentChangeMap

	// If enabled, keeps a revision map for all content.
	gitInfo *gitInfo
}

func (h *HugoSites) siteInfos() SiteInfos {
	infos := make(SiteInfos, len(h.Sites))
	for i, site := range h.Sites {
		infos[i] = &site.Info
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

func (h *HugoSites) LanguageSet() map[string]bool {
	set := make(map[string]bool)
	for _, s := range h.Sites {
		set[s.Language.Lang] = true
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

func (h *HugoSites) langSite() map[string]*Site {
	m := make(map[string]*Site)
	for _, s := range h.Sites {
		m[s.Language.Lang] = s
	}
	return m
}

// GetContentPage finds a Page with content given the absolute filename.
// Returns nil if none found.
func (h *HugoSites) GetContentPage(filename string) *Page {
	for _, s := range h.Sites {
		pos := s.rawAllPages.findPagePosByFilename(filename)
		if pos == -1 {
			continue
		}
		return s.rawAllPages[pos]
	}

	// If not found already, this may be bundled in another content file.
	dir := filepath.Dir(filename)

	for _, s := range h.Sites {
		pos := s.rawAllPages.findPagePosByFilnamePrefix(dir)
		if pos == -1 {
			continue
		}
		return s.rawAllPages[pos]
	}
	return nil
}

// NewHugoSites creates a new collection of sites given the input sites, building
// a language configuration based on those.
func newHugoSites(cfg deps.DepsCfg, sites ...*Site) (*HugoSites, error) {

	if cfg.Language != nil {
		return nil, errors.New("Cannot provide Language in Cfg when sites are provided")
	}

	langConfig, err := newMultiLingualFromSites(cfg.Cfg, sites...)

	if err != nil {
		return nil, err
	}

	var contentChangeTracker *contentChangeMap

	h := &HugoSites{
		running:      cfg.Running,
		multilingual: langConfig,
		multihost:    cfg.Cfg.GetBool("multihost"),
		Sites:        sites}

	for _, s := range sites {
		s.owner = h
	}

	if err := applyDeps(cfg, sites...); err != nil {
		return nil, err
	}

	h.Deps = sites[0].Deps

	// Only needed in server mode.
	// TODO(bep) clean up the running vs watching terms
	if cfg.Running {
		contentChangeTracker = &contentChangeMap{pathSpec: h.PathSpec, symContent: make(map[string]map[string]bool)}
		h.ContentChanges = contentChangeTracker
	}

	if err := h.initGitInfo(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *HugoSites) initGitInfo() error {
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
			s.publisher = publisher.NewDestinationPublisher(d.PathSpec.BaseFs.PublishFs, s.outputFormatsConfig, s.mediaTypesConfig, cfg.Cfg.GetBool("minify"))

			if err := s.initializeSiteInfo(); err != nil {
				return err
			}

			d.Site = &s.Info

			siteConfig, err := loadSiteConfig(s.Language)
			if err != nil {
				return err
			}
			s.siteConfig = siteConfig
			s.siteRefLinker, err = newSiteRefLinker(s.Language, s)
			return err
		}

		cfg.Language = s.Language
		cfg.MediaTypes = s.mediaTypesConfig
		cfg.OutputFormats = s.outputFormatsConfig

		if d == nil {
			cfg.WithTemplate = s.withSiteTemplates(cfg.WithTemplate)

			var err error
			d, err = deps.New(cfg)
			if err != nil {
				return err
			}

			d.OutputFormatsConfig = s.outputFormatsConfig

			if err := onCreated(d); err != nil {
				return err
			}

			if err = d.LoadResources(); err != nil {
				return err
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
		return nil, err
	}
	return newHugoSites(cfg, sites...)
}

func (s *Site) withSiteTemplates(withTemplates ...func(templ tpl.TemplateHandler) error) func(templ tpl.TemplateHandler) error {
	return func(templ tpl.TemplateHandler) error {
		if err := templ.LoadTemplates(""); err != nil {
			return err
		}

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

// Reset resets the sites and template caches, making it ready for a full rebuild.
func (h *HugoSites) reset() {
	for i, s := range h.Sites {
		h.Sites[i] = s.reset()
	}
}

// resetLogs resets the log counters etc. Used to do a new build on the same sites.
func (h *HugoSites) resetLogs() {
	h.Log.Reset()
	loggers.GlobalErrorCounter.Reset()
	for _, s := range h.Sites {
		s.Deps.DistinctErrorLog = helpers.NewDistinctLogger(h.Log.ERROR)
	}
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
		s.owner = h
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
		infos[i] = &s.Info
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

	// Recently visited URLs. This is used for partial re-rendering.
	RecentlyVisited map[string]bool
}

// shouldRender is used in the Fast Render Mode to determine if we need to re-render
// a Page: If it is recently visited (the home pages will always be in this set) or changed.
// Note that a page does not have to have a content page / file.
// For regular builds, this will allways return true.
// TODO(bep) rename/work this.
func (cfg *BuildCfg) shouldRender(p *Page) bool {
	if p.forceRender {
		p.forceRender = false
		return true
	}

	if len(cfg.RecentlyVisited) == 0 {
		return true
	}

	if cfg.RecentlyVisited[p.RelPermalink()] {
		if cfg.PartialReRender {
			_ = p.initMainOutputFormat()
		}
		return true
	}

	if cfg.whatChanged != nil && p.File != nil {
		return cfg.whatChanged.files[p.File.Filename()]
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

	// TODO(bep) DRY
	sitemapDefault := parseSitemap(h.Cfg.GetStringMap("sitemap"))

	s := h.Sites[0]

	smLayouts := []string{"sitemapindex.xml", "_default/sitemapindex.xml", "_internal/_default/sitemapindex.xml"}

	return s.renderAndWriteXML(&s.PathSpec.ProcessingStats.Sitemaps, "sitemapindex",
		sitemapDefault.Filename, h.toSiteInfos(), smLayouts...)
}

func (h *HugoSites) assignMissingTranslations() error {

	// This looks heavy, but it should be a small number of nodes by now.
	allPages := h.findAllPagesByKindNotIn(KindPage)
	for _, nodeType := range []string{KindHome, KindSection, KindTaxonomy, KindTaxonomyTerm} {
		nodes := h.findPagesByKindIn(nodeType, allPages)

		// Assign translations
		for _, t1 := range nodes {
			for _, t2 := range nodes {
				if t1.isNewTranslation(t2) {
					t1.translations = append(t1.translations, t2)
				}
			}
		}
	}

	// Now we can sort the translations.
	for _, p := range allPages {
		if len(p.translations) > 0 {
			pageBy(languagePageSort).Sort(p.translations)
		}
	}
	return nil

}

// createMissingPages creates home page, taxonomies etc. that isnt't created as an
// effect of having a content file.
func (h *HugoSites) createMissingPages() error {
	var newPages Pages

	for _, s := range h.Sites {
		if s.isEnabled(KindHome) {
			// home pages
			home := s.findPagesByKind(KindHome)
			if len(home) > 1 {
				panic("Too many homes")
			}
			if len(home) == 0 {
				n := s.newHomePage()
				s.Pages = append(s.Pages, n)
				newPages = append(newPages, n)
			}
		}

		// Will create content-less root sections.
		newSections := s.assembleSections()
		s.Pages = append(s.Pages, newSections...)
		newPages = append(newPages, newSections...)

		// taxonomy list and terms pages
		taxonomies := s.Language.GetStringMapString("taxonomies")
		if len(taxonomies) > 0 {
			taxonomyPages := s.findPagesByKind(KindTaxonomy)
			taxonomyTermsPages := s.findPagesByKind(KindTaxonomyTerm)
			for _, plural := range taxonomies {
				if s.isEnabled(KindTaxonomyTerm) {
					foundTaxonomyTermsPage := false
					for _, p := range taxonomyTermsPages {
						if p.sectionsPath() == plural {
							foundTaxonomyTermsPage = true
							break
						}
					}

					if !foundTaxonomyTermsPage {
						n := s.newTaxonomyTermsPage(plural)
						s.Pages = append(s.Pages, n)
						newPages = append(newPages, n)
					}
				}

				if s.isEnabled(KindTaxonomy) {
					for key := range s.Taxonomies[plural] {
						foundTaxonomyPage := false
						origKey := key

						if s.Info.preserveTaxonomyNames {
							key = s.PathSpec.MakePathSanitized(key)
						}
						for _, p := range taxonomyPages {
							sectionsPath := p.sectionsPath()

							if !strings.HasPrefix(sectionsPath, plural) {
								continue
							}

							singularKey := strings.TrimPrefix(sectionsPath, plural)
							singularKey = strings.TrimPrefix(singularKey, "/")

							// Some people may have /authors/MaxMustermann etc. as paths.
							// p.sections contains the raw values from the file system.
							// See https://github.com/gohugoio/hugo/issues/4238
							singularKey = s.PathSpec.MakePathSanitized(singularKey)

							if singularKey == key {
								foundTaxonomyPage = true
								break
							}
						}

						if !foundTaxonomyPage {
							n := s.newTaxonomyPage(plural, origKey)
							s.Pages = append(s.Pages, n)
							newPages = append(newPages, n)
						}
					}
				}
			}
		}
	}

	if len(newPages) > 0 {
		// This resorting is unfortunate, but it also needs to be sorted
		// when sections are created.
		first := h.Sites[0]

		first.AllPages = append(first.AllPages, newPages...)

		first.AllPages.sort()

		for _, s := range h.Sites {
			s.Pages.sort()
		}

		for i := 1; i < len(h.Sites); i++ {
			h.Sites[i].AllPages = first.AllPages
		}
	}

	return nil
}

func (h *HugoSites) removePageByFilename(filename string) {
	for _, s := range h.Sites {
		s.removePageFilename(filename)
	}
}

func (h *HugoSites) setupTranslations() {
	for _, s := range h.Sites {
		for _, p := range s.rawAllPages {
			if p.Kind == kindUnknown {
				p.Kind = p.kindFromSections()
			}

			if !p.s.isEnabled(p.Kind) {
				continue
			}

			shouldBuild := p.shouldBuild()
			s.updateBuildStats(p)
			if shouldBuild {
				if p.headless {
					s.headlessPages = append(s.headlessPages, p)
				} else {
					s.Pages = append(s.Pages, p)
				}
			}
		}
	}

	allPages := make(Pages, 0)

	for _, s := range h.Sites {
		allPages = append(allPages, s.Pages...)
	}

	allPages.sort()

	for _, s := range h.Sites {
		s.AllPages = allPages
	}

	// Pull over the collections from the master site
	for i := 1; i < len(h.Sites); i++ {
		h.Sites[i].Data = h.Sites[0].Data
	}

	if len(h.Sites) > 1 {
		allTranslations := pagesToTranslationsMap(allPages)
		assignTranslationsToPages(allTranslations, allPages)
	}
}

func (s *Site) preparePagesForRender(start bool) error {
	for _, p := range s.Pages {
		if err := p.prepareForRender(start); err != nil {
			return err
		}
	}

	for _, p := range s.headlessPages {
		if err := p.prepareForRender(start); err != nil {
			return err
		}
	}

	return nil
}

// Pages returns all pages for all sites.
func (h *HugoSites) Pages() Pages {
	return h.Sites[0].AllPages
}

func handleShortcodes(p *PageWithoutContent, rawContentCopy []byte) ([]byte, error) {
	if p.shortcodeState != nil && p.shortcodeState.contentShortcodes.Len() > 0 {
		p.s.Log.DEBUG.Printf("Replace %d shortcodes in %q", p.shortcodeState.contentShortcodes.Len(), p.BaseFileName())
		err := p.shortcodeState.executeShortcodesForDelta(p)

		if err != nil {

			return rawContentCopy, err
		}

		rawContentCopy, err = replaceShortcodeTokens(rawContentCopy, shortcodePlaceholderPrefix, p.shortcodeState.renderedShortcodes)

		if err != nil {
			p.s.Log.FATAL.Printf("Failed to replace shortcode tokens in %s:\n%s", p.BaseFileName(), err.Error())
		}
	}

	return rawContentCopy, nil
}

func (s *Site) updateBuildStats(page *Page) {
	if page.IsDraft() {
		s.draftCount++
	}

	if page.IsFuture() {
		s.futureCount++
	}

	if page.IsExpired() {
		s.expiredCount++
	}
}

func (h *HugoSites) findPagesByKindNotIn(kind string, inPages Pages) Pages {
	return h.Sites[0].findPagesByKindNotIn(kind, inPages)
}

func (h *HugoSites) findPagesByKindIn(kind string, inPages Pages) Pages {
	return h.Sites[0].findPagesByKindIn(kind, inPages)
}

func (h *HugoSites) findAllPagesByKind(kind string) Pages {
	return h.findPagesByKindIn(kind, h.Sites[0].AllPages)
}

func (h *HugoSites) findAllPagesByKindNotIn(kind string) Pages {
	return h.findPagesByKindNotIn(kind, h.Sites[0].AllPages)
}

func (h *HugoSites) findPagesByShortcode(shortcode string) Pages {
	var pages Pages
	for _, s := range h.Sites {
		pages = append(pages, s.findPagesByShortcode(shortcode)...)
	}
	return pages
}

// Used in partial reloading to determine if the change is in a bundle.
type contentChangeMap struct {
	mu       sync.RWMutex
	branches []string
	leafs    []string

	pathSpec *helpers.PathSpec

	// Hugo supports symlinked content (both directories and files). This
	// can lead to situations where the same file can be referenced from several
	// locations in /content -- which is really cool, but also means we have to
	// go an extra mile to handle changes.
	// This map is only used in watch mode.
	// It maps either file to files or the real dir to a set of content directories where it is in use.
	symContent   map[string]map[string]bool
	symContentMu sync.Mutex
}

func (m *contentChangeMap) add(filename string, tp bundleDirType) {
	m.mu.Lock()
	dir := filepath.Dir(filename) + helpers.FilePathSeparator
	dir = strings.TrimPrefix(dir, ".")
	switch tp {
	case bundleBranch:
		m.branches = append(m.branches, dir)
	case bundleLeaf:
		m.leafs = append(m.leafs, dir)
	default:
		panic("invalid bundle type")
	}
	m.mu.Unlock()
}

// Track the addition of bundle dirs.
func (m *contentChangeMap) handleBundles(b *bundleDirs) {
	for _, bd := range b.bundles {
		m.add(bd.fi.Path(), bd.tp)
	}
}

// resolveAndRemove resolves the given filename to the root folder of a bundle, if relevant.
// It also removes the entry from the map. It will be re-added again by the partial
// build if it still is a bundle.
func (m *contentChangeMap) resolveAndRemove(filename string) (string, string, bundleDirType) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Bundles share resources, so we need to start from the virtual root.
	relPath := m.pathSpec.RelContentDir(filename)
	dir, name := filepath.Split(relPath)
	if !strings.HasSuffix(dir, helpers.FilePathSeparator) {
		dir += helpers.FilePathSeparator
	}

	fileTp, isContent := classifyBundledFile(name)

	// This may be a member of a bundle. Start with branch bundles, the most specific.
	if fileTp == bundleBranch || (fileTp == bundleNot && !isContent) {
		for i, b := range m.branches {
			if b == dir {
				m.branches = append(m.branches[:i], m.branches[i+1:]...)
				return dir, b, bundleBranch
			}
		}
	}

	// And finally the leaf bundles, which can contain anything.
	for i, l := range m.leafs {
		if strings.HasPrefix(dir, l) {
			m.leafs = append(m.leafs[:i], m.leafs[i+1:]...)
			return dir, l, bundleLeaf
		}
	}

	if isContent && fileTp != bundleNot {
		// A new bundle.
		return dir, dir, fileTp
	}

	// Not part of any bundle
	return dir, filename, bundleNot
}

func (m *contentChangeMap) addSymbolicLinkMapping(from, to string) {
	m.symContentMu.Lock()
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
