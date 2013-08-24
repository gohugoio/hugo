package hugolib

import (
	"fmt"
	"strings"
	"testing"
)

var TEMPLATE_TITLE = "{{ .Title }}"
var PAGE_SIMPLE_TITLE = `---
title: simple template
---
content`

var TEMPLATE_MISSING_FUNC = "{{ .Title | funcdoesnotexists }}"
var TEMPLATE_FUNC = "{{ .Title | urlize }}"

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

func TestPrimeTempaltes(t *testing.T) {
	s := new(Site)
	s.prepTemplates()
	s.primeTemplates()
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
	}

	s := new(Site)
	s.prepTemplates()

	for i, test := range tests {
		p, err := ReadFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md")
		if err != nil {
			t.Fatalf("Error parsing buffer: %s", err)
		}
		templateName := fmt.Sprintf("foobar%d", i)
		err = s.addTemplate(templateName, test.template)
		if err != nil {
			t.Fatalf("Unable to add template")
		}

		html, err := s.RenderThing(p, templateName)
		if err != nil {
			t.Errorf("Unable to render html: %s", err)
		}

		if string(html.Bytes()) != test.expected {
			t.Errorf("Content does not match.  Expected '%s', got '%s'", test.expected, html)
		}
	}
}
