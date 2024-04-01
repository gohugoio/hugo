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

package hugolib

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"
)

func TestSitemapBasic(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
-- content/sect/doc1.md --
---
title: doc1
---
Doc1
-- content/sect/doc2.md --
---
title: doc2
---
Doc2
`

	b := Test(t, files)

	b.AssertFileContent("public/sitemap.xml", " <loc>https://example.com/sect/doc1/</loc>", "doc2")
}

func TestSitemapMultilingual(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
defaultContentLanguage = "en"
[languages]
[languages.en]
weight = 1
languageName = "English"
[languages.nn]
weight = 2
languageName = "Nynorsk"
-- content/sect/doc1.md --
---
title: doc1
---
Doc1
-- content/sect/doc2.md --
---
title: doc2
---
Doc2
-- content/sect/doc2.nn.md --
---
title: doc2
---
Doc2
`

	b := Test(t, files)

	b.AssertFileContent("public/sitemap.xml", "<loc>https://example.com/en/sitemap.xml</loc>", "<loc>https://example.com/nn/sitemap.xml</loc>")
	b.AssertFileContent("public/en/sitemap.xml", " <loc>https://example.com/sect/doc1/</loc>", "doc2")
	b.AssertFileContent("public/nn/sitemap.xml", " <loc>https://example.com/nn/sect/doc2/</loc>")
}

// https://github.com/gohugoio/hugo/issues/5910
func TestSitemapOutputFormats(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
-- content/blog/html-amp.md --
---
Title: AMP and HTML
outputs: [ "html", "amp" ]
---

`

	b := Test(t, files)

	// Should link to the HTML version.
	b.AssertFileContent("public/sitemap.xml", " <loc>https://example.com/blog/html-amp/</loc>")
}

func TestParseSitemap(t *testing.T) {
	t.Parallel()
	expected := config.SitemapConfig{ChangeFreq: "3", Disable: true, Filename: "doo.xml", Priority: 3.0}
	input := map[string]any{
		"changefreq": "3",
		"disable":    true,
		"filename":   "doo.xml",
		"priority":   3.0,
		"unknown":    "ignore",
	}
	result, err := config.DecodeSitemap(config.SitemapConfig{}, input)
	if err != nil {
		t.Fatalf("Failed to parse sitemap: %s", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Got \n%v expected \n%v", result, expected)
	}
}

func TestSitemapShouldNotUseListXML(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
[languages]
[languages.en]
weight = 1
languageName = "English"
[languages.nn]
weight = 2
-- layouts/_default/list.xml --
Site: {{ .Site.Title }}|
-- layouts/home --
Home.

`

	b := Test(t, files)

	b.AssertFileContent("public/sitemap.xml", "https://example.com/en/sitemap.xml")
}

func TestSitemapAndContentBundleNamedSitemap(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','taxonomy','term']
-- layouts/_default/single.html --
layouts/_default/single.html
-- layouts/sitemap/single.html --
layouts/sitemap/single.html
-- content/sitemap/index.md --
---
title: My sitemap
type: sitemap
---
`

	b := Test(t, files)

	b.AssertFileExists("public/sitemap.xml", true)
}

// Issue 12266
func TestSitemapIssue12266(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'https://example.org/'
disableKinds = ['rss','taxonomy','term']
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages.de]
[languages.en]
  `

	// Test A: multilingual with defaultContentLanguageInSubdir = true
	b := Test(t, files)

	b.AssertFileContent("public/sitemap.xml",
		"<loc>https://example.org/de/sitemap.xml</loc>",
		"<loc>https://example.org/en/sitemap.xml</loc>",
	)
	b.AssertFileContent("public/de/sitemap.xml", "<loc>https://example.org/de/</loc>")
	b.AssertFileContent("public/en/sitemap.xml", "<loc>https://example.org/en/</loc>")

	// Test B: multilingual with defaultContentLanguageInSubdir = false
	files = strings.ReplaceAll(files, "defaultContentLanguageInSubdir = true", "defaultContentLanguageInSubdir = false")

	b = Test(t, files)

	b.AssertFileContent("public/sitemap.xml",
		"<loc>https://example.org/de/sitemap.xml</loc>",
		"<loc>https://example.org/en/sitemap.xml</loc>",
	)
	b.AssertFileContent("public/de/sitemap.xml", "<loc>https://example.org/de/</loc>")
	b.AssertFileContent("public/en/sitemap.xml", "<loc>https://example.org/</loc>")

	// Test C: monolingual with defaultContentLanguageInSubdir = false
	files = strings.ReplaceAll(files, "[languages.de]", "")
	files = strings.ReplaceAll(files, "[languages.en]", "")

	b = Test(t, files)

	b.AssertFileExists("public/en/sitemap.xml", false)
	b.AssertFileContent("public/sitemap.xml", "<loc>https://example.org/</loc>")

	// Test D: monolingual with defaultContentLanguageInSubdir = true
	files = strings.ReplaceAll(files, "defaultContentLanguageInSubdir = false", "defaultContentLanguageInSubdir = true")

	b = Test(t, files)

	b.AssertFileContent("public/sitemap.xml", "<loc>https://example.org/en/sitemap.xml</loc>")
	b.AssertFileContent("public/en/sitemap.xml", "<loc>https://example.org/en/</loc>")
}
