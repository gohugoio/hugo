// Copyright 2024 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless requiredF by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
)

func TestRenderString(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithTemplates("index.html", `
{{ $p := site.GetPage "p1.md" }}
{{ $optBlock := dict "display" "block" }}
{{ $optOrg := dict "markup" "org" }}
RSTART:{{ "**Bold Markdown**" | $p.RenderString }}:REND
RSTART:{{  "**Bold Block Markdown**" | $p.RenderString  $optBlock }}:REND
RSTART:{{  "/italic org mode/" | $p.RenderString  $optOrg }}:REND
RSTART:{{ "## Header2" | $p.RenderString }}:REND


`, "_default/_markup/render-heading.html", "Hook Heading: {{ .Level }}")

	b.WithContent("p1.md", `---
title: "p1"
---
`,
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
RSTART:<strong>Bold Markdown</strong>:REND
RSTART:<p><strong>Bold Block Markdown</strong></p>
RSTART:<em>italic org mode</em>:REND
RSTART:Hook Heading: 2:REND
`)
}

// https://github.com/gohugoio/hugo/issues/6882
func TestRenderStringOnListPage(t *testing.T) {
	renderStringTempl := `
{{ .RenderString "**Hello**" }}
`
	b := newTestSitesBuilder(t)
	b.WithContent("mysection/p1.md", `FOO`)
	b.WithTemplates(
		"index.html", renderStringTempl,
		"_default/list.html", renderStringTempl,
		"_default/single.html", renderStringTempl,
	)

	b.Build(BuildCfg{})

	for _, filename := range []string{
		"index.html",
		"mysection/index.html",
		"categories/index.html",
		"tags/index.html",
		"mysection/p1/index.html",
	} {
		b.AssertFileContent("public/"+filename, `<strong>Hello</strong>`)
	}
}

// Issue 9433
func TestRenderStringOnPageNotBackedByAFile(t *testing.T) {
	t.Parallel()
	logger := loggers.NewDefault()
	b := newTestSitesBuilder(t).WithLogger(logger).WithConfigFile("toml", `
disableKinds = ["page", "section", "taxonomy", "term"]	
`)
	b.WithTemplates("index.html", `{{ .RenderString "**Hello**" }}`).WithContent("p1.md", "")
	b.BuildE(BuildCfg{})
	b.Assert(logger.LoggCount(logg.LevelWarn), qt.Equals, 0)
}

func TestRenderStringWithShortcode(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- config.toml --
title = "Hugo Rocks!"
enableInlineShortcodes = true
-- content/p1/index.md --
---
title: "P1"
---
## First
-- layouts/shortcodes/mark1.md --
{{ .Inner }}
-- layouts/shortcodes/mark2.md --
1. Item Mark2 1
1. Item Mark2 2
   1. Item Mark2 2-1
1. Item Mark2 3
-- layouts/shortcodes/myhthml.html --
Title: {{ .Page.Title }}
TableOfContents: {{ .Page.TableOfContents }}
Page Type: {{ printf "%T" .Page }}
-- layouts/_default/single.html --
{{ .RenderString "Markdown: {{% mark2 %}}|HTML: {{< myhthml >}}|Inline: {{< foo.inline >}}{{ site.Title }}{{< /foo.inline >}}|" }}
HasShortcode: mark2:{{ .HasShortcode "mark2" }}:true
HasShortcode: foo:{{ .HasShortcode "foo" }}:false

`

	t.Run("Basic", func(t *testing.T) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: filesTemplate,
			},
		).Build()

		b.AssertFileContent("public/p1/index.html",
			"<p>Markdown: 1. Item Mark2 1</p>\n<ol>\n<li>Item Mark2 2\n<ol>\n<li>Item Mark2 2-1</li>\n</ol>\n</li>\n<li>Item Mark2 3|",
			"<a href=\"#first\">First</a>", // ToC
			`
HTML: Title: P1
Inline: Hugo Rocks!
HasShortcode: mark2:true:true
HasShortcode: foo:false:false
Page Type: *hugolib.pageForShortcode`,
		)
	})

	t.Run("Edit shortcode", func(t *testing.T) {
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: filesTemplate,
				Running:     true,
			},
		).Build()

		b.EditFiles("layouts/shortcodes/myhthml.html", "Edit shortcode").Build()

		b.AssertFileContent("public/p1/index.html",
			`Edit shortcode`,
		)
	})
}

// Issue 9959
func TestRenderStringWithShortcodeInPageWithNoContentFile(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- layouts/shortcodes/myshort.html --
Page Kind: {{ .Page.Kind }}
-- layouts/index.html --
Short: {{ .RenderString "{{< myshort >}}" }}
Has myshort: {{ .HasShortcode "myshort" }}
Has other: {{ .HasShortcode "other" }}

	`

	b := Test(t, files)

	b.AssertFileContent("public/index.html",
		`
Page Kind: home
Has myshort: true
Has other: false
`)
}

func TestRenderStringWithShortcodeIssue10654(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
timeout = '300ms'
-- content/p1.md --
---
title: "P1"
---
{{< toc >}}

## Heading 1

{{< noop >}}
     {{ not a shortcode
{{< /noop >}}
}
-- layouts/shortcodes/noop.html --
{{ .Inner | $.Page.RenderString }}
-- layouts/shortcodes/toc.html --
{{ .Page.TableOfContents }}
-- layouts/_default/single.html --
{{ .Content }}
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", `TableOfContents`)
}
