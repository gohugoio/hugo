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
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/markup/blackfriday"

	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"

	"github.com/gohugoio/hugo/common/text"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"

	bp "github.com/gohugoio/hugo/bufferpool"
)

func sanitizeAnchorNameString(s string, idType string) string {
	return string(sanitizeAnchorName([]byte(s), idType))
}

func sanitizeAnchorName(b []byte, idType string) []byte {
	return sanitizeAnchorNameWithHook(b, idType, nil)
}

func sanitizeAnchorNameWithHook(b []byte, idType string, hook func(buf *bytes.Buffer)) []byte {
	buf := bp.GetBuffer()

	if idType == goldmark_config.AutoIDTypeBlackfriday {
		// TODO(bep) make it more efficient.
		buf.WriteString(blackfriday.SanitizedAnchorName(string(b)))
	} else {
		asciiOnly := idType == goldmark_config.AutoIDTypeGitHubAscii

		if asciiOnly {
			// Normalize it to preserve accents if possible.
			b = text.RemoveAccents(b)
		}

		b = bytes.TrimSpace(b)

		for len(b) > 0 {
			r, size := utf8.DecodeRune(b)
			switch {
			case asciiOnly && size != 1:
			case r == '-' || r == ' ':
				buf.WriteRune('-')
			case isAlphaNumeric(r):
				buf.WriteRune(unicode.ToLower(r))
			default:
			}

			b = b[size:]
		}
	}

	if hook != nil {
		hook(buf)
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())

	bp.PutBuffer(buf)

	return result
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

var _ parser.IDs = (*idFactory)(nil)

type idFactory struct {
	idType     string
	vals       map[string]struct{}
	duplicates []string
}

func newIDFactory(idType string) *idFactory {
	return &idFactory{
		vals:   make(map[string]struct{}),
		idType: idType,
	}
}

type stringValuesProvider interface {
	StringValues() []string
}

var _ stringValuesProvider = (*idFactory)(nil)

func (ids *idFactory) StringValues() []string {
	values := make([]string, 0, len(ids.vals))
	for k := range ids.vals {
		values = append(values, k)
	}
	values = append(values, ids.duplicates...)
	return values
}

func (ids *idFactory) Generate(value []byte, kind ast.NodeKind) []byte {
	return sanitizeAnchorNameWithHook(value, ids.idType, func(buf *bytes.Buffer) {
		if buf.Len() == 0 {
			if kind == ast.KindHeading {
				buf.WriteString("heading")
			} else if kind == east.KindDefinitionTerm {
				buf.WriteString("term")
			} else {
				buf.WriteString("id")
			}
		}

		if _, found := ids.vals[util.BytesToReadOnlyString(buf.Bytes())]; found {
			// Append a hyphen and a number, starting with 1.
			buf.WriteRune('-')
			pos := buf.Len()
			for i := 1; ; i++ {
				buf.WriteString(strconv.Itoa(i))
				if _, found := ids.vals[util.BytesToReadOnlyString(buf.Bytes())]; !found {
					break
				}
				buf.Truncate(pos)
			}
		}
		ids.put(buf.String())
	})
}

func (ids *idFactory) put(s string) {
	if _, found := ids.vals[s]; found {
		ids.duplicates = append(ids.duplicates, s)
	} else {
		ids.vals[s] = struct{}{}
	}
}

func (ids *idFactory) Put(value []byte) {
	ids.put(string(value))
}
