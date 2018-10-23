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

package hugolib

import (
	"bytes"
	"io"

	errors "github.com/pkg/errors"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
)

var (
	internalSummaryDivider = []byte("HUGOMORE42")
)

// The content related items on a Page.
type pageContent struct {
	renderable bool

	// workContent is a copy of rawContent that may be mutated during site build.
	workContent []byte

	shortcodeState *shortcodeHandler

	source rawPageContent
}

type rawPageContent struct {
	hasSummaryDivider bool

	// The AST of the parsed page. Contains information about:
	// shortcBackup3odes, front matter, summary indicators.
	// TODO(bep) 2errors add this to a new rawPagecContent struct
	// with frontMatterItem (pos) etc.
	// * RawContent, RawContentWithoutFrontMatter
	parsed pageparser.Result
}

// TODO(bep) lazy consolidate
func (p *Page) mapContent() error {
	p.shortcodeState = newShortcodeHandler(p)
	s := p.shortcodeState
	p.renderable = true

	result := bp.GetBuffer()
	defer bp.PutBuffer(result)

	iter := p.source.parsed.Iterator()

	fail := func(err error, i pageparser.Item) error {
		return parseError(err, iter.Input(), i.Pos)
	}

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currShortcode shortcode
	var ordinal int

Loop:
	for {
		it := iter.Next()

		switch {
		case it.Type == pageparser.TypeIgnore:
		case it.Type == pageparser.TypeHTMLComment:
			// Ignore. This is only a leading Front matter comment.
		case it.Type == pageparser.TypeHTMLDocument:
			// This is HTML only. No shortcode, front matter etc.
			p.renderable = false
			result.Write(it.Val)
		case it.IsFrontMatter():
			f := metadecoders.FormatFromFrontMatterType(it.Type)
			m, err := metadecoders.UnmarshalToMap(it.Val, f)
			if err != nil {
				if fe, ok := err.(herrors.FileError); ok {
					return herrors.ToFileErrorWithOffset(fe, iter.LineNumber()-1)
				} else {
					return err
				}
			}
			if err := p.updateMetaData(m); err != nil {
				return err
			}

			if !p.shouldBuild() {
				// Nothing more to do.
				return nil
			}

		case it.Type == pageparser.TypeLeadSummaryDivider, it.Type == pageparser.TypeSummaryDividerOrg:
			result.Write(internalSummaryDivider)
			p.source.hasSummaryDivider = true
			// Need to determine if the page is truncated.
			f := func(item pageparser.Item) bool {
				if item.IsNonWhitespace() {
					p.truncated = true

					// Done
					return false
				}
				return true
			}
			iter.PeekWalk(f)

		// Handle shortcode
		case it.IsLeftShortcodeDelim():
			// let extractShortcode handle left delim (will do so recursively)
			iter.Backup()

			currShortcode, err := s.extractShortcode(ordinal, iter, p)

			if currShortcode.name != "" {
				s.nameSet[currShortcode.name] = true
			}

			if err != nil {
				return fail(errors.Wrap(err, "failed to extract shortcode"), it)
			}

			if currShortcode.params == nil {
				currShortcode.params = make([]string, 0)
			}

			placeHolder := s.createShortcodePlaceholder()
			result.WriteString(placeHolder)
			ordinal++
			s.shortcodes.Add(placeHolder, currShortcode)
		case it.IsEOF():
			break Loop
		case it.IsError():
			err := fail(errors.WithStack(errors.New(it.ValStr())), it)
			currShortcode.err = err
			return err

		default:
			result.Write(it.Val)
		}
	}

	resultBytes := make([]byte, result.Len())
	copy(resultBytes, result.Bytes())
	p.workContent = resultBytes

	return nil
}

func (p *Page) parse(reader io.Reader) error {

	parseResult, err := pageparser.Parse(reader)
	if err != nil {
		return err
	}

	p.source = rawPageContent{
		parsed: parseResult,
	}

	// TODO(bep) 2errors
	p.lang = p.File.Lang()

	if p.s != nil && p.s.owner != nil {
		gi, enabled := p.s.owner.gitInfo.forPage(p)
		if gi != nil {
			p.GitInfo = gi
		} else if enabled {
			p.s.Log.WARN.Printf("Failed to find GitInfo for page %q", p.Path())
		}
	}

	return nil
}

func parseError(err error, input []byte, pos int) error {
	if herrors.UnwrapFileError(err) != nil {
		// Use the most specific location.
		return err
	}
	lf := []byte("\n")
	input = input[:pos]
	lineNumber := bytes.Count(input, lf) + 1
	endOfLastLine := bytes.LastIndex(input, lf)
	return herrors.NewFileError("md", -1, lineNumber, pos-endOfLastLine, err)

}
