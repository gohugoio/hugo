// Copyright 2026 The Hugo Authors. All rights reserved.
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

package css_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/internal/js/esbuild"
)

func TestCSSBuild(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/a/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- assets/css/main.css --
@import "tailwindcss";
@import url('https://example.org/foo.css');
@import "./foo.css";
@import "./bar.css" layer(mylayer);
@import "./bar.css" layer;
@import './baz.css' screen and (min-width: 800px);
@import 'relnodot.css' supports(display: grid);
@import 'relnodotboth.css';
@import '/relnodotboth2.css';
body {
  background-image: url("a/pixel.png");
}
.mask1 {
mask-image: url(a/pixel.png);
mask-repeat: no-repeat;
}
-- assets/css/foo.css --
p { background: red; }
-- assets/css/bar.css --
div { background: blue; }
-- assets/css/baz.css --
span { background: green; }
-- assets/css/relnodot.css --
article { background: yellow; }
-- assets/css/relnodotboth.css --
.relnodotbothrelative { background: green; }
-- assets/relnodotboth.css --
.relnodotbothroot { background: blue; }
-- assets/css/relnodotboth2.css --
.relnodotboth2relative { background: green; }
-- assets/relnodotboth2.css --
.relnodotboth2root { background: blue; }
-- layouts/home.html --
{{ with resources.Get "css/main.css"  }}
{{ $opts := (dict "minify" true "target" (slice "chrome108" "firefox116" "safari16.4" "edge115" "ios16.4" "opera101")  "loaders" (dict ".png" "dataurl") "externals" (slice  "tailwindcss"))}}
{{ with . | css.Build $opts  }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
	`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css",
		`-webkit-mask-repeat`,                      // browser prefix.
		"data:image/png;base64",                    // dataurl loader
		`@import"tailwindcss";`,                    // external
		`@import"https://example.org/foo.css";`,    // external
		"mylayer{div",                              // imported layer
		"@media screen and (min-width:800px){span", // imported media query
		"@supports (display: grid){article",        // imported supports
		"relnodotbothrelative",                     //@import 'relnodotboth.css'; file exists in both assets and project root.
		"relnodotboth2root",                        //@import '/relnodotboth2.css'; file exists in both assets and project root.
	)
}

func TestCSSBuildEdit(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss"]
-- assets/images/bar.svg --
barsvg
-- assets/images/foo.svg --
foosvg
-- assets/css/main.css --
@import "./foo.css";
@import "./bar.css" layer(mylayer);
.main { background: red; }
-- assets/css/foo.css --
.foo { background: green; }
-- assets/css/bar.css --
.bar { background-image: url("images/bar.svg"); }
-- layouts/all.html --
All. No CSS here.
-- layouts/home.html --
{{ with resources.Get "css/main.css"  }}
{{ $opts := (dict "minify" true) }}
{{ with . | css.Build $opts  }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`
	b := hugolib.TestRunning(t, files, hugolib.TestOptOsFs())

	b.AssertFileContent("public/css/main.css", `.foo{background:green}@layer mylayer{.bar{background-image:url("./bar-Y35ORVQM.svg")}}`)

	// Edit svg
	b.EditFileReplaceAll("assets/images/bar.svg", "barsvg", "newbarsvg").Build()
	b.AssertRenderCountPage(1)
	b.AssertFileContent("public/css/main.css", `bar-LVHHRPN5.svg`) // new hash.
	b.AssertFileContent("public/css/bar-LVHHRPN5.svg", "newbarsvg")

	// Edit foo.css
	b.EditFileReplaceAll("assets/css/foo.css", "green", "red").Build()
	b.AssertRenderCountPage(1)
	b.AssertFileContent("public/css/main.css", `.foo{background:red}`) // updated content.

	// Edit bar.css
	b.EditFileReplaceAll("assets/css/bar.css", "bar.svg", "foo.svg").Build()
	b.AssertRenderCountPage(1)
	b.AssertFileContent("public/css/main.css", `foo-52JTT5GU.svg`)
	b.AssertFileContent("public/css/foo-52JTT5GU.svg", "foosvg")

	// Edit main.css
	b.EditFileReplaceAll("assets/css/main.css", "red", "blue").Build()
	b.AssertRenderCountPage(1)
	b.AssertFileContent("public/css/main.css", `.main{background:#00f}`) // updated content, blue gets minified to #00f.
}

func TestCSSBuildEditOptions(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss"]
disableLiveReload = true
-- assets/css/main.css --
body {
 background: red;
}
-- layouts/home.html --
Home.
{{ with resources.Get "css/main.css"  }}
{{ $opts := dict "minify" false }}
{{ with . | css.Build $opts  }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	b := hugolib.TestRunning(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css", `  background: red;`)

	b.EditFileReplaceAll("layouts/home.html", `"minify" false`, `"minify" true`).Build()
	b.AssertFileContent("public/css/main.css", `{background:red}`)
}

func TestCSSBuildEditOptionsMultiple(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss"]
disableLiveReload = true
-- assets/css/main.css --
body {
 background: red;
}
-- layouts/_partials/css.html --
{{ with resources.Get "css/main.css"  }}
{{ $opts := dict "minify" false }}
{{ with . | css.Build $opts  }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
-- layouts/all.html --
All. {{ partial "css.html" . }}
-- content/p1.md --
-- content/p2.md --
-- content/p3.md --


`

	b := hugolib.TestRunning(t, files, hugolib.TestOptOsFs())

	for range 3 {
		b.AssertFileContent("public/css/main.css", `  background: red;`)
		b.EditFileReplaceAll("layouts/_partials/css.html", `"minify" false`, `"minify" true`).Build()
		b.AssertFileContent("public/css/main.css", `{background:red}`)
		b.EditFileReplaceAll("layouts/_partials/css.html", `"minify" true`, `"minify" false`).Build()
	}
}

func TestCSSBuildSourceMaps(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
-- assets/a/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- assets/css/main.css --
@import "./foo.css";
@import "./bar.css" layer(mylayer);
@import "./bar.css" layer;
@import './baz.css' screen and (min-width: 800px);

.main {
  background-color: red;
}
-- assets/css/foo.css --
p { background: red; }
-- assets/css/bar.css --
div { background: blue; }
-- assets/css/baz.css --
span { background: green; }
-- assets/css/qux.css --
article { background: yellow; }
-- layouts/home.html --
{{ with resources.Get "css/main.css"  }}
{{ $opts := (dict "minify" MINIFY "sourceMap" "SOURCE_MAP" "sourcesContent" SOURCES_CONTENT  "loaders" (dict ".png" "dataurl"))}}
{{ with . | css.Build $opts  }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
	`

	var (
		r     *strings.Replacer
		files string
		b     *hugolib.IntegrationTestBuilder
	)

	for _, minify := range []string{"true", "false"} {

		r = strings.NewReplacer(
			"SOURCE_MAP", "linked",
			"SOURCES_CONTENT", "true",
			"MINIFY", minify,
		)

		files = r.Replace(filesTemplate)

		b = hugolib.Test(t, files, hugolib.TestOptOsFs())
		b.AssertFileContent("public/css/main.css", "/*# sourceMappingURL=main.css.map */")
		b.AssertFileContent("public/css/main.css.map",
			`"sourcesContent":["`,
			`"mappings":"`,
		)

		sources := esbuild.SourcesFromSourceMap(b.FileContent("public/css/main.css.map"))
		// main.css + foo.css + bar.css + baz.css = 4 sources.
		b.Assert(len(sources), qt.Equals, 4)
	}

	r = strings.NewReplacer(
		"SOURCE_MAP", "external",
		"SOURCES_CONTENT", "true",
		"MINIFY", "false",
	)

	files = r.Replace(filesTemplate)

	b = hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css", "! sourceMappingURL")
	b.AssertFileContent("public/css/main.css.map",
		`"sourcesContent":["p { background: red; }"`,
		"AAAA;AAAI,cAAY;AAAK",
	)

	r = strings.NewReplacer(
		"SOURCE_MAP", "external",
		"SOURCES_CONTENT", "false",
		"MINIFY", "false",
	)

	files = r.Replace(filesTemplate)

	b = hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css.map",
		`"sourcesContent":null,"`,
		"AAAA;AAAI,cAAY;AAAK",
	)

	r = strings.NewReplacer(
		"SOURCE_MAP", "inline",
		"SOURCES_CONTENT", "false",
		"MINIFY", "false",
	)

	files = r.Replace(filesTemplate)

	b = hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileExists("public/css/main.css.map", false)
	b.AssertFileContent("public/css/main.css", "sourceMappingURL=data:application/json;base64,")
}

func TestCSSBuildMultihost(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
baseURL = "https://example.en"
weight = 1
[languages.fr]
baseURL = "https://example.fr"
weight = 2
-- assets/a/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- assets/css/main.css --
body {
  background-image: url("a/pixel.png");
  color: red;
}
-- layouts/home.html --
{{ with resources.Get "css/main.css"  }}
{{ with . | css.Build (dict "minify" true)  }}
CSS: {{ .RelPermalink }}|{{ .Content }}
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())

	for _, lang := range []string{"en", "fr"} {
		b.AssertFileContent("public/"+lang+"/css/main.css", `./pixel-NJRUOINY.png`)
		b.AssertFileExists("public/"+lang+"/css/pixel-NJRUOINY.png", true)
	}
}

func TestCSSBuildLoadersDefault(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/a/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- static/b/issue14619.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- assets/css/main.css --
body {
  background-image: url("a/pixel.png");
}
div {
  background-image: url("static/b/issue14619.png");
}
-- layouts/home.html --
{{ with resources.Get "css/main.css"  }}
{{ with . | css.Build (dict "minify" true)  }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css", `./pixel-NJRUOINY.png`)
	b.AssertFileExists("public/css/pixel-NJRUOINY.png", true)
	b.AssertFileContent("public/css/main.css", `url("./issue14619-NJRUOINY.png")`)
	b.AssertFileExists("public/css/issue14619-NJRUOINY.png", true)
}

// Issue #14623
func TestCSSBuildLoadersPartial(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/home.html --
{{ $opts := dict "loaders" (dict ".svg" "dataurl") }}
{{ with resources.Get "css/main.css" | css.Build $opts }}
  <link rel="stylesheet" href="{{ .RelPermalink }}">
{{ end }}
-- assets/css/main.css --
body { background-image: url('images/pixel.png'); }
div { background-image: url('images/foo.svg'); }
-- assets/images/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- assets/images/foo.svg --
SVG file.
`
	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileExists("public/css/pixel-NJRUOINY.png", true)
	b.AssertFileContent("public/css/main.css", `url(data:image/svg`) // svg should be inlined as dataurl.
	b.AssertFileContent("public/css/main.css", `pixel-NJRUOINY.png`) // png should be referenced as a hashed file, not inlined.
}

func TestCSSBuildBootstrapFromNPM(t *testing.T) {
	htesting.SkipSlowTestUnlessCI(t)
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
-- assets/css/main.css --
@import "bootstrap";
-- package.json --
{
	"devDependencies": {
		"bootstrap": "5.3.8"
	}
}
-- layouts/home.html --
{{ with resources.Get "css/main.css"  }}
{{ with . | css.Build (dict "minify" true)  }}
	 CSS size: {{ .Content | len }}|{{ .RelPermalink }}
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs(), hugolib.TestOptWithNpmInstall())
	b.AssertFileContent("public/css/main.css", `--bs-indigo: #6610f2;`)
}

func TestCSSBuildVars(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/css/main.css --
@import "hugo:vars";

body {
  background-color: var(--primary-color);
  font-size: var(--font-size);
  color: var(--text-color);
}
-- layouts/home.html --
{{ with resources.Get "css/main.css" }}
{{ $opts := dict
  "vars" (dict "primary-color" "blue" "font-size" "24px" "text-color" "#333" "--already-prefixed" "red")
}}
{{ with . | css.Build $opts }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css",
		"--primary-color: blue;",
		"--font-size: 24px;",
		"--text-color: #333;",
		"--already-prefixed: red;",
		"background-color: var(--primary-color)",
	)
}

func TestCSSBuildVarsNestedIssue14705(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/css/main.css --
@import "hugo:vars";
@import "hugo:vars/mobile" (max-width: 650px);

body {
  background-color: var(--primary-color);
}
-- layouts/home.html --
{{ with resources.Get "css/main.css" }}
{{ $opts := dict
  "vars" (dict "primary-color" "blue" "mobile" (dict "primary-color" "red" "font-size" "12px"))
}}
{{ with . | css.Build $opts }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css",
		"--primary-color: blue;",
		"@media (max-width: 650px)",
		"--primary-color: red;",
		"! --mobile:",
		"--font-size: 12px;",
		"background-color: var(--primary-color)",
	)
}

func TestCSSBuildVarsEmpty(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/css/main.css --
@import "hugo:vars";

body {
  background-color: red;
}
-- layouts/home.html --
{{ with resources.Get "css/main.css" }}
{{ with . | css.Build }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css", "background-color: red;")
}

func TestCSSBuildVarsNestedUpperCase(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
-- assets/css/main.css --
@import "hugo:vars";
@import "hugo:vars/MOBILE1" (max-width: 650px);

body {
  background-color: var(--primary-color);
}
-- layouts/home.html --
{{ with resources.Get "css/main.css" }}
{{ $opts := dict
  "vars" (dict "primary-color" "blue" "MOBILE2" (dict "primary-color" "red" "font-size" "12px" "MixedCaseKey" "value"))
}}
{{ with . | css.Build $opts }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	files := strings.ReplaceAll(filesTemplate, "MOBILE1", "Mobile")
	files = strings.ReplaceAll(files, "MOBILE2", "mobile")
	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css", "primary-color: red;")

	files = strings.ReplaceAll(filesTemplate, "MOBILE1", "mobile")
	files = strings.ReplaceAll(files, "MOBILE2", "MobilE")
	b = hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css", "primary-color: red;")
}

func TestCSSBuildVarsQuoted(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/css/main.css --
@import "hugo:vars";

body {
  font-family: var(--font-family);
  color: var(--brand-color);
}
-- layouts/home.html --
{{ with resources.Get "css/main.css" }}
{{ $opts := dict
  "vars" (dict "font-family" (css.Quoted "Arial, sans-serif") "brand-color" "hsl(0, 0%, 20%)")
}}
{{ with . | css.Build $opts }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css",
		`--font-family: "Arial, sans-serif";`,
		"--brand-color: hsl(0, 0%, 20%);",
	)
}

func TestCSSBuildVarsFromParams(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[params.styles]
primary-color = "blue"
font-size = "24px"
text-color = "#333"
-- assets/css/main.css --
@import "hugo:vars";
-- assets/css/main.css --
@import "hugo:vars";

body {
  background-color: var(--primary-color);
  font-size: var(--font-size);
  color: var(--text-color);
}
-- layouts/home.html --
{{ with resources.Get "css/main.css" }}
{{ $opts := dict
  "vars" site.Params.styles
}}
{{ with . | css.Build $opts }}
 <link rel="stylesheet" href="{{ .RelPermalink }}" />
{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())
	b.AssertFileContent("public/css/main.css",
		"--primary-color: blue;",
		"--font-size: 24px;",
		"--text-color: #333;",
	)
}
