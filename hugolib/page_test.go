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
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cast"
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

	simplePageRFC3339Date  = "---\ntitle: RFC3339 Date\ndate: \"2013-05-17T16:59:30Z\"\n---\nrfc3339 content"
	simplePageJSONMultiple = `
{
	"title": "foobar",
	"customData": { "foo": "bar" },
	"date": "2012-08-06"
}
Some text
`

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

	simplePageWithSummaryDelimiterOnlySummary = `---
title: Simple
---
Summary text

<!--more-->
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
	t.Parallel()
	s := newTestSite(t)
	_, err := s.NewPage("")
	if err == nil {
		t.Fatalf("A zero length page name must return an error")
	}

	checkError(t, err, "Zero length page name")
}

func TestDegenerateEmptyPage(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	_, err := s.NewPageFrom(strings.NewReader(emptyPage), "test")
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
	norm := c
	norm = strings.Replace(norm, "\n", " ", -1)
	norm = strings.Replace(norm, "    ", " ", -1)
	norm = strings.Replace(norm, "   ", " ", -1)
	norm = strings.Replace(norm, "  ", " ", -1)
	norm = strings.Replace(norm, "p> ", "p>", -1)
	norm = strings.Replace(norm, ">  <", "> <", -1)
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

func testAllMarkdownEnginesForPages(t *testing.T,
	assertFunc func(t *testing.T, ext string, pages Pages), settings map[string]interface{}, pageSources ...string) {

	engines := []struct {
		ext           string
		shouldExecute func() bool
	}{
		{"md", func() bool { return true }},
		{"mmark", func() bool { return true }},
		{"ad", func() bool { return helpers.HasAsciidoctor() || helpers.HasAsciidoc() }},
		// TODO(bep) figure a way to include this without too much work.{"html", func() bool { return true }},
		{"rst", func() bool { return helpers.HasRst() }},
	}

	for _, e := range engines {
		if !e.shouldExecute() {
			continue
		}

		cfg, fs := newTestCfg()

		if settings != nil {
			for k, v := range settings {
				cfg.Set(k, v)
			}
		}

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

		s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

		require.Len(t, s.RegularPages, len(pageSources))

		assertFunc(t, e.ext, s.RegularPages)

	}

}

func TestCreateNewPage(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]

		// issue #2290: Path is relative to the content dir and will continue to be so.
		require.Equal(t, filepath.FromSlash(fmt.Sprintf("p0.%s", ext)), p.Path())
		assert.False(t, p.IsHome())
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Simple Page</p>\n"))
		checkPageSummary(t, p, "Simple Page")
		checkPageType(t, p, "page")
		checkTruncation(t, p, false, "simple short page")
	}

	settings := map[string]interface{}{
		"contentDir": "mycontent",
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePage)
}

func TestSplitSummaryAndContent(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		markup          string
		content         string
		expectedSummary string
		expectedContent string
	}{
		{"markdown", `<p>Summary Same LineHUGOMORE42</p>

<p>Some more text</p>`, "<p>Summary Same Line</p>", "<p>Summary Same Line</p>\n\n<p>Some more text</p>"},
		{"asciidoc", `<div class="paragraph"><p>sn</p></div><div class="paragraph"><p>HUGOMORE42Some more text</p></div>`,
			"<div class=\"paragraph\"><p>sn</p></div>",
			"<div class=\"paragraph\"><p>sn</p></div><div class=\"paragraph\"><p>Some more text</p></div>"},
		{"rst",
			"<div class=\"document\"><p>Summary Next Line</p><p>HUGOMORE42Some more text</p></div>",
			"<div class=\"document\"><p>Summary Next Line</p></div>",
			"<div class=\"document\"><p>Summary Next Line</p><p>Some more text</p></div>"},
		{"markdown", "<p>a</p><p>b</p><p>HUGOMORE42c</p>", "<p>a</p><p>b</p>", "<p>a</p><p>b</p><p>c</p>"},
		{"markdown", "<p>a</p><p>b</p><p>cHUGOMORE42</p>", "<p>a</p><p>b</p><p>c</p>", "<p>a</p><p>b</p><p>c</p>"},
		{"markdown", "<p>a</p><p>bHUGOMORE42</p><p>c</p>", "<p>a</p><p>b</p>", "<p>a</p><p>b</p><p>c</p>"},
		{"markdown", "<p>aHUGOMORE42</p><p>b</p><p>c</p>", "<p>a</p>", "<p>a</p><p>b</p><p>c</p>"},
		{"markdown", "  HUGOMORE42 ", "", ""},
		{"markdown", "HUGOMORE42", "", ""},
		{"markdown", "<p>HUGOMORE42", "<p>", "<p>"},
		{"markdown", "HUGOMORE42<p>", "", "<p>"},
		{"markdown", "\n\n<p>HUGOMORE42</p>\n", "<p></p>", "<p></p>"},
		// Issue #2586
		// Note: Hugo will not split mid-sentence but will look for the closest
		// paragraph end marker. This may be a change from Hugo 0.16, but it makes sense.
		{"markdown", `<p>this is an example HUGOMORE42of the issue.</p>`,
			"<p>this is an example of the issue.</p>",
			"<p>this is an example of the issue.</p>"},
		// Issue: #2538
		{"markdown", fmt.Sprintf(` <p class="lead">%s</p>HUGOMORE42<p>%s</p>
`,
			strings.Repeat("A", 10), strings.Repeat("B", 31)),
			fmt.Sprintf(`<p class="lead">%s</p>`, strings.Repeat("A", 10)),
			fmt.Sprintf(`<p class="lead">%s</p><p>%s</p>`, strings.Repeat("A", 10), strings.Repeat("B", 31)),
		},
	} {

		sc, err := splitUserDefinedSummaryAndContent(this.markup, []byte(this.content))

		require.NoError(t, err)
		require.NotNil(t, sc, fmt.Sprintf("[%d] Nil %s", i, this.markup))
		require.Equal(t, this.expectedSummary, string(sc.summary), fmt.Sprintf("[%d] Summary markup %s", i, this.markup))
		require.Equal(t, this.expectedContent, string(sc.content), fmt.Sprintf("[%d] Content markup %s", i, this.markup))
	}
}

func TestPageWithDelimiter(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Next Line</p>\n\n<p>Some more text</p>\n"), ext)
		checkPageSummary(t, p, normalizeExpected(ext, "<p>Summary Next Line</p>"), ext)
		checkPageType(t, p, "page")
		checkTruncation(t, p, true, "page with summary delimiter")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithSummaryDelimiter)
}

// Issue #1076
func TestPageWithDelimiterForMarkdownThatCrossesBorder(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageWithSummaryDelimiterAndMarkdownThatCrossesBorder)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 1)

	p := s.RegularPages[0]

	if p.Summary != template.HTML("<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup class=\"footnote-ref\" id=\"fnref:1\"><a rel=\"footnote\" href=\"#fn:1\">1</a></sup>\n</p>") {
		t.Fatalf("Got summary:\n%q", p.Summary)
	}

	if p.Content != template.HTML("<p>The <a href=\"http://gohugo.io/\">best static site generator</a>.<sup class=\"footnote-ref\" id=\"fnref:1\"><a rel=\"footnote\" href=\"#fn:1\">1</a></sup>\n</p>\n<div class=\"footnotes\">\n\n<hr />\n\n<ol>\n<li id=\"fn:1\">Many people say so.\n <a class=\"footnote-return\" href=\"#fnref:1\"><sup>[return]</sup></a></li>\n</ol>\n</div>") {
		t.Fatalf("Got content:\n%q", p.Content)
	}
}

// Issue #3854
func TestPageWithDateFields(t *testing.T) {
	assert := require.New(t)
	pageWithDate := `---
title: P%d
weight: %d
%s: 2017-10-13
---
Simple Page With Some Date`

	hasBothDates := func(p *Page) bool {
		return p.Date.Year() == 2017 && p.PublishDate.Year() == 2017
	}

	datePage := func(field string, weight int) string {
		return fmt.Sprintf(pageWithDate, weight, weight, field)
	}

	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		assert.True(len(pages) > 0)
		for _, p := range pages {
			assert.True(hasBothDates(p))
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

	writeSource(t, fs, filepath.Join("content", "raw.md"), `---
title: Raw
---
**Raw**`)

	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{ .RawContent }}`)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 1)
	p := s.RegularPages[0]

	require.Contains(t, p.RawContent(), "**Raw**")

}

func TestPageWithShortCodeInSummary(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Next Line. \n<figure >\n    \n        <img src=\"/not/real\" />\n    \n    \n</figure>\n.\nMore text here.</p>\n\n<p>Some more text</p>\n"))
		checkPageSummary(t, p, "Summary Next Line.  . More text here. Some more text")
		checkPageType(t, p, "page")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithShortcodeInSummary)
}

func TestPageWithEmbeddedScriptTag(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		if ext == "ad" || ext == "rst" {
			// TOD(bep)
			return
		}
		checkPageContent(t, p, "<script type='text/javascript'>alert('the script tags are still there, right?');</script>\n", ext)
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithEmbeddedScript)
}

func TestPageWithAdditionalExtension(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageWithAdditionalExtension)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 1)

	p := s.RegularPages[0]

	checkPageContent(t, p, "<p>first line.<br />\nsecond line.</p>\n\n<p>fourth line.</p>\n")
}

func TestTableOfContents(t *testing.T) {

	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "tocpage.md"), pageWithToC)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 1)

	p := s.RegularPages[0]

	checkPageContent(t, p, "\n\n<p>For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.</p>\n\n<h2 id=\"aa\">AA</h2>\n\n<p>I have no idea, of course, how long it took me to reach the limit of the plain,\nbut at last I entered the foothills, following a pretty little canyon upward\ntoward the mountains. Beside me frolicked a laughing brooklet, hurrying upon\nits noisy way down to the silent sea. In its quieter pools I discovered many\nsmall fish, of four-or five-pound weight I should imagine. In appearance,\nexcept as to size and color, they were not unlike the whale of our own seas. As\nI watched them playing about I discovered, not only that they suckled their\nyoung, but that at intervals they rose to the surface to breathe as well as to\nfeed upon certain grasses and a strange, scarlet lichen which grew upon the\nrocks just above the water line.</p>\n\n<h3 id=\"aaa\">AAA</h3>\n\n<p>I remember I felt an extraordinary persuasion that I was being played with,\nthat presently, when I was upon the very verge of safety, this mysterious\ndeath&ndash;as swift as the passage of light&ndash;would leap after me from the pit about\nthe cylinder and strike me down. ## BB</p>\n\n<h3 id=\"bbb\">BBB</h3>\n\n<p>&ldquo;You&rsquo;re a great Granser,&rdquo; he cried delightedly, &ldquo;always making believe them little marks mean something.&rdquo;</p>\n")
	checkPageTOC(t, p, "<nav id=\"TableOfContents\">\n<ul>\n<li>\n<ul>\n<li><a href=\"#aa\">AA</a>\n<ul>\n<li><a href=\"#aaa\">AAA</a></li>\n<li><a href=\"#bbb\">BBB</a></li>\n</ul></li>\n</ul></li>\n</ul>\n</nav>")
}

func TestPageWithMoreTag(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		checkPageTitle(t, p, "Simple")
		checkPageContent(t, p, normalizeExpected(ext, "<p>Summary Same Line</p>\n\n<p>Some more text</p>\n"))
		checkPageSummary(t, p, normalizeExpected(ext, "<p>Summary Same Line</p>"))
		checkPageType(t, p, "page")

	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithSummaryDelimiterSameLine)
}

func TestPageWithMoreTagOnlySummary(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		checkTruncation(t, p, false, "page with summary delimiter at end")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithSummaryDelimiterOnlySummary)
}

// #2973
func TestSummaryWithHTMLTagsOnNextLine(t *testing.T) {

	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		require.Contains(t, p.Summary, "Happy new year everyone!")
		require.NotContains(t, p.Summary, "User interface")
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

	writeSource(t, fs, filepath.Join("content", "simple.md"), simplePageRFC3339Date)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 1)

	p := s.RegularPages[0]
	d, _ := time.Parse(time.RFC3339, "2013-05-17T16:59:30Z")

	checkPageDate(t, p, d)
}

func TestWordCountWithAllCJKRunesWithoutHasCJKLanguage(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		if p.WordCount() != 8 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 8, p.WordCount())
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithAllCJKRunes)
}

func TestWordCountWithAllCJKRunesHasCJKLanguage(t *testing.T) {
	t.Parallel()
	settings := map[string]interface{}{"hasCJKLanguage": true}

	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		if p.WordCount() != 15 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 15, p.WordCount())
		}
	}
	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithAllCJKRunes)
}

func TestWordCountWithMainEnglishWithCJKRunes(t *testing.T) {
	t.Parallel()
	settings := map[string]interface{}{"hasCJKLanguage": true}

	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		if p.WordCount() != 74 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 74, p.WordCount())
		}

		if p.Summary != simplePageWithMainEnglishWithCJKRunesSummary {
			t.Fatalf("[%s] incorrect Summary for content '%s'. expected %v, got %v", ext, p.plain,
				simplePageWithMainEnglishWithCJKRunesSummary, p.Summary)
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithMainEnglishWithCJKRunes)
}

func TestWordCountWithIsCJKLanguageFalse(t *testing.T) {
	t.Parallel()
	settings := map[string]interface{}{
		"hasCJKLanguage": true,
	}

	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		if p.WordCount() != 75 {
			t.Fatalf("[%s] incorrect word count for content '%s'. expected %v, got %v", ext, p.plain, 74, p.WordCount())
		}

		if p.Summary != simplePageWithIsCJKLanguageFalseSummary {
			t.Fatalf("[%s] incorrect Summary for content '%s'. expected %v, got %v", ext, p.plain,
				simplePageWithIsCJKLanguageFalseSummary, p.Summary)
		}
	}

	testAllMarkdownEnginesForPages(t, assertFunc, settings, simplePageWithIsCJKLanguageFalse)

}

func TestWordCount(t *testing.T) {
	t.Parallel()
	assertFunc := func(t *testing.T, ext string, pages Pages) {
		p := pages[0]
		if p.WordCount() != 483 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 483, p.WordCount())
		}

		if p.FuzzyWordCount() != 500 {
			t.Fatalf("[%s] incorrect word count. expected %v, got %v", ext, 500, p.WordCount())
		}

		if p.ReadingTime() != 3 {
			t.Fatalf("[%s] incorrect min read. expected %v, got %v", ext, 3, p.ReadingTime())
		}

		checkTruncation(t, p, true, "long page")
	}

	testAllMarkdownEnginesForPages(t, assertFunc, nil, simplePageWithLongContent)
}

func TestCreatePage(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		r string
	}{
		{simplePageJSON},
		{simplePageJSONMultiple},
		//{strings.NewReader(SIMPLE_PAGE_JSON_COMPACT)},
	}

	for i, test := range tests {
		s := newTestSite(t)
		p, _ := s.NewPage("page")
		if _, err := p.ReadFrom(strings.NewReader(test.r)); err != nil {
			t.Fatalf("[%d] Unable to parse page: %s", i, err)
		}
	}
}

func TestDegenerateInvalidFrontMatterShortDelim(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		r   string
		err string
	}{
		{invalidFrontmatterShortDelimEnding, "unable to read frontmatter at filepos 45: EOF"},
	}
	for _, test := range tests {
		s := newTestSite(t)
		p, _ := s.NewPage("invalid/front/matter/short/delim")
		_, err := p.ReadFrom(strings.NewReader(test.r))
		checkError(t, err, test.err)
	}
}

func TestShouldRenderContent(t *testing.T) {
	t.Parallel()
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
		s := newTestSite(t)
		p, _ := s.NewPage("render/front/matter")
		_, err := p.ReadFrom(strings.NewReader(test.text))
		p = pageMust(p, err)
		if p.IsRenderable() != test.render {
			t.Errorf("expected p.IsRenderable() == %t, got %t", test.render, p.IsRenderable())
		}
	}
}

// Issue #768
func TestCalendarParamsVariants(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	pageJSON, _ := s.NewPage("test/fileJSON.md")
	_, _ = pageJSON.ReadFrom(strings.NewReader(pageWithCalendarJSONFrontmatter))

	pageYAML, _ := s.NewPage("test/fileYAML.md")
	_, _ = pageYAML.ReadFrom(strings.NewReader(pageWithCalendarYAMLFrontmatter))

	pageTOML, _ := s.NewPage("test/fileTOML.md")
	_, _ = pageTOML.ReadFrom(strings.NewReader(pageWithCalendarTOMLFrontmatter))

	assert.True(t, compareObjects(pageJSON.Params, pageYAML.Params))
	assert.True(t, compareObjects(pageJSON.Params, pageTOML.Params))

}

func TestDifferentFrontMatterVarTypes(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	page, _ := s.NewPage("test/file1.md")
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
	t.Parallel()
	s := newTestSite(t)
	p, _ := s.NewPage("invalid/front/matter/leading/ws")
	_, err := p.ReadFrom(strings.NewReader(invalidFrontmatterLadingWs))
	if err != nil {
		t.Fatalf("Unable to parse front matter given leading whitespace: %s", err)
	}
}

func TestSectionEvaluation(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	page, _ := s.NewPage(filepath.FromSlash("blue/file1.md"))
	page.ReadFrom(strings.NewReader(simplePage))
	if page.Section() != "blue" {
		t.Errorf("Section should be %s, got: %s", "blue", page.Section())
	}
}

func TestSliceToLower(t *testing.T) {
	t.Parallel()
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

func TestReplaceDivider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		content           string
		from              string
		to                string
		expectedContent   string
		expectedTruncated bool
	}{
		{"none", "a", "b", "none", false},
		{"summary <!--more--> content", "<!--more-->", "HUGO", "summary HUGO content", true},
		{"summary\n\ndivider", "divider", "HUGO", "summary\n\nHUGO", false},
		{"summary\n\ndivider\n\r", "divider", "HUGO", "summary\n\nHUGO\n\r", false},
	}

	for i, test := range tests {
		replaced, truncated := replaceDivider([]byte(test.content), []byte(test.from), []byte(test.to))

		if truncated != test.expectedTruncated {
			t.Fatalf("[%d] Expected truncated to be %t, was %t", i, test.expectedTruncated, truncated)
		}

		if string(replaced) != test.expectedContent {
			t.Fatalf("[%d] Expected content to be %q, was %q", i, test.expectedContent, replaced)
		}
	}
}

func BenchmarkReplaceDivider(b *testing.B) {
	divider := "HUGO_DIVIDER"
	from, to := []byte(divider), []byte("HUGO_REPLACED")

	withDivider := make([][]byte, b.N)
	noDivider := make([][]byte, b.N)

	for i := 0; i < b.N; i++ {
		withDivider[i] = []byte(strings.Repeat("Summary ", 5) + "\n" + divider + "\n" + strings.Repeat("Word ", 300))
		noDivider[i] = []byte(strings.Repeat("Word ", 300))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, t1 := replaceDivider(withDivider[i], from, to)
		_, t2 := replaceDivider(noDivider[i], from, to)
		if !t1 {
			b.Fatal("Should be truncated")
		}
		if t2 {
			b.Fatal("Should not be truncated")
		}
	}
}

func TestPagePaths(t *testing.T) {
	t.Parallel()

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
		require.Len(t, s.RegularPages, 1)

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
	t.Parallel()
	s := newTestSite(t)
	_, err := s.NewPageFrom(strings.NewReader(pageWithDraftAndPublished), "content/post/broken.md")
	if err != ErrHasDraftAndPublished {
		t.Errorf("expected ErrHasDraftAndPublished, was %#v", err)
	}
}

var pagesWithPublishedFalse = `---
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
	t.Parallel()
	s := newTestSite(t)
	p, err := s.NewPageFrom(strings.NewReader(pagesWithPublishedFalse), "content/post/broken.md")
	if err != nil {
		t.Fatalf("err during parse: %s", err)
	}
	if !p.Draft {
		t.Errorf("expected true, got %t", p.Draft)
	}
	p, err = s.NewPageFrom(strings.NewReader(pageWithPublishedTrue), "content/post/broken.md")
	if err != nil {
		t.Fatalf("err during parse: %s", err)
	}
	if p.Draft {
		t.Errorf("expected false, got %t", p.Draft)
	}
}

var pagesDraftTemplate = []string{`---
title: "okay"
draft: %t
---
some content
`,
	`+++
title = "okay"
draft = %t
+++

some content
`,
}

func TestDraft(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	for _, draft := range []bool{true, false} {
		for i, templ := range pagesDraftTemplate {
			pageContent := fmt.Sprintf(templ, draft)
			p, err := s.NewPageFrom(strings.NewReader(pageContent), "content/post/broken.md")
			if err != nil {
				t.Fatalf("err during parse: %s", err)
			}
			if p.Draft != draft {
				t.Errorf("[%d] expected %t, got %t", i, draft, p.Draft)
			}
		}
	}
}

var pagesParamsTemplate = []string{`+++
title = "okay"
draft = false
tags = [ "hugo", "web" ]
social= [
  [ "a", "#" ],
  [ "b", "#" ],
]
+++
some content
`,
	`---
title: "okay"
draft: false
tags:
  - hugo
  - web
social:
  - - a
    - "#"
  - - b
    - "#"
---
some content
`,
	`{
	"title": "okay",
	"draft": false,
	"tags": [ "hugo", "web" ],
	"social": [
		[ "a", "#" ],
		[ "b", "#" ]
	]
}
some content
`,
}

func TestPageParams(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	wantedMap := map[string]interface{}{
		"tags": []string{"hugo", "web"},
		// Issue #2752
		"social": []interface{}{
			[]interface{}{"a", "#"},
			[]interface{}{"b", "#"},
		},
	}

	for i, c := range pagesParamsTemplate {
		p, err := s.NewPageFrom(strings.NewReader(c), "content/post/params.md")
		require.NoError(t, err, "err during parse", "#%d", i)
		for key := range wantedMap {
			assert.Equal(t, wantedMap[key], p.Params[key], "#%d", key)
		}
	}
}

func TestTraverse(t *testing.T) {
	exampleParams := `---
rating: "5 stars"
tags:
  - hugo
  - web
social:
  twitter: "@jxxf"
  facebook: "https://example.com"
---`
	t.Parallel()
	s := newTestSite(t)
	p, _ := s.NewPageFrom(strings.NewReader(exampleParams), "content/post/params.md")

	topLevelKeyValue, _ := p.Param("rating")
	assert.Equal(t, "5 stars", topLevelKeyValue)

	nestedStringKeyValue, _ := p.Param("social.twitter")
	assert.Equal(t, "@jxxf", nestedStringKeyValue)

	nonexistentKeyValue, _ := p.Param("doesn't.exist")
	assert.Nil(t, nonexistentKeyValue)
}

func TestPageSimpleMethods(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	for i, this := range []struct {
		assertFunc func(p *Page) bool
	}{
		{func(p *Page) bool { return !p.IsNode() }},
		{func(p *Page) bool { return p.IsPage() }},
		{func(p *Page) bool { return p.Plain() == "Do Be Do Be Do" }},
		{func(p *Page) bool { return strings.Join(p.PlainWords(), " ") == "Do Be Do Be Do" }},
	} {

		p, _ := s.NewPage("Test")
		p.Content = "<h1>Do Be Do Be Do</h1>"
		if !this.assertFunc(p) {
			t.Errorf("[%d] Page method error", i)
		}
	}
}

func TestIndexPageSimpleMethods(t *testing.T) {
	s := newTestSite(t)
	t.Parallel()
	for i, this := range []struct {
		assertFunc func(n *Page) bool
	}{
		{func(n *Page) bool { return n.IsNode() }},
		{func(n *Page) bool { return !n.IsPage() }},
		{func(n *Page) bool { return n.Scratch() != nil }},
		{func(n *Page) bool { return n.Hugo() != nil }},
		{func(n *Page) bool { return n.Now().Unix() == time.Now().Unix() }},
	} {

		n := s.newHomePage()

		if !this.assertFunc(n) {
			t.Errorf("[%d] Node method error", i)
		}
	}
}

func TestKind(t *testing.T) {
	t.Parallel()
	// Add tests for these constants to make sure they don't change
	require.Equal(t, "page", KindPage)
	require.Equal(t, "home", KindHome)
	require.Equal(t, "section", KindSection)
	require.Equal(t, "taxonomy", KindTaxonomy)
	require.Equal(t, "taxonomyTerm", KindTaxonomyTerm)

}

func TestTranslationKey(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", filepath.FromSlash("sect/simple.no.md")), "---\ntitle: \"A1\"\ntranslationKey: \"k1\"\n---\nContent\n")
	writeSource(t, fs, filepath.Join("content", filepath.FromSlash("sect/simple.en.md")), "---\ntitle: \"A2\"\n---\nContent\n")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 2)

	home, _ := s.Info.Home()
	assert.NotNil(home)
	assert.Equal("home", home.TranslationKey())
	assert.Equal("page/k1", s.RegularPages[0].TranslationKey())
	p2 := s.RegularPages[1]

	// This is a single language setup
	assert.Equal("page/sect/simple.en", p2.TranslationKey())

}

func TestChompBOM(t *testing.T) {
	t.Parallel()
	const utf8BOM = "\xef\xbb\xbf"

	cfg, fs := newTestCfg()

	writeSource(t, fs, filepath.Join("content", "simple.md"), utf8BOM+simplePage)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

	require.Len(t, s.RegularPages, 1)

	p := s.RegularPages[0]

	checkPageTitle(t, p, "Simple")
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
	t.Parallel()
	for _, disablePathToLower := range []bool{false, true} {
		for _, uglyURLs := range []bool{false, true} {
			t.Run(fmt.Sprintf("disablePathToLower=%t,uglyURLs=%t", disablePathToLower, uglyURLs), func(t *testing.T) {

				cfg, fs := newTestCfg()
				th := testHelper{cfg, fs, t}

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
				require.Len(t, s.RegularPages, 4)

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

				p := s.RegularPages[0]
				if uglyURLs {
					require.Equal(t, "/post/test0.dot.html", p.RelPermalink())
				} else {
					require.Equal(t, "/post/test0.dot/", p.RelPermalink())
				}

			})
		}
	}
}

func BenchmarkParsePage(b *testing.B) {
	s := newTestSite(b)
	f, _ := os.Open("testdata/redis.cn.md")
	var buf bytes.Buffer
	buf.ReadFrom(f)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		page, _ := s.NewPage("bench")
		page.ReadFrom(bytes.NewReader(buf.Bytes()))
	}
}
