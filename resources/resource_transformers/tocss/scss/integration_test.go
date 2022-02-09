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

package scss_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/scss"
)

func TestTransformIncludePaths(t *testing.T) {
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
			T:           c,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html", `T1: moo{color:#ccc}boo{color:green}zoo{color:pink}`)
}
