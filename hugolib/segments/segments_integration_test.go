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

package segments_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestSegmentsLegacy(t *testing.T) {
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
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .RelPermalink }}|
-- layouts/_default/list.html --
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

	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.ErrorMatches, ".* was deprecated and removed in v0.152.0.*")
}

// TODo1 check that we use dot vs / in matrix Globs.

func TestSegments(t *testing.T) {
	t.Skip("TODO1")
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
renderSegments = ["docs"]
[languages]
[languages.en]
weight = 1
[languages.no]
weight = 2
[languages.sv]
weight = 3

[module]
[[module.mounts]]
source = "content"
target = "content"
[module.mounts.sites.matrix]
languages = ['**']
[segments]
[segments.docs]
[[segments.docs.rules]]
[segments.docs.rules.sites.matrix]
languages = ["**"] # All languages except English
[[segments.docs.rules]]
output =  ["**"] # Only HTML output
[[segments.docs.rules]]
kind = "**" # Always render home page.
[[segments.docs.rules]]
path = "{/docs,/docs/**}"
-- layouts/all.html --
All: {{ .Title }}|{{ .RelPermalink }}|
-- content/_index.md --
-- content/docs/_index.md --
-- content/docs/section1/_index.md --
-- content/docs/section1/page1.md --
-- content/blog/_index.md --
-- content/blog/section1/page1.md --
`

	b := hugolib.Test(t, files)
	// b.Assert(b.H.Configs.Base.RootConfig.RenderSegments, qt.DeepEquals, []string{"docs"})

	b.AssertPublishDir("asdf")
}
