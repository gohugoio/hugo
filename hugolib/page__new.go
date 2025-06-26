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
	"iter"
	"strings"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/constants"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/page"
)

var (
	pageIDCounter       atomic.Uint64
	pageSourceIDCounter atomic.Uint64
)

func (h *HugoSites) newPages(m *pageMeta) (iter.Seq2[int, *pageState], *paths.Path, error) {
	panic("TODO1 remove me.")
}

func (h *HugoSites) newPage(m *pageMeta) (*pageState, *paths.Path, error) {
	p, pth, _, err := h.doNewPage(m)
	if err != nil {
		// Make sure that any partially created page part is marked as stale.
		m.MarkStale()
	}

	if p != nil && pth != nil && p.IsHome() && pth.IsLeafBundle() {
		msg := "Using %s in your content's root directory is usually incorrect for your home page. "
		msg += "You should use %s instead. If you don't rename this file, your home page will be "
		msg += "treated as a leaf bundle, meaning it won't be able to have any child pages or sections."
		h.Log.Warnidf(constants.WarnHomePageIsLeafBundle, msg, pth.PathNoLeadingSlash(), strings.ReplaceAll(pth.PathNoLeadingSlash(), "index", "_index"))
	}

	return p, pth, err
}

func (h *HugoSites) doNewPage(m *pageMeta) (*pageState, *paths.Path, []*Site, error) {
	panic("TODO1 remove me.")
}

func (h *HugoSites) doNewPageFromMeta(pid uint64, m *pageMeta) (*pageState, error) {
	panic("TODO1 remove me?")
}

// TODO1 move/rename.
func (s *Site) doNewPageFromMeta(pid uint64, m *pageMeta) (*pageState, error) {
	if err := m.initLate(s, pid); err != nil {
		return nil, m.wrapError(err, s.SourceFs)
	}
	// Parse the rest of the page content.
	var err error
	m.content, err = m.newCachedContent(s)
	if err != nil {
		return nil, m.wrapError(err, s.SourceFs)
	}
	ps := &pageState{
		pid:                               pid,
		s:                                 s,
		pageOutput:                        nopPageOutput,
		pageOutputTemplateVariationsState: &atomic.Uint32{},
		Staler:                            m,
		dependencyManager:                 s.Conf.NewIdentityManager(m.Path()),
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
			LanguageProvider:           s,
			RelatedDocsHandlerProvider: s,
			init:                       lazy.New(),
			m:                          m,
			s:                          s,
			sWrapped:                   s, // page.WrapSite(s), // TODO1 need this?
			pageContentConverter:       &pageContentConverter{},
		},
	}

	ps.pageMenus = &pageMenus{p: ps}
	ps.PageMenusProvider = ps.pageMenus
	ps.GetPageProvider = pageSiteAdapter{s: s, p: ps}
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
