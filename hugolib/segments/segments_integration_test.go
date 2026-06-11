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

package segments_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestSegments(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
renderSegments = ["docs"]
[languages]
[languages.en]
weight = 1
[languages.no]
weight = 2
[languages.nb]
weight = 3
[segments]
[segments.docs]
[[segments.docs.includes]]
kind = "{home,taxonomy,term}"
[[segments.docs.includes]]
path = "{/docs,/docs/**}"
[[segments.docs.excludes]]
path = "/blog/**"
[[segments.docs.excludes]]
lang = "n*"
output = "rss"
[[segments.docs.excludes]]
output = "json"
-- layouts/single.html --
Single: {{ .Title }}|{{ .RelPermalink }}|
-- layouts/list.html --
List: {{ .Title }}|{{ .RelPermalink }}|
-- content/docs/_index.md --
-- content/docs/section1/_index.md --
-- content/docs/section1/page1.md --
---
title: "Docs Page 1"
tags: ["tag1", "tag2"]
---
-- content/blog/_index.md --
-- content/blog/section1/page1.md --
---
title: "Blog Page 1"
tags: ["tag1", "tag2"]
---
`

	b := hugolib.Test(t, files, hugolib.TestOptInfo())
	b.AssertLogContains("deprecated") // lang => sites.matrix in v0.152.0
	b.Assert(b.H.Configs.Base.RootConfig.RenderSegments, qt.DeepEquals, []string{"docs"})

	b.AssertFileContent("public/docs/section1/page1/index.html", "Docs Page 1")
	b.AssertFileExists("public/blog/section1/page1/index.html", false)
	b.AssertFileExists("public/index.html", true)
	b.AssertFileExists("public/index.xml", true)
	b.AssertFileExists("public/no/index.html", true)
	b.AssertFileExists("public/no/index.xml", false)
}

// See issue 15024.
func TestSegmentsMultiple(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
renderSegments = SEGMENTS
disableKinds = ["home", "taxonomy", "term", "page"]
[outputs]
section = ['html', 'json']
[segments]
[segments.excludeallkinds]
[[segments.excludeallkinds.excludes]]
kind = "**"
[segments.blog]
[[segments.blog.includes]]
path = "{/blog,/blog/**}"
[[segments.blog.excludes]]
output = 'json'
[segments.news]
[[segments.news.includes]]
path = "{/news,/news/**}"
-- layouts/all.html --
{{ .Kind }}: {{ .Title }}|{{ .RelPermalink }}|
-- layouts/all.json --
{{ .Kind }}: {{ .Title }}|{{ .RelPermalink }}|
-- content/blog/_index.md --
-- content/blog/page1.md --
---
title: "Blog Page 1"
tags: ["tag1", "tag2"]
---
-- content/news/_index.md --
-- content/news/page1.md --
---
title: "News Page 1"
tags: ["tag1", "tag2"]
---
`
	files := strings.ReplaceAll(filesTemplate, "SEGMENTS", `["excludeallkinds", "blog", "news"]`)

	b := hugolib.Test(t, files, hugolib.TestOptInfo())

	b.AssertPublishDir(`
blog/index.html
! blog/index.json
news/index.html
news/index.json
`)

	files = strings.ReplaceAll(filesTemplate, "SEGMENTS", `["excludeallkinds"]`)

	b = hugolib.Test(t, files, hugolib.TestOptInfo())

	b.AssertPublishDir(`
! json
! html
`)
}

// See issue 14939.
func TestRenderSegmentsMergesHugoStatsJSON(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
renderSegments = ["docs"]
[build.buildStats]
enable = true
[segments]
[segments.docs]
[[segments.docs.includes]]
path = "{/docs,/docs/**}"
-- hugo_stats.json --
{
  "htmlElements": {
    "tags": [
      "section"
    ],
    "classes": [
      "from-previous-build"
    ],
    "ids": [
      "previous-id"
    ]
  }
}
-- layouts/single.html --
<div id="docs-id" class="docs-class">{{ .Title }}</div>
-- layouts/list.html --
<div id="list-id" class="list-class">{{ .Title }}</div>
-- content/docs/section1/page1.md --
---
title: "Docs Page 1"
---
-- content/blog/section1/page1.md --
---
title: "Blog Page 1"
---
`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())

	b.AssertFileContent("public/docs/section1/page1/index.html", "Docs Page 1")
	b.AssertFileExists("public/blog/section1/page1/index.html", false)

	b.AssertFileContent("hugo_stats.json",
		"from-previous-build",
		"previous-id",
		"section",
		"docs-class",
	)
}
