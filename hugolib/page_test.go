package hugolib

import (
	"html/template"
	"io"
	"strings"
	"testing"
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

	checkError(t, err, "unable to locate front matter")
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

func TestCreateNewPage(t *testing.T) {
	p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE), "simple")
	if err != nil {
		t.Fatalf("Unable to create a page with frontmatter and body content: %s", err)
	}
	checkPageTitle(t, p, "Simple")
	checkPageContent(t, p, "<p>Simple Page</p>\n")
	checkPageType(t, p, "page")
	checkPageLayout(t, p, "page/single.html")
}

func TestCreatePage(t *testing.T) {
	var tests = []struct {
		r io.Reader
	}{
		{strings.NewReader(SIMPLE_PAGE_JSON)},
		{strings.NewReader(SIMPLE_PAGE_JSON_MULTIPLE)},
		//{strings.NewReader(SIMPLE_PAGE_JSON_COMPACT)},
	}

	for _, test := range tests {
		_, err := ReadFrom(test.r, "page")
		if err != nil {
			t.Errorf("Unable to parse page: %s", err)
		}
	}
}

func TestDegenerateInvalidFrontMatterShortDelim(t *testing.T) {
	var tests = []struct {
		r   io.Reader
		err string
	}{
		{strings.NewReader(INVALID_FRONT_MATTER_SHORT_DELIM), "unable to match beginning front matter delimiter"},
		{strings.NewReader(INVALID_FRONT_MATTER_SHORT_DELIM_ENDING), "unable to match ending front matter delimiter"},
		{strings.NewReader(INVALID_FRONT_MATTER_MISSING), "unable to detect front matter"},
	}
	for _, test := range tests {
		_, err := ReadFrom(test.r, "invalid/front/matter/short/delim")
		checkError(t, err, test.err)
	}
}

func TestDegenerateInvalidFrontMatterLeadingWhitespace(t *testing.T) {
	_, err := ReadFrom(strings.NewReader(INVALID_FRONT_MATTER_LEADING_WS), "invalid/front/matter/leading/ws")
	if err != nil {
		t.Fatalf("Unable to parse front matter given leading whitespace: %s", err)
	}
}
