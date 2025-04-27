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
	"iter"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
)

var pageIDCounter atomic.Uint64

func (h *HugoSites) newPages(m *pageMeta) (iter.Seq2[int, *pageState], *paths.Path, error) {
	p, pth, sites, err := h.doNewPage(m)
	if err != nil {
		// Make sure that any partially created page part is marked as stale.
		m.MarkStale()
	}

	iter := func(yield func(i int, p *pageState) bool) {
		if !yield(0, p) {
			return
		}
		for i := 1; i < len(sites); i++ {
			s := sites[i]
			p, err = p.cloneForSite(s)
			if err != nil {
				return
			}
			if !yield(i, p) {
				return
			}
		}
	}

	return iter, pth, err
}

func (h *HugoSites) newPage(m *pageMeta) (*pageState, *paths.Path, error) {
	p, pth, _, err := h.doNewPage(m)
	if err != nil {
		// Make sure that any partially created page part is marked as stale.
		m.MarkStale()
	}
	return p, pth, err
}

func (h *HugoSites) doNewPage(m *pageMeta) (*pageState, *paths.Path, []*Site, error) {
	m.Staler = &resources.AtomicStaler{}
	if m.pageMetaParams == nil {
		m.pageMetaParams = &pageMetaParams{
			pageConfig: &pagemeta.PageConfig{},
		}
	}
	if m.pageConfig.Params == nil {
		m.pageConfig.Params = maps.Params{}
	}

	pid := pageIDCounter.Add(1)
	pi, err := m.parseFrontMatter(h, pid)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := m.setMetaPre(pi, h.Log, h.Conf); err != nil {
		return nil, nil, nil, m.wrapError(err, h.BaseFs.SourceFs)
	}
	pcfg := m.pageConfig
	if pcfg.Lang != "" {
		if h.Conf.IsLangDisabled(pcfg.Lang) {
			return nil, nil, nil, nil
		}
	}

	if pcfg.Path != "" {
		s := m.pageConfig.Path
		// Paths from content adapters should never have any extension.
		if pcfg.IsFromContentAdapter || !paths.HasExt(s) {
			var (
				isBranch    bool
				isBranchSet bool
				ext         string = m.pageConfig.ContentMediaType.FirstSuffix.Suffix
			)
			if pcfg.Kind != "" {
				isBranch = kinds.IsBranch(pcfg.Kind)
				isBranchSet = true
			}

			if !pcfg.IsFromContentAdapter {
				if m.pathInfo != nil {
					if !isBranchSet {
						isBranch = m.pathInfo.IsBranchBundle()
					}
					if m.pathInfo.Ext() != "" {
						ext = m.pathInfo.Ext()
					}
				} else if m.f != nil {
					pi := m.f.FileInfo().Meta().PathInfo
					if !isBranchSet {
						isBranch = pi.IsBranchBundle()
					}
					if pi.Ext() != "" {
						ext = pi.Ext()
					}
				}
			}

			if isBranch {
				s += "/_index." + ext
			} else {
				s += "/index." + ext
			}

		}
		m.pathInfo = h.Conf.PathParser().Parse(files.ComponentFolderContent, s)
	} else if m.pathInfo == nil {
		if m.f != nil {
			m.pathInfo = m.f.FileInfo().Meta().PathInfo
		}

		if m.pathInfo == nil {
			panic(fmt.Sprintf("missing pathInfo in %v", m))
		}
	}

	ps, sites, err := func() (*pageState, []*Site, error) {
		var sites []*Site
		if m.s == nil {
			// Identify the Site/language to associate this Page with.
			// TODO1 LanguagesCompiledMap
			/*var lang string
			if pcfg.Lang != "" {
				lang = pcfg.Lang
			} else if m.f != nil {
				meta := m.f.FileInfo().Meta()
				lang = meta.Lang
			} else {
				lang = m.pathInfo.Lang()
			}*/

			// TODO1 avoid allocating a new slice here.
			sites = h.resolveSites(pcfg.LanguagesCompiledMap, pcfg.VersionsCompiledMap, pcfg.RolesCompiledMap)
			if len(sites) == 0 {
				return nil, nil, fmt.Errorf("no site found for languages %v, versions %v and roles %v", pcfg.LanguagesCompiledMap, pcfg.VersionsCompiledMap, pcfg.RolesCompiledMap)
			}
			m.s = sites[0]

		}

		var tc viewName
		// Identify Page Kind.
		if m.pageConfig.Kind == "" {
			m.pageConfig.Kind = kinds.KindSection
			if m.pathInfo.Base() == "/" {
				m.pageConfig.Kind = kinds.KindHome
			} else if m.pathInfo.IsBranchBundle() {
				// A section, taxonomy or term.
				tc = m.s.pageMap.cfg.getTaxonomyConfig(m.Path())
				if !tc.IsZero() {
					// Either a taxonomy or a term.
					if tc.pluralTreeKey == m.Path() {
						m.pageConfig.Kind = kinds.KindTaxonomy
					} else {
						m.pageConfig.Kind = kinds.KindTerm
					}
				}
			} else if m.f != nil {
				m.pageConfig.Kind = kinds.KindPage
			}
		}

		if m.pageConfig.Kind == kinds.KindTerm || m.pageConfig.Kind == kinds.KindTaxonomy {
			if tc.IsZero() {
				tc = m.s.pageMap.cfg.getTaxonomyConfig(m.Path())
			}
			if tc.IsZero() {
				return nil, nil, fmt.Errorf("no taxonomy configuration found for %q", m.Path())
			}
			m.singular = tc.singular
			if m.pageConfig.Kind == kinds.KindTerm {
				m.term = paths.TrimLeading(strings.TrimPrefix(m.pathInfo.Unnormalized().Base(), tc.pluralTreeKey))
			}
		}

		if m.pageConfig.Kind == kinds.KindPage && !m.s.conf.IsKindEnabled(m.pageConfig.Kind) {
			return nil, nil, nil
		}

		// Parse the rest of the page content.
		m.content, err = m.newCachedContent(h, pi)
		if err != nil {
			return nil, nil, m.wrapError(err, h.SourceFs)
		}

		ps, err := h.doNewPageFromMeta(pid, m)
		if err != nil {
			return nil, nil, m.wrapError(err, h.SourceFs)
		}

		if m.f != nil {
			gi, err := m.s.h.gitInfoForPage(ps)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to load Git data: %w", err)
			}
			ps.gitInfo = gi
			owners, err := m.s.h.codeownersForPage(ps)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to load CODEOWNERS: %w", err)
			}
			ps.codeowners = owners
		}

		if err := ps.initLazyProviders(); err != nil {
			return nil, nil, ps.wrapError(err)
		}
		return ps, sites, nil
	}()

	if ps == nil {
		return nil, nil, nil, err
	}

	return ps, ps.PathInfo(), sites, err
}

func (h *HugoSites) doNewPageFromMeta(pid uint64, m *pageMeta) (*pageState, error) {
	ps := &pageState{
		pid:                               pid,
		pageOutput:                        nopPageOutput,
		pageOutputTemplateVariationsState: &atomic.Uint32{},
		resourcesPublishInit:              &sync.Once{},
		Staler:                            m,
		dependencyManager:                 m.s.Conf.NewIdentityManager(m.Path()),
		pageCommon: &pageCommon{
			store:                      maps.NewScratch(),
			Positioner:                 page.NopPage,
			InSectionPositioner:        page.NopPage,
			ResourceNameTitleProvider:  m,
			ResourceParamsProvider:     m,
			PageMetaProvider:           m,
			PageMetaInternalProvider:   m,
			FileProvider:               m,
			OutputFormatsProvider:      page.NopPage,
			ResourceTypeProvider:       pageTypesProvider,
			MediaTypeProvider:          pageTypesProvider,
			RefProvider:                page.NopPage,
			ShortcodeInfoProvider:      page.NopPage,
			LanguageProvider:           m.s,
			RelatedDocsHandlerProvider: m.s,
			init:                       lazy.New(),
			m:                          m,
			s:                          m.s,
			sWrapped:                   page.WrapSite(m.s), // TODO1 need this?
			pageContentConverter:       &pageContentConverter{},
		},
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

	return ps, nil
}
