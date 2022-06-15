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
	internalLinkLink     string
	internalLinkHasLabel bool
	internalLinkLabel    string
}

// Shortcode syntax
var (
	leftDelimInternalLink      = []byte("[[")
	rightDelimInternalLink     = []byte("]]")
	LabelDelimiterInternalLink = []byte("|")
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

	word := string(l.input[l.start:l.pos])
	l.internalLinkLabel = word

	if !l.internalLinkHasLabel {
		//Is both link&label
		l.internalLinkLink = word
		l.emit(tInternalLinkLinkLabel)
	} else {
		l.emit(tInternalLinkLabel)
	}

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
		return l.errorf("unclosed internal link")
	case r == '|':
		if l.internalLinkHasLabel {
			return l.errorf("internal link cannot have two or more pipes |")
		}
		l.internalLinkHasLabel = true
		word := string(l.input[l.start:(l.pos - 1)])
		l.internalLinkLink = word
		l.emit(tInternalLinkLink)
		return lexInsideInternalLink
	default:

		return lexInsideInternalLink
	}

}
