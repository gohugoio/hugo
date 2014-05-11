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
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/pkg/inflect"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
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
	Pages      Pages
	Tmpl       Template
	Taxonomies TaxonomyList
	Source     source.Input
	Sections   Taxonomy
	Info       SiteInfo
	Shortcodes map[string]ShortcodeFunc
	Menus      Menus
	timer      *nitro.B
	Target     target.Output
	Alias      target.AliasPublisher
	Completed  chan bool
	RunMode    runmode
	params     map[string]interface{}
}

type SiteInfo struct {
	BaseUrl         template.URL
	Taxonomies      TaxonomyList
	Indexes         *TaxonomyList // legacy, should be identical to Taxonomies
	Recent          *Pages
	Menus           *Menus
	Title           string
	Author          map[string]string
	LanguageCode    string
	DisqusShortname string
	Copyright       string
	LastChange      time.Time
	ConfigGet       func(key string) interface{}
	Permalinks      PermalinkOverrides
	Params          map[string]interface{}
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
	s.initTarget()
	s.Alias = &target.HTMLRedirectAlias{
		PublishDir: s.absPublishDir(),
	}
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
			jww.FEEDBACK.Println(p.FileName + " ")
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
	params, ok := viper.Get("Params").(map[string]interface{})
	if !ok {
		params = make(map[string]interface{})
	}

	permalinks := make(PermalinkOverrides)
	for k, v := range viper.GetStringMapString("Permalinks") {
		permalinks[k] = PathPattern(v)
	}

	s.Info = SiteInfo{
		BaseUrl:         template.URL(helpers.SanitizeUrl(viper.GetString("BaseUrl"))),
		Title:           viper.GetString("Title"),
		Author:          viper.GetStringMapString("author"),
		LanguageCode:    viper.GetString("languagecode"),
		Copyright:       viper.GetString("copyright"),
		DisqusShortname: viper.GetString("DisqusShortname"),
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
	if b, _ := helpers.DirExists(s.absContentDir()); !b {
		return fmt.Errorf("No source directory found, expecting to find it at " + s.absContentDir())
	}
	return
}

func (s *Site) CreatePages() (err error) {
	if s.Source == nil {
		panic(fmt.Sprintf("s.Source not set %s", s.absContentDir()))
	}
	if len(s.Source.Files()) < 1 {
		return fmt.Errorf("No source files found in %s", s.absContentDir())
	}

	var wg sync.WaitGroup
	for _, fi := range s.Source.Files() {
		wg.Add(1)
		go func(file *source.File) (err error) {
			defer wg.Done()

			page, err := NewPage(file.LogicalName)
			if err != nil {
				return err
			}
			err = page.ReadFrom(file.Contents)
			if err != nil {
				return err
			}
			page.Site = s.Info
			page.Tmpl = s.Tmpl
			page.Section = file.Section
			page.Dir = file.Dir

			//Handling short codes prior to Conversion to HTML
			page.ProcessShortcodes(s.Tmpl)

			err = page.Convert()
			if err != nil {
				return err
			}

			if viper.GetBool("BuildDrafts") || !page.Draft {
				s.Pages = append(s.Pages, page)
			}

			return
		}(fi)
	}

	wg.Wait()
	s.Pages.Sort()
	return
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
			flat[twoD{name, me.Name}] = me
		}
	}

	//creating flat hash
	for _, p := range s.Pages {
		for name, me := range p.Menus() {
			flat[twoD{name, me.Name}] = me
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
				v, ok := vals.([]string)
				if ok {
					for _, idx := range v {
						x := WeightedPage{weight.(int), p}

						s.Taxonomies[plural].Add(idx, x)
					}
				} else {
					jww.ERROR.Printf("Invalid %s in %s\n", plural, p.File.FileName)
				}
			}
		}
		for k := range s.Taxonomies[plural] {
			s.Taxonomies[plural][k].Sort()
		}
	}

	s.Info.Taxonomies = s.Taxonomies
	s.Info.Indexes = &s.Taxonomies
}

func (s *Site) assembleSections() {
	for i, p := range s.Pages {
		s.Sections.Add(p.Section, WeightedPage{s.Pages[i].Weight, s.Pages[i]})
	}

	for k := range s.Sections {
		s.Sections[k].Sort()
	}
}

func (s *Site) possibleTaxonomies() (taxonomies []string) {
	for _, p := range s.Pages {
		for k := range p.Params {
			if !inStringArray(taxonomies, k) {
				taxonomies = append(taxonomies, k)
			}
		}
	}
	return
}

func inStringArray(arr []string, el string) bool {
	for _, v := range arr {
		if v == el {
			return true
		}
	}
	return false
}

// Render shell pages that simply have a redirect in the header
func (s *Site) RenderAliases() error {
	for _, p := range s.Pages {
		for _, a := range p.Aliases {
			plink, err := p.Permalink()
			if err != nil {
				return err
			}
			if err := s.WriteAlias(a, template.HTML(plink)); err != nil {
				return err
			}
		}
	}
	return nil
}

// Render pages each corresponding to a markdown file
func (s *Site) RenderPages() (err error) {
	var wg sync.WaitGroup
	for _, page := range s.Pages {
		wg.Add(1)
		go func(p *Page) (err error) {
			var layouts []string
			defer wg.Done()

			if !p.IsRenderable() {
				self := "__" + p.TargetPath()
				_, err := s.Tmpl.New(self).Parse(string(p.Content))
				if err != nil {
					return err
				}
				layouts = append(layouts, self)
			} else {
				layouts = append(layouts, p.Layout()...)
				layouts = append(layouts, "_default/single.html")
			}

			return s.render(p, p.TargetPath(), s.appendThemeTemplates(layouts)...)
		}(page)
	}
	wg.Wait()

	if err != nil {
		return err
	}
	return nil
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

// Render the listing pages based on the meta data
// each unique term within a taxonomy will have a page created
func (s *Site) RenderTaxonomiesLists() (err error) {
	var wg sync.WaitGroup

	taxonomies := viper.GetStringMapString("Taxonomies")
	for sing, pl := range taxonomies {
		for key, oo := range s.Taxonomies[pl] {
			wg.Add(1)
			go func(k string, o WeightedPages, singular string, plural string) (err error) {
				defer wg.Done()
				base := plural + "/" + k
				n := s.NewNode()
				n.Title = strings.Title(k)
				s.setUrls(n, base)
				n.Date = o[0].Page.Date
				n.Data[singular] = o
				n.Data["Pages"] = o.Pages()
				layouts := []string{"taxonomy/" + singular + ".html", "indexes/" + singular + ".html", "_default/taxonomy.html", "_default/list.html"}
				err = s.render(n, base+".html", s.appendThemeTemplates(layouts)...)
				if err != nil {
					return err
				}

				if !viper.GetBool("DisableRSS") {
					// XML Feed
					s.setUrls(n, base+".xml")
					rssLayouts := []string{"taxonomy/" + singular + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}
					err := s.render(n, base+".xml", s.appendThemeTemplates(rssLayouts)...)
					if err != nil {
						return err
					}
				}
				return
			}(key, oo, sing, pl)
		}
	}
	wg.Wait()
	return nil
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
			err := s.render(n, plural+"/index.html", layouts...)
			if err != nil {
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
		n.Title = strings.Title(inflect.Pluralize(section))
		s.setUrls(n, section)
		n.Date = data[0].Page.Date
		n.Data["Pages"] = data.Pages()
		layouts := []string{"section/" + section + ".html", "_default/section.html", "_default/list.html", "indexes/" + section + ".html", "_default/indexes.html"}

		err := s.render(n, section, s.appendThemeTemplates(layouts)...)
		if err != nil {
			return err
		}

		if !viper.GetBool("DisableRSS") {
			// XML Feed
			rssLayouts := []string{"section/" + section + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"}
			s.setUrls(n, section+".xml")
			err = s.render(n, section+".xml", s.appendThemeTemplates(rssLayouts)...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Site) RenderHomePage() error {
	n := s.NewNode()
	n.Title = n.Site.Title
	s.setUrls(n, "/")
	n.Data["Pages"] = s.Pages
	layouts := []string{"index.html"}
	err := s.render(n, "/", s.appendThemeTemplates(layouts)...)
	if err != nil {
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

		if !viper.GetBool("DisableRSS") {
			rssLayouts := []string{"rss.xml", "_default/rss.xml", "_internal/_default/rss.xml"}
			err := s.render(n, ".xml", s.appendThemeTemplates(rssLayouts)...)
			if err != nil {
				return err
			}
		}
	}

	if a := s.Tmpl.Lookup("404.html"); a != nil {
		n.Url = helpers.Urlize("404.html")
		n.Title = "404 Page not found"
		n.Permalink = s.permalink("404.html")

		layouts := []string{"404.html"}
		return s.render(n, "404.html", s.appendThemeTemplates(layouts)...)
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
	page.Site = s.Info
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
	switch s.Target.(type) {
	case *target.Filesystem:
		s.Target.(*target.Filesystem).UglyUrls = true
		optChanged = true
	}

	smLayouts := []string{"sitemap.xml", "_default/sitemap.xml", "_internal/_default/sitemap.xml"}
	err := s.render(n, "sitemap.xml", s.appendThemeTemplates(smLayouts)...)
	if err != nil {
		return err
	}

	if optChanged {
		s.Target.(*target.Filesystem).UglyUrls = viper.GetBool("UglyUrls")
	}

	return nil
}

func (s *Site) Stats() {
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
	return template.HTML(helpers.MakePermalink(string(viper.GetString("BaseUrl")), s.prepUrl(plink)).String())
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
		Site: s.Info,
	}
}

func (s *Site) layoutExists(layouts ...string) bool {
	_, found := s.findFirstLayout(layouts...)

	return found
}

func (s *Site) render(d interface{}, out string, layouts ...string) (err error) {

	layout, found := s.findFirstLayout(layouts...)
	if found == false {
		jww.WARN.Printf("Unable to locate layout: %s\n", layouts)
		return
	}

	transformLinks := transform.NewEmptyTransforms()

	if viper.GetBool("CanonifyUrls") {
		absURL, err := transform.AbsURL(viper.GetString("BaseUrl"))
		if err != nil {
			return err
		}
		transformLinks = append(transformLinks, absURL...)
	}

	transformer := transform.NewChain(transformLinks...)

	var renderBuffer *bytes.Buffer

	if strings.HasSuffix(out, ".xml") {
		renderBuffer = s.NewXMLBuffer()
	} else {
		renderBuffer = new(bytes.Buffer)
	}

	err = s.renderThing(d, layout, renderBuffer)
	if err != nil {
		// Behavior here should be dependent on if running in server or watch mode.
		jww.ERROR.Println(fmt.Errorf("Rendering error: %v", err))
		if !s.Running() {
			os.Exit(-1)
		}
	}

	var outBuffer = new(bytes.Buffer)
	if strings.HasSuffix(out, ".xml") {
		outBuffer = renderBuffer
	} else {
		transformer.Apply(outBuffer, renderBuffer)
	}

	return s.WritePublic(out, outBuffer)
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

func (s *Site) initTarget() {
	if s.Target == nil {
		s.Target = &target.Filesystem{
			PublishDir: s.absPublishDir(),
			UglyUrls:   viper.GetBool("UglyUrls"),
		}
	}
}

func (s *Site) WritePublic(path string, reader io.Reader) (err error) {
	s.initTarget()

	jww.DEBUG.Println("writing to", path)
	return s.Target.Publish(path, reader)
}

func (s *Site) WriteAlias(path string, permalink template.HTML) (err error) {
	if s.Alias == nil {
		s.initTarget()
		s.Alias = &target.HTMLRedirectAlias{
			PublishDir: s.absPublishDir(),
		}
	}

	jww.DEBUG.Println("alias created at", path)

	return s.Alias.Publish(path, permalink)
}
