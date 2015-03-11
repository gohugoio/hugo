package hugolib

import (
	"bytes"
	"fmt"
	"github.com/spf13/hugo/parser"
	"html/template"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/viper"
	"reflect"
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

	s := new(Site)
	templatePrep(s)

	for i, test := range tests {
		p, err := NewPageFrom(strings.NewReader(test.content), "content/a/file.md")
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

	hugofs.DestinationFS = new(afero.MemMapFs)
	s := &Site{}
	templatePrep(s)

	for i, test := range tests {
		p, err := NewPageFrom(strings.NewReader(PAGE_SIMPLE_TITLE), "content/a/file.md")
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
		p := pageMust(NewPageFrom(strings.NewReader(test.content), helpers.AbsPathify(filepath.FromSlash(test.doc))))

		expected := filepath.FromSlash(test.expectedOutFile)

		if p.TargetPath() != expected {
			t.Errorf("%s => OutFile  expected: '%s', got: '%s'", test.doc, expected, p.TargetPath())
		}

		if p.Section() != test.expectedSection {
			t.Errorf("%s => p.Section expected: %s, got: %s", test.doc, test.expectedSection, p.Section())
		}
	}
}

func TestDraftAndFutureRender(t *testing.T) {
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

// Issue #939
func Test404ShouldAlwaysHaveUglyUrls(t *testing.T) {
	for _, uglyUrls := range []bool{true, false} {
		doTest404ShouldAlwaysHaveUglyUrls(t, uglyUrls)
	}
}

func doTest404ShouldAlwaysHaveUglyUrls(t *testing.T, uglyUrls bool) {
	viper.Set("verbose", true)
	viper.Set("baseurl", "http://auth/bub")
	viper.Set("DisableSitemap", false)
	viper.Set("DisableRSS", false)

	viper.Set("UglyUrls", uglyUrls)

	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
	}

	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		Targets: targetList{Page: &target.PagePub{UglyUrls: uglyUrls}},
	}

	s.initializeSiteInfo()
	templatePrep(s)

	must(s.addTemplate("index.html", "Home Sweet Home"))
	must(s.addTemplate("_default/single.html", "{{.Content}}"))
	must(s.addTemplate("404.html", "Page Not Found"))

	// make sure the XML files also end up with ugly urls
	must(s.addTemplate("rss.xml", "<root>RSS</root>"))
	must(s.addTemplate("sitemap.xml", "<root>SITEMAP</root>"))

	createAndRenderPages(t, s)
	s.RenderHomePage()
	s.RenderSitemap()

	var expectedPagePath string
	if uglyUrls {
		expectedPagePath = "sect/doc1.html"
	} else {
		expectedPagePath = "sect/doc1/index.html"
	}

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("index.html"), "Home Sweet Home"},
		{filepath.FromSlash(expectedPagePath), "\n\n<h1 id=\"title:5d74edbb89ef198cd37882b687940cda\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("404.html"), "Page Not Found"},
		{filepath.FromSlash("index.xml"), "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n<root>RSS</root>"},
		{filepath.FromSlash("sitemap.xml"), "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n<root>SITEMAP</root>"},
	}

	for _, test := range tests {
		file, err := hugofs.DestinationFS.Open(test.doc)
		if err != nil {
			t.Fatalf("Did not find %s in target: %s", test.doc, err)
		}
		content := helpers.ReaderToBytes(file)

		if !bytes.Equal(content, []byte(test.expected)) {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, string(content))
		}
	}

}

func TestSkipRender(t *testing.T) {
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

	viper.Set("verbose", true)
	viper.Set("CanonifyUrls", true)
	viper.Set("baseurl", "http://auth/bub")
	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		Targets: targetList{Page: &target.PagePub{UglyUrls: true}},
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
		content := helpers.ReaderToBytes(file)

		if !bytes.Equal(content, []byte(test.expected)) {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, string(content))
		}
	}
}

func TestAbsUrlify(t *testing.T) {
	hugofs.DestinationFS = new(afero.MemMapFs)
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>")},
		{filepath.FromSlash("content/blue/doc2.html"), []byte("---\nf: t\n---\n<!doctype html><html><body>more content</body></html>")},
	}
	for _, canonify := range []bool{true, false} {
		viper.Set("CanonifyUrls", canonify)
		viper.Set("BaseUrl", "http://auth/bub")
		s := &Site{
			Source:  &source.InMemorySource{ByteSource: sources},
			Targets: targetList{Page: &target.PagePub{UglyUrls: true}},
		}
		t.Logf("Rendering with BaseUrl %q and CanonifyUrls set %v", viper.GetString("baseUrl"), canonify)
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
			content := helpers.ReaderToBytes(file)

			expected := test.expected
			if !canonify {
				expected = strings.Replace(expected, viper.GetString("baseurl"), "", -1)
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

func TestDataDirJson(t *testing.T) {
	sources := []source.ByteSource{
		{filepath.FromSlash("test/foo.json"), []byte(`{ "bar": "foofoo"  }`)},
		{filepath.FromSlash("test.json"), []byte(`{ "hello": [ { "world": "foo" } ] }`)},
	}

	expected, err := parser.HandleJsonMetaData([]byte(`{ "test": { "hello": [{ "world": "foo"  }] , "foo": { "bar":"foofoo" } } }`))

	if err != nil {
		t.Fatalf("Error %s", err)
	}

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: sources}})
}

func TestDataDirToml(t *testing.T) {
	sources := []source.ByteSource{
		{filepath.FromSlash("test/kung.toml"), []byte("[foo]\nbar = 1")},
	}

	expected, err := parser.HandleTomlMetaData([]byte("[test]\n[test.kung]\n[test.kung.foo]\nbar = 1"))

	if err != nil {
		t.Fatalf("Error %s", err)
	}

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: sources}})
}

func TestDataDirYamlWithOverridenValue(t *testing.T) {
	sources := []source.ByteSource{
		// filepath.Walk walks the files in lexical order, '/' comes before '.'. Simulate this:
		{filepath.FromSlash("a.yaml"), []byte("a: 1")},
		{filepath.FromSlash("test/v1.yaml"), []byte("v1-2: 2")},
		{filepath.FromSlash("test/v2.yaml"), []byte("v2:\n- 2\n- 3")},
		{filepath.FromSlash("test.yaml"), []byte("v1: 1")},
	}

	expected := map[string]interface{}{"a": map[string]interface{}{"a": 1},
		"test": map[string]interface{}{"v1": map[string]interface{}{"v1-2": 2}, "v2": map[string]interface{}{"v2": []interface{}{2, 3}}}}

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: sources}})
}

// issue 892
func TestDataDirMultipleSources(t *testing.T) {
	s1 := []source.ByteSource{
		{filepath.FromSlash("test/first.toml"), []byte("bar = 1")},
	}

	s2 := []source.ByteSource{
		{filepath.FromSlash("test/first.toml"), []byte("bar = 2")},
		{filepath.FromSlash("test/second.toml"), []byte("tender = 2")},
	}

	expected, _ := parser.HandleTomlMetaData([]byte("[test.first]\nbar = 1\n[test.second]\ntender=2"))

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: s1}, &source.InMemorySource{ByteSource: s2}})

}

func TestDataDirUnknownFormat(t *testing.T) {
	sources := []source.ByteSource{
		{filepath.FromSlash("test.roml"), []byte("boo")},
	}
	s := &Site{}
	err := s.loadData([]source.Input{&source.InMemorySource{ByteSource: sources}})
	if err == nil {
		t.Fatalf("Should return an error")
	}
}

func doTestDataDir(t *testing.T, expected interface{}, sources []source.Input) {
	s := &Site{}
	err := s.loadData(sources)
	if err != nil {
		t.Fatalf("Error loading data: %s", err)
	}
	if !reflect.DeepEqual(expected, s.Data) {
		t.Errorf("Expected structure\n%#v got\n%#v", expected, s.Data)
	}
}
