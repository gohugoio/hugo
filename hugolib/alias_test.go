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

package hugolib

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
)

func TestAlias(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fileSuffix string
		urlPrefix  string
		urlSuffix  string
		settings   map[string]any
	}{
		{"/index.html", "http://example.com", "/", map[string]any{"baseURL": "http://example.com"}},
		{"/index.html", "http://example.com/some/path", "/", map[string]any{"baseURL": "http://example.com/some/path"}},
		{"/index.html", "http://example.com", "/", map[string]any{"baseURL": "http://example.com", "canonifyURLs": true}},
		{"/index.html", "../..", "/", map[string]any{"relativeURLs": true}},
		{".html", "", ".html", map[string]any{"uglyURLs": true}},
	}

	for _, test := range tests {
		files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "taxonomy", "term"]
CONFIG
-- content/blog/page.md --
---
title: Has Alias
aliases: ["/foo/bar/", "rel"]
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.
-- layouts/all.html --
Title: {{ .Title }}|Content: {{ .Content }}|
`
		files = strings.Replace(files, "CONFIG", config.FromMapToTOMLString(test.settings), 1)

		b := Test(t, files)

		// the real page
		b.AssertFileContent("public/blog/page"+test.fileSuffix, "For some moments the old man")

		// the alias redirectors
		b.AssertFileContent("public/foo/bar"+test.fileSuffix, "<meta http-equiv=\"refresh\" content=\"0; url="+test.urlPrefix+"/blog/page"+test.urlSuffix+"\">")
		b.AssertFileContent("public/blog/rel"+test.fileSuffix, "<meta http-equiv=\"refresh\" content=\"0; url="+test.urlPrefix+"/blog/page"+test.urlSuffix+"\">")
	}
}

func TestAliasMultipleOutputFormats(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com"
-- layouts/single.html --
{{ .Content }}
-- layouts/single.amp.html --
{{ .Content }}
-- layouts/single.json --
{{ .Content }}
-- content/blog/page.md --
---
title: Has Alias for HTML and AMP
aliases: ["/foo/bar/"]
outputs: ["html", "amp", "json"]
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.
`

	b := Test(t, files)

	// the real pages
	b.AssertFileContent("public/blog/page/index.html", "For some moments the old man")
	b.AssertFileContent("public/amp/blog/page/index.html", "For some moments the old man")
	b.AssertFileContent("public/blog/page/index.json", "For some moments the old man")

	// the alias redirectors
	b.AssertFileContent("public/foo/bar/index.html", "<meta http-equiv=\"refresh\" content=\"0; ")
	b.AssertFileContent("public/amp/foo/bar/index.html", "<meta http-equiv=\"refresh\" content=\"0; ")
	b.AssertFileExists("public/foo/bar/index.json", false)
}

func TestAliasTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com"
-- layouts/single.html --
Single.
-- layouts/home.html --
Home.
-- layouts/alias.html --
ALIASTEMPLATE
-- content/page.md --
---
title: "Page"
aliases: ["/foo/bar/"]
---
`

	b := Test(t, files)

	// the real page
	b.AssertFileContent("public/page/index.html", "Single.")
	// the alias redirector
	b.AssertFileContent("public/foo/bar/index.html", "ALIASTEMPLATE")
}

func TestTargetPathHTMLRedirectAlias(t *testing.T) {
	h := newAliasHandler(nil, loggers.NewDefault(), false)

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
		path, err := h.targetPathAlias(test.value)
		if (err == nil) != test.errIsNil {
			t.Errorf("Expected err == nil => %t, got: %t. err: %s", test.errIsNil, err == nil, err)
			continue
		}
		if err == nil && path != test.expected {
			t.Errorf("Expected: %q, got: %q", test.expected, path)
		}
	}
}

func TestAliasNIssue14053(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com"
-- layouts/all.html --
All.
-- content/page.md --
---
title: "Page"
aliases:
- n
- y
- no
- yes
---
`
	b := Test(t, files)

	b.AssertPublishDir("n/index.html", "yes/index.html", "no/index.html", "yes/index.html")
}

// Issue 14381
func TestIssue14381(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL      = 'https://example.org'
disableKinds = ['home', 'rss', 'sitemap', 'taxonomy', 'term']

[outputFormats.print]
  isHTML        = IS_HTML
  permalinkable = true
  mediaType     = 'text/html'
  path          = 'print'

[outputs]
  page    = ['html', 'print']
  section = ['html', 'print']
-- content/foo/s1/_index.md --
---
title: s1
aliases: [s1-alias]
---
-- content/foo/s1/p1.md --
---
title: p1
aliases: [p1-alias]
---
-- content/foo/s2/_index.md --
---
title: s2
aliases: [/s2-alias]
---
-- content/foo/s2/p2.md --
---
title: p2
aliases: [/p2-alias]
---
-- layouts/all.html --
{{ .Title }}
`
	// ------------------------------------------------------------------------
	// Test 1: Create aliases for the html output format only
	// ------------------------------------------------------------------------

	// public/
	// ├── foo/
	// │  ├── s1/
	// │  │  ├── p1/
	// │  │  │  └── index.html
	// │  │  ├── p1-alias/
	// │  │  │  └── index.html
	// │  │  └── index.html
	// │  ├── s1-alias/
	// │  │  └── index.html
	// │  ├── s2/
	// │  │  ├── p2/
	// │  │  │  └── index.html
	// │  │  └── index.html
	// │  └── index.html
	// ├── p2-alias/
	// │  └── index.html
	// ├── print/
	// │  └── foo/
	// │      ├── s1/
	// │      │  ├── p1/
	// │      │  │  └── index.html
	// │      │  └── index.html
	// │      ├── s2/
	// │      │  ├── p2/
	// │      │  │  └── index.html
	// │      │  └── index.html
	// │      └── index.html
	// └── s2-alias/
	//     └── index.html

	f := strings.ReplaceAll(files, "IS_HTML", "false")
	b := Test(t, f)

	// output format: html
	b.AssertFileContent("public/foo/s1-alias/index.html",
		`<title>https://example.org/foo/s1/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s1/">`,
	)
	b.AssertFileContent("public/foo/s1/p1-alias/index.html",
		`<title>https://example.org/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/s2-alias/index.html",
		`<title>https://example.org/foo/s2/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s2/">`,
	)
	b.AssertFileContent("public/p2-alias/index.html",
		`<title>https://example.org/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s2/p2/">`,
	)

	// output format: print
	b.AssertFileExists("public/print/foo/s1-alias/index.html", false)
	b.AssertFileExists("public/print/foo/s1/p1-alias/index.html", false)
	b.AssertFileExists("public/print/s2-alias/index.html", false)
	b.AssertFileExists("public/print/p2-alias/index.html", false)

	// ------------------------------------------------------------------------
	// Test 2: Create aliases for the html and print output formats
	// ------------------------------------------------------------------------

	// public/
	// ├── foo/
	// │  ├── s1/
	// │  │  ├── p1/
	// │  │  │  └── index.html
	// │  │  ├── p1-alias/
	// │  │  │  └── index.html
	// │  │  └── index.html
	// │  ├── s1-alias/
	// │  │  └── index.html
	// │  ├── s2/
	// │  │  ├── p2/
	// │  │  │  └── index.html
	// │  │  └── index.html
	// │  └── index.html
	// ├── p2-alias/
	// │  └── index.html
	// ├── print/
	// │  ├── foo/
	// │  │  ├── s1/
	// │  │  │  ├── p1/
	// │  │  │  │  └── index.html
	// │  │  │  ├── p1-alias/
	// │  │  │  │  └── index.html
	// │  │  │  └── index.html
	// │  │  ├── s1-alias/
	// │  │  │  └── index.html
	// │  │  ├── s2/
	// │  │  │  ├── p2/
	// │  │  │  │  └── index.html
	// │  │  │  └── index.html
	// │  │  └── index.html
	// │  ├── p2-alias/
	// │  │  └── index.html
	// │  └── s2-alias/
	// │      └── index.html
	// └── s2-alias/
	//     └── index.html

	f = strings.ReplaceAll(files, "IS_HTML", "true")
	b = Test(t, f)

	// output format: html
	b.AssertFileContent("public/foo/s1-alias/index.html",
		`<title>https://example.org/foo/s1/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s1/">`,
	)
	b.AssertFileContent("public/foo/s1/p1-alias/index.html",
		`<title>https://example.org/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/s2-alias/index.html",
		`<title>https://example.org/foo/s2/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s2/">`,
	)
	b.AssertFileContent("public/p2-alias/index.html",
		`<title>https://example.org/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/foo/s2/p2/">`,
	)

	// output format: print
	b.AssertFileContent("public/print/foo/s1-alias/index.html",
		`<title>https://example.org/print/foo/s1/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/print/foo/s1/">`,
	)
	b.AssertFileContent("public/print/foo/s1/p1-alias/index.html",
		`<title>https://example.org/print/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/print/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/print/s2-alias/index.html",
		`<title>https://example.org/print/foo/s2/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/print/foo/s2/">`,
	)
	b.AssertFileContent("public/print/p2-alias/index.html",
		`<title>https://example.org/print/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://example.org/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/print/foo/s2/p2/">`,
	)
}

// Issue 14388
func TestIssue14388(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL      = 'https://example.org'
disableKinds = ['home', 'rss', 'sitemap', 'taxonomy', 'term']

defaultContentLanguage = 'en'

defaultContentLanguageInSubdir = true
defaultContentRoleInSubdir     = true
defaultContentVersionInSubdir  = true

[languages.en]
  weight = 1
  # baseURL      = 'https://en.example.org/'

[languages.de]
  weight = 2
  # baseURL      = 'https://de.example.org/'

[outputFormats.print]
  isHTML    = true
  mediaType = 'text/html'
  path      = 'print'
  permalinkable = true

[outputs]
  page    = ['html', 'print']
  section = ['html', 'print']
-- content/foo/s1/_index.de.md --
---
title: s1 de
aliases: [s1-alias]
---
-- content/foo/s1/_index.en.md --
---
title: s1 en
aliases: [s1-alias]
---
-- content/foo/s1/p1.de.md --
---
title: p1 de
aliases: [p1-alias]
---
-- content/foo/s1/p1.en.md --
---
title: p1 en
aliases: [p1-alias]
---
-- content/foo/s2/_index.de.md --
---
title: s2 de
aliases: [/s2-alias]
---
-- content/foo/s2/_index.en.md --
---
title: s2 en
aliases: [/s2-alias]
---
-- content/foo/s2/p2.de.md --
---
title: p2 de
aliases: [/p2-alias]
---
-- content/foo/s2/p2.en.md --
---
title: p2 en
aliases: [/p2-alias]
---
-- layouts/all.html --
{{ .Title }}
`
	// ------------------------------------------------------------------------
	// Test 1: Multilingual single-host
	// ------------------------------------------------------------------------

	b := Test(t, files)

	// language: de, output format: html
	b.AssertFileContent("public/guest/v1.0.0/de/foo/s1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/foo/s1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/foo/s1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/de/foo/s1/p1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/de/s2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/foo/s2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/foo/s2/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/de/p2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/foo/s2/p2/">`,
	)

	// language: de, output format: print
	b.AssertFileContent("public/guest/v1.0.0/de/print/foo/s1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/print/foo/s1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/print/foo/s1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/de/print/foo/s1/p1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/print/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/print/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/de/print/s2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/print/foo/s2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/print/foo/s2/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/de/print/p2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/de/print/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/de/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/de/print/foo/s2/p2/">`,
	)

	// language: en, output format: html
	b.AssertFileContent("public/guest/v1.0.0/en/foo/s1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/foo/s1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/foo/s1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/en/foo/s1/p1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/en/s2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/foo/s2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/foo/s2/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/en/p2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/foo/s2/p2/">`,
	)

	// language: en, output format: print
	b.AssertFileContent("public/guest/v1.0.0/en/print/foo/s1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/print/foo/s1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/print/foo/s1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/en/print/foo/s1/p1-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/print/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/print/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/en/print/s2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/print/foo/s2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/print/foo/s2/">`,
	)
	b.AssertFileContent("public/guest/v1.0.0/en/print/p2-alias/index.html",
		`<title>https://example.org/guest/v1.0.0/en/print/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://example.org/guest/v1.0.0/en/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://example.org/guest/v1.0.0/en/print/foo/s2/p2/">`,
	)

	// ------------------------------------------------------------------------
	// Test 2: Multilingual multihost
	// ------------------------------------------------------------------------

	files = strings.ReplaceAll(files, "# baseURL", "baseURL")
	b = Test(t, files)

	// language: de, output format: html
	b.AssertFileContent("public/de/guest/v1.0.0/foo/s1-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/foo/s1/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/foo/s1/">`,
	)
	b.AssertFileContent("public/de/guest/v1.0.0/foo/s1/p1-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/de/guest/v1.0.0/s2-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/foo/s2/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/foo/s2/">`,
	)
	b.AssertFileContent("public/de/guest/v1.0.0/p2-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/foo/s2/p2/">`,
	)

	// language: de, output format: print
	b.AssertFileContent("public/de/guest/v1.0.0/print/foo/s1-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/print/foo/s1/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/print/foo/s1/">`,
	)
	b.AssertFileContent("public/de/guest/v1.0.0/print/foo/s1/p1-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/print/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/print/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/de/guest/v1.0.0/print/s2-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/print/foo/s2/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/print/foo/s2/">`,
	)
	b.AssertFileContent("public/de/guest/v1.0.0/print/p2-alias/index.html",
		`<title>https://de.example.org/guest/v1.0.0/print/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://de.example.org/guest/v1.0.0/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://de.example.org/guest/v1.0.0/print/foo/s2/p2/">`,
	)

	// language: en, output format: html
	b.AssertFileContent("public/en/guest/v1.0.0/foo/s1-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/foo/s1/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/foo/s1/">`,
	)
	b.AssertFileContent("public/en/guest/v1.0.0/foo/s1/p1-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/en/guest/v1.0.0/s2-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/foo/s2/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/foo/s2/">`,
	)
	b.AssertFileContent("public/en/guest/v1.0.0/p2-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/foo/s2/p2/">`,
	)

	// language: en, output format: print
	b.AssertFileContent("public/en/guest/v1.0.0/print/foo/s1-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/print/foo/s1/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s1/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/print/foo/s1/">`,
	)
	b.AssertFileContent("public/en/guest/v1.0.0/print/foo/s1/p1-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/print/foo/s1/p1/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s1/p1/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/print/foo/s1/p1/">`,
	)
	b.AssertFileContent("public/en/guest/v1.0.0/print/s2-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/print/foo/s2/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s2/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/print/foo/s2/">`,
	)
	b.AssertFileContent("public/en/guest/v1.0.0/print/p2-alias/index.html",
		`<title>https://en.example.org/guest/v1.0.0/print/foo/s2/p2/</title>`,
		`<link rel="canonical" href="https://en.example.org/guest/v1.0.0/foo/s2/p2/">`,
		`<meta http-equiv="refresh" content="0; url=https://en.example.org/guest/v1.0.0/print/foo/s2/p2/">`,
	)
}

// Issue #14402.
func TestComprehensiveAliasesRedirectsFile(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableAliases = true
defaultContentLanguageInSubdir = true
defaultContentRoleInSubdir     = true
defaultContentVersionInSubdir  = true
baseURL = "https://example.org/"
-- content/foo/p1.md --
---
aliases: ["/foo/p2/", "../p3/"]
---
-- content/foo/p2.md --
-- content/p3.md --
-- layouts/home.html --
Home.
{{ range $p := site.RegularPages }}{{ range .Aliases }}{{ . | printf "%-35s" }}=>{{ $p.RelPermalink -}}|{{ end -}}{{ end }}
`
	b := Test(t, files)

	b.AssertFileContent("public/guest/v1.0.0/en/index.html",
		"/guest/v1.0.0/en/foo/p2            =>/guest/v1.0.0/en/foo/p1/|",
		"/guest/v1.0.0/en/p3                =>/guest/v1.0.0/en/foo/p1/|",
	)
}

// Issue 14482
func TestOutputFormatIsHTMLWithMultilangAliases(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguageInSubdir = true
disableKinds = ['page', 'section', 'rss', 'sitemap', 'taxonomy', 'term']

[languages.en]
[languages.de]

[outputFormats.foo]
baseName = 'foo'
isHTML = false
mediaType = 'text/html'

[outputs]
home = ['html', 'foo']
-- content/index.de.md --
---
title: home de
---
-- content/index.en.md --
---
title: home en
---
-- layouts/home.html --
Output format: html|Title: {{ .Title }}
-- layouts/home.foo.html --
Output format: foo|Title: {{ .Title }}
`

	b := Test(t, files)

	b.AssertFileContent("public/en/index.html", `Output format: html|Title: home en`)
	b.AssertFileContent("public/en/foo.html", `Output format: foo|Title: home en`)
	b.AssertFileContent("public/de/index.html", `Output format: html|Title: home de`)
	b.AssertFileContent("public/de/foo.html", `Output format: foo|Title: home de`)
	b.AssertFileContent("public/index.html", `
<head>
  <title>/en/</title>
  <link rel="canonical" href="/en/">
  <meta charset="utf-8">
  <meta http-equiv="refresh" content="0; url=/en/">
</head>
	`)

	files = strings.ReplaceAll(files, "isHTML = false", "isHTML = true")

	b = Test(t, files)

	b.AssertFileContent("public/en/index.html", `Output format: html|Title: home en`)
	b.AssertFileContent("public/en/foo.html", `Output format: foo|Title: home en`)
	b.AssertFileContent("public/de/index.html", `Output format: html|Title: home de`)
	b.AssertFileContent("public/de/foo.html", `Output format: foo|Title: home de`)
	b.AssertFileContent("public/index.html", `
<head>
  <title>/en/</title>
  <link rel="canonical" href="/en/">
  <meta charset="utf-8">
  <meta http-equiv="refresh" content="0; url=/en/">
</head>
	`)
}
