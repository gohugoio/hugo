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

package compare

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Strings returns an integer comparing two strings lexicographically.
func Strings(s, t string) int {
	c := compareFold(s, t)

	if c == 0 {
		// "B" and "b" would be the same so we need a tiebreaker.
		return strings.Compare(s, t)
	}

	return c
}

// This function is derived from strings.EqualFold in Go's stdlib.
// https://github.com/golang/go/blob/ad4a58e31501bce5de2aad90a620eaecdc1eecb8/src/strings/strings.go#L893
func compareFold(s, t string) int {
	for s != "" && t != "" {
		var sr, tr rune
		if s[0] < utf8.RuneSelf {
			sr, s = rune(s[0]), s[1:]
		} else {
			r, size := utf8.DecodeRuneInString(s)
			sr, s = r, s[size:]
		}
		if t[0] < utf8.RuneSelf {
			tr, t = rune(t[0]), t[1:]
		} else {
			r, size := utf8.DecodeRuneInString(t)
			tr, t = r, t[size:]
		}

		if tr == sr {
			continue
		}

		c := 1
		if tr < sr {
			tr, sr = sr, tr
			c = -c
		}

		//  ASCII only.
		if tr < utf8.RuneSelf {
			if sr >= 'A' && sr <= 'Z' {
				if tr <= 'Z' {
					// Same case.
					return -c
				}

				diff := tr - (sr + 'a' - 'A')

				if diff == 0 {
					continue
				}

				if diff < 0 {
					return c
				}

				if diff > 0 {
					return -c
				}
			}
		}

		// Unicode.
		r := unicode.SimpleFold(sr)
		for r != sr && r < tr {
			r = unicode.SimpleFold(r)
		}

		if r == tr {
			continue
		}

		return -c
	}

	if s == "" && t == "" {
		return 0
	}

	if s == "" {
		return -1
	}

	return 1
}

// LessStrings returns whether s is less than t lexicographically.
func LessStrings(s, t string) bool {
	return Strings(s, t) < 0
}
