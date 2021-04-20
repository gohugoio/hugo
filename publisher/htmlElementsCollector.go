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

	"golang.org/x/net/html"

	"github.com/gohugoio/hugo/helpers"
)

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

type htmlElementsCollector struct {
	// Contains the raw HTML string. We will get the same element
	// several times, and want to avoid costly reparsing when this
	// is used for aggregated data only.
	elementSet map[string]bool

	elements []htmlElement

	mu sync.RWMutex
}

func newHTMLElementsCollector() *htmlElementsCollector {
	return &htmlElementsCollector{
		elementSet: make(map[string]bool),
	}
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
	buff      bytes.Buffer

	isCollecting bool
	inPreTag     string

	inQuote    bool
	quoteValue byte
}

func newHTMLElementsCollectorWriter(collector *htmlElementsCollector) *htmlElementsCollectorWriter {
	return &htmlElementsCollectorWriter{
		collector: collector,
	}
}

// Write splits the incoming stream into single html element.
func (w *htmlElementsCollectorWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	i := 0

	for i < len(p) {
		// If we are not collecting, cycle through byte stream until start bracket "<" is found.
		if !w.isCollecting {
			for ; i < len(p); i++ {
				b := p[i]
				if b == '<' {
					w.startCollecting()
					break
				}
			}
		}

		if w.isCollecting {
			// If we are collecting, cycle through byte stream until end bracket ">" is found,
			// disregard any ">" if within a quote,
			// write bytes until found to buffer.
			for ; i < len(p); i++ {
				b := p[i]
				w.toggleIfQuote(b)
				w.buff.WriteByte(b)

				if !w.inQuote && b == '>' {
					w.endCollecting()
					break
				}
			}
		}

		// If no end bracket ">" is found while collecting, but the stream ended
		// this could mean we received chunks of a stream from e.g. the minify functionality
		// next if loop will be skipped.

		// At this point we have collected an element line between angle brackets "<" and ">".
		if !w.isCollecting {
			if w.buff.Len() == 0 {
				continue
			}

			if w.inPreTag != "" { // within preformatted code block
				s := w.buff.String()
				w.buff.Reset()
				if tagName, isEnd := parseEndTag(s); isEnd && w.inPreTag == tagName {
					w.inPreTag = ""
				}
				continue
			}

			// First check if we have processed this element before.
			w.collector.mu.RLock()

			// Work with the bytes slice as long as it's practical,
			// to save memory allocations.
			b := w.buff.Bytes()

			// See https://github.com/dominikh/go-tools/issues/723
			//lint:ignore S1030 This construct avoids memory allocation for the string.
			seen := w.collector.elementSet[string(b)]
			w.collector.mu.RUnlock()
			if seen {
				w.buff.Reset()
				continue
			}

			// Filter out unwanted tags
			// if within preformatted code blocks <pre>, <textarea>, <script>, <style>
			// comments and doctype tags
			// end tags.
			switch {
			case bytes.HasPrefix(b, []byte("<!")): // comment or doctype tag
				w.buff.Reset()
				continue
			case bytes.HasPrefix(b, []byte("</")): // end tag
				w.buff.Reset()
				continue
			}

			s := w.buff.String()
			w.buff.Reset()

			// Check if a preformatted code block started.
			if tagName, isStart := parseStartTag(s); isStart && isPreFormatted(tagName) {
				w.inPreTag = tagName
			}

			// Parse each collected element.
			el, err := parseHTMLElement(s)
			if err != nil {
				return n, err
			}

			// Write this tag to the element set.
			w.collector.mu.Lock()
			w.collector.elementSet[s] = true
			w.collector.elements = append(w.collector.elements, el)
			w.collector.mu.Unlock()
		}
	}

	return
}

func (c *htmlElementsCollectorWriter) startCollecting() {
	c.isCollecting = true
}

func (c *htmlElementsCollectorWriter) endCollecting() {
	c.isCollecting = false
	c.inQuote = false
}

func (c *htmlElementsCollectorWriter) toggleIfQuote(b byte) {
	if isQuote(b) {
		if c.inQuote && b == c.quoteValue {
			c.inQuote = false
		} else if !c.inQuote {
			c.inQuote = true
			c.quoteValue = b
		}
	}
}

func isQuote(b byte) bool {
	return b == '"' || b == '\''
}

func parseStartTag(s string) (string, bool) {
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")

	spaceIndex := strings.Index(s, " ")
	if spaceIndex != -1 {
		s = s[:spaceIndex]
	}

	return strings.ToLower(strings.TrimSpace(s)), true
}

func parseEndTag(s string) (string, bool) {
	if !strings.HasPrefix(s, "</") {
		return "", false
	}

	s = strings.TrimPrefix(s, "</")
	s = strings.TrimSuffix(s, ">")

	return strings.ToLower(strings.TrimSpace(s)), true
}

// No need to look inside these for HTML elements.
func isPreFormatted(s string) bool {
	return s == "pre" || s == "textarea" || s == "script" || s == "style"
}

type htmlElement struct {
	Tag     string
	Classes []string
	IDs     []string
}

var (
	htmlJsonFixer = strings.NewReplacer(", ", "\n")
	jsonAttrRe    = regexp.MustCompile(`'?(.*?)'?:.*`)
	classAttrRe   = regexp.MustCompile(`(?i)^class$|transition`)

	exceptionList = map[string]bool{
		"thead": true,
		"tbody": true,
		"tfoot": true,
		"td":    true,
		"tr":    true,
	}
)

func parseHTMLElement(elStr string) (el htmlElement, err error) {
	var tagBuffer string = ""

	tagName, ok := parseStartTag(elStr)
	if !ok {
		return
	}

	// The net/html parser does not handle single table elements as input, e.g. tbody.
	// We only care about the element/class/ids, so just store away the original tag name
	// and pretend it's a <div>.
	if exceptionList[tagName] {
		tagBuffer = tagName
		elStr = strings.Replace(elStr, tagName, "div", 1)
	}

	n, err := html.Parse(strings.NewReader(elStr))
	if err != nil {
		return
	}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.Contains(elStr, n.Data) {
			el.Tag = n.Data

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

	// did we replaced the start tag?
	if tagBuffer != "" {
		el.Tag = tagBuffer
	}

	return
}
