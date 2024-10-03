// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/flect"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/publisher"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

func TestDraftAndFutureRender(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	sources := [][2]string{
		{filepath.FromSlash("sect/doc1.md"), "---\ntitle: doc1\ndraft: true\npublishdate: \"2414-05-29\"\n---\n# doc1\n*some content*"},
		{filepath.FromSlash("sect/doc2.md"), "---\ntitle: doc2\ndraft: true\npublishdate: \"2012-05-29\"\n---\n# doc2\n*some content*"},
		{filepath.FromSlash("sect/doc3.md"), "---\ntitle: doc3\ndraft: false\npublishdate: \"2414-05-29\"\n---\n# doc3\n*some content*"},
		{filepath.FromSlash("sect/doc4.md"), "---\ntitle: doc4\ndraft: false\npublishdate: \"2012-05-29\"\n---\n# doc4\n*some content*"},
	}

	siteSetup := func(t *testing.T, configKeyValues ...any) *Site {
		cfg, fs := newTestCfg()

		cfg.Set("baseURL", "http://auth/bub")

		for i := 0; i < len(configKeyValues); i += 2 {
			cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
		}
		configs, err := loadTestConfigFromProvider(cfg)
		c.Assert(err, qt.IsNil)

		for _, src := range sources {
			writeSource(t, fs, filepath.Join("content", src[0]), src[1])
		}

		return buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})
	}

	// Testing Defaults.. Only draft:true and publishDate in the past should be rendered
	s := siteSetup(t)
	if len(s.RegularPages()) != 1 {
		t.Fatal("Draft or Future dated content published unexpectedly")
	}

	// only publishDate in the past should be rendered
	s = siteSetup(t, "buildDrafts", true)
	if len(s.RegularPages()) != 2 {
		t.Fatal("Future Dated Posts published unexpectedly")
	}

	//  drafts should not be rendered, but all dates should
	s = siteSetup(t,
		"buildDrafts", false,
		"buildFuture", true)

	if len(s.RegularPages()) != 2 {
		t.Fatal("Draft posts published unexpectedly")
	}

	// all 4 should be included
	s = siteSetup(t,
		"buildDrafts", true,
		"buildFuture", true)

	if len(s.RegularPages()) != 4 {
		t.Fatal("Drafts or Future posts not included as expected")
	}
}

func TestFutureExpirationRender(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	sources := [][2]string{
		{filepath.FromSlash("sect/doc3.md"), "---\ntitle: doc1\nexpirydate: \"2400-05-29\"\n---\n# doc1\n*some content*"},
		{filepath.FromSlash("sect/doc4.md"), "---\ntitle: doc2\nexpirydate: \"2000-05-29\"\n---\n# doc2\n*some content*"},
	}

	siteSetup := func(t *testing.T) *Site {
		cfg, fs := newTestCfg()
		cfg.Set("baseURL", "http://auth/bub")

		configs, err := loadTestConfigFromProvider(cfg)
		c.Assert(err, qt.IsNil)

		for _, src := range sources {
			writeSource(t, fs, filepath.Join("content", src[0]), src[1])
		}

		return buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})
	}

	s := siteSetup(t)

	if len(s.AllPages()) != 1 {
		if len(s.RegularPages()) > 1 {
			t.Fatal("Expired content published unexpectedly")
		}

		if len(s.RegularPages()) < 1 {
			t.Fatal("Valid content expired unexpectedly")
		}
	}

	if s.AllPages()[0].Title() == "doc2" {
		t.Fatal("Expired content published unexpectedly")
	}
}

func TestLastChange(t *testing.T) {
	t.Parallel()

	cfg, fs := newTestCfg()
	c := qt.New(t)
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "sect/doc1.md"), "---\ntitle: doc1\nweight: 1\ndate: 2014-05-29\n---\n# doc1\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc2.md"), "---\ntitle: doc2\nweight: 2\ndate: 2015-05-29\n---\n# doc2\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc3.md"), "---\ntitle: doc3\nweight: 3\ndate: 2017-05-29\n---\n# doc3\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc4.md"), "---\ntitle: doc4\nweight: 4\ndate: 2016-05-29\n---\n# doc4\n*some content*")
	writeSource(t, fs, filepath.Join("content", "sect/doc5.md"), "---\ntitle: doc5\nweight: 3\n---\n# doc5\n*some content*")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

	c.Assert(s.Lastmod().IsZero(), qt.Equals, false)
	c.Assert(s.Lastmod().Year(), qt.Equals, 2017)
}

// Issue #_index
func TestPageWithUnderScoreIndexInFilename(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	cfg, fs := newTestCfg()
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "sect/my_index_file.md"), "---\ntitle: doc1\nweight: 1\ndate: 2014-05-29\n---\n# doc1\n*some content*")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)
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
	c := qt.New(t)

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

	sources := [][2]string{
		{
			filepath.FromSlash("sect/doc1.md"),
			fmt.Sprintf(`Ref 2: {{< %s "sect/doc2.md" >}}`, refShortcode),
		},
		// Issue #1148: Make sure that no P-tags is added around shortcodes.
		{
			filepath.FromSlash("sect/doc2.md"),
			fmt.Sprintf(`**Ref 1:**

{{< %s "sect/doc1.md" >}}

THE END.`, refShortcode),
		},
		// Issue #1753: Should not add a trailing newline after shortcode.
		{
			filepath.FromSlash("sect/doc3.md"),
			fmt.Sprintf(`**Ref 1:** {{< %s "sect/doc3.md" >}}.`, refShortcode),
		},
		// Issue #3703
		{
			filepath.FromSlash("sect/doc4.md"),
			fmt.Sprintf(`**Ref 1:** {{< %s "%s" >}}.`, refShortcode, doc3Slashed),
		},
	}

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", baseURL)
	cfg.Set("uglyURLs", uglyURLs)
	cfg.Set("verbose", true)
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	for _, src := range sources {
		writeSource(t, fs, filepath.Join("content", src[0]), src[1])
	}
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), "{{.Content}}")

	s := buildSingleSite(
		t,
		deps.DepsCfg{
			Fs:      fs,
			Configs: configs,
		},
		BuildCfg{})

	c.Assert(len(s.RegularPages()), qt.Equals, 4)

	th := newTestHelper(s.conf, s.Fs, t)

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc1%s", expectedPathSuffix)), fmt.Sprintf("<p>Ref 2: %s/sect/doc2%s</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc2%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong></p>\n%s/sect/doc1%s\n<p>THE END.</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc3%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong> %s/sect/doc3%s.</p>\n", expectedBase, expectedURLSuffix)},
		{filepath.FromSlash(fmt.Sprintf("public/sect/doc4%s", expectedPathSuffix)), fmt.Sprintf("<p><strong>Ref 1:</strong> %s/sect/doc3%s.</p>\n", expectedBase, expectedURLSuffix)},
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
	c := qt.New(t)

	cfg.Set("verbose", true)
	cfg.Set("baseURL", "http://auth/bub")
	cfg.Set("uglyURLs", uglyURLs)

	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	sources := [][2]string{
		{filepath.FromSlash("sect/doc1.md"), "---\nmarkup: markdown\n---\n# title\nsome *content*"},
		{filepath.FromSlash("sect/doc2.md"), "---\nurl: /ugly.html\nmarkup: markdown\n---\n# title\ndoc2 *content*"},
	}

	for _, src := range sources {
		writeSource(t, fs, filepath.Join("content", src[0]), src[1])
	}

	writeSource(t, fs, filepath.Join("layouts", "index.html"), "Home Sweet {{ if.IsHome  }}Home{{ end }}.")
	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "{{.Content}}{{ if.IsHome  }}This is not home!{{ end }}")
	writeSource(t, fs, filepath.Join("layouts", "404.html"), "Page Not Found.{{ if.IsHome  }}This is not home!{{ end }}")
	writeSource(t, fs, filepath.Join("layouts", "rss.xml"), "<root>RSS</root>")
	writeSource(t, fs, filepath.Join("layouts", "sitemap.xml"), "<root>SITEMAP</root>")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

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
		{filepath.FromSlash(expectedPagePath), "<h1 id=\"title\">title</h1>\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("public/404.html"), "Page Not Found."},
		{filepath.FromSlash("public/index.xml"), "<root>RSS</root>"},
		{filepath.FromSlash("public/sitemap.xml"), "<root>SITEMAP</root>"},
		// Issue #1923
		{filepath.FromSlash("public/ugly.html"), "<h1 id=\"title\">title</h1>\n<p>doc2 <em>content</em></p>\n"},
	}

	for _, p := range s.RegularPages() {
		c.Assert(p.IsHome(), qt.Equals, false)
	}

	for _, test := range tests {
		content := readWorkingDir(t, fs, test.doc)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}
}

// Issue #3355
func TestShouldNotWriteZeroLengthFilesToDestination(t *testing.T) {
	c := qt.New(t)

	cfg, fs := newTestCfg()
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "simple.html"), "simple")
	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "{{.Content}}")
	writeSource(t, fs, filepath.Join("layouts", "_default/list.html"), "")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})
	th := newTestHelper(s.conf, s.Fs, t)

	th.assertFileNotExist(filepath.Join("public", "index.html"))
}

func TestMainSections(t *testing.T) {
	c := qt.New(t)
	for _, paramSet := range []bool{false, true} {
		c.Run(fmt.Sprintf("param-%t", paramSet), func(c *qt.C) {
			v := config.New()
			if paramSet {
				v.Set("params", map[string]any{
					"mainSections": []string{"a1", "a2"},
				})
			}

			b := newTestSitesBuilder(c).WithViper(v)

			for i := 0; i < 20; i++ {
				b.WithContent(fmt.Sprintf("page%d.md", i), `---
title: "Page"
---
`)
			}

			for i := 0; i < 5; i++ {
				b.WithContent(fmt.Sprintf("blog/page%d.md", i), `---
title: "Page"
tags: ["a", "b", "c", "d", "e", "f", "g", "h", "i", "j"]
---
`)
			}

			for i := 0; i < 3; i++ {
				b.WithContent(fmt.Sprintf("docs/page%d.md", i), `---
title: "Page"
---
`)
			}

			b.WithTemplates("index.html", `
mainSections: {{ .Site.Params.mainSections }}

{{ range (where .Site.RegularPages "Type" "in" .Site.Params.mainSections) }}
Main section page: {{ .RelPermalink }}
{{ end }}
`)

			b.Build(BuildCfg{})

			if paramSet {
				b.AssertFileContent("public/index.html", "mainSections: [a1 a2]")
			} else {
				b.AssertFileContent("public/index.html", "mainSections: [blog]", "Main section page: /blog/page3/")
			}
		})
	}
}

func TestMainSectionsMoveToSite(t *testing.T) {
	t.Run("defined in params", func(t *testing.T) {
		t.Parallel()

		files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
[params]
mainSections=["a", "b"]
-- content/mysect/page1.md --
-- layouts/index.html --
{{/* Behaviour before Hugo 0.112.0. */}}
MainSections Params: {{ site.Params.mainSections }}|
MainSections Site method: {{ site.MainSections }}|


	`

		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
MainSections Params: [a b]|
MainSections Site method: [a b]|
	`)
	})

	t.Run("defined in top level config", func(t *testing.T) {
		t.Parallel()

		files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
mainSections=["a", "b"]
[params]
[params.sub]
mainSections=["c", "d"]
-- content/mysect/page1.md --
-- layouts/index.html --
{{/* Behaviour before Hugo 0.112.0. */}}
MainSections Params: {{ site.Params.mainSections }}|
MainSections Param sub: {{ site.Params.sub.mainSections }}|
MainSections Site method: {{ site.MainSections }}|


`

		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
MainSections Params: [a b]|
MainSections Param sub: [c d]|
MainSections Site method: [a b]|
`)
	})

	t.Run("guessed from pages", func(t *testing.T) {
		t.Parallel()

		files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
-- content/mysect/page1.md --
-- layouts/index.html --
MainSections Params: {{ site.Params.mainSections }}|
MainSections Site method: {{ site.MainSections }}|


	`

		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
MainSections Params: [mysect]|
MainSections Site method: [mysect]|
	`)
	})
}

// Issue #1176
func TestSectionNaming(t *testing.T) {
	for _, canonify := range []bool{true, false} {
		for _, uglify := range []bool{true, false} {
			for _, pluralize := range []bool{true, false} {
				canonify := canonify
				uglify := uglify
				pluralize := pluralize
				t.Run(fmt.Sprintf("canonify=%t,uglify=%t,pluralize=%t", canonify, uglify, pluralize), func(t *testing.T) {
					t.Parallel()
					doTestSectionNaming(t, canonify, uglify, pluralize)
				})
			}
		}
	}
}

func doTestSectionNaming(t *testing.T, canonify, uglify, pluralize bool) {
	c := qt.New(t)

	var expectedPathSuffix string

	if uglify {
		expectedPathSuffix = ".html"
	} else {
		expectedPathSuffix = "/index.html"
	}

	sources := [][2]string{
		{filepath.FromSlash("sect/doc1.html"), "doc1"},
		// Add one more page to sect to make sure sect is picked in mainSections
		{filepath.FromSlash("sect/sect.html"), "sect"},
		{filepath.FromSlash("Fish and Chips/doc2.html"), "doc2"},
		{filepath.FromSlash("ラーメン/doc3.html"), "doc3"},
	}

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", "http://auth/sub/")
	cfg.Set("uglyURLs", uglify)
	cfg.Set("pluralizeListTitles", pluralize)
	cfg.Set("canonifyURLs", canonify)

	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	for _, src := range sources {
		writeSource(t, fs, filepath.Join("content", src[0]), src[1])
	}

	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "{{.Content}}")
	writeSource(t, fs, filepath.Join("layouts", "_default/list.html"), "{{ .Kind }}|{{.Title}}")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

	c.Assert(s.MainSections(), qt.DeepEquals, []string{"sect"})

	th := newTestHelper(s.conf, s.Fs, t)
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
			test.expected = flect.Pluralize(test.expected)
		}

		th.assertFileContent(filepath.Join("public", test.doc), test.expected)
	}
}

var weightedPage1 = `+++
weight = "2"
title = "One"
my_param = "foo"
my_date = 1979-05-27T07:32:00Z
+++
Front Matter with Ordered Pages`

var weightedPage2 = `+++
weight = "6"
title = "Two"
publishdate = "2012-03-05"
my_param = "foo"
+++
Front Matter with Ordered Pages 2`

var weightedPage3 = `+++
weight = "4"
title = "Three"
date = "2012-04-06"
publishdate = "2012-04-06"
my_param = "bar"
only_one = "yes"
my_date = 2010-05-27T07:32:00Z
+++
Front Matter with Ordered Pages 3`

var weightedPage4 = `+++
weight = "4"
title = "Four"
date = "2012-01-01"
publishdate = "2012-01-01"
my_param = "baz"
my_date = 2010-05-27T07:32:00Z
summary = "A _custom_ summary"
categories = [ "hugo" ]
+++
Front Matter with Ordered Pages 4. This is longer content`

var weightedPage5 = `+++
weight = "5"
title = "Five"

[_build]
render = "never"
+++
Front Matter with Ordered Pages 5`

var weightedSources = [][2]string{
	{filepath.FromSlash("sect/doc1.md"), weightedPage1},
	{filepath.FromSlash("sect/doc2.md"), weightedPage2},
	{filepath.FromSlash("sect/doc3.md"), weightedPage3},
	{filepath.FromSlash("sect/doc4.md"), weightedPage4},
	{filepath.FromSlash("sect/doc5.md"), weightedPage5},
}

func TestOrderedPages(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub")
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	for _, src := range weightedSources {
		writeSource(t, fs, filepath.Join("content", src[0]), src[1])
	}

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

	if s.getPageOldVersion(kinds.KindSection, "sect").Pages()[1].Title() != "Three" || s.getPageOldVersion(kinds.KindSection, "sect").Pages()[2].Title() != "Four" {
		t.Error("Pages in unexpected order.")
	}

	bydate := s.RegularPages().ByDate()

	if bydate[0].Title() != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bydate[0].Title())
	}

	rev := bydate.Reverse()
	if rev[0].Title() != "Three" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Three", rev[0].Title())
	}

	bypubdate := s.RegularPages().ByPublishDate()

	if bypubdate[0].Title() != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bypubdate[0].Title())
	}

	rbypubdate := bypubdate.Reverse()
	if rbypubdate[0].Title() != "Three" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Three", rbypubdate[0].Title())
	}

	bylength := s.RegularPages().ByLength(context.Background())
	if bylength[0].Title() != "One" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "One", bylength[0].Title())
	}

	rbylength := bylength.Reverse()
	if rbylength[0].Title() != "Four" {
		t.Errorf("Pages in unexpected order. First should be '%s', got '%s'", "Four", rbylength[0].Title())
	}
}

var groupedSources = [][2]string{
	{filepath.FromSlash("sect1/doc1.md"), weightedPage1},
	{filepath.FromSlash("sect1/doc2.md"), weightedPage2},
	{filepath.FromSlash("sect2/doc3.md"), weightedPage3},
	{filepath.FromSlash("sect3/doc4.md"), weightedPage4},
}

func TestGroupedPages(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub")
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSourcesToSource(t, "content", fs, groupedSources...)
	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

	rbysection, err := s.RegularPages().GroupBy(context.Background(), "Section", "desc")
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
	if rbysection[0].Pages[0].Title() != "Four" {
		t.Errorf("PageGroup has an unexpected page. First group's pages should have '%s', got '%s'", "Four", rbysection[0].Pages[0].Title())
	}
	if len(rbysection[2].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. Third group should have '%d' pages, got '%d' pages", 2, len(rbysection[2].Pages))
	}

	bytype, err := s.RegularPages().GroupBy(context.Background(), "Type", "asc")
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
	if bytype[2].Pages[0].Title() != "Four" {
		t.Errorf("PageGroup has an unexpected page. Third group's data should have '%s', got '%s'", "Four", bytype[0].Pages[0].Title())
	}
	if len(bytype[0].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 2, len(bytype[2].Pages))
	}

	bydate, err := s.RegularPages().GroupByDate("2006-01", "asc")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if bydate[0].Key != "0001-01" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "0001-01", bydate[0].Key)
	}
	if bydate[1].Key != "2012-01" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "2012-01", bydate[1].Key)
	}

	bypubdate, err := s.RegularPages().GroupByPublishDate("2006")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if bypubdate[0].Key != "2012" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "2012", bypubdate[0].Key)
	}
	if bypubdate[1].Key != "0001" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "0001", bypubdate[1].Key)
	}
	if bypubdate[0].Pages[0].Title() != "Three" {
		t.Errorf("PageGroup has an unexpected page. Third group's pages should have '%s', got '%s'", "Three", bypubdate[0].Pages[0].Title())
	}
	if len(bypubdate[0].Pages) != 3 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 3, len(bypubdate[0].Pages))
	}

	byparam, err := s.RegularPages().GroupByParam("my_param", "desc")
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
	if byparam[2].Pages[0].Title() != "Three" {
		t.Errorf("PageGroup has an unexpected page. Third group's pages should have '%s', got '%s'", "Three", byparam[2].Pages[0].Title())
	}
	if len(byparam[0].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 2, len(byparam[0].Pages))
	}

	byNonExistentParam, err := s.RegularPages().GroupByParam("not_exist")
	if err != nil {
		t.Errorf("GroupByParam returned an error when it shouldn't")
	}
	if len(byNonExistentParam) != 0 {
		t.Errorf("PageGroup array has unexpected elements. Group length should be '%d', got '%d'", 0, len(byNonExistentParam))
	}

	byOnlyOneParam, err := s.RegularPages().GroupByParam("only_one")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if len(byOnlyOneParam) != 1 {
		t.Errorf("PageGroup array has unexpected elements. Group length should be '%d', got '%d'", 1, len(byOnlyOneParam))
	}
	if byOnlyOneParam[0].Key != "yes" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "yes", byOnlyOneParam[0].Key)
	}

	byParamDate, err := s.RegularPages().GroupByParamDate("my_date", "2006-01")
	if err != nil {
		t.Fatalf("Unable to make PageGroup array: %s", err)
	}
	if byParamDate[0].Key != "2010-05" {
		t.Errorf("PageGroup array in unexpected order. First group key should be '%s', got '%s'", "2010-05", byParamDate[0].Key)
	}
	if byParamDate[1].Key != "1979-05" {
		t.Errorf("PageGroup array in unexpected order. Second group key should be '%s', got '%s'", "1979-05", byParamDate[1].Key)
	}
	if byParamDate[1].Pages[0].Title() != "One" {
		t.Errorf("PageGroup has an unexpected page. Second group's pages should have '%s', got '%s'", "One", byParamDate[1].Pages[0].Title())
	}
	if len(byParamDate[0].Pages) != 2 {
		t.Errorf("PageGroup has unexpected number of pages. First group should have '%d' pages, got '%d' pages", 2, len(byParamDate[2].Pages))
	}
}

var pageWithWeightedTaxonomies1 = `+++
tags = [ "a", "b", "c" ]
tags_weight = 22
categories = ["d"]
title = "foo"
categories_weight = 44
+++
Front Matter with weighted tags and categories`

var pageWithWeightedTaxonomies2 = `+++
tags = "a"
tags_weight = 33
title = "bar"
categories = [ "d", "e" ]
categories_weight = 11.0
alias = "spf13"
date = 1979-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`

var pageWithWeightedTaxonomies3 = `+++
title = "bza"
categories = [ "e" ]
categories_weight = 11
alias = "spf13"
date = 2010-05-27T07:32:00Z
+++
Front Matter with weighted tags and categories`

func TestWeightedTaxonomies(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	sources := [][2]string{
		{filepath.FromSlash("sect/doc1.md"), pageWithWeightedTaxonomies2},
		{filepath.FromSlash("sect/doc2.md"), pageWithWeightedTaxonomies1},
		{filepath.FromSlash("sect/doc3.md"), pageWithWeightedTaxonomies3},
	}
	taxonomies := make(map[string]string)

	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", "http://auth/bub")
	cfg.Set("taxonomies", taxonomies)
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSourcesToSource(t, "content", fs, sources...)
	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

	if s.Taxonomies()["tags"]["a"][0].Page.Title() != "foo" {
		t.Errorf("Pages in unexpected order, 'foo' expected first, got '%v'", s.Taxonomies()["tags"]["a"][0].Page.Title())
	}

	if s.Taxonomies()["categories"]["d"][0].Page.Title() != "bar" {
		t.Errorf("Pages in unexpected order, 'bar' expected first, got '%v'", s.Taxonomies()["categories"]["d"][0].Page.Title())
	}

	if s.Taxonomies()["categories"]["e"][0].Page.Title() != "bza" {
		t.Errorf("Pages in unexpected order, 'bza' expected first, got '%v'", s.Taxonomies()["categories"]["e"][0].Page.Title())
	}
}

func setupLinkingMockSite(t *testing.T) *Site {
	sources := [][2]string{
		{filepath.FromSlash("level2/unique.md"), ""},
		{filepath.FromSlash("_index.md"), ""},
		{filepath.FromSlash("common.md"), ""},
		{filepath.FromSlash("rootfile.md"), ""},
		{filepath.FromSlash("root-image.png"), ""},

		{filepath.FromSlash("level2/2-root.md"), ""},
		{filepath.FromSlash("level2/common.md"), ""},

		{filepath.FromSlash("level2/2-image.png"), ""},
		{filepath.FromSlash("level2/common.png"), ""},

		{filepath.FromSlash("level2/level3/start.md"), ""},
		{filepath.FromSlash("level2/level3/_index.md"), ""},
		{filepath.FromSlash("level2/level3/3-root.md"), ""},
		{filepath.FromSlash("level2/level3/common.md"), ""},
		{filepath.FromSlash("level2/level3/3-image.png"), ""},
		{filepath.FromSlash("level2/level3/common.png"), ""},

		{filepath.FromSlash("level2/level3/embedded.dot.md"), ""},

		{filepath.FromSlash("leafbundle/index.md"), ""},
	}

	cfg, fs := newTestCfg()

	cfg.Set("baseURL", "http://auth/")
	cfg.Set("uglyURLs", false)
	cfg.Set("outputs", map[string]any{
		"page": []string{"HTML", "AMP"},
	})
	cfg.Set("pluralizeListTitles", false)
	cfg.Set("canonifyURLs", false)
	configs, err := loadTestConfigFromProvider(cfg)
	if err != nil {
		t.Fatal(err)
	}

	writeSourcesToSource(t, "content", fs, sources...)
	return buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})
}

func TestRefLinking(t *testing.T) {
	t.Parallel()
	site := setupLinkingMockSite(t)

	currentPage := site.getPageOldVersion(kinds.KindPage, "level2/level3/start.md")
	if currentPage == nil {
		t.Fatalf("failed to find current page in site")
	}

	for i, test := range []struct {
		link         string
		outputFormat string
		relative     bool
		expected     string
	}{
		// different refs resolving to the same unique filename:
		{"/level2/unique.md", "", true, "/level2/unique/"},
		{"../unique.md", "", true, "/level2/unique/"},
		{"unique.md", "", true, "/level2/unique/"},

		{"level2/common.md", "", true, "/level2/common/"},
		{"3-root.md", "", true, "/level2/level3/3-root/"},
		{"../..", "", true, "/"},

		// different refs resolving to the same ambiguous top-level filename:
		{"../../common.md", "", true, "/common/"},
		{"/common.md", "", true, "/common/"},

		// different refs resolving to the same ambiguous level-2 filename:
		{"/level2/common.md", "", true, "/level2/common/"},
		{"../common.md", "", true, "/level2/common/"},
		{"common.md", "", true, "/level2/level3/common/"},

		// different refs resolving to the same section:
		{"/level2", "", true, "/level2/"},
		{"..", "", true, "/level2/"},
		{"../", "", true, "/level2/"},

		// different refs resolving to the same subsection:
		{"/level2/level3", "", true, "/level2/level3/"},
		{"/level2/level3/_index.md", "", true, "/level2/level3/"},
		{".", "", true, "/level2/level3/"},
		{"./", "", true, "/level2/level3/"},

		// try to confuse parsing
		{"embedded.dot.md", "", true, "/level2/level3/embedded.dot/"},

		// test empty link, as well as fragment only link
		{"", "", true, ""},
	} {
		t.Run(fmt.Sprintf("t%dt", i), func(t *testing.T) {
			checkLinkCase(site, test.link, currentPage, test.relative, test.outputFormat, test.expected, t, i)

			// make sure fragment links are also handled
			checkLinkCase(site, test.link+"#intro", currentPage, test.relative, test.outputFormat, test.expected+"#intro", t, i)
		})
	}

	// TODO: and then the failure cases.
}

func TestRelRefWithTrailingSlash(t *testing.T) {
	files := `
-- hugo.toml --
-- content/docs/5.3/examples/_index.md --
---
title: "Examples"
---
-- content/_index.md --
---
title: "Home"
---

Examples: {{< relref "/docs/5.3/examples/" >}}
-- layouts/home.html --
Content: {{ .Content }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Examples: /docs/5.3/examples/")
}

func checkLinkCase(site *Site, link string, currentPage page.Page, relative bool, outputFormat string, expected string, t *testing.T, i int) {
	t.Helper()
	if out, err := site.refLink(link, currentPage, relative, outputFormat); err != nil || out != expected {
		t.Fatalf("[%d] Expected %q from %q to resolve to %q, got %q - error: %s", i, link, currentPage.Path(), expected, out, err)
	}
}

// https://github.com/gohugoio/hugo/issues/6952
func TestRefIssues(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithContent(
		"post/b1/index.md", "---\ntitle: pb1\n---\nRef: {{< ref \"b2\" >}}",
		"post/b2/index.md", "---\ntitle: pb2\n---\n",
		"post/nested-a/content-a.md", "---\ntitle: ca\n---\n{{< ref \"content-b\" >}}",
		"post/nested-b/content-b.md", "---\ntitle: ca\n---\n",
	)
	b.WithTemplates("index.html", `Home`)
	b.WithTemplates("_default/single.html", `Content: {{ .Content }}`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/post/b1/index.html", `Content: <p>Ref: http://example.com/post/b2/</p>`)
	b.AssertFileContent("public/post/nested-a/content-a/index.html", `Content: http://example.com/post/nested-b/content-b/`)
}

func TestClassCollector(t *testing.T) {
	for _, minify := range []bool{false, true} {
		t.Run(fmt.Sprintf("minify-%t", minify), func(t *testing.T) {
			statsFilename := "hugo_stats.json"
			defer os.Remove(statsFilename)

			b := newTestSitesBuilder(t)
			b.WithConfigFile("toml", fmt.Sprintf(`


minify = %t

[build]
  writeStats = true

`, minify))

			b.WithTemplates("index.html", `

<div id="el1" class="a b c">Foo</div>

Some text.

<div class="c d e [&>p]:text-red-600" id="el2">Foo</div>

<span class=z>FOO</span>

 <a class="text-base hover:text-gradient inline-block px-3 pb-1 rounded lowercase" href="{{ .RelPermalink }}">{{ .Title }}</a>


`)

			b.WithContent("p1.md", "")

			b.Build(BuildCfg{})

			b.AssertFileContent("hugo_stats.json", `
 {
          "htmlElements": {
            "tags": [
              "a",
              "div",
              "span"
            ],
            "classes": [
              "a",
              "b",
              "c",
              "d",
              "e",
              "hover:text-gradient",
			  "[&>p]:text-red-600",
              "inline-block",
              "lowercase",
              "pb-1",
              "px-3",
              "rounded",
              "text-base",
              "z"
            ],
            "ids": [
              "el1",
              "el2"
            ]
          }
        }
`)
		})
	}
}

func TestClassCollectorConfigWriteStats(t *testing.T) {
	r := func(writeStatsConfig string) *IntegrationTestBuilder {
		files := `
-- hugo.toml --
WRITE_STATS_CONFIG
-- layouts/_default/list.html --
<div id="myid" class="myclass">Foo</div>

`
		files = strings.Replace(files, "WRITE_STATS_CONFIG", writeStatsConfig, 1)

		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
				NeedsOsFS:   true,
			},
		).Build()

		return b
	}

	// Legacy config.
	b := r(`
[build]
writeStats = true
`)

	b.AssertFileContent("hugo_stats.json", "myclass", "div", "myid")

	b = r(`
[build]
writeStats = false
	`)

	b.AssertFileExists("public/hugo_stats.json", false)

	b = r(`
[build.buildStats]
enable = true
`)

	b.AssertFileContent("hugo_stats.json", "myclass", "div", "myid")

	b = r(`
[build.buildStats]
enable = true
disableids = true
`)

	b.AssertFileContent("hugo_stats.json", "myclass", "div", "! myid")

	b = r(`
[build.buildStats]
enable = true
disableclasses = true
`)

	b.AssertFileContent("hugo_stats.json", "! myclass", "div", "myid")

	b = r(`
[build.buildStats]
enable = true
disabletags = true
	`)

	b.AssertFileContent("hugo_stats.json", "myclass", "! div", "myid")

	b = r(`
[build.buildStats]
enable = true
disabletags = true
disableclasses = true
	`)

	b.AssertFileContent("hugo_stats.json", "! myclass", "! div", "myid")

	b = r(`
[build.buildStats]
enable = false
	`)
	b.AssertFileExists("public/hugo_stats.json", false)
}

func TestClassCollectorStress(t *testing.T) {
	statsFilename := "hugo_stats.json"
	defer os.Remove(statsFilename)

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `

disableKinds = ["home", "section", "term", "taxonomy" ]

[languages]
[languages.en]
[languages.nb]
[languages.no]
[languages.sv]


[build]
  writeStats = true

`)

	b.WithTemplates("_default/single.html", `
<div class="c d e" id="el2">Foo</div>

Some text.

{{ $n := index (shuffle (seq 1 20)) 0 }}

{{ "<span class=_a>Foo</span>" | strings.Repeat $n | safeHTML }}

<div class="{{ .Title }}">
ABC.
</div>

<div class="f"></div>

{{ $n := index (shuffle (seq 1 5)) 0 }}

{{ "<hr class=p-3>" | safeHTML }}

`)

	for _, lang := range []string{"en", "nb", "no", "sv"} {
		for i := 100; i <= 999; i++ {
			b.WithContent(fmt.Sprintf("p%d.%s.md", i, lang), fmt.Sprintf("---\ntitle: p%s%d\n---", lang, i))
		}
	}

	b.Build(BuildCfg{})

	contentMem := b.FileContent(statsFilename)
	cb, err := os.ReadFile(statsFilename)
	b.Assert(err, qt.IsNil)
	contentFile := string(cb)

	for _, content := range []string{contentMem, contentFile} {

		stats := &publisher.PublishStats{}
		b.Assert(json.Unmarshal([]byte(content), stats), qt.IsNil)

		els := stats.HTMLElements

		b.Assert(els.Classes, qt.HasLen, 3606) // (4 * 900) + 4 +2
		b.Assert(els.Tags, qt.HasLen, 8)
		b.Assert(els.IDs, qt.HasLen, 1)
	}
}
