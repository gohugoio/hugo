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

func lexIntroSection(l *pageLexer) stateFunc {
	l.summaryDivider = summaryDivider

LOOP:
	for {
		r := l.next()
		if r == eof {
			break
		}

		switch {
		case r == '+':
			return l.lexFrontMatterSection(TypeFrontMatterTOML, r, "TOML", delimTOML)
		case r == '-':
			return l.lexFrontMatterSection(TypeFrontMatterYAML, r, "YAML", delimYAML)
		case r == '{':
			return lexFrontMatterJSON
		case r == '#':
			return lexFrontMatterOrgMode
		case r == byteOrderMark:
			l.emit(TypeIgnore)
		case !isSpace(r) && !isEndOfLine(r):
			if r == '<' {
				l.backup()
				if l.hasPrefix(htmlCommentStart) {
					// This may be commented out front matter, which should
					// still be read.
					l.consumeToNextLine()
					l.isInHTMLComment = true
					l.emit(TypeIgnore)
					continue LOOP
				} else {
					return l.errorf("plain HTML documents not supported")
				}
			}
			break LOOP
		}
	}

	// Now move on to the shortcodes.
	return lexMainSection
}

func lexEndFromtMatterHTMLComment(l *pageLexer) stateFunc {
	l.isInHTMLComment = false
	right := l.index(htmlCommentEnd)
	if right == -1 {
		return l.errorf("starting HTML comment with no end")
	}
	l.pos += right + len(htmlCommentEnd)
	l.emit(TypeIgnore)

	// Now move on to the shortcodes.
	return lexMainSection
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
	l.emit(TypeFrontMatterJSON)

	return lexMainSection
}

func lexFrontMatterOrgMode(l *pageLexer) stateFunc {
	/*
		#+TITLE: Test File For chaseadamsio/goorgeous
		#+AUTHOR: Chase Adams
		#+DESCRIPTION: Just another golang parser for org content!
	*/

	l.summaryDivider = summaryDividerOrg

	l.backup()

	if !l.hasPrefix(delimOrg) {
		return lexMainSection
	}

	// Read lines until we no longer see a #+ prefix
LOOP:
	for {

		r := l.next()

		switch {
		case r == '\n':
			if !l.hasPrefix(delimOrg) {
				break LOOP
			}
		case r == eof:
			break LOOP

		}
	}

	l.emit(TypeFrontMatterORG)

	return lexMainSection

}

// Handle YAML or TOML front matter.
func (l *pageLexer) lexFrontMatterSection(tp ItemType, delimr rune, name string, delim []byte) stateFunc {

	for i := 0; i < 2; i++ {
		if r := l.next(); r != delimr {
			return l.errorf("invalid %s delimiter", name)
		}
	}

	// Let front matter start at line 1
	wasEndOfLine := l.consumeCRLF()
	// We don't care about the delimiters.
	l.ignore()

	var r rune

	for {
		if !wasEndOfLine {
			r = l.next()
			if r == eof {
				return l.errorf("EOF looking for end %s front matter delimiter", name)
			}
		}

		if wasEndOfLine || isEndOfLine(r) {
			if l.hasPrefix(delim) {
				l.emit(tp)
				l.pos += 3
				l.consumeCRLF()
				l.ignore()
				break
			}
		}

		wasEndOfLine = false
	}

	return lexMainSection
}
