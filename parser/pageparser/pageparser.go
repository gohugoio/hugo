// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"io"
	"io/ioutil"

	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/pkg/errors"
)

// Result holds the parse result.
type Result interface {
	// Iterator returns a new Iterator positioned at the beginning of the parse tree.
	Iterator() *Iterator
	// Input returns the input to Parse.
	Input() []byte
}

var _ Result = (*pageLexer)(nil)

// Parse parses the page in the given reader according to the given Config.
// TODO(bep) now that we have improved the "lazy order" init, it *may* be
// some potential saving in doing a buffered approach where the first pass does
// the frontmatter only.
func Parse(r io.Reader, cfg Config) (Result, error) {
	return parseSection(r, cfg, lexIntroSection)
}

type ContentFrontMatter struct {
	Content           []byte
	FrontMatter       map[string]interface{}
	FrontMatterFormat metadecoders.Format
}

// ParseFrontMatterAndContent is a convenience method to extract front matter
// and content from a content page.
func ParseFrontMatterAndContent(r io.Reader) (ContentFrontMatter, error) {
	var cf ContentFrontMatter

	psr, err := Parse(r, Config{})
	if err != nil {
		return cf, err
	}

	var frontMatterSource []byte

	iter := psr.Iterator()

	walkFn := func(item Item) bool {
		if frontMatterSource != nil {
			// The rest is content.
			cf.Content = psr.Input()[item.Pos:]
			// Done
			return false
		} else if item.IsFrontMatter() {
			cf.FrontMatterFormat = FormatFromFrontMatterType(item.Type)
			frontMatterSource = item.Val
		}
		return true

	}

	iter.PeekWalk(walkFn)

	cf.FrontMatter, err = metadecoders.Default.UnmarshalToMap(frontMatterSource, cf.FrontMatterFormat)
	return cf, err
}

func FormatFromFrontMatterType(typ ItemType) metadecoders.Format {
	switch typ {
	case TypeFrontMatterJSON:
		return metadecoders.JSON
	case TypeFrontMatterORG:
		return metadecoders.ORG
	case TypeFrontMatterTOML:
		return metadecoders.TOML
	case TypeFrontMatterYAML:
		return metadecoders.YAML
	default:
		return ""
	}
}

// ParseMain parses starting with the main section. Used in tests.
func ParseMain(r io.Reader, cfg Config) (Result, error) {
	return parseSection(r, cfg, lexMainSection)
}

func parseSection(r io.Reader, cfg Config, start stateFunc) (Result, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read page content")
	}
	return parseBytes(b, cfg, start)
}

func parseBytes(b []byte, cfg Config, start stateFunc) (Result, error) {
	lexer := newPageLexer(b, start, cfg)
	lexer.run()
	return lexer, nil
}

// An Iterator has methods to iterate a parsed page with support going back
// if needed.
type Iterator struct {
	l       *pageLexer
	lastPos int // position of the last item returned by nextItem
}

// consumes and returns the next item
func (t *Iterator) Next() Item {
	t.lastPos++
	return t.Current()
}

// Input returns the input source.
func (t *Iterator) Input() []byte {
	return t.l.Input()
}

var errIndexOutOfBounds = Item{tError, 0, []byte("no more tokens"), true}

// Current will repeatably return the current item.
func (t *Iterator) Current() Item {
	if t.lastPos >= len(t.l.items) {
		return errIndexOutOfBounds
	}
	return t.l.items[t.lastPos]
}

// backs up one token.
func (t *Iterator) Backup() {
	if t.lastPos < 0 {
		panic("need to go forward before going back")
	}
	t.lastPos--
}

// check for non-error and non-EOF types coming next
func (t *Iterator) IsValueNext() bool {
	i := t.Peek()
	return i.Type != tError && i.Type != tEOF
}

// look at, but do not consume, the next item
// repeated, sequential calls will return the same item
func (t *Iterator) Peek() Item {
	return t.l.items[t.lastPos+1]
}

// PeekWalk will feed the next items in the iterator to walkFn
// until it returns false.
func (t *Iterator) PeekWalk(walkFn func(item Item) bool) {
	for i := t.lastPos + 1; i < len(t.l.items); i++ {
		item := t.l.items[i]
		if !walkFn(item) {
			break
		}
	}
}

// Consume is a convencience method to consume the next n tokens,
// but back off Errors and EOF.
func (t *Iterator) Consume(cnt int) {
	for i := 0; i < cnt; i++ {
		token := t.Next()
		if token.Type == tError || token.Type == tEOF {
			t.Backup()
			break
		}
	}
}

// LineNumber returns the current line number. Used for logging.
func (t *Iterator) LineNumber() int {
	return bytes.Count(t.l.input[:t.Current().Pos], lf) + 1
}
