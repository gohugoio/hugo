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

package page_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/rst"
)

func TestPageMarkupMethods(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
summaryLength=2
-- content/p1.md --
---
title: "Post 1"
date: "2020-01-01"
---
{{% foo %}}
-- layouts/shortcodes/foo.html --
Two *words*.
{{/* Test that markup scope is set in all relevant constructs. */}}
{{ if eq hugo.Context.MarkupScope "foo" }}

## Heading 1
Sint ad mollit qui Lorem ut occaecat culpa officia. Et consectetur aute voluptate non sit ullamco adipisicing occaecat. Sunt deserunt amet sit ad. Deserunt enim voluptate proident ipsum dolore dolor ut sit velit esse est mollit irure esse. Mollit incididunt veniam laboris magna et excepteur sit duis. Magna adipisicing reprehenderit tempor irure.
### Heading 2
Exercitation quis est consectetur occaecat nostrud. Ullamco aute mollit aliqua est amet. Exercitation ullamco consectetur dolor labore et non irure eu cillum Lorem.
{{ end }}
-- layouts/index.html --
Home.
{{ .Content }}
-- layouts/_default/single.html --
Single.
Page.ContentWithoutSummmary: {{ .ContentWithoutSummary }}|
{{ template "render-scope" (dict "page" . "scope" "main") }}
{{ template "render-scope" (dict "page" . "scope" "foo") }}
{{ define "render-scope" }}
{{ $c := .page.Markup .scope }}
{{ with $c.Render }}
{{ $.scope }}: Content: {{ .Content }}|
 {{ $.scope }}: ContentWithoutSummary: {{ .ContentWithoutSummary }}|
{{ $.scope }}: Plain: {{ .Plain }}|
{{ $.scope }}: PlainWords: {{ .PlainWords }}|
{{ $.scope }}: WordCount: {{ .WordCount }}|
{{ $.scope }}: FuzzyWordCount: {{ .FuzzyWordCount }}|
{{ $.scope }}: ReadingTime: {{ .ReadingTime }}|
{{ $.scope }}: Len: {{ .Len }}|
{{ $.scope }}: Summary: {{ with .Summary }}{{ . }}{{ else }}nil{{ end }}|
{{ end }}
{{ $.scope }}: Fragments: {{ $c.Fragments.Identifiers }}|
{{ end }}



`

	b := hugolib.Test(t, files)

	// Main scope.
	b.AssertFileContent("public/p1/index.html",
		"Page.ContentWithoutSummmary: |",
		"main: Content: <p>Two <em>words</em>.</p>\n|",
		"main: ContentWithoutSummary: |",
		"main: Plain: Two words.\n|",
		"PlainWords: [Two words.]|\nmain: WordCount: 2|\nmain: FuzzyWordCount: 100|\nmain: ReadingTime: 1|",
		"main: Summary: <p>Two <em>words</em>.</p>|\n\nmain: Fragments: []|",
		"main: Len: 27|",
	)

	// Foo scope (has more content).
	b.AssertFileContent("public/p1/index.html",
		"foo: Content: <p>Two <em>words</em>.</p>\n<h2",
		"foo: ContentWithoutSummary: <h2",
		"Plain: Two words.\nHeading 1",
		"PlainWords: [Two words. Heading 1",
		"foo: WordCount: 81|\nfoo: FuzzyWordCount: 100|\nfoo: ReadingTime: 1|\nfoo: Len: 622|",
		"foo: Summary: <p>Two <em>words</em>.</p>|",
		"foo: Fragments: [heading-1 heading-2]|",
	)
}

func TestPageMarkupScope(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "section"]
-- content/p1.md --
---
title: "Post 1"
date: "2020-01-01"
---

# P1

{{< foo >}}

Begin:{{% includerendershortcodes "p2" %}}:End
Begin:{{< includecontent "p3" >}}:End

-- content/p2.md --
---
title: "Post 2"
date: "2020-01-02"
---

# P2
-- content/p3.md --
---
title: "Post 3"
date: "2020-01-03"
---

# P3

{{< foo >}}

-- layouts/index.html --
Home.
{{ with site.GetPage "p1" }}
	{{ with .Markup "home" }}
	 	{{ .Render.Content }}
	{{ end }}
{{ end }}
-- layouts/_default/single.html --
Single.
{{ with .Markup  }}
	{{ with .Render }}
	 	{{ .Content }}
	{{ end }}
{{ end }}
-- layouts/_default/_markup/render-heading.html --
Render heading: title: {{ .Text}} scope: {{ hugo.Context.MarkupScope }}|
-- layouts/shortcodes/foo.html --
Foo scope: {{ hugo.Context.MarkupScope }}|
-- layouts/shortcodes/includerendershortcodes.html --
{{ $p := site.GetPage (.Get 0) }}
includerendershortcodes: {{ hugo.Context.MarkupScope }}|{{ $p.Markup.RenderShortcodes }}|
-- layouts/shortcodes/includecontent.html --
{{ $p := site.GetPage (.Get 0) }}
includecontent: {{ hugo.Context.MarkupScope }}|{{ $p.Markup.Render.Content }}|

`

	b := hugolib.Test(t, files)

	b.AssertFileContentExact("public/p1/index.html", "Render heading: title: P1 scope: |", "Foo scope: |")

	b.AssertFileContentExact("public/index.html",
		"Begin:\nincludecontent: home|Render heading: title: P3 scope: home|Foo scope: home|\n|\n:End",
		"Render heading: title: P1 scope: home|",
		"Foo scope: home|",
		"Begin:\nincluderendershortcodes: home|</p>\nRender heading: title: P2 scope: home|<p>|:End",
	)
}

func TestPageContentWithoutSummary(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
summaryLength=5
-- content/p1.md --
---
title: "Post 1"
date: "2020-01-01"
---
This is summary.
<!--more-->
This is content.
-- content/p2.md --
---
title: "Post 2"
date: "2020-01-01"
---
This is some content about a summary and more.

Another paragraph.

Third paragraph.
-- content/p3.md --
---
title: "Post 3"
date: "2020-01-01"
summary: "This is summary in front matter."
---
This is content.
-- layouts/_default/single.html --
Single.
Page.Summary: {{ .Summary }}|
{{ with .Markup.Render }}
Content: {{ .Content }}|
ContentWithoutSummary: {{ .ContentWithoutSummary }}|
WordCount: {{ .WordCount }}|
FuzzyWordCount: {{ .FuzzyWordCount }}|
{{ with .Summary }}
Summary: {{ . }}|
Summary Type: {{ .Type }}|
Summary Truncated: {{ .Truncated }}|
{{ end }}
{{ end }}

`
	b := hugolib.Test(t, files)

	b.AssertFileContentExact("public/p1/index.html",
		"Content: <p>This is summary.</p>\n<p>This is content.</p>",
		"ContentWithoutSummary: <p>This is content.</p>|",
		"WordCount: 6|",
		"FuzzyWordCount: 100|",
		"Summary: <p>This is summary.</p>|",
		"Summary Type: manual|",
		"Summary Truncated: true|",
	)
	b.AssertFileContent("public/p2/index.html",
		"Summary: <p>This is some content about a summary and more.</p>|",
		"WordCount: 13|",
		"FuzzyWordCount: 100|",
		"Summary Type: auto",
		"Summary Truncated: true",
	)

	b.AssertFileContentExact("public/p3/index.html",
		"Summary: This is summary in front matter.|",
		"ContentWithoutSummary: <p>This is content.</p>\n|",
	)
}

func TestPageMarkupWithoutSummaryRST(t *testing.T) {
	t.Parallel()
	if !rst.Supports() {
		t.Skip("Skip RST test as not supported")
	}

	files := `
-- hugo.toml --
summaryLength=5
[security.exec]
allow = ["rst", "python"]

-- content/p1.rst --
This is a story about a summary and more.

Another paragraph.
-- content/p2.rst --
This is summary.
<!--more-->
This is content.
-- layouts/_default/single.html --
Single.
Page.Summary: {{ .Summary }}|
{{ with .Markup.Render }}
Content: {{ .Content }}|
ContentWithoutSummary: {{ .ContentWithoutSummary }}|
{{ with .Summary }}
Summary: {{ . }}|
Summary Type: {{ .Type }}|
Summary Truncated: {{ .Truncated }}|
{{ end }}
{{ end }}

`

	b := hugolib.Test(t, files)

	// Auto summary.
	b.AssertFileContentExact("public/p1/index.html",
		"Content: <div class=\"document\">\n\n\n<p>This is a story about a summary and more.</p>\n<p>Another paragraph.</p>\n</div>|",
		"Summary: <div class=\"document\">\n\n\n<p>This is a story about a summary and more.</p></div>|\nSummary Type: auto|\nSummary Truncated: true|",
		"ContentWithoutSummary: <div class=\"document\">\n<p>Another paragraph.</p>\n</div>|",
	)

	// Manual summary.
	b.AssertFileContentExact("public/p2/index.html",
		"Content: <div class=\"document\">\n\n\n<p>This is summary.</p>\n<p>This is content.</p>\n</div>|",
		"ContentWithoutSummary: <div class=\"document\"><p>This is content.</p>\n</div>|",
		"Summary: <div class=\"document\">\n\n\n<p>This is summary.</p>\n</div>|\nSummary Type: manual|\nSummary Truncated: true|",
	)
}

func TestPageMarkupWithoutSummaryAsciidoc(t *testing.T) {
	t.Parallel()
	if !asciidocext.Supports() {
		t.Skip("Skip asiidoc test as not supported")
	}

	files := `
-- hugo.toml --
summaryLength=5
[security.exec]
allow = ["asciidoc", "python"]

-- content/p1.ad --
This is a story about a summary and more.

Another paragraph.
-- content/p2.ad --
This is summary.
<!--more-->
This is content.
-- layouts/_default/single.html --
Single.
Page.Summary: {{ .Summary }}|
{{ with .Markup.Render }}
Content: {{ .Content }}|
ContentWithoutSummary: {{ .ContentWithoutSummary }}|
{{ with .Summary }}
Summary: {{ . }}|
Summary Type: {{ .Type }}|
Summary Truncated: {{ .Truncated }}|
{{ end }}
{{ end }}

`

	b := hugolib.Test(t, files)

	// Auto summary.
	b.AssertFileContentExact("public/p1/index.html",
		"Content: <div class=\"paragraph\">\n<p>This is a story about a summary and more.</p>\n</div>\n<div class=\"paragraph\">\n<p>Another paragraph.</p>\n</div>\n|",
		"Summary: <div class=\"paragraph\">\n<p>This is a story about a summary and more.</p>\n</div>|",
		"Summary Type: auto|\nSummary Truncated: true|",
		"ContentWithoutSummary: <div class=\"paragraph\">\n<p>Another paragraph.</p>\n</div>|",
	)

	// Manual summary.
	b.AssertFileContentExact("public/p2/index.html",
		"Content: <div class=\"paragraph\">\n<p>This is summary.</p>\n</div>\n<div class=\"paragraph\">\n<p>This is content.</p>\n</div>|",
		"ContentWithoutSummary: <div class=\"paragraph\">\n<p>This is content.</p>\n</div>|",
		"Summary: <div class=\"paragraph\">\n<p>This is summary.</p>\n</div>|\nSummary Type: manual|\nSummary Truncated: true|",
	)
}
