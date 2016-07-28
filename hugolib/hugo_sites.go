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
	"time"

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

		if !config.skipRender {
			if err := s.Render(); err != nil {
				return err
			}

		}

		if config.PrintStats {
			s.Stats()
		}

		// TODO(bep) ml lang in site.Info?
		// TODO(bep) ml Page sorting?
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

	for _, s := range h.Sites {

		if sourceChanged {
			if err := s.PostProcess(); err != nil {
				return err
			}
		}

		if !config.skipRender {
			if err := s.Render(); err != nil {
				return err
			}
		}

		if config.PrintStats {
			s.Stats()
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
