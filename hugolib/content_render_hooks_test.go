// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless requiredF by applicable law or agreed to in writing, software
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

func TestRenderHooks(t *testing.T) {
	config := `
baseURL="https://example.org"
workingDir="/mywork"
`
	b := newTestSitesBuilder(t).WithWorkingDir("/mywork").WithConfigFile("toml", config).Running()
	b.WithTemplatesAdded("_default/single.html", `{{ .Content }}`)
	b.WithTemplatesAdded("shortcodes/myshortcode1.html", `{{ partial "mypartial1" }}`)
	b.WithTemplatesAdded("shortcodes/myshortcode2.html", `{{ partial "mypartial2" }}`)
	b.WithTemplatesAdded("shortcodes/myshortcode3.html", `SHORT3|`)
	b.WithTemplatesAdded("shortcodes/myshortcode4.html", `
<div class="foo">
{{ .Inner | markdownify }}
</div>
`)
	b.WithTemplatesAdded("shortcodes/myshortcode5.html", `
Inner Inline: {{ .Inner | .Page.RenderString }}
Inner Block: {{ .Inner | .Page.RenderString (dict "display" "block" ) }}
`)

	b.WithTemplatesAdded("shortcodes/myshortcode6.html", `.Render: {{ .Page.Render "myrender" }}`)
	b.WithTemplatesAdded("partials/mypartial1.html", `PARTIAL1`)
	b.WithTemplatesAdded("partials/mypartial2.html", `PARTIAL2  {{ partial "mypartial3.html" }}`)
	b.WithTemplatesAdded("partials/mypartial3.html", `PARTIAL3`)
	b.WithTemplatesAdded("partials/mypartial4.html", `PARTIAL4`)
	b.WithTemplatesAdded("customview/myrender.html", `myrender: {{ .Title }}|P4: {{ partial "mypartial4" }}`)
	b.WithTemplatesAdded("_default/_markup/render-link.html", `{{ with .Page }}{{ .Title }}{{ end }}|{{ .Destination | safeURL }}|Title: {{ .Title | safeHTML }}|Text: {{ .Text | safeHTML }}|END`)
	b.WithTemplatesAdded("docs/_markup/render-link.html", `Link docs section: {{ .Text | safeHTML }}|END`)
	b.WithTemplatesAdded("_default/_markup/render-image.html", `IMAGE: {{ .Page.Title }}||{{ .Destination | safeURL }}|Title: {{ .Title | safeHTML }}|Text: {{ .Text | safeHTML }}|END`)

	b.WithContent("customview/p1.md", `---
title: Custom View
---

{{< myshortcode6 >}}

	`, "blog/p1.md", `---
title: Cool Page
---

[First Link](https://www.google.com "Google's Homepage")

{{< myshortcode3 >}}

[Second Link](https://www.google.com "Google's Homepage")

Image:

![Drag Racing](/images/Dragster.jpg "image title")


`, "blog/p2.md", `---
title: Cool Page2
layout: mylayout
---

{{< myshortcode1 >}}

[Some Text](https://www.google.com "Google's Homepage")

,[No Whitespace Please](https://gohugo.io),



`, "blog/p3.md", `---
title: Cool Page3
---

{{< myshortcode2 >}}


`, "docs/docs1.md", `---
title: Docs 1
---


[Docs 1](https://www.google.com "Google's Homepage")


`, "blog/p4.md", `---
title: Cool Page With Image
---

Image:

![Drag Racing](/images/Dragster.jpg "image title")


`, "blog/p5.md", `---
title: Cool Page With Markdownify
---

{{< myshortcode4 >}}
Inner Link: [Inner Link](https://www.google.com "Google's Homepage")
{{< /myshortcode4 >}}

`, "blog/p6.md", `---
title: With RenderString
---

{{< myshortcode5 >}}Inner Link: [Inner Link](https://www.gohugo.io "Hugo's Homepage"){{< /myshortcode5 >}}

`)

	for i := 1; i <= 30; i++ {
		// Add some content with no shortcodes or links, i.e no templates needed.
		b.WithContent(fmt.Sprintf("blog/notempl%d.md", i), `---
title: No Template
---

## Content
`)
	}
	counters := &testCounters{}
	b.Build(BuildCfg{testCounters: counters})
	b.Assert(int(counters.contentRenderCounter), qt.Equals, 43)

	b.AssertFileContent("public/blog/p1/index.html", `
<p>Cool Page|https://www.google.com|Title: Google's Homepage|Text: First Link|END</p>
Text: Second
SHORT3|
<p>IMAGE: Cool Page||/images/Dragster.jpg|Title: image title|Text: Drag Racing|END</p>
`)

	b.AssertFileContent("public/customview/p1/index.html", `.Render: myrender: Custom View|P4: PARTIAL4`)
	b.AssertFileContent("public/blog/p2/index.html",
		`PARTIAL
,Cool Page2|https://gohugo.io|Title: |Text: No Whitespace Please|END,`,
	)
	b.AssertFileContent("public/blog/p3/index.html", `PARTIAL3`)
	// We may add type template support later, keep this for then. b.AssertFileContent("public/docs/docs1/index.html", `Link docs section: Docs 1|END`)
	b.AssertFileContent("public/blog/p4/index.html", `<p>IMAGE: Cool Page With Image||/images/Dragster.jpg|Title: image title|Text: Drag Racing|END</p>`)
	// The regular markdownify func currently gets regular links.
	b.AssertFileContent("public/blog/p5/index.html", "Inner Link: <a href=\"https://www.google.com\" title=\"Google's Homepage\">Inner Link</a>\n</div>")

	b.AssertFileContent("public/blog/p6/index.html",
		"Inner Inline: Inner Link: With RenderString|https://www.gohugo.io|Title: Hugo's Homepage|Text: Inner Link|END",
		"Inner Block: <p>Inner Link: With RenderString|https://www.gohugo.io|Title: Hugo's Homepage|Text: Inner Link|END</p>",
	)

	b.EditFiles(
		"layouts/_default/_markup/render-link.html", `EDITED: {{ .Destination | safeURL }}|`,
		"layouts/_default/_markup/render-image.html", `IMAGE EDITED: {{ .Destination | safeURL }}|`,
		"layouts/docs/_markup/render-link.html", `DOCS EDITED: {{ .Destination | safeURL }}|`,
		"layouts/partials/mypartial1.html", `PARTIAL1_EDITED`,
		"layouts/partials/mypartial3.html", `PARTIAL3_EDITED`,
		"layouts/partials/mypartial4.html", `PARTIAL4_EDITED`,
		"layouts/shortcodes/myshortcode3.html", `SHORT3_EDITED|`,
	)

	counters = &testCounters{}
	b.Build(BuildCfg{testCounters: counters})
	// Make sure that only content using the changed templates are re-rendered.
	b.Assert(int(counters.contentRenderCounter), qt.Equals, 7)

	b.AssertFileContent("public/customview/p1/index.html", `.Render: myrender: Custom View|P4: PARTIAL4_EDITED`)
	b.AssertFileContent("public/blog/p1/index.html", `<p>EDITED: https://www.google.com|</p>`, "SHORT3_EDITED|")
	b.AssertFileContent("public/blog/p2/index.html", `PARTIAL1_EDITED`)
	b.AssertFileContent("public/blog/p3/index.html", `PARTIAL3_EDITED`)
	// We may add type template support later, keep this for then. b.AssertFileContent("public/docs/docs1/index.html", `DOCS EDITED: https://www.google.com|</p>`)
	b.AssertFileContent("public/blog/p4/index.html", `IMAGE EDITED: /images/Dragster.jpg|`)
	b.AssertFileContent("public/blog/p6/index.html", "<p>Inner Link: EDITED: https://www.gohugo.io|</p>")

}

func TestRenderHooksDeleteTemplate(t *testing.T) {
	config := `
baseURL="https://example.org"
workingDir="/mywork"
`
	b := newTestSitesBuilder(t).WithWorkingDir("/mywork").WithConfigFile("toml", config).Running()
	b.WithTemplatesAdded("_default/single.html", `{{ .Content }}`)
	b.WithTemplatesAdded("_default/_markup/render-link.html", `html-render-link`)

	b.WithContent("p1.md", `---
title: P1
---
[First Link](https://www.google.com "Google's Homepage")

`)
	b.Build(BuildCfg{})

	b.AssertFileContent("public/p1/index.html", `<p>html-render-link</p>`)

	b.RemoveFiles(
		"layouts/_default/_markup/render-link.html",
	)

	b.Build(BuildCfg{})
	b.AssertFileContent("public/p1/index.html", `<p><a href="https://www.google.com" title="Google's Homepage">First Link</a></p>`)

}

func TestRenderHookAddTemplate(t *testing.T) {
	config := `
baseURL="https://example.org"
workingDir="/mywork"
`
	b := newTestSitesBuilder(t).WithWorkingDir("/mywork").WithConfigFile("toml", config).Running()
	b.WithTemplatesAdded("_default/single.html", `{{ .Content }}`)

	b.WithContent("p1.md", `---
title: P1
---
[First Link](https://www.google.com "Google's Homepage")

`)
	b.Build(BuildCfg{})

	b.AssertFileContent("public/p1/index.html", `<p><a href="https://www.google.com" title="Google's Homepage">First Link</a></p>`)

	b.EditFiles("layouts/_default/_markup/render-link.html", `html-render-link`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/p1/index.html", `<p>html-render-link</p>`)

}

func TestRenderHooksRSS(t *testing.T) {

	b := newTestSitesBuilder(t)

	b.WithTemplates("index.html", `
{{ $p := site.GetPage "p1.md" }}

P1: {{ $p.Content }}
	
	`, "index.xml", `

{{ $p2 := site.GetPage "p2.md" }}
{{ $p3 := site.GetPage "p3.md" }}

P2: {{ $p2.Content }}
P3: {{ $p3.Content }}

	
	`,
		"_default/_markup/render-link.html", `html-link: {{ .Destination | safeURL }}|`,
		"_default/_markup/render-link.rss.xml", `xml-link: {{ .Destination | safeURL }}|`,
	)

	b.WithContent("p1.md", `---
title: "p1"
---
P1. [I'm an inline-style link](https://www.gohugo.io)


`, "p2.md", `---
title: "p2"
---
P1. [I'm an inline-style link](https://www.bep.is)


`,
		"p3.md", `---
title: "p2"
outputs: ["rss"]
---
P3. [I'm an inline-style link](https://www.example.org)

`,
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "P1: <p>P1. html-link: https://www.gohugo.io|</p>")
	b.AssertFileContent("public/index.xml", `
P2: <p>P1. xml-link: https://www.bep.is|</p>
P3: <p>P3. xml-link: https://www.example.org|</p>
`)

}

// https://github.com/gohugoio/hugo/issues/6629
func TestRenderLinkWithMarkupInText(t *testing.T) {

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `

baseURL="https://example.org"

[markup]
  [markup.goldmark]
    [markup.goldmark.renderer]
      unsafe = true
    
`)

	b.WithTemplates("index.html", `
{{ $p := site.GetPage "p1.md" }}
P1: {{ $p.Content }}

	`,
		"_default/_markup/render-link.html", `html-link: {{ .Destination | safeURL }}|Text: {{ .Text | safeHTML }}|Plain: {{ .PlainText | safeHTML }}`,
		"_default/_markup/render-image.html", `html-image: {{ .Destination | safeURL }}|Text: {{ .Text | safeHTML }}|Plain: {{ .PlainText | safeHTML }}`,
	)

	b.WithContent("p1.md", `---
title: "p1"
---

START: [**should be bold**](https://gohugo.io)END

Some regular **markup**.

Image:

![Hello<br> Goodbye](image.jpg)END

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
  P1: <p>START: html-link: https://gohugo.io|Text: <strong>should be bold</strong>|Plain: should be boldEND</p>
<p>Some regular <strong>markup</strong>.</p>
<p>html-image: image.jpg|Text: Hello<br> Goodbye|Plain: Hello GoodbyeEND</p>
`)

}

func TestRenderString(t *testing.T) {

	b := newTestSitesBuilder(t)

	b.WithTemplates("index.html", `
{{ $p := site.GetPage "p1.md" }}
{{ $optBlock := dict "display" "block" }}
{{ $optOrg := dict "markup" "org" }}
RSTART:{{ "**Bold Markdown**" | $p.RenderString }}:REND
RSTART:{{  "**Bold Block Markdown**" | $p.RenderString  $optBlock }}:REND
RSTART:{{  "/italic org mode/" | $p.RenderString  $optOrg }}:REND


`)

	b.WithContent("p1.md", `---
title: "p1"
---
`,
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
RSTART:<strong>Bold Markdown</strong>:REND
RSTART:<p><strong>Bold Block Markdown</strong></p>
RSTART:<em>italic org mode</em>:REND
`)

}

// https://github.com/gohugoio/hugo/issues/6882
func TestRenderStringOnListPage(t *testing.T) {
	renderStringTempl := `
{{ .RenderString "**Hello**" }}
`
	b := newTestSitesBuilder(t)
	b.WithContent("mysection/p1.md", `FOO`)
	b.WithTemplates(
		"index.html", renderStringTempl,
		"_default/list.html", renderStringTempl,
		"_default/single.html", renderStringTempl,
	)

	b.Build(BuildCfg{})

	for _, filename := range []string{
		"index.html",
		"mysection/index.html",
		"categories/index.html",
		"tags/index.html",
		"mysection/p1/index.html",
	} {
		b.AssertFileContent("public/"+filename, `<strong>Hello</strong>`)
	}

}
