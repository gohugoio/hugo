// Copyright 2018 The Hugo Authors. All rights reserved.
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

package urlreplacers

import (
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/transform"
)

type absurllexer struct {
	// the source to absurlify
	content []byte
	// the target for the new absurlified content
	w io.Writer

	// path may be set to a "." relative path
	path []byte

	pos   int // input position
	start int // item start position

	quotes [][]byte
}

type prefix struct {
	disabled bool
	b        []byte
	f        func(l *absurllexer)

	nextPos int
}

func (p *prefix) find(bs []byte, start int) bool {
	if p.disabled {
		return false
	}

	if p.nextPos == -1 {
		idx := bytes.Index(bs[start:], p.b)

		if idx == -1 {
			p.disabled = true
			// Find the closest match
			return false
		}

		p.nextPos = start + idx + len(p.b)
	}

	return true
}

func newPrefixState() []*prefix {
	return []*prefix{
		{b: []byte("src="), f: checkCandidateBase},
		{b: []byte("href="), f: checkCandidateBase},
		{b: []byte("action="), f: checkCandidateBase},
		{b: []byte("srcset="), f: checkCandidateSrcset},
	}
}

func (l *absurllexer) emit() {
	l.w.Write(l.content[l.start:l.pos])
	l.start = l.pos
}

var (
	relURLPrefix    = []byte("/")
	relURLPrefixLen = len(relURLPrefix)
)

func (l *absurllexer) consumeQuote() []byte {
	for _, q := range l.quotes {
		if bytes.HasPrefix(l.content[l.pos:], q) {
			l.pos += len(q)
			l.emit()
			return q
		}
	}
	return nil
}

// handle URLs in src and href.
func checkCandidateBase(l *absurllexer) {
	l.consumeQuote()

	if !bytes.HasPrefix(l.content[l.pos:], relURLPrefix) {
		return
	}

	// check for schemaless URLs
	posAfter := l.pos + relURLPrefixLen
	if posAfter >= len(l.content) {
		return
	}
	r, _ := utf8.DecodeRune(l.content[posAfter:])
	if r == '/' {
		// schemaless: skip
		return
	}
	if l.pos > l.start {
		l.emit()
	}
	l.pos += relURLPrefixLen
	l.w.Write(l.path)
	l.start = l.pos
}

func (l *absurllexer) posAfterURL(q []byte) int {
	if len(q) > 0 {
		// look for end quote
		return bytes.Index(l.content[l.pos:], q)
	}

	return bytes.IndexFunc(l.content[l.pos:], func(r rune) bool {
		return r == '>' || unicode.IsSpace(r)
	})

}

// handle URLs in srcset.
func checkCandidateSrcset(l *absurllexer) {
	q := l.consumeQuote()
	if q == nil {
		// srcset needs to be quoted.
		return
	}

	// special case, not frequent (me think)
	if !bytes.HasPrefix(l.content[l.pos:], relURLPrefix) {
		return
	}

	// check for schemaless URLs
	posAfter := l.pos + relURLPrefixLen
	if posAfter >= len(l.content) {
		return
	}
	r, _ := utf8.DecodeRune(l.content[posAfter:])
	if r == '/' {
		// schemaless: skip
		return
	}

	posEnd := l.posAfterURL(q)

	// safe guard
	if posEnd < 0 || posEnd > 2000 {
		return
	}

	if l.pos > l.start {
		l.emit()
	}

	section := l.content[l.pos : l.pos+posEnd+1]

	fields := bytes.Fields(section)
	for i, f := range fields {
		if f[0] == '/' {
			l.w.Write(l.path)
			l.w.Write(f[1:])

		} else {
			l.w.Write(f)
		}

		if i < len(fields)-1 {
			l.w.Write([]byte(" "))
		}
	}

	l.pos += len(section)
	l.start = l.pos

}

// main loop
func (l *absurllexer) replace() {
	contentLength := len(l.content)

	prefixes := newPrefixState()

	for {
		if l.pos >= contentLength {
			break
		}

		var match *prefix

		for _, p := range prefixes {
			if !p.find(l.content, l.pos) {
				continue
			}

			if match == nil || p.nextPos < match.nextPos {
				match = p
			}
		}

		if match == nil {
			// Done!
			l.pos = contentLength
			break
		} else {
			l.pos = match.nextPos
			match.nextPos = -1
			match.f(l)
		}
	}
	// Done!
	if l.pos > l.start {
		l.emit()
	}
}

func doReplace(path string, ct transform.FromTo, quotes [][]byte) {

	lexer := &absurllexer{
		content: ct.From().Bytes(),
		w:       ct.To(),
		path:    []byte(path),
		quotes:  quotes}

	lexer.replace()
}

type absURLReplacer struct {
	htmlQuotes [][]byte
	xmlQuotes  [][]byte
}

func newAbsURLReplacer() *absURLReplacer {
	return &absURLReplacer{
		htmlQuotes: [][]byte{[]byte("\""), []byte("'")},
		xmlQuotes:  [][]byte{[]byte("&#34;"), []byte("&#39;")}}
}

func (au *absURLReplacer) replaceInHTML(path string, ct transform.FromTo) {
	doReplace(path, ct, au.htmlQuotes)
}

func (au *absURLReplacer) replaceInXML(path string, ct transform.FromTo) {
	doReplace(path, ct, au.xmlQuotes)
}
