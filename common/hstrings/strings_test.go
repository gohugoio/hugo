// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"reflect"
	"regexp"
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

func TestUniqueStrings(t *testing.T) {
	in := []string{"a", "b", "a", "b", "c", "", "a", "", "d"}
	output := UniqueStrings(in)
	expected := []string{"a", "b", "c", "", "d"}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %#v, got %#v\n", expected, output)
	}
}

func TestUniqueStringsReuse(t *testing.T) {
	in := []string{"a", "b", "a", "b", "c", "", "a", "", "d"}
	output := UniqueStringsReuse(in)
	expected := []string{"a", "b", "c", "", "d"}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %#v, got %#v\n", expected, output)
	}
}

func TestUniqueStringsSorted(t *testing.T) {
	c := qt.New(t)
	in := []string{"a", "a", "b", "c", "b", "", "a", "", "d"}
	output := UniqueStringsSorted(in)
	expected := []string{"", "a", "b", "c", "d"}
	c.Assert(output, qt.DeepEquals, expected)
	c.Assert(UniqueStringsSorted(nil), qt.IsNil)
}

// Note that these cannot use b.Loop() because of golang/go#27217.
func BenchmarkUniqueStrings(b *testing.B) {
	input := []string{"a", "b", "d", "e", "d", "h", "a", "i"}

	b.Run("Safe", func(b *testing.B) {
		for b.Loop() {
			result := UniqueStrings(input)
			if len(result) != 6 {
				b.Fatalf("invalid count: %d", len(result))
			}
		}
	})

	b.Run("Reuse slice", func(b *testing.B) {
		inputs := make([][]string, b.N)
		for i := 0; i < b.N; i++ {
			inputc := make([]string, len(input))
			copy(inputc, input)
			inputs[i] = inputc
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			inputc := inputs[i]

			result := UniqueStringsReuse(inputc)
			if len(result) != 6 {
				b.Fatalf("invalid count: %d", len(result))
			}
		}
	})

	b.Run("Reuse slice sorted", func(b *testing.B) {
		inputs := make([][]string, b.N)
		for i := 0; i < b.N; i++ {
			inputc := make([]string, len(input))
			copy(inputc, input)
			inputs[i] = inputc
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			inputc := inputs[i]

			result := UniqueStringsSorted(inputc)
			if len(result) != 6 {
				b.Fatalf("invalid count: %d", len(result))
			}
		}
	})
}

func BenchmarkGetOrCompileRegexp(b *testing.B) {
	for b.Loop() {
		GetOrCompileRegexp(`\d+`)
	}
}

func BenchmarkCompileRegexp(b *testing.B) {
	for b.Loop() {
		regexp.MustCompile(`\d+`)
	}
}
