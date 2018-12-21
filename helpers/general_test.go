// Copyright 2018 The Hugo Authors. All rights reserved.
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

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuessType(t *testing.T) {
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
		{"excel", "unknown"},
	} {
		result := GuessType(this.in)
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

func TestReaderContains(t *testing.T) {
	for i, this := range append(containsBenchTestData, containsAdditionalTestData...) {
		result := ReaderContains(strings.NewReader(this.v1), this.v2)
		if result != this.expect {
			t.Errorf("[%d] got %t but expected %t", i, result, this.expect)
		}
	}

	assert.False(t, ReaderContains(nil, []byte("a")))
	assert.False(t, ReaderContains(nil, nil))
}

func TestGetTitleFunc(t *testing.T) {
	title := "somewhere over the rainbow"
	assert := require.New(t)

	assert.Equal("Somewhere Over The Rainbow", GetTitleFunc("go")(title))
	assert.Equal("Somewhere over the Rainbow", GetTitleFunc("chicago")(title), "Chicago style")
	assert.Equal("Somewhere over the Rainbow", GetTitleFunc("Chicago")(title), "Chicago style")
	assert.Equal("Somewhere Over the Rainbow", GetTitleFunc("ap")(title), "AP style")
	assert.Equal("Somewhere Over the Rainbow", GetTitleFunc("ap")(title), "AP style")
	assert.Equal("Somewhere Over the Rainbow", GetTitleFunc("")(title), "AP style")
	assert.Equal("Somewhere Over the Rainbow", GetTitleFunc("unknown")(title), "AP style")

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

func TestFindAvailablePort(t *testing.T) {
	addr, err := FindAvailablePort()
	assert.Nil(t, err)
	assert.NotNil(t, addr)
	assert.True(t, addr.Port > 0)
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

	req := require.New(t)

	sf1, err := fs.Open("small.txt")
	req.NoError(err)
	sf2, err := fs.Open("small2.txt")
	req.NoError(err)

	bf1, err := fs.Open("bigger.txt")
	req.NoError(err)
	bf2, err := fs.Open("bigger2.txt")
	req.NoError(err)

	defer sf1.Close()
	defer sf2.Close()
	defer bf1.Close()
	defer bf2.Close()

	m1, err := MD5FromFileFast(sf1)
	req.NoError(err)
	req.Equal("e9c8989b64b71a88b4efb66ad05eea96", m1)

	m2, err := MD5FromFileFast(sf2)
	req.NoError(err)
	req.NotEqual(m1, m2)

	m3, err := MD5FromFileFast(bf1)
	req.NoError(err)
	req.NotEqual(m2, m3)

	m4, err := MD5FromFileFast(bf2)
	req.NoError(err)
	req.NotEqual(m3, m4)

	m5, err := MD5FromReader(bf2)
	req.NoError(err)
	req.NotEqual(m4, m5)
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
