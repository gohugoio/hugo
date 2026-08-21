// Copyright 2017 The Hugo Authors. All rights reserved.
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

package metrics

import (
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/resources/page"

	qt "github.com/frankban/quicktest"
)

func TestSimilarPercentage(t *testing.T) {
	c := qt.New(t)

	sentence := "this is some words about nothing, Hugo!"
	words := strings.Fields(sentence)
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	sentenceReversed := strings.Join(words, " ")

	c.Assert(howSimilar("Hugo Rules", "Hugo Rules"), qt.Equals, 100)
	c.Assert(howSimilar("Hugo Rules", "Hugo Rocks"), qt.Equals, 50)
	c.Assert(howSimilar("The Hugo Rules", "The Hugo Rocks"), qt.Equals, 66)
	c.Assert(howSimilar("The Hugo Rules", "The Hugo"), qt.Equals, 66)
	c.Assert(howSimilar("The Hugo", "The Hugo Rules"), qt.Equals, 66)
	c.Assert(howSimilar("Totally different", "Not Same"), qt.Equals, 0)
	c.Assert(howSimilar(sentence, sentenceReversed), qt.Equals, 14)
	c.Assert(howSimilar(template.HTML("Hugo Rules"), template.HTML("Hugo Rules")), qt.Equals, 100)
	c.Assert(howSimilar(map[string]any{"a": 32, "b": 33}, map[string]any{"a": 32, "b": 33}), qt.Equals, 100)
	c.Assert(howSimilar(map[string]any{"a": 32, "b": 33}, map[string]any{"a": 32, "b": 34}), qt.Equals, 0)
	c.Assert(howSimilar("\n", ""), qt.Equals, 100)
}

type testStruct struct {
	Name string
}

func TestSimilarPercentageNonString(t *testing.T) {
	c := qt.New(t)
	c.Assert(howSimilar(page.NopPage, page.NopPage), qt.Equals, 100)
	c.Assert(howSimilar(page.Pages{}, page.Pages{}), qt.Equals, 90)
	c.Assert(howSimilar(testStruct{Name: "A"}, testStruct{Name: "B"}), qt.Equals, 0)
	c.Assert(howSimilar(testStruct{Name: "A"}, testStruct{Name: "A"}), qt.Equals, 100)
}

func BenchmarkHowSimilar(b *testing.B) {
	s1 := "Hugo is cool and " + strings.Repeat("fun ", 10) + "!"
	s2 := "Hugo is cool and " + strings.Repeat("cool ", 10) + "!"

	for b.Loop() {
		howSimilar(s1, s2)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{4*time.Second + 342*time.Millisecond, "4.34  s"}, // additional spacing between value and unit
		{170*time.Millisecond + 289*time.Microsecond, "170.29 ms"},
		{16*time.Microsecond + 90*time.Nanosecond, "16.09 µs"},
		{147 * time.Nanosecond, "147.00 ns"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.duration)
		if got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
		}
	}
}
