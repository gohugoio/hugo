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
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/markup/converter"

	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/bep/gitmap"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/source"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	_ page.Page           = (*pageState)(nil)
	_ collections.Grouper = (*pageState)(nil)
	_ collections.Slicer  = (*pageState)(nil)
)

var (
	pageTypesProvider = resource.NewResourceTypesProvider(media.OctetType, pageResourceType)
	nopPageOutput     = &pageOutput{
		pagePerOutputProviders:  nopPagePerOutput,
		ContentProvider:         page.NopPage,
		TableOfContentsProvider: page.NopPage,
	}
)

// pageContext provides contextual information about this page, for error
// logging and similar.
type pageContext interface {
	posOffset(offset int) text.Position
	wrapError(err error) error
	getContentConverter() converter.Converter
}

// wrapErr adds some context to the given error if possible.
func wrapErr(err error, ctx interface{}) error {
	if pc, ok := ctx.(pageContext); ok {
		return pc.wrapError(err)
	}
	return err
}

type pageSiteAdapter struct {
	p page.Page
	s *Site
}

func (pa pageSiteAdapter) GetPage(ref string) (page.Page, error) {
	p, err := pa.s.getPageNew(pa.p, ref)
	if p == nil {
		// The nil struct has meaning in some situations, mostly to avoid breaking
		// existing sites doing $nilpage.IsDescendant($p), which will always return
		// false.
		p = page.NilPage
	}
	return p, err
}

type pageState struct {
	// This slice will be of same length as the number of global slice of output
	// formats (for all sites).
	pageOutputs []*pageOutput

	// This will be shifted out when we start to render a new output format.
	*pageOutput

	// Common for all output formats.
	*pageCommon
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

func (p *pageState) GitInfo() *gitmap.GitInfo {
	return p.gitInfo
}

// GetTerms gets the terms defined on this page in the given taxonomy.
func (p *pageState) GetTerms(taxonomy string) page.Pages {
	taxonomy = strings.ToLower(taxonomy)
	m := p.s.pageMap
	prefix := cleanTreeKey(taxonomy)

	var self string
	if p.IsHome() {
		// TODO(bep) make this less magical, see taxonomyEntries.Insert.
		self = "/" + page.KindHome
	} else if p.treeRef != nil {
		self = p.treeRef.key
	}

	var pas page.Pages

	m.taxonomies.WalkQuery(pageMapQuery{Prefix: prefix}, func(s string, n *contentNode) bool {
		if _, found := m.taxonomyEntries.Get(s + self); found {
			pas = append(pas, n.p)
		}
		return false
	})

	page.SortByDefault(pas)

	return pas
}

func (p *pageState) MarshalJSON() ([]byte, error) {
	return page.MarshalPageToJSON(p)
}

func (p *pageState) getPages() page.Pages {
	b := p.bucket
	if b == nil {
		return nil
	}
	return b.getPages()
}

func (p *pageState) getPagesRecursive() page.Pages {
	b := p.bucket
	if b == nil {
		return nil
	}
	return b.getPagesRecursive()
}

func (p *pageState) getPagesAndSections() page.Pages {
	b := p.bucket
	if b == nil {
		return nil
	}
	return b.getPagesAndSections()
}

func (p *pageState) RegularPagesRecursive() page.Pages {
	p.regularPagesRecursiveInit.Do(func() {
		var pages page.Pages
		switch p.Kind() {
		case page.KindSection:
			pages = p.getPagesRecursive()
		default:
			pages = p.RegularPages()
		}
		p.regularPagesRecursive = pages
	})
	return p.regularPagesRecursive
}

func (p *pageState) PagesRecursive() page.Pages {
	return nil
}

func (p *pageState) RegularPages() page.Pages {
	p.regularPagesInit.Do(func() {
		var pages page.Pages

		switch p.Kind() {
		case page.KindPage:
		case page.KindSection, page.KindHome, page.KindTaxonomyTerm:
			pages = p.getPages()
		case page.KindTaxonomy:
			all := p.Pages()
			for _, p := range all {
				if p.IsPage() {
					pages = append(pages, p)
				}
			}
		default:
			pages = p.s.RegularPages()
		}

		p.regularPages = pages

	})

	return p.regularPages
}

func (p *pageState) Pages() page.Pages {
	p.pagesInit.Do(func() {
		var pages page.Pages

		switch p.Kind() {
		case page.KindPage:
		case page.KindSection, page.KindHome:
			pages = p.getPagesAndSections()
		case page.KindTaxonomy:
			pages = p.bucket.getTaxonomyEntries()
		case page.KindTaxonomyTerm:
			pages = p.bucket.getTaxonomies()
		default:
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

func (p *pageState) sortResources() {
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

		// Make sure not to use RelPermalink or any of the other methods that
		// trigger lazy publishing.
		return ri.Name() < rj.Name()
	})
}

func (p *pageState) Resources() resource.Resources {
	p.resourcesInit.Do(func() {
		p.sortResources()
		if len(p.m.resourcesMetadata) > 0 {
			resources.AssignMetadata(p.m.resourcesMetadata, p.resources...)
			p.sortResources()
		}
	})
	return p.resources
}

func (p *pageState) HasShortcode(name string) bool {
	if p.shortcodeState == nil {
		return false
	}

	return p.shortcodeState.nameSet[name]
}

func (p *pageState) Site() page.Site {
	return p.s.Info
}

func (p *pageState) String() string {
	if sourceRef := p.sourceRef(); sourceRef != "" {
		return fmt.Sprintf("Page(%s)", sourceRef)
	}
	return fmt.Sprintf("Page(%q)", p.Title())
}

// IsTranslated returns whether this content file is translated to
// other language(s).
func (p *pageState) IsTranslated() bool {
	p.s.h.init.translations.Do()
	return len(p.translations) > 0
}

// TranslationKey returns the key used to map language translations of this page.
// It will use the translationKey set in front matter if set, or the content path and
// filename (excluding any language code and extension), e.g. "about/index".
// The Page Kind is always prepended.
func (p *pageState) TranslationKey() string {
	p.translationKeyInit.Do(func() {
		if p.m.translationKey != "" {
			p.translationKey = p.Kind() + "/" + p.m.translationKey
		} else if p.IsPage() && !p.File().IsZero() {
			p.translationKey = path.Join(p.Kind(), filepath.ToSlash(p.File().Dir()), p.File().TranslationBaseName())
		} else if p.IsNode() {
			p.translationKey = path.Join(p.Kind(), p.SectionsPath())
		}

	})

	return p.translationKey

}

// AllTranslations returns all translations, including the current Page.
func (p *pageState) AllTranslations() page.Pages {
	p.s.h.init.translations.Do()
	return p.allTranslations
}

// Translations returns the translations excluding the current Page.
func (p *pageState) Translations() page.Pages {
	p.s.h.init.translations.Do()
	return p.translations
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
	ps.SitesProvider = ps.s.Info

	return nil
}

func (p *pageState) createRenderHooks(f output.Format) (*hooks.Render, error) {

	layoutDescriptor := p.getLayoutDescriptor()
	layoutDescriptor.RenderingHook = true
	layoutDescriptor.LayoutOverride = false
	layoutDescriptor.Layout = ""

	layoutDescriptor.Kind = "render-link"
	linkTempl, linkTemplFound, err := p.s.Tmpl().LookupLayout(layoutDescriptor, f)
	if err != nil {
		return nil, err
	}

	layoutDescriptor.Kind = "render-image"
	imgTempl, imgTemplFound, err := p.s.Tmpl().LookupLayout(layoutDescriptor, f)
	if err != nil {
		return nil, err
	}

	var linkRenderer hooks.LinkRenderer
	var imageRenderer hooks.LinkRenderer

	if linkTemplFound {
		linkRenderer = contentLinkRenderer{
			templateHandler: p.s.Tmpl(),
			Provider:        linkTempl.(tpl.Info),
			templ:           linkTempl,
		}
	}

	if imgTemplFound {
		imageRenderer = contentLinkRenderer{
			templateHandler: p.s.Tmpl(),
			Provider:        imgTempl.(tpl.Info),
			templ:           imgTempl,
		}
	}

	return &hooks.Render{
		LinkRenderer:  linkRenderer,
		ImageRenderer: imageRenderer,
	}, nil
}

func (p *pageState) getLayoutDescriptor() output.LayoutDescriptor {
	p.layoutDescriptorInit.Do(func() {
		var section string
		sections := p.SectionsEntries()

		switch p.Kind() {
		case page.KindSection:
			if len(sections) > 0 {
				section = sections[0]
			}
		case page.KindTaxonomyTerm, page.KindTaxonomy:
			b := p.getTreeRef().n
			section = b.viewInfo.name.singular
		default:
		}

		p.layoutDescriptor = output.LayoutDescriptor{
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

	if len(layouts) == 0 {
		selfLayout := p.selfLayoutForOutput(f)
		if selfLayout != "" {
			templ, found := p.s.Tmpl().Lookup(selfLayout)
			return templ, found, nil
		}
	}

	d := p.getLayoutDescriptor()

	if len(layouts) > 0 {
		d.Layout = layouts[0]
		d.LayoutOverride = true
	}

	return p.s.Tmpl().LookupLayout(d, f)
}

// This is serialized
func (p *pageState) initOutputFormat(isRenderingSite bool, idx int) error {
	if err := p.shiftToOutputFormat(isRenderingSite, idx); err != nil {
		return err
	}

	return nil

}

// Must be run after the site section tree etc. is built and ready.
func (p *pageState) initPage() error {
	if _, err := p.init.Do(); err != nil {
		return err
	}
	return nil
}

func (p *pageState) renderResources() (err error) {
	p.resourcesPublishInit.Do(func() {
		var toBeDeleted []int

		for i, r := range p.Resources() {

			if _, ok := r.(page.Page); ok {
				// Pages gets rendered with the owning page but we count them here.
				p.s.PathSpec.ProcessingStats.Incr(&p.s.PathSpec.ProcessingStats.Pages)
				continue
			}

			src, ok := r.(resource.Source)
			if !ok {
				err = errors.Errorf("Resource %T does not support resource.Source", src)
				return
			}

			if err := src.Publish(); err != nil {
				if os.IsNotExist(err) {
					// The resource has been deleted from the file system.
					// This should be extremely rare, but can happen on live reload in server
					// mode when the same resource is member of different page bundles.
					toBeDeleted = append(toBeDeleted, i)
				} else {
					p.s.Log.ERROR.Printf("Failed to publish Resource for page %q: %s", p.pathOrTitle(), err)
				}
			} else {
				p.s.PathSpec.ProcessingStats.Incr(&p.s.PathSpec.ProcessingStats.Files)
			}
		}

		for _, i := range toBeDeleted {
			p.deleteResource(i)
		}

	})

	return
}

func (p *pageState) deleteResource(i int) {
	p.resources = append(p.resources[:i], p.resources[i+1:]...)
}

func (p *pageState) getTargetPaths() page.TargetPaths {
	return p.targetPaths()
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

var defualtRenderStringOpts = renderStringOpts{
	Display: "inline",
	Markup:  "", // Will inherit the page's value when not set.
}

func (p *pageState) RenderString(args ...interface{}) (template.HTML, error) {
	if len(args) < 1 || len(args) > 2 {
		return "", errors.New("want 1 or 2 arguments")
	}

	var s string
	opts := defualtRenderStringOpts
	sidx := 1

	if len(args) == 1 {
		sidx = 0
	} else {
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return "", errors.New("first argument must be a map")
		}

		if err := mapstructure.WeakDecode(m, &opts); err != nil {
			return "", errors.WithMessage(err, "failed to decode options")
		}
	}

	var err error
	s, err = cast.ToStringE(args[sidx])
	if err != nil {
		return "", err
	}

	conv := p.getContentConverter()
	if opts.Markup != "" && opts.Markup != p.m.markup {
		var err error
		// TODO(bep) consider cache
		conv, err = p.m.newContentConverter(p, opts.Markup, nil)
		if err != nil {
			return "", p.wrapError(err)
		}
	}

	c, err := p.pageOutput.cp.renderContentWithConverter(conv, []byte(s), false)
	if err != nil {
		return "", p.wrapError(err)
	}

	b := c.Bytes()

	if opts.Display == "inline" {
		// We may have to rethink this in the future when we get other
		// renderers.
		b = p.s.ContentSpec.TrimShortHTML(b)
	}

	return template.HTML(string(b)), nil
}

func (p *pageState) addDependency(dep identity.Provider) {
	if !p.s.running() || p.pageOutput.cp == nil {
		return
	}
	p.pageOutput.cp.dependencyTracker.Add(dep)
}

func (p *pageState) RenderWithTemplateInfo(info tpl.Info, layout ...string) (template.HTML, error) {
	p.addDependency(info)
	return p.Render(layout...)
}

func (p *pageState) Render(layout ...string) (template.HTML, error) {
	templ, found, err := p.resolveTemplate(layout...)
	if err != nil {
		return "", p.wrapError(err)
	}

	if !found {
		return "", nil
	}

	p.addDependency(templ.(tpl.Info))
	res, err := executeToString(p.s.Tmpl(), templ, p)
	if err != nil {
		return "", p.wrapError(errors.Wrapf(err, "failed to execute template %q v", layout))
	}
	return template.HTML(res), nil

}

// wrapError adds some more context to the given error if possible/needed
func (p *pageState) wrapError(err error) error {
	if _, ok := err.(*herrors.ErrorWithFileContext); ok {
		// Preserve the first file context.
		return err
	}
	var filename string
	if !p.File().IsZero() {
		filename = p.File().Filename()
	}

	err, _ = herrors.WithFileContextForFile(
		err,
		filename,
		filename,
		p.s.SourceSpec.Fs.Source,
		herrors.SimpleLineMatcher)

	return err
}

func (p *pageState) getContentConverter() converter.Converter {
	var err error
	p.m.contentConverterInit.Do(func() {
		markup := p.m.markup
		if markup == "html" {
			// Only used for shortcode inner content.
			markup = "markdown"
		}
		p.m.contentConverter, err = p.m.newContentConverter(p, markup, p.m.renderingConfigOverrides)

	})

	if err != nil {
		p.s.Log.ERROR.Println("Failed to create content converter:", err)
	}
	return p.m.contentConverter
}

func (p *pageState) mapContent(bucket *pagesMapBucket, meta *pageMeta) error {

	s := p.shortcodeState

	rn := &pageContentMap{
		items: make([]interface{}, 0, 20),
	}

	iter := p.source.parsed.Iterator()

	fail := func(err error, i pageparser.Item) error {
		return p.parseError(err, iter.Input(), i.Pos)
	}

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currShortcode shortcode
	var ordinal int
	var frontMatterSet bool

Loop:
	for {
		it := iter.Next()

		switch {
		case it.Type == pageparser.TypeIgnore:
		case it.IsFrontMatter():
			f := pageparser.FormatFromFrontMatterType(it.Type)
			m, err := metadecoders.Default.UnmarshalToMap(it.Val, f)
			if err != nil {
				if fe, ok := err.(herrors.FileError); ok {
					return herrors.ToFileErrorWithOffset(fe, iter.LineNumber()-1)
				} else {
					return err
				}
			}

			if err := meta.setMetadata(bucket, p, m); err != nil {
				return err
			}

			frontMatterSet = true

			next := iter.Peek()
			if !next.IsDone() {
				p.source.posMainContent = next.Pos
			}

			if !p.s.shouldBuild(p) {
				// Nothing more to do.
				return nil
			}

		case it.Type == pageparser.TypeLeadSummaryDivider:
			posBody := -1
			f := func(item pageparser.Item) bool {
				if posBody == -1 && !item.IsDone() {
					posBody = item.Pos
				}

				if item.IsNonWhitespace() {
					p.truncated = true

					// Done
					return false
				}
				return true
			}
			iter.PeekWalk(f)

			p.source.posSummaryEnd = it.Pos
			p.source.posBodyStart = posBody
			p.source.hasSummaryDivider = true

			if meta.markup != "html" {
				// The content will be rendered by Blackfriday or similar,
				// and we need to track the summary.
				rn.AddReplacement(internalSummaryDividerPre, it)
			}

		// Handle shortcode
		case it.IsLeftShortcodeDelim():
			// let extractShortcode handle left delim (will do so recursively)
			iter.Backup()

			currShortcode, err := s.extractShortcode(ordinal, 0, iter)
			if err != nil {
				return fail(errors.Wrap(err, "failed to extract shortcode"), it)
			}

			currShortcode.pos = it.Pos
			currShortcode.length = iter.Current().Pos - it.Pos
			if currShortcode.placeholder == "" {
				currShortcode.placeholder = createShortcodePlaceholder("s", currShortcode.ordinal)
			}

			if currShortcode.name != "" {
				s.nameSet[currShortcode.name] = true
			}

			if currShortcode.params == nil {
				var s []string
				currShortcode.params = s
			}

			currShortcode.placeholder = createShortcodePlaceholder("s", ordinal)
			ordinal++
			s.shortcodes = append(s.shortcodes, currShortcode)

			rn.AddShortcode(currShortcode)

		case it.Type == pageparser.TypeEmoji:
			if emoji := helpers.Emoji(it.ValStr()); emoji != nil {
				rn.AddReplacement(emoji, it)
			} else {
				rn.AddBytes(it)
			}
		case it.IsEOF():
			break Loop
		case it.IsError():
			err := fail(errors.WithStack(errors.New(it.ValStr())), it)
			currShortcode.err = err
			return err

		default:
			rn.AddBytes(it)
		}
	}

	if !frontMatterSet {
		// Page content without front matter. Assign default front matter from
		// cascades etc.
		if err := meta.setMetadata(bucket, p, nil); err != nil {
			return err
		}
	}

	p.cmap = rn

	return nil
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

func (p *pageState) outputFormat() (f output.Format) {
	if p.pageOutput == nil {
		panic("no pageOutput")
	}
	return p.pageOutput.f
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
	if !p.File().IsZero() {
		return p.File().Filename()
	}

	if p.Path() != "" {
		return p.Path()
	}

	return p.Title()
}

func (p *pageState) posFromPage(offset int) text.Position {
	return p.posFromInput(p.source.parsed.Input(), offset)
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

func (p *pageState) posOffset(offset int) text.Position {
	return p.posFromInput(p.source.parsed.Input(), offset)
}

// shiftToOutputFormat is serialized. The output format idx refers to the
// full set of output formats for all sites.
func (p *pageState) shiftToOutputFormat(isRenderingSite bool, idx int) error {
	if err := p.initPage(); err != nil {
		return err
	}

	if len(p.pageOutputs) == 1 {
		idx = 0
	}

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
		cp := p.pageOutput.cp
		if cp == nil {

			// Look for content to reuse.
			for i := 0; i < len(p.pageOutputs); i++ {
				if i == idx {
					continue
				}
				po := p.pageOutputs[i]

				if po.cp != nil && po.cp.reuse {
					cp = po.cp
					break
				}
			}
		}

		if cp == nil {
			var err error
			cp, err = newPageContentOutput(p, p.pageOutput)
			if err != nil {
				return err
			}
		}
		p.pageOutput.initContentProvider(cp)
		p.pageOutput.cp = cp
	}

	return nil
}

// sourceRef returns the reference used by GetPage and ref/relref shortcodes to refer to
// this page. It is prefixed with a "/".
//
// For pages that have a source file, it is returns the path to this file as an
// absolute path rooted in this site's content dir.
// For pages that do not (sections witout content page etc.), it returns the
// virtual path, consistent with where you would add a source file.
func (p *pageState) sourceRef() string {
	if !p.File().IsZero() {
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

func (s *Site) sectionsFromFile(fi source.File) []string {
	dirname := fi.Dir()

	dirname = strings.Trim(dirname, helpers.FilePathSeparator)
	if dirname == "" {
		return nil
	}
	parts := strings.Split(dirname, helpers.FilePathSeparator)

	if fii, ok := fi.(*fileInfo); ok {
		if len(parts) > 0 && fii.FileInfo().Meta().Classifier() == files.ContentClassLeaf {
			// my-section/mybundle/index.md => my-section
			return parts[:len(parts)-1]
		}
	}

	return parts
}
