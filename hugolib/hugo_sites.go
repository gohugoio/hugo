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

	Multilingual *Multilingual
}

func NewHugoSites(sites ...*Site) (*HugoSites, error) {
	languages := make(Languages, len(sites))
	for i, s := range sites {
		if s.Language == nil {
			return nil, errors.New("Missing language for site")
		}
		languages[i] = s.Language
	}
	defaultLang := viper.GetString("DefaultContentLanguage")
	if defaultLang == "" {
		defaultLang = "en"
	}
	langConfig := &Multilingual{Languages: languages, DefaultLang: NewLanguage(defaultLang)}

	return &HugoSites{Multilingual: langConfig, Sites: sites}, nil
}

// Reset resets the sites, making it ready for a full rebuild.
// TODO(bep) multilingo
func (h HugoSites) Reset() {
	for i, s := range h.Sites {
		h.Sites[i] = s.Reset()
	}
}

type BuildCfg struct {
	// Whether we are in watch (server) mode
	Watching bool
	// Print build stats at the end of a build
	PrintStats bool
	// Skip rendering. Useful for testing.
	skipRender bool
	// Use this to add templates to use for rendering.
	// Useful for testing.
	withTemplate func(templ tpl.Template) error
}

// Build builds all sites.
func (h HugoSites) Build(config BuildCfg) error {

	if h.Sites == nil || len(h.Sites) == 0 {
		return errors.New("No site(s) to build")
	}

	t0 := time.Now()

	// We should probably refactor the Site and pull up most of the logic from there to here,
	// but that seems like a daunting task.
	// So for now, if there are more than one site (language),
	// we pre-process the first one, then configure all the sites based on that.
	firstSite := h.Sites[0]

	for _, s := range h.Sites {
		// TODO(bep) ml
		s.Multilingual = h.Multilingual
		s.RunMode.Watching = config.Watching
	}

	if err := firstSite.PreProcess(config); err != nil {
		return err
	}

	h.setupTranslations(firstSite)

	if len(h.Sites) > 1 {
		// Initialize the rest
		for _, site := range h.Sites[1:] {
			site.Tmpl = firstSite.Tmpl
			site.initializeSiteInfo()
		}
	}

	for _, s := range h.Sites {
		if err := s.PostProcess(); err != nil {
			return err
		}
	}

	if err := h.preRender(); err != nil {
		return err
	}

	for _, s := range h.Sites {

		if !config.skipRender {
			if err := s.Render(); err != nil {
				return err
			}

			if config.PrintStats {
				s.Stats()
			}
		}
		// TODO(bep) ml lang in site.Info?
	}

	if config.PrintStats {
		jww.FEEDBACK.Printf("total in %v ms\n", int(1000*time.Since(t0).Seconds()))
	}

	return nil

}

// Rebuild rebuilds all sites.
func (h HugoSites) Rebuild(config BuildCfg, events ...fsnotify.Event) error {
	t0 := time.Now()

	firstSite := h.Sites[0]

	for _, s := range h.Sites {
		s.resetBuildState()
	}

	sourceChanged, err := firstSite.ReBuild(events)

	if err != nil {
		return err
	}

	// Assign pages to sites per translation.
	h.setupTranslations(firstSite)

	if sourceChanged {
		for _, s := range h.Sites {
			if err := s.PostProcess(); err != nil {
				return err
			}
		}
	}

	if err := h.preRender(); err != nil {
		return err
	}

	if !config.skipRender {
		for _, s := range h.Sites {
			if err := s.Render(); err != nil {
				return err
			}
			if config.PrintStats {
				s.Stats()
			}
		}
	}

	if config.PrintStats {
		jww.FEEDBACK.Printf("total in %v ms\n", int(1000*time.Since(t0).Seconds()))
	}

	return nil

}

func (s *HugoSites) setupTranslations(master *Site) {

	for _, p := range master.rawAllPages {
		if p.Lang() == "" {
			panic("Page language missing: " + p.Title)
		}

		shouldBuild := p.shouldBuild()

		for i, site := range s.Sites {
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

		for i := 1; i < len(s.Sites); i++ {
			s.Sites[i].AllPages = s.Sites[0].AllPages
		}
	}

	if len(s.Sites) > 1 {
		pages := s.Sites[0].AllPages
		allTranslations := pagesToTranslationsMap(s.Multilingual, pages)
		assignTranslationsToPages(allTranslations, pages)
	}
}

// preRender performs build tasks that needs to be done as late as possible.
// Shortcode handling is the main task in here.
// TODO(bep) We need to look at the whole handler-chain construct witht he below in mind.
func (h *HugoSites) preRender() error {
	pageChan := make(chan *Page)

	wg := &sync.WaitGroup{}

	// We want all the pages, so just pick one.
	s := h.Sites[0]

	for i := 0; i < getGoMaxProcs()*4; i++ {
		wg.Add(1)
		go func(pages <-chan *Page, wg *sync.WaitGroup) {
			defer wg.Done()
			for p := range pages {
				if err := handleShortcodes(p, s.Tmpl); err != nil {
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

	for _, p := range s.AllPages {
		pageChan <- p
	}

	close(pageChan)

	wg.Wait()

	return nil
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
	sites, err := NewHugoSites(s)
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

	config := BuildCfg{skipRender: !render, withTemplate: addTemplates}
	return sites.Build(config)
}

// Convenience func used in tests.
func newHugoSitesFromSourceAndLanguages(input []source.ByteSource, languages Languages) (*HugoSites, error) {
	if len(languages) == 0 {
		panic("Must provide at least one language")
	}
	first := &Site{
		Source:   &source.InMemorySource{ByteSource: input},
		Language: languages[0],
	}
	if len(languages) == 1 {
		return NewHugoSites(first)
	}

	sites := make([]*Site, len(languages))
	sites[0] = first
	for i := 1; i < len(languages); i++ {
		sites[i] = &Site{Language: languages[i]}
	}

	return NewHugoSites(sites...)

}

// Convenience func used in tests.
func newHugoSitesFromLanguages(languages Languages) (*HugoSites, error) {
	return newHugoSitesFromSourceAndLanguages(nil, languages)
}
