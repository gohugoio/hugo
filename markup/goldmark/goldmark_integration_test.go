// Copyright 2021 The Hugo Authors. All rights reserved.
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

package goldmark_test

import (
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue 9463
func TestAttributeExclusion(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.goldmark.renderer]
	unsafe = false
[markup.goldmark.parser.attribute]
	block = true
	title = true
-- content/p1.md --
---
title: "p1"
---
## Heading {class="a" onclick="alert('heading')"}

> Blockquote
{class="b" ondblclick="alert('blockquote')"}

~~~bash {id="c" onmouseover="alert('code fence')" LINENOS=true}
foo
~~~
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
		<h2 class="a" id="heading">
		<blockquote class="b">
		<div class="highlight" id="c">
	`)
}

// Issue 9511
func TestAttributeExclusionWithRenderHook(t *testing.T) {
	t.Parallel()

	files := `
-- content/p1.md --
---
title: "p1"
---
## Heading {onclick="alert('renderhook')" data-foo="bar"}
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-heading.html --
<h{{ .Level }}
  {{- range $k, $v := .Attributes -}}
    {{- printf " %s=%q" $k $v | safeHTMLAttr -}}
  {{- end -}}
>{{ .Text }}</h{{ .Level }}>
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
		<h2 data-foo="bar" id="heading">Heading</h2>
	`)
}

func TestAttributesDefaultRenderer(t *testing.T) {
	t.Parallel()

	files := `
-- content/p1.md --
---
title: "p1"
---
## Heading Attribute Which Needs Escaping { class="a < b" }
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
class="a &lt; b"
	`)
}

// Issue 9558.
func TestAttributesHookNoEscape(t *testing.T) {
	t.Parallel()

	files := `
-- content/p1.md --
---
title: "p1"
---
## Heading Attribute Which Needs Escaping { class="Smith & Wesson" }
-- layouts/_default/_markup/render-heading.html --
plain: |{{- range $k, $v := .Attributes -}}{{ $k }}: {{ $v }}|{{ end }}|
safeHTML: |{{- range $k, $v := .Attributes -}}{{ $k }}: {{ $v | safeHTML }}|{{ end }}|
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
plain: |class: Smith &amp; Wesson|id: heading-attribute-which-needs-escaping|
safeHTML: |class: Smith & Wesson|id: heading-attribute-which-needs-escaping|
	`)
}

// Issue 9504
func TestLinkInTitle(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---
## Hello [Test](https://example.com)
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-heading.html --
<h{{ .Level }} id="{{ .Anchor | safeURL }}">
  {{ .Text }}
  <a class="anchor" href="#{{ .Anchor | safeURL }}">#</a>
</h{{ .Level }}>
-- layouts/_default/_markup/render-link.html --
<a href="{{ .Destination | safeURL }}"{{ with .Title}} title="{{ . }}"{{ end }}>{{ .Text }}</a>

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"<h2 id=\"hello-testhttpsexamplecom\">\n  Hello <a href=\"https://example.com\">Test</a>\n\n  <a class=\"anchor\" href=\"#hello-testhttpsexamplecom\">#</a>\n</h2>",
	)
}

func TestHighlight(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup]
[markup.highlight]
anchorLineNos = false
codeFences = true
guessSyntax = false
hl_Lines = ''
lineAnchors = ''
lineNoStart = 1
lineNos = false
lineNumbersInTable = true
noClasses = false
style = 'monokai'
tabWidth = 4
-- layouts/_default/single.html --
{{ .Content }}
-- content/p1.md --
---
title: "p1"
---

## Code Fences

§§§bash
LINE1
§§§

## Code Fences No Lexer

§§§moo
LINE1
§§§

## Code Fences Simple Attributes

§§A§bash { .myclass id="myid" }
LINE1
§§A§

## Code Fences Line Numbers

§§§bash {linenos=table,hl_lines=[8,"15-17"],linenostart=199}
LINE1
LINE2
LINE3
LINE4
LINE5
LINE6
LINE7
LINE8
§§§




`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"<div class=\"highlight\"><pre tabindex=\"0\" class=\"chroma\"><code class=\"language-bash\" data-lang=\"bash\"><span class=\"line\"><span class=\"cl\">LINE1\n</span></span></code></pre></div>",
		"Code Fences No Lexer</h2>\n<pre tabindex=\"0\"><code class=\"language-moo\" data-lang=\"moo\">LINE1\n</code></pre>",
		"lnt",
	)
}

func BenchmarkRenderHooks(b *testing.B) {
	files := `
-- config.toml --
-- layouts/_default/_markup/render-heading.html --
<h{{ .Level }} id="{{ .Anchor | safeURL }}">
	{{ .Text }}
	<a class="anchor" href="#{{ .Anchor | safeURL }}">#</a>
</h{{ .Level }}>
-- layouts/_default/_markup/render-link.html --
<a href="{{ .Destination | safeURL }}"{{ with .Title}} title="{{ . }}"{{ end }}>{{ .Text }}</a>
-- layouts/_default/single.html --
{{ .Content }}
`

	content := `

## Hello1 [Test](https://example.com)

A.

## Hello2 [Test](https://example.com)

B.

## Hello3 [Test](https://example.com)

C.

## Hello4 [Test](https://example.com)

D.

[Test](https://example.com)

## Hello5


`

	for i := 1; i < 100; i++ {
		files += fmt.Sprintf("\n-- content/posts/p%d.md --\n"+content, i+1)
	}

	cfg := hugolib.IntegrationTestConfig{
		T:           b,
		TxtarString: files,
	}
	builders := make([]*hugolib.IntegrationTestBuilder, b.N)

	for i := range builders {
		builders[i] = hugolib.NewIntegrationTestBuilder(cfg)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builders[i].Build()
	}
}

func BenchmarkCodeblocks(b *testing.B) {
	filesTemplate := `
-- config.toml --
[markup]
  [markup.highlight]
    anchorLineNos = false
    codeFences = true
    guessSyntax = false
    hl_Lines = ''
    lineAnchors = ''
    lineNoStart = 1
    lineNos = false
    lineNumbersInTable = true
    noClasses = true
    style = 'monokai'
    tabWidth = 4
-- layouts/_default/single.html --
{{ .Content }}
`

	content := `

FENCEgo
package main
import "fmt"
func main() {
    fmt.Println("hello world")
}
FENCE

FENCEunknownlexer
hello
FENCE
`

	content = strings.ReplaceAll(content, "FENCE", "```")

	for i := 1; i < 100; i++ {
		filesTemplate += fmt.Sprintf("\n-- content/posts/p%d.md --\n"+content, i+1)
	}

	runBenchmark := func(files string, b *testing.B) {
		cfg := hugolib.IntegrationTestConfig{
			T:           b,
			TxtarString: files,
		}
		builders := make([]*hugolib.IntegrationTestBuilder, b.N)

		for i := range builders {
			builders[i] = hugolib.NewIntegrationTestBuilder(cfg)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			builders[i].Build()
		}
	}

	b.Run("Default", func(b *testing.B) {
		runBenchmark(filesTemplate, b)
	})

	b.Run("Hook no higlight", func(b *testing.B) {
		files := filesTemplate + `
-- layouts/_default/_markup/render-codeblock.html --
{{ .Inner }}
`

		runBenchmark(files, b)
	})
}

// Iisse #8959
func TestHookInfiniteRecursion(t *testing.T) {
	t.Parallel()

	for _, renderFunc := range []string{"markdownify", ".Page.RenderString"} {
		t.Run(renderFunc, func(t *testing.T) {
			files := `
-- config.toml --
-- layouts/_default/_markup/render-link.html --
<a href="{{ .Destination | safeURL }}">{{ .Text | RENDERFUNC }}</a>
-- layouts/_default/single.html --
{{ .Content }}
-- content/p1.md --
---
title: "p1"
---

https://example.org

a@b.com


			`

			files = strings.ReplaceAll(files, "RENDERFUNC", renderFunc)

			b, err := hugolib.NewIntegrationTestBuilder(
				hugolib.IntegrationTestConfig{
					T:           t,
					TxtarString: files,
				},
			).BuildE()

			b.Assert(err, qt.IsNotNil)
			b.Assert(err.Error(), qt.Contains, "text is already rendered, repeating it may cause infinite recursion")
		})
	}
}

// Issue 9594
func TestQuotesInImgAltAttr(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.goldmark.extensions]
  typographer = false
-- content/p1.md --
---
title: "p1"
---
!["a"](b.jpg)
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
		<img src="b.jpg" alt="&quot;a&quot;">
	`)
}

func TestLinkifyProtocol(t *testing.T) {
	t.Parallel()

	runTest := func(protocol string, withHook bool) *hugolib.IntegrationTestBuilder {
		files := `
-- config.toml --
[markup.goldmark]
[markup.goldmark.extensions]
linkify = true
linkifyProtocol = "PROTOCOL"
-- content/p1.md --
---
title: "p1"
---
Link no procol: www.example.org
Link http procol: http://www.example.org
Link https procol: https://www.example.org

-- layouts/_default/single.html --
{{ .Content }}
`
		files = strings.ReplaceAll(files, "PROTOCOL", protocol)

		if withHook {
			files += `-- layouts/_default/_markup/render-link.html --
<a href="{{ .Destination | safeURL }}">{{ .Text }}</a>`
		}

		return hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()
	}

	for _, withHook := range []bool{false, true} {

		b := runTest("https", withHook)

		b.AssertFileContent("public/p1/index.html",
			"Link no procol: <a href=\"https://www.example.org\">www.example.org</a>",
			"Link http procol: <a href=\"http://www.example.org\">http://www.example.org</a>",
			"Link https procol: <a href=\"https://www.example.org\">https://www.example.org</a></p>",
		)

		b = runTest("http", withHook)

		b.AssertFileContent("public/p1/index.html",
			"Link no procol: <a href=\"http://www.example.org\">www.example.org</a>",
			"Link http procol: <a href=\"http://www.example.org\">http://www.example.org</a>",
			"Link https procol: <a href=\"https://www.example.org\">https://www.example.org</a></p>",
		)

		b = runTest("gopher", withHook)

		b.AssertFileContent("public/p1/index.html",
			"Link no procol: <a href=\"gopher://www.example.org\">www.example.org</a>",
			"Link http procol: <a href=\"http://www.example.org\">http://www.example.org</a>",
			"Link https procol: <a href=\"https://www.example.org\">https://www.example.org</a></p>",
		)

	}
}

func TestGoldmarkBugs(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.goldmark.renderer]
unsafe = true
-- content/p1.md --
---
title: "p1"
---

## Issue 9650

a <!-- b --> c

## Issue 9658

- This is a list item <!-- Comment: an innocent-looking comment -->


-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContentExact("public/p1/index.html",
		// Issue 9650
		"<p>a <!-- b --> c</p>",
		// Issue 9658 (crash)
		"<li>This is a list item <!-- Comment: an innocent-looking comment --></li>",
	)
}

// Issue #7332
// Issue #11587
func TestGoldmarkEmojiExtension(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
enableEmoji = true
-- content/p1.md --
---
title: "p1"
---
~~~text
:x:
~~~

{{% include "/p2" %}}

{{< sc1 >}}:smiley:{{< /sc1 >}}

{{< sc2 >}}:+1:{{< /sc2 >}}

{{% sc3 %}}:-1:{{% /sc3 %}}

-- content/p2.md --
---
title: "p2"
---
:heavy_check_mark:
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/sc1.html --
sc1_begin|{{ .Inner }}|sc1_end
-- layouts/shortcodes/sc2.html --
sc2_begin|{{ .Inner | .Page.RenderString }}|sc2_end
-- layouts/shortcodes/sc3.html --
sc3_begin|{{ .Inner }}|sc3_end
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContentExact("public/p1/index.html",
		// Issue #7332
		"<span>:x:\n</span>",
		// Issue #11587
		"<p>&#x2714;&#xfe0f;</p>",
		// Should not be converted to emoji
		"sc1_begin|:smiley:|sc1_end",
		// Should be converted to emoji
		"sc2_begin|&#x1f44d;|sc2_end",
		// Should be converted to emoji
		"sc3_begin|&#x1f44e;|sc3_end",
	)
}

func TestEmojiDisabled(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
enableEmoji = false
-- content/p1.md --
---
title: "p1"
---
:x:
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContentExact("public/p1/index.html", "<p>:x:</p>")
}

func TestEmojiDefaultConfig(t *testing.T) {
	t.Parallel()

	files := `
-- content/p1.md --
---
title: "p1"
---
:x:
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContentExact("public/p1/index.html", "<p>:x:</p>")
}

// Issue #5748
func TestGoldmarkTemplateDelims(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[minify]
  minifyOutput = true
[minify.tdewolff.html]
  templateDelims = ["<?php","?>"]
-- layouts/index.html --
<div class="foo">
{{ safeHTML "<?php" }}
echo "hello";
{{ safeHTML "?>" }}
</div>
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "<div class=foo><?php\necho \"hello\";\n?>\n</div>")
}

// Issue #10894
func TestPassthroughInlineFences(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
inline = [['$', '$'], ['\(', '\)']]
-- content/p1.md --
---
title: "p1"
---
## LaTeX test

Inline equation that would be mangled by default parser: $a^*=x-b^*$

-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html", `
		$a^*=x-b^*$
	`)
}

func TestPassthroughBlockFences(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
block = [['$$', '$$']]
-- content/p1.md --
---
title: "p1"
---
## LaTeX test

Block equation that would be mangled by default parser:

$$a^*=x-b^*$$

-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html", `
		$$a^*=x-b^*$$
	`)
}

func TestPassthroughWithAlternativeFences(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
inline = [['(((', ')))']]
block = [['%!%', '%!%']]
-- content/p1.md --
---
title: "p1"
---
## LaTeX test

Inline equation that would be mangled by default parser: (((a^*=x-b^*)))
Inline equation that should be mangled by default parser: $a^*=x-b^*$

Block equation that would be mangled by default parser:

%!%
a^*=x-b^*
%!%

-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html", `
		(((a^*=x-b^*)))
	`)
	b.AssertFileContent("public/p1/index.html", `
		$a^<em>=x-b^</em>$
	`)
	b.AssertFileContent("public/p1/index.html", `
%!%
a^*=x-b^*
%!%
	`)
}

func TestExtrasExtension(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
[markup.goldmark.extensions]
strikethrough = false
[markup.goldmark.extensions.extras.delete]
enable = false
[markup.goldmark.extensions.extras.insert]
enable = false
[markup.goldmark.extensions.extras.mark]
enable = false
[markup.goldmark.extensions.extras.subscript]
enable = false
[markup.goldmark.extensions.extras.superscript]
enable = false
-- layouts/index.html --
{{ .Content }}
-- content/_index.md --
---
title: home
---
~~delete~~

++insert++

==mark==

H~2~0

1^st^
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"<p>~~delete~~</p>",
		"<p>++insert++</p>",
		"<p>==mark==</p>",
		"<p>H~2~0</p>",
		"<p>1^st^</p>",
	)

	files = strings.ReplaceAll(files, "enable = false", "enable = true")

	b = hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"<p><del>delete</del></p>",
		"<p><ins>insert</ins></p>",
		"<p><mark>mark</mark></p>",
		"<p>H<sub>2</sub>0</p>",
		"<p>1<sup>st</sup></p>",
	)
}
