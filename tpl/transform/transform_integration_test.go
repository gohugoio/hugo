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

package transform_test

import (
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

// Issue #11698
func TestMarkdownifyIssue11698(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ['home','section','rss','sitemap','taxonomy','term']
[markup.goldmark.parser.attribute]
title = true
block = true
-- layouts/_default/single.html --
_{{ markdownify .RawContent }}_
-- content/p1.md --
---
title: p1
---
foo bar
-- content/p2.md --
---
title: p2
---
foo

**bar**
-- content/p3.md --
---
title: p3
---
## foo

bar
-- content/p4.md --
---
title: p4
---
foo
{#bar}
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "_foo bar_")
	b.AssertFileContent("public/p2/index.html", "_<p>foo</p>\n<p><strong>bar</strong></p>\n_")
	b.AssertFileContent("public/p3/index.html", "_<h2 id=\"foo\">foo</h2>\n<p>bar</p>\n_")
	b.AssertFileContent("public/p4/index.html", "_<p id=\"bar\">foo</p>\n_")
}

func TestXMLEscape(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ['section','sitemap','taxonomy','term']
-- content/p1.md --
---
title: p1
---
a **b** ` + "\v" + ` c
<!--more-->
  `
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.xml", `
	<description>&lt;p&gt;a &lt;strong&gt;b&lt;/strong&gt;  c&lt;/p&gt;</description>
	`)
}

// Issue #9642
func TestHighlightError(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ highlight "a" "b" 0 }}
  `
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)

	_, err := b.BuildE()
	b.Assert(err.Error(), qt.Contains, "error calling highlight: invalid Highlight option: 0")
}

// Issue #11884
func TestUnmarshalCSVLazyDecoding(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/pets.csv --
name,description,age
Spot,a nice dog,3
Rover,"a big dog",5
Felix,a "malicious" cat,7
Bella,"an "evil" cat",9
Scar,"a "dead cat",11
-- layouts/index.html --
{{ $opts := dict "lazyQuotes" true }}
{{ $data := resources.Get "pets.csv" | transform.Unmarshal $opts }}
{{ printf "%v" $data | safeHTML }}
  `
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
[[name description age] [Spot a nice dog 3] [Rover a big dog 5] [Felix a "malicious" cat 7] [Bella an "evil" cat 9] [Scar a "dead cat 11]]
	`)
}

func TestToMath(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ transform.ToMath "c = \\pm\\sqrt{a^2 + b^2}" }}
  `
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
<span class="katex"><math
	`)
}

func TestToMathError(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{  transform.ToMath "c = \\foo{a^2 + b^2}" }}
  `
		b, err := hugolib.TestE(t, files, hugolib.TestOptWarn())

		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, "KaTeX parse error: Undefined control sequence: \\foo")
	})

	t.Run("Disable ThrowOnError", func(t *testing.T) {
		files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ $opts := dict "throwOnError" false }}
{{  transform.ToMath "c = \\foo{a^2 + b^2}" $opts }}
  `
		b, err := hugolib.TestE(t, files, hugolib.TestOptWarn())

		b.Assert(err, qt.IsNil)
		b.AssertFileContent("public/index.html", `#cc0000`) // Error color
	})

	t.Run("Handle in template", func(t *testing.T) {
		files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ with transform.ToMath "c = \\foo{a^2 + b^2}" }}
	{{ with .Err }}
	 	{{ warnf "error: %s" . }}
	{{ else }}
		{{ . }}
	{{ end }}
{{ end }}
  `
		b, err := hugolib.TestE(t, files, hugolib.TestOptWarn())

		b.Assert(err, qt.IsNil)
		b.AssertLogContains("WARN  error: KaTeX parse error: Undefined control sequence: \\foo")
	})
}

func TestToMathBigAndManyExpressions(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
disableKinds = ['rss','section','sitemap','taxonomy','term']
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
block  = [['\[', '\]'], ['$$', '$$']]
inline = [['\(', '\)'], ['$', '$']]
-- content/p1.md --
P1_CONTENT
-- layouts/index.html --
Home.
-- layouts/_default/single.html --
Content: {{ .Content }}|
-- layouts/_default/_markup/render-passthrough.html --
{{ $opts := dict "throwOnError" false "displayMode" true }}
{{ transform.ToMath .Inner $opts }}
  `

	t.Run("Very large file with many complex KaTeX expressions", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "P1_CONTENT", "sourcefilename: testdata/large-katex.md")
		b := hugolib.Test(t, files)
		b.AssertFileContent("public/p1/index.html", `
		<span class="katex"><math
			`)
	})

	t.Run("Large and complex expression", func(t *testing.T) {
		// This is pulled from the file above, which times out for some reason.
		largeAndComplexeExpressions := `\begin{align*} \frac{\pi^2}{6}&=\frac{4}{3}\frac{(\arcsin 1)^2}{2}\\ &=\frac{4}{3}\int_0^1\frac{\arcsin x}{\sqrt{1-x^2}}\,dx\\ &=\frac{4}{3}\int_0^1\frac{x+\sum_{n=1}^{\infty}\frac{(2n-1)!!}{(2n)!!}\frac{x^{2n+1}}{2n+1}}{\sqrt{1-x^2}}\,dx\\ &=\frac{4}{3}\int_0^1\frac{x}{\sqrt{1-x^2}}\,dx +\frac{4}{3}\sum_{n=1}^{\infty}\frac{(2n-1)!!}{(2n)!!(2n+1)}\int_0^1x^{2n}\frac{x}{\sqrt{1-x^2}}\,dx\\ &=\frac{4}{3}+\frac{4}{3}\sum_{n=1}^{\infty}\frac{(2n-1)!!}{(2n)!!(2n+1)}\left[\frac{(2n)!!}{(2n+1)!!}\right]\\ &=\frac{4}{3}\sum_{n=0}^{\infty}\frac{1}{(2n+1)^2}\\ &=\frac{4}{3}\left(\sum_{n=1}^{\infty}\frac{1}{n^2}-\frac{1}{4}\sum_{n=1}^{\infty}\frac{1}{n^2}\right)\\ &=\sum_{n=1}^{\infty}\frac{1}{n^2} \end{align*}`
		files := strings.ReplaceAll(filesTemplate, "P1_CONTENT", fmt.Sprintf(`---
title: p1
---

$$%s$$
	`, largeAndComplexeExpressions))

		b := hugolib.Test(t, files)
		b.AssertFileContent("public/p1/index.html", `
		<span class="katex"><math
			`)
	})
}

func TestToMathMacros(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ $macros := dict 
    "\\addBar" "\\bar{#1}"
	"\\bold" "\\mathbf{#1}"
}}
{{ $opts := dict "macros" $macros }}
{{ transform.ToMath "\\addBar{y} + \\bold{H}" $opts }}
  `
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
<mi>y</mi>
	`)
}
