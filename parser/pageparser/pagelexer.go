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

// position (in bytes)
type pos int

const eof = -1

// returns the next state in scanner.
type stateFunc func(*pageLexer) stateFunc

type lexerShortcodeState struct {
	currLeftDelimItem  itemType
	currRightDelimItem itemType
	currShortcodeName  string          // is only set when a shortcode is in opened state
	closingState       int             // > 0 = on its way to be closed
	elementStepNum     int             // step number in element
	paramElements      int             // number of elements (name + value = 2) found first
	openShortcodes     map[string]bool // set of shortcodes in open state

}

type pageLexer struct {
	input      []byte
	stateStart stateFunc
	state      stateFunc
	pos        pos // input position
	start      pos // item start position
	width      pos // width of last element
	lastPos    pos // position of the last item returned by nextItem

	contentSections int

	lexerShortcodeState

	// items delivered to client
	items []Item
}

func Parse(s string) *Tokens {
	return ParseFrom(s, 0)
}

func ParseFrom(s string, from int) *Tokens {
	input := []byte(s)
	lexer := newPageLexer(input, pos(from), lexMainSection) // TODO(bep) 2errors
	lexer.run()
	return &Tokens{lexer: lexer}
}

// note: the input position here is normally 0 (start), but
// can be set if position of first shortcode is known
// TODO(bep) 2errors byte
func newPageLexer(input []byte, inputPosition pos, stateStart stateFunc) *pageLexer {
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

// main loop
func (l *pageLexer) run() *pageLexer {
	for l.state = l.stateStart; l.state != nil; {
		l.state = l.state(l)
	}
	return l
}

// Shortcode syntax
var (
	leftDelimScNoMarkup    = []byte("{{<")
	rightDelimScNoMarkup   = []byte(">}}")
	leftDelimScWithMarkup  = []byte("{{%")
	rightDelimScWithMarkup = []byte("%}}")
	leftComment            = []byte("/*") // comments in this context us used to to mark shortcodes as "not really a shortcode"
	rightComment           = []byte("*/")
)

// Page syntax
var (
	summaryDivider    = []byte("<!--more-->")
	summaryDividerOrg = []byte("# more")
	delimTOML         = []byte("+++")
	delimYAML         = []byte("---")
	delimOrg          = []byte("#+")
)

func (l *pageLexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}

	runeValue, runeWidth := utf8.DecodeRune(l.input[l.pos:])
	l.width = pos(runeWidth)
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
func (l *pageLexer) emit(t itemType) {
	l.items = append(l.items, Item{t, l.start, l.input[l.start:l.pos]})
	l.start = l.pos
}

// special case, do not send '\\' back to client
func (l *pageLexer) ignoreEscapesAndEmit(t itemType) {
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

// nice to have in error logs
func (l *pageLexer) lineNum() int {
	return bytes.Count(l.input[:l.lastPos], lf) + 1
}

// nil terminates the parser
func (l *pageLexer) errorf(format string, args ...interface{}) stateFunc {
	l.items = append(l.items, Item{tError, l.start, []byte(fmt.Sprintf(format, args...))})
	return nil
}

// consumes and returns the next item
func (l *pageLexer) nextItem() Item {
	item := l.items[0]
	l.items = l.items[1:]
	l.lastPos = item.pos
	return item
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

func lexMainSection(l *pageLexer) stateFunc {
	for {
		if l.isShortCodeStart() {
			if l.pos > l.start {
				l.emit(tText)
			}
			if bytes.HasPrefix(l.input[l.pos:], leftDelimScWithMarkup) {
				l.currLeftDelimItem = tLeftDelimScWithMarkup
				l.currRightDelimItem = tRightDelimScWithMarkup
			} else {
				l.currLeftDelimItem = tLeftDelimScNoMarkup
				l.currRightDelimItem = tRightDelimScNoMarkup
			}
			return lexShortcodeLeftDelim
		}

		if l.contentSections <= 1 {
			if bytes.HasPrefix(l.input[l.pos:], summaryDivider) {
				if l.pos > l.start {
					l.emit(tText)
				}
				l.contentSections++
				l.pos += pos(len(summaryDivider))
				l.emit(tSummaryDivider)
			} else if bytes.HasPrefix(l.input[l.pos:], summaryDividerOrg) {
				if l.pos > l.start {
					l.emit(tText)
				}
				l.contentSections++
				l.pos += pos(len(summaryDividerOrg))
				l.emit(tSummaryDividerOrg)
			}
		}

		r := l.next()
		if r == eof {
			break
		}

	}

	return lexDone

}

func (l *pageLexer) isShortCodeStart() bool {
	return bytes.HasPrefix(l.input[l.pos:], leftDelimScWithMarkup) || bytes.HasPrefix(l.input[l.pos:], leftDelimScNoMarkup)
}

func lexIntroSection(l *pageLexer) stateFunc {
LOOP:
	for {
		r := l.next()
		if r == eof {
			break
		}

		switch {
		case r == '+':
			return l.lexFrontMatterSection(tFrontMatterTOML, r, "TOML", delimTOML)
		case r == '-':
			return l.lexFrontMatterSection(tFrontMatterYAML, r, "YAML", delimYAML)
		case r == '{':
			return lexFrontMatterJSON
		case r == '#':
			return lexFrontMatterOrgMode
		case !isSpace(r) && !isEndOfLine(r):
			if r == '<' {
				l.emit(tHTMLLead)
				// Not need to look further. Hugo treats this as plain HTML,
				// no front matter, no shortcodes, no nothing.
				l.pos = pos(len(l.input))
				l.emit(tText)
				break LOOP
			}
			return l.errorf("failed to detect front matter type; got unknown identifier %q", r)
		}
	}

	l.contentSections = 1

	// Now move on to the shortcodes.
	return lexMainSection
}

func lexDone(l *pageLexer) stateFunc {

	// Done!
	if l.pos > l.start {
		l.emit(tText)
	}
	l.emit(tEOF)
	return nil
}

func lexFrontMatterJSON(l *pageLexer) stateFunc {
	// Include the left delimiter
	l.backup()

	var (
		inQuote bool
		level   int
	)

	for {

		r := l.next()

		switch {
		case r == eof:
			return l.errorf("unexpected EOF parsing JSON front matter")
		case r == '{':
			if !inQuote {
				level++
			}
		case r == '}':
			if !inQuote {
				level--
			}
		case r == '"':
			inQuote = !inQuote
		case r == '\\':
			// This may be an escaped quote. Make sure it's not marked as a
			// real one.
			l.next()
		}

		if level == 0 {
			break
		}
	}

	l.consumeCRLF()
	l.emit(tFrontMatterJSON)

	return lexMainSection
}

func lexFrontMatterOrgMode(l *pageLexer) stateFunc {
	/*
		#+TITLE: Test File For chaseadamsio/goorgeous
		#+AUTHOR: Chase Adams
		#+DESCRIPTION: Just another golang parser for org content!
	*/

	l.backup()

	if !bytes.HasPrefix(l.input[l.pos:], delimOrg) {
		// TODO(bep) consider error
		return lexMainSection
	}

	// Read lines until we no longer see a #+ prefix
LOOP:
	for {

		r := l.next()

		switch {
		case r == '\n':
			if !bytes.HasPrefix(l.input[l.pos:], delimOrg) {
				break LOOP
			}
		case r == eof:
			break LOOP

		}
	}

	l.emit(tFrontMatterORG)

	return lexMainSection

}

// Handle YAML or TOML front matter.
func (l *pageLexer) lexFrontMatterSection(tp itemType, delimr rune, name string, delim []byte) stateFunc {
	for i := 0; i < 2; i++ {
		if r := l.next(); r != delimr {
			return l.errorf("invalid %s delimiter", name)
		}
	}

	if !l.consumeCRLF() {
		return l.errorf("invalid %s delimiter", name)
	}

	// We don't care about the delimiters.
	l.ignore()

	for {
		r := l.next()
		if r == eof {
			return l.errorf("EOF looking for end %s front matter delimiter", name)
		}
		if isEndOfLine(r) {
			if bytes.HasPrefix(l.input[l.pos:], delim) {
				l.emit(tp)
				l.pos += 3
				l.consumeCRLF()
				l.ignore()
				break
			}
		}
	}

	return lexMainSection
}

func lexShortcodeLeftDelim(l *pageLexer) stateFunc {
	l.pos += pos(len(l.currentLeftShortcodeDelim()))
	if bytes.HasPrefix(l.input[l.pos:], leftComment) {
		return lexShortcodeComment
	}
	l.emit(l.currentLeftShortcodeDelimItem())
	l.elementStepNum = 0
	l.paramElements = 0
	return lexInsideShortcode
}

func lexShortcodeComment(l *pageLexer) stateFunc {
	posRightComment := bytes.Index(l.input[l.pos:], append(rightComment, l.currentRightShortcodeDelim()...))
	if posRightComment <= 1 {
		return l.errorf("comment must be closed")
	}
	// we emit all as text, except the comment markers
	l.emit(tText)
	l.pos += pos(len(leftComment))
	l.ignore()
	l.pos += pos(posRightComment - len(leftComment))
	l.emit(tText)
	l.pos += pos(len(rightComment))
	l.ignore()
	l.pos += pos(len(l.currentRightShortcodeDelim()))
	l.emit(tText)
	return lexMainSection
}

func lexShortcodeRightDelim(l *pageLexer) stateFunc {
	l.closingState = 0
	l.pos += pos(len(l.currentRightShortcodeDelim()))
	l.emit(l.currentRightShortcodeDelimItem())
	return lexMainSection
}

// either:
// 1. param
// 2. "param" or "param\"
// 3. param="123" or param="123\"
// 4. param="Some \"escaped\" text"
func lexShortcodeParam(l *pageLexer, escapedQuoteStart bool) stateFunc {

	first := true
	nextEq := false

	var r rune

	for {
		r = l.next()
		if first {
			if r == '"' {
				// a positional param with quotes
				if l.paramElements == 2 {
					return l.errorf("got quoted positional parameter. Cannot mix named and positional parameters")
				}
				l.paramElements = 1
				l.backup()
				return lexShortcodeQuotedParamVal(l, !escapedQuoteStart, tScParam)
			}
			first = false
		} else if r == '=' {
			// a named param
			l.backup()
			nextEq = true
			break
		}

		if !isAlphaNumericOrHyphen(r) {
			l.backup()
			break
		}
	}

	if l.paramElements == 0 {
		l.paramElements++

		if nextEq {
			l.paramElements++
		}
	} else {
		if nextEq && l.paramElements == 1 {
			return l.errorf("got named parameter '%s'. Cannot mix named and positional parameters", l.current())
		} else if !nextEq && l.paramElements == 2 {
			return l.errorf("got positional parameter '%s'. Cannot mix named and positional parameters", l.current())
		}
	}

	l.emit(tScParam)
	return lexInsideShortcode

}

func lexShortcodeQuotedParamVal(l *pageLexer, escapedQuotedValuesAllowed bool, typ itemType) stateFunc {
	openQuoteFound := false
	escapedInnerQuoteFound := false
	escapedQuoteState := 0

Loop:
	for {
		switch r := l.next(); {
		case r == '\\':
			if l.peek() == '"' {
				if openQuoteFound && !escapedQuotedValuesAllowed {
					l.backup()
					break Loop
				} else if openQuoteFound {
					// the coming quoute is inside
					escapedInnerQuoteFound = true
					escapedQuoteState = 1
				}
			}
		case r == eof, r == '\n':
			return l.errorf("unterminated quoted string in shortcode parameter-argument: '%s'", l.current())
		case r == '"':
			if escapedQuoteState == 0 {
				if openQuoteFound {
					l.backup()
					break Loop

				} else {
					openQuoteFound = true
					l.ignore()
				}
			} else {
				escapedQuoteState = 0
			}

		}
	}

	if escapedInnerQuoteFound {
		l.ignoreEscapesAndEmit(typ)
	} else {
		l.emit(typ)
	}

	r := l.next()

	if r == '\\' {
		if l.peek() == '"' {
			// ignore the escaped closing quote
			l.ignore()
			l.next()
			l.ignore()
		}
	} else if r == '"' {
		// ignore closing quote
		l.ignore()
	} else {
		// handled by next state
		l.backup()
	}

	return lexInsideShortcode
}

// scans an alphanumeric inside shortcode
func lexIdentifierInShortcode(l *pageLexer) stateFunc {
	lookForEnd := false
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumericOrHyphen(r):
		// Allow forward slash inside names to make it possible to create namespaces.
		case r == '/':
		default:
			l.backup()
			word := string(l.input[l.start:l.pos])
			if l.closingState > 0 && !l.openShortcodes[word] {
				return l.errorf("closing tag for shortcode '%s' does not match start tag", word)
			} else if l.closingState > 0 {
				l.openShortcodes[word] = false
				lookForEnd = true
			}

			l.closingState = 0
			l.currShortcodeName = word
			l.openShortcodes[word] = true
			l.elementStepNum++
			l.emit(tScName)
			break Loop
		}
	}

	if lookForEnd {
		return lexEndOfShortcode
	}
	return lexInsideShortcode
}

func lexEndOfShortcode(l *pageLexer) stateFunc {
	if bytes.HasPrefix(l.input[l.pos:], l.currentRightShortcodeDelim()) {
		return lexShortcodeRightDelim
	}
	switch r := l.next(); {
	case isSpace(r):
		l.ignore()
	default:
		return l.errorf("unclosed shortcode")
	}
	return lexEndOfShortcode
}

// scans the elements inside shortcode tags
func lexInsideShortcode(l *pageLexer) stateFunc {
	if bytes.HasPrefix(l.input[l.pos:], l.currentRightShortcodeDelim()) {
		return lexShortcodeRightDelim
	}
	switch r := l.next(); {
	case r == eof:
		// eol is allowed inside shortcodes; this may go to end of document before it fails
		return l.errorf("unclosed shortcode action")
	case isSpace(r), isEndOfLine(r):
		l.ignore()
	case r == '=':
		l.ignore()
		return lexShortcodeQuotedParamVal(l, l.peek() != '\\', tScParamVal)
	case r == '/':
		if l.currShortcodeName == "" {
			return l.errorf("got closing shortcode, but none is open")
		}
		l.closingState++
		l.emit(tScClose)
	case r == '\\':
		l.ignore()
		if l.peek() == '"' {
			return lexShortcodeParam(l, true)
		}
	case l.elementStepNum > 0 && (isAlphaNumericOrHyphen(r) || r == '"'): // positional params can have quotes
		l.backup()
		return lexShortcodeParam(l, false)
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifierInShortcode
	default:
		return l.errorf("unrecognized character in shortcode action: %#U. Note: Parameters with non-alphanumeric args must be quoted", r)
	}
	return lexInsideShortcode
}

// state helpers

func (l *pageLexer) currentLeftShortcodeDelimItem() itemType {
	return l.currLeftDelimItem
}

func (l *pageLexer) currentRightShortcodeDelimItem() itemType {
	return l.currRightDelimItem
}

func (l *pageLexer) currentLeftShortcodeDelim() []byte {
	if l.currLeftDelimItem == tLeftDelimScWithMarkup {
		return leftDelimScWithMarkup
	}
	return leftDelimScNoMarkup

}

func (l *pageLexer) currentRightShortcodeDelim() []byte {
	if l.currRightDelimItem == tRightDelimScWithMarkup {
		return rightDelimScWithMarkup
	}
	return rightDelimScNoMarkup
}

// helper functions

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
