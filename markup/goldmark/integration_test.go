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
## Heading {class="a" onclick="alert('heading')" linenos="inline"}

> Blockquote
{class="b" ondblclick="alert('blockquote')" LINENOS="inline"}

~~~bash {id="c" onmouseover="alert('code fence')"}
foo
~~~
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   false,
		},
	).Build()

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
>{{ .Text | safeHTML }}</h{{ .Level }}>
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   false,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", `
		<h2 data-foo="bar" id="heading">Heading</h2>
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
  {{ .Text | safeHTML }}
  <a class="anchor" href="#{{ .Anchor | safeURL }}">#</a>
</h{{ .Level }}>
-- layouts/_default/_markup/render-link.html --
<a href="{{ .Destination | safeURL }}"{{ with .Title}} title="{{ . }}"{{ end }}>{{ .Text | safeHTML }}</a>

`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   false,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html",
		"<h2 id=\"hello-testhttpsexamplecom\">\n  Hello <a href=\"https://example.com\">Test</a>\n\n  <a class=\"anchor\" href=\"#hello-testhttpsexamplecom\">#</a>\n</h2>",
	)
}

func BenchmarkRenderHooks(b *testing.B) {
	files := `
-- config.toml --
-- layouts/_default/_markup/render-heading.html --
<h{{ .Level }} id="{{ .Anchor | safeURL }}">
	{{ .Text | safeHTML }}
	<a class="anchor" href="#{{ .Anchor | safeURL }}">#</a>
</h{{ .Level }}>
-- layouts/_default/_markup/render-link.html --
<a href="{{ .Destination | safeURL }}"{{ with .Title}} title="{{ . }}"{{ end }}>{{ .Text | safeHTML }}</a>
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

FENCEbash
#!/bin/bash
# Usage: Hello World Bash Shell Script Using Variables
# Author: Vivek Gite
# -------------------------------------------------
 
# Define bash shell variable called var 
# Avoid spaces around the assignment operator (=)
var="Hello World"
 
# print it 
echo "$var"
 
# Another way of printing it
printf "%s\n" "$var"
FENCE
`

	content = strings.ReplaceAll(content, "FENCE", "```")

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
