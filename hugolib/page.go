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
	"encoding/json"
	"launchpad.net/goyaml"
	"fmt"
	"github.com/theplant/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var _ = filepath.Base("")

type Page struct {
	Status		string
	Images		[]string
	Content		template.HTML
	Summary		template.HTML
	RawMarkdown	string // TODO should be []byte
	Params		map[string]interface{}
	RenderedContent *bytes.Buffer
	contentType	string
	Draft		bool
	Tmpl		*template.Template
	Markup		string
	PageMeta
	File
	Position
	Node
}

const summaryLength = 70

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

func (p Pages) Len() int	   { return len(p) }
func (p Pages) Less(i, j int) bool { return p[i].Date.Unix() > p[j].Date.Unix() }
func (p Pages) Swap(i, j int)	   { p[i], p[j] = p[j], p[i] }

// TODO eliminate unnecessary things
func (p Pages) Sort()		  { sort.Sort(p) }
func (p Pages) Limit(n int) Pages { return p[0:n] }

func initializePage(filename string) (page Page) {
	page = Page{}
	page.Date, _ = time.Parse("20060102", "20080101")
	page.FileName = filename
	page.contentType = ""
	page.Extension = "html"
	page.Params = make(map[string]interface{})
	page.Keywords = make([]string, 10, 30)
	page.Markup = "md"
	page.setSection()

	return page
}

func (p *Page) setSection() {
	x := strings.Split(p.FileName, "/")

	if section := x[len(x)-2]; section != "content" {
		p.Section = section
	}
}

func (page *Page) Type() string {
	if page.contentType != "" {
		return page.contentType
	}

	if x := page.GetSection(); x != "" {
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

// TODO should return errors as well
// TODO new page should return just a page
// TODO initalize separately... load from reader (file, or []byte)
func NewPage(filename string) *Page {
	p := initializePage(filename)
	if err := p.buildPageFromFile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p.analyzePage()

	return &p
}

func (p *Page) analyzePage() {
	p.WordCount = TotalWords(p.RawMarkdown)
	p.FuzzyWordCount = int((p.WordCount+100)/100) * 100
}

// TODO //rewrite to use byte methods instead
func (page *Page) parseYamlMetaData(data []byte) ([]string, error) {
	var err error

	lines := strings.Split(string(data), "\n")
	datum := lines[0:]

	// go through content parse between "---" and "..."
	// must be on their own lines (for now)
	var found = 0
	for i, line := range lines {

		if strings.HasPrefix(line, "---") {
			found += 1
		}

		if strings.HasPrefix(line, "---") {
			found -= 1
		}

		if found == 0 {
			datum = lines[0: i+1]
			lines = lines[i+1:]
			break
		}
	}

	err = page.handleYamlMetaData([]byte(strings.Join(datum, "\n")))

	return lines, err
}

func (page *Page) parseJsonMetaData(data []byte) ([]string, error) {
	var err error

	lines := strings.Split(string(data), "\n")
	datum := lines[0:]

	// go through content parse between "{" and "}"
	// must be on their own lines (for now)
	var found = 0
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if line == "{" {
			found += 1
		}

		if line == "}" {
			found -= 1
		}

		if found == 0 {
			datum = lines[0 : i+1]
			lines = lines[i+1:]
			break
		}
	}

	err = page.handleJsonMetaData([]byte(strings.Join(datum, "\n")))

	return lines, err
}

func (p *Page) Permalink() template.HTML {
	if len(strings.TrimSpace(p.Slug)) > 0 {
		return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Section)+"/"+p.Slug))
	} else if len(strings.TrimSpace(p.Url)) > 2 {
		return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Url)))
	} else {
		_, t := filepath.Split(p.FileName)
		x := replaceExtension(strings.TrimSpace(t), p.Extension)
		return template.HTML(MakePermalink(string(p.Site.BaseUrl), strings.TrimSpace(p.Section)+"/"+x))
	}
}

func (page *Page) handleYamlMetaData(datum []byte) error {
	m := map[string]interface{}{}
	if err := goyaml.Unmarshal(datum, &m); err != nil {
		return fmt.Errorf("Invalid YAML in $v \nError parsing page meta data: %s", page.FileName, err)
	}

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
		case "status":
			page.Status = interfaceToString(v)
		default:
			// If not one of the explicit values, store in Params
			//fmt.Println(strings.ToLower(k))
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
	//Printer(page.Params)
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

func (page *Page) Err(message string) {
	fmt.Println(page.FileName + " : " + message)
}

// TODO return error on last line instead of nil
func (page *Page) parseFileHeading(data []byte) ([]string, error) {
	if len(data) == 0 {
		page.Err("Empty File, skipping")
	} else {
		if data[0] == '{' {
			return page.parseJsonMetaData(data)
		}
		return page.parseYamlMetaData(data)
	}
	return nil, nil
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

func (page *Page) readFile() []byte {
	var data, err = ioutil.ReadFile(page.FileName)
	if err != nil {
		PrintErr("Error Reading: " + page.FileName)
		return nil
	}
	return data
}

func (page *Page) buildPageFromFile() error {
	data := page.readFile()

	content, err := page.parseFileHeading(data)
	if err != nil {
		return err
	}

	if err := page.setOutFile(); err != nil {
		return err
	}

	switch page.Markup {
	case "md":
		page.convertMarkdown(content)
	case "rst":
		page.convertRestructuredText(content)
	}
	return nil
}

func (p *Page) setOutFile() error {
	if len(strings.TrimSpace(p.Slug)) > 0 {
		// Use Slug if provided
		p.OutFile = strings.TrimSpace(p.Slug + "." + p.Extension)
	} else if len(strings.TrimSpace(p.Url)) > 2 {
		// Use Url if provided & Slug missing
		p.OutFile = strings.TrimSpace(p.Url)
	} else {
		// Fall back to filename
		_, t := filepath.Split(p.FileName)
		p.OutFile = replaceExtension(strings.TrimSpace(t), p.Extension)
	}

	return nil
}

func (page *Page) convertMarkdown(lines []string) {

	page.RawMarkdown = strings.Join(lines, "\n")
	content := string(blackfriday.MarkdownCommon([]byte(page.RawMarkdown)))
	page.Content = template.HTML(content)
	page.Summary = template.HTML(TruncateWordsToWholeSentence(StripHTML(StripShortcodes(content)), summaryLength))
}

func (page *Page) convertRestructuredText(lines []string) {

	page.RawMarkdown = strings.Join(lines, " ")

	cmd := exec.Command("rst2html.py", "--template=/tmp/template.txt")
	cmd.Stdin = strings.NewReader(page.RawMarkdown)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		print(err)
	}

	content := out.String()
	page.Content = template.HTML(content)
	page.Summary = template.HTML(TruncateWordsToWholeSentence(StripHTML(StripShortcodes(content)), summaryLength))
}
