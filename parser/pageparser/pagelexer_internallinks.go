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

type lexerInternalLinkState struct {
	currLeftDelimItem  ItemType
	currRightDelimItem ItemType
	// isInline           bool
	// currShortcodeName  string          // is only set when a shortcode is in opened state
	// closingState       int             // > 0 = on its way to be closed
	// elementStepNum     int             // step number in element
	// paramElements      int             // number of elements (name + value = 2) found first
	// openShortcodes     map[string]bool // set of shortcodes in open state

}

// Shortcode syntax
var (
	leftDelimInternalLink  = []byte("[[")
	rightDelimInternalLink = []byte("]]")
	textDelimiter          = []byte("|")
)

func (l *pageLexer) isInternalLinkStart() bool {
	return l.hasPrefix(leftDelimInternalLink)
}

func lexInternalLinkLeftDelim(l *pageLexer) stateFunc {
	l.pos += len(leftDelimInternalLink)
	l.emit(tLeftDelimInternalLink)
	return lexInsideInternalLink
}

func lexInternalLinkRightDelim(l *pageLexer) stateFunc {
	l.pos += len(rightDelimInternalLink)
	l.emit(tRightDelimInternalLink)
	return lexMainSection
}

// scans the elements inside internal link [[]] tags
func lexInsideInternalLink(l *pageLexer) stateFunc {
	if l.hasPrefix(rightDelimInternalLink) {
		return lexInternalLinkRightDelim
	}
	switch r := l.next(); {
	case r == eof:
		// eol is allowed inside shortcodes; this may go to end of document before it fails
		return l.errorf("unclosed internal link")
	case isSpace(r), isEndOfLine(r):
		l.ignore()
		return lexInsideInternalLink
	case r == '|':
		/*	l.consumeSpace()
			l.ignore()
			peek := l.peek()
			if peek == '"' || peek == '\\' {
				return lexShortcodeQuotedParamVal(l, peek != '\\', tScParamVal)
			} else if peek == '`' {
				return lexShortCodeParamRawStringVal(l, tScParamVal)
			}
			return lexShortcodeParamVal*/
		l.ignore()
		return lexInsideInternalLink
	default:
		return l.errorf("unrecognized character in shortcode action: %#U. Note: Parameters with non-alphanumeric args must be quoted", r)
	}
	return lexInsideInternalLink
}
