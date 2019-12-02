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

package helpers

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/gohugoio/hugo/common/loggers"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestResolveMarkup(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	spec, err := NewContentSpec(cfg, loggers.NewErrorLogger(), afero.NewMemMapFs())
	c.Assert(err, qt.IsNil)

	for i, this := range []struct {
		in     string
		expect string
	}{
		{"md", "markdown"},
		{"markdown", "markdown"},
		{"mdown", "markdown"},
		{"asciidoc", "asciidoc"},
		{"adoc", "asciidoc"},
		{"ad", "asciidoc"},
		{"rst", "rst"},
		{"pandoc", "pandoc"},
		{"pdc", "pandoc"},
		{"mmark", "mmark"},
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
		result := FirstUpper(this.in)
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
		result := HasStringsPrefix(this.s, this.prefix)
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
		result := HasStringsSuffix(this.s, this.suffix)
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
		res := SliceToLower(test.value)
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
		result := ReaderContains(strings.NewReader(this.v1), this.v2)
		if result != this.expect {
			t.Errorf("[%d] got %t but expected %t", i, result, this.expect)
		}
	}

	c.Assert(ReaderContains(nil, []byte("a")), qt.Equals, false)
	c.Assert(ReaderContains(nil, nil), qt.Equals, false)
}

func TestGetTitleFunc(t *testing.T) {
	title := "somewhere over the rainbow"
	c := qt.New(t)

	c.Assert(GetTitleFunc("go")(title), qt.Equals, "Somewhere Over The Rainbow")
	c.Assert(GetTitleFunc("chicago")(title), qt.Equals, "Somewhere over the Rainbow")
	c.Assert(GetTitleFunc("Chicago")(title), qt.Equals, "Somewhere over the Rainbow")
	c.Assert(GetTitleFunc("ap")(title), qt.Equals, "Somewhere Over the Rainbow")
	c.Assert(GetTitleFunc("ap")(title), qt.Equals, "Somewhere Over the Rainbow")
	c.Assert(GetTitleFunc("")(title), qt.Equals, "Somewhere Over the Rainbow")
	c.Assert(GetTitleFunc("unknown")(title), qt.Equals, "Somewhere Over the Rainbow")

}

func BenchmarkReaderContains(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i, this := range containsBenchTestData {
			result := ReaderContains(strings.NewReader(this.v1), this.v2)
			if result != this.expect {
				b.Errorf("[%d] got %t but expected %t", i, result, this.expect)
			}
		}
	}
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

func TestFindAvailablePort(t *testing.T) {
	c := qt.New(t)
	addr, err := FindAvailablePort()
	c.Assert(err, qt.IsNil)
	c.Assert(addr, qt.Not(qt.IsNil))
	c.Assert(addr.Port > 0, qt.Equals, true)
}

func TestFastMD5FromFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := afero.WriteFile(fs, "small.txt", []byte("abc"), 0777); err != nil {
		t.Fatal(err)
	}

	if err := afero.WriteFile(fs, "small2.txt", []byte("abd"), 0777); err != nil {
		t.Fatal(err)
	}

	if err := afero.WriteFile(fs, "bigger.txt", []byte(strings.Repeat("a bc d e", 100)), 0777); err != nil {
		t.Fatal(err)
	}

	if err := afero.WriteFile(fs, "bigger2.txt", []byte(strings.Repeat("c d e f g", 100)), 0777); err != nil {
		t.Fatal(err)
	}

	c := qt.New(t)

	sf1, err := fs.Open("small.txt")
	c.Assert(err, qt.IsNil)
	sf2, err := fs.Open("small2.txt")
	c.Assert(err, qt.IsNil)

	bf1, err := fs.Open("bigger.txt")
	c.Assert(err, qt.IsNil)
	bf2, err := fs.Open("bigger2.txt")
	c.Assert(err, qt.IsNil)

	defer sf1.Close()
	defer sf2.Close()
	defer bf1.Close()
	defer bf2.Close()

	m1, err := MD5FromFileFast(sf1)
	c.Assert(err, qt.IsNil)
	c.Assert(m1, qt.Equals, "e9c8989b64b71a88b4efb66ad05eea96")

	m2, err := MD5FromFileFast(sf2)
	c.Assert(err, qt.IsNil)
	c.Assert(m2, qt.Not(qt.Equals), m1)

	m3, err := MD5FromFileFast(bf1)
	c.Assert(err, qt.IsNil)
	c.Assert(m3, qt.Not(qt.Equals), m2)

	m4, err := MD5FromFileFast(bf2)
	c.Assert(err, qt.IsNil)
	c.Assert(m4, qt.Not(qt.Equals), m3)

	m5, err := MD5FromReader(bf2)
	c.Assert(err, qt.IsNil)
	c.Assert(m5, qt.Not(qt.Equals), m4)
}

func BenchmarkMD5FromFileFast(b *testing.B) {
	fs := afero.NewMemMapFs()

	for _, full := range []bool{false, true} {
		b.Run(fmt.Sprintf("full=%t", full), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				if err := afero.WriteFile(fs, "file.txt", []byte(strings.Repeat("1234567890", 2000)), 0777); err != nil {
					b.Fatal(err)
				}
				f, err := fs.Open("file.txt")
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
				if full {
					if _, err := MD5FromReader(f); err != nil {
						b.Fatal(err)
					}
				} else {
					if _, err := MD5FromFileFast(f); err != nil {
						b.Fatal(err)
					}
				}
				f.Close()
			}
		})
	}

}

func BenchmarkUniqueStrings(b *testing.B) {
	input := []string{"a", "b", "d", "e", "d", "h", "a", "i"}

	b.Run("Safe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := UniqueStrings(input)
			if len(result) != 6 {
				b.Fatal(fmt.Sprintf("invalid count: %d", len(result)))
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

			result := UniqueStringsReuse(inputc)
			if len(result) != 6 {
				b.Fatal(fmt.Sprintf("invalid count: %d", len(result)))
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

			result := UniqueStringsSorted(inputc)
			if len(result) != 6 {
				b.Fatal(fmt.Sprintf("invalid count: %d", len(result)))
			}
		}
	})

}

func TestHashString(t *testing.T) {
	c := qt.New(t)

	c.Assert(HashString("a", "b"), qt.Equals, "2712570657419664240")
	c.Assert(HashString("ab"), qt.Equals, "590647783936702392")
}
