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
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/hugo/helpers"

	"github.com/spf13/viper"

	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	jww "github.com/spf13/jwalterweatherman"
)

// HugoSites represents the sites to build. Each site represents a language.
type HugoSites struct {
	Sites []*Site

	tmpl    tpl.Template
	runMode runmode

	multilingual *Multilingual
}

// NewHugoSites creates a new collection of sites given the input sites, building
// a language configuration based on those.
func newHugoSites(sites ...*Site) (*HugoSites, error) {
	langConfig, err := newMultiLingualFromSites(sites...)

	if err != nil {
		return nil, err
	}

	h := &HugoSites{multilingual: langConfig, Sites: sites}

	for _, s := range sites {
		s.owner = h
	}
	return h, nil
}

// NewHugoSitesFromConfiguration creates HugoSites from the global Viper config.
func NewHugoSitesFromConfiguration() (*HugoSites, error) {
	sites, err := createSitesFromConfig()
	if err != nil {
		return nil, err
	}
	return newHugoSites(sites...)
}

func createSitesFromConfig() ([]*Site, error) {
	var sites []*Site
	multilingual := viper.GetStringMap("languages")
	if len(multilingual) == 0 {
		sites = append(sites, newSite(helpers.NewDefaultLanguage()))
	}

	if len(multilingual) > 0 {
		var err error

		languages, err := toSortedLanguages(multilingual)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse multilingual config: %s", err)
		}

		for _, lang := range languages {
			sites = append(sites, newSite(lang))
		}

	}

	return sites, nil
}

// Reset resets the sites and template caches, making it ready for a full rebuild.
func (h *HugoSites) reset() {
	for i, s := range h.Sites {
		h.Sites[i] = s.reset()
	}

	tpl.ResetCaches()
}

func (h *HugoSites) createSitesFromConfig() error {

	sites, err := createSitesFromConfig()

	if err != nil {
		return err
	}

	langConfig, err := newMultiLingualFromSites(sites...)

	if err != nil {
		return err
	}

	h.Sites = sites

	for _, s := range sites {
		s.owner = h
	}

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
	// Use this to add templates to use for rendering.
	// Useful for testing.
	withTemplate func(templ tpl.Template) error
	// Use this to indicate what changed (for rebuilds).
	whatChanged *whatChanged
}

func (h *HugoSites) renderCrossSitesArtifacts() error {

	if !h.multilingual.enabled() {
		return nil
	}

	if viper.GetBool("disableSitemap") {
		return nil
	}

	// TODO(bep) DRY
	sitemapDefault := parseSitemap(viper.GetStringMap("sitemap"))

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

		// taxonomy list and terms pages
		taxonomies := s.Language.GetStringMapString("taxonomies")
		if len(taxonomies) > 0 {
			taxonomyPages := s.findPagesByKind(KindTaxonomy)
			taxonomyTermsPages := s.findPagesByKind(KindTaxonomyTerm)
			for _, plural := range taxonomies {
				tax := s.Taxonomies[plural]
				foundTaxonomyPage := false
				foundTaxonomyTermsPage := false
				for key := range tax {
					if s.Info.preserveTaxonomyNames {
						key = s.Info.pathSpec.MakePathSanitized(key)
					}
					for _, p := range taxonomyPages {
						if p.sections[0] == plural && p.sections[1] == key {
							foundTaxonomyPage = true
							break
						}
					}
					for _, p := range taxonomyTermsPages {
						if p.sections[0] == plural {
							foundTaxonomyTermsPage = true
							break
						}
					}
					if !foundTaxonomyPage {
						n := s.newTaxonomyPage(plural, key)
						s.Pages = append(s.Pages, n)
						newPages = append(newPages, n)
					}

					if !foundTaxonomyTermsPage {
						foundTaxonomyTermsPage = true
						n := s.newTaxonomyTermsPage(plural)
						s.Pages = append(s.Pages, n)
						newPages = append(newPages, n)
					}
				}
			}
		}

		sectionPages := s.findPagesByKind(KindSection)
		if len(sectionPages) < len(s.Sections) {
			for name, section := range s.Sections {
				// A section may be created for the root content folder if a
				// content file is placed there.
				// We cannot create a section node for that, because
				// that would overwrite the home page.
				if name == "" {
					continue
				}
				foundSection := false
				for _, sectionPage := range sectionPages {
					if sectionPage.sections[0] == name {
						foundSection = true
						break
					}
				}
				if !foundSection {
					n := s.newSectionPage(name, section)
					s.Pages = append(s.Pages, n)
					newPages = append(newPages, n)
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

func (h *HugoSites) setupTranslations() {

	master := h.Sites[0]

	for _, p := range master.rawAllPages {
		if p.Lang() == "" {
			panic("Page language missing: " + p.Title)
		}

		shouldBuild := p.shouldBuild()

		for i, site := range h.Sites {
			if strings.HasPrefix(site.Language.Lang, p.Lang()) {
				site.updateBuildStats(p)
				if shouldBuild {
					site.Pages = append(site.Pages, p)
					p.Site = &site.Info
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
		allTranslations := pagesToTranslationsMap(h.multilingual, pages)
		assignTranslationsToPages(allTranslations, pages)
	}
}

func (s *Site) preparePagesForRender(cfg *BuildCfg) {
	pageChan := make(chan *Page)
	wg := &sync.WaitGroup{}

	for i := 0; i < getGoMaxProcs()*4; i++ {
		wg.Add(1)
		go func(pages <-chan *Page, wg *sync.WaitGroup) {
			defer wg.Done()
			for p := range pages {

				if !cfg.whatChanged.other && p.rendered {
					// No need to process it again.
					continue
				}

				// If we got this far it means that this is either a new Page pointer
				// or a template or similar has changed so wee need to do a rerendering
				// of the shortcodes etc.

				// Mark it as rendered
				p.rendered = true

				// If in watch mode, we need to keep the original so we can
				// repeat this process on rebuild.
				var workContentCopy []byte
				if cfg.Watching {
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
				if workContentCopy, err = handleShortcodes(p, s.owner.tmpl, workContentCopy); err != nil {
					jww.ERROR.Printf("Failed to handle shortcodes for page %s: %s", p.BaseFileName(), err)
				}

				if p.Markup != "html" {

					// Now we know enough to create a summary of the page and count some words
					summaryContent, err := p.setUserDefinedSummaryIfProvided(workContentCopy)

					if err != nil {
						jww.ERROR.Printf("Failed to set user defined summary for page %q: %s", p.Path(), err)
					} else if summaryContent != nil {
						workContentCopy = summaryContent.content
					}

					p.Content = helpers.BytesToHTML(workContentCopy)

					if summaryContent == nil {
						p.setAutoSummary()
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

func handleShortcodes(p *Page, t tpl.Template, rawContentCopy []byte) ([]byte, error) {
	if len(p.contentShortCodes) > 0 {
		jww.DEBUG.Printf("Replace %d shortcodes in %q", len(p.contentShortCodes), p.BaseFileName())
		shortcodes, err := executeShortcodeFuncMap(p.contentShortCodes)

		if err != nil {
			return rawContentCopy, err
		}

		rawContentCopy, err = replaceShortcodeTokens(rawContentCopy, shortcodePlaceholderPrefix, shortcodes)

		if err != nil {
			jww.FATAL.Printf("Failed to replace short code tokens in %s:\n%s", p.BaseFileName(), err.Error())
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

// Convenience func used in tests to build a single site/language excluding render phase.
func buildSiteSkipRender(s *Site, additionalTemplates ...string) error {
	return doBuildSite(s, false, additionalTemplates...)
}

// Convenience func used in tests to build a single site/language including render phase.
func buildAndRenderSite(s *Site, additionalTemplates ...string) error {
	return doBuildSite(s, true, additionalTemplates...)
}

// Convenience func used in tests to build a single site/language.
func doBuildSite(s *Site, render bool, additionalTemplates ...string) error {
	if s.PageCollections == nil {
		s.PageCollections = newPageCollections()
	}
	sites, err := newHugoSites(s)
	if err != nil {
		return err
	}

	addTemplates := func(templ tpl.Template) error {
		for i := 0; i < len(additionalTemplates); i += 2 {
			err := templ.AddTemplate(additionalTemplates[i], additionalTemplates[i+1])
			if err != nil {
				return err
			}
		}
		return nil
	}

	config := BuildCfg{SkipRender: !render, withTemplate: addTemplates}
	return sites.Build(config)
}

// Convenience func used in tests.
func newHugoSitesFromSourceAndLanguages(input []source.ByteSource, languages helpers.Languages) (*HugoSites, error) {
	if len(languages) == 0 {
		panic("Must provide at least one language")
	}
	first := &Site{
		Source:   &source.InMemorySource{ByteSource: input},
		Language: languages[0],
	}
	if len(languages) == 1 {
		return newHugoSites(first)
	}

	sites := make([]*Site, len(languages))
	sites[0] = first
	for i := 1; i < len(languages); i++ {
		sites[i] = &Site{Language: languages[i]}
	}

	return newHugoSites(sites...)

}

// Convenience func used in tests.
func newHugoSitesDefaultLanguage() (*HugoSites, error) {
	return newHugoSitesFromSourceAndLanguages(nil, helpers.Languages{helpers.NewDefaultLanguage()})
}
