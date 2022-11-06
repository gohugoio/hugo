// Copyright 2020 The Hugo Authors. All rights reserved.
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

package publisher

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/html"

	"github.com/gohugoio/hugo/helpers"
)

const eof = -1

var (
	htmlJsonFixer = strings.NewReplacer(", ", "\n")
	jsonAttrRe    = regexp.MustCompile(`'?(.*?)'?:.*`)
	classAttrRe   = regexp.MustCompile(`(?i)^class$|transition`)

	skipInnerElementRe = regexp.MustCompile(`(?i)^(pre|textarea|script|style)`)
	skipAllElementRe   = regexp.MustCompile(`(?i)^!DOCTYPE`)
	endTagRe           = regexp.MustCompile(`(?i)<\/\s*([a-zA-Z]+)\s*>$`)

	exceptionList = map[string]bool{
		"thead": true,
		"tbody": true,
		"tfoot": true,
		"td":    true,
		"tr":    true,
	}
)

func newHTMLElementsCollector() *htmlElementsCollector {
	return &htmlElementsCollector{
		elementSet: make(map[string]bool),
	}
}

func newHTMLElementsCollectorWriter(collector *htmlElementsCollector) *htmlElementsCollectorWriter {
	w := &htmlElementsCollectorWriter{
		collector: collector,
		state:     htmlLexStart,
	}

	w.defaultLexElementInside = w.lexElementInside(htmlLexStart)

	return w
}

// HTMLElements holds lists of tags and attribute values for classes and id.
type HTMLElements struct {
	Tags    []string `json:"tags"`
	Classes []string `json:"classes"`
	IDs     []string `json:"ids"`
}

func (h *HTMLElements) Merge(other HTMLElements) {
	h.Tags = append(h.Tags, other.Tags...)
	h.Classes = append(h.Classes, other.Classes...)
	h.IDs = append(h.IDs, other.IDs...)

	h.Tags = helpers.UniqueStringsReuse(h.Tags)
	h.Classes = helpers.UniqueStringsReuse(h.Classes)
	h.IDs = helpers.UniqueStringsReuse(h.IDs)
}

func (h *HTMLElements) Sort() {
	sort.Strings(h.Tags)
	sort.Strings(h.Classes)
	sort.Strings(h.IDs)
}

type htmlElement struct {
	Tag     string
	Classes []string
	IDs     []string
}

type htmlElementsCollector struct {
	// Contains the raw HTML string. We will get the same element
	// several times, and want to avoid costly reparsing when this
	// is used for aggregated data only.
	elementSet map[string]bool

	elements []htmlElement

	mu sync.RWMutex
}

func (c *htmlElementsCollector) getHTMLElements() HTMLElements {
	var (
		classes []string
		ids     []string
		tags    []string
	)

	for _, el := range c.elements {
		classes = append(classes, el.Classes...)
		ids = append(ids, el.IDs...)
		tags = append(tags, el.Tag)
	}

	classes = helpers.UniqueStringsSorted(classes)
	ids = helpers.UniqueStringsSorted(ids)
	tags = helpers.UniqueStringsSorted(tags)

	els := HTMLElements{
		Classes: classes,
		IDs:     ids,
		Tags:    tags,
	}

	return els
}

type htmlElementsCollectorWriter struct {
	collector *htmlElementsCollector

	r     rune   // Current rune
	width int    // The width in bytes of r
	input []byte // The current slice written to Write
	pos   int    // The current position in input

	err error

	inQuote rune

	buff bytes.Buffer

	// Current state
	state htmlCollectorStateFunc

	// Precompiled state funcs
	defaultLexElementInside htmlCollectorStateFunc
}

// Write collects HTML elements from p, which must contain complete runes.
func (w *htmlElementsCollectorWriter) Write(p []byte) (int, error) {
	if p == nil {
		return 0, nil
	}

	w.input = p

	for {
		w.r = w.next()
		if w.r == eof || w.r == utf8.RuneError {
			break
		}
		w.state = w.state(w)
	}

	w.pos = 0
	w.input = nil

	return len(p), nil
}

func (l *htmlElementsCollectorWriter) backup() {
	l.pos -= l.width
	l.r, _ = utf8.DecodeRune(l.input[l.pos:])
}

func (w *htmlElementsCollectorWriter) consumeBuffUntil(condition func() bool, resolve htmlCollectorStateFunc) htmlCollectorStateFunc {
	var s htmlCollectorStateFunc
	s = func(*htmlElementsCollectorWriter) htmlCollectorStateFunc {
		w.buff.WriteRune(w.r)
		if condition() {
			w.buff.Reset()
			return resolve
		}
		return s
	}
	return s
}

func (w *htmlElementsCollectorWriter) consumeRuneUntil(condition func(r rune) bool, resolve htmlCollectorStateFunc) htmlCollectorStateFunc {
	var s htmlCollectorStateFunc
	s = func(*htmlElementsCollectorWriter) htmlCollectorStateFunc {
		if condition(w.r) {
			return resolve
		}
		return s
	}
	return s
}

// Starts with e.g. "<body " or "<div"
func (w *htmlElementsCollectorWriter) lexElementInside(resolve htmlCollectorStateFunc) htmlCollectorStateFunc {
	var s htmlCollectorStateFunc
	s = func(w *htmlElementsCollectorWriter) htmlCollectorStateFunc {
		w.buff.WriteRune(w.r)

		// Skip any text inside a quote.
		if w.r == '\'' || w.r == '"' {
			if w.inQuote == w.r {
				w.inQuote = 0
			} else if w.inQuote == 0 {
				w.inQuote = w.r
			}
		}

		if w.inQuote != 0 {
			return s
		}

		if w.r == '>' {

			// Work with the bytes slice as long as it's practical,
			// to save memory allocations.
			b := w.buff.Bytes()

			defer func() {
				w.buff.Reset()
			}()

			// First check if we have processed this element before.
			w.collector.mu.RLock()

			seen := w.collector.elementSet[string(b)]
			w.collector.mu.RUnlock()
			if seen {
				return resolve
			}

			s := w.buff.String()

			if s == "" {
				return resolve
			}

			// Parse each collected element.
			el, err := parseHTMLElement(s)
			if err != nil {
				w.err = err
				return resolve
			}

			// Write this tag to the element set.
			w.collector.mu.Lock()
			w.collector.elementSet[s] = true
			w.collector.elements = append(w.collector.elements, el)
			w.collector.mu.Unlock()

			return resolve

		}

		return s
	}

	return s
}

func (l *htmlElementsCollectorWriter) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	runeValue, runeWidth := utf8.DecodeRune(l.input[l.pos:])

	l.width = runeWidth
	l.pos += l.width
	return runeValue
}

// returns the next state in HTML element scanner.
type htmlCollectorStateFunc func(*htmlElementsCollectorWriter) htmlCollectorStateFunc

// At "<", buffer empty.
// Potentially starting a HTML element.
func htmlLexElementStart(w *htmlElementsCollectorWriter) htmlCollectorStateFunc {
	if w.r == '>' || unicode.IsSpace(w.r) {
		if w.buff.Len() < 2 || bytes.HasPrefix(w.buff.Bytes(), []byte("</")) {
			w.buff.Reset()
			return htmlLexStart
		}

		tagName := w.buff.Bytes()[1:]

		switch {
		case skipInnerElementRe.Match(tagName):
			// pre, script etc. We collect classes etc. on the surrounding
			// element, but skip the inner content.
			w.backup()

			// tagName will be overwritten, so make a copy.
			tagNameCopy := make([]byte, len(tagName))
			copy(tagNameCopy, tagName)

			return w.lexElementInside(
				w.consumeBuffUntil(
					func() bool {
						if w.r != '>' {
							return false
						}
						m := endTagRe.FindSubmatch(w.buff.Bytes())
						if m == nil {
							return false
						}
						return bytes.EqualFold(m[1], tagNameCopy)
					},
					htmlLexStart,
				))
		case skipAllElementRe.Match(tagName):
			// E.g. "<!DOCTYPE ..."
			w.buff.Reset()
			return w.consumeRuneUntil(func(r rune) bool {
				return r == '>'
			}, htmlLexStart)
		default:
			w.backup()
			return w.defaultLexElementInside
		}
	}

	w.buff.WriteRune(w.r)

	// If it's a comment, skip to its end.
	if w.r == '-' && bytes.Equal(w.buff.Bytes(), []byte("<!--")) {
		w.buff.Reset()
		return htmlLexToEndOfComment
	}

	return htmlLexElementStart
}

// Entry state func.
// Looks for a opening bracket, '<'.
func htmlLexStart(w *htmlElementsCollectorWriter) htmlCollectorStateFunc {
	if w.r == '<' {
		w.backup()
		w.buff.Reset()
		return htmlLexElementStart
	}

	return htmlLexStart
}

// After "<!--", buff empty.
func htmlLexToEndOfComment(w *htmlElementsCollectorWriter) htmlCollectorStateFunc {
	w.buff.WriteRune(w.r)

	if w.r == '>' && bytes.HasSuffix(w.buff.Bytes(), []byte("-->")) {
		// Done, start looking for HTML elements again.
		return htmlLexStart
	}

	return htmlLexToEndOfComment
}

func parseHTMLElement(elStr string) (el htmlElement, err error) {

	tagName := parseStartTag(elStr)

	el.Tag = strings.ToLower(tagName)
	tagNameToParse := el.Tag

	// The net/html parser does not handle single table elements as input, e.g. tbody.
	// We only care about the element/class/ids, so just store away the original tag name
	// and pretend it's a <div>.
	if exceptionList[el.Tag] {
		elStr = strings.Replace(elStr, tagName, "div", 1)
		tagNameToParse = "div"
	}

	n, err := html.Parse(strings.NewReader(elStr))
	if err != nil {
		return
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tagNameToParse {
			for _, a := range n.Attr {
				switch {
				case strings.EqualFold(a.Key, "id"):
					// There should be only one, but one never knows...
					el.IDs = append(el.IDs, a.Val)
				default:
					if classAttrRe.MatchString(a.Key) {
						el.Classes = append(el.Classes, strings.Fields(a.Val)...)
					} else {
						key := strings.ToLower(a.Key)
						val := strings.TrimSpace(a.Val)
						if strings.Contains(key, "class") && strings.HasPrefix(val, "{") {
							// This looks like a Vue or AlpineJS class binding.
							val = htmlJsonFixer.Replace(strings.Trim(val, "{}"))
							lines := strings.Split(val, "\n")
							for i, l := range lines {
								lines[i] = strings.TrimSpace(l)
							}
							val = strings.Join(lines, "\n")
							val = jsonAttrRe.ReplaceAllString(val, "$1")
							el.Classes = append(el.Classes, strings.Fields(val)...)
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(n)

	return
}

// Variants of s
//
//	<body class="b a">
//	<div>
func parseStartTag(s string) string {
	spaceIndex := strings.IndexFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	if spaceIndex == -1 {
		return s[1 : len(s)-1]
	}

	return s[1:spaceIndex]
}
