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

type lexerShortcodeState struct {
	currLeftDelimItem  ItemType
	currRightDelimItem ItemType
	isInline           bool
	currShortcodeName  string          // is only set when a shortcode is in opened state
	closingState       int             // > 0 = on its way to be closed
	elementStepNum     int             // step number in element
	paramElements      int             // number of elements (name + value = 2) found first
	openShortcodes     map[string]bool // set of shortcodes in open state
}

// Shortcode syntax
var (
	leftDelimSc            = []byte("{{")
	leftDelimScNoMarkup    = []byte("{{<")
	rightDelimScNoMarkup   = []byte(">}}")
	leftDelimScWithMarkup  = []byte("{{%")
	rightDelimScWithMarkup = []byte("%}}")
	leftComment            = []byte("/*") // comments in this context us used to to mark shortcodes as "not really a shortcode"
	rightComment           = []byte("*/")
)

func (l *pageLexer) isShortCodeStart() bool {
	return l.hasPrefix(leftDelimScWithMarkup) || l.hasPrefix(leftDelimScNoMarkup)
}

func lexShortcodeLeftDelim(l *pageLexer) stateFunc {
	l.pos += len(l.currentLeftShortcodeDelim())
	if l.hasPrefix(leftComment) {
		return lexShortcodeComment
	}
	l.emit(l.currentLeftShortcodeDelimItem())
	l.elementStepNum = 0
	l.paramElements = 0
	return lexInsideShortcode
}

func lexShortcodeComment(l *pageLexer) stateFunc {
	posRightComment := l.index(append(rightComment, l.currentRightShortcodeDelim()...))
	if posRightComment <= 1 {
		return l.errorf("comment must be closed")
	}
	// we emit all as text, except the comment markers
	l.emit(tText)
	l.pos += len(leftComment)
	l.ignore()
	l.pos += posRightComment - len(leftComment)
	l.emit(tText)
	l.pos += len(rightComment)
	l.ignore()
	l.pos += len(l.currentRightShortcodeDelim())
	l.emit(tText)
	return lexMainSection
}

func lexShortcodeRightDelim(l *pageLexer) stateFunc {
	l.closingState = 0
	l.pos += len(l.currentRightShortcodeDelim())
	l.emit(l.currentRightShortcodeDelimItem())
	return lexMainSection
}

// either:
// 1. param
// 2. "param" or "param\"
// 3. param="123" or param="123\"
// 4. param="Some \"escaped\" text"
// 5. `param`
// 6. param=`123`
func lexShortcodeParam(l *pageLexer, escapedQuoteStart bool) stateFunc {
	first := true
	nextEq := false

	var r rune

	for {
		r = l.next()
		if first {
			if r == '"' || (r == '`' && !escapedQuoteStart) {
				// a positional param with quotes
				if l.paramElements == 2 {
					return l.errorf("got quoted positional parameter. Cannot mix named and positional parameters")
				}
				l.paramElements = 1
				l.backup()
				if r == '"' {
					return lexShortcodeQuotedParamVal(l, !escapedQuoteStart, tScParam)
				}
				return lexShortCodeParamRawStringVal(l, tScParam)

			} else if r == '`' && escapedQuoteStart {
				return l.errorf("unrecognized escape character")
			}
			first = false
		} else if r == '=' {
			// a named param
			l.backup()
			nextEq = true
			break
		}

		if !isAlphaNumericOrHyphen(r) && r != '.' { // Floats have period
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

func lexShortcodeParamVal(l *pageLexer) stateFunc {
	l.consumeToSpace()
	l.emit(tScParamVal)
	return lexInsideShortcode
}

func lexShortCodeParamRawStringVal(l *pageLexer, typ ItemType) stateFunc {
	openBacktickFound := false

Loop:
	for {
		switch r := l.next(); {
		case r == '`':
			if openBacktickFound {
				l.backup()
				break Loop
			} else {
				openBacktickFound = true
				l.ignore()
			}
		case r == eof:
			return l.errorf("unterminated raw string in shortcode parameter-argument: '%s'", l.current())
		}
	}

	l.emitString(typ)
	l.next()
	l.ignore()

	return lexInsideShortcode
}

func lexShortcodeQuotedParamVal(l *pageLexer, escapedQuotedValuesAllowed bool, typ ItemType) stateFunc {
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
					// the coming quote is inside
					escapedInnerQuoteFound = true
					escapedQuoteState = 1
				}
			} else if l.peek() == '`' {
				return l.errorf("unrecognized escape character")
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
		l.ignoreEscapesAndEmit(typ, true)
	} else {
		l.emitString(typ)
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

// Inline shortcodes has the form {{< myshortcode.inline >}}
var inlineIdentifier = []byte("inline ")

// scans an alphanumeric inside shortcode
func lexIdentifierInShortcode(l *pageLexer) stateFunc {
	lookForEnd := false
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumericOrHyphen(r):
		// Allow forward slash inside names to make it possible to create namespaces.
		case r == '/':
		case r == '.':
			l.isInline = l.hasPrefix(inlineIdentifier)
			if !l.isInline {
				return l.errorf("period in shortcode name only allowed for inline identifiers")
			}
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
			if l.isInline {
				l.emit(tScNameInline)
			} else {
				l.emit(tScName)
			}
			break Loop
		}
	}

	if lookForEnd {
		return lexEndOfShortcode
	}
	return lexInsideShortcode
}

func lexEndOfShortcode(l *pageLexer) stateFunc {
	l.isInline = false
	if l.hasPrefix(l.currentRightShortcodeDelim()) {
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
	if l.hasPrefix(l.currentRightShortcodeDelim()) {
		return lexShortcodeRightDelim
	}
	switch r := l.next(); {
	case r == eof:
		// eol is allowed inside shortcodes; this may go to end of document before it fails
		return l.errorf("unclosed shortcode action")
	case isSpace(r), isEndOfLine(r):
		l.ignore()
	case r == '=':
		l.consumeSpace()
		l.ignore()
		peek := l.peek()
		if peek == '"' || peek == '\\' {
			return lexShortcodeQuotedParamVal(l, peek != '\\', tScParamVal)
		} else if peek == '`' {
			return lexShortCodeParamRawStringVal(l, tScParamVal)
		}
		return lexShortcodeParamVal
	case r == '/':
		if l.currShortcodeName == "" {
			return l.errorf("got closing shortcode, but none is open")
		}
		l.closingState++
		l.isInline = false
		l.emit(tScClose)
	case r == '\\':
		l.ignore()
		if l.peek() == '"' || l.peek() == '`' {
			return lexShortcodeParam(l, true)
		}
	case l.elementStepNum > 0 && (isAlphaNumericOrHyphen(r) || r == '"' || r == '`'): // positional params can have quotes
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

func (l *pageLexer) currentLeftShortcodeDelimItem() ItemType {
	return l.currLeftDelimItem
}

func (l *pageLexer) currentRightShortcodeDelimItem() ItemType {
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
