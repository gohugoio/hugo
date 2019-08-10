// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGroupFunc(t *testing.T) {
	c := qt.New(t)

	pageContent := `
---
title: "Page"
---

`
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().
		WithContent("page1.md", pageContent, "page2.md", pageContent).
		WithTemplatesAdded("index.html", `
{{ $cool := .Site.RegularPages | group "cool" }}
{{ $cool.Key }}: {{ len $cool.Pages }}

`)
	b.CreateSites().Build(BuildCfg{})

	c.Assert(len(b.H.Sites), qt.Equals, 1)
	c.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 2)

	b.AssertFileContent("public/index.html", "cool: 2")
}

func TestSliceFunc(t *testing.T) {
	c := qt.New(t)

	pageContent := `
---
title: "Page"
tags: ["blue", "green"]
tags_weight: %d
---

`
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().
		WithContent("page1.md", fmt.Sprintf(pageContent, 10), "page2.md", fmt.Sprintf(pageContent, 20)).
		WithTemplatesAdded("index.html", `
{{ $cool := first 1 .Site.RegularPages | group "cool" }}
{{ $blue := after 1 .Site.RegularPages | group "blue" }}
{{ $weightedPages := index (index .Site.Taxonomies "tags") "blue" }}

{{ $p1 := index .Site.RegularPages 0 }}{{ $p2 := index .Site.RegularPages 1 }}
{{ $wp1 := index $weightedPages 0 }}{{ $wp2 := index $weightedPages 1 }}

{{ $pages := slice $p1 $p2 }}
{{ $pageGroups := slice $cool $blue }}
{{ $weighted := slice $wp1 $wp2 }}

{{ printf "pages:%d:%T:%v/%v" (len $pages) $pages (index $pages 0) (index $pages 1) }}
{{ printf "pageGroups:%d:%T:%v/%v" (len $pageGroups) $pageGroups (index (index $pageGroups 0).Pages 0) (index (index $pageGroups 1).Pages 0)}}
{{ printf "weightedPages:%d::%T:%v" (len $weighted) $weighted $weighted | safeHTML }}

`)
	b.CreateSites().Build(BuildCfg{})

	c.Assert(len(b.H.Sites), qt.Equals, 1)
	c.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 2)

	b.AssertFileContent("public/index.html",
		"pages:2:page.Pages:Page(/page1.md)/Page(/page2.md)",
		"pageGroups:2:page.PagesGroup:Page(/page1.md)/Page(/page2.md)",
		`weightedPages:2::page.WeightedPages:[WeightedPage(10,"Page") WeightedPage(20,"Page")]`)
}

func TestUnionFunc(t *testing.T) {
	c := qt.New(t)

	pageContent := `
---
title: "Page"
tags: ["blue", "green"]
tags_weight: %d
---

`
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().
		WithContent("page1.md", fmt.Sprintf(pageContent, 10), "page2.md", fmt.Sprintf(pageContent, 20),
			"page3.md", fmt.Sprintf(pageContent, 30)).
		WithTemplatesAdded("index.html", `
{{ $unionPages := first 2 .Site.RegularPages | union .Site.RegularPages  }}
{{ $unionWeightedPages := .Site.Taxonomies.tags.blue | union .Site.Taxonomies.tags.green }}
{{ printf "unionPages: %T %d" $unionPages (len $unionPages) }} 
{{ printf "unionWeightedPages: %T %d" $unionWeightedPages (len $unionWeightedPages) }}
`)
	b.CreateSites().Build(BuildCfg{})

	c.Assert(len(b.H.Sites), qt.Equals, 1)
	c.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 3)

	b.AssertFileContent("public/index.html",
		"unionPages: page.Pages 3",
		"unionWeightedPages: page.WeightedPages 6")
}

func TestCollectionsFuncs(t *testing.T) {
	c := qt.New(t)

	pageContent := `
---
title: "Page %d"
tags: ["blue", "green"]
tags_weight: %d
---

`
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().
		WithContent("page1.md", fmt.Sprintf(pageContent, 10, 10), "page2.md", fmt.Sprintf(pageContent, 20, 20),
			"page3.md", fmt.Sprintf(pageContent, 30, 30)).
		WithTemplatesAdded("index.html", `
{{ $uniqPages := first 2 .Site.RegularPages | append .Site.RegularPages | uniq  }}
{{ $inTrue := in .Site.RegularPages (index .Site.RegularPages 1)  }}
{{ $inFalse := in .Site.RegularPages (.Site.Home)  }}

{{ printf "uniqPages: %T %d" $uniqPages (len $uniqPages) }}
{{ printf "inTrue: %t" $inTrue }}
{{ printf "inFalse: %t" $inFalse  }}
`)

	b.WithTemplatesAdded("_default/single.html", `
{{ $related := .Site.RegularPages.Related . }}
{{ $symdiff := $related | symdiff .Site.RegularPages }}
Related: {{ range $related }}{{ .RelPermalink }}|{{ end }}
Symdiff: {{ range $symdiff }}{{ .RelPermalink }}|{{ end }}
`)
	b.CreateSites().Build(BuildCfg{})

	c.Assert(len(b.H.Sites), qt.Equals, 1)
	c.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 3)

	b.AssertFileContent("public/index.html",
		"uniqPages: page.Pages 3",
		"inTrue: true",
		"inFalse: false",
	)

	b.AssertFileContent("public/page1/index.html", `Related: /page2/|/page3/|`, `Symdiff: /page1/|`)
}

func TestAppendFunc(t *testing.T) {
	c := qt.New(t)

	pageContent := `
---
title: "Page"
tags: ["blue", "green"]
tags_weight: %d
---

`
	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().
		WithContent("page1.md", fmt.Sprintf(pageContent, 10), "page2.md", fmt.Sprintf(pageContent, 20)).
		WithTemplatesAdded("index.html", `
{{ $p1 := index .Site.RegularPages 0 }}{{ $p2 := index .Site.RegularPages 1 }}

{{ $pages := slice }}

{{ if true }}
	{{ $pages = $pages | append $p2 $p1 }}
{{ end }}
{{ $appendPages := .Site.Pages | append .Site.RegularPages }}
{{ $appendStrings := slice "a" "b" | append "c" "d" "e" }}
{{ $appendStringsSlice := slice "a" "b" "c" | append (slice "c" "d") }}

{{ printf "pages:%d:%T:%v/%v" (len $pages) $pages (index $pages 0) (index $pages 1)  }}
{{ printf "appendPages:%d:%T:%v/%v" (len $appendPages) $appendPages (index $appendPages 0).Kind (index $appendPages 8).Kind  }}
{{ printf "appendStrings:%T:%v"  $appendStrings $appendStrings  }}
{{ printf "appendStringsSlice:%T:%v"  $appendStringsSlice $appendStringsSlice }}

{{/* add some slightly related funcs to check what types we get */}}
{{ $u :=  $appendStrings | union $appendStringsSlice }}
{{ $i :=  $appendStrings | intersect $appendStringsSlice }}
{{ printf "union:%T:%v" $u $u  }}
{{ printf "intersect:%T:%v" $i $i }}

`)
	b.CreateSites().Build(BuildCfg{})

	c.Assert(len(b.H.Sites), qt.Equals, 1)
	c.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 2)

	b.AssertFileContent("public/index.html",
		"pages:2:page.Pages:Page(/page2.md)/Page(/page1.md)",
		"appendPages:9:page.Pages:home/page",
		"appendStrings:[]string:[a b c d e]",
		"appendStringsSlice:[]string:[a b c c d]",
		"union:[]string:[a b c d e]",
		"intersect:[]string:[a b c d]",
	)
}
