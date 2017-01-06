// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bep/inflect"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"

	"github.com/spf13/hugo/target"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	pageSimpleTitle = `---
title: simple template
---
content`

	templateMissingFunc = "{{ .Title | funcdoesnotexists }}"
	templateWithURLAbs  = "<a href=\"/foobar.jpg\">Going</a>"
)

func init() {
	testMode = true
}

// Issue #1797
func TestReadPagesFromSourceWithEmptySource(t *testing.T) {
	testCommonResetState()

	viper.Set("defaultExtension", "html")
	viper.Set("verbose", true)
	viper.Set("baseURL", "http://auth/bub")

	sources := []source.ByteSource{}

	s := &Site{
		deps:    newDeps(DepsCfg{}),
		Source:  &source.InMemorySource{ByteSource: sources},
		targets: targetList{page: &target.PagePub{UglyURLs: true}},
	}

	var err error
	d := time.Second * 2
	ticker := time.NewTicker(d)
	select {
	case err = <-s.readPagesFromSource():
		break
	case <-ticker.C:
		err = fmt.Errorf("ReadPagesFromSource() never returns in %s", d.String())
	}
	ticker.Stop()
	if err != nil {
		t.Fatalf("Unable to read source: %s", err)
	}
}

func pageMust(p *Page, err error) *Page {
	if err != nil {
		panic(err)
	}
	return p
}

func TestDegenerateRenderThingMissingTemplate(t *testing.T) {
	s := newSiteFromSources("content/a/file.md", pageSimpleTitle)

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	require.Len(t, s.RegularPages, 1)

	p := s.RegularPages[0]

	err := s.renderThing(p, "foobar", nil)
	if err == nil {
		t.Errorf("Expected err to be returned when missing the template.")
	}
}

func TestRenderWithInvalidTemplate(t *testing.T) {
	jww.ResetLogCounters()

	s := NewSiteDefaultLang()
	if err := buildAndRenderSite(s, "missing", templateMissingFunc); err != nil {
		t.Fatalf("Got build error: %s", err)
	}

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) != 1 {
		t.Fatalf("Expecting the template to log an ERROR")
	}
}

func TestDraftAndFutureRender(t *testing.T) {
	testCommonResetState()

	hugofs.InitMemFs()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.md"), Content: []byte("---\ntitle: doc1\ndraft: true\npublishdate: \"2414-05-29\"\n---\n# doc1\n*some content*")},
		{Name: filepath.FromSlash("sect/doc2.md"), Content: []byte("---\ntitle: doc2\ndraft: true\npublishdate: \"2012-05-29\"\n---\n# doc2\n*some content*")},
		{Name: filepath.FromSlash("sect/doc3.md"), Content: []byte("---\ntitle: doc3\ndraft: false\npublishdate: \"2414-05-29\"\n---\n# doc3\n*some content*")},
		{Name: filepath.FromSlash("sect/doc4.md"), Content: []byte("---\ntitle: doc4\ndraft: false\npublishdate: \"2012-05-29\"\n---\n# doc4\n*some content*")},
	}

	siteSetup := func(t *testing.T) *Site {
		s := &Site{
			deps:     newDeps(DepsCfg{}),
			Source:   &source.InMemorySource{ByteSource: sources},
			Language: helpers.NewDefaultLanguage(),
		}

		if err := buildSiteSkipRender(s); err != nil {
			t.Fatalf("Failed to build site: %s", err)
		}

		return s
	}

	viper.Set("baseURL", "http://auth/bub")

	// Testing Defaults.. Only draft:true and publishDate in the past should be rendered
	s := siteSetup(t)
	if len(s.RegularPages) != 1 {
		t.Fatal("Draft or Future dated content published unexpectedly")
	}

	// only publishDate in the past should be rendered
	viper.Set("buildDrafts", true)
	s = siteSetup(t)
	if len(s.RegularPages) != 2 {
		t.Fatal("Future Dated Posts published unexpectedly")
	}

	//  drafts should not be rendered, but all dates should
	viper.Set("buildDrafts", false)
	viper.Set("buildFuture", true)
	s = siteSetup(t)
	if len(s.RegularPages) != 2 {
		t.Fatal("Draft posts published unexpectedly")
	}

	// all 4 should be included
	viper.Set("buildDrafts", true)
	viper.Set("buildFuture", true)
	s = siteSetup(t)
	if len(s.RegularPages) != 4 {
		t.Fatal("Drafts or Future posts not included as expected")
	}

	//setting defaults back
	viper.Set("buildDrafts", false)
	viper.Set("buildFuture", false)
}

func TestFutureExpirationRender(t *testing.T) {
	testCommonResetState()

	hugofs.InitMemFs()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc3.md"), Content: []byte("---\ntitle: doc1\nexpirydate: \"2400-05-29\"\n---\n# doc1\n*some content*")},
		{Name: filepath.FromSlash("sect/doc4.md"), Content: []byte("---\ntitle: doc2\nexpirydate: \"2000-05-29\"\n---\n# doc2\n*some content*")},
	}

	siteSetup := func(t *testing.T) *Site {
		s := &Site{
			deps:     newDeps(DepsCfg{}),
			Source:   &source.InMemorySource{ByteSource: sources},
			Language: helpers.NewDefaultLanguage(),
		}

		if err := buildSiteSkipRender(s); err != nil {
			t.Fatalf("Failed to build site: %s", err)
		}

		return s
	}

	viper.Set("baseURL", "http://auth/bub")

	s := siteSetup(t)

	if len(s.AllPages) != 1 {
		if len(s.RegularPages) > 1 {
			t.Fatal("Expired content published unexpectedly")
		}

		if len(s.RegularPages) < 1 {
			t.Fatal("Valid content expired unexpectedly")
		}
	}

	if s.AllPages[0].Title == "doc2" {
		t.Fatal("Expired content published unexpectedly")
	}
}

// Issue #957
func TestCrossrefs(t *testing.T) {
	for _, uglyURLs := range []bool{true, false} {
		for _, relative := range []bool{true, false} {
			doTestCrossrefs(t, relative, uglyURLs)
		}
	}
}

func doTestCrossrefs(t *testing.T, relative, uglyURLs bool) {
	testCommonResetState()

	baseURL := "http://foo/bar"
	viper.Set("defaultExtension", "html")
	viper.Set("baseURL", baseURL)
	viper.Set("uglyURLs", uglyURLs)
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
		{
			Name:    filepath.FromSlash("sect/doc1.md"),
			Content: []byte(fmt.Sprintf(`Ref 2: {{< %s "sect/doc2.md" >}}`, refShortcode)),
		},
		// Issue #1148: Make sure that no P-tags is added around shortcodes.
		{
			Name: filepath.FromSlash("sect/doc2.md"),
			Content: []byte(fmt.Sprintf(`**Ref 1:**

{{< %s "sect/doc1.md" >}}

THE END.`, refShortcode)),
		},
		// Issue #1753: Should not add a trailing newline after shortcode.
		{
			Name:    filepath.FromSlash("sect/doc3.md"),
			Content: []byte(fmt.Sprintf(`**Ref 1:**{{< %s "sect/doc3.md" >}}.`, refShortcode)),
		},
	}

	s := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: sources},
		targets:  targetList{page: &target.PagePub{UglyURLs: uglyURLs}},
		Language: helpers.NewDefaultLanguage(),
	}

	if err := buildAndRenderSite(s, "_default/single.html", "{{.Content}}"); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	if len(s.RegularPages) != 3 {
		t.Fatalf("Expected 3 got %d pages", len(s.AllPages))
	}

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash(fmt.Sprintf("sect/doc1%s", expectedPathSuffix)), fmt.Sprintf("<p>Ref 2: %s/sect/doc2%s</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("sect/doc2%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong></p>\n\n%s/sect/doc1%s\n\n<p>THE END.</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("sect/doc3%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong>%s/sect/doc3%s.</p>\n", expectedBase, expectedURLSuffix)},
	}

	for _, test := range tests {
		file, err := hugofs.Destination().Open(test.doc)

		if err != nil {
			t.Fatalf("Did not find %s in target: %s", test.doc, err)
		}

		content := helpers.ReaderToString(file)

		if content != test.expected {
			t.Fatalf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}

// Issue #939
// Issue #1923
func TestShouldAlwaysHaveUglyURLs(t *testing.T) {
	for _, uglyURLs := range []bool{true, false} {
		doTestShouldAlwaysHaveUglyURLs(t, uglyURLs)
	}
}

func doTestShouldAlwaysHaveUglyURLs(t *testing.T, uglyURLs bool) {
	testCommonResetState()

	viper.Set("defaultExtension", "html")
	viper.Set("verbose", true)
	viper.Set("baseURL", "http://auth/bub")
	viper.Set("disableSitemap", false)
	viper.Set("disableRSS", false)
	viper.Set("rssURI", "index.xml")
	viper.Set("blackfriday",
		map[string]interface{}{
			"plainIDAnchors": true})

	viper.Set("uglyURLs", uglyURLs)

	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.md"), Content: []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
		{Name: filepath.FromSlash("sect/doc2.md"), Content: []byte("---\nurl: /ugly.html\nmarkup: markdown\n---\n# title\ndoc2 *content*")},
	}

	s := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: sources},
		targets:  targetList{page: &target.PagePub{UglyURLs: uglyURLs, PublishDir: "public"}},
		Language: helpers.NewDefaultLanguage(),
	}

	if err := buildAndRenderSite(s,
		"index.html", "Home Sweet {{ if.IsHome  }}Home{{ end }}.",
		"_default/single.html", "{{.Content}}{{ if.IsHome  }}This is not home!{{ end }}",
		"404.html", "Page Not Found.{{ if.IsHome  }}This is not home!{{ end }}",
		"rss.xml", "<root>RSS</root>",
		"sitemap.xml", "<root>SITEMAP</root>"); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	var expectedPagePath string
	if uglyURLs {
		expectedPagePath = "public/sect/doc1.html"
	} else {
		expectedPagePath = "public/sect/doc1/index.html"
	}

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("public/index.html"), "Home Sweet Home."},
		{filepath.FromSlash(expectedPagePath), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("public/404.html"), "Page Not Found."},
		{filepath.FromSlash("public/index.xml"), "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n<root>RSS</root>"},
		{filepath.FromSlash("public/sitemap.xml"), "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n<root>SITEMAP</root>"},
		// Issue #1923
		{filepath.FromSlash("public/ugly.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>doc2 <em>content</em></p>\n"},
	}

	for _, p := range s.RegularPages {
		assert.False(t, p.IsHome())
	}

	for _, test := range tests {
		content := readDestination(t, test.doc)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}

// Issue #1176
func TestSectionNaming(t *testing.T) {
	//jww.SetStdoutThreshold(jww.LevelDebug)

	for _, canonify := range []bool{true, false} {
		for _, uglify := range []bool{true, false} {
			for _, pluralize := range []bool{true, false} {
				doTestSectionNaming(t, canonify, uglify, pluralize)
			}
		}
	}
}

func doTestSectionNaming(t *testing.T, canonify, uglify, pluralize bool) {
	testCommonResetState()

	viper.Set("baseURL", "http://auth/sub/")
	viper.Set("defaultExtension", "html")
	viper.Set("uglyURLs", uglify)
	viper.Set("pluralizeListTitles", pluralize)
	viper.Set("canonifyURLs", canonify)

	var expectedPathSuffix string

	if uglify {
		expectedPathSuffix = ".html"
	} else {
		expectedPathSuffix = "/index.html"
	}

	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.html"), Content: []byte("doc1")},
		{Name: filepath.FromSlash("Fish and Chips/doc2.html"), Content: []byte("doc2")},
		{Name: filepath.FromSlash("ラーメン/doc3.html"), Content: []byte("doc3")},
	}

	for _, source := range sources {
		writeSource(t, filepath.Join("content", source.Name), string(source.Content))
	}

	s := NewSiteDefaultLang()

	if err := buildAndRenderSite(s,
		"_default/single.html", "{{.Content}}",
		"_default/list.html", "{{ .Title }}"); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

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

		if test.pluralAware && pluralize {
			test.expected = inflect.Pluralize(test.expected)
		}

		assertFileContent(t, filepath.Join("public", test.doc), true, test.expected)
	}

}
func TestSkipRender(t *testing.T) {
	testCommonResetState()

	hugofs.InitMemFs()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.html"), Content: []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
		{Name: filepath.FromSlash("sect/doc2.html"), Content: []byte("<!doctype html><html><body>more content</body></html>")},
		{Name: filepath.FromSlash("sect/doc3.md"), Content: []byte("# doc3\n*some* content")},
		{Name: filepath.FromSlash("sect/doc4.md"), Content: []byte("---\ntitle: doc4\n---\n# doc4\n*some content*")},
		{Name: filepath.FromSlash("sect/doc5.html"), Content: []byte("<!doctype html><html>{{ template \"head\" }}<body>body5</body></html>")},
		{Name: filepath.FromSlash("sect/doc6.html"), Content: []byte("<!doctype html><html>{{ template \"head_abs\" }}<body>body5</body></html>")},
		{Name: filepath.FromSlash("doc7.html"), Content: []byte("<html><body>doc7 content</body></html>")},
		{Name: filepath.FromSlash("sect/doc8.html"), Content: []byte("---\nmarkup: md\n---\n# title\nsome *content*")},
	}

	viper.Set("defaultExtension", "html")
	viper.Set("verbose", true)
	viper.Set("canonifyURLs", true)
	viper.Set("baseURL", "http://auth/bub")
	s := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: sources},
		targets:  targetList{page: &target.PagePub{UglyURLs: true}},
		Language: helpers.NewDefaultLanguage(),
	}

	if err := buildAndRenderSite(s,
		"_default/single.html", "{{.Content}}",
		"head", "<head><script src=\"script.js\"></script></head>",
		"head_abs", "<head><script src=\"/script.js\"></script></head>"); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("sect/doc1.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("sect/doc2.html"), "<!doctype html><html><body>more content</body></html>"},
		{filepath.FromSlash("sect/doc3.html"), "\n\n<h1 id=\"doc3\">doc3</h1>\n\n<p><em>some</em> content</p>\n"},
		{filepath.FromSlash("sect/doc4.html"), "\n\n<h1 id=\"doc4\">doc4</h1>\n\n<p><em>some content</em></p>\n"},
		{filepath.FromSlash("sect/doc5.html"), "<!doctype html><html><head><script src=\"script.js\"></script></head><body>body5</body></html>"},
		{filepath.FromSlash("sect/doc6.html"), "<!doctype html><html><head><script src=\"http://auth/bub/script.js\"></script></head><body>body5</body></html>"},
		{filepath.FromSlash("doc7.html"), "<html><body>doc7 content</body></html>"},
		{filepath.FromSlash("sect/doc8.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
	}

	for _, test := range tests {
		file, err := hugofs.Destination().Open(test.doc)
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
	testCommonResetState()

	viper.Set("defaultExtension", "html")

	hugofs.InitMemFs()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.html"), Content: []byte("<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>")},
		{Name: filepath.FromSlash("blue/doc2.html"), Content: []byte("---\nf: t\n---\n<!doctype html><html><body>more content</body></html>")},
	}
	for _, baseURL := range []string{"http://auth/bub", "http://base", "//base"} {
		for _, canonify := range []bool{true, false} {
			viper.Set("canonifyURLs", canonify)
			viper.Set("baseURL", baseURL)
			s := &Site{
				deps:     newDeps(DepsCfg{}),
				Source:   &source.InMemorySource{ByteSource: sources},
				targets:  targetList{page: &target.PagePub{UglyURLs: true}},
				Language: helpers.NewDefaultLanguage(),
			}
			t.Logf("Rendering with baseURL %q and canonifyURLs set %v", viper.GetString("baseURL"), canonify)

			if err := buildAndRenderSite(s, "blue/single.html", templateWithURLAbs); err != nil {
				t.Fatalf("Failed to build site: %s", err)
			}

			tests := []struct {
				file, expected string
			}{
				{"blue/doc2.html", "<a href=\"%s/foobar.jpg\">Going</a>"},
				{"sect/doc1.html", "<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>"},
			}

			for _, test := range tests {

				file, err := hugofs.Destination().Open(filepath.FromSlash(test.file))
				if err != nil {
					t.Fatalf("Unable to locate rendered content: %s", test.file)
				}

				content := helpers.ReaderToString(file)

				expected := test.expected

				if strings.Contains(expected, "%s") {
					expected = fmt.Sprintf(expected, baseURL)
				}

				if !canonify {
					expected = strings.Replace(expected, baseURL, "", -1)
				}

				if content != expected {
					t.Errorf("AbsURLify with baseURL %q content expected:\n%q\ngot\n%q", baseURL, expected, content)
				}
			}
		}
	}
}

var weightedPage1 = []byte(`+++
weight = "2"
title = "One"
my_param = "foo"
my_date = 1979-05-27T07:32:00Z
+++
Front Matter with Ordered Pages`)

var weightedPage2 = []byte(`+++
weight = "6"
title = "Two"
publishdate = "2012-03-05"
my_param = "foo"
+++
Front Matter with Ordered Pages 2`)

var weightedPage3 = []byte(`+++
weight = "4"
title = "Three"
date = "2012-04-06"
publishdate = "2012-04-06"
my_param = "bar"
only_one = "yes"
my_date = 2010-05-27T07:32:00Z
+++
Front Matter with Ordered Pages 3`)

var weightedPage4 = []byte(`+++
weight = "4"
title = "Four"
date = "2012-01-01"
publishdate = "2012-01-01"
my_param = "baz"
my_date = 2010-05-27T07:32:00Z
categories = [ "hugo" ]
+++
Front Matter with Ordered Pages 4. This is longer content`)

var weightedSources = []source.ByteSource{
	{Name: filepath.FromSlash("sect/doc1.md"), Content: weightedPage1},
	{Name: filepath.FromSlash("sect/doc2.md"), Content: weightedPage2},
	{Name: filepath.FromSlash("sect/doc3.md"), Content: weightedPage3},
	{Name: filepath.FromSlash("sect/doc4.md"), Content: weightedPage4},
}

func TestOrderedPages(t *testing.T) {
	testCommonResetState()

	hugofs.InitMemFs()

	viper.Set("baseURL", "http://auth/bub")
	s := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: weightedSources},
		Language: helpers.NewDefaultLanguage(),
	}

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to process site: %s", err)
	}

	if s.Sections["sect"][0].Weight != 2 || s.Sections["sect"][3].Weight != 6 {
		t.Errorf("Pages in unexpected order. First should be '%d', got '%d'", 2, s.Sections["sect"][0].Weight)
	}

	if s.Sections["sect"][1].Page.Title != "Three" || s.Sections["sect"][2].Page.Title != "Four" {
		t.Errorf("Pages in unexpected order. Second should be '%s', got '%s'", "Three", s.Sections["sect"][1].Page.Title)
	}

	bydate := s.RegularPages.ByDate()

	if bydate[0].Title != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bydate[0].Title)
	}

	rev := bydate.Reverse()
	if rev[0].Title != "Three" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Three", rev[0].Title)
	}

	bypubdate := s.RegularPages.ByPublishDate()

	if bypubdate[0].Title != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bypubdate[0].Title)
	}

	rbypubdate := bypubdate.Reverse()
	if rbypubdate[0].Title != "Three" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Three", rbypubdate[0].Title)
	}

	bylength := s.RegularPages.ByLength()
	if bylength[0].Title != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bylength[0].Title)
	}

	rbylength := bylength.Reverse()
	if rbylength[0].Title != "Four" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Four", rbylength[0].Title)
	}
}

var groupedSources = []source.ByteSource{
	{Name: filepath.FromSlash("sect1/doc1.md"), Content: weightedPage1},
	{Name: filepath.FromSlash("sect1/doc2.md"), Content: weightedPage2},
	{Name: filepath.FromSlash("sect2/doc3.md"), Content: weightedPage3},
	{Name: filepath.FromSlash("sect3/doc4.md"), Content: weightedPage4},
}

func TestGroupedPages(t *testing.T) {
	testCommonResetState()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	hugofs.InitMemFs()

	viper.Set("baseURL", "http://auth/bub")
	s := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: groupedSources},
		Language: helpers.NewDefaultLanguage(),
	}

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	rbysection, err := s.RegularPages.GroupBy("Section", "desc")
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

	bytype, err := s.RegularPages.GroupBy("Type", "asc")
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

	bydate, err := s.RegularPages.GroupByDate("2006-01", "asc")
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

	bypubdate, err := s.RegularPages.GroupByPublishDate("2006")
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

	byparam, err := s.RegularPages.GroupByParam("my_param", "desc")
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

	_, err = s.RegularPages.GroupByParam("not_exist")
	if err == nil {
		t.Errorf("GroupByParam didn't return an expected error")
	}

	byOnlyOneParam, err := s.RegularPages.GroupByParam("only_one")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if len(byOnlyOneParam) != 1 {
		t.Errorf("PageGroup array has unexpected elements. Group length should be '%d', got '%d'", 1, len(byOnlyOneParam))
	}
	if byOnlyOneParam[0].Key != "yes" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "yes", byOnlyOneParam[0].Key)
	}

	byParamDate, err := s.RegularPages.GroupByParamDate("my_date", "2006-01")
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

var pageWithWeightedTaxonomies1 = []byte(`+++
tags = [ "a", "b", "c" ]
tags_weight = 22
categories = ["d"]
title = "foo"
categories_weight = 44
+++
Front Matter with weighted tags and categories`)

var pageWithWeightedTaxonomies2 = []byte(`+++
tags = "a"
tags_weight = 33
title = "bar"
categories = [ "d", "e" ]
categories_weight = 11
alias = "spf13"
date = 1979-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`)

var pageWithWeightedTaxonomies3 = []byte(`+++
title = "bza"
categories = [ "e" ]
categories_weight = 11
alias = "spf13"
date = 2010-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`)

func TestWeightedTaxonomies(t *testing.T) {
	testCommonResetState()

	hugofs.InitMemFs()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.md"), Content: pageWithWeightedTaxonomies2},
		{Name: filepath.FromSlash("sect/doc2.md"), Content: pageWithWeightedTaxonomies1},
		{Name: filepath.FromSlash("sect/doc3.md"), Content: pageWithWeightedTaxonomies3},
	}
	taxonomies := make(map[string]string)

	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	viper.Set("baseURL", "http://auth/bub")
	viper.Set("taxonomies", taxonomies)
	s := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: sources},
		Language: helpers.NewDefaultLanguage(),
	}

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to process site: %s", err)
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
	hugofs.InitMemFs()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("index.md"), Content: []byte("")},
		{Name: filepath.FromSlash("rootfile.md"), Content: []byte("")},
		{Name: filepath.FromSlash("root-image.png"), Content: []byte("")},

		{Name: filepath.FromSlash("level2/2-root.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/index.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/common.md"), Content: []byte("")},

		//		{Name: filepath.FromSlash("level2b/2b-root.md"), Content: []byte("")},
		//		{Name: filepath.FromSlash("level2b/index.md"), Content: []byte("")},
		//		{Name: filepath.FromSlash("level2b/common.md"), Content: []byte("")},

		{Name: filepath.FromSlash("level2/2-image.png"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/common.png"), Content: []byte("")},

		{Name: filepath.FromSlash("level2/level3/3-root.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/index.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/common.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/3-image.png"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/common.png"), Content: []byte("")},
	}

	viper.Set("baseURL", "http://auth/")
	viper.Set("defaultExtension", "html")
	viper.Set("uglyURLs", false)
	viper.Set("pluralizeListTitles", false)
	viper.Set("canonifyURLs", false)
	viper.Set("blackfriday",
		map[string]interface{}{
			"sourceRelativeLinksProjectFolder": "/docs"})

	site := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: sources},
		Language: helpers.NewDefaultLanguage(),
	}

	if err := buildSiteSkipRender(site); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	return site
}

func TestRefLinking(t *testing.T) {
	testCommonResetState()

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
	testCommonResetState()

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
			if out, err := site.Info.SourceRelativeLink(link, currentPage); err != nil || out != url {
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

func TestSourceRelativeLinkFileing(t *testing.T) {
	testCommonResetState()

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
			if out, err := site.Info.SourceRelativeLinkFile(link, currentPage); err != nil || out != url {
				t.Errorf("Expected %s to resolve to (%s), got (%s) - error: %s", link, url, out, err)
			} else {
				//t.Logf("tested ok %s maps to %s", link, out)
			}
		}
	}
}
