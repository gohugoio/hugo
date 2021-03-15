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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestSitesRebuild(t *testing.T) {
	configFile := `
baseURL = "https://example.com"
title = "Rebuild this"
contentDir = "content"
enableInlineShortcodes = true
timeout = "5s"


`

	var (
		contentFilename = "content/blog/page1.md"
		dataFilename    = "data/mydata.toml"
	)

	createSiteBuilder := func(t testing.TB) *sitesBuilder {
		b := newTestSitesBuilder(t).WithConfigFile("toml", configFile).Running()

		b.WithSourceFile(dataFilename, `hugo = "Rocks!"`)

		b.WithContent("content/_index.md", `---
title: Home, Sweet Home!
---

`)

		b.WithContent(contentFilename, `
---
title: "Page 1"
summary: "Initial summary"
paginate: 3
---

Content.

{{< badge.inline >}}
Data Inline: {{ site.Data.mydata.hugo }}
{{< /badge.inline >}}
`)

		// For .Page.Render tests
		b.WithContent("prender.md", `---
title: Page 1
---

Content for Page 1.

{{< dorender >}}

`)

		b.WithTemplatesAdded(
			"layouts/shortcodes/dorender.html", `
{{ $p := .Page }}
Render {{ $p.RelPermalink }}: {{ $p.Render "single" }}

`)

		b.WithTemplatesAdded("index.html", `
{{ range (.Paginate .Site.RegularPages).Pages }}
* Page Paginate: {{ .Title }}|Summary: {{ .Summary }}|Content: {{ .Content }}
{{ end }}
{{ range .Site.RegularPages }}
* Page Pages: {{ .Title }}|Summary: {{ .Summary }}|Content: {{ .Content }}
{{ end }}
Content: {{ .Content }}
Data: {{ site.Data.mydata.hugo }}
`)

		b.WithTemplatesAdded("layouts/partials/mypartial1.html", `Mypartial1`)
		b.WithTemplatesAdded("layouts/partials/mypartial2.html", `Mypartial2`)
		b.WithTemplatesAdded("layouts/partials/mypartial3.html", `Mypartial3`)
		b.WithTemplatesAdded("_default/single.html", `{{ define "main" }}Single Main: {{ .Title }}|Mypartial1: {{ partial "mypartial1.html" }}{{ end }}`)
		b.WithTemplatesAdded("_default/list.html", `{{ define "main" }}List Main: {{ .Title }}{{ end }}`)
		b.WithTemplatesAdded("_default/baseof.html", `Baseof:{{ block "main" . }}Baseof Main{{ end }}|Mypartial3: {{ partial "mypartial3.html" }}:END`)

		return b
	}

	t.Run("Refresh paginator on edit", func(t *testing.T) {
		b := createSiteBuilder(t)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", "* Page Paginate: Page 1|Summary: Initial summary|Content: <p>Content.</p>")

		b.EditFiles(contentFilename, `
---
title: "Page 1 edit"
summary: "Edited summary"
---

Edited content.

`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", "* Page Paginate: Page 1 edit|Summary: Edited summary|Content: <p>Edited content.</p>")
		// https://github.com/gohugoio/hugo/issues/5833
		b.AssertFileContent("public/index.html", "* Page Pages: Page 1 edit|Summary: Edited summary|Content: <p>Edited content.</p>")
	})

	// https://github.com/gohugoio/hugo/issues/6768
	t.Run("Edit data", func(t *testing.T) {
		b := createSiteBuilder(t)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
Data: Rocks!
Data Inline: Rocks!
`)

		b.EditFiles(dataFilename, `hugo = "Rules!"`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
Data: Rules!
Data Inline: Rules!`)
	})

	// https://github.com/gohugoio/hugo/issues/6968
	t.Run("Edit single.html with base", func(t *testing.T) {
		b := newTestSitesBuilder(t).Running()

		b.WithTemplates(
			"_default/single.html", `{{ define "main" }}Single{{ end }}`,
			"_default/baseof.html", `Base: {{ block "main"  .}}Block{{ end }}`,
		)

		b.WithContent("p1.md", "---\ntitle: Page\n---")

		b.Build(BuildCfg{})

		b.EditFiles("layouts/_default/single.html", `Single Edit: {{ define "main" }}Single{{ end }}`)

		counters := &testCounters{}

		b.Build(BuildCfg{testCounters: counters})

		b.Assert(int(counters.contentRenderCounter), qt.Equals, 0)
	})

	t.Run("Page.Render, edit baseof", func(t *testing.T) {
		b := createSiteBuilder(t)

		b.WithTemplatesAdded("index.html", `
{{ $p := site.GetPage "prender.md" }}
prender: {{ $p.Title }}|{{ $p.Content }}

`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
 Render /prender/: Baseof:Single Main: Page 1|Mypartial1: Mypartial1|Mypartial3: Mypartial3:END
`)

		b.EditFiles("layouts/_default/baseof.html", `Baseof Edited:{{ block "main" . }}Baseof Main{{ end }}:END`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
Render /prender/: Baseof Edited:Single Main: Page 1|Mypartial1: Mypartial1:END
`)
	})

	t.Run("Page.Render, edit partial in baseof", func(t *testing.T) {
		b := createSiteBuilder(t)

		b.WithTemplatesAdded("index.html", `
{{ $p := site.GetPage "prender.md" }}
prender: {{ $p.Title }}|{{ $p.Content }}

`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
 Render /prender/: Baseof:Single Main: Page 1|Mypartial1: Mypartial1|Mypartial3: Mypartial3:END
`)

		b.EditFiles("layouts/partials/mypartial3.html", `Mypartial3 Edited`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
Render /prender/: Baseof:Single Main: Page 1|Mypartial1: Mypartial1|Mypartial3: Mypartial3 Edited:END
`)
	})

	t.Run("Edit RSS shortcode", func(t *testing.T) {
		b := createSiteBuilder(t)

		b.WithContent("output.md", `---
title: Output
outputs: ["HTML", "AMP"]
layout: output
---

Content for Output.

{{< output >}}

`)

		b.WithTemplates(
			"layouts/_default/output.html", `Output HTML: {{ .RelPermalink }}|{{ .Content }}`,
			"layouts/_default/output.amp.html", `Output AMP: {{ .RelPermalink }}|{{ .Content }}`,
			"layouts/shortcodes/output.html", `Output Shortcode HTML`,
			"layouts/shortcodes/output.amp.html", `Output Shortcode AMP`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/output/index.html", `
Output Shortcode HTML
`)
		b.AssertFileContent("public/amp/output/index.html", `
Output Shortcode AMP
`)

		b.EditFiles("layouts/shortcodes/output.amp.html", `Output Shortcode AMP Edited`)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/amp/output/index.html", `
Output Shortcode AMP Edited
`)
	})
}

// Issues #7623 #7625
func TestSitesRebuildOnFilesIncludedWithGetPage(t *testing.T) {
	b := newTestSitesBuilder(t).Running()
	b.WithContent("pages/p1.md", `---
title: p1
---
P3: {{< GetPage "pages/p3" >}}
`)

	b.WithContent("pages/p2.md", `---
title: p2
---
P4: {{< site_GetPage "pages/p4" >}}
P5: {{< site_GetPage "p5" >}}
P6: {{< dot_site_GetPage "p6" >}}
`)

	b.WithContent("pages/p3/index.md", "---\ntitle: p3\nheadless: true\n---\nP3 content")
	b.WithContent("pages/p4/index.md", "---\ntitle: p4\nheadless: true\n---\nP4 content")
	b.WithContent("pages/p5.md", "---\ntitle: p5\n---\nP5 content")
	b.WithContent("pages/p6.md", "---\ntitle: p6\n---\nP6 content")

	b.WithTemplates(
		"_default/single.html", `{{ .Content }}`,
		"shortcodes/GetPage.html", `
{{ $arg := .Get 0 }}
{{ $p := .Page.GetPage $arg }}
{{ $p.Content }}
	`,
		"shortcodes/site_GetPage.html", `
{{ $arg := .Get 0 }}
{{ $p := site.GetPage $arg }}
{{ $p.Content }}
	`, "shortcodes/dot_site_GetPage.html", `
{{ $arg := .Get 0 }}
{{ $p := .Site.GetPage $arg }}
{{ $p.Content }}
	`,
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/pages/p1/index.html", "P3 content")
	b.AssertFileContent("public/pages/p2/index.html", `P4 content
P5 content
P6 content
`)

	b.EditFiles("content/pages/p3/index.md", "---\ntitle: p3\n---\nP3 changed content")
	b.EditFiles("content/pages/p4/index.md", "---\ntitle: p4\n---\nP4 changed content")
	b.EditFiles("content/pages/p5.md", "---\ntitle: p5\n---\nP5 changed content")
	b.EditFiles("content/pages/p6.md", "---\ntitle: p6\n---\nP6 changed content")

	b.Build(BuildCfg{})

	b.AssertFileContent("public/pages/p1/index.html", "P3 changed content")
	b.AssertFileContent("public/pages/p2/index.html", `P4 changed content
P5 changed content
P6 changed content
`)
}
