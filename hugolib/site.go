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
	"sync"
	"time"

	"github.com/spf13/hugo/config"

	"github.com/spf13/hugo/media"

	"github.com/bep/inflect"

	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/output"
	"github.com/spf13/hugo/parser"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/hugo/transform"
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

	Files      []*source.File
	Taxonomies TaxonomyList

	// Plural is what we get in the folder, so keep track of this mapping
	// to get the singular form from that value.
	taxonomiesPluralSingular map[string]string

	// This is temporary, see https://github.com/spf13/hugo/issues/2835
	// Maps 	"actors-gerard-depardieu" to "Gérard Depardieu" when preserveTaxonomyNames
	// is set.
	taxonomiesOrigKey map[string]string

	Source   source.Input
	Sections Taxonomy
	Info     SiteInfo
	Menus    Menus
	timer    *nitro.B

	layoutHandler *output.LayoutHandler

	draftCount   int
	futureCount  int
	expiredCount int
	Data         map[string]interface{}
	Language     *helpers.Language

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
	*deps.Deps `json:"-"`

	siteStats *siteStats
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

type siteStats struct {
	pageCount        int
	pageCountRegular int
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
		outputFormats:       s.outputFormats,
		outputFormatsConfig: s.outputFormatsConfig,
		mediaTypesConfig:    s.mediaTypesConfig,
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

	s := &Site{
		PageCollections:     c,
		layoutHandler:       output.NewLayoutHandler(cfg.Cfg.GetString("themesDir") != ""),
		Language:            cfg.Language,
		disabledKinds:       disabledKinds,
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
	loadDefaultSettingsFor(v)
	return newSiteForLang(helpers.NewDefaultLanguage(v), withTemplate...)
}

// NewEnglishSite creates a new site in English language.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
func NewEnglishSite(withTemplate ...func(templ tpl.TemplateHandler) error) (*Site, error) {
	v := viper.New()
	loadDefaultSettingsFor(v)
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
	// atomic requires 64-bit alignment for struct field access
	// According to the docs, " The first word in a global variable or in an
	// allocated struct or slice can be relied upon to be 64-bit aligned."
	// Moving paginationPageCount to the top of this struct didn't do the
	// magic, maybe due to the way SiteInfo is embedded.
	// Adding the 4 byte padding below does the trick.
	_                   [4]byte
	paginationPageCount uint64

	Taxonomies TaxonomyList
	Authors    AuthorList
	Social     SiteSocial
	*PageCollections
	Files                 *[]*source.File
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
	uglyURLs              bool
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

// Used in tests.

type siteBuilderCfg struct {
	language        *helpers.Language
	s               *Site
	pageCollections *PageCollections
	baseURL         string
}

// TODO(bep) get rid of this
func newSiteInfo(cfg siteBuilderCfg) SiteInfo {
	return SiteInfo{
		s:               cfg.s,
		multilingual:    newMultiLingualForLanguage(cfg.language),
		PageCollections: cfg.pageCollections,
		Params:          make(map[string]interface{}),
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

// SourceRelativeLink attempts to convert any source page relative links (like [../another.md]) into absolute links
func (s *SiteInfo) SourceRelativeLink(ref string, currentPage *Page) (string, error) {
	var refURL *url.URL
	var err error

	refURL, err = url.Parse(strings.TrimPrefix(ref, currentPage.getRenderingConfig().SourceRelativeLinksProjectFolder))
	if err != nil {
		return "", err
	}

	if refURL.Scheme != "" {
		// Not a relative source level path
		return ref, nil
	}

	var target *Page
	var link string

	if refURL.Path != "" {
		refPath := filepath.Clean(filepath.FromSlash(refURL.Path))

		if strings.IndexRune(refPath, os.PathSeparator) == 0 { // filepath.IsAbs fails to me.
			refPath = refPath[1:]
		} else {
			if currentPage != nil {
				refPath = filepath.Join(currentPage.Source.Dir(), refURL.Path)
			}
		}

		for _, page := range s.AllRegularPages {
			if page.Source.Path() == refPath {
				target = page
				break
			}
		}
		// need to exhaust the test, then try with the others :/
		// if the refPath doesn't end in a filename with extension `.md`, then try with `.md` , and then `/index.md`
		mdPath := strings.TrimSuffix(refPath, string(os.PathSeparator)) + ".md"
		for _, page := range s.AllRegularPages {
			if page.Source.Path() == mdPath {
				target = page
				break
			}
		}
		indexPath := filepath.Join(refPath, "index.md")
		for _, page := range s.AllRegularPages {
			if page.Source.Path() == indexPath {
				target = page
				break
			}
		}

		if target == nil {
			return "", fmt.Errorf("No page found for \"%s\" on page \"%s\".\n", ref, currentPage.Source.Path())
		}

		link = target.RelPermalink()

	}

	if refURL.Fragment != "" {
		link = link + "#" + refURL.Fragment

		if refURL.Path != "" && target != nil && !target.getRenderingConfig().PlainIDAnchors {
			link = link + ":" + target.UniqueID()
		} else if currentPage != nil && !currentPage.getRenderingConfig().PlainIDAnchors {
			link = link + ":" + currentPage.UniqueID()
		}
	}

	return link, nil
}

// SourceRelativeLinkFile attempts to convert any non-md source relative links (like [../another.gif]) into absolute links
func (s *SiteInfo) SourceRelativeLinkFile(ref string, currentPage *Page) (string, error) {
	var refURL *url.URL
	var err error

	refURL, err = url.Parse(strings.TrimPrefix(ref, currentPage.getRenderingConfig().SourceRelativeLinksProjectFolder))
	if err != nil {
		return "", err
	}

	if refURL.Scheme != "" {
		// Not a relative source level path
		return ref, nil
	}

	var target *source.File
	var link string

	if refURL.Path != "" {
		refPath := filepath.Clean(filepath.FromSlash(refURL.Path))

		if strings.IndexRune(refPath, os.PathSeparator) == 0 { // filepath.IsAbs fails to me.
			refPath = refPath[1:]
		} else {
			if currentPage != nil {
				refPath = filepath.Join(currentPage.Source.Dir(), refURL.Path)
			}
		}

		for _, file := range *s.Files {
			if file.Path() == refPath {
				target = file
				break
			}
		}

		if target == nil {
			return "", fmt.Errorf("No file found for \"%s\" on page \"%s\".\n", ref, currentPage.Source.Path())
		}

		link = target.Path()
		return "/" + filepath.ToSlash(link), nil
	}

	return "", fmt.Errorf("failed to find a file to match \"%s\" on page \"%s\"", ref, currentPage.Source.Path())
}

func (s *SiteInfo) addToPaginationPageCount(cnt uint64) {
	atomic.AddUint64(&s.paginationPageCount, cnt)
}

type runmode struct {
	Watching bool
}

func (s *Site) running() bool {
	return s.owner.runMode.Watching
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
}

// RegisterMediaTypes will register the Site's media types in the mime
// package, so it will behave correctly with Hugo's built-in server.
func (s *Site) RegisterMediaTypes() {
	for _, mt := range s.mediaTypesConfig {
		// The last one will win if there are any duplicates.
		_ = mime.AddExtensionType("."+mt.Suffix, mt.Type()+"; charset=utf-8")
	}
}

// reBuild partially rebuilds a site given the filesystem events.
// It returns whetever the content source was changed.
func (s *Site) reProcess(events []fsnotify.Event) (whatChanged, error) {
	s.Log.DEBUG.Printf("Rebuild for events %q", events)

	s.timerStep("initialize rebuild")

	// First we need to determine what changed

	sourceChanged := []fsnotify.Event{}
	sourceReallyChanged := []fsnotify.Event{}
	tmplChanged := []fsnotify.Event{}
	dataChanged := []fsnotify.Event{}
	i18nChanged := []fsnotify.Event{}
	shortcodesChanged := make(map[string]bool)
	// prevent spamming the log on changes
	logger := helpers.NewDistinctFeedbackLogger()
	seen := make(map[fsnotify.Event]bool)

	for _, ev := range events {
		// Avoid processing the same event twice.
		if seen[ev] {
			continue
		}
		seen[ev] = true

		if s.isContentDirEvent(ev) {
			logger.Println("Source changed", ev.Name)
			sourceChanged = append(sourceChanged, ev)
		}
		if s.isLayoutDirEvent(ev) {
			logger.Println("Template changed", ev.Name)
			tmplChanged = append(tmplChanged, ev)

			if strings.Contains(ev.Name, "shortcodes") {
				clearIsInnerShortcodeCache()
				shortcode := filepath.Base(ev.Name)
				shortcode = strings.TrimSuffix(shortcode, filepath.Ext(shortcode))
				shortcodesChanged[shortcode] = true
			}
		}
		if s.isDataDirEvent(ev) {
			logger.Println("Data changed", ev.Name)
			dataChanged = append(dataChanged, ev)
		}
		if s.isI18nEvent(ev) {
			logger.Println("i18n changed", ev.Name)
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

	// If a content file changes, we need to reload only it and re-render the entire site.

	// First step is to read the changed files and (re)place them in site.AllPages
	// This includes processing any meta-data for that content

	// The second step is to convert the content into HTML
	// This includes processing any shortcodes that may be present.

	// We do this in parallel... even though it's likely only one file at a time.
	// We need to process the reading prior to the conversion for each file, but
	// we can convert one file while another one is still reading.
	errs := make(chan error, 2)
	readResults := make(chan HandledResult)
	filechan := make(chan *source.File)
	convertResults := make(chan HandledResult)
	pageChan := make(chan *Page)
	fileConvChan := make(chan *source.File)
	coordinator := make(chan bool)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go sourceReader(s, filechan, readResults, wg)
	}

	wg2 := &sync.WaitGroup{}
	wg2.Add(4)
	for i := 0; i < 2; i++ {
		go fileConverter(s, fileConvChan, convertResults, wg2)
		go pageConverter(pageChan, convertResults, wg2)
	}

	sp := source.NewSourceSpec(s.Cfg, s.Fs)
	fs := sp.NewFilesystem("")

	for _, ev := range sourceChanged {
		// The incrementalReadCollator below will also make changes to the site's pages,
		// so we do this first to prevent races.
		if ev.Op&fsnotify.Remove == fsnotify.Remove {
			//remove the file & a create will follow
			path, _ := helpers.GetRelativePath(ev.Name, s.getContentDir(ev.Name))
			s.removePageByPath(path)
			continue
		}

		// Some editors (Vim) sometimes issue only a Rename operation when writing an existing file
		// Sometimes a rename operation means that file has been renamed other times it means
		// it's been updated
		if ev.Op&fsnotify.Rename == fsnotify.Rename {
			// If the file is still on disk, it's only been updated, if it's not, it's been moved
			if ex, err := afero.Exists(s.Fs.Source, ev.Name); !ex || err != nil {
				path, _ := helpers.GetRelativePath(ev.Name, s.getContentDir(ev.Name))
				s.removePageByPath(path)
				continue
			}
		}

		// ignore files shouldn't be proceed
		if fi, err := s.Fs.Source.Stat(ev.Name); err != nil {
			continue
		} else {
			if ok, err := fs.ShouldRead(ev.Name, fi); err != nil || !ok {
				continue
			}
		}

		sourceReallyChanged = append(sourceReallyChanged, ev)
	}

	go incrementalReadCollator(s, readResults, pageChan, fileConvChan, coordinator, errs)
	go converterCollator(convertResults, errs)

	for _, ev := range sourceReallyChanged {

		file, err := s.reReadFile(ev.Name)

		if err != nil {
			s.Log.ERROR.Println("Error reading file", ev.Name, ";", err)
		}

		if file != nil {
			filechan <- file
		}

	}

	for shortcode := range shortcodesChanged {
		// There are certain scenarios that, when a shortcode changes,
		// it isn't sufficient to just rerender the already parsed shortcode.
		// One example is if the user adds a new shortcode to the content file first,
		// and then creates the shortcode on the file system.
		// To handle these scenarios, we must do a full reprocessing of the
		// pages that keeps a reference to the changed shortcode.
		pagesWithShortcode := s.findPagesByShortcode(shortcode)
		for _, p := range pagesWithShortcode {
			p.rendered = false
			pageChan <- p
		}
	}

	// we close the filechan as we have sent everything we want to send to it.
	// this will tell the sourceReaders to stop iterating on that channel
	close(filechan)

	// waiting for the sourceReaders to all finish
	wg.Wait()
	// Now closing readResults as this will tell the incrementalReadCollator to
	// stop iterating over that.
	close(readResults)

	// once readResults is finished it will close coordinator and move along
	<-coordinator
	// allow that routine to finish, then close page & fileconvchan as we've sent
	// everything to them we need to.
	close(pageChan)
	close(fileConvChan)

	wg2.Wait()
	close(convertResults)

	s.timerStep("read & convert pages from source")

	for i := 0; i < 2; i++ {
		err := <-errs
		if err != nil {
			s.Log.ERROR.Println(err)
		}
	}

	changed := whatChanged{
		source: len(sourceChanged) > 0,
		other:  len(tmplChanged) > 0 || len(i18nChanged) > 0 || len(dataChanged) > 0,
	}

	return changed, nil

}

func (s *Site) loadData(sources []source.Input) (err error) {
	s.Log.DEBUG.Printf("Load Data from %d source(s)", len(sources))
	s.Data = make(map[string]interface{})
	var current map[string]interface{}
	for _, currentSource := range sources {
		for _, r := range currentSource.Files() {
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
				continue
			}

			if data == nil {
				continue
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
		}
	}

	return
}

func (s *Site) readData(f *source.File) (interface{}, error) {
	switch f.Extension() {
	case "yaml", "yml":
		return parser.HandleYAMLMetaData(f.Bytes())
	case "json":
		return parser.HandleJSONMetaData(f.Bytes())
	case "toml":
		return parser.HandleTOMLMetaData(f.Bytes())
	default:
		return nil, fmt.Errorf("Data not supported for extension '%s'", f.Extension())
	}
}

func (s *Site) readDataFromSourceFS() error {
	sp := source.NewSourceSpec(s.Cfg, s.Fs)
	dataSources := make([]source.Input, 0, 2)
	dataSources = append(dataSources, sp.NewFilesystem(s.absDataDir()))

	// have to be last - duplicate keys in earlier entries will win
	themeDataDir, err := s.PathSpec.GetThemeDataDirPath()
	if err == nil {
		dataSources = append(dataSources, sp.NewFilesystem(themeDataDir))
	}

	err = s.loadData(dataSources)
	s.timerStep("load data")
	return err
}

func (s *Site) process(config BuildCfg) (err error) {
	s.timerStep("Go initialization")
	if err = s.initialize(); err != nil {
		return
	}
	s.timerStep("initialize")

	if err = s.readDataFromSourceFS(); err != nil {
		return
	}

	s.timerStep("load i18n")
	return s.createPages()

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

func (s *Site) render(outFormatIdx int) (err error) {

	if outFormatIdx == 0 {
		if err = s.preparePages(); err != nil {
			return
		}
		s.timerStep("prepare pages")

		// Aliases must be rendered before pages.
		// Some sites, Hugo docs included, have faulty alias definitions that point
		// to itself or another real page. These will be overwritten in the next
		// step.
		if err = s.renderAliases(); err != nil {
			return
		}
		s.timerStep("render and write aliases")

	}

	if err = s.renderPages(); err != nil {
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

	// May be supplied in tests.
	if s.Source != nil && len(s.Source.Files()) > 0 {
		s.Log.DEBUG.Println("initialize: Source is already set")
		return
	}

	if err = s.checkDirectories(); err != nil {
		return err
	}

	staticDir := s.PathSpec.GetStaticDirPath() + "/"

	sp := source.NewSourceSpec(s.Cfg, s.Fs)
	s.Source = sp.NewFilesystem(s.absContentDir(), staticDir)

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
		uglyURLs:                       s.Cfg.GetBool("uglyURLs"),
		preserveTaxonomyNames:          lang.GetBool("preserveTaxonomyNames"),
		PageCollections:                s.PageCollections,
		Files:                          &s.Files,
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
	return s.PathSpec.AbsPathify(s.Cfg.GetString("contentDir"))
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

// reReadFile resets file to be read from disk again
func (s *Site) reReadFile(absFilePath string) (*source.File, error) {
	s.Log.INFO.Println("rereading", absFilePath)
	var file *source.File

	reader, err := source.NewLazyFileReader(s.Fs.Source, absFilePath)
	if err != nil {
		return nil, err
	}

	sp := source.NewSourceSpec(s.Cfg, s.Fs)
	file, err = sp.NewFileFromAbs(s.getContentDir(absFilePath), absFilePath, reader)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *Site) readPagesFromSource() chan error {
	if s.Source == nil {
		panic(fmt.Sprintf("s.Source not set %s", s.absContentDir()))
	}

	s.Log.DEBUG.Printf("Read %d pages from source", len(s.Source.Files()))

	errs := make(chan error)
	if len(s.Source.Files()) < 1 {
		close(errs)
		return errs
	}

	files := s.Source.Files()
	results := make(chan HandledResult)
	filechan := make(chan *source.File)
	wg := &sync.WaitGroup{}
	numWorkers := getGoMaxProcs() * 4
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go sourceReader(s, filechan, results, wg)
	}

	// we can only have exactly one result collator, since it makes changes that
	// must be synchronized.
	go readCollator(s, results, errs)

	for _, file := range files {
		filechan <- file
	}

	close(filechan)
	wg.Wait()
	close(results)

	return errs
}

func (s *Site) convertSource() chan error {
	errs := make(chan error)
	results := make(chan HandledResult)
	pageChan := make(chan *Page)
	fileConvChan := make(chan *source.File)
	numWorkers := getGoMaxProcs() * 4
	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(2)
		go fileConverter(s, fileConvChan, results, wg)
		go pageConverter(pageChan, results, wg)
	}

	go converterCollator(results, errs)

	for _, p := range s.rawAllPages {
		if p.shouldBuild() {
			pageChan <- p
		}
	}

	for _, f := range s.Files {
		fileConvChan <- f
	}

	close(pageChan)
	close(fileConvChan)
	wg.Wait()
	close(results)

	return errs
}

func (s *Site) createPages() error {
	readErrs := <-s.readPagesFromSource()
	s.timerStep("read pages from source")

	renderErrs := <-s.convertSource()
	s.timerStep("convert source")

	if renderErrs == nil && readErrs == nil {
		return nil
	}
	if renderErrs == nil {
		return readErrs
	}
	if readErrs == nil {
		return renderErrs
	}

	return fmt.Errorf("%s\n%s", readErrs, renderErrs)
}

func sourceReader(s *Site, files <-chan *source.File, results chan<- HandledResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		readSourceFile(s, file, results)
	}
}

func readSourceFile(s *Site, file *source.File, results chan<- HandledResult) {
	h := NewMetaHandler(file.Extension())
	if h != nil {
		h.Read(file, s, results)
	} else {
		s.Log.ERROR.Println("Unsupported File Type", file.Path())
	}
}

func pageConverter(pages <-chan *Page, results HandleResults, wg *sync.WaitGroup) {
	defer wg.Done()
	for page := range pages {
		var h *MetaHandle
		if page.Markup != "" {
			h = NewMetaHandler(page.Markup)
		} else {
			h = NewMetaHandler(page.File.Extension())
		}
		if h != nil {
			// Note that we convert pages from the site's rawAllPages collection
			// Which may contain pages from multiple sites, so we use the Page's site
			// for the conversion.
			h.Convert(page, page.s, results)
		}
	}
}

func fileConverter(s *Site, files <-chan *source.File, results HandleResults, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		h := NewMetaHandler(file.Extension())
		if h != nil {
			h.Convert(file, s, results)
		}
	}
}

func converterCollator(results <-chan HandledResult, errs chan<- error) {
	errMsgs := []string{}
	for r := range results {
		if r.err != nil {
			errMsgs = append(errMsgs, r.err.Error())
			continue
		}
	}
	if len(errMsgs) == 0 {
		errs <- nil
		return
	}
	errs <- fmt.Errorf("Errors rendering pages: %s", strings.Join(errMsgs, "\n"))
}

func (s *Site) replaceFile(sf *source.File) {
	for i, f := range s.Files {
		if f.Path() == sf.Path() {
			s.Files[i] = sf
			return
		}
	}

	// If a match isn't found, then append it
	s.Files = append(s.Files, sf)
}

func incrementalReadCollator(s *Site, results <-chan HandledResult, pageChan chan *Page, fileConvChan chan *source.File, coordinator chan bool, errs chan<- error) {
	errMsgs := []string{}
	for r := range results {
		if r.err != nil {
			errMsgs = append(errMsgs, r.Error())
			continue
		}

		if r.page == nil {
			s.replaceFile(r.file)
			fileConvChan <- r.file
		} else {
			s.replacePage(r.page)
			pageChan <- r.page
		}
	}

	s.rawAllPages.Sort()
	close(coordinator)

	if len(errMsgs) == 0 {
		errs <- nil
		return
	}
	errs <- fmt.Errorf("Errors reading pages: %s", strings.Join(errMsgs, "\n"))
}

func readCollator(s *Site, results <-chan HandledResult, errs chan<- error) {
	if s.PageCollections == nil {
		panic("No page collections")
	}
	errMsgs := []string{}
	for r := range results {
		if r.err != nil {
			errMsgs = append(errMsgs, r.Error())
			continue
		}

		// !page == file
		if r.page == nil {
			s.Files = append(s.Files, r.file)
		} else {
			s.addPage(r.page)
		}
	}

	s.rawAllPages.Sort()
	if len(errMsgs) == 0 {
		errs <- nil
		return
	}
	errs <- fmt.Errorf("Errors reading pages: %s", strings.Join(errMsgs, "\n"))
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
			weight := p.GetParam(plural + "_weight")
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

	s.PageCollections = newPageCollectionsFromPages(s.rawAllPages)
	// TODO(bep) get rid of this double
	s.Info.PageCollections = s.PageCollections

	s.Info.paginationPageCount = 0
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
	if _, isTaxonomy := s.Taxonomies[sections[0]]; isTaxonomy {
		if len(sections) == 1 {
			return KindTaxonomyTerm
		}
		return KindTaxonomy
	}
	return KindSection
}

func (s *Site) layouts(p *PageOutput) ([]string, error) {
	return s.layoutHandler.For(p.layoutDescriptor, "", p.outputFormat)
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

// Stats prints Hugo builds stats to the console.
// This is what you see after a successful hugo build.
func (s *Site) Stats() {

	s.Log.FEEDBACK.Printf("Built site for language %s:\n", s.Language.Lang)
	s.Log.FEEDBACK.Println(s.draftStats())
	s.Log.FEEDBACK.Println(s.futureStats())
	s.Log.FEEDBACK.Println(s.expiredStats())
	s.Log.FEEDBACK.Printf("%d regular pages created\n", s.siteStats.pageCountRegular)
	s.Log.FEEDBACK.Printf("%d other pages created\n", (s.siteStats.pageCount - s.siteStats.pageCountRegular))
	s.Log.FEEDBACK.Printf("%d non-page files copied\n", len(s.Files))
	s.Log.FEEDBACK.Printf("%d paginator pages created\n", s.Info.paginationPageCount)

	if s.isEnabled(KindTaxonomy) {
		taxonomies := s.Language.GetStringMapString("taxonomies")

		for _, pl := range taxonomies {
			s.Log.FEEDBACK.Printf("%d %s created\n", len(s.Taxonomies[pl]), pl)
		}
	}

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
	return s.permalinkForBaseURL(link, baseURL), nil
}

func (s *Site) permalink(link string) string {
	return s.permalinkForBaseURL(link, s.PathSpec.BaseURL.String())

}

func (s *Site) permalinkForBaseURL(link, baseURL string) string {
	link = strings.TrimPrefix(link, "/")
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + link
}

func (s *Site) renderAndWriteXML(name string, dest string, d interface{}, layouts ...string) error {
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
		s := s.Cfg.GetString("baseURL")
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

	return s.publish(dest, outBuffer)

}

func (s *Site) renderAndWritePage(name string, dest string, p *PageOutput, layouts ...string) error {
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
			transformLinks = append(transformLinks, transform.LiveReloadInject(s.Cfg.GetInt("port")))
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
		url := s.Cfg.GetString("baseURL")
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

	return s.publish(dest, outBuffer)
}

func (s *Site) renderForLayouts(name string, d interface{}, w io.Writer, layouts ...string) error {
	templ := s.findFirstTemplate(layouts...)
	if templ == nil {
		return fmt.Errorf("[%s] Unable to locate layout for %q: %s\n", s.Language.Lang, name, layouts)
	}

	if err := templ.Execute(w, d); err != nil {
		// Behavior here should be dependent on if running in server or watch mode.
		helpers.DistinctErrorLog.Printf("Error while rendering %q: %s", name, err)
		if !s.running() && !testMode {
			// TODO(bep) check if this can be propagated
			os.Exit(-1)
		} else if testMode {
			return err
		}
	}

	return nil
}

func (s *Site) findFirstTemplate(layouts ...string) tpl.Template {
	for _, layout := range layouts {
		if templ := s.Tmpl.Lookup(layout); templ != nil {
			return templ
		}
	}
	return nil
}

func (s *Site) publish(path string, r io.Reader) (err error) {
	path = filepath.Join(s.absPublishDir(), path)
	return helpers.WriteToDisk(path, r, s.Fs.Destination)
}

func (s *Site) draftStats() string {
	var msg string

	switch s.draftCount {
	case 0:
		return "0 draft content"
	case 1:
		msg = "1 draft rendered"
	default:
		msg = fmt.Sprintf("%d drafts rendered", s.draftCount)
	}

	if s.Cfg.GetBool("buildDrafts") {
		return fmt.Sprintf("%d of ", s.draftCount) + msg
	}

	return "0 of " + msg
}

func (s *Site) futureStats() string {
	var msg string

	switch s.futureCount {
	case 0:
		return "0 future content"
	case 1:
		msg = "1 future rendered"
	default:
		msg = fmt.Sprintf("%d futures rendered", s.futureCount)
	}

	if s.Cfg.GetBool("buildFuture") {
		return fmt.Sprintf("%d of ", s.futureCount) + msg
	}

	return "0 of " + msg
}

func (s *Site) expiredStats() string {
	var msg string

	switch s.expiredCount {
	case 0:
		return "0 expired content"
	case 1:
		msg = "1 expired rendered"
	default:
		msg = fmt.Sprintf("%d expired rendered", s.expiredCount)
	}

	if s.Cfg.GetBool("buildExpired") {
		return fmt.Sprintf("%d of ", s.expiredCount) + msg
	}

	return "0 of " + msg
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
		Data:     make(map[string]interface{}),
		Site:     &s.Info,
		sections: sections,
		s:        s}

	p.outputFormats = p.s.outputFormats[p.Kind]

	return p

}

func (s *Site) newHomePage() *Page {
	p := s.newNodePage(KindHome)
	p.Title = s.Info.Title
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
		p.Title = helpers.FirstUpper(key)
		key = s.PathSpec.MakePathSanitized(key)
	} else {
		p.Title = strings.Replace(strings.Title(key), "-", " ", -1)
	}

	return p
}

func (s *Site) newSectionPage(name string) *Page {
	p := s.newNodePage(KindSection, name)

	sectionName := helpers.FirstUpper(name)
	if s.Cfg.GetBool("pluralizeListTitles") {
		p.Title = inflect.Pluralize(sectionName)
	} else {
		p.Title = sectionName
	}
	return p
}

func (s *Site) newTaxonomyTermsPage(plural string) *Page {
	p := s.newNodePage(KindTaxonomyTerm, plural)
	p.Title = strings.Title(plural)
	return p
}
