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
	"html/template"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/parser"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/theplant/blackfriday"
)

type Page struct {
	Status            string
	Images            []string
	rawContent        []byte
	Content           template.HTML
	Summary           template.HTML
	TableOfContents   template.HTML
	Truncated         bool
	plain             string // TODO should be []byte
	Params            map[string]interface{}
	contentType       string
	Draft             bool
	Aliases           []string
	Tmpl              Template
	Markup            string
	renderable        bool
	layout            string
	linkTitle         string
	frontmatter       []byte
	sourceFrontmatter []byte
	sourceContent     []byte
	PageMeta
	File
	Position
	Node
}

type File struct {
	FileName, Extension, Dir string
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
		p.plain = helpers.StripHTML(StripShortcodes(string(p.renderBytes(p.rawContent))))
	}
	return p.plain
}

func (p *Page) setSummary() {
	if bytes.Contains(p.rawContent, summaryDivider) {
		// If user defines split:
		// Split then render
		p.Truncated = true // by definition
		header := bytes.Split(p.rawContent, summaryDivider)[0]
		p.Summary = bytesToHTML(p.renderBytes(header))
	} else {
		// If hugo defines split:
		// render, strip html, then split
		plain := strings.TrimSpace(p.Plain())
		p.Summary = bytesToHTML([]byte(TruncateWordsToWholeSentence(plain, summaryLength)))
		p.Truncated = len(p.Summary) != len(plain)
	}
}

func stripEmptyNav(in []byte) []byte {
	return bytes.Replace(in, []byte("<nav>\n</nav>\n\n"), []byte(``), -1)
}

func bytesToHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

func (p *Page) renderBytes(content []byte) []byte {
	return renderBytes(content, p.guessMarkupType())
}

func (p *Page) renderContent(content []byte) []byte {
	return renderBytesWithTOC(content, p.guessMarkupType())
}

func renderBytesWithTOC(content []byte, pagefmt string) []byte {
	switch pagefmt {
	default:
		return markdownRenderWithTOC(content)
	case "markdown":
		return markdownRenderWithTOC(content)
	case "rst":
		return []byte(getRstContent(content))
	}
}

func renderBytes(content []byte, pagefmt string) []byte {
	switch pagefmt {
	default:
		return markdownRender(content)
	case "markdown":
		return markdownRender(content)
	case "rst":
		return []byte(getRstContent(content))
	}
}

func newPage(filename string) *Page {
	page := Page{contentType: "",
		File:   File{FileName: filename, Extension: "html"},
		Node:   Node{Keywords: []string{}, Sitemap: Sitemap{Priority: -1}},
		Params: make(map[string]interface{})}

	jww.DEBUG.Println("Reading from", page.File.FileName)
	page.Date, _ = time.Parse("20060102", "20080101")
	page.guessSection()
	return &page
}

func (p *Page) IsRenderable() bool {
	return p.renderable
}

func (p *Page) guessSection() {
	if p.Section == "" {
		p.Section = helpers.GuessSection(p.FileName)
	}
}

func (page *Page) Type() string {
	if page.contentType != "" {
		return page.contentType
	}
	page.guessSection()
	if x := page.Section; x != "" {
		return x
	}

	return "page"
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
	for i := range t {
		search := t[:len(t)-i]
		layouts = append(layouts, fmt.Sprintf("%s/%s.html", strings.ToLower(path.Join(search...)), layout))
	}
	layouts = append(layouts, fmt.Sprintf("%s.html", layout))
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

	//analyze for raw stats
	p.analyzePage()

	return nil
}

func (p *Page) analyzePage() {
	p.WordCount = TotalWords(p.Plain())
	p.FuzzyWordCount = int((p.WordCount+100)/100) * 100
	p.ReadingTime = int((p.WordCount + 212) / 213)
}

func (p *Page) permalink() (*url.URL, error) {
	baseUrl := string(p.Site.BaseUrl)
	dir := strings.TrimSpace(p.Dir)
	pSlug := strings.TrimSpace(p.Slug)
	pUrl := strings.TrimSpace(p.Url)
	var permalink string
	var err error

	if len(pUrl) > 0 {
		return helpers.MakePermalink(baseUrl, pUrl), nil
	}

	if override, ok := p.Site.Permalinks[p.Section]; ok {
		permalink, err = override.Expand(p)

		if err != nil {
			return nil, err
		}
		// fmt.Printf("have a section override for %q in section %s → %s\n", p.Title, p.Section, permalink)
	} else {
		if len(pSlug) > 0 {
			permalink = helpers.UrlPrep(viper.GetBool("UglyUrls"), path.Join(dir, p.Slug+"."+p.Extension))
		} else {
			_, t := path.Split(p.FileName)
			permalink = helpers.UrlPrep(viper.GetBool("UglyUrls"), path.Join(dir, helpers.ReplaceExtension(strings.TrimSpace(t), p.Extension)))
		}
	}

	return helpers.MakePermalink(baseUrl, permalink), nil
}

func (p *Page) LinkTitle() string {
	if len(p.linkTitle) > 0 {
		return p.linkTitle
	} else {
		return p.Title
	}
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
		case "keywords":
			page.Keywords = cast.ToStringSlice(v)
		case "date", "pubdate":
			page.Date = cast.ToTime(v)
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

func (page *Page) HasMenuCurrent(menu string, me *MenuEntry) bool {
	menus := page.Menus()

	if m, ok := menus[menu]; ok {
		if me.HasChildren() {
			for _, child := range me.Children {
				if child.Name == m.Name {
					return true
				}
			}
		}
	}

	return false

}

func (page *Page) IsMenuCurrent(menu string, name string) bool {
	menus := page.Menus()

	if me, ok := menus[menu]; ok {
		return me.Name == name
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

	return bytesToHTML(p.ExecuteTemplate(curLayout).Bytes())
}

func (p *Page) ExecuteTemplate(layout string) *bytes.Buffer {
	l := p.Layout(layout)
	buffer := new(bytes.Buffer)
	for _, layout := range l {
		if p.Tmpl.Lookup(layout) != nil {
			p.Tmpl.ExecuteTemplate(buffer, layout, p)
			break
		}
	}
	return buffer
}

func (page *Page) guessMarkupType() string {
	// First try the explicitly set markup from the frontmatter
	if page.Markup != "" {
		format := guessType(page.Markup)
		if format != "unknown" {
			return format
		}
	}

	// Then try to guess from the extension
	ext := strings.ToLower(path.Ext(page.FileName))
	if strings.HasPrefix(ext, ".") {
		return guessType(ext[1:])
	}

	return "unknown"
}

func guessType(in string) string {
	switch strings.ToLower(in) {
	case "md", "markdown", "mdown":
		return "markdown"
	case "rst":
		return "rst"
	case "html", "htm":
		return "html"
	}
	return "unknown"
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
			jww.ERROR.Printf("Error parsing page meta data for %s", page.FileName)
			jww.ERROR.Println(err)
			return err
		}
		if err = page.update(meta); err != nil {
			return err
		}
	}

	page.rawContent = psr.Content()
	page.setSummary()

	return nil
}

func (page *Page) SetSourceContent(content []byte) {
	page.sourceContent = content
}

func (page *Page) SetSourceMetaData(in interface{}, mark rune) (err error) {
	by, err := parser.InterfaceToFrontMatter(in, mark)
	if err != nil {
		return err
	}
	by = append(by, '\n')

	page.sourceFrontmatter = by

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
	b.Write(page.sourceFrontmatter)
	b.Write(page.sourceContent)

	err := page.saveSource(b.Bytes(), path, safe)
	if err != nil {
		return err
	}
	return nil
}

func (page *Page) saveSource(by []byte, inpath string, safe bool) (err error) {
	if !path.IsAbs(inpath) {
		inpath = helpers.AbsPathify(inpath)
	}
	jww.INFO.Println("creating", inpath)

	if safe {
		err = helpers.SafeWriteToDisk(inpath, bytes.NewReader(by))
	} else {
		err = helpers.WriteToDisk(inpath, bytes.NewReader(by))
	}
	if err != nil {
		return
	}
	return nil
}

func (page *Page) SaveSource() error {
	return page.SaveSourceAs(page.FullFilePath())
}

func (p *Page) ProcessShortcodes(t Template) {
	p.rawContent = []byte(ShortcodesHandle(string(p.rawContent), p, t))
	p.Summary = template.HTML(ShortcodesHandle(string(p.Summary), p, t))
}

func (page *Page) Convert() error {
	markupType := page.guessMarkupType()
	switch markupType {
	case "markdown", "rst":
		tmpContent, tmpTableOfContents := extractTOC(page.renderContent(RemoveSummaryDivider(page.rawContent)))
		page.Content = bytesToHTML(tmpContent)
		page.TableOfContents = bytesToHTML(tmpTableOfContents)
	case "html":
		page.Content = bytesToHTML(page.rawContent)
	default:
		return fmt.Errorf("Error converting unsupported file type '%s' for page '%s'", markupType, page.FileName)
	}
	return nil
}

func markdownRender(content []byte) []byte {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_SKIP_SCRIPT
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	return blackfriday.Markdown(content, renderer, extensions)
}

func markdownRenderWithTOC(content []byte) []byte {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_SKIP_SCRIPT
	htmlFlags |= blackfriday.HTML_TOC
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	return blackfriday.Markdown(content, renderer, extensions)
}

func extractTOC(content []byte) (newcontent []byte, toc []byte) {
	origContent := make([]byte, len(content))
	copy(origContent, content)
	first := []byte(`<nav>
<ul>`)

	last := []byte(`</ul>
</nav>`)

	replacement := []byte(`<nav id="TableOfContents">
<ul>`)

	startOfTOC := bytes.Index(content, first)

	peekEnd := len(content)
	if peekEnd > 70+startOfTOC {
		peekEnd = 70 + startOfTOC
	}

	if startOfTOC < 0 {
		return stripEmptyNav(content), toc
	}
	// Need to peek ahead to see if this nav element is actually the right one.
	correctNav := bytes.Index(content[startOfTOC:peekEnd], []byte(`#toc_0`))
	if correctNav < 0 { // no match found
		return content, toc
	}
	lengthOfTOC := bytes.Index(content[startOfTOC:], last) + len(last)
	endOfTOC := startOfTOC + lengthOfTOC

	newcontent = append(content[:startOfTOC], content[endOfTOC:]...)
	toc = append(replacement, origContent[startOfTOC+len(first):endOfTOC]...)
	return
}

func ReaderToBytes(lines io.Reader) []byte {
	b := new(bytes.Buffer)
	b.ReadFrom(lines)
	return b.Bytes()
}

func (p *Page) FullFilePath() string {
	return path.Join(p.Dir, p.FileName)
}

func (p *Page) TargetPath() (outfile string) {

	// Always use Url if it's specified
	if len(strings.TrimSpace(p.Url)) > 2 {
		outfile = strings.TrimSpace(p.Url)

		if strings.HasSuffix(outfile, "/") {
			outfile = outfile + "index.html"
		}
		return
	}

	// If there's a Permalink specification, we use that
	if override, ok := p.Site.Permalinks[p.Section]; ok {
		var err error
		outfile, err = override.Expand(p)
		if err == nil {
			if strings.HasSuffix(outfile, "/") {
				outfile += "index.html"
			}
			return
		}
	}

	if len(strings.TrimSpace(p.Slug)) > 0 {
		outfile = strings.TrimSpace(p.Slug) + "." + p.Extension
	} else {
		// Fall back to filename
		_, t := path.Split(p.FileName)
		outfile = helpers.ReplaceExtension(strings.TrimSpace(t), p.Extension)
	}

	return path.Join(p.Dir, strings.TrimSpace(outfile))
}
