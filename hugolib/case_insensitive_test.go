// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/hugofs"
)

var (
	caseMixingSiteConfigTOML = `
Title = "In an Insensitive Mood"
DefaultContentLanguage = "nn"
defaultContentLanguageInSubdir = true

[Blackfriday]
AngledQuotes = true
HrefTargetBlank = true

[Params]
Search = true
Color = "green"
mood = "Happy"
[Params.Colors]
Blue = "blue"
Yellow = "yellow"

[Languages]
[Languages.nn]
title = "Nynorsk title"
languageName = "Nynorsk"
weight = 1

[Languages.en]
TITLE = "English title"
LanguageName = "English"
Mood = "Thoughtful"
Weight = 2
COLOR = "Pink"
[Languages.en.blackfriday]
angledQuotes = false
hrefTargetBlank = false
[Languages.en.Colors]
BLUE = "blues"
yellow = "golden"
`
	caseMixingPage1En = `
---
TITLE: Page1 En Translation
BlackFriday:
  AngledQuotes: false
Color: "black"
Search: true
mooD: "sad and lonely"
ColorS: 
  Blue: "bluesy"
  Yellow: "sunny"  
---
# "Hi"
{{< shortcode >}}
`

	caseMixingPage1 = `
---
titLe: Side 1
blackFriday:
  angledQuotes: true
color: "red"
search: false
MooD: "sad"
COLORS: 
  blue: "heavenly"
  yelloW: "Sunny"  
---
# "Hi"
{{< shortcode >}}
`

	caseMixingPage2 = `
---
TITLE: Page2 Title
BlackFriday:
  AngledQuotes: false
Color: "black"
search: true
MooD: "moody"
ColorS: 
  Blue: "sky"
  YELLOW: "flower"  
---
# Hi
{{< shortcode >}}
`
)

func caseMixingTestsWriteCommonSources(t *testing.T, fs *hugofs.Fs) {
	writeSource(t, fs, filepath.Join("content", "sect1", "page1.md"), caseMixingPage1)
	writeSource(t, fs, filepath.Join("content", "sect2", "page2.md"), caseMixingPage2)
	writeSource(t, fs, filepath.Join("content", "sect1", "page1.en.md"), caseMixingPage1En)

	writeSource(t, fs, "layouts/shortcodes/shortcode.html", `
Shortcode Page: {{ .Page.Params.COLOR }}|{{ .Page.Params.Colors.Blue  }}
Shortcode Site: {{ .Page.Site.Params.COLOR }}|{{ .Site.Params.COLORS.YELLOW  }}
`)

	writeSource(t, fs, "layouts/partials/partial.html", `
Partial Page: {{ .Params.COLOR }}|{{ .Params.Colors.Blue }}
Partial Site: {{ .Site.Params.COLOR }}|{{ .Site.Params.COLORS.YELLOW }}
`)

	writeSource(t, fs, "config.toml", caseMixingSiteConfigTOML)

}

func TestCaseInsensitiveConfigurationVariations(t *testing.T) {
	// See issues 2615, 1129, 2590 and maybe some others
	// Also see 2598
	//
	// Viper is now, at least for the Hugo part, case insensitive
	// So we need tests for all of it, with needed adjustments on the Hugo side.
	// Not sure what that will be. Let us see.

	// So all the below with case variations:
	// config: regular fields, blackfriday config, param with nested map
	// language: new and overridden values, in regular fields and nested paramsmap
	// page frontmatter: regular fields, blackfriday config, param with nested map

	testCommonResetState()

	depsCfg := newTestDepsConfig()
	viper.SetFs(depsCfg.Fs.Source)

	caseMixingTestsWriteCommonSources(t, depsCfg.Fs)

	if err := LoadGlobalConfig("", "config.toml"); err != nil {
		t.Fatalf("Failed to load config: %s", err)
	}

	writeSource(t, depsCfg.Fs, filepath.Join("layouts", "_default", "baseof.html"), `
Block Page Colors: {{ .Params.COLOR }}|{{ .Params.Colors.Blue }}	
{{ block "main" . }}default{{end}}`)

	writeSource(t, depsCfg.Fs, filepath.Join("layouts", "sect2", "single.html"), `
{{ define "main"}}
Page Colors: {{ .Params.CoLOR }}|{{ .Params.Colors.Blue }}
Site Colors: {{ .Site.Params.COlOR }}|{{ .Site.Params.COLORS.YELLOW }}
{{ .Content }}
{{ partial "partial.html" . }}
{{ end }}
`)

	writeSource(t, depsCfg.Fs, filepath.Join("layouts", "_default", "single.html"), `
Page Title: {{ .Title }}
Site Title: {{ .Site.Title }}
Site Lang Mood: {{ .Site.Language.Params.MOoD }}
Page Colors: {{ .Params.COLOR }}|{{ .Params.Colors.Blue }}
Site Colors: {{ .Site.Params.COLOR }}|{{ .Site.Params.COLORS.YELLOW }}
{{ .Content }}
{{ partial "partial.html" . }}
`)

	sites, err := NewHugoSitesFromConfiguration(depsCfg)

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	err = sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	assertFileContent(t, sites.Fs, filepath.Join("public", "nn", "sect1", "page1", "index.html"), true,
		"Page Colors: red|heavenly",
		"Site Colors: green|yellow",
		"Site Lang Mood: Happy",
		"Shortcode Page: red|heavenly",
		"Shortcode Site: green|yellow",
		"Partial Page: red|heavenly",
		"Partial Site: green|yellow",
		"Page Title: Side 1",
		"Site Title: Nynorsk title",
		"&laquo;Hi&raquo;", // angled quotes
	)

	assertFileContent(t, sites.Fs, filepath.Join("public", "en", "sect1", "page1", "index.html"), true,
		"Site Colors: Pink|golden",
		"Page Colors: black|bluesy",
		"Site Lang Mood: Thoughtful",
		"Page Title: Page1 En Translation",
		"Site Title: English title",
		"&ldquo;Hi&rdquo;",
	)

	assertFileContent(t, sites.Fs, filepath.Join("public", "nn", "sect2", "page2", "index.html"), true,
		"Page Colors: black|sky",
		"Site Colors: green|yellow",
		"Shortcode Page: black|sky",
		"Block Page Colors: black|sky",
		"Partial Page: black|sky",
		"Partial Site: green|yellow",
	)
}

func TestCaseInsensitiveConfigurationForAllTemplateEngines(t *testing.T) {
	noOp := func(s string) string {
		return s
	}

	amberFixer := func(s string) string {
		fixed := strings.Replace(s, "{{ .Site.Params", "{{ Site.Params", -1)
		fixed = strings.Replace(fixed, "{{ .Params", "{{ Params", -1)
		fixed = strings.Replace(fixed, ".Content", "Content", -1)
		fixed = strings.Replace(fixed, "{{", "#{", -1)
		fixed = strings.Replace(fixed, "}}", "}", -1)

		return fixed
	}

	for _, config := range []struct {
		suffix        string
		templateFixer func(s string) string
	}{
		{"amber", amberFixer},
		{"html", noOp},
		{"ace", noOp},
	} {
		doTestCaseInsensitiveConfigurationForTemplateEngine(t, config.suffix, config.templateFixer)

	}

}

func doTestCaseInsensitiveConfigurationForTemplateEngine(t *testing.T, suffix string, templateFixer func(s string) string) {

	testCommonResetState()

	fs := hugofs.NewMem()
	viper.SetFs(fs.Source)

	caseMixingTestsWriteCommonSources(t, fs)

	if err := LoadGlobalConfig("", "config.toml"); err != nil {
		t.Fatalf("Failed to load config: %s", err)
	}

	t.Log("Testing", suffix)

	templTemplate := `
p
	|
	| Page Colors: {{ .Params.CoLOR }}|{{ .Params.Colors.Blue }}
	| Site Colors: {{ .Site.Params.COlOR }}|{{ .Site.Params.COLORS.YELLOW }}
	| {{ .Content }}

`

	templ := templateFixer(templTemplate)

	t.Log(templ)

	writeSource(t, fs, filepath.Join("layouts", "_default", fmt.Sprintf("single.%s", suffix)), templ)

	sites, err := NewHugoSitesFromConfiguration(deps.DepsCfg{Fs: fs})

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	err = sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	assertFileContent(t, sites.Fs, filepath.Join("public", "nn", "sect1", "page1", "index.html"), true,
		"Page Colors: red|heavenly",
		"Site Colors: green|yellow",
		"Shortcode Page: red|heavenly",
		"Shortcode Site: green|yellow",
	)

}
