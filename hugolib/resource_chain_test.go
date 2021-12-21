// Copyright 2019 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/dartsass"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/htesting"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/scss"
)

func TestSCSSWithIncludePaths(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name     string
		supports func() bool
	}{
		{"libsass", func() bool { return scss.Supports() }},
		{"dartsass", func() bool { return dartsass.Supports() }},
	} {
		c.Run(test.name, func(c *qt.C) {
			if !test.supports() {
				c.Skip(fmt.Sprintf("Skip %s", test.name))
			}

			workDir, clean, err := htesting.CreateTempDir(hugofs.Os, fmt.Sprintf("hugo-scss-include-%s", test.name))
			c.Assert(err, qt.IsNil)
			defer clean()

			v := config.New()
			v.Set("workingDir", workDir)
			b := newTestSitesBuilder(c).WithLogger(loggers.NewErrorLogger())
			// Need to use OS fs for this.
			b.Fs = hugofs.NewDefault(v)
			b.WithWorkingDir(workDir)
			b.WithViper(v)

			fooDir := filepath.Join(workDir, "node_modules", "foo")
			scssDir := filepath.Join(workDir, "assets", "scss")
			c.Assert(os.MkdirAll(fooDir, 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "content", "sect"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "data"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "i18n"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "layouts", "shortcodes"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "layouts", "_default"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(scssDir), 0777), qt.IsNil)

			b.WithSourceFile(filepath.Join(fooDir, "_moo.scss"), `
$moolor: #fff;

moo {
  color: $moolor;
}
`)

			b.WithSourceFile(filepath.Join(scssDir, "main.scss"), `
@import "moo";

`)

			b.WithTemplatesAdded("index.html", fmt.Sprintf(`
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo") "transpiler" %q ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}
`, test.name))
			b.Build(BuildCfg{})

			b.AssertFileContent(filepath.Join(workDir, "public/index.html"), `T1: moo{color:#fff}`)
		})
	}
}

func TestSCSSWithRegularCSSImport(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name     string
		supports func() bool
	}{
		{"libsass", func() bool { return scss.Supports() }},
		{"dartsass", func() bool { return dartsass.Supports() }},
	} {
		c.Run(test.name, func(c *qt.C) {
			if !test.supports() {
				c.Skip(fmt.Sprintf("Skip %s", test.name))
			}

			workDir, clean, err := htesting.CreateTempDir(hugofs.Os, fmt.Sprintf("hugo-scss-include-regular-%s", test.name))
			c.Assert(err, qt.IsNil)
			defer clean()

			v := config.New()
			v.Set("workingDir", workDir)
			b := newTestSitesBuilder(c).WithLogger(loggers.NewErrorLogger())
			// Need to use OS fs for this.
			b.Fs = hugofs.NewDefault(v)
			b.WithWorkingDir(workDir)
			b.WithViper(v)

			scssDir := filepath.Join(workDir, "assets", "scss")
			c.Assert(os.MkdirAll(filepath.Join(workDir, "content", "sect"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "data"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "i18n"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "layouts", "shortcodes"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "layouts", "_default"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(scssDir), 0777), qt.IsNil)
			b.WithSourceFile(filepath.Join(scssDir, "regular.css"), ``)
			b.WithSourceFile(filepath.Join(scssDir, "another.css"), ``)
			b.WithSourceFile(filepath.Join(scssDir, "_moo.scss"), `
$moolor: #fff;

moo {
  color: $moolor;
}
`)

			b.WithSourceFile(filepath.Join(scssDir, "main.scss"), `
@import "moo";
@import "regular.css";
@import "moo";
@import "another.css";

/* foo */
`)

			b.WithTemplatesAdded("index.html", fmt.Sprintf(`
{{ $r := resources.Get "scss/main.scss" |  toCSS (dict "transpiler" %q)  }}
T1: {{ $r.Content | safeHTML }}
`, test.name))
			b.Build(BuildCfg{})

			if test.name == "libsass" {
				// LibSass does not support regular CSS imports. There
				// is an open bug about it that probably will never be resolved.
				// Hugo works around this by preserving them in place:
				b.AssertFileContent(filepath.Join(workDir, "public/index.html"), `
 T1: moo {
 color: #fff; }

@import "regular.css";
moo {
 color: #fff; }

@import "another.css";
/* foo */
        
`)
			} else {
				// Dart Sass does not follow regular CSS import, but they
				// get pulled to the top.
				b.AssertFileContent(filepath.Join(workDir, "public/index.html"), `T1: @import "regular.css";
@import "another.css";
moo {
  color: #fff;
}

moo {
  color: #fff;
}

/* foo */`)
			}
		})
	}
}

func TestSCSSWithThemeOverrides(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name     string
		supports func() bool
	}{
		{"libsass", func() bool { return scss.Supports() }},
		{"dartsass", func() bool { return dartsass.Supports() }},
	} {
		c.Run(test.name, func(c *qt.C) {
			if !test.supports() {
				c.Skip(fmt.Sprintf("Skip %s", test.name))
			}

			workDir, clean1, err := htesting.CreateTempDir(hugofs.Os, fmt.Sprintf("hugo-scss-include-theme-overrides-%s", test.name))
			c.Assert(err, qt.IsNil)
			defer clean1()

			theme := "mytheme"
			themesDir := filepath.Join(workDir, "themes")
			themeDirs := filepath.Join(themesDir, theme)
			v := config.New()
			v.Set("workingDir", workDir)
			v.Set("theme", theme)
			b := newTestSitesBuilder(c).WithLogger(loggers.NewErrorLogger())
			// Need to use OS fs for this.
			b.Fs = hugofs.NewDefault(v)
			b.WithWorkingDir(workDir)
			b.WithViper(v)

			fooDir := filepath.Join(workDir, "node_modules", "foo")
			scssDir := filepath.Join(workDir, "assets", "scss")
			scssThemeDir := filepath.Join(themeDirs, "assets", "scss")
			c.Assert(os.MkdirAll(fooDir, 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "content", "sect"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "data"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "i18n"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "layouts", "shortcodes"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(workDir, "layouts", "_default"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(scssDir, "components"), 0777), qt.IsNil)
			c.Assert(os.MkdirAll(filepath.Join(scssThemeDir, "components"), 0777), qt.IsNil)

			b.WithSourceFile(filepath.Join(scssThemeDir, "components", "_imports.scss"), `
@import "moo";
@import "_boo";
@import "_zoo";

`)

			b.WithSourceFile(filepath.Join(scssThemeDir, "components", "_moo.scss"), `
$moolor: #fff;

moo {
  color: $moolor;
}
`)

			// Only in theme.
			b.WithSourceFile(filepath.Join(scssThemeDir, "components", "_zoo.scss"), `
$zoolor: pink;

zoo {
  color: $zoolor;
}
`)

			b.WithSourceFile(filepath.Join(scssThemeDir, "components", "_boo.scss"), `
$boolor: orange;

boo {
  color: $boolor;
}
`)

			b.WithSourceFile(filepath.Join(scssThemeDir, "main.scss"), `
@import "components/imports";

`)

			b.WithSourceFile(filepath.Join(scssDir, "components", "_moo.scss"), `
$moolor: #ccc;

moo {
  color: $moolor;
}
`)

			b.WithSourceFile(filepath.Join(scssDir, "components", "_boo.scss"), `
$boolor: green;

boo {
  color: $boolor;
}
`)

			b.WithTemplatesAdded("index.html", fmt.Sprintf(`
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo" ) "transpiler" %q ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}
`, test.name))
			b.Build(BuildCfg{})

			b.AssertFileContent(
				filepath.Join(workDir, "public/index.html"),
				`T1: moo{color:#ccc}boo{color:green}zoo{color:pink}`,
			)
		})
	}
}

// https://github.com/gohugoio/hugo/issues/6274
func TestSCSSWithIncludePathsSass(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name     string
		supports func() bool
	}{
		{"libsass", func() bool { return scss.Supports() }},
		{"dartsass", func() bool { return dartsass.Supports() }},
	} {
		c.Run(test.name, func(c *qt.C) {
			if !test.supports() {
				c.Skip(fmt.Sprintf("Skip %s", test.name))
			}
		})
	}
	if !scss.Supports() {
		t.Skip("Skip SCSS")
	}
	workDir, clean1, err := htesting.CreateTempDir(hugofs.Os, "hugo-scss-includepaths")
	c.Assert(err, qt.IsNil)
	defer clean1()

	v := config.New()
	v.Set("workingDir", workDir)
	v.Set("theme", "mytheme")
	b := newTestSitesBuilder(t).WithLogger(loggers.NewErrorLogger())
	// Need to use OS fs for this.
	b.Fs = hugofs.NewDefault(v)
	b.WithWorkingDir(workDir)
	b.WithViper(v)

	hulmaDir := filepath.Join(workDir, "node_modules", "hulma")
	scssDir := filepath.Join(workDir, "themes/mytheme/assets", "scss")
	c.Assert(os.MkdirAll(hulmaDir, 0777), qt.IsNil)
	c.Assert(os.MkdirAll(scssDir, 0777), qt.IsNil)

	b.WithSourceFile(filepath.Join(scssDir, "main.scss"), `
@import "hulma/hulma";

`)

	b.WithSourceFile(filepath.Join(hulmaDir, "hulma.sass"), `
$hulma: #ccc;

foo
  color: $hulma;

`)

	b.WithTemplatesAdded("index.html", `
 {{ $scssOptions := (dict "targetPath" "css/styles.css" "enableSourceMap" false "includePaths" (slice "node_modules")) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $scssOptions  | minify  }}
T1: {{ $r.Content }}
`)
	b.Build(BuildCfg{})

	b.AssertFileContent(filepath.Join(workDir, "public/index.html"), `T1: foo{color:#ccc}`)
}

func TestResourceChainBasic(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
	t.Cleanup(func() {
		ts.Close()
	})

	b := newTestSitesBuilder(t)
	b.WithTemplatesAdded("index.html", fmt.Sprintf(`
{{ $hello := "<h1>     Hello World!   </h1>" | resources.FromString "hello.html" | fingerprint "sha512" | minify  | fingerprint }}
{{ $cssFingerprinted1 := "body {  background-color: lightblue; }" | resources.FromString "styles.css" |  minify  | fingerprint }}
{{ $cssFingerprinted2 := "body {  background-color: orange; }" | resources.FromString "styles2.css" |  minify  | fingerprint }}


HELLO: {{ $hello.Name }}|{{ $hello.RelPermalink }}|{{ $hello.Content | safeHTML }}

{{ $img := resources.Get "images/sunset.jpg" }}
{{ $fit := $img.Fit "200x200" }}
{{ $fit2 := $fit.Fit "100x200" }}
{{ $img = $img | fingerprint }}
SUNSET: {{ $img.Name }}|{{ $img.RelPermalink }}|{{ $img.Width }}|{{ len $img.Content }}
FIT: {{ $fit.Name }}|{{ $fit.RelPermalink }}|{{ $fit.Width }}
CSS integrity Data first: {{ $cssFingerprinted1.Data.Integrity }} {{ $cssFingerprinted1.RelPermalink }}
CSS integrity Data last:  {{ $cssFingerprinted2.RelPermalink }} {{ $cssFingerprinted2.Data.Integrity }}

{{ $rimg := resources.GetRemote "%[1]s/sunset.jpg" }}
{{ $remotenotfound := resources.GetRemote "%[1]s/notfound.jpg" }}
{{ $localnotfound := resources.Get "images/notfound.jpg" }}
{{ $gopherprotocol := resources.GetRemote "gopher://example.org" }}
{{ $rfit := $rimg.Fit "200x200" }}
{{ $rfit2 := $rfit.Fit "100x200" }}
{{ $rimg = $rimg | fingerprint }}
SUNSET REMOTE: {{ $rimg.Name }}|{{ $rimg.RelPermalink }}|{{ $rimg.Width }}|{{ len $rimg.Content }}
FIT REMOTE: {{ $rfit.Name }}|{{ $rfit.RelPermalink }}|{{ $rfit.Width }}
REMOTE NOT FOUND: {{ if $remotenotfound }}FAILED{{ else}}OK{{ end }}
LOCAL NOT FOUND: {{ if $localnotfound }}FAILED{{ else}}OK{{ end }}
PRINT PROTOCOL ERROR1: {{ with $gopherprotocol }}{{ . | safeHTML }}{{ end }}
PRINT PROTOCOL ERROR2: {{ with $gopherprotocol }}{{ .Err | safeHTML }}{{ end }}

`, ts.URL))

	fs := b.Fs.Source

	imageDir := filepath.Join("assets", "images")
	b.Assert(os.MkdirAll(imageDir, 0777), qt.IsNil)
	src, err := os.Open("testdata/sunset.jpg")
	b.Assert(err, qt.IsNil)
	out, err := fs.Create(filepath.Join(imageDir, "sunset.jpg"))
	b.Assert(err, qt.IsNil)
	_, err = io.Copy(out, src)
	b.Assert(err, qt.IsNil)
	out.Close()

	b.Running()

	for i := 0; i < 2; i++ {

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html",
			fmt.Sprintf(`
SUNSET: images/sunset.jpg|/images/sunset.a9bf1d944e19c0f382e0d8f51de690f7d0bc8fa97390c4242a86c3e5c0737e71.jpg|900|90587
FIT: images/sunset.jpg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x200_fit_q75_box.jpg|200
CSS integrity Data first: sha256-od9YaHw8nMOL8mUy97Sy8sKwMV3N4hI3aVmZXATxH&#43;8= /styles.min.a1df58687c3c9cc38bf26532f7b4b2f2c2b0315dcde212376959995c04f11fef.css
CSS integrity Data last:  /styles2.min.1cfc52986836405d37f9998a63fd6dd8608e8c410e5e3db1daaa30f78bc273ba.css sha256-HPxSmGg2QF03&#43;ZmKY/1t2GCOjEEOXj2x2qow94vCc7o=

SUNSET REMOTE: sunset_%[1]s.jpg|/sunset_%[1]s.a9bf1d944e19c0f382e0d8f51de690f7d0bc8fa97390c4242a86c3e5c0737e71.jpg|900|90587
FIT REMOTE: sunset_%[1]s.jpg|/sunset_%[1]s_hu59e56ffff1bc1d8d122b1403d34e039f_0_200x200_fit_q75_box.jpg|200
REMOTE NOT FOUND: OK
LOCAL NOT FOUND: OK
PRINT PROTOCOL ERROR1: error calling resources.GetRemote: Get "gopher://example.org": unsupported protocol scheme "gopher"
PRINT PROTOCOL ERROR2: error calling resources.GetRemote: Get "gopher://example.org": unsupported protocol scheme "gopher"


`, helpers.HashString(ts.URL+"/sunset.jpg", map[string]interface{}{})))

		b.AssertFileContent("public/styles.min.a1df58687c3c9cc38bf26532f7b4b2f2c2b0315dcde212376959995c04f11fef.css", "body{background-color:#add8e6}")
		b.AssertFileContent("public//styles2.min.1cfc52986836405d37f9998a63fd6dd8608e8c410e5e3db1daaa30f78bc273ba.css", "body{background-color:orange}")

		b.EditFiles("page1.md", `
---
title: "Page 1 edit"
summary: "Edited summary"
---

Edited content.

`)

		b.Assert(b.Fs.Destination.Remove("public"), qt.IsNil)
		b.H.ResourceSpec.ClearCaches()

	}
}

func TestResourceChainPostProcess(t *testing.T) {
	t.Parallel()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `[minify]
  minifyOutput = true
  [minify.tdewolff]
    [minify.tdewolff.html]
      keepQuotes = false
      keepWhitespace = false`)
	b.WithContent("page1.md", "---\ntitle: Page1\n---")
	b.WithContent("page2.md", "---\ntitle: Page2\n---")

	b.WithTemplates(
		"_default/single.html", `{{ $hello := "<h1>     Hello World!   </h1>" | resources.FromString "hello.html" | minify  | fingerprint "md5" | resources.PostProcess }}
HELLO: {{ $hello.RelPermalink }}	
`,
		"index.html", `Start.
{{ $hello := "<h1>     Hello World!   </h1>" | resources.FromString "hello.html" | minify  | fingerprint "md5" | resources.PostProcess }}

HELLO: {{ $hello.RelPermalink }}|Integrity: {{ $hello.Data.Integrity }}|MediaType: {{ $hello.MediaType.Type }}
HELLO2: Name: {{ $hello.Name }}|Content: {{ $hello.Content }}|Title: {{ $hello.Title }}|ResourceType: {{ $hello.ResourceType }}

// Issue #8884
<a href="hugo.rocks">foo</a>
<a href="{{ $hello.RelPermalink }}" integrity="{{ $hello.Data.Integrity}}">Hello</a>
`+strings.Repeat("a b", rnd.Intn(10)+1)+`


End.`)

	b.Running()
	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html",
		`Start.
HELLO: /hello.min.a2d1cb24f24b322a7dad520414c523e9.html|Integrity: md5-otHLJPJLMip9rVIEFMUj6Q==|MediaType: text/html
HELLO2: Name: hello.html|Content: <h1>Hello World!</h1>|Title: hello.html|ResourceType: text
<a href=hugo.rocks>foo</a>
<a href="/hello.min.a2d1cb24f24b322a7dad520414c523e9.html" integrity="md5-otHLJPJLMip9rVIEFMUj6Q==">Hello</a>
End.`)

	b.AssertFileContent("public/page1/index.html", `HELLO: /hello.min.a2d1cb24f24b322a7dad520414c523e9.html`)
	b.AssertFileContent("public/page2/index.html", `HELLO: /hello.min.a2d1cb24f24b322a7dad520414c523e9.html`)
}

func BenchmarkResourceChainPostProcess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := newTestSitesBuilder(b)
		for i := 0; i < 300; i++ {
			s.WithContent(fmt.Sprintf("page%d.md", i+1), "---\ntitle: Page\n---")
		}
		s.WithTemplates("_default/single.html", `Start.
Some text.


{{ $hello1 := "<h1>     Hello World 2!   </h1>" | resources.FromString "hello.html" | minify  | fingerprint "md5" | resources.PostProcess }}
{{ $hello2 := "<h1>     Hello World 2!   </h1>" | resources.FromString (printf "%s.html" .Path) | minify  | fingerprint "md5" | resources.PostProcess }}

Some more text.

HELLO: {{ $hello1.RelPermalink }}|Integrity: {{ $hello1.Data.Integrity }}|MediaType: {{ $hello1.MediaType.Type }}

Some more text.

HELLO2: Name: {{ $hello2.Name }}|Content: {{ $hello2.Content }}|Title: {{ $hello2.Title }}|ResourceType: {{ $hello2.ResourceType }}

Some more text.

HELLO2_2: Name: {{ $hello2.Name }}|Content: {{ $hello2.Content }}|Title: {{ $hello2.Title }}|ResourceType: {{ $hello2.ResourceType }}

End.
`)

		b.StartTimer()
		s.Build(BuildCfg{})

	}
}

func TestResourceChains(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/css/styles1.css":
			w.Header().Set("Content-Type", "text/css")
			w.Write([]byte(`h1 { 
				font-style: bold;
			}`))
			return

		case "/js/script1.js":
			w.Write([]byte(`var x; x = 5, document.getElementById("demo").innerHTML = x * 10`))
			return

		case "/mydata/json1.json":
			w.Write([]byte(`{
				"employees": [
					{
						"firstName": "John",
						"lastName": "Doe"
					},
					{
						"firstName": "Anna",
						"lastName": "Smith"
					},
					{
						"firstName": "Peter",
						"lastName": "Jones"
					}
				]
			}`))
			return

		case "/mydata/xml1.xml":
			w.Write([]byte(`
					<hello>
						<world>Hugo Rocks!</<world>
					</hello>`))
			return

		case "/mydata/svg1.svg":
			w.Header().Set("Content-Disposition", `attachment; filename="image.svg"`)
			w.Write([]byte(`
				<svg height="100" width="100">
					<path d="M1e2 1e2H3e2 2e2z"/>
				</svg>`))
			return

		case "/mydata/html1.html":
			w.Write([]byte(`
				<html>
					<a href=#>Cool</a>
				</html>`))
			return

		case "/authenticated/":
			w.Header().Set("Content-Type", "text/plain")
			if r.Header.Get("Authorization") == "Bearer abcd" {
				w.Write([]byte(`Welcome`))
				return
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
			return

		case "/post":
			w.Header().Set("Content-Type", "text/plain")
			if r.Method == http.MethodPost {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				w.Write(body)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}))
	t.Cleanup(func() {
		ts.Close()
	})

	tests := []struct {
		name      string
		shouldRun func() bool
		prepare   func(b *sitesBuilder)
		verify    func(b *sitesBuilder)
	}{
		{"tocss", func() bool { return scss.Supports() }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $scss := resources.Get "scss/styles2.scss" | toCSS }}
{{ $sass := resources.Get "sass/styles3.sass" | toCSS }}
{{ $scssCustomTarget := resources.Get "scss/styles2.scss" | toCSS (dict "targetPath" "styles/main.css") }}
{{ $scssCustomTargetString := resources.Get "scss/styles2.scss" | toCSS "styles/main.css" }}
{{ $scssMin := resources.Get "scss/styles2.scss" | toCSS | minify  }}
{{  $scssFromTempl :=  ".{{ .Kind }} { color: blue; }" | resources.FromString "kindofblue.templ"  | resources.ExecuteAsTemplate "kindofblue.scss" . | toCSS (dict "targetPath" "styles/templ.css") | minify }}
{{ $bundle1 := slice $scssFromTempl $scssMin  | resources.Concat "styles/bundle1.css" }}
T1: Len Content: {{ len $scss.Content }}|RelPermalink: {{ $scss.RelPermalink }}|Permalink: {{ $scss.Permalink }}|MediaType: {{ $scss.MediaType.Type }}
T2: Content: {{ $scssMin.Content }}|RelPermalink: {{ $scssMin.RelPermalink }}
T3: Content: {{ len $scssCustomTarget.Content }}|RelPermalink: {{ $scssCustomTarget.RelPermalink }}|MediaType: {{ $scssCustomTarget.MediaType.Type }}
T4: Content: {{ len $scssCustomTargetString.Content }}|RelPermalink: {{ $scssCustomTargetString.RelPermalink }}|MediaType: {{ $scssCustomTargetString.MediaType.Type }}
T5: Content: {{ $sass.Content }}|T5 RelPermalink: {{ $sass.RelPermalink }}|
T6: {{ $bundle1.Permalink }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `T1: Len Content: 24|RelPermalink: /scss/styles2.css|Permalink: http://example.com/scss/styles2.css|MediaType: text/css`)
			b.AssertFileContent("public/index.html", `T2: Content: body{color:#333}|RelPermalink: /scss/styles2.min.css`)
			b.AssertFileContent("public/index.html", `T3: Content: 24|RelPermalink: /styles/main.css|MediaType: text/css`)
			b.AssertFileContent("public/index.html", `T4: Content: 24|RelPermalink: /styles/main.css|MediaType: text/css`)
			b.AssertFileContent("public/index.html", `T5: Content: .content-navigation {`)
			b.AssertFileContent("public/index.html", `T5 RelPermalink: /sass/styles3.css|`)
			b.AssertFileContent("public/index.html", `T6: http://example.com/styles/bundle1.css`)

			c.Assert(b.CheckExists("public/styles/templ.min.css"), qt.Equals, false)
			b.AssertFileContent("public/styles/bundle1.css", `.home{color:blue}body{color:#333}`)
		}},

		{"minify", func() bool { return true }, func(b *sitesBuilder) {
			b.WithConfigFile("toml", `[minify]
  [minify.tdewolff]
    [minify.tdewolff.html]
      keepWhitespace = false
`)
			b.WithTemplates("home.html", fmt.Sprintf(`
Min CSS: {{ ( resources.Get "css/styles1.css" | minify ).Content }}
Min CSS Remote: {{ ( resources.GetRemote "%[1]s/css/styles1.css" | minify ).Content }}
Min JS: {{ ( resources.Get "js/script1.js" | resources.Minify ).Content | safeJS }}
Min JS Remote: {{ ( resources.GetRemote "%[1]s/js/script1.js" | minify ).Content }}
Min JSON: {{ ( resources.Get "mydata/json1.json" | resources.Minify ).Content | safeHTML }}
Min JSON Remote: {{ ( resources.GetRemote "%[1]s/mydata/json1.json" | resources.Minify ).Content | safeHTML }}
Min XML: {{ ( resources.Get "mydata/xml1.xml" | resources.Minify ).Content | safeHTML }}
Min XML Remote: {{ ( resources.GetRemote "%[1]s/mydata/xml1.xml" | resources.Minify ).Content | safeHTML }}
Min SVG: {{ ( resources.Get "mydata/svg1.svg" | resources.Minify ).Content | safeHTML }}
Min SVG Remote: {{ ( resources.GetRemote "%[1]s/mydata/svg1.svg" | resources.Minify ).Content | safeHTML }}
Min SVG again: {{ ( resources.Get "mydata/svg1.svg" | resources.Minify ).Content | safeHTML }}
Min HTML: {{ ( resources.Get "mydata/html1.html" | resources.Minify ).Content | safeHTML }}
Min HTML Remote: {{ ( resources.GetRemote "%[1]s/mydata/html1.html" | resources.Minify ).Content | safeHTML }}
`, ts.URL))
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `Min CSS: h1{font-style:bold}`)
			b.AssertFileContent("public/index.html", `Min CSS Remote: h1{font-style:bold}`)
			b.AssertFileContent("public/index.html", `Min JS: var x;x=5,document.getElementById(&#34;demo&#34;).innerHTML=x*10`)
			b.AssertFileContent("public/index.html", `Min JS Remote: var x;x=5,document.getElementById(&#34;demo&#34;).innerHTML=x*10`)
			b.AssertFileContent("public/index.html", `Min JSON: {"employees":[{"firstName":"John","lastName":"Doe"},{"firstName":"Anna","lastName":"Smith"},{"firstName":"Peter","lastName":"Jones"}]}`)
			b.AssertFileContent("public/index.html", `Min JSON Remote: {"employees":[{"firstName":"John","lastName":"Doe"},{"firstName":"Anna","lastName":"Smith"},{"firstName":"Peter","lastName":"Jones"}]}`)
			b.AssertFileContent("public/index.html", `Min XML: <hello><world>Hugo Rocks!</<world></hello>`)
			b.AssertFileContent("public/index.html", `Min XML Remote: <hello><world>Hugo Rocks!</<world></hello>`)
			b.AssertFileContent("public/index.html", `Min SVG: <svg height="100" width="100"><path d="M1e2 1e2H3e2 2e2z"/></svg>`)
			b.AssertFileContent("public/index.html", `Min SVG Remote: <svg height="100" width="100"><path d="M1e2 1e2H3e2 2e2z"/></svg>`)
			b.AssertFileContent("public/index.html", `Min SVG again: <svg height="100" width="100"><path d="M1e2 1e2H3e2 2e2z"/></svg>`)
			b.AssertFileContent("public/index.html", `Min HTML: <html><a href=#>Cool</a></html>`)
			b.AssertFileContent("public/index.html", `Min HTML Remote: <html><a href=#>Cool</a></html>`)
		}},

		{"remote", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", fmt.Sprintf(`
{{$js := resources.GetRemote "%[1]s/js/script1.js" }}
Remote Filename: {{ $js.RelPermalink }}
{{$svg := resources.GetRemote "%[1]s/mydata/svg1.svg" }}
Remote Content-Disposition: {{ $svg.RelPermalink }}
{{$auth := resources.GetRemote "%[1]s/authenticated/" (dict "headers" (dict "Authorization" "Bearer abcd")) }}
Remote Authorization: {{ $auth.Content }}
{{$post := resources.GetRemote "%[1]s/post" (dict "method" "post" "body" "Request body") }}
Remote POST: {{ $post.Content }}
`, ts.URL))
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `Remote Filename: /script1_`)
			b.AssertFileContent("public/index.html", `Remote Content-Disposition: /image_`)
			b.AssertFileContent("public/index.html", `Remote Authorization: Welcome`)
			b.AssertFileContent("public/index.html", `Remote POST: Request body`)
		}},

		{"concat", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $a := "A" | resources.FromString "a.txt"}}
{{ $b := "B" | resources.FromString "b.txt"}}
{{ $c := "C" | resources.FromString "c.txt"}}
{{ $textResources := .Resources.Match "*.txt" }}
{{ $combined := slice $a $b $c | resources.Concat "bundle/concat.txt" }}
T1: Content: {{ $combined.Content }}|RelPermalink: {{ $combined.RelPermalink }}|Permalink: {{ $combined.Permalink }}|MediaType: {{ $combined.MediaType.Type }}
{{ with $textResources }}
{{ $combinedText := . | resources.Concat "bundle/concattxt.txt" }}
T2: Content: {{ $combinedText.Content }}|{{ $combinedText.RelPermalink }}
{{ end }}
{{/* https://github.com/gohugoio/hugo/issues/5269 */}}
{{ $css := "body { color: blue; }" | resources.FromString "styles.css" }}
{{ $minified := resources.Get "css/styles1.css" | minify }}
{{ slice $css $minified | resources.Concat "bundle/mixed.css" }} 
{{/* https://github.com/gohugoio/hugo/issues/5403 */}}
{{ $d := "function D {} // A comment" | resources.FromString "d.js"}}
{{ $e := "(function E {})" | resources.FromString "e.js"}}
{{ $f := "(function F {})()" | resources.FromString "f.js"}}
{{ $jsResources := .Resources.Match "*.js" }}
{{ $combinedJs := slice $d $e $f | resources.Concat "bundle/concatjs.js" }}
T3: Content: {{ $combinedJs.Content }}|{{ $combinedJs.RelPermalink }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `T1: Content: ABC|RelPermalink: /bundle/concat.txt|Permalink: http://example.com/bundle/concat.txt|MediaType: text/plain`)
			b.AssertFileContent("public/bundle/concat.txt", "ABC")

			b.AssertFileContent("public/index.html", `T2: Content: t1t|t2t|`)
			b.AssertFileContent("public/bundle/concattxt.txt", "t1t|t2t|")

			b.AssertFileContent("public/index.html", `T3: Content: function D {} // A comment
;
(function E {})
;
(function F {})()|`)
			b.AssertFileContent("public/bundle/concatjs.js", `function D {} // A comment
;
(function E {})
;
(function F {})()`)
		}},

		{"concat and fingerprint", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $a := "A" | resources.FromString "a.txt"}}
{{ $b := "B" | resources.FromString "b.txt"}}
{{ $c := "C" | resources.FromString "c.txt"}}
{{ $combined := slice $a $b $c | resources.Concat "bundle/concat.txt" }}
{{ $fingerprinted := $combined | fingerprint }}
Fingerprinted: {{ $fingerprinted.RelPermalink }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", "Fingerprinted: /bundle/concat.b5d4045c3f466fa91fe2cc6abe79232a1a57cdf104f7a26e716e0a1e2789df78.txt")
			b.AssertFileContent("public/bundle/concat.b5d4045c3f466fa91fe2cc6abe79232a1a57cdf104f7a26e716e0a1e2789df78.txt", "ABC")
		}},

		{"fromstring", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $r := "Hugo Rocks!" | resources.FromString "rocks/hugo.txt" }}
{{ $r.Content }}|{{ $r.RelPermalink }}|{{ $r.Permalink }}|{{ $r.MediaType.Type }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `Hugo Rocks!|/rocks/hugo.txt|http://example.com/rocks/hugo.txt|text/plain`)
			b.AssertFileContent("public/rocks/hugo.txt", "Hugo Rocks!")
		}},
		{"execute-as-template", func() bool {
			return true
		}, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $var := "Hugo Page" }}
{{ if .IsHome }}
{{ $var = "Hugo Home" }}
{{ end }}
T1: {{ $var }}
{{ $result := "{{ .Kind | upper }}" | resources.FromString "mytpl.txt" | resources.ExecuteAsTemplate "result.txt" . }}
T2: {{ $result.Content }}|{{ $result.RelPermalink}}|{{$result.MediaType.Type }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `T2: HOME|/result.txt|text/plain`, `T1: Hugo Home`)
		}},
		{"fingerprint", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $r := "ab" | resources.FromString "rocks/hugo.txt" }}
{{ $result := $r | fingerprint }}
{{ $result512 := $r | fingerprint "sha512" }}
{{ $resultMD5 := $r | fingerprint "md5" }}
T1: {{ $result.Content }}|{{ $result.RelPermalink}}|{{$result.MediaType.Type }}|{{ $result.Data.Integrity }}|
T2: {{ $result512.Content }}|{{ $result512.RelPermalink}}|{{$result512.MediaType.Type }}|{{ $result512.Data.Integrity }}|
T3: {{ $resultMD5.Content }}|{{ $resultMD5.RelPermalink}}|{{$resultMD5.MediaType.Type }}|{{ $resultMD5.Data.Integrity }}|
{{ $r2 := "bc" | resources.FromString "rocks/hugo2.txt" | fingerprint }}
{{/* https://github.com/gohugoio/hugo/issues/5296 */}}
T4: {{ $r2.Data.Integrity }}|


`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `T1: ab|/rocks/hugo.fb8e20fc2e4c3f248c60c39bd652f3c1347298bb977b8b4d5903b85055620603.txt|text/plain|sha256-&#43;44g/C5MPySMYMOb1lLzwTRymLuXe4tNWQO4UFViBgM=|`)
			b.AssertFileContent("public/index.html", `T2: ab|/rocks/hugo.2d408a0717ec188158278a796c689044361dc6fdde28d6f04973b80896e1823975cdbf12eb63f9e0591328ee235d80e9b5bf1aa6a44f4617ff3caf6400eb172d.txt|text/plain|sha512-LUCKBxfsGIFYJ4p5bGiQRDYdxv3eKNbwSXO4CJbhgjl1zb8S62P54FkTKO4jXYDptb8apqRPRhf/PK9kAOsXLQ==|`)
			b.AssertFileContent("public/index.html", `T3: ab|/rocks/hugo.187ef4436122d1cc2f40dc2b92f0eba0.txt|text/plain|md5-GH70Q2Ei0cwvQNwrkvDroA==|`)
			b.AssertFileContent("public/index.html", `T4: sha256-Hgu9bGhroFC46wP/7txk/cnYCUf86CGrvl1tyNJSxaw=|`)
		}},
		// https://github.com/gohugoio/hugo/issues/5226
		{"baseurl-path", func() bool { return true }, func(b *sitesBuilder) {
			b.WithSimpleConfigFileAndBaseURL("https://example.com/hugo/")
			b.WithTemplates("home.html", `
{{ $r1 := "ab" | resources.FromString "rocks/hugo.txt" }}
T1: {{ $r1.Permalink }}|{{ $r1.RelPermalink }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `T1: https://example.com/hugo/rocks/hugo.txt|/hugo/rocks/hugo.txt`)
		}},

		// https://github.com/gohugoio/hugo/issues/4944
		{"Prevent resource publish on .Content only", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $cssInline := "body { color: green; }" | resources.FromString "inline.css" | minify }}
{{ $cssPublish1 := "body { color: blue; }" | resources.FromString "external1.css" | minify }}
{{ $cssPublish2 := "body { color: orange; }" | resources.FromString "external2.css" | minify }}

Inline: {{ $cssInline.Content }}
Publish 1: {{ $cssPublish1.Content }} {{ $cssPublish1.RelPermalink }}
Publish 2: {{ $cssPublish2.Permalink }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html",
				`Inline: body{color:green}`,
				"Publish 1: body{color:blue} /external1.min.css",
				"Publish 2: http://example.com/external2.min.css",
			)
			b.Assert(b.CheckExists("public/external2.css"), qt.Equals, false)
			b.Assert(b.CheckExists("public/external1.css"), qt.Equals, false)
			b.Assert(b.CheckExists("public/external2.min.css"), qt.Equals, true)
			b.Assert(b.CheckExists("public/external1.min.css"), qt.Equals, true)
			b.Assert(b.CheckExists("public/inline.min.css"), qt.Equals, false)
		}},

		{"unmarshal", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $toml := "slogan = \"Hugo Rocks!\"" | resources.FromString "slogan.toml" | transform.Unmarshal }}
{{ $csv1 := "\"Hugo Rocks\",\"Hugo is Fast!\"" | resources.FromString "slogans.csv" | transform.Unmarshal }}
{{ $csv2 := "a;b;c" | transform.Unmarshal (dict "delimiter" ";") }}
{{ $xml := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><note><to>You</to><from>Me</from><heading>Reminder</heading><body>Do not forget XML</body></note>" | transform.Unmarshal }}

Slogan: {{ $toml.slogan }}
CSV1: {{ $csv1 }} {{ len (index $csv1 0)  }}
CSV2: {{ $csv2 }}		
XML: {{ $xml.body }}
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html",
				`Slogan: Hugo Rocks!`,
				`[[Hugo Rocks Hugo is Fast!]] 2`,
				`CSV2: [[a b c]]`,
				`XML: Do not forget XML`,
			)
		}},
		{"resources.Get", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `NOT FOUND: {{ if (resources.Get "this-does-not-exist") }}FAILED{{ else }}OK{{ end }}`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", "NOT FOUND: OK")
		}},

		{"template", func() bool { return true }, func(b *sitesBuilder) {}, func(b *sitesBuilder) {
		}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if !test.shouldRun() {
				t.Skip()
			}
			t.Parallel()

			b := newTestSitesBuilder(t).WithLogger(loggers.NewErrorLogger())
			b.WithContent("_index.md", `
---
title: Home
---

Home.

`,
				"page1.md", `
---
title: Hello1
---

Hello1
`,
				"page2.md", `
---
title: Hello2
---

Hello2
`,
				"t1.txt", "t1t|",
				"t2.txt", "t2t|",
			)

			b.WithSourceFile(filepath.Join("assets", "css", "styles1.css"), `
h1 {
	 font-style: bold;
}
`)

			b.WithSourceFile(filepath.Join("assets", "js", "script1.js"), `
var x;
x = 5;
document.getElementById("demo").innerHTML = x * 10;
`)

			b.WithSourceFile(filepath.Join("assets", "mydata", "json1.json"), `
{
"employees":[
    {"firstName":"John", "lastName":"Doe"}, 
    {"firstName":"Anna", "lastName":"Smith"},
    {"firstName":"Peter", "lastName":"Jones"}
]
}
`)

			b.WithSourceFile(filepath.Join("assets", "mydata", "svg1.svg"), `
<svg height="100" width="100">
  <path d="M 100 100 L 300 100 L 200 100 z"/>
</svg> 
`)

			b.WithSourceFile(filepath.Join("assets", "mydata", "xml1.xml"), `
<hello>
<world>Hugo Rocks!</<world>
</hello>
`)

			b.WithSourceFile(filepath.Join("assets", "mydata", "html1.html"), `
<html>
<a  href="#">
Cool
</a >
</html>
`)

			b.WithSourceFile(filepath.Join("assets", "scss", "styles2.scss"), `
$color: #333;

body {
  color: $color;
}
`)

			b.WithSourceFile(filepath.Join("assets", "sass", "styles3.sass"), `
$color: #333;

.content-navigation
  border-color: $color

`)

			test.prepare(b)
			b.Build(BuildCfg{})
			test.verify(b)
		})
	}
}

func TestMultiSiteResource(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	b := newMultiSiteTestDefaultBuilder(t)

	b.CreateSites().Build(BuildCfg{})

	// This build is multilingual, but not multihost. There should be only one pipes.txt
	b.AssertFileContent("public/fr/index.html", "French Home Page", "String Resource: /blog/text/pipes.txt")
	c.Assert(b.CheckExists("public/fr/text/pipes.txt"), qt.Equals, false)
	c.Assert(b.CheckExists("public/en/text/pipes.txt"), qt.Equals, false)
	b.AssertFileContent("public/en/index.html", "Default Home Page", "String Resource: /blog/text/pipes.txt")
	b.AssertFileContent("public/text/pipes.txt", "Hugo Pipes")
}

func TestResourcesMatch(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)

	b.WithContent("page.md", "")

	b.WithSourceFile(
		"assets/jsons/data1.json", "json1 content",
		"assets/jsons/data2.json", "json2 content",
		"assets/jsons/data3.xml", "xml content",
	)

	b.WithTemplates("index.html", `
{{ $jsons := (resources.Match "jsons/*.json") }}
{{ $json := (resources.GetMatch "jsons/*.json") }}
{{ printf "JSONS: %d"  (len $jsons) }}
JSON: {{ $json.RelPermalink }}: {{ $json.Content }}
{{ range $jsons }}
{{- .RelPermalink }}: {{ .Content }}
{{ end }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		"JSON: /jsons/data1.json: json1 content",
		"JSONS: 2", "/jsons/data1.json: json1 content")
}

func TestExecuteAsTemplateWithLanguage(t *testing.T) {
	b := newMultiSiteTestDefaultBuilder(t)
	indexContent := `
Lang: {{ site.Language.Lang }}
{{ $templ := "{{T \"hello\"}}" | resources.FromString "f1.html" }}
{{ $helloResource := $templ | resources.ExecuteAsTemplate (print "f%s.html" .Lang) . }}
Hello1: {{T "hello"}}
Hello2: {{ $helloResource.Content }}
LangURL: {{ relLangURL "foo" }}
`
	b.WithTemplatesAdded("index.html", indexContent)
	b.WithTemplatesAdded("index.fr.html", indexContent)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/en/index.html", `
Hello1: Hello
Hello2: Hello
`)

	b.AssertFileContent("public/fr/index.html", `
Hello1: Bonjour
Hello2: Bonjour
`)
}

func TestResourceChainPostCSS(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)

	packageJSON := `{
  "scripts": {},

  "devDependencies": {
    "postcss-cli": "7.1.0",
    "tailwindcss": "1.2.0"
  }
}
`

	postcssConfig := `
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

	tailwindCss := `
@tailwind base;
@tailwind components;
@tailwind utilities;

@import "components/all.css";

h1 {
    @apply text-2xl font-bold;
}
  
`

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-test-postcss")
	c.Assert(err, qt.IsNil)
	defer clean()

	var logBuf bytes.Buffer

	newTestBuilder := func(v config.Provider) *sitesBuilder {
		v.Set("workingDir", workDir)
		v.Set("disableKinds", []string{"taxonomy", "term", "page"})
		logger := loggers.NewBasicLoggerForWriter(jww.LevelInfo, &logBuf)
		b := newTestSitesBuilder(t).WithLogger(logger)
		// Need to use OS fs for this.
		b.Fs = hugofs.NewDefault(v)
		b.WithWorkingDir(workDir)
		b.WithViper(v)

		b.WithContent("p1.md", "")
		b.WithTemplates("index.html", `
{{ $options := dict "inlineImports" true }}
{{ $styles := resources.Get "css/styles.css" | resources.PostCSS $options }}
Styles RelPermalink: {{ $styles.RelPermalink }}
{{ $cssContent := $styles.Content }}
Styles Content: Len: {{ len $styles.Content }}|

`)

		return b
	}

	b := newTestBuilder(config.New())

	cssDir := filepath.Join(workDir, "assets", "css", "components")
	b.Assert(os.MkdirAll(cssDir, 0777), qt.IsNil)

	b.WithSourceFile("assets/css/styles.css", tailwindCss)
	b.WithSourceFile("assets/css/components/all.css", `
@import "a.css";
@import "b.css";
`, "assets/css/components/a.css", `
class-in-a {
	color: blue;
}
`, "assets/css/components/b.css", `
@import "a.css";

class-in-b {
	color: blue;
}
`)

	b.WithSourceFile("package.json", packageJSON)
	b.WithSourceFile("postcss.config.js", postcssConfig)

	b.Assert(os.Chdir(workDir), qt.IsNil)
	cmd := b.NpmInstall()
	err = cmd.Run()
	b.Assert(err, qt.IsNil)
	b.Build(BuildCfg{})

	// Make sure Node sees this.
	b.Assert(logBuf.String(), qt.Contains, "Hugo Environment: production")
	b.Assert(logBuf.String(), qt.Contains, filepath.FromSlash(fmt.Sprintf("PostCSS Config File: %s/postcss.config.js", workDir)))
	b.Assert(logBuf.String(), qt.Contains, filepath.FromSlash(fmt.Sprintf("package.json: %s/package.json", workDir)))

	b.AssertFileContent("public/index.html", `
Styles RelPermalink: /css/styles.css
Styles Content: Len: 770878|
`)

	assertCss := func(b *sitesBuilder) {
		content := b.FileContent("public/css/styles.css")

		b.Assert(strings.Contains(content, "class-in-a"), qt.Equals, true)
		b.Assert(strings.Contains(content, "class-in-b"), qt.Equals, true)
	}

	assertCss(b)

	build := func(s string, shouldFail bool) error {
		b.Assert(os.RemoveAll(filepath.Join(workDir, "public")), qt.IsNil)

		v := config.New()
		v.Set("build", map[string]interface{}{
			"useResourceCacheWhen": s,
		})

		b = newTestBuilder(v)

		b.Assert(os.RemoveAll(filepath.Join(workDir, "public")), qt.IsNil)

		err := b.BuildE(BuildCfg{})
		if shouldFail {
			b.Assert(err, qt.Not(qt.IsNil))
		} else {
			b.Assert(err, qt.IsNil)
			assertCss(b)
		}

		return err
	}

	build("always", false)
	build("fallback", false)

	// Introduce a syntax error in an import
	b.WithSourceFile("assets/css/components/b.css", `@import "a.css";

class-in-b {
	@apply asdf;
}
`)

	err = build("never", true)

	err = herrors.UnwrapErrorWithFileContext(err)
	_, ok := err.(*herrors.ErrorWithFileContext)
	b.Assert(ok, qt.Equals, true)

	// TODO(bep) for some reason, we have starting to get
	// execute of template failed: template: index.html:5:25
	// on CI (GitHub action).
	// b.Assert(fe.Position().LineNumber, qt.Equals, 5)
	// b.Assert(fe.Error(), qt.Contains, filepath.Join(workDir, "assets/css/components/b.css:4:1"))

	// Remove PostCSS
	b.Assert(os.RemoveAll(filepath.Join(workDir, "node_modules")), qt.IsNil)

	build("always", false)
	build("fallback", false)
	build("never", true)

	// Remove cache
	b.Assert(os.RemoveAll(filepath.Join(workDir, "resources")), qt.IsNil)

	build("always", true)
	build("fallback", true)
	build("never", true)
}

func TestResourceMinifyDisabled(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).WithConfigFile("toml", `
baseURL = "https://example.org"

[minify]
disableXML=true


`)

	b.WithContent("page.md", "")

	b.WithSourceFile(
		"assets/xml/data.xml", "<root>   <foo> asdfasdf </foo> </root>",
	)

	b.WithTemplates("index.html", `
{{ $xml := resources.Get "xml/data.xml" | minify | fingerprint }}
XML: {{ $xml.Content | safeHTML }}|{{ $xml.RelPermalink }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
XML: <root>   <foo> asdfasdf </foo> </root>|/xml/data.min.3be4fddd19aaebb18c48dd6645215b822df74701957d6d36e59f203f9c30fd9f.xml
`)
}

// Issue 8954
func TestMinifyWithError(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()
	b.WithSourceFile(
		"assets/js/test.js", `
new Date(2002, 04, 11)
`,
	)
	b.WithTemplates("index.html", `
{{ $js := resources.Get "js/test.js" | minify | fingerprint }}
<script>
{{ $js.Content }}
</script>
`)
	b.WithContent("page.md", "")

	err := b.BuildE(BuildCfg{})

	if err == nil || !strings.Contains(err.Error(), "04") {
		t.Fatalf("expected a message about a legacy octal number, but got: %v", err)
	}
}
