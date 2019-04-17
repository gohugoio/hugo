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
	"html/template"
	"io"
	"log"
	"mime"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/common/text"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/publisher"
	_errors "github.com/pkg/errors"

	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/lazy"
	"golang.org/x/sync/errgroup"

	"github.com/gohugoio/hugo/media"

	"github.com/fsnotify/fsnotify"
	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"

	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Site contains all the information relevant for constructing a static
// site.  The basic flow of information is as follows:
//
// 1. A list of Files is parsed and then converted into Pages.
//
// 2. Pages contain sections (based on the file they were generated from),
//    aliases and slugs (included in a pages frontmatter) which are the
//    various targets that will get generated.  There will be canonical
//    listing.  The canonical path can be overruled based on a pattern.
//
// 3. Taxonomies are created via configuration and will present some aspect of
//    the final page and typically a perm url.
//
// 4. All Pages are passed through a template based on their desired
//    layout based on numerous different elements.
//
// 5. The entire collection of files is written to disk.
type Site struct {

	// The owning container. When multiple languages, there will be multiple
	// sites.
	h *HugoSites

	*PageCollections

	Taxonomies TaxonomyList

	taxonomyNodes *taxonomyNodeInfos

	Sections Taxonomy
	Info     SiteInfo

	layoutHandler *output.LayoutHandler

	buildStats *buildStats

	language *langs.Language

	siteCfg siteConfigHolder

	disabledKinds map[string]bool

	enableInlineShortcodes bool

	// Output formats defined in site config per Page Kind, or some defaults
	// if not set.
	// Output formats defined in Page front matter will override these.
	outputFormats map[string]output.Formats

	// All the output formats and media types available for this site.
	// These values will be merged from the Hugo defaults, the site config and,
	// finally, the language settings.
	outputFormatsConfig output.Formats
	mediaTypesConfig    media.Types

	siteConfigConfig SiteConfig

	// How to handle page front matter.
	frontmatterHandler pagemeta.FrontMatterHandler

	// We render each site for all the relevant output formats in serial with
	// this rendering context pointing to the current one.
	rc *siteRenderingContext

	// The output formats that we need to render this site in. This slice
	// will be fixed once set.
	// This will be the union of Site.Pages' outputFormats.
	// This slice will be sorted.
	renderFormats output.Formats

	// Logger etc.
	*deps.Deps `json:"-"`

	// The func used to title case titles.
	titleFunc func(s string) string

	relatedDocsHandler *page.RelatedDocsHandler
	siteRefLinker

	publisher publisher.Publisher

	menus navigation.Menus

	// Shortcut to the home page. Note that this may be nil if
	// home page, for some odd reason, is disabled.
	home *pageState

	// The last modification date of this site.
	lastmod time.Time

	// Lazily loaded site dependencies
	init *siteInit
}

type siteConfigHolder struct {
	sitemap          config.Sitemap
	taxonomiesConfig map[string]string
	timeout          time.Duration
	hasCJKLanguage   bool
	enableEmoji      bool
}

// Lazily loaded site dependencies.
type siteInit struct {
	prevNext          *lazy.Init
	prevNextInSection *lazy.Init
	menus             *lazy.Init
}

func (init *siteInit) Reset() {
	init.prevNext.Reset()
	init.prevNextInSection.Reset()
	init.menus.Reset()
}

func (s *Site) initInit(init *lazy.Init, pctx pageContext) {
	_, err := init.Do()
	if err != nil {
		s.h.FatalError(pctx.wrapError(err))
	}
}

func (s *Site) prepareInits() {
	s.init = &siteInit{}

	var init lazy.Init

	s.init.prevNext = init.Branch(func() (interface{}, error) {
		regularPages := s.findWorkPagesByKind(page.KindPage)
		for i, p := range regularPages {
			if p.posNextPrev == nil {
				continue
			}
			p.posNextPrev.nextPage = nil
			p.posNextPrev.prevPage = nil

			if i > 0 {
				p.posNextPrev.nextPage = regularPages[i-1]
			}

			if i < len(regularPages)-1 {
				p.posNextPrev.prevPage = regularPages[i+1]
			}
		}
		return nil, nil
	})

	s.init.prevNextInSection = init.Branch(func() (interface{}, error) {
		var rootSection []int
		for i, p1 := range s.workAllPages {
			if p1.IsPage() && p1.Section() == "" {
				rootSection = append(rootSection, i)
			}
			if p1.IsSection() && len(p1.SectionsEntries()) <= 1 {
				sectionPages := p1.Pages()
				for i, p2 := range sectionPages {
					p2s := p2.(*pageState)
					if p2s.posNextPrevSection == nil {
						continue
					}

					p2s.posNextPrevSection.nextPage = nil
					p2s.posNextPrevSection.prevPage = nil

					if i > 0 {
						p2s.posNextPrevSection.nextPage = sectionPages[i-1]
					}

					if i < len(sectionPages)-1 {
						p2s.posNextPrevSection.prevPage = sectionPages[i+1]
					}
				}
			}
		}

		for i, j := range rootSection {
			p := s.workAllPages[j]
			if i > 0 {
				p.posNextPrevSection.nextPage = s.workAllPages[rootSection[i-1]]
			}

			if i < len(rootSection)-1 {
				p.posNextPrevSection.prevPage = s.workAllPages[rootSection[i+1]]
			}
		}

		return nil, nil
	})

	s.init.menus = init.Branch(func() (interface{}, error) {
		s.assembleMenus()
		return nil, nil
	})

}

// Build stats for a given site.
type buildStats struct {
	draftCount   int
	futureCount  int
	expiredCount int
}

// TODO(bep) consolidate all site stats into this
func (b *buildStats) update(p page.Page) {
	if p.Draft() {
		b.draftCount++
	}

	if resource.IsFuture(p) {
		b.futureCount++
	}

	if resource.IsExpired(p) {
		b.expiredCount++
	}
}

type siteRenderingContext struct {
	output.Format
}

func (s *Site) Menus() navigation.Menus {
	s.init.menus.Do()
	return s.menus
}

func (s *Site) initRenderFormats() {
	formatSet := make(map[string]bool)
	formats := output.Formats{}
	for _, p := range s.workAllPages {
		for _, f := range p.m.configuredOutputFormats {
			if !formatSet[f.Name] {
				formats = append(formats, f)
				formatSet[f.Name] = true
			}
		}
	}

	// Add the per kind configured output formats
	for _, kind := range allKindsInPages {
		if siteFormats, found := s.outputFormats[kind]; found {
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

func (s *Site) isEnabled(kind string) bool {
	if kind == kindUnknown {
		panic("Unknown kind")
	}
	return !s.disabledKinds[kind]
}

// reset returns a new Site prepared for rebuild.
func (s *Site) reset() *Site {
	return &Site{Deps: s.Deps,
		layoutHandler:          output.NewLayoutHandler(),
		disabledKinds:          s.disabledKinds,
		titleFunc:              s.titleFunc,
		relatedDocsHandler:     s.relatedDocsHandler.Clone(),
		siteRefLinker:          s.siteRefLinker,
		outputFormats:          s.outputFormats,
		rc:                     s.rc,
		outputFormatsConfig:    s.outputFormatsConfig,
		frontmatterHandler:     s.frontmatterHandler,
		mediaTypesConfig:       s.mediaTypesConfig,
		language:               s.language,
		h:                      s.h,
		publisher:              s.publisher,
		siteConfigConfig:       s.siteConfigConfig,
		enableInlineShortcodes: s.enableInlineShortcodes,
		buildStats:             &buildStats{},
		init:                   s.init,
		PageCollections:        newPageCollections(),
		siteCfg:                s.siteCfg,
	}

}

// newSite creates a new site with the given configuration.
func newSite(cfg deps.DepsCfg) (*Site, error) {
	c := newPageCollections()

	if cfg.Language == nil {
		cfg.Language = langs.NewDefaultLanguage(cfg.Cfg)
	}

	disabledKinds := make(map[string]bool)
	for _, disabled := range cast.ToStringSlice(cfg.Language.Get("disableKinds")) {
		disabledKinds[disabled] = true
	}

	var (
		mediaTypesConfig    []map[string]interface{}
		outputFormatsConfig []map[string]interface{}

		siteOutputFormatsConfig output.Formats
		siteMediaTypesConfig    media.Types
		err                     error
	)

	// Add language last, if set, so it gets precedence.
	for _, cfg := range []config.Provider{cfg.Cfg, cfg.Language} {
		if cfg.IsSet("mediaTypes") {
			mediaTypesConfig = append(mediaTypesConfig, cfg.GetStringMap("mediaTypes"))
		}
		if cfg.IsSet("outputFormats") {
			outputFormatsConfig = append(outputFormatsConfig, cfg.GetStringMap("outputFormats"))
		}
	}

	siteMediaTypesConfig, err = media.DecodeTypes(mediaTypesConfig...)
	if err != nil {
		return nil, err
	}

	siteOutputFormatsConfig, err = output.DecodeFormats(siteMediaTypesConfig, outputFormatsConfig...)
	if err != nil {
		return nil, err
	}

	outputFormats, err := createSiteOutputFormats(siteOutputFormatsConfig, cfg.Language)
	if err != nil {
		return nil, err
	}

	taxonomies := cfg.Language.GetStringMapString("taxonomies")

	var relatedContentConfig related.Config

	if cfg.Language.IsSet("related") {
		relatedContentConfig, err = related.DecodeConfig(cfg.Language.Get("related"))
		if err != nil {
			return nil, err
		}
	} else {
		relatedContentConfig = related.DefaultConfig
		if _, found := taxonomies["tag"]; found {
			relatedContentConfig.Add(related.IndexConfig{Name: "tags", Weight: 80})
		}
	}

	titleFunc := helpers.GetTitleFunc(cfg.Language.GetString("titleCaseStyle"))

	frontMatterHandler, err := pagemeta.NewFrontmatterHandler(cfg.Logger, cfg.Cfg)
	if err != nil {
		return nil, err
	}

	siteConfig := siteConfigHolder{
		sitemap:          config.DecodeSitemap(config.Sitemap{Priority: -1, Filename: "sitemap.xml"}, cfg.Language.GetStringMap("sitemap")),
		taxonomiesConfig: taxonomies,
		timeout:          time.Duration(cfg.Language.GetInt("timeout")) * time.Millisecond,
		hasCJKLanguage:   cfg.Language.GetBool("hasCJKLanguage"),
		enableEmoji:      cfg.Language.Cfg.GetBool("enableEmoji"),
	}

	s := &Site{
		PageCollections:        c,
		layoutHandler:          output.NewLayoutHandler(),
		language:               cfg.Language,
		disabledKinds:          disabledKinds,
		titleFunc:              titleFunc,
		relatedDocsHandler:     page.NewRelatedDocsHandler(relatedContentConfig),
		outputFormats:          outputFormats,
		rc:                     &siteRenderingContext{output.HTMLFormat},
		outputFormatsConfig:    siteOutputFormatsConfig,
		mediaTypesConfig:       siteMediaTypesConfig,
		frontmatterHandler:     frontMatterHandler,
		buildStats:             &buildStats{},
		enableInlineShortcodes: cfg.Language.GetBool("enableInlineShortcodes"),
		siteCfg:                siteConfig,
	}

	s.prepareInits()

	return s, nil

}

// NewSite creates a new site with the given dependency configuration.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
func NewSite(cfg deps.DepsCfg) (*Site, error) {
	s, err := newSite(cfg)
	if err != nil {
		return nil, err
	}

	if err = applyDeps(cfg, s); err != nil {
		return nil, err
	}

	return s, nil
}

// NewSiteDefaultLang creates a new site in the default language.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
// TODO(bep) test refactor -- remove
func NewSiteDefaultLang(withTemplate ...func(templ tpl.TemplateHandler) error) (*Site, error) {
	v := viper.New()
	if err := loadDefaultSettingsFor(v); err != nil {
		return nil, err
	}
	return newSiteForLang(langs.NewDefaultLanguage(v), withTemplate...)
}

// NewEnglishSite creates a new site in English language.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
// TODO(bep) test refactor -- remove
func NewEnglishSite(withTemplate ...func(templ tpl.TemplateHandler) error) (*Site, error) {
	v := viper.New()
	if err := loadDefaultSettingsFor(v); err != nil {
		return nil, err
	}
	return newSiteForLang(langs.NewLanguage("en", v), withTemplate...)
}

// newSiteForLang creates a new site in the given language.
func newSiteForLang(lang *langs.Language, withTemplate ...func(templ tpl.TemplateHandler) error) (*Site, error) {
	withTemplates := func(templ tpl.TemplateHandler) error {
		for _, wt := range withTemplate {
			if err := wt(templ); err != nil {
				return err
			}
		}
		return nil
	}

	cfg := deps.DepsCfg{WithTemplate: withTemplates, Cfg: lang}

	return NewSiteForCfg(cfg)

}

// NewSiteForCfg creates a new site for the given configuration.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
func NewSiteForCfg(cfg deps.DepsCfg) (*Site, error) {
	h, err := NewHugoSites(cfg)
	if err != nil {
		return nil, err
	}
	return h.Sites[0], nil

}

type SiteInfo struct {
	Authors page.AuthorList
	Social  SiteSocial

	hugoInfo     hugo.Info
	title        string
	RSSLink      string
	Author       map[string]interface{}
	LanguageCode string
	Copyright    string

	permalinks map[string]string

	LanguagePrefix string
	Languages      langs.Languages

	BuildDrafts bool

	canonifyURLs bool
	relativeURLs bool
	uglyURLs     func(p page.Page) bool

	owner                          *HugoSites
	s                              *Site
	language                       *langs.Language
	defaultContentLanguageInSubdir bool
	sectionPagesMenu               string
}

func (s *SiteInfo) Pages() page.Pages {
	return s.s.Pages()

}

func (s *SiteInfo) RegularPages() page.Pages {
	return s.s.RegularPages()

}

func (s *SiteInfo) AllPages() page.Pages {
	return s.s.AllPages()
}

func (s *SiteInfo) AllRegularPages() page.Pages {
	return s.s.AllRegularPages()
}

func (s *SiteInfo) Permalinks() map[string]string {
	// Remove in 0.57
	helpers.Deprecated("Site", ".Permalinks", "", false)
	return s.permalinks
}

func (s *SiteInfo) LastChange() time.Time {
	return s.s.lastmod
}

func (s *SiteInfo) Title() string {
	return s.title
}

func (s *SiteInfo) Site() page.Site {
	return s
}

func (s *SiteInfo) Menus() navigation.Menus {
	return s.s.Menus()
}

// TODO(bep) type
func (s *SiteInfo) Taxonomies() interface{} {
	return s.s.Taxonomies
}

func (s *SiteInfo) Params() map[string]interface{} {
	return s.s.Language().Params()
}

func (s *SiteInfo) Data() map[string]interface{} {
	return s.s.h.Data()
}

func (s *SiteInfo) Language() *langs.Language {
	return s.language
}

func (s *SiteInfo) Config() SiteConfig {
	return s.s.siteConfigConfig
}

func (s *SiteInfo) Hugo() hugo.Info {
	return s.hugoInfo
}

// Sites is a convenience method to get all the Hugo sites/languages configured.
func (s *SiteInfo) Sites() page.Sites {
	return s.s.h.siteInfos()
}

func (s *SiteInfo) String() string {
	return fmt.Sprintf("Site(%q)", s.title)
}

func (s *SiteInfo) BaseURL() template.URL {
	return template.URL(s.s.PathSpec.BaseURL.String())
}

// ServerPort returns the port part of the BaseURL, 0 if none found.
func (s *SiteInfo) ServerPort() int {
	ps := s.s.PathSpec.BaseURL.URL().Port()
	if ps == "" {
		return 0
	}
	p, err := strconv.Atoi(ps)
	if err != nil {
		return 0
	}
	return p
}

// GoogleAnalytics is kept here for historic reasons.
func (s *SiteInfo) GoogleAnalytics() string {
	return s.Config().Services.GoogleAnalytics.ID

}

// DisqusShortname is kept here for historic reasons.
func (s *SiteInfo) DisqusShortname() string {
	return s.Config().Services.Disqus.Shortname
}

// SiteSocial is a place to put social details on a site level. These are the
// standard keys that themes will expect to have available, but can be
// expanded to any others on a per site basis
// github
// facebook
// facebook_admin
// twitter
// twitter_domain
// googleplus
// pinterest
// instagram
// youtube
// linkedin
type SiteSocial map[string]string

// Param is a convenience method to do lookups in SiteInfo's Params map.
//
// This method is also implemented on Page and Node.
func (s *SiteInfo) Param(key interface{}) (interface{}, error) {
	keyStr, err := cast.ToStringE(key)
	if err != nil {
		return nil, err
	}
	keyStr = strings.ToLower(keyStr)
	return s.Params()[keyStr], nil
}

func (s *SiteInfo) IsMultiLingual() bool {
	return len(s.Languages) > 1
}

func (s *SiteInfo) IsServer() bool {
	return s.owner.running
}

type siteRefLinker struct {
	s *Site

	errorLogger *log.Logger
	notFoundURL string
}

func newSiteRefLinker(cfg config.Provider, s *Site) (siteRefLinker, error) {
	logger := s.Log.ERROR

	notFoundURL := cfg.GetString("refLinksNotFoundURL")
	errLevel := cfg.GetString("refLinksErrorLevel")
	if strings.EqualFold(errLevel, "warning") {
		logger = s.Log.WARN
	}
	return siteRefLinker{s: s, errorLogger: logger, notFoundURL: notFoundURL}, nil
}

func (s siteRefLinker) logNotFound(ref, what string, p page.Page, position text.Position) {
	if position.IsValid() {
		s.errorLogger.Printf("[%s] REF_NOT_FOUND: Ref %q: %s: %s", s.s.Lang(), ref, position.String(), what)
	} else if p == nil {
		s.errorLogger.Printf("[%s] REF_NOT_FOUND: Ref %q: %s", s.s.Lang(), ref, what)
	} else {
		s.errorLogger.Printf("[%s] REF_NOT_FOUND: Ref %q from page %q: %s", s.s.Lang(), ref, p.Path(), what)
	}
}

func (s *siteRefLinker) refLink(ref string, source interface{}, relative bool, outputFormat string) (string, error) {

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
		target, err := s.s.getPageNew(p, refURL.Path)
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
		if pctx, ok := target.(pageContext); ok && !target.File().IsZero() && !pctx.getRenderingConfig().PlainIDAnchors {
			if refURL.Path != "" {
				link = link + ":" + target.File().UniqueID()
			}
		} else if pctx, ok := p.(pageContext); ok && !p.File().IsZero() && !pctx.getRenderingConfig().PlainIDAnchors {
			link = link + ":" + p.File().UniqueID()
		}

	}
	return link, nil
}

// Ref will give an absolute URL to ref in the given Page.
func (s *SiteInfo) Ref(ref string, page page.Page, options ...string) (string, error) {
	// Remove in Hugo 0.54
	helpers.Deprecated("Site", ".Ref", "Use .Site.GetPage", true)
	outputFormat := ""
	if len(options) > 0 {
		outputFormat = options[0]
	}

	return s.s.refLink(ref, page, false, outputFormat)
}

// RelRef will give an relative URL to ref in the given Page.
func (s *SiteInfo) RelRef(ref string, page page.Page, options ...string) (string, error) {
	// Remove in Hugo 0.54
	helpers.Deprecated("Site", ".RelRef", "Use .Site.GetPage", true)
	outputFormat := ""
	if len(options) > 0 {
		outputFormat = options[0]
	}

	return s.s.refLink(ref, page, true, outputFormat)
}

func (s *Site) running() bool {
	return s.h != nil && s.h.running
}

func (s *Site) multilingual() *Multilingual {
	return s.h.multilingual
}

type whatChanged struct {
	source bool
	other  bool
	files  map[string]bool
}

// RegisterMediaTypes will register the Site's media types in the mime
// package, so it will behave correctly with Hugo's built-in server.
func (s *Site) RegisterMediaTypes() {
	for _, mt := range s.mediaTypesConfig {
		for _, suffix := range mt.Suffixes {
			_ = mime.AddExtensionType(mt.Delimiter+suffix, mt.Type()+"; charset=utf-8")
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
		if err != nil && os.IsNotExist(err) && (ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Rename == fsnotify.Rename) {
			// Force keep of event
			isRegular = true
		}
		if !isRegular {
			continue
		}

		filtered = append(filtered, ev)
	}

	return filtered
}

func (s *Site) translateFileEvents(events []fsnotify.Event) []fsnotify.Event {
	var filtered []fsnotify.Event

	eventMap := make(map[string][]fsnotify.Event)

	// We often get a Remove etc. followed by a Create, a Create followed by a Write.
	// Remove the superflous events to mage the update logic simpler.
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
// It returns whetever the content source was changed.
// TODO(bep) clean up/rewrite this method.
func (s *Site) processPartial(events []fsnotify.Event) (whatChanged, error) {

	events = s.filterFileEvents(events)
	events = s.translateFileEvents(events)

	s.Log.DEBUG.Printf("Rebuild for events %q", events)

	h := s.h

	// First we need to determine what changed

	var (
		sourceChanged       = []fsnotify.Event{}
		sourceReallyChanged = []fsnotify.Event{}
		contentFilesChanged []string
		tmplChanged         = []fsnotify.Event{}
		dataChanged         = []fsnotify.Event{}
		i18nChanged         = []fsnotify.Event{}
		shortcodesChanged   = make(map[string]bool)
		sourceFilesChanged  = make(map[string]bool)

		// prevent spamming the log on changes
		logger = helpers.NewDistinctFeedbackLogger()
	)

	cachePartitions := make([]string, len(events))

	for i, ev := range events {
		cachePartitions[i] = resources.ResourceKeyPartition(ev.Name)

		if s.isContentDirEvent(ev) {
			logger.Println("Source changed", ev)
			sourceChanged = append(sourceChanged, ev)
		}
		if s.isLayoutDirEvent(ev) {
			logger.Println("Template changed", ev)
			tmplChanged = append(tmplChanged, ev)

			if strings.Contains(ev.Name, "shortcodes") {
				shortcode := filepath.Base(ev.Name)
				shortcode = strings.TrimSuffix(shortcode, filepath.Ext(shortcode))
				shortcodesChanged[shortcode] = true
			}
		}
		if s.isDataDirEvent(ev) {
			logger.Println("Data changed", ev)
			dataChanged = append(dataChanged, ev)
		}
		if s.isI18nEvent(ev) {
			logger.Println("i18n changed", ev)
			i18nChanged = append(dataChanged, ev)
		}
	}

	// These in memory resource caches will be rebuilt on demand.
	for _, s := range s.h.Sites {
		s.ResourceSpec.ResourceCache.DeletePartitions(cachePartitions...)
	}

	if len(tmplChanged) > 0 || len(i18nChanged) > 0 {
		sites := s.h.Sites
		first := sites[0]

		s.h.init.Reset()

		// TOD(bep) globals clean
		if err := first.Deps.LoadResources(); err != nil {
			return whatChanged{}, err
		}

		for i := 1; i < len(sites); i++ {
			site := sites[i]
			var err error
			depsCfg := deps.DepsCfg{
				Language:      site.language,
				MediaTypes:    site.mediaTypesConfig,
				OutputFormats: site.outputFormatsConfig,
			}
			site.Deps, err = first.Deps.ForLanguage(depsCfg, func(d *deps.Deps) error {
				d.Site = &site.Info
				return nil
			})
			if err != nil {
				return whatChanged{}, err
			}
		}
	}

	if len(dataChanged) > 0 {
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
		if removed && IsContentFile(ev.Name) {
			h.removePageByFilename(ev.Name)
		}

		sourceReallyChanged = append(sourceReallyChanged, ev)
		sourceFilesChanged[ev.Name] = true
	}

	for shortcode := range shortcodesChanged {
		// There are certain scenarios that, when a shortcode changes,
		// it isn't sufficient to just rerender the already parsed shortcode.
		// One example is if the user adds a new shortcode to the content file first,
		// and then creates the shortcode on the file system.
		// To handle these scenarios, we must do a full reprocessing of the
		// pages that keeps a reference to the changed shortcode.
		pagesWithShortcode := h.findPagesByShortcode(shortcode)
		for _, p := range pagesWithShortcode {
			contentFilesChanged = append(contentFilesChanged, p.File().Filename())
		}
	}

	if len(sourceReallyChanged) > 0 || len(contentFilesChanged) > 0 {
		var filenamesChanged []string
		for _, e := range sourceReallyChanged {
			filenamesChanged = append(filenamesChanged, e.Name)
		}
		if len(contentFilesChanged) > 0 {
			filenamesChanged = append(filenamesChanged, contentFilesChanged...)
		}

		filenamesChanged = helpers.UniqueStrings(filenamesChanged)

		if err := s.readAndProcessContent(filenamesChanged...); err != nil {
			return whatChanged{}, err
		}

	}

	changed := whatChanged{
		source: len(sourceChanged) > 0 || len(shortcodesChanged) > 0,
		other:  len(tmplChanged) > 0 || len(i18nChanged) > 0 || len(dataChanged) > 0,
		files:  sourceFilesChanged,
	}

	return changed, nil

}

func (s *Site) process(config BuildCfg) (err error) {
	if err = s.initialize(); err != nil {
		return
	}
	if err := s.readAndProcessContent(); err != nil {
		return err
	}
	return err

}

func (s *Site) setupSitePages() {
	var homeDates *resource.Dates
	if s.home != nil {
		// If the home page has no dates set, we fall back to the site dates.
		homeDates = &s.home.m.Dates
	}

	if !s.lastmod.IsZero() && (homeDates == nil || !resource.IsZeroDates(homeDates)) {
		return
	}

	if homeDates != nil && !s.lastmod.IsZero() {
		homeDates.FDate = s.lastmod
		homeDates.FLastmod = s.lastmod
		return

	}

	var siteLastmod time.Time
	var siteLastDate time.Time

	for _, page := range s.workAllPages {
		if !page.IsPage() {
			continue
		}
		// Determine Site.Info.LastChange
		// Note that the logic to determine which date to use for Lastmod
		// is already applied, so this is *the* date to use.
		// We cannot just pick the last page in the default sort, because
		// that may not be ordered by date.
		// TODO(bep) check if this can be done earlier
		if page.Lastmod().After(siteLastmod) {
			siteLastmod = page.Lastmod()
		}
		if page.Date().After(siteLastDate) {
			siteLastDate = page.Date()
		}
	}

	s.lastmod = siteLastmod

	if homeDates != nil && resource.IsZeroDates(homeDates) {
		homeDates.FDate = siteLastDate
		homeDates.FLastmod = s.lastmod
	}

}

func (s *Site) render(ctx *siteRenderContext) (err error) {

	if err := page.Clear(); err != nil {
		return err
	}

	if ctx.outIdx == 0 {
		// Note that even if disableAliases is set, the aliases themselves are
		// preserved on page. The motivation with this is to be able to generate
		// 301 redirects in a .htacess file and similar using a custom output format.
		if !s.Cfg.GetBool("disableAliases") {
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

		if err = s.renderRobotsTXT(); err != nil {
			return
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

func (s *Site) Initialise() (err error) {
	return s.initialize()
}

func (s *Site) initialize() (err error) {
	return s.initializeSiteInfo()
}

// HomeAbsURL is a convenience method giving the absolute URL to the home page.
func (s *SiteInfo) HomeAbsURL() string {
	base := ""
	if s.IsMultiLingual() {
		base = s.Language().Lang
	}
	return s.owner.AbsURL(base, false)
}

// SitemapAbsURL is a convenience method giving the absolute URL to the sitemap.
func (s *SiteInfo) SitemapAbsURL() string {
	p := s.HomeAbsURL()
	if !strings.HasSuffix(p, "/") {
		p += "/"
	}
	p += s.s.siteCfg.sitemap.Filename
	return p
}

func (s *Site) initializeSiteInfo() error {
	var (
		lang      = s.language
		languages langs.Languages
	)

	if s.h != nil && s.h.multilingual != nil {
		languages = s.h.multilingual.Languages
	}

	permalinks := s.Cfg.GetStringMapString("permalinks")

	defaultContentInSubDir := s.Cfg.GetBool("defaultContentLanguageInSubdir")
	defaultContentLanguage := s.Cfg.GetString("defaultContentLanguage")

	languagePrefix := ""
	if s.multilingualEnabled() && (defaultContentInSubDir || lang.Lang != defaultContentLanguage) {
		languagePrefix = "/" + lang.Lang
	}

	var uglyURLs = func(p page.Page) bool {
		return false
	}

	v := s.Cfg.Get("uglyURLs")
	if v != nil {
		switch vv := v.(type) {
		case bool:
			uglyURLs = func(p page.Page) bool {
				return vv
			}
		case string:
			// Is what be get from CLI (--uglyURLs)
			vvv := cast.ToBool(vv)
			uglyURLs = func(p page.Page) bool {
				return vvv
			}
		default:
			m := cast.ToStringMapBool(v)
			uglyURLs = func(p page.Page) bool {
				return m[p.Section()]
			}
		}
	}

	s.Info = SiteInfo{
		title:                          lang.GetString("title"),
		Author:                         lang.GetStringMap("author"),
		Social:                         lang.GetStringMapString("social"),
		LanguageCode:                   lang.GetString("languageCode"),
		Copyright:                      lang.GetString("copyright"),
		language:                       lang,
		LanguagePrefix:                 languagePrefix,
		Languages:                      languages,
		defaultContentLanguageInSubdir: defaultContentInSubDir,
		sectionPagesMenu:               lang.GetString("sectionPagesMenu"),
		BuildDrafts:                    s.Cfg.GetBool("buildDrafts"),
		canonifyURLs:                   s.Cfg.GetBool("canonifyURLs"),
		relativeURLs:                   s.Cfg.GetBool("relativeURLs"),
		uglyURLs:                       uglyURLs,
		permalinks:                     permalinks,
		owner:                          s.h,
		s:                              s,
		hugoInfo:                       hugo.NewInfo(s.Cfg.GetString("environment")),
	}

	rssOutputFormat, found := s.outputFormats[page.KindHome].GetByName(output.RSSFormat.Name)

	if found {
		s.Info.RSSLink = s.permalink(rssOutputFormat.BaseFilename())
	}

	return nil
}

func (s *Site) isI18nEvent(e fsnotify.Event) bool {
	return s.BaseFs.SourceFilesystems.IsI18n(e.Name)
}

func (s *Site) isDataDirEvent(e fsnotify.Event) bool {
	return s.BaseFs.SourceFilesystems.IsData(e.Name)
}

func (s *Site) isLayoutDirEvent(e fsnotify.Event) bool {
	return s.BaseFs.SourceFilesystems.IsLayout(e.Name)
}

func (s *Site) isContentDirEvent(e fsnotify.Event) bool {
	return s.BaseFs.IsContent(e.Name)
}

type contentCaptureResultHandler struct {
	defaultContentProcessor *siteContentProcessor
	contentProcessors       map[string]*siteContentProcessor
}

func (c *contentCaptureResultHandler) getContentProcessor(lang string) *siteContentProcessor {
	proc, found := c.contentProcessors[lang]
	if found {
		return proc
	}
	return c.defaultContentProcessor
}

func (c *contentCaptureResultHandler) handleSingles(fis ...*fileInfo) {
	for _, fi := range fis {
		proc := c.getContentProcessor(fi.Lang())
		proc.processSingle(fi)
	}
}
func (c *contentCaptureResultHandler) handleBundles(d *bundleDirs) {
	for _, b := range d.bundles {
		proc := c.getContentProcessor(b.fi.Lang())
		proc.processBundle(b)
	}
}

func (c *contentCaptureResultHandler) handleCopyFile(f pathLangFile) {
	proc := c.getContentProcessor(f.Lang())
	proc.processAsset(f)
}

func (s *Site) readAndProcessContent(filenames ...string) error {

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	defaultContentLanguage := s.SourceSpec.DefaultContentLanguage

	contentProcessors := make(map[string]*siteContentProcessor)
	var defaultContentProcessor *siteContentProcessor
	sites := s.h.langSite()
	for k, v := range sites {
		if v.language.Disabled {
			continue
		}
		proc := newSiteContentProcessor(ctx, len(filenames) > 0, v)
		contentProcessors[k] = proc
		if k == defaultContentLanguage {
			defaultContentProcessor = proc
		}
		g.Go(func() error {
			return proc.process(ctx)
		})
	}

	var (
		handler   captureResultHandler
		bundleMap *contentChangeMap
	)

	mainHandler := &contentCaptureResultHandler{contentProcessors: contentProcessors, defaultContentProcessor: defaultContentProcessor}

	sourceSpec := source.NewSourceSpec(s.PathSpec, s.BaseFs.Content.Fs)

	if s.running() {
		// Need to track changes.
		bundleMap = s.h.ContentChanges
		handler = &captureResultHandlerChain{handlers: []captureBundlesHandler{mainHandler, bundleMap}}

	} else {
		handler = mainHandler
	}

	c := newCapturer(s.Log, sourceSpec, handler, bundleMap, filenames...)

	err1 := c.capture()

	for _, proc := range contentProcessors {
		proc.closeInput()
	}

	err2 := g.Wait()

	if err1 != nil {
		return err1
	}
	return err2
}

func (s *Site) getMenusFromConfig() navigation.Menus {

	ret := navigation.Menus{}

	if menus := s.language.GetStringMap("menus"); menus != nil {
		for name, menu := range menus {
			m, err := cast.ToSliceE(menu)
			if err != nil {
				s.Log.ERROR.Printf("unable to process menus in site config\n")
				s.Log.ERROR.Println(err)
			} else {
				for _, entry := range m {
					s.Log.DEBUG.Printf("found menu: %q, in site config\n", name)

					menuEntry := navigation.MenuEntry{Menu: name}
					ime, err := cast.ToStringMapE(entry)
					if err != nil {
						s.Log.ERROR.Printf("unable to process menus in site config\n")
						s.Log.ERROR.Println(err)
					}

					menuEntry.MarshallMap(ime)
					menuEntry.URL = s.Info.createNodeMenuEntryURL(menuEntry.URL)

					if ret[name] == nil {
						ret[name] = navigation.Menu{}
					}
					ret[name] = ret[name].Add(&menuEntry)
				}
			}
		}
		return ret
	}
	return ret
}

func (s *SiteInfo) createNodeMenuEntryURL(in string) string {

	if !strings.HasPrefix(in, "/") {
		return in
	}
	// make it match the nodes
	menuEntryURL := in
	menuEntryURL = helpers.SanitizeURLKeepTrailingSlash(s.s.PathSpec.URLize(menuEntryURL))
	if !s.canonifyURLs {
		menuEntryURL = helpers.AddContextRoot(s.s.PathSpec.BaseURL.String(), menuEntryURL)
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
	menuConfig := s.getMenusFromConfig()
	for name, menu := range menuConfig {
		for _, me := range menu {
			flat[twoD{name, me.KeyName()}] = me
		}
	}

	sectionPagesMenu := s.Info.sectionPagesMenu

	if sectionPagesMenu != "" {
		for _, p := range s.workAllPages {
			if p.Kind() == page.KindSection {
				// From Hugo 0.22 we have nested sections, but until we get a
				// feel of how that would work in this setting, let us keep
				// this menu for the top level only.
				id := p.Section()
				if _, ok := flat[twoD{sectionPagesMenu, id}]; ok {
					continue
				}

				me := navigation.MenuEntry{Identifier: id,
					Name:   p.LinkTitle(),
					Weight: p.Weight(),
					URL:    p.RelPermalink()}
				flat[twoD{sectionPagesMenu, me.KeyName()}] = &me
			}
		}
	}

	// Add menu entries provided by pages
	for _, p := range s.workAllPages {
		for name, me := range p.pageMenus.menus() {
			if _, ok := flat[twoD{name, me.KeyName()}]; ok {
				s.SendError(p.wrapError(errors.Errorf("duplicate menu entry with identifier %q in menu %q", me.KeyName(), name)))
				continue
			}
			flat[twoD{name, me.KeyName()}] = me
		}
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
			flat[twoD{p.MenuName, p.EntryName}] = &navigation.MenuEntry{Name: p.EntryName, URL: ""}
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

// get any lanaguagecode to prefix the target file path with.
func (s *Site) getLanguageTargetPathLang(alwaysInSubDir bool) string {
	if s.h.IsMultihost() {
		return s.Language().Lang
	}

	return s.getLanguagePermalinkLang(alwaysInSubDir)
}

// get any lanaguagecode to prefix the relative permalink with.
func (s *Site) getLanguagePermalinkLang(alwaysInSubDir bool) string {

	if !s.Info.IsMultiLingual() || s.h.IsMultihost() {
		return ""
	}

	if alwaysInSubDir {
		return s.Language().Lang
	}

	isDefault := s.Language().Lang == s.multilingual().DefaultLang.Lang

	if !isDefault || s.Info.defaultContentLanguageInSubdir {
		return s.Language().Lang
	}

	return ""
}

func (s *Site) getTaxonomyKey(key string) string {
	if s.PathSpec.DisablePathToLower {
		return s.PathSpec.MakePath(key)
	}
	return strings.ToLower(s.PathSpec.MakePath(key))
}

func (s *Site) assembleTaxonomies() error {
	s.Taxonomies = make(TaxonomyList)
	taxonomies := s.siteCfg.taxonomiesConfig
	for _, plural := range taxonomies {
		s.Taxonomies[plural] = make(Taxonomy)
	}

	s.taxonomyNodes = &taxonomyNodeInfos{
		m:      make(map[string]*taxonomyNodeInfo),
		getKey: s.getTaxonomyKey,
	}

	s.Log.INFO.Printf("found taxonomies: %#v\n", taxonomies)

	for singular, plural := range taxonomies {
		parent := s.taxonomyNodes.GetOrCreate(plural, "")
		parent.singular = singular

		addTaxonomy := func(plural, term string, weight int, p page.Page) {
			key := s.getTaxonomyKey(term)

			n := s.taxonomyNodes.GetOrCreate(plural, term)
			n.parent = parent

			w := page.NewWeightedPage(weight, p, n.owner)

			s.Taxonomies[plural].add(key, w)

			n.UpdateFromPage(w.Page)
			parent.UpdateFromPage(w.Page)
		}

		for _, p := range s.workAllPages {
			vals := getParam(p, plural, false)

			w := getParamToLower(p, plural+"_weight")
			weight, err := cast.ToIntE(w)
			if err != nil {
				s.Log.ERROR.Printf("Unable to convert taxonomy weight %#v to int for %q", w, p.pathOrTitle())
				// weight will equal zero, so let the flow continue
			}

			if vals != nil {
				if v, ok := vals.([]string); ok {
					for _, idx := range v {
						addTaxonomy(plural, idx, weight, p)
					}
				} else if v, ok := vals.(string); ok {
					addTaxonomy(plural, v, weight, p)
				} else {
					s.Log.ERROR.Printf("Invalid %s in %q\n", plural, p.pathOrTitle())
				}
			}
		}

		for k := range s.Taxonomies[plural] {
			s.Taxonomies[plural][k].Sort()
		}
	}

	return nil
}

// Prepare site for a new full build.
func (s *Site) resetBuildState() {
	s.relatedDocsHandler = s.relatedDocsHandler.Clone()
	s.PageCollections = newPageCollectionsFromPages(s.rawAllPages)
	s.buildStats = &buildStats{}
	s.init.Reset()

	for _, p := range s.rawAllPages {
		p.pagePages = &pagePages{}
		p.subSections = page.Pages{}
		p.parent = nil
		p.Scratcher = maps.NewScratcher()
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
func (s *SiteInfo) GetPage(ref ...string) (page.Page, error) {
	p, err := s.s.getPageOldVersion(ref...)

	if p == nil {
		// The nil struct has meaning in some situations, mostly to avoid breaking
		// existing sites doing $nilpage.IsDescendant($p), which will always return
		// false.
		p = page.NilPage
	}

	return p, err
}

func (s *Site) permalink(link string) string {
	return s.PathSpec.PermalinkForBaseURL(link, s.PathSpec.BaseURL.String())

}

func (s *Site) renderAndWriteXML(statCounter *uint64, name string, targetPath string, d interface{}, layouts ...string) error {
	s.Log.DEBUG.Printf("Render XML for %q to %q", name, targetPath)
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	if err := s.renderForLayouts(name, "", d, renderBuffer, layouts...); err != nil {
		return err
	}

	var path string
	if s.Info.relativeURLs {
		path = helpers.GetDottedRelativePath(targetPath)
	} else {
		s := s.PathSpec.BaseURL.String()
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		path = s
	}

	pd := publisher.Descriptor{
		Src:         renderBuffer,
		TargetPath:  targetPath,
		StatCounter: statCounter,
		// For the minification part of XML,
		// we currently only use the MIME type.
		OutputFormat: output.RSSFormat,
		AbsURLPath:   path,
	}

	return s.publisher.Publish(pd)

}

func (s *Site) renderAndWritePage(statCounter *uint64, name string, targetPath string, p *pageState, layouts ...string) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	of := p.outputFormat()

	if err := s.renderForLayouts(p.Kind(), of.Name, p, renderBuffer, layouts...); err != nil {
		return err
	}

	if renderBuffer.Len() == 0 {
		return nil
	}

	isHTML := of.IsHTML
	isRSS := of.Name == "RSS"

	var path string

	if s.Info.relativeURLs {
		path = helpers.GetDottedRelativePath(targetPath)
	} else if isRSS || s.Info.canonifyURLs {
		url := s.PathSpec.BaseURL.String()
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		path = url
	}

	pd := publisher.Descriptor{
		Src:          renderBuffer,
		TargetPath:   targetPath,
		StatCounter:  statCounter,
		OutputFormat: p.outputFormat(),
	}

	if isRSS {
		// Always canonify URLs in RSS
		pd.AbsURLPath = path
	} else if isHTML {
		if s.Info.relativeURLs || s.Info.canonifyURLs {
			pd.AbsURLPath = path
		}

		if s.running() && s.Cfg.GetBool("watch") && !s.Cfg.GetBool("disableLiveReload") {
			pd.LiveReloadPort = s.Cfg.GetInt("liveReloadPort")
		}

		// For performance reasons we only inject the Hugo generator tag on the home page.
		if p.IsHome() {
			pd.AddHugoGeneratorTag = !s.Cfg.GetBool("disableHugoGeneratorInject")
		}

	}

	return s.publisher.Publish(pd)
}

var infoOnMissingLayout = map[string]bool{
	// The 404 layout is very much optional in Hugo, but we do look for it.
	"404": true,
}

func (s *Site) renderForLayouts(name, outputFormat string, d interface{}, w io.Writer, layouts ...string) (err error) {
	templ := s.findFirstTemplate(layouts...)
	if templ == nil {
		log := s.Log.WARN
		if infoOnMissingLayout[name] {
			log = s.Log.INFO
		}

		errMsg := "You should create a template file which matches Hugo Layouts Lookup Rules for this combination."
		var args []interface{}
		msg := "found no layout file for"
		if outputFormat != "" {
			msg += " %q"
			args = append(args, outputFormat)
		}
		if name != "" {
			msg += " for %q"
			args = append(args, name)
		}

		msg += ": " + errMsg

		log.Printf(msg, args...)

		return nil
	}

	if err = templ.Execute(w, d); err != nil {
		return _errors.Wrapf(err, "render of %q failed", name)
	}
	return
}

func (s *Site) findFirstTemplate(layouts ...string) tpl.Template {
	for _, layout := range layouts {
		if templ, found := s.Tmpl.Lookup(layout); found {
			return templ
		}
	}
	return nil
}

func (s *Site) publish(statCounter *uint64, path string, r io.Reader) (err error) {
	s.PathSpec.ProcessingStats.Incr(statCounter)

	return helpers.WriteToDisk(filepath.Clean(path), r, s.BaseFs.PublishFs)
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
	if len(sections) == 0 || len(s.siteCfg.taxonomiesConfig) == 0 {
		return page.KindSection
	}

	sectionPath := path.Join(sections...)

	for _, plural := range s.siteCfg.taxonomiesConfig {
		if plural == sectionPath {
			return page.KindTaxonomyTerm
		}

		if strings.HasPrefix(sectionPath, plural) {
			return page.KindTaxonomy
		}

	}

	return page.KindSection
}

func (s *Site) newTaxonomyPage(title string, sections ...string) *pageState {
	p, err := newPageFromMeta(&pageMeta{
		title:    title,
		s:        s,
		kind:     page.KindTaxonomy,
		sections: sections,
	})

	if err != nil {
		panic(err)
	}

	return p

}

func (s *Site) newPage(kind string, sections ...string) *pageState {
	p, err := newPageFromMeta(&pageMeta{
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
	return shouldBuild(s.BuildFuture, s.BuildExpired,
		s.BuildDrafts, p.Draft(), p.PublishDate(), p.ExpiryDate())
}

func shouldBuild(buildFuture bool, buildExpired bool, buildDrafts bool, Draft bool,
	publishDate time.Time, expiryDate time.Time) bool {
	if !(buildDrafts || !Draft) {
		return false
	}
	if !buildFuture && !publishDate.IsZero() && publishDate.After(time.Now()) {
		return false
	}
	if !buildExpired && !expiryDate.IsZero() && expiryDate.Before(time.Now()) {
		return false
	}
	return true
}
