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

package helpers

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/miekg/mmark"
	"github.com/russross/blackfriday"
	"github.com/stretchr/testify/assert"
)

const tstHTMLContent = "<!DOCTYPE html><html><head><script src=\"http://two/foobar.js\"></script></head><body><nav><ul><li hugo-nav=\"section_0\"></li><li hugo-nav=\"section_1\"></li></ul></nav><article>content <a href=\"http://two/foobar\">foobar</a>. Follow up</article><p>This is some text.<br>And some more.</p></body></html>"

func TestStripHTML(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"<h1>strip h1 tag <h1>", "strip h1 tag "},
		{"<p> strip p tag </p>", " strip p tag "},
		{"</br> strip br<br>", " strip br\n"},
		{"</br> strip br2<br />", " strip br2\n"},
		{"This <strong>is</strong> a\nnewline", "This is a newline"},
		{"No Tags", "No Tags"},
		{`<p>Summary Next Line.
<figure >

        <img src="/not/real" />


</figure>
.
More text here.</p>

<p>Some more text</p>`, "Summary Next Line.  . More text here.\nSome more text\n"},
	}
	for i, d := range data {
		output := StripHTML(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func BenchmarkStripHTML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StripHTML(tstHTMLContent)
	}
}

func TestStripEmptyNav(t *testing.T) {
	cleaned := stripEmptyNav([]byte("do<nav>\n</nav>\n\nbedobedo"))
	assert.Equal(t, []byte("dobedobedo"), cleaned)
}

func TestBytesToHTML(t *testing.T) {
	assert.Equal(t, template.HTML("dobedobedo"), BytesToHTML([]byte("dobedobedo")))
}

var benchmarkTruncateString = strings.Repeat("This is a sentence about nothing.", 20)

func BenchmarkTestTruncateWordsToWholeSentence(b *testing.B) {
	c := newTestContentSpec()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.TruncateWordsToWholeSentence(benchmarkTruncateString)
	}
}

func BenchmarkTestTruncateWordsToWholeSentenceOld(b *testing.B) {
	c := newTestContentSpec()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.truncateWordsToWholeSentenceOld(benchmarkTruncateString)
	}
}

func TestTruncateWordsToWholeSentence(t *testing.T) {
	c := newTestContentSpec()
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"This is a sentence.", "This is a sentence.", 5, false},
		{"This is also a sentence!", "This is also a sentence!", 1, false},
		{"To be. Or not to be. That's the question.", "To be.", 1, true},
		{" \nThis is not a sentence\nAnd this is another", "This is not a sentence", 4, true},
		{"", "", 10, false},
		{"This... is a more difficult test?", "This... is a more difficult test?", 1, false},
	}
	for i, d := range data {
		c.summaryLength = d.max
		output, truncated := c.TruncateWordsToWholeSentence(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}

func TestTruncateWordsByRune(t *testing.T) {
	c := newTestContentSpec()
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"", "", 1, false},
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"Hello 中国", "", 0, true},
		{"这是中文，全中文。", "这是中文，", 5, true},
		{"Hello 中国", "Hello 中", 2, true},
		{"Hello 中国", "Hello 中国", 3, false},
		{"Hello中国 Good 好的", "Hello中国 Good 好", 9, true},
		{"This is a sentence.", "This is", 2, true},
		{"This is also a sentence!", "This", 1, true},
		{"To be. Or not to be. That's the question.", "To be. Or not", 4, true},
		{" \nThis is    not a sentence\n ", "This is not", 3, true},
	}
	for i, d := range data {
		c.summaryLength = d.max
		output, truncated := c.TruncateWordsByRune(strings.Fields(d.input))
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}

func TestGetHTMLRendererFlags(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	renderer := c.getHTMLRenderer(blackfriday.HTML_USE_XHTML, ctx)
	flags := renderer.GetFlags()
	if flags&blackfriday.HTML_USE_XHTML != blackfriday.HTML_USE_XHTML {
		t.Errorf("Test flag: %d was not found amongs set flags:%d; Result: %d", blackfriday.HTML_USE_XHTML, flags, flags&blackfriday.HTML_USE_XHTML)
	}
}

func TestGetHTMLRendererAllFlags(t *testing.T) {
	c := newTestContentSpec()

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
		{blackfriday.HTML_SMARTYPANTS_DASHES},
		{blackfriday.HTML_SMARTYPANTS_LATEX_DASHES},
	}
	defaultFlags := blackfriday.HTML_USE_XHTML
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Config.AngledQuotes = true
	ctx.Config.Fractions = true
	ctx.Config.HrefTargetBlank = true
	ctx.Config.LatexDashes = true
	ctx.Config.PlainIDAnchors = true
	ctx.Config.SmartDashes = true
	ctx.Config.Smartypants = true
	ctx.Config.SmartypantsQuotesNBSP = true
	renderer := c.getHTMLRenderer(defaultFlags, ctx)
	actualFlags := renderer.GetFlags()
	var expectedFlags int
	//OR-ing flags together...
	for _, d := range allFlags {
		expectedFlags |= d.testFlag
	}
	if expectedFlags != actualFlags {
		t.Errorf("Expected flags (%d) did not equal actual (%d) flags.", expectedFlags, actualFlags)
	}
}

func TestGetHTMLRendererAnchors(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.DocumentID = "testid"
	ctx.Config.PlainIDAnchors = false

	actualRenderer := c.getHTMLRenderer(0, ctx)
	headerBuffer := &bytes.Buffer{}
	footnoteBuffer := &bytes.Buffer{}
	expectedFootnoteHref := []byte("href=\"#fn:testid:href\"")
	expectedHeaderID := []byte("<h1 id=\"id:testid\"></h1>\n")

	actualRenderer.Header(headerBuffer, func() bool { return true }, 1, "id")
	actualRenderer.FootnoteRef(footnoteBuffer, []byte("href"), 1)

	if !bytes.Contains(footnoteBuffer.Bytes(), expectedFootnoteHref) {
		t.Errorf("Footnote anchor prefix not applied. Actual:%s Expected:%s", footnoteBuffer.String(), expectedFootnoteHref)
	}

	if !bytes.Equal(headerBuffer.Bytes(), expectedHeaderID) {
		t.Errorf("Header Id Postfix not applied. Actual:%s Expected:%s", headerBuffer.String(), expectedHeaderID)
	}
}

func TestGetMmarkHTMLRenderer(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.DocumentID = "testid"
	ctx.Config.PlainIDAnchors = false
	actualRenderer := c.getMmarkHTMLRenderer(0, ctx)

	headerBuffer := &bytes.Buffer{}
	footnoteBuffer := &bytes.Buffer{}
	expectedFootnoteHref := []byte("href=\"#fn:testid:href\"")
	expectedHeaderID := []byte("<h1 id=\"id\"></h1>")

	actualRenderer.FootnoteRef(footnoteBuffer, []byte("href"), 1)
	actualRenderer.Header(headerBuffer, func() bool { return true }, 1, "id")

	if !bytes.Contains(footnoteBuffer.Bytes(), expectedFootnoteHref) {
		t.Errorf("Footnote anchor prefix not applied. Actual:%s Expected:%s", footnoteBuffer.String(), expectedFootnoteHref)
	}

	if bytes.Equal(headerBuffer.Bytes(), expectedHeaderID) {
		t.Errorf("Header Id Postfix applied. Actual:%s Expected:%s", headerBuffer.String(), expectedHeaderID)
	}
}

func TestGetMarkdownExtensionsMasksAreRemovedFromExtensions(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Config.Extensions = []string{"headerId"}
	ctx.Config.ExtensionsMask = []string{"noIntraEmphasis"}

	actualFlags := getMarkdownExtensions(ctx)
	if actualFlags&blackfriday.EXTENSION_NO_INTRA_EMPHASIS == blackfriday.EXTENSION_NO_INTRA_EMPHASIS {
		t.Errorf("Masked out flag {%v} found amongst returned extensions.", blackfriday.EXTENSION_NO_INTRA_EMPHASIS)
	}
}

func TestGetMarkdownExtensionsByDefaultAllExtensionsAreEnabled(t *testing.T) {
	type data struct {
		testFlag int
	}
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Config.Extensions = []string{""}
	ctx.Config.ExtensionsMask = []string{""}
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

	actualFlags := getMarkdownExtensions(ctx)
	for _, e := range allExtensions {
		if actualFlags&e.testFlag != e.testFlag {
			t.Errorf("Flag %v was not found in the list of extensions.", e)
		}
	}
}

func TestGetMarkdownExtensionsAddingFlagsThroughRenderingContext(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Config.Extensions = []string{"definitionLists"}
	ctx.Config.ExtensionsMask = []string{""}

	actualFlags := getMarkdownExtensions(ctx)
	if actualFlags&blackfriday.EXTENSION_DEFINITION_LISTS != blackfriday.EXTENSION_DEFINITION_LISTS {
		t.Errorf("Masked out flag {%v} found amongst returned extensions.", blackfriday.EXTENSION_DEFINITION_LISTS)
	}
}

func TestGetMarkdownRenderer(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Content = []byte("testContent")
	actualRenderedMarkdown := c.markdownRender(ctx)
	expectedRenderedMarkdown := []byte("<p>testContent</p>\n")
	if !bytes.Equal(actualRenderedMarkdown, expectedRenderedMarkdown) {
		t.Errorf("Actual rendered Markdown (%s) did not match expected markdown (%s)", actualRenderedMarkdown, expectedRenderedMarkdown)
	}
}

func TestGetMarkdownRendererWithTOC(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{RenderTOC: true, Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Content = []byte("testContent")
	actualRenderedMarkdown := c.markdownRender(ctx)
	expectedRenderedMarkdown := []byte("<nav>\n</nav>\n\n<p>testContent</p>\n")
	if !bytes.Equal(actualRenderedMarkdown, expectedRenderedMarkdown) {
		t.Errorf("Actual rendered Markdown (%s) did not match expected markdown (%s)", actualRenderedMarkdown, expectedRenderedMarkdown)
	}
}

func TestGetMmarkExtensions(t *testing.T) {
	//TODO: This is doing the same just with different marks...
	type data struct {
		testFlag int
	}
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Config.Extensions = []string{"tables"}
	ctx.Config.ExtensionsMask = []string{""}
	allExtensions := []data{
		{mmark.EXTENSION_TABLES},
		{mmark.EXTENSION_FENCED_CODE},
		{mmark.EXTENSION_AUTOLINK},
		{mmark.EXTENSION_SPACE_HEADERS},
		{mmark.EXTENSION_CITATION},
		{mmark.EXTENSION_TITLEBLOCK_TOML},
		{mmark.EXTENSION_HEADER_IDS},
		{mmark.EXTENSION_AUTO_HEADER_IDS},
		{mmark.EXTENSION_UNIQUE_HEADER_IDS},
		{mmark.EXTENSION_FOOTNOTES},
		{mmark.EXTENSION_SHORT_REF},
		{mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK},
		{mmark.EXTENSION_INCLUDE},
	}

	actualFlags := getMmarkExtensions(ctx)
	for _, e := range allExtensions {
		if actualFlags&e.testFlag != e.testFlag {
			t.Errorf("Flag %v was not found in the list of extensions.", e)
		}
	}
}

func TestMmarkRender(t *testing.T) {
	c := newTestContentSpec()
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	ctx.Content = []byte("testContent")
	actualRenderedMarkdown := c.mmarkRender(ctx)
	expectedRenderedMarkdown := []byte("<p>testContent</p>\n")
	if !bytes.Equal(actualRenderedMarkdown, expectedRenderedMarkdown) {
		t.Errorf("Actual rendered Markdown (%s) did not match expected markdown (%s)", actualRenderedMarkdown, expectedRenderedMarkdown)
	}
}

func TestExtractTOCNormalContent(t *testing.T) {
	content := []byte("<nav>\n<ul>\nTOC<li><a href=\"#")

	actualTocLessContent, actualToc := ExtractTOC(content)
	expectedTocLess := []byte("TOC<li><a href=\"#")
	expectedToc := []byte("<nav id=\"TableOfContents\">\n<ul>\n")

	if !bytes.Equal(actualTocLessContent, expectedTocLess) {
		t.Errorf("Actual tocless (%s) did not equal expected (%s) tocless content", actualTocLessContent, expectedTocLess)
	}

	if !bytes.Equal(actualToc, expectedToc) {
		t.Errorf("Actual toc (%s) did not equal expected (%s) toc content", actualToc, expectedToc)
	}
}

func TestExtractTOCGreaterThanSeventy(t *testing.T) {
	content := []byte("<nav>\n<ul>\nTOC This is a very long content which will definitely be greater than seventy, I promise you that.<li><a href=\"#")

	actualTocLessContent, actualToc := ExtractTOC(content)
	//Because the start of Toc is greater than 70+startpoint of <li> content and empty TOC will be returned
	expectedToc := []byte("")

	if !bytes.Equal(actualTocLessContent, content) {
		t.Errorf("Actual tocless (%s) did not equal expected (%s) tocless content", actualTocLessContent, content)
	}

	if !bytes.Equal(actualToc, expectedToc) {
		t.Errorf("Actual toc (%s) did not equal expected (%s) toc content", actualToc, expectedToc)
	}
}

func TestExtractNoTOC(t *testing.T) {
	content := []byte("TOC")

	actualTocLessContent, actualToc := ExtractTOC(content)
	expectedToc := []byte("")

	if !bytes.Equal(actualTocLessContent, content) {
		t.Errorf("Actual tocless (%s) did not equal expected (%s) tocless content", actualTocLessContent, content)
	}

	if !bytes.Equal(actualToc, expectedToc) {
		t.Errorf("Actual toc (%s) did not equal expected (%s) toc content", actualToc, expectedToc)
	}
}

var totalWordsBenchmarkString = strings.Repeat("Hugo Rocks ", 200)

func TestTotalWords(t *testing.T) {

	for i, this := range []struct {
		s     string
		words int
	}{
		{"Two, Words!", 2},
		{"Word", 1},
		{"", 0},
		{"One, Two,      Three", 3},
		{totalWordsBenchmarkString, 400},
	} {
		actualWordCount := TotalWords(this.s)

		if actualWordCount != this.words {
			t.Errorf("[%d] Actual word count (%d) for test string (%s) did not match %d", i, actualWordCount, this.s, this.words)
		}
	}
}

func BenchmarkTotalWords(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wordCount := TotalWords(totalWordsBenchmarkString)
		if wordCount != 400 {
			b.Fatal("Wordcount error")
		}
	}
}

func BenchmarkTotalWordsOld(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wordCount := totalWordsOld(totalWordsBenchmarkString)
		if wordCount != 400 {
			b.Fatal("Wordcount error")
		}
	}
}
