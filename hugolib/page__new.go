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

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/lazy"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
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
			ResourceTypesProvider:   pageTypesProvider,
			RefProvider:             page.NopPage,
			ShortcodeInfoProvider:   page.NopPage,
			LanguageProvider:        s,
			pagePages:               &pagePages{},

			InternalDependencies: s,
			init:                 lazy.New(),
			m:                    metaProvider,
			s:                    s},
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

func newPageFromMeta(meta map[string]interface{}, metaProvider *pageMeta) (*pageState, error) {
	if metaProvider.f == nil {
		metaProvider.f = page.NewZeroFile(metaProvider.s.DistinctWarningLog)
	}

	ps, err := newPageBase(metaProvider)
	if err != nil {
		return nil, err
	}

	initMeta := func(bucket *pagesMapBucket) error {
		if meta != nil || bucket != nil {
			if err := metaProvider.setMetadata(bucket, ps, meta); err != nil {
				return ps.wrapError(err)
			}
		}

		if err := metaProvider.applyDefaultValues(); err != nil {
			return err
		}

		return nil
	}

	if metaProvider.standalone {
		initMeta(nil)
	} else {
		// Because of possible cascade keywords, we need to delay this
		// until we have the complete page graph.
		ps.metaInitFn = initMeta
	}

	ps.init.Add(func() (interface{}, error) {
		pp, err := newPagePaths(metaProvider.s, ps, metaProvider)
		if err != nil {
			return nil, err
		}

		makeOut := func(f output.Format, render bool) *pageOutput {
			return newPageOutput(nil, ps, pp, f, render)
		}

		if ps.m.standalone {
			ps.pageOutput = makeOut(ps.m.outputFormats()[0], true)
		} else {
			ps.pageOutputs = make([]*pageOutput, len(ps.s.h.renderFormats))
			created := make(map[string]*pageOutput)
			outputFormatsForPage := ps.m.outputFormats()
			for i, f := range ps.s.h.renderFormats {
				po, found := created[f.Name]
				if !found {
					_, shouldRender := outputFormatsForPage.GetByName(f.Name)
					po = makeOut(f, shouldRender)
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
	p, err := newPageFromMeta(nil, m)

	if err != nil {
		return nil, err
	}

	if err := p.initPage(); err != nil {
		return nil, err
	}

	return p, nil

}

func newPageWithContent(f *fileInfo, s *Site, bundled bool, content resource.OpenReadSeekCloser) (*pageState, error) {
	sections := s.sectionsFromFile(f)
	kind := s.kindFromFileInfoOrSections(f, sections)
	if kind == page.KindTaxonomy {
		s.PathSpec.MakePathsSanitized(sections)
	}

	metaProvider := &pageMeta{kind: kind, sections: sections, bundled: bundled, s: s, f: f}

	ps, err := newPageBase(metaProvider)
	if err != nil {
		return nil, err
	}

	gi, err := s.h.gitInfoForPage(ps)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load Git data")
	}
	ps.gitInfo = gi

	r, err := content()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	parseResult, err := pageparser.Parse(
		r,
		pageparser.Config{EnableEmoji: s.siteCfg.enableEmoji},
	)
	if err != nil {
		return nil, err
	}

	ps.pageContent = pageContent{
		source: rawPageContent{
			parsed:         parseResult,
			posMainContent: -1,
			posSummaryEnd:  -1,
			posBodyStart:   -1,
		},
	}

	ps.shortcodeState = newShortcodeHandler(ps, ps.s, nil)

	ps.metaInitFn = func(bucket *pagesMapBucket) error {
		if err := ps.mapContent(bucket, metaProvider); err != nil {
			return ps.wrapError(err)
		}

		if err := metaProvider.applyDefaultValues(); err != nil {
			return err
		}

		return nil
	}

	ps.init.Add(func() (interface{}, error) {
		reuseContent := ps.renderable && !ps.shortcodeState.hasShortcodes()

		// Creates what's needed for each output format.
		contentPerOutput := newPageContentOutput(ps)

		pp, err := newPagePaths(s, ps, metaProvider)
		if err != nil {
			return nil, err
		}

		// Prepare output formats for all sites.
		ps.pageOutputs = make([]*pageOutput, len(ps.s.h.renderFormats))
		created := make(map[string]*pageOutput)
		outputFormatsForPage := ps.m.outputFormats()

		for i, f := range ps.s.h.renderFormats {
			if po, found := created[f.Name]; found {
				ps.pageOutputs[i] = po
				continue
			}

			_, render := outputFormatsForPage.GetByName(f.Name)
			var contentProvider *pageContentOutput
			if reuseContent && i > 0 {
				contentProvider = ps.pageOutputs[0].cp
			} else {
				var err error
				contentProvider, err = contentPerOutput(f)
				if err != nil {
					return nil, err
				}
			}

			po := newPageOutput(contentProvider, ps, pp, f, render)
			ps.pageOutputs[i] = po
			created[f.Name] = po
		}

		if err := ps.initCommonProviders(pp); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return ps, nil
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
