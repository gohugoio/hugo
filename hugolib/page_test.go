package hugolib

import (
	"html/template"
	"path"
	"strings"
	"testing"
	"time"
)

var EMPTY_PAGE = ""

var SIMPLE_PAGE = `---
title: Simple
---
Simple Page
`

var INVALID_FRONT_MATTER_MISSING = `This is a test`

var INVALID_FRONT_MATTER_SHORT_DELIM = `
--
title: Short delim start
---
Short Delim
`

var INVALID_FRONT_MATTER_SHORT_DELIM_ENDING = `
---
title: Short delim ending
--
Short Delim
`

var INVALID_FRONT_MATTER_LEADING_WS = `

 ---
title: Leading WS
---
Leading
`

var SIMPLE_PAGE_JSON = `
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
var SIMPLE_PAGE_RFC3339_DATE = "---\ntitle: RFC3339 Date\ndate: \"2013-05-17T16:59:30Z\"\n---\nrfc3339 content"
var SIMPLE_PAGE_JSON_MULTIPLE = `
{
	"title": "foobar",
	"customData": { "foo": "bar" },
	"date": "2012-08-06"
}
Some text
`

var SIMPLE_PAGE_JSON_COMPACT = `
{"title":"foobar","customData":{"foo":"bar"},"date":"2012-08-06"}
Text
`

var SIMPLE_PAGE_NOLAYOUT = `---
title: simple_no_layout
---
No Layout called out`

var SIMPLE_PAGE_LAYOUT_FOOBAR = `---
title: simple layout foobar
layout: foobar
---
Layout foobar`

var SIMPLE_PAGE_TYPE_FOOBAR = `---
type: foobar
---
type foobar`

var SIMPLE_PAGE_TYPE_LAYOUT = `---
type: barfoo
layout: buzfoo
---
type and layout set`

var SIMPLE_PAGE_WITH_SUMMARY_DELIMITER = `---
title: Simple
---
Simple Page

<!--more-->
Some more text
`

var SIMPLE_PAGE_WITH_SUMMARY_DELIMITER_SAME_LINE = `---
title: Simple
---
Simple Page<!--more-->

Some more text
`

func checkError(t *testing.T, err error, expected string) {
	if err == nil {
		t.Fatalf("err is nil")
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
	if err == nil {
		t.Fatalf("Expected ReadFrom to return an error when an empty buffer is passed.")
	}

	checkError(t, err, "EOF")
}

func checkPageTitle(t *testing.T, page *Page, title string) {
	if page.Title != title {
		t.Fatalf("Page title is: %s.  Expected %s", page.Title, title)
	}
}

func checkPageContent(t *testing.T, page *Page, content string) {
	if page.Content != template.HTML(content) {
		t.Fatalf("Page content is: %s.  Expected %s", page.Content, content)
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

func checkPageLayout(t *testing.T, page *Page, layout string) {
	if page.Layout() != layout {
		t.Fatalf("Page layout is: %s.  Expected: %s", page.Layout(), layout)
	}
}

func checkPageDate(t *testing.T, page *Page, time time.Time) {
	if page.Date != time {
		t.Fatalf("Page date is: %s.  Expected: %s", page.Date, time)
	}
}

func TestCreateNewPage(t *testing.T) {
	p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE), "simple")
	if err != nil {
		t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
	}
	checkPageTitle(t, p, "Simple")
	checkPageContent(t, p, "<p>Simple Page</p>\n")
	checkPageSummary(t, p, "Simple Page")
	checkPageType(t, p, "page")
	checkPageLayout(t, p, "page/single.html")
}

func TestPageWithDelimiter(t *testing.T) {
	p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_WITH_SUMMARY_DELIMITER), "simple")
	if err != nil {
		t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
	}
	checkPageTitle(t, p, "Simple")
	checkPageContent(t, p, "<p>Simple Page</p>\n\n<p>Some more text</p>\n")
	checkPageSummary(t, p, "<p>Simple Page</p>\n")
	checkPageType(t, p, "page")
	checkPageLayout(t, p, "page/single.html")

}

func TestPageWithMoreTag(t *testing.T) {
	p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_WITH_SUMMARY_DELIMITER_SAME_LINE), "simple")
	if err != nil {
		t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
	}
	checkPageTitle(t, p, "Simple")
	checkPageContent(t, p, "<p>Simple Page</p>\n\n<p>Some more text</p>\n")
	checkPageSummary(t, p, "<p>Simple Page</p>\n")
	checkPageType(t, p, "page")
	checkPageLayout(t, p, "page/single.html")
}

func TestPageWithDate(t *testing.T) {
	p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_RFC3339_DATE), "simple")
	if err != nil {
		t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
	}
	d, err := time.Parse(time.RFC3339, "2013-05-17T16:59:30Z")
	if err != nil {
		t.Fatalf("Unable to prase page.")
	}
	checkPageDate(t, p, d)
}

func TestCreatePage(t *testing.T) {
	var tests = []struct {
		r string
	}{
		{SIMPLE_PAGE_JSON},
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
		{INVALID_FRONT_MATTER_SHORT_DELIM, "Unable to locate frontmatter"},
		{INVALID_FRONT_MATTER_SHORT_DELIM_ENDING, "EOF"},
		{INVALID_FRONT_MATTER_MISSING, "Unable to locate frontmatter"},
	}
	for _, test := range tests {
		_, err := ReadFrom(strings.NewReader(test.r), "invalid/front/matter/short/delim")
		checkError(t, err, test.err)
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

func TestLayoutOverride(t *testing.T) {
	var (
		path_content_one_dir = path.Join("content", "gub", "file1.md")
		path_content_two_dir = path.Join("content", "dub", "sub", "file1.md")
		path_content_no_dir  = path.Join("content", "file1")
		path_one_directory   = path.Join("fub", "file1.md")
		path_no_directory    = path.Join("file1.md")
	)
	tests := []struct {
		content        string
		path           string
		expectedLayout string
	}{
		{SIMPLE_PAGE_NOLAYOUT, path_content_two_dir, "sub/single.html"},
		{SIMPLE_PAGE_NOLAYOUT, path_content_one_dir, "gub/single.html"},
		{SIMPLE_PAGE_NOLAYOUT, path_content_no_dir, "page/single.html"},
		{SIMPLE_PAGE_NOLAYOUT, path_one_directory, "fub/single.html"},
		{SIMPLE_PAGE_NOLAYOUT, path_no_directory, "page/single.html"},
		{SIMPLE_PAGE_LAYOUT_FOOBAR, path_content_two_dir, "foobar"},
		{SIMPLE_PAGE_LAYOUT_FOOBAR, path_content_one_dir, "foobar"},
		{SIMPLE_PAGE_LAYOUT_FOOBAR, path_one_directory, "foobar"},
		{SIMPLE_PAGE_LAYOUT_FOOBAR, path_no_directory, "foobar"},
		{SIMPLE_PAGE_TYPE_FOOBAR, path_content_two_dir, "foobar/single.html"},
		{SIMPLE_PAGE_TYPE_FOOBAR, path_content_one_dir, "foobar/single.html"},
		{SIMPLE_PAGE_TYPE_FOOBAR, path_content_no_dir, "foobar/single.html"},
		{SIMPLE_PAGE_TYPE_FOOBAR, path_one_directory, "foobar/single.html"},
		{SIMPLE_PAGE_TYPE_FOOBAR, path_no_directory, "foobar/single.html"},
		{SIMPLE_PAGE_TYPE_LAYOUT, path_content_two_dir, "buzfoo"},
		{SIMPLE_PAGE_TYPE_LAYOUT, path_content_one_dir, "buzfoo"},
		{SIMPLE_PAGE_TYPE_LAYOUT, path_content_no_dir, "buzfoo"},
		{SIMPLE_PAGE_TYPE_LAYOUT, path_one_directory, "buzfoo"},
		{SIMPLE_PAGE_TYPE_LAYOUT, path_no_directory, "buzfoo"},
	}
	for _, test := range tests {
		p, err := ReadFrom(strings.NewReader(test.content), test.path)
		if err != nil {
			t.Fatalf("Unable to parse content:\n%s\n", test.content)
		}
		if p.Layout() != test.expectedLayout {
			t.Errorf("Layout mismatch. Expected: %s, got: %s", test.expectedLayout, p.Layout())
		}
	}
}
