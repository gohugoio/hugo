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

package text

import (
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var transliteratePool = &sync.Pool{
	New: func() any {
		return transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)),
			runes.Map(func(r rune) rune {
				switch r {
				case 'ą':
					return 'a'
				case 'ć':
					return 'c'
				case 'ę':
					return 'e'
				case 'ł':
					return 'l'
				case 'ń':
					return 'n'
				case 'ó':
					return 'o'
				case 'ś':
					return 's'
				case 'ż':
					return 'z'
				case 'ź':
					return 'z'
				case 'ø':
					return 'o'
				}
				return r
			}),
			norm.NFC)
	},
}

var transliterateMap = map[rune]rune{
	'ą': 'a',
	'ć': 'c',
	'ę': 'e',
	'ł': 'l',
	'ń': 'n',
	'ó': 'o',
	'ś': 's',
	'ż': 'z',
	'ź': 'z',
	'ø': 'o',
}

var transliteratePoolMap = &sync.Pool{
	New: func() any {
		return transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)),
			runes.Map(func(r rune) rune {
				if rr, ok := transliterateMap[r]; ok {
					return rr
				}
				return r
			}),
			norm.NFC)
	},
}

var accentTransformerPool = &sync.Pool{
	New: func() any {
		return transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	},
}

// RemoveAccents removes all accents from b.
func RemoveAccents(b []byte) []byte {
	t := accentTransformerPool.Get().(transform.Transformer)
	b, _, _ = transform.Bytes(t, b)
	t.Reset()
	accentTransformerPool.Put(t)
	return b
}

// RemoveAccentsString removes all accents from s.
func RemoveAccentsString(s string) string {
	t := accentTransformerPool.Get().(transform.Transformer)
	s, _, _ = transform.String(t, s)
	t.Reset()
	accentTransformerPool.Put(t)
	return s
}

func TransliterateString(s string) string {
	t := transliteratePool.Get().(transform.Transformer)
	s, _, _ = transform.String(t, s)
	t.Reset()
	transliteratePool.Put(t)
	return s
}

func TransliterateStringMap(s string) string {
	t := transliteratePoolMap.Get().(transform.Transformer)
	s, _, _ = transform.String(t, s)
	t.Reset()
	transliteratePoolMap.Put(t)
	return s
}

// Chomp removes trailing newline characters from s.
func Chomp(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
}

// Puts adds a trailing \n none found.
func Puts(s string) string {
	if s == "" || s[len(s)-1] == '\n' {
		return s
	}
	return s + "\n"
}

// VisitLinesAfter calls the given function for each line, including newlines, in the given string.
func VisitLinesAfter(s string, fn func(line string)) {
	high := strings.IndexRune(s, '\n')
	for high != -1 {
		fn(s[:high+1])
		s = s[high+1:]

		high = strings.IndexRune(s, '\n')
	}

	if s != "" {
		fn(s)
	}
}
