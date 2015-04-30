// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// The lexical scanning below is highly inspired by the great talk given by
// Rob Pike called "Lexical Scanning in Go" (it's on YouTube, Google it!).
// See slides here: http://cuddle.googlecode.com/hg/talk/lex.html

// parsing

type pageTokens struct {
	lexer     *pagelexer
	token     [3]item // 3-item look-ahead is what we currently need
	peekCount int
}

func (t *pageTokens) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lexer.nextItem()
	}
	return t.token[t.peekCount]
}

// backs up one token.
func (t *pageTokens) backup() {
	t.peekCount++
}

// backs up two tokens.
func (t *pageTokens) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}

// backs up three tokens.
func (t *pageTokens) backup3(t2, t1 item) {
	t.token[1] = t1
	t.token[2] = t2
	t.peekCount = 3
}

// check for non-error and non-EOF types coming next
func (t *pageTokens) isValueNext() bool {
	i := t.peek()
	return i.typ != tError && i.typ != tEOF
}

// look at, but do not consume, the next item
// repeated, sequential calls will return the same item
func (t *pageTokens) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lexer.nextItem()
	return t.token[0]
}

// convencience method to consume the next n tokens, but back off Errors and EOF
func (t *pageTokens) consume(cnt int) {
	for i := 0; i < cnt; i++ {
		token := t.next()
		if token.typ == tError || token.typ == tEOF {
			t.backup()
			break
		}
	}
}

// lexical scanning

// position (in bytes)
type pos int

type item struct {
	typ itemType
	pos pos
	val string
}

func (i item) String() string {
	switch {
	case i.typ == tEOF:
		return "EOF"
	case i.typ == tError:
		return i.val
	case i.typ > tKeywordMarker:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 20:
		return fmt.Sprintf("%.20q...", i.val)
	}
	return fmt.Sprintf("[%s]", i.val)
}

type itemType int

// named params in shortcodes
type namedParam struct {
	name  string
	value string
}

// for testing
func (np namedParam) String() string {
	return fmt.Sprintf("%s=%s", np.name, np.value)
}

const (
	tError itemType = iota
	tEOF

	// shortcode items
	tLeftDelimScNoMarkup
	tRightDelimScNoMarkup
	tLeftDelimScWithMarkup
	tRightDelimScWithMarkup
	tScClose
	tScName
	tScParam
	tScParamVal

	//itemIdentifier
	tText // plain text, used for everything outside the shortcodes

	// preserved for later - keywords come after this
	tKeywordMarker
)

const eof = -1

// returns the next state in scanner.
type stateFunc func(*pagelexer) stateFunc

type pagelexer struct {
	name    string
	input   string
	state   stateFunc
	pos     pos // input position
	start   pos // item start position
	width   pos // width of last element
	lastPos pos // position of the last item returned by nextItem

	// shortcode state
	currLeftDelimItem  itemType
	currRightDelimItem itemType
	currShortcodeName  string          // is only set when a shortcode is in opened state
	closingState       int             // > 0 = on it's way to be closed
	elementStepNum     int             // step number in element
	paramElements      int             // number of elements (name + value = 2) found first
	openShortcodes     map[string]bool // set of shortcodes in open state

	// items delivered to client
	items chan item
}

// note: the input position here is normally 0 (start), but
// can be set if position of first shortcode is known
func newShortcodeLexer(name, input string, inputPosition pos) *pagelexer {
	lexer := &pagelexer{
		name:               name,
		input:              input,
		currLeftDelimItem:  tLeftDelimScNoMarkup,
		currRightDelimItem: tRightDelimScNoMarkup,
		pos:                inputPosition,
		openShortcodes:     make(map[string]bool),
		items:              make(chan item),
	}
	go lexer.runShortcodeLexer()
	return lexer
}

// main loop
// this looks kind of funky, but it works
func (l *pagelexer) runShortcodeLexer() {
	for l.state = lexTextOutsideShortcodes; l.state != nil; {
		l.state = l.state(l)
	}

	close(l.items)
}

// state functions

const (
	leftDelimScNoMarkup    = "{{<"
	rightDelimScNoMarkup   = ">}}"
	leftDelimScWithMarkup  = "{{%"
	rightDelimScWithMarkup = "%}}"
	leftComment            = "/*" // comments in this context us used to to mark shortcodes as "not really a shortcode"
	rightComment           = "*/"
)

func (l *pagelexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}

	// looks expensive, but should produce the same iteration sequence as the string range loop
	// see: http://blog.golang.org/strings
	runeValue, runeWidth := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = pos(runeWidth)
	l.pos += l.width
	return runeValue
}

// peek, but no consume
func (l *pagelexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// steps back one
func (l *pagelexer) backup() {
	l.pos -= l.width
}

// sends an item back to the client.
func (l *pagelexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// special case, do not send '\\' back to client
func (l *pagelexer) ignoreEscapesAndEmit(t itemType) {
	val := strings.Map(func(r rune) rune {
		if r == '\\' {
			return -1
		}
		return r
	}, l.input[l.start:l.pos])
	l.items <- item{t, l.start, val}
	l.start = l.pos
}

// gets the current value (for debugging and error handling)
func (l *pagelexer) current() string {
	return l.input[l.start:l.pos]
}

// ignore current element
func (l *pagelexer) ignore() {
	l.start = l.pos
}

// nice to have in error logs
func (l *pagelexer) lineNum() int {
	return strings.Count(l.input[:l.lastPos], "\n") + 1
}

// nil terminates the parser
func (l *pagelexer) errorf(format string, args ...interface{}) stateFunc {
	l.items <- item{tError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// consumes and returns the next item
func (l *pagelexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

// scans until an opening shortcode opening bracket.
// if no shortcodes, it will keep on scanning until EOF
func lexTextOutsideShortcodes(l *pagelexer) stateFunc {
	for {
		if strings.HasPrefix(l.input[l.pos:], leftDelimScWithMarkup) || strings.HasPrefix(l.input[l.pos:], leftDelimScNoMarkup) {
			if l.pos > l.start {
				l.emit(tText)
			}
			if strings.HasPrefix(l.input[l.pos:], leftDelimScWithMarkup) {
				l.currLeftDelimItem = tLeftDelimScWithMarkup
				l.currRightDelimItem = tRightDelimScWithMarkup
			} else {
				l.currLeftDelimItem = tLeftDelimScNoMarkup
				l.currRightDelimItem = tRightDelimScNoMarkup
			}
			return lexShortcodeLeftDelim

		}
		if l.next() == eof {
			break
		}
	}
	// Done!
	if l.pos > l.start {
		l.emit(tText)
	}
	l.emit(tEOF)
	return nil
}

func lexShortcodeLeftDelim(l *pagelexer) stateFunc {
	l.pos += pos(len(l.currentLeftShortcodeDelim()))
	if strings.HasPrefix(l.input[l.pos:], leftComment) {
		return lexShortcodeComment
	}
	l.emit(l.currentLeftShortcodeDelimItem())
	l.elementStepNum = 0
	l.paramElements = 0
	return lexInsideShortcode
}

func lexShortcodeComment(l *pagelexer) stateFunc {
	posRightComment := strings.Index(l.input[l.pos:], rightComment)
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
	if !strings.HasPrefix(l.input[l.pos:], l.currentRightShortcodeDelim()) {
		return l.errorf("comment ends before the right shortcode delimiter")
	}
	l.pos += pos(len(l.currentRightShortcodeDelim()))
	l.emit(tText)
	return lexTextOutsideShortcodes
}

func lexShortcodeRightDelim(l *pagelexer) stateFunc {
	l.pos += pos(len(l.currentRightShortcodeDelim()))
	l.emit(l.currentRightShortcodeDelimItem())
	return lexTextOutsideShortcodes
}

// either:
// 1. param
// 2. "param" or "param\"
// 3. param="123" or param="123\"
// 4. param="Some \"escaped\" text"
func lexShortcodeParam(l *pagelexer, escapedQuoteStart bool) stateFunc {

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

func lexShortcodeQuotedParamVal(l *pagelexer, escapedQuotedValuesAllowed bool, typ itemType) stateFunc {
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
func lexIdentifierInShortcode(l *pagelexer) stateFunc {
	lookForEnd := false
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumericOrHyphen(r):
		default:
			l.backup()
			word := l.input[l.start:l.pos]
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

func lexEndOfShortcode(l *pagelexer) stateFunc {
	if strings.HasPrefix(l.input[l.pos:], l.currentRightShortcodeDelim()) {
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
func lexInsideShortcode(l *pagelexer) stateFunc {
	if strings.HasPrefix(l.input[l.pos:], l.currentRightShortcodeDelim()) {
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

func (l *pagelexer) currentLeftShortcodeDelimItem() itemType {
	return l.currLeftDelimItem
}

func (l *pagelexer) currentRightShortcodeDelimItem() itemType {
	return l.currRightDelimItem
}

func (l *pagelexer) currentLeftShortcodeDelim() string {
	if l.currLeftDelimItem == tLeftDelimScWithMarkup {
		return leftDelimScWithMarkup
	}
	return leftDelimScNoMarkup

}

func (l *pagelexer) currentRightShortcodeDelim() string {
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

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
