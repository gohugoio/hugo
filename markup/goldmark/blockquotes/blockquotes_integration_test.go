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

package blockquotes_test

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestBlockquoteHook(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[markup]
  [markup.goldmark]
    [markup.goldmark.parser]
      [markup.goldmark.parser.attribute]
        block = true
        title = true
-- layouts/_default/_markup/render-blockquote.html --
Blockquote: |{{ .Text | safeHTML }}|{{ .Type }}|
-- layouts/_default/_markup/render-blockquote-alert.html --
{{ $text := .Text | safeHTML }}
Blockquote Alert: |{{ $text }}|{{ .Type }}|
Blockquote Alert Attributes: |{{ $text }}|{{ .Attributes }}|
Blockquote Alert Page: |{{ $text }}|{{ .Page.Title }}|{{ .PageInner.Title }}|
{{ if .Attributes.showpos }}
Blockquote Alert Position: |{{ $text }}|{{ .Position | safeHTML }}|
{{ end }}
-- layouts/_default/single.html --
Content: {{ .Content }}
-- content/p1.md --
---
title: "p1"
---

> [!NOTE]  
> This is a note with some whitespace after the alert type.


> [!TIP]
> This is a tip.

>   [!CAUTION]
> This is a caution with some whitespace before the alert type.

> A regular blockquote.

> [!TIP]
> This is a tip with attributes.
{class="foo bar" id="baz"}

> [!NOTE]  
> Note triggering showing the position.
{showpos="true"}

`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html",
		"Blockquote Alert: |<p>This is a note with some whitespace after the alert type.</p>\n|alert|",
		"Blockquote Alert: |<p>This is a tip.</p>",
		"Blockquote Alert: |<p>This is a caution with some whitespace before the alert type.</p>\n|alert|",
		"Blockquote: |<p>A regular blockquote.</p>\n|regular|",
		"Blockquote Alert Attributes: |<p>This is a tip with attributes.</p>\n|map[class:foo bar id:baz]|",
		filepath.FromSlash("/content/p1.md:20:3"),
		"Blockquote Alert Page: |<p>This is a tip with attributes.</p>\n|p1|p1|",
	)
}
