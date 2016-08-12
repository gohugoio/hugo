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
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/hugo/helpers"

	"github.com/spf13/viper"

	"github.com/fsnotify/fsnotify"
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

	// Maps internalID to a set of nodes.
	nodeMap   map[string]Nodes
	nodeMapMu sync.Mutex
}

// NewHugoSites creates a new collection of sites given the input sites, building
// a language configuration based on those.
func newHugoSites(sites ...*Site) (*HugoSites, error) {
	langConfig, err := newMultiLingualFromSites(sites...)

	if err != nil {
		return nil, err
	}

	h := &HugoSites{multilingual: langConfig, Sites: sites, nodeMap: make(map[string]Nodes)}

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
	multilingual := viper.GetStringMap("Languages")
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

func (h *HugoSites) addNode(nodeID string, node *Node) {
	h.nodeMapMu.Lock()

	if nodes, ok := h.nodeMap[nodeID]; ok {
		h.nodeMap[nodeID] = append(nodes, node)
	} else {
		h.nodeMap[nodeID] = Nodes{node}
	}
	h.nodeMapMu.Unlock()
}

func (h *HugoSites) getNodes(nodeID string) Nodes {
	// At this point it is read only, so no need to lock.
	if nodeID != "" {
		if nodes, ok := h.nodeMap[nodeID]; ok {
			return nodes
		}
	}
	// Paginator pages will not have related nodes.
	return Nodes{}
}

// Reset resets the sites, making it ready for a full rebuild.
func (h *HugoSites) reset() {
	h.nodeMap = make(map[string]Nodes)
	for i, s := range h.Sites {
		h.Sites[i] = s.reset()
	}
}

func (h *HugoSites) reCreateFromConfig() error {
	h.nodeMap = make(map[string]Nodes)

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
}

// Build builds all sites.
func (h *HugoSites) Build(config BuildCfg) error {

	t0 := time.Now()

	if config.ResetState {
		h.reset()
	}

	if config.CreateSitesFromConfig {
		if err := h.reCreateFromConfig(); err != nil {
			return err
		}
	}

	h.runMode.Watching = config.Watching

	// We should probably refactor the Site and pull up most of the logic from there to here,
	// but that seems like a daunting task.
	// So for now, if there are more than one site (language),
	// we pre-process the first one, then configure all the sites based on that.
	firstSite := h.Sites[0]

	if err := firstSite.preProcess(config); err != nil {
		return err
	}

	h.setupTranslations(firstSite)

	if len(h.Sites) > 1 {
		// Initialize the rest
		for _, site := range h.Sites[1:] {
			site.initializeSiteInfo()
		}
	}

	for _, s := range h.Sites {
		if err := s.postProcess(); err != nil {
			return err
		}
	}

	if err := h.preRender(); err != nil {
		return err
	}

	if !config.SkipRender {
		for _, s := range h.Sites {

			if err := s.render(); err != nil {
				return err
			}

			if config.PrintStats {
				s.Stats()
			}
		}

		if err := h.render(); err != nil {
			return err
		}
	}

	if config.PrintStats {
		jww.FEEDBACK.Printf("total in %v ms\n", int(1000*time.Since(t0).Seconds()))
	}

	return nil

}

// Rebuild rebuilds all sites.
func (h *HugoSites) Rebuild(config BuildCfg, events ...fsnotify.Event) error {
	t0 := time.Now()

	if config.CreateSitesFromConfig {
		return errors.New("Rebuild does not support 'CreateSitesFromConfig'. Use Build.")
	}

	if config.ResetState {
		return errors.New("Rebuild does not support 'ResetState'. Use Build.")
	}

	h.runMode.Watching = config.Watching

	firstSite := h.Sites[0]

	h.nodeMap = make(map[string]Nodes)
	for _, s := range h.Sites {
		s.resetBuildState()
	}

	sourceChanged, err := firstSite.reBuild(events)

	if err != nil {
		return err
	}

	// Assign pages to sites per translation.
	h.setupTranslations(firstSite)

	if sourceChanged {
		for _, s := range h.Sites {
			if err := s.postProcess(); err != nil {
				return err
			}
		}
	}

	if err := h.preRender(); err != nil {
		return err
	}

	if !config.SkipRender {
		for _, s := range h.Sites {
			if err := s.render(); err != nil {
				return err
			}
			if config.PrintStats {
				s.Stats()
			}
		}

		if err := h.render(); err != nil {
			return err
		}
	}

	if config.PrintStats {
		jww.FEEDBACK.Printf("total in %v ms\n", int(1000*time.Since(t0).Seconds()))
	}

	return nil

}

// Analyze prints a build report to Stdout.
// Useful for debugging.
func (h *HugoSites) Analyze() error {
	if err := h.Build(BuildCfg{SkipRender: true}); err != nil {
		return err
	}
	s := h.Sites[0]
	return s.ShowPlan(os.Stdout)
}

// Render the cross-site artifacts.
func (h *HugoSites) render() error {

	if !h.multilingual.enabled() {
		return nil
	}

	// TODO(bep) DRY
	sitemapDefault := parseSitemap(viper.GetStringMap("Sitemap"))

	s := h.Sites[0]

	smLayouts := []string{"sitemapindex.xml", "_default/sitemapindex.xml", "_internal/_default/sitemapindex.xml"}

	if err := s.renderAndWriteXML("sitemapindex", sitemapDefault.Filename,
		h.toSiteInfos(), s.appendThemeTemplates(smLayouts)...); err != nil {
		return err
	}

	return nil
}

func (h *HugoSites) setupTranslations(master *Site) {

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

// preRender performs build tasks that need to be done as late as possible.
// Shortcode handling is the main task in here.
// TODO(bep) We need to look at the whole handler-chain construct with he below in mind.
func (h *HugoSites) preRender() error {

	for _, s := range h.Sites {
		if err := s.setCurrentLanguageConfig(); err != nil {
			return err
		}
		// Run "render prepare"
		if err := s.renderHomePage(true); err != nil {
			return err
		}
		if err := s.renderTaxonomiesLists(true); err != nil {
			return err
		}
		if err := s.renderListsOfTaxonomyTerms(true); err != nil {
			return err
		}
		if err := s.renderSectionLists(true); err != nil {
			return err
		}
	}

	for _, s := range h.Sites {
		if err := s.setCurrentLanguageConfig(); err != nil {
			return err
		}
		renderShortcodesForSite(s)
	}

	return nil
}

func renderShortcodesForSite(s *Site) {
	pageChan := make(chan *Page)
	wg := &sync.WaitGroup{}

	for i := 0; i < getGoMaxProcs()*4; i++ {
		wg.Add(1)
		go func(pages <-chan *Page, wg *sync.WaitGroup) {
			defer wg.Done()
			for p := range pages {
				if err := handleShortcodes(p, s.owner.tmpl); err != nil {
					jww.ERROR.Printf("Failed to handle shortcodes for page %s: %s", p.BaseFileName(), err)
				}

				if p.Markup == "markdown" {
					tmpContent, tmpTableOfContents := helpers.ExtractTOC(p.rawContent)
					p.TableOfContents = helpers.BytesToHTML(tmpTableOfContents)
					p.rawContent = tmpContent
				}

				if p.Markup != "html" {

					// Now we know enough to create a summary of the page and count some words
					summaryContent, err := p.setUserDefinedSummaryIfProvided()

					if err != nil {
						jww.ERROR.Printf("Failed to set use defined summary: %s", err)
					} else if summaryContent != nil {
						p.rawContent = summaryContent.content
					}

					p.Content = helpers.BytesToHTML(p.rawContent)
					p.rendered = true

					if summaryContent == nil {
						p.setAutoSummary()
					}
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

func handleShortcodes(p *Page, t tpl.Template) error {
	if len(p.contentShortCodes) > 0 {
		jww.DEBUG.Printf("Replace %d shortcodes in %q", len(p.contentShortCodes), p.BaseFileName())
		shortcodes, err := executeShortcodeFuncMap(p.contentShortCodes)

		if err != nil {
			return err
		}

		p.rawContent, err = replaceShortcodeTokens(p.rawContent, shortcodePlaceholderPrefix, shortcodes)

		if err != nil {
			jww.FATAL.Printf("Failed to replace short code tokens in %s:\n%s", p.BaseFileName(), err.Error())
		}
	}

	return nil
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
func newHugoSitesFromLanguages(languages helpers.Languages) (*HugoSites, error) {
	return newHugoSitesFromSourceAndLanguages(nil, languages)
}

func newHugoSitesDefaultLanguage() (*HugoSites, error) {
	return newHugoSitesFromSourceAndLanguages(nil, helpers.Languages{helpers.NewDefaultLanguage()})
}
