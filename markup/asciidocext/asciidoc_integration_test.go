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

package asciidocext_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/asciidocext"
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

type testCase struct {
	name                           string
	baseURL                        string
	defaultContentLanguageInSubdir bool
	langEnBaseURL                  string
	langDeBaseURL                  string
	langDeDisabled                 bool
	uglyURLsS1                     bool
	validate                       func(b *hugolib.IntegrationTestBuilder)
}

// Issue 9202, Issue 10183, Issue 10473
func TestAsciiDocDiagrams(t *testing.T) {
	t.Skip() // see if this affects the "No space left on device" GiHub runner problem

	if !htesting.IsRealCI() {
		t.Skip()
	}
	if ok, err := asciidocext.Supports(); !ok {
		t.Skip(err)
	}
	if ok, err := asciidocext.SupportsGoATDiagrams(); !ok {
		t.Skip(err)
	}

	t.Cleanup(func() {
		resetDefaultAsciiDocExtConfig()
	})

	diagramCacheDir := t.TempDir() // we want this to persist for all tests

	files := createContentFiles() + contentAdapter + layouts

	testCases := []testCase{
		// Test 1 - Monolingual, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 1", "https://example.org/", false, "", "", true, false, validatePublishedSite_1},
		// Test 2 - Monolingual, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 2", "https://example.org/subdir/", false, "", "", true, false, validatePublishedSite_2},
		// // Test 3 - Monolingual, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 3", "https://example.org/", false, "", "", true, true, validatePublishedSite_3},
		// // Test 4 - Monolingual, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 4", "https://example.org/subdir/", false, "", "", true, true, validatePublishedSite_4},
		// // Test 5 - Multilingual, single-host, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 5", "https://example.org/", false, "", "", false, false, validatePublishedSite_5},
		// Test 6 - Multilingual, single-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 6", "https://example.org/subdir/", false, "", "", false, false, validatePublishedSite_6},
		// Test 7 - Multilingual, single-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 7", "https://example.org/", false, "", "", false, true, validatePublishedSite_7},
		// Test 8 - Multilingual, single-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 8", "https://example.org/subdir/", false, "", "", false, true, validatePublishedSite_8},
		// Test 9 - Multilingual, single-host, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 9", "https://example.org/", true, "", "", false, false, validatePublishedSite_9},
		// Test 10 - Multilingual, single-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 10", "https://example.org/subdir/", true, "", "", false, false, validatePublishedSite_10},
		// Test 11 - Multilingual, single-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 11", "https://example.org/", true, "", "", false, true, validatePublishedSite_11},
		// Test 12 - Multilingual, single-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 12", "https://example.org/subdir/", true, "", "", false, true, validatePublishedSite_12},
		// Test 13 - Multilingual, multi-host, uglyURLs = false, without subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 13", "", false, `baseURL = "https://.en.example.org/"`, `baseURL = "https://.de.example.org/"`, false, false, validatePublishedSite_13},
		// Test 14 - Multilingual, multi-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 14", "", false, `baseURL = "https://.en.example.org/subdir/"`, `baseURL = "https://.de.example.org/subdir/"`, false, false, validatePublishedSite_14},
		// Test 15 - Multilingual, multi-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 15", "", false, `baseURL = "https://.en.example.org/"`, `baseURL = "https://.de.example.org/"`, false, true, validatePublishedSite_15},
		// Test 16 - Multilingual, multi-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = false
		{"Test 16", "", false, `baseURL = "https://.en.example.org/subdir/"`, `baseURL = "https://.de.example.org/subdir/"`, false, true, validatePublishedSite_16},
		// Test 17 - Multilingual, multi-host, uglyURLs = false for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 17", "", true, `baseURL = "https://.en.example.org/"`, `baseURL = "https://.de.example.org/"`, false, false, validatePublishedSite_17},
		// Test 18 - Multilingual, multi-host, uglyURLs = false for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 18", "", true, `baseURL = "https://.en.example.org/subdir/"`, `baseURL = "https://.de.example.org/subdir/"`, false, false, validatePublishedSite_18},
		// Test 19 - Multilingual, multi-host, uglyURLs = true for s1, without subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 19", "", true, `baseURL = "https://.en.example.org/"`, `baseURL = "https://.de.example.org/"`, false, true, validatePublishedSite_19},
		// Test 20 - Multilingual, multi-host, uglyURLs = true for s1, with subdir in base URL, defaultContentLanguageInSubdir = true
		{"Test 20", "", true, `baseURL = "https://.en.example.org/subdir/"`, `baseURL = "https://.de.example.org/subdir/"`, false, true, validatePublishedSite_20},
	}

	for _, tc := range testCases {
		// Asciidoctor is really, really slow on Windows. Only run a few tests.
		isExcludedTest := htesting.IsRealCI() && runtime.GOOS == "windows" && !slices.Contains([]string{"Test 4", "Test 8", "Test 16"}, tc.name)
		if isExcludedTest {
			continue
		}

		t.Run(tc.name, func(t *testing.T) {
			baseURLKV := ""
			if tc.baseURL != "" {
				baseURLKV = fmt.Sprintf(`baseURL = %q`, tc.baseURL)
			}

			replacer := strings.NewReplacer(
				`BASE_URL_KEY_VALUE_PAIR`, baseURLKV,
				`DEFAULT_CONTENT_LANGUAGE_IN_SUBDIR`, fmt.Sprintf("%v", tc.defaultContentLanguageInSubdir),
				`LANGUAGES.EN.BASE_URL_KEY_VALUE_PAIR`, tc.langEnBaseURL,
				`LANGUAGES.DE.BASE_URL_KEY_VALUE_PAIR`, tc.langDeBaseURL,
				`LANGUAGES.DE.DISABLED`, fmt.Sprintf("%v", tc.langDeDisabled),
				`UGLYURLS_S1`, fmt.Sprintf("%v", tc.uglyURLsS1),
				`DIAGRAM_CACHEDIR`, diagramCacheDir,
			)

			f := files + replacer.Replace(configFileWithPlaceholders)

			tempDir := t.TempDir()
			t.Chdir(tempDir)

			b := hugolib.Test(t, f, hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
				cfg.NeedsOsFS = true
				cfg.WorkingDir = tempDir
			}))

			// Verify that asciidoctor-diagram is caching in the correct
			// location. Checking one file is sufficient.
			err := fileExistsOsFs(filepath.Join(diagramCacheDir, "filecache/misc/asciidoctor-diagram/home_en.svg"))
			if err != nil {
				t.Fatalf("unable to locate file in diagram cache: %v", err.Error())
			}

			tc.validate(b)
		})
	}
}

// fileExistsOsFs checks if a file exists on the OS file system.
func fileExistsOsFs(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the path %q does not exist", path)
		}
		return fmt.Errorf("unable to get file status for %q: %v", path, err)
	}
	if fi.IsDir() {
		return fmt.Errorf("the path %q is a directory, not a file", path)
	}
	return nil
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
title: %[2]s
---
[goat,%[2]s]
....
%[2]s
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
{{ $markup := printf "[goat,%[1]s]\n....\n%[1]s\n...." $title }}
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
disableKinds = ['rss', 'sitemap', 'taxonomy', 'term']
defaultContentLanguage = 'en'
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
  extensions           = ['asciidoctor-diagram']
  workingFolderCurrent = true

[caches.misc]
dir = 'DIAGRAM_CACHEDIR'

[security.exec]
  allow = ['^((dart-)?sass|git|go|npx|postcss|tailwindcss|asciidoctor)$']`

func validatePublishedSite_1(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_2(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_3(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_4(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_5(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_6(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_7(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_8(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/index.html`, `src="home_en.svg"`},
		{`public/s1/index.html`, `src="s1_en.svg"`},
		{`public/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/s2/index.html`, `src="s2_en.svg"`},
		{`public/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/home_en.svg`,
		`public/s1/p1/s1_p1_en.svg`,
		`public/s1/p2/s1_p2_en.svg`,
		`public/s1/s1_en.svg`,
		`public/s2/p1/s2_p1_en.svg`,
		`public/s2/s2_en.svg`,
		`public/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_9(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_10(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_11(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_12(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_13(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_14(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_15(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_16(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_17(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_18(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1/index.html`, `src="s1_p1_de.svg"`},
		{`public/de/s1/p2/index.html`, `src="s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1/index.html`, `src="s1_p1_en.svg"`},
		{`public/en/s1/p2/index.html`, `src="s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_19(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}

func validatePublishedSite_20(b *hugolib.IntegrationTestBuilder) {
	for _, file := range []publishedContentFile{
		{`public/de/index.html`, `src="home_de.svg"`},
		{`public/de/s1/index.html`, `src="s1_de.svg"`},
		{`public/de/s1/p1.html`, `src="p1/s1_p1_de.svg"`},
		{`public/de/s1/p2.html`, `src="p2/s1_p2_de.svg"`},
		{`public/de/s2/index.html`, `src="s2_de.svg"`},
		{`public/de/s2/p1/index.html`, `src="s2_p1_de.svg"`},
		{`public/de/straßen/frühling/index.html`, `src="Straßen_Frühling_de.svg"`},
		{`public/de/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_de.svg"`},
		{`public/de/straßen/index.html`, `src="Straßen_de.svg"`},
		{`public/en/index.html`, `src="home_en.svg"`},
		{`public/en/s1/index.html`, `src="s1_en.svg"`},
		{`public/en/s1/p1.html`, `src="p1/s1_p1_en.svg"`},
		{`public/en/s1/p2.html`, `src="p2/s1_p2_en.svg"`},
		{`public/en/s2/index.html`, `src="s2_en.svg"`},
		{`public/en/s2/p1/index.html`, `src="s2_p1_en.svg"`},
		{`public/en/straßen/frühling/index.html`, `src="Straßen_Frühling_en.svg"`},
		{`public/en/straßen/frühling/müll-brücke/index.html`, `src="Straßen_Frühling_Müll_Brücke_en.svg"`},
		{`public/en/straßen/index.html`, `src="Straßen_en.svg"`},
	} {
		b.AssertFileContent(file.path, file.match)
	}

	for _, path := range []string{
		`public/de/home_de.svg`,
		`public/de/s1/p1/s1_p1_de.svg`,
		`public/de/s1/p2/s1_p2_de.svg`,
		`public/de/s1/s1_de.svg`,
		`public/de/s2/p1/s2_p1_de.svg`,
		`public/de/s2/s2_de.svg`,
		`public/de/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_de.svg`,
		`public/de/straßen/frühling/Straßen_Frühling_de.svg`,
		`public/de/straßen/Straßen_de.svg`,
		`public/en/home_en.svg`,
		`public/en/s1/p1/s1_p1_en.svg`,
		`public/en/s1/p2/s1_p2_en.svg`,
		`public/en/s1/s1_en.svg`,
		`public/en/s2/p1/s2_p1_en.svg`,
		`public/en/s2/s2_en.svg`,
		`public/en/straßen/frühling/müll-brücke/Straßen_Frühling_Müll_Brücke_en.svg`,
		`public/en/straßen/frühling/Straßen_Frühling_en.svg`,
		`public/en/straßen/Straßen_en.svg`,
	} {
		b.AssertFileExists(path, true)
	}
}
