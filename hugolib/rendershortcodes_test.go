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

package hugolib

import (
	"strings"
	"testing"
)

func TestRenderShortcodesBasic(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "term"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### p2-h1
{{< withhtml >}}
### p2-h2
{{% withmarkdown %}}
### p2-h3
{{% include "p3" %}}
-- content/p3.md --
---
title: "p3"
---
### p3-h1
{{< withhtml >}}
### p3-h2
{{% withmarkdown %}}
{{< level3 >}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/withhtml.html --
<div>{{ .Page.Title }} withhtml</div>
-- layouts/shortcodes/withmarkdown.html --
#### {{ .Page.Title }} withmarkdown
-- layouts/shortcodes/level3.html --
Level 3: {{ .Page.Title }}
-- layouts/_default/single.html --
Fragments: {{ .Fragments.Identifiers }}|
HasShortcode Level 1: {{ .HasShortcode "include" }}|
HasShortcode Level 2: {{ .HasShortcode "withmarkdown" }}|
HasShortcode Level 3: {{ .HasShortcode "level3" }}|
HasShortcode not found: {{ .HasShortcode "notfound" }}|
Content: {{ .Content }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Fragments: [p1-h1 p2-h1 p2-h2 p2-h3 p2-withmarkdown p3-h1 p3-h2 p3-withmarkdown]|",
		"HasShortcode Level 1: true|",
		"HasShortcode Level 2: true|",
		"HasShortcode Level 3: true|",
		"HasShortcode not found: false|",
	)
}

func TestRenderShortcodesNestedMultipleOutputFormatTemplates(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "term", "section", "rss", "sitemap", "robotsTXT", "404"]
[outputs]
page = ["html", "json"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### p2-h1
{{% myshort %}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/myshort.html --
Myshort HTML.
-- layouts/shortcodes/myshort.json --
Myshort JSON.
-- layouts/_default/single.html --
HTML: {{ .Content }}
-- layouts/_default/single.json --
JSON: {{ .Content }}


`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Myshort HTML")
	b.AssertFileContent("public/p1/index.json", "Myshort JSON")
}

func TestRenderShortcodesEditNested(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term", "section", "rss", "sitemap", "robotsTXT", "404"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### p2-h1
{{% myshort %}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/myshort.html --
Myshort Original.
-- layouts/_default/single.html --
 {{ .Content }}
`
	b := TestRunning(t, files)
	b.AssertFileContent("public/p1/index.html", "Myshort Original.")

	b.EditFileReplaceAll("layouts/shortcodes/myshort.html", "Original", "Edited").Build()
	b.AssertFileContent("public/p1/index.html", "Myshort Edited.")
}

func TestRenderShortcodesEditIncludedPage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
disableKinds = ["home", "taxonomy", "term", "section", "rss", "sitemap", "robotsTXT", "404"]
-- content/p1.md --
---
title: "p1"
---
## p1-h1
{{% include "p2" %}}
-- content/p2.md --
---
title: "p2"
---
### Original
{{% myshort %}}
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
-- layouts/shortcodes/myshort.html --
Myshort Original.
-- layouts/_default/single.html --
 {{ .Content }}



`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			Running:     true,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", "Original")

	b.EditFileReplaceFunc("content/p2.md", func(s string) string {
		return strings.Replace(s, "Original", "Edited", 1)
	})
	b.Build()
	b.AssertFileContent("public/p1/index.html", "Edited")
}
