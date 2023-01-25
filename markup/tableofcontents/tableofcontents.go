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

// Headings holds the top level headings.
type Headings []Heading

// Heading holds the data about a heading and its children.
type Heading struct {
	ID   string
	Text string

	Headings Headings
}

// IsZero is true when no ID or Text is set.
func (h Heading) IsZero() bool {
	return h.ID == "" && h.Text == ""
}

// Root implements AddAt, which can be used to build the
// data structure for the ToC.
type Root struct {
	Headings Headings
}

// AddAt adds the heading into the given location.
func (toc *Root) AddAt(h Heading, row, level int) {
	for i := len(toc.Headings); i <= row; i++ {
		toc.Headings = append(toc.Headings, Heading{})
	}

	if level == 0 {
		toc.Headings[row] = h
		return
	}

	heading := &toc.Headings[row]

	for i := 1; i < level; i++ {
		if len(heading.Headings) == 0 {
			heading.Headings = append(heading.Headings, Heading{})
		}
		heading = &heading.Headings[len(heading.Headings)-1]
	}
	heading.Headings = append(heading.Headings, h)
}

// ToHTML renders the ToC as HTML.
func (toc Root) ToHTML(startLevel, stopLevel int, ordered bool) string {
	b := &tocBuilder{
		s:          strings.Builder{},
		h:          toc.Headings,
		startLevel: startLevel,
		stopLevel:  stopLevel,
		ordered:    ordered,
	}
	b.Build()
	return b.s.String()
}

type tocBuilder struct {
	s strings.Builder
	h Headings

	startLevel int
	stopLevel  int
	ordered    bool
}

func (b *tocBuilder) Build() {
	b.writeNav(b.h)
}

func (b *tocBuilder) writeNav(h Headings) {
	if len(h) == 0 {
		return
	}
	b.s.WriteString("<nav id=\"TableOfContents\">")
	b.writeHeadings(1, 0, h)
	b.s.WriteString("</nav>")
}

func (b *tocBuilder) writeHeadings(level, indent int, h Headings) {
	if level < b.startLevel {
		for _, h := range h {
			b.writeHeadings(level+1, indent, h.Headings)
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
		b.writeHeading(level+1, indent+2, h)
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

func (b *tocBuilder) writeHeading(level, indent int, h Heading) {
	b.indent(indent)
	b.s.WriteString("<li>")
	if !h.IsZero() {
		b.s.WriteString("<a href=\"#" + h.ID + "\">" + h.Text + "</a>")
	}
	b.writeHeadings(level, indent, h.Headings)
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
