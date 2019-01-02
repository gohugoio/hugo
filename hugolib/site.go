// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/common/text"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/publisher"
	_errors "github.com/pkg/errors"

	"github.com/gohugoio/hugo/langs"

	src "github.com/gohugoio/hugo/source"

	"golang.org/x/sync/errgroup"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/markbates/inflect"

	"github.com/fsnotify/fsnotify"
	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib/pagemeta"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/nitro"
	"github.com/spf13/viper"
)

// used to indicate if run as a test.
var testMode bool

var defaultTimer *nitro.B

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
	owner *HugoSites

	*PageCollections

	Taxonomies TaxonomyList

	// Plural is what we get in the folder, so keep track of this mapping
	// to get the singular form from that value.
	taxonomiesPluralSingular map[string]string

	// This is temporary, see https://github.com/gohugoio/hugo/issues/2835
	// Maps 	"actors-gerard-depardieu" to "Gérard Depardieu" when preserveTaxonomyNames
	// is set.
	taxonomiesOrigKey map[string]string

	Sections Taxonomy
	Info     SiteInfo
	Menus    Menus
	timer    *nitro.B

	layoutHandler *output.LayoutHandler

	draftCount   int
	futureCount  int
	expiredCount int

	Data     map[string]interface{}
	Language *langs.Language

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

	siteConfig SiteConfig

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

	relatedDocsHandler *relatedDocsHandler
	siteRefLinker
	// Set in some tests
	shortcodePlaceholderFunc func() string

	publisher publisher.Publisher
}

type siteRenderingContext struct {
	output.Format
}

func (s *Site) initRenderFormats() {
	formatSet := make(map[string]bool)
	formats := output.Formats{}
	for _, p := range s.Pages {
		for _, f := range p.outputFormats {
			if !formatSet[f.Name] {
				formats = append(formats, f)
				formatSet[f.Name] = true
			}
		}
	}

	sort.Sort(formats)
	s.renderFormats = formats
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
		relatedDocsHandler:     newSearchIndexHandler(s.relatedDocsHandler.cfg),
		siteRefLinker:          s.siteRefLinker,
		outputFormats:          s.outputFormats,
		rc:                     s.rc,
		outputFormatsConfig:    s.outputFormatsConfig,
		frontmatterHandler:     s.frontmatterHandler,
		mediaTypesConfig:       s.mediaTypesConfig,
		Language:               s.Language,
		owner:                  s.owner,
		publisher:              s.publisher,
		siteConfig:             s.siteConfig,
		enableInlineShortcodes: s.enableInlineShortcodes,
		PageCollections:        newPageCollections()}

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

	var relatedContentConfig related.Config

	if cfg.Language.IsSet("related") {
		relatedContentConfig, err = related.DecodeConfig(cfg.Language.Get("related"))
		if err != nil {
			return nil, err
		}
	} else {
		relatedContentConfig = related.DefaultConfig
		taxonomies := cfg.Language.GetStringMapString("taxonomies")
		if _, found := taxonomies["tag"]; found {
			relatedContentConfig.Add(related.IndexConfig{Name: "tags", Weight: 80})
		}
	}

	titleFunc := helpers.GetTitleFunc(cfg.Language.GetString("titleCaseStyle"))

	frontMatterHandler, err := pagemeta.NewFrontmatterHandler(cfg.Logger, cfg.Cfg)
	if err != nil {
		return nil, err
	}

	s := &Site{
		PageCollections:        c,
		layoutHandler:          output.NewLayoutHandler(),
		Language:               cfg.Language,
		disabledKinds:          disabledKinds,
		titleFunc:              titleFunc,
		relatedDocsHandler:     newSearchIndexHandler(relatedContentConfig),
		outputFormats:          outputFormats,
		rc:                     &siteRenderingContext{output.HTMLFormat},
		outputFormatsConfig:    siteOutputFormatsConfig,
		mediaTypesConfig:       siteMediaTypesConfig,
		frontmatterHandler:     frontMatterHandler,
		enableInlineShortcodes: cfg.Language.GetBool("enableInlineShortcodes"),
	}

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

type SiteInfos []*SiteInfo

// First is a convenience method to get the first Site, i.e. the main language.
func (s SiteInfos) First() *SiteInfo {
	if len(s) == 0 {
		return nil
	}
	return s[0]
}

type SiteInfo struct {
	Taxonomies TaxonomyList
	Authors    AuthorList
	Social     SiteSocial
	*PageCollections
	Menus                          *Menus
	hugoInfo                       hugo.Info
	Title                          string
	RSSLink                        string
	Author                         map[string]interface{}
	LanguageCode                   string
	Copyright                      string
	LastChange                     time.Time
	Permalinks                     PermalinkOverrides
	Params                         map[string]interface{}
	BuildDrafts                    bool
	canonifyURLs                   bool
	relativeURLs                   bool
	uglyURLs                       func(p *Page) bool
	preserveTaxonomyNames          bool
	Data                           *map[string]interface{}
	owner                          *HugoSites
	s                              *Site
	language                       *langs.Language
	LanguagePrefix                 string
	Languages                      langs.Languages
	defaultContentLanguageInSubdir bool
	sectionPagesMenu               string
}

func (s *SiteInfo) Language() *langs.Language {
	return s.language
}

func (s *SiteInfo) Config() SiteConfig {
	return s.s.siteConfig
}

func (s *SiteInfo) Hugo() hugo.Info {
	return s.hugoInfo
}

// Sites is a convenience method to get all the Hugo sites/languages configured.
func (s *SiteInfo) Sites() SiteInfos {
	return s.s.owner.siteInfos()
}
func (s *SiteInfo) String() string {
	return fmt.Sprintf("Site(%q)", s.Title)
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
	return s.Params[keyStr], nil
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

func (s siteRefLinker) logNotFound(ref, what string, p *Page, position text.Position) {
	if position.IsValid() {
		s.errorLogger.Printf("[%s] REF_NOT_FOUND: Ref %q: %s: %s", s.s.Lang(), ref, position.String(), what)
	} else if p == nil {
		s.errorLogger.Printf("[%s] REF_NOT_FOUND: Ref %q: %s", s.s.Lang(), ref, what)
	} else {
		s.errorLogger.Printf("[%s] REF_NOT_FOUND: Ref %q from page %q: %s", s.s.Lang(), ref, p.pathOrTitle(), what)
	}
}

func (s *siteRefLinker) refLink(ref string, source interface{}, relative bool, outputFormat string) (string, error) {

	var page *Page
	switch v := source.(type) {
	case *Page:
		page = v
	case pageContainer:
		page = v.page()
	}

	var refURL *url.URL
	var err error

	ref = filepath.ToSlash(ref)

	refURL, err = url.Parse(ref)

	if err != nil {
		return s.notFoundURL, err
	}

	var target *Page
	var link string

	if refURL.Path != "" {
		target, err := s.s.getPageNew(page, refURL.Path)
		var pos text.Position
		if err != nil || target == nil {
			if p, ok := source.(text.Positioner); ok {
				pos = p.Position()

			}
		}

		if err != nil {
			s.logNotFound(refURL.Path, err.Error(), page, pos)
			return s.notFoundURL, nil
		}

		if target == nil {
			s.logNotFound(refURL.Path, "page not found", page, pos)
			return s.notFoundURL, nil
		}

		var permalinker Permalinker = target

		if outputFormat != "" {
			o := target.OutputFormats().Get(outputFormat)

			if o == nil {
				s.logNotFound(refURL.Path, fmt.Sprintf("output format %q", outputFormat), page, pos)
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
		link = link + "#" + refURL.Fragment

		if refURL.Path != "" && target != nil && !target.getRenderingConfig().PlainIDAnchors {
			link = link + ":" + target.UniqueID()
		} else if page != nil && !page.getRenderingConfig().PlainIDAnchors {
			link = link + ":" + page.UniqueID()
		}
	}

	return link, nil
}

// Ref will give an absolute URL to ref in the given Page.
func (s *SiteInfo) Ref(ref string, page *Page, options ...string) (string, error) {
	// Remove in Hugo 0.53
	helpers.Deprecated("Site", ".Ref", "Use .Site.GetPage", false)
	outputFormat := ""
	if len(options) > 0 {
		outputFormat = options[0]
	}

	return s.s.refLink(ref, page, false, outputFormat)
}

// RelRef will give an relative URL to ref in the given Page.
func (s *SiteInfo) RelRef(ref string, page *Page, options ...string) (string, error) {
	// Remove in Hugo 0.53
	helpers.Deprecated("Site", ".RelRef", "Use .Site.GetPage", false)
	outputFormat := ""
	if len(options) > 0 {
		outputFormat = options[0]
	}

	return s.s.refLink(ref, page, true, outputFormat)
}

func (s *Site) running() bool {
	return s.owner != nil && s.owner.running
}

func (s *Site) multilingual() *Multilingual {
	return s.owner.multilingual
}

func init() {
	defaultTimer = nitro.Initalize()
}

func (s *Site) timerStep(step string) {
	if s.timer == nil {
		s.timer = defaultTimer
	}
	s.timer.Step(step)
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

	h := s.owner

	s.timerStep("initialize rebuild")

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
				clearIsInnerShortcodeCache()
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
	for _, s := range s.owner.Sites {
		s.ResourceSpec.ResourceCache.DeletePartitions(cachePartitions...)
	}

	if len(tmplChanged) > 0 || len(i18nChanged) > 0 {
		sites := s.owner.Sites
		first := sites[0]

		// TOD(bep) globals clean
		if err := first.Deps.LoadResources(); err != nil {
			return whatChanged{}, err
		}

		for i := 1; i < len(sites); i++ {
			site := sites[i]
			var err error
			depsCfg := deps.DepsCfg{
				Language:      site.Language,
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

		s.timerStep("template prep")
	}

	if len(dataChanged) > 0 {
		if err := s.readDataFromSourceFS(); err != nil {
			return whatChanged{}, err
		}
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
			contentFilesChanged = append(contentFilesChanged, p.File.Filename())
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

func (s *Site) loadData(fs afero.Fs) (err error) {
	spec := src.NewSourceSpec(s.PathSpec, fs)
	fileSystem := spec.NewFilesystem("")
	s.Data = make(map[string]interface{})
	for _, r := range fileSystem.Files() {
		if err := s.handleDataFile(r); err != nil {
			return err
		}
	}

	return
}

func (s *Site) errWithFileContext(err error, f source.File) error {
	rfi, ok := f.FileInfo().(hugofs.RealFilenameInfo)
	if !ok {
		return err
	}

	realFilename := rfi.RealFilename()

	err, _ = herrors.WithFileContextForFile(
		err,
		realFilename,
		realFilename,
		s.SourceSpec.Fs.Source,
		herrors.SimpleLineMatcher)

	return err
}

func (s *Site) handleDataFile(r source.ReadableFile) error {
	var current map[string]interface{}

	f, err := r.Open()
	if err != nil {
		return _errors.Wrapf(err, "Failed to open data file %q:", r.LogicalName())
	}
	defer f.Close()

	// Crawl in data tree to insert data
	current = s.Data
	keyParts := strings.Split(r.Dir(), helpers.FilePathSeparator)
	// The first path element is the virtual folder (typically theme name), which is
	// not part of the key.
	if len(keyParts) > 1 {
		for _, key := range keyParts[1:] {
			if key != "" {
				if _, ok := current[key]; !ok {
					current[key] = make(map[string]interface{})
				}
				current = current[key].(map[string]interface{})
			}
		}
	}

	data, err := s.readData(r)
	if err != nil {
		return s.errWithFileContext(err, r)
	}

	if data == nil {
		return nil
	}

	// filepath.Walk walks the files in lexical order, '/' comes before '.'
	// this warning could happen if
	// 1. A theme uses the same key; the main data folder wins
	// 2. A sub folder uses the same key: the sub folder wins
	higherPrecedentData := current[r.BaseFileName()]

	switch data.(type) {
	case nil:
		// hear the crickets?

	case map[string]interface{}:

		switch higherPrecedentData.(type) {
		case nil:
			current[r.BaseFileName()] = data
		case map[string]interface{}:
			// merge maps: insert entries from data for keys that
			// don't already exist in higherPrecedentData
			higherPrecedentMap := higherPrecedentData.(map[string]interface{})
			for key, value := range data.(map[string]interface{}) {
				if _, exists := higherPrecedentMap[key]; exists {
					s.Log.WARN.Printf("Data for key '%s' in path '%s' is overridden by higher precedence data already in the data tree", key, r.Path())
				} else {
					higherPrecedentMap[key] = value
				}
			}
		default:
			// can't merge: higherPrecedentData is not a map
			s.Log.WARN.Printf("The %T data from '%s' overridden by "+
				"higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
		}

	case []interface{}:
		if higherPrecedentData == nil {
			current[r.BaseFileName()] = data
		} else {
			// we don't merge array data
			s.Log.WARN.Printf("The %T data from '%s' overridden by "+
				"higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
		}

	default:
		s.Log.ERROR.Printf("unexpected data type %T in file %s", data, r.LogicalName())
	}

	return nil
}

func (s *Site) readData(f source.ReadableFile) (interface{}, error) {
	file, err := f.Open()
	if err != nil {
		return nil, _errors.Wrap(err, "readData: failed to open data file")
	}
	defer file.Close()
	content := helpers.ReaderToBytes(file)

	format := metadecoders.FormatFromString(f.Extension())
	return metadecoders.Default.Unmarshal(content, format)
}

func (s *Site) readDataFromSourceFS() error {
	err := s.loadData(s.PathSpec.BaseFs.Data.Fs)
	s.timerStep("load data")
	return err
}

func (s *Site) process(config BuildCfg) (err error) {
	if err = s.initialize(); err != nil {
		return
	}
	s.timerStep("initialize")

	if err = s.readDataFromSourceFS(); err != nil {
		return
	}

	s.timerStep("load i18n")

	if err := s.readAndProcessContent(); err != nil {
		return err
	}
	s.timerStep("read and convert pages from source")

	return err

}

func (s *Site) setupSitePages() {
	var siteLastChange time.Time

	for i, page := range s.RegularPages {
		if i > 0 {
			page.NextPage = s.RegularPages[i-1]
		}

		if i < len(s.RegularPages)-1 {
			page.PrevPage = s.RegularPages[i+1]
		}

		// Determine Site.Info.LastChange
		// Note that the logic to determine which date to use for Lastmod
		// is already applied, so this is *the* date to use.
		// We cannot just pick the last page in the default sort, because
		// that may not be ordered by date.
		if page.Lastmod.After(siteLastChange) {
			siteLastChange = page.Lastmod
		}
	}

	s.Info.LastChange = siteLastChange
}

func (s *Site) render(config *BuildCfg, outFormatIdx int) (err error) {
	// Clear the global page cache.
	spc.clear()

	if outFormatIdx == 0 {
		if err = s.preparePages(); err != nil {
			return
		}
		s.timerStep("prepare pages")

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
			s.timerStep("render and write aliases")
		}

	}

	if err = s.renderPages(config); err != nil {
		return
	}

	s.timerStep("render and write pages")

	// TODO(bep) render consider this, ref. render404 etc.
	if outFormatIdx > 0 {
		return
	}

	if err = s.renderSitemap(); err != nil {
		return
	}
	s.timerStep("render and write Sitemap")

	if err = s.renderRobotsTXT(); err != nil {
		return
	}
	s.timerStep("render and write robots.txt")

	if err = s.render404(); err != nil {
		return
	}
	s.timerStep("render and write 404")

	return
}

func (s *Site) Initialise() (err error) {
	return s.initialize()
}

func (s *Site) initialize() (err error) {
	s.Menus = Menus{}

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
	sitemapDefault := parseSitemap(s.s.Cfg.GetStringMap("sitemap"))
	p := s.HomeAbsURL()
	if !strings.HasSuffix(p, "/") {
		p += "/"
	}
	p += sitemapDefault.Filename
	return p
}

func (s *Site) initializeSiteInfo() error {
	var (
		lang      = s.Language
		languages langs.Languages
	)

	if s.owner != nil && s.owner.multilingual != nil {
		languages = s.owner.multilingual.Languages
	}

	params := lang.Params()

	permalinks := make(PermalinkOverrides)
	for k, v := range s.Cfg.GetStringMapString("permalinks") {
		permalinks[k] = pathPattern(v)
	}

	defaultContentInSubDir := s.Cfg.GetBool("defaultContentLanguageInSubdir")
	defaultContentLanguage := s.Cfg.GetString("defaultContentLanguage")

	languagePrefix := ""
	if s.multilingualEnabled() && (defaultContentInSubDir || lang.Lang != defaultContentLanguage) {
		languagePrefix = "/" + lang.Lang
	}

	var uglyURLs = func(p *Page) bool {
		return false
	}

	v := s.Cfg.Get("uglyURLs")
	if v != nil {
		switch vv := v.(type) {
		case bool:
			uglyURLs = func(p *Page) bool {
				return vv
			}
		case string:
			// Is what be get from CLI (--uglyURLs)
			vvv := cast.ToBool(vv)
			uglyURLs = func(p *Page) bool {
				return vvv
			}
		default:
			m := cast.ToStringMapBool(v)
			uglyURLs = func(p *Page) bool {
				return m[p.Section()]
			}
		}
	}

	s.Info = SiteInfo{
		Title:                          lang.GetString("title"),
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
		preserveTaxonomyNames:          lang.GetBool("preserveTaxonomyNames"),
		PageCollections:                s.PageCollections,
		Menus:                          &s.Menus,
		Params:                         params,
		Permalinks:                     permalinks,
		Data:                           &s.Data,
		owner:                          s.owner,
		s:                              s,
		hugoInfo:                       hugo.NewInfo(s.Cfg.GetString("environment")),
		// TODO(bep) make this Menu and similar into delegate methods on SiteInfo
		Taxonomies: s.Taxonomies,
	}

	rssOutputFormat, found := s.outputFormats[KindHome].GetByName(output.RSSFormat.Name)

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

func (s *Site) absContentDir() string {
	return s.PathSpec.AbsPathify(s.PathSpec.ContentDir)
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

func (c *contentCaptureResultHandler) handleCopyFiles(files ...pathLangFile) {
	for _, proc := range c.contentProcessors {
		proc.processAssets(files)
	}
}

func (s *Site) readAndProcessContent(filenames ...string) error {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	defaultContentLanguage := s.SourceSpec.DefaultContentLanguage

	contentProcessors := make(map[string]*siteContentProcessor)
	var defaultContentProcessor *siteContentProcessor
	sites := s.owner.langSite()
	for k, v := range sites {
		if v.Language.Disabled {
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
		bundleMap = s.owner.ContentChanges
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

func (s *Site) buildSiteMeta() (err error) {
	defer s.timerStep("build Site meta")

	if len(s.Pages) == 0 {
		return
	}

	s.assembleTaxonomies()

	for _, p := range s.AllPages {
		// this depends on taxonomies
		p.setValuesForKind(s)
	}

	return
}

func (s *Site) getMenusFromConfig() Menus {

	ret := Menus{}

	if menus := s.Language.GetStringMap("menus"); menus != nil {
		for name, menu := range menus {
			m, err := cast.ToSliceE(menu)
			if err != nil {
				s.Log.ERROR.Printf("unable to process menus in site config\n")
				s.Log.ERROR.Println(err)
			} else {
				for _, entry := range m {
					s.Log.DEBUG.Printf("found menu: %q, in site config\n", name)

					menuEntry := MenuEntry{Menu: name}
					ime, err := cast.ToStringMapE(entry)
					if err != nil {
						s.Log.ERROR.Printf("unable to process menus in site config\n")
						s.Log.ERROR.Println(err)
					}

					menuEntry.marshallMap(ime)
					menuEntry.URL = s.Info.createNodeMenuEntryURL(menuEntry.URL)

					if ret[name] == nil {
						ret[name] = &Menu{}
					}
					*ret[name] = ret[name].add(&menuEntry)
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
	s.Menus = Menus{}

	type twoD struct {
		MenuName, EntryName string
	}
	flat := map[twoD]*MenuEntry{}
	children := map[twoD]Menu{}

	// add menu entries from config to flat hash
	menuConfig := s.getMenusFromConfig()
	for name, menu := range menuConfig {
		for _, me := range *menu {
			flat[twoD{name, me.KeyName()}] = me
		}
	}

	sectionPagesMenu := s.Info.sectionPagesMenu
	pages := s.Pages

	if sectionPagesMenu != "" {
		for _, p := range pages {
			if p.Kind == KindSection {
				// From Hugo 0.22 we have nested sections, but until we get a
				// feel of how that would work in this setting, let us keep
				// this menu for the top level only.
				id := p.Section()
				if _, ok := flat[twoD{sectionPagesMenu, id}]; ok {
					continue
				}

				me := MenuEntry{Identifier: id,
					Name:   p.LinkTitle(),
					Weight: p.Weight,
					URL:    p.RelPermalink()}
				flat[twoD{sectionPagesMenu, me.KeyName()}] = &me
			}
		}
	}

	// Add menu entries provided by pages
	for _, p := range pages {
		for name, me := range p.Menus() {
			if _, ok := flat[twoD{name, me.KeyName()}]; ok {
				s.SendError(p.errWithFileContext(errors.Errorf("duplicate menu entry with identifier %q in menu %q", me.KeyName(), name)))
				continue
			}
			flat[twoD{name, me.KeyName()}] = me
		}
	}

	// Create Children Menus First
	for _, e := range flat {
		if e.Parent != "" {
			children[twoD{e.Menu, e.Parent}] = children[twoD{e.Menu, e.Parent}].add(e)
		}
	}

	// Placing Children in Parents (in flat)
	for p, childmenu := range children {
		_, ok := flat[twoD{p.MenuName, p.EntryName}]
		if !ok {
			// if parent does not exist, create one without a URL
			flat[twoD{p.MenuName, p.EntryName}] = &MenuEntry{Name: p.EntryName, URL: ""}
		}
		flat[twoD{p.MenuName, p.EntryName}].Children = childmenu
	}

	// Assembling Top Level of Tree
	for menu, e := range flat {
		if e.Parent == "" {
			_, ok := s.Menus[menu.MenuName]
			if !ok {
				s.Menus[menu.MenuName] = &Menu{}
			}
			*s.Menus[menu.MenuName] = s.Menus[menu.MenuName].add(e)
		}
	}
}

func (s *Site) getTaxonomyKey(key string) string {
	if s.Info.preserveTaxonomyNames {
		// Keep as is
		return key
	}
	return s.PathSpec.MakePathSanitized(key)
}

// We need to create the top level taxonomy early in the build process
// to be able to determine the page Kind correctly.
func (s *Site) createTaxonomiesEntries() {
	s.Taxonomies = make(TaxonomyList)
	taxonomies := s.Language.GetStringMapString("taxonomies")
	for _, plural := range taxonomies {
		s.Taxonomies[plural] = make(Taxonomy)
	}
}

func (s *Site) assembleTaxonomies() {
	s.taxonomiesPluralSingular = make(map[string]string)
	s.taxonomiesOrigKey = make(map[string]string)

	taxonomies := s.Language.GetStringMapString("taxonomies")

	s.Log.INFO.Printf("found taxonomies: %#v\n", taxonomies)

	for singular, plural := range taxonomies {
		s.taxonomiesPluralSingular[plural] = singular

		for _, p := range s.Pages {
			vals := p.getParam(plural, !s.Info.preserveTaxonomyNames)

			w := p.getParamToLower(plural + "_weight")
			weight, err := cast.ToIntE(w)
			if err != nil {
				s.Log.ERROR.Printf("Unable to convert taxonomy weight %#v to int for %s", w, p.File.Path())
				// weight will equal zero, so let the flow continue
			}

			if vals != nil {
				if v, ok := vals.([]string); ok {
					for _, idx := range v {
						x := WeightedPage{weight, p}
						s.Taxonomies[plural].add(s.getTaxonomyKey(idx), x)
						if s.Info.preserveTaxonomyNames {
							// Need to track the original
							s.taxonomiesOrigKey[fmt.Sprintf("%s-%s", plural, s.PathSpec.MakePathSanitized(idx))] = idx
						}
					}
				} else if v, ok := vals.(string); ok {
					x := WeightedPage{weight, p}
					s.Taxonomies[plural].add(s.getTaxonomyKey(v), x)
					if s.Info.preserveTaxonomyNames {
						// Need to track the original
						s.taxonomiesOrigKey[fmt.Sprintf("%s-%s", plural, s.PathSpec.MakePathSanitized(v))] = v
					}
				} else {
					s.Log.ERROR.Printf("Invalid %s in %s\n", plural, p.File.Path())
				}
			}
		}
		for k := range s.Taxonomies[plural] {
			s.Taxonomies[plural][k].Sort()
		}
	}

	s.Info.Taxonomies = s.Taxonomies
}

// Prepare site for a new full build.
func (s *Site) resetBuildState() {

	s.relatedDocsHandler = newSearchIndexHandler(s.relatedDocsHandler.cfg)
	s.PageCollections = newPageCollectionsFromPages(s.rawAllPages)
	// TODO(bep) get rid of this double
	s.Info.PageCollections = s.PageCollections

	s.draftCount = 0
	s.futureCount = 0

	s.expiredCount = 0

	for _, p := range s.rawAllPages {
		p.subSections = Pages{}
		p.parent = nil
		p.scratch = maps.NewScratch()
		p.mainPageOutput = nil
	}
}

func (s *Site) layouts(p *PageOutput) ([]string, error) {
	return s.layoutHandler.For(p.layoutDescriptor, p.outputFormat)
}

func (s *Site) preparePages() error {
	var errors []error

	for _, p := range s.Pages {
		if err := p.prepareLayouts(); err != nil {
			errors = append(errors, err)
		}
		if err := p.prepareData(s); err != nil {
			errors = append(errors, err)
		}
	}

	return s.owner.pickOneAndLogTheRest(errors)
}

func (s *Site) errorCollator(results <-chan error, errs chan<- error) {
	var errors []error
	for e := range results {
		errors = append(errors, e)
	}

	errs <- s.owner.pickOneAndLogTheRest(errors)

	close(errs)
}

// GetPage looks up a page of a given type for the given ref.
// In Hugo <= 0.44 you had to add Page Kind (section, home) etc. as the first
// argument and then either a unix styled path (with or without a leading slash))
// or path elements separated.
// When we now remove the Kind from this API, we need to make the transition as painless
// as possible for existing sites. Most sites will use {{ .Site.GetPage "section" "my/section" }},
// i.e. 2 arguments, so we test for that.
func (s *SiteInfo) GetPage(ref ...string) (*Page, error) {
	return s.getPageOldVersion(ref...)
}

func (s *Site) permalinkForOutputFormat(link string, f output.Format) (string, error) {
	var (
		baseURL string
		err     error
	)

	if f.Protocol != "" {
		baseURL, err = s.PathSpec.BaseURL.WithProtocol(f.Protocol)
		if err != nil {
			return "", err
		}
	} else {
		baseURL = s.PathSpec.BaseURL.String()
	}
	return s.PathSpec.PermalinkForBaseURL(link, baseURL), nil
}

func (s *Site) permalink(link string) string {
	return s.PathSpec.PermalinkForBaseURL(link, s.PathSpec.BaseURL.String())

}

func (s *Site) renderAndWriteXML(statCounter *uint64, name string, targetPath string, d interface{}, layouts ...string) error {
	s.Log.DEBUG.Printf("Render XML for %q to %q", name, targetPath)
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)
	renderBuffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n")

	if err := s.renderForLayouts(name, d, renderBuffer, layouts...); err != nil {
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

func (s *Site) renderAndWritePage(statCounter *uint64, name string, targetPath string, p *PageOutput, layouts ...string) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	if err := s.renderForLayouts(p.Kind, p, renderBuffer, layouts...); err != nil {

		return err
	}

	if renderBuffer.Len() == 0 {
		return nil
	}

	isHTML := p.outputFormat.IsHTML

	var path string

	if s.Info.relativeURLs {
		path = helpers.GetDottedRelativePath(targetPath)
	} else if s.Info.canonifyURLs {
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
		OutputFormat: p.outputFormat,
	}

	if isHTML {
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

func (s *Site) renderForLayouts(name string, d interface{}, w io.Writer, layouts ...string) (err error) {
	var templ tpl.Template

	templ = s.findFirstTemplate(layouts...)
	if templ == nil {
		log := s.Log.WARN
		if infoOnMissingLayout[name] {
			log = s.Log.INFO
		}

		if p, ok := d.(*PageOutput); ok {
			log.Printf("Found no layout for %q, language %q, output format %q: create a template below /layouts with one of these filenames: %s\n", name, s.Language.Lang, p.outputFormat.Name, layoutsLogFormat(layouts))
		} else {
			log.Printf("Found no layout for %q, language %q: create a template below /layouts with one of these filenames: %s\n", name, s.Language.Lang, layoutsLogFormat(layouts))
		}
		return nil
	}

	if err = templ.Execute(w, d); err != nil {
		return _errors.Wrapf(err, "render of %q failed", name)
	}
	return
}

func layoutsLogFormat(layouts []string) string {
	var filtered []string
	for _, l := range layouts {
		// This is  a technical prefix of no interest to the user.
		lt := strings.TrimPrefix(l, "_text/")
		// We have this in the lookup path for historical reasons.
		lt = strings.TrimPrefix(lt, "page/")
		filtered = append(filtered, lt)
	}

	filtered = helpers.UniqueStrings(filtered)
	return strings.Join(filtered, ", ")
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

func getGoMaxProcs() int {
	if gmp := os.Getenv("GOMAXPROCS"); gmp != "" {
		if p, err := strconv.Atoi(gmp); err != nil {
			return p
		}
	}
	return 1
}

func (s *Site) newNodePage(typ string, sections ...string) *Page {
	p := &Page{
		language:        s.Language,
		pageInit:        &pageInit{},
		pageContentInit: &pageContentInit{},
		Kind:            typ,
		File:            &source.FileInfo{},
		data:            make(map[string]interface{}),
		Site:            &s.Info,
		sections:        sections,
		s:               s}

	p.outputFormats = p.s.outputFormats[p.Kind]

	return p

}

func (s *Site) newHomePage() *Page {
	p := s.newNodePage(KindHome)
	p.title = s.Info.Title
	pages := Pages{}
	p.data["Pages"] = pages
	p.Pages = pages
	return p
}

func (s *Site) newTaxonomyPage(plural, key string) *Page {

	p := s.newNodePage(KindTaxonomy, plural, key)

	if s.Info.preserveTaxonomyNames {
		p.title = key
	} else {
		p.title = strings.Replace(s.titleFunc(key), "-", " ", -1)
	}

	return p
}

func (s *Site) newSectionPage(name string) *Page {
	p := s.newNodePage(KindSection, name)

	sectionName := helpers.FirstUpper(name)
	if s.Cfg.GetBool("pluralizeListTitles") {
		p.title = inflect.Pluralize(sectionName)
	} else {
		p.title = sectionName
	}
	return p
}

func (s *Site) newTaxonomyTermsPage(plural string) *Page {
	p := s.newNodePage(KindTaxonomyTerm, plural)
	p.title = s.titleFunc(plural)
	return p
}
