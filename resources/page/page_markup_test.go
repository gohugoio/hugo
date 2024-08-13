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

package page

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/media"
)

func TestExtractSummaryFromHTML(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		mt                          media.Type
		input                       string
		isCJK                       bool
		numWords                    int
		expectSummary               string
		expectContentWithoutSummary string
	}{
		{media.Builtin.ReStructuredTextType, "<div class=\"document\">\n\n\n<p>Simple Page</p>\n</div>", false, 70, "<div class=\"document\">\n\n\n<p>Simple Page</p>\n</div>", ""},
		{media.Builtin.ReStructuredTextType, "<div class=\"document\"><p>First paragraph</p><p>Second paragraph</p></div>", false, 2, `<div class="document"><p>First paragraph</p></div>`, "<div class=\"document\"><p>Second paragraph</p></div>"},
		{media.Builtin.MarkdownType, "<p>First paragraph</p>", false, 10, "<p>First paragraph</p>", ""},
		{media.Builtin.MarkdownType, "<p>First paragraph</p><p>Second paragraph</p>", false, 2, "<p>First paragraph</p>", "<p>Second paragraph</p>"},
		{media.Builtin.MarkdownType, "<p>First paragraph</p><p>Second paragraph</p><p>Third paragraph</p>", false, 3, "<p>First paragraph</p><p>Second paragraph</p>", "<p>Third paragraph</p>"},
		{media.Builtin.AsciiDocType, "<div><p>First paragraph</p></div><div><p>Second paragraph</p></div>", false, 2, "<div><p>First paragraph</p></div>", "<div><p>Second paragraph</p></div>"},
		{media.Builtin.MarkdownType, "<p>这是中文，全中文</p><p>a这是中文，全中文</p>", true, 5, "<p>这是中文，全中文</p>", "<p>a这是中文，全中文</p>"},
	}

	for i, test := range tests {
		summary := ExtractSummaryFromHTML(test.mt, test.input, test.numWords, test.isCJK)
		c.Assert(summary.Summary(), qt.Equals, test.expectSummary, qt.Commentf("Summary %d", i))
		c.Assert(summary.ContentWithoutSummary(), qt.Equals, test.expectContentWithoutSummary, qt.Commentf("ContentWithoutSummary %d", i))
	}
}

func TestExtractSummaryFromHTMLWithDivider(t *testing.T) {
	c := qt.New(t)

	const divider = "FOOO"

	tests := []struct {
		mt                          media.Type
		input                       string
		expectSummary               string
		expectContentWithoutSummary string
		expectContent               string
	}{
		{media.Builtin.MarkdownType, "<p>First paragraph</p><p>FOOO</p><p>Second paragraph</p>", "<p>First paragraph</p>", "<p>Second paragraph</p>", "<p>First paragraph</p><p>Second paragraph</p>"},
		{media.Builtin.MarkdownType, "<p>First paragraph</p>\n<p>FOOO</p>\n<p>Second paragraph</p>", "<p>First paragraph</p>", "<p>Second paragraph</p>", "<p>First paragraph</p>\n<p>Second paragraph</p>"},
		{media.Builtin.MarkdownType, "<p>FOOO</p>\n<p>First paragraph</p>", "", "<p>First paragraph</p>", "<p>First paragraph</p>"},
		{media.Builtin.MarkdownType, "<p>First paragraph</p><p>Second paragraphFOOO</p><p>Third paragraph</p>", "<p>First paragraph</p><p>Second paragraph</p>", "<p>Third paragraph</p>", "<p>First paragraph</p><p>Second paragraph</p><p>Third paragraph</p>"},
		{media.Builtin.MarkdownType, "<p>这是中文，全中文FOOO</p><p>a这是中文，全中文</p>", "<p>这是中文，全中文</p>", "<p>a这是中文，全中文</p>", "<p>这是中文，全中文</p><p>a这是中文，全中文</p>"},
		{media.Builtin.MarkdownType, `<p>a <strong>b</strong>` + "\v" + ` c</p>` + "\n<p>FOOO</p>", "<p>a <strong>b</strong>\v c</p>", "", "<p>a <strong>b</strong>\v c</p>"},

		{media.Builtin.HTMLType, "<p>First paragraph</p>FOOO<p>Second paragraph</p>", "<p>First paragraph</p>", "<p>Second paragraph</p>", "<p>First paragraph</p><p>Second paragraph</p>"},

		{media.Builtin.ReStructuredTextType, "<div class=\"document\">\n\n\n<p>This is summary.</p>\n<p>FOOO</p>\n<p>This is content.</p>\n</div>", "<div class=\"document\">\n\n\n<p>This is summary.</p>\n</div>", "<div class=\"document\"><p>This is content.</p>\n</div>", "<div class=\"document\">\n\n\n<p>This is summary.</p>\n<p>This is content.</p>\n</div>"},
		{media.Builtin.ReStructuredTextType, "<div class=\"document\"><p>First paragraphFOOO</p><p>Second paragraph</p></div>", "<div class=\"document\"><p>First paragraph</p></div>", "<div class=\"document\"><p>Second paragraph</p></div>", `<div class="document"><p>First paragraph</p><p>Second paragraph</p></div>`},

		{media.Builtin.AsciiDocType, "<div class=\"paragraph\"><p>Summary Next Line</p></div><div class=\"paragraph\"><p>FOOO</p></div><div class=\"paragraph\"><p>Some more text</p></div>", "<div class=\"paragraph\"><p>Summary Next Line</p></div>", "<div class=\"paragraph\"><p>Some more text</p></div>", "<div class=\"paragraph\"><p>Summary Next Line</p></div><div class=\"paragraph\"><p>Some more text</p></div>"},
		{media.Builtin.AsciiDocType, "<div class=\"paragraph\">\n<p>Summary Next Line</p>\n</div>\n<div class=\"paragraph\">\n<p>FOOO</p>\n</div>\n<div class=\"paragraph\">\n<p>Some more text</p>\n</div>\n", "<div class=\"paragraph\">\n<p>Summary Next Line</p>\n</div>", "<div class=\"paragraph\">\n<p>Some more text</p>\n</div>", "<div class=\"paragraph\">\n<p>Summary Next Line</p>\n</div>\n<div class=\"paragraph\">\n<p>Some more text</p>\n</div>"},
		{media.Builtin.AsciiDocType, "<div><p>FOOO</p></div><div><p>First paragraph</p></div>", "", "<div><p>First paragraph</p></div>", "<div><p>First paragraph</p></div>"},
		{media.Builtin.AsciiDocType, "<div><p>First paragraphFOOO</p></div><div><p>Second paragraph</p></div>", "<div><p>First paragraph</p></div>", "<div><p>Second paragraph</p></div>", "<div><p>First paragraph</p></div><div><p>Second paragraph</p></div>"},
	}

	for i, test := range tests {
		summary := ExtractSummaryFromHTMLWithDivider(test.mt, test.input, divider)
		c.Assert(summary.Summary(), qt.Equals, test.expectSummary, qt.Commentf("Summary %d", i))
		c.Assert(summary.ContentWithoutSummary(), qt.Equals, test.expectContentWithoutSummary, qt.Commentf("ContentWithoutSummary %d", i))
		c.Assert(summary.Content(), qt.Equals, test.expectContent, qt.Commentf("Content %d", i))
	}
}

func TestExpandDivider(t *testing.T) {
	c := qt.New(t)

	for i, test := range []struct {
		input           string
		divider         string
		ptag            tagReStartEnd
		expect          string
		expectEndMarkup string
	}{
		{"<p>First paragraph</p>\n<p>FOOO</p>\n<p>Second paragraph</p>", "FOOO", startEndP, "<p>FOOO</p>\n", ""},
		{"<div class=\"paragraph\">\n<p>FOOO</p>\n</div>", "FOOO", startEndDiv, "<div class=\"paragraph\">\n<p>FOOO</p>\n</div>", ""},
		{"<div><p>FOOO</p></div><div><p>Second paragraph</p></div>", "FOOO", startEndDiv, "<div><p>FOOO</p></div>", ""},
		{"<div><p>First paragraphFOOO</p></div><div><p>Second paragraph</p></div>", "FOOO", startEndDiv, "FOOO", "</p></div>"},
		{"   <p> abc FOOO  </p>  ", "FOOO", startEndP, "FOOO", "  </p>"},
		{"   <p>  FOOO  </p>  ", "FOOO", startEndP, "<p>  FOOO  </p>", ""},
		{"   <p>\n  \nFOOO  </p>  ", "FOOO", startEndP, "<p>\n  \nFOOO  </p>", ""},
		{"   <div>  FOOO  </div>  ", "FOOO", startEndDiv, "<div>  FOOO  </div>", ""},
	} {

		l := types.LowHigh[string]{Low: strings.Index(test.input, test.divider), High: strings.Index(test.input, test.divider) + len(test.divider)}
		e, t := expandSummaryDivider(test.input, test.ptag, l)
		c.Assert(test.input[e.Low:e.High], qt.Equals, test.expect, qt.Commentf("[%d] Test.expect %q", i, test.input))
		c.Assert(test.input[t.Low:t.High], qt.Equals, test.expectEndMarkup, qt.Commentf("[%d] Test.expectEndMarkup %q", i, test.input))
	}
}

func BenchmarkSummaryFromHTML(b *testing.B) {
	b.StopTimer()
	input := "<p>First paragraph</p><p>Second paragraph</p>"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		summary := ExtractSummaryFromHTML(media.Builtin.MarkdownType, input, 2, false)
		if s := summary.Content(); s != input {
			b.Fatalf("unexpected content: %q", s)
		}
		if s := summary.ContentWithoutSummary(); s != "<p>Second paragraph</p>" {
			b.Fatalf("unexpected content without summary: %q", s)
		}
		if s := summary.Summary(); s != "<p>First paragraph</p>" {
			b.Fatalf("unexpected summary: %q", s)
		}
	}
}

func BenchmarkSummaryFromHTMLWithDivider(b *testing.B) {
	b.StopTimer()
	input := "<p>First paragraph</p><p>FOOO</p><p>Second paragraph</p>"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		summary := ExtractSummaryFromHTMLWithDivider(media.Builtin.MarkdownType, input, "FOOO")
		if s := summary.Content(); s != "<p>First paragraph</p><p>Second paragraph</p>" {
			b.Fatalf("unexpected content: %q", s)
		}
		if s := summary.ContentWithoutSummary(); s != "<p>Second paragraph</p>" {
			b.Fatalf("unexpected content without summary: %q", s)
		}
		if s := summary.Summary(); s != "<p>First paragraph</p>" {
			b.Fatalf("unexpected summary: %q", s)
		}
	}
}
