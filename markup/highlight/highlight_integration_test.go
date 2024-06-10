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

package highlight_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestHighlightInline(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup]
[markup.highlight]
codeFences = true
noClasses = false
-- content/p1.md --
---
title: "p1"
---

## Inline in Shortcode

Inline:{{< highlight emacs "hl_inline=true" >}}(message "this highlight shortcode"){{< /highlight >}}:End.
Inline Unknown:{{< highlight foo "hl_inline=true" >}}(message "this highlight shortcode"){{< /highlight >}}:End.

## Inline in code block

Not sure if this makes sense, but add a test for it:

§§§bash {hl_inline=true}
(message "highlight me")
§§§

## HighlightCodeBlock in hook

§§§html
(message "highlight me 2")
§§§

## Unknown lexer

§§§foo {hl_inline=true}
(message "highlight me 3")
§§§


-- layouts/_default/_markup/render-codeblock-html.html --
{{ $opts := dict "hl_inline" true }}
{{ $result := transform.HighlightCodeBlock . $opts }}
HighlightCodeBlock: Wrapped:{{ $result.Wrapped  }}|Inner:{{ $result.Inner }}
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Inline:<code class=\"code-inline language-emacs\"><span class=\"p\">(</span><span class=\"nf\">message</span> <span class=\"s\">&#34;this highlight shortcode&#34;</span><span class=\"p\">)</span></code>:End.",
		"Inline Unknown:<code class=\"code-inline language-foo\">(message &#34;this highlight shortcode&#34;)</code>:End.",
		"Not sure if this makes sense, but add a test for it:</p>\n<code class=\"code-inline language-bash\"><span class=\"o\">(</span>message <span class=\"s2\">&#34;highlight me&#34;</span><span class=\"o\">)</span>\n</code>",
		"HighlightCodeBlock: Wrapped:<code class=\"code-inline language-html\">(message &#34;highlight me 2&#34;)</code>|Inner:<code class=\"code-inline language-html\">(message &#34;highlight me 2&#34;)</code>",
		"<code class=\"code-inline language-foo\">(message &#34;highlight me 3&#34;)\n</code>",
	)
}

// Issue #11311
func TestIssue11311(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.highlight]
noClasses = false
-- content/_index.md --
---
title: home
---
§§§go
xəx := 0
§§§
-- layouts/index.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
		<span class="nx">xəx</span>
	`)
}
