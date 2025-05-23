// Copyright 2024 The Hugo Authors. All rights reserved.
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

package hstrings

import (
	"regexp"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestStringEqualFold(t *testing.T) {
	c := qt.New(t)

	s1 := "A"
	s2 := "a"

	c.Assert(StringEqualFold(s1).EqualFold(s2), qt.Equals, true)
	c.Assert(StringEqualFold(s1).EqualFold(s1), qt.Equals, true)
	c.Assert(StringEqualFold(s2).EqualFold(s1), qt.Equals, true)
	c.Assert(StringEqualFold(s2).EqualFold(s2), qt.Equals, true)
	c.Assert(StringEqualFold(s1).EqualFold("b"), qt.Equals, false)
	c.Assert(StringEqualFold(s1).Eq(s2), qt.Equals, true)
	c.Assert(StringEqualFold(s1).Eq("b"), qt.Equals, false)
}

func TestGetOrCompileRegexp(t *testing.T) {
	c := qt.New(t)

	re, err := GetOrCompileRegexp(`\d+`)
	c.Assert(err, qt.IsNil)
	c.Assert(re.MatchString("123"), qt.Equals, true)
}

func BenchmarkGetOrCompileRegexp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetOrCompileRegexp(`\d+`)
	}
}

func BenchmarkCompileRegexp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		regexp.MustCompile(`\d+`)
	}
}

func BenchmarkCountWordsASCII(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	sLong := strings.Repeat(s, 400) // 9 x 400 = 3600 words
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountWords(sLong)
	}
}

func BenchmarkCountWordsMixed(b *testing.B) {
	s := "The 素早い brown fox jumps over the lazy 犬"
	sLong := strings.Repeat(s, 400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountWords(sLong)
	}
}

func BenchmarkCountWordsChinese(b *testing.B) {
	s := "敏捷的棕狐狸跳过懒惰的狗"
	sLong := strings.Repeat(s, 400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountWords(sLong)
	}
}
