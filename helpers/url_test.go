// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
)

func TestURLize(t *testing.T) {

	v := newTestCfg()
	l := langs.NewDefaultLanguage(v)
	p, _ := NewPathSpec(hugofs.NewMem(v), l, nil)

	tests := []struct {
		input    string
		expected string
	}{
		{"  foo bar  ", "foo-bar"},
		{"foo.bar/foo_bar-foo", "foo.bar/foo_bar-foo"},
		{"foo,bar:foobar", "foobarfoobar"},
		{"foo/bar.html", "foo/bar.html"},
		{"трям/трям", "%D1%82%D1%80%D1%8F%D0%BC/%D1%82%D1%80%D1%8F%D0%BC"},
		{"100%-google", "100-google"},
	}

	for _, test := range tests {
		output := p.URLize(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestAbsURL(t *testing.T) {
	for _, defaultInSubDir := range []bool{true, false} {
		for _, addLanguage := range []bool{true, false} {
			for _, m := range []bool{true, false} {
				for _, l := range []string{"en", "fr"} {
					doTestAbsURL(t, defaultInSubDir, addLanguage, m, l)
				}
			}
		}
	}
}

func doTestAbsURL(t *testing.T, defaultInSubDir, addLanguage, multilingual bool, lang string) {
	v := newTestCfg()
	v.Set("multilingual", multilingual)
	v.Set("defaultContentLanguage", "en")
	v.Set("defaultContentLanguageInSubdir", defaultInSubDir)

	tests := []struct {
		input    string
		baseURL  string
		expected string
	}{
		{"/test/foo", "http://base/", "http://base/MULTItest/foo"},
		{"/" + lang + "/test/foo", "http://base/", "http://base/" + lang + "/test/foo"},
		{"", "http://base/ace/", "http://base/ace/MULTI"},
		{"/test/2/foo/", "http://base", "http://base/MULTItest/2/foo/"},
		{"http://abs", "http://base/", "http://abs"},
		{"schema://abs", "http://base/", "schema://abs"},
		{"//schemaless", "http://base/", "//schemaless"},
		{"test/2/foo/", "http://base/path", "http://base/path/MULTItest/2/foo/"},
		{lang + "/test/2/foo/", "http://base/path", "http://base/path/" + lang + "/test/2/foo/"},
		{"/test/2/foo/", "http://base/path", "http://base/MULTItest/2/foo/"},
		{"http//foo", "http://base/path", "http://base/path/MULTIhttp/foo"},
	}

	for _, test := range tests {
		v.Set("baseURL", test.baseURL)
		v.Set("contentDir", "content")
		l := langs.NewLanguage(lang, v)
		p, _ := NewPathSpec(hugofs.NewMem(v), l, nil)

		output := p.AbsURL(test.input, addLanguage)
		expected := test.expected
		if multilingual && addLanguage {
			if !defaultInSubDir && lang == "en" {
				expected = strings.Replace(expected, "MULTI", "", 1)
			} else {
				expected = strings.Replace(expected, "MULTI", lang+"/", 1)
			}

		} else {
			expected = strings.Replace(expected, "MULTI", "", 1)
		}
		if output != expected {
			t.Fatalf("Expected %#v, got %#v\n", expected, output)
		}
	}
}

func TestIsAbsURL(t *testing.T) {
	c := qt.New(t)

	for _, this := range []struct {
		a string
		b bool
	}{
		{"http://gohugo.io", true},
		{"https://gohugo.io", true},
		{"//gohugo.io", true},
		{"http//gohugo.io", false},
		{"/content", false},
		{"content", false},
	} {
		c.Assert(IsAbsURL(this.a) == this.b, qt.Equals, true)
	}
}

func TestRelURL(t *testing.T) {
	for _, defaultInSubDir := range []bool{true, false} {
		for _, addLanguage := range []bool{true, false} {
			for _, m := range []bool{true, false} {
				for _, l := range []string{"en", "fr"} {
					doTestRelURL(t, defaultInSubDir, addLanguage, m, l)
				}
			}
		}
	}
}

func doTestRelURL(t *testing.T, defaultInSubDir, addLanguage, multilingual bool, lang string) {
	v := newTestCfg()
	v.Set("multilingual", multilingual)
	v.Set("defaultContentLanguage", "en")
	v.Set("defaultContentLanguageInSubdir", defaultInSubDir)

	tests := []struct {
		input    string
		baseURL  string
		canonify bool
		expected string
	}{
		{"/test/foo", "http://base/", false, "MULTI/test/foo"},
		{"/" + lang + "/test/foo", "http://base/", false, "/" + lang + "/test/foo"},
		{lang + "/test/foo", "http://base/", false, "/" + lang + "/test/foo"},
		{"test.css", "http://base/sub", false, "/subMULTI/test.css"},
		{"test.css", "http://base/sub", true, "MULTI/test.css"},
		{"/test/", "http://base/", false, "MULTI/test/"},
		{"/test/", "http://base/sub/", false, "/subMULTI/test/"},
		{"/test/", "http://base/sub/", true, "MULTI/test/"},
		{"", "http://base/ace/", false, "/aceMULTI/"},
		{"", "http://base/ace", false, "/aceMULTI"},
		{"http://abs", "http://base/", false, "http://abs"},
		{"//schemaless", "http://base/", false, "//schemaless"},
	}

	for i, test := range tests {
		v.Set("baseURL", test.baseURL)
		v.Set("canonifyURLs", test.canonify)
		l := langs.NewLanguage(lang, v)
		p, _ := NewPathSpec(hugofs.NewMem(v), l, nil)

		output := p.RelURL(test.input, addLanguage)

		expected := test.expected
		if multilingual && addLanguage {
			if !defaultInSubDir && lang == "en" {
				expected = strings.Replace(expected, "MULTI", "", 1)
			} else {
				expected = strings.Replace(expected, "MULTI", "/"+lang, 1)
			}
		} else {
			expected = strings.Replace(expected, "MULTI", "", 1)
		}

		if output != expected {
			t.Errorf("[%d][%t] Expected %#v, got %#v\n", i, test.canonify, expected, output)
		}
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://foo.bar/", "http://foo.bar"},
		{"http://foo.bar", "http://foo.bar"},          // issue #1105
		{"http://foo.bar/zoo/", "http://foo.bar/zoo"}, // issue #931
	}

	for i, test := range tests {
		o1 := SanitizeURL(test.input)
		o2 := SanitizeURLKeepTrailingSlash(test.input)

		expected2 := test.expected

		if strings.HasSuffix(test.input, "/") && !strings.HasSuffix(expected2, "/") {
			expected2 += "/"
		}

		if o1 != test.expected {
			t.Errorf("[%d] 1: Expected %#v, got %#v\n", i, test.expected, o1)
		}
		if o2 != expected2 {
			t.Errorf("[%d] 2: Expected %#v, got %#v\n", i, expected2, o2)
		}
	}
}

func TestMakePermalink(t *testing.T) {
	type test struct {
		host, link, output string
	}

	data := []test{
		{"http://abc.com/foo", "post/bar", "http://abc.com/foo/post/bar"},
		{"http://abc.com/foo/", "post/bar", "http://abc.com/foo/post/bar"},
		{"http://abc.com", "post/bar", "http://abc.com/post/bar"},
		{"http://abc.com", "bar", "http://abc.com/bar"},
		{"http://abc.com/foo/bar", "post/bar", "http://abc.com/foo/bar/post/bar"},
		{"http://abc.com/foo/bar", "post/bar/", "http://abc.com/foo/bar/post/bar/"},
	}

	for i, d := range data {
		output := MakePermalink(d.host, d.link).String()
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}
}

func TestURLPrep(t *testing.T) {
	type test struct {
		ugly   bool
		input  string
		output string
	}

	data := []test{
		{false, "/section/name.html", "/section/name/"},
		{true, "/section/name/index.html", "/section/name.html"},
	}

	for i, d := range data {
		v := newTestCfg()
		v.Set("uglyURLs", d.ugly)
		l := langs.NewDefaultLanguage(v)
		p, _ := NewPathSpec(hugofs.NewMem(v), l, nil)

		output := p.URLPrep(d.input)
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}

}

func TestAddContextRoot(t *testing.T) {
	tests := []struct {
		baseURL  string
		url      string
		expected string
	}{
		{"http://example.com/sub/", "/foo", "/sub/foo"},
		{"http://example.com/sub/", "/foo/index.html", "/sub/foo/index.html"},
		{"http://example.com/sub1/sub2", "/foo", "/sub1/sub2/foo"},
		{"http://example.com", "/foo", "/foo"},
		// cannot guess that the context root is already added int the example below
		{"http://example.com/sub/", "/sub/foo", "/sub/sub/foo"},
		{"http://example.com/тря", "/трям/", "/тря/трям/"},
		{"http://example.com", "/", "/"},
		{"http://example.com/bar", "//", "/bar/"},
	}

	for _, test := range tests {
		output := AddContextRoot(test.baseURL, test.url)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestPretty(t *testing.T) {
	c := qt.New(t)
	c.Assert("/section/name/index.html", qt.Equals, PrettifyURLPath("/section/name.html"))
	c.Assert("/section/sub/name/index.html", qt.Equals, PrettifyURLPath("/section/sub/name.html"))
	c.Assert("/section/name/index.html", qt.Equals, PrettifyURLPath("/section/name/"))
	c.Assert("/section/name/index.html", qt.Equals, PrettifyURLPath("/section/name/index.html"))
	c.Assert("/index.html", qt.Equals, PrettifyURLPath("/index.html"))
	c.Assert("/name/index.xml", qt.Equals, PrettifyURLPath("/name.xml"))
	c.Assert("/", qt.Equals, PrettifyURLPath("/"))
	c.Assert("/", qt.Equals, PrettifyURLPath(""))
	c.Assert("/section/name", qt.Equals, PrettifyURL("/section/name.html"))
	c.Assert("/section/sub/name", qt.Equals, PrettifyURL("/section/sub/name.html"))
	c.Assert("/section/name", qt.Equals, PrettifyURL("/section/name/"))
	c.Assert("/section/name", qt.Equals, PrettifyURL("/section/name/index.html"))
	c.Assert("/", qt.Equals, PrettifyURL("/index.html"))
	c.Assert("/name/index.xml", qt.Equals, PrettifyURL("/name.xml"))
	c.Assert("/", qt.Equals, PrettifyURL("/"))
	c.Assert("/", qt.Equals, PrettifyURL(""))
}

func TestUgly(t *testing.T) {
	c := qt.New(t)
	c.Assert("/section/name.html", qt.Equals, Uglify("/section/name.html"))
	c.Assert("/section/sub/name.html", qt.Equals, Uglify("/section/sub/name.html"))
	c.Assert("/section/name.html", qt.Equals, Uglify("/section/name/"))
	c.Assert("/section/name.html", qt.Equals, Uglify("/section/name/index.html"))
	c.Assert("/index.html", qt.Equals, Uglify("/index.html"))
	c.Assert("/name.xml", qt.Equals, Uglify("/name.xml"))
	c.Assert("/", qt.Equals, Uglify("/"))
	c.Assert("/", qt.Equals, Uglify(""))
}
