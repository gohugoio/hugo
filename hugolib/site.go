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
	"errors"
	"fmt"
	"github.com/spf13/hugo/target"
	"github.com/spf13/nitro"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var DefaultTimer = nitro.Initalize()

type Site struct {
	Config     Config
	Pages      Pages
	Tmpl       Template
	Indexes    IndexList
	Files      []string
	Sections   Index
	Info       SiteInfo
	Shortcodes map[string]ShortcodeFunc
	timer      *nitro.B
	Target     target.Publisher
}

type SiteInfo struct {
	BaseUrl    URL
	Indexes    OrderedIndexList
	Recent     *Pages
	LastChange time.Time
	Title      string
	Config     *Config
}

func (s *Site) getFromIndex(kind string, name string) Pages {
	return s.Indexes[kind][name]
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
		fmt.Printf("Error rendering site: %s\n", err)
		fmt.Printf("Available templates:")
		for _, tpl := range s.Tmpl.Templates() {
			fmt.Printf("\t%s\n", tpl.Name())
		}
		return
	}
	s.Write()
	return nil
}

func (s *Site) Analyze() {
	s.Process()
	s.checkDescriptions()
}

func (s *Site) prepTemplates() {
	s.Tmpl = NewTemplate()
	s.Tmpl.LoadTemplates(s.absLayoutDir())
}

func (s *Site) addTemplate(name, data string) error {
	return s.Tmpl.AddTemplate(name, data)
}

func (s *Site) Process() (err error) {
	s.initialize()
	s.prepTemplates()
	s.timerStep("initialize & template prep")
	s.CreatePages()
	s.setupPrevNext()
	s.timerStep("import pages")
	if err = s.BuildSiteMeta(); err != nil {
		return
	}
	s.timerStep("build indexes")
	return
}

func (s *Site) Render() (err error) {
	s.RenderAliases()
	s.timerStep("render and write aliases")
	s.ProcessShortcodes()
	s.timerStep("render shortcodes")
	s.AbsUrlify()
	s.timerStep("absolute URLify")
	if err = s.RenderIndexes(); err != nil {
		return
	}
	s.RenderIndexesIndexes()
	s.timerStep("render and write indexes")
	s.RenderLists()
	s.timerStep("render and write lists")
	if err = s.RenderPages(); err != nil {
		return
	}
	s.timerStep("render pages")
	if err = s.RenderHomePage(); err != nil {
		return
	}
	s.timerStep("render and write homepage")
	return
}

func (s *Site) Write() {
	s.WritePages()
	s.timerStep("write pages")
}

func (s *Site) checkDescriptions() {
	for _, p := range s.Pages {
		if len(p.Description) < 60 {
			fmt.Print(p.FileName + " ")
		}
	}
}

func (s *Site) initialize() {
	s.checkDirectories()

	staticDir := s.Config.GetAbsPath(s.Config.StaticDir + "/")

	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			PrintErr("Walker: ", err)
			return nil
		}

		if fi.IsDir() {
			if path == staticDir {
				return filepath.SkipDir
			}
			return nil
		} else {
			if ignoreDotFile(path) {
				return nil
			}
			s.Files = append(s.Files, path)
			return nil
		}
	}

	filepath.Walk(s.absContentDir(), walker)
	s.Info = SiteInfo{
		BaseUrl: URL(s.Config.BaseUrl),
		Title:   s.Config.Title,
		Recent:  &s.Pages,
		Config:  &s.Config,
	}

	s.Shortcodes = make(map[string]ShortcodeFunc)
}

func ignoreDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
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

func (s *Site) checkDirectories() {
	if b, _ := dirExists(s.absLayoutDir()); !b {
		FatalErr("No layout directory found, expecting to find it at " + s.absLayoutDir())
	}
	if b, _ := dirExists(s.absContentDir()); !b {
		FatalErr("No source directory found, expecting to find it at " + s.absContentDir())
	}
	mkdirIf(s.absPublishDir())
}

func (s *Site) ProcessShortcodes() {
	for _, page := range s.Pages {
		page.Content = HTML(ShortcodesHandle(string(page.Content), page, s.Tmpl))
	}
}

func (s *Site) AbsUrlify() {
	baseWithoutTrailingSlash := strings.TrimRight(s.Config.BaseUrl, "/")
	baseWithSlash := baseWithoutTrailingSlash + "/"
	for _, page := range s.Pages {
		content := string(page.Content)
		content = strings.Replace(content, " src=\"/", " src=\""+baseWithSlash, -1)
		content = strings.Replace(content, " src='/", " src='"+baseWithSlash, -1)
		content = strings.Replace(content, " href='/", " href='"+baseWithSlash, -1)
		content = strings.Replace(content, " href=\"/", " href=\""+baseWithSlash, -1)
		content = strings.Replace(content, baseWithoutTrailingSlash+"//", baseWithSlash, -1)
		page.Content = HTML(content)
	}
}

func (s *Site) CreatePages() {
	for _, fileName := range s.Files {
		page := NewPage(fileName)
		page.Site = s.Info
		page.Tmpl = s.Tmpl
		_ = s.setUrlPath(page)
		page.Initalize()
		s.setOutFile(page)
		if s.Config.BuildDrafts || !page.Draft {
			s.Pages = append(s.Pages, page)
		}
	}

	s.Pages.Sort()
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

func (s *Site) setUrlPath(p *Page) error {
	y := strings.TrimPrefix(p.FileName, s.Config.GetAbsPath(s.Config.ContentDir))
	x := strings.Split(y, string(os.PathSeparator))

	if len(x) <= 1 {
		return errors.New("Zero length page name")
	}

	p.Section = strings.Trim(x[1], "/\\")
	p.Path = strings.Trim(strings.Join(x[:len(x)-1], string(os.PathSeparator)), "/\\")
	return nil
}

// If Url is provided it is assumed to be the complete relative path
// and will override everything
// Otherwise path + slug is used if provided
// Lastly path + filename is used if provided
func (s *Site) setOutFile(p *Page) {
	// Always use Url if it's specified
	if len(strings.TrimSpace(p.Url)) > 2 {
		p.OutFile = strings.TrimSpace(p.Url)

		if strings.HasSuffix(p.OutFile, "/") {
			p.OutFile = p.OutFile + "index.html"
		}
		return
	}

	var outfile string
	if len(strings.TrimSpace(p.Slug)) > 0 {
		// Use Slug if provided
		if s.Config.UglyUrls {
			outfile = strings.TrimSpace(p.Slug) + "." + p.Extension
		} else {
			outfile = filepath.Join(strings.TrimSpace(p.Slug), "index."+p.Extension)
		}
	} else {
		// Fall back to filename
		_, t := filepath.Split(p.FileName)
		if s.Config.UglyUrls {
			outfile = replaceExtension(strings.TrimSpace(t), p.Extension)
		} else {
			file, _ := fileExt(strings.TrimSpace(t))
			outfile = filepath.Join(file, "index."+p.Extension)
		}
	}

	p.OutFile = p.Path + string(os.PathSeparator) + strings.TrimSpace(outfile)
}

func (s *Site) BuildSiteMeta() (err error) {
	s.Indexes = make(IndexList)
	s.Sections = make(Index)

	for _, plural := range s.Config.Indexes {
		s.Indexes[plural] = make(Index)
		for _, p := range s.Pages {
			vals := p.GetParam(plural)

			if vals != nil {
				v, ok := vals.([]string)
				if ok {
					for _, idx := range v {
						s.Indexes[plural].Add(idx, p)
					}
				} else {
					PrintErr("Invalid " + plural + " in " + p.File.FileName)
				}
			}
		}
		for k, _ := range s.Indexes[plural] {
			s.Indexes[plural][k].Sort()
		}
	}

	for _, p := range s.Pages {
		s.Sections.Add(p.Section, p)
	}

	for k, _ := range s.Sections {
		s.Sections[k].Sort()
	}

	s.Info.Indexes = s.Indexes.BuildOrderedIndexList()

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
		for k, _ := range p.Params {
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
			t := "alias"
			if strings.HasSuffix(a, ".xhtml") {
				t = "alias-xhtml"
			}
			content, err := s.RenderThing(p, t)
			if strings.HasSuffix(a, "/") {
				a = a + "index.html"
			}
			if err != nil {
				return err
			}
			err = s.WritePublic(a, content.Bytes())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Site) RenderPages() error {
	for _, p := range s.Pages {
		content, err := s.RenderThingOrDefault(p, p.Layout(), "_default/single.html")
		if err != nil {
			return err
		}
		p.RenderedContent = content
	}
	return nil
}

func (s *Site) WritePages() (err error) {
	for _, p := range s.Pages {
		err = s.WritePublic(p.OutFile, p.RenderedContent.Bytes())
		if err != nil {
			return
		}
	}
	return
}

func (s *Site) RenderIndexes() error {
	for singular, plural := range s.Config.Indexes {
		for k, o := range s.Indexes[plural] {
			n := s.NewNode()
			n.Title = strings.Title(k)
			url := Urlize(plural + "/" + k)
			plink := url
			if s.Config.UglyUrls {
				n.Url = url + ".html"
				plink = n.Url
			} else {
				n.Url = url + "/index.html"
			}
			n.Permalink = permalink(s, plink)
			n.RSSlink = permalink(s, url+".xml")
			n.Date = o[0].Date
			n.Data[singular] = o
			n.Data["Pages"] = o
			layout := "indexes/" + singular + ".html"
			x, err := s.RenderThing(n, layout)
			if err != nil {
				return err
			}

			var base string
			if s.Config.UglyUrls {
				base = plural + "/" + k
			} else {
				base = plural + "/" + k + "/" + "index"
			}

			err = s.WritePublic(base+".html", x.Bytes())
			if err != nil {
				return err
			}

			if a := s.Tmpl.Lookup("rss.xml"); a != nil {
				// XML Feed
				y := s.NewXMLBuffer()
				if s.Config.UglyUrls {
					n.Url = Urlize(plural + "/" + k + ".xml")
				} else {
					n.Url = Urlize(plural + "/" + k + "/" + "index.xml")
				}
				n.Permalink = permalink(s, n.Url)
				s.Tmpl.ExecuteTemplate(y, "rss.xml", n)
				err = s.WritePublic(base+".xml", y.Bytes())
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
			url := Urlize(plural)
			n.Url = url + "/index.html"
			n.Permalink = permalink(s, n.Url)
			n.Data["Singular"] = singular
			n.Data["Plural"] = plural
			n.Data["Index"] = s.Indexes[plural]
			n.Data["OrderedIndex"] = s.Info.Indexes[plural]

			x, err := s.RenderThing(n, layout)
			if err != nil {
				return err
			}

			err = s.WritePublic(plural+"/index.html", x.Bytes())
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
		n.Url = Urlize(section + "/" + "index.html")
		n.Permalink = permalink(s, n.Url)
		n.RSSlink = permalink(s, section+".xml")
		n.Date = data[0].Date
		n.Data["Pages"] = data
		layout := "indexes/" + section + ".html"

		content, err := s.RenderThingOrDefault(n, layout, "_default/index.html")
		if err != nil {
			return err
		}
		err = s.WritePublic(section+"/index.html", content.Bytes())
		if err != nil {
			return err
		}

		if a := s.Tmpl.Lookup("rss.xml"); a != nil {
			// XML Feed
			if s.Config.UglyUrls {
				n.Url = Urlize(section + ".xml")
			} else {
				n.Url = Urlize(section + "/" + "index.xml")
			}
			n.Permalink = HTML(string(n.Site.BaseUrl) + n.Url)
			y := s.NewXMLBuffer()
			s.Tmpl.ExecuteTemplate(y, "rss.xml", n)
			err = s.WritePublic(section+"/index.xml", y.Bytes())
			return err
		}
	}
	return nil
}

func (s *Site) RenderHomePage() error {
	n := s.NewNode()
	n.Title = n.Site.Title
	n.Url = Urlize(string(n.Site.BaseUrl))
	n.RSSlink = permalink(s, "index.xml")
	n.Permalink = permalink(s, "")
	if len(s.Pages) > 0 {
		n.Date = s.Pages[0].Date
		if len(s.Pages) < 9 {
			n.Data["Pages"] = s.Pages
		} else {
			n.Data["Pages"] = s.Pages[:9]
		}
	}
	x, err := s.RenderThing(n, "index.html")
	if err != nil {
		return err
	}
	err = s.WritePublic("index.html", x.Bytes())
	if err != nil {
		return err
	}

	if a := s.Tmpl.Lookup("rss.xml"); a != nil {
		// XML Feed
		n.Url = Urlize("index.xml")
		n.Title = "Recent Content"
		n.Permalink = permalink(s, "index.xml")
		y := s.NewXMLBuffer()
		s.Tmpl.ExecuteTemplate(y, "rss.xml", n)
		err = s.WritePublic("index.xml", y.Bytes())
		return err
	}

	if a := s.Tmpl.Lookup("404.html"); a != nil {
		n.Url = Urlize("404.html")
		n.Title = "404 Page not found"
		n.Permalink = permalink(s, "404.html")
		x, err := s.RenderThing(n, "404.html")
		if err != nil {
			return err
		}
		err = s.WritePublic("404.html", x.Bytes())
		return err
	}

	return nil
}

func (s *Site) Stats() {
	fmt.Printf("%d pages created \n", len(s.Pages))
	for _, pl := range s.Config.Indexes {
		fmt.Printf("%d %s index created\n", len(s.Indexes[pl]), pl)
	}
}

func permalink(s *Site, plink string) HTML {
	return HTML(MakePermalink(string(s.Info.BaseUrl), plink))
}

func (s *Site) NewNode() *Node {
	return &Node{
		Data: make(map[string]interface{}),
		Site: s.Info,
	}
}

func (s *Site) RenderThing(d interface{}, layout string) (*bytes.Buffer, error) {
	if s.Tmpl.Lookup(layout) == nil {
		return nil, errors.New(fmt.Sprintf("Layout not found: %s", layout))
	}
	buffer := new(bytes.Buffer)
	err := s.Tmpl.ExecuteTemplate(buffer, layout, d)
	return buffer, err
}

func (s *Site) RenderThingOrDefault(d interface{}, layout string, defaultLayout string) (*bytes.Buffer, error) {
	content, err := s.RenderThing(d, layout)
	if err != nil {
		var err2 error
		content, err2 = s.RenderThing(d, defaultLayout)
		if err2 == nil {
			return content, err2
		}
	}
	return content, err
}

func (s *Site) NewXMLBuffer() *bytes.Buffer {
	header := "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n"
	return bytes.NewBufferString(header)
}

func (s *Site) WritePublic(path string, content []byte) (err error) {

	if s.Target != nil {
		return s.Target.Publish(path, bytes.NewReader(content))
	}

	if s.Config.Verbose {
		fmt.Println(path)
	}

	path, filename := filepath.Split(path)

	path = filepath.FromSlash(s.Config.GetAbsPath(filepath.Join(s.Config.PublishDir, path)))
	err = mkdirIf(path)
	if err != nil {
		return
	}

	file, _ := os.Create(filepath.Join(path, filename))
	defer file.Close()

	_, err = file.Write(content)
	return
}
