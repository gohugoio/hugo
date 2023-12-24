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
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

var pageIDCounter atomic.Uint64

func (h *HugoSites) newPage(m *pageMeta) (*pageState, error) {
	if m.pathInfo == nil {
		if m.f != nil {
			m.pathInfo = m.f.FileInfo().Meta().PathInfo
		}
		if m.pathInfo == nil {
			panic(fmt.Sprintf("missing pathInfo in %v", m))
		}
	}

	m.Staler = &resources.AtomicStaler{}

	ps, err := func() (*pageState, error) {
		if m.s == nil {
			// Identify the Site/language to associate this Page with.
			var lang string
			if m.f != nil {
				meta := m.f.FileInfo().Meta()
				lang = meta.Lang
				m.s = h.Sites[meta.LangIndex]
			} else {
				lang = m.pathInfo.Lang()
			}
			var found bool
			for _, ss := range h.Sites {
				if ss.Lang() == lang {
					m.s = ss
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("no site found for language %q", lang)
			}

		}

		// Identify Page Kind.
		if m.kind == "" {
			m.kind = kinds.KindSection
			if m.pathInfo.Base() == "/" {
				m.kind = kinds.KindHome
			} else if m.pathInfo.IsBranchBundle() {
				// A section, taxonomy or term.
				tc := m.s.pageMap.cfg.getTaxonomyConfig(m.Path())
				if !tc.IsZero() {
					// Either a taxonomy or a term.
					if tc.pluralTreeKey == m.Path() {
						m.kind = kinds.KindTaxonomy
					} else {
						m.kind = kinds.KindTerm
					}
				}
			} else if m.f != nil {
				m.kind = kinds.KindPage
			}
		}

		if m.kind == kinds.KindPage && !m.s.conf.IsKindEnabled(m.kind) {
			return nil, nil
		}

		pid := pageIDCounter.Add(1)

		// Parse page content.
		cachedContent, err := newCachedContent(m, pid)
		if err != nil {
			return nil, m.wrapError(err)
		}

		var dependencyManager identity.Manager = identity.NopManager

		if m.s.conf.Internal.Watch {
			dependencyManager = identity.NewManager(m.Path())
		}

		ps := &pageState{
			pid:                               pid,
			pageOutput:                        nopPageOutput,
			pageOutputTemplateVariationsState: &atomic.Uint32{},
			resourcesPublishInit:              &sync.Once{},
			Staler:                            m,
			dependencyManager:                 dependencyManager,
			pageCommon: &pageCommon{
				content:                   cachedContent,
				FileProvider:              m,
				AuthorProvider:            m,
				Scratcher:                 maps.NewScratcher(),
				store:                     maps.NewScratch(),
				Positioner:                page.NopPage,
				InSectionPositioner:       page.NopPage,
				ResourceNameTitleProvider: m,
				ResourceParamsProvider:    m,
				PageMetaProvider:          m,
				RelatedKeywordsProvider:   m,
				OutputFormatsProvider:     page.NopPage,
				ResourceTypeProvider:      pageTypesProvider,
				MediaTypeProvider:         pageTypesProvider,
				RefProvider:               page.NopPage,
				ShortcodeInfoProvider:     page.NopPage,
				LanguageProvider:          m.s,

				InternalDependencies: m.s,
				init:                 lazy.New(),
				m:                    m,
				s:                    m.s,
				sWrapped:             page.WrapSite(m.s),
			},
		}

		if m.f != nil {
			gi, err := m.s.h.gitInfoForPage(ps)
			if err != nil {
				return nil, fmt.Errorf("failed to load Git data: %w", err)
			}
			ps.gitInfo = gi
			owners, err := m.s.h.codeownersForPage(ps)
			if err != nil {
				return nil, fmt.Errorf("failed to load CODEOWNERS: %w", err)
			}
			ps.codeowners = owners
		}

		ps.pageMenus = &pageMenus{p: ps}
		ps.PageMenusProvider = ps.pageMenus
		ps.GetPageProvider = pageSiteAdapter{s: m.s, p: ps}
		ps.GitInfoProvider = ps
		ps.TranslationsProvider = ps
		ps.ResourceDataProvider = &pageData{pageState: ps}
		ps.RawContentProvider = ps
		ps.ChildCareProvider = ps
		ps.TreeProvider = pageTree{p: ps}
		ps.Eqer = ps
		ps.TranslationKeyProvider = ps
		ps.ShortcodeInfoProvider = ps
		ps.AlternativeOutputFormatsProvider = ps

		if err := ps.setMetaPre(); err != nil {
			return nil, ps.wrapError(err)
		}

		if err := ps.initLazyProviders(); err != nil {
			return nil, ps.wrapError(err)
		}
		return ps, nil
	}()
	// Make sure to evict any cached and now stale data.
	if err != nil {
		m.MarkStale()
	}

	return ps, err
}
