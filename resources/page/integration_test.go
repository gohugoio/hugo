// Copyright 2021 The Hugo Authors. All rights reserved.
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

package page_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestGroupByLocalizedDate(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
title = 'My blog'
weight = 1
[languages.fr]
title = 'Mon blogue'
weight = 2
[languages.nn]
title = 'Bloggen min'
weight = 3
-- content/p1.md --
---
title: "Post 1"
date: "2020-01-01"
---
-- content/p2.md --
---
title: "Post 2"
date: "2020-02-01"
---
-- content/p1.fr.md --
---
title: "Post 1"
date: "2020-01-01"
---
-- content/p2.fr.md --
---
title: "Post 2"
date: "2020-02-01"
---
-- layouts/index.html --
{{ range $k, $v := site.RegularPages.GroupByDate "January, 2006" }}{{ $k }}|{{ $v.Key }}|{{ $v.Pages }}{{ end }}

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/en/index.html", "0|February, 2020|Pages(1)1|January, 2020|Pages(1)")
	b.AssertFileContent("public/fr/index.html", "0|février, 2020|Pages(1)1|janvier, 2020|Pages(1)")
}

func TestPagesSortCollation(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
title = 'My blog'
weight = 1
[languages.fr]
title = 'Mon blogue'
weight = 2
[languages.nn]
title = 'Bloggen min'
weight = 3
-- content/p1.md --
---
title: "zulu"
date: "2020-01-01"
param1: "xylophone"
tags: ["xylophone", "éclair", "zulu", "emma"]
---
-- content/p2.md --
---
title: "émotion"
date: "2020-01-01"
param1: "violin"
---
-- content/p3.md --
---
title: "alpha"
date: "2020-01-01"
param1: "éclair"
---
-- layouts/index.html --
ByTitle: {{ range site.RegularPages.ByTitle }}{{ .Title }}|{{ end }}
ByLinkTitle: {{ range site.RegularPages.ByLinkTitle }}{{ .Title }}|{{ end }}
ByParam: {{ range site.RegularPages.ByParam "param1" }}{{ .Params.param1 }}|{{ end }}
Tags Alphabetical: {{  range site.Taxonomies.tags.Alphabetical }}{{ .Term }}|{{ end }}
GroupBy: {{ range site.RegularPages.GroupBy "Title" }}{{ .Key }}|{{ end }}
{{ with (site.GetPage "p1").Params.tags }}
Sort: {{  sort . }}
ByWeight: {{ range site.RegularPages.ByWeight }}{{ .Title }}|{{ end }}
{{ end }}

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/en/index.html", `
ByTitle: alpha|émotion|zulu|
ByLinkTitle: alpha|émotion|zulu|
ByParam: éclair|violin|xylophone
Tags Alphabetical: éclair|emma|xylophone|zulu|
GroupBy: alpha|émotion|zulu|
Sort: [éclair emma xylophone zulu]
ByWeight: alpha|émotion|zulu|
`)
}

// See #10377
func TestPermalinkExpansionSectionsRepeated(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "taxonomyTerm", "sitemap"]
[outputs]
home = ["HTML"]
page = ["HTML"]
section = ["HTML"]
[outputFormats]
[permalinks]
posts = '/:sections[1]/:sections[last]/:slug'
-- content/posts/_index.md --
-- content/posts/a/_index.md --
-- content/posts/a/b/_index.md --
-- content/posts/a/b/c/_index.md --
-- content/posts/a/b/c/d.md --
---
title: "D"
slug: "d"
---
D
-- layouts/_default/single.html --
RelPermalink: {{ .RelPermalink }}

`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		}).Build()

	b.AssertFileContent("public/a/c/d/index.html", "RelPermalink: /a/c/d/")

}
