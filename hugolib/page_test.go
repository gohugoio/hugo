package hugolib

import (
    "html/template"
    "path"
    "strings"
    "testing"
    "time"
)

var EMPTY_PAGE = ""

const (
    SIMPLE_PAGE                      = "---\ntitle: Simple\n---\nSimple Page\n"
    INVALID_FRONT_MATTER_MISSING     = "This is a test"
    RENDER_NO_FRONT_MATTER           = "<!doctype><html><head></head><body>This is a test</body></html>"
    INVALID_FRONT_MATTER_SHORT_DELIM = `
--
title: Short delim start
---
Short Delim
`

    INVALID_FRONT_MATTER_SHORT_DELIM_ENDING = `
---
title: Short delim ending
--
Short Delim
`

    INVALID_FRONT_MATTER_LEADING_WS = `

 ---
title: Leading WS
---
Leading
`

    SIMPLE_PAGE_JSON = `
{
"title": "spf13-vim 3.0 release and new website",
"description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim.",
"tags": [ ".vimrc", "plugins", "spf13-vim", "vim" ],
"date": "2012-04-06",
"categories": [
    "Development",
    "VIM"
],
"slug": "spf13-vim-3-0-release-and-new-website"
}

Content of the file goes Here
`
    SIMPLE_PAGE_JSON_LOOSE = `
{
"title": "spf13-vim 3.0 release and new website"
"description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim."
"tags": [ ".vimrc", "plugins", "spf13-vim", "vim" ]
"date": "2012-04-06"
"categories": [
    "Development"
    "VIM"
],
"slug": "spf13-vim-3-0-release-and-new-website"
}

Content of the file goes Here
`
    SIMPLE_PAGE_RFC3339_DATE  = "---\ntitle: RFC3339 Date\ndate: \"2013-05-17T16:59:30Z\"\n---\nrfc3339 content"
    SIMPLE_PAGE_JSON_MULTIPLE = `
{
	"title": "foobar",
	"customData": { "foo": "bar" },
	"date": "2012-08-06"
}
Some text
`

    SIMPLE_PAGE_JSON_COMPACT = `
{"title":"foobar","customData":{"foo":"bar"},"date":"2012-08-06"}
Text
`

    SIMPLE_PAGE_NOLAYOUT = `---
title: simple_no_layout
---
No Layout called out`

    SIMPLE_PAGE_LAYOUT_FOOBAR = `---
title: simple layout foobar
layout: foobar
---
Layout foobar`

    SIMPLE_PAGE_TYPE_FOOBAR = `---
type: foobar
---
type foobar`

    SIMPLE_PAGE_TYPE_LAYOUT = `---
type: barfoo
layout: buzfoo
---
type and layout set`

    SIMPLE_PAGE_WITH_SUMMARY_DELIMITER = `---
title: Simple
---
Summary Next Line

<!--more-->
Some more text
`
    SIMPLE_PAGE_WITH_SHORTCODE_IN_SUMMARY = `---
title: Simple
---
Summary Next Line. {{% img src="/not/real" %}}.
More text here.

Some more text
`

    SIMPLE_PAGE_WITH_SUMMARY_DELIMITER_SAME_LINE = `---
title: Simple
---
Summary Same Line<!--more-->

Some more text
`

    SIMPLE_PAGE_WITH_LONG_CONTENT = `---
title: Simple
---
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
`
)

var PAGE_WITH_VARIOUS_FRONTMATTER_TYPES = `+++
a_string = "bar"
an_integer = 1
a_float = 1.3
a_date = 1979-05-27T07:32:00Z
+++
Front Matter with various frontmatter types`

func checkError(t *testing.T, err error, expected string) {
    if err == nil {
        t.Fatalf("err is nil.  Expected: %s", expected)
    }
    if err.Error() != expected {
        t.Errorf("err.Error() returned: '%s'.  Expected: '%s'", err.Error(), expected)
    }
}

func TestDegenerateEmptyPageZeroLengthName(t *testing.T) {
    _, err := ReadFrom(strings.NewReader(EMPTY_PAGE), "")
    if err == nil {
        t.Fatalf("A zero length page name must return an error")
    }

    checkError(t, err, "Zero length page name")
}

func TestDegenerateEmptyPage(t *testing.T) {
    _, err := ReadFrom(strings.NewReader(EMPTY_PAGE), "test")
    if err != nil {
        t.Fatalf("Empty files should not trigger an error. Should be able to touch a file while watching without erroring out.")
    }

    //checkError(t, err, "EOF")
}

func checkPageTitle(t *testing.T, page *Page, title string) {
    if page.Title != title {
        t.Fatalf("Page title is: %s.  Expected %s", page.Title, title)
    }
}

func checkPageContent(t *testing.T, page *Page, content string) {
    if page.Content != template.HTML(content) {
        t.Fatalf("Page content mismatch\nexp: %q\ngot: %q", content, page.Content)
    }
}

func checkPageSummary(t *testing.T, page *Page, summary string) {
    if page.Summary != template.HTML(summary) {
        t.Fatalf("Page summary is: `%s`.  Expected `%s`", page.Summary, summary)
    }
}

func checkPageType(t *testing.T, page *Page, pageType string) {
    if page.Type() != pageType {
        t.Fatalf("Page type is: %s.  Expected: %s", page.Type(), pageType)
    }
}

func checkPageLayout(t *testing.T, page *Page, layout ...string) {
    if !listEqual(page.Layout(), layout) {
        t.Fatalf("Page layout is: %s.  Expected: %s", page.Layout(), layout)
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

func TestCreateNewPage(t *testing.T) {
    p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE), "simple.md")
    p.Convert()

    if err != nil {
        t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
    }
    checkPageTitle(t, p, "Simple")
    checkPageContent(t, p, "<p>Simple Page</p>\n")
    checkPageSummary(t, p, "Simple Page")
    checkPageType(t, p, "page")
    checkPageLayout(t, p, "page/single.html", "single.html")
    checkTruncation(t, p, false, "simple short page")
}

func TestPageWithDelimiter(t *testing.T) {
    p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_WITH_SUMMARY_DELIMITER), "simple.md")
    p.Convert()
    if err != nil {
        t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
    }
    checkPageTitle(t, p, "Simple")
    checkPageContent(t, p, "<p>Summary Next Line</p>\n\n<p>Some more text</p>\n")
    checkPageSummary(t, p, "<p>Summary Next Line</p>\n")
    checkPageType(t, p, "page")
    checkPageLayout(t, p, "page/single.html", "single.html")
    checkTruncation(t, p, true, "page with summary delimiter")
}

func TestPageWithShortCodeInSummary(t *testing.T) {
    p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_WITH_SHORTCODE_IN_SUMMARY), "simple.md")
    p.Convert()
    if err != nil {
        t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
    }
    checkPageTitle(t, p, "Simple")
    checkPageContent(t, p, "<p>Summary Next Line. {{% img src=&ldquo;/not/real&rdquo; %}}.\nMore text here.</p>\n\n<p>Some more text</p>\n")
    checkPageSummary(t, p, "Summary Next Line. . More text here. Some more text")
    checkPageType(t, p, "page")
    checkPageLayout(t, p, "page/single.html", "single.html")
}

func TestPageWithMoreTag(t *testing.T) {
    p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_WITH_SUMMARY_DELIMITER_SAME_LINE), "simple.md")
    p.Convert()
    if err != nil {
        t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
    }
    checkPageTitle(t, p, "Simple")
    checkPageContent(t, p, "<p>Summary Same Line</p>\n\n<p>Some more text</p>\n")
    checkPageSummary(t, p, "<p>Summary Same Line</p>\n")
    checkPageType(t, p, "page")
    checkPageLayout(t, p, "page/single.html", "single.html")
}

func TestPageWithDate(t *testing.T) {
    p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_RFC3339_DATE), "simple")
    p.Convert()
    if err != nil {
        t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
    }
    d, err := time.Parse(time.RFC3339, "2013-05-17T16:59:30Z")
    if err != nil {
        t.Fatalf("Unable to prase page.")
    }
    checkPageDate(t, p, d)
}

func TestWordCount(t *testing.T) {
    p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_WITH_LONG_CONTENT), "simple.md")
    p.Convert()
    if err != nil {
        t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
    }

    if p.WordCount != 483 {
        t.Fatalf("incorrect word count. expected %v, got %v", 483, p.WordCount)
    }

    if p.FuzzyWordCount != 500 {
        t.Fatalf("incorrect word count. expected %v, got %v", 500, p.WordCount)
    }

    if p.ReadingTime != 3 {
        t.Fatalf("incorrect min read. expected %v, got %v", 3, p.ReadingTime)
    }

    checkTruncation(t, p, true, "long page")
}

func TestCreatePage(t *testing.T) {
    var tests = []struct {
        r string
    }{
        {SIMPLE_PAGE_JSON},
        {SIMPLE_PAGE_JSON_LOOSE},
        {SIMPLE_PAGE_JSON_MULTIPLE},
        //{strings.NewReader(SIMPLE_PAGE_JSON_COMPACT)},
    }

    for _, test := range tests {
        if _, err := ReadFrom(strings.NewReader(test.r), "page"); err != nil {
            t.Errorf("Unable to parse page: %s", err)
        }
    }
}

func TestDegenerateInvalidFrontMatterShortDelim(t *testing.T) {
    var tests = []struct {
        r   string
        err string
    }{
        {INVALID_FRONT_MATTER_SHORT_DELIM_ENDING, "Unable to read frontmatter at filepos 45: EOF"},
    }
    for _, test := range tests {
        _, err := ReadFrom(strings.NewReader(test.r), "invalid/front/matter/short/delim")
        checkError(t, err, test.err)
    }
}

func TestShouldRenderContent(t *testing.T) {
    var tests = []struct {
        text   string
        render bool
    }{
        {INVALID_FRONT_MATTER_MISSING, true},
        // TODO how to deal with malformed frontmatter.  In this case it'll be rendered as markdown.
        {INVALID_FRONT_MATTER_SHORT_DELIM, true},
        {RENDER_NO_FRONT_MATTER, false},
    }

    for _, test := range tests {
        p := pageMust(ReadFrom(strings.NewReader(test.text), "render/front/matter"))
        if p.IsRenderable() != test.render {
            t.Errorf("expected p.IsRenderable() == %t, got %t", test.render, p.IsRenderable())
        }
    }
}

func TestDifferentFrontMatterVarTypes(t *testing.T) {
    page, _ := ReadFrom(strings.NewReader(PAGE_WITH_VARIOUS_FRONTMATTER_TYPES), "test/file1.md")

    dateval, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
    if page.GetParam("a_string") != "bar" {
        t.Errorf("frontmatter not handling strings correctly should be %s, got: %s", "bar", page.GetParam("a_string"))
    }
    if page.GetParam("an_integer") != 1 {
        t.Errorf("frontmatter not handling ints correctly should be %s, got: %s", "1", page.GetParam("an_integer"))
    }
    if page.GetParam("a_float") != 1.3 {
        t.Errorf("frontmatter not handling floats correctly should be %s, got: %s", 1.3, page.GetParam("a_float"))
    }
    if page.GetParam("a_date") != dateval {
        t.Errorf("frontmatter not handling dates correctly should be %s, got: %s", dateval, page.GetParam("a_date"))
    }
}

func TestDegenerateInvalidFrontMatterLeadingWhitespace(t *testing.T) {
    _, err := ReadFrom(strings.NewReader(INVALID_FRONT_MATTER_LEADING_WS), "invalid/front/matter/leading/ws")
    if err != nil {
        t.Fatalf("Unable to parse front matter given leading whitespace: %s", err)
    }
}

func TestSectionEvaluation(t *testing.T) {
    page, _ := ReadFrom(strings.NewReader(SIMPLE_PAGE), "blue/file1.md")
    if page.Section != "blue" {
        t.Errorf("Section should be %s, got: %s", "blue", page.Section)
    }
}

func L(s ...string) []string {
    return s
}

func TestLayoutOverride(t *testing.T) {
    var (
        path_content_two_dir = path.Join("content", "dub", "sub", "file1.md")
        path_content_one_dir = path.Join("content", "gub", "file1.md")
        path_content_no_dir  = path.Join("content", "file1")
        path_one_directory   = path.Join("fub", "file1.md")
        path_no_directory    = path.Join("file1.md")
    )
    tests := []struct {
        content        string
        path           string
        expectedLayout []string
    }{
        {SIMPLE_PAGE_NOLAYOUT, path_content_two_dir, L("dub/sub/single.html", "dub/single.html", "single.html")},
        {SIMPLE_PAGE_NOLAYOUT, path_content_one_dir, L("gub/single.html", "single.html")},
        {SIMPLE_PAGE_NOLAYOUT, path_content_no_dir, L("page/single.html", "single.html")},
        {SIMPLE_PAGE_NOLAYOUT, path_one_directory, L("fub/single.html", "single.html")},
        {SIMPLE_PAGE_NOLAYOUT, path_no_directory, L("page/single.html", "single.html")},
        {SIMPLE_PAGE_LAYOUT_FOOBAR, path_content_two_dir, L("dub/sub/foobar.html", "dub/foobar.html", "foobar.html")},
        {SIMPLE_PAGE_LAYOUT_FOOBAR, path_content_one_dir, L("gub/foobar.html", "foobar.html")},
        {SIMPLE_PAGE_LAYOUT_FOOBAR, path_one_directory, L("fub/foobar.html", "foobar.html")},
        {SIMPLE_PAGE_LAYOUT_FOOBAR, path_no_directory, L("page/foobar.html", "foobar.html")},
        {SIMPLE_PAGE_TYPE_FOOBAR, path_content_two_dir, L("foobar/single.html", "single.html")},
        {SIMPLE_PAGE_TYPE_FOOBAR, path_content_one_dir, L("foobar/single.html", "single.html")},
        {SIMPLE_PAGE_TYPE_FOOBAR, path_content_no_dir, L("foobar/single.html", "single.html")},
        {SIMPLE_PAGE_TYPE_FOOBAR, path_one_directory, L("foobar/single.html", "single.html")},
        {SIMPLE_PAGE_TYPE_FOOBAR, path_no_directory, L("foobar/single.html", "single.html")},
        {SIMPLE_PAGE_TYPE_LAYOUT, path_content_two_dir, L("barfoo/buzfoo.html", "buzfoo.html")},
        {SIMPLE_PAGE_TYPE_LAYOUT, path_content_one_dir, L("barfoo/buzfoo.html", "buzfoo.html")},
        {SIMPLE_PAGE_TYPE_LAYOUT, path_content_no_dir, L("barfoo/buzfoo.html", "buzfoo.html")},
        {SIMPLE_PAGE_TYPE_LAYOUT, path_one_directory, L("barfoo/buzfoo.html", "buzfoo.html")},
        {SIMPLE_PAGE_TYPE_LAYOUT, path_no_directory, L("barfoo/buzfoo.html", "buzfoo.html")},
    }
    for _, test := range tests {
        p, err := ReadFrom(strings.NewReader(test.content), test.path)
        if err != nil {
            t.Fatalf("Unable to parse content:\n%s\n", test.content)
        }
        if !listEqual(p.Layout(), test.expectedLayout) {
            t.Errorf("Layout mismatch. Expected: %s, got: %s", test.expectedLayout, p.Layout())
        }
    }
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
