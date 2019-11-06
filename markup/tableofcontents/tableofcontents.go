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

type Headers []Header

type Header struct {
	ID   string
	Text string

	Headers Headers
}

func (h Header) IsZero() bool {
	return h.ID == "" && h.Text == ""
}

type Root struct {
	Headers Headers
}

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

func (toc Root) ToHTML(startLevel, stopLevel int) string {
	b := &tocBuilder{
		s:          strings.Builder{},
		h:          toc.Headers,
		startLevel: startLevel,
		stopLevel:  stopLevel,
	}
	b.Build()
	return b.s.String()
}

type tocBuilder struct {
	s strings.Builder
	h Headers

	startLevel int
	stopLevel  int
}

func (b *tocBuilder) Build() {
	b.buildHeaders2(b.h)
}

func (b *tocBuilder) buildHeaders2(h Headers) {
	b.s.WriteString("<nav id=\"TableOfContents\">")
	b.buildHeaders(1, 0, b.h)
	b.s.WriteString("</nav>")
}

func (b *tocBuilder) buildHeaders(level, indent int, h Headers) {
	if level < b.startLevel {
		for _, h := range h {
			b.buildHeaders(level+1, indent, h.Headers)
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
		b.s.WriteString("<ul>\n")
	}

	for _, h := range h {
		b.buildHeader(level+1, indent+2, h)
	}

	if hasChildren {
		b.indent(indent + 1)
		b.s.WriteString("</ul>")
		b.s.WriteString("\n")
		b.indent(indent)
	}

}
func (b *tocBuilder) buildHeader(level, indent int, h Header) {
	b.indent(indent)
	b.s.WriteString("<li>")
	if !h.IsZero() {
		b.s.WriteString("<a href=\"#" + h.ID + "\">" + h.Text + "</a>")
	}
	b.buildHeaders(level, indent, h.Headers)
	b.s.WriteString("</li>\n")
}

func (b *tocBuilder) indent(n int) {
	for i := 0; i < n; i++ {
		b.s.WriteString("  ")
	}
}

var DefaultConfig = Config{
	StartLevel: 2,
	EndLevel:   3,
}

type Config struct {
	// Heading start level to include in the table of contents, starting
	// at h1 (inclusive).
	StartLevel int

	// Heading end level, inclusive, to include in the table of contents.
	// Default is 3, a value of -1 will include everything.
	EndLevel int
}
