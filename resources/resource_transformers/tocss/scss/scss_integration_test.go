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

package scss_test

import (
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/scss"
)

func TestTransformIncludePaths(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}
	c := qt.New(t)

	files := `
-- assets/scss/main.scss --
@import "moo";
-- node_modules/foo/_moo.scss --
$moolor: #fff;

moo {
  color: $moolor;
}
-- config.toml --
-- layouts/index.html --
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo") ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           c,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", `T1: moo{color:#fff}`)
}

func TestTransformImportRegularCSS(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}

	c := qt.New(t)

	files := `
-- assets/scss/_moo.scss --
$moolor: #fff;

moo {
	color: $moolor;
}
-- assets/scss/another.css --

-- assets/scss/main.scss --
@import "moo";
@import "regular.css";
@import "moo";
@import "another.css";

/* foo */
-- assets/scss/regular.css --

-- config.toml --
-- layouts/index.html --
{{ $r := resources.Get "scss/main.scss" |  toCSS }}
T1: {{ $r.Content | safeHTML }}

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           c,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	// LibSass does not support regular CSS imports. There
	// is an open bug about it that probably will never be resolved.
	// Hugo works around this by preserving them in place:
	b.AssertFileContent("public/index.html", `
 T1: moo {
 color: #fff; }

@import "regular.css";
moo {
 color: #fff; }

@import "another.css";
/* foo */

`)
}

func TestTransformThemeOverrides(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}

	c := qt.New(t)

	files := `
-- assets/scss/components/_boo.scss --
$boolor: green;

boo {
	color: $boolor;
}
-- assets/scss/components/_moo.scss --
$moolor: #ccc;

moo {
	color: $moolor;
}
-- config.toml --
theme = 'mytheme'
-- layouts/index.html --
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo" ) ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}
-- themes/mytheme/assets/scss/components/_boo.scss --
$boolor: orange;

boo {
	color: $boolor;
}
-- themes/mytheme/assets/scss/components/_imports.scss --
@import "moo";
@import "_boo";
@import "_zoo";
-- themes/mytheme/assets/scss/components/_moo.scss --
$moolor: #fff;

moo {
	color: $moolor;
}
-- themes/mytheme/assets/scss/components/_zoo.scss --
$zoolor: pink;

zoo {
	color: $zoolor;
}
-- themes/mytheme/assets/scss/main.scss --
@import "components/imports";
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           c,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", `T1: moo{color:#ccc}boo{color:green}zoo{color:pink}`)
}

func TestTransformErrors(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}

	c := qt.New(t)

	const filesTemplate = `
-- config.toml --
theme = 'mytheme'
-- assets/scss/components/_foo.scss --
/* comment line 1 */
$foocolor: #ccc;

foo {
	color: $foocolor;
}
-- themes/mytheme/assets/scss/main.scss --
/* comment line 1 */
/* comment line 2 */
@import "components/foo";
/* comment line 4 */

$maincolor: #eee;

body {
	color: $maincolor;
}

-- layouts/index.html --
{{ $cssOpts := dict }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}

	`

	c.Run("error in main", func(c *qt.C) {
		b, err := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           c,
				TxtarString: strings.Replace(filesTemplate, "$maincolor: #eee;", "$maincolor #eee;", 1),
				NeedsOsFS:   true,
			}).BuildE()

		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, filepath.FromSlash(`themes/mytheme/assets/scss/main.scss:6:1": expected ':' after $maincolor in assignment statement`))
		fe := b.AssertIsFileError(err)
		b.Assert(fe.ErrorContext(), qt.IsNotNil)
		b.Assert(fe.ErrorContext().Lines, qt.DeepEquals, []string{"/* comment line 4 */", "", "$maincolor #eee;", "", "body {"})
		b.Assert(fe.ErrorContext().ChromaLexer, qt.Equals, "scss")
	})

	c.Run("error in import", func(c *qt.C) {
		b, err := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           c,
				TxtarString: strings.Replace(filesTemplate, "$foocolor: #ccc;", "$foocolor #ccc;", 1),
				NeedsOsFS:   true,
			}).BuildE()

		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, `assets/scss/components/_foo.scss:2:1": expected ':' after $foocolor in assignment statement`)
		fe := b.AssertIsFileError(err)
		b.Assert(fe.ErrorContext(), qt.IsNotNil)
		b.Assert(fe.ErrorContext().Lines, qt.DeepEquals, []string{"/* comment line 1 */", "$foocolor #ccc;", "", "foo {"})
		b.Assert(fe.ErrorContext().ChromaLexer, qt.Equals, "scss")
	})
}

func TestOptionVars(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}

	files := `
-- assets/scss/main.scss --
@import "hugo:vars";

body {
	body {
		background: url($image) no-repeat center/cover;
		font-family: $font;
	  }
}

p {
	color: $color1;
	font-size: var$font_size;
}

b {
	color: $color2;
}
-- layouts/index.html --
{{ $image := "images/hero.jpg" }}
{{ $font := "Hugo's New Roman" }}
{{ $vars := dict "$color1" "blue" "$color2" "green" "font_size" "24px" "image" $image "font" $font }}
{{ $cssOpts := (dict "transpiler" "libsass" "outputStyle" "compressed" "vars" $vars ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts }}
T1: {{ $r.Content }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", `T1: body body{background:url(images/hero.jpg) no-repeat center/cover;font-family:Hugo&#39;s New Roman}p{color:blue;font-size:var 24px}b{color:green}`)
}

// Note: This test is more or less duplicated in both of the SCSS packages (libsass and dartsass).
func TestBootstrap(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}
	if !htesting.IsCI() {
		t.Skip("skip (slow) test in non-CI environment")
	}

	files := `
-- hugo.toml --
disableKinds = ["term", "taxonomy", "section", "page"]
[module]
[[module.imports]]
path="github.com/gohugoio/hugo-mod-bootstrap-scss/v5"
-- go.mod --
module github.com/gohugoio/tests/testHugoModules
-- assets/scss/main.scss --
@import "bootstrap/bootstrap";
-- layouts/index.html --
{{ $cssOpts := (dict "transpiler" "libsass" ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts }}
Styles: {{ $r.RelPermalink }}
		`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", "Styles: /scss/main.css")
}

// Issue #1239.
func TestRebuildAssetGetMatch(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}

	files := `
-- assets/scss/main.scss --
b {
	color: red;
}
-- layouts/index.html --
{{ $r := resources.GetMatch "scss/main.scss" |  toCSS  }}
T1: {{ $r.Content }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Running:     true,
		}).Build()

	b.AssertFileContent("public/index.html", `color: red`)

	b.EditFiles("assets/scss/main.scss", `b { color: blue; }`).Build()

	b.AssertFileContent("public/index.html", `color: blue`)
}

func TestRebuildAssetMatchIssue12456(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}

	files := `
-- hugo.toml --
disableKinds = ["term", "taxonomy", "section", "page"]
disableLiveReload = true
-- assets/a.scss --
h1 {
	color: red;
}
-- assets/dir/b.scss --
h2 {
	color: blue;
}
-- assets/dir/c.scss --
h3 {
	color: green;
}
-- layouts/index.html --
{{ $a := slice (resources.Get "a.scss") }}
{{ $b := resources.Match "dir/*.scss" }}

{{/* Add styles in a specific order. */}}
{{ $styles := slice $a $b }}

{{ $stylesheets := slice }}
  {{ range $styles }}
  {{ $stylesheets = $stylesheets | collections.Append . }}
{{ end }}


{{ range $stylesheets }}
  {{ with . | css.Sass | fingerprint }}
    <link as="style"  href="{{ .RelPermalink }}" rel="preload stylesheet">
  {{ end }}
{{ end }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Running:     true,
			// LogLevel:    logg.LevelTrace,
		}).Build()

	b.AssertFileContent("public/index.html", `b.60a9f3bdc189ee8a857afd5b7e1b93ad1644de0873761a7c9bc84f781a821942.css`)

	b.EditFiles("assets/dir/b.scss", `h2 { color: orange; }`).Build()

	b.AssertFileContent("public/index.html", `b.46b2d77c7ffe37ee191678f72df991ecb1319f849957151654362f09b0ef467f.css`)
}

// Issue 12851
func TestDirectoryIndexes(t *testing.T) {
	t.Parallel()
	if !scss.Supports() {
		t.Skip()
	}

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','sitemap','taxonomy','term']

[[module.mounts]]
source = 'assets'
target = 'assets'
[[module.mounts]]
source = "miscellaneous/sass"
target = "assets/sass"
-- layouts/index.html --
{{ $opts := dict "transpiler" "libsass" "outputStyle" "compressed" }}
{{ (resources.Get "sass/main.scss" | toCSS $opts).Content }}
-- assets/sass/main.scss --
@import "foo1"; // directory with _index file from OS file system
@import "bar1"; // directory with _index file from module mount
@import "foo2"; // directory with index file from OS file system
@import "bar2"; // directory with index file from module mount
-- assets/sass/foo1/_index.scss --
.foo1 {color: red;}
-- miscellaneous/sass/bar1/_index.scss --
.bar1 {color: blue;}
-- assets/sass/foo2/index.scss --
.foo2 {color: red;}
-- miscellaneous/sass/bar2/index.scss --
.bar2 {color: blue;}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			NeedsOsFS:   true,
			TxtarString: files,
		}).Build()

	b.AssertFileContent("public/index.html", ".foo1{color:red}.bar1{color:blue}.foo2{color:red}.bar2{color:blue}")
}
