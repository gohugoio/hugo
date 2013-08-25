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
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/theplant/blackfriday"
	"github.com/noahcampbell/akebia/parser"
	"html/template"
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var _ = filepath.Base("")

type Page struct {
	Status          string
	Images          []string
	Content         template.HTML
	Summary         template.HTML
	RawMarkdown     string // TODO should be []byte
	Params          map[string]interface{}
	RenderedContent *bytes.Buffer
	contentType     string
	Draft           bool
	Aliases         []string
	Tmpl            *template.Template
	Markup          string
	PageMeta
	File
	Position
	Node
}

type File struct {
	FileName, OutFile, Extension string
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

// TODO abstract further to support loading from more
// than just files on disk. Should load reader (file, []byte)
func NewPage(filename string) *Page {
	page := Page{contentType: "",
		File:   File{FileName: filename, Extension: "html"},
		Node:   Node{Keywords: make([]string, 10, 30)},
		Params: make(map[string]interface{}),
		Markup: "md"}
	page.Date, _ = time.Parse("20060102", "20080101")
	return &page
}

func (page *Page) Initalize() error {
	err := page.buildPageFromFile()
	if err != nil {
		return err
	}
	page.analyzePage()
	return nil
}

func (p *Page) guessSection() {
	if p.Section == "" {
		x := strings.Split(p.FileName, string(os.PathSeparator))
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

	p := NewPage(name)

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

func (p *Page) Permalink() template.HTML {
	if len(strings.TrimSpace(p.Slug)) > 0 {
		if p.Site.Config.UglyUrls {
			return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Section)+"/"+p.Slug+"."+p.Extension))
		} else {
			return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Section)+"/"+p.Slug+"/"))
		}
	} else if len(strings.TrimSpace(p.Url)) > 2 {
		return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Url)))
	} else {
		_, t := filepath.Split(p.FileName)
		if p.Site.Config.UglyUrls {
			x := replaceExtension(strings.TrimSpace(t), p.Extension)
			return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Section)+"/"+x))
		} else {
			file, _ := fileExt(strings.TrimSpace(t))
			return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Section)+"/"+file))
		}
	}
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
			page.Slug = Urlize(interfaceToString(v))
		case "url":
			if url := interfaceToString(v); strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				return fmt.Errorf("Only relative urls are supported, %v provided", url)
			}
			page.Url = Urlize(interfaceToString(v))
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

func (page *Page) detectFrontMatter(mark rune) (f *frontmatterType) {
	switch mark {
	case '-':
		return &frontmatterType{[]byte{'-', '-', '-'}, []byte{'-', '-', '-'}, page.handleYamlMetaData, false}
	case '+':
		return &frontmatterType{[]byte{'+', '+', '+'}, []byte{'+', '+', '+'}, page.handleTomlMetaData, false}
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

func (page *Page) readFile() (data []byte, err error) {
	data, err = ioutil.ReadFile(page.FileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (page *Page) buildPageFromFile() error {
	f, err := os.Open(page.FileName)
	if err != nil {
		return err
	}
	return page.parse(bufio.NewReader(f))
}

func (page *Page) parse(reader io.Reader) error {
	p, err := parser.ReadFrom(reader)
	if err != nil {
		return err
	}

	front := p.FrontMatter()
	if len(front) == 0 {
		return errors.New("Unable to locate frontmatter")
	}
	fm := page.detectFrontMatter(rune(front[0]))
	meta, err := fm.parse(front)
	if err != nil {
		return err
	}

	if err = page.update(meta); err != nil {
		return err
	}

	switch page.Markup {
	case "md":
		page.convertMarkdown(bytes.NewReader(p.Content()))
	case "rst":
		page.convertRestructuredText(bytes.NewReader(p.Content()))
	}
	return nil
}

func (page *Page) convertMarkdown(lines io.Reader) {
	b := new(bytes.Buffer)
	b.ReadFrom(lines)
	content := b.Bytes()
	page.Content = template.HTML(string(blackfriday.MarkdownCommon(content)))
	summary, plain := getSummaryString(content)
	if plain {
		page.Summary = template.HTML(string(summary))
	} else {
		page.Summary = template.HTML(string(blackfriday.MarkdownCommon(summary)))
	}
}

func (page *Page) convertRestructuredText(lines io.Reader) {
	b := new(bytes.Buffer)
	b.ReadFrom(lines)
	content := b.Bytes()
	page.Content = template.HTML(getRstContent(content))
	summary, plain := getSummaryString(content)
	if plain {
		page.Summary = template.HTML(string(summary))
	} else {
		page.Summary = template.HTML(getRstContent(summary))
	}
}
