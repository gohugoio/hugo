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

package tableofcontents

import (
	"strings"
)

// Headers holds the top level (h1) headers.
type Headers []Header

// Header holds the data about a header and its children.
type Header struct {
	ID   string
	Text string

	Headers Headers
}

// IsZero is true when no ID or Text is set.
func (h Header) IsZero() bool {
	return h.ID == "" && h.Text == ""
}

// Root implements AddAt, which can be used to build the
// data structure for the ToC.
type Root struct {
	Headers Headers
}

// AddAt adds the header into the given location.
func (toc *Root) AddAt(h Header, y, x int) {
	for i := len(toc.Headers); i <= y; i++ {
		toc.Headers = append(toc.Headers, Header{})
	}

	if x == 0 {
		toc.Headers[y] = h
		return
	}

	header := &toc.Headers[y]

	for i := 1; i < x; i++ {
		if len(header.Headers) == 0 {
			header.Headers = append(header.Headers, Header{})
		}
		header = &header.Headers[len(header.Headers)-1]
	}
	header.Headers = append(header.Headers, h)
}

// ToHTML renders the ToC as HTML.
func (toc Root) ToHTML(startLevel, stopLevel int, ordered bool) string {
	b := &tocBuilder{
		s:          strings.Builder{},
		h:          toc.Headers,
		startLevel: startLevel,
		stopLevel:  stopLevel,
		ordered:    ordered,
	}
	b.Build()
	return b.s.String()
}

type tocBuilder struct {
	s strings.Builder
	h Headers

	startLevel int
	stopLevel  int
	ordered    bool
}

func (b *tocBuilder) Build() {
	b.writeNav(b.h)
}

func (b *tocBuilder) writeNav(h Headers) {
	b.s.WriteString("<nav id=\"TableOfContents\">")
	b.writeHeaders(1, 0, b.h)
	b.s.WriteString("</nav>")
}

func (b *tocBuilder) writeHeaders(level, indent int, h Headers) {
	if level < b.startLevel {
		for _, h := range h {
			b.writeHeaders(level+1, indent, h.Headers)
		}
		return
	}

	if b.stopLevel != -1 && level > b.stopLevel {
		return
	}

	hasChildren := len(h) > 0

	if hasChildren {
		b.s.WriteString("\n")
		b.indent(indent + 1)
		if b.ordered {
			b.s.WriteString("<ol>\n")
		} else {
			b.s.WriteString("<ul>\n")
		}
	}

	for _, h := range h {
		b.writeHeader(level+1, indent+2, h)
	}

	if hasChildren {
		b.indent(indent + 1)
		if b.ordered {
			b.s.WriteString("</ol>")
		} else {
			b.s.WriteString("</ul>")
		}
		b.s.WriteString("\n")
		b.indent(indent)
	}

}
func (b *tocBuilder) writeHeader(level, indent int, h Header) {
	b.indent(indent)
	b.s.WriteString("<li>")
	if !h.IsZero() {
		b.s.WriteString("<a href=\"#" + h.ID + "\">" + h.Text + "</a>")
	}
	b.writeHeaders(level, indent, h.Headers)
	b.s.WriteString("</li>\n")
}

func (b *tocBuilder) indent(n int) {
	for i := 0; i < n; i++ {
		b.s.WriteString("  ")
	}
}

// DefaultConfig is the default ToC configuration.
var DefaultConfig = Config{
	StartLevel: 2,
	EndLevel:   3,
	Ordered:    false,
}

type Config struct {
	// Heading start level to include in the table of contents, starting
	// at h1 (inclusive).
	StartLevel int

	// Heading end level, inclusive, to include in the table of contents.
	// Default is 3, a value of -1 will include everything.
	EndLevel int

	// Whether to produce a ordered list or not.
	Ordered bool
}
