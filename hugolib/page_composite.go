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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/pkg/errors"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/compare"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	cjk = regexp.MustCompile(`\p{Han}|\p{Hangul}|\p{Hiragana}|\p{Katakana}`)

	// This is all the kinds we can expect to find in .Site.Pages.
	allKindsInPages = []string{page.KindPage, page.KindHome, page.KindSection, page.KindTaxonomy, page.KindTaxonomyTerm}

	allKinds = append(allKindsInPages, []string{kindRSS, kindSitemap, kindRobotsTXT, kind404}...)
)

const (

	// Temporary state.
	kindUnknown = "unknown"

	// The following are (currently) temporary nodes,
	// i.e. nodes we create just to render in isolation.
	kindRSS       = "RSS"
	kindSitemap   = "sitemap"
	kindRobotsTXT = "robotsTXT"
	kind404       = "404"

	pageResourceType = "page"
)

var (
	_ page.Page           = (*pageState)(nil)
	_ collections.Grouper = (*pageState)(nil)
	_ collections.Slicer  = (*pageState)(nil)
)

var (
	pageTypesProvider = resource.NewResourceTypesProvider(media.OctetType, pageResourceType)
	zeroFile          = &source.FileInfo{} // TODO(bep) page check vs with
)

func newBuildState(metaProvider *pageMeta) (*pageState, error) {
	if metaProvider.s == nil {
		panic("must provide a Site")
	}

	if metaProvider.f == nil {
		metaProvider.f = zeroFile
	}

	s := metaProvider.s

	ps := &pageState{
		pagePerOutputProviders: nopPagePerOutput,
		PaginatorProvider:      page.NopPage,
		FileProvider:           metaProvider,
		AuthorProvider:         metaProvider,
		Scratcher:              maps.NewScratcher(),
		Positioner:             page.NopPage,
		InSectionPositioner:    page.NopPage,
		ResourceMetaProvider:   metaProvider,
		ResourceParamsProvider: metaProvider,
		PageMetaProvider:       metaProvider,
		OutputFormatsProvider:  page.NopPage,
		ResourceTypesProvider:  pageTypesProvider,
		ResourcePathsProvider:  page.NopPage,
		RefProvider:            page.NopPage,
		ShortcodeInfoProvider:  page.NopPage,
		LanguageProvider:       s,

		TODOProvider: page.NopPage,

		InternalDependencies: s,

		lateInit: lazy.New(),

		m: metaProvider,
		s: s,
	}

	// TODO(bep) page

	siteAdapter := pageSiteAdapter{s: s, p: ps}

	ps.pageMenus = &pageMenus{p: ps}
	ps.PageMenusProvider = ps.pageMenus
	ps.GetPageProvider = siteAdapter
	ps.GitInfoProvider = ps
	ps.TranslationsProvider = ps
	ps.ResourceDataProvider = &dataProvider{pageState: ps}
	ps.RawContentProvider = ps
	ps.ChildCareProvider = ps
	ps.TreeProvider = pageTreeProvider{p: ps}
	ps.Eqer = ps
	ps.TranslationKeyProvider = ps
	ps.ShortcodeInfoProvider = ps

	return ps, nil

}

// Used by the legacy 404, sitemap and robots.txt rendering
func newStandalonePage(m *pageMeta, f output.Format) (*pageState, error) {
	m.configuredOutputFormats = output.Formats{f}
	p, err := newBuildStatePageFromMeta(m)

	if err != nil {
		return nil, err
	}

	if err := p.initPage(); err != nil {
		return nil, err
	}

	return p, p.initOutputFormat(f, true)

}

func newBuildStatePageFromMeta(metaProvider *pageMeta) (*pageState, error) {
	ps, err := newBuildState(metaProvider)
	if err != nil {
		return nil, err
	}

	if err := metaProvider.applyDefaultValues(); err != nil {
		return nil, err
	}

	ps.lateInit.Add(func() (interface{}, error) {
		pp, err := newPagePaths(metaProvider.s, ps, metaProvider)
		if err != nil {
			return nil, err
		}

		provdidersPerOutput := func(f output.Format) (pagePerOutputProviders, error) {
			var targetPath targetPathString

			if ft, found := pp.targetPaths[f.Name]; found {
				targetPath = ft
			}

			return struct {
				page.ContentProvider
				page.PageRenderProvider
				page.AlternativeOutputFormatsProvider
				targetPather
			}{
				page.NopPage,
				page.NopPage,
				page.NopPage,
				targetPath,
			}, nil
		}

		ps.createOutputFormatProvider = provdidersPerOutput

		if err := ps.initCommonProviders(pp); err != nil {
			return nil, err
		}

		return nil, nil

	})

	return ps, err

}

func newBuildStatePageWithContent(f *fileInfo, s *Site, content resource.OpenReadSeekCloser) (*pageState, error) {

	sections := sectionsFromFile(f)
	kind := s.kindFromFileInfoOrSections(f, sections)

	metaProvider := &pageMeta{kind: kind, sections: sections, s: s, f: f}

	ps, err := newBuildState(metaProvider)
	if err != nil {
		return nil, err
	}

	gi, err := s.h.gitInfoForPage(ps)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load Git data")
	}
	ps.gitInfo = gi

	// TODO(bep) page simplify
	metaSetter := func(frontmatter map[string]interface{}) error {
		if err := metaProvider.setMetadata(ps, frontmatter); err != nil {
			return err
		}

		return nil
	}

	r, err := content()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	parseResult, err := pageparser.Parse(
		r,
		pageparser.Config{EnableEmoji: s.Cfg.GetBool("enableEmoji")},
	)
	if err != nil {
		return nil, err
	}

	ps.pageContent = pageContent{
		source: rawPageContent{
			parsed: parseResult,
		},
	}

	ps.shortcodeState = ps.newShortcodeHandler()

	if err := ps.mapContent(metaSetter); err != nil {
		return nil, ps.errWithFileContext(err)
	}

	if err := metaProvider.applyDefaultValues(); err != nil {
		return nil, err
	}

	ps.lateInit.Add(func() (interface{}, error) {
		// Provides content and render func per output format.
		perOutputFormatFn := newPageContentProvider(ps)

		// TODO(bep) page check permalink vs headless
		pp, err := newPagePaths(s, ps, metaProvider)
		if err != nil {
			fmt.Println("error:", err)
			return nil, err
		}

		provdidersPerOutput := func(f output.Format) (pagePerOutputProviders, error) {
			var targetPath targetPathString

			contentProvider, err := perOutputFormatFn(f)
			if err != nil {
				return nil, err
			}

			if ft, found := pp.targetPaths[f.Name]; found {
				targetPath = ft
			}

			return struct {
				page.ContentProvider
				page.PageRenderProvider
				targetPather
				page.AlternativeOutputFormatsProvider
			}{
				contentProvider,
				contentProvider,
				targetPath,
				contentProvider,
			}, nil
		}

		ps.createOutputFormatProvider = provdidersPerOutput

		if ps.IsNode() {
			ps.paginator = &pagePaginator{source: ps}
			ps.PaginatorProvider = ps.paginator
		}

		if err := ps.initCommonProviders(pp); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return ps, nil
}

type pageMenus struct {
	p *pageState

	q navigation.MenyQueryProvider

	pmInit sync.Once
	pm     navigation.PageMenus
}

func (p *pageMenus) init() {
	p.pmInit.Do(func() {
		p.q = navigation.NewMenuQueryProvider(
			p.p.s.Info.sectionPagesMenu,
			p,
			p.p.s,
			p.p,
		)

		// TODO(bep) page error handling
		p.pm, _ = navigation.PageMenusFromPage(p.p)

	})

}
func (p *pageMenus) menus() navigation.PageMenus {
	p.init()
	return p.pm

}

func (p *pageMenus) Menus() navigation.PageMenus {
	// There is a reverse dependency here. initMenus will, once, build the
	// site menus and update any relevant page.
	p.p.s.init.menus.Do()

	return p.menus()
}

func (p *pageMenus) HasMenuCurrent(menuID string, me *navigation.MenuEntry) bool {
	p.p.s.init.menus.Do()
	p.init()
	return p.q.HasMenuCurrent(menuID, me)
}

func (p *pageMenus) IsMenuCurrent(menuID string, inme *navigation.MenuEntry) bool {
	p.p.s.init.menus.Do()
	p.init()
	return p.q.IsMenuCurrent(menuID, inme)
}

func (s *Site) newPage(kind string, sections ...string) *pageState {
	p, err := newBuildStatePageFromMeta(&pageMeta{
		s:        s,
		kind:     kind,
		sections: sections,
	})

	if err != nil {
		panic(err)
	}

	return p
}

type dataProvider struct {
	*pageState

	dataInit sync.Once
	data     page.Data
}

type pageSiteAdapter struct {
	p page.Page
	s *Site
}

func (p pageSiteAdapter) GetPage(ref string) (page.Page, error) {
	return p.s.getPageNew(p.p, ref)
}

// TODO(bep) page name etc.
type pageState struct {
	s *Site

	m *pageMeta

	lateInit *lazy.Init

	currentOutputFormat        output.Format
	createOutputFormatProvider func(f output.Format) (pagePerOutputProviders, error)

	targetPathDescriptor page.TargetPathDescriptor
	relTargetPathBase    string

	pageContent

	gitInfo *gitmap.GitInfo

	// All of these represents a page.Page
	compare.Eqer
	pagePerOutputProviders
	page.AuthorProvider
	page.FileProvider
	page.GitInfoProvider
	page.GetPageProvider
	maps.Scratcher
	page.SitesProvider
	page.OutputFormatsProvider
	page.ChildCareProvider
	page.PageMetaProvider
	page.PaginatorProvider
	page.Positioner
	page.RefProvider
	page.InSectionPositioner
	page.ShortcodeInfoProvider
	page.RawContentProvider
	page.TODOProvider
	page.TranslationsProvider
	page.TreeProvider
	resource.LanguageProvider
	resource.ResourceDataProvider
	resource.ResourceMetaProvider
	resource.ResourceParamsProvider
	resource.ResourcePathsProvider
	resource.ResourceTypesProvider
	resource.TranslationKeyProvider
	navigation.PageMenusProvider

	paginator *pagePaginator

	// Positional navigation
	posNextPrev        *nextPrev
	posNextPrevSection *nextPrev

	// Menus
	pageMenus *pageMenus

	// Internal use
	page.InternalDependencies

	pagesInit sync.Once
	pages     page.Pages

	// Any bundled resources
	resourcesInit sync.Once
	resources     resource.Resources

	translations    page.Pages
	allTranslations page.Pages

	// Calculated an cached translation mapping key
	translationKey     string
	translationKeyInit sync.Once

	// Will only be set for sections and regular pages.
	parent *pageState

	// Will only be set for section pages and the home page.
	subSections page.Pages

	// Set in fast render mode to force render a given page.
	forceRender bool
}

// AllTranslations returns all translations, including the current Page.
func (p *pageState) AllTranslations() page.Pages {
	p.s.h.init.translations.Do()
	return p.allTranslations
}

func (p *dataProvider) Data() interface{} {
	p.dataInit.Do(func() {
		p.data = make(page.Data)

		if p.Kind() == page.KindPage {
			return
		}

		switch p.Kind() {
		case page.KindTaxonomy:
			plural := p.SectionsEntries()[0]
			term := p.SectionsEntries()[1]

			if p.s.Info.preserveTaxonomyNames {
				if v, ok := p.s.taxonomiesOrigKey[fmt.Sprintf("%s-%s", plural, term)]; ok {
					term = v
				}
			}

			singular := p.s.taxonomiesPluralSingular[plural]
			taxonomy := p.s.Taxonomies[plural].Get(term)

			p.data[singular] = taxonomy
			p.data["Singular"] = singular
			p.data["Plural"] = plural
			p.data["Term"] = term
		case page.KindTaxonomyTerm:
			plural := p.SectionsEntries()[0]
			singular := p.s.taxonomiesPluralSingular[plural]

			p.data["Singular"] = singular
			p.data["Plural"] = plural
			p.data["Terms"] = p.s.Taxonomies[plural]
			// keep the following just for legacy reasons
			p.data["OrderedIndex"] = p.data["Terms"]
			p.data["Index"] = p.data["Terms"]
		}

		// Assign the function to the map to make sure it is lazily initialized
		p.data["pages"] = p.Pages

	})

	return p.data
}

func (p *pageState) GitInfo() *gitmap.GitInfo {
	return p.gitInfo
}

// Eq returns whether the current page equals the given page.
// This is what's invoked when doing `{{ if eq $page $otherPage }}`
func (p *pageState) Eq(other interface{}) bool {
	pp, err := unwrapPage(other)
	if err != nil {
		return false
	}

	return p == pp
}

func (p *pageState) HasShortcode(name string) bool {
	if p.shortcodeState == nil {
		return false
	}

	return p.shortcodeState.nameSet[name]
}

func (p *pageState) Hugo() hugo.Info {
	return p.s.Info.hugoInfo
}

// IsTranslated returns whether this content file is translated to
// other language(s).
func (p *pageState) IsTranslated() bool {
	p.s.h.init.translations.Do()
	return len(p.translations) > 0
}

func (p *pageState) LanguagePrefix() string {
	return p.s.Info.LanguagePrefix
}

func (p *pageState) Pages() page.Pages {
	p.pagesInit.Do(func() {
		if p.pages != nil {
			return
		}

		var pages page.Pages

		switch p.Kind() {
		case page.KindPage:
		// No pages for you.
		case page.KindHome:
			pages = p.s.RegularPages()
		case page.KindTaxonomy:
			plural := p.SectionsEntries()[0]
			term := p.SectionsEntries()[1]

			if p.s.Info.preserveTaxonomyNames {
				if v, ok := p.s.taxonomiesOrigKey[fmt.Sprintf("%s-%s", plural, term)]; ok {
					term = v
				}
			}

			taxonomy := p.s.Taxonomies[plural].Get(term)
			pages = taxonomy.Pages()

		case page.KindTaxonomyTerm:
			plural := p.SectionsEntries()[0]
			// A list of all page.KindTaxonomy pages with matching plural
			// TODO(bep) page
			for _, p := range p.s.findPagesByKind(page.KindTaxonomy) {
				if p.SectionsEntries()[0] == plural {
					pages = append(pages, p)
				}
			}
		case kind404, kindSitemap, kindRobotsTXT:
			pages = p.s.Pages()
		}

		p.pages = pages
	})

	return p.pages
}

// RawContent returns the un-rendered source content without
// any leading front matter.
func (p *pageState) RawContent() string {
	if p.source.parsed == nil {
		return ""
	}
	start := p.source.posMainContent
	if start == -1 {
		start = 0
	}
	return string(p.source.parsed.Input()[start:])
}

func (p *pageState) Resources() resource.Resources {
	p.resourcesInit.Do(func() {

		sort := func() {
			sort.SliceStable(p.resources, func(i, j int) bool {
				ri, rj := p.resources[i], p.resources[j]
				if ri.ResourceType() < rj.ResourceType() {
					return true
				}

				p1, ok1 := ri.(page.Page)
				p2, ok2 := rj.(page.Page)

				if ok1 != ok2 {
					return ok2
				}

				if ok1 {
					return page.DefaultPageSort(p1, p2)
				}

				return ri.RelPermalink() < rj.RelPermalink()
			})
		}

		sort()

		if len(p.m.resourcesMetadata) > 0 {
			resources.AssignMetadata(p.m.resourcesMetadata, p.resources...)
			sort()
		}

	})
	return p.resources
}

func (p *pageState) Site() page.Site {
	return &p.s.Info
}

func (p *pageState) String() string {
	if sourceRef := p.sourceRef(); sourceRef != "" {
		return fmt.Sprintf("Page(%s)", sourceRef)
	}
	return fmt.Sprintf("Page(%q)", p.Title())
}

// TranslationKey returns the key used to map language translations of this page.
// It will use the translationKey set in front matter if set, or the content path and
// filename (excluding any language code and extension), e.g. "about/index".
// The Page Kind is always prepended.
func (p *pageState) TranslationKey() string {
	p.translationKeyInit.Do(func() {
		if p.m.translationKey != "" {
			p.translationKey = p.Kind() + "/" + p.m.translationKey
		} else if p.IsPage() && p.File() != nil {
			p.translationKey = path.Join(p.Kind(), filepath.ToSlash(p.File().Dir()), p.File().TranslationBaseName())
		} else if p.IsNode() {
			p.translationKey = path.Join(p.Kind(), p.SectionsPath())
		}

	})

	return p.translationKey

}

// Translations returns the translations excluding the current Page.
func (p *pageState) Translations() page.Pages {
	p.s.h.init.translations.Do()
	return p.translations
}

func (p *pageState) addResources(r ...resource.Resource) {
	p.resources = append(p.resources, r...)
}

func (p *pageState) addSectionToParent() {
	if p.parent == nil {
		return
	}
	p.parent.subSections = append(p.parent.subSections, p)
}

func (p *pageState) contentMarkupType() string {
	if p.m.markup != "" {
		return p.m.markup

	}
	return p.File().Ext()
}

func (p *pageState) createLayoutDescriptor() output.LayoutDescriptor {
	var section string
	sections := p.SectionsEntries()

	switch p.Kind() {
	case page.KindSection:
		section = sections[0]
	case page.KindTaxonomy, page.KindTaxonomyTerm:
		section = p.s.taxonomiesPluralSingular[sections[0]]
	default:
	}

	return output.LayoutDescriptor{
		Kind:    p.Kind(),
		Type:    p.Type(),
		Lang:    p.Language().Lang,
		Layout:  p.Layout(),
		Section: section,
	}
}

func (p *pageState) errWithFileContext(err error) error {

	err, _ = herrors.WithFileContextForFile(
		err,
		p.File().Filename(),
		p.File().Filename(),
		p.s.SourceSpec.Fs.Source,
		herrors.SimpleLineMatcher)

	return err
}

func (p *pageState) errorf(err error, format string, a ...interface{}) error {
	if herrors.UnwrapErrorWithFileContext(err) != nil {
		// More isn't always better.
		return err
	}
	args := append([]interface{}{p.Language().Lang, p.pathOrTitle()}, a...)
	format = "[%s] page %q: " + format
	if err == nil {
		errors.Errorf(format, args...)
		return fmt.Errorf(format, args...)
	}
	return errors.Wrapf(err, format, args...)
}

func (p *pageState) getLayouts(f output.Format, layouts ...string) ([]string, error) {

	if len(layouts) == 0 {
		selfLayout := p.selfLayoutForOutput(f)
		if selfLayout != "" {
			return []string{selfLayout}, nil
		}
	}

	// TODO(bep) page cache
	layoutDescriptor := p.createLayoutDescriptor()

	if len(layouts) > 0 {
		layoutDescriptor.Layout = layouts[0]
		layoutDescriptor.LayoutOverride = true
	}

	return p.s.layoutHandler.For(
		layoutDescriptor,
		f)
}

func (ps *pageState) initCommonProviders(pp pagePaths) error {
	if ps.IsPage() {
		ps.posNextPrev = &nextPrev{init: ps.s.init.prevNext}
		ps.posNextPrevSection = &nextPrev{init: ps.s.init.prevNextInSection}
		ps.InSectionPositioner = newPagePositionInSection(ps.posNextPrevSection)
		ps.Positioner = newPagePosition(ps.posNextPrev)
	}

	ps.ResourcePathsProvider = pp
	ps.OutputFormatsProvider = pp
	ps.targetPathDescriptor = pp.targetPathDescriptor
	ps.relTargetPathBase = pp.relTargetPathBase
	ps.RefProvider = newPageRef(ps)
	ps.SitesProvider = &ps.s.Info

	return nil
}

func (s *Site) kindFromSections(sections []string) string {
	if len(sections) == 0 || len(s.siteConfigHolder.taxonomiesConfig) == 0 {
		return page.KindSection
	}

	sectionPath := path.Join(sections...)

	for _, plural := range s.siteConfigHolder.taxonomiesConfig {
		if plural == sectionPath {
			return page.KindTaxonomyTerm
		}

		if strings.HasPrefix(sectionPath, plural) {
			return page.KindTaxonomy
		}

	}

	return page.KindSection
}

func (p *pageState) mapContent(
	metaSetter func(frontmatter map[string]interface{}) error) error {

	s := p.shortcodeState

	p.renderable = true
	p.source.posMainContent = -1

	result := bp.GetBuffer()
	defer bp.PutBuffer(result)

	iter := p.source.parsed.Iterator()

	fail := func(err error, i pageparser.Item) error {
		return p.parseError(err, iter.Input(), i.Pos)
	}

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currShortcode shortcode
	var ordinal int

Loop:
	for {
		it := iter.Next()

		switch {
		case it.Type == pageparser.TypeIgnore:
		case it.Type == pageparser.TypeHTMLStart:
			// This is HTML without front matter. It can still have shortcodes.
			p.selfLayout = "__" + p.File().Filename()
			p.renderable = false
			result.Write(it.Val)
		case it.IsFrontMatter():
			f := metadecoders.FormatFromFrontMatterType(it.Type)
			m, err := metadecoders.Default.UnmarshalToMap(it.Val, f)
			if err != nil {
				if fe, ok := err.(herrors.FileError); ok {
					return herrors.ToFileErrorWithOffset(fe, iter.LineNumber()-1)
				} else {
					return err
				}
			}

			if err := metaSetter(m); err != nil {
				return err
			}

			next := iter.Peek()
			if !next.IsDone() {
				p.source.posMainContent = next.Pos
			}

			if !p.s.shouldBuild(p) {
				// Nothing more to do.
				return nil
			}

		case it.Type == pageparser.TypeLeadSummaryDivider:
			result.Write(internalSummaryDividerPre)
			p.source.hasSummaryDivider = true
			// Need to determine if the page is truncated.
			f := func(item pageparser.Item) bool {
				if item.IsNonWhitespace() {
					p.truncated = true

					// Done
					return false
				}
				return true
			}
			iter.PeekWalk(f)

		// Handle shortcode
		case it.IsLeftShortcodeDelim():
			// let extractShortcode handle left delim (will do so recursively)
			iter.Backup()

			currShortcode, err := s.extractShortcode(ordinal, iter, p)

			if currShortcode.name != "" {
				s.nameSet[currShortcode.name] = true
			}

			if err != nil {
				return fail(errors.Wrap(err, "failed to extract shortcode"), it)
			}

			if currShortcode.params == nil {
				currShortcode.params = make([]string, 0)
			}

			placeHolder := s.createShortcodePlaceholder()
			result.WriteString(placeHolder)
			ordinal++
			s.shortcodes.Add(placeHolder, currShortcode)
		case it.Type == pageparser.TypeEmoji:
			if emoji := helpers.Emoji(it.ValStr()); emoji != nil {
				result.Write(emoji)
			} else {
				result.Write(it.Val)
			}
		case it.IsEOF():
			break Loop
		case it.IsError():
			err := fail(errors.WithStack(errors.New(it.ValStr())), it)
			currShortcode.err = err
			return err

		default:
			result.Write(it.Val)
		}
	}

	resultBytes := make([]byte, result.Len())
	copy(resultBytes, result.Bytes())
	p.workContent = resultBytes

	return nil
}

func (p *pageState) newShortcodeHandler() *shortcodeHandler {

	s := &shortcodeHandler{
		p:                      p,
		s:                      p.s,
		enableInlineShortcodes: p.s.enableInlineShortcodes,
		contentShortcodes:      newOrderedMap(),
		shortcodes:             newOrderedMap(),
		nameSet:                make(map[string]bool),
		renderedShortcodes:     make(map[string]string),
	}

	var placeholderFunc func() string // TODO(bep) page p.s.shortcodePlaceholderFunc
	if placeholderFunc == nil {
		placeholderFunc = func() string {
			return fmt.Sprintf("HAHA%s-%p-%d-HBHB", shortcodePlaceholderPrefix, p, s.nextPlaceholderID())
		}

	}

	s.placeholderFunc = placeholderFunc

	return s
}

func (p *pageState) outputFormat() output.Format {
	return p.currentOutputFormat
}

func (p *pageState) parseError(err error, input []byte, offset int) error {
	if herrors.UnwrapFileError(err) != nil {
		// Use the most specific location.
		return err
	}
	pos := p.posFromInput(input, offset)
	return herrors.NewFileError("md", -1, pos.LineNumber, pos.ColumnNumber, err)

}

func (p *pageState) pathOrTitle() string {
	if p.File() != nil {
		return p.File().Filename()
	}

	if p.Path() != "" {
		return p.Path()
	}

	return p.Title()
}

func (p *pageState) posFromInput(input []byte, offset int) text.Position {
	lf := []byte("\n")
	input = input[:offset]
	lineNumber := bytes.Count(input, lf) + 1
	endOfLastLine := bytes.LastIndex(input, lf)

	return text.Position{
		Filename:     p.pathOrTitle(),
		LineNumber:   lineNumber,
		ColumnNumber: offset - endOfLastLine,
		Offset:       offset,
	}
}

// TODO(bep) page name
type pageInternals interface {
	posOffset(offset int) text.Position
}

func (p *pageState) posOffset(offset int) text.Position {
	return p.posFromInput(p.source.parsed.Input(), offset)
}

func (p *pageState) subResourceTargetPathFactory(base string) string {
	return path.Join(p.relTargetPathBase, base)
}

func (p *pageState) setPages(pages page.Pages) {
	page.SortByDefault(pages)
	p.pages = pages
}

func (p *pageState) setTranslations(pages page.Pages) {
	p.allTranslations = pages
	page.SortByLanguage(p.allTranslations)
	translations := make(page.Pages, 0)
	for _, t := range p.allTranslations {
		if !t.Eq(p) {
			translations = append(translations, t)
		}
	}
	p.translations = translations
}

// Must be run after the site section tree etc. is built and ready.
func (p *pageState) initPage() error {
	if _, err := p.lateInit.Do(); err != nil {
		return err
	}
	return nil
}

// shiftToOutputFormat is serialized.
func (p *pageState) shiftToOutputFormat(f output.Format) error {
	if err := p.initPage(); err != nil {
		return err
	}

	providers, err := p.createOutputFormatProvider(f)
	if err != nil {
		return err
	}

	p.pagePerOutputProviders = providers
	p.currentOutputFormat = f

	for _, r := range p.Resources().ByType(pageResourceType) {
		rp := r.(*pageState)
		if err := rp.shiftToOutputFormat(f); err != nil {
			return errors.Wrap(err, "failed to shift outputformat in Page resource")
		}
	}

	return nil
}

func (p *pageState) renderResources() error {
	for _, r := range p.Resources() {
		src, ok := r.(resource.Source)
		if !ok {
			// Pages gets rendered with the owning page.
			continue
		}

		if err := src.Publish(); err != nil {
			if os.IsNotExist(err) {
				// The resource has been deleted from the file system.
				// This should be extremely rare, but can happen on live reload in server
				// mode when the same resource is member of different page bundles.
				// TODO(bep) page p.deleteResource(i)
			} else {
				p.s.Log.ERROR.Printf("Failed to publish Resource for page %q: %s", p.pathOrTitle(), err)
			}
		} else {
			p.s.PathSpec.ProcessingStats.Incr(&p.s.PathSpec.ProcessingStats.Files)
		}
	}
	return nil
}

// is serialized
func (p *pageState) initOutputFormat(f output.Format, start bool) error {
	if err := p.shiftToOutputFormat(f); err != nil {
		return err
	}

	if start {
		if !p.renderable {
			if _, err := p.Content(); err != nil {
				return err
			}
		}

		if p.IsNode() {
			p.paginator = newPagePaginator(p)
			p.PaginatorProvider = p.paginator
		}
	}

	return nil

}

func (p *pageState) sortParentSections() {
	if p.parent == nil {
		return
	}
	page.SortByDefault(p.parent.subSections)
}

// sourceRef returns the canonical, absolute fully-qualifed logical reference used by
// methods such as GetPage and ref/relref shortcodes to refer to
// this page. It is prefixed with a "/".
//
// For pages that have a source file, it is returns the path to this file as an
// absolute path rooted in this site's content dir.
// For pages that do not (sections witout content page etc.), it returns the
// virtual path, consistent with where you would add a source file.
func (p *pageState) sourceRef() string {
	if p.File() != nil {
		sourcePath := p.File().Path()
		if sourcePath != "" {
			return "/" + filepath.ToSlash(sourcePath)
		}
	}

	if len(p.SectionsEntries()) > 0 {
		// no backing file, return the virtual source path
		return "/" + p.SectionsPath()
	}

	return ""
}

// TODO(bep) page
func (p *pageState) MarshalJSON() ([]byte, error) {
	s := struct {
		Title string
	}{
		Title: p.Title(),
	}

	return json.Marshal(&s)

}

// TODO(bep) page
func (p *pageState) updatePageDates() {
	// TODO(bep) there is a potential issue with page sorting for home pages
	// etc. without front matter dates set, but let us wrap the head around
	// that in another time.
	if true {
		return
	}
	/*
		if !p.Date().IsZero() {
			if p.Lastmod().IsZero() {
				updater.FLastmod = p.Date()
			}
			return
		} else if !p.Lastmod().IsZero() {
			if p.Date().IsZero() {
				updater.FDate = p.Lastmod()
			}
			return
		}

		// Set it to the first non Zero date in children
		var foundDate, foundLastMod bool

		for _, child := range p.Pages() {
			if !child.Date().IsZero() {
				updater.FDate = child.Date()
				foundDate = true
			}
			if !child.Lastmod().IsZero() {
				updater.FLastmod = child.Lastmod()
				foundLastMod = true
			}

			if foundDate && foundLastMod {
				break
			}
		}
	*/
}

type pageStatePages []*pageState

// Implement sorting.
func (ps pageStatePages) Len() int { return len(ps) }

func (ps pageStatePages) Less(i, j int) bool { return page.DefaultPageSort(ps[i], ps[j]) }

func (ps pageStatePages) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

// findPagePos Given a page, it will find the position in Pages
// will return -1 if not found
func (ps pageStatePages) findPagePos(page *pageState) int {
	for i, x := range ps {
		if x.File().Filename() == page.File().Filename() {
			return i
		}
	}
	return -1
}

func (ps pageStatePages) findPagePosByFilename(filename string) int {
	for i, x := range ps {
		if x.File().Filename() == filename {
			return i
		}
	}
	return -1
}

func (ps pageStatePages) findPagePosByFilnamePrefix(prefix string) int {
	if prefix == "" {
		return -1
	}

	lenDiff := -1
	currPos := -1
	prefixLen := len(prefix)

	// Find the closest match
	for i, x := range ps {
		if strings.HasPrefix(x.File().Filename(), prefix) {
			diff := len(x.File().Filename()) - prefixLen
			if lenDiff == -1 || diff < lenDiff {
				lenDiff = diff
				currPos = i
			}
		}
	}
	return currPos
}

func content(c resource.ContentProvider) string {
	cc, err := c.Content()
	if err != nil {
		panic(err)
	}

	ccs, err := cast.ToStringE(cc)
	if err != nil {
		panic(err)
	}
	return ccs
}

func (s *Site) kindFromFileInfoOrSections(fi *fileInfo, sections []string) string {
	if fi.TranslationBaseName() == "_index" {
		if fi.Dir() == "" {
			return page.KindHome
		}

		return s.kindFromSections(sections)

	}
	return page.KindPage
}

func sectionsFromFile(fi source.File) []string {
	dirname := fi.Dir()
	dirname = strings.Trim(dirname, helpers.FilePathSeparator)
	if dirname == "" {
		return nil
	}
	parts := strings.Split(dirname, helpers.FilePathSeparator)

	if fii, ok := fi.(*fileInfo); ok {
		if fii.bundleTp == bundleLeaf && len(parts) > 0 {
			// my-section/mybundle/index.md => my-section
			return parts[:len(parts)-1]
		}
	}

	return parts
}

func stackTrace(length int) string {
	trace := make([]byte, length)
	runtime.Stack(trace, true)
	return string(trace)
}
