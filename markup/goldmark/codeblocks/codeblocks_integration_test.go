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

package codeblocks_test

import (
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

func TestCodeblocks(t *testing.T) {
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
-- layouts/_default/_markup/render-codeblock-goat.html --
{{ $diagram := diagrams.Goat .Inner }}
Goat SVG:{{ substr $diagram.Wrapped 0 100 | safeHTML }}  }}|
Goat Attribute: {{ .Attributes.width}}|
-- layouts/_default/_markup/render-codeblock-go.html --
Go Code: {{ .Inner | safeHTML }}|
Go Language: {{ .Type }}|
-- layouts/_default/single.html --
{{ .Content }}
-- content/p1.md --
---
title: "p1"
---

## Ascii Diagram

§§§goat { width="600" }
--->
§§§

## Go Code

§§§go
fmt.Println("Hello, World!");
§§§

## Golang Code

§§§golang
fmt.Println("Hello, Golang!");
§§§

## Bash Code

§§§bash { linenos=inline,hl_lines=[2,"5-6"],linenostart=32 class=blue }
echo "l1";
echo "l2";
echo "l3";
echo "l4";
echo "l5";
echo "l6";
echo "l7";
echo "l8";
§§§
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
Goat SVG:<svg class='diagram'
Goat Attribute: 600|

Go Language: go|
Go Code: fmt.Println("Hello, World!");

Go Code: fmt.Println("Hello, Golang!");
Go Language: golang|


	`,
		"Goat SVG:<svg class='diagram' xmlns='http://www.w3.org/2000/svg' version='1.1' height='25' width='40'",
		"Goat Attribute: 600|",
		"<h2 id=\"go-code\">Go Code</h2>\nGo Code: fmt.Println(\"Hello, World!\");\n|\nGo Language: go|",
		"<h2 id=\"golang-code\">Golang Code</h2>\nGo Code: fmt.Println(\"Hello, Golang!\");\n|\nGo Language: golang|",
		"<h2 id=\"bash-code\">Bash Code</h2>\n<div class=\"highlight blue\"><pre tabindex=\"0\" class=\"chroma\"><code class=\"language-bash\" data-lang=\"bash\"><span class=\"line\"><span class=\"ln\">32</span><span class=\"cl\"><span class=\"nb\">echo</span> <span class=\"s2\">&#34;l1&#34;</span><span class=\"p\">;</span>\n</span></span><span class=\"line hl\"><span class=\"ln\">33</span>",
	)
}

func TestHighlightCodeblock(t *testing.T) {
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
-- layouts/_default/_markup/render-codeblock.html --
{{ $result := transform.HighlightCodeBlock . }}
Inner: |{{ $result.Inner | safeHTML }}|
Wrapped: |{{ $result.Wrapped | safeHTML }}|
-- layouts/_default/single.html --
{{ .Content }}
-- content/p1.md --
---
title: "p1"
---

## Go Code

§§§go
fmt.Println("Hello, World!");
§§§

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Inner: |<span class=\"line\"><span class=\"cl\"><span class=\"nx\">fmt</span><span class=\"p\">.</span><span class=\"nf\">Println</span><span class=\"p\">(</span><span class=\"s\">&#34;Hello, World!&#34;</span><span class=\"p\">);</span></span></span>|",
		"Wrapped: |<div class=\"highlight\"><pre tabindex=\"0\" class=\"chroma\"><code class=\"language-go\" data-lang=\"go\"><span class=\"line\"><span class=\"cl\"><span class=\"nx\">fmt</span><span class=\"p\">.</span><span class=\"nf\">Println</span><span class=\"p\">(</span><span class=\"s\">&#34;Hello, World!&#34;</span><span class=\"p\">);</span></span></span></code></pre></div>|",
	)
}

func TestCodeblocksBugs(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- layouts/_default/_markup/render-codeblock.html --
{{ .Position | safeHTML }}
-- layouts/_default/single.html --
{{ .Content }}
-- content/p1.md --
---
title: "p1"
---

## Issue 9627

§§§text
{{</* foo */>}}
§§§

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
# Issue 9627: For the Position in code blocks we try to match the .Inner with the original source. This isn't always possible.
p1.md:0:0
	`,
	)
}

func TestCodeChomp(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---

§§§bash
echo "p1";
§§§
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-codeblock.html --
|{{ .Inner | safeHTML }}|

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "|echo \"p1\";|")
}

func TestCodePosition(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---

##   Code

§§§
echo "p1";
§§§
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-codeblock.html --
Position: {{ .Position | safeHTML }}


`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", filepath.FromSlash("Position: \"/content/p1.md:7:1\""))
}

// Issue 10118
func TestAttributes(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---

## Issue 10118

§§§ {foo="bar"}
Hello, World!
§§§

-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-codeblock.html --
Attributes: {{ .Attributes }}|Type: {{ .Type }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "<h2 id=\"issue-10118\">Issue 10118</h2>\nAttributes: map[foo:bar]|Type: |")
}

// Issue 9571
func TestAttributesChroma(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---

##   Code

§§§LANGUAGE {style=monokai}
echo "p1";
§§§
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-codeblock.html --
Attributes: {{ .Attributes }}|Options: {{ .Options }}|


`
	testLanguage := func(language, expect string) {
		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           t,
				TxtarString: strings.ReplaceAll(files, "LANGUAGE", language),
			},
		).Build()

		b.AssertFileContent("public/p1/index.html", expect)
	}

	testLanguage("bash", "Attributes: map[]|Options: map[style:monokai]|")
	testLanguage("hugo", "Attributes: map[style:monokai]|Options: map[]|")
}

func TestPanics(t *testing.T) {
	files := `
-- config.toml --
[markup]
[markup.goldmark]
[markup.goldmark.parser]
autoHeadingID = true
autoHeadingIDType = "github"
[markup.goldmark.parser.attribute]
block = true
title = true
-- content/p1.md --
---
title: "p1"
---

BLOCK

Common

-- layouts/_default/single.html --
{{ .Content }}


`

	for _, test := range []struct {
		name     string
		markdown string
	}{
		{"issue-9819", "asdf\n: {#myid}"},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			b := hugolib.NewIntegrationTestBuilder(
				hugolib.IntegrationTestConfig{
					T:           t,
					TxtarString: strings.ReplaceAll(files, "BLOCK", test.markdown),
				},
			).Build()

			b.AssertFileContent("public/p1/index.html", "Common")
		})
	}
}

// Issue 10835
func TestAttributesValidation(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
-- content/p1.md --
---
title: "p1"
---

## Issue 10835

§§§bash { color=red dimensions=300x200 }
Hello, World!
§§§

-- layouts/index.html --
-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-codeblock.html --
Attributes: {{ .Attributes }}|Type: {{ .Type }}|
`

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, "p1.md:7:9\": failed to parse Markdown attributes; you may need to quote the values")
}
