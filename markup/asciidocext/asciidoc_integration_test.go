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

// Package asciidocext converts AsciiDoc to HTML using Asciidoctor
// external binary. The `asciidoc` module is reserved for a future golang
// implementation.

package asciidocext_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/markup_config"
)

type contentFile struct {
	// The path to the content file, relative to the project directory.
	path string
	// The title of the content file.
	title string
}

type publishedContentFile struct {
	// The path to the published content file, relative to the project directory.
	path string
	// The expected string in the content file used for content test assertions.
	match string
}

// Issue 9202
// Issue 10183
// Issue 10473
func TestWorkingFolderCurrentWithDiagrams(t *testing.T) {
	if ok, err := asciidocext.Supports(); !ok {
		t.Skip(err)
	}
	if ok, err := asciidocext.SupportsPlantUML(); !ok {
		t.Skip(err)
	}

	c := qt.New(t)

	defaultAsciiDocConfig := markup_config.Default.AsciiDocExt // shallow copy
	t.Cleanup(func() {
		resetAsciiDocConfig(defaultAsciiDocConfig)
	})

	t.Chdir(t.TempDir())

	// A functional test site with config files matching each of the tests
	// below can be found in the https://github.com/jmooring/hugo-testing
	// repository in the hugo-github-issue-14094 branch.

	// Each of the tests below verifies both the published path of the diagram
	// and the src attribute within the published content file.

	files := createContentFiles() + contentAdapter + layouts

	// Test 1 - Monolingual, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
	replacer := strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `true`,
		`UGLYURLS_S1`, `false`,
	)
	f := files + replacer.Replace(configFileWithPlaceholders)
	b := hugolib.Test(t, f)
	validatePublishedSite_1(b, c)

	// Test 2 - Monolingual, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/subdir/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `true`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_2(b, c)

	// Test 3 - Monolingual, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `true`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_3(b, c)

	// Test 4 - Monolingual, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `true`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_4(b, c)

	// Test 5 - Multilingual, single-host, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_5(b, c)

	// Test 6 - Multilingual, single-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/subdir/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_6(b, c)

	// Test 7 - Multilingual, single-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_7(b, c)

	// Test 8 - Multilingual, single-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/subdir/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_8(b, c)

	// Test 9 - Multilingual, single-host, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_9(b, c)

	// Test 10 - Multilingual, single-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/subdir/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_10(b, c)

	// Test 11 - Multilingual, single-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_11(b, c)

	// Test 12 - Multilingual, single-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://example.org/subdir/"`,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, ``,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_12(b, c)

	// Test 13 - Multilingual, multi-host, uglyURLs = false, without subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_13(b, c)

	// Test 14 - Multilingual, multi-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/subdir/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/subdir/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_14(b, c)

	// Test 15 - Multilingual, multi-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_15(b, c)

	// Test 16 - Multilingual, multi-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `false`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/subdir/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/subdir/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_16(b, c)

	// Test 17 - Multilingual, multi-host, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_17(b, c)

	// Test 18 - Multilingual, multi-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/subdir/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/subdir/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `false`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_18(b, c)

	// Test 19 - Multilingual, multi-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_19(b, c)

	// Test 20 - Multilingual, multi-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
	replacer = strings.NewReplacer(
		`BASE_URL_KEY_VALUE_PAIR`, ``,
		`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, `true`,
		`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.en.example.org/subdir/"`,
		`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, `baseURL = "https://.de.example.org/subdir/"`,
		`LANGUAGES.DE.DISABLED`, `false`,
		`UGLYURLS_S1`, `true`,
	)
	f = files + replacer.Replace(configFileWithPlaceholders)
	b = hugolib.Test(t, f)
	validatePublishedSite_20(b, c)
}

// osFileExists reports whether a file with the given path exists on the OS FS.
// We need to use os.Stat for existence checks on Asciidoctor-generated
// diagrams. They are written directly to the OS file system, making them
// invisible to the afero FS used by b.AssertFileExists.
func osFileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// createContentFiles returns a txtar representation of the content files.
func createContentFiles() string {
	contentFiles := []contentFile{
		{`content/_index.de.adoc`, `home_de`},
		{`content/_index.en.adoc`, `home_en`},
		{`content/s1/_index.de.adoc`, `s1_de`},
		{`content/s1/_index.en.adoc`, `s1_en`},
		{`content/s1/p1/index.de.adoc`, `s1_p1_de`},
		{`content/s1/p1/index.en.adoc`, `s1_p1_en`},
		{`content/s1/p2.de.adoc`, `s1_p2_de`},
		{`content/s1/p2.en.adoc`, `s1_p2_en`},
		{`content/s2/_content.gotmpl`, `s2_content_gotmpl`},
		{`content/s2/_index.de.adoc`, `s2_de`},
		{`content/s2/_index.en.adoc`, `s2_en`},
		{`content/Straßen/Frühling/_index.de.adoc`, `Straßen_Frühling_de`},
		{`content/Straßen/Frühling/_index.en.adoc`, `Straßen_Frühling_en`},
		{`content/Straßen/Frühling/Müll Brücke.de.adoc`, `Straßen_Frühling_Müll_Brücke_de`},
		{`content/Straßen/Frühling/Müll Brücke.en.adoc`, `Straßen_Frühling_Müll_Brücke_en`},
		{`content/Straßen/_index.de.adoc`, `Straßen_de`},
		{`content/Straßen/_index.en.adoc`, `Straßen_en`},
	}

	formatString := `
-- %[1]s --
---
title: %[2]q
---
[plantuml,%[2]q]
....
@startuml
%[2]q
@enduml
....`

	var txtarContentFiles strings.Builder
	for _, f := range contentFiles {
		txtarContentFiles.WriteString(fmt.Sprintf(formatString, f.path, f.title))
	}

	return txtarContentFiles.String()
}

var contentAdapter = `
-- content/s2/_content.gotmpl --
{{ $title := printf "s2_p1_%s" site.Language.Lang }}
{{ $markup := printf "[plantuml,%[1]q]\n....\n@startuml\n%[1]q\n@enduml\n...." $title }}
{{ $content := dict
  "mediaType" "text/asciidoc"
  "value" $markup
}}
{{ $page := dict
  "content" $content
  "kind" "page"
  "path" "p1"
  "title" $title
}}
{{ .EnableAllLanguages }}
{{ .AddPage $page }}`

var layouts = `
-- layouts/all.html --
{{ .Content }}|`

var configFileWithPlaceholders = `
-- hugo.toml --
BASE_URL_KEY_VALUE_PAIR
disableKinds = ["rss", "sitemap", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR

[languages.en]
  LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR
  disabled = false
  weight = 1

[languages.de]
  LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR
  disabled = LANGUAGES.DE.DISABLED
  weight = 2

[uglyurls]
  s1 = UGLYURLS_S1

[markup.asciidocext]
  extensions           = ["asciidoctor-diagram"]
  workingFolderCurrent = true

  [markup.asciidocext.attributes]
    diagram-cachedir = "resources/_gen/asciidoctor/diagrams"

[security.exec]
  allow = ["^((dart-)?sass|git|go|npx|postcss|tailwindcss|asciidoctor)$"]`

// validatePublishedSite_1 validates the published site created by Test 1.
func validatePublishedSite_1(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_2(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_3(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_4(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_5(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_6(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_7(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_8(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/index.html`, `src="home_en.png"`},
		{`public/s1/index.html`, `src="s1_en.png"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/s2/index.html`, `src="s2_en.png"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/home_en.png`,
		`public/s1/p1/s1_p1_en.png`,
		`public/s1/p2/s1_p2_en.png`,
		`public/s1/s1_en.png`,
		`public/s2/p1/s2_p1_en.png`,
		`public/s2/s2_en.png`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/straßen/frühling/Straßen_Frühling_en.png`,
		`public/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_9(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_10(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_11(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_12(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_13(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_14(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_15(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_16(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_17(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_18(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.png"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.png"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_19(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}

func validatePublishedSite_20(b *hugolib.IntegrationTestBuilder, c *qt.C) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.png"`},
		{`public/de/s1/index.html`, `src="s1_de.png"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.png"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.png"`},
		{`public/de/s2/index.html`, `src="s2_de.png"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.png"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.png"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.png"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.png"`},
		{`public/en/index.html`, `src="home_en.png"`},
		{`public/en/s1/index.html`, `src="s1_en.png"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.png"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.png"`},
		{`public/en/s2/index.html`, `src="s2_en.png"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.png"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.png"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.png"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.png"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.png`,
		`public/de/s1/p1/s1_p1_de.png`,
		`public/de/s1/p2/s1_p2_de.png`,
		`public/de/s1/s1_de.png`,
		`public/de/s2/p1/s2_p1_de.png`,
		`public/de/s2/s2_de.png`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.png`,
		`public/de/straßen/frühling/Straßen_Frühling_de.png`,
		`public/de/straßen/Straßen_de.png`,
		`public/en/home_en.png`,
		`public/en/s1/p1/s1_p1_en.png`,
		`public/en/s1/p2/s1_p2_en.png`,
		`public/en/s1/s1_en.png`,
		`public/en/s2/p1/s2_p1_en.png`,
		`public/en/s2/s2_en.png`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.png`,
		`public/en/straßen/frühling/Straßen_Frühling_en.png`,
		`public/en/straßen/Straßen_en.png`,
	} {
		c.Assert(osFileExists(path), qt.IsTrue)
	}
}
