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

-- layouts/shortcodes/see-also.html --
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

-- layouts/_default/single.html --
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

	b.AssertFileContent("public/p1/index.html",
		"Related 1:  0: p2: h0: First title|ref1|::END",
		"Related 2:  0: p5: h0: Common p3, p4, p5|common-p3-p4-p5|::END 1: p4: h0: Common p3, p4, p5|common-p3-p4-p5|::END",
	)
}

func BenchmarkRelatedSite(b *testing.B) {
	files := `
-- config.toml --
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
-- layouts/_default/single.html --
Len related: {{ site.RegularPages.Related . | len }}
`

	createContent := func(n int) string {
		base := `---
title: "Page %d"
keywords: ['k%d']
---
`

		for i := 0; i < 32; i++ {
			base += fmt.Sprintf("\n## Title %d", rand.Intn(100))
		}

		return fmt.Sprintf(base, n, rand.Intn(32))
	}

	for i := 1; i < 100; i++ {
		files += fmt.Sprintf("\n-- content/posts/p%d.md --\n"+createContent(i+1), i+1)
	}

	cfg := hugolib.IntegrationTestConfig{
		T:           b,
		TxtarString: files,
	}
	builders := make([]*hugolib.IntegrationTestBuilder, b.N)

	for i := range builders {
		builders[i] = hugolib.NewIntegrationTestBuilder(cfg)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builders[i].Build()
	}
}
