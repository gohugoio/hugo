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
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var emptyPage = ""

const (
	simplePage                           = "---\ntitle: Simple\n---\nSimple Page\n"
	invalidFrontMatterMissing            = "This is a test"
	renderNoFrontmatter                  = "<!doctype><html><head></head><body>This is a test</body></html>"
	contentWithCommentedFrontmatter      = "<!--\n+++\ntitle = \"Network configuration\"\ndescription = \"Docker networking\"\nkeywords = [\"network\"]\n[menu.main]\nparent= \"smn_administrate\"\n+++\n-->\n\n# Network configuration\n\n##\nSummary"
	contentWithCommentedTextFrontmatter  = "<!--[metaData]>\n+++\ntitle = \"Network configuration\"\ndescription = \"Docker networking\"\nkeywords = [\"network\"]\n[menu.main]\nparent= \"smn_administrate\"\n+++\n<![end-metadata]-->\n\n# Network configuration\n\n##\nSummary"
	contentWithCommentedLongFrontmatter  = "<!--[metaData123456789012345678901234567890]>\n+++\ntitle = \"Network configuration\"\ndescription = \"Docker networking\"\nkeywords = [\"network\"]\n[menu.main]\nparent= \"smn_administrate\"\n+++\n<![end-metadata]-->\n\n# Network configuration\n\n##\nSummary"
	contentWithCommentedLong2Frontmatter = "<!--[metaData]>\n+++\ntitle = \"Network configuration\"\ndescription = \"Docker networking\"\nkeywords = [\"network\"]\n[menu.main]\nparent= \"smn_administrate\"\n+++\n<![end-metadata123456789012345678901234567890]-->\n\n# Network configuration\n\n##\nSummary"
	invalidFrontmatterShortDelim         = `
--
title: Short delim start
---
Short Delim
`

	invalidFrontmatterShortDelimEnding = `
---
title: Short delim ending
--
Short Delim
`

	invalidFrontmatterLadingWs = `

 ---
title: Leading WS
---
Leading
`

	simplePageJSON = `
{
"title": "spf13-vim 3.0 release and new website",
"description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim.",
"tags": [ ".vimrc", "plugins", "spf13-vim", "VIm" ],
"date": "2012-04-06",
"categories": [
    "Development",
    "VIM"
],
"slug": "-spf13-vim-3-0-release-and-new-website-"
}

Content of the file goes Here
`
	simplePageJSONLoose = `
{
"title": "spf13-vim 3.0 release and new website"
"description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim."
"tags": [ ".vimrc", "plugins", "spf13-vim", "VIm" ]
"date": "2012-04-06"
"categories": [
    "Development"
    "VIM"
],
"slug": "spf13-vim-3-0-release-and-new-website"
}

Content of the file goes Here
`
	simplePageRFC3339Date  = "---\ntitle: RFC3339 Date\ndate: \"2013-05-17T16:59:30Z\"\n---\nrfc3339 content"
	simplePageJSONMultiple = `
{
	"title": "foobar",
	"customData": { "foo": "bar" },
	"date": "2012-08-06"
}
Some text
`

	simplePageNoLayout = `---
title: simple_no_layout
---
No Layout called out`

	simplePageLayoutFoobar = `---
title: simple layout foobar
layout: foobar
---
Layout foobar`

	simplePageTypeFoobar = `---
type: foobar
---
type foobar`

	simplePageTypeLayout = `---
type: barfoo
layout: buzfoo
---
type and layout set`

	simplePageWithSummaryDelimiter = `---
title: Simple
---
Summary Next Line

<!--more-->
Some more text
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

	simplePageWithEmbeddedScript = `---
title: Simple
---
<script type='text/javascript'>alert('the script tags are still there, right?');</script>
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

var pageWithVariousFrontmatterTypes = `+++
a_string = "bar"
an_integer = 1
a_float = 1.3
a_bool = false
a_date = 1979-05-27T07:32:00Z

[a_table]
a_key = "a_value"
+++
Front Matter with various frontmatter types`

var pageWithCalendarYAMLFrontmatter = `---
type: calendar
weeks:
  -
    start: "Jan 5"
    days:
      - activity: class
        room: EN1000
      - activity: lab
      - activity: class
      - activity: lab
      - activity: class
  -
    start: "Jan 12"
    days:
      - activity: class
      - activity: lab
      - activity: class
      - activity: lab
      - activity: exam
---

Hi.
`

var pageWithCalendarJSONFrontmatter = `{
  "type": "calendar",
  "weeks": [
    {
      "start": "Jan 5",
      "days": [
        { "activity": "class", "room": "EN1000" },
        { "activity": "lab" },
        { "activity": "class" },
        { "activity": "lab" },
        { "activity": "class" }
      ]
    },
    {
      "start": "Jan 12",
      "days": [
        { "activity": "class" },
        { "activity": "lab" },
        { "activity": "class" },
        { "activity": "lab" },
        { "activity": "exam" }
      ]
    }
  ]
}

Hi.
`

var pageWithCalendarTOMLFrontmatter = `+++
type = "calendar"

[[weeks]]
start = "Jan 5"

[[weeks.days]]
activity = "class"
room = "EN1000"

[[weeks.days]]
activity = "lab"

[[weeks.days]]
activity = "class"

[[weeks.days]]
activity = "lab"

[[weeks.days]]
activity = "class"

[[weeks]]
start = "Jan 12"

[[weeks.days]]
activity = "class"

[[weeks.days]]
activity = "lab"

[[weeks.days]]
activity = "class"

[[weeks.days]]
activity = "lab"

[[weeks.days]]
activity = "exam"
+++

Hi.
`

func checkError(t *testing.T, err error, expected string) {
	if err == nil {
		t.Fatalf("err is nil.  Expected: %s", expected)
	}
	if err.Error() != expected {
		t.Errorf("err.Error() returned: '%s'.  Expected: '%s'", err.Error(), expected)
	}
}

func TestDegenerateEmptyPageZeroLengthName(t *testing.T) {

	_, err := NewPage("")
	if err == nil {
		t.Fatalf("A zero length page name must return an error")
	}

	checkError(t, err, "Zero length page name")
}

func TestDegenerateEmptyPage(t *testing.T) {
	_, err := NewPageFrom(strings.NewReader(emptyPage), "test")
	if err != nil {
		t.Fatalf("Empty files should not trigger an error. Should be able to touch a file while watching without erroring out.")
	}
}

func checkPageTitle(t *testing.T, page *Page, title string) {
	if page.Title != title {
		t.Fatalf("Page title is: %s.  Expected %s", page.Title, title)
	}
}

func checkPageContent(t *testing.T, page *Page, content string, msg ...interface{}) {
	a := normalizeContent(content)
	b := normalizeContent(string(page.Content))
	if a != b {
		t.Fatalf("Page content is:\n%q\nExpected:\n%q (%q)", b, a, msg)
	}
}

func normalizeContent(c string) string {
	norm := strings.Replace(c, "\n", "", -1)
	norm = strings.Replace(norm, "    ", " ", -1)
	norm = strings.Replace(norm, "   ", " ", -1)
	norm = strings.Replace(norm, "  ", " ", -1)
	return strings.TrimSpace(norm)
}

func checkPageTOC(t *testing.T, page *Page, toc string) {
	if page.TableOfContents != template.HTML(toc) {
		t.Fatalf("Page TableOfContents is: %q.\nExpected %q", page.TableOfContents, toc)
	}
}

func checkPageSummary(t *testing.T, page *Page, summary string, msg ...interface{}) {
	a := normalizeContent(string(page.Summary))
	b := normalizeContent(summary)
	if a != b {
		t.Fatalf("Page summary is:\n%q.\nExpected\n%q (%q)", a, b, msg)
	}
}

func checkPageType(t *testing.T, page *Page, pageType string) {
	if page.Type() != pageType {
		t.Fatalf("Page type is: %s.  Expected: %s", page.Type(), pageType)
	}
}

func checkPageLayout(t *testing.T, page *Page, layout ...string) {
	if !listEqual(page.layouts(), layout) {
		t.Fatalf("Page layout is: %s.  Expected: %s", page.layouts(), layout)
	}
}

func checkPageDate(t *testing.T, page *Page, time time.Time) {
	if page.Date != time {
		t.Fatalf("Page date is: %s.  Expected: %s", page.Date, time)
	}
}

func checkTruncation(t *testing.T, page *Page, shouldBe bool, msg string) {
	if page.Summary == "" {
		t.Fatal("page has no summary, can not check truncation")
	}
	if page.Truncated != shouldBe {
		if shouldBe {
			t.Fatalf("page wasn't truncated: %s", msg)
		} else {
			t.Fatalf("page was truncated: %s", msg)
		}
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

func testAllMarkdownEnginesForPage(t *testing.T,
	assertFunc func(t *testing.T, ext string, p *Page), baseFilename, pageContent string) {

	engines := []struct {
		ext           string
		shouldExecute func() bool
	}{
		{"ad", func() bool { return helpers.HasAsciidoc() }},
		{"md", func() bool { return true }},
		{"mmark", func() bool { return true }},
		// TODO(bep) figure a way to include this without too much work.{"html", func() bool { return true }},
		{"rst", func() bool { return helpers.HasRst() }},
	}

	for _, e := range engines {
		if !e.shouldExecute() {
			continue
		}

		filename := baseFilename + "." + e.ext

		s := newSiteFromSources(filename, pageContent)

		if err := buildSiteSkipRender(s); err != nil {
			t.Fatalf("Failed to build site: %s", err)
		}

		require.Len(t, s.Pages, 1)

		p := s.Pages[0]

		assertFunc(t, e.ext, p)

	}

}

func TestCreateNewPage(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, p *Page) {
		assert.False(t, p.IsHome)
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Simple Page</p>\n"))
		checkPageSummary(t, p, "Simple Page")
		checkPageType(t, p, "page")
		checkPageLayout(t, p, "page/single.html", "_default/single.html", "theme/page/single.html", "theme/_default/single.html")
		checkTruncation(t, p, false, "simple short page")
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePage)
}

func TestSplitSummaryAndContent(t *testing.T) {

	for i, this := range []struct {
		markup                        string
		content                       string
		expectedSummary               string
		expectedContent               string
		expectedContentWithoutSummary string
	}{
		{"markdown", `<p>Summary Same LineHUGOMORE42</p>

<p>Some more text</p>`, "<p>Summary Same Line</p>", "<p>Summary Same Line</p>\n\n<p>Some more text</p>", "<p>Some more text</p>"},
		{"asciidoc", `<div class="paragraph"><p>sn</p></div><div class="paragraph"><p>HUGOMORE42Some more text</p></div>`,
			"<div class=\"paragraph\"><p>sn</p></div>",
			"<div class=\"paragraph\"><p>sn</p></div><div class=\"paragraph\"><p>Some more text</p></div>",
			"<div class=\"paragraph\"><p>Some more text</p></div>"},
		{"rst",
			"<div class=\"document\"><p>Summary Next Line</p><p>HUGOMORE42Some more text</p></div>",
			"<div class=\"document\"><p>Summary Next Line</p></div>",
			"<div class=\"document\"><p>Summary Next Line</p><p>Some more text</p></div>",
			"<div class=\"document\"><p>Some more text</p></div>"},
		{"markdown", "<p>a</p><p>b</p><p>HUGOMORE42c</p>", "<p>a</p><p>b</p>", "<p>a</p><p>b</p><p>c</p>", "<p>c</p>"},
		{"markdown", "<p>a</p><p>b</p><p>cHUGOMORE42</p>", "<p>a</p><p>b</p><p>c</p>", "<p>a</p><p>b</p><p>c</p>", ""},
		{"markdown", "<p>a</p><p>bHUGOMORE42</p><p>c</p>", "<p>a</p><p>b</p>", "<p>a</p><p>b</p><p>c</p>", "<p>c</p>"},
		{"markdown", "<p>aHUGOMORE42</p><p>b</p><p>c</p>", "<p>a</p>", "<p>a</p><p>b</p><p>c</p>", "<p>b</p><p>c</p>"},
		{"markdown", "  HUGOMORE42 ", "", "", ""},
		{"markdown", "HUGOMORE42", "", "", ""},
		{"markdown", "<p>HUGOMORE42", "<p>", "<p>", ""},
		{"markdown", "HUGOMORE42<p>", "", "<p>", "<p>"},
		{"markdown", "\n\n<p>HUGOMORE42</p>\n", "<p></p>", "<p></p>", ""},
	} {

		sc := splitUserDefinedSummaryAndContent(this.markup, []byte(this.content))

		require.NotNil(t, sc, fmt.Sprintf("[%d] Nil %s", i, this.markup))
		require.Equal(t, this.expectedSummary, string(sc.summary), fmt.Sprintf("[%d] Summary markup %s", i, this.markup))
		require.Equal(t, this.expectedContent, string(sc.content), fmt.Sprintf("[%d] Content markup %s", i, this.markup))
		require.Equal(t, this.expectedContentWithoutSummary, string(sc.contentWithoutSummary), fmt.Sprintf("[%d] Content without summary, markup %s", i, this.markup))
	}
}

func TestPageWithDelimiter(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, p *Page) {
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Next Line</p>\n\n<p>Some more text</p>\n"), ext)
		checkPageSummary(t, p, normalizeExpected(ext, "<p>Summary Next Line</p>"), ext)
		checkPageType(t, p, "page")
		checkPageLayout(t, p, "page/single.html", "_default/single.html", "theme/page/single.html", "theme/_default/single.html")
		checkTruncation(t, p, true, "page with summary delimiter")
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithSummaryDelimiter)
}

// Issue #1076
func TestPageWithDelimiterForMarkdownThatCrossesBorder(t *testing.T) {
	s := newSiteFromSources("simple.md", simplePageWithSummaryDelimiterAndMarkdownThatCrossesBorder)

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	require.Len(t, s.Pages, 1)

	p := s.Pages[0]

	if p.Summary != template.HTML("<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup class=\"footnote-ref\" id=\"fnref:1\"><a rel=\"footnote\" href=\"#fn:1\">1</a></sup>\n</p>") {
		t.Fatalf("Got summary:\n%q", p.Summary)
	}

	if p.Content != template.HTML("<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup class=\"footnote-ref\" id=\"fnref:1\"><a rel=\"footnote\" href=\"#fn:1\">1</a></sup>\n</p>\n<div class=\"footnotes\">\n\n<hr />\n\n<ol>\n<li id=\"fn:1\">Many people say so.\n <a class=\"footnote-return\" href=\"#fnref:1\"><sup>[return]</sup></a></li>\n</ol>\n</div>") {
		t.Fatalf("Got content:\n%q", p.Content)
	}
}

func TestPageWithShortCodeInSummary(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, p *Page) {
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Next Line. <figure > <img src=\"/not/real\" /> </figure>.\nMore text here.</p><p>Some more text</p>"), ext)
		checkPageSummary(t, p, "Summary Next Line. . More text here. Some more text", ext)
		checkPageType(t, p, "page")
		checkPageLayout(t, p, "page/single.html", "_default/single.html", "theme/page/single.html", "theme/_default/single.html")
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithShortcodeInSummary)
}

func TestPageWithEmbeddedScriptTag(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, p *Page) {
		if ext == "ad" || ext == "rst" {
			// TOD(bep)
			return
		}
		checkPageContent(t, p, "<script type='text/javascript'>alert('the script tags are still there, right?');</script>\n", ext)
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithEmbeddedScript)
}

func TestPageWithAdditionalExtension(t *testing.T) {
	s := newSiteFromSources("simple.md", simplePageWithAdditionalExtension)

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	require.Len(t, s.Pages, 1)

	p := s.Pages[0]

	checkPageContent(t, p, "<p>first line.<br />\nsecond line.</p>\n\n<p>fourth line.</p>\n")
}

func TestTableOfContents(t *testing.T) {
	s := newSiteFromSources("tocpage.md", pageWithToC)

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	require.Len(t, s.Pages, 1)

	p := s.Pages[0]

	checkPageContent(t, p, "\n\n<p>For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.</p>\n\n<h2 id=\"aa\">AA</h2>\n\n<p>I have no idea, of course, how long it took me to reach the limit of the plain,\nbut at last I entered the foothills, following a pretty little canyon upward\ntoward the mountains. Beside me frolicked a laughing brooklet, hurrying upon\nits noisy way down to the silent sea. In its quieter pools I discovered many\nsmall fish, of four-or five-pound weight I should imagine. In appearance,\nexcept as to size and color, they were not unlike the whale of our own seas. As\nI watched them playing about I discovered, not only that they suckled their\nyoung, but that at intervals they rose to the surface to breathe as well as to\nfeed upon certain grasses and a strange, scarlet lichen which grew upon the\nrocks just above the water line.</p>\n\n<h3 id=\"aaa\">AAA</h3>\n\n<p>I remember I felt an extraordinary persuasion that I was being played with,\nthat presently, when I was upon the very verge of safety, this mysterious\ndeath&ndash;as swift as the passage of light&ndash;would leap after me from the pit about\nthe cylinder and strike me down. ## BB</p>\n\n<h3 id=\"bbb\">BBB</h3>\n\n<p>&ldquo;You&rsquo;re a great Granser,&rdquo; he cried delightedly, &ldquo;always making believe them little marks mean something.&rdquo;</p>\n")
	checkPageTOC(t, p, "<nav id=\"TableOfContents\">\n<ul>\n<li>\n<ul>\n<li><a href=\"#aa\">AA</a>\n<ul>\n<li><a href=\"#aaa\">AAA</a></li>\n<li><a href=\"#bbb\">BBB</a></li>\n</ul></li>\n</ul></li>\n</ul>\n</nav>")
}

func TestPageWithMoreTag(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, p *Page) {
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Same Line</p>\n\n<p>Some more text</p>\n"))
		checkPageSummary(t, p, normalizeExpected(ext, "<p>Summary Same Line</p>"))
		checkPageType(t, p, "page")
		checkPageLayout(t, p, "page/single.html", "_default/single.html", "theme/page/single.html", "theme/_default/single.html")
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithSummaryDelimiterSameLine)
}

func TestPageWithDate(t *testing.T) {
	s := newSiteFromSources("simple.md", simplePageRFC3339Date)

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	require.Len(t, s.Pages, 1)

	p := s.Pages[0]
	d, _ := time.Parse(time.RFC3339, "2013-05-17T16:59:30Z")

	checkPageDate(t, p, d)
}

func TestWordCountWithAllCJKRunesWithoutHasCJKLanguage(t *testing.T) {
	testCommonResetState()

	assertFunc := func(t *testing.T, ext string, p *Page) {
		if p.WordCount != 8 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 8, p.WordCount)
		}
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithAllCJKRunes)
}

func TestWordCountWithAllCJKRunesHasCJKLanguage(t *testing.T) {
	testCommonResetState()
	viper.Set("HasCJKLanguage", true)

	assertFunc := func(t *testing.T, ext string, p *Page) {
		if p.WordCount != 15 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 15, p.WordCount)
		}
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithAllCJKRunes)
}

func TestWordCountWithMainEnglishWithCJKRunes(t *testing.T) {
	testCommonResetState()

	viper.Set("HasCJKLanguage", true)

	assertFunc := func(t *testing.T, ext string, p *Page) {
		if p.WordCount != 74 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 74, p.WordCount)
		}

		if p.Summary != simplePageWithMainEnglishWithCJKRunesSummary {
			t.Fatalf("[%s] incorrect Summary for content '%s'. expected %v, got %v", ext, p.plain,
				simplePageWithMainEnglishWithCJKRunesSummary, p.Summary)
		}

	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithMainEnglishWithCJKRunes)
}

func TestWordCountWithIsCJKLanguageFalse(t *testing.T) {
	testCommonResetState()
	viper.Set("HasCJKLanguage", true)

	assertFunc := func(t *testing.T, ext string, p *Page) {
		if p.WordCount != 75 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 74, p.WordCount)
		}

		if p.Summary != simplePageWithIsCJKLanguageFalseSummary {
			t.Fatalf("[%s] incorrect Summary for content '%s'. expected %v, got %v", ext, p.plain,
				simplePageWithIsCJKLanguageFalseSummary, p.Summary)
		}

	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithIsCJKLanguageFalse)

}

func TestWordCount(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, p *Page) {
		if p.WordCount != 483 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 483, p.WordCount)
		}

		if p.FuzzyWordCount != 500 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 500, p.WordCount)
		}

		if p.ReadingTime != 3 {
			t.Fatalf("[%s] incorrect min read. expected %v, got %v", ext, 3, p.ReadingTime)
		}

		checkTruncation(t, p, true, "long page")
	}

	testAllMarkdownEnginesForPage(t, assertFunc, "simple", simplePageWithLongContent)
}

func TestCreatePage(t *testing.T) {
	var tests = []struct {
		r string
	}{
		{simplePageJSON},
		{simplePageJSONLoose},
		{simplePageJSONMultiple},
		//{strings.NewReader(SIMPLE_PAGE_JSON_COMPACT)},
	}

	for _, test := range tests {
		p, _ := NewPage("page")
		if _, err := p.ReadFrom(strings.NewReader(test.r)); err != nil {
			t.Errorf("Unable to parse page: %s", err)
		}
	}
}

func TestDegenerateInvalidFrontMatterShortDelim(t *testing.T) {
	var tests = []struct {
		r   string
		err string
	}{
		{invalidFrontmatterShortDelimEnding, "unable to read frontmatter at filepos 45: EOF"},
	}
	for _, test := range tests {

		p, _ := NewPage("invalid/front/matter/short/delim")
		_, err := p.ReadFrom(strings.NewReader(test.r))
		checkError(t, err, test.err)
	}
}

func TestShouldRenderContent(t *testing.T) {
	var tests = []struct {
		text   string
		render bool
	}{
		{invalidFrontMatterMissing, true},
		// TODO how to deal with malformed frontmatter.  In this case it'll be rendered as markdown.
		{invalidFrontmatterShortDelim, true},
		{renderNoFrontmatter, false},
		{contentWithCommentedFrontmatter, true},
		{contentWithCommentedTextFrontmatter, true},
		{contentWithCommentedLongFrontmatter, false},
		{contentWithCommentedLong2Frontmatter, true},
	}

	for _, test := range tests {

		p, _ := NewPage("render/front/matter")
		_, err := p.ReadFrom(strings.NewReader(test.text))
		p = pageMust(p, err)
		if p.IsRenderable() != test.render {
			t.Errorf("expected p.IsRenderable() == %t, got %t", test.render, p.IsRenderable())
		}
	}
}

// Issue #768
func TestCalendarParamsVariants(t *testing.T) {
	pageJSON, _ := NewPage("test/fileJSON.md")
	_, _ = pageJSON.ReadFrom(strings.NewReader(pageWithCalendarJSONFrontmatter))

	pageYAML, _ := NewPage("test/fileYAML.md")
	_, _ = pageYAML.ReadFrom(strings.NewReader(pageWithCalendarYAMLFrontmatter))

	pageTOML, _ := NewPage("test/fileTOML.md")
	_, _ = pageTOML.ReadFrom(strings.NewReader(pageWithCalendarTOMLFrontmatter))

	assert.True(t, compareObjects(pageJSON.Params, pageYAML.Params))
	assert.True(t, compareObjects(pageJSON.Params, pageTOML.Params))

}

func TestDifferentFrontMatterVarTypes(t *testing.T) {
	page, _ := NewPage("test/file1.md")
	_, _ = page.ReadFrom(strings.NewReader(pageWithVariousFrontmatterTypes))

	dateval, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
	if page.GetParam("a_string") != "bar" {
		t.Errorf("frontmatter not handling strings correctly should be %s, got: %s", "bar", page.GetParam("a_string"))
	}
	if page.GetParam("an_integer") != 1 {
		t.Errorf("frontmatter not handling ints correctly should be %s, got: %s", "1", page.GetParam("an_integer"))
	}
	if page.GetParam("a_float") != 1.3 {
		t.Errorf("frontmatter not handling floats correctly should be %f, got: %s", 1.3, page.GetParam("a_float"))
	}
	if page.GetParam("a_bool") != false {
		t.Errorf("frontmatter not handling bools correctly should be %t, got: %s", false, page.GetParam("a_bool"))
	}
	if page.GetParam("a_date") != dateval {
		t.Errorf("frontmatter not handling dates correctly should be %s, got: %s", dateval, page.GetParam("a_date"))
	}
	param := page.GetParam("a_table")
	if param == nil {
		t.Errorf("frontmatter not handling tables correctly should be type of %v, got: type of %v", reflect.TypeOf(page.Params["a_table"]), reflect.TypeOf(param))
	}
	if cast.ToStringMap(param)["a_key"] != "a_value" {
		t.Errorf("frontmatter not handling values inside a table correctly should be %s, got: %s", "a_value", cast.ToStringMap(page.Params["a_table"])["a_key"])
	}
}

func TestDegenerateInvalidFrontMatterLeadingWhitespace(t *testing.T) {
	p, _ := NewPage("invalid/front/matter/leading/ws")
	_, err := p.ReadFrom(strings.NewReader(invalidFrontmatterLadingWs))
	if err != nil {
		t.Fatalf("Unable to parse front matter given leading whitespace: %s", err)
	}
}

func TestSectionEvaluation(t *testing.T) {
	page, _ := NewPage(filepath.FromSlash("blue/file1.md"))
	page.ReadFrom(strings.NewReader(simplePage))
	if page.Section() != "blue" {
		t.Errorf("Section should be %s, got: %s", "blue", page.Section())
	}
}

func L(s ...string) []string {
	return s
}

func TestLayoutOverride(t *testing.T) {
	var (
		pathContentTwoDir = filepath.Join("content", "dub", "sub", "file1.md")
		pathContentOneDir = filepath.Join("content", "gub", "file1.md")
		pathContentNoDir  = filepath.Join("content", "file1")
		pathOneDirectory  = filepath.Join("fub", "file1.md")
		pathNoDirectory   = filepath.Join("file1.md")
	)
	tests := []struct {
		content        string
		path           string
		expectedLayout []string
	}{
		{simplePageNoLayout, pathContentTwoDir, L("dub/single.html", "_default/single.html")},
		{simplePageNoLayout, pathContentOneDir, L("gub/single.html", "_default/single.html")},
		{simplePageNoLayout, pathContentNoDir, L("page/single.html", "_default/single.html")},
		{simplePageNoLayout, pathOneDirectory, L("fub/single.html", "_default/single.html")},
		{simplePageNoLayout, pathNoDirectory, L("page/single.html", "_default/single.html")},
		{simplePageLayoutFoobar, pathContentTwoDir, L("dub/foobar.html", "_default/foobar.html")},
		{simplePageLayoutFoobar, pathContentOneDir, L("gub/foobar.html", "_default/foobar.html")},
		{simplePageLayoutFoobar, pathOneDirectory, L("fub/foobar.html", "_default/foobar.html")},
		{simplePageLayoutFoobar, pathNoDirectory, L("page/foobar.html", "_default/foobar.html")},
		{simplePageTypeFoobar, pathContentTwoDir, L("foobar/single.html", "_default/single.html")},
		{simplePageTypeFoobar, pathContentOneDir, L("foobar/single.html", "_default/single.html")},
		{simplePageTypeFoobar, pathContentNoDir, L("foobar/single.html", "_default/single.html")},
		{simplePageTypeFoobar, pathOneDirectory, L("foobar/single.html", "_default/single.html")},
		{simplePageTypeFoobar, pathNoDirectory, L("foobar/single.html", "_default/single.html")},
		{simplePageTypeLayout, pathContentTwoDir, L("barfoo/buzfoo.html", "_default/buzfoo.html")},
		{simplePageTypeLayout, pathContentOneDir, L("barfoo/buzfoo.html", "_default/buzfoo.html")},
		{simplePageTypeLayout, pathContentNoDir, L("barfoo/buzfoo.html", "_default/buzfoo.html")},
		{simplePageTypeLayout, pathOneDirectory, L("barfoo/buzfoo.html", "_default/buzfoo.html")},
		{simplePageTypeLayout, pathNoDirectory, L("barfoo/buzfoo.html", "_default/buzfoo.html")},
	}
	for _, test := range tests {
		p, _ := NewPage(test.path)
		_, err := p.ReadFrom(strings.NewReader(test.content))
		if err != nil {
			t.Fatalf("Unable to parse content:\n%s\n", test.content)
		}

		for _, y := range test.expectedLayout {
			test.expectedLayout = append(test.expectedLayout, "theme/"+y)
		}
		if !listEqual(p.layouts(), test.expectedLayout) {
			t.Errorf("Layout mismatch. Expected: %s, got: %s", test.expectedLayout, p.layouts())
		}
	}
}

func TestSliceToLower(t *testing.T) {
	tests := []struct {
		value    []string
		expected []string
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "B", "c"}, []string{"a", "b", "c"}},
		{[]string{"A", "B", "C"}, []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		res := helpers.SliceToLower(test.value)
		for i, val := range res {
			if val != test.expected[i] {
				t.Errorf("Case mismatch. Expected %s, got %s", test.expected[i], res[i])
			}
		}
	}
}

func TestPagePaths(t *testing.T) {
	testCommonResetState()

	viper.Set("DefaultExtension", "html")
	siteParmalinksSetting := PermalinkOverrides{
		"post": ":year/:month/:day/:title/",
	}

	tests := []struct {
		content      string
		path         string
		hasPermalink bool
		expected     string
	}{
		{simplePage, "content/post/x.md", false, "content/post/x.html"},
		{simplePageWithURL, "content/post/x.md", false, "simple/url/index.html"},
		{simplePageWithSlug, "content/post/x.md", false, "content/post/simple-slug.html"},
		{simplePageWithDate, "content/post/x.md", true, "2013/10/15/simple/index.html"},
		{UTF8Page, "content/post/x.md", false, "content/post/x.html"},
		{UTF8PageWithURL, "content/post/x.md", false, "ラーメン/url/index.html"},
		{UTF8PageWithSlug, "content/post/x.md", false, "content/post/ラーメン-slug.html"},
		{UTF8PageWithDate, "content/post/x.md", true, "2013/10/15/ラーメン/index.html"},
	}

	for _, test := range tests {
		p, _ := NewPageFrom(strings.NewReader(test.content), filepath.FromSlash(test.path))
		p.Node.Site = newSiteInfoDefaultLanguage("")

		if test.hasPermalink {
			p.Node.Site.Permalinks = siteParmalinksSetting
		}

		expectedTargetPath := filepath.FromSlash(test.expected)
		expectedFullFilePath := filepath.FromSlash(test.path)

		if p.TargetPath() != expectedTargetPath {
			t.Errorf("%s => TargetPath  expected: '%s', got: '%s'", test.content, expectedTargetPath, p.TargetPath())
		}

		if p.FullFilePath() != expectedFullFilePath {
			t.Errorf("%s => FullFilePath  expected: '%s', got: '%s'", test.content, expectedFullFilePath, p.FullFilePath())
		}
	}
}

var pageWithDraftAndPublished = `---
title: broken
published: false
draft: true
---
some content
`

func TestDraftAndPublishedFrontMatterError(t *testing.T) {
	_, err := NewPageFrom(strings.NewReader(pageWithDraftAndPublished), "content/post/broken.md")
	if err != ErrHasDraftAndPublished {
		t.Errorf("expected ErrHasDraftAndPublished, was %#v", err)
	}
}

var pageWithPublishedFalse = `---
title: okay
published: false
---
some content
`
var pageWithPublishedTrue = `---
title: okay
published: true
---
some content
`

func TestPublishedFrontMatter(t *testing.T) {
	p, err := NewPageFrom(strings.NewReader(pageWithPublishedFalse), "content/post/broken.md")
	if err != nil {
		t.Fatalf("err during parse: %s", err)
	}
	if !p.Draft {
		t.Errorf("expected true, got %t", p.Draft)
	}
	p, err = NewPageFrom(strings.NewReader(pageWithPublishedTrue), "content/post/broken.md")
	if err != nil {
		t.Fatalf("err during parse: %s", err)
	}
	if p.Draft {
		t.Errorf("expected false, got %t", p.Draft)
	}
}

func TestPageSimpleMethods(t *testing.T) {
	for i, this := range []struct {
		assertFunc func(p *Page) bool
	}{
		{func(p *Page) bool { return !p.IsNode() }},
		{func(p *Page) bool { return p.IsPage() }},
		{func(p *Page) bool { return p.Plain() == "Do Be Do Be Do" }},
		{func(p *Page) bool { return strings.Join(p.PlainWords(), " ") == "Do Be Do Be Do" }},
	} {

		p, _ := NewPage("Test")
		p.Content = "<h1>Do Be Do Be Do</h1>"
		if !this.assertFunc(p) {
			t.Errorf("[%d] Page method error", i)
		}
	}
}

func TestChompBOM(t *testing.T) {
	const utf8BOM = "\xef\xbb\xbf"

	s := newSiteFromSources("simple.md", utf8BOM+simplePage)

	if err := buildSiteSkipRender(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	require.Len(t, s.Pages, 1)

	p := s.Pages[0]

	checkPageTitle(t, p, "Simple")
}

func listEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}

// TODO(bep) this may be useful for other tests.
func compareObjects(a interface{}, b interface{}) bool {
	aStr := strings.Split(fmt.Sprintf("%v", a), "")
	sort.Strings(aStr)

	bStr := strings.Split(fmt.Sprintf("%v", b), "")
	sort.Strings(bStr)

	return strings.Join(aStr, "") == strings.Join(bStr, "")
}

func TestShouldBuild(t *testing.T) {
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
