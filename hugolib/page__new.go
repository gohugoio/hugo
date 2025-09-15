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
	"strings"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/constants"
	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/page"
)

var (
	pageIDCounter       atomic.Uint64
	pageSourceIDCounter atomic.Uint64
)

func (s *Site) newPageFromPageMeta(m *pageMeta) (*pageState, error) {
	if err := m.initLate(s); err != nil {
		return nil, m.wrapError(err, s.SourceFs)
	}
	pid := pageIDCounter.Add(1)
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
			sWrapped:                   page.WrapSite(s),
			pageContentConverter:       &pageContentConverter{},
		},
	}

	if ps.IsHome() && ps.PathInfo().IsLeafBundle() {
		msg := "Using %s in your content's root directory is usually incorrect for your home page. "
		msg += "You should use %s instead. If you don't rename this file, your home page will be "
		msg += "treated as a leaf bundle, meaning it won't be able to have any child pages or sections."
		ps.s.Log.Warnidf(constants.WarnHomePageIsLeafBundle, msg, ps.PathInfo().PathNoLeadingSlash(), strings.ReplaceAll(ps.PathInfo().PathNoLeadingSlash(), "index", "_index"))
	}

	if m.f != nil {
		gi, err := s.h.gitInfoForPage(ps)
		if err != nil {
			return nil, fmt.Errorf("failed to load Git data: %w", err)
		}
		ps.gitInfo = gi
		owners, err := s.h.codeownersForPage(ps)
		if err != nil {
			return nil, fmt.Errorf("failed to load CODEOWNERS: %w", err)
		}
		ps.codeowners = owners
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

	if err := ps.initLazyProviders(); err != nil {
		return nil, ps.wrapError(err)
	}

	return ps, nil
}
