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
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bep/clocks"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/rst"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
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

func checkPageContent(t *testing.T, page page.Page, expected string, msg ...any) {
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
	if page.TableOfContents(context.Background()) != template.HTML(toc) {
		t.Fatalf("Page TableOfContents is:\n%q.\nExpected %q", page.TableOfContents(context.Background()), toc)
	}
}

func checkPageSummary(t *testing.T, page page.Page, summary string, msg ...any) {
	s := string(page.Summary(context.Background()))
	a := normalizeContent(s)
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
		return strings.Trim(tpl.StripHTML(str), " ")
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
		if str == "" {
			return "<div class=\"document\"></div>"
		}
		return fmt.Sprintf("<div class=\"document\">\n\n\n%s</div>", str)
	}
}

func testAllMarkdownEnginesForPages(t *testing.T,
	assertFunc func(t *testing.T, ext string, pages page.Pages), settings map[string]any, pageSources ...string,
) {
	engines := []struct {
		ext           string
		shouldExecute func() bool
	}{
		{"md", func() bool { return true }},
		{"ad", func() bool { return asciidocext.Supports() }},
		{"rst", func() bool { return rst.Supports() }},
	}

	for _, e := range engines {
		if !e.shouldExecute() {
			continue
		}

		t.Run(e.ext, func(t *testing.T) {
			cfg := config.New()
			for k, v := range settings {
				cfg.Set(k, v)
			}

			if s := cfg.GetString("contentDir"); s != "" && s != "content" {
				panic("contentDir must be set to 'content' for this test")
			}

			files := `
-- hugo.toml --
[security]
[security.exec]
allow = ['^python$', '^rst2html.*', '^asciidoctor$']
`

			for i, source := range pageSources {
				files += fmt.Sprintf("-- content/p%d.%s --\n%s\n", i, e.ext, source)
			}
			homePath := fmt.Sprintf("_index.%s", e.ext)
			files += fmt.Sprintf("-- content/%s --\n%s\n", homePath, homePage)

			b := NewIntegrationTestBuilder(
				IntegrationTestConfig{
					T:           t,
					TxtarString: files,
					NeedsOsFS:   true,
					BaseCfg:     cfg,
				},
			).Build()

			s := b.H.Sites[0]

			b.Assert(len(s.RegularPages()), qt.Equals, len(pageSources))

			assertFunc(t, e.ext, s.RegularPages())

			home := s.Home()
			b.Assert(home, qt.Not(qt.IsNil))
			b.Assert(home.File().Path(), qt.Equals, homePath)
			b.Assert(content(home), qt.Contains, "Home Page Content")
		})

	}
}

// Issue #1076
func TestPageWithDelimiterForMarkdownThatCrossesBorder(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	c := qt.New(t)
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageWithSummaryDelimiterAndMarkdownThatCrossesBorder)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]

	if p.Summary(context.Background()) != template.HTML(
		"<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup id=\"fnref:1\"><a href=\"#fn:1\" class=\"footnote-ref\" role=\"doc-noteref\">1</a></sup></p>") {
		t.Fatalf("Got summary:\n%q", p.Summary(context.Background()))
	}

	cnt := content(p)
	if cnt != "<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup id=\"fnref:1\"><a href=\"#fn:1\" class=\"footnote-ref\" role=\"doc-noteref\">1</a></sup></p>\n<div class=\"footnotes\" role=\"doc-endnotes\">\n<hr>\n<ol>\n<li id=\"fn:1\">\n<p>Many people say so.&#160;<a href=\"#fnref:1\" class=\"footnote-backref\" role=\"doc-backlink\">&#x21a9;&#xfe0e;</a></p>\n</li>\n</ol>\n</div>" {
		t.Fatalf("Got content:\n%q", cnt)
	}
}

func TestPageDatesTerms(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/p1.md --
---
title: p1
date: 2022-01-15
lastMod: 2022-01-16
tags: ["a", "b"]
categories: ["c", "d"]
---
p1
-- content/p2.md --
---
title: p2
date: 2017-01-16
lastMod: 2017-01-17
tags: ["a", "c"]
categories: ["c", "e"]
---
p2
-- layouts/_default/list.html --
{{ .Title }}|Date: {{ .Date.Format "2006-01-02" }}|Lastmod: {{ .Lastmod.Format "2006-01-02" }}|

`
	b := Test(t, files)

	b.AssertFileContent("public/categories/index.html", "Categories|Date: 2022-01-15|Lastmod: 2022-01-16|")
	b.AssertFileContent("public/categories/c/index.html", "C|Date: 2022-01-15|Lastmod: 2022-01-16|")
	b.AssertFileContent("public/categories/e/index.html", "E|Date: 2017-01-16|Lastmod: 2017-01-17|")
	b.AssertFileContent("public/tags/index.html", "Tags|Date: 2022-01-15|Lastmod: 2022-01-16|")
	b.AssertFileContent("public/tags/a/index.html", "A|Date: 2022-01-15|Lastmod: 2022-01-16|")
	b.AssertFileContent("public/tags/c/index.html", "C|Date: 2017-01-16|Lastmod: 2017-01-17|")
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
		b.Helper()
		b.Assert(t.Year(), qt.Equals, 2017, qt.Commentf(msg))
	}

	checkDated := func(d resource.Dated, msg string) {
		b.Helper()
		checkDate(d.Date(), "date: "+msg)
		checkDate(d.Lastmod(), "lastmod: "+msg)
	}
	for _, p := range s.Pages() {
		checkDated(p, p.Kind())
	}
	checkDate(s.Lastmod(), "site")
}

func TestPageDatesSections(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithContent("no-index/page.md", `
---
title: Page
date: 2017-01-15
---
`, "with-index-no-date/_index.md", `---
title: No Date
---

`,
		// https://github.com/gohugoio/hugo/issues/5854
		"with-index-date/_index.md", `---
title: Date
date: 2018-01-15
---

`, "with-index-date/p1.md", `---
title: Date
date: 2018-01-15
---

`, "with-index-date/p1.md", `---
title: Date
date: 2018-01-15
---

`)

	for i := 1; i <= 20; i++ {
		b.WithContent(fmt.Sprintf("main-section/p%d.md", i), `---
title: Date
date: 2012-01-12
---

`)
	}

	b.CreateSites().Build(BuildCfg{})

	b.Assert(len(b.H.Sites), qt.Equals, 1)
	s := b.H.Sites[0]

	checkDate := func(p page.Page, year int) {
		b.Assert(p.Date().Year(), qt.Equals, year)
		b.Assert(p.Lastmod().Year(), qt.Equals, year)
	}

	checkDate(s.getPageOldVersion("/"), 2018)
	checkDate(s.getPageOldVersion("/no-index"), 2017)
	b.Assert(s.getPageOldVersion("/with-index-no-date").Date().IsZero(), qt.Equals, true)
	checkDate(s.getPageOldVersion("/with-index-date"), 2018)

	b.Assert(s.Site().Lastmod().Year(), qt.Equals, 2018)
}

func TestPageSummary(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		checkPageTitle(t, p, "SimpleWithoutSummaryDelimiter")
		// Source is not Asciidoctor- or RST-compatible so don't test them
		if ext != "ad" && ext != "rst" {
			checkPageContent(t, p, normalizeExpected(ext, "<p><a href=\"https://lipsum.com/\">Lorem ipsum</a> dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>\n\n<p>Additional text.</p>\n\n<p>Further text.</p>\n"), ext)
			checkPageSummary(t, p, normalizeExpected(ext, "<p><a href=\"https://lipsum.com/\">Lorem ipsum</a> dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p><p>Additional text.</p>"), ext)
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

func TestPageRawContent(t *testing.T) {
	files := `
-- hugo.toml --
-- content/basic.md --
---
title: "basic"
---
**basic**
-- content/empty.md --
---
title: "empty"
---
-- layouts/_default/single.html --
|{{ .RawContent }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/basic/index.html", "|**basic**|")
	b.AssertFileContent("public/empty/index.html", "! title")
}

func TestTableOfContents(t *testing.T) {
	c := qt.New(t)
	cfg, fs := newTestCfg()
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "tocpage.md"), pageWithToC)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

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

func TestSummaryInFrontMatter(t *testing.T) {
	t.Parallel()
	Test(t, `
-- hugo.toml --
-- content/simple.md --
---
title: Simple
summary: "Front **matter** summary"
---
Simple Page
-- layouts/_default/single.html --
Summary: {{ .Summary }}|Truncated: {{ .Truncated }}|

`).AssertFileContent("public/simple/index.html", "Summary: Front <strong>matter</strong> summary|", "Truncated: false")
}

func TestSummaryManualSplit(t *testing.T) {
	t.Parallel()
	Test(t, `
-- hugo.toml --
-- content/simple.md --
---
title: Simple
---
This is **summary**.
<!--more-->
This is **content**.
-- layouts/_default/single.html --
Summary: {{ .Summary }}|Truncated: {{ .Truncated }}|
Content: {{ .Content }}|

`).AssertFileContent("public/simple/index.html",
		"Summary: <p>This is <strong>summary</strong>.</p>|",
		"Truncated: true|",
		"Content: <p>This is <strong>summary</strong>.</p>\n<p>This is <strong>content</strong>.</p>|",
	)
}

func TestSummaryManualSplitHTML(t *testing.T) {
	t.Parallel()
	Test(t, `
-- hugo.toml --
-- content/simple.html --
---
title: Simple
---
<div>
This is <b>summary</b>.
</div>
<!--more-->
<div>
This is <b>content</b>.
</div>
-- layouts/_default/single.html --
Summary: {{ .Summary }}|Truncated: {{ .Truncated }}|
Content: {{ .Content }}|

`).AssertFileContent("public/simple/index.html", "Summary: <div>\nThis is <b>summary</b>.\n</div>\n|Truncated: true|\nContent: \n\n<div>\nThis is <b>content</b>.\n</div>|")
}

func TestSummaryAuto(t *testing.T) {
	t.Parallel()
	Test(t, `
-- hugo.toml --
summaryLength = 10
-- content/simple.md --
---
title: Simple
---
This is **summary**.
This is **more summary**.
This is *even more summary**.
This is **more summary**.

This is **content**.
-- layouts/_default/single.html --
Summary: {{ .Summary }}|Truncated: {{ .Truncated }}|
Content: {{ .Content }}|

`).AssertFileContent("public/simple/index.html",
		"Summary: <p>This is <strong>summary</strong>.\nThis is <strong>more summary</strong>.\nThis is <em>even more summary</em>*.\nThis is <strong>more summary</strong>.</p>|",
		"Truncated: true|",
		"Content: <p>This is <strong>summary</strong>.")
}

// #2973
func TestSummaryWithHTMLTagsOnNextLine(t *testing.T) {
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		c := qt.New(t)
		p := pages[0]
		s := string(p.Summary(context.Background()))
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

// Issue 9383
func TestRenderStringForRegularPageTranslations(t *testing.T) {
	c := qt.New(t)
	b := newTestSitesBuilder(t)
	b.WithLogger(loggers.NewDefault())

	b.WithConfigFile("toml",
		`baseurl = "https://example.org/"
title = "My Site"

defaultContentLanguage = "ru"
defaultContentLanguageInSubdir = true

[languages.ru]
contentDir = 'content/ru'
weight = 1

[languages.en]
weight = 2
contentDir = 'content/en'

[outputs]
home = ["HTML", "JSON"]`)

	b.WithTemplates("index.html", `
{{- range .Site.Home.Translations -}}
	<p>{{- .RenderString "foo" -}}</p>
{{- end -}}
{{- range .Site.Home.AllTranslations -}}
	<p>{{- .RenderString "bar" -}}</p>
{{- end -}}
`, "_default/single.html",
		`{{ .Content }}`,
		"index.json",
		`{"Title": "My Site"}`,
	)

	b.WithContent(
		"ru/a.md",
		"",
		"en/a.md",
		"",
	)

	err := b.BuildE(BuildCfg{})
	c.Assert(err, qt.Equals, nil)

	b.AssertFileContent("public/ru/index.html", `
<p>foo</p>
<p>foo</p>
<p>bar</p>
<p>bar</p>
`)

	b.AssertFileContent("public/en/index.html", `
<p>foo</p>
<p>foo</p>
<p>bar</p>
<p>bar</p>
`)
}

// Issue 8919
func TestContentProviderWithCustomOutputFormat(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithLogger(loggers.NewDefault())
	b.WithConfigFile("toml", `baseURL = 'http://example.org/'
title = 'My New Hugo Site'

timeout = 600000 # ten minutes in case we want to pause and debug

defaultContentLanguage = "en"

[languages]
	[languages.en]
	title = "Repro"
	languageName = "English"
	contentDir = "content/en"

	[languages.zh_CN]
	title = "Repro"
	languageName = "简体中文"
	contentDir = "content/zh_CN"

[outputFormats]
	[outputFormats.metadata]
	baseName = "metadata"
	mediaType = "text/html"
	isPlainText = true
	notAlternative = true

[outputs]
	home = ["HTML", "metadata"]`)

	b.WithTemplates("home.metadata.html", `<h2>Translations metadata</h2>
<ul>
{{ $p := .Page }}
{{ range $p.Translations}}
<li>Title: {{ .Title }}, {{ .Summary }}</li>
<li>Content: {{ .Content }}</li>
<li>Plain: {{ .Plain }}</li>
<li>PlainWords: {{ .PlainWords }}</li>
<li>Summary: {{ .Summary }}</li>
<li>Truncated: {{ .Truncated }}</li>
<li>FuzzyWordCount: {{ .FuzzyWordCount }}</li>
<li>ReadingTime: {{ .ReadingTime }}</li>
<li>Len: {{ .Len }}</li>
{{ end }}
</ul>`)

	b.WithTemplates("_default/baseof.html", `<html>

<body>
	{{ block "main" . }}{{ end }}
</body>

</html>`)

	b.WithTemplates("_default/home.html", `{{ define "main" }}
<h2>Translations</h2>
<ul>
{{ $p := .Page }}
{{ range $p.Translations}}
<li>Title: {{ .Title }}, {{ .Summary }}</li>
<li>Content: {{ .Content }}</li>
<li>Plain: {{ .Plain }}</li>
<li>PlainWords: {{ .PlainWords }}</li>
<li>Summary: {{ .Summary }}</li>
<li>Truncated: {{ .Truncated }}</li>
<li>FuzzyWordCount: {{ .FuzzyWordCount }}</li>
<li>ReadingTime: {{ .ReadingTime }}</li>
<li>Len: {{ .Len }}</li>
{{ end }}
</ul>
{{ end }}`)

	b.WithContent("en/_index.md", `---
title: Title (en)
summary: Summary (en)
---

Here is some content.
`)

	b.WithContent("zh_CN/_index.md", `---
title: Title (zh)
summary: Summary (zh)
---

这是一些内容
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `<html>

<body>

<h2>Translations</h2>
<ul>


<li>Title: Title (zh), Summary (zh)</li>
<li>Content: <p>这是一些内容</p>
</li>
<li>Plain: 这是一些内容
</li>
<li>PlainWords: [这是一些内容]</li>
<li>Summary: Summary (zh)</li>
<li>Truncated: false</li>
<li>FuzzyWordCount: 100</li>
<li>ReadingTime: 1</li>
<li>Len: 26</li>

</ul>

</body>

</html>`)
	b.AssertFileContent("public/metadata.html", `<h2>Translations metadata</h2>
<ul>


<li>Title: Title (zh), Summary (zh)</li>
<li>Content: <p>这是一些内容</p>
</li>
<li>Plain: 这是一些内容
</li>
<li>PlainWords: [这是一些内容]</li>
<li>Summary: Summary (zh)</li>
<li>Truncated: false</li>
<li>FuzzyWordCount: 100</li>
<li>ReadingTime: 1</li>
<li>Len: 26</li>

</ul>`)
	b.AssertFileContent("public/zh_cn/index.html", `<html>

<body>

<h2>Translations</h2>
<ul>


<li>Title: Title (en), Summary (en)</li>
<li>Content: <p>Here is some content.</p>
</li>
<li>Plain: Here is some content.
</li>
<li>PlainWords: [Here is some content.]</li>
<li>Summary: Summary (en)</li>
<li>Truncated: false</li>
<li>FuzzyWordCount: 100</li>
<li>ReadingTime: 1</li>
<li>Len: 29</li>

</ul>

</body>

</html>`)
	b.AssertFileContent("public/zh_cn/metadata.html", `<h2>Translations metadata</h2>
<ul>


<li>Title: Title (en), Summary (en)</li>
<li>Content: <p>Here is some content.</p>
</li>
<li>Plain: Here is some content.
</li>
<li>PlainWords: [Here is some content.]</li>
<li>Summary: Summary (en)</li>
<li>Truncated: false</li>
<li>FuzzyWordCount: 100</li>
<li>ReadingTime: 1</li>
<li>Len: 29</li>

</ul>`)
}

func TestPageWithDate(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	cfg, fs := newTestCfg()
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageRFC3339Date)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]
	d, _ := time.Parse(time.RFC3339, "2013-05-17T16:59:30Z")

	checkPageDate(t, p, d)
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

			cfg.Set("frontmatter", map[string]any{
				"date": []string{dateHandler, "date"},
			})
			configs, err := loadTestConfigFromProvider(cfg)
			c.Assert(err, qt.IsNil)

			c1 := filepath.Join("content", "section", "2012-02-21-noslug.md")
			c2 := filepath.Join("content", "section", "2012-02-22-slug.md")

			writeSource(t, fs, c1, fmt.Sprintf(pageTemplate, 1, ""))
			writeSource(t, fs, c2, fmt.Sprintf(pageTemplate, 2, "slug: aslug"))

			c1fi, err := fs.Source.Stat(c1)
			c.Assert(err, qt.IsNil)
			c2fi, err := fs.Source.Stat(c2)
			c.Assert(err, qt.IsNil)

			b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{Fs: fs, Configs: configs}).WithNothingAdded()
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
		if p.WordCount(context.Background()) != 8 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 8, p.WordCount(context.Background()))
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithAllCJKRunes)
}

func TestWordCountWithAllCJKRunesHasCJKLanguage(t *testing.T) {
	t.Parallel()
	settings := map[string]any{"hasCJKLanguage": true}

	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount(context.Background()) != 15 {
			t.Fatalf("[%s] incorrect word count, expected %v, got %v", ext, 15, p.WordCount(context.Background()))
		}
	}
	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithAllCJKRunes)
}

func TestWordCountWithMainEnglishWithCJKRunes(t *testing.T) {
	t.Parallel()
	settings := map[string]any{"hasCJKLanguage": true}

	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount(context.Background()) != 74 {
			t.Fatalf("[%s] incorrect word count, expected %v, got %v", ext, 74, p.WordCount(context.Background()))
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithMainEnglishWithCJKRunes)
}

func TestWordCountWithIsCJKLanguageFalse(t *testing.T) {
	t.Parallel()
	settings := map[string]any{
		"hasCJKLanguage": true,
	}

	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount(context.Background()) != 75 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.Plain(context.Background()), 74, p.WordCount(context.Background()))
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithIsCJKLanguageFalse)
}

func TestWordCount(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages page.Pages) {
		p := pages[0]
		if p.WordCount(context.Background()) != 483 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 483, p.WordCount(context.Background()))
		}

		if p.FuzzyWordCount(context.Background()) != 500 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 500, p.FuzzyWordCount(context.Background()))
		}

		if p.ReadingTime(context.Background()) != 3 {
			t.Fatalf("[%s] incorrect min read. expected %v, got %v", ext, 3, p.ReadingTime(context.Background()))
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
		configs, err := loadTestConfigFromProvider(cfg)
		c.Assert(err, qt.IsNil)

		if test.hasPermalink {
			cfg.Set("permalinks", siteParmalinksSetting)
		}

		writeSource(t, fs, filepath.Join("content", filepath.FromSlash(test.path)), test.content)

		s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})
		c.Assert(len(s.RegularPages()), qt.Equals, 1)

	}
}

func TestTranslationKey(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/sect/p1.en.md --
---
translationkey: "adfasdf"
title: "p1 en"
---
-- content/sect/p1.nn.md --
---
translationkey: "adfasdf"
title: "p1 nn"
---
-- layouts/_default/single.html --
Title: {{ .Title }}|TranslationKey: {{ .TranslationKey }}|
Translations: {{ range .Translations }}{{ .Language.Lang }}|{{ end }}|
AllTranslations: {{ range .AllTranslations }}{{ .Language.Lang }}|{{ end }}|

`
	b := Test(t, files)
	b.AssertFileContent("public/en/sect/p1/index.html",
		"TranslationKey: adfasdf|",
		"AllTranslations: en|nn||",
		"Translations: nn||",
	)

	b.AssertFileContent("public/nn/sect/p1/index.html",
		"TranslationKey: adfasdf|",
		"Translations: en||",
		"AllTranslations: en|nn||",
	)
}

func TestTranslationKeyTermPages(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy']
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages.en]
weight = 1
[languages.pt]
weight = 2
[taxonomies]
category = 'categories'
-- layouts/_default/list.html --
{{ .IsTranslated }}|{{ range .Translations }}{{ .RelPermalink }}|{{ end }}
-- layouts/_default/single.html --
{{ .Title }}|
-- content/p1.en.md --
---
title: p1 (en)
categories: [music]
---
-- content/p1.pt.md --
---
title: p1 (pt)
categories: [música]
---
-- content/categories/music/_index.en.md --
---
title: music
translationKey: foo
---
-- content/categories/música/_index.pt.md --
---
title: música
translationKey: foo
---
`

	b := Test(t, files)
	b.AssertFileContent("public/en/categories/music/index.html", "true|/pt/categories/m%C3%BAsica/|")
	b.AssertFileContent("public/pt/categories/música/index.html", "true|/en/categories/music/|")
}

// Issue #11540.
func TestTranslationKeyResourceSharing(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/sect/mybundle_en/index.en.md --
---
translationkey: "adfasdf"
title: "mybundle en"
---
-- content/sect/mybundle_en/f1.txt --
f1.en
-- content/sect/mybundle_en/f2.txt --
f2.en
-- content/sect/mybundle_nn/index.nn.md --
---
translationkey: "adfasdf"
title: "mybundle nn"
---
-- content/sect/mybundle_nn/f2.nn.txt --
f2.nn
-- layouts/_default/single.html --
Title: {{ .Title }}|TranslationKey: {{ .TranslationKey }}|
Resources: {{ range .Resources }}{{ .RelPermalink }}|{{ .Content }}|{{ end }}|

`
	b := Test(t, files)
	b.AssertFileContent("public/en/sect/mybundle_en/index.html",
		"TranslationKey: adfasdf|",
		"Resources: /en/sect/mybundle_en/f1.txt|f1.en|/en/sect/mybundle_en/f2.txt|f2.en||",
	)

	b.AssertFileContent("public/nn/sect/mybundle_nn/index.html",
		"TranslationKey: adfasdf|",
		"Title: mybundle nn|TranslationKey: adfasdf|\nResources: /en/sect/mybundle_en/f1.txt|f1.en|/nn/sect/mybundle_nn/f2.nn.txt|f2.nn||",
	)
}

func TestChompBOM(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	const utf8BOM = "\xef\xbb\xbf"

	cfg, fs := newTestCfg()
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "simple.md"), utf8BOM+simplePage)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	p := s.RegularPages()[0]

	checkPageTitle(t, p, "Simple")
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
	past := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	future := time.Date(2037, 11, 17, 20, 34, 58, 651387237, time.UTC)
	zero := time.Time{}

	publishSettings := []struct {
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

func TestShouldBuildWithClock(t *testing.T) {
	htime.Clock = clocks.Start(time.Date(2021, 11, 17, 20, 34, 58, 651387237, time.UTC))
	t.Cleanup(func() { htime.Clock = clocks.System() })
	past := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	future := time.Date(2037, 11, 17, 20, 34, 58, 651387237, time.UTC)
	zero := time.Time{}

	publishSettings := []struct {
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
			t.Errorf("AssertShouldBuildWithClock unexpected output with params: %+v", ps)
		}
	}
}

// See https://github.com/gohugoio/hugo/issues/9171
// We redefined disablePathToLower in v0.121.0.
func TestPagePathDisablePathToLower(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "http://example.com"
disablePathToLower = true
[permalinks]
sect2 = "/:section/:filename/"
sect3 = "/:section/:title/"
-- content/sect/p1.md --
---
title: "Page1"
---
p1.
-- content/sect/p2.md --
---
title: "Page2"
slug: "PaGe2"
---
p2.
-- content/sect2/PaGe3.md --
---
title: "Page3"
---
-- content/seCt3/p4.md --
---
title: "Pag.E4"
slug: "PaGe4"
---
p4.
-- layouts/_default/single.html --
Single: {{ .Title}}|{{ .RelPermalink }}|{{ .Path }}|
`
	b := Test(t, files)
	b.AssertFileContent("public/sect/p1/index.html", "Single: Page1|/sect/p1/|/sect/p1")
	b.AssertFileContent("public/sect/PaGe2/index.html", "Single: Page2|/sect/PaGe2/|/sect/p2")
	b.AssertFileContent("public/sect2/PaGe3/index.html", "Single: Page3|/sect2/PaGe3/|/sect2/page3|")
	b.AssertFileContent("public/sect3/Pag.E4/index.html", "Single: Pag.E4|/sect3/Pag.E4/|/sect3/p4|")
}

func TestScratch(t *testing.T) {
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
<code class="language-bash" data-lang="bash"><span class="line hl"><span class="cl">SHORT
<code class="language-bash" data-lang="bash"><span class="line hl"><span class="cl">MARKDOWN
<p><a href="https://google.com">https://google.com</a></p>
`)
}

func TestPageHashString(t *testing.T) {
	files := `
-- config.toml --
baseURL = "https://example.org"
[languages]
[languages.en]
weight = 1
title = "English"
[languages.no]
weight = 2
title = "Norsk"
-- content/p1.md --
---
title: "p1"
---
-- content/p2.md --
---
title: "p2"
---
`

	b := NewIntegrationTestBuilder(IntegrationTestConfig{
		T:           t,
		TxtarString: files,
	}).Build()

	p1 := b.H.Sites[0].RegularPages()[0]
	p2 := b.H.Sites[0].RegularPages()[1]
	sites := p1.Sites()

	b.Assert(p1, qt.Not(qt.Equals), p2)

	b.Assert(hashing.HashString(p1), qt.Not(qt.Equals), hashing.HashString(p2))
	b.Assert(hashing.HashString(sites[0]), qt.Not(qt.Equals), hashing.HashString(sites[1]))
}

// Issue #11243
func TestRenderWithoutArgument(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/index.html --
{{ .Render }}
`

	b, err := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.IsNotNil)
}
