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
	"html/template"
	"strings"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/output"

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
			pagePages:               &pagePages{},

			InternalDependencies: s,
			init:                 lazy.New(),
			m:                    metaProvider,
			s:                    s,
		},
	}

	siteAdapter := pageSiteAdapter{s: s, p: ps}

	deprecatedWarningPage := struct {
		source.FileWithoutOverlap
		page.DeprecatedWarningPageMethods1
	}{
		FileWithoutOverlap:            metaProvider.File(),
		DeprecatedWarningPageMethods1: &pageDeprecatedWarning{p: ps},
	}

	ps.DeprecatedWarningPageMethods = page.NewDeprecatedWarningPage(deprecatedWarningPage)
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

func newPageBucket(p *pageState) *pagesMapBucket {
	return &pagesMapBucket{owner: p, pagesMapBucketPages: &pagesMapBucketPages{}}
}

func newPageFromMeta(
	n *contentNode,
	parentBucket *pagesMapBucket,
	meta map[string]interface{},
	metaProvider *pageMeta) (*pageState, error) {

	if metaProvider.f == nil {
		metaProvider.f = page.NewZeroFile(metaProvider.s.DistinctWarningLog)
	}

	ps, err := newPageBase(metaProvider)
	if err != nil {
		return nil, err
	}

	bucket := parentBucket

	if ps.IsNode() {
		ps.bucket = newPageBucket(ps)
	}

	if meta != nil || parentBucket != nil {
		if err := metaProvider.setMetadata(bucket, ps, meta); err != nil {
			return nil, ps.wrapError(err)
		}
	}

	if err := metaProvider.applyDefaultValues(n); err != nil {
		return nil, err
	}

	ps.init.Add(func() (interface{}, error) {
		pp, err := newPagePaths(metaProvider.s, ps, metaProvider)
		if err != nil {
			return nil, err
		}

		makeOut := func(f output.Format, render bool) *pageOutput {
			return newPageOutput(ps, pp, f, render)
		}

		shouldRenderPage := !ps.m.noRender()

		if ps.m.standalone {
			ps.pageOutput = makeOut(ps.m.outputFormats()[0], shouldRenderPage)
		} else {
			outputFormatsForPage := ps.m.outputFormats()

			// Prepare output formats for all sites.
			// We do this even if this page does not get rendered on
			// its own. It may be referenced via .Site.GetPage and
			// it will then need an output format.
			ps.pageOutputs = make([]*pageOutput, len(ps.s.h.renderFormats))
			created := make(map[string]*pageOutput)
			for i, f := range ps.s.h.renderFormats {
				po, found := created[f.Name]
				if !found {
					render := shouldRenderPage
					if render {
						_, render = outputFormatsForPage.GetByName(f.Name)
					}
					po = makeOut(f, render)
					created[f.Name] = po
				}
				ps.pageOutputs[i] = po
			}
		}

		if err := ps.initCommonProviders(pp); err != nil {
			return nil, err
		}

		return nil, nil

	})

	return ps, err

}

// Used by the legacy 404, sitemap and robots.txt rendering
func newPageStandalone(m *pageMeta, f output.Format) (*pageState, error) {
	m.configuredOutputFormats = output.Formats{f}
	m.standalone = true
	p, err := newPageFromMeta(nil, nil, nil, m)

	if err != nil {
		return nil, err
	}

	if err := p.initPage(); err != nil {
		return nil, err
	}

	return p, nil

}

type pageDeprecatedWarning struct {
	p *pageState
}

func (p *pageDeprecatedWarning) IsDraft() bool          { return p.p.m.draft }
func (p *pageDeprecatedWarning) Hugo() hugo.Info        { return p.p.s.Info.Hugo() }
func (p *pageDeprecatedWarning) LanguagePrefix() string { return p.p.s.Info.LanguagePrefix }
func (p *pageDeprecatedWarning) GetParam(key string) interface{} {
	return p.p.m.params[strings.ToLower(key)]
}
func (p *pageDeprecatedWarning) RSSLink() template.URL {
	f := p.p.OutputFormats().Get("RSS")
	if f == nil {
		return ""
	}
	return template.URL(f.Permalink())
}
func (p *pageDeprecatedWarning) URL() string {
	if p.p.IsPage() && p.p.m.urlPaths.URL != "" {
		// This is the url set in front matter
		return p.p.m.urlPaths.URL
	}
	// Fall back to the relative permalink.
	return p.p.RelPermalink()

}
