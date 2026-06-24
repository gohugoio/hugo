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

package cssjs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/herrors"
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

-- hugo.toml --
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
-- layouts/home.html --
{{ $options := dict "inlineImports" true }}
{{ $styles := resources.Get "css/styles.css" | css.PostCSS $options }}
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

		b := hugolib.Test(c, files, hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall(), hugolib.TestOptInfo(), hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
			cfg.WorkingDir = tempDir
		}))

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
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	if runtime.GOOS == "windows" {
		// TODO(bep) This has started to fail on Windows with Go 1.19 on GitHub Actions for some mysterious reason.
		t.Skip("Skip on Windows")
	}

	c := qt.New(t)

	b, err := hugolib.TestE(c, strings.ReplaceAll(postCSSIntegrationTestFiles, "color: blue;", "@apply foo;"), hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall())

	ferrs := herrors.UnwrapFileErrors(err)
	b.Assert(len(ferrs), qt.Equals, 2)
	b.Assert(err.Error(), qt.Contains, "a.css:4:2")
}

// #9895
func TestTransformPostCSSImportError(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)

	_, err := hugolib.TestE(c, strings.ReplaceAll(postCSSIntegrationTestFiles, `@import "components/all.css";`, `@import "components/doesnotexist.css";`), hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall(), hugolib.TestOptInfo())
	ferrs := herrors.UnwrapFileErrors(err)
	c.Assert(len(ferrs), qt.Equals, 2)
	c.Assert(err.Error(), qt.Contains, "styles.css:4:3")
	c.Assert(err.Error(), qt.Contains, filepath.FromSlash(`failed to resolve CSS @import "/css/components/doesnotexist.css"`))
}

func TestTransformPostCSSImporSkipInlineImportsNotFound(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)

	files := strings.ReplaceAll(postCSSIntegrationTestFiles, `@import "components/all.css";`, `@import "components/doesnotexist.css";`)
	files = strings.ReplaceAll(files, `{{ $options := dict "inlineImports" true }}`, `{{ $options := dict "inlineImports" true "skipInlineImportsNotFound" true }}`)

	s := hugolib.Test(c, files, hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall(), hugolib.TestOptInfo())

	s.AssertFileContent("public/css/styles.css", `@import "components/doesnotexist.css";`)
}

// Issue 9787
func TestTransformPostCSSResourceCacheWithPathInBaseURL(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)
	tempDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
	c.Assert(err, qt.IsNil)
	c.Cleanup(clean)

	for i := range 2 {
		files := postCSSIntegrationTestFiles

		if i == 1 {
			files = strings.ReplaceAll(files, "https://example.com", "https://example.com/foo")
			files = strings.ReplaceAll(files, "useResourceCacheWhen = 'never'", "	useResourceCacheWhen = 'always'")
		}

		b := hugolib.Test(c, files, hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall(), hugolib.TestOptInfo(), hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
			cfg.WorkingDir = tempDir
		}))

		b.AssertFileContent("public/index.html", `
Styles Content: Len: 770917
`)

	}
}

// See Issue 15039.
// See Issue 15040.
func TestTransformPostCSSConfigResolution(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
[[module.imports]]
path = "github.com/bep/hugo-mod-nop"
-- assets/css/styles.css --
body { color: red }
-- layouts/home.html --
{{ $styles := resources.Get "css/styles.css" | css.PostCSS }}
RelPermalink: {{ $styles.RelPermalink }}|HasBody: {{ in $styles.Content "color:" }}|
-- package.json --
{
  "devDependencies": {
    "postcss-cli": "11.0.0"
  }
}
-- go.mod --
module github.com/example/project

go 1.26

replace github.com/bep/hugo-mod-nop => ../external-module
-- ../external-module/go.mod --
module github.com/bep/hugo-mod-nop

go 1.26
-- ../external-module/CONFIG_FILE_NAME --
CONFIG_FILE_CONTENT
	`

	tests := []struct {
		name              string
		configFileName    string
		configFileContent string
	}{
		{
			name:              "mjs in module",
			configFileName:    "postcss.config.mjs",
			configFileContent: "export default {};\n",
		},
		{
			name:              "cjs in module",
			configFileName:    "postcss.config.cjs",
			configFileContent: "module.exports = {};\n",
		},
		{
			name:              "js in module",
			configFileName:    "postcss.config.js",
			configFileContent: "module.exports = {};\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			rootDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
			c.Assert(err, qt.IsNil)
			c.Cleanup(clean)

			projectDir := filepath.Join(rootDir, "project")
			moduleDir := filepath.Join(rootDir, "external-module")
			c.Assert(os.MkdirAll(projectDir, 0o755), qt.IsNil)
			c.Assert(os.MkdirAll(moduleDir, 0o755), qt.IsNil)

			f := strings.ReplaceAll(files, "CONFIG_FILE_NAME", tt.configFileName)
			f = strings.ReplaceAll(f, "CONFIG_FILE_CONTENT", tt.configFileContent)

			b := hugolib.Test(c, f,
				hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
					cfg.WorkingDir = projectDir
					cfg.NeedsOsFS = true
					cfg.NeedsNpmInstall = true
				}),
				hugolib.TestOptInfo(),
			)

			b.AssertFileContent("public/index.html",
				"RelPermalink: /css/styles.css|HasBody: true|",
			)
			b.AssertLogContains(tt.configFileName)
		})
	}
}

// See Issue 13987.
func TestTransformPostCSSESMConfigInModule(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)
	// Use htesting.CreateTempDir to get canonical paths on macOS
	// (/private/var/...); Node's --permission model rejects the symlinked
	// /var/folders/... form when crossing the project boundary.
	rootDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
	c.Assert(err, qt.IsNil)
	c.Cleanup(clean)

	projectDir := filepath.Join(rootDir, "project")
	moduleDir := filepath.Join(rootDir, "external-module")
	c.Assert(os.MkdirAll(projectDir, 0o755), qt.IsNil)
	c.Assert(os.MkdirAll(moduleDir, 0o755), qt.IsNil)

	files := `
-- hugo.toml --
disableKinds = ['taxonomy', 'term', 'page']
baseURL = "https://example.com"
[[module.imports]]
path = "github.com/bep/hugo-mod-nop"
-- assets/css/styles.css --
body { color: red }
-- layouts/home.html --
{{ $styles := resources.Get "css/styles.css" | css.PostCSS }}
RelPermalink: {{ $styles.RelPermalink }}|HasBody: {{ in $styles.Content "color:" }}|
-- content/_index.md --
---
title: home
---
-- package.json --
{
  "devDependencies": {
    "postcss-cli": "11.0.0",
    "postcss-import": "16.0.0"
  }
}
-- go.mod --
module github.com/example/project

go 1.20

replace github.com/bep/hugo-mod-nop => ../external-module
-- ../external-module/go.mod --
module github.com/bep/hugo-mod-nop

go 1.20
-- ../external-module/postcss.config.js --
import postcssImport from "postcss-import";
export default { plugins: [postcssImport()] };

`

	b := hugolib.Test(c, files,
		hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
			cfg.WorkingDir = projectDir
			cfg.NeedsOsFS = true
			cfg.NeedsNpmInstall = true
		}),
	)

	b.AssertFileContent("public/index.html",
		"RelPermalink: /css/styles.css|HasBody: true|",
	)
}

// Netlify stores its node_modules cache in the same tree as the Hugo file
// cache, so Node's resolver can walk up from a module's postcss.config.js and
// hit a node_modules outside the permission allow-list, aborting with
// ERR_ACCESS_DENIED instead of falling through to NODE_PATH. The restricted
// postcss-import below (an ancestor of the external module, never installed by
// us) forces that walk to fail; the build must still succeed by resolving the
// real postcss-import via NODE_PATH.
//
// See issue 15041.
func TestTransformPostCSSESMConfigAccessDenied(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)
	// Use htesting.CreateTempDir to get canonical paths on macOS
	// (/private/var/...); Node's --permission model rejects the symlinked
	// /var/folders/... form when crossing the project boundary.
	rootDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
	c.Assert(err, qt.IsNil)
	c.Cleanup(clean)

	projectDir := filepath.Join(rootDir, "project")
	moduleDir := filepath.Join(rootDir, "external-module")
	c.Assert(os.MkdirAll(projectDir, 0o755), qt.IsNil)
	c.Assert(os.MkdirAll(moduleDir, 0o755), qt.IsNil)

	files := `
-- hugo.toml --
disableKinds = ['taxonomy', 'term', 'page']
baseURL = "https://example.com"
[[module.imports]]
path = "github.com/bep/hugo-mod-nop"
-- assets/css/styles.css --
body { color: red }
-- layouts/home.html --
{{ $styles := resources.Get "css/styles.css" | css.PostCSS }}
RelPermalink: {{ $styles.RelPermalink }}|HasBody: {{ in $styles.Content "color:" }}|
-- content/_index.md --
---
title: home
---
-- package.json --
{
  "devDependencies": {
    "postcss-cli": "11.0.0",
    "postcss-import": "16.0.0"
  }
}
-- go.mod --
module github.com/example/project

go 1.20

replace github.com/bep/hugo-mod-nop => ../external-module
-- ../node_modules/postcss-import/package.json --
{ "name": "postcss-import", "version": "0.0.0-RESTRICTED", "main": "index.js" }
-- ../external-module/go.mod --
module github.com/bep/hugo-mod-nop

go 1.20
-- ../external-module/postcss.config.js --
import postcssImport from "postcss-import";
export default { plugins: [postcssImport()] };

`

	b := hugolib.Test(c, files,
		hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
			cfg.WorkingDir = projectDir
			cfg.NeedsOsFS = true
			cfg.NeedsNpmInstall = true
		}),
	)

	b.AssertFileContent("public/index.html",
		"RelPermalink: /css/styles.css|HasBody: true|",
	)
}
