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
	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/page"
)

func newPageBase(metaProvider *pageMeta) (*pageState, error) {
	if metaProvider.s == nil {
		panic("must provide a Site")
	}

	s := metaProvider.s

	ps := &pageState{
		pageOutput: nopPageOutput,
		pageCommon: &pageCommon{
			FileProvider:            metaProvider,
			AuthorProvider:          metaProvider,
			Scratcher:               maps.NewScratcher(),
			Positioner:              page.NopPage,
			InSectionPositioner:     page.NopPage,
			ResourceMetaProvider:    metaProvider,
			ResourceParamsProvider:  metaProvider,
			PageMetaProvider:        metaProvider,
			RelatedKeywordsProvider: metaProvider,
			OutputFormatsProvider:   page.NopPage,
			ResourceTypeProvider:    pageTypesProvider,
			MediaTypeProvider:       pageTypesProvider,
			RefProvider:             page.NopPage,
			ShortcodeInfoProvider:   page.NopPage,
			LanguageProvider:        s,

			InternalDependencies: s,
			init:                 lazy.New(),
			m:                    metaProvider,
			s:                    s,
		},
	}

	siteAdapter := pageSiteAdapter{s: s, p: ps}

	ps.pageMenus = &pageMenus{p: ps}
	ps.PageMenusProvider = ps.pageMenus
	ps.GetPageProvider = siteAdapter
	ps.GitInfoProvider = ps
	ps.TranslationsProvider = ps
	ps.ResourceDataProvider = &pageData{pageState: ps}
	ps.RawContentProvider = ps
	ps.ChildCareProvider = ps
	ps.TreeProvider = pageTree{p: ps}
	ps.Eqer = ps
	ps.TranslationKeyProvider = ps
	ps.ShortcodeInfoProvider = ps
	ps.PageRenderProvider = ps
	ps.AlternativeOutputFormatsProvider = ps

	return ps, nil
}

func newPageBucket(parent *pagesMapBucket, self *pageState) *pagesMapBucket {
	return &pagesMapBucket{parent: parent, self: self, pagesMapBucketPages: &pagesMapBucketPages{}}
}
