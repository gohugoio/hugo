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

package related_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestRelatedFragments(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT"]
[related]
  includeNewer = false
  threshold = 80
  toLower = false
[[related.indices]]
  name = 'pagerefs'
  type = 'fragments'
  applyFilter = true
  weight = 90
[[related.indices]]
  name = 'keywords'
  weight = 80
-- content/p1.md --
---
title: p1
pagerefs: ['ref1']
---
{{< see-also >}}

## P1 title

-- content/p2.md --
---
title: p2
---

## P2 title 1

## P2 title 2

## First title {#ref1}
{{< see-also "ref1" >}}
-- content/p3.md --
---
title: p3
keywords: ['foo']
---

## P3 title 1

## P3 title 2

## Common p3, p4, p5
-- content/p4.md --
---
title: p4
---

## Common p3, p4, p5

## P4 title 1

-- content/p5.md --
---
title: p5
keywords: ['foo']
---

## P5 title 1

## Common p3, p4, p5

-- layouts/_shortcodes/see-also.html --
{{ $p1 := site.GetPage "p1" }}
{{ $p2 := site.GetPage "p2" }}
{{ $p3 := site.GetPage "p3" }}
P1 Fragments: {{ $p1.Fragments.Identifiers }}
P2 Fragments: {{ $p2.Fragments.Identifiers }}
Contains ref1: {{ $p2.Fragments.Identifiers.Contains "ref1" }}
Count ref1: {{ $p2.Fragments.Identifiers.Count "ref1" }}
{{ $opts := dict "document" .Page "fragments" $.Params }}
{{ $related1 := site.RegularPages.Related $opts }}
{{ $related2 := site.RegularPages.Related $p3 }}
Len Related 1: {{ len $related1 }}
Len Related 2: {{ len $related2 }}
Related 1: {{ template "list-related" $related1 }}
Related 2: {{ template "list-related" $related2 }}

{{ define "list-related" }}{{ range $i, $e := . }} {{ $i }}: {{ .Title }}: {{ with .HeadingsFiltered}}{{ range $i, $e := .}}h{{ $i }}: {{ .Title }}|{{ .ID }}|{{ end }}{{ end }}::END{{ end }}{{ end }}

-- layouts/single.html --
Content: {{ .Content }}


`

	b := hugolib.Test(t, files)

	expect := `
P1 Fragments: [p1-title]
P2 Fragments: [p2-title-1 p2-title-2 ref1]
Len Related 1: 1
Related 2: 2
`

	for _, p := range []string{"p1", "p2"} {
		b.AssertFileContent("public/"+p+"/index.html", expect)
	}

	b.AssertFileContent(
		"public/p1/index.html",
		"Related 1:  0: p2: h0: First title|ref1|::END",
		"Related 2:  0: p5: h0: Common p3, p4, p5|common-p3-p4-p5|::END 1: p4: h0: Common p3, p4, p5|common-p3-p4-p5|::END",
	)
}

// See issue 15199.
func TestRelatedFragmentsTokenize(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home', 'section', 'taxonomy', 'term', 'rss', 'sitemap']
[related]
  includeNewer = true
  threshold = 80
  toLower = true
[[related.indices]]
  type = 'fragments'
  weight = 100
  tokenize = true
-- content/p1.md --
---
title: p1
---
## Hugo Templates Overview
-- content/p2.md --
---
title: p2
---
## Templates in Hugo
-- content/p3.md --
---
title: p3
---
## JavaScript and CSS
-- layouts/page.html --
Related: {{ range site.RegularPages.Related . }}{{ .Title }} {{ end }}
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Related: p2")
	b.AssertFileContent("public/p2/index.html", "Related: p1")
	b.AssertFileContent("public/p3/index.html", "Related:  ")
}

// See issue 15199.
func TestRelatedFragmentsTokenizeHTMLEntities(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
disableKinds = ["home", "section", "taxonomy", "term", "rss", "sitemap"]
[related]
  includeNewer = true
  threshold = 80
  toLower = true
[[related.indices]]
  type = 'fragments'
  weight = 100
  tokenize = true
-- content/p1.md --
---
title: p1
---
## Types of apples
-- content/p2.md --
---
title: p2
---
## Something about "apples"
-- layouts/page.html --
Related: {{ range site.RegularPages.Related . }}{{ .Title }} {{ end }}
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Related: p2")
	b.AssertFileContent("public/p2/index.html", "Related: p1")
}

// See issue 15199.
func TestRelatedBasicTokenize(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home', 'section', 'taxonomy', 'term', 'rss', 'sitemap']
[related]
  includeNewer = true
  threshold = 80
  toLower = true
[[related.indices]]
  name = 'title'
  weight = 100
  tokenize = true
-- content/p1.md --
---
title: Hugo Templates Overview
---
-- content/p2.md --
---
title: Templates in Hugo
---
-- content/p3.md --
---
title: JavaScript and CSS
---
-- layouts/page.html --
Related: {{ range site.RegularPages.Related . }}{{ .Title }} {{ end }}
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Related: Templates in Hugo")
	b.AssertFileContent("public/p2/index.html", "Related: Hugo Templates Overview")
	b.AssertFileContent("public/p3/index.html", "Related:  ")
}

// See issue 15199.
func TestRelatedFragmentsTokenizeApplyFilter(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home', 'section', 'taxonomy', 'term', 'rss', 'sitemap']
[related]
  includeNewer = true
  threshold = 80
  toLower = true
[[related.indices]]
  type = 'fragments'
  weight = 100
  applyFilter = true
  tokenize = true
-- content/p1.md --
---
title: p1
---
## Hugo Templates Overview
-- content/p2.md --
---
title: p2
---
## Templates in Hugo

## Unrelated Topic
-- layouts/page.html --
Related: {{ range site.RegularPages.Related . }}{{ .Title }}:{{ with .HeadingsFiltered }}{{ range . }}{{ .Title }}|{{ end }}{{ end }}::END {{ end }}
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Related: p2:Templates in Hugo|::END")
	b.AssertFileContent("public/p2/index.html", "Related: p1:Hugo Templates Overview|::END")
}

// See issue 15199.
func TestRelatedTokenizeCaseSensitive(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home', 'section', 'taxonomy', 'term', 'rss', 'sitemap']
[related]
  includeNewer = true
  threshold = 80
  toLower = false
[[related.indices]]
  type = 'fragments'
  weight = 100
  tokenize = true
-- content/p1.md --
---
title: p1
---
## Hugo Templates
-- content/p2.md --
---
title: p2
---
## hugo templates
-- layouts/page.html --
Related: {{ range site.RegularPages.Related . }}{{ .Title }} {{ end }}
`
	b := hugolib.Test(t, files)

	// toLower=false: "Hugo" != "hugo", so no match.
	b.AssertFileContent("public/p1/index.html", "Related:  ")
	b.AssertFileContent("public/p2/index.html", "Related:  ")
}

// See issue 15238.
func TestRelatedTokenizeMinTokenLength(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home', 'section', 'taxonomy', 'term', 'rss', 'sitemap']
[related]
  includeNewer = true
  threshold = 80
  toLower = true
[[related.indices]]
  name = 'title'
  weight = 100
  tokenize = true
  minTokenLength = 4
-- content/p1.md --
---
title: The long journey home
---
-- content/p2.md --
---
title: Planning your next big journey
---
-- content/p3.md --
---
title: Go to the top
---
-- layouts/page.html --
Related: {{ range site.RegularPages.Related . }}{{ .Title }} {{ end }}
`
	b := hugolib.Test(t, files)

	// p1 and p2 share "journey" (>= minTokenLength=4); short words are filtered.
	b.AssertFileContent("public/p1/index.html", "Related: Planning your next big journey")
	b.AssertFileContent("public/p2/index.html", "Related: The long journey home")
	// All tokens in p3's title ("go", "to", "the", "top") are < minTokenLength=4.
	b.AssertFileContent("public/p3/index.html", "Related:  ")
}

// Repetitive headings must not inflate a document's rank above a document
// that matches on an additional index.
func TestRelatedFragmentsTokenizeDeduplication(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home', 'section', 'taxonomy', 'term', 'rss', 'sitemap']
[related]
  includeNewer = true
  threshold = 80
  toLower = true
[[related.indices]]
  type = 'fragments'
  weight = 100
  tokenize = true
[[related.indices]]
  name = 'keywords'
  weight = 100
-- content/p1.md --
---
title: p1
keywords: ['fruit']
---
## I love apples, apples and more apples
-- content/p2.md --
---
title: p2
keywords: ['fruit']
---
## I love apples, apples and more apples
## I love apples, apples and more apples
## I love apples, apples and more apples
-- content/p3.md --
---
title: p3
---
## I love apples, apples and more apples
## I love apples, apples and more apples
## I love apples, apples and more apples
-- layouts/page.html --
Related: {{ range site.RegularPages.Related . }}{{ .Title }} {{ end }}
`
	b := hugolib.Test(t, files)

	// p1 and p3 share the same unique fragment tokens; p1 also matches on
	// keywords ("fruit"). Without deduplication, p3's repeated headings
	// inflate its weight above p1's combined fragment+keywords weight.
	b.AssertFileContent("public/p2/index.html", "Related: p1 p3")
}

func BenchmarkRelatedSite(b *testing.B) {
	var files strings.Builder
	files.WriteString(`
-- hugo.toml --
baseURL = "http://example.com/"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT"]
[related]
  includeNewer = false
  threshold = 80
  toLower = false
[[related.indices]]
  name = 'keywords'
  weight = 70
[[related.indices]]
  name = 'pagerefs'
  type = 'fragments'
  weight = 30
-- layouts/single.html --
Len related: {{ site.RegularPages.Related . | len }}
`)

	createContent := func(n int) string {
		var base strings.Builder
		base.WriteString(`---
title: "Page %d"
keywords: ['k%d']
---
`)

		for range 32 {
			base.WriteString(fmt.Sprintf("\n## Title %d", rand.Intn(100)))
		}

		return fmt.Sprintf(base.String(), n, rand.Intn(32))
	}

	for i := 1; i < 100; i++ {
		files.WriteString(fmt.Sprintf("\n-- content/posts/p%d.md --\n"+createContent(i+1), i+1))
	}

	cfg := hugolib.IntegrationTestConfig{
		T:           b,
		TxtarString: files.String(),
	}

	for b.Loop() {
		b.StopTimer()
		bb := hugolib.NewIntegrationTestBuilder(cfg)
		b.StartTimer()
		bb.Build()
	}
}
