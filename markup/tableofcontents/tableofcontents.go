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
	"strconv"
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
func (toc *Root) AddAt(h Header, row, level int) {
	for i := len(toc.Headers); i <= row; i++ {
		toc.Headers = append(toc.Headers, Header{})
	}

	if level == 0 {
		toc.Headers[row] = h
		return
	}

	header := &toc.Headers[row]

	for i := 1; i < level; i++ {
		if len(header.Headers) == 0 {
			header.Headers = append(header.Headers, Header{})
		}
		header = &header.Headers[len(header.Headers)-1]
	}
	header.Headers = append(header.Headers, h)
}

func GetDefault(val string, defVal string) string {
	if len(val) == 0 {
		return defVal
	}
	return val
}

// ToHTML renders the ToC as HTML.
func (toc Root) ToHTML(startLevel, stopLevel int, ordered bool, writeLevels bool, wrapperElement string, wrapperId string, wrapperClass string) string {
	b := &tocBuilder{
		s:              strings.Builder{},
		h:              toc.Headers,
		startLevel:     startLevel,
		stopLevel:      stopLevel,
		ordered:        ordered,
		writeLevels:    writeLevels,
		wrapperElement: GetDefault(wrapperElement, "nav"),
		wrapperId:      GetDefault(wrapperId, "TableOfContents"),
		wrapperClass:   wrapperClass,
	}
	b.Build()
	return b.s.String()
}

type tocBuilder struct {
	s strings.Builder
	h Headers

	startLevel     int
	stopLevel      int
	ordered        bool
	writeLevels    bool
	wrapperElement string
	wrapperId      string
	wrapperClass   string
}

func (b *tocBuilder) Build() {
	b.writeNav(b.h)
}

func (b *tocBuilder) writeNav(h Headers) {
	wrapperClass := ""
	if len(b.wrapperClass) > 0 {
		wrapperClass = " class=\"" + b.wrapperClass + "\""
	}
	b.s.WriteString("<" + b.wrapperElement + wrapperClass + " id=\"" + b.wrapperId + "\">")
	b.writeHeaders(1, 0, b.h)
	b.s.WriteString("</" + b.wrapperElement + ">")
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
		levelAttr := ""
		if b.writeLevels {
			levelAttr = " data-level=\"" + strconv.Itoa(level) + "\""
		}
		if b.ordered {
			b.s.WriteString("<ol" + levelAttr + ">\n")
		} else {
			b.s.WriteString("<ul" + levelAttr + ">\n")
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
	StartLevel:     2,
	EndLevel:       3,
	Ordered:        false,
	WriteLevels:    false,
	WrapperElement: "nav",
	WrapperId:      "TableOfContents",
	WrapperClass:   "",
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

	// If set to true each <ul> element will provide an attribute 'data-level'
	// with the corresponding level (eg. h2 -> 2, h3 -> 3, ...)
	WriteLevels bool

	// Allows to specify a different element than <nav> to be used for the toc
	WrapperElement string

	// The ID to be used (default: 'TableOfContents')
	WrapperId string

	// CSS class which will be applied to the wrapper
	WrapperClass string
}
