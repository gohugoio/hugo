// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"bitbucket.org/pkg/inflect"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

func createAndRenderPages(t *testing.T, s *Site) {
	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if err := s.RenderPages(); err != nil {
		t.Fatalf("Unable to render pages. %s", err)
	}
}

func templatePrep(s *Site) {
	s.Tmpl = tpl.New()
	s.Tmpl.LoadTemplates(s.absLayoutDir())
	if s.hasTheme() {
		s.Tmpl.LoadTemplatesWithPrefix(s.absThemeDir()+"/layouts", "theme")
	}
}

func pageMust(p *Page, err error) *Page {
	if err != nil {
		panic(err)
	}
	return p
}

func TestDegenerateRenderThingMissingTemplate(t *testing.T) {
	p, _ := NewPageFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md")
	p.Convert()
	s := new(Site)
	templatePrep(s)
	err := s.renderThing(p, "foobar", nil)
	if err == nil {
		t.Errorf("Expected err to be returned when missing the template.")
	}
}

func TestAddInvalidTemplate(t *testing.T) {
	s := new(Site)
	templatePrep(s)
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
		{PAGE_WITH_MD, TEMPLATE_CONTENT, "\n\n<h1 id=\"heading-1:91b5c4a22fc6103c73bb91e4a40568f8\">heading 1</h1>\n\n<p>text</p>\n\n<h2 id=\"heading-2:91b5c4a22fc6103c73bb91e4a40568f8\">heading 2</h2>\n\n<p>more text</p>\n"},
		{SIMPLE_PAGE_RFC3339_DATE, TEMPLATE_DATE, "2013-05-17 16:59:30 &#43;0000 UTC"},
	}

	for i, test := range tests {

		s := new(Site)
		templatePrep(s)

		p, err := NewPageFrom(strings.NewReader(test.content), "content/a/file.md")
		p.Convert()
		if err != nil {
			t.Fatalf("Error parsing buffer: %s", err)
		}
		templateName := fmt.Sprintf("foobar%d", i)
		err = s.addTemplate(templateName, test.template)
		if err != nil {
			t.Fatalf("Unable to add template: %s", err)
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

	hugofs.DestinationFS = new(afero.MemMapFs)

	for i, test := range tests {

		s := &Site{}
		templatePrep(s)

		p, err := NewPageFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md")
		if err != nil {
			t.Fatalf("Error parsing buffer: %s", err)
		}
		templateName := fmt.Sprintf("default%d", i)
		err = s.addTemplate(templateName, test.template)
		if err != nil {
			t.Fatalf("Unable to add template: %s", err)
		}

		var err2 error

		if test.missing {
			err2 = s.renderAndWritePage("name", "out", p, "missing", templateName)
		} else {
			err2 = s.renderAndWritePage("name", "out", p, templateName, "missing_default")
		}

		if err2 != nil {
			t.Errorf("Unable to render html: %s", err)
		}

		file, err := hugofs.DestinationFS.Open(filepath.FromSlash("out/index.html"))
		if err != nil {
			t.Errorf("Unable to open html: %s", err)
		}
		if helpers.ReaderToString(file) != test.expected {
			t.Errorf("Content does not match. Expected '%s', got '%s'", test.expected, helpers.ReaderToString(file))
		}
	}
}

func TestDraftAndFutureRender(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.DestinationFS = new(afero.MemMapFs)
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.md"), []byte("---\ntitle: doc1\ndraft: true\npublishdate: \"2414-05-29\"\n---\n# doc1\n*some content*")},
		{filepath.FromSlash("sect/doc2.md"), []byte("---\ntitle: doc2\ndraft: true\npublishdate: \"2012-05-29\"\n---\n# doc2\n*some content*")},
		{filepath.FromSlash("sect/doc3.md"), []byte("---\ntitle: doc3\ndraft: false\npublishdate: \"2414-05-29\"\n---\n# doc3\n*some content*")},
		{filepath.FromSlash("sect/doc4.md"), []byte("---\ntitle: doc4\ndraft: false\npublishdate: \"2012-05-29\"\n---\n# doc4\n*some content*")},
	}

	siteSetup := func() *Site {
		s := &Site{
			Source: &source.InMemorySource{ByteSource: sources},
		}

		s.initializeSiteInfo()

		if err := s.CreatePages(); err != nil {
			t.Fatalf("Unable to create pages: %s", err)
		}
		return s
	}

	viper.Set("baseurl", "http://auth/bub")

	// Testing Defaults.. Only draft:true and publishDate in the past should be rendered
	s := siteSetup()
	if len(s.Pages) != 1 {
		t.Fatal("Draft or Future dated content published unexpectedly")
	}

	// only publishDate in the past should be rendered
	viper.Set("BuildDrafts", true)
	s = siteSetup()
	if len(s.Pages) != 2 {
		t.Fatal("Future Dated Posts published unexpectedly")
	}

	//  drafts should not be rendered, but all dates should
	viper.Set("BuildDrafts", false)
	viper.Set("BuildFuture", true)
	s = siteSetup()
	if len(s.Pages) != 2 {
		t.Fatal("Draft posts published unexpectedly")
	}

	// all 4 should be included
	viper.Set("BuildDrafts", true)
	viper.Set("BuildFuture", true)
	s = siteSetup()
	if len(s.Pages) != 4 {
		t.Fatal("Drafts or Future posts not included as expected")
	}

	//setting defaults back
	viper.Set("BuildDrafts", false)
	viper.Set("BuildFuture", false)
}

// Issue #957
func TestCrossrefs(t *testing.T) {
	hugofs.DestinationFS = new(afero.MemMapFs)
	for _, uglyURLs := range []bool{true, false} {
		for _, relative := range []bool{true, false} {
			doTestCrossrefs(t, relative, uglyURLs)
		}
	}
}

func doTestCrossrefs(t *testing.T, relative, uglyURLs bool) {
	viper.Reset()
	defer viper.Reset()

	baseURL := "http://foo/bar"
	viper.Set("DefaultExtension", "html")
	viper.Set("baseurl", baseURL)
	viper.Set("UglyURLs", uglyURLs)
	viper.Set("verbose", true)

	var refShortcode string
	var expectedBase string
	var expectedURLSuffix string
	var expectedPathSuffix string

	if relative {
		refShortcode = "relref"
		expectedBase = "/bar"
	} else {
		refShortcode = "ref"
		expectedBase = baseURL
	}

	if uglyURLs {
		expectedURLSuffix = ".html"
		expectedPathSuffix = ".html"
	} else {
		expectedURLSuffix = "/"
		expectedPathSuffix = "/index.html"
	}

	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.md"),
			[]byte(fmt.Sprintf(`Ref 2: {{< %s "sect/doc2.md" >}}`, refShortcode))},
		// Issue #1148: Make sure that no P-tags is added around shortcodes.
		{filepath.FromSlash("sect/doc2.md"),
			[]byte(fmt.Sprintf(`**Ref 1:**

{{< %s "sect/doc1.md" >}}

THE END.`, refShortcode))},
		// Issue #1753: Should not add a trailing newline after shortcode.
		{filepath.FromSlash("sect/doc3.md"),
			[]byte(fmt.Sprintf(`**Ref 1:**{{< %s "sect/doc3.md" >}}.`, refShortcode))},
	}

	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		Targets: targetList{Page: &target.PagePub{UglyURLs: uglyURLs}},
	}

	s.initializeSiteInfo()
	templatePrep(s)

	must(s.addTemplate("_default/single.html", "{{.Content}}"))

	createAndRenderPages(t, s)

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash(fmt.Sprintf("sect/doc1%s", expectedPathSuffix)), fmt.Sprintf("<p>Ref 2: %s/sect/doc2%s</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("sect/doc2%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong></p>\n\n%s/sect/doc1%s\n\n<p>THE END.</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("sect/doc3%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong>%s/sect/doc3%s.</p>\n", expectedBase, expectedURLSuffix)},
	}

	for _, test := range tests {
		file, err := hugofs.DestinationFS.Open(test.doc)

		if err != nil {
			t.Fatalf("Did not find %s in target: %s", test.doc, err)
		}

		content := helpers.ReaderToString(file)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}

// Issue #939
func Test404ShouldAlwaysHaveUglyURLs(t *testing.T) {
	hugofs.DestinationFS = new(afero.MemMapFs)
	for _, uglyURLs := range []bool{true, false} {
		doTest404ShouldAlwaysHaveUglyURLs(t, uglyURLs)
	}
}

func doTest404ShouldAlwaysHaveUglyURLs(t *testing.T, uglyURLs bool) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("DefaultExtension", "html")
	viper.Set("verbose", true)
	viper.Set("baseurl", "http://auth/bub")
	viper.Set("DisableSitemap", false)
	viper.Set("DisableRSS", false)
	viper.Set("RSSUri", "index.xml")

	viper.Set("UglyURLs", uglyURLs)

	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
	}

	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		Targets: targetList{Page: &target.PagePub{UglyURLs: uglyURLs}},
	}

	s.initializeSiteInfo()
	templatePrep(s)

	must(s.addTemplate("index.html", "Home Sweet Home. IsHome={{ .IsHome  }}"))
	must(s.addTemplate("_default/single.html", "{{.Content}} IsHome={{ .IsHome  }}"))
	must(s.addTemplate("404.html", "Page Not Found. IsHome={{ .IsHome  }}"))

	// make sure the XML files also end up with ugly urls
	must(s.addTemplate("rss.xml", "<root>RSS</root>"))
	must(s.addTemplate("sitemap.xml", "<root>SITEMAP</root>"))

	createAndRenderPages(t, s)
	s.RenderHomePage()
	s.RenderSitemap()

	var expectedPagePath string
	if uglyURLs {
		expectedPagePath = "sect/doc1.html"
	} else {
		expectedPagePath = "sect/doc1/index.html"
	}

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("index.html"), "Home Sweet Home. IsHome=true"},
		{filepath.FromSlash(expectedPagePath), "\n\n<h1 id=\"title:5d74edbb89ef198cd37882b687940cda\">title</h1>\n\n<p>some <em>content</em></p>\n IsHome=false"},
		{filepath.FromSlash("404.html"), "Page Not Found. IsHome=false"},
		{filepath.FromSlash("index.xml"), "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n<root>RSS</root>"},
		{filepath.FromSlash("sitemap.xml"), "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n<root>SITEMAP</root>"},
	}

	for _, p := range s.Pages {
		assert.False(t, p.IsHome)
	}

	for _, test := range tests {
		file, err := hugofs.DestinationFS.Open(test.doc)
		if err != nil {
			t.Fatalf("Did not find %s in target: %s", test.doc, err)
		}

		content := helpers.ReaderToString(file)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}

// Issue #1176
func TestSectionNaming(t *testing.T) {

	for _, canonify := range []bool{true, false} {
		for _, uglify := range []bool{true, false} {
			for _, pluralize := range []bool{true, false} {
				doTestSectionNaming(t, canonify, uglify, pluralize)
			}
		}
	}
}

func doTestSectionNaming(t *testing.T, canonify, uglify, pluralize bool) {
	hugofs.DestinationFS = new(afero.MemMapFs)
	viper.Reset()
	defer viper.Reset()
	viper.Set("baseurl", "http://auth/sub/")
	viper.Set("DefaultExtension", "html")
	viper.Set("UglyURLs", uglify)
	viper.Set("PluralizeListTitles", pluralize)
	viper.Set("CanonifyURLs", canonify)

	var expectedPathSuffix string

	if uglify {
		expectedPathSuffix = ".html"
	} else {
		expectedPathSuffix = "/index.html"
	}

	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("doc1")},
		{filepath.FromSlash("Fish and Chips/doc2.html"), []byte("doc2")},
		{filepath.FromSlash("ラーメン/doc3.html"), []byte("doc3")},
	}

	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		Targets: targetList{Page: &target.PagePub{UglyURLs: uglify}},
	}

	s.initializeSiteInfo()
	templatePrep(s)

	must(s.addTemplate("_default/single.html", "{{.Content}}"))
	must(s.addTemplate("_default/list.html", "{{ .Title }}"))

	createAndRenderPages(t, s)
	s.RenderSectionLists()

	tests := []struct {
		doc         string
		pluralAware bool
		expected    string
	}{
		{filepath.FromSlash(fmt.Sprintf("sect/doc1%s", expectedPathSuffix)), false, "doc1"},
		{filepath.FromSlash(fmt.Sprintf("sect%s", expectedPathSuffix)), true, "Sect"},
		{filepath.FromSlash(fmt.Sprintf("fish-and-chips/doc2%s", expectedPathSuffix)), false, "doc2"},
		{filepath.FromSlash(fmt.Sprintf("fish-and-chips%s", expectedPathSuffix)), true, "Fish and Chips"},
		{filepath.FromSlash(fmt.Sprintf("ラーメン/doc3%s", expectedPathSuffix)), false, "doc3"},
		{filepath.FromSlash(fmt.Sprintf("ラーメン%s", expectedPathSuffix)), true, "ラーメン"},
	}

	for _, test := range tests {
		file, err := hugofs.DestinationFS.Open(test.doc)
		if err != nil {
			t.Fatalf("Did not find %s in target: %s", test.doc, err)
		}

		content := helpers.ReaderToString(file)

		if test.pluralAware && pluralize {
			test.expected = inflect.Pluralize(test.expected)
		}

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}
func TestSkipRender(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.DestinationFS = new(afero.MemMapFs)
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
		{filepath.FromSlash("sect/doc2.html"), []byte("<!doctype html><html><body>more content</body></html>")},
		{filepath.FromSlash("sect/doc3.md"), []byte("# doc3\n*some* content")},
		{filepath.FromSlash("sect/doc4.md"), []byte("---\ntitle: doc4\n---\n# doc4\n*some content*")},
		{filepath.FromSlash("sect/doc5.html"), []byte("<!doctype html><html>{{ template \"head\" }}<body>body5</body></html>")},
		{filepath.FromSlash("sect/doc6.html"), []byte("<!doctype html><html>{{ template \"head_abs\" }}<body>body5</body></html>")},
		{filepath.FromSlash("doc7.html"), []byte("<html><body>doc7 content</body></html>")},
		{filepath.FromSlash("sect/doc8.html"), []byte("---\nmarkup: md\n---\n# title\nsome *content*")},
	}

	viper.Set("DefaultExtension", "html")
	viper.Set("verbose", true)
	viper.Set("CanonifyURLs", true)
	viper.Set("baseurl", "http://auth/bub")
	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		Targets: targetList{Page: &target.PagePub{UglyURLs: true}},
	}

	s.initializeSiteInfo()
	templatePrep(s)

	must(s.addTemplate("_default/single.html", "{{.Content}}"))
	must(s.addTemplate("head", "<head><script src=\"script.js\"></script></head>"))
	must(s.addTemplate("head_abs", "<head><script src=\"/script.js\"></script></head>"))

	createAndRenderPages(t, s)

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("sect/doc1.html"), "\n\n<h1 id=\"title:5d74edbb89ef198cd37882b687940cda\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("sect/doc2.html"), "<!doctype html><html><body>more content</body></html>"},
		{filepath.FromSlash("sect/doc3.html"), "\n\n<h1 id=\"doc3:28c75a9e2162b8eccda73a1ab9ce80b4\">doc3</h1>\n\n<p><em>some</em> content</p>\n"},
		{filepath.FromSlash("sect/doc4.html"), "\n\n<h1 id=\"doc4:f8e6806123f341b8975509637645a4d3\">doc4</h1>\n\n<p><em>some content</em></p>\n"},
		{filepath.FromSlash("sect/doc5.html"), "<!doctype html><html><head><script src=\"script.js\"></script></head><body>body5</body></html>"},
		{filepath.FromSlash("sect/doc6.html"), "<!doctype html><html><head><script src=\"http://auth/bub/script.js\"></script></head><body>body5</body></html>"},
		{filepath.FromSlash("doc7.html"), "<html><body>doc7 content</body></html>"},
		{filepath.FromSlash("sect/doc8.html"), "\n\n<h1 id=\"title:0ae308ad73e2f37bd09874105281b5d8\">title</h1>\n\n<p>some <em>content</em></p>\n"},
	}

	for _, test := range tests {
		file, err := hugofs.DestinationFS.Open(test.doc)
		if err != nil {
			t.Fatalf("Did not find %s in target.", test.doc)
		}

		content := helpers.ReaderToString(file)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}
}

func TestAbsURLify(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("DefaultExtension", "html")

	hugofs.DestinationFS = new(afero.MemMapFs)
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>")},
		{filepath.FromSlash("content/blue/doc2.html"), []byte("---\nf: t\n---\n<!doctype html><html><body>more content</body></html>")},
	}
	for _, canonify := range []bool{true, false} {
		viper.Set("CanonifyURLs", canonify)
		viper.Set("BaseURL", "http://auth/bub")
		s := &Site{
			Source:  &source.InMemorySource{ByteSource: sources},
			Targets: targetList{Page: &target.PagePub{UglyURLs: true}},
		}
		t.Logf("Rendering with BaseURL %q and CanonifyURLs set %v", viper.GetString("baseURL"), canonify)
		s.initializeSiteInfo()
		templatePrep(s)
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

			file, err := hugofs.DestinationFS.Open(filepath.FromSlash(test.file))
			if err != nil {
				t.Fatalf("Unable to locate rendered content: %s", test.file)
			}

			content := helpers.ReaderToString(file)

			expected := test.expected

			if !canonify {
				expected = strings.Replace(expected, viper.GetString("baseurl"), "", -1)
			}

			if content != expected {
				t.Errorf("AbsURLify content expected:\n%q\ngot\n%q", expected, content)
			}
		}
	}
}

var WEIGHTED_PAGE_1 = []byte(`+++
weight = "2"
title = "One"
my_param = "foo"
my_date = 1979-05-27T07:32:00Z
+++
Front Matter with Ordered Pages`)

var WEIGHTED_PAGE_2 = []byte(`+++
weight = "6"
title = "Two"
publishdate = "2012-03-05"
my_param = "foo"
+++
Front Matter with Ordered Pages 2`)

var WEIGHTED_PAGE_3 = []byte(`+++
weight = "4"
title = "Three"
date = "2012-04-06"
publishdate = "2012-04-06"
my_param = "bar"
only_one = "yes"
my_date = 2010-05-27T07:32:00Z
+++
Front Matter with Ordered Pages 3`)

var WEIGHTED_PAGE_4 = []byte(`+++
weight = "4"
title = "Four"
date = "2012-01-01"
publishdate = "2012-01-01"
my_param = "baz"
my_date = 2010-05-27T07:32:00Z
+++
Front Matter with Ordered Pages 4. This is longer content`)

var WEIGHTED_SOURCES = []source.ByteSource{
	{filepath.FromSlash("sect/doc1.md"), WEIGHTED_PAGE_1},
	{filepath.FromSlash("sect/doc2.md"), WEIGHTED_PAGE_2},
	{filepath.FromSlash("sect/doc3.md"), WEIGHTED_PAGE_3},
	{filepath.FromSlash("sect/doc4.md"), WEIGHTED_PAGE_4},
}

func TestOrderedPages(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.DestinationFS = new(afero.MemMapFs)

	viper.Set("baseurl", "http://auth/bub")
	s := &Site{
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

	bypubdate := s.Pages.ByPublishDate()

	if bypubdate[0].Title != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bypubdate[0].Title)
	}

	rbypubdate := bypubdate.Reverse()
	if rbypubdate[0].Title != "Three" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Three", rbypubdate[0].Title)
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

var GROUPED_SOURCES = []source.ByteSource{
	{filepath.FromSlash("sect1/doc1.md"), WEIGHTED_PAGE_1},
	{filepath.FromSlash("sect1/doc2.md"), WEIGHTED_PAGE_2},
	{filepath.FromSlash("sect2/doc3.md"), WEIGHTED_PAGE_3},
	{filepath.FromSlash("sect3/doc4.md"), WEIGHTED_PAGE_4},
}

func TestGroupedPages(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	hugofs.DestinationFS = new(afero.MemMapFs)

	viper.Set("baseurl", "http://auth/bub")
	s := &Site{
		Source: &source.InMemorySource{ByteSource: GROUPED_SOURCES},
	}
	s.initializeSiteInfo()

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	rbysection, err := s.Pages.GroupBy("Section", "desc")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if rbysection[0].Key != "sect3" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "sect3", rbysection[0].Key)
	}
	if rbysection[1].Key != "sect2" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "sect2", rbysection[1].Key)
	}
	if rbysection[2].Key != "sect1" {
		t.Errorf("PageGroup array in unexpected order. Third group key should be '%s', got '%s'", "sect1", rbysection[2].Key)
	}
	if rbysection[0].Pages[0].Title != "Four" {
		t.Errorf("PageGroup has an unexpected page. First group's pages should have '%s', got '%s'", "Four", rbysection[0].Pages[0].Title)
	}
	if len(rbysection[2].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. Third group should have '%d' pages, got '%d' pages", 2, len(rbysection[2].Pages))
	}

	bytype, err := s.Pages.GroupBy("Type", "asc")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if bytype[0].Key != "sect1" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "sect1", bytype[0].Key)
	}
	if bytype[1].Key != "sect2" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "sect2", bytype[1].Key)
	}
	if bytype[2].Key != "sect3" {
		t.Errorf("PageGroup array in unexpected order. Third group key should be '%s', got '%s'", "sect3", bytype[2].Key)
	}
	if bytype[2].Pages[0].Title != "Four" {
		t.Errorf("PageGroup has an unexpected page. Third group's data should have '%s', got '%s'", "Four", bytype[0].Pages[0].Title)
	}
	if len(bytype[0].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 2, len(bytype[2].Pages))
	}

	bydate, err := s.Pages.GroupByDate("2006-01", "asc")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if bydate[0].Key != "0001-01" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "0001-01", bydate[0].Key)
	}
	if bydate[1].Key != "2012-01" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "2012-01", bydate[1].Key)
	}
	if bydate[2].Key != "2012-04" {
		t.Errorf("PageGroup array in unexpected order. Third group key should be '%s', got '%s'", "2012-04", bydate[2].Key)
	}
	if bydate[2].Pages[0].Title != "Three" {
		t.Errorf("PageGroup has an unexpected page. Third group's pages should have '%s', got '%s'", "Three", bydate[2].Pages[0].Title)
	}
	if len(bydate[0].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 2, len(bydate[2].Pages))
	}

	bypubdate, err := s.Pages.GroupByPublishDate("2006")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if bypubdate[0].Key != "2012" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "2012", bypubdate[0].Key)
	}
	if bypubdate[1].Key != "0001" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "0001", bypubdate[1].Key)
	}
	if bypubdate[0].Pages[0].Title != "Three" {
		t.Errorf("PageGroup has an unexpected page. Third group's pages should have '%s', got '%s'", "Three", bypubdate[0].Pages[0].Title)
	}
	if len(bypubdate[0].Pages) != 3 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 3, len(bypubdate[0].Pages))
	}

	byparam, err := s.Pages.GroupByParam("my_param", "desc")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if byparam[0].Key != "foo" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "foo", byparam[0].Key)
	}
	if byparam[1].Key != "baz" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "baz", byparam[1].Key)
	}
	if byparam[2].Key != "bar" {
		t.Errorf("PageGroup array in unexpected order. Third group key should be '%s', got '%s'", "bar", byparam[2].Key)
	}
	if byparam[2].Pages[0].Title != "Three" {
		t.Errorf("PageGroup has an unexpected page. Third group's pages should have '%s', got '%s'", "Three", byparam[2].Pages[0].Title)
	}
	if len(byparam[0].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 2, len(byparam[0].Pages))
	}

	_, err = s.Pages.GroupByParam("not_exist")
	if err == nil {
		t.Errorf("GroupByParam didn't return an expected error")
	}

	byOnlyOneParam, err := s.Pages.GroupByParam("only_one")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if len(byOnlyOneParam) != 1 {
		t.Errorf("PageGroup array has unexpected elements. Group length should be '%d', got '%d'", 1, len(byOnlyOneParam))
	}
	if byOnlyOneParam[0].Key != "yes" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "yes", byOnlyOneParam[0].Key)
	}

	byParamDate, err := s.Pages.GroupByParamDate("my_date", "2006-01")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if byParamDate[0].Key != "2010-05" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "2010-05", byParamDate[0].Key)
	}
	if byParamDate[1].Key != "1979-05" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "1979-05", byParamDate[1].Key)
	}
	if byParamDate[1].Pages[0].Title != "One" {
		t.Errorf("PageGroup has an unexpected page. Second group's pages should have '%s', got '%s'", "One", byParamDate[1].Pages[0].Title)
	}
	if len(byParamDate[0].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 2, len(byParamDate[2].Pages))
	}
}

var PAGE_WITH_WEIGHTED_TAXONOMIES_2 = []byte(`+++
tags = [ "a", "b", "c" ]
tags_weight = 22
categories = ["d"]
title = "foo"
categories_weight = 44
+++
Front Matter with weighted tags and categories`)

var PAGE_WITH_WEIGHTED_TAXONOMIES_1 = []byte(`+++
tags = "a"
tags_weight = 33
title = "bar"
categories = [ "d", "e" ]
categories_weight = 11
alias = "spf13"
date = 1979-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`)

var PAGE_WITH_WEIGHTED_TAXONOMIES_3 = []byte(`+++
title = "bza"
categories = [ "e" ]
categories_weight = 11
alias = "spf13"
date = 2010-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`)

func TestWeightedTaxonomies(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.DestinationFS = new(afero.MemMapFs)
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.md"), PAGE_WITH_WEIGHTED_TAXONOMIES_1},
		{filepath.FromSlash("sect/doc2.md"), PAGE_WITH_WEIGHTED_TAXONOMIES_2},
		{filepath.FromSlash("sect/doc3.md"), PAGE_WITH_WEIGHTED_TAXONOMIES_3},
	}
	taxonomies := make(map[string]string)

	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	viper.Set("baseurl", "http://auth/bub")
	viper.Set("taxonomies", taxonomies)
	s := &Site{
		Source: &source.InMemorySource{ByteSource: sources},
	}
	s.initializeSiteInfo()

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if s.Taxonomies["tags"]["a"][0].Page.Title != "foo" {
		t.Errorf("Pages in unexpected order, 'foo' expected first, got '%v'", s.Taxonomies["tags"]["a"][0].Page.Title)
	}

	if s.Taxonomies["categories"]["d"][0].Page.Title != "bar" {
		t.Errorf("Pages in unexpected order, 'bar' expected first, got '%v'", s.Taxonomies["categories"]["d"][0].Page.Title)
	}

	if s.Taxonomies["categories"]["e"][0].Page.Title != "bza" {
		t.Errorf("Pages in unexpected order, 'bza' expected first, got '%v'", s.Taxonomies["categories"]["e"][0].Page.Title)
	}
}

func findPage(site *Site, f string) *Page {
	// TODO: it seems that filepath.FromSlash results in page.Source.Path() returning windows backslash - which means refLinking's string compare is totally busted.
	// TODO: Not used for non-fragment linking (SVEN thinks this is a bug)
	currentPath := source.NewFile(filepath.FromSlash(f))
	//t.Logf("looking for currentPath: %s", currentPath.Path())

	for _, page := range site.Pages {
		//t.Logf("page: %s", page.Source.Path())
		if page.Source.Path() == currentPath.Path() {
			return page
		}
	}
	return nil
}

func setupLinkingMockSite(t *testing.T) *Site {
	hugofs.DestinationFS = new(afero.MemMapFs)
	sources := []source.ByteSource{
		{filepath.FromSlash("index.md"), []byte("")},
		{filepath.FromSlash("rootfile.md"), []byte("")},
		{filepath.FromSlash("root-image.png"), []byte("")},

		{filepath.FromSlash("level2/2-root.md"), []byte("")},
		{filepath.FromSlash("level2/index.md"), []byte("")},
		{filepath.FromSlash("level2/common.md"), []byte("")},

		//		{filepath.FromSlash("level2b/2b-root.md"), []byte("")},
		//		{filepath.FromSlash("level2b/index.md"), []byte("")},
		//		{filepath.FromSlash("level2b/common.md"), []byte("")},

		{filepath.FromSlash("level2/2-image.png"), []byte("")},
		{filepath.FromSlash("level2/common.png"), []byte("")},

		{filepath.FromSlash("level2/level3/3-root.md"), []byte("")},
		{filepath.FromSlash("level2/level3/index.md"), []byte("")},
		{filepath.FromSlash("level2/level3/common.md"), []byte("")},
		{filepath.FromSlash("level2/level3/3-image.png"), []byte("")},
		{filepath.FromSlash("level2/level3/common.png"), []byte("")},
	}

	site := &Site{
		Source: &source.InMemorySource{ByteSource: sources},
	}

	site.initializeSiteInfo()

	if err := site.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	viper.Set("baseurl", "http://auth/bub")
	viper.Set("DefaultExtension", "html")
	viper.Set("UglyURLs", false)
	viper.Set("PluralizeListTitles", false)
	viper.Set("CanonifyURLs", false)

	return site
}

func TestRefLinking(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	site := setupLinkingMockSite(t)

	currentPage := findPage(site, "level2/level3/index.md")
	if currentPage == nil {
		t.Fatalf("failed to find current page in site")
	}

	// refLink doesn't use the location of the current page to work out reflinks
	okresults := map[string]string{
		"index.md":  "/",
		"common.md": "/level2/common/",
		"3-root.md": "/level2/level3/3-root/",
	}
	for link, url := range okresults {
		if out, err := site.Info.refLink(link, currentPage, true); err != nil || out != url {
			t.Errorf("Expected %s to resolve to (%s), got (%s) - error: %s", link, url, out, err)
		}
	}
	// TODO: and then the failure cases.
}

func TestSourceRelativeLinksing(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	site := setupLinkingMockSite(t)

	type resultMap map[string]string

	okresults := map[string]resultMap{
		"index.md": map[string]string{
			"/docs/rootfile.md":             "/rootfile/",
			"/docs/index.md":                "/",
			"rootfile.md":                   "/rootfile/",
			"index.md":                      "/",
			"level2/2-root.md":              "/level2/2-root/",
			"level2/index.md":               "/level2/",
			"/docs/level2/2-root.md":        "/level2/2-root/",
			"/docs/level2/index.md":         "/level2/",
			"level2/level3/3-root.md":       "/level2/level3/3-root/",
			"level2/level3/index.md":        "/level2/level3/",
			"/docs/level2/level3/3-root.md": "/level2/level3/3-root/",
			"/docs/level2/level3/index.md":  "/level2/level3/",
			"/docs/level2/2-root/":          "/level2/2-root/",
			"/docs/level2/":                 "/level2/",
			"/docs/level2/2-root":           "/level2/2-root/",
			"/docs/level2":                  "/level2/",
			"/level2/2-root/":               "/level2/2-root/",
			"/level2/":                      "/level2/",
			"/level2/2-root":                "/level2/2-root/",
			"/level2":                       "/level2/",
		}, "rootfile.md": map[string]string{
			"/docs/rootfile.md":             "/rootfile/",
			"/docs/index.md":                "/",
			"rootfile.md":                   "/rootfile/",
			"index.md":                      "/",
			"level2/2-root.md":              "/level2/2-root/",
			"level2/index.md":               "/level2/",
			"/docs/level2/2-root.md":        "/level2/2-root/",
			"/docs/level2/index.md":         "/level2/",
			"level2/level3/3-root.md":       "/level2/level3/3-root/",
			"level2/level3/index.md":        "/level2/level3/",
			"/docs/level2/level3/3-root.md": "/level2/level3/3-root/",
			"/docs/level2/level3/index.md":  "/level2/level3/",
		}, "level2/2-root.md": map[string]string{
			"../rootfile.md":                "/rootfile/",
			"../index.md":                   "/",
			"/docs/rootfile.md":             "/rootfile/",
			"/docs/index.md":                "/",
			"2-root.md":                     "/level2/2-root/",
			"index.md":                      "/level2/",
			"../level2/2-root.md":           "/level2/2-root/",
			"../level2/index.md":            "/level2/",
			"./2-root.md":                   "/level2/2-root/",
			"./index.md":                    "/level2/",
			"/docs/level2/index.md":         "/level2/",
			"/docs/level2/2-root.md":        "/level2/2-root/",
			"level3/3-root.md":              "/level2/level3/3-root/",
			"level3/index.md":               "/level2/level3/",
			"../level2/level3/index.md":     "/level2/level3/",
			"../level2/level3/3-root.md":    "/level2/level3/3-root/",
			"/docs/level2/level3/index.md":  "/level2/level3/",
			"/docs/level2/level3/3-root.md": "/level2/level3/3-root/",
		}, "level2/index.md": map[string]string{
			"../rootfile.md":                "/rootfile/",
			"../index.md":                   "/",
			"/docs/rootfile.md":             "/rootfile/",
			"/docs/index.md":                "/",
			"2-root.md":                     "/level2/2-root/",
			"index.md":                      "/level2/",
			"../level2/2-root.md":           "/level2/2-root/",
			"../level2/index.md":            "/level2/",
			"./2-root.md":                   "/level2/2-root/",
			"./index.md":                    "/level2/",
			"/docs/level2/index.md":         "/level2/",
			"/docs/level2/2-root.md":        "/level2/2-root/",
			"level3/3-root.md":              "/level2/level3/3-root/",
			"level3/index.md":               "/level2/level3/",
			"../level2/level3/index.md":     "/level2/level3/",
			"../level2/level3/3-root.md":    "/level2/level3/3-root/",
			"/docs/level2/level3/index.md":  "/level2/level3/",
			"/docs/level2/level3/3-root.md": "/level2/level3/3-root/",
		}, "level2/level3/3-root.md": map[string]string{
			"../../rootfile.md":      "/rootfile/",
			"../../index.md":         "/",
			"/docs/rootfile.md":      "/rootfile/",
			"/docs/index.md":         "/",
			"../2-root.md":           "/level2/2-root/",
			"../index.md":            "/level2/",
			"/docs/level2/2-root.md": "/level2/2-root/",
			"/docs/level2/index.md":  "/level2/",
			"3-root.md":              "/level2/level3/3-root/",
			"index.md":               "/level2/level3/",
			"./3-root.md":            "/level2/level3/3-root/",
			"./index.md":             "/level2/level3/",
			//			"../level2/level3/3-root.md":    "/level2/level3/3-root/",
			//			"../level2/level3/index.md":     "/level2/level3/",
			"/docs/level2/level3/3-root.md": "/level2/level3/3-root/",
			"/docs/level2/level3/index.md":  "/level2/level3/",
		}, "level2/level3/index.md": map[string]string{
			"../../rootfile.md":      "/rootfile/",
			"../../index.md":         "/",
			"/docs/rootfile.md":      "/rootfile/",
			"/docs/index.md":         "/",
			"../2-root.md":           "/level2/2-root/",
			"../index.md":            "/level2/",
			"/docs/level2/2-root.md": "/level2/2-root/",
			"/docs/level2/index.md":  "/level2/",
			"3-root.md":              "/level2/level3/3-root/",
			"index.md":               "/level2/level3/",
			"./3-root.md":            "/level2/level3/3-root/",
			"./index.md":             "/level2/level3/",
			//			"../level2/level3/3-root.md":    "/level2/level3/3-root/",
			//			"../level2/level3/index.md":     "/level2/level3/",
			"/docs/level2/level3/3-root.md": "/level2/level3/3-root/",
			"/docs/level2/level3/index.md":  "/level2/level3/",
		},
	}

	for currentFile, results := range okresults {
		currentPage := findPage(site, currentFile)
		if currentPage == nil {
			t.Fatalf("failed to find current page in site")
		}
		for link, url := range results {
			if out, err := site.Info.githubLink(link, currentPage, true); err != nil || out != url {
				t.Errorf("Expected %s to resolve to (%s), got (%s) - error: %s", link, url, out, err)
			} else {
				//t.Logf("tested ok %s maps to %s", link, out)
			}
		}
	}
	// TODO: and then the failure cases.
	// 			"https://docker.com":           "",
	// site_test.go:1094: Expected https://docker.com to resolve to (), got () - error: Not a plain filepath link (https://docker.com)

}

func TestGitHubFileLinking(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	site := setupLinkingMockSite(t)

	type resultMap map[string]string

	okresults := map[string]resultMap{
		"index.md": map[string]string{
			"/root-image.png": "/root-image.png",
			"root-image.png":  "/root-image.png",
		}, "rootfile.md": map[string]string{
			"/root-image.png": "/root-image.png",
		}, "level2/2-root.md": map[string]string{
			"/root-image.png": "/root-image.png",
			"common.png":      "/level2/common.png",
		}, "level2/index.md": map[string]string{
			"/root-image.png": "/root-image.png",
			"common.png":      "/level2/common.png",
			"./common.png":    "/level2/common.png",
		}, "level2/level3/3-root.md": map[string]string{
			"/root-image.png": "/root-image.png",
			"common.png":      "/level2/level3/common.png",
			"../common.png":   "/level2/common.png",
		}, "level2/level3/index.md": map[string]string{
			"/root-image.png": "/root-image.png",
			"common.png":      "/level2/level3/common.png",
			"../common.png":   "/level2/common.png",
		},
	}

	for currentFile, results := range okresults {
		currentPage := findPage(site, currentFile)
		if currentPage == nil {
			t.Fatalf("failed to find current page in site")
		}
		for link, url := range results {
			if out, err := site.Info.githubFileLink(link, currentPage, false); err != nil || out != url {
				t.Errorf("Expected %s to resolve to (%s), got (%s) - error: %s", link, url, out, err)
			} else {
				//t.Logf("tested ok %s maps to %s", link, out)
			}
		}
	}
	// TODO: and then the failure cases.
	// 			"https://docker.com":           "",
	// site_test.go:1094: Expected https://docker.com to resolve to (), got () - error: Not a plain filepath link (https://docker.com)
}

func TestMultilingualSwitch(t *testing.T) {
	// General settings
	viper.Set("DefaultExtension", "html")
	viper.Set("baseurl", "http://example.com/blog")
	viper.Set("DisableSitemap", false)
	viper.Set("DisableRSS", false)
	viper.Set("RSSUri", "index.xml")
	viper.Set("Taxonomies", map[string]string{"tag": "tags"})
	viper.Set("Permalinks", map[string]string{"other": "/somewhere/else/:filename"})

	// Multilingual settings
	viper.Set("Multilingual", true)
	viper.Set("RenderLanguage", "en")
	viper.Set("DefaultContentLang", "fr")
	viper.Set("paginate", "2")

	// Sources
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.en.md"), []byte(`---
title: doc1
slug: doc1-slug
tags:
 - tag1
publishdate: "2000-01-01"
---
# doc1
*some content*
NOTE: slug should be used as URL
`)},
		{filepath.FromSlash("sect/doc1.fr.md"), []byte(`---
title: doc1
tags:
 - tag1
 - tag2
publishdate: "2000-01-04"
---
# doc1
*quelque contenu*
NOTE: should be in the 'en' Page's 'Translations' field.
NOTE: date is after "doc3"
`)},
		{filepath.FromSlash("sect/doc2.en.md"), []byte(`---
title: doc2
publishdate: "2000-01-02"
---
# doc2
*some content*
NOTE: without slug, "doc2" should be used, without ".en" as URL
`)},
		{filepath.FromSlash("sect/doc3.en.md"), []byte(`---
title: doc3
publishdate: "2000-01-03"
tags:
 - tag2
url: /superbob
---
# doc3
*some content*
NOTE: third 'en' doc, should trigger pagination on home page.
`)},
		{filepath.FromSlash("sect/doc4.md"), []byte(`---
title: doc4
tags:
 - tag1
publishdate: "2000-01-05"
---
# doc4
*du contenu francophone*
NOTE: should use the DefaultContentLang and mark this doc as 'fr'.
NOTE: doesn't have any corresponding translation in 'en'
`)},
		{filepath.FromSlash("other/doc5.fr.md"), []byte(`---
title: doc5
publishdate: "2000-01-06"
---
# doc5
*autre contenu francophone*
NOTE: should use the "permalinks" configuration with :filename
`)},
	}

	hugofs.DestinationFS = new(afero.MemMapFs)

	s := &Site{
		Source: &source.InMemorySource{ByteSource: sources},
	}
	templatePrep(s)
	s.initializeSiteInfo()

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	s.setupTranslations()
	s.setupPrevNext()

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	assert.Len(t, s.Source.Files(), 6, "should have 6 source files")
	assert.Len(t, s.Pages, 3, "should have 3 pages")
	assert.Len(t, s.TranslatedPages, 3, "should have 3 translated pages")

	doc1en := s.Pages[0]
	permalink, err := doc1en.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/en/sect/doc1-slug", permalink, "invalid doc1.en permalink")

	assert.Len(t, doc1en.Translations, 1, "doc1-en should have one translation")

	doc2 := s.Pages[1]
	permalink, err = doc2.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/en/sect/doc2", permalink, "invalid doc2 permalink")

	doc3 := s.Pages[2]
	permalink, err = doc3.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/superbob", permalink, "invalid doc3 permalink")
	assert.Equal(t, "/superbob", doc3.URL, "invalid url, was specified on doc3")

	assert.Equal(t, doc2.Next, doc3, "doc3 should follow doc2, in .Next")

	doc1fr := s.TranslatedPages[0]
	permalink, err = doc1fr.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/fr/sect/doc1", permalink, "invalid doc1fr permalink")

	assert.Equal(t, doc1en.Translations["fr"], doc1fr, "doc1-en should have doc1-fr as translation")
	assert.Equal(t, doc1fr.Translations["en"], doc1en, "doc1-fr should have doc1-en as translation")

	doc4 := s.TranslatedPages[1]
	permalink, err = doc4.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/fr/sect/doc4", permalink, "invalid doc4 permalink")
	assert.Len(t, doc4.Translations, 0, "found translations for doc4")

	doc5 := s.TranslatedPages[2]
	permalink, err = doc5.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/fr/somewhere/else/doc5", permalink, "invalid doc5 permalink")

	// Taxonomies and their URLs
	assert.Len(t, s.Taxonomies, 1, "should have 1 taxonomy")
	tags := s.Taxonomies["tags"]
	assert.Len(t, tags, 2, "should have 2 different tags")
	assert.Equal(t, tags["tag1"][0].Page, doc1en, "first tag1 page should be doc1")

	// Expect the tags locations to be in certain places, with the /en/ prefixes, etc..
}

func assertFileContent(t *testing.T, path string, content string) {
	fl, err := hugofs.DestinationFS.Open(path)
	assert.NoError(t, err, "file content not found when asserting on content of %s", path)

	cnt, err := ioutil.ReadAll(fl)
	assert.NoError(t, err, "cannot read file content when asserting on content of %s", path)

	assert.Equal(t, content, string(cnt))
}
