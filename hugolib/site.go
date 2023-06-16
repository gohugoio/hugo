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
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/types"
	"golang.org/x/text/unicode/norm"

	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/markup/converter/hooks"

	"github.com/gohugoio/hugo/markup/converter"

	"github.com/gohugoio/hugo/hugofs/files"
	hglob "github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/common/text"

	"github.com/gohugoio/hugo/publisher"

	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/lazy"

	"github.com/fsnotify/fsnotify"
	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"

	"github.com/spf13/afero"
)

func (s *Site) Taxonomies() page.TaxonomyList {
	s.init.taxonomies.Do(context.Background())
	return s.taxonomies
}

type taxonomiesConfig map[string]string

func (t taxonomiesConfig) Values() []viewName {
	var vals []viewName
	for k, v := range t {
		vals = append(vals, viewName{singular: k, plural: v})
	}
	sort.Slice(vals, func(i, j int) bool {
		return vals[i].plural < vals[j].plural
	})

	return vals
}

type siteConfigHolder struct {
	sitemap          config.SitemapConfig
	taxonomiesConfig taxonomiesConfig
	timeout          time.Duration
	hasCJKLanguage   bool
	enableEmoji      bool
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

func (s *Site) initInit(ctx context.Context, init *lazy.Init, pctx pageContext) bool {
	_, err := init.Do(ctx)

	if err != nil {
		s.h.FatalError(pctx.wrapError(err))
	}
	return err == nil
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
		var sections page.Pages
		s.home.treeRef.m.collectSectionsRecursiveIncludingSelf(pageMapQuery{Prefix: s.home.treeRef.key}, func(n *contentNode) {
			sections = append(sections, n.p)
		})

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

		for _, sect := range sections {
			treeRef := sect.(treeRefProvider).getTreeRef()

			var pas page.Pages
			treeRef.m.collectPages(pageMapQuery{Prefix: treeRef.key + cmBranchSeparator}, func(c *contentNode) {
				pas = append(pas, c.p)
			})
			page.SortByDefault(pas)

			setNextPrev(pas)
		}

		// The root section only goes one level down.
		treeRef := s.home.getTreeRef()

		var pas page.Pages
		treeRef.m.collectPages(pageMapQuery{Prefix: treeRef.key + cmBranchSeparator}, func(c *contentNode) {
			pas = append(pas, c.p)
		})
		page.SortByDefault(pas)

		setNextPrev(pas)

		return nil, nil
	})

	s.init.menus = init.Branch(func(context.Context) (any, error) {
		s.assembleMenus()
		return nil, nil
	})

	s.init.taxonomies = init.Branch(func(context.Context) (any, error) {
		err := s.pageMap.assembleTaxonomies()
		return nil, err
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
	rssDisabled := !s.conf.IsKindEnabled("rss")
	s.pageMap.pageTrees.WalkRenderable(func(s string, n *contentNode) bool {
		for _, f := range n.p.m.configuredOutputFormats {
			if rssDisabled && f.Name == "rss" {
				// legacy
				continue
			}
			if !formatSet[f.Name] {
				formats = append(formats, f)
				formatSet[f.Name] = true
			}
		}
		return false
	})

	// Add the per kind configured output formats
	for _, kind := range allKindsInPages {
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

func (s *Site) isEnabled(kind string) bool {
	if kind == kindUnknown {
		panic("Unknown kind")
	}
	return s.conf.IsKindEnabled(kind)
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
		s.errorLogger.Logf("[%s] REF_NOT_FOUND: Ref %q from page %q: %s", s.s.Lang(), ref, p.Pathc(), what)
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
	source bool
	files  map[string]bool
}

// RegisterMediaTypes will register the Site's media types in the mime
// package, so it will behave correctly with Hugo's built-in server.
func (s *Site) RegisterMediaTypes() {
	for _, mt := range s.conf.MediaTypes.Config {
		for _, suffix := range mt.Suffixes() {
			_ = mime.AddExtensionType(mt.Delimiter+suffix, mt.Type+"; charset=utf-8")
		}
	}
}

func (s *Site) filterFileEvents(events []fsnotify.Event) []fsnotify.Event {
	var filtered []fsnotify.Event
	seen := make(map[fsnotify.Event]bool)

	for _, ev := range events {
		// Avoid processing the same event twice.
		if seen[ev] {
			continue
		}
		seen[ev] = true

		if s.SourceSpec.IgnoreFile(ev.Name) {
			continue
		}

		// Throw away any directories
		isRegular, err := s.SourceSpec.IsRegularSourceFile(ev.Name)
		if err != nil && herrors.IsNotExist(err) && (ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Rename == fsnotify.Rename) {
			// Force keep of event
			isRegular = true
		}
		if !isRegular {
			continue
		}

		if runtime.GOOS == "darwin" { // When a file system is HFS+, its filepath is in NFD form.
			ev.Name = norm.NFC.String(ev.Name)
		}

		filtered = append(filtered, ev)
	}

	return filtered
}

func (s *Site) translateFileEvents(events []fsnotify.Event) []fsnotify.Event {
	var filtered []fsnotify.Event

	eventMap := make(map[string][]fsnotify.Event)

	// We often get a Remove etc. followed by a Create, a Create followed by a Write.
	// Remove the superfluous events to mage the update logic simpler.
	for _, ev := range events {
		eventMap[ev.Name] = append(eventMap[ev.Name], ev)
	}

	for _, ev := range events {
		mapped := eventMap[ev.Name]

		// Keep one
		found := false
		var kept fsnotify.Event
		for i, ev2 := range mapped {
			if i == 0 {
				kept = ev2
			}

			if ev2.Op&fsnotify.Write == fsnotify.Write {
				kept = ev2
				found = true
			}

			if !found && ev2.Op&fsnotify.Create == fsnotify.Create {
				kept = ev2
			}
		}

		filtered = append(filtered, kept)
	}

	return filtered
}

// reBuild partially rebuilds a site given the filesystem events.
// It returns whatever the content source was changed.
// TODO(bep) clean up/rewrite this method.
func (s *Site) processPartial(config *BuildCfg, init func(config *BuildCfg) error, events []fsnotify.Event) error {
	events = s.filterFileEvents(events)
	events = s.translateFileEvents(events)

	changeIdentities := make(identity.Identities)

	s.Log.Debugf("Rebuild for events %q", events)

	h := s.h

	// First we need to determine what changed

	var (
		sourceChanged       = []fsnotify.Event{}
		sourceReallyChanged = []fsnotify.Event{}
		contentFilesChanged []string

		tmplChanged bool
		tmplAdded   bool
		dataChanged bool
		i18nChanged bool

		sourceFilesChanged = make(map[string]bool)
	)

	var cacheBusters []func(string) bool
	bcfg := s.conf.Build

	for _, ev := range events {
		component, relFilename := s.BaseFs.MakePathRelative(ev.Name)
		if relFilename != "" {
			p := hglob.NormalizePath(path.Join(component, relFilename))
			g, err := bcfg.MatchCacheBuster(s.Log, p)
			if err == nil && g != nil {
				cacheBusters = append(cacheBusters, g)
			}
		}

		id, found := s.eventToIdentity(ev)
		if found {
			changeIdentities[id] = id

			switch id.Type {
			case files.ComponentFolderContent:
				s.Log.Println("Source changed", ev)
				sourceChanged = append(sourceChanged, ev)
			case files.ComponentFolderLayouts:
				tmplChanged = true
				if !s.Tmpl().HasTemplate(id.Path) {
					tmplAdded = true
				}
				if tmplAdded {
					s.Log.Println("Template added", ev)
				} else {
					s.Log.Println("Template changed", ev)
				}

			case files.ComponentFolderData:
				s.Log.Println("Data changed", ev)
				dataChanged = true
			case files.ComponentFolderI18n:
				s.Log.Println("i18n changed", ev)
				i18nChanged = true

			}
		}
	}

	changed := &whatChanged{
		source: len(sourceChanged) > 0,
		files:  sourceFilesChanged,
	}

	config.whatChanged = changed

	if err := init(config); err != nil {
		return err
	}

	var cacheBusterOr func(string) bool
	if len(cacheBusters) > 0 {
		cacheBusterOr = func(s string) bool {
			for _, cb := range cacheBusters {
				if cb(s) {
					return true
				}
			}
			return false
		}
	}

	// These in memory resource caches will be rebuilt on demand.
	if len(cacheBusters) > 0 {
		s.h.ResourceSpec.ResourceCache.DeleteMatches(cacheBusterOr)
	}

	if tmplChanged || i18nChanged {
		s.h.init.Reset()
		var prototype *deps.Deps
		for i, s := range s.h.Sites {
			if err := s.Deps.Compile(prototype); err != nil {
				return err
			}
			if i == 0 {
				prototype = s.Deps
			}
		}
	}

	if dataChanged {
		s.h.init.data.Reset()
	}

	for _, ev := range sourceChanged {
		removed := false

		if ev.Op&fsnotify.Remove == fsnotify.Remove {
			removed = true
		}

		// Some editors (Vim) sometimes issue only a Rename operation when writing an existing file
		// Sometimes a rename operation means that file has been renamed other times it means
		// it's been updated
		if ev.Op&fsnotify.Rename == fsnotify.Rename {
			// If the file is still on disk, it's only been updated, if it's not, it's been moved
			if ex, err := afero.Exists(s.Fs.Source, ev.Name); !ex || err != nil {
				removed = true
			}
		}

		if removed && files.IsContentFile(ev.Name) {
			h.removePageByFilename(ev.Name)
		}

		sourceReallyChanged = append(sourceReallyChanged, ev)
		sourceFilesChanged[ev.Name] = true
	}

	if config.ErrRecovery || tmplAdded || dataChanged {
		h.resetPageState()
	} else {
		h.resetPageStateFromEvents(changeIdentities)
	}

	if len(sourceReallyChanged) > 0 || len(contentFilesChanged) > 0 {
		var filenamesChanged []string
		for _, e := range sourceReallyChanged {
			filenamesChanged = append(filenamesChanged, e.Name)
		}
		if len(contentFilesChanged) > 0 {
			filenamesChanged = append(filenamesChanged, contentFilesChanged...)
		}

		filenamesChanged = helpers.UniqueStringsReuse(filenamesChanged)

		if err := s.readAndProcessContent(*config, filenamesChanged...); err != nil {
			return err
		}

	}

	return nil
}

func (s *Site) process(config BuildCfg) (err error) {
	if err = s.readAndProcessContent(config); err != nil {
		err = fmt.Errorf("readAndProcessContent: %w", err)
		return
	}
	return err
}

func (s *Site) render(ctx *siteRenderContext) (err error) {
	if err := page.Clear(); err != nil {
		return err
	}

	if ctx.outIdx == 0 {
		// Note that even if disableAliases is set, the aliases themselves are
		// preserved on page. The motivation with this is to be able to generate
		// 301 redirects in a .htacess file and similar using a custom output format.
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

	if ctx.outIdx == 0 {
		if err = s.renderSitemap(); err != nil {
			return
		}

		if ctx.multihost {
			if err = s.renderRobotsTXT(); err != nil {
				return
			}
		}

		if err = s.render404(); err != nil {
			return
		}
	}

	if !ctx.renderSingletonPages() {
		return
	}

	if err = s.renderMainLanguageRedirect(); err != nil {
		return
	}

	return
}

// HomeAbsURL is a convenience method giving the absolute URL to the home page.
func (s *Site) HomeAbsURL() string {
	base := ""
	if len(s.conf.Languages) > 1 {
		base = s.Language().Lang
	}
	return s.AbsURL(base, false)
}

// SitemapAbsURL is a convenience method giving the absolute URL to the sitemap.
func (s *Site) SitemapAbsURL() string {
	p := s.HomeAbsURL()
	if !strings.HasSuffix(p, "/") {
		p += "/"
	}
	p += s.conf.Sitemap.Filename
	return p
}

func (s *Site) eventToIdentity(e fsnotify.Event) (identity.PathIdentity, bool) {
	for _, fs := range s.BaseFs.SourceFilesystems.FileSystems() {
		if p := fs.Path(e.Name); p != "" {
			return identity.NewPathIdentity(fs.Name, filepath.ToSlash(p)), true
		}
	}
	return identity.PathIdentity{}, false
}

func (s *Site) readAndProcessContent(buildConfig BuildCfg, filenames ...string) error {
	if s.Deps == nil {
		panic("nil deps on site")
	}

	sourceSpec := source.NewSourceSpec(s.PathSpec, buildConfig.ContentInclusionFilter, s.BaseFs.Content.Fs)

	proc := newPagesProcessor(s.h, sourceSpec)

	c := newPagesCollector(sourceSpec, s.h.getContentMaps(), s.Log, s.h.ContentChanges, proc, filenames...)

	if err := c.Collect(); err != nil {
		return err
	}

	return nil
}

func (s *Site) createNodeMenuEntryURL(in string) string {
	if !strings.HasPrefix(in, "/") {
		return in
	}
	// make it match the nodes
	menuEntryURL := in
	menuEntryURL = helpers.SanitizeURLKeepTrailingSlash(s.s.PathSpec.URLize(menuEntryURL))
	if !s.conf.CanonifyURLs {
		menuEntryURL = paths.AddContextRoot(s.s.PathSpec.Cfg.BaseURL().String(), menuEntryURL)
	}
	return menuEntryURL
}

func (s *Site) assembleMenus() {
	s.menus = make(navigation.Menus)

	type twoD struct {
		MenuName, EntryName string
	}
	flat := map[twoD]*navigation.MenuEntry{}
	children := map[twoD]navigation.Menu{}

	// add menu entries from config to flat hash
	for name, menu := range s.conf.Menus.Config {
		for _, me := range menu {
			if types.IsNil(me.Page) {
				if me.PageRef != "" {
					// Try to resolve the page.
					p, _ := s.getPageNew(nil, me.PageRef)
					if !types.IsNil(p) {
						navigation.SetPageValues(me, p)
					}
				}
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
		s.pageMap.sections.Walk(func(s string, v any) bool {
			p := v.(*contentNode).p
			if p.IsHome() {
				return false
			}
			// From Hugo 0.22 we have nested sections, but until we get a
			// feel of how that would work in this setting, let us keep
			// this menu for the top level only.
			id := p.Section()
			if _, ok := flat[twoD{sectionPagesMenu, id}]; ok {
				return false
			}

			me := navigation.MenuEntry{
				MenuConfig: navigation.MenuConfig{
					Identifier: id,
					Name:       p.LinkTitle(),
					Weight:     p.Weight(),
				},
			}
			navigation.SetPageValues(&me, p)
			flat[twoD{sectionPagesMenu, me.KeyName()}] = &me

			return false
		})
	}

	// Add menu entries provided by pages
	s.pageMap.pageTrees.WalkRenderable(func(ss string, n *contentNode) bool {
		p := n.p

		for name, me := range p.pageMenus.menus() {
			if _, ok := flat[twoD{name, me.KeyName()}]; ok {
				err := p.wrapError(fmt.Errorf("duplicate menu entry with identifier %q in menu %q", me.KeyName(), name))
				s.Log.Warnln(err)
				continue
			}
			flat[twoD{name, me.KeyName()}] = me
		}

		return false
	})

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
}

// get any language code to prefix the target file path with.
func (s *Site) getLanguageTargetPathLang(alwaysInSubDir bool) string {
	if s.h.Conf.IsMultihost() {
		return s.Language().Lang
	}

	return s.getLanguagePermalinkLang(alwaysInSubDir)
}

// get any lanaguagecode to prefix the relative permalink with.
func (s *Site) getLanguagePermalinkLang(alwaysInSubDir bool) string {
	if !s.h.isMultiLingual() || s.h.Conf.IsMultihost() {
		return ""
	}

	if alwaysInSubDir {
		return s.Language().Lang
	}

	isDefault := s.Language().Lang == s.conf.DefaultContentLanguage

	if !isDefault || s.conf.DefaultContentLanguageInSubdir {
		return s.Language().Lang
	}

	return ""
}

func (s *Site) getTaxonomyKey(key string) string {
	if s.conf.DisablePathToLower {
		return s.PathSpec.MakePath(key)
	}
	return strings.ToLower(s.PathSpec.MakePath(key))
}

// Prepare site for a new full build.
func (s *Site) resetBuildState(sourceChanged bool) {
	s.relatedDocsHandler = s.relatedDocsHandler.Clone()
	s.init.Reset()

	if sourceChanged {
		s.pageMap.contentMap.pageReverseIndex.Reset()
		s.PageCollections = newPageCollections(s.pageMap)
		s.pageMap.withEveryBundlePage(func(p *pageState) bool {
			p.pagePages = &pagePages{}
			if p.bucket != nil {
				p.bucket.pagesMapBucketPages = &pagesMapBucketPages{}
			}
			p.parent = nil
			p.Scratcher = maps.NewScratcher()
			return false
		})
	} else {
		s.pageMap.withEveryBundlePage(func(p *pageState) bool {
			p.Scratcher = maps.NewScratcher()
			return false
		})
	}
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
	p, err := s.s.getPageOldVersion(ref...)

	if p == nil {
		// The nil struct has meaning in some situations, mostly to avoid breaking
		// existing sites doing $nilpage.IsDescendant($p), which will always return
		// false.
		p = page.NilPage
	}

	return p, err
}

func (s *Site) GetPageWithTemplateInfo(info tpl.Info, ref ...string) (page.Page, error) {
	p, err := s.GetPage(ref...)
	if p != nil {
		// Track pages referenced by templates/shortcodes
		// when in server mode.
		if im, ok := info.(identity.Manager); ok {
			im.Add(p)
		}
	}
	return p, err
}

func (s *Site) permalink(link string) string {
	return s.PathSpec.PermalinkForBaseURL(link, s.PathSpec.Cfg.BaseURL().String())
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

func (s *Site) lookupLayouts(layouts ...string) tpl.Template {
	for _, l := range layouts {
		if templ, found := s.Tmpl().Lookup(l); found {
			return templ
		}
	}

	return nil
}

func (s *Site) renderAndWriteXML(ctx context.Context, statCounter *uint64, name string, targetPath string, d any, templ tpl.Template) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	if err := s.renderForTemplate(ctx, name, "", d, renderBuffer, templ); err != nil {
		return err
	}

	pd := publisher.Descriptor{
		Src:         renderBuffer,
		TargetPath:  targetPath,
		StatCounter: statCounter,
		// For the minification part of XML,
		// we currently only use the MIME type.
		OutputFormat: output.RSSFormat,
		AbsURLPath:   s.absURLPath(targetPath),
	}

	return s.publisher.Publish(pd)
}

func (s *Site) renderAndWritePage(statCounter *uint64, name string, targetPath string, p *pageState, templ tpl.Template) error {
	s.h.IncrPageRender()
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	of := p.outputFormat()
	ctx := tpl.SetPageInContext(context.Background(), p)

	if err := s.renderForTemplate(ctx, p.Kind(), of.Name, p, renderBuffer, templ); err != nil {
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

		if s.watching() && s.conf.Internal.Running && !s.conf.Internal.DisableLiveReload {
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
	identity.SearchProvider
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

func (s *Site) lookupTemplate(layouts ...string) (tpl.Template, bool) {
	for _, l := range layouts {
		if templ, found := s.Tmpl().Lookup(l); found {
			return templ, true
		}
	}

	return nil, false
}

func (s *Site) publish(statCounter *uint64, path string, r io.Reader, fs afero.Fs) (err error) {
	s.PathSpec.ProcessingStats.Incr(statCounter)

	return helpers.WriteToDisk(filepath.Clean(path), r, fs)
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

func (s *Site) kindFromSections(sections []string) string {
	if len(sections) == 0 {
		return page.KindHome
	}

	return s.kindFromSectionPath(path.Join(sections...))
}

func (s *Site) kindFromSectionPath(sectionPath string) string {
	var taxonomiesConfig taxonomiesConfig = s.conf.Taxonomies
	for _, plural := range taxonomiesConfig {
		if plural == sectionPath {
			return page.KindTaxonomy
		}

		if strings.HasPrefix(sectionPath, plural) {
			return page.KindTerm
		}

	}

	return page.KindSection
}

func (s *Site) newPage(
	n *contentNode,
	parentbBucket *pagesMapBucket,
	kind, title string,
	sections ...string) *pageState {
	m := map[string]any{}
	if title != "" {
		m["title"] = title
	}

	p, err := newPageFromMeta(
		n,
		parentbBucket,
		m,
		&pageMeta{
			s:        s,
			kind:     kind,
			sections: sections,
		})
	if err != nil {
		panic(err)
	}

	return p
}

func (s *Site) shouldBuild(p page.Page) bool {
	return shouldBuild(s.Conf.BuildFuture(), s.Conf.BuildExpired(),
		s.Conf.BuildDrafts(), p.Draft(), p.PublishDate(), p.ExpiryDate())
}

func shouldBuild(buildFuture bool, buildExpired bool, buildDrafts bool, Draft bool,
	publishDate time.Time, expiryDate time.Time) bool {
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
