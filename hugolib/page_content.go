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
	"fmt"
	"io"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
)

// The content related items on a Page.
type pageContent struct {
	renderable bool

	frontmatter []byte

	// rawContent is the raw content read from the content file.
	rawContent []byte

	// workContent is a copy of rawContent that may be mutated during site build.
	workContent []byte

	shortcodeState *shortcodeHandler

	source rawPageContent
}

type rawPageContent struct {
	// The AST of the parsed page. Contains information about:
	// shortcBackup3odes, front matter, summary indicators.
	// TODO(bep) 2errors add this to a new rawPagecContent struct
	// with frontMatterItem (pos) etc.
	// * also Result.Iterator, Result.Source
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

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currShortcode shortcode
	var ordinal int

Loop:
	for {
		it := iter.Next()

		switch {
		case it.Typ == pageparser.TypeIgnore:
		case it.Typ == pageparser.TypeHTMLComment:
			// Ignore. This is only a leading Front matter comment.
		case it.Typ == pageparser.TypeHTMLDocument:
			// This is HTML only. No shortcode, front matter etc.
			p.renderable = false
			result.Write(it.Val)
			// TODO(bep) 2errors commented out frontmatter
		case it.IsFrontMatter():
			f := metadecoders.FormatFromFrontMatterType(it.Typ)
			m, err := metadecoders.UnmarshalToMap(it.Val, f)
			if err != nil {
				return err
			}
			if err := p.updateMetaData(m); err != nil {
				return err
			}

			if !p.shouldBuild() {
				// Nothing more to do.
				return nil

			}

		//case it.Typ == pageparser.TypeLeadSummaryDivider, it.Typ == pageparser.TypeSummaryDividerOrg:
		// TODO(bep) 2errors store if divider is there and use that to determine if replace or not
		// Handle shortcode
		case it.IsLeftShortcodeDelim():
			// let extractShortcode handle left delim (will do so recursively)
			iter.Backup()

			currShortcode, err := s.extractShortcode(ordinal, iter, p)

			if currShortcode.name != "" {
				s.nameSet[currShortcode.name] = true
			}

			if err != nil {
				return err
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
			err := fmt.Errorf("%s:shortcode:%d: %s",
				p.pathOrTitle(), iter.LineNumber(), it)
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
	p.lang = p.Source.File.Lang()

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
