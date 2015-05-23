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
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/parser"

	"html/template"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type Page struct {
	Params          map[string]interface{}
	Content         template.HTML
	Summary         template.HTML
	Aliases         []string
	Status          string
	Images          []Image
	Videos          []Video
	TableOfContents template.HTML
	Truncated       bool
	Draft           bool
	PublishDate     time.Time
	Tmpl            tpl.Template
	Markup          string

	extension           string
	contentType         string
	renderable          bool
	layout              string
	linkTitle           string
	frontmatter         []byte
	rawContent          []byte
	contentShortCodes   map[string]string
	plain               string // TODO should be []byte
	plainWords          []string
	plainInit           sync.Once
	renderingConfig     *helpers.Blackfriday
	renderingConfigInit sync.Once
	PageMeta
	Source
	Position `json:"-"`
	Node
	pageMenus     PageMenus
	pageMenusInit sync.Once
}

type Source struct {
	Frontmatter []byte
	Content     []byte
	source.File
}
type PageMeta struct {
	WordCount      int
	FuzzyWordCount int
	ReadingTime    int
	Weight         int
}

type Position struct {
	Prev          *Page
	Next          *Page
	PrevInSection *Page
	NextInSection *Page
}

type Pages []*Page

func (p *Page) Plain() string {
	p.initPlain()
	return p.plain
}

func (p *Page) PlainWords() []string {
	p.initPlain()
	return p.plainWords
}

func (p *Page) initPlain() {
	p.plainInit.Do(func() {
		p.plain = helpers.StripHTML(string(p.Content))
		p.plainWords = strings.Fields(p.plain)
	})
}

func (p *Page) IsNode() bool {
	return false
}

func (p *Page) IsPage() bool {
	return true
}

func (p *Page) Author() Author {
	authors := p.Authors()

	for _, author := range authors {
		return author
	}
	return Author{}
}

func (p *Page) Authors() AuthorList {
	authorKeys, ok := p.Params["authors"]
	authors := authorKeys.([]string)
	if !ok || len(authors) < 1 || len(p.Site.Authors) < 1 {
		return AuthorList{}
	}

	al := make(AuthorList)
	for _, author := range authors {
		a, ok := p.Site.Authors[author]
		if ok {
			al[author] = a
		}
	}
	return al
}

func (p *Page) UniqueID() string {
	return p.Source.UniqueID()
}

func (p *Page) Ref(ref string) (string, error) {
	return p.Node.Site.Ref(ref, p)
}

func (p *Page) RelRef(ref string) (string, error) {
	return p.Node.Site.RelRef(ref, p)
}

// for logging
func (p *Page) lineNumRawContentStart() int {
	return bytes.Count(p.frontmatter, []byte("\n")) + 1
}

func (p *Page) setSummary() {

	// at this point, p.rawContent contains placeholders for the short codes,
	// rendered and ready in p.contentShortcodes

	if bytes.Contains(p.rawContent, helpers.SummaryDivider) {
		sections := bytes.Split(p.rawContent, helpers.SummaryDivider)
		header := sections[0]
		p.Truncated = true
		if len(sections[1]) < 20 {
			// only whitespace?
			p.Truncated = len(bytes.Trim(sections[1], " \n\r")) > 0
		}

		renderedHeader := p.renderBytes(header)
		if len(p.contentShortCodes) > 0 {
			tmpContentWithTokensReplaced, err :=
				replaceShortcodeTokens(renderedHeader, shortcodePlaceholderPrefix, true, p.contentShortCodes)
			if err != nil {
				jww.FATAL.Printf("Failed to replace short code tokens in Summary for %s:\n%s", p.BaseFileName(), err.Error())
			} else {
				renderedHeader = tmpContentWithTokensReplaced
			}
		}
		p.Summary = helpers.BytesToHTML(renderedHeader)
	} else {
		// If hugo defines split:
		// render, strip html, then split
		summary, truncated := helpers.TruncateWordsToWholeSentence(p.PlainWords(), helpers.SummaryLength)
		p.Summary = template.HTML(summary)
		p.Truncated = truncated

	}
}

func (p *Page) renderBytes(content []byte) []byte {
	return helpers.RenderBytes(
		&helpers.RenderingContext{Content: content, PageFmt: p.guessMarkupType(),
			DocumentID: p.UniqueID(), Config: p.getRenderingConfig()})
}

func (p *Page) renderContent(content []byte) []byte {
	return helpers.RenderBytesWithTOC(&helpers.RenderingContext{Content: content, PageFmt: p.guessMarkupType(),
		DocumentID: p.UniqueID(), Config: p.getRenderingConfig()})
}

func (p *Page) getRenderingConfig() *helpers.Blackfriday {

	p.renderingConfigInit.Do(func() {
		pageParam := p.GetParam("blackfriday")
		siteParam := viper.GetStringMap("blackfriday")

		combinedParam := siteParam

		if pageParam != nil {
			combinedParam = make(map[string]interface{})

			for k, v := range siteParam {
				combinedParam[k] = v
			}

			pageConfig := cast.ToStringMap(pageParam)

			for key, value := range pageConfig {
				combinedParam[key] = value
			}
		}
		p.renderingConfig = helpers.NewBlackfriday()
		if err := mapstructure.Decode(combinedParam, p.renderingConfig); err != nil {
			jww.FATAL.Printf("Failed to get rendering config for %s:\n%s", p.BaseFileName(), err.Error())
		}
	})

	return p.renderingConfig
}

func newPage(filename string) *Page {
	page := Page{contentType: "",
		Source: Source{File: *source.NewFile(filename)},
		Node:   Node{Keywords: []string{}, Sitemap: Sitemap{Priority: -1}},
		Params: make(map[string]interface{})}

	jww.DEBUG.Println("Reading from", page.File.Path())
	return &page
}

func (p *Page) IsRenderable() bool {
	return p.renderable
}

func (p *Page) Type() string {
	if p.contentType != "" {
		return p.contentType
	}

	if x := p.Section(); x != "" {
		return x
	}

	return "page"
}

func (p *Page) Section() string {
	return p.Source.Section()
}

func (p *Page) Layout(l ...string) []string {
	if p.layout != "" {
		return layouts(p.Type(), p.layout)
	}

	layout := ""
	if len(l) == 0 {
		layout = "single"
	} else {
		layout = l[0]
	}

	return layouts(p.Type(), layout)
}

func layouts(types string, layout string) (layouts []string) {
	t := strings.Split(types, "/")

	// Add type/layout.html
	for i := range t {
		search := t[:len(t)-i]
		layouts = append(layouts, fmt.Sprintf("%s/%s.html", strings.ToLower(path.Join(search...)), layout))
	}

	// Add _default/layout.html
	layouts = append(layouts, fmt.Sprintf("_default/%s.html", layout))

	// Add theme/type/layout.html & theme/_default/layout.html
	for _, l := range layouts {
		layouts = append(layouts, "theme/"+l)
	}

	return
}

func NewPageFrom(buf io.Reader, name string) (*Page, error) {
	p, err := NewPage(name)
	if err != nil {
		return p, err
	}
	_, err = p.ReadFrom(buf)

	return p, err
}

func NewPage(name string) (*Page, error) {
	if len(name) == 0 {
		return nil, errors.New("Zero length page name")
	}

	// Create new page
	p := newPage(name)

	return p, nil
}

func (p *Page) ReadFrom(buf io.Reader) (int64, error) {
	// Parse for metadata & body
	if err := p.parse(buf); err != nil {
		jww.ERROR.Print(err)
		return 0, err
	}

	return int64(len(p.rawContent)), nil
}

func (p *Page) analyzePage() {
	p.WordCount = len(p.PlainWords())
	p.FuzzyWordCount = int((p.WordCount+100)/100) * 100
	p.ReadingTime = int((p.WordCount + 212) / 213)
}

func (p *Page) permalink() (*url.URL, error) {
	baseURL := string(p.Site.BaseURL)
	dir := strings.TrimSpace(filepath.ToSlash(p.Source.Dir()))
	pSlug := strings.TrimSpace(helpers.URLize(p.Slug))
	pURL := strings.TrimSpace(helpers.URLize(p.URL))
	var permalink string
	var err error

	if len(pURL) > 0 {
		return helpers.MakePermalink(baseURL, pURL), nil
	}

	if override, ok := p.Site.Permalinks[p.Section()]; ok {
		permalink, err = override.Expand(p)

		if err != nil {
			return nil, err
		}
		// fmt.Printf("have a section override for %q in section %s → %s\n", p.Title, p.Section, permalink)
	} else {
		if len(pSlug) > 0 {
			permalink = helpers.URLPrep(viper.GetBool("UglyURLs"), path.Join(dir, p.Slug+"."+p.Extension()))
		} else {
			_, t := filepath.Split(p.Source.LogicalName())
			permalink = helpers.URLPrep(viper.GetBool("UglyURLs"), path.Join(dir, helpers.ReplaceExtension(strings.TrimSpace(t), p.Extension())))
		}
	}

	return helpers.MakePermalink(baseURL, permalink), nil
}

func (p *Page) Extension() string {
	if p.extension != "" {
		return p.extension
	}
	return viper.GetString("DefaultExtension")
}

func (p *Page) LinkTitle() string {
	if len(p.linkTitle) > 0 {
		return p.linkTitle
	}
	return p.Title
}

func (p *Page) ShouldBuild() bool {
	if viper.GetBool("BuildFuture") || p.PublishDate.IsZero() || p.PublishDate.Before(time.Now()) {
		if viper.GetBool("BuildDrafts") || !p.Draft {
			return true
		}
	}
	return false
}

func (p *Page) IsDraft() bool {
	return p.Draft
}

func (p *Page) IsFuture() bool {
	if p.PublishDate.Before(time.Now()) {
		return false
	}
	return true
}

func (p *Page) Permalink() (string, error) {
	link, err := p.permalink()
	if err != nil {
		return "", err
	}
	return link.String(), nil
}

func (p *Page) RelPermalink() (string, error) {
	link, err := p.permalink()
	if err != nil {
		return "", err
	}

	if viper.GetBool("CanonifyURLs") {
		// replacements for relpermalink with baseURL on the form http://myhost.com/sub/ will fail later on
		// have to return the URL relative from baseURL
		relpath, err := helpers.GetRelativePath(link.String(), string(p.Site.BaseURL))
		if err != nil {
			return "", err
		}
		return "/" + filepath.ToSlash(relpath), nil
	}

	link.Scheme = ""
	link.Host = ""
	link.User = nil
	link.Opaque = ""
	return link.String(), nil
}

func (p *Page) update(f interface{}) error {
	if f == nil {
		return fmt.Errorf("no metadata found")
	}
	m := f.(map[string]interface{})
	var err error
	for k, v := range m {
		loki := strings.ToLower(k)
		switch loki {
		case "title":
			p.Title = cast.ToString(v)
		case "linktitle":
			p.linkTitle = cast.ToString(v)
		case "description":
			p.Description = cast.ToString(v)
		case "slug":
			p.Slug = cast.ToString(v)
		case "url":
			if url := cast.ToString(v); strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				return fmt.Errorf("Only relative URLs are supported, %v provided", url)
			}
			p.URL = cast.ToString(v)
		case "type":
			p.contentType = cast.ToString(v)
		case "extension", "ext":
			p.extension = cast.ToString(v)
		case "keywords":
			p.Keywords = cast.ToStringSlice(v)
		case "date":
			p.Date, err = cast.ToTimeE(v)
			if err != nil {
				jww.ERROR.Printf("Failed to parse date '%v' in page %s", v, p.File.Path())
			}
		case "lastmod":
			p.Lastmod, err = cast.ToTimeE(v)
			if err != nil {
				jww.ERROR.Printf("Failed to parse lastmod '%v' in page %s", v, p.File.Path())
			}
		case "publishdate", "pubdate":
			p.PublishDate, err = cast.ToTimeE(v)
			if err != nil {
				jww.ERROR.Printf("Failed to parse publishdate '%v' in page %s", v, p.File.Path())
			}
		case "draft":
			p.Draft = cast.ToBool(v)
		case "layout":
			p.layout = cast.ToString(v)
		case "markup":
			p.Markup = cast.ToString(v)
		case "weight":
			p.Weight = cast.ToInt(v)
		case "aliases":
			p.Aliases = cast.ToStringSlice(v)
			for _, alias := range p.Aliases {
				if strings.HasPrefix(alias, "http://") || strings.HasPrefix(alias, "https://") {
					return fmt.Errorf("Only relative aliases are supported, %v provided", alias)
				}
			}
		case "status":
			p.Status = cast.ToString(v)
		case "sitemap":
			p.Sitemap = parseSitemap(cast.ToStringMap(v))
		default:
			// If not one of the explicit values, store in Params
			switch vv := v.(type) {
			case bool:
				p.Params[loki] = vv
			case string:
				p.Params[loki] = vv
			case int64, int32, int16, int8, int:
				p.Params[loki] = vv
			case float64, float32:
				p.Params[loki] = vv
			case time.Time:
				p.Params[loki] = vv
			default: // handle array of strings as well
				switch vvv := vv.(type) {
				case []interface{}:
					var a = make([]string, len(vvv))
					for i, u := range vvv {
						a[i] = cast.ToString(u)
					}
					p.Params[loki] = a
				default:
					p.Params[loki] = vv
				}
			}
		}
	}

	if p.Lastmod.IsZero() {
		p.Lastmod = p.Date
	}

	return nil

}

func (p *Page) GetParam(key string) interface{} {
	v := p.Params[strings.ToLower(key)]

	if v == nil {
		return nil
	}

	switch v.(type) {
	case bool:
		return cast.ToBool(v)
	case string:
		return strings.ToLower(cast.ToString(v))
	case int64, int32, int16, int8, int:
		return cast.ToInt(v)
	case float64, float32:
		return cast.ToFloat64(v)
	case time.Time:
		return cast.ToTime(v)
	case []string:
		return helpers.SliceToLower(v.([]string))
	case map[string]interface{}: // JSON and TOML
		return v
	case map[interface{}]interface{}: // YAML
		return v
	}

	jww.ERROR.Printf("GetParam(\"%s\"): Unknown type %s\n", key, reflect.TypeOf(v))
	return nil
}

func (p *Page) HasMenuCurrent(menu string, me *MenuEntry) bool {
	menus := p.Menus()
	sectionPagesMenu := viper.GetString("SectionPagesMenu")

	// page is labeled as "shadow-member" of the menu with the same identifier as the section
	if sectionPagesMenu != "" && p.Section() != "" && sectionPagesMenu == menu && p.Section() == me.Identifier {
		return true
	}

	if m, ok := menus[menu]; ok {
		if me.HasChildren() {
			for _, child := range me.Children {
				if child.IsEqual(m) {
					return true
				}
			}
		}
	}

	return false

}

func (p *Page) IsMenuCurrent(menu string, inme *MenuEntry) bool {
	menus := p.Menus()

	if me, ok := menus[menu]; ok {
		return me.IsEqual(inme)
	}

	return false
}

func (p *Page) Menus() PageMenus {
	p.pageMenusInit.Do(func() {
		p.pageMenus = PageMenus{}

		if ms, ok := p.Params["menu"]; ok {
			link, _ := p.RelPermalink()

			me := MenuEntry{Name: p.LinkTitle(), Weight: p.Weight, URL: link}

			// Could be the name of the menu to attach it to
			mname, err := cast.ToStringE(ms)

			if err == nil {
				me.Menu = mname
				p.pageMenus[mname] = &me
				return
			}

			// Could be a slice of strings
			mnames, err := cast.ToStringSliceE(ms)

			if err == nil {
				for _, mname := range mnames {
					me.Menu = mname
					p.pageMenus[mname] = &me
					return
				}
			}

			// Could be a structured menu entry
			menus, err := cast.ToStringMapE(ms)

			if err != nil {
				jww.ERROR.Printf("unable to process menus for %q\n", p.Title)
			}

			for name, menu := range menus {
				menuEntry := MenuEntry{Name: p.LinkTitle(), URL: link, Weight: p.Weight, Menu: name}
				jww.DEBUG.Printf("found menu: %q, in %q\n", name, p.Title)

				ime, err := cast.ToStringMapE(menu)
				if err != nil {
					jww.ERROR.Printf("unable to process menus for %q\n", p.Title)
				}

				menuEntry.MarshallMap(ime)
				p.pageMenus[name] = &menuEntry
			}
		}
	})

	return p.pageMenus
}

func (p *Page) Render(layout ...string) template.HTML {
	curLayout := ""

	if len(layout) > 0 {
		curLayout = layout[0]
	}

	return tpl.ExecuteTemplateToHTML(p, p.Layout(curLayout)...)
}

func (p *Page) guessMarkupType() string {
	// First try the explicitly set markup from the frontmatter
	if p.Markup != "" {
		format := helpers.GuessType(p.Markup)
		if format != "unknown" {
			return format
		}
	}

	return helpers.GuessType(p.Source.Ext())
}

func (p *Page) detectFrontMatter() (f *parser.FrontmatterType) {
	return parser.DetectFrontMatter(rune(p.frontmatter[0]))
}

func (p *Page) parse(reader io.Reader) error {
	psr, err := parser.ReadFrom(reader)
	if err != nil {
		return err
	}

	p.renderable = psr.IsRenderable()
	p.frontmatter = psr.FrontMatter()
	meta, err := psr.Metadata()
	if meta != nil {
		if err != nil {
			jww.ERROR.Printf("Error parsing page meta data for %s", p.File.Path())
			jww.ERROR.Println(err)
			return err
		}
		if err = p.update(meta); err != nil {
			return err
		}
	}

	p.rawContent = psr.Content()

	return nil
}

func (p *Page) SetSourceContent(content []byte) {
	p.Source.Content = content
}

func (p *Page) SetSourceMetaData(in interface{}, mark rune) (err error) {
	by, err := parser.InterfaceToFrontMatter(in, mark)
	if err != nil {
		return err
	}
	by = append(by, '\n')

	p.Source.Frontmatter = by

	return nil
}

func (p *Page) SafeSaveSourceAs(path string) error {
	return p.saveSourceAs(path, true)
}

func (p *Page) SaveSourceAs(path string) error {
	return p.saveSourceAs(path, false)
}

func (p *Page) saveSourceAs(path string, safe bool) error {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)

	b.Write(p.Source.Frontmatter)
	b.Write(p.Source.Content)

	bc := make([]byte, b.Len(), b.Len())
	copy(bc, b.Bytes())

	err := p.saveSource(bc, path, safe)
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) saveSource(by []byte, inpath string, safe bool) (err error) {
	if !filepath.IsAbs(inpath) {
		inpath = helpers.AbsPathify(inpath)
	}
	jww.INFO.Println("creating", inpath)

	if safe {
		err = helpers.SafeWriteToDisk(inpath, bytes.NewReader(by), hugofs.SourceFs)
	} else {
		err = helpers.WriteToDisk(inpath, bytes.NewReader(by), hugofs.SourceFs)
	}
	if err != nil {
		return
	}
	return nil
}

func (p *Page) SaveSource() error {
	return p.SaveSourceAs(p.FullFilePath())
}

func (p *Page) ProcessShortcodes(t tpl.Template) {

	// these short codes aren't used until after Page render,
	// but processed here to avoid coupling
	tmpContent, tmpContentShortCodes, _ := extractAndRenderShortcodes(string(p.rawContent), p, t)
	p.rawContent = []byte(tmpContent)
	p.contentShortCodes = tmpContentShortCodes

}

// TODO(spf13): Remove this entirely
// Here for backwards compatibility & testing. Only works in isolation
func (p *Page) Convert() error {
	var h Handler
	if p.Markup != "" {
		h = FindHandler(p.Markup)
	} else {
		h = FindHandler(p.File.Extension())
	}
	if h != nil {
		h.PageConvert(p, tpl.T())
	}

	//// now we know enough to create a summary of the page and count some words
	p.setSummary()
	//analyze for raw stats
	p.analyzePage()

	return nil
}

func (p *Page) FullFilePath() string {
	return filepath.Join(p.Source.Dir(), p.Source.Path())
}

func (p *Page) TargetPath() (outfile string) {

	// Always use URL if it's specified
	if len(strings.TrimSpace(p.URL)) > 2 {
		outfile = strings.TrimSpace(p.URL)

		if strings.HasSuffix(outfile, "/") {
			outfile = outfile + "index.html"
		}
		outfile = filepath.FromSlash(outfile)
		return
	}

	// If there's a Permalink specification, we use that
	if override, ok := p.Site.Permalinks[p.Section()]; ok {
		var err error
		outfile, err = override.Expand(p)
		if err == nil {
			outfile, _ = url.QueryUnescape(outfile)
			if strings.HasSuffix(outfile, "/") {
				outfile += "index.html"
			}
			outfile = filepath.FromSlash(outfile)
			return
		}
	}

	if len(strings.TrimSpace(p.Slug)) > 0 {
		outfile = strings.TrimSpace(p.Slug) + "." + p.Extension()
	} else {
		// Fall back to filename
		outfile = helpers.ReplaceExtension(p.Source.LogicalName(), p.Extension())
	}

	return filepath.Join(p.Source.Dir(), strings.TrimSpace(outfile))
}
