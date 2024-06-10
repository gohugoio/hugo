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

	// Contains lexers for shortcodes and other main section
	// elements.
	sectionHandlers *sectionHandlers

	cfg Config

	// The summary divider to look for.
	summaryDivider []byte
	// Set when we have parsed any summary divider
	summaryDividerChecked bool

	lexerShortcodeState

	// items delivered to client
	items Items

	// error delivered to the client
	err error
}

// Implement the Result interface
func (l *pageLexer) Iterator() *Iterator {
	return NewIterator(l.items)
}

func (l *pageLexer) Input() []byte {
	return l.input
}

type Config struct {
	NoFrontMatter    bool
	NoSummaryDivider bool
}

// note: the input position here is normally 0 (start), but
// can be set if position of first shortcode is known
func newPageLexer(input []byte, stateStart stateFunc, cfg Config) *pageLexer {
	lexer := &pageLexer{
		input:          input,
		stateStart:     stateStart,
		summaryDivider: summaryDivider,
		cfg:            cfg,
		lexerShortcodeState: lexerShortcodeState{
			currLeftDelimItem:  tLeftDelimScNoMarkup,
			currRightDelimItem: tRightDelimScNoMarkup,
			openShortcodes:     make(map[string]bool),
		},
		items: make([]Item, 0, 5),
	}

	lexer.sectionHandlers = createSectionHandlers(lexer)

	return lexer
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
)

func (l *pageLexer) next() rune {
	if l.pos >= len(l.input) {
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

func (l *pageLexer) append(item Item) {
	if item.Pos() < len(l.input) {
		item.firstByte = l.input[item.Pos()]
	}
	l.items = append(l.items, item)
}

// sends an item back to the client.
func (l *pageLexer) emit(t ItemType) {
	defer func() {
		l.start = l.pos
	}()

	if t == tText {
		// Identify any trailing whitespace/intendation.
		// We currently only care about the last one.
		for i := l.pos - 1; i >= l.start; i-- {
			b := l.input[i]
			if b != ' ' && b != '\t' && b != '\r' && b != '\n' {
				break
			}
			if i == l.start && b != '\n' {
				l.append(Item{Type: tIndentation, low: l.start, high: l.pos})
				return
			} else if b == '\n' && i < l.pos-1 {
				l.append(Item{Type: t, low: l.start, high: i + 1})
				l.append(Item{Type: tIndentation, low: i + 1, high: l.pos})
				return
			} else if b == '\n' && i == l.pos-1 {
				break
			}

		}
	}

	l.append(Item{Type: t, low: l.start, high: l.pos})
}

// sends a string item back to the client.
func (l *pageLexer) emitString(t ItemType) {
	l.append(Item{Type: t, low: l.start, high: l.pos, isString: true})
	l.start = l.pos
}

func (l *pageLexer) isEOF() bool {
	return l.pos >= len(l.input)
}

// special case, do not send '\\' back to client
func (l *pageLexer) ignoreEscapesAndEmit(t ItemType, isString bool) {
	i := l.start
	k := i

	var segments []lowHigh

	for i < l.pos {
		r, w := utf8.DecodeRune(l.input[i:l.pos])
		if r == '\\' {
			if i > k {
				segments = append(segments, lowHigh{k, i})
			}
			// See issue #10236.
			// We don't send the backslash back to the client,
			// which makes the end parsing simpler.
			// This means that we cannot render the AST back to be
			// exactly the same as the input,
			// but that was also the situation before we introduced the issue in #10236.
			k = i + w
		}
		i += w
	}

	if k < l.pos {
		segments = append(segments, lowHigh{k, l.pos})
	}

	if len(segments) > 0 {
		l.append(Item{Type: t, segments: segments})
	}

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
func (l *pageLexer) errorf(format string, args ...any) stateFunc {
	l.append(Item{Type: tError, Err: fmt.Errorf(format, args...), low: l.start, high: l.pos})
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

func (l *pageLexer) consumeToSpace() {
	for {
		r := l.next()
		if r == eof || unicode.IsSpace(r) {
			l.backup()
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

type sectionHandlers struct {
	l *pageLexer

	// Set when none of the sections are found so we
	// can safely stop looking and skip to the end.
	skipAll bool

	handlers    []*sectionHandler
	skipIndexes []int
}

func (s *sectionHandlers) skip() int {
	if s.skipAll {
		return -1
	}

	s.skipIndexes = s.skipIndexes[:0]
	var shouldSkip bool
	for _, skipper := range s.handlers {
		idx := skipper.skip()
		if idx != -1 {
			shouldSkip = true
			s.skipIndexes = append(s.skipIndexes, idx)
		}
	}

	if !shouldSkip {
		s.skipAll = true
		return -1
	}

	return minIndex(s.skipIndexes...)
}

func createSectionHandlers(l *pageLexer) *sectionHandlers {
	handlers := make([]*sectionHandler, 0, 2)

	shortCodeHandler := &sectionHandler{
		l: l,
		skipFunc: func(l *pageLexer) int {
			return l.index(leftDelimSc)
		},
		lexFunc: func(origin stateFunc, l *pageLexer) (stateFunc, bool) {
			if !l.isShortCodeStart() {
				return origin, false
			}

			if l.isInline {
				// If we're inside an inline shortcode, the only valid shortcode markup is
				// the markup which closes it.
				b := l.input[l.pos+3:]
				end := indexNonWhiteSpace(b, '/')
				if end != len(l.input)-1 {
					b = bytes.TrimSpace(b[end+1:])
					if end == -1 || !bytes.HasPrefix(b, []byte(l.currShortcodeName+" ")) {
						return l.errorf("inline shortcodes do not support nesting"), true
					}
				}
			}

			if l.hasPrefix(leftDelimScWithMarkup) {
				l.currLeftDelimItem = tLeftDelimScWithMarkup
				l.currRightDelimItem = tRightDelimScWithMarkup
			} else {
				l.currLeftDelimItem = tLeftDelimScNoMarkup
				l.currRightDelimItem = tRightDelimScNoMarkup
			}

			return lexShortcodeLeftDelim, true
		},
	}

	handlers = append(handlers, shortCodeHandler)

	if !l.cfg.NoSummaryDivider {
		summaryDividerHandler := &sectionHandler{
			l: l,
			skipFunc: func(l *pageLexer) int {
				if l.summaryDividerChecked {
					return -1
				}
				return l.index(l.summaryDivider)
			},
			lexFunc: func(origin stateFunc, l *pageLexer) (stateFunc, bool) {
				if !l.hasPrefix(l.summaryDivider) {
					return origin, false
				}

				l.summaryDividerChecked = true
				l.pos += len(l.summaryDivider)
				// This makes it a little easier to reason about later.
				l.consumeSpace()
				l.emit(TypeLeadSummaryDivider)

				return origin, true
			},
		}

		handlers = append(handlers, summaryDividerHandler)

	}

	return &sectionHandlers{
		l:           l,
		handlers:    handlers,
		skipIndexes: make([]int, len(handlers)),
	}
}

func (s *sectionHandlers) lex(origin stateFunc) stateFunc {
	if s.skipAll {
		return nil
	}

	if s.l.pos > s.l.start {
		s.l.emit(tText)
	}

	for _, handler := range s.handlers {
		if handler.skipAll {
			continue
		}

		next, handled := handler.lexFunc(origin, handler.l)
		if next == nil || handled {
			return next
		}
	}

	// Not handled by the above.
	s.l.pos++

	return origin
}

type sectionHandler struct {
	l *pageLexer

	// No more sections of this type.
	skipAll bool

	// Returns the index of the next match, -1 if none found.
	skipFunc func(l *pageLexer) int

	// Lex lexes the current section and returns the next state func and
	// a bool telling if this section was handled.
	// Note that returning nil as the next state will terminate the
	// lexer.
	lexFunc func(origin stateFunc, l *pageLexer) (stateFunc, bool)
}

func (s *sectionHandler) skip() int {
	if s.skipAll {
		return -1
	}

	idx := s.skipFunc(s.l)
	if idx == -1 {
		s.skipAll = true
	}
	return idx
}

func lexMainSection(l *pageLexer) stateFunc {
	if l.isEOF() {
		return lexDone
	}

	// Fast forward as far as possible.
	skip := l.sectionHandlers.skip()

	if skip == -1 {
		l.pos = len(l.input)
		return lexDone
	} else if skip > 0 {
		l.pos += skip
	}

	next := l.sectionHandlers.lex(lexMainSection)
	if next != nil {
		return next
	}

	l.pos = len(l.input)
	return lexDone
}

func lexDone(l *pageLexer) stateFunc {
	// Done!
	if l.pos > l.start {
		l.emit(tText)
	}
	l.emit(tEOF)
	return nil
}

//lint:ignore U1000 useful for debugging
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
