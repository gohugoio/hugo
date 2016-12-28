// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bep/inflect"

	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/parser"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/hugo/transform"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/nitro"
	"github.com/spf13/viper"
)

var _ = transform.AbsURL

// used to indicate if run as a test.
var testMode bool

var defaultTimer *nitro.B

var (
	distinctErrorLogger    = helpers.NewDistinctErrorLogger()
	distinctFeedbackLogger = helpers.NewDistinctFeedbackLogger()
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

	Source         source.Input
	Sections       Taxonomy
	Info           SiteInfo
	Menus          Menus
	timer          *nitro.B
	targets        targetList
	targetListInit sync.Once
	draftCount     int
	futureCount    int
	expiredCount   int
	Data           map[string]interface{}
	Language       *helpers.Language
}

// reset returns a new Site prepared for rebuild.
func (s *Site) reset() *Site {
	return &Site{Language: s.Language, owner: s.owner, PageCollections: newPageCollections()}
}

// newSite creates a new site in the given language.
func newSite(lang *helpers.Language) *Site {
	c := newPageCollections()
	return &Site{Language: lang, PageCollections: c, Info: newSiteInfo(siteBuilderCfg{pageCollections: c, language: lang})}

}

// newSite creates a new site in the default language.
func newSiteDefaultLang() *Site {
	return newSite(helpers.NewDefaultLanguage())
}

// Convenience func used in tests.
func newSiteFromSources(pathContentPairs ...string) *Site {
	if len(pathContentPairs)%2 != 0 {
		panic("pathContentPairs must come in pairs")
	}

	sources := make([]source.ByteSource, 0)

	for i := 0; i < len(pathContentPairs); i += 2 {
		path := pathContentPairs[i]
		content := pathContentPairs[i+1]
		sources = append(sources, source.ByteSource{Name: filepath.FromSlash(path), Content: []byte(content)})
	}

	lang := helpers.NewDefaultLanguage()

	return &Site{
		PageCollections: newPageCollections(),
		Source:          &source.InMemorySource{ByteSource: sources},
		Language:        lang,
		Info:            newSiteInfo(siteBuilderCfg{language: lang}),
	}
}

type targetList struct {
	page          target.Output
	pageUgly      target.Output
	file          target.Output
	alias         target.AliasPublisher
	languageAlias target.AliasPublisher
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

	BaseURL    template.URL
	Taxonomies TaxonomyList
	Authors    AuthorList
	Social     SiteSocial
	Sections   Taxonomy
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
	preserveTaxonomyNames bool
	Data                  *map[string]interface{}

	owner                          *HugoSites
	multilingual                   *Multilingual
	Language                       *helpers.Language
	LanguagePrefix                 string
	Languages                      helpers.Languages
	defaultContentLanguageInSubdir bool
	sectionPagesMenu               string

	pathSpec *helpers.PathSpec
}

func (s *SiteInfo) String() string {
	return fmt.Sprintf("Site(%q)", s.Title)
}

// Used in tests.

type siteBuilderCfg struct {
	language        *helpers.Language
	pageCollections *PageCollections
	baseURL         string
}

func newSiteInfo(cfg siteBuilderCfg) SiteInfo {
	return SiteInfo{
		BaseURL:         template.URL(cfg.baseURL),
		pathSpec:        helpers.NewPathSpecFromConfig(cfg.language),
		multilingual:    newMultiLingualForLanguage(cfg.language),
		PageCollections: cfg.pageCollections,
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

// Param is a convenience method to do lookups in Site's Params map.
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

// GetParam gets a site parameter value if found, nil if not.
func (s *SiteInfo) GetParam(key string) interface{} {
	helpers.Deprecated("SiteInfo", ".GetParam", ".Param", false)
	v := s.Params[strings.ToLower(key)]

	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val
	case int64, int32, int16, int8, int:
		return cast.ToInt(v)
	case float64, float32:
		return cast.ToFloat64(v)
	case time.Time:
		return val
	case []string:
		return v
	}
	return nil
}

func (s *SiteInfo) IsMultiLingual() bool {
	return len(s.Languages) > 1
}

func (s *SiteInfo) refLink(ref string, page *Page, relative bool) (string, error) {
	var refURL *url.URL
	var err error

	refURL, err = url.Parse(ref)

	if err != nil {
		return "", err
	}

	var target *Page
	var link string

	if refURL.Path != "" {
		for _, page := range s.AllRegularPages {
			refPath := filepath.FromSlash(refURL.Path)
			if page.Source.Path() == refPath || page.Source.LogicalName() == refPath {
				target = page
				break
			}
		}

		if target == nil {
			return "", fmt.Errorf("No page found with path or logical name \"%s\".\n", refURL.Path)
		}

		if relative {
			link = target.RelPermalink()
		} else {
			link = target.Permalink()
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
func (s *SiteInfo) Ref(ref string, page *Page) (string, error) {
	return s.refLink(ref, page, false)
}

// RelRef will give an relative URL to ref in the given Page.
func (s *SiteInfo) RelRef(ref string, page *Page) (string, error) {
	return s.refLink(ref, page, true)
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

// reBuild partially rebuilds a site given the filesystem events.
// It returns whetever the content source was changed.
func (s *Site) reProcess(events []fsnotify.Event) (whatChanged, error) {

	jww.DEBUG.Printf("Rebuild for events %q", events)

	s.timerStep("initialize rebuild")

	// First we need to determine what changed

	sourceChanged := []fsnotify.Event{}
	sourceReallyChanged := []fsnotify.Event{}
	tmplChanged := []fsnotify.Event{}
	dataChanged := []fsnotify.Event{}
	i18nChanged := []fsnotify.Event{}

	// prevent spamming the log on changes
	logger := helpers.NewDistinctFeedbackLogger()

	for _, ev := range events {
		if s.isContentDirEvent(ev) {
			logger.Println("Source changed", ev.Name)
			sourceChanged = append(sourceChanged, ev)
		}
		if s.isLayoutDirEvent(ev) {
			logger.Println("Template changed", ev.Name)
			tmplChanged = append(tmplChanged, ev)
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

	if len(tmplChanged) > 0 {
		s.prepTemplates(nil)
		s.owner.tmpl.PrintErrors()
		s.timerStep("template prep")
	}

	if len(dataChanged) > 0 {
		s.readDataFromSourceFS()
	}

	if len(i18nChanged) > 0 {
		if err := s.readI18nSources(); err != nil {
			jww.ERROR.Println(err)
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
		go pageConverter(s, pageChan, convertResults, wg2)
	}

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
			if ex, err := afero.Exists(hugofs.Source(), ev.Name); !ex || err != nil {
				path, _ := helpers.GetRelativePath(ev.Name, s.getContentDir(ev.Name))
				s.removePageByPath(path)
				continue
			}
		}

		sourceReallyChanged = append(sourceReallyChanged, ev)
	}

	go incrementalReadCollator(s, readResults, pageChan, fileConvChan, coordinator, errs)
	go converterCollator(s, convertResults, errs)

	for _, ev := range sourceReallyChanged {

		file, err := s.reReadFile(ev.Name)

		if err != nil {
			jww.ERROR.Println("Error reading file", ev.Name, ";", err)
		}

		if file != nil {
			filechan <- file
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
			jww.ERROR.Println(err)
		}
	}

	changed := whatChanged{
		source: len(sourceChanged) > 0,
		other:  len(tmplChanged) > 0 || len(i18nChanged) > 0 || len(dataChanged) > 0,
	}

	return changed, nil

}

func (s *Site) loadTemplates() {
	s.owner.tmpl = tpl.InitializeT()
	s.owner.tmpl.LoadTemplates(s.absLayoutDir())
	if s.hasTheme() {
		s.owner.tmpl.LoadTemplatesWithPrefix(s.absThemeDir()+"/layouts", "theme")
	}
}

func (s *Site) prepTemplates(withTemplate func(templ tpl.Template) error) error {
	s.loadTemplates()

	if withTemplate != nil {
		if err := withTemplate(s.owner.tmpl); err != nil {
			return err
		}
	}

	s.owner.tmpl.MarkReady()

	return nil
}

func (s *Site) loadData(sources []source.Input) (err error) {
	jww.DEBUG.Printf("Load Data from %q", sources)
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

			data, err := readData(r)
			if err != nil {
				return fmt.Errorf("Failed to read data from %s: %s", filepath.Join(r.Path(), r.LogicalName()), err)
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
						jww.WARN.Printf("Data for key '%s' in path '%s' is overridden in subfolder", key, r.Path())
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

func readData(f *source.File) (interface{}, error) {
	switch f.Extension() {
	case "yaml", "yml":
		return parser.HandleYAMLMetaData(f.Bytes())
	case "json":
		return parser.HandleJSONMetaData(f.Bytes())
	case "toml":
		return parser.HandleTOMLMetaData(f.Bytes())
	default:
		jww.WARN.Printf("Data not supported for extension '%s'", f.Extension())
		return nil, nil
	}
}

func (s *Site) readI18nSources() error {

	i18nSources := []source.Input{&source.Filesystem{Base: s.absI18nDir()}}

	themeI18nDir, err := helpers.GetThemeI18nDirPath()
	if err == nil {
		i18nSources = []source.Input{&source.Filesystem{Base: themeI18nDir}, i18nSources[0]}
	}

	if err = loadI18n(i18nSources); err != nil {
		return err
	}

	return nil
}

func (s *Site) readDataFromSourceFS() error {
	dataSources := make([]source.Input, 0, 2)
	dataSources = append(dataSources, &source.Filesystem{Base: s.absDataDir()})

	// have to be last - duplicate keys in earlier entries will win
	themeDataDir, err := helpers.GetThemeDataDirPath()
	if err == nil {
		dataSources = append(dataSources, &source.Filesystem{Base: themeDataDir})
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
	s.prepTemplates(config.withTemplate)
	s.owner.tmpl.PrintErrors()
	s.timerStep("initialize & template prep")

	if err = s.readDataFromSourceFS(); err != nil {
		return
	}

	if err = s.readI18nSources(); err != nil {
		return
	}

	s.timerStep("load i18n")
	return s.createPages()

}

func (s *Site) setupPrevNext() {
	for i, page := range s.Pages {
		if i < len(s.Pages)-1 {
			page.Next = s.Pages[i+1]
		}

		if i > 0 {
			page.Prev = s.Pages[i-1]
		}
	}
}

func (s *Site) setCurrentLanguageConfig() error {
	// There are sadly some global template funcs etc. that need the language information.
	viper.Set("multilingual", s.multilingualEnabled())
	viper.Set("currentContentLanguage", s.Language)
	// Cache the current config.
	helpers.InitConfigProviderForCurrentContentLanguage()
	s.Info.pathSpec = helpers.CurrentPathSpec()
	return tpl.SetTranslateLang(s.Language)
}

func (s *Site) render() (err error) {
	if err = s.setCurrentLanguageConfig(); err != nil {
		return
	}

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

	if err = s.renderPages(); err != nil {
		return
	}
	s.timerStep("render and write pages")

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
		jww.DEBUG.Println("initialize: Source is already set")
		return
	}

	if err = s.checkDirectories(); err != nil {
		return err
	}

	staticDir := helpers.AbsPathify(viper.GetString("staticDir") + "/")

	s.Source = &source.Filesystem{
		AvoidPaths: []string{staticDir},
		Base:       s.absContentDir(),
	}

	return
}

// HomeAbsURL is a convenience method giving the absolute URL to the home page.
func (s *SiteInfo) HomeAbsURL() string {
	base := ""
	if s.IsMultiLingual() {
		base = s.Language.Lang
	}
	return s.pathSpec.AbsURL(base, false)
}

// SitemapAbsURL is a convenience method giving the absolute URL to the sitemap.
func (s *SiteInfo) SitemapAbsURL() string {
	sitemapDefault := parseSitemap(viper.GetStringMap("sitemap"))
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
	for k, v := range viper.GetStringMapString("permalinks") {
		permalinks[k] = pathPattern(v)
	}

	defaultContentInSubDir := viper.GetBool("defaultContentLanguageInSubdir")
	defaultContentLanguage := viper.GetString("defaultContentLanguage")

	languagePrefix := ""
	if s.multilingualEnabled() && (defaultContentInSubDir || lang.Lang != defaultContentLanguage) {
		languagePrefix = "/" + lang.Lang
	}

	var multilingual *Multilingual
	if s.owner != nil {
		multilingual = s.owner.multilingual
	}

	s.Info = SiteInfo{
		BaseURL:                        template.URL(helpers.SanitizeURLKeepTrailingSlash(viper.GetString("baseURL"))),
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
		BuildDrafts:                    viper.GetBool("buildDrafts"),
		canonifyURLs:                   viper.GetBool("canonifyURLs"),
		preserveTaxonomyNames:          lang.GetBool("preserveTaxonomyNames"),
		PageCollections:                s.PageCollections,
		Files:                          &s.Files,
		Menus:                          &s.Menus,
		Params:                         params,
		Permalinks:                     permalinks,
		Data:                           &s.Data,
		owner:                          s.owner,
		pathSpec:                       helpers.NewPathSpecFromConfig(lang),
	}

	s.Info.RSSLink = s.Info.permalinkStr(lang.GetString("rssURI"))
}

func (s *Site) hasTheme() bool {
	return viper.GetString("theme") != ""
}

func (s *Site) dataDir() string {
	return viper.GetString("dataDir")
}
func (s *Site) absDataDir() string {
	return helpers.AbsPathify(s.dataDir())
}

func (s *Site) i18nDir() string {
	return viper.GetString("i18nDir")
}

func (s *Site) absI18nDir() string {
	return helpers.AbsPathify(s.i18nDir())
}

func (s *Site) isI18nEvent(e fsnotify.Event) bool {
	if s.getI18nDir(e.Name) != "" {
		return true
	}
	return s.getThemeI18nDir(e.Name) != ""
}

func (s *Site) getI18nDir(path string) string {
	return getRealDir(s.absI18nDir(), path)
}

func (s *Site) getThemeI18nDir(path string) string {
	if !s.hasTheme() {
		return ""
	}
	return getRealDir(helpers.AbsPathify(filepath.Join(s.themeDir(), s.i18nDir())), path)
}

func (s *Site) isDataDirEvent(e fsnotify.Event) bool {
	if s.getDataDir(e.Name) != "" {
		return true
	}
	return s.getThemeDataDir(e.Name) != ""
}

func (s *Site) getDataDir(path string) string {
	return getRealDir(s.absDataDir(), path)
}

func (s *Site) getThemeDataDir(path string) string {
	if !s.hasTheme() {
		return ""
	}
	return getRealDir(helpers.AbsPathify(filepath.Join(s.themeDir(), s.dataDir())), path)
}

func (s *Site) themeDir() string {
	return viper.GetString("themesDir") + "/" + viper.GetString("theme")
}

func (s *Site) absThemeDir() string {
	return helpers.AbsPathify(s.themeDir())
}

func (s *Site) layoutDir() string {
	return viper.GetString("layoutDir")
}

func (s *Site) absLayoutDir() string {
	return helpers.AbsPathify(s.layoutDir())
}

func (s *Site) isLayoutDirEvent(e fsnotify.Event) bool {
	if s.getLayoutDir(e.Name) != "" {
		return true
	}
	return s.getThemeLayoutDir(e.Name) != ""
}

func (s *Site) getLayoutDir(path string) string {
	return getRealDir(s.absLayoutDir(), path)
}

func (s *Site) getThemeLayoutDir(path string) string {
	if !s.hasTheme() {
		return ""
	}
	return getRealDir(helpers.AbsPathify(filepath.Join(s.themeDir(), s.layoutDir())), path)
}

func (s *Site) absContentDir() string {
	return helpers.AbsPathify(viper.GetString("contentDir"))
}

func (s *Site) isContentDirEvent(e fsnotify.Event) bool {
	return s.getContentDir(e.Name) != ""
}

func (s *Site) getContentDir(path string) string {
	return getRealDir(s.absContentDir(), path)
}

// getRealDir gets the base path of the given path, also handling the case where
// base is a symlinked folder.
func getRealDir(base, path string) string {

	if strings.HasPrefix(path, base) {
		return base
	}

	realDir, err := helpers.GetRealPath(hugofs.Source(), base)

	if err != nil {
		if !os.IsNotExist(err) {
			jww.ERROR.Printf("Failed to get real path for %s: %s", path, err)
		}
		return ""
	}

	if strings.HasPrefix(path, realDir) {
		return realDir
	}

	return ""
}

func (s *Site) absPublishDir() string {
	return helpers.AbsPathify(viper.GetString("publishDir"))
}

func (s *Site) checkDirectories() (err error) {
	if b, _ := helpers.DirExists(s.absContentDir(), hugofs.Source()); !b {
		return errors.New("No source directory found, expecting to find it at " + s.absContentDir())
	}
	return
}

// reReadFile resets file to be read from disk again
func (s *Site) reReadFile(absFilePath string) (*source.File, error) {
	jww.INFO.Println("rereading", absFilePath)
	var file *source.File

	reader, err := source.NewLazyFileReader(hugofs.Source(), absFilePath)
	if err != nil {
		return nil, err
	}
	file, err = source.NewFileFromAbs(s.getContentDir(absFilePath), absFilePath, reader)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *Site) readPagesFromSource() chan error {
	if s.Source == nil {
		panic(fmt.Sprintf("s.Source not set %s", s.absContentDir()))
	}

	jww.DEBUG.Printf("Read %d pages from source", len(s.Source.Files()))

	errs := make(chan error)
	if len(s.Source.Files()) < 1 {
		close(errs)
		return errs
	}

	files := s.Source.Files()
	results := make(chan HandledResult)
	filechan := make(chan *source.File)
	procs := getGoMaxProcs()
	wg := &sync.WaitGroup{}

	wg.Add(procs * 4)
	for i := 0; i < procs*4; i++ {
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
	procs := getGoMaxProcs()
	wg := &sync.WaitGroup{}

	wg.Add(2 * procs * 4)
	for i := 0; i < procs*4; i++ {
		go fileConverter(s, fileConvChan, results, wg)
		go pageConverter(s, pageChan, results, wg)
	}

	go converterCollator(s, results, errs)

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
		jww.ERROR.Println("Unsupported File Type", file.Path())
	}
}

func pageConverter(s *Site, pages <-chan *Page, results HandleResults, wg *sync.WaitGroup) {
	defer wg.Done()
	for page := range pages {
		var h *MetaHandle
		if page.Markup != "" {
			h = NewMetaHandler(page.Markup)
		} else {
			h = NewMetaHandler(page.File.Extension())
		}
		if h != nil {
			h.Convert(page, s, results)
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

func converterCollator(s *Site, results <-chan HandledResult, errs chan<- error) {
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

	s.assembleSections()

	s.assembleMenus()

	s.Info.LastChange = s.Pages[0].Lastmod

	return
}

func (s *Site) getMenusFromConfig() Menus {

	ret := Menus{}

	if menus := s.Language.GetStringMap("menu"); menus != nil {
		for name, menu := range menus {
			m, err := cast.ToSliceE(menu)
			if err != nil {
				jww.ERROR.Printf("unable to process menus in site config\n")
				jww.ERROR.Println(err)
			} else {
				for _, entry := range m {
					jww.DEBUG.Printf("found menu: %q, in site config\n", name)

					menuEntry := MenuEntry{Menu: name}
					ime, err := cast.ToStringMapE(entry)
					if err != nil {
						jww.ERROR.Printf("unable to process menus in site config\n")
						jww.ERROR.Println(err)
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
	menuEntryURL = helpers.SanitizeURLKeepTrailingSlash(s.pathSpec.URLize(menuEntryURL))
	if !s.canonifyURLs {
		menuEntryURL = helpers.AddContextRoot(string(s.BaseURL), menuEntryURL)
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

	menuConfig := s.getMenusFromConfig()
	for name, menu := range menuConfig {
		for _, me := range *menu {
			flat[twoD{name, me.KeyName()}] = me
		}
	}

	sectionPagesMenu := s.Info.sectionPagesMenu
	sectionPagesMenus := make(map[string]interface{})
	//creating flat hash
	pages := s.Pages
	for _, p := range pages {

		if sectionPagesMenu != "" {
			if _, ok := sectionPagesMenus[p.Section()]; !ok {
				if p.Section() != "" {
					me := MenuEntry{Identifier: p.Section(),
						Name: helpers.MakeTitle(helpers.FirstUpper(p.Section())),
						URL:  s.Info.createNodeMenuEntryURL(p.addLangPathPrefix("/"+p.Section()) + "/")}
					if _, ok := flat[twoD{sectionPagesMenu, me.KeyName()}]; ok {
						// menu with same id defined in config, let that one win
						continue
					}
					flat[twoD{sectionPagesMenu, me.KeyName()}] = &me
					sectionPagesMenus[p.Section()] = true
				}
			}
		}

		for name, me := range p.Menus() {
			if _, ok := flat[twoD{name, me.KeyName()}]; ok {
				jww.ERROR.Printf("Two or more menu items have the same name/identifier in Menu %q: %q.\nRename or set an unique identifier.\n", name, me.KeyName())
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

func (s *Site) assembleTaxonomies() {
	s.Taxonomies = make(TaxonomyList)
	s.taxonomiesPluralSingular = make(map[string]string)
	s.taxonomiesOrigKey = make(map[string]string)

	taxonomies := s.Language.GetStringMapString("taxonomies")

	jww.INFO.Printf("found taxonomies: %#v\n", taxonomies)

	for singular, plural := range taxonomies {
		s.Taxonomies[plural] = make(Taxonomy)
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
						s.Taxonomies[plural].add(idx, x, s.Info.preserveTaxonomyNames)
						if s.Info.preserveTaxonomyNames {
							// Need to track the original
							s.taxonomiesOrigKey[fmt.Sprintf("%s-%s", plural, kp(idx))] = idx
						}
					}
				} else if v, ok := vals.(string); ok {
					x := WeightedPage{weight.(int), p}
					s.Taxonomies[plural].add(v, x, s.Info.preserveTaxonomyNames)
					if s.Info.preserveTaxonomyNames {
						// Need to track the original
						s.taxonomiesOrigKey[fmt.Sprintf("%s-%s", plural, kp(v))] = v
					}
				} else {
					jww.ERROR.Printf("Invalid %s in %s\n", plural, p.File.Path())
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

	s.Info.paginationPageCount = 0
	s.draftCount = 0
	s.futureCount = 0
	s.expiredCount = 0

	for _, p := range s.rawAllPages {
		p.scratch = newScratch()
	}
}

func (s *Site) assembleSections() {
	s.Sections = make(Taxonomy)
	s.Info.Sections = s.Sections

	regularPages := s.findPagesByKind(KindPage)
	sectionPages := s.findPagesByKind(KindSection)

	for i, p := range regularPages {
		s.Sections.add(p.Section(), WeightedPage{regularPages[i].Weight, regularPages[i]}, s.Info.preserveTaxonomyNames)
	}

	// Add sections without regular pages, but with a content page
	for _, sectionPage := range sectionPages {
		if _, ok := s.Sections[sectionPage.sections[0]]; !ok {
			s.Sections[sectionPage.sections[0]] = WeightedPages{}
		}
	}

	for k := range s.Sections {
		s.Sections[k].Sort()

		for i, wp := range s.Sections[k] {
			if i > 0 {
				wp.Page.NextInSection = s.Sections[k][i-1].Page
			}
			if i < len(s.Sections[k])-1 {
				wp.Page.PrevInSection = s.Sections[k][i+1].Page
			}
		}
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
	if !s.hasTheme() {
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
	jww.FEEDBACK.Printf("Built site for language %s:\n", s.Language.Lang)
	jww.FEEDBACK.Println(s.draftStats())
	jww.FEEDBACK.Println(s.futureStats())
	jww.FEEDBACK.Println(s.expiredStats())
	jww.FEEDBACK.Printf("%d regular pages created\n", len(s.RegularPages))
	jww.FEEDBACK.Printf("%d other pages created\n", (len(s.Pages) - len(s.RegularPages)))
	jww.FEEDBACK.Printf("%d non-page files copied\n", len(s.Files))
	jww.FEEDBACK.Printf("%d paginator pages created\n", s.Info.paginationPageCount)
	taxonomies := s.Language.GetStringMapString("taxonomies")

	for _, pl := range taxonomies {
		jww.FEEDBACK.Printf("%d %s created\n", len(s.Taxonomies[pl]), pl)
	}

}

// GetPage looks up a index page of a given type in the path given.
// This method may support regular pages in the future, but currently it is a
// convenient way of getting the home page or
// a section from a template:
//    {{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}
//
// This will return nil when no page could be found.
//
// The valid page types are: home, section, taxonomy and taxonomyTerm
func (s *SiteInfo) GetPage(typ string, path ...string) *Page {
	return s.getPage(typ, path...)
}

func (s *SiteInfo) permalink(plink string) string {
	return s.permalinkStr(plink)
}

func (s *SiteInfo) permalinkStr(plink string) string {
	return helpers.MakePermalink(
		viper.GetString("baseURL"),
		s.pathSpec.URLizeAndPrep(plink)).String()
}

func (s *Site) renderAndWriteXML(name string, dest string, d interface{}, layouts ...string) error {
	jww.DEBUG.Printf("Render XML for %q to %q", name, dest)
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)
	renderBuffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n")

	err := s.renderForLayouts(name, d, renderBuffer, layouts...)

	if err != nil {
		return err
	}

	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)

	var path []byte
	if viper.GetBool("relativeURLs") {
		path = []byte(helpers.GetDottedRelativePath(dest))
	} else {
		s := viper.GetString("baseURL")
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		path = []byte(s)
	}
	transformer := transform.NewChain(transform.AbsURLInXML)
	transformer.Apply(outBuffer, renderBuffer, path)

	return s.writeDestFile(dest, outBuffer)

}

func (s *Site) renderAndWritePage(name string, dest string, d interface{}, layouts ...string) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	err := s.renderForLayouts(name, d, renderBuffer, layouts...)

	if err != nil {
		return err
	}

	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)

	var pageTarget target.Output

	if p, ok := d.(*Page); ok && p.IsPage() && path.Ext(p.URLPath.URL) != "" {
		// user has explicitly set a URL with extension for this page
		// make sure it sticks even if "ugly URLs" are turned off.
		pageTarget = s.pageUglyTarget()
	} else {
		pageTarget = s.pageTarget()
	}

	transformLinks := transform.NewEmptyTransforms()

	if viper.GetBool("relativeURLs") || viper.GetBool("canonifyURLs") {
		transformLinks = append(transformLinks, transform.AbsURL)
	}

	if s.running() && viper.GetBool("watch") && !viper.GetBool("disableLiveReload") {
		transformLinks = append(transformLinks, transform.LiveReloadInject)
	}

	// For performance reasons we only inject the Hugo generator tag on the home page.
	if n, ok := d.(*Page); ok && n.IsHome() {
		if !viper.GetBool("disableHugoGeneratorInject") {
			transformLinks = append(transformLinks, transform.HugoGeneratorInject)
		}
	}

	var path []byte

	if viper.GetBool("relativeURLs") {
		translated, err := pageTarget.(target.OptionalTranslator).TranslateRelative(dest)
		if err != nil {
			return err
		}
		path = []byte(helpers.GetDottedRelativePath(translated))
	} else if viper.GetBool("canonifyURLs") {
		s := viper.GetString("baseURL")
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		path = []byte(s)
	}

	transformer := transform.NewChain(transformLinks...)
	transformer.Apply(outBuffer, renderBuffer, path)

	if outBuffer.Len() == 0 {

		jww.WARN.Printf("%s is rendered empty\n", dest)
		if dest == "/" {
			debugAddend := ""
			if !viper.GetBool("verbose") {
				debugAddend = "* For more debugging information, run \"hugo -v\""
			}
			distinctFeedbackLogger.Printf(`=============================================================
Your rendered home page is blank: /index.html is zero-length
 * Did you specify a theme on the command-line or in your
   %q file?  (Current theme: %q)
 %s
=============================================================`,
				filepath.Base(viper.ConfigFileUsed()),
				viper.GetString("theme"),
				debugAddend)
		}

		// Avoid writing empty files to disk.
		return nil

	}

	if err = s.writeDestPage(dest, pageTarget, outBuffer); err != nil {
		return err
	}

	return nil
}

func (s *Site) renderForLayouts(name string, d interface{}, w io.Writer, layouts ...string) error {
	layout, found := s.findFirstLayout(layouts...)
	if !found {
		jww.WARN.Printf("Unable to locate layout for %s: %s\n", name, layouts)
		return nil
	}

	if err := s.renderThing(d, layout, w); err != nil {

		// Behavior here should be dependent on if running in server or watch mode.
		distinctErrorLogger.Printf("Error while rendering %q: %s", name, err)
		if !s.running() && !testMode {
			// TODO(bep) check if this can be propagated
			os.Exit(-1)
		} else if testMode {
			return err
		}
	}

	return nil
}

func (s *Site) findFirstLayout(layouts ...string) (string, bool) {
	for _, layout := range layouts {
		if s.owner.tmpl.Lookup(layout) != nil {
			return layout, true
		}
	}
	return "", false
}

func (s *Site) renderThing(d interface{}, layout string, w io.Writer) error {

	// If the template doesn't exist, then return, but leave the Writer open
	if templ := s.owner.tmpl.Lookup(layout); templ != nil {
		return templ.Execute(w, d)
	}
	return fmt.Errorf("Layout not found: %s", layout)

}

func (s *Site) pageTarget() target.Output {
	s.initTargetList()
	return s.targets.page
}

func (s *Site) pageUglyTarget() target.Output {
	s.initTargetList()
	return s.targets.pageUgly
}

func (s *Site) fileTarget() target.Output {
	s.initTargetList()
	return s.targets.file
}

func (s *Site) aliasTarget() target.AliasPublisher {
	s.initTargetList()
	return s.targets.alias
}

func (s *Site) languageAliasTarget() target.AliasPublisher {
	s.initTargetList()
	return s.targets.languageAlias
}

func (s *Site) initTargetList() {
	s.targetListInit.Do(func() {
		langDir := ""
		if s.Language.Lang != s.Info.multilingual.DefaultLang.Lang || s.Info.defaultContentLanguageInSubdir {
			langDir = s.Language.Lang
		}
		if s.targets.page == nil {
			s.targets.page = &target.PagePub{
				PublishDir: s.absPublishDir(),
				UglyURLs:   viper.GetBool("uglyURLs"),
				LangDir:    langDir,
			}
		}
		if s.targets.pageUgly == nil {
			s.targets.pageUgly = &target.PagePub{
				PublishDir: s.absPublishDir(),
				UglyURLs:   true,
				LangDir:    langDir,
			}
		}
		if s.targets.file == nil {
			s.targets.file = &target.Filesystem{
				PublishDir: s.absPublishDir(),
			}
		}
		if s.targets.alias == nil {
			s.targets.alias = &target.HTMLRedirectAlias{
				PublishDir: s.absPublishDir(),
				Templates:  s.owner.tmpl.Lookup("alias.html"),
			}
		}
		if s.targets.languageAlias == nil {
			s.targets.languageAlias = &target.HTMLRedirectAlias{
				PublishDir: s.absPublishDir(),
				AllowRoot:  true,
			}
		}
	})
}

func (s *Site) writeDestFile(path string, reader io.Reader) (err error) {
	jww.DEBUG.Println("creating file:", path)
	return s.fileTarget().Publish(path, reader)
}

func (s *Site) writeDestPage(path string, publisher target.Publisher, reader io.Reader) (err error) {
	jww.DEBUG.Println("creating page:", path)
	return publisher.Publish(path, reader)
}

// AliasPublisher
func (s *Site) writeDestAlias(path, permalink string, p *Page) (err error) {
	return s.publishDestAlias(s.aliasTarget(), path, permalink, p)
}

func (s *Site) publishDestAlias(aliasPublisher target.AliasPublisher, path, permalink string, p *Page) (err error) {
	if viper.GetBool("relativeURLs") {
		// convert `permalink` into URI relative to location of `path`
		baseURL := helpers.SanitizeURLKeepTrailingSlash(viper.GetString("baseURL"))
		if strings.HasPrefix(permalink, baseURL) {
			permalink = "/" + strings.TrimPrefix(permalink, baseURL)
		}
		permalink, err = helpers.GetRelativePath(permalink, path)
		if err != nil {
			jww.ERROR.Println("Failed to make a RelativeURL alias:", path, "redirecting to", permalink)
		}
		permalink = filepath.ToSlash(permalink)
	}
	jww.DEBUG.Println("creating alias:", path, "redirecting to", permalink)
	return aliasPublisher.Publish(path, permalink, p)
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

	if viper.GetBool("buildDrafts") {
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

	if viper.GetBool("buildFuture") {
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

	if viper.GetBool("buildExpired") {
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

func (s *Site) newNodePage(typ string) *Page {
	return &Page{
		pageInit: &pageInit{},
		Kind:     typ,
		Data:     make(map[string]interface{}),
		Site:     &s.Info,
		language: s.Language,
		site:     s}
}

func (s *Site) newHomePage() *Page {
	p := s.newNodePage(KindHome)
	p.Title = s.Info.Title
	pages := Pages{}
	p.Data["Pages"] = pages
	p.Pages = pages
	s.setPageURLs(p, "/")
	return p
}

func (s *Site) setPageURLs(p *Page, in string) {
	p.URLPath.URL = s.Info.pathSpec.URLizeAndPrep(in)
	p.URLPath.Permalink = s.Info.permalink(p.URLPath.URL)
	p.RSSLink = template.HTML(s.Info.permalink(in + ".xml"))
}

func (s *Site) newTaxonomyPage(plural, key string) *Page {

	p := s.newNodePage(KindTaxonomy)

	p.sections = []string{plural, key}

	if s.Info.preserveTaxonomyNames {
		key = s.Info.pathSpec.MakePathSanitized(key)
	}

	if s.Info.preserveTaxonomyNames {
		// keep as is in the title
		p.Title = key
	} else {
		p.Title = strings.Replace(strings.Title(key), "-", " ", -1)
	}

	s.setPageURLs(p, path.Join(plural, key))

	return p
}

func (s *Site) newSectionPage(name string, section WeightedPages) *Page {

	p := s.newNodePage(KindSection)
	p.sections = []string{name}

	sectionName := name
	if !s.Info.preserveTaxonomyNames && len(section) > 0 {
		sectionName = section[0].Page.Section()
	}

	sectionName = helpers.FirstUpper(sectionName)
	if viper.GetBool("pluralizeListTitles") {
		p.Title = inflect.Pluralize(sectionName)
	} else {
		p.Title = sectionName
	}
	s.setPageURLs(p, name)
	return p
}

func (s *Site) newTaxonomyTermsPage(plural string) *Page {
	p := s.newNodePage(KindTaxonomyTerm)
	p.sections = []string{plural}
	p.Title = strings.Title(plural)
	s.setPageURLs(p, plural)
	return p
}
