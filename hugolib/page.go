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
	"iter"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/segments"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/tpl/tplimpl"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/tableofcontents"

	"github.com/gohugoio/hugo/common/herrors"
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

	s *Site

	// This slice will be of same length as the number of global slice of output
	// formats (for all sites).
	pageOutputs []*pageOutput

	// Used to determine if we can reuse content across output formats.
	pageOutputTemplateVariationsState atomic.Uint32

	// This will be shifted out when we start to render a new output format.
	pageOutputIdx int
	*pageOutput

	// Common for all output formats.
	*pageCommon

	resource.Staler
	dependencyManager identity.Manager
}

// This is not accurate and only used for progress reporting.
// We can do better, but this will do for now.
func (p *pageState) hasRenderableOutput() bool {
	for _, po := range p.pageOutputs {
		if po.render {
			return true
		}
	}
	return false
}

func (p *pageState) incrPageOutputTemplateVariation() {
	p.pageOutputTemplateVariationsState.Add(1)
}

func (ps *pageState) canReusePageOutputContent() bool {
	return ps.pageOutputTemplateVariationsState.Load() == 1
}

func (ps *pageState) IdentifierBase() string {
	return ps.Path()
}

func (ps *pageState) GetIdentity() identity.Identity {
	return ps
}

func (ps *pageState) ForEeachIdentity(f func(identity.Identity) bool) bool {
	return f(ps)
}

func (ps *pageState) GetDependencyManager() identity.Manager {
	return ps.dependencyManager
}

func (ps *pageState) GetDependencyManagerForScope(scope int) identity.Manager {
	switch scope {
	case pageDependencyScopeDefault:
		return ps.dependencyManagerOutput
	case pageDependencyScopeGlobal:
		return ps.dependencyManager
	default:
		return identity.NopManager
	}
}

func (ps *pageState) GetDependencyManagerForScopesAll() []identity.Manager {
	return []identity.Manager{ps.dependencyManager, ps.dependencyManagerOutput}
}

// Param is a convenience method to do lookups in Page's and Site's Params map,
// in that order.
//
// This method is also implemented on SiteInfo.
func (ps *pageState) Param(key any) (any, error) {
	return resource.Param(ps, ps.s.Params(), key)
}

func (ps *pageState) Key() string {
	return "page-" + strconv.FormatUint(ps.pid, 10)
}

// RelatedKeywords implements the related.Document interface needed for fast page searches.
func (ps *pageState) RelatedKeywords(cfg related.IndexConfig) ([]related.Keyword, error) {
	v, found, err := page.NamedPageMetaValue(ps, cfg.Name)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return cfg.ToKeywords(v)
}

func (ps *pageState) resetBuildState() {
	ps.m.prepareRebuild()
}

func (ps *pageState) skipRender() bool {
	b := ps.s.conf.Segments.Config.SegmentFilter.ShouldExcludeFine(
		segments.SegmentQuery{
			Path:   ps.Path(),
			Kind:   ps.Kind(),
			Site:   ps.s.siteVector,
			Output: ps.pageOutput.f.Name,
		},
	)

	return b
}

func (ps *pageState) isRenderedAny() bool {
	for _, o := range ps.pageOutputs {
		if o.isRendered() {
			return true
		}
	}
	return false
}

// Implements contentNode.

func (ps *pageState) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	return f(ps.s.siteVector, ps)
}

func (ps *pageState) contentWeight() int {
	if ps.m.f == nil {
		return 0
	}
	return ps.m.f.FileInfo().Meta().Weight
}

func (ps *pageState) nodeSourceEntryID() any {
	return ps.m.nodeSourceEntryID()
}

func (ps *pageState) nodeCategoryPage() {
	// Marker method.
}

func (m *pageState) nodeCategorySingle() {
	// Marker method.
}

func (ps *pageState) lookupContentNode(v sitesmatrix.Vector) contentNode {
	pc := ps.m.pageConfigSource
	if pc.MatchSiteVector(v) {
		return ps
	}
	return nil
}

func (ps *pageState) lookupContentNodes(vec sitesmatrix.Vector, fallback bool) iter.Seq[contentNodeForSite] {
	nop := func(yield func(n contentNodeForSite) bool) {}
	pc := ps.m.pageConfigSource
	if !fallback {
		if !pc.MatchSiteVector(vec) {
			return nop
		}
		return func(yield func(n contentNodeForSite) bool) {
			yield(ps)
		}
	}

	if !pc.MatchLanguageCoarse(vec) {
		return nop
	}

	if !pc.MatchVersionCoarse(vec) {
		return nop
	}
	if !pc.MatchRoleCoarse(vec) {
		return nop
	}

	return func(yield func(n contentNodeForSite) bool) {
		yield(ps)
	}
}

// Eq returns whether the current page equals the given page.
// This is what's invoked when doing `{{ if eq $page $otherPage }}`
func (ps *pageState) Eq(other any) bool {
	pp, err := unwrapPage(other)
	if err != nil {
		return false
	}

	return ps == pp
}

func (ps *pageState) HeadingsFiltered(context.Context) tableofcontents.Headings {
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
func (ps *pageState) ApplyFilterToHeadings(ctx context.Context, fn func(*tableofcontents.Heading) bool) related.Document {
	fragments := ps.pageOutput.pco.c().Fragments(ctx)
	headings := fragments.Headings.FilterBy(fn)
	return &pageHeadingsFiltered{
		pageState: ps,
		headings:  headings,
	}
}

func (ps *pageState) GitInfo() *source.GitInfo {
	return ps.gitInfo
}

func (ps *pageState) CodeOwners() []string {
	return ps.codeowners
}

// GetTerms gets the terms defined on this page in the given taxonomy.
// The pages returned will be ordered according to the front matter.
func (ps *pageState) GetTerms(taxonomy string) page.Pages {
	return ps.s.pageMap.getTermsForPageInTaxonomy(ps.Path(), taxonomy)
}

func (ps *pageState) MarshalJSON() ([]byte, error) {
	return page.MarshalPageToJSON(ps)
}

func (ps *pageState) RegularPagesRecursive() page.Pages {
	switch ps.Kind() {
	case kinds.KindSection, kinds.KindHome:
		return ps.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    ps.Path(),
					Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindPage).BoolFunc(),
				},
				Recursive: true,
			},
		)
	default:
		return ps.RegularPages()
	}
}

func (ps *pageState) PagesRecursive() page.Pages {
	return nil
}

func (ps *pageState) RegularPages() page.Pages {
	switch ps.Kind() {
	case kinds.KindPage:
	case kinds.KindSection, kinds.KindHome, kinds.KindTaxonomy:
		return ps.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    ps.Path(),
					Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindPage).BoolFunc(),
				},
			},
		)
	case kinds.KindTerm:
		return ps.s.pageMap.getPagesWithTerm(
			pageMapQueryPagesBelowPath{
				Path:    ps.Path(),
				Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindPage).BoolFunc(),
			},
		)
	default:
		return ps.s.RegularPages()
	}
	return nil
}

func (ps *pageState) Pages() page.Pages {
	switch ps.Kind() {
	case kinds.KindPage:
	case kinds.KindSection, kinds.KindHome:
		return ps.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    ps.Path(),
					KeyPart: "page-section",
					Include: pagePredicates.ShouldListLocal.And(
						pagePredicates.KindPage.Or(pagePredicates.KindSection),
					).BoolFunc(),
				},
			},
		)
	case kinds.KindTerm:
		return ps.s.pageMap.getPagesWithTerm(
			pageMapQueryPagesBelowPath{
				Path: ps.Path(),
			},
		)
	case kinds.KindTaxonomy:
		return ps.s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    ps.Path(),
					KeyPart: "term",
					Include: pagePredicates.ShouldListLocal.And(pagePredicates.KindTerm).BoolFunc(),
				},
				Recursive: true,
			},
		)
	default:
		return ps.s.Pages()
	}
	return nil
}

// RawContent returns the un-rendered source content without
// any leading front matter.
func (ps *pageState) RawContent() string {
	if ps.m.content.pi.itemsStep2 == nil {
		return ""
	}
	start := ps.m.content.pi.posMainContent
	if start == -1 {
		start = 0
	}
	source, err := ps.m.content.pi.contentSource(ps.m.content)
	if err != nil {
		panic(err)
	}
	return string(source[start:])
}

func (ps *pageState) Resources() resource.Resources {
	return ps.s.pageMap.getOrCreateResourcesForPage(ps)
}

func (ps *pageState) HasShortcode(name string) bool {
	p := ps.m.content.hasShortcode.Load()
	if p == nil {
		return false
	}
	hasShortcode := *p
	return hasShortcode(name)
}

func (ps *pageState) Site() page.Site {
	return ps.s.siteWrapped
}

func (ps *pageState) String() string {
	var sb strings.Builder
	if ps.File() != nil {
		// The forward slashes even on Windows is motivated by
		// getting stable tests.
		// This information is meant for getting positional information in logs,
		// so the direction of the slashes should not matter.
		sb.WriteString(filepath.ToSlash(ps.File().Filename()))
		if ps.File().IsContentAdapter() {
			// Also include the path.
			sb.WriteString(":")
			sb.WriteString(ps.Path())
		}
	} else {
		sb.WriteString(ps.Path())
	}
	return sb.String()
}

// IsTranslated returns whether this content file is translated to
// other language(s).
func (ps *pageState) IsTranslated() bool {
	return len(ps.Translations()) > 0
}

// TranslationKey returns the key used to identify a translation of this content.
func (ps *pageState) TranslationKey() string {
	if ps.m.pageConfig.TranslationKey != "" {
		return ps.m.pageConfig.TranslationKey
	}
	return ps.Path()
}

// AllTranslations returns all translations, including the current Page.
func (ps *pageState) AllTranslations() page.Pages {
	key := ps.Path() + "/" + "translations-all"
	// This is called from Translations, so we need to use a different partition, cachePages2,
	// to avoid potential deadlocks.
	pages, err := ps.s.pageMap.getOrCreatePagesFromCache(ps.s.pageMap.cachePages2, key, func(string) (page.Pages, error) {
		if ps.m.pageConfig.TranslationKey != "" {
			// translationKey set by user.
			pas, _ := ps.s.h.translationKeyPages.Get(ps.m.pageConfig.TranslationKey)
			pasc := slices.Clone(pas)
			page.SortByLanguage(pasc)
			return pasc, nil
		}
		var pas page.Pages

		ps.s.pageMap.treePages.ForEeachInDimension(ps.Path(), ps.s.siteVector, sitesmatrix.Language,
			func(n contentNode) bool {
				if n != nil {
					pas = append(pas, n.(page.Page))
				}
				return true
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

// For internal use only.
func (ps *pageState) SiteVector() sitesmatrix.Vector {
	return ps.s.siteVector
}

func (ps *pageState) siteVector() sitesmatrix.Vector {
	return ps.s.siteVector
}

func (ps *pageState) siteVectors() sitesmatrix.VectorIterator {
	return ps.s.siteVector
}

// Rotate returns all pages in the given dimension for this page.
func (ps *pageState) Rotate(dimensionStr string) (page.Pages, error) {
	dimensionStr = strings.ToLower(dimensionStr)
	key := ps.Path() + "/" + "rotate-" + dimensionStr
	d, err := sitesmatrix.ParseDimension(dimensionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dimension %q: %w", dimensionStr, err)
	}

	pages, err := ps.s.pageMap.getOrCreatePagesFromCache(ps.s.pageMap.cachePages2, key, func(string) (page.Pages, error) {
		var pas page.Pages
		ps.s.pageMap.treePages.ForEeachInDimension(ps.Path(), ps.s.siteVector, d,
			func(n contentNode) bool {
				if n != nil {
					p := n.(page.Page)
					pas = append(pas, p)
				}
				return true
			},
		)

		if dimensionStr == "language" && ps.m.pageConfig.TranslationKey != "" {
			// translationKey set by user.
			// This is an old construct back from when languages were introduced.
			// We keep it for backward compatibility.
			// Also see AllTranslations.
			pas, _ := ps.s.h.translationKeyPages.Get(ps.m.pageConfig.TranslationKey)
			pasc := slices.Clone(pas)
			page.SortByLanguage(pasc)
			return pasc, nil
		}

		pas = pagePredicates.ShouldLink.Filter(pas)
		page.SortByDims(pas)
		return pas, nil
	})

	return pages, err
}

// Translations returns the translations excluding the current Page.
func (ps *pageState) Translations() page.Pages {
	key := ps.Path() + "/" + "translations"
	pages, err := ps.s.pageMap.getOrCreatePagesFromCache(nil, key, func(string) (page.Pages, error) {
		var pas page.Pages
		for _, pp := range ps.AllTranslations() {
			if !pp.Eq(ps) {
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
		if ps.s == nil {
			panic("no site")
		}
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

// Exported so it can be used in integration tests.
func (po *pageOutput) GetInternalTemplateBasePathAndDescriptor() (string, tplimpl.TemplateDescriptor) {
	p := po.p
	f := po.f

	base := p.PathInfo().BaseReTyped(p.m.pageConfig.Type)
	return base, tplimpl.TemplateDescriptor{
		Kind:           p.Kind(),
		LayoutFromUser: p.Layout(),
		OutputFormat:   f.Name,
		MediaType:      f.MediaType.Type,
		IsPlainText:    f.IsPlainText,
	}
}

func (ps *pageState) resolveTemplate(layouts ...string) (*tplimpl.TemplInfo, bool, error) {
	dir, d := ps.GetInternalTemplateBasePathAndDescriptor()

	if len(layouts) > 0 {
		d.LayoutFromUser = layouts[0]
		d.LayoutFromUserMustMatch = true
	}

	q := tplimpl.TemplateQuery{
		Path:     dir,
		Category: tplimpl.CategoryLayout,
		Sites:    ps.s.siteVector,
		Desc:     d,
	}

	tinfo := ps.s.TemplateStore.LookupPagesLayout(q)
	if tinfo == nil {
		return nil, false, nil
	}

	return tinfo, true, nil
}

// Must be run after the site section tree etc. is built and ready.
func (ps *pageState) initPage() error {
	var initErr error
	ps.pageInit.Do(func() {
		var pp pagePaths
		pp, initErr = newPagePaths(ps)
		if initErr != nil {
			return
		}

		var outputFormatsForPage output.Formats
		var renderFormats output.Formats

		if ps.m.standaloneOutputFormat.IsZero() {
			outputFormatsForPage = ps.outputFormats()
			renderFormats = ps.s.h.renderFormats
		} else {
			// One of the fixed output format pages, e.g. 404.
			outputFormatsForPage = output.Formats{ps.m.standaloneOutputFormat}
			renderFormats = outputFormatsForPage
		}

		// Prepare output formats for all sites.
		// We do this even if this page does not get rendered on
		// its own. It may be referenced via one of the site collections etc.
		// it will then need an output format.
		ps.pageOutputs = make([]*pageOutput, len(renderFormats))
		created := make(map[string]*pageOutput)
		shouldRenderPage := !ps.m.noRender()

		for i, f := range renderFormats {

			if po, found := created[f.Name]; found {
				ps.pageOutputs[i] = po
				continue
			}

			render := shouldRenderPage
			if render {
				_, render = outputFormatsForPage.GetByName(f.Name)
			}

			po := newPageOutput(ps, pp, f, render)

			// Create a content provider for the first,
			// we may be able to reuse it.
			if i == 0 {
				var contentProvider *pageContentOutput
				contentProvider, initErr = newPageContentOutput(po)
				if initErr != nil {
					return
				}
				po.setContentProvider(contentProvider)
			}

			ps.pageOutputs[i] = po
			created[f.Name] = po

		}

		if initErr = ps.initCommonProviders(pp); initErr != nil {
			return
		}
	})

	return initErr
}

func (ps *pageState) renderResources() error {
	for _, r := range ps.Resources() {
		if _, ok := r.(page.Page); ok {
			if ps.s.h.buildCounter.Load() == 0 {
				// Pages gets rendered with the owning page but we count them here.
				ps.s.PathSpec.ProcessingStats.Incr(&ps.s.PathSpec.ProcessingStats.Pages)
			}
			continue
		}

		if resources.IsPublished(r) {
			continue
		}

		src, ok := r.(resource.Source)
		if !ok {
			return fmt.Errorf("resource %T does not support resource.Source", r)
		}

		if err := src.Publish(); err != nil {
			if !herrors.IsNotExist(err) {
				ps.s.Log.Errorf("Failed to publish Resource for page %q: %s", ps.pathOrTitle(), err)
			}
		} else {
			ps.s.PathSpec.ProcessingStats.Incr(&ps.s.PathSpec.ProcessingStats.Files)
		}
	}

	return nil
}

func (ps *pageState) AlternativeOutputFormats() page.OutputFormats {
	f := ps.outputFormat()
	var o page.OutputFormats
	for _, of := range ps.OutputFormats() {
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

func (m *pageMetaSource) wrapError(err error, sourceFs afero.Fs) error {
	if err == nil {
		panic("wrapError with nil")
	}

	if m.f == nil {
		// No more details to add.
		return fmt.Errorf("%q: %w", m.Path(), err)
	}

	return hugofs.AddFileInfoToError(err, m.f.FileInfo(), sourceFs)
}

// wrapError adds some more context to the given error if possible/needed
func (ps *pageState) wrapError(err error) error {
	return ps.m.wrapError(err, ps.s.h.SourceFs)
}

func (ps *pageState) getPageInfoForError() string {
	s := fmt.Sprintf("kind: %q, path: %q", ps.Kind(), ps.Path())
	if ps.File() != nil {
		s += fmt.Sprintf(", file: %q", ps.File().Filename())
	}
	return s
}

func (ps *pageState) getContentConverter() converter.Converter {
	var err error
	ps.contentConverterInit.Do(func() {
		if ps.m.pageConfigSource.ContentMediaType.IsZero() {
			panic("ContentMediaType not set")
		}
		markup := ps.m.pageConfigSource.ContentMediaType.SubType

		if markup == "html" {
			// Only used for shortcode inner content.
			markup = "markdown"
		}
		ps.contentConverter, err = ps.m.newContentConverter(ps, markup)
	})

	if err != nil {
		ps.s.Log.Errorln("Failed to create content converter:", err)
	}
	return ps.contentConverter
}

func (ps *pageState) errorf(err error, format string, a ...any) error {
	if herrors.UnwrapFileError(err) != nil {
		// More isn't always better.
		return err
	}
	args := append([]any{ps.Language().Lang, ps.pathOrTitle()}, a...)
	args = append(args, err)
	format = "[%s] page %q: " + format + ": %w"
	if err == nil {
		return fmt.Errorf(format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (ps *pageState) outputFormat() (f output.Format) {
	if ps.pageOutput == nil {
		panic("no pageOutput")
	}
	return ps.pageOutput.f
}

func (ps *pageState) parseError(err error, input []byte, offset int) error {
	pos := posFromInput("", input, offset)
	return herrors.NewFileErrorFromName(err, ps.File().Filename()).UpdatePosition(pos)
}

func (ps *pageState) pathOrTitle() string {
	if ps.File() != nil {
		return ps.File().Filename()
	}

	if ps.Path() != "" {
		return ps.Path()
	}

	return ps.Title()
}

func (ps *pageState) posFromInput(input []byte, offset int) text.Position {
	return posFromInput(ps.pathOrTitle(), input, offset)
}

func (ps *pageState) posOffset(offset int) text.Position {
	return ps.posFromInput(ps.m.content.mustSource(), offset)
}

// shiftToOutputFormat is serialized. The output format idx refers to the
// full set of output formats for all sites.
// This is serialized.
func (ps *pageState) shiftToOutputFormat(isRenderingSite bool, idx int) error {
	if err := ps.initPage(); err != nil {
		return err
	}

	if len(ps.pageOutputs) == 1 {
		idx = 0
	}

	ps.pageOutputIdx = idx
	ps.pageOutput = ps.pageOutputs[idx]
	if ps.pageOutput == nil {
		panic(fmt.Sprintf("pageOutput is nil for output idx %d", idx))
	}

	// Reset any built paginator. This will trigger when re-rendering pages in
	// server mode.
	if isRenderingSite && ps.pageOutput.paginator != nil && ps.pageOutput.paginator.current != nil {
		ps.pageOutput.paginator.reset()
	}

	if isRenderingSite {
		cp := ps.pageOutput.pco
		if cp == nil && ps.canReusePageOutputContent() {
			// Look for content to reuse.
			for i := range ps.pageOutputs {
				if i == idx {
					continue
				}
				po := ps.pageOutputs[i]

				if po.pco != nil {
					cp = po.pco
					break
				}
			}
		}

		if cp == nil {
			var err error
			cp, err = newPageContentOutput(ps.pageOutput)
			if err != nil {
				return err
			}
		}
		ps.pageOutput.setContentProvider(cp)
	} else {
		// We attempt to assign pageContentOutputs while preparing each site
		// for rendering and before rendering each site. This lets us share
		// content between page outputs to conserve resources. But if a template
		// unexpectedly calls a method of a ContentProvider that is not yet
		// initialized, we assign a LazyContentProvider that performs the
		// initialization just in time.
		if lcp, ok := (ps.pageOutput.ContentProvider.(*page.LazyContentProvider)); ok {
			lcp.Reset()
		} else {
			lcp = page.NewLazyContentProvider(func(context.Context) page.OutputFormatContentProvider {
				cp, err := newPageContentOutput(ps.pageOutput)
				if err != nil {
					panic(err)
				}
				return cp
			})

			ps.pageOutput.contentRenderer = lcp
			ps.pageOutput.ContentProvider = lcp
			ps.pageOutput.MarkupProvider = lcp
			ps.pageOutput.PageRenderProvider = lcp
			ps.pageOutput.TableOfContentsProvider = lcp
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
