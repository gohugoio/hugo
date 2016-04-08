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
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"sync/atomic"

	"path"

	"github.com/bep/inflect"
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
	"gopkg.in/fsnotify.v1"
)

var _ = transform.AbsURL

// used to indicate if run as a test.
var testMode bool

var defaultTimer *nitro.B

var distinctErrorLogger = helpers.NewDistinctErrorLogger()

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
	Pages          Pages
	Files          []*source.File
	Tmpl           tpl.Template
	Taxonomies     TaxonomyList
	Source         source.Input
	Sections       Taxonomy
	Info           SiteInfo
	Menus          Menus
	timer          *nitro.B
	targets        targetList
	targetListInit sync.Once
	RunMode        runmode
	draftCount     int
	futureCount    int
	Data           map[string]interface{}
}

type targetList struct {
	page     target.Output
	pageUgly target.Output
	file     target.Output
	alias    target.AliasPublisher
}

type SiteInfo struct {
	BaseURL               template.URL
	Taxonomies            TaxonomyList
	Authors               AuthorList
	Social                SiteSocial
	Sections              Taxonomy
	Pages                 *Pages
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
	paginationPageCount   uint64
	Data                  *map[string]interface{}
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

// GetParam gets a site parameter value if found, nil if not.
func (s *SiteInfo) GetParam(key string) interface{} {
	v := s.Params[strings.ToLower(key)]

	if v == nil {
		return nil
	}

	switch v.(type) {
	case bool:
		return cast.ToBool(v)
	case string:
		return cast.ToString(v)
	case int64, int32, int16, int8, int:
		return cast.ToInt(v)
	case float64, float32:
		return cast.ToFloat64(v)
	case time.Time:
		return cast.ToTime(v)
	case []string:
		return v
	}
	return nil
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
		for _, page := range []*Page(*s.Pages) {
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
			link, err = target.RelPermalink()
		} else {
			link, err = target.Permalink()
		}

		if err != nil {
			return "", err
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

		for _, page := range []*Page(*s.Pages) {
			if page.Source.Path() == refPath {
				target = page
				break
			}
		}
		// need to exhaust the test, then try with the others :/
		// if the refPath doesn't end in a filename with extension `.md`, then try with `.md` , and then `/index.md`
		mdPath := strings.TrimSuffix(refPath, string(os.PathSeparator)) + ".md"
		for _, page := range []*Page(*s.Pages) {
			if page.Source.Path() == mdPath {
				target = page
				break
			}
		}
		indexPath := filepath.Join(refPath, "index.md")
		for _, page := range []*Page(*s.Pages) {
			if page.Source.Path() == indexPath {
				target = page
				break
			}
		}

		if target == nil {
			return "", fmt.Errorf("No page found for \"%s\" on page \"%s\".\n", ref, currentPage.Source.Path())
		}

		link, err = target.RelPermalink()

		if err != nil {
			return "", err
		}
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
	return s.RunMode.Watching
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

func (s *Site) Build() (err error) {
	if err = s.Process(); err != nil {
		return
	}

	if err = s.Render(); err != nil {
		// Better reporting when the template is missing (commit 2bbecc7b)
		jww.ERROR.Printf("Error rendering site: %s", err)

		jww.ERROR.Printf("Available templates:")
		var keys []string
		for _, template := range s.Tmpl.Templates() {
			if name := template.Name(); name != "" {
				keys = append(keys, name)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			jww.ERROR.Printf("\t%s\n", k)
		}

		return
	}

	return nil
}

func (s *Site) ReBuild(events []fsnotify.Event) error {
	s.timerStep("initialize rebuild")
	// First we need to determine what changed

	sourceChanged := []fsnotify.Event{}
	tmplChanged := []fsnotify.Event{}
	dataChanged := []fsnotify.Event{}

	var err error

	// prevent spamming the log on changes
	logger := helpers.NewDistinctFeedbackLogger()

	for _, ev := range events {
		// Need to re-read source
		if strings.HasPrefix(ev.Name, s.absContentDir()) {
			sourceChanged = append(sourceChanged, ev)
		}
		if strings.HasPrefix(ev.Name, s.absLayoutDir()) || strings.HasPrefix(ev.Name, s.absThemeDir()) {
			logger.Println("Template changed", ev.Name)
			tmplChanged = append(tmplChanged, ev)
		}
		if strings.HasPrefix(ev.Name, s.absDataDir()) {
			logger.Println("Data changed", ev.Name)
			dataChanged = append(dataChanged, ev)
		}
	}

	if len(tmplChanged) > 0 {
		s.prepTemplates()
		s.Tmpl.PrintErrors()
		s.timerStep("template prep")
	}

	if len(dataChanged) > 0 {
		s.readDataFromSourceFS()
	}

	// we reuse the state, so have to do some cleanup before we can rebuild.
	s.resetPageBuildState()

	// If a content file changes, we need to reload only it and re-render the entire site.

	// First step is to read the changed files and (re)place them in site.Pages
	// This includes processing any meta-data for that content

	// The second step is to convert the content into HTML
	// This includes processing any shortcodes that may be present.

	// We do this in parallel... even though it's likely only one file at a time.
	// We need to process the reading prior to the conversion for each file, but
	// we can convert one file while another one is still reading.
	errs := make(chan error)
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

	go incrementalReadCollator(s, readResults, pageChan, fileConvChan, coordinator, errs)
	go converterCollator(s, convertResults, errs)

	if len(tmplChanged) > 0 || len(dataChanged) > 0 {
		// Do not need to read the files again, but they need conversion
		// for shortocde re-rendering.
		for _, p := range s.Pages {
			pageChan <- p
		}
	}

	for _, ev := range sourceChanged {

		if ev.Op&fsnotify.Remove == fsnotify.Remove {
			//remove the file & a create will follow
			path, _ := helpers.GetRelativePath(ev.Name, s.absContentDir())
			s.removePageByPath(path)
			continue
		}

		// Some editors (Vim) sometimes issue only a Rename operation when writing an existing file
		// Sometimes a rename operation means that file has been renamed other times it means
		// it's been updated
		if ev.Op&fsnotify.Rename == fsnotify.Rename {
			// If the file is still on disk, it's only been updated, if it's not, it's been moved
			if ex, err := afero.Exists(hugofs.Source(), ev.Name); !ex || err != nil {
				path, _ := helpers.GetRelativePath(ev.Name, s.absContentDir())
				s.removePageByPath(path)
				continue
			}
		}

		file, err := s.reReadFile(ev.Name)
		if err != nil {
			errs <- err
		}

		filechan <- file
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

	if len(sourceChanged) > 0 {
		s.setupPrevNext()
		if err = s.buildSiteMeta(); err != nil {
			return err
		}
		s.timerStep("build taxonomies")
	}

	// Once the appropriate prep step is done we render the entire site
	if err = s.Render(); err != nil {
		// Better reporting when the template is missing (commit 2bbecc7b)
		jww.ERROR.Printf("Error rendering site: %s", err)
		jww.ERROR.Printf("Available templates:")
		var keys []string
		for _, template := range s.Tmpl.Templates() {
			if name := template.Name(); name != "" {
				keys = append(keys, name)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			jww.ERROR.Printf("\t%s\n", k)
		}

		return nil
	}
	return err

}

func (s *Site) Analyze() error {
	if err := s.Process(); err != nil {
		return err
	}
	return s.ShowPlan(os.Stdout)
}

func (s *Site) loadTemplates() {
	s.Tmpl = tpl.InitializeT()
	s.Tmpl.LoadTemplates(s.absLayoutDir())
	if s.hasTheme() {
		s.Tmpl.LoadTemplatesWithPrefix(s.absThemeDir()+"/layouts", "theme")
	}
}

func (s *Site) prepTemplates(additionalNameValues ...string) error {
	s.loadTemplates()

	for i := 0; i < len(additionalNameValues); i += 2 {
		err := s.Tmpl.AddTemplate(additionalNameValues[i], additionalNameValues[i+1])
		if err != nil {
			return err
		}
	}
	s.Tmpl.MarkReady()

	return nil
}

func (s *Site) loadData(sources []source.Input) (err error) {
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

func (s *Site) readDataFromSourceFS() error {
	dataSources := make([]source.Input, 0, 2)
	dataSources = append(dataSources, &source.Filesystem{Base: s.absDataDir()})

	// have to be last - duplicate keys in earlier entries will win
	themeStaticDir, err := helpers.GetThemeDataDirPath()
	if err == nil {
		dataSources = append(dataSources, &source.Filesystem{Base: themeStaticDir})
	}

	err = s.loadData(dataSources)
	s.timerStep("load data")
	return err
}

func (s *Site) Process() (err error) {
	s.timerStep("Go initialization")
	if err = s.initialize(); err != nil {
		return
	}
	s.prepTemplates()
	s.Tmpl.PrintErrors()
	s.timerStep("initialize & template prep")

	if err = s.readDataFromSourceFS(); err != nil {
		return
	}

	if err = s.createPages(); err != nil {
		return
	}
	s.setupPrevNext()
	if err = s.buildSiteMeta(); err != nil {
		return
	}
	s.timerStep("build taxonomies")
	return
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

func (s *Site) Render() (err error) {
	if err = s.renderAliases(); err != nil {
		return
	}
	s.timerStep("render and write aliases")
	if err = s.renderTaxonomiesLists(); err != nil {
		return
	}
	s.timerStep("render and write taxonomies")
	s.renderListsOfTaxonomyTerms()
	s.timerStep("render & write taxonomy lists")
	if err = s.renderSectionLists(); err != nil {
		return
	}
	s.timerStep("render and write lists")
	if err = s.renderPages(); err != nil {
		return
	}
	s.timerStep("render and write pages")
	if err = s.renderHomePage(); err != nil {
		return
	}
	s.timerStep("render and write homepage")
	if err = s.renderSitemap(); err != nil {
		return
	}
	s.timerStep("render and write Sitemap")

	if err = s.renderRobotsTXT(); err != nil {
		return
	}
	s.timerStep("render and write robots.txt")

	return
}

func (s *Site) Initialise() (err error) {
	return s.initialize()
}

func (s *Site) initialize() (err error) {
	if err = s.checkDirectories(); err != nil {
		return err
	}

	staticDir := helpers.AbsPathify(viper.GetString("StaticDir") + "/")

	s.Source = &source.Filesystem{
		AvoidPaths: []string{staticDir},
		Base:       s.absContentDir(),
	}

	s.Menus = Menus{}

	s.initializeSiteInfo()

	return
}

func (s *Site) initializeSiteInfo() {
	params := viper.GetStringMap("Params")

	permalinks := make(PermalinkOverrides)
	for k, v := range viper.GetStringMapString("Permalinks") {
		permalinks[k] = pathPattern(v)
	}

	s.Info = SiteInfo{
		BaseURL:               template.URL(helpers.SanitizeURLKeepTrailingSlash(viper.GetString("BaseURL"))),
		Title:                 viper.GetString("Title"),
		Author:                viper.GetStringMap("author"),
		Social:                viper.GetStringMapString("social"),
		LanguageCode:          viper.GetString("languagecode"),
		Copyright:             viper.GetString("copyright"),
		DisqusShortname:       viper.GetString("DisqusShortname"),
		GoogleAnalytics:       viper.GetString("GoogleAnalytics"),
		RSSLink:               s.permalinkStr(viper.GetString("RSSUri")),
		BuildDrafts:           viper.GetBool("BuildDrafts"),
		canonifyURLs:          viper.GetBool("CanonifyURLs"),
		preserveTaxonomyNames: viper.GetBool("PreserveTaxonomyNames"),
		Pages:      &s.Pages,
		Files:      &s.Files,
		Menus:      &s.Menus,
		Params:     params,
		Permalinks: permalinks,
		Data:       &s.Data,
	}
}

func (s *Site) hasTheme() bool {
	return viper.GetString("theme") != ""
}

func (s *Site) absDataDir() string {
	return helpers.AbsPathify(viper.GetString("DataDir"))
}

func (s *Site) absThemeDir() string {
	return helpers.AbsPathify(viper.GetString("themesDir") + "/" + viper.GetString("theme"))
}

func (s *Site) absLayoutDir() string {
	return helpers.AbsPathify(viper.GetString("LayoutDir"))
}

func (s *Site) absContentDir() string {
	return helpers.AbsPathify(viper.GetString("ContentDir"))
}

func (s *Site) absPublishDir() string {
	return helpers.AbsPathify(viper.GetString("PublishDir"))
}

func (s *Site) checkDirectories() (err error) {
	if b, _ := helpers.DirExists(s.absContentDir(), hugofs.Source()); !b {
		return fmt.Errorf("No source directory found, expecting to find it at " + s.absContentDir())
	}
	return
}

// reReadFile resets file to be read from disk again
func (s *Site) reReadFile(absFilePath string) (*source.File, error) {
	jww.INFO.Println("rereading", absFilePath)
	var file *source.File

	reader, err := source.NewLazyFileReader(absFilePath)
	if err != nil {
		return nil, err
	}

	file, err = source.NewFileFromAbs(s.absContentDir(), absFilePath, reader)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *Site) readPagesFromSource() chan error {
	if s.Source == nil {
		panic(fmt.Sprintf("s.Source not set %s", s.absContentDir()))
	}

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

	for _, p := range s.Pages {
		pageChan <- p
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

func (s *Site) addPage(page *Page) {
	if page.ShouldBuild() {
		s.Pages = append(s.Pages, page)
	}

	if page.IsDraft() {
		s.draftCount++
	}

	if page.IsFuture() {
		s.futureCount++
	}
}

func (s *Site) removePageByPath(path string) {
	if i := s.Pages.FindPagePosByFilePath(path); i >= 0 {
		page := s.Pages[i]

		if page.IsDraft() {
			s.draftCount--
		}

		if page.IsFuture() {
			s.futureCount--
		}

		s.Pages = append(s.Pages[:i], s.Pages[i+1:]...)
	}
}

func (s *Site) removePage(page *Page) {
	if i := s.Pages.FindPagePos(page); i >= 0 {
		if page.IsDraft() {
			s.draftCount--
		}

		if page.IsFuture() {
			s.futureCount--
		}

		s.Pages = append(s.Pages[:i], s.Pages[i+1:]...)
	}
}

func (s *Site) replacePage(page *Page) {
	// will find existing page that matches filepath and remove it
	s.removePage(page)
	s.addPage(page)
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

	s.Pages.Sort()
	close(coordinator)

	if len(errMsgs) == 0 {
		errs <- nil
		return
	}
	errs <- fmt.Errorf("Errors reading pages: %s", strings.Join(errMsgs, "\n"))
}

func readCollator(s *Site, results <-chan HandledResult, errs chan<- error) {
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

	s.Pages.Sort()
	if len(errMsgs) == 0 {
		errs <- nil
		return
	}
	errs <- fmt.Errorf("Errors reading pages: %s", strings.Join(errMsgs, "\n"))
}

func (s *Site) buildSiteMeta() (err error) {
	s.assembleMenus()

	if len(s.Pages) == 0 {
		return
	}

	s.assembleTaxonomies()
	s.assembleSections()
	s.Info.LastChange = s.Pages[0].Lastmod

	return
}

func (s *Site) getMenusFromConfig() Menus {

	ret := Menus{}

	if menus := viper.GetStringMap("menu"); menus != nil {
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
	menuEntryURL = helpers.SanitizeURLKeepTrailingSlash(helpers.URLize(menuEntryURL))
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

	sectionPagesMenu := viper.GetString("SectionPagesMenu")
	sectionPagesMenus := make(map[string]interface{})
	//creating flat hash
	for _, p := range s.Pages {

		if sectionPagesMenu != "" {
			if _, ok := sectionPagesMenus[p.Section()]; !ok {
				if p.Section() != "" {
					me := MenuEntry{Identifier: p.Section(), Name: helpers.MakeTitle(helpers.FirstUpper(p.Section())), URL: s.Info.createNodeMenuEntryURL("/" + p.Section())}
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
	s.Sections = make(Taxonomy)

	taxonomies := viper.GetStringMapString("Taxonomies")
	jww.INFO.Printf("found taxonomies: %#v\n", taxonomies)

	for _, plural := range taxonomies {
		s.Taxonomies[plural] = make(Taxonomy)
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
					}
				} else if v, ok := vals.(string); ok {
					x := WeightedPage{weight.(int), p}
					s.Taxonomies[plural].add(v, x, s.Info.preserveTaxonomyNames)
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
	s.Info.Sections = s.Sections
}

// Prepare pages for a new full build.
func (s *Site) resetPageBuildState() {

	s.Info.paginationPageCount = 0

	for _, p := range s.Pages {
		p.scratch = newScratch()

	}
}

func (s *Site) assembleSections() {
	for i, p := range s.Pages {
		s.Sections.add(p.Section(), WeightedPage{s.Pages[i].Weight, s.Pages[i]}, s.Info.preserveTaxonomyNames)
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

func (s *Site) possibleTaxonomies() (taxonomies []string) {
	for _, p := range s.Pages {
		for k := range p.Params {
			if !helpers.InStringArray(taxonomies, k) {
				taxonomies = append(taxonomies, k)
			}
		}
	}
	return
}

// renderAliases renders shell pages that simply have a redirect in the header.
func (s *Site) renderAliases() error {
	for _, p := range s.Pages {
		if len(p.Aliases) == 0 {
			continue
		}

		plink, err := p.Permalink()
		if err != nil {
			return err
		}

		for _, a := range p.Aliases {
			if err := s.writeDestAlias(a, plink); err != nil {
				return err
			}
		}
	}
	return nil
}

// renderPages renders pages each corresponding to a markdown file.
func (s *Site) renderPages() error {

	results := make(chan error)
	pages := make(chan *Page)
	errs := make(chan error)

	go errorCollator(results, errs)

	procs := getGoMaxProcs()

	// this cannot be fanned out to multiple Go routines
	// See issue #1601
	// TODO(bep): Check the IsRenderable logic.
	for _, p := range s.Pages {
		var layouts []string
		if !p.IsRenderable() {
			self := "__" + p.TargetPath()
			_, err := s.Tmpl.GetClone().New(self).Parse(string(p.Content))
			if err != nil {
				results <- err
				continue
			}
			layouts = append(layouts, self)
		} else {
			layouts = append(layouts, p.layouts()...)
			layouts = append(layouts, "_default/single.html")
		}
		p.layoutsCalculated = layouts
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < procs*4; i++ {
		wg.Add(1)
		go pageRenderer(s, pages, results, wg)
	}

	for _, page := range s.Pages {
		pages <- page
	}

	close(pages)

	wg.Wait()

	close(results)

	err := <-errs
	if err != nil {
		return fmt.Errorf("Error(s) rendering pages: %s", err)
	}
	return nil
}

func pageRenderer(s *Site, pages <-chan *Page, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range pages {
		err := s.renderAndWritePage("page "+p.FullFilePath(), p.TargetPath(), p, s.appendThemeTemplates(p.layouts())...)
		if err != nil {
			results <- err
		}
	}
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

type taxRenderInfo struct {
	key      string
	pages    WeightedPages
	singular string
	plural   string
}

// renderTaxonomiesLists renders the listing pages based on the meta data
// each unique term within a taxonomy will have a page created
func (s *Site) renderTaxonomiesLists() error {
	wg := &sync.WaitGroup{}

	taxes := make(chan taxRenderInfo)
	results := make(chan error)

	procs := getGoMaxProcs()

	for i := 0; i < procs*4; i++ {
		wg.Add(1)
		go taxonomyRenderer(s, taxes, results, wg)
	}

	errs := make(chan error)

	go errorCollator(results, errs)

	taxonomies := viper.GetStringMapString("Taxonomies")
	for singular, plural := range taxonomies {
		for key, pages := range s.Taxonomies[plural] {
			taxes <- taxRenderInfo{key, pages, singular, plural}
		}
	}
	close(taxes)

	wg.Wait()

	close(results)

	err := <-errs
	if err != nil {
		return fmt.Errorf("Error(s) rendering taxonomies: %s", err)
	}
	return nil
}

func (s *Site) newTaxonomyNode(t taxRenderInfo) (*Node, string) {
	key := t.key
	n := s.newNode()
	if s.Info.preserveTaxonomyNames {
		key = helpers.MakePathSanitized(key)
		// keep as is in the title
		n.Title = t.key
	} else {
		n.Title = strings.Replace(strings.Title(t.key), "-", " ", -1)
	}
	base := t.plural + "/" + key
	s.setURLs(n, base)
	if len(t.pages) > 0 {
		n.Date = t.pages[0].Page.Date
		n.Lastmod = t.pages[0].Page.Lastmod
	}
	n.Data[t.singular] = t.pages
	n.Data["Singular"] = t.singular
	n.Data["Plural"] = t.plural
	n.Data["Pages"] = t.pages.Pages()
	return n, base
}

func taxonomyRenderer(s *Site, taxes <-chan taxRenderInfo, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	var n *Node

	for t := range taxes {

		var base string
		layouts := s.appendThemeTemplates(
			[]string{"taxonomy/" + t.singular + ".html", "indexes/" + t.singular + ".html", "_default/taxonomy.html", "_default/list.html"})

		n, base = s.newTaxonomyNode(t)

		dest := base
		if viper.GetBool("UglyURLs") {
			dest = helpers.Uglify(base + ".html")
		} else {
			dest = helpers.PrettifyPath(base + "/index.html")
		}

		if err := s.renderAndWritePage("taxonomy "+t.singular, dest, n, layouts...); err != nil {
			results <- err
			continue
		}

		if n.paginator != nil {

			paginatePath := viper.GetString("paginatePath")

			// write alias for page 1
			s.writeDestAlias(helpers.PaginateAliasPath(base, 1), s.permalink(base))

			pagers := n.paginator.Pagers()

			for i, pager := range pagers {
				if i == 0 {
					// already created
					continue
				}

				taxonomyPagerNode, _ := s.newTaxonomyNode(t)
				taxonomyPagerNode.paginator = pager
				if pager.TotalPages() > 0 {
					first, _ := pager.page(0)
					taxonomyPagerNode.Date = first.Date
					taxonomyPagerNode.Lastmod = first.Lastmod
				}
				pageNumber := i + 1
				htmlBase := fmt.Sprintf("/%s/%s/%d", base, paginatePath, pageNumber)
				if err := s.renderAndWritePage(fmt.Sprintf("taxonomy %s", t.singular), htmlBase, taxonomyPagerNode, layouts...); err != nil {
					results <- err
					continue
				}
			}
		}

		if !viper.GetBool("DisableRSS") {
			// XML Feed
			rssuri := viper.GetString("RSSUri")
			n.URL = s.permalinkStr(base + "/" + rssuri)
			n.Permalink = s.permalink(base)
			rssLayouts := []string{"taxonomy/" + t.singular + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}

			if err := s.renderAndWriteXML("taxonomy "+t.singular+" rss", base+"/"+rssuri, n, s.appendThemeTemplates(rssLayouts)...); err != nil {
				results <- err
				continue
			}
		}
	}
}

// renderListsOfTaxonomyTerms renders a page per taxonomy that lists the terms for that taxonomy
func (s *Site) renderListsOfTaxonomyTerms() (err error) {
	taxonomies := viper.GetStringMapString("Taxonomies")
	for singular, plural := range taxonomies {
		n := s.newNode()
		n.Title = strings.Title(plural)
		s.setURLs(n, plural)
		n.Data["Singular"] = singular
		n.Data["Plural"] = plural
		n.Data["Terms"] = s.Taxonomies[plural]
		// keep the following just for legacy reasons
		n.Data["OrderedIndex"] = n.Data["Terms"]
		n.Data["Index"] = n.Data["Terms"]
		layouts := []string{"taxonomy/" + singular + ".terms.html", "_default/terms.html", "indexes/indexes.html"}
		layouts = s.appendThemeTemplates(layouts)
		if s.layoutExists(layouts...) {
			if err := s.renderAndWritePage("taxonomy terms for "+singular, plural+"/index.html", n, layouts...); err != nil {
				return err
			}
		}
	}

	return
}

func (s *Site) newSectionListNode(sectionName, section string, data WeightedPages) *Node {
	n := s.newNode()
	sectionName = helpers.FirstUpper(sectionName)
	if viper.GetBool("PluralizeListTitles") {
		n.Title = inflect.Pluralize(sectionName)
	} else {
		n.Title = sectionName
	}
	s.setURLs(n, section)
	n.Date = data[0].Page.Date
	n.Lastmod = data[0].Page.Lastmod
	n.Data["Pages"] = data.Pages()

	return n
}

// renderSectionLists renders a page for each section
func (s *Site) renderSectionLists() error {
	for section, data := range s.Sections {
		// section keys can be lower case (depending on site.pathifyTaxonomyKeys)
		// extract the original casing from the first page to get sensible titles.
		sectionName := section
		if !s.Info.preserveTaxonomyNames && len(data) > 0 {
			sectionName = data[0].Page.Section()
		}
		layouts := s.appendThemeTemplates(
			[]string{"section/" + section + ".html", "_default/section.html", "_default/list.html", "indexes/" + section + ".html", "_default/indexes.html"})

		if s.Info.preserveTaxonomyNames {
			section = helpers.MakePathSanitized(section)
		}

		n := s.newSectionListNode(sectionName, section, data)
		if err := s.renderAndWritePage(fmt.Sprintf("section %s", section), section, n, s.appendThemeTemplates(layouts)...); err != nil {
			return err
		}

		if n.paginator != nil {

			paginatePath := viper.GetString("paginatePath")

			// write alias for page 1
			s.writeDestAlias(helpers.PaginateAliasPath(section, 1), s.permalink(section))

			pagers := n.paginator.Pagers()

			for i, pager := range pagers {
				if i == 0 {
					// already created
					continue
				}

				sectionPagerNode := s.newSectionListNode(sectionName, section, data)
				sectionPagerNode.paginator = pager
				if pager.TotalPages() > 0 {
					first, _ := pager.page(0)
					sectionPagerNode.Date = first.Date
					sectionPagerNode.Lastmod = first.Lastmod
				}
				pageNumber := i + 1
				htmlBase := fmt.Sprintf("/%s/%s/%d", section, paginatePath, pageNumber)
				if err := s.renderAndWritePage(fmt.Sprintf("section %s", section), filepath.FromSlash(htmlBase), sectionPagerNode, layouts...); err != nil {
					return err
				}
			}
		}

		if !viper.GetBool("DisableRSS") && section != "" {
			// XML Feed
			rssuri := viper.GetString("RSSUri")
			n.URL = s.permalinkStr(section + "/" + rssuri)
			n.Permalink = s.permalink(section)
			rssLayouts := []string{"section/" + section + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}
			if err := s.renderAndWriteXML("section "+section+" rss", section+"/"+rssuri, n, s.appendThemeTemplates(rssLayouts)...); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Site) newHomeNode() *Node {
	n := s.newNode()
	n.Title = n.Site.Title
	n.IsHome = true
	s.setURLs(n, "/")
	n.Data["Pages"] = s.Pages
	if len(s.Pages) != 0 {
		n.Date = s.Pages[0].Date
		n.Lastmod = s.Pages[0].Lastmod
	}
	return n
}

func (s *Site) renderHomePage() error {
	n := s.newHomeNode()
	layouts := s.appendThemeTemplates([]string{"index.html", "_default/list.html"})

	if err := s.renderAndWritePage("homepage", helpers.FilePathSeparator, n, layouts...); err != nil {
		return err
	}

	if n.paginator != nil {

		paginatePath := viper.GetString("paginatePath")

		// write alias for page 1
		s.writeDestAlias(helpers.PaginateAliasPath("", 1), s.permalink("/"))

		pagers := n.paginator.Pagers()

		for i, pager := range pagers {
			if i == 0 {
				// already created
				continue
			}

			homePagerNode := s.newHomeNode()
			homePagerNode.paginator = pager
			if pager.TotalPages() > 0 {
				first, _ := pager.page(0)
				homePagerNode.Date = first.Date
				homePagerNode.Lastmod = first.Lastmod
			}
			pageNumber := i + 1
			htmlBase := fmt.Sprintf("/%s/%d", paginatePath, pageNumber)
			if err := s.renderAndWritePage(fmt.Sprintf("homepage"), filepath.FromSlash(htmlBase), homePagerNode, layouts...); err != nil {
				return err
			}
		}
	}

	if !viper.GetBool("DisableRSS") {
		// XML Feed
		n.URL = s.permalinkStr(viper.GetString("RSSUri"))
		n.Title = ""
		high := 50
		if len(s.Pages) < high {
			high = len(s.Pages)
		}
		n.Data["Pages"] = s.Pages[:high]
		if len(s.Pages) > 0 {
			n.Date = s.Pages[0].Date
			n.Lastmod = s.Pages[0].Lastmod
		}

		rssLayouts := []string{"rss.xml", "_default/rss.xml", "_internal/_default/rss.xml"}

		if err := s.renderAndWriteXML("homepage rss", viper.GetString("RSSUri"), n, s.appendThemeTemplates(rssLayouts)...); err != nil {
			return err
		}
	}

	// TODO(bep) reusing the Home Node smells trouble
	n.URL = helpers.URLize("404.html")
	n.IsHome = false
	n.Title = "404 Page not found"
	n.Permalink = s.permalink("404.html")
	n.scratch = newScratch()

	nfLayouts := []string{"404.html"}
	if nfErr := s.renderAndWritePage("404 page", "404.html", n, s.appendThemeTemplates(nfLayouts)...); nfErr != nil {
		return nfErr
	}

	return nil
}

func (s *Site) renderSitemap() error {
	if viper.GetBool("DisableSitemap") {
		return nil
	}

	sitemapDefault := parseSitemap(viper.GetStringMap("Sitemap"))

	n := s.newNode()

	// Prepend homepage to the list of pages
	pages := make(Pages, 0)

	page := &Page{}
	page.Date = s.Info.LastChange
	page.Lastmod = s.Info.LastChange
	page.Site = &s.Info
	page.URL = "/"
	page.Sitemap.ChangeFreq = sitemapDefault.ChangeFreq
	page.Sitemap.Priority = sitemapDefault.Priority

	pages = append(pages, page)
	pages = append(pages, s.Pages...)

	n.Data["Pages"] = pages

	for _, page := range pages {
		if page.Sitemap.ChangeFreq == "" {
			page.Sitemap.ChangeFreq = sitemapDefault.ChangeFreq
		}

		if page.Sitemap.Priority == -1 {
			page.Sitemap.Priority = sitemapDefault.Priority
		}

		if page.Sitemap.Filename == "" {
			page.Sitemap.Filename = sitemapDefault.Filename
		}
	}

	smLayouts := []string{"sitemap.xml", "_default/sitemap.xml", "_internal/_default/sitemap.xml"}

	if err := s.renderAndWriteXML("sitemap", page.Sitemap.Filename, n, s.appendThemeTemplates(smLayouts)...); err != nil {
		return err
	}

	return nil
}

func (s *Site) renderRobotsTXT() error {
	if !viper.GetBool("EnableRobotsTXT") {
		return nil
	}

	n := s.newNode()
	n.Data["Pages"] = s.Pages

	rLayouts := []string{"robots.txt", "_default/robots.txt", "_internal/_default/robots.txt"}
	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)
	err := s.render("robots", n, outBuffer, s.appendThemeTemplates(rLayouts)...)

	if err == nil {
		err = s.writeDestFile("robots.txt", outBuffer)
	}

	return err
}

// Stats prints Hugo builds stats to the console.
// This is what you see after a successful hugo build.
func (s *Site) Stats() {
	jww.FEEDBACK.Println(s.draftStats())
	jww.FEEDBACK.Println(s.futureStats())
	jww.FEEDBACK.Printf("%d pages created\n", len(s.Pages))
	jww.FEEDBACK.Printf("%d non-page files copied\n", len(s.Files))
	jww.FEEDBACK.Printf("%d paginator pages created\n", s.Info.paginationPageCount)
	taxonomies := viper.GetStringMapString("Taxonomies")

	for _, pl := range taxonomies {
		jww.FEEDBACK.Printf("%d %s created\n", len(s.Taxonomies[pl]), pl)
	}
}

func (s *Site) setURLs(n *Node, in string) {
	n.URL = helpers.URLizeAndPrep(in)
	n.Permalink = s.permalink(n.URL)
	n.RSSLink = template.HTML(s.permalink(in + ".xml"))
}

func (s *Site) permalink(plink string) string {
	return s.permalinkStr(plink)
}

func (s *Site) permalinkStr(plink string) string {
	return helpers.MakePermalink(viper.GetString("BaseURL"), helpers.URLizeAndPrep(plink)).String()
}

func (s *Site) newNode() *Node {
	return &Node{
		Data: make(map[string]interface{}),
		Site: &s.Info,
	}
}

func (s *Site) layoutExists(layouts ...string) bool {
	_, found := s.findFirstLayout(layouts...)

	return found
}

func (s *Site) renderAndWriteXML(name string, dest string, d interface{}, layouts ...string) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)
	renderBuffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n")

	err := s.render(name, d, renderBuffer, layouts...)

	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)

	var path []byte
	if viper.GetBool("RelativeURLs") {
		path = []byte(helpers.GetDottedRelativePath(dest))
	} else {
		s := viper.GetString("BaseURL")
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		path = []byte(s)
	}
	transformer := transform.NewChain(transform.AbsURLInXML)
	transformer.Apply(outBuffer, renderBuffer, path)

	if err == nil {
		err = s.writeDestFile(dest, outBuffer)
	}

	return err
}

func (s *Site) renderAndWritePage(name string, dest string, d interface{}, layouts ...string) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	err := s.render(name, d, renderBuffer, layouts...)

	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)

	var pageTarget target.Output

	if p, ok := d.(*Page); ok && path.Ext(p.URL) != "" {
		// user has explicitly set a URL with extension for this page
		// make sure it sticks even if "ugly URLs" are turned off.
		pageTarget = s.pageUglyTarget()
	} else {
		pageTarget = s.pageTarget()
	}

	transformLinks := transform.NewEmptyTransforms()

	if viper.GetBool("RelativeURLs") || viper.GetBool("CanonifyURLs") {
		transformLinks = append(transformLinks, transform.AbsURL)
	}

	if s.running() && viper.GetBool("watch") && !viper.GetBool("DisableLiveReload") {
		transformLinks = append(transformLinks, transform.LiveReloadInject)
	}

	var path []byte

	if viper.GetBool("RelativeURLs") {
		translated, err := pageTarget.(target.OptionalTranslator).TranslateRelative(dest)
		if err != nil {
			return err
		}
		path = []byte(helpers.GetDottedRelativePath(translated))
	} else if viper.GetBool("CanonifyURLs") {
		s := viper.GetString("BaseURL")
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		path = []byte(s)
	}

	transformer := transform.NewChain(transformLinks...)
	transformer.Apply(outBuffer, renderBuffer, path)

	if outBuffer.Len() == 0 {
		jww.WARN.Printf("%q is rendered empty\n", dest)
		if dest == "/" {
			jww.FEEDBACK.Println("=============================================================")
			jww.FEEDBACK.Println("Your rendered home page is blank: /index.html is zero-length")
			jww.FEEDBACK.Println(" * Did you specify a theme on the command-line or in your")
			jww.FEEDBACK.Printf("   %q file?  (Current theme: %q)\n", filepath.Base(viper.ConfigFileUsed()), viper.GetString("Theme"))
			if !viper.GetBool("Verbose") {
				jww.FEEDBACK.Println(" * For more debugging information, run \"hugo -v\"")
			}
			jww.FEEDBACK.Println("=============================================================")
		}
	}

	if err == nil {
		if err = s.writeDestPage(dest, pageTarget, outBuffer); err != nil {
			return err
		}
	}
	return err
}

func (s *Site) render(name string, d interface{}, w io.Writer, layouts ...string) error {
	layout, found := s.findFirstLayout(layouts...)
	if found == false {
		jww.WARN.Printf("Unable to locate layout for %s: %s\n", name, layouts)
		return nil
	}

	if err := s.renderThing(d, layout, w); err != nil {
		// Behavior here should be dependent on if running in server or watch mode.
		distinctErrorLogger.Printf("Error while rendering %s: %v", name, err)
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
		if s.Tmpl.Lookup(layout) != nil {
			return layout, true
		}
	}
	return "", false
}

func (s *Site) renderThing(d interface{}, layout string, w io.Writer) error {
	// If the template doesn't exist, then return, but leave the Writer open
	if templ := s.Tmpl.Lookup(layout); templ != nil {
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

func (s *Site) initTargetList() {
	s.targetListInit.Do(func() {
		if s.targets.page == nil {
			s.targets.page = &target.PagePub{
				PublishDir: s.absPublishDir(),
				UglyURLs:   viper.GetBool("UglyURLs"),
			}
		}
		if s.targets.pageUgly == nil {
			s.targets.pageUgly = &target.PagePub{
				PublishDir: s.absPublishDir(),
				UglyURLs:   true,
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

func (s *Site) writeDestAlias(path string, permalink string) (err error) {
	jww.DEBUG.Println("creating alias:", path)
	return s.aliasTarget().Publish(path, permalink)
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

	if viper.GetBool("BuildDrafts") {
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
		msg = fmt.Sprintf("%d future rendered", s.draftCount)
	}

	if viper.GetBool("BuildFuture") {
		return fmt.Sprintf("%d of ", s.futureCount) + msg
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
