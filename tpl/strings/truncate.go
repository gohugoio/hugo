// Copyright 2016 The Hugo Authors. All rights reserved.
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

package strings

import (
	"errors"
	"html"
	"html/template"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/common/hhtml"
	"github.com/spf13/cast"
)

// Truncate truncates the string in s to the specified length.
func (ns *Namespace) Truncate(s any, options ...any) (template.HTML, error) {
	length, err := cast.ToIntE(s)
	if err != nil {
		return "", err
	}
	var textParam any
	var ellipsis string

	switch len(options) {
	case 0:
		return "", errors.New("truncate requires a length and a string")
	case 1:
		textParam = options[0]
		ellipsis = " …"
	case 2:
		textParam = options[1]
		ellipsis, err = cast.ToStringE(options[0])
		if err != nil {
			return "", errors.New("ellipsis must be a string")
		}
		if _, ok := options[0].(template.HTML); !ok {
			ellipsis = html.EscapeString(ellipsis)
		}
	default:
		return "", errors.New("too many arguments passed to truncate")
	}

	text, err := cast.ToStringE(textParam)
	if err != nil {
		return "", errors.New("text must be a string")
	}

	_, isHTML := textParam.(template.HTML)

	if utf8.RuneCountInString(text) <= length {
		if isHTML {
			return template.HTML(text), nil
		}
		return template.HTML(html.EscapeString(text)), nil
	}

	t := truncator{length: length}

	if !isHTML {
		if end, ok := t.scanText(text, 0, len(text)); ok {
			return template.HTML(html.EscapeString(text[:end]) + ellipsis), nil
		}
		return template.HTML(html.EscapeString(text)), nil
	}

	sc := hhtml.NewScanner(text, 0)
	for {
		textStart := sc.Pos()
		tok, ok := sc.Next(len(text))
		textEnd := len(text)
		if ok {
			textEnd = tok.Start
		}
		if end, cut := t.scanText(text, textStart, textEnd); cut {
			return template.HTML(truncateHTML(text, end, ellipsis)), nil
		}
		if !ok {
			return template.HTML(text), nil
		}
		t.lastWordIndex = t.lastNonSpace
		if tok.Type == hhtml.StartTag && !tok.SelfClosing && hhtml.IsRawTextElement(tok.Name) {
			sc.SkipRawText(len(text), tok.Name)
		}
	}
}

// truncateHTML returns s[:end] followed by ellipsis and the end tags of
// the elements left open.
func truncateHTML(s string, end int, ellipsis string) string {
	var buf [16]string
	stack := buf[:0]
	sc := hhtml.NewScanner(s, 0)
	for {
		tok, ok := sc.Next(end)
		if !ok {
			break
		}
		stack = hhtml.UpdateElementStack(&sc, stack, tok, end)
	}

	var b strings.Builder
	b.WriteString(s[:end])
	b.WriteString(ellipsis)
	for _, name := range slices.Backward(stack) {
		b.WriteString("</")
		b.WriteString(name)
		b.WriteString(">")
	}
	return b.String()
}

type truncator struct {
	length        int
	currentLen    int
	lastWordIndex int
	lastNonSpace  int
}

// scanText counts the runes in the text s[start:end].
// If length is exceeded, it returns the position to cut s at and true.
func (t *truncator) scanText(s string, start, end int) (int, bool) {
	for i, r := range s[start:end] {
		i += start
		t.currentLen++
		if unicode.IsSpace(r) {
			t.lastWordIndex = t.lastNonSpace
		} else if unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana) {
			t.lastWordIndex = t.lastNonSpace
			t.lastNonSpace = i + utf8.RuneLen(r)
			if t.currentLen <= t.length {
				t.lastWordIndex = t.lastNonSpace
			}
		} else {
			t.lastNonSpace = i + utf8.RuneLen(r)
		}

		if t.currentLen > t.length {
			if t.lastWordIndex == 0 {
				return i, true
			}
			return t.lastWordIndex, true
		}
	}
	return 0, false
}
