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

type matchState int

const (
	matchStateNone matchState = iota
	matchStateWhitespace
	matchStatePartial
	matchStateFull
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
	width int // width of last element

	quotes [][]byte

	ms      matchState
	matches [3]bool // track matches of the 3 prefixes
	idx     int     // last index in matches checked

}

type stateFunc func(*absurllexer) stateFunc

// prefix is how to identify and which func to handle the replacement.
type prefix struct {
	r []rune
	f func(l *absurllexer)
}

// new prefixes can be added below, but note:
// - the matches array above must be expanded.
// - the prefix must with the current logic end with '='
var prefixes = []*prefix{
	{r: []rune{'s', 'r', 'c', '='}, f: checkCandidateBase},
	{r: []rune{'h', 'r', 'e', 'f', '='}, f: checkCandidateBase},
	{r: []rune{'s', 'r', 'c', 's', 'e', 't', '='}, f: checkCandidateSrcset},
}

type absURLMatcher struct {
	match []byte
	quote []byte
}

// match check rune inside word. Will be != ' '.
func (l *absurllexer) match(r rune) {

	var found bool

	// note, the prefixes can start off on the same foot, i.e.
	// src and srcset.
	if l.ms == matchStateWhitespace {
		l.idx = 0
		for j, p := range prefixes {
			if r == p.r[l.idx] {
				l.matches[j] = true
				found = true
				// checkMatchState will only return true when r=='=', so
				// we can safely ignore the return value here.
				l.checkMatchState(r, j)
			}
		}

		if !found {
			l.ms = matchStateNone
		}

		return
	}

	l.idx++
	for j, m := range l.matches {
		// still a match?
		if m {
			if prefixes[j].r[l.idx] == r {
				found = true
				if l.checkMatchState(r, j) {
					return
				}
			} else {
				l.matches[j] = false
			}
		}
	}

	if !found {
		l.ms = matchStateNone
	}
}

func (l *absurllexer) checkMatchState(r rune, idx int) bool {
	if r == '=' {
		l.ms = matchStateFull
		for k := range l.matches {
			if k != idx {
				l.matches[k] = false
			}
		}
		return true
	}

	l.ms = matchStatePartial

	return false
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
	var r rune

	for {
		if l.pos >= contentLength {
			l.width = 0
			break
		}

		var width = 1
		r = rune(l.content[l.pos])
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRune(l.content[l.pos:])
		}
		l.width = width
		l.pos += l.width
		if r == ' ' {
			l.ms = matchStateWhitespace
		} else if l.ms != matchStateNone {
			l.match(r)
			if l.ms == matchStateFull {
				var p *prefix
				for i, m := range l.matches {
					if m {
						p = prefixes[i]
						l.matches[i] = false
					}
				}
				l.ms = matchStateNone
				p.f(l)
			}
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
