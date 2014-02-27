// Copyright © 2013 Steve Francia <spf@spf13.com>.
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
	"bitbucket.org/pkg/inflect"
	"bytes"
	"fmt"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/hugo/template/bundle"
	"github.com/spf13/hugo/transform"
	"github.com/spf13/nitro"
	"html/template"
	"io"
	"os"
	"strings"
	"time"
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
// 3. Indexes are created via configuration and will present some aspect of
//    the final page and typically a perm url.
//
// 4. All Pages are passed through a template based on their desired
//    layout based on numerous different elements.
//
// 5. The entire collection of files is written to disk.
type Site struct {
	Config     Config
	Pages      Pages
	Tmpl       bundle.Template
	Indexes    IndexList
	Source     source.Input
	Sections   Index
	Info       SiteInfo
	Shortcodes map[string]ShortcodeFunc
	timer      *nitro.B
	Target     target.Output
	Alias      target.AliasPublisher
	Completed  chan bool
	RunMode    runmode
	params     map[string]interface{}
}

type SiteInfo struct {
	BaseUrl    template.URL
	Indexes    IndexList
	Recent     *Pages
	LastChange time.Time
	Title      string
	Config     *Config
	Permalinks PermalinkOverrides
	Params     map[string]interface{}
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
		fmt.Printf("Error rendering site: %s\nAvailable templates:\n", err)
		for _, template := range s.Tmpl.Templates() {
			fmt.Printf("\t%s\n", template.Name())
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
	s.Tmpl = bundle.NewTemplate()
	s.Tmpl.LoadTemplates(s.absLayoutDir())
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
	s.timerStep("build indexes")
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
	if err = s.RenderIndexes(); err != nil {
		return
	}
	s.timerStep("render and write indexes")
	s.RenderIndexesIndexes()
	s.timerStep("render & write index indexes")
	if err = s.RenderLists(); err != nil {
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
	return
}

func (s *Site) checkDescriptions() {
	for _, p := range s.Pages {
		if len(p.Description) < 60 {
			fmt.Println(p.FileName + " ")
		}
	}
}

func (s *Site) initialize() (err error) {
	if err = s.checkDirectories(); err != nil {
		return err
	}

	staticDir := s.Config.GetAbsPath(s.Config.StaticDir + "/")

	s.Source = &source.Filesystem{
		AvoidPaths: []string{staticDir},
		Base:       s.absContentDir(),
	}

	s.initializeSiteInfo()

	s.Shortcodes = make(map[string]ShortcodeFunc)
	return
}

func (s *Site) initializeSiteInfo() {
	s.Info = SiteInfo{
		BaseUrl:    template.URL(s.Config.BaseUrl),
		Title:      s.Config.Title,
		Recent:     &s.Pages,
		Config:     &s.Config,
		Params:     s.Config.Params,
		Permalinks: s.Config.Permalinks,
	}
}

func (s *Site) absLayoutDir() string {
	return s.Config.GetAbsPath(s.Config.LayoutDir)
}

func (s *Site) absContentDir() string {
	return s.Config.GetAbsPath(s.Config.ContentDir)
}

func (s *Site) absPublishDir() string {
	return s.Config.GetAbsPath(s.Config.PublishDir)
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
	for _, file := range s.Source.Files() {
		page, err := ReadFrom(file.Contents, file.LogicalName)
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

		if s.Config.BuildDrafts || !page.Draft {
			s.Pages = append(s.Pages, page)
		}
	}

	s.Pages.Sort()
	return
}

func (s *Site) BuildSiteMeta() (err error) {
	s.Indexes = make(IndexList)
	s.Sections = make(Index)

	for _, plural := range s.Config.Indexes {
		s.Indexes[plural] = make(Index)
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

						s.Indexes[plural].Add(idx, x)
					}
				} else {
					if s.Config.Verbose {
						fmt.Fprintf(os.Stderr, "Invalid %s in %s\n", plural, p.File.FileName)
					}
				}
			}
		}
		for k := range s.Indexes[plural] {
			s.Indexes[plural][k].Sort()
		}
	}

	for i, p := range s.Pages {
		s.Sections.Add(p.Section, WeightedPage{s.Pages[i].Weight, s.Pages[i]})
	}

	for k := range s.Sections {
		s.Sections[k].Sort()
	}

	s.Info.Indexes = s.Indexes

	if len(s.Pages) == 0 {
		return
	}
	s.Info.LastChange = s.Pages[0].Date

	// populate pages with site metadata
	for _, p := range s.Pages {
		p.Site = s.Info
	}

	return
}

func (s *Site) possibleIndexes() (indexes []string) {
	for _, p := range s.Pages {
		for k := range p.Params {
			if !inStringArray(indexes, k) {
				indexes = append(indexes, k)
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

func (s *Site) RenderPages() (err error) {
	for _, p := range s.Pages {
		var layout []string

		if !p.IsRenderable() {
			self := "__" + p.TargetPath()
			_, err := s.Tmpl.New(self).Parse(string(p.Content))
			if err != nil {
				return err
			}
			layout = append(layout, self)
		} else {
			layout = append(layout, p.Layout()...)
			layout = append(layout, "_default/single.html")
		}

		err := s.render(p, p.TargetPath(), layout...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Site) RenderIndexes() error {
	for singular, plural := range s.Config.Indexes {
		for k, o := range s.Indexes[plural] {
			n := s.NewNode()
			n.Title = strings.Title(k)
			base := plural + "/" + k
			s.setUrls(n, base)
			n.Date = o[0].Page.Date
			n.Data[singular] = o
			n.Data["Pages"] = o.Pages()
			layout := "indexes/" + singular + ".html"
			err := s.render(n, base+".html", layout)
			if err != nil {
				return err
			}

			if a := s.Tmpl.Lookup("rss.xml"); a != nil {
				// XML Feed
				s.setUrls(n, base+".xml")
				err := s.render(n, base+".xml", "rss.xml")
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Site) RenderIndexesIndexes() (err error) {
	layout := "indexes/indexes.html"
	if s.Tmpl.Lookup(layout) != nil {
		for singular, plural := range s.Config.Indexes {
			n := s.NewNode()
			n.Title = strings.Title(plural)
			s.setUrls(n, plural)
			n.Data["Singular"] = singular
			n.Data["Plural"] = plural
			n.Data["Index"] = s.Indexes[plural]
			// keep the following just for legacy reasons
			n.Data["OrderedIndex"] = s.Indexes[plural]

			err := s.render(n, plural+"/index.html", layout)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (s *Site) RenderLists() error {
	for section, data := range s.Sections {
		n := s.NewNode()
		n.Title = strings.Title(inflect.Pluralize(section))
		s.setUrls(n, section)
		n.Date = data[0].Page.Date
		n.Data["Pages"] = data.Pages()
		layout := "indexes/" + section + ".html"

		err := s.render(n, section, layout, "_default/indexes.html")
		if err != nil {
			return err
		}

		if a := s.Tmpl.Lookup("rss.xml"); a != nil {
			// XML Feed
			s.setUrls(n, section+".xml")
			err = s.render(n, section+".xml", "rss.xml")
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
	err := s.render(n, "/", "index.html")
	if err != nil {
		return err
	}

	if a := s.Tmpl.Lookup("rss.xml"); a != nil {
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
		err := s.render(n, ".xml", "rss.xml")
		if err != nil {
			return err
		}
	}

	if a := s.Tmpl.Lookup("404.html"); a != nil {
		n.Url = helpers.Urlize("404.html")
		n.Title = "404 Page not found"
		n.Permalink = s.permalink("404.html")
		return s.render(n, "404.html", "404.html")
	}

	return nil
}

func (s *Site) Stats() {
	fmt.Printf("%d pages created \n", len(s.Pages))
	for _, pl := range s.Config.Indexes {
		fmt.Printf("%d %s index created\n", len(s.Indexes[pl]), pl)
	}
}

func (s *Site) setUrls(n *Node, in string) {
	n.Url = s.prepUrl(in)
	n.Permalink = s.permalink(n.Url)
	n.RSSLink = s.permalink(in + ".xml")
}

func (s *Site) permalink(plink string) template.HTML {
	return template.HTML(helpers.MakePermalink(string(s.Config.BaseUrl), s.prepUrl(plink)).String())
}

func (s *Site) prepUrl(in string) string {
	return helpers.Urlize(s.PrettifyUrl(in))
}

func (s *Site) PrettifyUrl(in string) string {
	return helpers.UrlPrep(s.Config.UglyUrls, in)
}

func (s *Site) PrettifyPath(in string) string {
	return helpers.PathPrep(s.Config.UglyUrls, in)
}

func (s *Site) NewNode() *Node {
	return &Node{
		Data: make(map[string]interface{}),
		Site: s.Info,
	}
}

func (s *Site) render(d interface{}, out string, layouts ...string) (err error) {

	layout := s.findFirstLayout(layouts...)
	if layout == "" {
		if s.Config.Verbose {
			fmt.Printf("Unable to locate layout: %s\n", layouts)
		}
		return
	}

	transformLinks := transform.NewEmptyTransforms()

	if s.Config.CanonifyUrls {
		absURL, err := transform.AbsURL(s.Config.BaseUrl)
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
		fmt.Println(fmt.Errorf("Rendering error: %v", err))
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

func (s *Site) findFirstLayout(layouts ...string) (layout string) {
	for _, layout = range layouts {
		if s.Tmpl.Lookup(layout) != nil {
			return
		}
	}
	return ""
}

func (s *Site) renderThing(d interface{}, layout string, w io.Writer) error {
	// If the template doesn't exist, then return, but leave the Writer open
	if s.Tmpl.Lookup(layout) == nil {
		return fmt.Errorf("Layout not found: %s", layout)
	}
	//defer w.Close()
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
			UglyUrls:   s.Config.UglyUrls,
		}
	}
}

func (s *Site) WritePublic(path string, reader io.Reader) (err error) {
	s.initTarget()

	if s.Config.Verbose {
		fmt.Println(path)
	}
	return s.Target.Publish(path, reader)
}

func (s *Site) WriteAlias(path string, permalink template.HTML) (err error) {
	if s.Alias == nil {
		s.initTarget()
		s.Alias = &target.HTMLRedirectAlias{
			PublishDir: s.absPublishDir(),
		}
	}

	if s.Config.Verbose {
		fmt.Println(path)
	}

	return s.Alias.Publish(path, permalink)
}
