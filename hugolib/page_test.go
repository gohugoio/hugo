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
	"fmt"
	"html/template"
	"os"

	"github.com/gohugoio/hugo/markup/rst"

	"github.com/gohugoio/hugo/markup/asciidoc"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/common/loggers"

	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
)

const (
	homePage   = "---\ntitle: Home\n---\nHome Page Content\n"
	simplePage = "---\ntitle: Simple\n---\nSimple Page\n"

	simplePageRFC3339Date = "---\ntitle: RFC3339 Date\ndate: \"2013-05-17T16:59:30Z\"\n---\nrfc3339 content"

	simplePageWithoutSummaryDelimiter = `---
title: SimpleWithoutSummaryDelimiter
---
[Lorem ipsum](https://lipsum.com/) dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

Additional text.

Further text.
`

	simplePageWithSummaryDelimiter = `---
title: Simple
---
Summary Next Line

<!--more-->
Some more text
`

	simplePageWithSummaryParameter = `---
title: SimpleWithSummaryParameter
summary: "Page with summary parameter and [a link](http://www.example.com/)"
---

Some text.

Some more text.
`

	simplePageWithSummaryDelimiterAndMarkdownThatCrossesBorder = `---
title: Simple
---
The [best static site generator][hugo].[^1]
<!--more-->
[hugo]: http://gohugo.io/
[^1]: Many people say so.
`
	simplePageWithShortcodeInSummary = `---
title: Simple
---
Summary Next Line. {{<figure src="/not/real" >}}.
More text here.

Some more text
`

	simplePageWithSummaryDelimiterSameLine = `---
title: Simple
---
Summary Same Line<!--more-->

Some more text
`

	simplePageWithAllCJKRunes = `---
title: Simple
---


€ € € € €
你好
도형이
カテゴリー


`

	simplePageWithMainEnglishWithCJKRunes = `---
title: Simple
---


In Chinese, 好 means good.  In Chinese, 好 means good.
In Chinese, 好 means good.  In Chinese, 好 means good.
In Chinese, 好 means good.  In Chinese, 好 means good.
In Chinese, 好 means good.  In Chinese, 好 means good.
In Chinese, 好 means good.  In Chinese, 好 means good.
In Chinese, 好 means good.  In Chinese, 好 means good.
In Chinese, 好 means good.  In Chinese, 好 means good.
More then 70 words.


`
	simplePageWithMainEnglishWithCJKRunesSummary = "In Chinese, 好 means good. In Chinese, 好 means good. " +
		"In Chinese, 好 means good. In Chinese, 好 means good. " +
		"In Chinese, 好 means good. In Chinese, 好 means good. " +
		"In Chinese, 好 means good. In Chinese, 好 means good. " +
		"In Chinese, 好 means good. In Chinese, 好 means good. " +
		"In Chinese, 好 means good. In Chinese, 好 means good. " +
		"In Chinese, 好 means good. In Chinese, 好 means good."

	simplePageWithIsCJKLanguageFalse = `---
title: Simple
isCJKLanguage: false
---

In Chinese, 好的啊 means good.  In Chinese, 好的呀 means good.
In Chinese, 好的啊 means good.  In Chinese, 好的呀 means good.
In Chinese, 好的啊 means good.  In Chinese, 好的呀 means good.
In Chinese, 好的啊 means good.  In Chinese, 好的呀 means good.
In Chinese, 好的啊 means good.  In Chinese, 好的呀 means good.
In Chinese, 好的啊 means good.  In Chinese, 好的呀 means good.
In Chinese, 好的啊 means good.  In Chinese, 好的呀呀 means good enough.
More then 70 words.


`
	simplePageWithIsCJKLanguageFalseSummary = "In Chinese, 好的啊 means good. In Chinese, 好的呀 means good. " +
		"In Chinese, 好的啊 means good. In Chinese, 好的呀 means good. " +
		"In Chinese, 好的啊 means good. In Chinese, 好的呀 means good. " +
		"In Chinese, 好的啊 means good. In Chinese, 好的呀 means good. " +
		"In Chinese, 好的啊 means good. In Chinese, 好的呀 means good. " +
		"In Chinese, 好的啊 means good. In Chinese, 好的呀 means good. " +
		"In Chinese, 好的啊 means good. In Chinese, 好的呀呀 means good enough."

	simplePageWithLongContent = `---
title: Simple
---

Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit
amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore
et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation
ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui
officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet,
consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et
dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in
reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia
deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur
adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi
ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim
id est laborum. Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed
do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem
ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit
amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore
et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation
ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui
officia deserunt mollit anim id est laborum.`

	pageWithToC = `---
title: TOC
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.

## AA

I have no idea, of course, how long it took me to reach the limit of the plain,
but at last I entered the foothills, following a pretty little canyon upward
toward the mountains. Beside me frolicked a laughing brooklet, hurrying upon
its noisy way down to the silent sea. In its quieter pools I discovered many
small fish, of four-or five-pound weight I should imagine. In appearance,
except as to size and color, they were not unlike the whale of our own seas. As
I watched them playing about I discovered, not only that they suckled their
young, but that at intervals they rose to the surface to breathe as well as to
feed upon certain grasses and a strange, scarlet lichen which grew upon the
rocks just above the water line.

### AAA

I remember I felt an extraordinary persuasion that I was being played with,
that presently, when I was upon the very verge of safety, this mysterious
death--as swift as the passage of light--would leap after me from the pit about
the cylinder and strike me down. ## BB

### BBB

"You're a great Granser," he cried delightedly, "always making believe them little marks mean something."
`

	simplePageWithAdditionalExtension = `+++
[blackfriday]
  extensions = ["hardLineBreak"]
+++
first line.
second line.

fourth line.
`

	simplePageWithURL = `---
title: Simple
url: simple/url/
---
Simple Page With URL`

	simplePageWithSlug = `---
title: Simple
slug: simple-slug
---
Simple Page With Slug`

	simplePageWithDate = `---
title: Simple
date: '2013-10-15T06:16:13'
---
Simple Page With Date`

	UTF8Page = `---
title: ラーメン
---
UTF8 Page`

	UTF8PageWithURL = `---
title: ラーメン
url: ラーメン/url/
---
UTF8 Page With URL`

	UTF8PageWithSlug = `---
title: ラーメン
slug: ラーメン-slug
---
UTF8 Page With Slug`

	UTF8PageWithDate = `---
title: ラーメン
date: '2013-10-15T06:16:13'
---
UTF8 Page With Date`
)

func checkPageTitle(t *testing.T, page page.Page, title string) {
	if page.Title() != title {
		t.Fatalf("Page title is: %s.  Expected %s", page.Title(), title)
	}
}

func checkPageContent(t *testing.T, page page.Page, expected string, msg ...interface{}) {
	t.Helper()
	a := normalizeContent(expected)
	b := normalizeContent(content(page))
	if a != b {
		t.Fatalf("Page content is:\n%q\nExpected:\n%q (%q)", b, a, msg)
	}
}

func normalizeContent(c string) string {
	norm := c
	norm = strings.Replace(norm, "\n", " ", -1)
	norm = strings.Replace(norm, "    ", " ", -1)
	norm = strings.Replace(norm, "   ", " ", -1)
	norm = strings.Replace(norm, "  ", " ", -1)
	norm = strings.Replace(norm, "p> ", "p>", -1)
	norm = strings.Replace(norm, ">  <", "> <", -1)
	return strings.TrimSpace(norm)
}

func checkPageTOC(t *testing.T, page page.Page, toc string) {
	t.Helper()
	if page.TableOfContents() != template.HTML(toc) {
		t.Fatalf("Page TableOfContents is:\n%q.\nExpected %q", page.TableOfContents(), toc)
	}
}

func checkPageSummary(t *testing.T, page page.Page, summary string, msg ...interface{}) {
	a := normalizeContent(string(page.Summary()))
	b := normalizeContent(summary)
	if a != b {
		t.Fatalf("Page summary is:\n%q.\nExpected\n%q (%q)", a, b, msg)
	}
}

func checkPageType(t *testing.T, page page.Page, pageType string) {
	if page.Type() != pageType {
		t.Fatalf("Page type is: %s.  Expected: %s", page.Type(), pageType)
	}
}

func checkPageDate(t *testing.T, page page.Page, time time.Time) {
	if page.Date() != time {
		t.Fatalf("Page date is: %s.  Expected: %s", page.Date(), time)
	}
}

func normalizeExpected(ext, str string) string {
	str = normalizeContent(str)
	switch ext {
	default:
		return str
	case "html":
		return strings.Trim(helpers.StripHTML(str), " ")
	case "ad":
		paragraphs := strings.Split(str, "</p>")
		expected := ""
		for _, para := range paragraphs {
			if para == "" {
				continue
			}
			expected += fmt.Sprintf("<div class=\"paragraph\">\n%s</p></div>\n", para)
		}

		return expected
	case "rst":
		return fmt.Sprintf("<div class=\"document\">\n\n\n%s</div>", str)
	}
}

func testAllMarkdownEnginesForPages(t *testing.T,
	assertFunc func(t *testing.T, ext string, pages page.Pages), settings map[string]interface{}, pageSources ...string) {

	engines := []struct {
		ext           string
		shouldExecute func() bool
	}{
		{"md", func() bool { return true }},
		{"mmark", func() bool { return true }},
		{"ad", func() bool { return asciidoc.Supports() }},
		{"rst", func() bool { return rst.Supports() }},
	}

	for _, e := range engines {
		if !e.shouldExecute() {
			continue
		}

		cfg, fs := newTestCfg(func(cfg config.Provider) error {
			for k, v := range settings {
				cfg.Set(k, v)
			}
			return nil

		})

		contentDir := "content"

		if s := cfg.GetString("contentDir"); s != "" {
			contentDir = s
		}

		var fileSourcePairs []string

		for i, source := range pageSources {
			fileSourcePairs = append(fileSourcePairs, fmt.Sprintf("p%d.%s", i, e.ext), source)
		}

		for i := 0; i < len(fileSourcePairs); i += 2 {
			writeSource(t, fs, filepath.Join(contentDir, fileSourcePairs[i]), fileSourcePairs[i+1])
		}

		// Add a content page for the home page
		homePath := fmt.Sprintf("_index.%s", e.ext)
		writeSource(t, fs, filepath.Join(contentDir, homePath), homePage)

		b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{Fs: fs, Cfg: cfg}).WithNothingAdded()
		b.Build(BuildCfg{SkipRender: true})

		s := b.H.Sites[0]

		b.Assert(len(s.RegularPages()), qt.Equals, len(pageSources))

		assertFunc(t, e.ext, s.RegularPages())

		home, err := s.Info.Home()
		b.Assert(err, qt.IsNil)
		b.Assert(home, qt.Not(qt.IsNil))
		b.Assert(home.File().Path(), qt.Equals, homePath)
		b.Assert(content(home), qt.Contains, "Home Page Content")

	}

}

// Issue #1076
func TestPageWithDelimiterForMarkdownThatCrossesBorder(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	c := qt.New(t)

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageWithSummaryDelimiterAndMarkdownThatCrossesBorder)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]

	if p.Summary() != template.HTML(
		"<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup id=\"fnref:1\"><a href=\"#fn:1\" class=\"footnote-ref\" role=\"doc-noteref\">1</a></sup></p>") {
		t.Fatalf("Got summary:\n%q", p.Summary())
	}

	cnt := content(p)
	if cnt != "<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup id=\"fnref:1\"><a href=\"#fn:1\" class=\"footnote-ref\" role=\"doc-noteref\">1</a></sup></p>\n<section class=\"footnotes\" role=\"doc-endnotes\">\n<hr>\n<ol>\n<li id=\"fn:1\" role=\"doc-endnote\">\n<p>Many people say so. <a href=\"#fnref:1\" class=\"footnote-backref\" role=\"doc-backlink\">&#x21a9;&#xfe0e;</a></p>\n</li>\n</ol>\n</section>" {
		t.Fatalf("Got content:\n%q", cnt)
	}
}

func TestPageDatesAllKinds(t *testing.T) {
	t.Parallel()

	pageContent := `
---
title: Page
date: 2017-01-15
tags: ["hugo"]
categories: ["cool stuff"]
---
`

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithContent("page.md", pageContent)
	b.WithContent("blog/page.md", pageContent)

	b.CreateSites().Build(BuildCfg{})

	b.Assert(len(b.H.Sites), qt.Equals, 1)
	s := b.H.Sites[0]

	checkDate := func(t time.Time, msg string) {
		b.Assert(t.Year(), qt.Equals, 2017, qt.Commentf(msg))
	}

	checkDated := func(d resource.Dated, msg string) {
		checkDate(d.Date(), "date: "+msg)
		checkDate(d.Lastmod(), "lastmod: "+msg)
	}
	for _, p := range s.Pages() {
		checkDated(p, p.Kind())
	}
	checkDate(s.Info.LastChange(), "site")

}

func TestPageDatesSections(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithContent("no-index/page.md", `
---
title: Page
date: 2017-01-15
---
`)
	b.WithSimpleConfigFile().WithContent("with-index-no-date/_index.md", `---
title: No Date
---

`)

	// https://github.com/gohugoio/hugo/issues/5854
	b.WithSimpleConfigFile().WithContent("with-index-date/_index.md", `---
title: Date
date: 2018-01-15
---

`)

	b.CreateSites().Build(BuildCfg{})

	b.Assert(len(b.H.Sites), qt.Equals, 1)
	s := b.H.Sites[0]

	b.Assert(s.getPage("/").Date().Year(), qt.Equals, 2018)
	b.Assert(s.getPage("/no-index").Date().Year(), qt.Equals, 2017)
	b.Assert(s.getPage("/with-index-no-date").Date().IsZero(), qt.Equals, true)
	b.Assert(s.getPage("/with-index-date").Date().Year(), qt.Equals, 2018)
}

func TestCreateNewPage(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]

		// issue #2290: Path is relative to the content dir and will continue to be so.
		c.Assert(p.File().Path(), qt.Equals, fmt.Sprintf("p0.%s", ext))
		c.Assert(p.IsHome(), qt.Equals, false)
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Simple Page</p>\n"))
		checkPageSummary(t, p, "Simple Page")
		checkPageType(t, p, "page")
	}

	settings := map[string]interface{}{
		"contentDir": "mycontent",
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePage)
}

func TestPageSummary(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		checkPageTitle(t, p, "SimpleWithoutSummaryDelimiter")
		// Source is not Asciidoctor- or RST-compatible so don't test them
		if ext != "ad" && ext != "rst" {
			checkPageContent(t, p, normalizeExpected(ext, "<p><a href=\"https://lipsum.com/\">Lorem ipsum</a> dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>\n\n<p>Additional text.</p>\n\n<p>Further text.</p>\n"), ext)
			checkPageSummary(t, p, normalizeExpected(ext, "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Additional text."), ext)
		}
		checkPageType(t, p, "page")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithoutSummaryDelimiter)
}

func TestPageWithDelimiter(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Next Line</p>\n\n<p>Some more text</p>\n"), ext)
		checkPageSummary(t, p, normalizeExpected(ext, "<p>Summary Next Line</p>"), ext)
		checkPageType(t, p, "page")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithSummaryDelimiter)
}

func TestPageWithSummaryParameter(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		checkPageTitle(t, p, "SimpleWithSummaryParameter")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Some text.</p>\n\n<p>Some more text.</p>\n"), ext)
		// Summary is not Asciidoctor- or RST-compatible so don't test them
		if ext != "ad" && ext != "rst" {
			checkPageSummary(t, p, normalizeExpected(ext, "Page with summary parameter and <a href=\"http://www.example.com/\">a link</a>"), ext)
		}
		checkPageType(t, p, "page")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithSummaryParameter)
}

// Issue #3854
// Also see https://github.com/gohugoio/hugo/issues/3977
func TestPageWithDateFields(t *testing.T) {
	c := qt.New(t)
	pageWithDate := `---
title: P%d
weight: %d
%s: 2017-10-13
---
Simple Page With Some Date`

	hasDate := func(p page.Page) bool {
		return p.Date().Year() == 2017
	}

	datePage := func(field string, weight int) string {
		return fmt.Sprintf(pageWithDate, weight, weight, field)
	}

	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		c.Assert(len(pages) > 0, qt.Equals, true)
		for _, p := range pages {
			c.Assert(hasDate(p), qt.Equals, true)
		}

	}

	fields := []string{"date", "publishdate", "pubdate", "published"}
	pageContents := make([]string, len(fields))
	for i, field := range fields {
		pageContents[i] = datePage(field, i+1)
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, pageContents...)
}

// Issue #2601
func TestPageRawContent(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()
	c := qt.New(t)

	writeSource(t, fs, filepath.Join("content", "raw.md"), `---
title: Raw
---
**Raw**`)

	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .RawContent }}`)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)
	p := s.RegularPages()[0]

	c.Assert("**Raw**", qt.Equals, p.RawContent())

}

func TestPageWithShortCodeInSummary(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Next Line. <figure> <img src=\"/not/real\"/> </figure> . More text here.</p><p>Some more text</p>"))
		checkPageSummary(t, p, "Summary Next Line.  . More text here. Some more text")
		checkPageType(t, p, "page")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithShortcodeInSummary)
}

func TestPageWithAdditionalExtension(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()
	cfg.Set("markup", map[string]interface{}{
		"defaultMarkdownHandler": "blackfriday", // TODO(bep)
	})

	c := qt.New(t)

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageWithAdditionalExtension)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]

	checkPageContent(t, p, "<p>first line.<br />\nsecond line.</p>\n\n<p>fourth line.</p>\n")
}

func TestTableOfContents(t *testing.T) {

	cfg, fs := newTestCfg()
	c := qt.New(t)

	writeSource(t, fs, filepath.Join("content", "tocpage.md"), pageWithToC)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]

	checkPageContent(t, p, "<p>For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.</p><h2 id=\"aa\">AA</h2> <p>I have no idea, of course, how long it took me to reach the limit of the plain, but at last I entered the foothills, following a pretty little canyon upward toward the mountains. Beside me frolicked a laughing brooklet, hurrying upon its noisy way down to the silent sea. In its quieter pools I discovered many small fish, of four-or five-pound weight I should imagine. In appearance, except as to size and color, they were not unlike the whale of our own seas. As I watched them playing about I discovered, not only that they suckled their young, but that at intervals they rose to the surface to breathe as well as to feed upon certain grasses and a strange, scarlet lichen which grew upon the rocks just above the water line.</p><h3 id=\"aaa\">AAA</h3> <p>I remember I felt an extraordinary persuasion that I was being played with, that presently, when I was upon the very verge of safety, this mysterious death&ndash;as swift as the passage of light&ndash;would leap after me from the pit about the cylinder and strike me down. ## BB</p><h3 id=\"bbb\">BBB</h3> <p>&ldquo;You&rsquo;re a great Granser,&rdquo; he cried delightedly, &ldquo;always making believe them little marks mean something.&rdquo;</p>")
	checkPageTOC(t, p, "<nav id=\"TableOfContents\">\n  <ul>\n    <li><a href=\"#aa\">AA</a>\n      <ul>\n        <li><a href=\"#aaa\">AAA</a></li>\n        <li><a href=\"#bbb\">BBB</a></li>\n      </ul>\n    </li>\n  </ul>\n</nav>")
}

func TestPageWithMoreTag(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Same Line</p>\n\n<p>Some more text</p>\n"))
		checkPageSummary(t, p, normalizeExpected(ext, "<p>Summary Same Line</p>"))
		checkPageType(t, p, "page")

	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithSummaryDelimiterSameLine)
}

// #2973
func TestSummaryWithHTMLTagsOnNextLine(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		c := qt.New(t)
		p := pages[0]
		s := string(p.Summary())
		c.Assert(s, qt.Contains, "Happy new year everyone!")
		c.Assert(s, qt.Not(qt.Contains), "User interface")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, `---
title: Simple
---
Happy new year everyone!

Here is the last report for commits in the year 2016. It covers hrev50718-hrev50829.

<!--more-->

<h3>User interface</h3>

`)
}

func TestPageWithDate(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()
	c := qt.New(t)

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageRFC3339Date)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]
	d, _ := time.Parse(time.RFC3339, "2013-05-17T16:59:30Z")

	checkPageDate(t, p, d)
}

func TestPageWithLastmodFromGitInfo(t *testing.T) {
	c := qt.New(t)

	// We need to use the OS fs for this.
	cfg := viper.New()
	fs := hugofs.NewFrom(hugofs.Os, cfg)
	fs.Destination = &afero.MemMapFs{}

	wd, err := os.Getwd()
	c.Assert(err, qt.IsNil)

	cfg.Set("frontmatter", map[string]interface{}{
		"lastmod": []string{":git", "lastmod"},
	})
	cfg.Set("defaultContentLanguage", "en")

	langConfig := map[string]interface{}{
		"en": map[string]interface{}{
			"weight":       1,
			"languageName": "English",
			"contentDir":   "content",
		},
		"nn": map[string]interface{}{
			"weight":       2,
			"languageName": "Nynorsk",
			"contentDir":   "content_nn",
		},
	}

	cfg.Set("languages", langConfig)
	cfg.Set("enableGitInfo", true)

	cfg.Set("workingDir", filepath.Join(wd, "testsite"))

	b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{Fs: fs, Cfg: cfg}).WithNothingAdded()

	b.Build(BuildCfg{SkipRender: true})
	h := b.H

	c.Assert(len(h.Sites), qt.Equals, 2)

	enSite := h.Sites[0]
	c.Assert(len(enSite.RegularPages()), qt.Equals, 1)

	// 2018-03-11 is the Git author date for testsite/content/first-post.md
	c.Assert(enSite.RegularPages()[0].Lastmod().Format("2006-01-02"), qt.Equals, "2018-03-11")

	nnSite := h.Sites[1]
	c.Assert(len(nnSite.RegularPages()), qt.Equals, 1)

	// 2018-08-11 is the Git author date for testsite/content_nn/first-post.md
	c.Assert(nnSite.RegularPages()[0].Lastmod().Format("2006-01-02"), qt.Equals, "2018-08-11")

}

func TestPageWithFrontMatterConfig(t *testing.T) {
	for _, dateHandler := range []string{":filename", ":fileModTime"} {
		dateHandler := dateHandler
		t.Run(fmt.Sprintf("dateHandler=%q", dateHandler), func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)
			cfg, fs := newTestCfg()

			pageTemplate := `
---
title: Page
weight: %d
lastMod: 2018-02-28
%s
---
Content
`

			cfg.Set("frontmatter", map[string]interface{}{
				"date": []string{dateHandler, "date"},
			})

			c1 := filepath.Join("content", "section", "2012-02-21-noslug.md")
			c2 := filepath.Join("content", "section", "2012-02-22-slug.md")

			writeSource(t, fs, c1, fmt.Sprintf(pageTemplate, 1, ""))
			writeSource(t, fs, c2, fmt.Sprintf(pageTemplate, 2, "slug: aslug"))

			c1fi, err := fs.Source.Stat(c1)
			c.Assert(err, qt.IsNil)
			c2fi, err := fs.Source.Stat(c2)
			c.Assert(err, qt.IsNil)

			b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{Fs: fs, Cfg: cfg}).WithNothingAdded()
			b.Build(BuildCfg{SkipRender: true})

			s := b.H.Sites[0]
			c.Assert(len(s.RegularPages()), qt.Equals, 2)

			noSlug := s.RegularPages()[0]
			slug := s.RegularPages()[1]

			c.Assert(noSlug.Lastmod().Day(), qt.Equals, 28)

			switch strings.ToLower(dateHandler) {
			case ":filename":
				c.Assert(noSlug.Date().IsZero(), qt.Equals, false)
				c.Assert(slug.Date().IsZero(), qt.Equals, false)
				c.Assert(noSlug.Date().Year(), qt.Equals, 2012)
				c.Assert(slug.Date().Year(), qt.Equals, 2012)
				c.Assert(noSlug.Slug(), qt.Equals, "noslug")
				c.Assert(slug.Slug(), qt.Equals, "aslug")
			case ":filemodtime":
				c.Assert(noSlug.Date().Year(), qt.Equals, c1fi.ModTime().Year())
				c.Assert(slug.Date().Year(), qt.Equals, c2fi.ModTime().Year())
				fallthrough
			default:
				c.Assert(noSlug.Slug(), qt.Equals, "")
				c.Assert(slug.Slug(), qt.Equals, "aslug")

			}
		})
	}

}

func TestWordCountWithAllCJKRunesWithoutHasCJKLanguage(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount() != 8 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 8, p.WordCount())
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithAllCJKRunes)
}

func TestWordCountWithAllCJKRunesHasCJKLanguage(t *testing.T) {
	t.Parallel()
	settings := map[string]interface{}{"hasCJKLanguage": true}

	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount() != 15 {
			t.Fatalf("[%s] incorrect word count, expected %v, got %v", ext, 15, p.WordCount())
		}
	}
	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithAllCJKRunes)
}

func TestWordCountWithMainEnglishWithCJKRunes(t *testing.T) {
	t.Parallel()
	settings := map[string]interface{}{"hasCJKLanguage": true}

	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount() != 74 {
			t.Fatalf("[%s] incorrect word count, expected %v, got %v", ext, 74, p.WordCount())
		}

		if p.Summary() != simplePageWithMainEnglishWithCJKRunesSummary {
			t.Fatalf("[%s] incorrect Summary for content '%s'. expected %v, got %v", ext, p.Plain(),
				simplePageWithMainEnglishWithCJKRunesSummary, p.Summary())
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithMainEnglishWithCJKRunes)
}

func TestWordCountWithIsCJKLanguageFalse(t *testing.T) {
	t.Parallel()
	settings := map[string]interface{}{
		"hasCJKLanguage": true,
	}

	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount() != 75 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.Plain(), 74, p.WordCount())
		}

		if p.Summary() != simplePageWithIsCJKLanguageFalseSummary {
			t.Fatalf("[%s] incorrect Summary for content '%s'. expected %v, got %v", ext, p.Plain(),
				simplePageWithIsCJKLanguageFalseSummary, p.Summary())
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithIsCJKLanguageFalse)

}

func TestWordCount(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount() != 483 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 483, p.WordCount())
		}

		if p.FuzzyWordCount() != 500 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 500, p.FuzzyWordCount())
		}

		if p.ReadingTime() != 3 {
			t.Fatalf("[%s] incorrect min read. expected %v, got %v", ext, 3, p.ReadingTime())
		}

	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithLongContent)
}

func TestPagePaths(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	siteParmalinksSetting := map[string]string{
		"post": ":year/:month/:day/:title/",
	}

	tests := []struct {
		content      string
		path         string
		hasPermalink bool
		expected     string
	}{
		{simplePage, "post/x.md", false, "post/x.html"},
		{simplePageWithURL, "post/x.md", false, "simple/url/index.html"},
		{simplePageWithSlug, "post/x.md", false, "post/simple-slug.html"},
		{simplePageWithDate, "post/x.md", true, "2013/10/15/simple/index.html"},
		{UTF8Page, "post/x.md", false, "post/x.html"},
		{UTF8PageWithURL, "post/x.md", false, "ラーメン/url/index.html"},
		{UTF8PageWithSlug, "post/x.md", false, "post/ラーメン-slug.html"},
		{UTF8PageWithDate, "post/x.md", true, "2013/10/15/ラーメン/index.html"},
	}

	for _, test := range tests {
		cfg, fs := newTestCfg()

		if test.hasPermalink {
			cfg.Set("permalinks", siteParmalinksSetting)
		}

		writeSource(t, fs, filepath.Join("content", filepath.FromSlash(test.path)), test.content)

		s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})
		c.Assert(len(s.RegularPages()), qt.Equals, 1)

	}
}

func TestTranslationKey(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", filepath.FromSlash("sect/simple.no.md")), "---\ntitle: \"A1\"\ntranslationKey: \"k1\"\n---\nContent\n")
	writeSource(t, fs, filepath.Join("content", filepath.FromSlash("sect/simple.en.md")), "---\ntitle: \"A2\"\n---\nContent\n")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 2)

	home, _ := s.Info.Home()
	c.Assert(home, qt.Not(qt.IsNil))
	c.Assert(home.TranslationKey(), qt.Equals, "home")
	c.Assert(s.RegularPages()[0].TranslationKey(), qt.Equals, "page/k1")
	p2 := s.RegularPages()[1]

	c.Assert(p2.TranslationKey(), qt.Equals, "page/sect/simple")

}

func TestChompBOM(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	const utf8BOM = "\xef\xbb\xbf"

	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "simple.md"), utf8BOM+simplePage)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]

	checkPageTitle(t, p, "Simple")
}

func TestPageWithEmoji(t *testing.T) {
	for _, enableEmoji := range []bool{true, false} {
		v := viper.New()
		v.Set("enableEmoji", enableEmoji)

		b := newTestSitesBuilder(t).WithViper(v)

		b.WithContent("page-emoji.md", `---
title: "Hugo Smile"
---
This is a :smile:.
<!--more--> 

Another :smile: This is :not: :an: :emoji:.

O :christmas_tree:

Write me an :e-mail: or :email:?

Too many colons: :: ::: :::: :?: :!: :.:

If you dislike this video, you can hit that :-1: button :stuck_out_tongue_winking_eye:,
but if you like it, hit :+1: and get subscribed!
`)

		b.CreateSites().Build(BuildCfg{})

		if enableEmoji {
			b.AssertFileContent("public/page-emoji/index.html",
				"This is a 😄",
				"Another 😄",
				"This is :not: :an: :emoji:.",
				"O 🎄",
				"Write me an 📧 or ✉️?",
				"Too many colons: :: ::: :::: :?: :!: :.:",
				"you can hit that 👎 button 😜,",
				"hit 👍 and get subscribed!",
			)
		} else {
			b.AssertFileContent("public/page-emoji/index.html",
				"This is a :smile:",
				"Another :smile:",
				"This is :not: :an: :emoji:.",
				"O :christmas_tree:",
				"Write me an :e-mail: or :email:?",
				"Too many colons: :: ::: :::: :?: :!: :.:",
				"you can hit that :-1: button :stuck_out_tongue_winking_eye:,",
				"hit :+1: and get subscribed!",
			)
		}

	}

}

func TestPageHTMLContent(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile()

	frontmatter := `---
title: "HTML Content"
---
`
	b.WithContent("regular.html", frontmatter+`<h1>Hugo</h1>`)
	b.WithContent("noblackfridayforyou.html", frontmatter+`**Hugo!**`)
	b.WithContent("manualsummary.html", frontmatter+`
<p>This is summary</p>
<!--more-->
<p>This is the main content.</p>`)

	b.Build(BuildCfg{})

	b.AssertFileContent(
		"public/regular/index.html",
		"Single: HTML Content|Hello|en|RelPermalink: /regular/|",
		"Summary: Hugo|Truncated: false")

	b.AssertFileContent(
		"public/noblackfridayforyou/index.html",
		"Permalink: http://example.com/noblackfridayforyou/|**Hugo!**|",
	)

	// https://github.com/gohugoio/hugo/issues/5723
	b.AssertFileContent(
		"public/manualsummary/index.html",
		"Single: HTML Content|Hello|en|RelPermalink: /manualsummary/|",
		"Summary: \n<p>This is summary</p>\n|Truncated: true",
		"|<p>This is the main content.</p>|",
	)

}

// https://github.com/gohugoio/hugo/issues/5381
func TestPageManualSummary(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile()

	b.WithContent("page-md-shortcode.md", `---
title: "Hugo"
---
This is a {{< sc >}}.
<!--more--> 
Content.
`)

	// https://github.com/gohugoio/hugo/issues/5464
	b.WithContent("page-md-only-shortcode.md", `---
title: "Hugo"
---
{{< sc >}}
<!--more--> 
{{< sc >}}
`)

	b.WithContent("page-md-shortcode-same-line.md", `---
title: "Hugo"
---
This is a {{< sc >}}<!--more-->Same line.
`)

	b.WithContent("page-md-shortcode-same-line-after.md", `---
title: "Hugo"
---
Summary<!--more-->{{< sc >}}
`)

	b.WithContent("page-org-shortcode.org", `#+TITLE: T1
#+AUTHOR: A1
#+DESCRIPTION: D1
This is a {{< sc >}}.
# more
Content.	
`)

	b.WithContent("page-org-variant1.org", `#+TITLE: T1
Summary.

# more

Content.	
`)

	b.WithTemplatesAdded("layouts/shortcodes/sc.html", "a shortcode")
	b.WithTemplatesAdded("layouts/_default/single.html", `
SUMMARY:{{ .Summary }}:END
--------------------------
CONTENT:{{ .Content }}
`)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/page-md-shortcode/index.html",
		"SUMMARY:<p>This is a a shortcode.</p>:END",
		"CONTENT:<p>This is a a shortcode.</p>\n\n<p>Content.</p>\n",
	)

	b.AssertFileContent("public/page-md-shortcode-same-line/index.html",
		"SUMMARY:<p>This is a a shortcode</p>:END",
		"CONTENT:<p>This is a a shortcode</p>\n\n<p>Same line.</p>\n",
	)

	b.AssertFileContent("public/page-md-shortcode-same-line-after/index.html",
		"SUMMARY:<p>Summary</p>:END",
		"CONTENT:<p>Summary</p>\n\na shortcode",
	)

	b.AssertFileContent("public/page-org-shortcode/index.html",
		"SUMMARY:<p>\nThis is a a shortcode.\n</p>:END",
		"CONTENT:<p>\nThis is a a shortcode.\n</p>\n<p>\nContent.\t\n</p>\n",
	)
	b.AssertFileContent("public/page-org-variant1/index.html",
		"SUMMARY:<p>\nSummary.\n</p>:END",
		"CONTENT:<p>\nSummary.\n</p>\n<p>\nContent.\t\n</p>\n",
	)

	b.AssertFileContent("public/page-md-only-shortcode/index.html",
		"SUMMARY:a shortcode:END",
		"CONTENT:a shortcode\n\na shortcode\n",
	)
}

// https://github.com/gohugoio/hugo/issues/5478
func TestPageWithCommentedOutFrontMatter(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile()

	b.WithContent("page.md", `<!--
+++
title = "hello"
+++
-->
This is the content.
`)

	b.WithTemplatesAdded("layouts/_default/single.html", `
Title: {{ .Title }}
Content:{{ .Content }}
`)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/page/index.html",
		"Title: hello",
		"Content:<p>This is the content.</p>",
	)

}

// https://github.com/gohugoio/hugo/issues/5781
func TestPageWithZeroFile(t *testing.T) {
	newTestSitesBuilder(t).WithLogger(loggers.NewWarningLogger()).WithSimpleConfigFile().
		WithTemplatesAdded("index.html", "{{ .File.Filename }}{{ with .File }}{{ .Dir }}{{ end }}").Build(BuildCfg{})
}

func TestHomePageWithNoTitle(t *testing.T) {
	b := newTestSitesBuilder(t).WithConfigFile("toml", `
title = "Site Title"
`)
	b.WithTemplatesAdded("index.html", "Title|{{ with .Title }}{{ . }}{{ end }}|")
	b.WithContent("_index.md", `---
description: "No title for you!"
---

Content.
`)

	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html", "Title||")
}

func TestShouldBuild(t *testing.T) {
	t.Parallel()
	var past = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	var future = time.Date(2037, 11, 17, 20, 34, 58, 651387237, time.UTC)
	var zero = time.Time{}

	var publishSettings = []struct {
		buildFuture  bool
		buildExpired bool
		buildDrafts  bool
		draft        bool
		publishDate  time.Time
		expiryDate   time.Time
		out          bool
	}{
		// publishDate and expiryDate
		{false, false, false, false, zero, zero, true},
		{false, false, false, false, zero, future, true},
		{false, false, false, false, past, zero, true},
		{false, false, false, false, past, future, true},
		{false, false, false, false, past, past, false},
		{false, false, false, false, future, future, false},
		{false, false, false, false, future, past, false},

		// buildFuture and buildExpired
		{false, true, false, false, past, past, true},
		{true, true, false, false, past, past, true},
		{true, false, false, false, past, past, false},
		{true, false, false, false, future, future, true},
		{true, true, false, false, future, future, true},
		{false, true, false, false, future, past, false},

		// buildDrafts and draft
		{true, true, false, true, past, future, false},
		{true, true, true, true, past, future, true},
		{true, true, true, true, past, future, true},
	}

	for _, ps := range publishSettings {
		s := shouldBuild(ps.buildFuture, ps.buildExpired, ps.buildDrafts, ps.draft,
			ps.publishDate, ps.expiryDate)
		if s != ps.out {
			t.Errorf("AssertShouldBuild unexpected output with params: %+v", ps)
		}
	}
}

// "dot" in path: #1885 and #2110
// disablePathToLower regression: #3374
func TestPathIssues(t *testing.T) {
	for _, disablePathToLower := range []bool{false, true} {
		for _, uglyURLs := range []bool{false, true} {
			disablePathToLower := disablePathToLower
			uglyURLs := uglyURLs
			t.Run(fmt.Sprintf("disablePathToLower=%t,uglyURLs=%t", disablePathToLower, uglyURLs), func(t *testing.T) {
				t.Parallel()
				cfg, fs := newTestCfg()
				th := newTestHelper(cfg, fs, t)
				c := qt.New(t)

				cfg.Set("permalinks", map[string]string{
					"post": ":section/:title",
				})

				cfg.Set("uglyURLs", uglyURLs)
				cfg.Set("disablePathToLower", disablePathToLower)
				cfg.Set("paginate", 1)

				writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), "<html><body>{{.Content}}</body></html>")
				writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"),
					"<html><body>P{{.Paginator.PageNumber}}|URL: {{.Paginator.URL}}|{{ if .Paginator.HasNext }}Next: {{.Paginator.Next.URL }}{{ end }}</body></html>")

				for i := 0; i < 3; i++ {
					writeSource(t, fs, filepath.Join("content", "post", fmt.Sprintf("doc%d.md", i)),
						fmt.Sprintf(`---
title: "test%d.dot"
tags:
- ".net"
---
# doc1
*some content*`, i))
				}

				writeSource(t, fs, filepath.Join("content", "Blog", "Blog1.md"),
					fmt.Sprintf(`---
title: "testBlog"
tags:
- "Blog"
---
# doc1
*some blog content*`))

				s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

				c.Assert(len(s.RegularPages()), qt.Equals, 4)

				pathFunc := func(s string) string {
					if uglyURLs {
						return strings.Replace(s, "/index.html", ".html", 1)
					}
					return s
				}

				blog := "blog"

				if disablePathToLower {
					blog = "Blog"
				}

				th.assertFileContent(pathFunc("public/"+blog+"/"+blog+"1/index.html"), "some blog content")

				th.assertFileContent(pathFunc("public/post/test0.dot/index.html"), "some content")

				if uglyURLs {
					th.assertFileContent("public/post/page/1.html", `canonical" href="/post.html"/`)
					th.assertFileContent("public/post.html", `<body>P1|URL: /post.html|Next: /post/page/2.html</body>`)
					th.assertFileContent("public/post/page/2.html", `<body>P2|URL: /post/page/2.html|Next: /post/page/3.html</body>`)
				} else {
					th.assertFileContent("public/post/page/1/index.html", `canonical" href="/post/"/`)
					th.assertFileContent("public/post/index.html", `<body>P1|URL: /post/|Next: /post/page/2/</body>`)
					th.assertFileContent("public/post/page/2/index.html", `<body>P2|URL: /post/page/2/|Next: /post/page/3/</body>`)
					th.assertFileContent("public/tags/.net/index.html", `<body>P1|URL: /tags/.net/|Next: /tags/.net/page/2/</body>`)

				}

				p := s.RegularPages()[0]
				if uglyURLs {
					c.Assert(p.RelPermalink(), qt.Equals, "/post/test0.dot.html")
				} else {
					c.Assert(p.RelPermalink(), qt.Equals, "/post/test0.dot/")
				}

			})
		}
	}
}

// https://github.com/gohugoio/hugo/issues/4675
func TestWordCountAndSimilarVsSummary(t *testing.T) {

	t.Parallel()
	c := qt.New(t)

	single := []string{"_default/single.html", `
WordCount: {{ .WordCount }}
FuzzyWordCount: {{ .FuzzyWordCount }}
ReadingTime: {{ .ReadingTime }}
Len Plain: {{ len .Plain }}
Len PlainWords: {{ len .PlainWords }}
Truncated: {{ .Truncated }}
Len Summary: {{ len .Summary }}
Len Content: {{ len .Content }}

SUMMARY:{{ .Summary }}:{{ len .Summary }}:END
`}

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded(single...).WithContent("p1.md", fmt.Sprintf(`---
title: p1	
---

%s

`, strings.Repeat("word ", 510)),

		"p2.md", fmt.Sprintf(`---
title: p2
---
This is a summary.

<!--more-->

%s

`, strings.Repeat("word ", 310)),
		"p3.md", fmt.Sprintf(`---
title: p3
isCJKLanguage: true
---
Summary: In Chinese, 好 means good.

<!--more-->

%s

`, strings.Repeat("好", 200)),
		"p4.md", fmt.Sprintf(`---
title: p4
isCJKLanguage: false
---
Summary: In Chinese, 好 means good.

<!--more-->

%s

`, strings.Repeat("好", 200)),

		"p5.md", fmt.Sprintf(`---
title: p4
isCJKLanguage: true
---
Summary: In Chinese, 好 means good.

%s

`, strings.Repeat("好", 200)),
		"p6.md", fmt.Sprintf(`---
title: p4
isCJKLanguage: false
---
Summary: In Chinese, 好 means good.

%s

`, strings.Repeat("好", 200)),
	)

	b.CreateSites().Build(BuildCfg{})

	c.Assert(len(b.H.Sites), qt.Equals, 1)
	c.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 6)

	b.AssertFileContent("public/p1/index.html", "WordCount: 510\nFuzzyWordCount: 600\nReadingTime: 3\nLen Plain: 2550\nLen PlainWords: 510\nTruncated: false\nLen Summary: 2549\nLen Content: 2557")

	b.AssertFileContent("public/p2/index.html", "WordCount: 314\nFuzzyWordCount: 400\nReadingTime: 2\nLen Plain: 1569\nLen PlainWords: 314\nTruncated: true\nLen Summary: 25\nLen Content: 1582")

	b.AssertFileContent("public/p3/index.html", "WordCount: 206\nFuzzyWordCount: 300\nReadingTime: 1\nLen Plain: 638\nLen PlainWords: 7\nTruncated: true\nLen Summary: 43\nLen Content: 651")
	b.AssertFileContent("public/p4/index.html", "WordCount: 7\nFuzzyWordCount: 100\nReadingTime: 1\nLen Plain: 638\nLen PlainWords: 7\nTruncated: true\nLen Summary: 43\nLen Content: 651")
	b.AssertFileContent("public/p5/index.html", "WordCount: 206\nFuzzyWordCount: 300\nReadingTime: 1\nLen Plain: 638\nLen PlainWords: 7\nTruncated: true\nLen Summary: 229\nLen Content: 652")
	b.AssertFileContent("public/p6/index.html", "WordCount: 7\nFuzzyWordCount: 100\nReadingTime: 1\nLen Plain: 638\nLen PlainWords: 7\nTruncated: false\nLen Summary: 637\nLen Content: 652")

}

func TestScratchSite(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded("index.html", `
{{ .Scratch.Set "b" "bv" }}
B: {{ .Scratch.Get "b" }}
`,
		"shortcodes/scratch.html", `
{{ .Scratch.Set "c" "cv" }}
C: {{ .Scratch.Get "c" }}
`,
	)

	b.WithContentAdded("scratchme.md", `
---
title: Scratch Me!
---

{{< scratch >}}
`)
	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "B: bv")
	b.AssertFileContent("public/scratchme/index.html", "C: cv")
}

func TestPageParam(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).WithConfigFile("toml", `

baseURL = "https://example.org"

[params]
[params.author]
  name = "Kurt Vonnegut"

`)
	b.WithTemplatesAdded("index.html", `

{{ $withParam := .Site.GetPage "withparam" }}
{{ $noParam := .Site.GetPage "noparam" }}
{{ $withStringParam := .Site.GetPage "withstringparam" }}

Author page: {{ $withParam.Param "author.name" }}
Author name page string: {{ $withStringParam.Param "author.name" }}|
Author page string: {{ $withStringParam.Param "author" }}|
Author site config:  {{ $noParam.Param "author.name" }}

`,
	)

	b.WithContent("withparam.md", `
+++
title = "With Param!"
[author]
  name = "Ernest Miller Hemingway"

+++

`,

		"noparam.md", `
---
title: "No Param!"
---
`, "withstringparam.md", `
+++
title = "With string Param!"
author = "Jo Nesbø"

+++

`)
	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		"Author page: Ernest Miller Hemingway",
		"Author name page string: Kurt Vonnegut|",
		"Author page string: Jo Nesbø|",
		"Author site config:  Kurt Vonnegut")

}

func TestGoldmark(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).WithConfigFile("toml", `
baseURL = "https://example.org"

[markup]
defaultMarkdownHandler="goldmark"
[markup.goldmark]
[markup.goldmark.renderer]
unsafe = false
[markup.highlight]
noClasses=false


`)
	b.WithTemplatesAdded("_default/single.html", `
Title: {{ .Title }}
ToC: {{ .TableOfContents }}
Content: {{ .Content }}

`, "shortcodes/t.html", `T-SHORT`, "shortcodes/s.html", `## Code
{{ .Inner }}
`)

	content := `
+++
title = "A Page!"
+++

## Shortcode {{% t %}} in header

## Code Fense in Shortcode

{{% s %}}
$$$bash {hl_lines=[1]}
SHORT
$$$
{{% /s %}}

## Code Fence

$$$bash {hl_lines=[1]}
MARKDOWN
$$$

Link with URL as text

[https://google.com](https://google.com)


`
	content = strings.ReplaceAll(content, "$$$", "```")

	b.WithContent("page.md", content)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/page/index.html",
		`<nav id="TableOfContents">
<li><a href="#shortcode-t-short-in-header">Shortcode T-SHORT in header</a></li>
<code class="language-bash" data-lang="bash"><span class="hl">SHORT
<code class="language-bash" data-lang="bash"><span class="hl">MARKDOWN
<p><a href="https://google.com">https://google.com</a></p>
`)
}

func TestBlackfridayDefault(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).WithConfigFile("toml", `
baseURL = "https://example.org"

[markup]
defaultMarkdownHandler="blackfriday"
[markup.highlight]
noClasses=false
[markup.goldmark]
[markup.goldmark.renderer]
unsafe=true


`)
	// Use the new attribute syntax to make sure it's not Goldmark.
	b.WithTemplatesAdded("_default/single.html", `
Title: {{ .Title }}
Content: {{ .Content }}

`, "shortcodes/s.html", `## Code
{{ .Inner }}
`)

	content := `
+++
title = "A Page!"
+++


## Code Fense in Shortcode

{{% s %}}
S:
{{% s %}}
$$$bash {hl_lines=[1]}
SHORT
$$$
{{% /s %}}
{{% /s %}}

## Code Fence

$$$bash {hl_lines=[1]}
MARKDOWN
$$$

`
	content = strings.ReplaceAll(content, "$$$", "```")

	for i, ext := range []string{"md", "html"} {
		b.WithContent(fmt.Sprintf("page%d.%s", i+1, ext), content)

	}

	b.Build(BuildCfg{})

	// Blackfriday does not support this extended attribute syntax.
	b.AssertFileContent("public/page1/index.html",
		`<pre><code class="language-bash {hl_lines=[1]}" data-lang="bash {hl_lines=[1]}">SHORT</code></pre>`,
		`<pre><code class="language-bash {hl_lines=[1]}" data-lang="bash {hl_lines=[1]}">MARKDOWN`,
	)

	b.AssertFileContent("public/page2/index.html",
		`<pre><code class="language-bash {hl_lines=[1]}" data-lang="bash {hl_lines=[1]}">SHORT`,
	)
}
