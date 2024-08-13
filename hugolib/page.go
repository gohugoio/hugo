// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/hugolib/segments"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/output/layouts"
	"github.com/gohugoio/hugo/related"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"

	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	_ page.Page                                = (*pageState)(nil)
	_ collections.Grouper                      = (*pageState)(nil)
	_ collections.Slicer                       = (*pageState)(nil)
	_ identity.DependencyManagerScopedProvider = (*pageState)(nil)
	_ contentNodeI                             = (*pageState)(nil)
	_ pageContext                              = (*pageState)(nil)
)

var (
	pageTypesProvider = resource.NewResourceTypesProvider(media.Builtin.OctetType, pageResourceType)
	nopPageOutput     = &pageOutput{
		pagePerOutputProviders: nopPagePerOutput,
		MarkupProvider:         page.NopPage,
		ContentProvider:        page.NopPage,
	}
)

// pageContext provides contextual information about this page, for error
// logging and similar.
type pageContext interface {
	posOffset(offset int) text.Position
	wrapError(err error) error
	getContentConverter() converter.Converter
}

type pageSiteAdapter struct {
	p page.Page
	s *Site
}

func (pa pageSiteAdapter) GetPage(ref string) (page.Page, error) {
	p, err := pa.s.getPage(pa.p, ref)

	if p == nil {
		// The nil struct has meaning in some situations, mostly to avoid breaking
		// existing sites doing $nilpage.IsDescendant($p), which will always return
		// false.
		p = page.NilPage
	}
	return p, err
}

type pageState struct {
	// Incremented for each new page created.
	// Note that this will change between builds for a given Page.
	pid uint64

	// This slice will be of same length as the number of global slice of output
	// formats (for all sites).
	pageOutputs []*pageOutput

	// Used to determine if we can reuse content across output formats.
	pageOutputTemplateVariationsState *atomic.Uint32

	// This will be shifted out when we start to render a new output format.
	pageOutputIdx int
	*pageOutput

	// Common for all output formats.
	*pageCommon

	resource.Staler
	dependencyManager    identity.Manager
	resourcesPublishInit *sync.Once
}

func (p *pageState) IdentifierBase() string {
	return p.Path()
}

func (p *pageState) GetIdentity() identity.Identity {
	return p
}

func (p *pageState) ForEeachIdentity(f func(identity.Identity) bool) bool {
	return f(p)
}

func (p *pageState) GetDependencyManager() identity.Manager {
	return p.dependencyManager
}

func (p *pageState) GetDependencyManagerForScope(scope int) identity.Manager {
	switch scope {
	case pageDependencyScopeDefault:
		return p.dependencyManagerOutput
	case pageDependencyScopeGlobal:
		return p.dependencyManager
	default:
		return identity.NopManager
	}
}

func (p *pageState) Key() string {
	return "page-" + strconv.FormatUint(p.pid, 10)
}

func (p *pageState) resetBuildState() {
	p.Scratcher = maps.NewScratcher()
}

func (p *pageState) reusePageOutputContent() bool {
	return p.pageOutputTemplateVariationsState.Load() == 1
}

func (p *pageState) skipRender() bool {
	b := p.s.conf.C.SegmentFilter.ShouldExcludeFine(
		segments.SegmentMatcherFields{
			Path:   p.Path(),
			Kind:   p.Kind(),
			Lang:   p.Lang(),
			Output: p.pageOutput.f.Name,
		},
	)

	return b
}

func (po *pageState) isRenderedAny() bool {
	for _, o := range po.pageOutputs {
		if o.isRendered() {
			return true
		}
	}
	return false
}

func (p *pageState) isContentNodeBranch() bool {
	return p.IsNode()
}

func (p *pageState) Err() resource.ResourceError {
	return nil
}

// Eq returns whether the current page equals the given page.
// This is what's invoked when doing `{{ if eq $page $otherPage }}`
func (p *pageState) Eq(other any) bool {
	pp, err := unwrapPage(other)
	if err != nil {
		return false
	}

	return p == pp
}

func (p *pageState) HeadingsFiltered(context.Context) tableofcontents.Headings {
	return nil
}

type pageHeadingsFiltered struct {
	*pageState
	headings tableofcontents.Headings
}

func (p *pageHeadingsFiltered) HeadingsFiltered(context.Context) tableofcontents.Headings {
	return p.headings
}

func (p *pageHeadingsFiltered) page() page.Page {
	return p.pageState
}

// For internal use by the related content feature.
func (p *pageState) ApplyFilterToHeadings(ctx context.Context, fn func(*tableofcontents.Heading) bool) related.Document {
	fragments := p.pageOutput.pco.c().Fragments(ctx)
	headings := fragments.Headings.FilterBy(fn)
	return &pageHeadingsFiltered{
		pageState: p,
		headings:  headings,
	}
}

func (p *pageState) GitInfo() source.GitInfo {
	return p.gitInfo
}

func (p *pageState) CodeOwners() []string {
	return p.codeowners
}

// GetTerms gets the terms defined on this page in the given taxonomy.
// The pages returned will be ordered according to the front matter.
func (p *pageState) GetTerms(taxonomy string) page.Pages {
	return p.s.pageMap.getTermsForPageInTaxonomy(p.Path(), taxonomy)
}

func (p *pageState) MarshalJSON() ([]byte, error) {
	return page.MarshalPageToJSON(p)
}

func (p *pageState) RegularPagesRecursive() page.Pages {
	switch p.Kind() {
	case kinds.KindSection, kinds.KindHome:
		return p.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    p.Path(),
					Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindPage),
				},
				Recursive: true,
			},
		)
	default:
		return p.RegularPages()
	}
}

func (p *pageState) PagesRecursive() page.Pages {
	return nil
}

func (p *pageState) RegularPages() page.Pages {
	switch p.Kind() {
	case kinds.KindPage:
	case kinds.KindSection, kinds.KindHome, kinds.KindTaxonomy:
		return p.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    p.Path(),
					Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindPage),
				},
			},
		)
	case kinds.KindTerm:
		return p.s.pageMap.getPagesWithTerm(
			pageMapQueryPagesBelowPath{
				Path:    p.Path(),
				Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindPage),
			},
		)
	default:
		return p.s.RegularPages()
	}
	return nil
}

func (p *pageState) Pages() page.Pages {
	switch p.Kind() {
	case kinds.KindPage:
	case kinds.KindSection, kinds.KindHome:
		return p.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    p.Path(),
					KeyPart: "page-section",
					Include: pagePredicates.ShouldListLocal.And(
						pagePredicates.KindPage.Or(pagePredicates.KindSection),
					),
				},
			},
		)
	case kinds.KindTerm:
		return p.s.pageMap.getPagesWithTerm(
			pageMapQueryPagesBelowPath{
				Path: p.Path(),
			},
		)
	case kinds.KindTaxonomy:
		return p.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    p.Path(),
					KeyPart: "term",
					Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindTerm),
				},
				Recursive: true,
			},
		)
	default:
		return p.s.Pages()
	}
	return nil
}

// RawContent returns the un-rendered source content without
// any leading front matter.
func (p *pageState) RawContent() string {
	if p.m.content.pi.itemsStep2 == nil {
		return ""
	}
	start := p.m.content.pi.posMainContent
	if start == -1 {
		start = 0
	}
	source, err := p.m.content.pi.contentSource(p.m.content)
	if err != nil {
		panic(err)
	}
	return string(source[start:])
}

func (p *pageState) Resources() resource.Resources {
	return p.s.pageMap.getOrCreateResourcesForPage(p)
}

func (p *pageState) HasShortcode(name string) bool {
	if p.m.content.shortcodeState == nil {
		return false
	}

	return p.m.content.shortcodeState.hasName(name)
}

func (p *pageState) Site() page.Site {
	return p.sWrapped
}

func (p *pageState) String() string {
	return fmt.Sprintf("Page(%s)", p.Path())
}

// IsTranslated returns whether this content file is translated to
// other language(s).
func (p *pageState) IsTranslated() bool {
	return len(p.Translations()) > 0
}

// TranslationKey returns the key used to identify a translation of this content.
func (p *pageState) TranslationKey() string {
	if p.m.pageConfig.TranslationKey != "" {
		return p.m.pageConfig.TranslationKey
	}
	return p.Path()
}

// AllTranslations returns all translations, including the current Page.
func (p *pageState) AllTranslations() page.Pages {
	key := p.Path() + "/" + "translations-all"
	// This is called from Translations, so we need to use a different partition, cachePages2,
	// to avoid potential deadlocks.
	pages, err := p.s.pageMap.getOrCreatePagesFromCache(p.s.pageMap.cachePages2, key, func(string) (page.Pages, error) {
		if p.m.pageConfig.TranslationKey != "" {
			// translationKey set by user.
			pas, _ := p.s.h.translationKeyPages.Get(p.m.pageConfig.TranslationKey)
			pasc := make(page.Pages, len(pas))
			copy(pasc, pas)
			page.SortByLanguage(pasc)
			return pasc, nil
		}
		var pas page.Pages
		p.s.pageMap.treePages.ForEeachInDimension(p.Path(), doctree.DimensionLanguage.Index(),
			func(n contentNodeI) bool {
				if n != nil {
					pas = append(pas, n.(page.Page))
				}
				return false
			},
		)

		pas = pagePredicates.ShouldLink.Filter(pas)
		page.SortByLanguage(pas)
		return pas, nil
	})
	if err != nil {
		panic(err)
	}

	return pages
}

// Translations returns the translations excluding the current Page.
func (p *pageState) Translations() page.Pages {
	key := p.Path() + "/" + "translations"
	pages, err := p.s.pageMap.getOrCreatePagesFromCache(nil, key, func(string) (page.Pages, error) {
		var pas page.Pages
		for _, pp := range p.AllTranslations() {
			if !pp.Eq(p) {
				pas = append(pas, pp)
			}
		}
		return pas, nil
	})
	if err != nil {
		panic(err)
	}
	return pages
}

func (ps *pageState) initCommonProviders(pp pagePaths) error {
	if ps.IsPage() {
		ps.posNextPrev = &nextPrev{init: ps.s.init.prevNext}
		ps.posNextPrevSection = &nextPrev{init: ps.s.init.prevNextInSection}
		ps.InSectionPositioner = newPagePositionInSection(ps.posNextPrevSection)
		ps.Positioner = newPagePosition(ps.posNextPrev)
	}

	ps.OutputFormatsProvider = pp
	ps.targetPathDescriptor = pp.targetPathDescriptor
	ps.RefProvider = newPageRef(ps)
	ps.SitesProvider = ps.s

	return nil
}

func (p *pageState) getLayoutDescriptor() layouts.LayoutDescriptor {
	p.layoutDescriptorInit.Do(func() {
		var section string
		sections := p.SectionsEntries()

		switch p.Kind() {
		case kinds.KindSection:
			if len(sections) > 0 {
				section = sections[0]
			}
		case kinds.KindTaxonomy, kinds.KindTerm:

			if p.m.singular != "" {
				section = p.m.singular
			} else if len(sections) > 0 {
				section = sections[0]
			}
		default:
		}

		p.layoutDescriptor = layouts.LayoutDescriptor{
			Kind:    p.Kind(),
			Type:    p.Type(),
			Lang:    p.Language().Lang,
			Layout:  p.Layout(),
			Section: section,
		}
	})

	return p.layoutDescriptor
}

func (p *pageState) resolveTemplate(layouts ...string) (tpl.Template, bool, error) {
	f := p.outputFormat()

	d := p.getLayoutDescriptor()

	if len(layouts) > 0 {
		d.Layout = layouts[0]
		d.LayoutOverride = true
	}

	return p.s.Tmpl().LookupLayout(d, f)
}

// Must be run after the site section tree etc. is built and ready.
func (p *pageState) initPage() error {
	if _, err := p.init.Do(context.Background()); err != nil {
		return err
	}
	return nil
}

func (p *pageState) renderResources() error {
	var initErr error

	p.resourcesPublishInit.Do(func() {
		for _, r := range p.Resources() {
			if _, ok := r.(page.Page); ok {
				// Pages gets rendered with the owning page but we count them here.
				p.s.PathSpec.ProcessingStats.Incr(&p.s.PathSpec.ProcessingStats.Pages)
				continue
			}

			if _, isWrapper := r.(resource.ResourceWrapper); isWrapper {
				// Skip resources that are wrapped.
				// These gets published on its own.
				continue
			}

			src, ok := r.(resource.Source)
			if !ok {
				initErr = fmt.Errorf("resource %T does not support resource.Source", r)
				return
			}

			if err := src.Publish(); err != nil {
				if !herrors.IsNotExist(err) {
					p.s.Log.Errorf("Failed to publish Resource for page %q: %s", p.pathOrTitle(), err)
				}
			} else {
				p.s.PathSpec.ProcessingStats.Incr(&p.s.PathSpec.ProcessingStats.Files)
			}
		}
	})

	return initErr
}

func (p *pageState) AlternativeOutputFormats() page.OutputFormats {
	f := p.outputFormat()
	var o page.OutputFormats
	for _, of := range p.OutputFormats() {
		if of.Format.NotAlternative || of.Format.Name == f.Name {
			continue
		}

		o = append(o, of)
	}
	return o
}

type renderStringOpts struct {
	Display string
	Markup  string
}

var defaultRenderStringOpts = renderStringOpts{
	Display: "inline",
	Markup:  "", // Will inherit the page's value when not set.
}

func (p *pageMeta) wrapError(err error, sourceFs afero.Fs) error {
	if err == nil {
		panic("wrapError with nil")
	}

	if p.File() == nil {
		// No more details to add.
		return fmt.Errorf("%q: %w", p.Path(), err)
	}

	return hugofs.AddFileInfoToError(err, p.File().FileInfo(), sourceFs)
}

// wrapError adds some more context to the given error if possible/needed
func (p *pageState) wrapError(err error) error {
	return p.m.wrapError(err, p.s.h.SourceFs)
}

func (p *pageState) getPageInfoForError() string {
	s := fmt.Sprintf("kind: %q, path: %q", p.Kind(), p.Path())
	if p.File() != nil {
		s += fmt.Sprintf(", file: %q", p.File().Filename())
	}
	return s
}

func (p *pageState) getContentConverter() converter.Converter {
	var err error
	p.contentConverterInit.Do(func() {
		if p.m.pageConfig.ContentMediaType.IsZero() {
			panic("ContentMediaType not set")
		}
		markup := p.m.pageConfig.ContentMediaType.SubType

		if markup == "html" {
			// Only used for shortcode inner content.
			markup = "markdown"
		}
		p.contentConverter, err = p.m.newContentConverter(p, markup)
	})

	if err != nil {
		p.s.Log.Errorln("Failed to create content converter:", err)
	}
	return p.contentConverter
}

func (p *pageState) errorf(err error, format string, a ...any) error {
	if herrors.UnwrapFileError(err) != nil {
		// More isn't always better.
		return err
	}
	args := append([]any{p.Language().Lang, p.pathOrTitle()}, a...)
	args = append(args, err)
	format = "[%s] page %q: " + format + ": %w"
	if err == nil {
		return fmt.Errorf(format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (p *pageState) outputFormat() (f output.Format) {
	if p.pageOutput == nil {
		panic("no pageOutput")
	}
	return p.pageOutput.f
}

func (p *pageState) parseError(err error, input []byte, offset int) error {
	pos := posFromInput("", input, offset)
	return herrors.NewFileErrorFromName(err, p.File().Filename()).UpdatePosition(pos)
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
	return posFromInput(p.pathOrTitle(), input, offset)
}

func (p *pageState) posOffset(offset int) text.Position {
	return p.posFromInput(p.m.content.mustSource(), offset)
}

// shiftToOutputFormat is serialized. The output format idx refers to the
// full set of output formats for all sites.
// This is serialized.
func (p *pageState) shiftToOutputFormat(isRenderingSite bool, idx int) error {
	if err := p.initPage(); err != nil {
		return err
	}

	if len(p.pageOutputs) == 1 {
		idx = 0
	}

	p.pageOutputIdx = idx
	p.pageOutput = p.pageOutputs[idx]
	if p.pageOutput == nil {
		panic(fmt.Sprintf("pageOutput is nil for output idx %d", idx))
	}

	// Reset any built paginator. This will trigger when re-rendering pages in
	// server mode.
	if isRenderingSite && p.pageOutput.paginator != nil && p.pageOutput.paginator.current != nil {
		p.pageOutput.paginator.reset()
	}

	if isRenderingSite {
		cp := p.pageOutput.pco
		if cp == nil && p.reusePageOutputContent() {
			// Look for content to reuse.
			for i := 0; i < len(p.pageOutputs); i++ {
				if i == idx {
					continue
				}
				po := p.pageOutputs[i]

				if po.pco != nil {
					cp = po.pco
					break
				}
			}
		}

		if cp == nil {
			var err error
			cp, err = newPageContentOutput(p.pageOutput)
			if err != nil {
				return err
			}
		}
		p.pageOutput.setContentProvider(cp)
	} else {
		// We attempt to assign pageContentOutputs while preparing each site
		// for rendering and before rendering each site. This lets us share
		// content between page outputs to conserve resources. But if a template
		// unexpectedly calls a method of a ContentProvider that is not yet
		// initialized, we assign a LazyContentProvider that performs the
		// initialization just in time.
		if lcp, ok := (p.pageOutput.ContentProvider.(*page.LazyContentProvider)); ok {
			lcp.Reset()
		} else {
			lcp = page.NewLazyContentProvider(func() (page.OutputFormatContentProvider, error) {
				cp, err := newPageContentOutput(p.pageOutput)
				if err != nil {
					return nil, err
				}
				return cp, nil
			})
			p.pageOutput.contentRenderer = lcp
			p.pageOutput.ContentProvider = lcp
			p.pageOutput.MarkupProvider = lcp
			p.pageOutput.PageRenderProvider = lcp
			p.pageOutput.TableOfContentsProvider = lcp
		}
	}

	return nil
}

var (
	_ page.Page         = (*pageWithOrdinal)(nil)
	_ collections.Order = (*pageWithOrdinal)(nil)
	_ pageWrapper       = (*pageWithOrdinal)(nil)
)

type pageWithOrdinal struct {
	ordinal int
	*pageState
}

func (p pageWithOrdinal) Ordinal() int {
	return p.ordinal
}

func (p pageWithOrdinal) page() page.Page {
	return p.pageState
}

type pageWithWeight0 struct {
	weight0 int
	*pageState
}

func (p pageWithWeight0) Weight0() int {
	return p.weight0
}

func (p pageWithWeight0) page() page.Page {
	return p.pageState
}

var _ types.Unwrapper = (*pageWithWeight0)(nil)

func (p pageWithWeight0) Unwrapv() any {
	return p.pageState
}
