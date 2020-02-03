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

package blackfriday

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/gohugoio/hugo/markup/converter"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/markup/blackfriday/blackfriday_config"
	"github.com/russross/blackfriday"
)

func TestGetMarkdownExtensionsMasksAreRemovedFromExtensions(t *testing.T) {
	b := blackfriday_config.Default
	b.Extensions = []string{"headerId"}
	b.ExtensionsMask = []string{"noIntraEmphasis"}

	actualFlags := getMarkdownExtensions(b)
	if actualFlags&blackfriday.EXTENSION_NO_INTRA_EMPHASIS == blackfriday.EXTENSION_NO_INTRA_EMPHASIS {
		t.Errorf("Masked out flag {%v} found amongst returned extensions.", blackfriday.EXTENSION_NO_INTRA_EMPHASIS)
	}
}

func TestGetMarkdownExtensionsByDefaultAllExtensionsAreEnabled(t *testing.T) {
	type data struct {
		testFlag int
	}

	b := blackfriday_config.Default

	b.Extensions = []string{""}
	b.ExtensionsMask = []string{""}
	allExtensions := []data{
		{blackfriday.EXTENSION_NO_INTRA_EMPHASIS},
		{blackfriday.EXTENSION_TABLES},
		{blackfriday.EXTENSION_FENCED_CODE},
		{blackfriday.EXTENSION_AUTOLINK},
		{blackfriday.EXTENSION_STRIKETHROUGH},
		// {blackfriday.EXTENSION_LAX_HTML_BLOCKS},
		{blackfriday.EXTENSION_SPACE_HEADERS},
		// {blackfriday.EXTENSION_HARD_LINE_BREAK},
		// {blackfriday.EXTENSION_TAB_SIZE_EIGHT},
		{blackfriday.EXTENSION_FOOTNOTES},
		// {blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK},
		{blackfriday.EXTENSION_HEADER_IDS},
		// {blackfriday.EXTENSION_TITLEBLOCK},
		{blackfriday.EXTENSION_AUTO_HEADER_IDS},
		{blackfriday.EXTENSION_BACKSLASH_LINE_BREAK},
		{blackfriday.EXTENSION_DEFINITION_LISTS},
	}

	actualFlags := getMarkdownExtensions(b)
	for _, e := range allExtensions {
		if actualFlags&e.testFlag != e.testFlag {
			t.Errorf("Flag %v was not found in the list of extensions.", e)
		}
	}
}

func TestGetMarkdownExtensionsAddingFlagsThroughRenderingContext(t *testing.T) {
	b := blackfriday_config.Default

	b.Extensions = []string{"definitionLists"}
	b.ExtensionsMask = []string{""}

	actualFlags := getMarkdownExtensions(b)
	if actualFlags&blackfriday.EXTENSION_DEFINITION_LISTS != blackfriday.EXTENSION_DEFINITION_LISTS {
		t.Errorf("Masked out flag {%v} found amongst returned extensions.", blackfriday.EXTENSION_DEFINITION_LISTS)
	}
}

func TestGetFlags(t *testing.T) {
	b := blackfriday_config.Default
	flags := getFlags(false, b)
	if flags&blackfriday.HTML_USE_XHTML != blackfriday.HTML_USE_XHTML {
		t.Errorf("Test flag: %d was not found amongs set flags:%d; Result: %d", blackfriday.HTML_USE_XHTML, flags, flags&blackfriday.HTML_USE_XHTML)
	}
}

func TestGetAllFlags(t *testing.T) {
	c := qt.New(t)

	b := blackfriday_config.Default

	type data struct {
		testFlag int
	}

	allFlags := []data{
		{blackfriday.HTML_USE_XHTML},
		{blackfriday.HTML_FOOTNOTE_RETURN_LINKS},
		{blackfriday.HTML_USE_SMARTYPANTS},
		{blackfriday.HTML_SMARTYPANTS_QUOTES_NBSP},
		{blackfriday.HTML_SMARTYPANTS_ANGLED_QUOTES},
		{blackfriday.HTML_SMARTYPANTS_FRACTIONS},
		{blackfriday.HTML_HREF_TARGET_BLANK},
		{blackfriday.HTML_NOFOLLOW_LINKS},
		{blackfriday.HTML_NOREFERRER_LINKS},
		{blackfriday.HTML_SMARTYPANTS_DASHES},
		{blackfriday.HTML_SMARTYPANTS_LATEX_DASHES},
	}

	b.AngledQuotes = true
	b.Fractions = true
	b.HrefTargetBlank = true
	b.NofollowLinks = true
	b.NoreferrerLinks = true
	b.LatexDashes = true
	b.PlainIDAnchors = true
	b.SmartDashes = true
	b.Smartypants = true
	b.SmartypantsQuotesNBSP = true

	actualFlags := getFlags(false, b)

	var expectedFlags int
	//OR-ing flags together...
	for _, d := range allFlags {
		expectedFlags |= d.testFlag
	}

	c.Assert(actualFlags, qt.Equals, expectedFlags)
}

func TestConvert(t *testing.T) {
	c := qt.New(t)
	p, err := Provider.New(converter.ProviderConfig{
		Cfg: viper.New(),
	})
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, "<p>testContent</p>\n")
}

func TestGetHTMLRendererAnchors(t *testing.T) {
	c := qt.New(t)
	p, err := Provider.New(converter.ProviderConfig{
		Cfg: viper.New(),
	})
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{
		DocumentID: "testid",
		ConfigOverrides: map[string]interface{}{
			"plainIDAnchors": false,
			"footnotes":      true,
		},
	})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{Src: []byte(`# Header

This is a footnote.[^1] And then some.


[^1]: Footnote text.

`)})

	c.Assert(err, qt.IsNil)
	s := string(b.Bytes())
	c.Assert(s, qt.Contains, "<h1 id=\"header:testid\">Header</h1>")
	c.Assert(s, qt.Contains, "This is a footnote.<sup class=\"footnote-ref\" id=\"fnref:testid:1\"><a href=\"#fn:testid:1\">1</a></sup>")
	c.Assert(s, qt.Contains, "<a class=\"footnote-return\" href=\"#fnref:testid:1\"><sup>[return]</sup></a>")
}

// Tests borrowed from https://github.com/russross/blackfriday/blob/a925a152c144ea7de0f451eaf2f7db9e52fa005a/block_test.go#L1817
func TestSanitizedAnchorName(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{
			text: "This is a header",
			want: "this-is-a-header",
		},
		{
			text: "This is also          a header",
			want: "this-is-also-a-header",
		},
		{
			text: "main.go",
			want: "main-go",
		},
		{
			text: "Article 123",
			want: "article-123",
		},
		{
			text: "<- Let's try this, shall we?",
			want: "let-s-try-this-shall-we",
		},
		{
			text: "        ",
			want: "",
		},
		{
			text: "Hello, 世界",
			want: "hello-世界",
		},
	}
	for _, test := range tests {
		if got := SanitizedAnchorName(test.text); got != test.want {
			t.Errorf("SanitizedAnchorName(%q):\ngot %q\nwant %q", test.text, got, test.want)
		}
	}
}
