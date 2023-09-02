// Copyright 2021 The Hugo Authors. All rights reserved.
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

package postcss_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
)

const postCSSIntegrationTestFiles = `
-- assets/css/components/a.css --
/* A comment. */
/* Another comment. */
class-in-a {
	color: blue;
}

-- assets/css/components/all.css --
@import "a.css";
@import "b.css";
-- assets/css/components/b.css --
@import "a.css";

class-in-b {
	color: blue;
}

-- assets/css/styles.css --
@tailwind base;
@tailwind components;
@tailwind utilities;
  @import "components/all.css";
h1 {
	@apply text-2xl font-bold;
}

-- config.toml --
disablekinds = ['taxonomy', 'term', 'page']
baseURL = "https://example.com"
[build]
useResourceCacheWhen = 'never'
-- content/p1.md --
-- data/hugo.toml --
slogan = "Hugo Rocks!"
-- i18n/en.yaml --
hello:
   other: "Hello"
-- i18n/fr.yaml --
hello:
   other: "Bonjour"
-- layouts/index.html --
{{ $options := dict "inlineImports" true }}
{{ $styles := resources.Get "css/styles.css" | resources.PostCSS $options }}
Styles RelPermalink: {{ $styles.RelPermalink }}
{{ $cssContent := $styles.Content }}
Styles Content: Len: {{ len $styles.Content }}|
-- package.json --
{
	"scripts": {},

	"devDependencies": {
	"postcss-cli": "7.1.0",
	"tailwindcss": "1.2.0"
	}
}
-- postcss.config.js --
console.error("Hugo Environment:", process.env.HUGO_ENVIRONMENT );
console.error("Hugo PublishDir:", process.env.HUGO_PUBLISHDIR );
// https://github.com/gohugoio/hugo/issues/7656
console.error("package.json:", process.env.HUGO_FILE_PACKAGE_JSON );
console.error("PostCSS Config File:", process.env.HUGO_FILE_POSTCSS_CONFIG_JS );

module.exports = {
	plugins: [
	require('tailwindcss')
	]
}

`

func TestTransformPostCSS(t *testing.T) {
	t.Parallel()
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)
	tempDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
	c.Assert(err, qt.IsNil)
	c.Cleanup(clean)

	for _, s := range []string{"never", "always"} {

		repl := strings.NewReplacer(
			"https://example.com",
			"https://example.com/foo",
			"useResourceCacheWhen = 'never'",
			fmt.Sprintf("useResourceCacheWhen = '%s'", s),
		)

		files := repl.Replace(postCSSIntegrationTestFiles)

		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:               c,
				NeedsOsFS:       true,
				NeedsNpmInstall: true,
				LogLevel:        logg.LevelInfo,
				WorkingDir:      tempDir,
				TxtarString:     files,
			}).Build()

		b.AssertFileContent("public/index.html", `
Styles RelPermalink: /foo/css/styles.css
Styles Content: Len: 770917|
`)

		if s == "never" {
			b.AssertLogContains("Hugo Environment: production")
			b.AssertLogContains("Hugo PublishDir: " + filepath.Join(tempDir, "public"))
		}
	}

}

// 9880
func TestTransformPostCSSError(t *testing.T) {
	t.Parallel()
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	if runtime.GOOS == "windows" {
		//TODO(bep) This has started to fail on Windows with Go 1.19 on GitHub Actions for some mysterious reason.
		t.Skip("Skip on Windows")
	}

	c := qt.New(t)

	s, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               c,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			TxtarString:     strings.ReplaceAll(postCSSIntegrationTestFiles, "color: blue;", "@apply foo;"), // Syntax error
		}).BuildE()

	s.AssertIsFileError(err)
	c.Assert(err.Error(), qt.Contains, "a.css:4:2")

}

func TestTransformPostCSSNotInstalledError(t *testing.T) {
	t.Parallel()
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)

	s, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           c,
			NeedsOsFS:   true,
			TxtarString: postCSSIntegrationTestFiles,
		}).BuildE()

	s.AssertIsFileError(err)
	c.Assert(err.Error(), qt.Contains, `binary with name "npx" not found`)

}

// #9895
func TestTransformPostCSSImportError(t *testing.T) {
	t.Parallel()
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)

	s, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               c,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			LogLevel:        logg.LevelInfo,
			TxtarString:     strings.ReplaceAll(postCSSIntegrationTestFiles, `@import "components/all.css";`, `@import "components/doesnotexist.css";`),
		}).BuildE()

	s.AssertIsFileError(err)
	c.Assert(err.Error(), qt.Contains, "styles.css:4:3")
	c.Assert(err.Error(), qt.Contains, filepath.FromSlash(`failed to resolve CSS @import "css/components/doesnotexist.css"`))

}

func TestTransformPostCSSImporSkipInlineImportsNotFound(t *testing.T) {
	t.Parallel()
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)

	files := strings.ReplaceAll(postCSSIntegrationTestFiles, `@import "components/all.css";`, `@import "components/doesnotexist.css";`)
	files = strings.ReplaceAll(files, `{{ $options := dict "inlineImports" true }}`, `{{ $options := dict "inlineImports" true "skipInlineImportsNotFound" true }}`)

	s := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               c,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			LogLevel:        logg.LevelInfo,
			TxtarString:     files,
		}).Build()

	s.AssertFileContent("public/css/styles.css", `@import "components/doesnotexist.css";`)

}

// Issue 9787
func TestTransformPostCSSResourceCacheWithPathInBaseURL(t *testing.T) {
	t.Parallel()
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)
	tempDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
	c.Assert(err, qt.IsNil)
	c.Cleanup(clean)

	for i := 0; i < 2; i++ {
		files := postCSSIntegrationTestFiles

		if i == 1 {
			files = strings.ReplaceAll(files, "https://example.com", "https://example.com/foo")
			files = strings.ReplaceAll(files, "useResourceCacheWhen = 'never'", "	useResourceCacheWhen = 'always'")
		}

		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:               c,
				NeedsOsFS:       true,
				NeedsNpmInstall: true,
				LogLevel:        logg.LevelInfo,
				TxtarString:     files,
				WorkingDir:      tempDir,
			}).Build()

		b.AssertFileContent("public/index.html", `
Styles Content: Len: 770917
`)

	}

}
