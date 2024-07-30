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

package helpers_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
)

func TestResolveMarkup(t *testing.T) {
	spec := newTestContentSpec(nil)

	for i, this := range []struct {
		in     string
		expect string
	}{
		{"md", "markdown"},
		{"markdown", "markdown"},
		{"mdown", "markdown"},
		{"asciidocext", "asciidoc"},
		{"adoc", "asciidoc"},
		{"ad", "asciidoc"},
		{"rst", "rst"},
		{"pandoc", "pandoc"},
		{"pdc", "pandoc"},
		{"html", "html"},
		{"htm", "html"},
		{"org", "org"},
		{"excel", ""},
	} {
		result := spec.ResolveMarkup(this.in)
		if result != this.expect {
			t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
		}
	}
}

func TestFirstUpper(t *testing.T) {
	for i, this := range []struct {
		in     string
		expect string
	}{
		{"foo", "Foo"},
		{"foo bar", "Foo bar"},
		{"Foo Bar", "Foo Bar"},
		{"", ""},
		{"å", "Å"},
	} {
		result := helpers.FirstUpper(this.in)
		if result != this.expect {
			t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
		}
	}
}

func TestHasStringsPrefix(t *testing.T) {
	for i, this := range []struct {
		s      []string
		prefix []string
		expect bool
	}{
		{[]string{"a"}, []string{"a"}, true},
		{[]string{}, []string{}, true},
		{[]string{"a", "b", "c"}, []string{"a", "b"}, true},
		{[]string{"d", "a", "b", "c"}, []string{"a", "b"}, false},
		{[]string{"abra", "ca", "dabra"}, []string{"abra", "ca"}, true},
		{[]string{"abra", "ca"}, []string{"abra", "ca", "dabra"}, false},
	} {
		result := helpers.HasStringsPrefix(this.s, this.prefix)
		if result != this.expect {
			t.Fatalf("[%d] got %t but expected %t", i, result, this.expect)
		}
	}
}

func TestHasStringsSuffix(t *testing.T) {
	for i, this := range []struct {
		s      []string
		suffix []string
		expect bool
	}{
		{[]string{"a"}, []string{"a"}, true},
		{[]string{}, []string{}, true},
		{[]string{"a", "b", "c"}, []string{"b", "c"}, true},
		{[]string{"abra", "ca", "dabra"}, []string{"abra", "ca"}, false},
		{[]string{"abra", "ca", "dabra"}, []string{"ca", "dabra"}, true},
	} {
		result := helpers.HasStringsSuffix(this.s, this.suffix)
		if result != this.expect {
			t.Fatalf("[%d] got %t but expected %t", i, result, this.expect)
		}
	}
}

var containsTestText = (`На берегу пустынных волн
Стоял он, дум великих полн,
И вдаль глядел. Пред ним широко
Река неслася; бедный чёлн
По ней стремился одиноко.
По мшистым, топким берегам
Чернели избы здесь и там,
Приют убогого чухонца;
И лес, неведомый лучам
В тумане спрятанного солнца,
Кругом шумел.

Τη γλώσσα μου έδωσαν ελληνική
το σπίτι φτωχικό στις αμμουδιές του Ομήρου.
Μονάχη έγνοια η γλώσσα μου στις αμμουδιές του Ομήρου.

από το Άξιον Εστί
του Οδυσσέα Ελύτη

Sîne klâwen durh die wolken sint geslagen,
er stîget ûf mit grôzer kraft,
ich sih in grâwen tägelîch als er wil tagen,
den tac, der im geselleschaft
erwenden wil, dem werden man,
den ich mit sorgen în verliez.
ich bringe in hinnen, ob ich kan.
sîn vil manegiu tugent michz leisten hiez.
`)

var containsBenchTestData = []struct {
	v1     string
	v2     []byte
	expect bool
}{
	{"abc", []byte("a"), true},
	{"abc", []byte("b"), true},
	{"abcdefg", []byte("efg"), true},
	{"abc", []byte("d"), false},
	{containsTestText, []byte("стремился"), true},
	{containsTestText, []byte(containsTestText[10:80]), true},
	{containsTestText, []byte(containsTestText[100:111]), true},
	{containsTestText, []byte(containsTestText[len(containsTestText)-100 : len(containsTestText)-10]), true},
	{containsTestText, []byte(containsTestText[len(containsTestText)-20:]), true},
	{containsTestText, []byte("notfound"), false},
}

// some corner cases
var containsAdditionalTestData = []struct {
	v1     string
	v2     []byte
	expect bool
}{
	{"", nil, false},
	{"", []byte("a"), false},
	{"a", []byte(""), false},
	{"", []byte(""), false},
}

func TestSliceToLower(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value    []string
		expected []string
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "B", "c"}, []string{"a", "b", "c"}},
		{[]string{"A", "B", "C"}, []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		res := helpers.SliceToLower(test.value)
		for i, val := range res {
			if val != test.expected[i] {
				t.Errorf("Case mismatch. Expected %s, got %s", test.expected[i], res[i])
			}
		}
	}
}

func TestReaderContains(t *testing.T) {
	c := qt.New(t)
	for i, this := range append(containsBenchTestData, containsAdditionalTestData...) {
		result := helpers.ReaderContains(strings.NewReader(this.v1), this.v2)
		if result != this.expect {
			t.Errorf("[%d] got %t but expected %t", i, result, this.expect)
		}
	}

	c.Assert(helpers.ReaderContains(nil, []byte("a")), qt.Equals, false)
	c.Assert(helpers.ReaderContains(nil, nil), qt.Equals, false)
}

func TestGetTitleFunc(t *testing.T) {
	title := "somewhere over the Rainbow"
	c := qt.New(t)

	c.Assert(helpers.GetTitleFunc("go")(title), qt.Equals, "Somewhere Over The Rainbow")
	c.Assert(helpers.GetTitleFunc("chicago")(title), qt.Equals, "Somewhere over the Rainbow")
	c.Assert(helpers.GetTitleFunc("Chicago")(title), qt.Equals, "Somewhere over the Rainbow")
	c.Assert(helpers.GetTitleFunc("ap")(title), qt.Equals, "Somewhere Over the Rainbow")
	c.Assert(helpers.GetTitleFunc("ap")(title), qt.Equals, "Somewhere Over the Rainbow")
	c.Assert(helpers.GetTitleFunc("")(title), qt.Equals, "Somewhere Over the Rainbow")
	c.Assert(helpers.GetTitleFunc("unknown")(title), qt.Equals, "Somewhere Over the Rainbow")
	c.Assert(helpers.GetTitleFunc("none")(title), qt.Equals, title)
	c.Assert(helpers.GetTitleFunc("firstupper")(title), qt.Equals, "Somewhere over the Rainbow")
}

func BenchmarkReaderContains(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i, this := range containsBenchTestData {
			result := helpers.ReaderContains(strings.NewReader(this.v1), this.v2)
			if result != this.expect {
				b.Errorf("[%d] got %t but expected %t", i, result, this.expect)
			}
		}
	}
}

func TestUniqueStrings(t *testing.T) {
	in := []string{"a", "b", "a", "b", "c", "", "a", "", "d"}
	output := helpers.UniqueStrings(in)
	expected := []string{"a", "b", "c", "", "d"}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %#v, got %#v\n", expected, output)
	}
}

func TestUniqueStringsReuse(t *testing.T) {
	in := []string{"a", "b", "a", "b", "c", "", "a", "", "d"}
	output := helpers.UniqueStringsReuse(in)
	expected := []string{"a", "b", "c", "", "d"}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %#v, got %#v\n", expected, output)
	}
}

func TestUniqueStringsSorted(t *testing.T) {
	c := qt.New(t)
	in := []string{"a", "a", "b", "c", "b", "", "a", "", "d"}
	output := helpers.UniqueStringsSorted(in)
	expected := []string{"", "a", "b", "c", "d"}
	c.Assert(output, qt.DeepEquals, expected)
	c.Assert(helpers.UniqueStringsSorted(nil), qt.IsNil)
}

func BenchmarkUniqueStrings(b *testing.B) {
	input := []string{"a", "b", "d", "e", "d", "h", "a", "i"}

	b.Run("Safe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := helpers.UniqueStrings(input)
			if len(result) != 6 {
				b.Fatalf("invalid count: %d", len(result))
			}
		}
	})

	b.Run("Reuse slice", func(b *testing.B) {
		b.StopTimer()
		inputs := make([][]string, b.N)
		for i := 0; i < b.N; i++ {
			inputc := make([]string, len(input))
			copy(inputc, input)
			inputs[i] = inputc
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			inputc := inputs[i]

			result := helpers.UniqueStringsReuse(inputc)
			if len(result) != 6 {
				b.Fatalf("invalid count: %d", len(result))
			}
		}
	})

	b.Run("Reuse slice sorted", func(b *testing.B) {
		b.StopTimer()
		inputs := make([][]string, b.N)
		for i := 0; i < b.N; i++ {
			inputc := make([]string, len(input))
			copy(inputc, input)
			inputs[i] = inputc
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			inputc := inputs[i]

			result := helpers.UniqueStringsSorted(inputc)
			if len(result) != 6 {
				b.Fatalf("invalid count: %d", len(result))
			}
		}
	})
}
