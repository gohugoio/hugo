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
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugofs"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"
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
Yellow = "golden"
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

func caseMixingTestsWriteCommonSources(t *testing.T, fs afero.Fs) {
	writeToFs(t, fs, filepath.Join("content", "sect1", "page1.md"), caseMixingPage1)
	writeToFs(t, fs, filepath.Join("content", "sect2", "page2.md"), caseMixingPage2)
	writeToFs(t, fs, filepath.Join("content", "sect1", "page1.en.md"), caseMixingPage1En)

	writeToFs(t, fs, "layouts/shortcodes/shortcode.html", `
Shortcode Page: {{ .Page.Params.COLOR }}|{{ .Page.Params.Colors.Blue  }}
Shortcode Site: {{ .Page.Site.Params.COLOR }}|{{ .Site.Params.COLORS.YELLOW  }}
`)

	writeToFs(t, fs, "layouts/partials/partial.html", `
Partial Page: {{ .Params.COLOR }}|{{ .Params.Colors.Blue }}
Partial Site: {{ .Site.Params.COLOR }}|{{ .Site.Params.COLORS.YELLOW }}
Partial Site Global: {{ site.Params.COLOR }}|{{ site.Params.COLORS.YELLOW }}
`)

	writeToFs(t, fs, "config.toml", caseMixingSiteConfigTOML)

}

func TestCaseInsensitiveConfigurationVariations(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	mm := afero.NewMemMapFs()

	caseMixingTestsWriteCommonSources(t, mm)

	cfg, _, err := LoadConfig(ConfigSourceDescriptor{Fs: mm, Filename: "config.toml"})
	c.Assert(err, qt.IsNil)

	fs := hugofs.NewFrom(mm, cfg)

	th := newTestHelper(cfg, fs, t)

	writeSource(t, fs, filepath.Join("layouts", "_default", "baseof.html"), `
Block Page Colors: {{ .Params.COLOR }}|{{ .Params.Colors.Blue }}	
{{ block "main" . }}default{{end}}`)

	writeSource(t, fs, filepath.Join("layouts", "sect2", "single.html"), `
{{ define "main"}}
Page Colors: {{ .Params.CoLOR }}|{{ .Params.Colors.Blue }}
Site Colors: {{ .Site.Params.COlOR }}|{{ .Site.Params.COLORS.YELLOW }}
{{ template "index-color" (dict "name" "Page" "params" .Params) }}
{{ template "index-color" (dict "name" "Site" "params" .Site.Params) }}

{{ .Content }}
{{ partial "partial.html" . }}
{{ end }}
{{ define "index-color" }}
{{ $yellow := index .params "COLoRS" "yELLOW" }}
{{ $colors := index .params "COLoRS" }}
{{ $yellow2 := index $colors "yEllow" }}
index1|{{ .name }}: {{ $yellow }}|
index2|{{ .name }}: {{ $yellow2 }}|
{{ end }}
`)

	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `
Page Title: {{ .Title }}
Site Title: {{ .Site.Title }}
Site Lang Mood: {{ .Site.Language.Params.MOoD }}
Page Colors: {{ .Params.COLOR }}|{{ .Params.Colors.Blue }}|{{ index .Params "ColOR" }}
Site Colors: {{ .Site.Params.COLOR }}|{{ .Site.Params.COLORS.YELLOW }}|{{ index .Site.Params "ColOR" }}
{{ $page2 := .Site.GetPage "/sect2/page2" }}
{{ if $page2 }}
Page2: {{ $page2.Params.ColoR }} 
{{ end }}
{{ .Content }}
{{ partial "partial.html" . }}
`)

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	err = sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	th.assertFileContent(filepath.Join("public", "nn", "sect1", "page1", "index.html"),
		"Page Colors: red|heavenly|red",
		"Site Colors: green|yellow|green",
		"Site Lang Mood: Happy",
		"Shortcode Page: red|heavenly",
		"Shortcode Site: green|yellow",
		"Partial Page: red|heavenly",
		"Partial Site: green|yellow",
		"Partial Site Global: green|yellow",
		"Page Title: Side 1",
		"Site Title: Nynorsk title",
		"Page2: black ",
	)

	th.assertFileContent(filepath.Join("public", "en", "sect1", "page1", "index.html"),
		"Site Colors: Pink|golden",
		"Page Colors: black|bluesy",
		"Site Lang Mood: Thoughtful",
		"Page Title: Page1 En Translation",
		"Site Title: English title",
		"&ldquo;Hi&rdquo;",
	)

	th.assertFileContent(filepath.Join("public", "nn", "sect2", "page2", "index.html"),
		"Page Colors: black|sky",
		"Site Colors: green|yellow",
		"Shortcode Page: black|sky",
		"Block Page Colors: black|sky",
		"Partial Page: black|sky",
		"Partial Site: green|yellow",
		"index1|Page: flower|",
		"index1|Site: yellow|",
		"index2|Page: flower|",
		"index2|Site: yellow|",
	)
}
