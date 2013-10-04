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
	"github.com/BurntSushi/toml"
	"github.com/spf13/hugo/parser"
	helper "github.com/spf13/hugo/template"
	"github.com/spf13/hugo/template/bundle"
	"github.com/theplant/blackfriday"
	"html/template"
	"io"
	"launchpad.net/goyaml"
	json "launchpad.net/rjson"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

type Page struct {
	Status      string
	Images      []string
	Content     template.HTML
	Summary     template.HTML
	RawMarkdown string // TODO should be []byte
	Params      map[string]interface{}
	contentType string
	Draft       bool
	Aliases     []string
	Tmpl        bundle.Template
	Markup      string
	renderable  bool
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
}

type Position struct {
	Prev *Page
	Next *Page
}

type Pages []*Page

func (p Pages) Len() int           { return len(p) }
func (p Pages) Less(i, j int) bool { return p[i].Date.Unix() > p[j].Date.Unix() }
func (p Pages) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// TODO eliminate unnecessary things
func (p Pages) Sort()             { sort.Sort(p) }
func (p Pages) Limit(n int) Pages { return p[0:n] }

func getSummaryString(content []byte, fmt string) []byte {
	if bytes.Contains(content, summaryDivider) {
		// If user defines split:
		// Split then render
		return renderBytes(bytes.Split(content, summaryDivider)[0], fmt)
	} else {
		// If hugo defines split:
		// render, strip html, then split
		plainContent := StripHTML(StripShortcodes(string(renderBytes(content, fmt))))
		return []byte(TruncateWordsToWholeSentence(plainContent, summaryLength))
	}
}

func renderBytes(content []byte, fmt string) []byte {
	switch fmt {
	default:
		return blackfriday.MarkdownCommon(content)
	case "markdown":
		return blackfriday.MarkdownCommon(content)
	case "rst":
		return []byte(getRstContent(content))
	}
}

// TODO abstract further to support loading from more
// than just files on disk. Should load reader (file, []byte)
func newPage(filename string) *Page {
	page := Page{contentType: "",
		File:   File{FileName: filename, Extension: "html"},
		Node:   Node{Keywords: make([]string, 10, 30)},
		Params: make(map[string]interface{})}
	page.Date, _ = time.Parse("20060102", "20080101")
	page.guessSection()
	return &page
}

func StripHTML(s string) string {
	output := ""

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {
		s = strings.Replace(s, "\n", " ", -1)
		s = strings.Replace(s, "</p>", " \n", -1)
		s = strings.Replace(s, "<br>", " \n", -1)
		s = strings.Replace(s, "</br>", " \n", -1)

		// Walk through the string removing all tags
		b := new(bytes.Buffer)
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}
	return output
}

func (p *Page) IsRenderable() bool {
	return p.renderable
}

func (p *Page) guessSection() {
	if p.Section == "" {
		x := strings.Split(p.FileName, "/")
		if len(x) > 1 {
			if section := x[len(x)-2]; section != "content" {
				p.Section = section
			}
		}
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

func (page *Page) Layout(l ...string) string {
	layout := ""
	if len(l) == 0 {
		layout = "single"
	} else {
		layout = l[0]
	}

	if x := page.layout; x != "" {
		return x
	}

	return strings.ToLower(page.Type()) + "/" + layout + ".html"
}

func ReadFrom(buf io.Reader, name string) (page *Page, err error) {
	if len(name) == 0 {
		return nil, errors.New("Zero length page name")
	}

	p := newPage(name)

	if err = p.parse(buf); err != nil {
		return
	}

	p.analyzePage()

	return p, nil
}

func (p *Page) analyzePage() {
	p.WordCount = TotalWords(p.RawMarkdown)
	p.FuzzyWordCount = int((p.WordCount+100)/100) * 100
}

func (p *Page) Permalink() (string, error) {
	baseUrl := string(p.Site.BaseUrl)
	section := strings.TrimSpace(p.Section)
	pSlug := strings.TrimSpace(p.Slug)
	pUrl := strings.TrimSpace(p.Url)
	var permalink string
	if len(pSlug) > 0 {
		if p.Site.Config.UglyUrls {
			permalink = section + "/" + p.Slug + "." + p.Extension
		} else {
			permalink = section + "/" + p.Slug + "/"
		}
	} else if len(pUrl) > 2 {
		permalink = pUrl
	} else {
		_, t := path.Split(p.FileName)
		if p.Site.Config.UglyUrls {
			x := replaceExtension(strings.TrimSpace(t), p.Extension)
			permalink = section + "/" + x
		} else {
			file, _ := fileExt(strings.TrimSpace(t))
			permalink = section + "/" + file
		}
	}

	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	path, err := url.Parse(permalink)
	if err != nil {
		return "", err
	}

	return MakePermalink(base, path).String(), nil
}

func (page *Page) handleTomlMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	datum = removeTomlIdentifier(datum)
	if _, err := toml.Decode(string(datum), &m); err != nil {
		return m, fmt.Errorf("Invalid TOML in %s \nError parsing page meta data: %s", page.FileName, err)
	}
	return m, nil
}

func removeTomlIdentifier(datum []byte) []byte {
	return bytes.Replace(datum, []byte("+++"), []byte(""), -1)
}

func (page *Page) handleYamlMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	if err := goyaml.Unmarshal(datum, &m); err != nil {
		return m, fmt.Errorf("Invalid YAML in %s \nError parsing page meta data: %s", page.FileName, err)
	}
	return m, nil
}

func (page *Page) handleJsonMetaData(datum []byte) (interface{}, error) {
	var f interface{}
	if err := json.Unmarshal(datum, &f); err != nil {
		return f, fmt.Errorf("Invalid JSON in %v \nError parsing page meta data: %s", page.FileName, err)
	}
	return f, nil
}

func (page *Page) update(f interface{}) error {
	m := f.(map[string]interface{})

	for k, v := range m {
		switch strings.ToLower(k) {
		case "title":
			page.Title = interfaceToString(v)
		case "description":
			page.Description = interfaceToString(v)
		case "slug":
			page.Slug = helper.Urlize(interfaceToString(v))
		case "url":
			if url := interfaceToString(v); strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				return fmt.Errorf("Only relative urls are supported, %v provided", url)
			}
			page.Url = helper.Urlize(interfaceToString(v))
		case "type":
			page.contentType = interfaceToString(v)
		case "keywords":
			page.Keywords = interfaceArrayToStringArray(v)
		case "date", "pubdate":
			page.Date = interfaceToStringToDate(v)
		case "draft":
			page.Draft = interfaceToBool(v)
		case "layout":
			page.layout = interfaceToString(v)
		case "markup":
			page.Markup = interfaceToString(v)
		case "aliases":
			page.Aliases = interfaceArrayToStringArray(v)
			for _, alias := range page.Aliases {
				if strings.HasPrefix(alias, "http://") || strings.HasPrefix(alias, "https://") {
					return fmt.Errorf("Only relative aliases are supported, %v provided", alias)
				}
			}
		case "status":
			page.Status = interfaceToString(v)
		default:
			// If not one of the explicit values, store in Params
			switch vv := v.(type) {
			case string: // handle string values
				page.Params[strings.ToLower(k)] = vv
			default: // handle array of strings as well
				switch vvv := vv.(type) {
				case []interface{}:
					var a = make([]string, len(vvv))
					for i, u := range vvv {
						a[i] = interfaceToString(u)
					}
					page.Params[strings.ToLower(k)] = a
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
	case string:
		return interfaceToString(v)
	case []string:
		return v
	}
	return nil
}

type frontmatterType struct {
	markstart, markend []byte
	parse              func([]byte) (interface{}, error)
	includeMark        bool
}

const YAML_DELIM = "---"
const TOML_DELIM = "+++"

func (page *Page) detectFrontMatter(mark rune) (f *frontmatterType) {
	switch mark {
	case '-':
		return &frontmatterType{[]byte(YAML_DELIM), []byte(YAML_DELIM), page.handleYamlMetaData, false}
	case '+':
		return &frontmatterType{[]byte(TOML_DELIM), []byte(TOML_DELIM), page.handleTomlMetaData, false}
	case '{':
		return &frontmatterType{[]byte{'{'}, []byte{'}'}, page.handleJsonMetaData, true}
	default:
		return nil
	}
}

func (p *Page) Render(layout ...string) template.HTML {
	curLayout := ""

	if len(layout) > 0 {
		curLayout = layout[0]
	}

	return template.HTML(string(p.ExecuteTemplate(curLayout).Bytes()))
}

func (p *Page) ExecuteTemplate(layout string) *bytes.Buffer {
	l := p.Layout(layout)
	buffer := new(bytes.Buffer)
	p.Tmpl.ExecuteTemplate(buffer, l, p)
	return buffer
}

func (page *Page) guessMarkupType() string {
	if page.Markup != "" {
		return page.Markup
	}

	if strings.HasSuffix(page.FileName, ".md") {
		return "md"
	}

	return "unknown"
}

func (page *Page) parse(reader io.Reader) error {
	p, err := parser.ReadFrom(reader)
	if err != nil {
		return err
	}

	page.renderable = p.IsRenderable()

	front := p.FrontMatter()

	if len(front) != 0 {
		fm := page.detectFrontMatter(rune(front[0]))
		meta, err := fm.parse(front)
		if err != nil {
			return err
		}

		if err = page.update(meta); err != nil {
			return err
		}
	}

	switch page.guessMarkupType() {
	case "md", "markdown", "mdown":
		page.convertMarkdown(bytes.NewReader(p.Content()))
	case "rst":
		page.convertRestructuredText(bytes.NewReader(p.Content()))
	case "html":
		fallthrough
	default:
		page.Content = template.HTML(p.Content())
	}
	return nil
}

func (page *Page) convertMarkdown(lines io.Reader) {
	b := new(bytes.Buffer)
	b.ReadFrom(lines)
	content := b.Bytes()
	page.Content = template.HTML(string(blackfriday.MarkdownCommon(RemoveSummaryDivider(content))))
	summary := getSummaryString(content, "markdown")
	page.Summary = template.HTML(string(summary))
}

func (page *Page) convertRestructuredText(lines io.Reader) {
	b := new(bytes.Buffer)
	b.ReadFrom(lines)
	content := b.Bytes()
	page.Content = template.HTML(getRstContent(content))
	summary := getSummaryString(content, "rst")
	page.Summary = template.HTML(string(summary))
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

	if len(strings.TrimSpace(p.Slug)) > 0 {
		outfile = strings.TrimSpace(p.Slug) + "." + p.Extension
	} else {
		// Fall back to filename
		_, t := path.Split(p.FileName)
		outfile = replaceExtension(strings.TrimSpace(t), p.Extension)
	}

	return path.Join(p.Dir, strings.TrimSpace(outfile))
}
