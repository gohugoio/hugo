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
	"regexp"

	"github.com/gohugoio/hugo/helpers"
	"golang.org/x/net/html"

	"bytes"
	"sort"
	"strings"
	"sync"
)

func newHTMLElementsCollector() *htmlElementsCollector {
	return &htmlElementsCollector{
		elementSet: make(map[string]bool),
	}
}

func newHTMLElementsCollectorWriter(collector *htmlElementsCollector) *cssClassCollectorWriter {
	return &cssClassCollectorWriter{
		collector: collector,
	}
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

type cssClassCollectorWriter struct {
	collector *htmlElementsCollector
	buff      bytes.Buffer

	isCollecting bool
	dropValue    bool

	inQuote    bool
	quoteValue byte
}

func (w *cssClassCollectorWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	i := 0

	for i < len(p) {
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
			for ; i < len(p); i++ {
				b := p[i]
				w.toggleIfQuote(b)
				if !w.inQuote && b == '>' {
					w.endCollecting(false)
					break
				}
				w.buff.WriteByte(b)
			}

			if !w.isCollecting {
				if w.dropValue {
					w.buff.Reset()
				} else {
					// First check if we have processed this element before.
					w.collector.mu.RLock()

					// See https://github.com/dominikh/go-tools/issues/723
					//lint:ignore S1030 This construct avoids memory allocation for the string.
					seen := w.collector.elementSet[string(w.buff.Bytes())]
					w.collector.mu.RUnlock()
					if seen {
						w.buff.Reset()
						continue
					}

					s := w.buff.String()

					w.buff.Reset()

					if strings.HasPrefix(s, "</") {
						continue
					}

					key := s

					s, tagName := w.insertStandinHTMLElement(s)
					el := parseHTMLElement(s)
					el.Tag = tagName

					w.collector.mu.Lock()
					w.collector.elementSet[key] = true
					if el.Tag != "" {
						w.collector.elements = append(w.collector.elements, el)
					}
					w.collector.mu.Unlock()
				}
			}
		}
	}

	return
}

// The net/html parser does not handle single table elemnts as input, e.g. tbody.
// We only care about the element/class/ids, so just store away the original tag name
// and pretend it's a <div>.
func (c *cssClassCollectorWriter) insertStandinHTMLElement(el string) (string, string) {
	tag := el[1:]
	spacei := strings.Index(tag, " ")
	if spacei != -1 {
		tag = tag[:spacei]
	}
	newv := strings.Replace(el, tag, "div", 1)
	return newv, strings.ToLower(tag)

}

func (c *cssClassCollectorWriter) endCollecting(drop bool) {
	c.isCollecting = false
	c.inQuote = false
	c.dropValue = drop
}

func (c *cssClassCollectorWriter) startCollecting() {
	c.isCollecting = true
	c.dropValue = false
}

func (c *cssClassCollectorWriter) toggleIfQuote(b byte) {
	if isQuote(b) {
		if c.inQuote && b == c.quoteValue {
			c.inQuote = false
		} else if !c.inQuote {
			c.inQuote = true
			c.quoteValue = b
		}
	}
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

func isQuote(b byte) bool {
	return b == '"' || b == '\''
}

var (
	htmlJsonFixer = strings.NewReplacer(", ", "\n")
	jsonAttrRe    = regexp.MustCompile(`'?(.*?)'?:.*`)
	classAttrRe   = regexp.MustCompile(`(?i)^class$|transition`)
)

func parseHTMLElement(elStr string) (el htmlElement) {
	elStr = strings.TrimSpace(elStr)
	if !strings.HasSuffix(elStr, ">") {
		elStr += ">"
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

	return
}
