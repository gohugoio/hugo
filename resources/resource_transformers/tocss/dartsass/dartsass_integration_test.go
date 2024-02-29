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

package dartsass_test

import (
	"strings"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/dartsass"
)

func TestTransformIncludePaths(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

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
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo") "transpiler" "dartsass" ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", `T1: moo{color:#fff}`)
}

func TestTransformImportRegularCSS(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

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
{{ $r := resources.Get "scss/main.scss" |  toCSS (dict "transpiler" "dartsass")  }}
T1: {{ $r.Content | safeHTML }}

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	// Dart Sass does not follow regular CSS import, but they
	// get pulled to the top.
	b.AssertFileContent("public/index.html", `T1: @import "regular.css";
		@import "another.css";
		moo {
		  color: #fff;
		}

		moo {
		  color: #fff;
		}

		/* foo */`)
}

func TestTransformImportIndentedSASS(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

	files := `
-- assets/scss/_moo.sass --
#main
	color: blue
-- assets/scss/main.scss --
@import "moo";

/* foo */
-- config.toml --
-- layouts/index.html --
{{ $r := resources.Get "scss/main.scss" |  toCSS (dict "transpiler" "dartsass")  }}
T1: {{ $r.Content | safeHTML }}

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", "T1: #main {\n  color: blue;\n}\n\n/* foo */")
}

// Issue 10592
func TestTransformImportMountedCSS(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

	files := `
-- assets/main.scss --
@import "import-this-file.css";
@import "foo/import-this-mounted-file.css";
@import "compile-this-file";
@import "foo/compile-this-mounted-file";
a {color: main-scss;}
-- assets/_compile-this-file.css --
a {color: compile-this-file-css;}
-- assets/_import-this-file.css --
a {color: import-this-file-css;}
-- foo/_compile-this-mounted-file.css --
a {color: compile-this-mounted-file-css;}
-- foo/_import-this-mounted-file.css --
a {color: import-this-mounted-file-css;}
-- layouts/index.html --
{{- $opts := dict "transpiler" "dartsass" }}
{{- with resources.Get "main.scss" | toCSS $opts }}{{ .Content | safeHTML }}{{ end }}
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term','page','section']

[[module.mounts]]
source = 'assets'
target = 'assets'

[[module.mounts]]
source = 'foo'
target = 'assets/foo'
	`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
		@import "import-this-file.css";
		@import "foo/import-this-mounted-file.css";
		a {
			color: compile-this-file-css;
		}

		a {
			color: compile-this-mounted-file-css;
		}

		a {
			color: main-scss;
		}
	`)
}

func TestTransformThemeOverrides(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

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
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo" ) "transpiler" "dartsass" ) }}
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
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", `T1: moo{color:#ccc}boo{color:green}zoo{color:pink}`)
}

func TestTransformLogging(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

	files := `
-- assets/scss/main.scss --
@warn "foo";
@debug "bar";

-- config.toml --
disableKinds = ["term", "taxonomy", "section", "page"]
-- layouts/index.html --
{{ $cssOpts := (dict  "transpiler" "dartsass" ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts   }}
T1: {{ $r.Content }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			LogLevel:    logg.LevelInfo,
		}).Build()

	b.AssertLogMatches(`Dart Sass: foo`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:1:0: bar`)
}

func TestTransformErrors(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

	c := qt.New(t)

	const filesTemplate = `
-- config.toml --
-- assets/scss/components/_foo.scss --
/* comment line 1 */
$foocolor: #ccc;

foo {
	color: $foocolor;
}
-- assets/scss/main.scss --
/* comment line 1 */
/* comment line 2 */
@import "components/foo";
/* comment line 4 */

  $maincolor: #eee;

body {
	color: $maincolor;
}

-- layouts/index.html --
{{ $cssOpts := dict "transpiler" "dartsass" }}
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
		b.Assert(err.Error(), qt.Contains, `main.scss:8:13":`)
		b.Assert(err.Error(), qt.Contains, `: expected ":".`)
		fe := b.AssertIsFileError(err)
		b.Assert(fe.ErrorContext(), qt.IsNotNil)
		b.Assert(fe.ErrorContext().Lines, qt.DeepEquals, []string{"  $maincolor #eee;", "", "body {", "\tcolor: $maincolor;", "}"})
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
		b.Assert(err.Error(), qt.Contains, `_foo.scss:2:10":`)
		b.Assert(err.Error(), qt.Contains, `: expected ":".`)
		fe := b.AssertIsFileError(err)
		b.Assert(fe.ErrorContext(), qt.IsNotNil)
		b.Assert(fe.ErrorContext().Lines, qt.DeepEquals, []string{"/* comment line 1 */", "$foocolor #ccc;", "", "foo {"})
		b.Assert(fe.ErrorContext().ChromaLexer, qt.Equals, "scss")
	})
}

func TestOptionVars(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

	files := `
-- assets/scss/main.scss --
@use "hugo:vars";

body {
	body {
		background: url(vars.$image) no-repeat center/cover;
		font-family: vars.$font;
	  }
}

p {
	color: vars.$color1;
	font-size: vars.$font_size;
}

b {
	color: vars.$color2;
}
-- layouts/index.html --
{{ $image := "images/hero.jpg" }}
{{ $font := "Hugo's New Roman" }}
{{ $vars := dict "$color1" "blue" "$color2" "green" "font_size" "24px" "image" $image "font" $font }}
{{ $cssOpts := (dict "transpiler" "dartsass" "outputStyle" "compressed" "vars" $vars ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts }}
T1: {{ $r.Content }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", `T1: body body{background:url(images/hero.jpg) no-repeat center/cover;font-family:Hugo&#39;s New Roman}p{color:blue;font-size:24px}b{color:green}`)
}

func TestOptionVarsParams(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

	files := `
-- config.toml --
[params]
[params.sassvars]
color1 = "blue"
color2 = "green"
font_size = "24px"
image = "images/hero.jpg"
-- assets/scss/main.scss --
@use "hugo:vars";

body {
	body {
		background: url(vars.$image) no-repeat center/cover;
	  }
}

p {
	color: vars.$color1;
	font-size: vars.$font_size;
}

b {
	color: vars.$color2;
}
-- layouts/index.html --
{{ $vars := site.Params.sassvars}}
{{ $cssOpts := (dict "transpiler" "dartsass" "outputStyle" "compressed" "vars" $vars ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts }}
T1: {{ $r.Content }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", `T1: body body{background:url(images/hero.jpg) no-repeat center/cover}p{color:blue;font-size:24px}b{color:green}`)
}

func TestVarsCasting(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
		t.Skip()
	}

	files := `
-- config.toml --
disableKinds = ["term", "taxonomy", "section", "page"]

[params]
[params.sassvars]
color_hex = "#fff"
color_rgb = "rgb(255, 255, 255)"
color_hsl = "hsl(0, 0%, 100%)"
dimension = "24px"
percentage = "10%"
flex = "5fr"
name = "Hugo"
url = "https://gohugo.io"
integer = 32
float = 3.14
-- assets/scss/main.scss --
@use "hugo:vars";
@use "sass:meta";

@debug meta.type-of(vars.$color_hex);
@debug meta.type-of(vars.$color_rgb);
@debug meta.type-of(vars.$color_hsl);
@debug meta.type-of(vars.$dimension);
@debug meta.type-of(vars.$percentage);
@debug meta.type-of(vars.$flex);
@debug meta.type-of(vars.$name);
@debug meta.type-of(vars.$url);
@debug meta.type-of(vars.$not_a_number);
@debug meta.type-of(vars.$integer);
@debug meta.type-of(vars.$float);
@debug meta.type-of(vars.$a_number);
-- layouts/index.html --
{{ $vars := site.Params.sassvars}}
{{ $vars = merge $vars (dict "not_a_number" ("32xxx" | css.Quoted) "a_number" ("234" | css.Unquoted) )}}
{{ $cssOpts := (dict "transpiler" "dartsass" "vars" $vars ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts }}
T1: {{ $r.Content }}
		`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			LogLevel:    logg.LevelInfo,
		}).Build()

	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:3:0: color`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:4:0: color`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:5:0: color`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:6:0: number`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:7:0: number`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:8:0: number`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:9:0: string`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:10:0: string`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:11:0: string`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:12:0: number`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:13:0: number`)
	b.AssertLogMatches(`Dart Sass: .*assets.*main.scss:14:0: number`)
}

// Note: This test is more or less duplicated in both of the SCSS packages (libsass and dartsass).
func TestBootstrap(t *testing.T) {
	t.Parallel()
	if !dartsass.Supports() {
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
{{ $cssOpts := (dict "transpiler" "dartsass" ) }}
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
