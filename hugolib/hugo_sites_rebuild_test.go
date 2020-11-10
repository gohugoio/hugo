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
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
)

func TestRebuildAddPageToSection(t *testing.T) {
	c := qt.New(t)

	files := `
-- config.toml --
disableKinds=["home", "taxonomy", "term", "sitemap", "robotsTXT"]
[outputs]
	section = ['HTML']
	page = ['HTML']
-- content/blog/b1.md --
-- content/blog/b3.md --
-- content/doc/d1.md --
-- content/doc/d3.md --
-- layouts/_default/single.html --
{{ .Pathc }}
-- layouts/_default/list.html --
List:
{{ range $i, $e := .RegularPages }}
{{ $i }}: {{ .Pathc }}
{{ end }}

`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           c,
			TxtarString: files,
			Running:     true,
		},
	).Build()

	b.AssertRenderCountPage(6)
	b.AssertFileContent("public/blog/index.html", `
0: /blog/b1
1: /blog/b3
`)

	b.AddFiles("content/blog/b2.md", "").Build()
	b.AssertFileContent("public/blog/index.html", `
0: /blog/b1
1: /blog/b2
2: /blog/b3
`)

	// The 3 sections.
	b.AssertRenderCountPage(3)
}

func TestRebuildAddPageToSectionListItFromAnotherSection(t *testing.T) {
	c := qt.New(t)

	files := `
-- config.toml --
disableKinds=["home", "taxonomy", "term", "sitemap", "robotsTXT"]
[outputs]
	section = ['HTML']
	page = ['HTML']
-- content/blog/b1.md --
-- content/blog/b3.md --
-- content/doc/d1.md --
-- content/doc/d3.md --
-- layouts/_default/single.html --
{{ .Pathc }}
-- layouts/_default/list.html --
List Default
-- layouts/doc/list.html --
{{ $blog := site.GetPage "blog" }}
List Doc:
{{ range $i, $e := $blog.RegularPages }}
{{ $i }}: {{ .Pathc }}
{{ end }}

`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           c,
			TxtarString: files,
			Running:     true,
		},
	).Build()

	b.AssertRenderCountPage(6)
	b.AssertFileContent("public/doc/index.html", `
0: /blog/b1
1: /blog/b3
`)

	b.AddFiles("content/blog/b2.md", "").Build()
	b.AssertFileContent("public/doc/index.html", `
0: /blog/b1
1: /blog/b2
2: /blog/b3
`)

	// Just the 3 sections.
	b.AssertRenderCountPage(3)
}

func TestRebuildChangePartialUsedInShortcode(t *testing.T) {
	c := qt.New(t)

	files := `
-- config.toml --
disableKinds=["home", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
[outputs]
	page = ['HTML']
-- content/blog/p1.md --
Shortcode: {{< c >}}
-- content/blog/p2.md --
CONTENT
-- layouts/_default/single.html --
{{ .Pathc }}: {{ .Content }}
-- layouts/shortcodes/c.html --
{{ partial "p.html" . }}
-- layouts/partials/p.html --
MYPARTIAL

`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           c,
			TxtarString: files,
			Running:     true,
		},
	).Build()

	b.AssertRenderCountPage(2)
	b.AssertFileContent("public/blog/p1/index.html", `/blog/p1: <p>Shortcode: MYPARTIAL`)

	b.EditFiles("layouts/partials/p.html", "MYPARTIAL CHANGED").Build()

	b.AssertRenderCountPage(1)
	b.AssertFileContent("public/blog/p1/index.html", `/blog/p1: <p>Shortcode: MYPARTIAL CHANGED`)
}

func TestRebuildEditPartials(t *testing.T) {
	c := qt.New(t)

	files := `
-- config.toml --
disableKinds=["home", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
[outputs]
	page = ['HTML']
-- content/blog/p1.md --
Shortcode: {{< c >}}
-- content/blog/p2.md --
CONTENT
-- content/blog/p3.md --
Shortcode: {{< d >}}
-- content/blog/p4.md --
Shortcode: {{< d >}}
-- content/blog/p5.md --
Shortcode: {{< d >}}
-- content/blog/p6.md --
Shortcode: {{< d >}}
-- content/blog/p7.md --
Shortcode: {{< d >}}
-- layouts/_default/single.html --
{{ .Pathc }}: {{ .Content }}
-- layouts/shortcodes/c.html --
{{ partial "p.html" . }}
-- layouts/shortcodes/d.html --
{{ partialCached "p.html" . }}
-- layouts/partials/p.html --
MYPARTIAL

`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           c,
			TxtarString: files,
			Running:     true,
		},
	).Build()

	b.AssertRenderCountPage(7)
	b.AssertFileContent("public/blog/p1/index.html", `/blog/p1: <p>Shortcode: MYPARTIAL`)
	b.AssertFileContent("public/blog/p3/index.html", `/blog/p3: <p>Shortcode: MYPARTIAL`)

	b.EditFiles("layouts/partials/p.html", "MYPARTIAL CHANGED").Build()

	b.AssertRenderCountPage(6)
	b.AssertFileContent("public/blog/p1/index.html", `/blog/p1: <p>Shortcode: MYPARTIAL CHANGED`)
	b.AssertFileContent("public/blog/p3/index.html", `/blog/p3: <p>Shortcode: MYPARTIAL CHANGED`)
	b.AssertFileContent("public/blog/p4/index.html", `/blog/p4: <p>Shortcode: MYPARTIAL CHANGED`)
}

// bookmark1
func TestRebuildBasic(t *testing.T) {
	// TODO1
	pinnedTestCase := ""
	tt := htesting.NewPinnedRunner(t, pinnedTestCase)

	var (
		twoPagesAndHomeDataInP1 = `
-- config.toml --
disableKinds=["section", "taxonomy", "term", "sitemap", "robotsTXT"]
[permalinks]
"/"="/:filename/"
[outputs]
  home = ['HTML']
  page = ['HTML']
-- data/mydata.toml --
hugo="Rocks!"
-- content/p1.md --
---
includeData: true
---
CONTENT
-- content/p2.md --
CONTENT
-- layouts/_default/single.html --
{{ if .Params.includeData }}
Hugo {{ site.Data.mydata.hugo }}
{{ else }}
NO DATA USED
{{ end }}
Title: {{ .Title }}|Content Start: {{ .Content }}:End:
-- layouts/index.html --
Home: Len site.Pages: {{ len site.Pages}}|Len site.RegularPages: {{ len site.RegularPages}}|Len site.AllPages: {{ len site.AllPages}}:End:
`

		twoPagesDataInShortcodeInP2HTMLAndRSS = `
-- config.toml --
disableKinds=["home", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
[outputs]
  page = ['HTML', 'RSS']
-- data/mydata.toml --
hugo="Rocks!"
-- content/p1.md --
---
slug: p1
---
CONTENT
-- content/p2.md --
---
slug: p2
---
{{< foo >}}
CONTENT
-- layouts/_default/single.html --
HTML: {{ .Slug }}: {{ .Content }}
-- layouts/_default/single.xml --
XML: {{ .Slug }}: {{ .Content }}
-- layouts/shortcodes/foo.html --
Hugo {{ site.Data.mydata.hugo }}
-- layouts/shortcodes/foo.xml --
No Data
`

		twoPagesDataInRenderHookInP2 = `
-- config.toml --
disableKinds=["home", "section", "taxonomy", "term", "sitemap", "robotsTXT"]
-- data/mydata.toml --
hugo="Rocks!"
-- content/p1.md --
---
slug: p1
---
-- content/p2.md --
---
slug: p2
---
[Text](https://www.gohugo.io "Title")
-- layouts/_default/single.html --
{{ .Slug }}: {{ .Content }}
-- layouts/_default/_markup/render-link.html --
Hugo {{ site.Data.mydata.hugo }}
`

		twoPagesAndHomeWithBaseTemplate = `
-- config.toml --
disableKinds=[ "section", "taxonomy", "term", "sitemap", "robotsTXT"]
[outputs]
  home = ['HTML']
  page = ['HTML']
-- data/mydata.toml --
hugo="Rocks!"
-- content/_index.md --
---
title: MyHome
---
-- content/p1.md --
---
slug: p1
---
-- content/p2.md --
---
slug: p2
---
-- layouts/_default/baseof.html --
Block Main Start:{{ block "main" . }}{{ end }}:End:
-- layouts/_default/single.html --
{{ define "main" }}Single Main Start:{{ .Slug }}: {{ .Content }}:End:{{ end }}
-- layouts/_default/list.html --
{{ define "main" }}List Main Start:{{ .Title }}: {{ .Content }}:End{{ end }}
`
	)

	// * Remove doc
	// * Add
	// * Rename file
	// * Change doc
	// * Change a template
	// * Change language file
	// OK * Site.LastChange - mod, no mod

	// Tests for  Site.LastChange
	for _, changeSiteLastChanged := range []bool{false, true} {
		name := "Site.LastChange"
		if changeSiteLastChanged {
			name += " Changed"
		} else {
			name += " Not Changed"
		}

		const files = `
-- config.toml --
disableKinds=["section", "taxonomy", "term", "sitemap", "robotsTXT", "404"]
[outputs]
	home = ['HTML']
	page = ['HTML']
-- content/_index.md --
---
title: Home
lastMod: 2020-02-01
---
-- content/p1.md --
---
title: P1
lastMod: 2020-03-01
---
CONTENT
-- content/p2.md --
---
title: P2
lastMod: 2020-03-02
---
CONTENT
-- layouts/_default/single.html --
Title: {{ .Title }}|Lastmod: {{ .Lastmod.Format "2006-01-02" }}|Content Start: {{ .Content }}:End:
-- layouts/index.html --
Home: Lastmod: {{ .Lastmod.Format "2006-01-02" }}|site.LastChange: {{ site.LastChange.Format "2006-01-02" }}:End:
		`

		tt.Run(name, func(c *qt.C) {
			b := NewIntegrationTestBuilder(
				IntegrationTestConfig{
					T:           c,
					TxtarString: files,
					Running:     true,
				},
			).Build()

			b.AssertFileContent("public/p1/index.html", "Title: P1|Lastmod: 2020-03-01")
			b.AssertFileContent("public/index.html", "Home: Lastmod: 2020-02-01|site.LastChange: 2020-03-02")
			b.AssertRenderCountPage(3)

			if changeSiteLastChanged {
				b.EditFileReplace("content/p1.md", func(s string) string { return strings.ReplaceAll(s, "lastMod: 2020-03-01", "lastMod: 2020-05-01") })
			} else {
				b.EditFileReplace("content/p1.md", func(s string) string { return strings.ReplaceAll(s, "CONTENT", "Content Changed") })
			}

			b.Build()

			if changeSiteLastChanged {
				b.AssertFileContent("public/p1/index.html", "Title: P1|Lastmod: 2020-05-01")
				b.AssertFileContent("public/index.html", "Home: Lastmod: 2020-02-01|site.LastChange: 2020-05-01")
				b.AssertRenderCountPage(2)
			} else {
				b.AssertRenderCountPage(1)
				b.AssertFileContent("public/p1/index.html", "Content Changed")

			}
		})
	}

	tt.Run("Content Edit, Add, Rename, Remove", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesAndHomeDataInP1,
				Running:     true,
			},
		).Build()

		b.AssertFileContent("public/p1/index.html", "Hugo Rocks!")
		b.AssertFileContent("public/index.html", `Home: Len site.Pages: 3|Len site.RegularPages: 2|Len site.AllPages: 3:End:`)
		b.AssertRenderCountPage(3)
		b.AssertBuildCountData(1)
		b.AssertBuildCountLayouts(1)

		// Edit
		b.EditFileReplace("content/p1.md", func(s string) string { return strings.ReplaceAll(s, "CONTENT", "Changed Content") }).Build()

		b.AssertFileContent("public/p1/index.html", "Changed Content")
		b.AssertRenderCountPage(1)
		b.AssertRenderCountContent(1)
		b.AssertBuildCountData(1)
		b.AssertBuildCountLayouts(1)

		b.AddFiles("content/p3.md", `ADDED`).Build()
		b.AssertFileContent("public/index.html", `Home: Len site.Pages: 4|Len site.RegularPages: 3|Len site.AllPages: 4:End:`)

		// Remove
		b.RemoveFiles("content/p1.md").Build()

		b.AssertFileContent("public/index.html", `Home: Len site.Pages: 3|Len site.RegularPages: 2|Len site.AllPages: 3:End:`)
		b.AssertRenderCountPage(1)
		b.AssertRenderCountContent(0)
		b.AssertBuildCountData(1)
		b.AssertBuildCountLayouts(1)

		// Rename
		b.RenameFile("content/p2.md", "content/p2n.md").Build()

		b.AssertFileContent("public/index.html", `Home: Len site.Pages: 3|Len site.RegularPages: 2|Len site.AllPages: 3:End:`)
		b.AssertFileContent("public/p2n/index.html", "NO DATA USED")
		b.AssertRenderCountPage(2)
		b.AssertRenderCountContent(1)
		b.AssertBuildCountData(1)
		b.AssertBuildCountLayouts(1)
	})

	tt.Run("Data in page template", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesAndHomeDataInP1,
				Running:     true,
			},
		).Build()

		b.AssertFileContent("public/p1/index.html", "Hugo Rocks!")
		b.AssertFileContent("public/p2/index.html", "NO DATA USED")
		b.AssertRenderCountPage(3)
		b.AssertBuildCountData(1)
		b.AssertBuildCountLayouts(1)

		b.EditFiles("data/mydata.toml", `hugo="Rules!"`).Build()

		b.AssertFileContent("public/p1/index.html", "Hugo Rules!")

		b.AssertBuildCountData(2)
		b.AssertBuildCountLayouts(1)
		b.AssertRenderCountPage(1) // We only need to re-render the one page that uses site.Data.
	})

	tt.Run("Data in shortcode", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesDataInShortcodeInP2HTMLAndRSS,
				Running:     true,
			},
		).Build()

		b.AssertFileContent("public/p2/index.html", "Hugo Rocks!")
		b.AssertFileContent("public/p2/index.xml", "No Data")

		b.AssertRenderCountContent(3) // p2 (2 variants), p1
		b.AssertRenderCountPage(4)    // p2 (2), p1 (2)
		b.AssertBuildCountData(1)
		b.AssertBuildCountLayouts(1)

		b.EditFiles("data/mydata.toml", `hugo="Rules!"`).Build()

		b.AssertFileContent("public/p2/index.html", "Hugo Rules!")
		b.AssertFileContent("public/p2/index.xml", "No Data")

		// We only need to re-render the one page that uses the shortcode with site.Data (p2)
		b.AssertRenderCountContent(1)
		b.AssertRenderCountPage(1)
		b.AssertBuildCountData(2)
		b.AssertBuildCountLayouts(1)
	})

	// TODO1 site date(s).

	tt.Run("Layout Shortcode", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesDataInShortcodeInP2HTMLAndRSS,
				Running:     true,
			},
		).Build()

		b.AssertBuildCountLayouts(1)
		b.AssertBuildCountData(1)

		b.EditFiles("layouts/shortcodes/foo.html", `Shortcode changed"`).Build()

		b.AssertFileContent("public/p2/index.html", "Shortcode changed")
		b.AssertRenderCountContent(1)
		b.AssertRenderCountPage(1)
		b.AssertBuildCountLayouts(2)
		b.AssertBuildCountData(1)
	})

	tt.Run("Data in Render Hook", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesDataInRenderHookInP2,
				Running:     true,
			},
		).Build()

		b.AssertFileContent("public/p2/index.html", "Hugo Rocks!")
		b.AssertBuildCountData(1)

		b.EditFiles("data/mydata.toml", `hugo="Rules!"`).Build()

		b.AssertFileContent("public/p2/index.html", "Hugo Rules!")
		// We only need to re-render the one page that contains a link (p2)
		b.AssertRenderCountContent(1)
		b.AssertRenderCountPage(1)
		b.AssertBuildCountData(2)
	})

	tt.Run("Layout Single", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesAndHomeWithBaseTemplate,
				Running:     true,
			},
		).Build()

		b.EditFiles("layouts/_default/single.html", `Single template changed"`).Build()
		b.AssertFileContent("public/p1/index.html", "Single template changed")
		b.AssertFileContent("public/p2/index.html", "Single template changed")
		b.AssertRenderCountContent(0) // Reuse .Content
		b.AssertRenderCountPage(2)    // Re-render both pages using single.html
	})

	tt.Run("Layout List", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesAndHomeWithBaseTemplate,
				Running:     true,
			},
		).Build()

		b.EditFiles("layouts/_default/list.html", `List template changed"`).Build()
		b.AssertFileContent("public/index.html", "List template changed")
		b.AssertFileContent("public/p2/index.html", "Block Main Start:Single Main Start:p2: :End::End:")
		b.AssertRenderCountContent(0) // Reuse .Content
		b.AssertRenderCountPage(1)    // Re-render home page only
	})

	tt.Run("Layout Base", func(c *qt.C) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           c,
				TxtarString: twoPagesAndHomeWithBaseTemplate,
				Running:     true,
			},
		).Build()

		b.AssertFileContent("public/index.html", "Block Main Start:List Main Start:MyHome: :End:End:")
		b.EditFiles("layouts/_default/baseof.html", `Block Main Changed Start:{{ block "main" . }}{{ end }}:End:"`).Build()
		b.AssertFileContent("public/index.html", "Block Main Changed Start:List Main Start:MyHome: :End:End:")
		b.AssertFileContent("public/p2/index.html", "Block Main Changed Start:Single Main Start:p2: :End::End:")
		b.AssertRenderCountContent(0) // Reuse .Content
		b.AssertRenderCountPage(3)    // Re-render all 3 pages
	})
}
