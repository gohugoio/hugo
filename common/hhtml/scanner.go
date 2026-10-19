// Copyright 2026 The Hugo Authors. All rights reserved.
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

// Package hhtml contains small, allocation-free HTML helpers.
package hhtml

import (
	"strings"
)

// TokenType is the type of a Token.
type TokenType int

const (
	// StartTag is a start tag, e.g. <p>.
	StartTag TokenType = iota + 1
	// EndTag is an end tag, e.g. </p>.
	EndTag
	// Markup is a comment, doctype, CDATA section or processing instruction.
	Markup
)

// Token is a tag or other markup found by Scanner.
type Token struct {
	Type TokenType
	// Name is the lowercased tag name. It is empty for Markup.
	Name string
	// SelfClosing is set for start tags ending in "/>".
	SelfClosing bool
	// Start is the position of the '<'.
	Start int
}

// Scanner is a minimal, allocation-free scanner for HTML tags.
// It handles quoted attribute values containing '>' and treats a '<'
// that cannot start a tag (e.g. "I <3 Hugo") as text.
type Scanner struct {
	s   string
	pos int
}

// NewScanner creates a new Scanner for s starting at pos.
func NewScanner(s string, pos int) Scanner {
	return Scanner{s: s, pos: pos}
}

// Pos returns the position right after the last token returned by Next.
func (sc *Scanner) Pos() int {
	return sc.pos
}

// Next returns the next token starting before limit.
// ok is false when no more tokens are found.
// A tag not closed before limit ends at limit.
func (sc *Scanner) Next(limit int) (t Token, ok bool) {
	for {
		i := strings.IndexByte(sc.s[sc.pos:limit], '<')
		if i == -1 || sc.pos+i+1 >= limit {
			sc.pos = limit
			return Token{}, false
		}
		t = Token{Start: sc.pos + i}
		sc.pos = t.Start + 1

		switch c := sc.s[sc.pos]; {
		case c == '!' || c == '?':
			rest := sc.s[sc.pos:limit]
			switch {
			case strings.HasPrefix(rest, "!--"):
				sc.pos += 3
				sc.skipPast(limit, "-->")
			case strings.HasPrefix(rest, "![CDATA["):
				sc.skipPast(limit, "]]>")
			default:
				sc.skipPast(limit, ">")
			}
			t.Type = Markup
			return t, true
		case c == '/':
			sc.pos++
			if sc.pos < limit && !isASCIILetter(sc.s[sc.pos]) {
				// A bogus comment, e.g. "</ x>".
				sc.skipPast(limit, ">")
				t.Type = Markup
				return t, true
			}
			t.Type = EndTag
		case isASCIILetter(c):
			t.Type = StartTag
		default:
			continue
		}

		start := sc.pos
		for sc.pos < limit && !isTagNameEnd(sc.s[sc.pos]) {
			sc.pos++
		}
		t.Name = strings.ToLower(sc.s[start:sc.pos])

		for sc.pos < limit {
			c := sc.s[sc.pos]
			sc.pos++
			switch c {
			case '>':
				return t, true
			case '/':
				t.SelfClosing = t.Type == StartTag && sc.pos < limit && sc.s[sc.pos] == '>'
			case '=':
				for sc.pos < limit && isSpace(sc.s[sc.pos]) {
					sc.pos++
				}
				if sc.pos < limit && (sc.s[sc.pos] == '"' || sc.s[sc.pos] == '\'') {
					q := sc.s[sc.pos]
					sc.pos++
					if j := strings.IndexByte(sc.s[sc.pos:limit], q); j == -1 {
						sc.pos = limit
					} else {
						sc.pos += j + 1
					}
				}
			}
		}
		return t, true
	}
}

// SkipRawText moves the position to the start of the end tag of the
// raw text element name (e.g. script), or to limit if there is none.
func (sc *Scanner) SkipRawText(limit int, name string) {
	for {
		sc.skipPast(limit, "</")
		if sc.pos >= limit {
			return
		}
		if rest := sc.s[sc.pos:limit]; len(rest) > len(name) && strings.EqualFold(rest[:len(name)], name) && isTagNameEnd(rest[len(name)]) {
			sc.pos -= 2
			return
		}
	}
}

func (sc *Scanner) skipPast(limit int, sep string) {
	if i := strings.Index(sc.s[sc.pos:limit], sep); i == -1 {
		sc.pos = limit
	} else {
		sc.pos += i + len(sep)
	}
}

// IsRawTextElement reports whether the element with the given lowercased name
// contains raw text, i.e. no tags.
func IsRawTextElement(name string) bool {
	switch name {
	case "script", "style", "textarea", "title", "iframe", "xmp", "noembed", "noframes":
		return true
	}
	return false
}

// IsVoidElement reports whether the element with the given lowercased name
// has no end tag.
func IsVoidElement(name string) bool {
	switch name {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "source", "track", "wbr":
		return true
	}
	return false
}

// IsBlockElement reports whether the element with the given lowercased name
// is rendered as a block (or table cell), separating its content from
// the surrounding text.
func IsBlockElement(name string) bool {
	switch name {
	case "address", "article", "aside", "blockquote", "caption", "dd", "details", "dialog", "div", "dl", "dt",
		"fieldset", "figcaption", "figure", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "header",
		"hgroup", "hr", "li", "main", "nav", "ol", "p", "pre", "section", "summary", "table", "tbody", "td",
		"tfoot", "th", "thead", "tr", "ul":
		return true
	}
	return false
}

func isASCIILetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

func isTagNameEnd(c byte) bool {
	return isSpace(c) || c == '/' || c == '>'
}
