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

package pageparser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/parser/metadecoders"
)

// Result holds the parse result.
type Result interface {
	// Iterator returns a new Iterator positioned at the beginning of the parse tree.
	Iterator() *Iterator
	// Input returns the input to Parse.
	Input() []byte
}

var _ Result = (*pageLexer)(nil)

// ParseBytes parses the page in b according to the given Config.
func ParseBytes(b []byte, cfg Config) (Items, error) {
	startLexer := lexIntroSection
	if cfg.NoFrontMatter {
		startLexer = lexMainSection
	}
	l, err := parseBytes(b, cfg, startLexer)
	if err != nil {
		return nil, err
	}
	return l.items, l.err
}

type ContentFrontMatter struct {
	Content           []byte
	FrontMatter       map[string]any
	FrontMatterFormat metadecoders.Format
}

// ParseFrontMatterAndContent is a convenience method to extract front matter
// and content from a content page.
func ParseFrontMatterAndContent(r io.Reader) (ContentFrontMatter, error) {
	var cf ContentFrontMatter

	input, err := io.ReadAll(r)
	if err != nil {
		return cf, fmt.Errorf("failed to read page content: %w", err)
	}

	psr, err := ParseBytes(input, Config{})
	if err != nil {
		return cf, err
	}

	var frontMatterSource []byte

	iter := NewIterator(psr)

	walkFn := func(item Item) bool {
		if frontMatterSource != nil {
			// The rest is content.
			cf.Content = input[item.low:]
			// Done
			return false
		} else if item.IsFrontMatter() {
			cf.FrontMatterFormat = FormatFromFrontMatterType(item.Type)
			frontMatterSource = item.Val(input)
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
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read page content: %w", err)
	}
	return parseBytes(b, cfg, start)
}

func parseBytes(b []byte, cfg Config, start stateFunc) (*pageLexer, error) {
	lexer := newPageLexer(b, start, cfg)
	lexer.run()
	return lexer, nil
}

// NewIterator creates a new Iterator.
func NewIterator(items Items) *Iterator {
	return &Iterator{items: items, lastPos: -1}
}

// An Iterator has methods to iterate a parsed page with support going back
// if needed.
type Iterator struct {
	items   Items
	lastPos int // position of the last item returned by nextItem
}

// consumes and returns the next item
func (t *Iterator) Next() Item {
	t.lastPos++
	return t.Current()
}

var errIndexOutOfBounds = Item{Type: tError, Err: errors.New("no more tokens")}

// Current will repeatably return the current item.
func (t *Iterator) Current() Item {
	if t.lastPos >= len(t.items) {
		return errIndexOutOfBounds
	}
	return t.items[t.lastPos]
}

// backs up one token.
func (t *Iterator) Backup() {
	if t.lastPos < 0 {
		panic("need to go forward before going back")
	}
	t.lastPos--
}

// Pos returns the current position in the input.
func (t *Iterator) Pos() int {
	return t.lastPos
}

// check for non-error and non-EOF types coming next
func (t *Iterator) IsValueNext() bool {
	i := t.Peek()
	return i.Type != tError && i.Type != tEOF
}

// look at, but do not consume, the next item
// repeated, sequential calls will return the same item
func (t *Iterator) Peek() Item {
	return t.items[t.lastPos+1]
}

// PeekWalk will feed the next items in the iterator to walkFn
// until it returns false.
func (t *Iterator) PeekWalk(walkFn func(item Item) bool) {
	for i := t.lastPos + 1; i < len(t.items); i++ {
		item := t.items[i]
		if !walkFn(item) {
			break
		}
	}
}

// Consume is a convenience method to consume the next n tokens,
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
func (t *Iterator) LineNumber(source []byte) int {
	return bytes.Count(source[:t.Current().low], lf) + 1
}

// IsProbablySourceOfItems returns true if the given source looks like original
// source of the items.
// There may be some false positives, but that is highly unlikely and good enough
// for the planned purpose.
// It will also return false if the last item is not EOF (error situations) and
// true if both source and items are empty.
func IsProbablySourceOfItems(source []byte, items Items) bool {
	if len(source) == 0 && len(items) == 0 {
		return false
	}
	if len(items) == 0 {
		return false
	}

	last := items[len(items)-1]
	if last.Type != tEOF {
		return false
	}

	if last.Pos() != len(source) {
		return false
	}

	for _, item := range items {
		if item.Type == tError {
			return false
		}
		if item.Type == tEOF {
			return true
		}

		if item.Pos() >= len(source) {
			return false
		}

		if item.firstByte != source[item.Pos()] {
			return false
		}
	}

	return true
}

var hasShortcodeRe = regexp.MustCompile(`{{[%,<][^\/]`)

// HasShortcode returns true if the given string contains a shortcode.
func HasShortcode(s string) bool {
	// Fast path for the common case.
	if !strings.Contains(s, "{{") {
		return false
	}
	return hasShortcodeRe.MatchString(s)
}
