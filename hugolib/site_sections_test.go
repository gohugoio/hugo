// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"testing"
)

func TestNextInSectionNested(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
-- content/blog/page1.md --
---
weight: 1
---
-- content/blog/page2.md --
---
weight: 2
---
-- content/blog/cool/_index.md --
---
weight: 1
---
-- content/blog/cool/cool1.md --
---
weight: 1
---
-- content/blog/cool/cool2.md --
---
weight: 2
---
-- content/root1.md --
---
weight: 1
---
-- content/root2.md --
---
weight: 2
---
-- layouts/single.html --
Prev: {{ with .PrevInSection }}{{ .RelPermalink }}{{ end }}|
Next: {{ with .NextInSection }}{{ .RelPermalink }}{{ end }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/root1/index.html",
		"Prev: /root2/|", "Next: |")
	b.AssertFileContent("public/root2/index.html",
		"Prev: |", "Next: /root1/|")
	b.AssertFileContent("public/blog/page1/index.html",
		"Prev: /blog/page2/|", "Next: |")
	b.AssertFileContent("public/blog/page2/index.html",
		"Prev: |", "Next: /blog/page1/|")
	b.AssertFileContent("public/blog/cool/cool1/index.html",
		"Prev: /blog/cool/cool2/|", "Next: |")
	b.AssertFileContent("public/blog/cool/cool2/index.html",
		"Prev: |", "Next: /blog/cool/cool1/|")
}

func TestSectionEntries(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
-- content/myfirstsection/p1.md --
---
title: "P1"
---
P1
-- content/a/b/c/_index.md --
---
title: "C"
---
C
-- content/a/b/c/mybundle/index.md --
---
title: "My Bundle"
---
-- layouts/_default/list.html --
Kind: {{ .Kind }}|RelPermalink: {{ .RelPermalink }}|SectionsPath: {{ .SectionsPath }}|SectionsEntries: {{ .SectionsEntries }}|Len: {{ len .SectionsEntries }}|
-- layouts/_default/single.html --
Kind: {{ .Kind }}|RelPermalink: {{ .RelPermalink }}|SectionsPath: {{ .SectionsPath }}|SectionsEntries: {{ .SectionsEntries }}|Len: {{ len .SectionsEntries }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/myfirstsection/p1/index.html", "RelPermalink: /myfirstsection/p1/|SectionsPath: /myfirstsection|SectionsEntries: [myfirstsection]|Len: 1")
	b.AssertFileContent("public/a/b/c/index.html", "RelPermalink: /a/b/c/|SectionsPath: /a/b/c|SectionsEntries: [a b c]|Len: 3")
	b.AssertFileContent("public/a/b/c/mybundle/index.html", "Kind: page|RelPermalink: /a/b/c/mybundle/|SectionsPath: /a/b/c|SectionsEntries: [a b c]|Len: 3")
	b.AssertFileContent("public/index.html", "Kind: home|RelPermalink: /|SectionsPath: /|SectionsEntries: []|Len: 0")
}

func TestParentWithPageOverlap(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com/"
-- content/docs/_index.md --
-- content/docs/logs/_index.md --
-- content/docs/logs/sdk.md --
-- content/docs/logs/sdk_exporters/stdout.md --
-- layouts/_default/list.html --
{{ .RelPermalink }}|{{ with .Parent}}{{ .RelPermalink }}{{ end }}|
-- layouts/_default/single.html --
{{ .RelPermalink }}|{{ with .Parent}}{{ .RelPermalink }}{{ end }}|

`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "/||")
	b.AssertFileContent("public/docs/index.html", "/docs/|/|")
	b.AssertFileContent("public/docs/logs/index.html", "/docs/logs/|/docs/|")
	b.AssertFileContent("public/docs/logs/sdk/index.html", "/docs/logs/sdk/|/docs/logs/|")
	b.AssertFileContent("public/docs/logs/sdk_exporters/stdout/index.html", "/docs/logs/sdk_exporters/stdout/|/docs/logs/|")
}

func TestNestedSectionsEmpty(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
-- content/a/b/c/_index.md --
---
title: "C"
---
-- layouts/all.html --
All: {{ .Title }}|{{ .Kind }}|
`
	b := Test(t, files)

	b.AssertFileContent("public/a/index.html", "All: As|section|")
	b.AssertFileExists("public/a/b/index.html", false)
	b.AssertFileContent("public/a/b/c/index.html", "All: C|section|")
}
