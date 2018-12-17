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

// Package pageparser provides a parser for Hugo content files (Markdown, HTML etc.) in Hugo.
// This implementation is highly inspired by the great talk given by Rob Pike called "Lexical Scanning in Go"
// It's on YouTube, Google it!.
// See slides here: http://cuddle.googlecode.com/hg/talk/lex.html
package pageparser

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

const eof = -1

// returns the next state in scanner.
type stateFunc func(*pageLexer) stateFunc

type pageLexer struct {
	input      []byte
	stateStart stateFunc
	state      stateFunc
	pos        int // input position
	start      int // item start position
	width      int // width of last element

	// The summary divider to look for.
	summaryDivider []byte
	// Set when we have parsed any summary divider
	summaryDividerChecked bool
	// Whether we're in a HTML comment.
	isInHTMLComment bool

	lexerShortcodeState

	// items delivered to client
	items Items
}

// Implement the Result interface
func (l *pageLexer) Iterator() *Iterator {
	return l.newIterator()
}

func (l *pageLexer) Input() []byte {
	return l.input

}

// note: the input position here is normally 0 (start), but
// can be set if position of first shortcode is known
func newPageLexer(input []byte, inputPosition int, stateStart stateFunc) *pageLexer {
	lexer := &pageLexer{
		input:      input,
		pos:        inputPosition,
		stateStart: stateStart,
		lexerShortcodeState: lexerShortcodeState{
			currLeftDelimItem:  tLeftDelimScNoMarkup,
			currRightDelimItem: tRightDelimScNoMarkup,
			openShortcodes:     make(map[string]bool),
		},
		items: make([]Item, 0, 5),
	}

	return lexer
}

func (l *pageLexer) newIterator() *Iterator {
	return &Iterator{l: l, lastPos: -1}
}

// main loop
func (l *pageLexer) run() *pageLexer {
	for l.state = l.stateStart; l.state != nil; {
		l.state = l.state(l)
	}
	return l
}

// Page syntax
var (
	byteOrderMark     = '\ufeff'
	summaryDivider    = []byte("<!--more-->")
	summaryDividerOrg = []byte("# more")
	delimTOML         = []byte("+++")
	delimYAML         = []byte("---")
	delimOrg          = []byte("#+")
	htmlCommentStart  = []byte("<!--")
	htmlCommentEnd    = []byte("-->")
)

func (l *pageLexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}

	runeValue, runeWidth := utf8.DecodeRune(l.input[l.pos:])
	l.width = runeWidth
	l.pos += l.width
	return runeValue
}

// peek, but no consume
func (l *pageLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// steps back one
func (l *pageLexer) backup() {
	l.pos -= l.width
}

// sends an item back to the client.
func (l *pageLexer) emit(t ItemType) {
	l.items = append(l.items, Item{t, l.start, l.input[l.start:l.pos]})
	l.start = l.pos
}

// special case, do not send '\\' back to client
func (l *pageLexer) ignoreEscapesAndEmit(t ItemType) {
	val := bytes.Map(func(r rune) rune {
		if r == '\\' {
			return -1
		}
		return r
	}, l.input[l.start:l.pos])
	l.items = append(l.items, Item{t, l.start, val})
	l.start = l.pos
}

// gets the current value (for debugging and error handling)
func (l *pageLexer) current() []byte {
	return l.input[l.start:l.pos]
}

// ignore current element
func (l *pageLexer) ignore() {
	l.start = l.pos
}

var lf = []byte("\n")

// nil terminates the parser
func (l *pageLexer) errorf(format string, args ...interface{}) stateFunc {
	l.items = append(l.items, Item{tError, l.start, []byte(fmt.Sprintf(format, args...))})
	return nil
}

func (l *pageLexer) consumeCRLF() bool {
	var consumed bool
	for _, r := range crLf {
		if l.next() != r {
			l.backup()
		} else {
			consumed = true
		}
	}
	return consumed
}

func (l *pageLexer) consumeToNextLine() {
	for {
		r := l.next()
		if r == eof || isEndOfLine(r) {
			return
		}
	}
}

func (l *pageLexer) consumeSpace() {
	for {
		r := l.next()
		if r == eof || !unicode.IsSpace(r) {
			l.backup()
			return
		}
	}
}

func lexMainSection(l *pageLexer) stateFunc {
	if l.isInHTMLComment {
		return lexEndFromtMatterHTMLComment
	}

	// Fast forward as far as possible.
	var l1, l2 int

	if !l.summaryDividerChecked && l.summaryDivider != nil {
		l1 = l.index(l.summaryDivider)
		if l1 == -1 {
			l.summaryDividerChecked = true
		}
	}

	l2 = l.index(leftDelimSc)
	skip := minIndex(l1, l2)

	if skip > 0 {
		l.pos += skip
	}

	for {
		if l.isShortCodeStart() {
			if l.isInline {
				// If we're inside an inline shortcode, the only valid shortcode markup is
				// the markup which closes it.
				b := l.input[l.pos+3:]
				end := indexNonWhiteSpace(b, '/')
				if end != len(l.input)-1 {
					b = bytes.TrimSpace(b[end+1:])
					if end == -1 || !bytes.HasPrefix(b, []byte(l.currShortcodeName+" ")) {
						return l.errorf("inline shortcodes do not support nesting")
					}
				}
			}

			if l.pos > l.start {
				l.emit(tText)
			}
			if l.hasPrefix(leftDelimScWithMarkup) {
				l.currLeftDelimItem = tLeftDelimScWithMarkup
				l.currRightDelimItem = tRightDelimScWithMarkup
			} else {
				l.currLeftDelimItem = tLeftDelimScNoMarkup
				l.currRightDelimItem = tRightDelimScNoMarkup
			}
			return lexShortcodeLeftDelim
		}

		if !l.summaryDividerChecked && l.summaryDivider != nil {
			if l.hasPrefix(l.summaryDivider) {
				if l.pos > l.start {
					l.emit(tText)
				}
				l.summaryDividerChecked = true
				l.pos += len(l.summaryDivider)
				// This makes it a little easier to reason about later.
				l.consumeSpace()
				l.emit(TypeLeadSummaryDivider)

				// We have already moved to the next.
				continue
			}
		}

		r := l.next()
		if r == eof {
			break
		}

	}

	return lexDone

}

func (l *pageLexer) posFirstNonWhiteSpace() int {
	f := func(c rune) bool {
		return !unicode.IsSpace(c)
	}
	return bytes.IndexFunc(l.input[l.pos:], f)
}

func lexDone(l *pageLexer) stateFunc {

	// Done!
	if l.pos > l.start {
		l.emit(tText)
	}
	l.emit(tEOF)
	return nil
}

func (l *pageLexer) printCurrentInput() {
	fmt.Printf("input[%d:]: %q", l.pos, string(l.input[l.pos:]))
}

// state helpers

func (l *pageLexer) index(sep []byte) int {
	return bytes.Index(l.input[l.pos:], sep)
}

func (l *pageLexer) hasPrefix(prefix []byte) bool {
	return bytes.HasPrefix(l.input[l.pos:], prefix)
}

// helper functions

// returns the min index >= 0
func minIndex(indices ...int) int {
	min := -1

	for _, j := range indices {
		if j < 0 {
			continue
		}
		if min == -1 {
			min = j
		} else if j < min {
			min = j
		}
	}
	return min
}

func indexNonWhiteSpace(s []byte, in rune) int {
	idx := bytes.IndexFunc(s, func(r rune) bool {
		return !unicode.IsSpace(r)
	})

	if idx == -1 {
		return -1
	}

	r, _ := utf8.DecodeRune(s[idx:])
	if r == in {
		return idx
	}
	return -1
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isAlphaNumericOrHyphen(r rune) bool {
	// let unquoted YouTube ids as positional params slip through (they contain hyphens)
	return isAlphaNumeric(r) || r == '-'
}

var crLf = []rune{'\r', '\n'}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
