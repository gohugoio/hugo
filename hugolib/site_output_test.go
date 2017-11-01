// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"

	"fmt"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/output"
	"github.com/spf13/viper"
)

func TestDefaultOutputFormats(t *testing.T) {
	t.Parallel()
	defs, err := createDefaultOutputFormats(output.DefaultFormats, viper.New())

	require.NoError(t, err)

	tests := []struct {
		name string
		kind string
		want output.Formats
	}{
		{"RSS not for regular pages", KindPage, output.Formats{output.HTMLFormat}},
		{"Home Sweet Home", KindHome, output.Formats{output.HTMLFormat, output.RSSFormat}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defs[tt.kind]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDefaultOutputFormats(%v) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}

func TestDefaultOutputFormatsWithOverrides(t *testing.T) {
	t.Parallel()

	htmlOut := output.HTMLFormat
	htmlOut.BaseName = "htmlindex"
	rssOut := output.RSSFormat
	rssOut.BaseName = "feed"

	defs, err := createDefaultOutputFormats(output.Formats{htmlOut, rssOut}, viper.New())

	homeDefs := defs[KindHome]

	rss, found := homeDefs.GetByName("RSS")
	require.True(t, found)
	require.Equal(t, rss.BaseName, "feed")

	html, found := homeDefs.GetByName("HTML")
	require.True(t, found)
	require.Equal(t, html.BaseName, "htmlindex")

	require.NoError(t, err)

}

func TestSiteWithPageOutputs(t *testing.T) {
	for _, outputs := range [][]string{{"html", "json", "calendar"}, {"json"}} {
		t.Run(fmt.Sprintf("%v", outputs), func(t *testing.T) {
			doTestSiteWithPageOutputs(t, outputs)
		})
	}
}

func doTestSiteWithPageOutputs(t *testing.T, outputs []string) {
	t.Parallel()

	outputsStr := strings.Replace(fmt.Sprintf("%q", outputs), " ", ", ", -1)

	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 1
defaultContentLanguage = "en"

disableKinds = ["page", "section", "taxonomy", "taxonomyTerm", "RSS", "sitemap", "robotsTXT", "404"]

[Taxonomies]
tag = "tags"
category = "categories"

defaultContentLanguage = "en"

[languages]

[languages.en]
title = "Title in English"
languageName = "English"
weight = 1

[languages.nn]
languageName = "Nynorsk"
weight = 2
title = "Tittel på Nynorsk"

`

	pageTemplate := `---
title: "%s"
outputs: %s
---
# Doc

{{< myShort >}}

{{< myOtherShort >}}

`

	mf := afero.NewMemMapFs()

	writeToFs(t, mf, "i18n/en.toml", `
[elbow]
other = "Elbow"
`)
	writeToFs(t, mf, "i18n/nn.toml", `
[elbow]
other = "Olboge"
`)

	th, h := newTestSitesFromConfig(t, mf, siteConfig,

		// Case issue partials #3333
		"layouts/partials/GoHugo.html", `Go Hugo Partial`,
		"layouts/_default/baseof.json", `START JSON:{{block "main" .}}default content{{ end }}:END JSON`,
		"layouts/_default/baseof.html", `START HTML:{{block "main" .}}default content{{ end }}:END HTML`,
		"layouts/shortcodes/myOtherShort.html", `OtherShort: {{ "<h1>Hi!</h1>" | safeHTML }}`,
		"layouts/shortcodes/myShort.html", `ShortHTML`,
		"layouts/shortcodes/myShort.json", `ShortJSON`,

		"layouts/_default/list.json", `{{ define "main" }}
List JSON|{{ .Title }}|{{ .Content }}|Alt formats: {{ len .AlternativeOutputFormats -}}|
{{- range .AlternativeOutputFormats -}}
Alt Output: {{ .Name -}}|
{{- end -}}|
{{- range .OutputFormats -}}
Output/Rel: {{ .Name -}}/{{ .Rel }}|{{ .MediaType }}
{{- end -}}
 {{ with .OutputFormats.Get "JSON" }}
<atom:link href={{ .Permalink }} rel="self" type="{{ .MediaType }}" />
{{ end }}
{{ .Site.Language.Lang }}: {{ T "elbow" -}}
{{ end }}
`,
		"layouts/_default/list.html", `{{ define "main" }}
List HTML|{{.Title }}|
{{- with .OutputFormats.Get "HTML" -}}
<atom:link href={{ .Permalink }} rel="self" type="{{ .MediaType }}" />
{{- end -}}
{{ .Site.Language.Lang }}: {{ T "elbow" -}}
Partial Hugo 1: {{ partial "GoHugo.html" . }}
Partial Hugo 2: {{ partial "GoHugo" . -}}
Content: {{ .Content }}
{{ end }}
`,
	)
	require.Len(t, h.Sites, 2)

	fs := th.Fs

	writeSource(t, fs, "content/_index.md", fmt.Sprintf(pageTemplate, "JSON Home", outputsStr))
	writeSource(t, fs, "content/_index.nn.md", fmt.Sprintf(pageTemplate, "JSON Nynorsk Heim", outputsStr))

	err := h.Build(BuildCfg{})

	require.NoError(t, err)

	s := h.Sites[0]
	require.Equal(t, "en", s.Language.Lang)

	home := s.getPage(KindHome)

	require.NotNil(t, home)

	lenOut := len(outputs)

	require.Len(t, home.outputFormats, lenOut)

	// There is currently always a JSON output to make it simpler ...
	altFormats := lenOut - 1
	hasHTML := helpers.InStringArray(outputs, "html")
	th.assertFileContent("public/index.json",
		"List JSON",
		fmt.Sprintf("Alt formats: %d", altFormats),
	)

	if hasHTML {
		th.assertFileContent("public/index.json",
			"Alt Output: HTML",
			"Output/Rel: JSON/alternate|",
			"Output/Rel: HTML/canonical|",
			"en: Elbow",
			"ShortJSON",
			"OtherShort: <h1>Hi!</h1>",
		)

		th.assertFileContent("public/index.html",
			// The HTML entity is a deliberate part of this test: The HTML templates are
			// parsed with html/template.
			`List HTML|JSON Home|<atom:link href=http://example.com/blog/ rel="self" type="text/html&#43;html" />`,
			"en: Elbow",
			"ShortHTML",
			"OtherShort: <h1>Hi!</h1>",
		)
		th.assertFileContent("public/nn/index.html",
			"List HTML|JSON Nynorsk Heim|",
			"nn: Olboge")
	} else {
		th.assertFileContent("public/index.json",
			"Output/Rel: JSON/canonical|",
			// JSON is plain text, so no need to safeHTML this and that
			`<atom:link href=http://example.com/blog/index.json rel="self" type="application/json+json" />`,
			"ShortJSON",
			"OtherShort: <h1>Hi!</h1>",
		)
		th.assertFileContent("public/nn/index.json",
			"List JSON|JSON Nynorsk Heim|",
			"nn: Olboge",
			"ShortJSON",
		)
	}

	of := home.OutputFormats()
	require.Len(t, of, lenOut)
	require.Nil(t, of.Get("Hugo"))
	require.NotNil(t, of.Get("json"))
	json := of.Get("JSON")
	_, err = home.AlternativeOutputFormats()
	require.Error(t, err)
	require.NotNil(t, json)
	require.Equal(t, "/blog/index.json", json.RelPermalink())
	require.Equal(t, "http://example.com/blog/index.json", json.Permalink())

	if helpers.InStringArray(outputs, "cal") {
		cal := of.Get("calendar")
		require.NotNil(t, cal)
		require.Equal(t, "/blog/index.ics", cal.RelPermalink())
		require.Equal(t, "webcal://example.com/blog/index.ics", cal.Permalink())
	}

	require.True(t, home.HasShortcode("myShort"))
	require.False(t, home.HasShortcode("doesNotExist"))

}

// Issue #3447
func TestRedefineRSSOutputFormat(t *testing.T) {
	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 1
defaultContentLanguage = "en"

disableKinds = ["page", "section", "taxonomy", "taxonomyTerm", "sitemap", "robotsTXT", "404"]

[outputFormats]
[outputFormats.RSS]
mediatype = "application/rss"
baseName = "feed"

`

	mf := afero.NewMemMapFs()
	writeToFs(t, mf, "content/foo.html", `foo`)

	th, h := newTestSitesFromConfig(t, mf, siteConfig)

	err := h.Build(BuildCfg{})

	require.NoError(t, err)

	th.assertFileContent("public/feed.xml", "Recent content on")

	s := h.Sites[0]

	//Issue #3450
	require.Equal(t, "http://example.com/blog/feed.xml", s.Info.RSSLink)

}

// Issue #3614
func TestDotLessOutputFormat(t *testing.T) {
	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 1
defaultContentLanguage = "en"

disableKinds = ["page", "section", "taxonomy", "taxonomyTerm", "sitemap", "robotsTXT", "404"]

[mediaTypes]
[mediaTypes."text/nodot"]
suffix = ""
delimiter = ""
[mediaTypes."text/defaultdelim"]
suffix = "defd"
[mediaTypes."text/nosuffix"]
suffix = ""
[mediaTypes."text/customdelim"]
suffix = "del"
delimiter = "_"

[outputs]
home = [ "DOTLESS", "DEF", "NOS", "CUS" ]

[outputFormats]
[outputFormats.DOTLESS]
mediatype = "text/nodot"
baseName = "_redirects" # This is how Netlify names their redirect files.
[outputFormats.DEF]
mediatype = "text/defaultdelim"
baseName = "defaultdelimbase"
[outputFormats.NOS]
mediatype = "text/nosuffix"
baseName = "nosuffixbase"
[outputFormats.CUS]
mediatype = "text/customdelim"
baseName = "customdelimbase"

`

	mf := afero.NewMemMapFs()
	writeToFs(t, mf, "content/foo.html", `foo`)
	writeToFs(t, mf, "layouts/_default/list.dotless", `a dotless`)
	writeToFs(t, mf, "layouts/_default/list.def.defd", `default delimim`)
	writeToFs(t, mf, "layouts/_default/list.nos", `no suffix`)
	writeToFs(t, mf, "layouts/_default/list.cus.del", `custom delim`)

	th, h := newTestSitesFromConfig(t, mf, siteConfig)

	err := h.Build(BuildCfg{})

	require.NoError(t, err)

	th.assertFileContent("public/_redirects", "a dotless")
	th.assertFileContent("public/defaultdelimbase.defd", "default delimim")
	// This looks weird, but the user has chosen this definition.
	th.assertFileContent("public/nosuffixbase.", "no suffix")
	th.assertFileContent("public/customdelimbase_del", "custom delim")

	s := h.Sites[0]
	home := s.getPage(KindHome)
	require.NotNil(t, home)

	outputs := home.OutputFormats()

	require.Equal(t, "/blog/_redirects", outputs.Get("DOTLESS").RelPermalink())
	require.Equal(t, "/blog/defaultdelimbase.defd", outputs.Get("DEF").RelPermalink())
	require.Equal(t, "/blog/nosuffixbase.", outputs.Get("NOS").RelPermalink())
	require.Equal(t, "/blog/customdelimbase_del", outputs.Get("CUS").RelPermalink())

}
