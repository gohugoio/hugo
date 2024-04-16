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
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"golang.org/x/text/unicode/norm"

	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/markup/converter/hooks"

	"github.com/gohugoio/hugo/markup/converter"

	"github.com/gohugoio/hugo/common/text"

	"github.com/gohugoio/hugo/publisher"

	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/lazy"

	"github.com/fsnotify/fsnotify"
	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/tpl"
)

func (s *Site) Taxonomies() page.TaxonomyList {
	s.init.taxonomies.Do(context.Background())
	return s.taxonomies
}

type (
	taxonomiesConfig       map[string]string
	taxonomiesConfigValues struct {
		views          []viewName
		viewsByTreeKey map[string]viewName
	}
)

func (t taxonomiesConfig) Values() taxonomiesConfigValues {
	var views []viewName
	for k, v := range t {
		views = append(views, viewName{singular: k, plural: v, pluralTreeKey: cleanTreeKey(v)})
	}
	sort.Slice(views, func(i, j int) bool {
		return views[i].plural < views[j].plural
	})

	viewsByTreeKey := make(map[string]viewName)
	for _, v := range views {
		viewsByTreeKey[v.pluralTreeKey] = v
	}

	return taxonomiesConfigValues{
		views:          views,
		viewsByTreeKey: viewsByTreeKey,
	}
}

// Lazily loaded site dependencies.
type siteInit struct {
	prevNext          *lazy.Init
	prevNextInSection *lazy.Init
	menus             *lazy.Init
	taxonomies        *lazy.Init
}

func (init *siteInit) Reset() {
	init.prevNext.Reset()
	init.prevNextInSection.Reset()
	init.menus.Reset()
	init.taxonomies.Reset()
}

func (s *Site) prepareInits() {
	s.init = &siteInit{}

	var init lazy.Init

	s.init.prevNext = init.Branch(func(context.Context) (any, error) {
		regularPages := s.RegularPages()
		for i, p := range regularPages {
			np, ok := p.(nextPrevProvider)
			if !ok {
				continue
			}

			pos := np.getNextPrev()
			if pos == nil {
				continue
			}

			pos.nextPage = nil
			pos.prevPage = nil

			if i > 0 {
				pos.nextPage = regularPages[i-1]
			}

			if i < len(regularPages)-1 {
				pos.prevPage = regularPages[i+1]
			}
		}
		return nil, nil
	})

	s.init.prevNextInSection = init.Branch(func(context.Context) (any, error) {
		setNextPrev := func(pas page.Pages) {
			for i, p := range pas {
				np, ok := p.(nextPrevInSectionProvider)
				if !ok {
					continue
				}

				pos := np.getNextPrevInSection()
				if pos == nil {
					continue
				}

				pos.nextPage = nil
				pos.prevPage = nil

				if i > 0 {
					pos.nextPage = pas[i-1]
				}

				if i < len(pas)-1 {
					pos.prevPage = pas[i+1]
				}
			}
		}

		sections := s.pageMap.getPagesInSection(
			pageMapQueryPagesInSection{
				pageMapQueryPagesBelowPath: pageMapQueryPagesBelowPath{
					Path:    "",
					KeyPart: "sectionorhome",
					Include: pagePredicates.KindSection.Or(pagePredicates.KindHome),
				},
				IncludeSelf: true,
				Recursive:   true,
			},
		)

		for _, section := range sections {
			setNextPrev(section.RegularPages())
		}

		return nil, nil
	})

	s.init.menus = init.Branch(func(context.Context) (any, error) {
		err := s.assembleMenus()
		return nil, err
	})

	s.init.taxonomies = init.Branch(func(ctx context.Context) (any, error) {
		if err := s.pageMap.CreateSiteTaxonomies(ctx); err != nil {
			return nil, err
		}
		return s.taxonomies, nil
	})
}

type siteRenderingContext struct {
	output.Format
}

func (s *Site) Menus() navigation.Menus {
	s.init.menus.Do(context.Background())
	return s.menus
}

func (s *Site) initRenderFormats() {
	formatSet := make(map[string]bool)
	formats := output.Formats{}

	w := &doctree.NodeShiftTreeWalker[contentNodeI]{
		Tree: s.pageMap.treePages,
		Handle: func(key string, n contentNodeI, match doctree.DimensionFlag) (bool, error) {
			if p, ok := n.(*pageState); ok {
				for _, f := range p.m.configuredOutputFormats {
					if !formatSet[f.Name] {
						formats = append(formats, f)
						formatSet[f.Name] = true
					}
				}
			}
			return false, nil
		},
	}

	if err := w.Walk(context.TODO()); err != nil {
		panic(err)
	}

	// Add the per kind configured output formats
	for _, kind := range kinds.AllKindsInPages {
		if siteFormats, found := s.conf.C.KindOutputFormats[kind]; found {
			for _, f := range siteFormats {
				if !formatSet[f.Name] {
					formats = append(formats, f)
					formatSet[f.Name] = true
				}
			}
		}
	}

	sort.Sort(formats)
	s.renderFormats = formats
}

func (s *Site) GetRelatedDocsHandler() *page.RelatedDocsHandler {
	return s.relatedDocsHandler
}

func (s *Site) Language() *langs.Language {
	return s.language
}

func (s *Site) Languages() langs.Languages {
	return s.h.Configs.Languages
}

type siteRefLinker struct {
	s *Site

	errorLogger logg.LevelLogger
	notFoundURL string
}

func newSiteRefLinker(s *Site) (siteRefLinker, error) {
	logger := s.Log.Error()

	notFoundURL := s.conf.RefLinksNotFoundURL
	errLevel := s.conf.RefLinksErrorLevel
	if strings.EqualFold(errLevel, "warning") {
		logger = s.Log.Warn()
	}
	return siteRefLinker{s: s, errorLogger: logger, notFoundURL: notFoundURL}, nil
}

func (s siteRefLinker) logNotFound(ref, what string, p page.Page, position text.Position) {
	if position.IsValid() {
		s.errorLogger.Logf("[%s] REF_NOT_FOUND: Ref %q: %s: %s", s.s.Lang(), ref, position.String(), what)
	} else if p == nil {
		s.errorLogger.Logf("[%s] REF_NOT_FOUND: Ref %q: %s", s.s.Lang(), ref, what)
	} else {
		s.errorLogger.Logf("[%s] REF_NOT_FOUND: Ref %q from page %q: %s", s.s.Lang(), ref, p.Path(), what)
	}
}

func (s *siteRefLinker) refLink(ref string, source any, relative bool, outputFormat string) (string, error) {
	p, err := unwrapPage(source)
	if err != nil {
		return "", err
	}

	var refURL *url.URL

	ref = filepath.ToSlash(ref)

	refURL, err = url.Parse(ref)
	if err != nil {
		return s.notFoundURL, err
	}

	var target page.Page
	var link string

	if refURL.Path != "" {
		var err error
		target, err = s.s.getPageRef(p, refURL.Path)
		var pos text.Position
		if err != nil || target == nil {
			if p, ok := source.(text.Positioner); ok {
				pos = p.Position()
			}
		}

		if err != nil {
			s.logNotFound(refURL.Path, err.Error(), p, pos)
			return s.notFoundURL, nil
		}

		if target == nil {
			s.logNotFound(refURL.Path, "page not found", p, pos)
			return s.notFoundURL, nil
		}

		var permalinker Permalinker = target

		if outputFormat != "" {
			o := target.OutputFormats().Get(outputFormat)

			if o == nil {
				s.logNotFound(refURL.Path, fmt.Sprintf("output format %q", outputFormat), p, pos)
				return s.notFoundURL, nil
			}
			permalinker = o
		}

		if relative {
			link = permalinker.RelPermalink()
		} else {
			link = permalinker.Permalink()
		}
	}

	if refURL.Fragment != "" {
		_ = target
		link = link + "#" + refURL.Fragment

		if pctx, ok := target.(pageContext); ok {
			if refURL.Path != "" {
				if di, ok := pctx.getContentConverter().(converter.DocumentInfo); ok {
					link = link + di.AnchorSuffix()
				}
			}
		} else if pctx, ok := p.(pageContext); ok {
			if di, ok := pctx.getContentConverter().(converter.DocumentInfo); ok {
				link = link + di.AnchorSuffix()
			}
		}

	}

	return link, nil
}

func (s *Site) watching() bool {
	return s.h != nil && s.h.Configs.Base.Internal.Watch
}

type whatChanged struct {
	mu sync.Mutex

	contentChanged bool
	identitySet    identity.Identities
}

func (w *whatChanged) Add(ids ...identity.Identity) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, id := range ids {
		w.identitySet[id] = true
	}
}

func (w *whatChanged) Changes() []identity.Identity {
	if w == nil || w.identitySet == nil {
		return nil
	}
	return w.identitySet.AsSlice()
}

// RegisterMediaTypes will register the Site's media types in the mime
// package, so it will behave correctly with Hugo's built-in server.
func (s *Site) RegisterMediaTypes() {
	for _, mt := range s.conf.MediaTypes.Config {
		for _, suffix := range mt.Suffixes() {
			_ = mime.AddExtensionType(mt.Delimiter+suffix, mt.Type)
		}
	}
}

func (h *HugoSites) fileEventsFilter(events []fsnotify.Event) []fsnotify.Event {
	seen := make(map[fsnotify.Event]bool)

	n := 0
	for _, ev := range events {
		// Avoid processing the same event twice.
		if seen[ev] {
			continue
		}
		seen[ev] = true

		if h.SourceSpec.IgnoreFile(ev.Name) {
			continue
		}

		if runtime.GOOS == "darwin" { // When a file system is HFS+, its filepath is in NFD form.
			ev.Name = norm.NFC.String(ev.Name)
		}

		events[n] = ev
		n++
	}
	events = events[:n]

	eventOrdinal := func(e fsnotify.Event) int {
		// Pull the structural changes to the top.
		if e.Op.Has(fsnotify.Create) {
			return 1
		}
		if e.Op.Has(fsnotify.Remove) {
			return 2
		}
		if e.Op.Has(fsnotify.Rename) {
			return 3
		}
		if e.Op.Has(fsnotify.Write) {
			return 4
		}
		return 5
	}

	sort.Slice(events, func(i, j int) bool {
		// First sort by event type.
		if eventOrdinal(events[i]) != eventOrdinal(events[j]) {
			return eventOrdinal(events[i]) < eventOrdinal(events[j])
		}
		// Then sort by name.
		return events[i].Name < events[j].Name
	})

	return events
}

type fileEventInfo struct {
	fsnotify.Event
	fi           os.FileInfo
	added        bool
	removed      bool
	isChangedDir bool
}

func (h *HugoSites) fileEventsApplyInfo(events []fsnotify.Event) []fileEventInfo {
	var infos []fileEventInfo
	for _, ev := range events {
		removed := false
		added := false

		if ev.Op&fsnotify.Remove == fsnotify.Remove {
			removed = true
		}

		fi, statErr := h.Fs.Source.Stat(ev.Name)

		// Some editors (Vim) sometimes issue only a Rename operation when writing an existing file
		// Sometimes a rename operation means that file has been renamed other times it means
		// it's been updated.
		if ev.Op.Has(fsnotify.Rename) {
			// If the file is still on disk, it's only been updated, if it's not, it's been moved
			if statErr != nil {
				removed = true
			}
		}
		if ev.Op.Has(fsnotify.Create) {
			added = true
		}

		isChangedDir := statErr == nil && fi.IsDir()

		infos = append(infos, fileEventInfo{
			Event:        ev,
			fi:           fi,
			added:        added,
			removed:      removed,
			isChangedDir: isChangedDir,
		})
	}

	n := 0

	for _, ev := range infos {
		// Remove any directories that's also represented by a file.
		keep := true
		if ev.isChangedDir {
			for _, ev2 := range infos {
				if ev2.fi != nil && !ev2.fi.IsDir() && filepath.Dir(ev2.Name) == ev.Name {
					keep = false
					break
				}
			}
		}
		if keep {
			infos[n] = ev
			n++
		}
	}
	infos = infos[:n]

	return infos
}

func (h *HugoSites) fileEventsTrim(events []fsnotify.Event) []fsnotify.Event {
	seen := make(map[string]bool)
	n := 0
	for _, ev := range events {
		if seen[ev.Name] {
			continue
		}
		seen[ev.Name] = true
		events[n] = ev
		n++
	}
	return events
}

func (h *HugoSites) fileEventsContentPaths(p []pathChange) []pathChange {
	var bundles []pathChange
	var dirs []pathChange
	var regular []pathChange

	var others []pathChange
	for _, p := range p {
		if p.isDir {
			dirs = append(dirs, p)
		} else {
			others = append(others, p)
		}
	}

	// Remove all files below dir.
	if len(dirs) > 0 {
		n := 0
		for _, d := range dirs {
			dir := d.p.Path() + "/"
			for _, o := range others {
				if !strings.HasPrefix(o.p.Path(), dir) {
					others[n] = o
					n++
				}
			}

		}
		others = others[:n]
	}

	for _, p := range others {
		if p.p.IsBundle() {
			bundles = append(bundles, p)
		} else {
			regular = append(regular, p)
		}
	}

	// Remove any files below leaf bundles.
	// Remove any files in the same folder as branch bundles.
	var keepers []pathChange

	for _, o := range regular {
		keep := true
		for _, b := range bundles {
			prefix := b.p.Base() + "/"
			if b.p.IsLeafBundle() && strings.HasPrefix(o.p.Path(), prefix) {
				keep = false
				break
			} else if b.p.IsBranchBundle() && o.p.Dir() == b.p.Dir() {
				keep = false
				break
			}
		}

		if keep {
			keepers = append(keepers, o)
		}
	}

	keepers = append(dirs, keepers...)
	keepers = append(bundles, keepers...)

	return keepers
}

// SitemapAbsURL is a convenience method giving the absolute URL to the sitemap.
func (s *Site) SitemapAbsURL() string {
	base := ""
	if len(s.conf.Languages) > 1 || s.Conf.DefaultContentLanguageInSubdir() {
		base = s.Language().Lang
	}
	p := s.AbsURL(base, false)
	if !strings.HasSuffix(p, "/") {
		p += "/"
	}
	p += s.conf.Sitemap.Filename
	return p
}

func (s *Site) createNodeMenuEntryURL(in string) string {
	if !strings.HasPrefix(in, "/") {
		return in
	}
	// make it match the nodes
	menuEntryURL := in
	menuEntryURL = s.s.PathSpec.URLize(menuEntryURL)
	if !s.conf.CanonifyURLs {
		menuEntryURL = paths.AddContextRoot(s.s.PathSpec.Cfg.BaseURL().String(), menuEntryURL)
	}
	return menuEntryURL
}

func (s *Site) assembleMenus() error {
	s.menus = make(navigation.Menus)

	type twoD struct {
		MenuName, EntryName string
	}
	flat := map[twoD]*navigation.MenuEntry{}
	children := map[twoD]navigation.Menu{}

	// add menu entries from config to flat hash
	for name, menu := range s.conf.Menus.Config {
		for _, me := range menu {
			if types.IsNil(me.Page) && me.PageRef != "" {
				// Try to resolve the page.
				me.Page, _ = s.getPage(nil, me.PageRef)
			}

			// If page is still nill, we must make sure that we have a URL that considers baseURL etc.
			if types.IsNil(me.Page) {
				me.ConfiguredURL = s.createNodeMenuEntryURL(me.MenuConfig.URL)
			}

			flat[twoD{name, me.KeyName()}] = me
		}
	}

	sectionPagesMenu := s.conf.SectionPagesMenu

	if sectionPagesMenu != "" {
		if err := s.pageMap.forEachPage(pagePredicates.ShouldListGlobal, func(p *pageState) (bool, error) {
			if p.IsHome() || !p.m.shouldBeCheckedForMenuDefinitions() {
				return false, nil
			}

			// The section pages menus are attached to the top level section.
			id := p.Section()
			if id == "" {
				id = "/"
			}

			if _, ok := flat[twoD{sectionPagesMenu, id}]; ok {
				return false, nil
			}
			me := navigation.MenuEntry{
				MenuConfig: navigation.MenuConfig{
					Identifier: id,
					Name:       p.LinkTitle(),
					Weight:     p.Weight(),
				},
				Page: p,
			}

			navigation.SetPageValues(&me, p)
			flat[twoD{sectionPagesMenu, me.KeyName()}] = &me
			return false, nil
		}); err != nil {
			return err
		}
	}

	// Add menu entries provided by pages
	if err := s.pageMap.forEachPage(pagePredicates.ShouldListGlobal, func(p *pageState) (bool, error) {
		for name, me := range p.pageMenus.menus() {
			if _, ok := flat[twoD{name, me.KeyName()}]; ok {
				err := p.wrapError(fmt.Errorf("duplicate menu entry with identifier %q in menu %q", me.KeyName(), name))
				s.Log.Warnln(err)
				continue
			}
			flat[twoD{name, me.KeyName()}] = me
		}
		return false, nil
	}); err != nil {
		return err
	}

	// Create Children Menus First
	for _, e := range flat {
		if e.Parent != "" {
			children[twoD{e.Menu, e.Parent}] = children[twoD{e.Menu, e.Parent}].Add(e)
		}
	}

	// Placing Children in Parents (in flat)
	for p, childmenu := range children {
		_, ok := flat[twoD{p.MenuName, p.EntryName}]
		if !ok {
			// if parent does not exist, create one without a URL
			flat[twoD{p.MenuName, p.EntryName}] = &navigation.MenuEntry{
				MenuConfig: navigation.MenuConfig{
					Name: p.EntryName,
				},
			}
		}
		flat[twoD{p.MenuName, p.EntryName}].Children = childmenu
	}

	// Assembling Top Level of Tree
	for menu, e := range flat {
		if e.Parent == "" {
			_, ok := s.menus[menu.MenuName]
			if !ok {
				s.menus[menu.MenuName] = navigation.Menu{}
			}
			s.menus[menu.MenuName] = s.menus[menu.MenuName].Add(e)
		}
	}

	return nil
}

// get any language code to prefix the target file path with.
func (s *Site) getLanguageTargetPathLang(alwaysInSubDir bool) string {
	if s.h.Conf.IsMultihost() {
		return s.Language().Lang
	}

	return s.getLanguagePermalinkLang(alwaysInSubDir)
}

// get any language code to prefix the relative permalink with.
func (s *Site) getLanguagePermalinkLang(alwaysInSubDir bool) string {
	if s.h.Conf.IsMultihost() {
		return ""
	}

	if s.h.Conf.IsMultilingual() && alwaysInSubDir {
		return s.Language().Lang
	}

	return s.GetLanguagePrefix()
}

// Prepare site for a new full build.
func (s *Site) resetBuildState(sourceChanged bool) {
	s.relatedDocsHandler = s.relatedDocsHandler.Clone()
	s.init.Reset()
}

func (s *Site) errorCollator(results <-chan error, errs chan<- error) {
	var errors []error
	for e := range results {
		errors = append(errors, e)
	}

	errs <- s.h.pickOneAndLogTheRest(errors)

	close(errs)
}

// GetPage looks up a page of a given type for the given ref.
// In Hugo <= 0.44 you had to add Page Kind (section, home) etc. as the first
// argument and then either a unix styled path (with or without a leading slash))
// or path elements separated.
// When we now remove the Kind from this API, we need to make the transition as painless
// as possible for existing sites. Most sites will use {{ .Site.GetPage "section" "my/section" }},
// i.e. 2 arguments, so we test for that.
func (s *Site) GetPage(ref ...string) (page.Page, error) {
	p, err := s.s.getPageForRefs(ref...)

	if p == nil {
		// The nil struct has meaning in some situations, mostly to avoid breaking
		// existing sites doing $nilpage.IsDescendant($p), which will always return
		// false.
		p = page.NilPage
	}

	return p, err
}

func (s *Site) absURLPath(targetPath string) string {
	var path string
	if s.conf.RelativeURLs {
		path = helpers.GetDottedRelativePath(targetPath)
	} else {
		url := s.PathSpec.Cfg.BaseURL().String()
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		path = url
	}

	return path
}

const (
	pageDependencyScopeDefault int = iota
	pageDependencyScopeGlobal
)

func (s *Site) renderAndWritePage(statCounter *uint64, name string, targetPath string, p *pageState, d any, templ tpl.Template) error {
	s.h.buildCounters.pageRenderCounter.Add(1)
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	of := p.outputFormat()
	p.incrRenderState()

	ctx := tpl.Context.Page.Set(context.Background(), p)
	ctx = tpl.Context.DependencyManagerScopedProvider.Set(ctx, p)

	if err := s.renderForTemplate(ctx, p.Kind(), of.Name, d, renderBuffer, templ); err != nil {
		return err
	}

	if renderBuffer.Len() == 0 {
		return nil
	}

	isHTML := of.IsHTML
	isRSS := of.Name == "rss"

	pd := publisher.Descriptor{
		Src:          renderBuffer,
		TargetPath:   targetPath,
		StatCounter:  statCounter,
		OutputFormat: p.outputFormat(),
	}

	if isRSS {
		// Always canonify URLs in RSS
		pd.AbsURLPath = s.absURLPath(targetPath)
	} else if isHTML {
		if s.conf.RelativeURLs || s.conf.CanonifyURLs {
			pd.AbsURLPath = s.absURLPath(targetPath)
		}

		if s.watching() && s.conf.Internal.Running && !s.conf.DisableLiveReload {
			pd.LiveReloadBaseURL = s.Conf.BaseURLLiveReload().URL()
		}

		// For performance reasons we only inject the Hugo generator tag on the home page.
		if p.IsHome() {
			pd.AddHugoGeneratorTag = !s.conf.DisableHugoGeneratorInject
		}

	}

	return s.publisher.Publish(pd)
}

var infoOnMissingLayout = map[string]bool{
	// The 404 layout is very much optional in Hugo, but we do look for it.
	"404": true,
}

// hookRendererTemplate is the canonical implementation of all hooks.ITEMRenderer,
// where ITEM is the thing being hooked.
type hookRendererTemplate struct {
	templateHandler tpl.TemplateHandler
	templ           tpl.Template
	resolvePosition func(ctx any) text.Position
}

func (hr hookRendererTemplate) RenderLink(cctx context.Context, w io.Writer, ctx hooks.LinkContext) error {
	return hr.templateHandler.ExecuteWithContext(cctx, hr.templ, w, ctx)
}

func (hr hookRendererTemplate) RenderHeading(cctx context.Context, w io.Writer, ctx hooks.HeadingContext) error {
	return hr.templateHandler.ExecuteWithContext(cctx, hr.templ, w, ctx)
}

func (hr hookRendererTemplate) RenderCodeblock(cctx context.Context, w hugio.FlexiWriter, ctx hooks.CodeblockContext) error {
	return hr.templateHandler.ExecuteWithContext(cctx, hr.templ, w, ctx)
}

func (hr hookRendererTemplate) ResolvePosition(ctx any) text.Position {
	return hr.resolvePosition(ctx)
}

func (hr hookRendererTemplate) IsDefaultCodeBlockRenderer() bool {
	return false
}

func (s *Site) renderForTemplate(ctx context.Context, name, outputFormat string, d any, w io.Writer, templ tpl.Template) (err error) {
	if templ == nil {
		s.logMissingLayout(name, "", "", outputFormat)
		return nil
	}

	if ctx == nil {
		panic("nil context")
	}

	if err = s.Tmpl().ExecuteWithContext(ctx, templ, w, d); err != nil {
		return fmt.Errorf("render of %q failed: %w", name, err)
	}
	return
}

func (s *Site) shouldBuild(p page.Page) bool {
	if !s.conf.IsKindEnabled(p.Kind()) {
		return false
	}
	return shouldBuild(s.Conf.BuildFuture(), s.Conf.BuildExpired(),
		s.Conf.BuildDrafts(), p.Draft(), p.PublishDate(), p.ExpiryDate())
}

func shouldBuild(buildFuture bool, buildExpired bool, buildDrafts bool, Draft bool,
	publishDate time.Time, expiryDate time.Time,
) bool {
	if !(buildDrafts || !Draft) {
		return false
	}
	hnow := htime.Now()
	if !buildFuture && !publishDate.IsZero() && publishDate.After(hnow) {
		return false
	}
	if !buildExpired && !expiryDate.IsZero() && expiryDate.Before(hnow) {
		return false
	}
	return true
}

func (s *Site) render(ctx *siteRenderContext) (err error) {
	if err := page.Clear(); err != nil {
		return err
	}

	if ctx.outIdx == 0 {
		// Note that even if disableAliases is set, the aliases themselves are
		// preserved on page. The motivation with this is to be able to generate
		// 301 redirects in a .htaccess file and similar using a custom output format.
		if !s.conf.DisableAliases {
			// Aliases must be rendered before pages.
			// Some sites, Hugo docs included, have faulty alias definitions that point
			// to itself or another real page. These will be overwritten in the next
			// step.
			if err = s.renderAliases(); err != nil {
				return
			}
		}
	}

	if err = s.renderPages(ctx); err != nil {
		return
	}

	if !ctx.shouldRenderStandalonePage("") {
		return
	}

	if err = s.renderMainLanguageRedirect(); err != nil {
		return
	}

	return
}
