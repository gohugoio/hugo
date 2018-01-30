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
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gohugoio/hugo/resource"

	"golang.org/x/sync/errgroup"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/media"

	"github.com/markbates/inflect"
	"golang.org/x/net/context"

	"github.com/fsnotify/fsnotify"
	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"
	"github.com/gohugoio/hugo/transform"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/nitro"
	"github.com/spf13/viper"
)

var _ = transform.AbsURL

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
	Language *helpers.Language

	disabledKinds map[string]bool

	// Output formats defined in site config per Page Kind, or some defaults
	// if not set.
	// Output formats defined in Page front matter will override these.
	outputFormats map[string]output.Formats

	// All the output formats and media types available for this site.
	// These values will be merged from the Hugo defaults, the site config and,
	// finally, the language settings.
	outputFormatsConfig output.Formats
	mediaTypesConfig    media.Types

	// We render each site for all the relevant output formats in serial with
	// this rendering context pointing to the current one.
	rc *siteRenderingContext

	// The output formats that we need to render this site in. This slice
	// will be fixed once set.
	// This will be the union of Site.Pages' outputFormats.
	// This slice will be sorted.
	renderFormats output.Formats

	// Logger etc.
	*deps.Deps   `json:"-"`
	resourceSpec *resource.Spec

	// The func used to title case titles.
	titleFunc func(s string) string

	relatedDocsHandler *relatedDocsHandler
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
		layoutHandler:       output.NewLayoutHandler(s.PathSpec.ThemeSet()),
		disabledKinds:       s.disabledKinds,
		titleFunc:           s.titleFunc,
		relatedDocsHandler:  newSearchIndexHandler(s.relatedDocsHandler.cfg),
		outputFormats:       s.outputFormats,
		outputFormatsConfig: s.outputFormatsConfig,
		mediaTypesConfig:    s.mediaTypesConfig,
		resourceSpec:        s.resourceSpec,
		Language:            s.Language,
		owner:               s.owner,
		PageCollections:     newPageCollections()}
}

// newSite creates a new site with the given configuration.
func newSite(cfg deps.DepsCfg) (*Site, error) {
	c := newPageCollections()

	if cfg.Language == nil {
		cfg.Language = helpers.NewDefaultLanguage(cfg.Cfg)
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

	s := &Site{
		PageCollections:     c,
		layoutHandler:       output.NewLayoutHandler(cfg.Cfg.GetString("themesDir") != ""),
		Language:            cfg.Language,
		disabledKinds:       disabledKinds,
		titleFunc:           titleFunc,
		relatedDocsHandler:  newSearchIndexHandler(relatedContentConfig),
		outputFormats:       outputFormats,
		outputFormatsConfig: siteOutputFormatsConfig,
		mediaTypesConfig:    siteMediaTypesConfig,
	}

	s.Info = newSiteInfo(siteBuilderCfg{s: s, pageCollections: c, language: s.Language})

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

	if err = applyDepsIfNeeded(cfg, s); err != nil {
		return nil, err
	}

	return s, nil
}

// NewSiteDefaultLang creates a new site in the default language.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
func NewSiteDefaultLang(withTemplate ...func(templ tpl.TemplateHandler) error) (*Site, error) {
	v := viper.New()
	if err := loadDefaultSettingsFor(v); err != nil {
		return nil, err
	}
	return newSiteForLang(helpers.NewDefaultLanguage(v), withTemplate...)
}

// NewEnglishSite creates a new site in English language.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
func NewEnglishSite(withTemplate ...func(templ tpl.TemplateHandler) error) (*Site, error) {
	v := viper.New()
	if err := loadDefaultSettingsFor(v); err != nil {
		return nil, err
	}
	return newSiteForLang(helpers.NewLanguage("en", v), withTemplate...)
}

// newSiteForLang creates a new site in the given language.
func newSiteForLang(lang *helpers.Language, withTemplate ...func(templ tpl.TemplateHandler) error) (*Site, error) {
	withTemplates := func(templ tpl.TemplateHandler) error {
		for _, wt := range withTemplate {
			if err := wt(templ); err != nil {
				return err
			}
		}
		return nil
	}

	cfg := deps.DepsCfg{WithTemplate: withTemplates, Language: lang, Cfg: lang}

	return NewSiteForCfg(cfg)

}

// NewSiteForCfg creates a new site for the given configuration.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
func NewSiteForCfg(cfg deps.DepsCfg) (*Site, error) {
	s, err := newSite(cfg)

	if err != nil {
		return nil, err
	}

	if err := applyDepsIfNeeded(cfg, s); err != nil {
		return nil, err
	}
	return s, nil
}

type SiteInfo struct {
	Taxonomies TaxonomyList
	Authors    AuthorList
	Social     SiteSocial
	*PageCollections
	Menus                 *Menus
	Hugo                  *HugoInfo
	Title                 string
	RSSLink               string
	Author                map[string]interface{}
	LanguageCode          string
	DisqusShortname       string
	GoogleAnalytics       string
	Copyright             string
	LastChange            time.Time
	Permalinks            PermalinkOverrides
	Params                map[string]interface{}
	BuildDrafts           bool
	canonifyURLs          bool
	relativeURLs          bool
	uglyURLs              func(p *Page) bool
	preserveTaxonomyNames bool
	Data                  *map[string]interface{}

	owner                          *HugoSites
	s                              *Site
	multilingual                   *Multilingual
	Language                       *helpers.Language
	LanguagePrefix                 string
	Languages                      helpers.Languages
	defaultContentLanguageInSubdir bool
	sectionPagesMenu               string
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

// Used in tests.

type siteBuilderCfg struct {
	language        *helpers.Language
	s               *Site
	pageCollections *PageCollections
}

// TODO(bep) get rid of this
func newSiteInfo(cfg siteBuilderCfg) SiteInfo {
	return SiteInfo{
		s:               cfg.s,
		multilingual:    newMultiLingualForLanguage(cfg.language),
		PageCollections: cfg.pageCollections,
		Params:          make(map[string]interface{}),
		uglyURLs: func(p *Page) bool {
			return false
		},
	}
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

func (s *SiteInfo) refLink(ref string, page *Page, relative bool, outputFormat string) (string, error) {
	var refURL *url.URL
	var err error

	ref = filepath.ToSlash(ref)
	ref = strings.TrimPrefix(ref, "/")

	refURL, err = url.Parse(ref)

	if err != nil {
		return "", err
	}

	var target *Page
	var link string

	if refURL.Path != "" {
		target := s.getPage(KindPage, refURL.Path)

		if target == nil {
			return "", fmt.Errorf("No page found with path or logical name \"%s\".\n", refURL.Path)
		}

		var permalinker Permalinker = target

		if outputFormat != "" {
			o := target.OutputFormats().Get(outputFormat)

			if o == nil {
				return "", fmt.Errorf("Output format %q not found for page %q", outputFormat, refURL.Path)
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
	outputFormat := ""
	if len(options) > 0 {
		outputFormat = options[0]
	}

	return s.refLink(ref, page, false, outputFormat)
}

// RelRef will give an relative URL to ref in the given Page.
func (s *SiteInfo) RelRef(ref string, page *Page, options ...string) (string, error) {
	outputFormat := ""
	if len(options) > 0 {
		outputFormat = options[0]
	}

	return s.refLink(ref, page, true, outputFormat)
}

func (s *Site) running() bool {
	return s.owner.running
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
		// The last one will win if there are any duplicates.
		_ = mime.AddExtensionType("."+mt.Suffix, mt.Type()+"; charset=utf-8")
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

	for _, ev := range events {
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

	if len(tmplChanged) > 0 || len(i18nChanged) > 0 {
		sites := s.owner.Sites
		first := sites[0]

		// TOD(bep) globals clean
		if err := first.Deps.LoadResources(); err != nil {
			s.Log.ERROR.Println(err)
		}

		s.TemplateHandler().PrintErrors()

		for i := 1; i < len(sites); i++ {
			site := sites[i]
			var err error
			site.Deps, err = first.Deps.ForLanguage(site.Language)
			if err != nil {
				return whatChanged{}, err
			}
		}

		s.timerStep("template prep")
	}

	if len(dataChanged) > 0 {
		if err := s.readDataFromSourceFS(); err != nil {
			s.Log.ERROR.Println(err)
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
		if removed && isContentFile(ev.Name) {
			path, _ := helpers.GetRelativePath(ev.Name, s.getContentDir(ev.Name))

			h.removePageByPath(path)
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
		source: len(sourceChanged) > 0,
		other:  len(tmplChanged) > 0 || len(i18nChanged) > 0 || len(dataChanged) > 0,
		files:  sourceFilesChanged,
	}

	return changed, nil

}

func (s *Site) loadData(sourceDirs []string) (err error) {
	s.Log.DEBUG.Printf("Load Data from %d source(s)", len(sourceDirs))
	s.Data = make(map[string]interface{})
	for _, sourceDir := range sourceDirs {
		fs := s.SourceSpec.NewFilesystem(sourceDir)
		for _, r := range fs.Files() {
			if err := s.handleDataFile(r); err != nil {
				return err
			}
		}
	}

	return
}

func (s *Site) handleDataFile(r source.ReadableFile) error {
	var current map[string]interface{}

	f, err := r.Open()
	if err != nil {
		return fmt.Errorf("Failed to open data file %q: %s", r.LogicalName(), err)
	}
	defer f.Close()

	// Crawl in data tree to insert data
	current = s.Data
	for _, key := range strings.Split(r.Dir(), helpers.FilePathSeparator) {
		if key != "" {
			if _, ok := current[key]; !ok {
				current[key] = make(map[string]interface{})
			}
			current = current[key].(map[string]interface{})
		}
	}

	data, err := s.readData(r)
	if err != nil {
		s.Log.WARN.Printf("Failed to read data from %s: %s", filepath.Join(r.Path(), r.LogicalName()), err)
		return nil
	}

	if data == nil {
		return nil
	}

	// Copy content from current to data when needed
	if _, ok := current[r.BaseFileName()]; ok {
		data := data.(map[string]interface{})

		for key, value := range current[r.BaseFileName()].(map[string]interface{}) {
			if _, override := data[key]; override {
				// filepath.Walk walks the files in lexical order, '/' comes before '.'
				// this warning could happen if
				// 1. A theme uses the same key; the main data folder wins
				// 2. A sub folder uses the same key: the sub folder wins
				s.Log.WARN.Printf("Data for key '%s' in path '%s' is overridden in subfolder", key, r.Path())
			}
			data[key] = value
		}
	}

	// Insert data
	current[r.BaseFileName()] = data

	return nil
}

func (s *Site) readData(f source.ReadableFile) (interface{}, error) {
	file, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content := helpers.ReaderToBytes(file)

	switch f.Extension() {
	case "yaml", "yml":
		return parser.HandleYAMLMetaData(content)
	case "json":
		return parser.HandleJSONMetaData(content)
	case "toml":
		return parser.HandleTOMLMetaData(content)
	default:
		return nil, fmt.Errorf("Data not supported for extension '%s'", f.Extension())
	}
}

func (s *Site) readDataFromSourceFS() error {
	var dataSourceDirs []string

	// have to be last - duplicate keys in earlier entries will win
	themeDataDir, err := s.PathSpec.GetThemeDataDirPath()
	if err == nil {
		dataSourceDirs = []string{s.absDataDir(), themeDataDir}
	} else {
		dataSourceDirs = []string{s.absDataDir()}

	}

	err = s.loadData(dataSourceDirs)
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
		if i < len(s.RegularPages)-1 {
			page.Next = s.RegularPages[i+1]
		}

		if i > 0 {
			page.Prev = s.RegularPages[i-1]
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
	defer s.initializeSiteInfo()
	s.Menus = Menus{}

	if err = s.checkDirectories(); err != nil {
		return err
	}

	return
}

// HomeAbsURL is a convenience method giving the absolute URL to the home page.
func (s *SiteInfo) HomeAbsURL() string {
	base := ""
	if s.IsMultiLingual() {
		base = s.Language.Lang
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

func (s *Site) initializeSiteInfo() {
	var (
		lang      = s.Language
		languages helpers.Languages
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

	var multilingual *Multilingual
	if s.owner != nil {
		multilingual = s.owner.multilingual
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
		DisqusShortname:                lang.GetString("disqusShortname"),
		multilingual:                   multilingual,
		Language:                       lang,
		LanguagePrefix:                 languagePrefix,
		Languages:                      languages,
		defaultContentLanguageInSubdir: defaultContentInSubDir,
		sectionPagesMenu:               lang.GetString("sectionPagesMenu"),
		GoogleAnalytics:                lang.GetString("googleAnalytics"),
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
	}

	rssOutputFormat, found := s.outputFormats[KindHome].GetByName(output.RSSFormat.Name)

	if found {
		s.Info.RSSLink = s.permalink(rssOutputFormat.BaseFilename())
	}
}

func (s *Site) dataDir() string {
	return s.Cfg.GetString("dataDir")
}

func (s *Site) absDataDir() string {
	return s.PathSpec.AbsPathify(s.dataDir())
}

func (s *Site) i18nDir() string {
	return s.Cfg.GetString("i18nDir")
}

func (s *Site) absI18nDir() string {
	return s.PathSpec.AbsPathify(s.i18nDir())
}

func (s *Site) isI18nEvent(e fsnotify.Event) bool {
	if s.getI18nDir(e.Name) != "" {
		return true
	}
	return s.getThemeI18nDir(e.Name) != ""
}

func (s *Site) getI18nDir(path string) string {
	return s.getRealDir(s.absI18nDir(), path)
}

func (s *Site) getThemeI18nDir(path string) string {
	if !s.PathSpec.ThemeSet() {
		return ""
	}
	return s.getRealDir(filepath.Join(s.PathSpec.GetThemeDir(), s.i18nDir()), path)
}

func (s *Site) isDataDirEvent(e fsnotify.Event) bool {
	if s.getDataDir(e.Name) != "" {
		return true
	}
	return s.getThemeDataDir(e.Name) != ""
}

func (s *Site) getDataDir(path string) string {
	return s.getRealDir(s.absDataDir(), path)
}

func (s *Site) getThemeDataDir(path string) string {
	if !s.PathSpec.ThemeSet() {
		return ""
	}
	return s.getRealDir(filepath.Join(s.PathSpec.GetThemeDir(), s.dataDir()), path)
}

func (s *Site) layoutDir() string {
	return s.Cfg.GetString("layoutDir")
}

func (s *Site) isLayoutDirEvent(e fsnotify.Event) bool {
	if s.getLayoutDir(e.Name) != "" {
		return true
	}
	return s.getThemeLayoutDir(e.Name) != ""
}

func (s *Site) getLayoutDir(path string) string {
	return s.getRealDir(s.PathSpec.GetLayoutDirPath(), path)
}

func (s *Site) getThemeLayoutDir(path string) string {
	if !s.PathSpec.ThemeSet() {
		return ""
	}
	return s.getRealDir(filepath.Join(s.PathSpec.GetThemeDir(), s.layoutDir()), path)
}

func (s *Site) absContentDir() string {
	return s.PathSpec.AbsPathify(s.PathSpec.ContentDir())
}

func (s *Site) isContentDirEvent(e fsnotify.Event) bool {
	return s.getContentDir(e.Name) != ""
}

func (s *Site) getContentDir(path string) string {
	return s.getRealDir(s.absContentDir(), path)
}

// getRealDir gets the base path of the given path, also handling the case where
// base is a symlinked folder.
func (s *Site) getRealDir(base, path string) string {

	if strings.HasPrefix(path, base) {
		return base
	}

	realDir, err := helpers.GetRealPath(s.Fs.Source, base)

	if err != nil {
		if !os.IsNotExist(err) {
			s.Log.ERROR.Printf("Failed to get real path for %s: %s", path, err)
		}
		return ""
	}

	if strings.HasPrefix(path, realDir) {
		return realDir
	}

	return ""
}

func (s *Site) absPublishDir() string {
	return s.PathSpec.AbsPathify(s.Cfg.GetString("publishDir"))
}

func (s *Site) checkDirectories() (err error) {
	if b, _ := helpers.DirExists(s.absContentDir(), s.Fs.Source); !b {
		return errors.New("No source directory found, expecting to find it at " + s.absContentDir())
	}
	return
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
		proc.fileSinglesChan <- fi
	}
}
func (c *contentCaptureResultHandler) handleBundles(d *bundleDirs) {
	for _, b := range d.bundles {
		proc := c.getContentProcessor(b.fi.Lang())
		proc.fileBundlesChan <- b
	}
}

func (c *contentCaptureResultHandler) handleCopyFiles(filenames ...string) {
	for _, proc := range c.contentProcessors {
		proc.fileAssetsChan <- filenames
	}
}

func (s *Site) readAndProcessContent(filenames ...string) error {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	sourceSpec := source.NewSourceSpec(s.owner.Cfg, s.Fs)
	baseDir := s.absContentDir()
	defaultContentLanguage := s.SourceSpec.DefaultContentLanguage

	contentProcessors := make(map[string]*siteContentProcessor)
	var defaultContentProcessor *siteContentProcessor
	sites := s.owner.langSite()
	for k, v := range sites {
		proc := newSiteContentProcessor(baseDir, len(filenames) > 0, v)
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

	if s.running() {
		// Need to track changes.
		bundleMap = s.owner.ContentChanges
		handler = &captureResultHandlerChain{handlers: []captureBundlesHandler{mainHandler, bundleMap}}

	} else {
		handler = mainHandler
	}

	c := newCapturer(s.Log, sourceSpec, handler, bundleMap, baseDir, filenames...)

	if err := c.capture(); err != nil {
		return err
	}

	for _, proc := range contentProcessors {
		proc.closeInput()
	}

	return g.Wait()
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

	if menus := s.Language.GetStringMap("menu"); menus != nil {
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
				s.Log.ERROR.Printf("Two or more menu items have the same name/identifier in Menu %q: %q.\nRename or set an unique identifier.\n", name, me.KeyName())
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
			weight := p.getParamToLower(plural + "_weight")
			if weight == nil {
				weight = 0
			}
			if vals != nil {
				if v, ok := vals.([]string); ok {
					for _, idx := range v {
						x := WeightedPage{weight.(int), p}
						s.Taxonomies[plural].add(s.getTaxonomyKey(idx), x)
						if s.Info.preserveTaxonomyNames {
							// Need to track the original
							s.taxonomiesOrigKey[fmt.Sprintf("%s-%s", plural, s.PathSpec.MakePathSanitized(idx))] = idx
						}
					}
				} else if v, ok := vals.(string); ok {
					x := WeightedPage{weight.(int), p}
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
		p.scratch = newScratch()
		p.subSections = Pages{}
		p.parent = nil
	}
}

func (s *Site) kindFromSections(sections []string) string {
	if len(sections) == 0 {
		return KindSection
	}

	if _, isTaxonomy := s.Taxonomies[sections[0]]; isTaxonomy {
		if len(sections) == 1 {
			return KindTaxonomyTerm
		}
		return KindTaxonomy
	}
	return KindSection
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

	if len(errors) != 0 {
		return fmt.Errorf("Prepare pages failed: %.100q…", errors)
	}

	return nil
}

func errorCollator(results <-chan error, errs chan<- error) {
	errMsgs := []string{}
	for err := range results {
		if err != nil {
			errMsgs = append(errMsgs, err.Error())
		}
	}
	if len(errMsgs) == 0 {
		errs <- nil
	} else {
		errs <- errors.New(strings.Join(errMsgs, "\n"))
	}
	close(errs)
}

func (s *Site) appendThemeTemplates(in []string) []string {
	if !s.PathSpec.ThemeSet() {
		return in
	}

	out := []string{}
	// First place all non internal templates
	for _, t := range in {
		if !strings.HasPrefix(t, "_internal/") {
			out = append(out, t)
		}
	}

	// Then place theme templates with the same names
	for _, t := range in {
		if !strings.HasPrefix(t, "_internal/") {
			out = append(out, "theme/"+t)
		}
	}

	// Lastly place internal templates
	for _, t := range in {
		if strings.HasPrefix(t, "_internal/") {
			out = append(out, t)
		}
	}
	return out

}

// GetPage looks up a page of a given type in the path given.
//    {{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}
//
// This will return nil when no page could be found, and will return the
// first page found if the key is ambigous.
func (s *SiteInfo) GetPage(typ string, path ...string) (*Page, error) {
	return s.getPage(typ, path...), nil
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

func (s *Site) renderAndWriteXML(statCounter *uint64, name string, dest string, d interface{}, layouts ...string) error {
	s.Log.DEBUG.Printf("Render XML for %q to %q", name, dest)
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)
	renderBuffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n")

	if err := s.renderForLayouts(name, d, renderBuffer, layouts...); err != nil {
		helpers.DistinctWarnLog.Println(err)
		return nil
	}

	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)

	var path []byte
	if s.Info.relativeURLs {
		path = []byte(helpers.GetDottedRelativePath(dest))
	} else {
		s := s.PathSpec.BaseURL.String()
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		path = []byte(s)
	}
	transformer := transform.NewChain(transform.AbsURLInXML)
	if err := transformer.Apply(outBuffer, renderBuffer, path); err != nil {
		helpers.DistinctErrorLog.Println(err)
		return nil
	}

	return s.publish(statCounter, dest, outBuffer)

}

func (s *Site) renderAndWritePage(statCounter *uint64, name string, dest string, p *PageOutput, layouts ...string) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	if err := s.renderForLayouts(p.Kind, p, renderBuffer, layouts...); err != nil {
		helpers.DistinctWarnLog.Println(err)
		return nil
	}

	if renderBuffer.Len() == 0 {
		return nil
	}

	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)

	transformLinks := transform.NewEmptyTransforms()

	isHTML := p.outputFormat.IsHTML

	if isHTML {
		if s.Info.relativeURLs || s.Info.canonifyURLs {
			transformLinks = append(transformLinks, transform.AbsURL)
		}

		if s.running() && s.Cfg.GetBool("watch") && !s.Cfg.GetBool("disableLiveReload") {
			transformLinks = append(transformLinks, transform.LiveReloadInject(s.Cfg.GetInt("liveReloadPort")))
		}

		// For performance reasons we only inject the Hugo generator tag on the home page.
		if p.IsHome() {
			if !s.Cfg.GetBool("disableHugoGeneratorInject") {
				transformLinks = append(transformLinks, transform.HugoGeneratorInject)
			}
		}
	}

	var path []byte

	if s.Info.relativeURLs {
		path = []byte(helpers.GetDottedRelativePath(dest))
	} else if s.Info.canonifyURLs {
		url := s.PathSpec.BaseURL.String()
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		path = []byte(url)
	}

	transformer := transform.NewChain(transformLinks...)
	if err := transformer.Apply(outBuffer, renderBuffer, path); err != nil {
		helpers.DistinctErrorLog.Println(err)
		return nil
	}

	return s.publish(statCounter, dest, outBuffer)
}

func (s *Site) renderForLayouts(name string, d interface{}, w io.Writer, layouts ...string) (err error) {
	var templ tpl.Template

	defer func() {
		if r := recover(); r != nil {
			templName := ""
			if templ != nil {
				templName = templ.Name()
			}
			helpers.DistinctErrorLog.Printf("Failed to render %q: %s", templName, r)
			// TOD(bep) we really need to fix this. Also see below.
			if !s.running() && !testMode {
				os.Exit(-1)
			}
		}
	}()

	templ = s.findFirstTemplate(layouts...)
	if templ == nil {
		return fmt.Errorf("[%s] Unable to locate layout for %q: %s\n", s.Language.Lang, name, layouts)
	}

	if err = templ.Execute(w, d); err != nil {
		// Behavior here should be dependent on if running in server or watch mode.
		if p, ok := d.(*PageOutput); ok {
			if p.File != nil {
				helpers.DistinctErrorLog.Printf("Error while rendering %q in %q: %s", name, p.File.Dir(), err)
			} else {
				helpers.DistinctErrorLog.Printf("Error while rendering %q: %s", name, err)
			}
		} else {
			helpers.DistinctErrorLog.Printf("Error while rendering %q: %s", name, err)
		}
		if !s.running() && !testMode {
			// TODO(bep) check if this can be propagated
			os.Exit(-1)
		} else if testMode {
			return
		}
	}

	return
}

func (s *Site) findFirstTemplate(layouts ...string) tpl.Template {
	for _, layout := range layouts {
		if templ := s.Tmpl.Lookup(layout); templ != nil {
			return templ
		}
	}
	return nil
}

func (s *Site) publish(statCounter *uint64, path string, r io.Reader) (err error) {
	s.PathSpec.ProcessingStats.Incr(statCounter)

	path = filepath.Join(s.absPublishDir(), path)

	return helpers.WriteToDisk(path, r, s.Fs.Destination)
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
		language: s.Language,
		pageInit: &pageInit{},
		Kind:     typ,
		Source:   Source{File: &source.FileInfo{}},
		Data:     make(map[string]interface{}),
		Site:     &s.Info,
		sections: sections,
		s:        s}

	p.outputFormats = p.s.outputFormats[p.Kind]

	return p

}

func (s *Site) newHomePage() *Page {
	p := s.newNodePage(KindHome)
	p.title = s.Info.Title
	pages := Pages{}
	p.Data["Pages"] = pages
	p.Pages = pages
	return p
}

func (s *Site) newTaxonomyPage(plural, key string) *Page {

	p := s.newNodePage(KindTaxonomy, plural, key)

	if s.Info.preserveTaxonomyNames {
		// Keep (mostly) as is in the title
		// We make the first character upper case, mostly because
		// it is easier to reason about in the tests.
		p.title = helpers.FirstUpper(key)
		key = s.PathSpec.MakePathSanitized(key)
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
