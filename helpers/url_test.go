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
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
)

func TestURLize(t *testing.T) {
	p := newTestPathSpec()

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
	c := qt.New(t)

	tests := []struct {
		input    string
		baseURL  string
		expected string
	}{
		// Issue 9994
		{"foo/bar", "https://example.org/foo/", "https://example.org/foo/MULTIfoo/bar"},
		{"/foo/bar", "https://example.org/foo/", "https://example.org/MULTIfoo/bar"},

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
		c.Run(fmt.Sprintf("%v/%t-%t-%t/%s", test, defaultInSubDir, addLanguage, multilingual, lang), func(c *qt.C) {
			v := config.New()
			if multilingual {
				v.Set("languages", map[string]any{
					"fr": map[string]interface{}{
						"weight": 20,
					},
					"en": map[string]interface{}{
						"weight": 10,
					},
				})
				v.Set("defaultContentLanguage", "en")
			} else {
				v.Set("defaultContentLanguage", lang)
				v.Set("languages", map[string]any{
					lang: map[string]interface{}{
						"weight": 10,
					},
				})
			}

			v.Set("defaultContentLanguageInSubdir", defaultInSubDir)
			v.Set("baseURL", test.baseURL)

			var configLang string
			if multilingual {
				configLang = lang
			}
			defaultContentLanguage := lang
			if multilingual {
				defaultContentLanguage = "en"
			}

			p := newTestPathSpecFromCfgAndLang(v, configLang)

			output := p.AbsURL(test.input, addLanguage)
			expected := test.expected
			if addLanguage {
				addLanguage = defaultInSubDir && lang == defaultContentLanguage
				addLanguage = addLanguage || (lang != defaultContentLanguage && multilingual)
			}
			if addLanguage {
				expected = strings.Replace(expected, "MULTI", lang+"/", 1)
			} else {
				expected = strings.Replace(expected, "MULTI", "", 1)
			}

			c.Assert(output, qt.Equals, expected)
		})
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

func doTestRelURL(t testing.TB, defaultInSubDir, addLanguage, multilingual bool, lang string) {
	t.Helper()
	c := qt.New(t)
	v := config.New()
	if multilingual {
		v.Set("languages", map[string]any{
			"fr": map[string]interface{}{
				"weight": 20,
			},
			"en": map[string]interface{}{
				"weight": 10,
			},
		})
		v.Set("defaultContentLanguage", "en")
	} else {
		v.Set("defaultContentLanguage", lang)
		v.Set("languages", map[string]any{
			lang: map[string]interface{}{
				"weight": 10,
			},
		})
	}

	v.Set("defaultContentLanguageInSubdir", defaultInSubDir)

	tests := []struct {
		input    string
		baseURL  string
		canonify bool
		expected string
	}{
		// Issue 9994
		{"/foo/bar", "https://example.org/foo/", false, "MULTI/foo/bar"},
		{"foo/bar", "https://example.org/foo/", false, "/fooMULTI/foo/bar"},

		// Issue 11080
		{"mailto:a@b.com", "http://base/", false, "mailto:a@b.com"},
		{"ftp://b.com/a.txt", "http://base/", false, "ftp://b.com/a.txt"},

		{"/test/foo", "http://base/", false, "MULTI/test/foo"},
		{"/" + lang + "/test/foo", "http://base/", false, "/" + lang + "/test/foo"},
		{lang + "/test/foo", "http://base/", false, "/" + lang + "/test/foo"},
		{"test.css", "http://base/sub", false, "/subMULTI/test.css"},
		{"test.css", "http://base/sub", true, "MULTI/test.css"},
		{"/test/", "http://base/", false, "MULTI/test/"},
		{"test/", "http://base/sub/", false, "/subMULTI/test/"},
		{"/test/", "http://base/sub/", true, "MULTI/test/"},
		{"", "http://base/ace/", false, "/aceMULTI/"},
		{"", "http://base/ace", false, "/aceMULTI/"},
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
		c.Run(fmt.Sprintf("%v/defaultInSubDir=%t;addLanguage=%t;multilingual=%t/%s", test, defaultInSubDir, addLanguage, multilingual, lang), func(c *qt.C) {
			v.Set("baseURL", test.baseURL)
			v.Set("canonifyURLs", test.canonify)
			defaultContentLanguage := lang
			if multilingual {
				defaultContentLanguage = "en"
			}
			p := newTestPathSpecFromCfgAndLang(v, lang)

			output := p.RelURL(test.input, addLanguage)

			expected := test.expected
			if addLanguage {
				addLanguage = defaultInSubDir && lang == defaultContentLanguage
				addLanguage = addLanguage || (lang != defaultContentLanguage && multilingual)
			}
			if addLanguage {
				expected = strings.Replace(expected, "MULTI", "/"+lang, 1)
			} else {
				expected = strings.Replace(expected, "MULTI", "", 1)
			}

			c.Assert(output, qt.Equals, expected, qt.Commentf("[%d] %s", i, test.input))
		})
	}
}

func BenchmarkRelURL(b *testing.B) {
	v := config.New()
	v.Set("baseURL", "https://base/")
	p := newTestPathSpecFromCfgAndLang(v, "")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.RelURL("https://base/foo/bar", false)
	}
}

func BenchmarkAbsURL(b *testing.B) {
	v := config.New()
	v.Set("baseURL", "https://base/")
	p := newTestPathSpecFromCfgAndLang(v, "")
	b.ResetTimer()
	b.Run("relurl", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = p.AbsURL("foo/bar", false)
		}
	})
	b.Run("absurl", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = p.AbsURL("https://base/foo/bar", false)
		}
	})
}
