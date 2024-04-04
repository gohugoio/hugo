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

package strings

import (
	"html/template"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/cast"
)

var ns = New(&deps.Deps{
	Conf: testconfig.GetTestConfig(nil, nil),
})

type tstNoStringer struct{}

func TestChomp(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect any
	}{
		{"\n a\n", "\n a"},
		{"\n a\n\n", "\n a"},
		{"\n a\r\n", "\n a"},
		{"\n a\n\r\n", "\n a"},
		{"\n a\r\r", "\n a"},
		{"\n a\r", "\n a"},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Chomp(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)

		// repeat the check with template.HTML input
		result, err = ns.Chomp(template.HTML(cast.ToString(test.s)))
		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, template.HTML(cast.ToString(test.expect)))
	}
}

func TestContains(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		substr any
		expect bool
		isErr  bool
	}{
		{"", "", true, false},
		{"123", "23", true, false},
		{"123", "234", false, false},
		{"123", "", true, false},
		{"", "a", false, false},
		{123, "23", true, false},
		{123, "234", false, false},
		{123, "", true, false},
		{template.HTML("123"), []byte("23"), true, false},
		{template.HTML("123"), []byte("234"), false, false},
		{template.HTML("123"), []byte(""), true, false},
		// errors
		{"", tstNoStringer{}, false, true},
		{tstNoStringer{}, "", false, true},
	} {

		result, err := ns.Contains(test.s, test.substr)

		if test.isErr {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestContainsAny(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		substr any
		expect bool
		isErr  bool
	}{
		{"", "", false, false},
		{"", "1", false, false},
		{"", "123", false, false},
		{"1", "", false, false},
		{"1", "1", true, false},
		{"111", "1", true, false},
		{"123", "789", false, false},
		{"123", "729", true, false},
		{"a☺b☻c☹d", "uvw☻xyz", true, false},
		{1, "", false, false},
		{1, "1", true, false},
		{111, "1", true, false},
		{123, "789", false, false},
		{123, "729", true, false},
		{[]byte("123"), template.HTML("789"), false, false},
		{[]byte("123"), template.HTML("729"), true, false},
		{[]byte("a☺b☻c☹d"), template.HTML("uvw☻xyz"), true, false},
		// errors
		{"", tstNoStringer{}, false, true},
		{tstNoStringer{}, "", false, true},
	} {

		result, err := ns.ContainsAny(test.s, test.substr)

		if test.isErr {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestContainsNonSpace(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect bool
		isErr  bool
	}{
		{"", false, false},
		{" ", false, false},
		{"        ", false, false},
		{"\t", false, false},
		{"\r", false, false},
		{"a", true, false},
		{"    a", true, false},
		{"a\n", true, false},
		// error
		{tstNoStringer{}, false, true},
	} {

		result, err := ns.ContainsNonSpace(test.s)

		if test.isErr {
			c.Assert(err, qt.IsNotNil)
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)

	}
}

func TestCountRunes(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect any
	}{
		{"foo bar", 6},
		{"旁边", 2},
		{`<div class="test">旁边</div>`, 2},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.CountRunes(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestRuneCount(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect any
	}{
		{"foo bar", 7},
		{"旁边", 2},
		{`<div class="test">旁边</div>`, 26},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.RuneCount(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestCountWords(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect any
	}{
		{"Do Be Do Be Do", 5},
		{"旁边", 2},
		{`<div class="test">旁边</div>`, 2},
		{"Here's to you...", 3},
		{"Here’s to you...", 3},
		{"Here’s to you…", 3},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.CountWords(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		prefix any
		expect any
		isErr  bool
	}{
		{"abcd", "ab", true, false},
		{"abcd", "cd", false, false},
		{template.HTML("abcd"), "ab", true, false},
		{template.HTML("abcd"), "cd", false, false},
		{template.HTML("1234"), 12, true, false},
		{template.HTML("1234"), 34, false, false},
		{[]byte("abcd"), "ab", true, false},
		// errors
		{"", tstNoStringer{}, false, true},
		{tstNoStringer{}, "", false, true},
	} {

		result, err := ns.HasPrefix(test.s, test.prefix)

		if test.isErr {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestHasSuffix(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		suffix any
		expect any
		isErr  bool
	}{
		{"abcd", "cd", true, false},
		{"abcd", "ab", false, false},
		{template.HTML("abcd"), "cd", true, false},
		{template.HTML("abcd"), "ab", false, false},
		{template.HTML("1234"), 34, true, false},
		{template.HTML("1234"), 12, false, false},
		{[]byte("abcd"), "cd", true, false},
		// errors
		{"", tstNoStringer{}, false, true},
		{tstNoStringer{}, "", false, true},
	} {

		result, err := ns.HasSuffix(test.s, test.suffix)

		if test.isErr {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestReplace(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		old    any
		new    any
		limit  any
		expect any
	}{
		{"aab", "a", "b", nil, "bbb"},
		{"11a11", 1, 2, nil, "22a22"},
		{12345, 1, 2, nil, "22345"},
		{"aab", "a", "b", 1, "bab"},
		{"11a11", 1, 2, 2, "22a11"},
		// errors
		{tstNoStringer{}, "a", "b", nil, false},
		{"a", tstNoStringer{}, "b", nil, false},
		{"a", "b", tstNoStringer{}, nil, false},
	} {

		var (
			result string
			err    error
		)

		if test.limit != nil {
			result, err = ns.Replace(test.s, test.old, test.new, test.limit)
		} else {
			result, err = ns.Replace(test.s, test.old, test.new)
		}

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestSliceString(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var err error
	for _, test := range []struct {
		v1     any
		v2     any
		v3     any
		expect any
	}{
		{"abc", 1, 2, "b"},
		{"abc", 1, 3, "bc"},
		{"abcdef", 1, int8(3), "bc"},
		{"abcdef", 1, int16(3), "bc"},
		{"abcdef", 1, int32(3), "bc"},
		{"abcdef", 1, int64(3), "bc"},
		{"abc", 0, 1, "a"},
		{"abcdef", nil, nil, "abcdef"},
		{"abcdef", 0, 6, "abcdef"},
		{"abcdef", 0, 2, "ab"},
		{"abcdef", 2, nil, "cdef"},
		{"abcdef", int8(2), nil, "cdef"},
		{"abcdef", int16(2), nil, "cdef"},
		{"abcdef", int32(2), nil, "cdef"},
		{"abcdef", int64(2), nil, "cdef"},
		{123, 1, 3, "23"},
		{"abcdef", 6, nil, false},
		{"abcdef", 4, 7, false},
		{"abcdef", -1, nil, false},
		{"abcdef", -1, 7, false},
		{"abcdef", 1, -1, false},
		{tstNoStringer{}, 0, 1, false},
		{"ĀĀĀ", 0, 1, "Ā"}, // issue #1333
		{"a", t, nil, false},
		{"a", 1, t, false},
	} {

		var result string
		if test.v2 == nil {
			result, err = ns.SliceString(test.v1)
		} else if test.v3 == nil {
			result, err = ns.SliceString(test.v1, test.v2)
		} else {
			result, err = ns.SliceString(test.v1, test.v2, test.v3)
		}

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}

	// Too many arguments
	_, err = ns.SliceString("a", 1, 2, 3)
	if err == nil {
		t.Errorf("Should have errored")
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		v1     any
		v2     string
		expect any
	}{
		{"a, b", ", ", []string{"a", "b"}},
		{"a & b & c", " & ", []string{"a", "b", "c"}},
		{"http://example.com", "http://", []string{"", "example.com"}},
		{123, "2", []string{"1", "3"}},
		{tstNoStringer{}, ",", false},
	} {

		result, err := ns.Split(test.v1, test.v2)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.DeepEquals, test.expect)
	}
}

func TestSubstr(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var err error
	for _, test := range []struct {
		v1     any
		v2     any
		v3     any
		expect any
	}{
		{"abc", 1, 2, "bc"},
		{"abc", 0, 1, "a"},
		{"abcdef", 0, 0, ""},
		{"abcdef", 1, 0, ""},
		{"abcdef", -1, 0, ""},
		{"abcdef", -1, 2, "f"},
		{"abcdef", -3, 3, "def"},
		{"abcdef", -1, nil, "f"},
		{"abcdef", -2, nil, "ef"},
		{"abcdef", -3, 1, "d"},
		{"abcdef", 0, -1, "abcde"},
		{"abcdef", 2, -1, "cde"},
		{"abcdef", 4, -4, ""},
		{"abcdef", 7, 1, ""},
		{"abcdef", 6, nil, ""},
		{"abcdef", 1, 100, "bcdef"},
		{"abcdef", -100, 3, "abc"},
		{"abcdef", -3, -1, "de"},
		{"abcdef", 2, nil, "cdef"},
		{"abcdef", int8(2), nil, "cdef"},
		{"abcdef", int16(2), nil, "cdef"},
		{"abcdef", int32(2), nil, "cdef"},
		{"abcdef", int64(2), nil, "cdef"},
		{"abcdef", 2, int8(3), "cde"},
		{"abcdef", 2, int16(3), "cde"},
		{"abcdef", 2, int32(3), "cde"},
		{"abcdef", 2, int64(3), "cde"},
		{123, 1, 3, "23"},
		{1.2e3, 0, 4, "1200"},
		{tstNoStringer{}, 0, 1, false},
		{"abcdef", 2.0, nil, "cdef"},
		{"abcdef", 2.0, 2, "cd"},
		{"abcdef", 2, 2.0, "cd"},
		{"ĀĀĀ", 1, 2, "ĀĀ"}, // # issue 1333
		{"abcdef", "doo", nil, false},
		{"abcdef", "doo", "doo", false},
		{"abcdef", 1, "doo", false},
		{"", 0, nil, ""},
	} {

		var result string

		if test.v3 == nil {
			result, err = ns.Substr(test.v1, test.v2)
		} else {
			result, err = ns.Substr(test.v1, test.v2, test.v3)
		}

		if b, ok := test.expect.(bool); ok && !b {
			c.Check(err, qt.Not(qt.IsNil), qt.Commentf("%v", test))
			continue
		}

		c.Assert(err, qt.IsNil, qt.Commentf("%v", test))
		c.Check(result, qt.Equals, test.expect, qt.Commentf("%v", test))
	}

	_, err = ns.Substr("abcdef")
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.Substr("abcdef", 1, 2, 3)
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestTitle(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect any
	}{
		{"test", "Test"},
		{template.HTML("hypertext"), "Hypertext"},
		{[]byte("bytes"), "Bytes"},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.Title(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestToLower(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect any
	}{
		{"TEST", "test"},
		{template.HTML("LoWeR"), "lower"},
		{[]byte("BYTES"), "bytes"},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.ToLower(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestToUpper(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		expect any
	}{
		{"test", "TEST"},
		{template.HTML("UpPeR"), "UPPER"},
		{[]byte("bytes"), "BYTES"},
		// errors
		{tstNoStringer{}, false},
	} {

		result, err := ns.ToUpper(test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestTrim(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		cutset any
		expect any
	}{
		{"abba", "a", "bb"},
		{"abba", "ab", ""},
		{"<tag>", "<>", "tag"},
		{`"quote"`, `"`, "quote"},
		{1221, "1", "22"},
		{1221, "12", ""},
		{template.HTML("<tag>"), "<>", "tag"},
		{[]byte("<tag>"), "<>", "tag"},
		// errors
		{"", tstNoStringer{}, false},
		{tstNoStringer{}, "", false},
	} {

		result, err := ns.Trim(test.s, test.cutset)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestTrimLeft(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		cutset any
		expect any
	}{
		{"abba", "a", "bba"},
		{"abba", "ab", ""},
		{"<tag>", "<>", "tag>"},
		{`"quote"`, `"`, `quote"`},
		{1221, "1", "221"},
		{1221, "12", ""},
		{"007", "0", "7"},
		{template.HTML("<tag>"), "<>", "tag>"},
		{[]byte("<tag>"), "<>", "tag>"},
		// errors
		{"", tstNoStringer{}, false},
		{tstNoStringer{}, "", false},
	} {

		result, err := ns.TrimLeft(test.cutset, test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestTrimPrefix(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		prefix any
		expect any
	}{
		{"aabbaa", "a", "abbaa"},
		{"aabb", "b", "aabb"},
		{1234, "12", "34"},
		{1234, "34", "1234"},
		// errors
		{"", tstNoStringer{}, false},
		{tstNoStringer{}, "", false},
	} {

		result, err := ns.TrimPrefix(test.prefix, test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestTrimRight(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		cutset any
		expect any
	}{
		{"abba", "a", "abb"},
		{"abba", "ab", ""},
		{"<tag>", "<>", "<tag"},
		{`"quote"`, `"`, `"quote`},
		{1221, "1", "122"},
		{1221, "12", ""},
		{"007", "0", "007"},
		{template.HTML("<tag>"), "<>", "<tag"},
		{[]byte("<tag>"), "<>", "<tag"},
		// errors
		{"", tstNoStringer{}, false},
		{tstNoStringer{}, "", false},
	} {

		result, err := ns.TrimRight(test.cutset, test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestTrimSuffix(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		suffix any
		expect any
	}{
		{"aabbaa", "a", "aabba"},
		{"aabb", "b", "aab"},
		{1234, "12", "1234"},
		{1234, "34", "12"},
		// errors
		{"", tstNoStringer{}, false},
		{tstNoStringer{}, "", false},
	} {

		result, err := ns.TrimSuffix(test.suffix, test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestRepeat(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		s      any
		n      any
		expect any
	}{
		{"yo", "2", "yoyo"},
		{"~", "16", "~~~~~~~~~~~~~~~~"},
		{"<tag>", "0", ""},
		{"yay", "1", "yay"},
		{1221, "1", "1221"},
		{1221, 2, "12211221"},
		{template.HTML("<tag>"), "2", "<tag><tag>"},
		{[]byte("<tag>"), 2, "<tag><tag>"},
		// errors
		{"", tstNoStringer{}, false},
		{tstNoStringer{}, "", false},
		{"ab", -1, false},
	} {

		result, err := ns.Repeat(test.n, test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestDiff(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		oldname string
		old     any
		newname string
		new     any
		expect  any
	}{
		{"old", "foo\n", "new", "bar\n", "diff old new\n--- old\n+++ new\n@@ -1,1 +1,1 @@\n-foo\n+bar\n"},
		{"old", "foo\n", "new", "foo\n", ""},
		{"old", "foo\n", "new", "", "diff old new\n--- old\n+++ new\n@@ -1,1 +0,0 @@\n-foo\n"},
		{"old", "foo\n", "new", nil, "diff old new\n--- old\n+++ new\n@@ -1,1 +0,0 @@\n-foo\n"},
		{"old", "", "new", "", ""},
		// errors
		{"old", tstNoStringer{}, "new", "foo", false},
		{"old", "foo", "new", tstNoStringer{}, false},
	} {

		result, err := ns.Diff(test.oldname, test.old, test.newname, test.new)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)

	}
}
