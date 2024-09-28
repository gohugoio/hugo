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

package tables_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestTableHook(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
[markup.goldmark.parser.attribute]
block = true
title = true
-- content/p1.md --
## Table 1

| Item              | In Stock | Price |
| :---------------- | :------: | ----: |
| Python Hat        |   True   | 23.99 |
| SQL **Hat**       |   True   | 23.99 |
| Codecademy Tee    |  False   | 19.99 |
| Codecademy Hoodie |  False   | 42.99 |
{.foo foo="bar"}

## Table 2

| Month | Savings |
| -------- | ------- |
| January | $250 |
| February | $80 |
| March | $420 |

-- layouts/_default/single.html --
{{ .Content }}
-- layouts/_default/_markup/render-table.html --
Attributes: {{ .Attributes }}|
{{ template "print" (dict "what" (printf "table-%d-thead" $.Ordinal) "rows" .THead) }}
{{ template "print" (dict "what" (printf "table-%d-tbody" $.Ordinal)  "rows" .TBody) }}
{{ define "print" }}
 {{ .what }}:{{ range $i, $a := .rows }} {{ $i }}:{{ range $j, $b := . }} {{ $j }}: {{ .Alignment }}: {{ .Text }}|{{ end }}{{ end }}$
{{ end }}

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Attributes: map[class:foo foo:bar]|",
		"table-0-thead: 0: 0: left: Item| 1: center: In Stock| 2: right: Price|$",
		"table-0-tbody: 0: 0: left: Python Hat| 1: center: True| 2: right: 23.99| 1: 0: left: SQL <strong>Hat</strong>| 1: center: True| 2: right: 23.99| 2: 0: left: Codecademy Tee| 1: center: False| 2: right: 19.99| 3: 0: left: Codecademy Hoodie| 1: center: False| 2: right: 42.99|$",
	)

	b.AssertFileContent("public/p1/index.html",
		"table-1-thead: 0: 0: : Month| 1: : Savings|$",
		"table-1-tbody: 0: 0: : January| 1: : $250| 1: 0: : February| 1: : $80| 2: 0: : March| 1: : $420|$",
	)
}

func TestTableDefault(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
[markup.goldmark.parser.attribute]
block = true
title = true
-- content/p1.md --

## Table 1

| Item              | In Stock | Price |
| :---------------- | :------: | ----: |
| Python Hat        |   True   | 23.99 |
| SQL Hat           |   True   | 23.99 |
| Codecademy Tee    |  False   | 19.99 |
| Codecademy Hoodie |  False   | 42.99 |
{.foo}


-- layouts/_default/single.html --
Summary: {{ .Summary }}
Content: {{ .Content }}

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "<table class=\"foo\">")
}

// Issue 12811.
func TestTableDefaultRSSAndHTML(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
[outputFormats]
  [outputFormats.rss]
    weight = 30
 [outputFormats.html]
    weight = 20
-- content/_index.md --
---
title: "Home"
output: ["rss", "html"]
---

| Item              | In Stock | Price |
| :---------------- | :------: | ----: |
| Python Hat        |   True   | 23.99 |
| SQL Hat           |   True   | 23.99 |
| Codecademy Tee    |  False   | 19.99 |
| Codecademy Hoodie |  False   | 42.99 |

{{< foo >}}

-- layouts/index.html --
Content: {{ .Content }}
-- layouts/index.xml --
Content: {{ .Content }}
-- layouts/shortcodes/foo.xml --
foo xml
-- layouts/shortcodes/foo.html --
foo html

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.xml", "<table>")
	b.AssertFileContent("public/index.html", "<table>")
}

func TestTableDefaultRSSOnly(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
[outputs]
  home = ['rss']
  section = ['rss']
  taxonomy = ['rss']
  term = ['rss']
  page = ['rss']
disableKinds = ["taxonomy", "term", "page", "section"]
-- content/_index.md --
---
title: "Home"
---

## Table 1

| Item              | In Stock | Price |
| :---------------- | :------: | ----: |
| Python Hat        |   True   | 23.99 |
| SQL Hat           |   True   | 23.99 |
| Codecademy Tee    |  False   | 19.99 |
| Codecademy Hoodie |  False   | 42.99 |





-- layouts/index.xml --
Content: {{ .Content }}


`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.xml", "<table>")
}
