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

// The lexical scanning below

type Tokens struct {
	lexer     *pageLexer
	token     [3]Item // 3-item look-ahead is what we currently need
	peekCount int
}

func (t *Tokens) Next() Item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lexer.nextItem()
	}
	return t.token[t.peekCount]
}

// backs up one token.
func (t *Tokens) Backup() {
	t.peekCount++
}

// backs up two tokens.
func (t *Tokens) Backup2(t1 Item) {
	t.token[1] = t1
	t.peekCount = 2
}

// backs up three tokens.
func (t *Tokens) Backup3(t2, t1 Item) {
	t.token[1] = t1
	t.token[2] = t2
	t.peekCount = 3
}

// check for non-error and non-EOF types coming next
func (t *Tokens) IsValueNext() bool {
	i := t.Peek()
	return i.typ != tError && i.typ != tEOF
}

// look at, but do not consume, the next item
// repeated, sequential calls will return the same item
func (t *Tokens) Peek() Item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lexer.nextItem()
	return t.token[0]
}

// Consume is a convencience method to consume the next n tokens,
// but back off Errors and EOF.
func (t *Tokens) Consume(cnt int) {
	for i := 0; i < cnt; i++ {
		token := t.Next()
		if token.typ == tError || token.typ == tEOF {
			t.Backup()
			break
		}
	}
}

// LineNumber returns the current line number. Used for logging.
func (t *Tokens) LineNumber() int {
	return t.lexer.lineNum()
}
