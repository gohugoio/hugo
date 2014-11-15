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
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"bitbucket.org/pkg/inflect"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/hugo/transform"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/nitro"
	"github.com/spf13/viper"
)

var _ = transform.AbsURL

var DefaultTimer *nitro.B

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
	Pages       Pages
	Files       []*source.File
	Tmpl        Template
	Taxonomies  TaxonomyList
	Source      source.Input
	Sections    Taxonomy
	Info        SiteInfo
	Shortcodes  map[string]ShortcodeFunc
	Menus       Menus
	timer       *nitro.B
	Targets     targetList
	Completed   chan bool
	RunMode     runmode
	params      map[string]interface{}
	draftCount  int
	futureCount int
}

type targetList struct {
	Page  target.Output
	File  target.Output
	Alias target.AliasPublisher
}

type SiteInfo struct {
	BaseUrl         template.URL
	Taxonomies      TaxonomyList
	Indexes         *TaxonomyList // legacy, should be identical to Taxonomies
	Sections        Taxonomy
	Pages           *Pages
	Files           []*source.File
	Recent          *Pages // legacy, should be identical to Pages
	Menus           *Menus
	Title           string
	Author          map[string]interface{}
	LanguageCode    string
	DisqusShortname string
	Copyright       string
	LastChange      time.Time
	Permalinks      PermalinkOverrides
	Params          map[string]interface{}
	BuildDrafts     bool
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
	s.Tmpl = NewTemplate()
	s.Tmpl.LoadTemplates(s.absLayoutDir())
	if s.hasTheme() {
		s.Tmpl.LoadTemplatesWithPrefix(s.absThemeDir()+"/layouts", "theme")
	}
}

func (s *Site) addTemplate(name, data string) error {
	return s.Tmpl.AddTemplate(name, data)
}

func (s *Site) Process() (err error) {
	if err = s.initialize(); err != nil {
		return
	}
	s.prepTemplates()
	s.timerStep("initialize & template prep")
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

	s.Shortcodes = make(map[string]ShortcodeFunc)
	return
}

func (s *Site) initializeSiteInfo() {
	params := viper.GetStringMap("Params")

	permalinks := make(PermalinkOverrides)
	for k, v := range viper.GetStringMapString("Permalinks") {
		permalinks[k] = PathPattern(v)
	}

	s.Info = SiteInfo{
		BaseUrl:         template.URL(helpers.SanitizeUrl(viper.GetString("BaseUrl"))),
		Title:           viper.GetString("Title"),
		Author:          viper.GetStringMap("author"),
		LanguageCode:    viper.GetString("languagecode"),
		Copyright:       viper.GetString("copyright"),
		DisqusShortname: viper.GetString("DisqusShortname"),
		BuildDrafts:     viper.GetBool("BuildDrafts"),
		Pages:           &s.Pages,
		Recent:          &s.Pages,
		Menus:           &s.Menus,
		Params:          params,
		Permalinks:      permalinks,
	}
}

func (s *Site) hasTheme() bool {
	return viper.GetString("theme") != ""
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

	for i := 0; i < procs*4; i++ {
		wg.Add(1)
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

	for i := 0; i < procs*4; i++ {
		wg.Add(1)
		go fileConverter(s, fileConvChan, results, wg)
	}

	wg = &sync.WaitGroup{}
	for i := 0; i < procs*4; i++ {
		wg.Add(1)
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
		h := FindHandler(file.Extension())
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
		var h Handler
		if page.Markup != "" {
			h = FindHandler(page.Markup)
		} else {
			h = FindHandler(page.File.Extension())
		}
		if h != nil {
			h.Convert(page, s, results)
		}
	}
}

func fileConverter(s *Site, files <-chan *source.File, results HandleResults, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		h := FindHandler(file.Extension())
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
			errMsgs = append(errMsgs, r.err.Error())
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
				s.draftCount += 1
			}

			if r.page.IsFuture() {
				s.futureCount += 1
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
	s.Info.LastChange = s.Pages[0].Date

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
					if strings.HasPrefix(menuEntry.Url, "/") {
						// make it absolute so it matches the nodes
						menuEntry.Url = s.permalinkStr(menuEntry.Url)
					}
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

	//creating flat hash
	for _, p := range s.Pages {
		for name, me := range p.Menus() {
			if _, ok := flat[twoD{name, me.KeyName()}]; ok {
				jww.ERROR.Printf("Two or more menu items have the same name/identifier in %q Menu. Identified as %q.\n Rename or set a unique identifier. \n", name, me.KeyName())
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
			// if parent does not exist, create one without a url
			flat[twoD{p.MenuName, p.EntryName}] = &MenuEntry{Name: p.EntryName, Url: ""}
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
	s.Info.Indexes = &s.Taxonomies
	s.Info.Sections = s.Sections
}

func (s *Site) assembleSections() {
	for i, p := range s.Pages {
		s.Sections.Add(p.Section(), WeightedPage{s.Pages[i].Weight, s.Pages[i]})
	}

	for k := range s.Sections {
		s.Sections[k].Sort()
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

// Render shell pages that simply have a redirect in the header
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

// Render pages each corresponding to a markdown file
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

		b, err := s.renderPage("page "+p.FullFilePath(), p, s.appendThemeTemplates(layouts)...)
		if err != nil {
			results <- err
		} else {
			results <- s.WriteDestPage(p.TargetPath(), b)
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
			if !strings.HasPrefix("_internal/", t) {
				out = append(out, t)
			}
		}

		// Then place theme templates with the same names
		for _, t := range in {
			if !strings.HasPrefix("_internal/", t) {
				out = append(out, "theme/"+t)
			}
		}
		// Lastly place internal templates
		for _, t := range in {
			if strings.HasPrefix("_internal/", t) {
				out = append(out, "theme/"+t)
			}
		}
		return out
	} else {
		return in
	}
}

type taxRenderInfo struct {
	key      string
	pages    WeightedPages
	singular string
	plural   string
}

// Render the listing pages based on the meta data
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
	s.setUrls(n, base)
	if len(t.pages) > 0 {
		n.Date = t.pages[0].Page.Date
	}
	n.Data[t.singular] = t.pages
	n.Data["Singular"] = t.singular
	n.Data["Plural"] = t.plural
	n.Data["Pages"] = t.pages.Pages()
	return n, base
}

func taxonomyRenderer(s *Site, taxes <-chan taxRenderInfo, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range taxes {
		n, base := s.newTaxonomyNode(t)
		layouts := []string{"taxonomy/" + t.singular + ".html", "indexes/" + t.singular + ".html", "_default/taxonomy.html", "_default/list.html"}
		b, err := s.renderPage("taxononomy "+t.singular, n, s.appendThemeTemplates(layouts)...)
		if err != nil {
			results <- err
			continue
		} else {
			err := s.WriteDestPage(base+".html", b)
			if err != nil {
				results <- err
			}
		}

		if !viper.GetBool("DisableRSS") {
			// XML Feed
			s.setUrls(n, base+".xml")
			rssLayouts := []string{"taxonomy/" + t.singular + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}
			b, err := s.renderXML("taxonomy "+t.singular+" rss", n, s.appendThemeTemplates(rssLayouts)...)
			if err != nil {
				results <- err
				continue
			} else {
				err := s.WriteDestFile(base+".xml", b)
				if err != nil {
					results <- err
				}
			}
		}
	}
}

// Render a page per taxonomy that lists the terms for that taxonomy
func (s *Site) RenderListsOfTaxonomyTerms() (err error) {
	taxonomies := viper.GetStringMapString("Taxonomies")
	for singular, plural := range taxonomies {
		n := s.NewNode()
		n.Title = strings.Title(plural)
		s.setUrls(n, plural)
		n.Data["Singular"] = singular
		n.Data["Plural"] = plural
		n.Data["Terms"] = s.Taxonomies[plural]
		// keep the following just for legacy reasons
		n.Data["OrderedIndex"] = n.Data["Terms"]
		n.Data["Index"] = n.Data["Terms"]
		layouts := []string{"taxonomy/" + singular + ".terms.html", "_default/terms.html", "indexes/indexes.html"}
		layouts = s.appendThemeTemplates(layouts)
		if s.layoutExists(layouts...) {
			b, err := s.renderPage("taxonomy terms for "+singular, n, layouts...)
			if err != nil {
				return err
			}
			if err := s.WriteDestPage(plural+"/index.html", b); err != nil {
				return err
			}
		}
	}

	return
}

// Render a page for each section
func (s *Site) RenderSectionLists() error {
	for section, data := range s.Sections {
		n := s.NewNode()
		if viper.GetBool("PluralizeListTitles") {
			n.Title = strings.Title(inflect.Pluralize(section))
		} else {
			n.Title = strings.Title(section)
		}
		s.setUrls(n, section)
		n.Date = data[0].Page.Date
		n.Data["Pages"] = data.Pages()
		layouts := []string{"section/" + section + ".html", "_default/section.html", "_default/list.html", "indexes/" + section + ".html", "_default/indexes.html"}

		b, err := s.renderPage("section "+section, n, s.appendThemeTemplates(layouts)...)
		if err != nil {
			return err
		}
		if err := s.WriteDestPage(section, b); err != nil {
			return err
		}

		if !viper.GetBool("DisableRSS") {
			// XML Feed
			rssLayouts := []string{"section/" + section + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}
			s.setUrls(n, section+".xml")
			b, err = s.renderXML("section "+section+" rss", n, s.appendThemeTemplates(rssLayouts)...)
			if err != nil {
				return err
			}
			if err := s.WriteDestFile(section+".xml", b); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Site) newHomeNode() *Node {
	n := s.NewNode()
	n.Title = n.Site.Title
	s.setUrls(n, "/")
	n.Data["Pages"] = s.Pages
	return n
}

func (s *Site) RenderHomePage() error {
	n := s.newHomeNode()
	layouts := []string{"index.html", "_default/list.html", "_default/single.html"}
	b, err := s.renderPage("homepage", n, s.appendThemeTemplates(layouts)...)
	if err != nil {
		return err
	}
	if err := s.WriteDestPage("/", b); err != nil {
		return err
	}

	if !viper.GetBool("DisableRSS") {
		// XML Feed
		n.Url = helpers.Urlize("index.xml")
		n.Title = "Recent Content"
		n.Permalink = s.permalink("index.xml")
		high := 50
		if len(s.Pages) < high {
			high = len(s.Pages)
		}
		n.Data["Pages"] = s.Pages[:high]
		if len(s.Pages) > 0 {
			n.Date = s.Pages[0].Date
		}

		rssLayouts := []string{"rss.xml", "_default/rss.xml", "_internal/_default/rss.xml"}
		b, err := s.renderXML("homepage rss", n, s.appendThemeTemplates(rssLayouts)...)
		if err != nil {
			return err
		}
		if err := s.WriteDestFile("index.xml", b); err != nil {
			return err
		}
	}

	n.Url = helpers.Urlize("404.html")
	n.Title = "404 Page not found"
	n.Permalink = s.permalink("404.html")

	nfLayouts := []string{"404.html"}
	b, nfErr := s.renderPage("404 page", n, s.appendThemeTemplates(nfLayouts)...)
	if nfErr != nil {
		return nfErr
	}
	if err := s.WriteDestFile("404.html", b); err != nil {
		return err
	}

	return nil
}

func (s *Site) RenderSitemap() error {
	if viper.GetBool("DisableSitemap") {
		return nil
	}

	sitemapDefault := parseSitemap(viper.GetStringMap("Sitemap"))

	optChanged := false

	n := s.NewNode()

	// Prepend homepage to the list of pages
	pages := make(Pages, 0)

	page := &Page{}
	page.Date = s.Info.LastChange
	page.Site = &s.Info
	page.Url = "/"

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

	// Force `UglyUrls` option to force `sitemap.xml` file name
	switch s.PageTarget().(type) {
	case *target.Filesystem:
		s.PageTarget().(*target.PagePub).UglyUrls = true
		optChanged = true
	}

	smLayouts := []string{"sitemap.xml", "_default/sitemap.xml", "_internal/_default/sitemap.xml"}
	b, err := s.renderXML("sitemap", n, s.appendThemeTemplates(smLayouts)...)
	if err != nil {
		return err
	}
	if err := s.WriteDestFile("sitemap.xml", b); err != nil {
		return err
	}

	if optChanged {
		s.PageTarget().(*target.PagePub).UglyUrls = viper.GetBool("UglyUrls")
	}

	return nil
}

func (s *Site) Stats() {
	jww.FEEDBACK.Println(s.draftStats())
	jww.FEEDBACK.Println(s.futureStats())
	jww.FEEDBACK.Printf("%d pages created \n", len(s.Pages))

	taxonomies := viper.GetStringMapString("Taxonomies")

	for _, pl := range taxonomies {
		jww.FEEDBACK.Printf("%d %s created\n", len(s.Taxonomies[pl]), pl)
	}
}

func (s *Site) setUrls(n *Node, in string) {
	n.Url = s.prepUrl(in)
	n.Permalink = s.permalink(n.Url)
	n.RSSLink = s.permalink(in + ".xml")
}

func (s *Site) permalink(plink string) template.HTML {
	return template.HTML(s.permalinkStr(plink))
}

func (s *Site) permalinkStr(plink string) string {
	return helpers.MakePermalink(string(viper.GetString("BaseUrl")), s.prepUrl(plink)).String()
}

func (s *Site) prepUrl(in string) string {
	return helpers.Urlize(s.PrettifyUrl(in))
}

func (s *Site) PrettifyUrl(in string) string {
	return helpers.UrlPrep(viper.GetBool("UglyUrls"), in)
}

func (s *Site) PrettifyPath(in string) string {
	return helpers.PathPrep(viper.GetBool("UglyUrls"), in)
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

func (s *Site) renderXML(name string, d interface{}, layouts ...string) (io.Reader, error) {
	renderBuffer := s.NewXMLBuffer()
	err := s.render(name, d, renderBuffer, layouts...)
	return renderBuffer, err
}

func (s *Site) renderPage(name string, d interface{}, layouts ...string) (io.Reader, error) {
	renderBuffer := new(bytes.Buffer)
	err := s.render(name, d, renderBuffer, layouts...)

	var outBuffer = new(bytes.Buffer)

	transformLinks := transform.NewEmptyTransforms()

	if viper.GetBool("CanonifyUrls") {
		absURL, err := transform.AbsURL(viper.GetString("BaseUrl"))
		if err != nil {
			return nil, err
		}
		transformLinks = append(transformLinks, absURL...)
	}

	if viper.GetBool("watch") && !viper.GetBool("DisableLiveReload") {
		transformLinks = append(transformLinks, transform.LiveReloadInject)
	}

	transformer := transform.NewChain(transformLinks...)
	transformer.Apply(outBuffer, renderBuffer)
	return outBuffer, err
}

func (s *Site) render(name string, d interface{}, renderBuffer *bytes.Buffer, layouts ...string) error {
	layout, found := s.findFirstLayout(layouts...)
	if found == false {
		jww.WARN.Printf("Unable to locate layout for %s: %s\n", name, layouts)
		return nil
	}

	if err := s.renderThing(d, layout, renderBuffer); err != nil {
		// Behavior here should be dependent on if running in server or watch mode.
		jww.ERROR.Println(fmt.Errorf("Error while rendering %s: %v", name, err))
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
	if s.Targets.Page == nil {
		s.Targets.Page = &target.PagePub{
			PublishDir: s.absPublishDir(),
			UglyUrls:   viper.GetBool("UglyUrls"),
		}
	}
	return s.Targets.Page
}

func (s *Site) FileTarget() target.Output {
	if s.Targets.File == nil {
		s.Targets.File = &target.Filesystem{
			PublishDir: s.absPublishDir(),
		}
	}
	return s.Targets.File
}

func (s *Site) AliasTarget() target.AliasPublisher {
	if s.Targets.Alias == nil {
		s.Targets.Alias = &target.HTMLRedirectAlias{
			PublishDir: s.absPublishDir(),
		}

	}
	return s.Targets.Alias
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
		return "0 draft content "
	case 1:
		msg = "1 draft rendered "
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
