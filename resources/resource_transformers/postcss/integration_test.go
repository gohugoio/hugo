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
	"strings"
	"testing"

	jww "github.com/spf13/jwalterweatherman"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
)

func TestTransformPostCSS(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("Skip long running test when running locally")
	}

	c := qt.New(t)

	files := `
-- assets/css/components/a.css --
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
// https://github.com/gohugoio/hugo/issues/7656
console.error("package.json:", process.env.HUGO_FILE_PACKAGE_JSON );
console.error("PostCSS Config File:", process.env.HUGO_FILE_POSTCSS_CONFIG_JS );

module.exports = {
	plugins: [
	require('tailwindcss')
	]
}

`

	c.Run("Success", func(c *qt.C) {
		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:               c,
				NeedsOsFS:       true,
				NeedsNpmInstall: true,
				LogLevel:        jww.LevelInfo,
				TxtarString:     files,
			}).Build()

		b.AssertLogContains("Hugo Environment: production")
		b.AssertLogContains(filepath.FromSlash(fmt.Sprintf("PostCSS Config File: %s/postcss.config.js", b.Cfg.WorkingDir)))
		b.AssertLogContains(filepath.FromSlash(fmt.Sprintf("package.json: %s/package.json", b.Cfg.WorkingDir)))

		b.AssertFileContent("public/index.html", `
Styles RelPermalink: /css/styles.css
Styles Content: Len: 770875|
`)
	})

	c.Run("Error", func(c *qt.C) {
		s, err := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:               c,
				NeedsOsFS:       true,
				NeedsNpmInstall: true,
				TxtarString:     strings.ReplaceAll(files, "color: blue;", "@apply foo;"), // Syntax error
			}).BuildE()
		s.AssertIsFileError(err)
	})
}

// bookmark2
func TestIntegrationTestTemplate(t *testing.T) {
	c := qt.New(t)

	files := ``

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:               c,
			NeedsOsFS:       false,
			NeedsNpmInstall: false,
			TxtarString:     files,
		}).Build()

	b.Assert(true, qt.IsTrue)
}
