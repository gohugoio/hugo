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

package page_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestThatPageIsAvailableEverywhere(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- config.toml --
baseURL = 'http://example.com/'
disableKinds = ["taxonomy", "term"]
enableInlineShortcodes = true
enableRobotsTXT = true
[pagination]
pagerSize = 1
LANG_CONFIG
-- content/_index.md --
---
title: "Home"
aliases: ["/homealias/"]
---
{{< shortcode "Angled Brackets" >}}
{{% shortcode "Percentage" %}}

{{< outer >}}
{{< inner >}}
{{< /outer >}}

{{< foo.inline >}}{{ if page.IsHome }}Shortcode Inline OK.{{ end }}{{< /foo.inline >}}

## Heading

[I'm an inline-style link](https://www.google.com)

![alt text](https://github.com/adam-p/markdown-here/raw/master/src/common/images/icon48.png "Logo Title Text 1")

$$$bash
echo "hello";
$$$

-- content/p1.md --
-- content/p2/index.md --
-- content/p2/p2_1.md --
---
title: "P2_1"
---
{{< foo.inline >}}{{ if page.IsHome }}Shortcode in bundled page OK.{{ else}}Failed.{{ end }}{{< /foo.inline >}}
-- content/p3.md --
-- layouts/_default/_markup/render-heading.html --
{{ if page.IsHome }}
Heading OK.
{{ end }}
-- layouts/_default/_markup/render-image.html --
{{ if page.IsHome }}
Image OK.
{{ end }}
-- layouts/_default/_markup/render-link.html --
{{ if page.IsHome }}
Link OK.
{{ end }}
-- layouts/_default/myview.html
{{ if page.IsHome }}
Render OK.
{{ end }}
-- layouts/_default/_markup/render-codeblock.html --
{{ if page.IsHome }}
Codeblock OK.
{{ end }}
-- layouts/_default/single.html --
Single.
-- layouts/index.html --
{{ if eq page . }}Page OK.{{ end }}
{{ $r := "{{ if page.IsHome }}ExecuteAsTemplate OK.{{ end }}" | resources.FromString "foo.html" |  resources.ExecuteAsTemplate "foo.html" . }}
{{ $r.Content }}
{{ .RenderString "{{< renderstring.inline >}}{{ if page.IsHome }}RenderString OK.{{ end }}{{< /renderstring.inline >}}}}"}}
{{ .Render "myview" }}
{{ .Content }}
partial: {{ partials.Include "foo.html" . }}
{{ $pag := (.Paginate site.RegularPages) }}
PageNumber: {{ $pag.PageNumber }}/{{ $pag.TotalPages }}|
{{ $p2 := site.GetPage "p2" }}
{{ $p2_1 := index $p2.Resources 0 }}
Bundled page: {{ $p2_1.Content }}
-- layouts/alias.html --
{{ if eq page .Page }}Alias OK.{{ else }}Failed.{{ end }}
-- layouts/404.html --
{{ if eq page . }}404 Page OK.{{ else }}Failed.{{ end }}
-- layouts/partials/foo.html --
{{ if page.IsHome }}Partial OK.{{ else }}Failed.{{ end }}
-- layouts/shortcodes/outer.html --
{{ .Inner }}
-- layouts/shortcodes/inner.html --
{{ if page.IsHome }}Shortcode Inner OK.{{ else }}Failed.{{ end }}
-- layouts/shortcodes/shortcode.html --
{{ if page.IsHome }}Shortcode {{ .Get 0 }} OK.{{ else }}Failed.{{ end }}
-- layouts/sitemap.xml --
{{ if eq page . }}Sitemap OK.{{ else }}Failed.{{ end }}
-- layouts/robots.txt --
{{ if eq page . }}Robots OK.{{ else }}Failed.{{ end }}
-- layouts/sitemapindex.xml --
{{ with page }}SitemapIndex OK: {{ .Kind }}{{ else }}Failed.{{ end }}

  `

	for _, multilingual := range []bool{false, true} {
		t.Run(fmt.Sprintf("multilingual-%t", multilingual), func(t *testing.T) {
			// Fenced code blocks.
			files := strings.ReplaceAll(filesTemplate, "$$$", "```")

			if multilingual {
				files = strings.ReplaceAll(files, "LANG_CONFIG", `
[languages]
[languages.en]
weight = 1
[languages.no]
weight = 2
`)
			} else {
				files = strings.ReplaceAll(files, "LANG_CONFIG", "")
			}

			b := hugolib.NewIntegrationTestBuilder(
				hugolib.IntegrationTestConfig{
					T:           t,
					TxtarString: files,
				},
			).Build()

			b.AssertFileContent("public/index.html", `
Heading OK.
Image OK.
Link OK.
Codeblock OK.
Page OK.
Partial OK.
Shortcode Angled Brackets OK.
Shortcode Percentage OK.
Shortcode Inner OK.
Shortcode Inline OK.
ExecuteAsTemplate OK.
RenderString OK.
Render OK.
Shortcode in bundled page OK.
	`)

			b.AssertFileContent("public/404.html", `404 Page OK.`)
			b.AssertFileContent("public/robots.txt", `Robots OK.`)
			b.AssertFileContent("public/homealias/index.html", `Alias OK.`)
			b.AssertFileContent("public/page/1/index.html", `Alias OK.`)
			b.AssertFileContent("public/page/2/index.html", `Page OK.`)
			if multilingual {
				b.AssertFileContent("public/sitemap.xml", `SitemapIndex OK: sitemapindex`)
			} else {
				b.AssertFileContent("public/sitemap.xml", `Sitemap OK.`)
			}
		})
	}
}

// Issue 10791.
func TestPageTableOfContentsInShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
disableKinds = ["taxonomy", "term"]
-- content/p1.md --
---
title: "P1"
---
{{< toc >}}

# Heading 1
-- layouts/shortcodes/toc.html --
{{ page.TableOfContents }}
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "<nav id=\"TableOfContents\"></nav> \n<h1 id=\"heading-1\">Heading 1</h1>")
}

func TestFromStringRunning(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
-- layouts/index.html --
{{ with resources.FromString "foo" "{{ seq 3 }}" }}
{{ with resources.ExecuteAsTemplate "bar" $ . }}
	{{ .Content | safeHTML }}
{{ end }}
{{ end }}
  `

	b := hugolib.TestRunning(t, files)

	b.AssertFileContent("public/index.html", "1\n2\n3")
}
