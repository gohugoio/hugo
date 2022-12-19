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

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/dartsass"
	jww "github.com/spf13/jwalterweatherman"
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
			LogLevel:    jww.LevelInfo,
		}).Build()

	b.AssertLogMatches(`WARN.*Dart Sass: foo`)
	b.AssertLogMatches(`INFO.*Dart Sass: .*assets.*main.scss:1:0: bar`)

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
{{ $vars := dict "$color1" "blue" "$color2" "green" "font_size" "24px" "image" $image }}
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
