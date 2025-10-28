// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"
)

func TestRenderShortcodesBasic(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "term"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### p2-h1
{{< withhtml >}}
### p2-h2
{{% withmarkdown %}}
### p2-h3
{{% include "p3" %}}
-- content/p3.md --
---
title: "p3"
---
### p3-h1
{{< withhtml >}}
### p3-h2
{{% withmarkdown %}}
{{< level3 >}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/withhtml.html --
<div>{{ .Page.Title }} withhtml</div>
-- layouts/shortcodes/withmarkdown.html --
#### {{ .Page.Title }} withmarkdown
-- layouts/shortcodes/level3.html --
Level 3: {{ .Page.Title }}
-- layouts/_default/single.html --
Fragments: {{ .Fragments.Identifiers }}|
HasShortcode Level 1: {{ .HasShortcode "include" }}|
HasShortcode Level 2: {{ .HasShortcode "withmarkdown" }}|
HasShortcode Level 3: {{ .HasShortcode "level3" }}|
HasShortcode not found: {{ .HasShortcode "notfound" }}|
Content: {{ .Content }}|
`

	b := Test(t, files)

	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/p1/index.html",
		"Fragments: [p1-h1 p2-h1 p2-h2 p2-h3 p2-withmarkdown p3-h1 p3-h2 p3-withmarkdown]|",
		"HasShortcode Level 1: true|",
		"HasShortcode Level 2: true|",
		"HasShortcode Level 3: true|",
		"HasShortcode not found: false|",
	)
}

func TestRenderShortcodesNestedMultipleOutputFormatTemplates(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "term", "section", "rss", "sitemap", "robotsTXT", "404"]
[outputs]
page = ["html", "json"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### p2-h1
{{% myshort %}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/myshort.html --
Myshort HTML.
-- layouts/shortcodes/myshort.json --
Myshort JSON.
-- layouts/_default/single.html --
HTML: {{ .Content }}
-- layouts/_default/single.json --
JSON: {{ .Content }}


`

	b := Test(t, files)

	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/p1/index.html", "Myshort HTML")
	b.AssertFileContent("public/p1/index.json", "Myshort JSON")
}

func TestRenderShortcodesEditNested(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term", "section", "rss", "sitemap", "robotsTXT", "404"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### p2-h1
{{% myshort %}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/myshort.html --
Myshort Original.
-- layouts/_default/single.html --
 {{ .Content }}
`
	b := TestRunning(t, files)
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/p1/index.html", "Myshort Original.")

	b.EditFileReplaceAll("layouts/shortcodes/myshort.html", "Original", "Edited").Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/p1/index.html", "Myshort Edited.")
}

func TestRenderShortcodesEditIncludedPage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term", "section", "rss", "sitemap", "robotsTXT", "404"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### Original
{{% myshort %}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/myshort.html --
Myshort Original.
-- layouts/_default/single.html --
 {{ .Content }}



`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			Running:     true,
		},
	).Build()

	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/p1/index.html", "Original")

	b.EditFileReplaceFunc("content/p2.md", func(s string) string {
		return strings.Replace(s, "Original", "Edited", 1)
	})
	b.Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/p1/index.html", "Edited")
}

func TestRenderShortcodesEditSectionContentWithShortcodeInIncludedPageIssue12458(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term", "rss", "sitemap", "robotsTXT", "404"]
-- content/mysection/_index.md --
---
title: "My Section"
---
## p1-h1
{{% include "p2" %}}
-- content/mysection/p2.md --
---
title: "p2"
---
### Original
{{% myshort %}}
-- layouts/shortcodes/include.html --
{{ $p := .Page.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/myshort.html --
Myshort Original.
-- layouts/_default/list.html --
 {{ .Content }}



`
	b := TestRunning(t, files)

	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/mysection/index.html", "p1-h1")
	b.EditFileReplaceAll("content/mysection/_index.md", "p1-h1", "p1-h1 Edited").Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/mysection/index.html", "p1-h1 Edited")
}

func TestRenderShortcodesNestedPageContextIssue12356(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "robotsTXT", "404"]
-- layouts/_default/_markup/render-image.html --
{{- with .PageInner.Resources.Get .Destination -}}Image: {{ .RelPermalink }}|{{- end -}}
-- layouts/_default/_markup/render-link.html --
{{- with .PageInner.GetPage .Destination -}}Link: {{ .RelPermalink }}|{{- end -}}
-- layouts/_default/_markup/render-heading.html --
Heading: {{ .PageInner.Title }}: {{ .PlainText }}|
-- layouts/_default/_markup/render-codeblock.html --
CodeBlock: {{ .PageInner.Title }}: {{ .Type }}|
-- layouts/_default/list.html --
Content:{{ .Content }}|
Fragments: {{ with .Fragments }}{{.Identifiers }}{{ end }}|
-- layouts/_default/single.html --
Content:{{ .Content }}|
-- layouts/shortcodes/include.html --
{{ with site.GetPage (.Get 0) }}
  {{ .RenderShortcodes }}
{{ end }}
-- content/markdown/_index.md --
---
title: "Markdown"
---
# H1
|{{% include "/posts/p1" %}}|
![kitten](pixel3.png "Pixel 3")

§§§go
fmt.Println("Hello")
§§§

-- content/markdown2/_index.md --
---
title: "Markdown 2"
---
|{{< include "/posts/p1" >}}|
-- content/html/_index.html --
---
title: "HTML"
---
|{{% include "/posts/p1" %}}|

-- content/posts/p1/index.md --
---
title: "p1"
---
## H2-p1
![kitten](pixel1.png "Pixel 1")
![kitten](pixel2.png "Pixel 2")
[p2](p2)

§§§bash
echo "Hello"
§§§

-- content/posts/p2/index.md --
---
title: "p2"
---
-- content/posts/p1/pixel1.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- content/posts/p1/pixel2.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- content/markdown/pixel3.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- content/html/pixel4.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==

`

	b := Test(t, files)

	b.AssertNoRenderShortcodesArtifacts()

	b.AssertFileContent("public/markdown/index.html",
		// Images.
		"Image: /posts/p1/pixel1.png|\nImage: /posts/p1/pixel2.png|\n|\nImage: /markdown/pixel3.png|</p>\n|",
		// Links.
		"Link: /posts/p2/|",
		// Code blocks
		"CodeBlock: p1: bash|", "CodeBlock: Markdown: go|",
		// Headings.
		"Heading: Markdown: H1|", "Heading: p1: H2-p1|",
		// Fragments.
		"Fragments: [h1 h2-p1]|",
		// Check that the special context markup is not rendered.
		"! hugo_ctx",
	)

	b.AssertFileContent("public/markdown2/index.html", "! hugo_ctx", "Content:<p>|\n  ![kitten](pixel1.png \"Pixel 1\")\n![kitten](pixel2.png \"Pixel 2\")\n|</p>\n|")

	b.AssertFileContent("public/html/index.html", "! hugo_ctx")
}

// Issue 12854.
func TestRenderShortcodesWithHTML(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term"]
markup.goldmark.renderer.unsafe = true
-- content/p1.md --
---
title: "p1"
---
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
Hello <b>world</b>. Some **bold** text. Some Unicode: 神真美好.
-- layouts/shortcodes/include.html --
{{ with site.GetPage (.Get 0) }}
<div>{{ .RenderShortcodes }}</div>
{{ end }}
-- layouts/_default/single.html --
{{ .Content }}
`

	b := TestRunning(t, files, TestOptWarn())

	b.AssertNoRenderShortcodesArtifacts()
	b.AssertLogContains(filepath.ToSlash("WARN  .RenderShortcodes detected inside HTML block in \"/content/p1.md\"; this may not be what you intended, see https://gohugo.io/docs/reference/methods/page/rendershortcodes/#limitations\nYou can suppress this warning by adding the following to your site configuration:\nignoreLogs = ['warning-rendershortcodes-in-html']"))
	b.AssertFileContent("public/p1/index.html", "<div>Hello <b>world</b>. Some **bold** text. Some Unicode: 神真美好.\n</div>")
	b.EditFileReplaceAll("content/p2.md", "Hello", "Hello Edited").Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContent("public/p1/index.html", "<div>Hello Edited <b>world</b>. Some **bold** text. Some Unicode: 神真美好.\n</div>")
}

func TestRenderShortcodesIncludeMarkdownFileWithoutTrailingNewline(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term"]
markup.goldmark.renderer.unsafe = true
-- content/p1.md --
---
title: "p1"
---
Content p1 id-1000.{{% include "p2" %}}{{% include "p3" %}}

§§§ go
code_p1
§§§
§§§ go
code_p1_2
§§§

§§§ go
code_p1_3
§§§
-- content/p2.md --
---
title: "p2"
---
§§§ bash
code_p2
§§§
Foo.
-- content/p3.md --
---
title: "p3"
---
§§§ php
code_p3
§§§
-- layouts/shortcodes/include.html --
{{ with site.GetPage (.Get 0) -}}
{{ .RenderShortcodes -}}
{{ end -}}
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-codeblock.html --
<code>{{ .Inner | safeHTML }}</code>
`

	b := TestRunning(t, files, TestOptWarn())

	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContentEquals("public/p1/index.html", "<p>Content p1 id-1000.</p>\n<code>code_p2</code><p>Foo.</p>\n<code>code_p3</code><code>code_p1</code><code>code_p1_2</code><code>code_p1_3</code>")
	b.EditFileReplaceAll("content/p1.md", "id-1000.", "id-100.").Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContentEquals("public/p1/index.html", "<p>Content p1 id-100.</p>\n<code>code_p2</code><p>Foo.</p>\n<code>code_p3</code><code>code_p1</code><code>code_p1_2</code><code>code_p1_3</code>")
	b.EditFileReplaceAll("content/p2.md", "code_p2", "codep2").Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContentEquals("public/p1/index.html", "<p>Content p1 id-100.</p>\n<code>codep2</code><p>Foo.</p>\n<code>code_p3</code><code>code_p1</code><code>code_p1_2</code><code>code_p1_3</code>")
	b.EditFileReplaceAll("content/p3.md", "code_p3", "code_p3_edited").Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContentEquals("public/p1/index.html", "<p>Content p1 id-100.</p>\n<code>codep2</code><p>Foo.</p>\n<code>code_p3_edited</code><code>code_p1</code><code>code_p1_2</code><code>code_p1_3</code>")
}

// Issue 13004.
func TestRenderShortcodesIncludeShortRefEdit(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term", "section", "rss", "sitemap", "robotsTXT", "404"]
-- content/first/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/second/p2.md --
---
title: "p2"
---
### p2-h1

This is some **markup**.
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) -}}
{{ $p.RenderShortcodes -}}
-- layouts/_default/single.html --
{{ .Content }}
`
	b := TestRunning(t, files)
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContentEquals("public/first/p1/index.html", "<h2 id=\"p1-h1\">p1-h1</h2>\n<h3 id=\"p2-h1\">p2-h1</h3>\n<p>This is some <strong>markup</strong>.</p>\n")
	b.EditFileReplaceAll("content/second/p2.md", "p2-h1", "p2-h1-edited").Build()
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContentEquals("public/first/p1/index.html", "<h2 id=\"p1-h1\">p1-h1</h2>\n<h3 id=\"p2-h1-edited\">p2-h1-edited</h3>\n<p>This is some <strong>markup</strong>.</p>\n")
}

// Issue 13051.
func TestRenderShortcodesEmptyParagraph(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['section','rss','sitemap','taxonomy','term']
-- layouts/_default/home.html --
{{ .Content }}
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/shortcodes/include.html --
 {{ with site.GetPage (.Get 0) }}
  {{ .RenderShortcodes }}
{{ end }}
-- content/_index.md --
---
title: home
---

a

{{% include "/snippet" %}}

b

-- content/snippet.md --
---
title: snippet
build:
  render: never
  list: never
---

_emphasized_

not emphasized

`

	b := Test(t, files)
	b.AssertNoRenderShortcodesArtifacts()
	b.AssertFileContentEquals("public/index.html",
		"<p>a</p>\n<p><em>emphasized</em></p>\n<p>not emphasized</p>\n<p>b</p>\n",
	)
}
