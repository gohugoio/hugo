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
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/parser"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
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

	extension         string
	contentType       string
	renderable        bool
	layout            string
	linkTitle         string
	frontmatter       []byte
	rawContent        []byte
	contentShortCodes map[string]string
	plain             string // TODO should be []byte
	RelatedPages      Pages
	relevance         Relevance
	PageMeta
	Source
	Position
	Node
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
	Prev *Page
	Next *Page
}

type Pages []*Page

func (p *Page) Plain() string {
	if len(p.plain) == 0 {
		p.plain = helpers.StripHTML(string(p.Content))
	}
	return p.plain
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

func (p *Page) UniqueId() string {
	return p.Source.UniqueId()
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
		// If user defines split:
		// Split, replace shortcode tokens, then render
		p.Truncated = true // by definition
		header := bytes.Split(p.rawContent, helpers.SummaryDivider)[0]
		renderedHeader := p.renderBytes(header)
		numShortcodesInHeader := bytes.Count(header, []byte(shortcodePlaceholderPrefix))
		if len(p.contentShortCodes) > 0 {
			tmpContentWithTokensReplaced, err :=
				replaceShortcodeTokens(renderedHeader, shortcodePlaceholderPrefix, numShortcodesInHeader, true, p.contentShortCodes)
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
		plain := strings.TrimSpace(p.Plain())
		p.Summary = helpers.BytesToHTML([]byte(helpers.TruncateWordsToWholeSentence(plain, helpers.SummaryLength)))
		p.Truncated = len(p.Summary) != len(plain)
	}
}

func (p *Page) renderBytes(content []byte) []byte {
	return helpers.RenderBytes(
		helpers.RenderingContext{Content: content, PageFmt: p.guessMarkupType(),
			DocumentId: p.UniqueId(), ConfigFlags: p.getRenderingConfigFlags()})
}

func (p *Page) renderContent(content []byte) []byte {
	return helpers.RenderBytesWithTOC(helpers.RenderingContext{Content: content, PageFmt: p.guessMarkupType(),
		DocumentId: p.UniqueId(), ConfigFlags: p.getRenderingConfigFlags()})
}

func (p *Page) getRenderingConfigFlags() map[string]bool {
	flags := make(map[string]bool)

	pageParam := p.GetParam("blackfriday")
	siteParam := viper.GetStringMap("blackfriday")

	flags = cast.ToStringMapBool(siteParam)

	if pageParam != nil {
		pageFlags := cast.ToStringMapBool(pageParam)
		for key, value := range pageFlags {
			flags[key] = value
		}
	}

	return flags
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

func (page *Page) Type() string {
	if page.contentType != "" {
		return page.contentType
	}

	if x := page.Section(); x != "" {
		return x
	}

	return "page"
}

func (page *Page) Section() string {
	return page.Source.Section()
}

func (page *Page) Layout(l ...string) []string {
	if page.layout != "" {
		return layouts(page.Type(), page.layout)
	}

	layout := ""
	if len(l) == 0 {
		layout = "single"
	} else {
		layout = l[0]
	}

	return layouts(page.Type(), layout)
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

func NewPageFrom(buf io.Reader, name string) (page *Page, err error) {
	p, err := NewPage(name)
	if err != nil {
		return p, err
	}
	err = p.ReadFrom(buf)

	return p, err
}

func NewPage(name string) (page *Page, err error) {
	if len(name) == 0 {
		return nil, errors.New("Zero length page name")
	}

	// Create new page
	p := newPage(name)

	return p, nil
}

func (p *Page) ReadFrom(buf io.Reader) (err error) {
	// Parse for metadata & body
	if err = p.parse(buf); err != nil {
		jww.ERROR.Print(err)
		return
	}

	return nil
}

func (p *Page) analyzePage() {
	p.WordCount = helpers.TotalWords(p.Plain())
	p.FuzzyWordCount = int((p.WordCount+100)/100) * 100
	p.ReadingTime = int((p.WordCount + 212) / 213)
}

func (p *Page) permalink() (*url.URL, error) {
	baseUrl := string(p.Site.BaseUrl)
	dir := strings.TrimSpace(filepath.ToSlash(p.Source.Dir()))
	pSlug := strings.TrimSpace(p.Slug)
	pUrl := strings.TrimSpace(p.Url)
	var permalink string
	var err error

	if len(pUrl) > 0 {
		return helpers.MakePermalink(baseUrl, pUrl), nil
	}

	if override, ok := p.Site.Permalinks[p.Section()]; ok {
		permalink, err = override.Expand(p)

		if err != nil {
			return nil, err
		}
		// fmt.Printf("have a section override for %q in section %s → %s\n", p.Title, p.Section, permalink)
	} else {
		if len(pSlug) > 0 {
			permalink = helpers.UrlPrep(viper.GetBool("UglyUrls"), path.Join(dir, p.Slug+"."+p.Extension()))
		} else {
			_, t := filepath.Split(p.Source.LogicalName())
			permalink = helpers.UrlPrep(viper.GetBool("UglyUrls"), path.Join(dir, helpers.ReplaceExtension(strings.TrimSpace(t), p.Extension())))
		}
	}

	return helpers.MakePermalink(baseUrl, permalink), nil
}

func (p *Page) Extension() string {
	if p.extension != "" {
		return p.extension
	} else {
		return viper.GetString("DefaultExtension")
	}
}

func (p *Page) LinkTitle() string {
	if len(p.linkTitle) > 0 {
		return p.linkTitle
	} else {
		return p.Title
	}
}

func (page *Page) ShouldBuild() bool {
	if viper.GetBool("BuildFuture") || page.PublishDate.IsZero() || page.PublishDate.Before(time.Now()) {
		if viper.GetBool("BuildDrafts") || !page.Draft {
			return true
		}
	}
	return false
}

func (page *Page) IsDraft() bool {
	return page.Draft
}

func (page *Page) IsFuture() bool {
	if page.PublishDate.Before(time.Now()) {
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

	link.Scheme = ""
	link.Host = ""
	link.User = nil
	link.Opaque = ""
	return link.String(), nil
}

func (page *Page) update(f interface{}) error {
	if f == nil {
		return fmt.Errorf("no metadata found")
	}
	m := f.(map[string]interface{})

	for k, v := range m {
		loki := strings.ToLower(k)
		switch loki {
		case "title":
			page.Title = cast.ToString(v)
		case "linktitle":
			page.linkTitle = cast.ToString(v)
		case "description":
			page.Description = cast.ToString(v)
		case "slug":
			page.Slug = helpers.Urlize(cast.ToString(v))
		case "url":
			if url := cast.ToString(v); strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				return fmt.Errorf("Only relative urls are supported, %v provided", url)
			}
			page.Url = helpers.Urlize(cast.ToString(v))
		case "type":
			page.contentType = cast.ToString(v)
		case "extension", "ext":
			page.extension = cast.ToString(v)
		case "keywords":
			page.Keywords = cast.ToStringSlice(v)
		case "date":
			page.Date = cast.ToTime(v)
		case "publishdate", "pubdate":
			page.PublishDate = cast.ToTime(v)
		case "draft":
			page.Draft = cast.ToBool(v)
		case "layout":
			page.layout = cast.ToString(v)
		case "markup":
			page.Markup = cast.ToString(v)
		case "weight":
			page.Weight = cast.ToInt(v)
		case "aliases":
			page.Aliases = cast.ToStringSlice(v)
			for _, alias := range page.Aliases {
				if strings.HasPrefix(alias, "http://") || strings.HasPrefix(alias, "https://") {
					return fmt.Errorf("Only relative aliases are supported, %v provided", alias)
				}
			}
		case "status":
			page.Status = cast.ToString(v)
		case "sitemap":
			page.Sitemap = parseSitemap(cast.ToStringMap(v))
		default:
			// If not one of the explicit values, store in Params
			switch vv := v.(type) {
			case bool:
				page.Params[loki] = vv
			case string:
				page.Params[loki] = vv
			case int64, int32, int16, int8, int:
				page.Params[loki] = vv
			case float64, float32:
				page.Params[loki] = vv
			case time.Time:
				page.Params[loki] = vv
			default: // handle array of strings as well
				switch vvv := vv.(type) {
				case []interface{}:
					var a = make([]string, len(vvv))
					for i, u := range vvv {
						a[i] = cast.ToString(u)
					}
					page.Params[loki] = a
				default:
					page.Params[loki] = vv
				}
			}
		}
	}
	return nil

}

func (page *Page) GetParam(key string) interface{} {
	v := page.Params[strings.ToLower(key)]

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
	case map[interface{}]interface{}:
		return v
	}
	return nil
}

func (page *Page) HasMenuCurrent(menu string, me *MenuEntry) bool {
	menus := page.Menus()

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

func (page *Page) IsMenuCurrent(menu string, inme *MenuEntry) bool {
	menus := page.Menus()

	if me, ok := menus[menu]; ok {
		return me.IsEqual(inme)
	}

	return false
}

func (page *Page) Menus() PageMenus {
	ret := PageMenus{}

	if ms, ok := page.Params["menu"]; ok {
		link, _ := page.Permalink()

		me := MenuEntry{Name: page.LinkTitle(), Weight: page.Weight, Url: link}

		// Could be the name of the menu to attach it to
		mname, err := cast.ToStringE(ms)

		if err == nil {
			me.Menu = mname
			ret[mname] = &me
			return ret
		}

		// Could be an slice of strings
		mnames, err := cast.ToStringSliceE(ms)

		if err == nil {
			for _, mname := range mnames {
				me.Menu = mname
				ret[mname] = &me
				return ret
			}
		}

		// Could be a structured menu entry
		menus, err := cast.ToStringMapE(ms)

		if err != nil {
			jww.ERROR.Printf("unable to process menus for %q\n", page.Title)
		}

		for name, menu := range menus {
			menuEntry := MenuEntry{Name: page.LinkTitle(), Url: link, Weight: page.Weight, Menu: name}
			jww.DEBUG.Printf("found menu: %q, in %q\n", name, page.Title)

			ime, err := cast.ToStringMapE(menu)
			if err != nil {
				jww.ERROR.Printf("unable to process menus for %q\n", page.Title)
			}

			menuEntry.MarshallMap(ime)
			ret[name] = &menuEntry
		}
		return ret
	}

	return nil
}

func (p *Page) Render(layout ...string) template.HTML {
	curLayout := ""

	if len(layout) > 0 {
		curLayout = layout[0]
	}

	return tpl.ExecuteTemplateToHTML(p, p.Layout(curLayout)...)
}

func (page *Page) guessMarkupType() string {
	// First try the explicitly set markup from the frontmatter
	if page.Markup != "" {
		format := helpers.GuessType(page.Markup)
		if format != "unknown" {
			return format
		}
	}

	return helpers.GuessType(page.Source.Ext())
}

func (page *Page) detectFrontMatter() (f *parser.FrontmatterType) {
	return parser.DetectFrontMatter(rune(page.frontmatter[0]))
}

func (page *Page) parse(reader io.Reader) error {
	psr, err := parser.ReadFrom(reader)
	if err != nil {
		return err
	}

	page.renderable = psr.IsRenderable()
	page.frontmatter = psr.FrontMatter()
	meta, err := psr.Metadata()
	if meta != nil {
		if err != nil {
			jww.ERROR.Printf("Error parsing page meta data for %s", page.File.Path())
			jww.ERROR.Println(err)
			return err
		}
		if err = page.update(meta); err != nil {
			return err
		}
	}

	page.rawContent = psr.Content()

	return nil
}

func (page *Page) SetSourceContent(content []byte) {
	page.Source.Content = content
}

func (page *Page) SetSourceMetaData(in interface{}, mark rune) (err error) {
	by, err := parser.InterfaceToFrontMatter(in, mark)
	if err != nil {
		return err
	}
	by = append(by, '\n')

	page.Source.Frontmatter = by

	return nil
}

func (page *Page) SafeSaveSourceAs(path string) error {
	return page.saveSourceAs(path, true)
}

func (page *Page) SaveSourceAs(path string) error {
	return page.saveSourceAs(path, false)
}

func (page *Page) saveSourceAs(path string, safe bool) error {
	b := new(bytes.Buffer)
	b.Write(page.Source.Frontmatter)
	b.Write(page.Source.Content)

	err := page.saveSource(b.Bytes(), path, safe)
	if err != nil {
		return err
	}
	return nil
}

func (page *Page) saveSource(by []byte, inpath string, safe bool) (err error) {
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

func (page *Page) SaveSource() error {
	return page.SaveSourceAs(page.FullFilePath())
}

func (p *Page) ProcessShortcodes(t tpl.Template) {

	// these short codes aren't used until after Page render,
	// but processed here to avoid coupling
	tmpContent, tmpContentShortCodes := extractAndRenderShortcodes(string(p.rawContent), p, t)
	p.rawContent = []byte(tmpContent)
	p.contentShortCodes = tmpContentShortCodes

}

// TODO(spf13): Remove this entirely
// Here for backwards compatibility & testing. Only works in isolation
func (page *Page) Convert() error {
	var h Handler
	if page.Markup != "" {
		h = FindHandler(page.Markup)
	} else {
		h = FindHandler(page.File.Extension())
	}
	if h != nil {
		h.PageConvert(page, tpl.T())
	}

	//// now we know enough to create a summary of the page and count some words
	page.setSummary()
	//analyze for raw stats
	page.analyzePage()

	return nil
}

func (p *Page) FullFilePath() string {
	return filepath.Join(p.Source.Dir(), p.Source.Path())
}

func (p *Page) TargetPath() (outfile string) {

	// Always use Url if it's specified
	if len(strings.TrimSpace(p.Url)) > 2 {
		outfile = strings.TrimSpace(p.Url)

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
