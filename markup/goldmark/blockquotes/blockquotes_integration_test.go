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
	b.AssertFileContentExact("public/p1/index.html",
		"Blockquote Alert: |<p>This is a note with some whitespace after the alert type.</p>|alert|",
		"Blockquote Alert: |<p>This is a tip.</p>",
		"Blockquote Alert: |<p>This is a caution with some whitespace before the alert type.</p>|alert|",
		"Blockquote: |<p>A regular blockquote.</p>|regular|",
		"Blockquote Alert Attributes: |<p>This is a tip with attributes.</p>|map[class:foo bar id:baz]|",
		filepath.FromSlash("/content/p1.md:19:3"),
		"Blockquote Alert Page: |<p>This is a tip with attributes.</p>|p1|p1|",

		// Issue 12767.
		"Blockquote Alert: |<p>Mixed case alert type.</p>|alert",
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
	b.AssertFileContentExact("public/p1/index.html", "Content: <blockquote>\n</blockquote>\n")
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

> [!info] Can callouts be nested?
> > [!important] Yes!, they can.
> > > [!tip] You can even use multiple layers of nesting.

-- layouts/index.html --
{{ .Content }}
-- layouts/_default/_markup/render-blockquote.html --
AlertType: {{ .AlertType }}|AlertTitle: {{ .AlertTitle }}|AlertSign: {{ .AlertSign | safeHTML }}|Text: {{ .Text }}|

	`

	b := hugolib.Test(t, files)
	b.AssertFileContentExact("public/index.html",
		"AlertType: tip|AlertTitle: Callouts can have custom titles|AlertSign: |",
		"AlertType: tip|AlertTitle: Title-only callout|AlertSign: |",
		"AlertType: faq|AlertTitle: Foldable negated callout|AlertSign: -|Text: <p>Yes! In a foldable callout, the contents are hidden when the callout is collapsed</p>|",
		"AlertType: faq|AlertTitle: Foldable callout|AlertSign: +|Text: <p>Yes! In a foldable callout, the contents are hidden when the callout is collapsed</p>|",
		"AlertType: danger|AlertTitle: |AlertSign: |Text: <p>Do not approach or handle without protective gear.</p>|",
		"AlertTitle: Can callouts be nested?|",
		"AlertTitle: You can even use multiple layers of nesting.|",
	)
}

// Issue 12913
// Issue 13119
// Issue 13302
func TestBlockquoteRenderHookTextParsing(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ .Content }}
-- layouts/_default/_markup/render-blockquote.html --
AlertType: {{ .AlertType }}|AlertTitle: {{ .AlertTitle }}|Text: {{ .Text }}|
-- content/_index.md --
---
title: home
---

> [!one]

> [!two] title

> [!three]
> line 1

> [!four] title
> line 1

> [!five]
> line 1
> line 2

> [!six] title
> line 1
> line 2

> [!seven]
> - list item

> [!eight] title
> - list item

> [!nine]
> line 1
> - list item

> [!ten] title
> line 1
> - list item

> [!eleven]
> line 1
> - list item
>
> line 2

> [!twelve] title
> line 1
> - list item
>
> line 2

> [!thirteen]
> ![alt](a.jpg)

> [!fourteen] title
> ![alt](a.jpg)

> [!fifteen] _title_

> [!sixteen] _title_
> line one

> seventeen
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"AlertType: one|AlertTitle: |Text: |",
		"AlertType: two|AlertTitle: title|Text: |",
		"AlertType: three|AlertTitle: |Text: <p>line 1</p>|",
		"AlertType: four|AlertTitle: title|Text: <p>line 1</p>|",
		"AlertType: five|AlertTitle: |Text: <p>line 1\nline 2</p>|",
		"AlertType: six|AlertTitle: title|Text: <p>line 1\nline 2</p>|",
		"AlertType: seven|AlertTitle: |Text: <ul>\n<li>list item</li>\n</ul>|",
		"AlertType: eight|AlertTitle: title|Text: <ul>\n<li>list item</li>\n</ul>|",
		"AlertType: nine|AlertTitle: |Text: <p>line 1</p>\n<ul>\n<li>list item</li>\n</ul>|",
		"AlertType: ten|AlertTitle: title|Text: <p>line 1</p>\n<ul>\n<li>list item</li>\n</ul>|",
		"AlertType: eleven|AlertTitle: |Text: <p>line 1</p>\n<ul>\n<li>list item</li>\n</ul>\n<p>line 2</p>|",
		"AlertType: twelve|AlertTitle: title|Text: <p>line 1</p>\n<ul>\n<li>list item</li>\n</ul>\n<p>line 2</p>|",
		"AlertType: thirteen|AlertTitle: |Text: <p><img src=\"a.jpg\" alt=\"alt\"></p>|",
		"AlertType: fourteen|AlertTitle: title|Text: <p><img src=\"a.jpg\" alt=\"alt\"></p>|",
		"AlertType: fifteen|AlertTitle: <em>title</em>|Text: |",
		"AlertType: sixteen|AlertTitle: <em>title</em>|Text: <p>line one</p>|",
		"AlertType: |AlertTitle: |Text: <p>seventeen</p>|",
	)
}
