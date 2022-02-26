// Copyright 2022 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"

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
{{ $diagram := diagrams.Goat .Code }}
Goat SVG:{{ substr $diagram.SVG 0 100 | safeHTML }}  }}|
Goat Attribute: {{ .Attributes.width}}|
-- layouts/_default/_markup/render-codeblock-go.html --
Go Code: {{ .Code | safeHTML }}|
Go Language: {{ .Lang }}|
-- layouts/_default/single.html --
{{ .Content }}
-- content/p1.md --
---
title: "p1"
---

## Ascii Diagram

CODE_FENCEgoat { width="600" }
--->
CODE_FENCE

## Go Code

CODE_FENCEgo
fmt.Println("Hello, World!");
CODE_FENCE

## Golang Code

CODE_FENCEgolang
fmt.Println("Hello, Golang!");
CODE_FENCE

## Bash Code

CODE_FENCEbash { linenos=inline,hl_lines=[2,"5-6"],linenostart=32 class=blue }
echo "l1";
echo "l2";
echo "l3";
echo "l4";
echo "l5";
echo "l6";
echo "l7";
echo "l8";
CODE_FENCE
`

	files = strings.ReplaceAll(files, "CODE_FENCE", "```")

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   false,
		},
	).Build()

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
|{{ .Code | safeHTML }}|

`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   false,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", "|echo \"p1\";|")
}
