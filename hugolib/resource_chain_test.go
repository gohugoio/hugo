// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/scss"
)

func TestSCSSWithIncludePaths(t *testing.T) {
	if !scss.Supports() {
		t.Skip("Skip SCSS")
	}
	assert := require.New(t)
	workDir, clean, err := createTempDir("hugo-scss-include")
	assert.NoError(err)
	defer clean()

	v := viper.New()
	v.Set("workingDir", workDir)
	b := newTestSitesBuilder(t).WithLogger(loggers.NewWarningLogger())
	b.WithViper(v)
	b.WithWorkingDir(workDir)
	// Need to use OS fs for this.
	b.Fs = hugofs.NewDefault(v)

	fooDir := filepath.Join(workDir, "node_modules", "foo")
	scssDir := filepath.Join(workDir, "assets", "scss")
	assert.NoError(os.MkdirAll(fooDir, 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "content", "sect"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "data"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "i18n"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "layouts", "shortcodes"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "layouts", "_default"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(scssDir), 0777))

	b.WithSourceFile(filepath.Join(fooDir, "_moo.scss"), `
$moolor: #fff;

moo {
  color: $moolor;
}
`)

	b.WithSourceFile(filepath.Join(scssDir, "main.scss"), `
@import "moo";

`)

	b.WithTemplatesAdded("index.html", `
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo" ) ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}
`)
	b.Build(BuildCfg{})

	b.AssertFileContent(filepath.Join(workDir, "public/index.html"), `T1: moo{color:#fff}`)

}

func TestSCSSWithThemeOverrides(t *testing.T) {
	if !scss.Supports() {
		t.Skip("Skip SCSS")
	}
	assert := require.New(t)
	workDir, clean, err := createTempDir("hugo-scss-include")
	assert.NoError(err)
	defer clean()

	theme := "mytheme"
	themesDir := filepath.Join(workDir, "themes")
	themeDirs := filepath.Join(themesDir, theme)
	v := viper.New()
	v.Set("workingDir", workDir)
	v.Set("theme", theme)
	b := newTestSitesBuilder(t).WithLogger(loggers.NewWarningLogger())
	b.WithViper(v)
	b.WithWorkingDir(workDir)
	// Need to use OS fs for this.
	b.Fs = hugofs.NewDefault(v)

	fooDir := filepath.Join(workDir, "node_modules", "foo")
	scssDir := filepath.Join(workDir, "assets", "scss")
	scssThemeDir := filepath.Join(themeDirs, "assets", "scss")
	assert.NoError(os.MkdirAll(fooDir, 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "content", "sect"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "data"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "i18n"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "layouts", "shortcodes"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(workDir, "layouts", "_default"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(scssDir, "components"), 0777))
	assert.NoError(os.MkdirAll(filepath.Join(scssThemeDir, "components"), 0777))

	b.WithSourceFile(filepath.Join(scssThemeDir, "components", "_imports.scss"), `
@import "moo";
@import "_boo";
`)

	b.WithSourceFile(filepath.Join(scssThemeDir, "components", "_moo.scss"), `
$moolor: #fff;

moo {
  color: $moolor;
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

	b.WithTemplatesAdded("index.html", `
{{ $cssOpts := (dict "includePaths" (slice "node_modules/foo" ) ) }}
{{ $r := resources.Get "scss/main.scss" |  toCSS $cssOpts  | minify  }}
T1: {{ $r.Content }}
`)
	b.Build(BuildCfg{})

	b.AssertFileContent(filepath.Join(workDir, "public/index.html"), `T1: moo{color:#ccc}boo{color:green}`)

}

func TestResourceChain(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

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

			assert.False(b.CheckExists("public/styles/templ.min.css"))
			b.AssertFileContent("public/styles/bundle1.css", `.home{color:blue}body{color:#333}`)

		}},

		{"minify", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
Min CSS: {{ ( resources.Get "css/styles1.css" | minify ).Content }}
Min JS: {{ ( resources.Get "js/script1.js" | resources.Minify ).Content | safeJS }}
Min JSON: {{ ( resources.Get "mydata/json1.json" | resources.Minify ).Content | safeHTML }}
Min XML: {{ ( resources.Get "mydata/xml1.xml" | resources.Minify ).Content | safeHTML }}
Min SVG: {{ ( resources.Get "mydata/svg1.svg" | resources.Minify ).Content | safeHTML }}
Min SVG again: {{ ( resources.Get "mydata/svg1.svg" | resources.Minify ).Content | safeHTML }}
Min HTML: {{ ( resources.Get "mydata/html1.html" | resources.Minify ).Content | safeHTML }}


`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `Min CSS: h1{font-style:bold}`)
			b.AssertFileContent("public/index.html", `Min JS: var x;x=5;document.getElementById(&#34;demo&#34;).innerHTML=x*10;`)
			b.AssertFileContent("public/index.html", `Min JSON: {"employees":[{"firstName":"John","lastName":"Doe"},{"firstName":"Anna","lastName":"Smith"},{"firstName":"Peter","lastName":"Jones"}]}`)
			b.AssertFileContent("public/index.html", `Min XML: <hello><world>Hugo Rocks!</<world></hello>`)
			b.AssertFileContent("public/index.html", `Min SVG: <svg height="100" width="100"><path d="M5 10 20 40z"/></svg>`)
			b.AssertFileContent("public/index.html", `Min SVG again: <svg height="100" width="100"><path d="M5 10 20 40z"/></svg>`)
			b.AssertFileContent("public/index.html", `Min HTML: <html><a href=#>Cool</a></html>`)
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
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html", `T1: Content: ABC|RelPermalink: /bundle/concat.txt|Permalink: http://example.com/bundle/concat.txt|MediaType: text/plain`)
			b.AssertFileContent("public/bundle/concat.txt", "ABC")

			b.AssertFileContent("public/index.html", `T2: Content: t1t|t2t|`)
			b.AssertFileContent("public/bundle/concattxt.txt", "t1t|t2t|")
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
			// TODO(bep) eventually remove
			return isGo111()
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
			assert.True(b.CheckExists("public/external2.min.css"), "Referenced content should be copied to /public")
			assert.True(b.CheckExists("public/external1.min.css"), "Referenced content should be copied to /public")

			assert.False(b.CheckExists("public/inline.min.css"), "Inline content should not be copied to /public")
		}},

		{"unmarshal", func() bool { return true }, func(b *sitesBuilder) {
			b.WithTemplates("home.html", `
{{ $toml := "slogan = \"Hugo Rocks!\"" | resources.FromString "slogan.toml" | transform.Unmarshal }}
{{ $csv1 := "\"Hugo Rocks\",\"Hugo is Fast!\"" | resources.FromString "slogans.csv" | transform.Unmarshal }}
{{ $csv2 := "a;b;c" | transform.Unmarshal (dict "delimiter" ";") }}

Slogan: {{ $toml.slogan }}
CSV1: {{ $csv1 }} {{ len (index $csv1 0)  }}
CSV2: {{ $csv2 }}		
`)
		}, func(b *sitesBuilder) {
			b.AssertFileContent("public/index.html",
				`Slogan: Hugo Rocks!`,
				`[[Hugo Rocks Hugo is Fast!]] 2`,
				`CSV2: [[a b c]]`,
			)
		}},

		{"template", func() bool { return true }, func(b *sitesBuilder) {}, func(b *sitesBuilder) {
		}},
	}

	for _, test := range tests {
		if !test.shouldRun() {
			t.Log("Skip", test.name)
			continue
		}

		b := newTestSitesBuilder(t).WithLogger(loggers.NewWarningLogger())
		b.WithSimpleConfigFile()
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
  <line x1="5" y1="10" x2="20" y2="40"/>
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

		t.Log("Test", test.name)
		test.prepare(b)
		b.Build(BuildCfg{})
		test.verify(b)
	}
}

func TestMultiSiteResource(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	b := newMultiSiteTestDefaultBuilder(t)

	b.CreateSites().Build(BuildCfg{})

	// This build is multilingual, but not multihost. There should be only one pipes.txt
	b.AssertFileContent("public/fr/index.html", "French Home Page", "String Resource: /blog/text/pipes.txt")
	assert.False(b.CheckExists("public/fr/text/pipes.txt"))
	assert.False(b.CheckExists("public/en/text/pipes.txt"))
	b.AssertFileContent("public/en/index.html", "Default Home Page", "String Resource: /blog/text/pipes.txt")
	b.AssertFileContent("public/text/pipes.txt", "Hugo Pipes")

}
