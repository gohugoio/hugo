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

	"github.com/gohugoio/hugo/common/maps"

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

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/text"
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
	nopPageOutput     = &pageOutput{pagePerOutputProviders: nopPagePerOutput}
)

// pageContext provides contextual information about this page, for error
// logging and similar.
type pageContext interface {
	posOffset(offset int) text.Position
	wrapError(err error) error
	getRenderingConfig() *helpers.BlackFriday
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

func (p *pageState) getPagesAndSections() page.Pages {
	b := p.bucket
	if b == nil {
		return nil
	}
	return b.getPagesAndSections()
}

// TODO(bep) cm add a test
func (p *pageState) RegularPages() page.Pages {
	p.regularPagesInit.Do(func() {
		var pages page.Pages

		switch p.Kind() {
		case page.KindPage:
		case page.KindHome:
			pages = p.s.RegularPages()
		case page.KindSection, page.KindTaxonomyTerm:
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
		case page.KindHome:
			// See https://github.com/gohugoio/hugo/issues/6238
			// Note: When making the change below, also remember RegularPages.
			helpers.DistinctWarnLog.Println(`In the next Hugo version (0.58.0) we will change how $home.Pages behaves. If you want to list all regular pages, replace .Pages or .Data.Pages with .Site.RegularPages in your home page template.`)
			pages = p.s.RegularPages()
		case page.KindSection:
			pages = p.getPagesAndSections()
		case page.KindTaxonomy:
			termInfo := p.bucket
			plural := maps.GetString(termInfo.meta, "plural")
			term := maps.GetString(termInfo.meta, "termKey")
			taxonomy := p.s.Taxonomies[plural].Get(term)
			pages = taxonomy.Pages()
		case page.KindTaxonomyTerm:
			pages = p.getPagesAndSections()
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

func (p *pageState) HasShortcode(name string) bool {
	if p.shortcodeState == nil {
		return false
	}

	return p.shortcodeState.nameSet[name]
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

func (p *pageState) getRenderingConfig() *helpers.BlackFriday {
	if p.m.renderingConfig == nil {
		return p.s.ContentSpec.BlackFriday
	}
	return p.m.renderingConfig
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
	ps.SitesProvider = &ps.s.Info

	return nil
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
			section = maps.GetString(p.bucket.meta, "singular")

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

func (p *pageState) getLayouts(layouts ...string) ([]string, error) {
	f := p.outputFormat()

	if len(layouts) == 0 {
		selfLayout := p.selfLayoutForOutput(f)
		if selfLayout != "" {
			return []string{selfLayout}, nil
		}
	}

	layoutDescriptor := p.getLayoutDescriptor()

	if len(layouts) > 0 {
		layoutDescriptor.Layout = layouts[0]
		layoutDescriptor.LayoutOverride = true
	}

	return p.s.layoutHandler.For(layoutDescriptor, f)
}

// This is serialized
func (p *pageState) initOutputFormat(isRenderingSite bool, idx int) error {
	if err := p.shiftToOutputFormat(isRenderingSite, idx); err != nil {
		return err
	}

	if !p.renderable {
		if _, err := p.Content(); err != nil {
			return err
		}
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

func (p *pageState) Render(layout ...string) template.HTML {
	l, err := p.getLayouts(layout...)
	if err != nil {
		p.s.SendError(p.wrapError(errors.Errorf(".Render: failed to resolve layout %v", layout)))
		return ""
	}

	for _, layout := range l {
		templ, found := p.s.Tmpl.Lookup(layout)
		if !found {
			// This is legacy from when we had only one output format and
			// HTML templates only. Some have references to layouts without suffix.
			// We default to good old HTML.
			templ, _ = p.s.Tmpl.Lookup(layout + ".html")
		}
		if templ != nil {
			res, err := executeToString(templ, p)
			if err != nil {
				p.s.SendError(p.wrapError(errors.Wrapf(err, ".Render: failed to execute template %q v", layout)))
				return ""
			}
			return template.HTML(res)
		}
	}

	return ""

}

// wrapError adds some more context to the given error if possible
func (p *pageState) wrapError(err error) error {

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

func (p *pageState) addResources(r ...resource.Resource) {
	p.resources = append(p.resources, r...)
}

func (p *pageState) mapContent(bucket *pagesMapBucket, meta *pageMeta) error {

	s := p.shortcodeState

	p.renderable = true

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

Loop:
	for {
		it := iter.Next()

		switch {
		case it.Type == pageparser.TypeIgnore:
		case it.Type == pageparser.TypeHTMLStart:
			// This is HTML without front matter. It can still have shortcodes.
			p.selfLayout = "__" + p.File().Filename()
			p.renderable = false
			rn.AddBytes(it)
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

			if err := meta.setMetadata(bucket, p, m); err != nil {
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

	if idx >= len(p.pageOutputs) {
		panic(fmt.Sprintf("invalid page state for %q: got output format index %d, have %d", p.pathOrTitle(), idx, len(p.pageOutputs)))
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

	if idx > 0 {
		// Check if we can reuse content from one of the previous formats.
		for i := idx - 1; i >= 0; i-- {
			po := p.pageOutputs[i]
			if po.cp != nil && po.cp.reuse {
				p.pageOutput.cp = po.cp
				break
			}
		}
	}

	for _, r := range p.Resources().ByType(pageResourceType) {
		rp := r.(*pageState)
		if err := rp.shiftToOutputFormat(isRenderingSite, idx); err != nil {
			return errors.Wrap(err, "failed to shift outputformat in Page resource")
		}
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
