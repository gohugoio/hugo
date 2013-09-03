package hugolib

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"html/template"
)

var TEMPLATE_TITLE = "{{ .Title }}"
var PAGE_SIMPLE_TITLE = `---
title: simple template
---
content`

var TEMPLATE_MISSING_FUNC = "{{ .Title | funcdoesnotexists }}"
var TEMPLATE_FUNC = "{{ .Title | urlize }}"
var TEMPLATE_CONTENT = "{{ .Content }}"
var TEMPLATE_DATE = "{{ .Date }}"
var INVALID_TEMPLATE_FORMAT_DATE = "{{ .Date.Format time.RFC3339 }}"

var PAGE_URL_SPECIFIED = `---
title: simple template
url: "mycategory/my-whatever-content/"
---
content`

var PAGE_WITH_MD = `---
title: page with md
---
# heading 1
text
## heading 2
more text
`

func pageMust(p *Page, err error) *Page {
	if err != nil {
		panic(err)
	}
	return p
}

func TestDegenerateRenderThingMissingTemplate(t *testing.T) {
	p, _ := ReadFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md")
	s := new(Site)
	s.prepTemplates()
	_, err := s.RenderThing(p, "foobar")
	if err == nil {
		t.Errorf("Expected err to be returned when missing the template.")
	}
}

func TestPrimeTemplates(t *testing.T) {
	s := new(Site)
	s.prepTemplates()
	if s.Tmpl.Lookup("alias") == nil {
		t.Fatalf("alias template not created.")
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

func matchRender(t *testing.T, s *Site, p *Page, tmplName string, expected string) {
	content, err := s.RenderThing(p, tmplName)
	if err != nil {
		t.Fatalf("Unable to render template.")
	}

	if string(content.Bytes()) != expected {
		t.Fatalf("Content did not match expected: %s. got: %s", expected, content)
	}
}

func _TestAddSameTemplateTwice(t *testing.T) {
	p := pageMust(ReadFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md"))
	s := new(Site)
	s.prepTemplates()
	err := s.addTemplate("foo", TEMPLATE_TITLE)
	if err != nil {
		t.Fatalf("Unable to add template foo")
	}

	matchRender(t, s, p, "foo", "simple template")

	err = s.addTemplate("foo", "NEW {{ .Title }}")
	if err != nil {
		t.Fatalf("Unable to add template foo: %s", err)
	}

	matchRender(t, s, p, "foo", "NEW simple template")
}

func TestRenderThing(t *testing.T) {
	tests := []struct {
		content  string
		template string
		expected string
	}{
		{PAGE_SIMPLE_TITLE, TEMPLATE_TITLE, "simple template"},
		{PAGE_SIMPLE_TITLE, TEMPLATE_FUNC, "simple-template"},
		{PAGE_WITH_MD, TEMPLATE_CONTENT, "<h1>heading 1</h1>\n\n<p>text</p>\n\n<h2>heading 2</h2>\n\n<p>more text</p>\n"},
		{SIMPLE_PAGE_RFC3339_DATE, TEMPLATE_DATE, "2013-05-17 16:59:30 &#43;0000 UTC"},
	}

	s := new(Site)
	s.prepTemplates()

	for i, test := range tests {
		p, err := ReadFrom(strings.NewReader(test.content), "content/a/file.md")
		if err != nil {
			t.Fatalf("Error parsing buffer: %s", err)
		}
		templateName := fmt.Sprintf("foobar%d", i)
		err = s.addTemplate(templateName, test.template)
		if err != nil {
			t.Fatalf("Unable to add template")
		}

		p.Content = template.HTML(p.Content)
		html, err2 := s.RenderThing(p, templateName)
		if err2 != nil {
			t.Errorf("Unable to render html: %s", err)
		}

		if string(html.Bytes()) != test.expected {
			t.Errorf("Content does not match.\nExpected\n\t'%q'\ngot\n\t'%q'", test.expected, html)
		}
	}
}

func TestRenderThingOrDefault(t *testing.T) {
	tests := []struct {
		content  string
		missing  bool
		template string
		expected string
	}{
		{PAGE_SIMPLE_TITLE, true, TEMPLATE_TITLE, "simple template"},
		{PAGE_SIMPLE_TITLE, true, TEMPLATE_FUNC, "simple-template"},
		{PAGE_SIMPLE_TITLE, false, TEMPLATE_TITLE, "simple template"},
		{PAGE_SIMPLE_TITLE, false, TEMPLATE_FUNC, "simple-template"},
	}

	s := new(Site)
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

		var html *bytes.Buffer
		var err2 error
		if test.missing {
			html, err2 = s.RenderThingOrDefault(p, "missing", templateName)
		} else {
			html, err2 = s.RenderThingOrDefault(p, templateName, "missing_default")
		}

		if err2 != nil {
			t.Errorf("Unable to render html: %s", err)
		}

		if string(html.Bytes()) != test.expected {
			t.Errorf("Content does not match.  Expected '%s', got '%s'", test.expected, html)
		}
	}
}

func TestSetOutFile(t *testing.T) {
	s := new(Site)
	p := pageMust(ReadFrom(strings.NewReader(PAGE_URL_SPECIFIED), "content/a/file.md"))
	s.setOutFile(p)

	expected := "mycategory/my-whatever-content/index.html"

	if p.OutFile != "mycategory/my-whatever-content/index.html" {
		t.Errorf("Outfile does not match.  Expected '%s', got '%s'", expected, p.OutFile)
	}
}
