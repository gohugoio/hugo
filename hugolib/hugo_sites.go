// Copyright 2016-present The Hugo Authors. All rights reserved.
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
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"

	"github.com/spf13/hugo/i18n"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/hugo/tpl/tplimpl"
)

// HugoSites represents the sites to build. Each site represents a language.
type HugoSites struct {
	Sites []*Site

	runMode runmode

	multilingual *Multilingual

	*deps.Deps
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

	h := &HugoSites{
		multilingual: langConfig,
		Sites:        sites}

	for _, s := range sites {
		s.owner = h
	}

	// TODO(bep)
	cfg.Cfg.Set("multilingual", sites[0].multilingualEnabled())

	if err := applyDepsIfNeeded(cfg, sites...); err != nil {
		return nil, err
	}

	h.Deps = sites[0].Deps

	return h, nil
}

func applyDepsIfNeeded(cfg deps.DepsCfg, sites ...*Site) error {
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

		if d == nil {
			cfg.Language = s.Language
			cfg.WithTemplate = s.withSiteTemplates(cfg.WithTemplate)

			var err error
			d, err = deps.New(cfg)
			if err != nil {
				return err
			}

			d.OutputFormatsConfig = s.outputFormatsConfig
			s.Deps = d

			if err = d.LoadResources(); err != nil {
				return err
			}

		} else {
			d, err = d.ForLanguage(s.Language)
			if err != nil {
				return err
			}
			d.OutputFormatsConfig = s.outputFormatsConfig
			s.Deps = d
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
		templ.LoadTemplates(s.PathSpec.GetLayoutDirPath(), "")
		if s.PathSpec.ThemeSet() {
			templ.LoadTemplates(s.PathSpec.GetThemeDir()+"/layouts", "theme")
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

	multilingual := cfg.Cfg.GetStringMap("languages")

	if len(multilingual) == 0 {
		l := helpers.NewDefaultLanguage(cfg.Cfg)
		cfg.Language = l
		s, err := newSite(cfg)
		if err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}

	if len(multilingual) > 0 {
		var err error

		languages, err := toSortedLanguages(cfg.Cfg, multilingual)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse multilingual config: %s", err)
		}

		for _, lang := range languages {
			var s *Site
			var err error
			cfg.Language = lang
			s, err = newSite(cfg)

			if err != nil {
				return nil, err
			}

			sites = append(sites, s)
		}
	}

	return sites, nil
}

// Reset resets the sites and template caches, making it ready for a full rebuild.
func (h *HugoSites) reset() {
	for i, s := range h.Sites {
		h.Sites[i] = s.reset()
	}
}

func (h *HugoSites) createSitesFromConfig() error {

	depsCfg := deps.DepsCfg{Fs: h.Fs, Cfg: h.Cfg}
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

	if err := applyDepsIfNeeded(depsCfg, sites...); err != nil {
		return err
	}

	h.Deps = sites[0].Deps

	h.multilingual = langConfig

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
	// Whether we are in watch (server) mode
	Watching bool
	// Print build stats at the end of a build
	PrintStats bool
	// Reset site state before build. Use to force full rebuilds.
	ResetState bool
	// Re-creates the sites from configuration before a build.
	// This is needed if new languages are added.
	CreateSitesFromConfig bool
	// Skip rendering. Useful for testing.
	SkipRender bool
	// Use this to indicate what changed (for rebuilds).
	whatChanged *whatChanged
}

func (h *HugoSites) renderCrossSitesArtifacts() error {

	if !h.multilingual.enabled() {
		return nil
	}

	if h.Cfg.GetBool("disableSitemap") {
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

	return s.renderAndWriteXML("sitemapindex",
		sitemapDefault.Filename, h.toSiteInfos(), s.appendThemeTemplates(smLayouts)...)
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
						if p.sections[0] == plural {
							foundTaxonomyTermsPage = true
							break
						}
					}

					if !foundTaxonomyTermsPage {
						foundTaxonomyTermsPage = true
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
							if p.sections[0] == plural && p.sections[1] == key {
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

		first.AllPages.Sort()

		for _, s := range h.Sites {
			s.Pages.Sort()
		}

		for i := 1; i < len(h.Sites); i++ {
			h.Sites[i].AllPages = first.AllPages
		}
	}

	return nil
}

func (s *Site) assignSiteByLanguage(p *Page) {

	pageLang := p.Lang()

	if pageLang == "" {
		panic("Page language missing: " + p.Title)
	}

	for _, site := range s.owner.Sites {
		if strings.HasPrefix(site.Language.Lang, pageLang) {
			p.s = site
			p.Site = &site.Info
			return
		}
	}

}

func (h *HugoSites) setupTranslations() {

	master := h.Sites[0]

	for _, p := range master.rawAllPages {
		if p.Lang() == "" {
			panic("Page language missing: " + p.Title)
		}

		if p.Kind == kindUnknown {
			p.Kind = p.s.kindFromSections(p.sections)
		}

		if !p.s.isEnabled(p.Kind) {
			continue
		}

		shouldBuild := p.shouldBuild()

		for i, site := range h.Sites {
			// The site is assigned by language when read.
			if site == p.s {
				site.updateBuildStats(p)
				if shouldBuild {
					site.Pages = append(site.Pages, p)
				}
			}

			if !shouldBuild {
				continue
			}

			if i == 0 {
				site.AllPages = append(site.AllPages, p)
			}
		}

	}

	// Pull over the collections from the master site
	for i := 1; i < len(h.Sites); i++ {
		h.Sites[i].AllPages = h.Sites[0].AllPages
		h.Sites[i].Data = h.Sites[0].Data
	}

	if len(h.Sites) > 1 {
		pages := h.Sites[0].AllPages
		allTranslations := pagesToTranslationsMap(pages)
		assignTranslationsToPages(allTranslations, pages)
	}
}

func (s *Site) preparePagesForRender(cfg *BuildCfg) {

	pageChan := make(chan *Page)
	wg := &sync.WaitGroup{}
	numWorkers := getGoMaxProcs() * 4

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(pages <-chan *Page, wg *sync.WaitGroup) {
			defer wg.Done()
			for p := range pages {
				if !p.shouldRenderTo(s.rc.Format) {
					// No need to prepare
					continue
				}
				var shortcodeUpdate bool
				if p.shortcodeState != nil {
					shortcodeUpdate = p.shortcodeState.updateDelta()
				}

				if !shortcodeUpdate && !cfg.whatChanged.other && p.rendered {
					// No need to process it again.
					continue
				}

				// If we got this far it means that this is either a new Page pointer
				// or a template or similar has changed so wee need to do a rerendering
				// of the shortcodes etc.

				// Mark it as rendered
				p.rendered = true

				// If in watch mode or if we have multiple output formats,
				// we need to keep the original so we can
				// potentially repeat this process on rebuild.
				needsACopy := cfg.Watching || len(p.outputFormats) > 1
				var workContentCopy []byte
				if needsACopy {
					workContentCopy = make([]byte, len(p.workContent))
					copy(workContentCopy, p.workContent)
				} else {
					// Just reuse the same slice.
					workContentCopy = p.workContent
				}

				if p.Markup == "markdown" {
					tmpContent, tmpTableOfContents := helpers.ExtractTOC(workContentCopy)
					p.TableOfContents = helpers.BytesToHTML(tmpTableOfContents)
					workContentCopy = tmpContent
				}

				var err error
				if workContentCopy, err = handleShortcodes(p, workContentCopy); err != nil {
					s.Log.ERROR.Printf("Failed to handle shortcodes for page %s: %s", p.BaseFileName(), err)
				}

				if p.Markup != "html" {

					// Now we know enough to create a summary of the page and count some words
					summaryContent, err := p.setUserDefinedSummaryIfProvided(workContentCopy)

					if err != nil {
						s.Log.ERROR.Printf("Failed to set user defined summary for page %q: %s", p.Path(), err)
					} else if summaryContent != nil {
						workContentCopy = summaryContent.content
					}

					p.Content = helpers.BytesToHTML(workContentCopy)

					if summaryContent == nil {
						if err := p.setAutoSummary(); err != nil {
							s.Log.ERROR.Printf("Failed to set user auto summary for page %q: %s", p.pathOrTitle(), err)
						}
					}

				} else {
					p.Content = helpers.BytesToHTML(workContentCopy)
				}

				//analyze for raw stats
				p.analyzePage()

			}
		}(pageChan, wg)
	}

	for _, p := range s.Pages {
		pageChan <- p
	}

	close(pageChan)

	wg.Wait()

}

// Pages returns all pages for all sites.
func (h *HugoSites) Pages() Pages {
	return h.Sites[0].AllPages
}

func handleShortcodes(p *Page, rawContentCopy []byte) ([]byte, error) {
	if p.shortcodeState != nil && len(p.shortcodeState.contentShortcodes) > 0 {
		p.s.Log.DEBUG.Printf("Replace %d shortcodes in %q", len(p.shortcodeState.contentShortcodes), p.BaseFileName())
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
