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
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/resources/kinds"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/output"
)

func TestSiteWithPageOutputs(t *testing.T) {
	for _, outputs := range [][]string{{"html", "json", "calendar"}, {"json"}} {
		outputs := outputs
		t.Run(fmt.Sprintf("%v", outputs), func(t *testing.T) {
			t.Parallel()
			doTestSiteWithPageOutputs(t, outputs)
		})
	}
}

func doTestSiteWithPageOutputs(t *testing.T, outputs []string) {
	outputsStr := strings.Replace(fmt.Sprintf("%q", outputs), " ", ", ", -1)

	siteConfig := `
baseURL = "http://example.com/blog"

defaultContentLanguage = "en"

disableKinds = ["section", "term", "taxonomy", "RSS", "sitemap", "robotsTXT", "404"]

[pagination]
pagerSize = 1

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

	b := newTestSitesBuilder(t).WithConfigFile("toml", siteConfig)
	b.WithI18n("en.toml", `
[elbow]
other = "Elbow"
`, "nn.toml", `
[elbow]
other = "Olboge"
`)

	b.WithTemplates(
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
Len Pages: {{ .Kind }} {{ len .Site.RegularPages }} Page Number: {{ .Paginator.PageNumber }}
{{ end }}
`,
		"layouts/_default/single.html", `{{ define "main" }}{{ .Content }}{{ end }}`,
	)

	b.WithContent("_index.md", fmt.Sprintf(pageTemplate, "JSON Home", outputsStr))
	b.WithContent("_index.nn.md", fmt.Sprintf(pageTemplate, "JSON Nynorsk Heim", outputsStr))

	for i := 1; i <= 10; i++ {
		b.WithContent(fmt.Sprintf("p%d.md", i), fmt.Sprintf(pageTemplate, fmt.Sprintf("Page %d", i), outputsStr))
	}

	b.Build(BuildCfg{})

	s := b.H.Sites[0]
	b.Assert(s.language.Lang, qt.Equals, "en")

	home := s.getPageOldVersion(kinds.KindHome)

	b.Assert(home, qt.Not(qt.IsNil))

	lenOut := len(outputs)

	b.Assert(len(home.OutputFormats()), qt.Equals, lenOut)

	// There is currently always a JSON output to make it simpler ...
	altFormats := lenOut - 1
	hasHTML := hstrings.InSlice(outputs, "html")
	b.AssertFileContent("public/index.json",
		"List JSON",
		fmt.Sprintf("Alt formats: %d", altFormats),
	)

	if hasHTML {
		b.AssertFileContent("public/index.json",
			"Alt Output: html",
			"Output/Rel: json/alternate|",
			"Output/Rel: html/canonical|",
			"en: Elbow",
			"ShortJSON",
			"OtherShort: <h1>Hi!</h1>",
		)

		b.AssertFileContent("public/index.html",
			// The HTML entity is a deliberate part of this test: The HTML templates are
			// parsed with html/template.
			`List HTML|JSON Home|<atom:link href=http://example.com/blog/ rel="self" type="text/html" />`,
			"en: Elbow",
			"ShortHTML",
			"OtherShort: <h1>Hi!</h1>",
			"Len Pages: home 10",
		)
		b.AssertFileContent("public/page/2/index.html", "Page Number: 2")
		b.Assert(b.CheckExists("public/page/2/index.json"), qt.Equals, false)

		b.AssertFileContent("public/nn/index.html",
			"List HTML|JSON Nynorsk Heim|",
			"nn: Olboge")
	} else {
		b.AssertFileContent("public/index.json",
			"Output/Rel: json/canonical|",
			// JSON is plain text, so no need to safeHTML this and that
			`<atom:link href=http://example.com/blog/index.json rel="self" type="application/json" />`,
			"ShortJSON",
			"OtherShort: <h1>Hi!</h1>",
		)
		b.AssertFileContent("public/nn/index.json",
			"List JSON|JSON Nynorsk Heim|",
			"nn: Olboge",
			"ShortJSON",
		)
	}

	of := home.OutputFormats()

	json := of.Get("JSON")
	b.Assert(json, qt.Not(qt.IsNil))
	b.Assert(json.RelPermalink(), qt.Equals, "/blog/index.json")
	b.Assert(json.Permalink(), qt.Equals, "http://example.com/blog/index.json")

	if hstrings.InSlice(outputs, "cal") {
		cal := of.Get("calendar")
		b.Assert(cal, qt.Not(qt.IsNil))
		b.Assert(cal.RelPermalink(), qt.Equals, "/blog/index.ics")
		b.Assert(cal.Permalink(), qt.Equals, "webcal://example.com/blog/index.ics")
	}

	b.Assert(home.HasShortcode("myShort"), qt.Equals, true)
	b.Assert(home.HasShortcode("doesNotExist"), qt.Equals, false)
}

// Issue #3447
func TestRedefineRSSOutputFormat(t *testing.T) {
	siteConfig := `
baseURL = "http://example.com/blog"

defaultContentLanguage = "en"

disableKinds = ["page", "section", "term", "taxonomy", "sitemap", "robotsTXT", "404"]

[pagination]
pagerSize = 1

[outputFormats]
[outputFormats.RSS]
mediatype = "application/rss"
baseName = "feed"

`

	c := qt.New(t)

	mf := afero.NewMemMapFs()
	writeToFs(t, mf, "content/foo.html", `foo`)

	th, h := newTestSitesFromConfig(t, mf, siteConfig)

	err := h.Build(BuildCfg{})

	c.Assert(err, qt.IsNil)

	th.assertFileContent("public/feed.xml", "Recent content on")

	s := h.Sites[0]

	// Issue #3450
	c.Assert(s.Home().OutputFormats().Get("rss").Permalink(), qt.Equals, "http://example.com/blog/feed.xml")
}

// Issue #3614
func TestDotLessOutputFormat(t *testing.T) {
	siteConfig := `
baseURL = "http://example.com/blog"

defaultContentLanguage = "en"

disableKinds = ["page", "section", "term", "taxonomy", "sitemap", "robotsTXT", "404"]

[pagination]
pagerSize = 1

[mediaTypes]
[mediaTypes."text/nodot"]
delimiter = ""
[mediaTypes."text/defaultdelim"]
suffixes = ["defd"]
[mediaTypes."text/nosuffix"]
[mediaTypes."text/customdelim"]
suffixes = ["del"]
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

	c := qt.New(t)

	mf := afero.NewMemMapFs()
	writeToFs(t, mf, "content/foo.html", `foo`)
	writeToFs(t, mf, "layouts/_default/list.dotless", `a dotless`)
	writeToFs(t, mf, "layouts/_default/list.def.defd", `default delimim`)
	writeToFs(t, mf, "layouts/_default/list.nos", `no suffix`)
	writeToFs(t, mf, "layouts/_default/list.cus.del", `custom delim`)

	th, h := newTestSitesFromConfig(t, mf, siteConfig)

	err := h.Build(BuildCfg{})

	c.Assert(err, qt.IsNil)

	s := h.Sites[0]

	th.assertFileContent("public/_redirects", "a dotless")
	th.assertFileContent("public/defaultdelimbase.defd", "default delimim")
	// This looks weird, but the user has chosen this definition.
	th.assertFileContent("public/nosuffixbase", "no suffix")
	th.assertFileContent("public/customdelimbase_del", "custom delim")

	home := s.getPageOldVersion(kinds.KindHome)
	c.Assert(home, qt.Not(qt.IsNil))

	outputs := home.OutputFormats()

	c.Assert(outputs.Get("DOTLESS").RelPermalink(), qt.Equals, "/blog/_redirects")
	c.Assert(outputs.Get("DEF").RelPermalink(), qt.Equals, "/blog/defaultdelimbase.defd")
	c.Assert(outputs.Get("NOS").RelPermalink(), qt.Equals, "/blog/nosuffixbase")
	c.Assert(outputs.Get("CUS").RelPermalink(), qt.Equals, "/blog/customdelimbase_del")
}

// Issue 8030
func TestGetOutputFormatRel(t *testing.T) {
	b := newTestSitesBuilder(t).
		WithSimpleConfigFileAndSettings(map[string]any{
			"outputFormats": map[string]any{
				"HUMANS": map[string]any{
					"mediaType":   "text/plain",
					"baseName":    "humans",
					"isPlainText": true,
					"rel":         "author",
				},
			},
		}).WithTemplates("index.html", `
{{- with ($.Site.GetPage "humans").OutputFormats.Get "humans" -}}
<link rel="{{ .Rel }}" type="{{ .MediaType.String }}" href="{{ .Permalink }}">
{{- end -}}
`).WithContent("humans.md", `---
outputs:
- HUMANS
---
This is my content.
`)

	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html", `
<link rel="author" type="text/plain" href="/humans.txt">
`)
}

func TestCreateSiteOutputFormats(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		c := qt.New(t)

		outputsConfig := map[string]any{
			kinds.KindHome:    []string{"HTML", "JSON"},
			kinds.KindSection: []string{"JSON"},
		}

		cfg := config.New()
		cfg.Set("outputs", outputsConfig)

		outputs, err := createSiteOutputFormats(output.DefaultFormats, cfg.GetStringMap("outputs"), false)
		c.Assert(err, qt.IsNil)
		c.Assert(outputs[kinds.KindSection], deepEqualsOutputFormats, output.Formats{output.JSONFormat})
		c.Assert(outputs[kinds.KindHome], deepEqualsOutputFormats, output.Formats{output.HTMLFormat, output.JSONFormat})

		// Defaults
		c.Assert(outputs[kinds.KindTerm], deepEqualsOutputFormats, output.Formats{output.HTMLFormat, output.RSSFormat})
		c.Assert(outputs[kinds.KindTaxonomy], deepEqualsOutputFormats, output.Formats{output.HTMLFormat, output.RSSFormat})
		c.Assert(outputs[kinds.KindPage], deepEqualsOutputFormats, output.Formats{output.HTMLFormat})

		// These aren't (currently) in use when rendering in Hugo,
		// but the pages needs to be assigned an output format,
		// so these should also be correct/sensible.
		c.Assert(outputs[kinds.KindRSS], deepEqualsOutputFormats, output.Formats{output.RSSFormat})
		c.Assert(outputs[kinds.KindSitemap], deepEqualsOutputFormats, output.Formats{output.SitemapFormat})
		c.Assert(outputs[kinds.KindRobotsTXT], deepEqualsOutputFormats, output.Formats{output.RobotsTxtFormat})
		c.Assert(outputs[kinds.KindStatus404], deepEqualsOutputFormats, output.Formats{output.HTMLFormat})
	})

	// Issue #4528
	t.Run("Mixed case", func(t *testing.T) {
		c := qt.New(t)
		cfg := config.New()

		outputsConfig := map[string]any{
			// Note that we in Hugo 0.53.0 renamed this Kind to "taxonomy",
			// but keep this test to test the legacy mapping.
			"taxonomyterm": []string{"JSON"},
		}
		cfg.Set("outputs", outputsConfig)

		outputs, err := createSiteOutputFormats(output.DefaultFormats, cfg.GetStringMap("outputs"), false)
		c.Assert(err, qt.IsNil)
		c.Assert(outputs[kinds.KindTaxonomy], deepEqualsOutputFormats, output.Formats{output.JSONFormat})
	})
}

func TestCreateSiteOutputFormatsInvalidConfig(t *testing.T) {
	c := qt.New(t)

	outputsConfig := map[string]any{
		kinds.KindHome: []string{"FOO", "JSON"},
	}

	cfg := config.New()
	cfg.Set("outputs", outputsConfig)

	_, err := createSiteOutputFormats(output.DefaultFormats, cfg.GetStringMap("outputs"), false)
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestCreateSiteOutputFormatsEmptyConfig(t *testing.T) {
	c := qt.New(t)

	outputsConfig := map[string]any{
		kinds.KindHome: []string{},
	}

	cfg := config.New()
	cfg.Set("outputs", outputsConfig)

	outputs, err := createSiteOutputFormats(output.DefaultFormats, cfg.GetStringMap("outputs"), false)
	c.Assert(err, qt.IsNil)
	c.Assert(outputs[kinds.KindHome], deepEqualsOutputFormats, output.Formats{output.HTMLFormat, output.RSSFormat})
}

func TestCreateSiteOutputFormatsCustomFormats(t *testing.T) {
	c := qt.New(t)

	outputsConfig := map[string]any{
		kinds.KindHome: []string{},
	}

	cfg := config.New()
	cfg.Set("outputs", outputsConfig)

	var (
		customRSS  = output.Format{Name: "RSS", BaseName: "customRSS"}
		customHTML = output.Format{Name: "HTML", BaseName: "customHTML"}
	)

	outputs, err := createSiteOutputFormats(output.Formats{customRSS, customHTML}, cfg.GetStringMap("outputs"), false)
	c.Assert(err, qt.IsNil)
	c.Assert(outputs[kinds.KindHome], deepEqualsOutputFormats, output.Formats{customHTML, customRSS})
}

// https://github.com/gohugoio/hugo/issues/5849
func TestOutputFormatPermalinkable(t *testing.T) {
	config := `
baseURL = "https://example.com"



# DAMP is similar to AMP, but not permalinkable.
[outputFormats]
[outputFormats.damp]
mediaType = "text/html"
path = "damp"
[outputFormats.ramp]
mediaType = "text/html"
path = "ramp"
permalinkable = true
[outputFormats.base]
mediaType = "text/html"
isHTML = true
baseName = "that"
permalinkable = true
[outputFormats.nobase]
mediaType = "application/json"
permalinkable = true

`

	b := newTestSitesBuilder(t).WithConfigFile("toml", config)
	b.WithContent("_index.md", `
---
Title: Home Sweet Home
outputs: [ "html", "amp", "damp", "base" ]
---

`)

	b.WithContent("blog/html-amp.md", `
---
Title: AMP and HTML
outputs: [ "html", "amp" ]
---

`)

	b.WithContent("blog/html-damp.md", `
---
Title: DAMP and HTML
outputs: [ "html", "damp" ]
---

`)

	b.WithContent("blog/html-ramp.md", `
---
Title: RAMP and HTML
outputs: [ "html", "ramp" ]
---

`)

	b.WithContent("blog/html.md", `
---
Title: HTML only
outputs: [ "html" ]
---

`)

	b.WithContent("blog/amp.md", `
---
Title: AMP only
outputs: [ "amp" ]
---

`)

	b.WithContent("blog/html-base-nobase.md", `
---
Title: HTML, Base and Nobase
outputs: [ "html", "base", "nobase" ]
---

`)

	const commonTemplate = `
This RelPermalink: {{ .RelPermalink }}
Output Formats: {{ len .OutputFormats }};{{ range .OutputFormats }}{{ .Name }};{{ .RelPermalink }}|{{ end }}

`

	b.WithTemplatesAdded("index.html", commonTemplate)
	b.WithTemplatesAdded("_default/single.html", commonTemplate)
	b.WithTemplatesAdded("_default/single.json", commonTemplate)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		"This RelPermalink: /",
		"Output Formats: 4;html;/|amp;/amp/|damp;/damp/|base;/that.html|",
	)

	b.AssertFileContent("public/amp/index.html",
		"This RelPermalink: /amp/",
		"Output Formats: 4;html;/|amp;/amp/|damp;/damp/|base;/that.html|",
	)

	b.AssertFileContent("public/blog/html-amp/index.html",
		"Output Formats: 2;html;/blog/html-amp/|amp;/amp/blog/html-amp/|",
		"This RelPermalink: /blog/html-amp/")

	b.AssertFileContent("public/amp/blog/html-amp/index.html",
		"Output Formats: 2;html;/blog/html-amp/|amp;/amp/blog/html-amp/|",
		"This RelPermalink: /amp/blog/html-amp/")

	// Damp is not permalinkable
	b.AssertFileContent("public/damp/blog/html-damp/index.html",
		"This RelPermalink: /blog/html-damp/",
		"Output Formats: 2;html;/blog/html-damp/|damp;/damp/blog/html-damp/|")

	b.AssertFileContent("public/blog/html-ramp/index.html",
		"This RelPermalink: /blog/html-ramp/",
		"Output Formats: 2;html;/blog/html-ramp/|ramp;/ramp/blog/html-ramp/|")

	b.AssertFileContent("public/ramp/blog/html-ramp/index.html",
		"This RelPermalink: /ramp/blog/html-ramp/",
		"Output Formats: 2;html;/blog/html-ramp/|ramp;/ramp/blog/html-ramp/|")

	// https://github.com/gohugoio/hugo/issues/5877
	outputFormats := "Output Formats: 3;html;/blog/html-base-nobase/|base;/blog/html-base-nobase/that.html|nobase;/blog/html-base-nobase/index.json|"

	b.AssertFileContent("public/blog/html-base-nobase/index.json",
		"This RelPermalink: /blog/html-base-nobase/index.json",
		outputFormats,
	)

	b.AssertFileContent("public/blog/html-base-nobase/that.html",
		"This RelPermalink: /blog/html-base-nobase/that.html",
		outputFormats,
	)

	b.AssertFileContent("public/blog/html-base-nobase/index.html",
		"This RelPermalink: /blog/html-base-nobase/",
		outputFormats,
	)
}

func TestSiteWithPageNoOutputs(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `
baseURL = "https://example.com"

[outputFormats.o1]
mediaType = "text/html"



`)
	b.WithContent("outputs-empty.md", `---
title: "Empty Outputs"
outputs: []
---

Word1. Word2.

`,
		"outputs-string.md", `---
title: "Outputs String"
outputs: "o1"
---

Word1. Word2.

`)

	b.WithTemplates("index.html", `
{{ range .Site.RegularPages }}
WordCount: {{ .WordCount }}
{{ end }}
`)

	b.WithTemplates("_default/single.html", `HTML: {{ .Content }}`)
	b.WithTemplates("_default/single.o1.html", `O1: {{ .Content }}`)

	b.Build(BuildCfg{})

	b.AssertFileContent(
		"public/index.html",
		" WordCount: 2")

	b.AssertFileContent("public/outputs-empty/index.html", "HTML:", "Word1. Word2.")
	b.AssertFileContent("public/outputs-string/index.html", "O1:", "Word1. Word2.")
}

func TestOuputFormatFrontMatterTermIssue12275(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','page','rss','section','sitemap','taxonomy']
-- content/p1.md --
---
title: p1
tags:
  - tag-a
  - tag-b
---
-- content/tags/tag-a/_index.md --
---
title: tag-a
outputs:
  - html
  - json
---
-- content/tags/tag-b/_index.md --
---
title: tag-b
---
-- layouts/_default/term.html --
{{ .Title }}
-- layouts/_default/term.json --
{{ jsonify (dict "title" .Title) }}
`

	b := Test(t, files)

	b.AssertFileContent("public/tags/tag-a/index.html", "tag-a")
	b.AssertFileContent("public/tags/tag-b/index.html", "tag-b")
	b.AssertFileContent("public/tags/tag-a/index.json", `{"title":"tag-a"}`) // failing test
}
