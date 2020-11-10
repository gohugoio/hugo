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
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/paths"

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

// TODO1 remove this.
func BenchmarkURLEscape(b *testing.B) {
	const (
		input                   = "трям/трям"
		expect                  = "%D1%82%D1%80%D1%8F%D0%BC/%D1%82%D1%80%D1%8F%D0%BC"
		forwardSlashReplacement = "ABC"
	)

	fn1 := func(s string) string {
		ss, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		return ss.EscapedPath()
	}

	fn2 := func(s string) string {
		s = strings.ReplaceAll(s, "/", forwardSlashReplacement)
		s = url.PathEscape(s)
		s = strings.ReplaceAll(s, forwardSlashReplacement, "/")

		return s
	}

	fn3 := func(s string) string {
		parts := paths.FieldsSlash(s)
		for i, part := range parts {
			parts[i] = url.PathEscape(part)
		}

		return path.Join(parts...)
	}

	benchFunc := func(b *testing.B, fn func(s string) string) {
		for i := 0; i < b.N; i++ {
			res := fn(input)
			if res != expect {
				b.Fatal(res)
			}
		}
	}

	b.Run("url.Parse", func(b *testing.B) {
		benchFunc(b, fn1)
	})

	b.Run("url.PathEscape_replace", func(b *testing.B) {
		benchFunc(b, fn2)
	})

	b.Run("url.PathEscape_fields", func(b *testing.B) {
		benchFunc(b, fn3)
	})

	b.Run("url.PathEscape", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := url.PathEscape(input)
			// url.PathEscape also escapes forward slash.
			if res != "%D1%82%D1%80%D1%8F%D0%BC%2F%D1%82%D1%80%D1%8F%D0%BC" {
				panic(res)
			}
		}
	})

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

	if multilingual && addLanguage && defaultInSubDir {
		newTests := []struct {
			input    string
			baseURL  string
			expected string
		}{
			{lang + "test", "http://base/", "http://base/" + lang + "/" + lang + "test"},
			{"/" + lang + "test", "http://base/", "http://base/" + lang + "/" + lang + "test"},
		}

		tests = append(tests, newTests...)

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

	if multilingual && addLanguage && defaultInSubDir {
		newTests := []struct {
			input    string
			baseURL  string
			canonify bool
			expected string
		}{
			{lang + "test", "http://base/", false, "/" + lang + "/" + lang + "test"},
			{"/" + lang + "test", "http://base/", false, "/" + lang + "/" + lang + "test"},
		}
		tests = append(tests, newTests...)
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
