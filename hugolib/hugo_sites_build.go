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
	"bytes"
	"fmt"

	"errors"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/helpers"
)

// Build builds all sites. If filesystem events are provided,
// this is considered to be a potential partial rebuild.
func (h *HugoSites) Build(config BuildCfg, events ...fsnotify.Event) error {

	if h.Metrics != nil {
		h.Metrics.Reset()
	}

	//t0 := time.Now()

	// Need a pointer as this may be modified.
	conf := &config

	if conf.whatChanged == nil {
		// Assume everything has changed
		conf.whatChanged = &whatChanged{source: true, other: true}
	}

	for _, s := range h.Sites {
		s.Deps.BuildStartListeners.Notify()
	}

	if len(events) > 0 {
		// Rebuild
		if err := h.initRebuild(conf); err != nil {
			return err
		}
	} else {
		if err := h.init(conf); err != nil {
			return err
		}
	}

	if err := h.process(conf, events...); err != nil {
		return err
	}

	if err := h.assemble(conf); err != nil {
		return err
	}

	if err := h.render(conf); err != nil {
		return err
	}

	if h.Metrics != nil {
		var b bytes.Buffer
		h.Metrics.WriteMetrics(&b)

		h.Log.FEEDBACK.Printf("\nTemplate Metrics:\n\n")
		h.Log.FEEDBACK.Print(b.String())
		h.Log.FEEDBACK.Println()
	}

	errorCount := h.Log.LogCountForLevel(jww.LevelError)
	if errorCount > 0 {
		return fmt.Errorf("logged %d error(s)", errorCount)
	}

	return nil

}

// Build lifecycle methods below.
// The order listed matches the order of execution.

func (h *HugoSites) init(config *BuildCfg) error {

	for _, s := range h.Sites {
		if s.PageCollections == nil {
			s.PageCollections = newPageCollections()
		}
	}

	if config.ResetState {
		h.reset()
	}

	if config.CreateSitesFromConfig {
		if err := h.createSitesFromConfig(); err != nil {
			return err
		}
	}

	return nil
}

func (h *HugoSites) initRebuild(config *BuildCfg) error {
	if config.CreateSitesFromConfig {
		return errors.New("Rebuild does not support 'CreateSitesFromConfig'.")
	}

	if config.ResetState {
		return errors.New("Rebuild does not support 'ResetState'.")
	}

	if !h.running {
		return errors.New("Rebuild called when not in watch mode")
	}

	if config.whatChanged.source {
		// This is for the non-renderable content pages (rarely used, I guess).
		// We could maybe detect if this is really needed, but it should be
		// pretty fast.
		h.TemplateHandler().RebuildClone()
	}

	for _, s := range h.Sites {
		s.resetBuildState()
	}

	h.resetLogs()
	helpers.InitLoggers()

	return nil
}

func (h *HugoSites) process(config *BuildCfg, events ...fsnotify.Event) error {
	// We should probably refactor the Site and pull up most of the logic from there to here,
	// but that seems like a daunting task.
	// So for now, if there are more than one site (language),
	// we pre-process the first one, then configure all the sites based on that.

	firstSite := h.Sites[0]

	if len(events) > 0 {
		// This is a rebuild
		changed, err := firstSite.processPartial(events)
		config.whatChanged = &changed
		return err
	}

	return firstSite.process(*config)

}

func (h *HugoSites) assemble(config *BuildCfg) error {
	if config.whatChanged.source {
		for _, s := range h.Sites {
			s.createTaxonomiesEntries()
		}
	}

	// TODO(bep) we could probably wait and do this in one go later
	h.setupTranslations()

	if len(h.Sites) > 1 {
		// The first is initialized during process; initialize the rest
		for _, site := range h.Sites[1:] {
			if err := site.initializeSiteInfo(); err != nil {
				return err
			}
		}
	}

	if config.whatChanged.source {
		for _, s := range h.Sites {
			if err := s.buildSiteMeta(); err != nil {
				return err
			}
		}
	}

	if err := h.createMissingPages(); err != nil {
		return err
	}

	for _, s := range h.Sites {
		for _, pages := range []Pages{s.Pages, s.headlessPages} {
			for _, p := range pages {
				// May have been set in front matter
				if len(p.outputFormats) == 0 {
					p.outputFormats = s.outputFormats[p.Kind]
				}

				if p.headless {
					// headless = 1 output format only
					p.outputFormats = p.outputFormats[:1]
				}
				for _, r := range p.Resources.ByType(pageResourceType) {
					r.(*Page).outputFormats = p.outputFormats
				}

				if err := p.initPaths(); err != nil {
					return err
				}

			}
		}
		s.assembleMenus()
		s.refreshPageCaches()
		s.setupSitePages()
	}

	if err := h.assignMissingTranslations(); err != nil {
		return err
	}

	return nil

}

func (h *HugoSites) render(config *BuildCfg) error {
	for _, s := range h.Sites {
		s.initRenderFormats()
	}

	for _, s := range h.Sites {
		for i, rf := range s.renderFormats {
			for _, s2 := range h.Sites {
				// We render site by site, but since the content is lazily rendered
				// and a site can "borrow" content from other sites, every site
				// needs this set.
				s2.rc = &siteRenderingContext{Format: rf}

				isRenderingSite := s == s2

				if err := s2.preparePagesForRender(isRenderingSite && i == 0); err != nil {
					return err
				}

			}

			if !config.SkipRender {
				if err := s.render(config, i); err != nil {
					return err
				}
			}
		}
	}

	if !config.SkipRender {
		if err := h.renderCrossSitesArtifacts(); err != nil {
			return err
		}
	}

	return nil
}
