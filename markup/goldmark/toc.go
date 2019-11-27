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

package goldmark

import (
	"bytes"

	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	tocResultKey     = parser.NewContextKey()
	renderContextKey = parser.NewContextKey()
)

// TODO1 revert the content spec changes

type tocAstHook struct {
	reader text.Reader
	pc     parser.Context

	// ToC state
	toc         tableofcontents.Root
	header      tableofcontents.Header
	level       int
	row         int
	inHeading   bool
	headingText bytes.Buffer
}

func newTocAstHook(reader text.Reader, pc parser.Context) *tocAstHook {
	return &tocAstHook{
		reader: reader,
		pc:     pc,
		row:    -1,
	}
}

func (h *tocAstHook) Visit(n ast.Node, entering bool) (ast.WalkStatus, error) {
	s := ast.WalkStatus(ast.WalkContinue)
	if n.Kind() == ast.KindHeading {
		if h.inHeading && !entering {
			h.header.Text = h.headingText.String()
			h.headingText.Reset()
			h.toc.AddAt(h.header, h.row, h.level-1)
			h.header = tableofcontents.Header{}
			h.inHeading = false
			return s, nil
		}

		h.inHeading = true
	}

	if !(h.inHeading && entering) {
		return s, nil
	}

	switch n.Kind() {
	case ast.KindHeading:
		heading := n.(*ast.Heading)
		h.level = heading.Level

		if h.level == 1 || h.row == -1 {
			h.row++
		}

		id, found := heading.AttributeString("id")
		if found {
			h.header.ID = string(id.([]byte))
		}
	case ast.KindText:
		textNode := n.(*ast.Text)
		h.headingText.Write(textNode.Text(h.reader.Source()))
	}

	return s, nil
}

func (h *tocAstHook) Done() {
	h.pc.Set(tocResultKey, h.toc)
}
