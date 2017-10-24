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

	"github.com/markbates/inflect"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
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

func pageMust(p *Page, err error) *Page {
	if err != nil {
		panic(err)
	}
	return p
}

func TestRenderWithInvalidTemplate(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "foo.md"), "foo")

	withTemplate := createWithTemplateFromNameValues("missing", templateMissingFunc)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg, WithTemplate: withTemplate}, BuildCfg{})

	errCount := s.Log.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError)

	// TODO(bep) clean up the template error handling
	// The template errors are stored in a slice etc. so we get 4 log entries
	// When we should get only 1
	if errCount == 0 {
		t.Fatalf("Expecting the template to log 1 ERROR, got %d", errCount)
	}
}

func TestDraftAndFutureRender(t *testing.T) {
	t.Parallel()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.md"), Content: []byte("---\ntitle: doc1\ndraft: true\npublishdate: \"2414-05-29\"\n---\n# doc1\n*some content*")},
		{Name: filepath.FromSlash("sect/doc2.md"), Content: []byte("---\ntitle: doc2\ndraft: true\npublishdate: \"2012-05-29\"\n---\n# doc2\n*some content*")},
		{Name: filepath.FromSlash("sect/doc3.md"), Content: []byte("---\ntitle: doc3\ndraft: false\npublishdate: \"2414-05-29\"\n---\n# doc3\n*some content*")},
		{Name: filepath.FromSlash("sect/doc4.md"), Content: []byte("---\ntitle: doc4\ndraft: false\npublishdate: \"2012-05-29\"\n---\n# doc4\n*some content*")},
	}

	siteSetup := func(t *testing.T, configKeyValues ...interface{}) *Site {
		cfg, fs := newTestCfg()

		cfg.Set("baseURL", "http://auth/bub")

		for i := 0; i < len(configKeyValues); i += 2 {
			cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
		}

		for _, src := range sources {
			writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

		}

		return buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
	}

	// Testing Defaults.. Only draft:true and publishDate in the past should be rendered
	s := siteSetup(t)
	if len(s.RegularPages) != 1 {
		t.Fatal("Draft or Future dated content published unexpectedly")
	}

	// only publishDate in the past should be rendered
	s = siteSetup(t, "buildDrafts", true)
	if len(s.RegularPages) != 2 {
		t.Fatal("Future Dated Posts published unexpectedly")
	}

	//  drafts should not be rendered, but all dates should
	s = siteSetup(t,
		"buildDrafts", false,
		"buildFuture", true)

	if len(s.RegularPages) != 2 {
		t.Fatal("Draft posts published unexpectedly")
	}

	// all 4 should be included
	s = siteSetup(t,
		"buildDrafts", true,
		"buildFuture", true)

	if len(s.RegularPages) != 4 {
		t.Fatal("Drafts or Future posts not included as expected")
	}

}

func TestFutureExpirationRender(t *testing.T) {
	t.Parallel()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc3.md"), Content: []byte("---\ntitle: doc1\nexpirydate: \"2400-05-29\"\n---\n# doc1\n*some content*")},
		{Name: filepath.FromSlash("sect/doc4.md"), Content: []byte("---\ntitle: doc2\nexpirydate: \"2000-05-29\"\n---\n# doc2\n*some content*")},
	}

	siteSetup := func(t *testing.T) *Site {
		cfg, fs := newTestCfg()
		cfg.Set("baseURL", "http://auth/bub")

		for _, src := range sources {
			writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

		}

		return buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
	}

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

func TestLastChange(t *testing.T) {
	t.Parallel()

	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "sect/doc1.md"), "---\ntitle: doc1\nweight: 1\ndate: 2014-05-29\n---\n# doc1\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc2.md"), "---\ntitle: doc2\nweight: 2\ndate: 2015-05-29\n---\n# doc2\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc3.md"), "---\ntitle: doc3\nweight: 3\ndate: 2017-05-29\n---\n# doc3\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc4.md"), "---\ntitle: doc4\nweight: 4\ndate: 2016-05-29\n---\n# doc4\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc5.md"), "---\ntitle: doc5\nweight: 3\n---\n# doc5\n*some content*")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.False(t, s.Info.LastChange.IsZero(), "Site.LastChange is zero")
	require.Equal(t, 2017, s.Info.LastChange.Year(), "Site.LastChange should be set to the page with latest Lastmod (year 2017)")
}

// Issue #_index
func TestPageWithUnderScoreIndexInFilename(t *testing.T) {
	t.Parallel()

	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "sect/my_index_file.md"), "---\ntitle: doc1\nweight: 1\ndate: 2014-05-29\n---\n# doc1\n*some content*")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 1)

}

// Issue #957
func TestCrossrefs(t *testing.T) {
	t.Parallel()
	for _, uglyURLs := range []bool{true, false} {
		for _, relative := range []bool{true, false} {
			doTestCrossrefs(t, relative, uglyURLs)
		}
	}
}

func doTestCrossrefs(t *testing.T, relative, uglyURLs bool) {

	baseURL := "http://foo/bar"

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

	doc3Slashed := filepath.FromSlash("/sect/doc3.md")

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
		// Issue #3703
		{
			Name:    filepath.FromSlash("sect/doc4.md"),
			Content: []byte(fmt.Sprintf(`**Ref 1:**{{< %s "%s" >}}.`, refShortcode, doc3Slashed)),
		},
	}

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", baseURL)
	cfg.Set("uglyURLs", uglyURLs)
	cfg.Set("verbose", true)

	for _, src := range sources {
		writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))
	}

	s := buildSingleSite(
		t,
		deps.DepsCfg{
			Fs:           fs,
			Cfg:          cfg,
			WithTemplate: createWithTemplateFromNameValues("_default/single.html", "{{.Content}}")},
		BuildCfg{})

	require.Len(t, s.RegularPages, 4)

	th := testHelper{s.Cfg, s.Fs, t}

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc1%s", expectedPathSuffix)), fmt.Sprintf("<p>Ref 2: %s/sect/doc2%s</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc2%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong></p>\n\n%s/sect/doc1%s\n\n<p>THE END.</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc3%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong>%s/sect/doc3%s.</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc4%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong>%s/sect/doc3%s.</p>\n", expectedBase, expectedURLSuffix)},
	}

	for _, test := range tests {
		th.assertFileContent(test.doc, test.expected)

	}

}

// Issue #939
// Issue #1923
func TestShouldAlwaysHaveUglyURLs(t *testing.T) {
	t.Parallel()
	for _, uglyURLs := range []bool{true, false} {
		doTestShouldAlwaysHaveUglyURLs(t, uglyURLs)
	}
}

func doTestShouldAlwaysHaveUglyURLs(t *testing.T, uglyURLs bool) {

	cfg, fs := newTestCfg()

	cfg.Set("verbose", true)
	cfg.Set("baseURL", "http://auth/bub")
	cfg.Set("disableSitemap", false)
	cfg.Set("disableRSS", false)
	cfg.Set("rssURI", "index.xml")
	cfg.Set("blackfriday",
		map[string]interface{}{
			"plainIDAnchors": true})

	cfg.Set("uglyURLs", uglyURLs)

	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.md"), Content: []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
		{Name: filepath.FromSlash("sect/doc2.md"), Content: []byte("---\nurl: /ugly.html\nmarkup: markdown\n---\n# title\ndoc2 *content*")},
	}

	for _, src := range sources {
		writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))
	}

	writeSource(t, fs, filepath.Join("layouts", "index.html"), "Home Sweet {{ if.IsHome  }}Home{{ end }}.")
	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "{{.Content}}{{ if.IsHome  }}This is not home!{{ end }}")
	writeSource(t, fs, filepath.Join("layouts", "404.html"), "Page Not Found.{{ if.IsHome  }}This is not home!{{ end }}")
	writeSource(t, fs, filepath.Join("layouts", "rss.xml"), "<root>RSS</root>")
	writeSource(t, fs, filepath.Join("layouts", "sitemap.xml"), "<root>SITEMAP</root>")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

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
		content := readDestination(t, fs, test.doc)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}

func TestNewSiteDefaultLang(t *testing.T) {
	t.Parallel()
	s, err := NewSiteDefaultLang()
	require.NoError(t, err)
	require.Equal(t, hugofs.Os, s.Fs.Source)
	require.Equal(t, hugofs.Os, s.Fs.Destination)
}

// Issue #3355
func TestShouldNotWriteZeroLengthFilesToDestination(t *testing.T) {
	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "simple.html"), "simple")
	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "{{.Content}}")
	writeSource(t, fs, filepath.Join("layouts", "_default/list.html"), "")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
	th := testHelper{s.Cfg, s.Fs, t}

	th.assertFileNotExist(filepath.Join("public", "index.html"))
}

// Issue #1176
func TestSectionNaming(t *testing.T) {
	t.Parallel()
	for _, canonify := range []bool{true, false} {
		for _, uglify := range []bool{true, false} {
			for _, pluralize := range []bool{true, false} {
				doTestSectionNaming(t, canonify, uglify, pluralize)
			}
		}
	}
}

func doTestSectionNaming(t *testing.T, canonify, uglify, pluralize bool) {

	var expectedPathSuffix string

	if uglify {
		expectedPathSuffix = ".html"
	} else {
		expectedPathSuffix = "/index.html"
	}

	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.html"), Content: []byte("doc1")},
		// Add one more page to sect to make sure sect is picked in mainSections
		{Name: filepath.FromSlash("sect/sect.html"), Content: []byte("sect")},
		{Name: filepath.FromSlash("Fish and Chips/doc2.html"), Content: []byte("doc2")},
		{Name: filepath.FromSlash("ラーメン/doc3.html"), Content: []byte("doc3")},
	}

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", "http://auth/sub/")
	cfg.Set("uglyURLs", uglify)
	cfg.Set("pluralizeListTitles", pluralize)
	cfg.Set("canonifyURLs", canonify)

	for _, source := range sources {
		writeSource(t, fs, filepath.Join("content", source.Name), string(source.Content))
	}

	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "{{.Content}}")
	writeSource(t, fs, filepath.Join("layouts", "_default/list.html"), "{{.Title}}")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	mainSections, err := s.Info.Param("mainSections")
	require.NoError(t, err)
	require.Equal(t, []string{"sect"}, mainSections)

	th := testHelper{s.Cfg, s.Fs, t}
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

		th.assertFileContent(filepath.Join("public", test.doc), test.expected)
	}

}
func TestSkipRender(t *testing.T) {
	t.Parallel()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.html"), Content: []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
		{Name: filepath.FromSlash("sect/doc2.html"), Content: []byte("<!doctype html><html><body>more content</body></html>")},
		{Name: filepath.FromSlash("sect/doc3.md"), Content: []byte("# doc3\n*some* content")},
		{Name: filepath.FromSlash("sect/doc4.md"), Content: []byte("---\ntitle: doc4\n---\n# doc4\n*some content*")},
		{Name: filepath.FromSlash("sect/doc5.html"), Content: []byte("<!doctype html><html>{{ template \"head\" }}<body>body5</body></html>")},
		{Name: filepath.FromSlash("sect/doc6.html"), Content: []byte("<!doctype html><html>{{ template \"head_abs\" }}<body>body5</body></html>")},
		{Name: filepath.FromSlash("doc7.html"), Content: []byte("<html><body>doc7 content</body></html>")},
		{Name: filepath.FromSlash("sect/doc8.html"), Content: []byte("---\nmarkup: md\n---\n# title\nsome *content*")},
		// Issue #3021
		{Name: filepath.FromSlash("doc9.html"), Content: []byte("<html><body>doc9: {{< myshortcode >}}</body></html>")},
	}

	cfg, fs := newTestCfg()

	cfg.Set("verbose", true)
	cfg.Set("canonifyURLs", true)
	cfg.Set("uglyURLs", true)
	cfg.Set("baseURL", "http://auth/bub")

	for _, src := range sources {
		writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

	}

	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "{{.Content}}")
	writeSource(t, fs, filepath.Join("layouts", "head"), "<head><script src=\"script.js\"></script></head>")
	writeSource(t, fs, filepath.Join("layouts", "head_abs"), "<head><script src=\"/script.js\"></script></head>")
	writeSource(t, fs, filepath.Join("layouts", "shortcodes", "myshortcode.html"), "SHORT")

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("public/sect/doc1.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("public/sect/doc2.html"), "<!doctype html><html><body>more content</body></html>"},
		{filepath.FromSlash("public/sect/doc3.html"), "\n\n<h1 id=\"doc3\">doc3</h1>\n\n<p><em>some</em> content</p>\n"},
		{filepath.FromSlash("public/sect/doc4.html"), "\n\n<h1 id=\"doc4\">doc4</h1>\n\n<p><em>some content</em></p>\n"},
		{filepath.FromSlash("public/sect/doc5.html"), "<!doctype html><html><head><script src=\"script.js\"></script></head><body>body5</body></html>"},
		{filepath.FromSlash("public/sect/doc6.html"), "<!doctype html><html><head><script src=\"http://auth/bub/script.js\"></script></head><body>body5</body></html>"},
		{filepath.FromSlash("public/doc7.html"), "<html><body>doc7 content</body></html>"},
		{filepath.FromSlash("public/sect/doc8.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("public/doc9.html"), "<html><body>doc9: SHORT</body></html>"},
	}

	for _, test := range tests {
		file, err := fs.Destination.Open(test.doc)
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
	t.Parallel()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.html"), Content: []byte("<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>")},
		{Name: filepath.FromSlash("blue/doc2.html"), Content: []byte("---\nf: t\n---\n<!doctype html><html><body>more content</body></html>")},
	}
	for _, baseURL := range []string{"http://auth/bub", "http://base", "//base"} {
		for _, canonify := range []bool{true, false} {

			cfg, fs := newTestCfg()

			cfg.Set("uglyURLs", true)
			cfg.Set("canonifyURLs", canonify)
			cfg.Set("baseURL", baseURL)

			for _, src := range sources {
				writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

			}

			writeSource(t, fs, filepath.Join("layouts", "blue/single.html"), templateWithURLAbs)

			s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
			th := testHelper{s.Cfg, s.Fs, t}

			tests := []struct {
				file, expected string
			}{
				{"public/blue/doc2.html", "<a href=\"%s/foobar.jpg\">Going</a>"},
				{"public/sect/doc1.html", "<!doctype html><html><head></head><body><a href=\"#frag1\">link</a></body></html>"},
			}

			for _, test := range tests {

				expected := test.expected

				if strings.Contains(expected, "%s") {
					expected = fmt.Sprintf(expected, baseURL)
				}

				if !canonify {
					expected = strings.Replace(expected, baseURL, "", -1)
				}

				th.assertFileContent(test.file, expected)

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
	t.Parallel()
	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub")

	for _, src := range weightedSources {
		writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

	}

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	if s.getPage(KindSection, "sect").Pages[1].Title != "Three" || s.getPage(KindSection, "sect").Pages[2].Title != "Four" {
		t.Error("Pages in unexpected order.")
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
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub")

	writeSourcesToSource(t, "content", fs, groupedSources...)
	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

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
	t.Parallel()
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("sect/doc1.md"), Content: pageWithWeightedTaxonomies2},
		{Name: filepath.FromSlash("sect/doc2.md"), Content: pageWithWeightedTaxonomies1},
		{Name: filepath.FromSlash("sect/doc3.md"), Content: pageWithWeightedTaxonomies3},
	}
	taxonomies := make(map[string]string)

	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", "http://auth/bub")
	cfg.Set("taxonomies", taxonomies)

	writeSourcesToSource(t, "content", fs, sources...)
	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

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
	sp := source.NewSourceSpec(site.Cfg, site.Fs)
	currentPath := sp.NewFile(filepath.FromSlash(f))
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
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("level2/unique.md"), Content: []byte("")},
		{Name: filepath.FromSlash("index.md"), Content: []byte("")},
		{Name: filepath.FromSlash("rootfile.md"), Content: []byte("")},
		{Name: filepath.FromSlash("root-image.png"), Content: []byte("")},

		{Name: filepath.FromSlash("level2/2-root.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/index.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/common.md"), Content: []byte("")},

		{Name: filepath.FromSlash("level2/2-image.png"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/common.png"), Content: []byte("")},

		{Name: filepath.FromSlash("level2/level3/3-root.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/index.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/common.md"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/3-image.png"), Content: []byte("")},
		{Name: filepath.FromSlash("level2/level3/common.png"), Content: []byte("")},
	}

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", "http://auth/")
	cfg.Set("uglyURLs", false)
	cfg.Set("outputs", map[string]interface{}{
		"page": []string{"HTML", "AMP"},
	})
	cfg.Set("pluralizeListTitles", false)
	cfg.Set("canonifyURLs", false)
	cfg.Set("blackfriday",
		map[string]interface{}{})
	writeSourcesToSource(t, "content", fs, sources...)
	return buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

}

func TestRefLinking(t *testing.T) {
	t.Parallel()
	site := setupLinkingMockSite(t)

	currentPage := findPage(site, "level2/level3/index.md")
	if currentPage == nil {
		t.Fatalf("failed to find current page in site")
	}

	for i, test := range []struct {
		link         string
		outputFormat string
		relative     bool
		expected     string
	}{
		{"unique.md", "", true, "/level2/unique/"},
		{"level2/common.md", "", true, "/level2/common/"},
		{"3-root.md", "", true, "/level2/level3/3-root/"},
		{"level2/level3/index.md", "amp", true, "/amp/level2/level3/"},
		{"level2/index.md", "amp", false, "http://auth/amp/level2/"},
	} {
		if out, err := site.Info.refLink(test.link, currentPage, test.relative, test.outputFormat); err != nil || out != test.expected {
			t.Errorf("[%d] Expected %s to resolve to (%s), got (%s) - error: %s", i, test.link, test.expected, out, err)
		}
	}

	// TODO: and then the failure cases.
}
