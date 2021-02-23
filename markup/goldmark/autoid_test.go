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
	"strings"
	"testing"

	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"

	qt "github.com/frankban/quicktest"
)

func TestSanitizeAnchorName(t *testing.T) {
	c := qt.New(t)

	// Tests generated manually on github.com
	tests := `
God is good: 神真美好
Number 32
Question?
1+2=3
Special !"#$%&(parens)=?´* chars
Resumé
One-Hyphen
Multiple--Hyphens
Trailing hyphen-
Many   spaces  here
Forward/slash
Backward\slash
Under_score
Nonbreaking Space
Tab	Space
`

	expect := `
god-is-good-神真美好
number-32
question
123
special-parens-chars
resumé
one-hyphen
multiple--hyphens
trailing-hyphen-
many---spaces--here
forwardslash
backwardslash
under_score
nonbreakingspace
tabspace
`

	tests, expect = strings.TrimSpace(tests), strings.TrimSpace(expect)

	testlines, expectlines := strings.Split(tests, "\n"), strings.Split(expect, "\n")

	testlines = append(testlines, "Trailing Space ")
	expectlines = append(expectlines, "trailing-space")

	if len(testlines) != len(expectlines) {
		panic("test setup failed")
	}

	for i, input := range testlines {
		input := input
		expect := expectlines[i]
		c.Run(input, func(c *qt.C) {
			b := []byte(input)
			got := string(sanitizeAnchorName(b, goldmark_config.AutoHeadingIDTypeGitHub))
			c.Assert(got, qt.Equals, expect)
			c.Assert(sanitizeAnchorNameString(input, goldmark_config.AutoHeadingIDTypeGitHub), qt.Equals, expect)
			c.Assert(string(b), qt.Equals, input)
		})
	}
}

func TestSanitizeAnchorNameAsciiOnly(t *testing.T) {
	c := qt.New(t)

	c.Assert(sanitizeAnchorNameString("god is神真美好 good", goldmark_config.AutoHeadingIDTypeGitHubAscii), qt.Equals, "god-is-good")
	c.Assert(sanitizeAnchorNameString("Resumé", goldmark_config.AutoHeadingIDTypeGitHubAscii), qt.Equals, "resume")
}

func TestSanitizeAnchorNameBlackfriday(t *testing.T) {
	c := qt.New(t)
	c.Assert(sanitizeAnchorNameString("Let's try this, shall we?", goldmark_config.AutoHeadingIDTypeBlackfriday), qt.Equals, "let-s-try-this-shall-we")
}

func BenchmarkSanitizeAnchorName(b *testing.B) {
	input := []byte("God is good: 神真美好")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := sanitizeAnchorName(input, goldmark_config.AutoHeadingIDTypeGitHub)
		if len(result) != 24 {
			b.Fatalf("got %d", len(result))
		}
	}
}

func BenchmarkSanitizeAnchorNameAsciiOnly(b *testing.B) {
	input := []byte("God is good: 神真美好")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := sanitizeAnchorName(input, goldmark_config.AutoHeadingIDTypeGitHubAscii)
		if len(result) != 12 {
			b.Fatalf("got %d", len(result))
		}
	}
}

func BenchmarkSanitizeAnchorNameBlackfriday(b *testing.B) {
	input := []byte("God is good: 神真美好")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := sanitizeAnchorName(input, goldmark_config.AutoHeadingIDTypeBlackfriday)
		if len(result) != 24 {
			b.Fatalf("got %d", len(result))
		}
	}
}

func BenchmarkSanitizeAnchorNameString(b *testing.B) {
	input := "God is good: 神真美好"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := sanitizeAnchorNameString(input, goldmark_config.AutoHeadingIDTypeGitHub)
		if len(result) != 24 {
			b.Fatalf("got %d", len(result))
		}
	}
}
