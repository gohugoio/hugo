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
Blockquote: |{{ .Text }}|{{ .Type }}|
-- layouts/_default/_markup/render-blockquote-alert.html --
{{ $text := .Text }}
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


> [!nOtE]  
> Mixed case alert type.


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

		// Issue 12767.
		"Blockquote Alert: |<p>Mixed case alert type.</p>\n|alert",
	)
}

func TestBlockquoteEmptyIssue12756(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/p1.md --
---
title: "p1"
---

>
-- layouts/_default/single.html --
Content: {{ .Content }}

`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html", "Content: <blockquote>\n</blockquote>\n")
}

func TestBlockquObsidianWithTitleAndSign(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/_index.md --
---
title: "Home"
---

> [!danger]
> Do not approach or handle without protective gear.


> [!tip] Callouts can have custom titles
> Like this one.

> [!tip] Title-only callout

> [!faq]- Foldable negated callout
> Yes! In a foldable callout, the contents are hidden when the callout is collapsed

> [!faq]+ Foldable callout
> Yes! In a foldable callout, the contents are hidden when the callout is collapsed

-- layouts/index.html --
{{ .Content }}
-- layouts/_default/_markup/render-blockquote.html --
AlertType: {{ .AlertType }}|
AlertTitle: {{ .AlertTitle }}|
AlertSign: {{ .AlertSign | safeHTML }}|
Text: {{ .Text }}|
	
	`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html",
		"AlertType: tip|\nAlertTitle: Callouts can have custom titles|\nAlertSign: |",
		"AlertType: tip|\nAlertTitle: Title-only callout</p>|\nAlertSign: |",
		"AlertType: faq|\nAlertTitle: Foldable negated callout|\nAlertSign: -|\nText: <p>Yes!",
		"AlertType: faq|\nAlertTitle: Foldable callout|\nAlertSign: +|",
		"AlertType: danger|\nAlertTitle: |\nAlertSign: |\nText: <p>Do not approach or handle without protective gear.</p>\n|",
	)
}
