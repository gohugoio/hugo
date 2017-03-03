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

package hugolib

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestTargetPathHTMLRedirectAlias(t *testing.T) {
	w := siteWriter{log: newErrorLogger()}

	errIsNilForThisOS := runtime.GOOS != "windows"

	tests := []struct {
		value    string
		expected string
		errIsNil bool
	}{
		{"", "", false},
		{"s", filepath.FromSlash("s/index.html"), true},
		{"/", "", false},
		{"alias 1", filepath.FromSlash("alias 1/index.html"), true},
		{"alias 2/", filepath.FromSlash("alias 2/index.html"), true},
		{"alias 3.html", "alias 3.html", true},
		{"alias4.html", "alias4.html", true},
		{"/alias 5.html", "alias 5.html", true},
		{"/трям.html", "трям.html", true},
		{"../../../../tmp/passwd", "", false},
		{"/foo/../../../../tmp/passwd", filepath.FromSlash("tmp/passwd/index.html"), true},
		{"foo/../../../../tmp/passwd", "", false},
		{"C:\\Windows", filepath.FromSlash("C:\\Windows/index.html"), errIsNilForThisOS},
		{"/trailing-space /", filepath.FromSlash("trailing-space /index.html"), errIsNilForThisOS},
		{"/trailing-period./", filepath.FromSlash("trailing-period./index.html"), errIsNilForThisOS},
		{"/tab\tseparated/", filepath.FromSlash("tab\tseparated/index.html"), errIsNilForThisOS},
		{"/chrome/?p=help&ctx=keyboard#topic=3227046", filepath.FromSlash("chrome/?p=help&ctx=keyboard#topic=3227046/index.html"), errIsNilForThisOS},
		{"/LPT1/Printer/", filepath.FromSlash("LPT1/Printer/index.html"), errIsNilForThisOS},
	}

	for _, test := range tests {
		path, err := w.targetPathAlias(test.value)
		if (err == nil) != test.errIsNil {
			t.Errorf("Expected err == nil => %t, got: %t. err: %s", test.errIsNil, err == nil, err)
			continue
		}
		if err == nil && path != test.expected {
			t.Errorf("Expected: \"%s\", got: \"%s\"", test.expected, path)
		}
	}
}

func TestTargetPathPage(t *testing.T) {
	w := siteWriter{log: newErrorLogger()}

	tests := []struct {
		content  string
		expected string
	}{
		{"/", "index.html"},
		{"index.html", "index.html"},
		{"bar/index.html", "bar/index.html"},
		{"foo", "foo/index.html"},
		{"foo.html", "foo/index.html"},
		{"foo.xhtml", "foo/index.xhtml"},
		{"section", "section/index.html"},
		{"section/", "section/index.html"},
		{"section/foo", "section/foo/index.html"},
		{"section/foo.html", "section/foo/index.html"},
		{"section/foo.rss", "section/foo/index.rss"},
	}

	for _, test := range tests {
		dest, err := w.targetPathPage(filepath.FromSlash(test.content))
		expected := filepath.FromSlash(test.expected)
		if err != nil {
			t.Fatalf("Translate returned and unexpected err: %s", err)
		}

		if dest != expected {
			t.Errorf("Translate expected return: %s, got: %s", expected, dest)
		}
	}
}

func TestTargetPathPageBase(t *testing.T) {
	w := siteWriter{log: newErrorLogger()}

	tests := []struct {
		content  string
		expected string
	}{
		{"/", "a/base/index.html"},
	}

	for _, test := range tests {

		for _, pd := range []string{"a/base", "a/base/"} {
			w.publishDir = pd
			dest, err := w.targetPathPage(test.content)
			if err != nil {
				t.Fatalf("Translated returned and err: %s", err)
			}

			if dest != filepath.FromSlash(test.expected) {
				t.Errorf("Translate expected: %s, got: %s", test.expected, dest)
			}
		}
	}
}

func TestTargetPathUglyURLs(t *testing.T) {
	w := siteWriter{log: newErrorLogger(), uglyURLs: true}

	tests := []struct {
		content  string
		expected string
	}{
		{"foo.html", "foo.html"},
		{"/", "index.html"},
		{"section", "section.html"},
		{"index.html", "index.html"},
	}

	for _, test := range tests {
		dest, err := w.targetPathPage(filepath.FromSlash(test.content))
		if err != nil {
			t.Fatalf("Translate returned an unexpected err: %s", err)
		}

		if dest != test.expected {
			t.Errorf("Translate expected return: %s, got: %s", test.expected, dest)
		}
	}
}
