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

import "testing"

func TestTempT(t *testing.T) {
	config := `
baseURL="https://example.org"
defaultContentLanguageInSubDir=true

[params]
[params.COLORS]
BLUE="nice"

[languages]
[languages.en]
weight=1
[languages.nn]
weight=2
`
	b := newTestSitesBuilder(t).WithConfigFile("toml", config)
	b.WithTemplates(
		"index.html", `
{{ $params := .Site.Params }}
{{ $colors := $params.Colors }}
{{ $blue := $colors.Blue }}
Len: {{ len $params.Colors }}

Params: {{ $colors }}
Site en: {{ site.Language.Lang }}|{{ .Site.Language.Lang }}|Blue: {{ $blue }}`,
		"index.nn.html", `Site nn: {{ site.Language.Lang }}|{{ .Site.Language.Lang }}`)
	b.WithContent("p1.md", "asdf")
	b.Build(BuildCfg{})

	b.AssertFileContent("public/nn/index.html", "Site nn: nn|nn")
	b.AssertFileContent("public/en/index.html", "Blue: nice")

}
func TestRenderHooks(t *testing.T) {
	// TODO1 markdownify
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
<div class="foo">
{{ .Inner | .Page.RenderString }}
</div>
`)

	b.WithTemplatesAdded("partials/mypartial1.html", `PARTIAL1`)
	b.WithTemplatesAdded("partials/mypartial2.html", `PARTIAL2  {{ partial "mypartial3.html" }}`)
	b.WithTemplatesAdded("partials/mypartial3.html", `PARTIAL3`)
	b.WithTemplatesAdded("_default/_markup/render-link.html", `{{ with .Page }}{{ .Title }}{{ end }}|{{ .Destination | safeURL }}|Title: {{ .Title | safeHTML }}|Text: {{ .Text | safeHTML }}|END`)
	b.WithTemplatesAdded("docs/_markup/render-link.html", `Link docs section: {{ .Text | safeHTML }}|END`)
	b.WithTemplatesAdded("_default/_markup/render-image.html", `IMAGE: {{ .Page.Title }}||{{ .Destination | safeURL }}|Title: {{ .Title | safeHTML }}|Text: {{ .Text | safeHTML }}|END`)

	b.WithContent("blog/p1.md", `---
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
	b.Build(BuildCfg{})
	b.AssertFileContent("public/blog/p1/index.html", `
<p>Cool Page|https://www.google.com|Title: Google's Homepage|Text: First Link|END</p>
Text: Second
SHORT3|
<p>IMAGE: Cool Page||/images/Dragster.jpg|Title: image title|Text: Drag Racing|END</p>
`)
	b.AssertFileContent("public/blog/p2/index.html", `PARTIAL`)
	b.AssertFileContent("public/blog/p3/index.html", `PARTIAL3`)
	b.AssertFileContent("public/docs/docs1/index.html", `Link docs section: Docs 1|END`)
	b.AssertFileContent("public/blog/p4/index.html", `<p>IMAGE: Cool Page With Image||/images/Dragster.jpg|Title: image title|Text: Drag Racing|END</p>`)
	// The regular markdownify func currently gets regular links.
	b.AssertFileContent("public/blog/p5/index.html", "Inner Link: <a href=\"https://www.google.com\" title=\"Google's Homepage\">Inner Link</a>\n</div>")

	b.AssertFileContent("public/blog/p6/index.html", "<div class=\"foo\">\n<p>Inner Link: With RenderString|https://www.gohugo.io|Title: Hugo's Homepage|Text: Inner Link|END</p>\n\n</div>")

	b.EditFiles(
		"layouts/_default/_markup/render-link.html", `EDITED: {{ .Destination | safeURL }}|`,
		"layouts/_default/_markup/render-image.html", `IMAGE EDITED: {{ .Destination | safeURL }}|`,
		"layouts/docs/_markup/render-link.html", `DOCS EDITED: {{ .Destination | safeURL }}|`,
		"layouts/partials/mypartial1.html", `PARTIAL1_EDITED`,
		"layouts/partials/mypartial3.html", `PARTIAL3_EDITED`,
		"layouts/shortcodes/myshortcode3.html", `SHORT3_EDITED|`,
	)
	b.Build(BuildCfg{})
	b.AssertFileContent("public/blog/p1/index.html", `<p>EDITED: https://www.google.com|</p>`, "SHORT3_EDITED|")
	b.AssertFileContent("public/blog/p2/index.html", `PARTIAL1_EDITED`)
	b.AssertFileContent("public/blog/p3/index.html", `PARTIAL3_EDITED`)
	b.AssertFileContent("public/docs/docs1/index.html", `DOCS EDITED: https://www.google.com|</p>`)
	b.AssertFileContent("public/blog/p4/index.html", `IMAGE EDITED: /images/Dragster.jpg|`)
	b.AssertFileContent("public/blog/p6/index.html", "<p>Inner Link: EDITED: https://www.gohugo.io|</p>")

}
