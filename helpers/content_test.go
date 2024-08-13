// Copyright 2024 The Hugo Authors. All rights reserved.
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

package helpers_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/helpers"
)

func TestTrimShortHTML(t *testing.T) {
	tests := []struct {
		markup string
		input  []byte
		output []byte
	}{
		{"markdown", []byte(""), []byte("")},
		{"markdown", []byte("Plain text"), []byte("Plain text")},
		{"markdown", []byte("<p>Simple paragraph</p>"), []byte("Simple paragraph")},
		{"markdown", []byte("\n  \n \t  <p> \t Whitespace\nHTML  \n\t </p>\n\t"), []byte("Whitespace\nHTML")},
		{"markdown", []byte("<p>Multiple</p><p>paragraphs</p>"), []byte("<p>Multiple</p><p>paragraphs</p>")},
		{"markdown", []byte("<p>Nested<p>paragraphs</p></p>"), []byte("<p>Nested<p>paragraphs</p></p>")},
		{"markdown", []byte("<p>Hello</p>\n<ul>\n<li>list1</li>\n<li>list2</li>\n</ul>"), []byte("<p>Hello</p>\n<ul>\n<li>list1</li>\n<li>list2</li>\n</ul>")},
		// Issue 11698
		{"markdown", []byte("<h2 id=`a`>b</h2>\n\n<p>c</p>"), []byte("<h2 id=`a`>b</h2>\n\n<p>c</p>")},
		// Issue 12369
		{"markdown", []byte("<div class=\"paragraph\">\n<p>foo</p>\n</div>"), []byte("<div class=\"paragraph\">\n<p>foo</p>\n</div>")},
		{"asciidoc", []byte("<div class=\"paragraph\">\n<p>foo</p>\n</div>"), []byte("foo")},
	}

	c := newTestContentSpec(nil)
	for i, test := range tests {
		output := c.TrimShortHTML(test.input, test.markup)
		if !bytes.Equal(test.output, output) {
			t.Errorf("Test %d failed. Expected %q got %q", i, test.output, output)
		}
	}
}

func BenchmarkTrimShortHTML(b *testing.B) {
	c := newTestContentSpec(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.TrimShortHTML([]byte("<p>Simple paragraph</p>"), "markdown")
	}
}

func TestBytesToHTML(t *testing.T) {
	c := qt.New(t)
	c.Assert(helpers.BytesToHTML([]byte("dobedobedo")), qt.Equals, template.HTML("dobedobedo"))
}

func TestExtractTOCNormalContent(t *testing.T) {
	content := []byte("<nav>\n<ul>\nTOC<li><a href=\"#")

	actualTocLessContent, actualToc := helpers.ExtractTOC(content)
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

	actualTocLessContent, actualToc := helpers.ExtractTOC(content)
	// Because the start of Toc is greater than 70+startpoint of <li> content and empty TOC will be returned
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

	actualTocLessContent, actualToc := helpers.ExtractTOC(content)
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
		actualWordCount := helpers.TotalWords(this.s)

		if actualWordCount != this.words {
			t.Errorf("[%d] Actual word count (%d) for test string (%s) did not match %d", i, actualWordCount, this.s, this.words)
		}
	}
}

func BenchmarkTotalWords(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wordCount := helpers.TotalWords(totalWordsBenchmarkString)
		if wordCount != 400 {
			b.Fatal("Wordcount error")
		}
	}
}
