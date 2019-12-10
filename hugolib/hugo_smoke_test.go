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
)

// The most basic build test.
func TestHello(t *testing.T) {
	t.Parallel()
	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `
baseURL="https://example.org"
disableKinds = ["taxonomy", "taxonomyTerm", "section", "page"]
`)
	b.WithContent("p1", `
---
title: Page
---

`)
	b.WithTemplates("index.html", `Site: {{ .Site.Language.Lang | upper }}`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `Site: EN`)
}

func TestSmoke(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	const configFile = `
baseURL = "https://example.com"
title = "Simple Site"
rssLimit = 3
defaultContentLanguage = "en"
enableRobotsTXT = true

[languages]
[languages.en]
weight = 1
title = "In English"
[languages.no]
weight = 2
title = "På norsk"

[params]
hugo = "Rules!"

[outputs]
  home = ["HTML", "JSON", "CSV", "RSS"]

`

	const pageContentAndSummaryDivider = `---
title: Page with outputs
hugo: "Rocks!"
outputs: ["HTML", "JSON"]
tags: [ "hugo" ]
aliases: [ "/a/b/c" ]
---

This is summary.

<!--more--> 

This is content with some shortcodes.

Shortcode 1: {{< sc >}}.
Shortcode 2: {{< sc >}}.

`

	const pageContentWithMarkdownShortcodes = `---
title: Page with markdown shortcode
hugo: "Rocks!"
outputs: ["HTML", "JSON"]
---

This is summary.

<!--more--> 

This is content[^a].

# Header above

{{% markdown-shortcode %}}
# Header inside

Some **markdown**.[^b]

{{% /markdown-shortcode %}}

# Heder below

Some more content[^c].

Footnotes:

[^a]: Fn 1
[^b]: Fn 2
[^c]: Fn 3

`

	var pageContentAutoSummary = strings.Replace(pageContentAndSummaryDivider, "<!--more-->", "", 1)

	b := newTestSitesBuilder(t).WithConfigFile("toml", configFile)
	b.WithTemplatesAdded("shortcodes/markdown-shortcode.html", `
Some **Markdown** in shortcode.

{{ .Inner }}


	
`)

	b.WithTemplatesAdded("shortcodes/markdown-shortcode.json", `
Some **Markdown** in JSON shortcode.
{{ .Inner }}

`)

	for i := 1; i <= 11; i++ {
		if i%2 == 0 {
			b.WithContent(fmt.Sprintf("blog/page%d.md", i), pageContentAndSummaryDivider)
			b.WithContent(fmt.Sprintf("blog/page%d.no.md", i), pageContentAndSummaryDivider)
		} else {
			b.WithContent(fmt.Sprintf("blog/page%d.md", i), pageContentAutoSummary)
		}
	}

	for i := 1; i <= 5; i++ {
		// Root section pages
		b.WithContent(fmt.Sprintf("root%d.md", i), pageContentAutoSummary)
	}

	// https://github.com/gohugoio/hugo/issues/4695
	b.WithContent("blog/markyshort.md", pageContentWithMarkdownShortcodes)

	// Add one bundle
	b.WithContent("blog/mybundle/index.md", pageContentAndSummaryDivider)
	b.WithContent("blog/mybundle/mydata.csv", "Bundled CSV")

	const (
		commonPageTemplate            = `|{{ .Kind }}|{{ .Title }}|{{ .Path }}|{{ .Summary }}|{{ .Content }}|RelPermalink: {{ .RelPermalink }}|WordCount: {{ .WordCount }}|Pages: {{ .Pages }}|Data Pages: Pages({{ len .Data.Pages }})|Resources: {{ len .Resources }}|Summary: {{ .Summary }}`
		commonPaginatorTemplate       = `|Paginator: {{ with .Paginator }}{{ .PageNumber }}{{ else }}NIL{{ end }}`
		commonListTemplateNoPaginator = `|{{ $pages := .Pages }}{{ if .IsHome }}{{ $pages = .Site.RegularPages }}{{ end }}{{ range $i, $e := ($pages | first 1) }}|Render {{ $i }}: {{ .Kind }}|{{ .Render "li" }}|{{ end }}|Site params: {{ $.Site.Params.hugo }}|RelPermalink: {{ .RelPermalink }}`
		commonListTemplate            = commonPaginatorTemplate + `|{{ $pages := .Pages }}{{ if .IsHome }}{{ $pages = .Site.RegularPages }}{{ end }}{{ range $i, $e := ($pages | first 1) }}|Render {{ $i }}: {{ .Kind }}|{{ .Render "li" }}|{{ end }}|Site params: {{ $.Site.Params.hugo }}|RelPermalink: {{ .RelPermalink }}`
		commonShortcodeTemplate       = `|{{ .Name }}|{{ .Ordinal }}|{{ .Page.Summary }}|{{ .Page.Content }}|WordCount: {{ .Page.WordCount }}`
		prevNextTemplate              = `|Prev: {{ with .Prev }}{{ .RelPermalink }}{{ end }}|Next: {{ with .Next }}{{ .RelPermalink }}{{ end }}`
		prevNextInSectionTemplate     = `|PrevInSection: {{ with .PrevInSection }}{{ .RelPermalink }}{{ end }}|NextInSection: {{ with .NextInSection }}{{ .RelPermalink }}{{ end }}`
		paramsTemplate                = `|Params: {{ .Params.hugo }}`
		treeNavTemplate               = `|CurrentSection: {{ .CurrentSection }}`
	)

	b.WithTemplates(
		"_default/list.html", "HTML: List"+commonPageTemplate+commonListTemplate+"|First Site: {{ .Sites.First.Title }}",
		"_default/list.json", "JSON: List"+commonPageTemplate+commonListTemplateNoPaginator,
		"_default/list.csv", "CSV: List"+commonPageTemplate+commonListTemplateNoPaginator,
		"_default/single.html", "HTML: Single"+commonPageTemplate+prevNextTemplate+prevNextInSectionTemplate+treeNavTemplate,
		"_default/single.json", "JSON: Single"+commonPageTemplate,

		// For .Render test
		"_default/li.html", `HTML: LI|{{ strings.Contains .Content "HTML: Shortcode: sc" }}`+paramsTemplate,
		"_default/li.json", `JSON: LI|{{ strings.Contains .Content "JSON: Shortcode: sc" }}`+paramsTemplate,
		"_default/li.csv", `CSV: LI|{{ strings.Contains .Content "CSV: Shortcode: sc" }}`+paramsTemplate,

		"404.html", "{{ .Kind }}|{{ .Title }}|Page not found",

		"shortcodes/sc.html", "HTML: Shortcode: "+commonShortcodeTemplate,
		"shortcodes/sc.json", "JSON: Shortcode: "+commonShortcodeTemplate,
		"shortcodes/sc.csv", "CSV: Shortcode: "+commonShortcodeTemplate,
	)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/blog/page1/index.html",
		"This is content with some shortcodes.",
		"Page with outputs",
		"Pages: Pages(0)",
		"RelPermalink: /blog/page1/|",
		"Shortcode 1: HTML: Shortcode: |sc|0|||WordCount: 0.",
		"Shortcode 2: HTML: Shortcode: |sc|1|||WordCount: 0.",
		"Prev: /blog/page10/|Next: /blog/mybundle/",
		"PrevInSection: /blog/page10/|NextInSection: /blog/mybundle/",
		"Summary: This is summary.",
		"CurrentSection: Page(/blog)",
	)

	b.AssertFileContent("public/blog/page1/index.json",
		"JSON: Single|page|Page with outputs|",
		"SON: Shortcode: |sc|0||")

	b.AssertFileContent("public/index.html",
		"home|In English",
		"Site params: Rules",
		"Pages: Pages(6)|Data Pages: Pages(6)",
		"Paginator: 1",
		"First Site: In English",
		"RelPermalink: /",
	)

	b.AssertFileContent("public/no/index.html", "home|På norsk", "RelPermalink: /no/")

	// Check RSS
	rssHome := b.FileContent("public/index.xml")
	c.Assert(rssHome, qt.Contains, `<atom:link href="https://example.com/index.xml" rel="self" type="application/rss+xml" />`)
	c.Assert(strings.Count(rssHome, "<item>"), qt.Equals, 3) // rssLimit = 3

	// .Render should use template/content from the current output format
	// even if that output format isn't configured for that page.
	b.AssertFileContent(
		"public/index.json",
		"Render 0: page|JSON: LI|false|Params: Rocks!",
	)

	b.AssertFileContent(
		"public/index.html",
		"Render 0: page|HTML: LI|false|Params: Rocks!|",
	)

	b.AssertFileContent(
		"public/index.csv",
		"Render 0: page|CSV: LI|false|Params: Rocks!|",
	)

	// Check bundled resources
	b.AssertFileContent(
		"public/blog/mybundle/index.html",
		"Resources: 1",
	)

	// Check pages in root section
	b.AssertFileContent(
		"public/root3/index.html",
		"Single|page|Page with outputs|root3.md|",
		"Prev: /root4/|Next: /root2/|PrevInSection: /root4/|NextInSection: /root2/",
	)

	b.AssertFileContent(
		"public/root3/index.json", "Shortcode 1: JSON:")

	// Paginators
	b.AssertFileContent("public/page/1/index.html", `rel="canonical" href="https://example.com/"`)
	b.AssertFileContent("public/page/2/index.html", "HTML: List|home|In English|", "Paginator: 2")

	// 404
	b.AssertFileContent("public/404.html", "404|404 Page not found")

	// Sitemaps
	b.AssertFileContent("public/en/sitemap.xml", "<loc>https://example.com/blog/</loc>")
	b.AssertFileContent("public/no/sitemap.xml", `hreflang="no"`)

	b.AssertFileContent("public/sitemap.xml", "<loc>https://example.com/en/sitemap.xml</loc>", "<loc>https://example.com/no/sitemap.xml</loc>")

	// robots.txt
	b.AssertFileContent("public/robots.txt", `User-agent: *`)

	// Aliases
	b.AssertFileContent("public/a/b/c/index.html", `refresh`)

	// Markdown vs shortcodes
	// Check that all footnotes are grouped (even those from inside the shortcode)
	b.AssertFileContentRe("public/blog/markyshort/index.html", `Footnotes:.*<ol>.*Fn 1.*Fn 2.*Fn 3.*</ol>`)

}

// https://github.com/golang/go/issues/30286
func TestDataRace(t *testing.T) {

	const page = `
---
title: "The Page"
outputs: ["HTML", "JSON"]
---	

The content.
	

	`

	b := newTestSitesBuilder(t).WithSimpleConfigFile()
	for i := 1; i <= 50; i++ {
		b.WithContent(fmt.Sprintf("blog/page%d.md", i), page)
	}

	b.WithContent("_index.md", `
---
title: "The Home"
outputs: ["HTML", "JSON", "CSV", "RSS"]
---	

The content.
	

`)

	commonTemplate := `{{ .Data.Pages }}`

	b.WithTemplatesAdded("_default/single.html", "HTML Single: "+commonTemplate)
	b.WithTemplatesAdded("_default/list.html", "HTML List: "+commonTemplate)

	b.CreateSites().Build(BuildCfg{})
}
