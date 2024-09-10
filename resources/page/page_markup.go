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
	"context"
	"html/template"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/tpl"
)

type Content interface {
	Content(context.Context) (template.HTML, error)
	ContentWithoutSummary(context.Context) (template.HTML, error)
	Summary(context.Context) (Summary, error)
	Plain(context.Context) string
	PlainWords(context.Context) []string
	WordCount(context.Context) int
	FuzzyWordCount(context.Context) int
	ReadingTime(context.Context) int
	Len(context.Context) int
}

type Markup interface {
	Render(context.Context) (Content, error)
	RenderString(ctx context.Context, args ...any) (template.HTML, error)
	RenderShortcodes(context.Context) (template.HTML, error)
	Fragments(context.Context) *tableofcontents.Fragments
}

var _ types.PrintableValueProvider = Summary{}

const (
	SummaryTypeAuto        = "auto"
	SummaryTypeManual      = "manual"
	SummaryTypeFrontMatter = "frontmatter"
)

type Summary struct {
	Text      template.HTML
	Type      string // "auto", "manual" or "frontmatter"
	Truncated bool
}

func (s Summary) IsZero() bool {
	return s.Text == ""
}

func (s Summary) PrintableValue() any {
	return s.Text
}

var _ types.PrintableValueProvider = (*Summary)(nil)

type HtmlSummary struct {
	source         string
	SummaryLowHigh types.LowHigh[string]
	SummaryEndTag  types.LowHigh[string]
	WrapperStart   types.LowHigh[string]
	WrapperEnd     types.LowHigh[string]
	Divider        types.LowHigh[string]
}

func (s HtmlSummary) wrap(ss string) string {
	if s.WrapperStart.IsZero() {
		return ss
	}
	return s.source[s.WrapperStart.Low:s.WrapperStart.High] + ss + s.source[s.WrapperEnd.Low:s.WrapperEnd.High]
}

func (s HtmlSummary) wrapLeft(ss string) string {
	if s.WrapperStart.IsZero() {
		return ss
	}

	return s.source[s.WrapperStart.Low:s.WrapperStart.High] + ss
}

func (s HtmlSummary) Value(l types.LowHigh[string]) string {
	return s.source[l.Low:l.High]
}

func (s HtmlSummary) trimSpace(ss string) string {
	return strings.TrimSpace(ss)
}

func (s HtmlSummary) Content() string {
	if s.Divider.IsZero() {
		return s.source
	}
	ss := s.source[:s.Divider.Low]
	ss += s.source[s.Divider.High:]
	return s.trimSpace(ss)
}

func (s HtmlSummary) Summary() string {
	if s.Divider.IsZero() {
		return s.trimSpace(s.wrap(s.Value(s.SummaryLowHigh)))
	}
	ss := s.source[s.SummaryLowHigh.Low:s.Divider.Low]
	if s.SummaryLowHigh.High > s.Divider.High {
		ss += s.source[s.Divider.High:s.SummaryLowHigh.High]
	}
	if !s.SummaryEndTag.IsZero() {
		ss += s.Value(s.SummaryEndTag)
	}
	return s.trimSpace(s.wrap(ss))
}

func (s HtmlSummary) ContentWithoutSummary() string {
	if s.Divider.IsZero() {
		if s.SummaryLowHigh.Low == s.WrapperStart.High && s.SummaryLowHigh.High == s.WrapperEnd.Low {
			return ""
		}
		return s.trimSpace(s.wrapLeft(s.source[s.SummaryLowHigh.High:]))
	}
	if s.SummaryEndTag.IsZero() {
		return s.trimSpace(s.wrapLeft(s.source[s.Divider.High:]))
	}
	return s.trimSpace(s.wrapLeft(s.source[s.SummaryEndTag.High:]))
}

func (s HtmlSummary) Truncated() bool {
	return s.SummaryLowHigh.High < len(s.source)
}

func (s *HtmlSummary) resolveParagraphTagAndSetWrapper(mt media.Type) tagReStartEnd {
	ptag := startEndP

	switch mt.SubType {
	case media.DefaultContentTypes.AsciiDoc.SubType:
		ptag = startEndDiv
	case media.DefaultContentTypes.ReStructuredText.SubType:
		const markerStart = "<div class=\"document\">"
		const markerEnd = "</div>"
		i1 := strings.Index(s.source, markerStart)
		i2 := strings.LastIndex(s.source, markerEnd)
		if i1 > -1 && i2 > -1 {
			s.WrapperStart = types.LowHigh[string]{Low: 0, High: i1 + len(markerStart)}
			s.WrapperEnd = types.LowHigh[string]{Low: i2, High: len(s.source)}
		}
	}
	return ptag
}

// Avoid counting words that are most likely HTML tokens.
var (
	isProbablyHTMLTag      = regexp.MustCompile(`^<\/?[A-Za-z]+>?$`)
	isProablyHTMLAttribute = regexp.MustCompile(`^[A-Za-z]+=["']`)
)

func isProbablyHTMLToken(s string) bool {
	return s == ">" || isProbablyHTMLTag.MatchString(s) || isProablyHTMLAttribute.MatchString(s)
}

// ExtractSummaryFromHTML extracts a summary from the given HTML content.
func ExtractSummaryFromHTML(mt media.Type, input string, numWords int, isCJK bool) (result HtmlSummary) {
	result.source = input
	ptag := result.resolveParagraphTagAndSetWrapper(mt)

	if numWords <= 0 {
		return result
	}

	var count int

	countWord := func(word string) int {
		word = strings.TrimSpace(word)
		if len(word) == 0 {
			return 0
		}
		if isProbablyHTMLToken(word) {
			return 0
		}

		if isCJK {
			word = tpl.StripHTML(word)
			runeCount := utf8.RuneCountInString(word)
			if len(word) == runeCount {
				return 1
			} else {
				return runeCount
			}
		}

		return 1
	}

	high := len(input)
	if result.WrapperEnd.Low > 0 {
		high = result.WrapperEnd.Low
	}

	for j := result.WrapperStart.High; j < high; {
		s := input[j:]
		closingIndex := strings.Index(s, "</"+ptag.tagName+">")

		if closingIndex == -1 {
			break
		}

		s = s[:closingIndex]

		// Count the words in the current paragraph.
		var wi int

		for i, r := range s {
			if unicode.IsSpace(r) || (i+utf8.RuneLen(r) == len(s)) {
				word := s[wi:i]
				count += countWord(word)
				wi = i
				if count >= numWords {
					break
				}
			}
		}

		if count >= numWords {
			result.SummaryLowHigh = types.LowHigh[string]{
				Low:  result.WrapperStart.High,
				High: j + closingIndex + len(ptag.tagName) + 3,
			}
			return
		}

		j += closingIndex + len(ptag.tagName) + 2

	}

	result.SummaryLowHigh = types.LowHigh[string]{
		Low:  result.WrapperStart.High,
		High: high,
	}

	return
}

// ExtractSummaryFromHTMLWithDivider extracts a summary from the given HTML content with
// a manual summary divider.
func ExtractSummaryFromHTMLWithDivider(mt media.Type, input, divider string) (result HtmlSummary) {
	result.source = input
	result.Divider.Low = strings.Index(input, divider)
	result.Divider.High = result.Divider.Low + len(divider)

	if result.Divider.Low == -1 {
		// No summary.
		return
	}

	ptag := result.resolveParagraphTagAndSetWrapper(mt)

	if !mt.IsHTML() {
		result.Divider, result.SummaryEndTag = expandSummaryDivider(result.source, ptag, result.Divider)
	}

	result.SummaryLowHigh = types.LowHigh[string]{
		Low:  result.WrapperStart.High,
		High: result.Divider.Low,
	}

	return
}

var (
	pOrDiv = regexp.MustCompile(`<p[^>]?>|<div[^>]?>$`)

	startEndDiv = tagReStartEnd{
		startEndOfString: regexp.MustCompile(`<div[^>]*?>$`),
		endEndOfString:   regexp.MustCompile(`</div>$`),
		tagName:          "div",
	}

	startEndP = tagReStartEnd{
		startEndOfString: regexp.MustCompile(`<p[^>]*?>$`),
		endEndOfString:   regexp.MustCompile(`</p>$`),
		tagName:          "p",
	}
)

type tagReStartEnd struct {
	startEndOfString *regexp.Regexp
	endEndOfString   *regexp.Regexp
	tagName          string
}

func expandSummaryDivider(s string, re tagReStartEnd, divider types.LowHigh[string]) (types.LowHigh[string], types.LowHigh[string]) {
	var endMarkup types.LowHigh[string]

	if divider.IsZero() {
		return divider, endMarkup
	}

	lo, hi := divider.Low, divider.High

	var preserveEndMarkup bool

	// Find the start of the paragraph.

	for i := lo - 1; i >= 0; i-- {
		if s[i] == '>' {
			if match := re.startEndOfString.FindString(s[:i+1]); match != "" {
				lo = i - len(match) + 1
				break
			}
			if match := pOrDiv.FindString(s[:i+1]); match != "" {
				i -= len(match) - 1
				continue
			}
		}

		r, _ := utf8.DecodeRuneInString(s[i:])
		if !unicode.IsSpace(r) {
			preserveEndMarkup = true
			break
		}
	}

	divider.Low = lo

	// Now walk forward to the end of the paragraph.
	for ; hi < len(s); hi++ {
		if s[hi] != '>' {
			continue
		}
		if match := re.endEndOfString.FindString(s[:hi+1]); match != "" {
			hi++
			break
		}
	}

	if preserveEndMarkup {
		endMarkup.Low = divider.High
		endMarkup.High = hi
	} else {
		divider.High = hi
	}

	// Consume trailing newline if any.
	if divider.High < len(s) && s[divider.High] == '\n' {
		divider.High++
	}

	return divider, endMarkup
}
