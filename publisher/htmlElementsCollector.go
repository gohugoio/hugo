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
	"github.com/gohugoio/hugo/helpers"
	"golang.org/x/net/html"
	yaml "gopkg.in/yaml.v2"

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
	inQuote      bool
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

					el := parseHTMLElement(s)

					w.collector.mu.Lock()
					w.collector.elementSet[s] = true
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
		c.inQuote = !c.inQuote
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

var htmlJsonFixer = strings.NewReplacer(", ", "\n")

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
					if strings.EqualFold(a.Key, "class") {
						el.Classes = append(el.Classes, strings.Fields(a.Val)...)
					} else {
						key := strings.ToLower(a.Key)
						val := strings.TrimSpace(a.Val)
						if strings.Contains(key, "class") && strings.HasPrefix(val, "{") {
							// This looks like a Vue or AlpineJS class binding.
							// Try to unmarshal it as YAML and pull the keys.
							// This may look odd, as the source is (probably) JS (JSON), but the YAML
							// parser is much more lenient with simple JS input, it seems.
							m := make(map[string]interface{})
							val = htmlJsonFixer.Replace(strings.Trim(val, "{}"))
							// Remove leading space to make it look like YAML.
							lines := strings.Split(val, "\n")
							for i, l := range lines {
								lines[i] = strings.TrimSpace(l)
							}
							val = strings.Join(lines, "\n")
							err := yaml.Unmarshal([]byte(val), &m)
							if err == nil {
								for k := range m {
									el.Classes = append(el.Classes, strings.Fields(k)...)
								}
							} else {
								// Just insert the raw values. This is used for CSS class pruning
								// so, it's important not to leave out values that may be a CSS class.
								el.Classes = append(el.Classes, strings.Fields(val)...)
							}
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
