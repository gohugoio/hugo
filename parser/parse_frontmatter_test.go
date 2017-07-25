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

package parser

// TODO Support Mac Encoding (\r)

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	contentNoFrontmatter                        = "a page with no front matter"
	contentWithFrontmatter                      = "---\ntitle: front matter\n---\nContent with front matter"
	contentHTMLNoDoctype                        = "<html>\n\t<body>\n\t</body>\n</html>"
	contentHTMLWithDoctype                      = "<!doctype html><html><body></body></html>"
	contentHTMLWithFrontmatter                  = "---\ntitle: front matter\n---\n<!doctype><html><body></body></html>"
	contentHTML                                 = "    <html><body></body></html>"
	contentLinefeedAndHTML                      = "\n<html><body></body></html>"
	contentIncompleteEndFrontmatterDelim        = "---\ntitle: incomplete end fm delim\n--\nincomplete frontmatter delim"
	contentMissingEndFrontmatterDelim           = "---\ntitle: incomplete end fm delim\nincomplete frontmatter delim"
	contentSlugWorking                          = "---\ntitle: slug doc 2\nslug: slug-doc-2\n\n---\nslug doc 2 content"
	contentSlugWorkingVariation                 = "---\ntitle: slug doc 3\nslug: slug-doc 3\n---\nslug doc 3 content"
	contentSlugBug                              = "---\ntitle: slug doc 2\nslug: slug-doc-2\n---\nslug doc 2 content"
	contentSlugWithJSONFrontMatter              = "{\n  \"categories\": \"d\",\n  \"tags\": [\n    \"a\", \n    \"b\", \n    \"c\"\n  ]\n}\nJSON Front Matter with tags and categories"
	contentWithJSONLooseFrontmatter             = "{\n  \"categories\": \"d\"\n  \"tags\": [\n    \"a\" \n    \"b\" \n    \"c\"\n  ]\n}\nJSON Front Matter with tags and categories"
	contentSlugWithJSONFrontMatterFirstLineOnly = "{\"categories\":\"d\",\"tags\":[\"a\",\"b\",\"c\"]}\nJSON Front Matter with tags and categories"
	contentSlugWithJSONFrontMatterFirstLine     = "{\"categories\":\"d\",\n  \"tags\":[\"a\",\"b\",\"c\"]}\nJSON Front Matter with tags and categories"
)

var lineEndings = []string{"\n", "\r\n"}
var delimiters = []string{"---", "+++"}

func pageMust(p Page, err error) *page {
	if err != nil {
		panic(err)
	}
	return p.(*page)
}

func TestDegenerateCreatePageFrom(t *testing.T) {
	tests := []struct {
		content string
	}{
		{contentMissingEndFrontmatterDelim},
		{contentIncompleteEndFrontmatterDelim},
	}

	for _, test := range tests {
		for _, ending := range lineEndings {
			test.content = strings.Replace(test.content, "\n", ending, -1)
			_, err := ReadFrom(strings.NewReader(test.content))
			if err == nil {
				t.Errorf("Content should return an err:\n%q\n", test.content)
			}
		}
	}
}

func checkPageRender(t *testing.T, p *page, expected bool) {
	if p.render != expected {
		t.Errorf("page.render should be %t, got: %t", expected, p.render)
	}
}

func checkPageFrontMatterIsNil(t *testing.T, p *page, content string, expected bool) {
	if bool(p.frontmatter == nil) != expected {
		t.Logf("\n%q\n", content)
		t.Errorf("page.frontmatter == nil? %t, got %t", expected, p.frontmatter == nil)
	}
}

func checkPageFrontMatterContent(t *testing.T, p *page, frontMatter string) {
	if p.frontmatter == nil {
		return
	}
	if !bytes.Equal(p.frontmatter, []byte(frontMatter)) {
		t.Errorf("frontmatter mismatch\nexp: %q\ngot: %q", frontMatter, p.frontmatter)
	}
}

func checkPageContent(t *testing.T, p *page, expected string) {
	if !bytes.Equal(p.content, []byte(expected)) {
		t.Errorf("content mismatch\nexp: %q\ngot: %q", expected, p.content)
	}
}

func TestStandaloneCreatePageFrom(t *testing.T) {
	tests := []struct {
		content            string
		expectedMustRender bool
		frontMatterIsNil   bool
		frontMatter        string
		bodycontent        string
	}{

		{contentNoFrontmatter, true, true, "", "a page with no front matter"},
		{contentWithFrontmatter, true, false, "---\ntitle: front matter\n---\n", "Content with front matter"},
		{contentHTMLNoDoctype, false, true, "", "<html>\n\t<body>\n\t</body>\n</html>"},
		{contentHTMLWithDoctype, false, true, "", "<!doctype html><html><body></body></html>"},
		{contentHTMLWithFrontmatter, true, false, "---\ntitle: front matter\n---\n", "<!doctype><html><body></body></html>"},
		{contentHTML, false, true, "", "<html><body></body></html>"},
		{contentLinefeedAndHTML, false, true, "", "<html><body></body></html>"},
		{contentSlugWithJSONFrontMatter, true, false, "{\n  \"categories\": \"d\",\n  \"tags\": [\n    \"a\", \n    \"b\", \n    \"c\"\n  ]\n}", "JSON Front Matter with tags and categories"},
		{contentWithJSONLooseFrontmatter, true, false, "{\n  \"categories\": \"d\"\n  \"tags\": [\n    \"a\" \n    \"b\" \n    \"c\"\n  ]\n}", "JSON Front Matter with tags and categories"},
		{contentSlugWithJSONFrontMatterFirstLineOnly, true, false, "{\"categories\":\"d\",\"tags\":[\"a\",\"b\",\"c\"]}", "JSON Front Matter with tags and categories"},
		{contentSlugWithJSONFrontMatterFirstLine, true, false, "{\"categories\":\"d\",\n  \"tags\":[\"a\",\"b\",\"c\"]}", "JSON Front Matter with tags and categories"},
		{contentSlugWorking, true, false, "---\ntitle: slug doc 2\nslug: slug-doc-2\n\n---\n", "slug doc 2 content"},
		{contentSlugWorkingVariation, true, false, "---\ntitle: slug doc 3\nslug: slug-doc 3\n---\n", "slug doc 3 content"},
		{contentSlugBug, true, false, "---\ntitle: slug doc 2\nslug: slug-doc-2\n---\n", "slug doc 2 content"},
	}

	for _, test := range tests {
		for _, ending := range lineEndings {
			test.content = strings.Replace(test.content, "\n", ending, -1)
			test.frontMatter = strings.Replace(test.frontMatter, "\n", ending, -1)
			test.bodycontent = strings.Replace(test.bodycontent, "\n", ending, -1)

			p := pageMust(ReadFrom(strings.NewReader(test.content)))

			checkPageRender(t, p, test.expectedMustRender)
			checkPageFrontMatterIsNil(t, p, test.content, test.frontMatterIsNil)
			checkPageFrontMatterContent(t, p, test.frontMatter)
			checkPageContent(t, p, test.bodycontent)
		}
	}
}

func BenchmarkLongFormRender(b *testing.B) {

	tests := []struct {
		filename string
		buf      []byte
	}{
		{filename: "long_text_test.md"},
	}
	for i, test := range tests {
		path := filepath.FromSlash(test.filename)
		f, err := os.Open(path)
		if err != nil {
			b.Fatalf("Unable to open %s: %s", path, err)
		}
		defer f.Close()
		membuf := new(bytes.Buffer)
		if _, err := io.Copy(membuf, f); err != nil {
			b.Fatalf("Unable to read %s: %s", path, err)
		}
		tests[i].buf = membuf.Bytes()
	}

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		for _, test := range tests {
			ReadFrom(bytes.NewReader(test.buf))
		}
	}
}

func TestPageShouldRender(t *testing.T) {
	tests := []struct {
		content  []byte
		expected bool
	}{
		{[]byte{}, false},
		{[]byte{'<'}, false},
		{[]byte{'-'}, true},
		{[]byte("--"), true},
		{[]byte("---"), true},
		{[]byte("---\n"), true},
		{[]byte{'a'}, true},
	}

	for _, test := range tests {
		for _, ending := range lineEndings {
			test.content = bytes.Replace(test.content, []byte("\n"), []byte(ending), -1)
			if render := shouldRender(test.content); render != test.expected {

				t.Errorf("Expected %s to shouldRender = %t, got: %t", test.content, test.expected, render)
			}
		}
	}
}

func TestPageHasFrontMatter(t *testing.T) {
	tests := []struct {
		content  []byte
		expected bool
	}{
		{[]byte{'-'}, false},
		{[]byte("--"), false},
		{[]byte("---"), false},
		{[]byte("---\n"), true},
		{[]byte("---\n"), true},
		{[]byte("--- \n"), true},
		{[]byte("---  \n"), true},
		{[]byte{'a'}, false},
		{[]byte{'{'}, true},
		{[]byte("{\n  "), true},
		{[]byte{'}'}, false},
	}
	for _, test := range tests {
		for _, ending := range lineEndings {
			test.content = bytes.Replace(test.content, []byte("\n"), []byte(ending), -1)
			if isFrontMatterDelim := isFrontMatterDelim(test.content); isFrontMatterDelim != test.expected {
				t.Errorf("Expected %q isFrontMatterDelim = %t,  got: %t", test.content, test.expected, isFrontMatterDelim)
			}
		}
	}
}

func TestExtractFrontMatter(t *testing.T) {

	tests := []struct {
		frontmatter string
		extracted   []byte
		errIsNil    bool
	}{
		{"", nil, false},
		{"-", nil, false},
		{"---\n", nil, false},
		{"---\nfoobar", nil, false},
		{"---\nfoobar\nbarfoo\nfizbaz\n", nil, false},
		{"---\nblar\n-\n", nil, false},
		{"---\nralb\n---\n", []byte("---\nralb\n---\n"), true},
		{"---\neof\n---", []byte("---\neof\n---"), true},
		{"--- \neof\n---", []byte("---\neof\n---"), true},
		{"---\nminc\n---\ncontent", []byte("---\nminc\n---\n"), true},
		{"---\nminc\n---    \ncontent", []byte("---\nminc\n---\n"), true},
		{"---  \nminc\n--- \ncontent", []byte("---\nminc\n---\n"), true},
		{"---\ncnim\n---\ncontent\n", []byte("---\ncnim\n---\n"), true},
		{"---\ntitle: slug doc 2\nslug: slug-doc-2\n---\ncontent\n", []byte("---\ntitle: slug doc 2\nslug: slug-doc-2\n---\n"), true},
		{"---\npermalink: '/blog/title---subtitle.html'\n---\ncontent\n", []byte("---\npermalink: '/blog/title---subtitle.html'\n---\n"), true},
	}

	for _, test := range tests {
		for _, ending := range lineEndings {
			test.frontmatter = strings.Replace(test.frontmatter, "\n", ending, -1)
			test.extracted = bytes.Replace(test.extracted, []byte("\n"), []byte(ending), -1)
			for _, delim := range delimiters {
				test.frontmatter = strings.Replace(test.frontmatter, "---", delim, -1)
				test.extracted = bytes.Replace(test.extracted, []byte("---"), []byte(delim), -1)
				line, err := peekLine(bufio.NewReader(strings.NewReader(test.frontmatter)))
				if err != nil {
					continue
				}
				l, r := determineDelims(line)
				fm, err := extractFrontMatterDelims(bufio.NewReader(strings.NewReader(test.frontmatter)), l, r)
				if (err == nil) != test.errIsNil {
					t.Logf("\n%q\n", string(test.frontmatter))
					t.Errorf("Expected err == nil => %t, got: %t. err: %s", test.errIsNil, err == nil, err)
					continue
				}
				if !bytes.Equal(fm, test.extracted) {
					t.Errorf("Frontmatter did not match:\nexp: %q\ngot: %q", string(test.extracted), fm)
				}
			}
		}
	}
}

func TestExtractFrontMatterDelim(t *testing.T) {
	var (
		noErrExpected = true
		errExpected   = false
	)
	tests := []struct {
		frontmatter string
		extracted   string
		errIsNil    bool
	}{
		{"", "", errExpected},
		{"{", "", errExpected},
		{"{}", "{}", noErrExpected},
		{"{} ", "{}", noErrExpected},
		{"{ } ", "{ }", noErrExpected},
		{"{ { }", "", errExpected},
		{"{ { } }", "{ { } }", noErrExpected},
		{"{ { } { } }", "{ { } { } }", noErrExpected},
		{"{\n{\n}\n}\n", "{\n{\n}\n}", noErrExpected},
		{"{\n  \"categories\": \"d\",\n  \"tags\": [\n    \"a\", \n    \"b\", \n    \"c\"\n  ]\n}\nJSON Front Matter with tags and categories", "{\n  \"categories\": \"d\",\n  \"tags\": [\n    \"a\", \n    \"b\", \n    \"c\"\n  ]\n}", noErrExpected},
		{"{\n  \"categories\": \"d\"\n  \"tags\": [\n    \"a\" \n    \"b\" \n    \"c\"\n  ]\n}\nJSON Front Matter with tags and categories", "{\n  \"categories\": \"d\"\n  \"tags\": [\n    \"a\" \n    \"b\" \n    \"c\"\n  ]\n}", noErrExpected},
		// Issue #3511
		{`{ "title": "{" }`, `{ "title": "{" }`, noErrExpected},
		{`{ "title": "{}" }`, `{ "title": "{}" }`, noErrExpected},
		// Issue #3661
		{`{ "title": "\"" }`, `{ "title": "\"" }`, noErrExpected},
		{`{ "title": "\"{", "other": "\"{}" }`, `{ "title": "\"{", "other": "\"{}" }`, noErrExpected},
		{`{ "title": "\"Foo\"" }`, `{ "title": "\"Foo\"" }`, noErrExpected},
		{`{ "title": "\"Foo\"\"" }`, `{ "title": "\"Foo\"\"" }`, noErrExpected},
		{`{ "url": "http:\/\/example.com\/play\/url?id=1" }`, `{ "url": "http:\/\/example.com\/play\/url?id=1" }`, noErrExpected},
		{`{ "test": "\"New\r\nString\"" }`, `{ "test": "\"New\r\nString\"" }`, noErrExpected},
		{`{ "test": "RTS\/RPG" }`, `{ "test": "RTS\/RPG" }`, noErrExpected},
	}

	for i, test := range tests {
		fm, err := extractFrontMatterDelims(bufio.NewReader(strings.NewReader(test.frontmatter)), []byte("{"), []byte("}"))
		if (err == nil) != test.errIsNil {
			t.Logf("\n%q\n", string(test.frontmatter))
			t.Errorf("[%d] Expected err == nil => %t, got: %t. err: %s", i, test.errIsNil, err == nil, err)
			continue
		}
		if !bytes.Equal(fm, []byte(test.extracted)) {
			t.Logf("\n%q\n", string(test.frontmatter))
			t.Errorf("[%d] Frontmatter did not match:\nexp: %q\ngot:  %q", i, string(test.extracted), fm)
		}
	}
}
