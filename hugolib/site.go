// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"sync/atomic"

	"bitbucket.org/pkg/inflect"
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

var DefaultTimer *nitro.B

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
	Targets        targetList
	targetListInit sync.Once
	RunMode        runmode
	params         map[string]interface{}
	draftCount     int
	futureCount    int
	Data           map[string]interface{}
}

type targetList struct {
	Page  target.Output
	File  target.Output
	Alias target.AliasPublisher
}

type SiteInfo struct {
	BaseURL             template.URL
	Taxonomies          TaxonomyList
	Authors             AuthorList
	Social              SiteSocial
	Sections            Taxonomy
	Pages               *Pages
	Files               []*source.File
	Menus               *Menus
	Hugo                *HugoInfo
	Title               string
	Author              map[string]interface{}
	LanguageCode        string
	DisqusShortname     string
	Copyright           string
	LastChange          time.Time
	Permalinks          PermalinkOverrides
	Params              map[string]interface{}
	BuildDrafts         bool
	canonifyURLs        bool
	paginationPageCount uint64
	Data                *map[string]interface{}
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

// BaseUrl is deprecated. Will be removed in 0.15.
func (s *SiteInfo) BaseUrl() template.URL {
	helpers.Deprecated("Site", ".BaseUrl", ".BaseURL")
	return s.BaseURL
}

// Recent is deprecated. Will be removed in 0.15.
func (s *SiteInfo) Recent() *Pages {
	helpers.Deprecated("Site", ".Recent", ".Pages")
	return s.Pages
}

// Indexes is deprecated. Will be removed in 0.15.
func (s *SiteInfo) Indexes() *TaxonomyList {
	helpers.Deprecated("Site", ".Indexes", ".Taxonomies")
	return &s.Taxonomies
}

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

func (s *SiteInfo) Ref(ref string, page *Page) (string, error) {
	return s.refLink(ref, page, false)
}

func (s *SiteInfo) RelRef(ref string, page *Page) (string, error) {
	return s.refLink(ref, page, true)
}

func (s *SiteInfo) addToPaginationPageCount(cnt uint64) {
	atomic.AddUint64(&s.paginationPageCount, cnt)
}

type runmode struct {
	Watching bool
}

func (s *Site) Running() bool {
	return s.RunMode.Watching
}

func init() {
	DefaultTimer = nitro.Initalize()
}

func (s *Site) timerStep(step string) {
	if s.timer == nil {
		s.timer = DefaultTimer
	}
	s.timer.Step(step)
}

func (s *Site) Build() (err error) {
	if err = s.Process(); err != nil {
		return
	}
	if err = s.Render(); err != nil {
		jww.ERROR.Printf("Error rendering site: %s\nAvailable templates:\n", err)
		for _, template := range s.Tmpl.Templates() {
			jww.ERROR.Printf("\t%s\n", template.Name())
		}
		return
	}
	return nil
}

func (s *Site) Analyze() {
	s.Process()
	s.ShowPlan(os.Stdout)
}

func (s *Site) prepTemplates() {
	s.Tmpl = tpl.InitializeT()
	s.Tmpl.LoadTemplates(s.absLayoutDir())
	if s.hasTheme() {
		s.Tmpl.LoadTemplatesWithPrefix(s.absThemeDir()+"/layouts", "theme")
	}
}

func (s *Site) addTemplate(name, data string) error {
	return s.Tmpl.AddTemplate(name, data)
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

func (s *Site) Process() (err error) {
	if err = s.initialize(); err != nil {
		return
	}
	s.prepTemplates()
	s.Tmpl.PrintErrors()
	s.timerStep("initialize & template prep")

	dataSources := make([]source.Input, 0, 2)

	dataSources = append(dataSources, &source.Filesystem{Base: s.absDataDir()})

	// have to be last - duplicate keys in earlier entries will win
	themeStaticDir, err := helpers.GetThemeDataDirPath()
	if err == nil {
		dataSources = append(dataSources, &source.Filesystem{Base: themeStaticDir})
	}

	if err = s.loadData(dataSources); err != nil {
		return
	}
	s.timerStep("load data")

	if err = s.CreatePages(); err != nil {
		return
	}
	s.setupPrevNext()
	s.timerStep("import pages")
	if err = s.BuildSiteMeta(); err != nil {
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
	if err = s.RenderAliases(); err != nil {
		return
	}
	s.timerStep("render and write aliases")
	if err = s.RenderTaxonomiesLists(); err != nil {
		return
	}
	s.timerStep("render and write taxonomies")
	s.RenderListsOfTaxonomyTerms()
	s.timerStep("render & write taxonomy lists")
	if err = s.RenderSectionLists(); err != nil {
		return
	}
	s.timerStep("render and write lists")
	if err = s.RenderPages(); err != nil {
		return
	}
	s.timerStep("render and write pages")
	if err = s.RenderHomePage(); err != nil {
		return
	}
	s.timerStep("render and write homepage")
	if err = s.RenderSitemap(); err != nil {
		return
	}
	s.timerStep("render and write Sitemap")
	return
}

func (s *Site) checkDescriptions() {
	for _, p := range s.Pages {
		if len(p.Description) < 60 {
			jww.FEEDBACK.Println(p.Source.Path() + " ")
		}
	}
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
		permalinks[k] = PathPattern(v)
	}

	s.Info = SiteInfo{
		BaseURL:         template.URL(helpers.SanitizeURLKeepTrailingSlash(viper.GetString("BaseURL"))),
		Title:           viper.GetString("Title"),
		Author:          viper.GetStringMap("author"),
		LanguageCode:    viper.GetString("languagecode"),
		Copyright:       viper.GetString("copyright"),
		DisqusShortname: viper.GetString("DisqusShortname"),
		BuildDrafts:     viper.GetBool("BuildDrafts"),
		canonifyURLs:    viper.GetBool("CanonifyURLs"),
		Pages:           &s.Pages,
		Menus:           &s.Menus,
		Params:          params,
		Permalinks:      permalinks,
		Data:            &s.Data,
	}
}

func (s *Site) hasTheme() bool {
	return viper.GetString("theme") != ""
}

func (s *Site) absDataDir() string {
	return helpers.AbsPathify(viper.GetString("DataDir"))
}

func (s *Site) absThemeDir() string {
	return helpers.AbsPathify("themes/" + viper.GetString("theme"))
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
	if b, _ := helpers.DirExists(s.absContentDir(), hugofs.SourceFs); !b {
		return fmt.Errorf("No source directory found, expecting to find it at " + s.absContentDir())
	}
	return
}

type pageResult struct {
	page *Page
	err  error
}

func (s *Site) CreatePages() error {
	if s.Source == nil {
		panic(fmt.Sprintf("s.Source not set %s", s.absContentDir()))
	}
	if len(s.Source.Files()) < 1 {
		return nil
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

	errs := make(chan error)

	// we can only have exactly one result collator, since it makes changes that
	// must be synchronized.
	go readCollator(s, results, errs)

	for _, file := range files {
		filechan <- file
	}

	close(filechan)

	wg.Wait()

	close(results)

	readErrs := <-errs

	results = make(chan HandledResult)
	pageChan := make(chan *Page)
	fileConvChan := make(chan *source.File)

	wg = &sync.WaitGroup{}

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

	renderErrs := <-errs

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
		h := NewMetaHandler(file.Extension())
		if h != nil {
			h.Read(file, s, results)
		} else {
			jww.ERROR.Println("Unsupported File Type", file.Path())
		}
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
			if r.page.ShouldBuild() {
				s.Pages = append(s.Pages, r.page)
			}

			if r.page.IsDraft() {
				s.draftCount++
			}

			if r.page.IsFuture() {
				s.futureCount++
			}
		}
	}

	s.Pages.Sort()
	if len(errMsgs) == 0 {
		errs <- nil
		return
	}
	errs <- fmt.Errorf("Errors reading pages: %s", strings.Join(errMsgs, "\n"))
}

func (s *Site) BuildSiteMeta() (err error) {

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

					menuEntry.MarshallMap(ime)
					menuEntry.URL = s.Info.createNodeMenuEntryURL(menuEntry.URL)

					if ret[name] == nil {
						ret[name] = &Menu{}
					}
					*ret[name] = ret[name].Add(&menuEntry)
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
	menuEntryURL = helpers.URLizeAndPrep(menuEntryURL)
	if !s.canonifyURLs {
		menuEntryURL = helpers.AddContextRoot(string(s.BaseURL), menuEntryURL)
	}
	return menuEntryURL
}

func (s *Site) assembleMenus() {

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
					me := MenuEntry{Identifier: p.Section(), Name: helpers.MakeTitle(p.Section()), URL: s.Info.createNodeMenuEntryURL("/" + p.Section())}
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
			children[twoD{e.Menu, e.Parent}] = children[twoD{e.Menu, e.Parent}].Add(e)
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
			*s.Menus[menu.MenuName] = s.Menus[menu.MenuName].Add(e)
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
			vals := p.GetParam(plural)
			weight := p.GetParam(plural + "_weight")
			if weight == nil {
				weight = 0
			}

			if vals != nil {
				if v, ok := vals.([]string); ok {
					for _, idx := range v {
						x := WeightedPage{weight.(int), p}
						s.Taxonomies[plural].Add(idx, x)
					}
				} else if v, ok := vals.(string); ok {
					x := WeightedPage{weight.(int), p}
					s.Taxonomies[plural].Add(v, x)
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

func (s *Site) assembleSections() {
	for i, p := range s.Pages {
		s.Sections.Add(p.Section(), WeightedPage{s.Pages[i].Weight, s.Pages[i]})
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

// RenderAliases renders shell pages that simply have a redirect in the header
func (s *Site) RenderAliases() error {
	for _, p := range s.Pages {
		for _, a := range p.Aliases {
			plink, err := p.Permalink()
			if err != nil {
				return err
			}
			if err := s.WriteDestAlias(a, template.HTML(plink)); err != nil {
				return err
			}
		}
	}
	return nil
}

// RenderPages renders pages each corresponding to a markdown file
func (s *Site) RenderPages() error {

	results := make(chan error)
	pages := make(chan *Page)

	procs := getGoMaxProcs()

	wg := &sync.WaitGroup{}

	for i := 0; i < procs*4; i++ {
		wg.Add(1)
		go pageRenderer(s, pages, results, wg)
	}

	errs := make(chan error)

	go errorCollator(results, errs)

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
		var layouts []string

		if !p.IsRenderable() {
			self := "__" + p.TargetPath()
			_, err := s.Tmpl.New(self).Parse(string(p.Content))
			if err != nil {
				results <- err
				continue
			}
			layouts = append(layouts, self)
		} else {
			layouts = append(layouts, p.Layout()...)
			layouts = append(layouts, "_default/single.html")
		}

		err := s.renderAndWritePage("page "+p.FullFilePath(), p.TargetPath(), p, s.appendThemeTemplates(layouts)...)
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
	if s.hasTheme() {
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
	return in
}

type taxRenderInfo struct {
	key      string
	pages    WeightedPages
	singular string
	plural   string
}

// RenderTaxonomiesLists renders the listing pages based on the meta data
// each unique term within a taxonomy will have a page created
func (s *Site) RenderTaxonomiesLists() error {
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
	base := t.plural + "/" + t.key
	n := s.NewNode()
	n.Title = strings.Replace(strings.Title(t.key), "-", " ", -1)
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

		if err := s.renderAndWritePage("taxononomy "+t.singular, base, n, layouts...); err != nil {
			results <- err
			continue
		}

		if n.paginator != nil {

			paginatePath := viper.GetString("paginatePath")

			// write alias for page 1
			s.WriteDestAlias(helpers.PaginateAliasPath(base, 1), s.permalink(base))

			pagers := n.paginator.Pagers()

			for i, pager := range pagers {
				if i == 0 {
					// already created
					continue
				}

				taxonomyPagerNode, _ := s.newTaxonomyNode(t)
				taxonomyPagerNode.paginator = pager
				if pager.TotalPages() > 0 {
					taxonomyPagerNode.Date = pager.Pages()[0].Date
					taxonomyPagerNode.Lastmod = pager.Pages()[0].Lastmod
				}
				pageNumber := i + 1
				htmlBase := fmt.Sprintf("/%s/%s/%d", base, paginatePath, pageNumber)
				if err := s.renderAndWritePage(fmt.Sprintf("taxononomy %s", t.singular), htmlBase, taxonomyPagerNode, layouts...); err != nil {
					results <- err
					continue
				}
			}
		}

		if !viper.GetBool("DisableRSS") {
			// XML Feed
			n.URL = s.permalinkStr(base + "/index.xml")
			n.Permalink = s.permalink(base)
			rssLayouts := []string{"taxonomy/" + t.singular + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}

			if err := s.renderAndWriteXML("taxonomy "+t.singular+" rss", base+"/index.xml", n, s.appendThemeTemplates(rssLayouts)...); err != nil {
				results <- err
				continue
			}
		}
	}
}

// RenderListsOfTaxonomyTerms renders a page per taxonomy that lists the terms for that taxonomy
func (s *Site) RenderListsOfTaxonomyTerms() (err error) {
	taxonomies := viper.GetStringMapString("Taxonomies")
	for singular, plural := range taxonomies {
		n := s.NewNode()
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

func (s *Site) newSectionListNode(section string, data WeightedPages) *Node {
	n := s.NewNode()
	if viper.GetBool("PluralizeListTitles") {
		n.Title = strings.Title(inflect.Pluralize(section))
	} else {
		n.Title = strings.Title(section)
	}
	s.setURLs(n, section)
	n.Date = data[0].Page.Date
	n.Lastmod = data[0].Page.Lastmod
	n.Data["Pages"] = data.Pages()

	return n
}

// RenderSectionLists renders a page for each section
func (s *Site) RenderSectionLists() error {
	for section, data := range s.Sections {

		layouts := s.appendThemeTemplates(
			[]string{"section/" + section + ".html", "_default/section.html", "_default/list.html", "indexes/" + section + ".html", "_default/indexes.html"})

		n := s.newSectionListNode(section, data)

		if err := s.renderAndWritePage(fmt.Sprintf("section %s", section), fmt.Sprintf("/%s", section), n, s.appendThemeTemplates(layouts)...); err != nil {
			return err
		}

		if n.paginator != nil {

			paginatePath := viper.GetString("paginatePath")

			// write alias for page 1
			s.WriteDestAlias(helpers.PaginateAliasPath(section, 1), s.permalink(section))

			pagers := n.paginator.Pagers()

			for i, pager := range pagers {
				if i == 0 {
					// already created
					continue
				}

				sectionPagerNode := s.newSectionListNode(section, data)
				sectionPagerNode.paginator = pager
				if pager.TotalPages() > 0 {
					sectionPagerNode.Date = pager.Pages()[0].Date
					sectionPagerNode.Lastmod = pager.Pages()[0].Lastmod
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
			n.URL = s.permalinkStr(section + "/index.xml")
			n.Permalink = s.permalink(section)
			rssLayouts := []string{"section/" + section + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}
			if err := s.renderAndWriteXML("section "+section+" rss", section+"/index.xml", n, s.appendThemeTemplates(rssLayouts)...); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Site) newHomeNode() *Node {
	n := s.NewNode()
	n.Title = n.Site.Title
	s.setURLs(n, "/")
	n.Data["Pages"] = s.Pages
	return n
}

func (s *Site) RenderHomePage() error {
	n := s.newHomeNode()
	layouts := s.appendThemeTemplates([]string{"index.html", "_default/list.html"})

	if err := s.renderAndWritePage("homepage", helpers.FilePathSeparator, n, layouts...); err != nil {
		return err
	}

	if n.paginator != nil {

		paginatePath := viper.GetString("paginatePath")

		// write alias for page 1
		s.WriteDestAlias(helpers.PaginateAliasPath("", 1), s.permalink("/"))

		pagers := n.paginator.Pagers()

		for i, pager := range pagers {
			if i == 0 {
				// already created
				continue
			}

			homePagerNode := s.newHomeNode()
			homePagerNode.paginator = pager
			if pager.TotalPages() > 0 {
				homePagerNode.Date = pager.Pages()[0].Date
				homePagerNode.Lastmod = pager.Pages()[0].Lastmod
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

	n.URL = helpers.URLize("404.html")
	n.Title = "404 Page not found"
	n.Permalink = s.permalink("404.html")

	nfLayouts := []string{"404.html"}
	if nfErr := s.renderAndWritePage("404 page", "404.html", n, s.appendThemeTemplates(nfLayouts)...); nfErr != nil {
		return nfErr
	}

	return nil
}

func (s *Site) RenderSitemap() error {
	if viper.GetBool("DisableSitemap") {
		return nil
	}

	sitemapDefault := parseSitemap(viper.GetStringMap("Sitemap"))

	n := s.NewNode()

	// Prepend homepage to the list of pages
	pages := make(Pages, 0)

	page := &Page{}
	page.Date = s.Info.LastChange
	page.Lastmod = s.Info.LastChange
	page.Site = &s.Info
	page.URL = "/"

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
	}

	smLayouts := []string{"sitemap.xml", "_default/sitemap.xml", "_internal/_default/sitemap.xml"}

	if err := s.renderAndWriteXML("sitemap", "sitemap.xml", n, s.appendThemeTemplates(smLayouts)...); err != nil {
		return err
	}

	return nil
}

func (s *Site) Stats() {
	jww.FEEDBACK.Println(s.draftStats())
	jww.FEEDBACK.Println(s.futureStats())
	jww.FEEDBACK.Printf("%d pages created\n", len(s.Pages))
	jww.FEEDBACK.Printf("%d paginator pages created\n", s.Info.paginationPageCount)
	taxonomies := viper.GetStringMapString("Taxonomies")

	for _, pl := range taxonomies {
		jww.FEEDBACK.Printf("%d %s created\n", len(s.Taxonomies[pl]), pl)
	}
}

func (s *Site) setURLs(n *Node, in string) {
	n.URL = helpers.URLizeAndPrep(in)
	n.Permalink = s.permalink(n.URL)
	n.RSSLink = s.permalink(in + ".xml")
}

func (s *Site) permalink(plink string) template.HTML {
	return template.HTML(s.permalinkStr(plink))
}

func (s *Site) permalinkStr(plink string) string {
	return helpers.MakePermalink(string(viper.GetString("BaseURL")), helpers.URLizeAndPrep(plink)).String()
}

func (s *Site) NewNode() *Node {
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
		err = s.WriteDestFile(dest, outBuffer)
	}

	return err
}

func (s *Site) renderAndWritePage(name string, dest string, d interface{}, layouts ...string) error {
	renderBuffer := bp.GetBuffer()
	defer bp.PutBuffer(renderBuffer)

	err := s.render(name, d, renderBuffer, layouts...)

	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)

	transformLinks := transform.NewEmptyTransforms()

	if viper.GetBool("RelativeURLs") || viper.GetBool("CanonifyURLs") {
		transformLinks = append(transformLinks, transform.AbsURL)
	}

	if viper.GetBool("watch") && !viper.GetBool("DisableLiveReload") {
		transformLinks = append(transformLinks, transform.LiveReloadInject)
	}

	var path []byte

	if viper.GetBool("RelativeURLs") {
		translated, err := s.PageTarget().(target.OptionalTranslator).TranslateRelative(dest)
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

	if err == nil {
		if err = s.WriteDestPage(dest, outBuffer); err != nil {
			return err
		}
	}
	return err
}

func (s *Site) render(name string, d interface{}, renderBuffer *bytes.Buffer, layouts ...string) error {
	layout, found := s.findFirstLayout(layouts...)
	if found == false {
		jww.WARN.Printf("Unable to locate layout for %s: %s\n", name, layouts)
		return nil
	}

	if err := s.renderThing(d, layout, renderBuffer); err != nil {
		// Behavior here should be dependent on if running in server or watch mode.
		distinctErrorLogger.Printf("Error while rendering %s: %v", name, err)
		if !s.Running() {
			os.Exit(-1)
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
	if s.Tmpl.Lookup(layout) == nil {
		return fmt.Errorf("Layout not found: %s", layout)
	}
	return s.Tmpl.ExecuteTemplate(w, layout, d)
}

func (s *Site) NewXMLBuffer() *bytes.Buffer {
	header := "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n"
	return bytes.NewBufferString(header)
}

func (s *Site) PageTarget() target.Output {
	s.initTargetList()
	return s.Targets.Page
}

func (s *Site) FileTarget() target.Output {
	s.initTargetList()
	return s.Targets.File
}

func (s *Site) AliasTarget() target.AliasPublisher {
	s.initTargetList()
	return s.Targets.Alias
}

func (s *Site) initTargetList() {
	s.targetListInit.Do(func() {
		if s.Targets.Page == nil {
			s.Targets.Page = &target.PagePub{
				PublishDir: s.absPublishDir(),
				UglyURLs:   viper.GetBool("UglyURLs"),
			}
		}
		if s.Targets.File == nil {
			s.Targets.File = &target.Filesystem{
				PublishDir: s.absPublishDir(),
			}
		}
		if s.Targets.Alias == nil {
			s.Targets.Alias = &target.HTMLRedirectAlias{
				PublishDir: s.absPublishDir(),
			}
		}
	})
}

func (s *Site) WriteDestFile(path string, reader io.Reader) (err error) {
	jww.DEBUG.Println("creating file:", path)
	return s.FileTarget().Publish(path, reader)
}

func (s *Site) WriteDestPage(path string, reader io.Reader) (err error) {
	jww.DEBUG.Println("creating page:", path)
	return s.PageTarget().Publish(path, reader)
}

func (s *Site) WriteDestAlias(path string, permalink template.HTML) (err error) {
	jww.DEBUG.Println("alias created at:", path)
	return s.AliasTarget().Publish(path, permalink)
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
		return "0 future content "
	case 1:
		msg = "1 future rendered "
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
