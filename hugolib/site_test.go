package hugolib

import (
	"bytes"
	"fmt"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"html/template"
	"io"
	"strings"
	"testing"
)

const (
	TEMPLATE_TITLE    = "{{ .Title }}"
	PAGE_SIMPLE_TITLE = `---
title: simple template
---
content`

	TEMPLATE_MISSING_FUNC        = "{{ .Title | funcdoesnotexists }}"
	TEMPLATE_FUNC                = "{{ .Title | urlize }}"
	TEMPLATE_CONTENT             = "{{ .Content }}"
	TEMPLATE_DATE                = "{{ .Date }}"
	INVALID_TEMPLATE_FORMAT_DATE = "{{ .Date.Format time.RFC3339 }}"
	TEMPLATE_WITH_URL_REL        = "<a href=\"foobar.jpg\">Going</a>"
	TEMPLATE_WITH_URL_ABS        = "<a href=\"/foobar.jpg\">Going</a>"
	PAGE_URL_SPECIFIED           = `---
title: simple template
url: "mycategory/my-whatever-content/"
---
content`

	PAGE_WITH_MD = `---
title: page with md
---
# heading 1
text
## heading 2
more text
`
)

func pageMust(p *Page, err error) *Page {
	if err != nil {
		panic(err)
	}
	return p
}

func TestDegenerateRenderThingMissingTemplate(t *testing.T) {
	p, _ := ReadFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md")
	p.Convert()
	s := new(Site)
	s.prepTemplates()
	err := s.renderThing(p, "foobar", nil)
	if err == nil {
		t.Errorf("Expected err to be returned when missing the template.")
	}
}

func TestAddInvalidTemplate(t *testing.T) {
	s := new(Site)
	s.prepTemplates()
	err := s.addTemplate("missing", TEMPLATE_MISSING_FUNC)
	if err == nil {
		t.Fatalf("Expecting the template to return an error")
	}
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

func matchRender(t *testing.T, s *Site, p *Page, tmplName string, expected string) {
	content := new(bytes.Buffer)
	err := s.renderThing(p, tmplName, NopCloser(content))
	if err != nil {
		t.Fatalf("Unable to render template.")
	}

	if string(content.Bytes()) != expected {
		t.Fatalf("Content did not match expected: %s. got: %s", expected, content)
	}
}

func TestRenderThing(t *testing.T) {
	tests := []struct {
		content  string
		template string
		expected string
	}{
		{PAGE_SIMPLE_TITLE, TEMPLATE_TITLE, "simple template"},
		{PAGE_SIMPLE_TITLE, TEMPLATE_FUNC, "simple-template"},
		{PAGE_WITH_MD, TEMPLATE_CONTENT, "\n\n<h1 id=\"toc_0\">heading 1</h1>\n\n<p>text</p>\n\n<h2 id=\"toc_1\">heading 2</h2>\n\n<p>more text</p>\n"},
		{SIMPLE_PAGE_RFC3339_DATE, TEMPLATE_DATE, "2013-05-17 16:59:30 &#43;0000 UTC"},
	}

	s := new(Site)
	s.prepTemplates()

	for i, test := range tests {
		p, err := ReadFrom(strings.NewReader(test.content), "content/a/file.md")
		p.Convert()
		if err != nil {
			t.Fatalf("Error parsing buffer: %s", err)
		}
		templateName := fmt.Sprintf("foobar%d", i)
		err = s.addTemplate(templateName, test.template)
		if err != nil {
			t.Fatalf("Unable to add template")
		}

		p.Content = template.HTML(p.Content)
		html := new(bytes.Buffer)
		err = s.renderThing(p, templateName, NopCloser(html))
		if err != nil {
			t.Errorf("Unable to render html: %s", err)
		}

		if string(html.Bytes()) != test.expected {
			t.Errorf("Content does not match.\nExpected\n\t'%q'\ngot\n\t'%q'", test.expected, html)
		}
	}
}

func HTML(in string) string {
	return in
}

func TestRenderThingOrDefault(t *testing.T) {
	tests := []struct {
		content  string
		missing  bool
		template string
		expected string
	}{
		{PAGE_SIMPLE_TITLE, true, TEMPLATE_TITLE, HTML("simple template")},
		{PAGE_SIMPLE_TITLE, true, TEMPLATE_FUNC, HTML("simple-template")},
		{PAGE_SIMPLE_TITLE, false, TEMPLATE_TITLE, HTML("simple template")},
		{PAGE_SIMPLE_TITLE, false, TEMPLATE_FUNC, HTML("simple-template")},
	}

	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}
	s := &Site{
		Target: target,
	}
	s.prepTemplates()

	for i, test := range tests {
		p, err := ReadFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md")
		if err != nil {
			t.Fatalf("Error parsing buffer: %s", err)
		}
		templateName := fmt.Sprintf("default%d", i)
		err = s.addTemplate(templateName, test.template)
		if err != nil {
			t.Fatalf("Unable to add template")
		}

		var err2 error
		if test.missing {
			err2 = s.render(p, "out", "missing", templateName)
		} else {
			err2 = s.render(p, "out", templateName, "missing_default")
		}

		if err2 != nil {
			t.Errorf("Unable to render html: %s", err)
		}

		if string(files["out"]) != test.expected {
			t.Errorf("Content does not match. Expected '%s', got '%s'", test.expected, files["out"])
		}
	}
}

func TestTargetPath(t *testing.T) {
	tests := []struct {
		doc             string
		content         string
		expectedOutFile string
		expectedSection string
	}{
		{"content/a/file.md", PAGE_URL_SPECIFIED, "mycategory/my-whatever-content/index.html", "a"},
		{"content/x/y/deepfile.md", SIMPLE_PAGE, "x/y/deepfile.html", "x/y"},
		{"content/x/y/z/deeperfile.md", SIMPLE_PAGE, "x/y/z/deeperfile.html", "x/y/z"},
		{"content/b/file.md", SIMPLE_PAGE, "b/file.html", "b"},
		{"a/file.md", SIMPLE_PAGE, "a/file.html", "a"},
		{"file.md", SIMPLE_PAGE, "file.html", ""},
	}

	if true {
		return
	}
	for _, test := range tests {
		s := &Site{
			Config: Config{ContentDir: "content"},
		}
		p := pageMust(ReadFrom(strings.NewReader(test.content), s.Config.GetAbsPath(test.doc)))

		expected := test.expectedOutFile

		if p.TargetPath() != expected {
			t.Errorf("%s => OutFile  expected: '%s', got: '%s'", test.doc, expected, p.TargetPath())
		}

		if p.Section != test.expectedSection {
			t.Errorf("%s => p.Section expected: %s, got: %s", test.doc, test.expectedSection, p.Section)
		}
	}
}

func TestSkipRender(t *testing.T) {
	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}
	sources := []source.ByteSource{
		{"sect/doc1.html", []byte("---\nmarkup: markdown\n---\n# title\nsome *content*"), "sect"},
		{"sect/doc2.html", []byte("<!doctype html><html><body>more content</body></html>"), "sect"},
		{"sect/doc3.md", []byte("# doc3\n*some* content"), "sect"},
		{"sect/doc4.md", []byte("---\ntitle: doc4\n---\n# doc4\n*some content*"), "sect"},
		{"sect/doc5.html", []byte("<!doctype html><html>{{ template \"head\" }}<body>body5</body></html>"), "sect"},
		{"sect/doc6.html", []byte("<!doctype html><html>{{ template \"head_abs\" }}<body>body5</body></html>"), "sect"},
		{"doc7.html", []byte("<html><body>doc7 content</body></html>"), ""},
		{"sect/doc8.html", []byte("---\nmarkup: md\n---\n# title\nsome *content*"), "sect"},
	}

	s := &Site{
		Target: target,
		Config: Config{
			Verbose:      true,
			BaseUrl:      "http://auth/bub",
			CanonifyUrls: true,
		},
		Source: &source.InMemorySource{ByteSource: sources},
	}

	s.initializeSiteInfo()
	s.prepTemplates()

	must(s.addTemplate("_default/single.html", "{{.Content}}"))
	must(s.addTemplate("head", "<head><script src=\"script.js\"></script></head>"))
	must(s.addTemplate("head_abs", "<head><script src=\"/script.js\"></script></head>"))

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if err := s.RenderPages(); err != nil {
		t.Fatalf("Unable to render pages. %s", err)
	}

	tests := []struct {
		doc      string
		expected string
	}{
		{"sect/doc1.html", "\n\n<h1 id=\"toc_0\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{"sect/doc2.html", "<!doctype html><html><body>more content</body></html>"},
		{"sect/doc3.html", "\n\n<h1 id=\"toc_0\">doc3</h1>\n\n<p><em>some</em> content</p>\n"},
		{"sect/doc4.html", "\n\n<h1 id=\"toc_0\">doc4</h1>\n\n<p><em>some content</em></p>\n"},
		{"sect/doc5.html", "<!doctype html><html><head><script src=\"script.js\"></script></head><body>body5</body></html>"},
		{"sect/doc6.html", "<!doctype html><html><head><script src=\"http://auth/bub/script.js\"></script></head><body>body5</body></html>"},
		{"doc7.html", "<html><body>doc7 content</body></html>"},
		{"sect/doc8.html", "\n\n<h1 id=\"toc_0\">title</h1>\n\n<p>some <em>content</em></p>\n"},
	}

	for _, test := range tests {
		content, ok := target.Files[test.doc]
		if !ok {
			t.Fatalf("Did not find %s in target. %v", test.doc, target.Files)
		}

		if !bytes.Equal(content, []byte(test.expected)) {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, string(content))
		}
	}
}

func TestAbsUrlify(t *testing.T) {
	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}
	sources := []source.ByteSource{
		{"sect/doc1.html", []byte("<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>"), "sect"},
		{"content/blue/doc2.html", []byte("---\nf: t\n---\n<!doctype html><html><body>more content</body></html>"), "blue"},
	}
	for _, canonify := range []bool{true, false} {
		s := &Site{
			Target: target,
			Config: Config{
				BaseUrl:      "http://auth/bub",
				CanonifyUrls: canonify,
			},
			Source: &source.InMemorySource{ByteSource: sources},
		}
		t.Logf("Rendering with BaseUrl %q and CanonifyUrls set %v", s.Config.BaseUrl, canonify)
		s.initializeSiteInfo()
		s.prepTemplates()
		must(s.addTemplate("blue/single.html", TEMPLATE_WITH_URL_ABS))

		if err := s.CreatePages(); err != nil {
			t.Fatalf("Unable to create pages: %s", err)
		}

		if err := s.BuildSiteMeta(); err != nil {
			t.Fatalf("Unable to build site metadata: %s", err)
		}

		if err := s.RenderPages(); err != nil {
			t.Fatalf("Unable to render pages. %s", err)
		}

		tests := []struct {
			file, expected string
		}{
			{"content/blue/doc2.html", "<a href=\"http://auth/bub/foobar.jpg\">Going</a>"},
			{"sect/doc1.html", "<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>"},
		}

		for _, test := range tests {
			content, ok := target.Files[test.file]
			if !ok {
				t.Fatalf("Unable to locate rendered content: %s", test.file)
			}

			expected := test.expected
			if !canonify {
				expected = strings.Replace(expected, s.Config.BaseUrl, "", -1)
			}
			if string(content) != expected {
				t.Errorf("AbsUrlify content expected:\n%q\ngot\n%q", expected, string(content))
			}
		}
	}
}

var WEIGHTED_PAGE_1 = []byte(`+++
weight = "2"
title = "One"
+++
Front Matter with Ordered Pages`)

var WEIGHTED_PAGE_2 = []byte(`+++
weight = "6"
title = "Two"
+++
Front Matter with Ordered Pages 2`)

var WEIGHTED_PAGE_3 = []byte(`+++
weight = "4"
title = "Three"
date = "2012-04-06"
+++
Front Matter with Ordered Pages 3`)

var WEIGHTED_PAGE_4 = []byte(`+++
weight = "4"
title = "Four"
date = "2012-01-01"
+++
Front Matter with Ordered Pages 4. This is longer content`)

var WEIGHTED_SOURCES = []source.ByteSource{
	{"sect/doc1.md", WEIGHTED_PAGE_1, "sect"},
	{"sect/doc2.md", WEIGHTED_PAGE_2, "sect"},
	{"sect/doc3.md", WEIGHTED_PAGE_3, "sect"},
	{"sect/doc4.md", WEIGHTED_PAGE_4, "sect"},
}

func TestOrderedPages(t *testing.T) {
	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}
	s := &Site{
		Target: target,
		Config: Config{BaseUrl: "http://auth/bub/"},
		Source: &source.InMemorySource{ByteSource: WEIGHTED_SOURCES},
	}
	s.initializeSiteInfo()

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if s.Sections["sect"][0].Weight != 2 || s.Sections["sect"][3].Weight != 6 {
		t.Errorf("Pages in unexpected order. First should be '%d', got '%d'", 2, s.Sections["sect"][0].Weight)
	}

	if s.Sections["sect"][1].Page.Title != "Three" || s.Sections["sect"][2].Page.Title != "Four" {
		t.Errorf("Pages in unexpected order. Second should be '%s', got '%s'", "Three", s.Sections["sect"][1].Page.Title)
	}

	bydate := s.Pages.ByDate()

	if bydate[0].Title != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bydate[0].Title)
	}

	rev := bydate.Reverse()
	if rev[0].Title != "Three" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Three", rev[0].Title)
	}

	bylength := s.Pages.ByLength()
	if bylength[0].Title != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bylength[0].Title)
	}

	rbylength := bylength.Reverse()
	if rbylength[0].Title != "Four" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Four", rbylength[0].Title)
	}
}

var PAGE_WITH_WEIGHTED_INDEXES_2 = []byte(`+++
tags = [ "a", "b", "c" ]
tags_weight = 22
categories = ["d"]
title = "foo"
categories_weight = 44
+++
Front Matter with weighted tags and categories`)

var PAGE_WITH_WEIGHTED_INDEXES_1 = []byte(`+++
tags = [ "a" ]
tags_weight = 33
title = "bar"
categories = [ "d", "e" ]
categories_weight = 11
alias = "spf13"
date = 1979-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`)

var PAGE_WITH_WEIGHTED_INDEXES_3 = []byte(`+++
title = "bza"
categories = [ "e" ]
categories_weight = 11
alias = "spf13"
date = 2010-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`)

func TestWeightedIndexes(t *testing.T) {
	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}
	sources := []source.ByteSource{
		{"sect/doc1.md", PAGE_WITH_WEIGHTED_INDEXES_1, "sect"},
		{"sect/doc2.md", PAGE_WITH_WEIGHTED_INDEXES_2, "sect"},
		{"sect/doc3.md", PAGE_WITH_WEIGHTED_INDEXES_3, "sect"},
	}
	indexes := make(map[string]string)

	indexes["tag"] = "tags"
	indexes["category"] = "categories"
	s := &Site{
		Target: target,
		Config: Config{BaseUrl: "http://auth/bub/", Indexes: indexes},
		Source: &source.InMemorySource{ByteSource: sources},
	}
	s.initializeSiteInfo()

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if s.Indexes["tags"]["a"][0].Page.Title != "foo" {
		t.Errorf("Pages in unexpected order, 'foo' expected first, got '%v'", s.Indexes["tags"]["a"][0].Page.Title)
	}

	if s.Indexes["categories"]["d"][0].Page.Title != "bar" {
		t.Errorf("Pages in unexpected order, 'bar' expected first, got '%v'", s.Indexes["categories"]["d"][0].Page.Title)
	}

	if s.Indexes["categories"]["e"][0].Page.Title != "bza" {
		t.Errorf("Pages in unexpected order, 'bza' expected first, got '%v'", s.Indexes["categories"]["e"][0].Page.Title)
	}
}
