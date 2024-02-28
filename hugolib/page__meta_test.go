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

// Issue 9793
// Issue 12115
func TestListTitles(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','sitemap']
capitalizeListTitles = true
pluralizeListTitles = true
[taxonomies]
tag = 'tags'
-- content/section-1/page-1.md --
---
title: page-1
tags: 'tag-a'
---
-- layouts/_default/list.html --
{{ .Title }}
-- layouts/_default/single.html --
{{ .Title }}
	`

	b := Test(t, files)

	b.AssertFileContent("public/section-1/index.html", "Section-1s")
	b.AssertFileContent("public/tags/index.html", "Tags")
	b.AssertFileContent("public/tags/tag-a/index.html", "Tag-A")

	files = strings.Replace(files, "true", "false", -1)

	b = Test(t, files)

	b.AssertFileContent("public/section-1/index.html", "section-1")
	b.AssertFileContent("public/tags/index.html", "tags")
	b.AssertFileContent("public/tags/tag-a/index.html", "tag-a")
}

func TestDraftNonDefaultContentLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "en"
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/p1.md --
-- content/p2.nn.md --
---
title: "p2"
draft: true
---
-- layouts/_default/single.html --
{{ .Title }}
`
	b := Test(t, files)

	b.AssertFileExists("public/nn/p2/index.html", false)
}
