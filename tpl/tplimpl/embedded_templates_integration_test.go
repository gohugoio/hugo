// Copyright 2026 The Hugo Authors. All rights reserved.
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

package tplimpl_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

func TestDisqusTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','sitemap','taxonomy','term']
[services.disqus]
shortname = 'foo'
[privacy.disqus]
disable = false
-- layouts/home.html --
{{ .Title }}
{{ partial "disqus.html" . }}
-- content/_index.md --
+++
title = 'Home'
[params]
disqus_identifier = 'my-identifier'
disqus_title = 'My Title'
disqus_url = 'https://example.com/my-page/'
+++
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		`s.src = '//' + "foo" + '.disqus.com/embed.js';`,
		`this.page.identifier = 'my-identifier';`,
		`this.page.title = 'My Title';`,
		`this.page.url = 'https:\/\/example.com\/my-page\/';`,
	)

	// disable = true: no output
	f := strings.ReplaceAll(files, "disable = false", "disable = true")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", `! disqus_thread`)

	// No shortname: no output
	f = strings.ReplaceAll(files, "shortname = 'foo'", "shortname = ''")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", `! disqus_thread`)
}

func TestGoogleAnalyticsTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','sitemap','taxonomy','term']
[privacy.googleAnalytics]
disable = false
DNT
[services.googleAnalytics]
id = 'G-0123456789'
-- layouts/home.html --
{{ .Title }}
{{ partial "google_analytics.html" . }}
-- content/_index.md --
---
title: Home
---
`

	// Default respectDoNotTrack value
	f := strings.ReplaceAll(files, "DNT", "")
	b := hugolib.Test(t, f)
	b.AssertFileContent("public/index.html",
		`<script async src="https://www.googletagmanager.com/gtag/js?id=G-0123456789"></script>`,
		`if ( true )`,
	)

	// respectDoNotTrack = true
	f = strings.ReplaceAll(files, "DNT", "respectDoNotTrack = true")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html",
		`<script async src="https://www.googletagmanager.com/gtag/js?id=G-0123456789"></script>`,
		`if ( true )`,
	)

	// respectDoNotTrack = false
	f = strings.ReplaceAll(files, "DNT", "respectDoNotTrack = false")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html",
		`<script async src="https://www.googletagmanager.com/gtag/js?id=G-0123456789"></script>`,
		`if ( false )`,
	)

	// disable = true: no output
	f = strings.ReplaceAll(files, "DNT", "")
	f = strings.ReplaceAll(f, "disable = false", "disable = true")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", `! googletagmanager`)

	// No ID: no output
	f = strings.ReplaceAll(files, "DNT", "")
	f = strings.ReplaceAll(f, "id = 'G-0123456789'", "id = ''")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", `! googletagmanager`)

	// UA- prefix ID: warning logged, no script output
	f = strings.ReplaceAll(files, "DNT", "")
	f = strings.ReplaceAll(f, "G-0123456789", "UA-12345678-1")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", `! googletagmanager`)
}

func TestOpengraphTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
title = 'My Site'
capitalizeListTitles = false
disableKinds = ['section','rss','sitemap','taxonomy']
locale = 'en-US'
[markup.goldmark.renderer]
unsafe = true
[params]
description = "m <em>n</em> and **o** can't."
[params.social]
facebook_admin = 'foo'
[taxonomies]
series = 'series'
tag = 'tags'
-- layouts/home.html --
{{ partial "opengraph.html" . }}
-- layouts/page.html --
{{ partial "opengraph.html" . }}
-- layouts/term.html --
{{ partial "opengraph.html" . }}
-- content/_index.md --
---
title: Home
---
-- content/s1/p1.md --
---
title: p1
date: 2024-04-24T08:00:00-07:00
lastmod: 2024-04-24T11:00:00-07:00
images: [a.jpg,b.jpg]
audio: [c.mp3,d.mp3]
videos: [e.mp4,f.mp4]
series: [series-1]
tags: [t1,t2]
---
a <em>b</em> and **c** can't.
-- content/s1/p2.md --
---
title: p2
series: [series-1]
---
d <em>e</em> and **f** can't.
<!--more-->
-- content/s1/p3.md --
---
title: p3
series: [series-1]
summary: g <em>h</em> and **i** can't.
---
-- content/s1/p4.md --
---
title: p4
series: [series-1]
description: j <em>k</em> and **l** can't.
---
-- content/s1/p5.md --
---
title: p5
series: [series-1]
---
-- content/s1/p6.md --
---
title: p6
series: [series-1]
locale: fr-FR
---
-- content/s1/p7.md --
---
title: ''
---
-- content/p0.md --
---
title: p0
---
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/p1/index.html", `
		<meta property="og:url" content="/s1/p1/">
		<meta property="og:site_name" content="My Site">
		<meta property="og:title" content="p1">
		<meta property="og:description" content="a b and c can’t.">
		<meta property="og:locale" content="en_US">
		<meta property="og:type" content="article">
		<meta property="article:section" content="s1">
		<meta property="article:published_time" content="2024-04-24T08:00:00-07:00">
		<meta property="article:modified_time" content="2024-04-24T11:00:00-07:00">
		<meta property="article:tag" content="t1">
		<meta property="article:tag" content="t2">
		<meta property="og:image" content="/a.jpg">
		<meta property="og:image" content="/b.jpg">
		<meta property="og:audio" content="/c.mp3">
		<meta property="og:audio" content="/d.mp3">
		<meta property="og:video" content="/e.mp4">
		<meta property="og:video" content="/f.mp4">
		<meta property="og:see_also" content="/s1/p2/">
		<meta property="og:see_also" content="/s1/p3/">
		<meta property="og:see_also" content="/s1/p4/">
		<meta property="og:see_also" content="/s1/p5/">
		<meta property="fb:admins" content="foo">
		`,
	)

	// Description from content excerpt
	b.AssertFileContent("public/s1/p2/index.html",
		`<meta property="og:description" content="d e and f can’t.">`,
	)

	// Description from front matter summary
	b.AssertFileContent("public/s1/p3/index.html",
		`<meta property="og:description" content="g h and i can’t.">`,
	)

	// The markdown is intentionally not rendered to HTML.
	b.AssertFileContent("public/s1/p4/index.html",
		`<meta property="og:description" content="j k and **l** can&#39;t.">`,
	)

	// The markdown is intentionally not rendered to HTML.
	b.AssertFileContent("public/s1/p5/index.html",
		`<meta property="og:description" content="m n and **o** can&#39;t.">`,
	)

	// og:locale from page front matter
	b.AssertFileContent("public/s1/p6/index.html",
		`<meta property="og:locale" content="fr_FR">`,
	)

	// No article:section for root-level page
	b.AssertFileContent("public/p0/index.html",
		`<meta property="og:type" content="article">`,
		`! article:section`,
	)

	// Term page
	b.AssertFileContent("public/series/series-1/index.html",
		`<meta property="og:url" content="/series/series-1/">`,
		`<meta property="og:title" content="series-1">`,
		`<meta property="og:type" content="website">`,
	)

	// site.Title fallback when page title is empty
	b.AssertFileContent("public/s1/p7/index.html", `<meta property="og:title" content="My Site">`)

	// facebook_app_id takes precedence over facebook_admin
	f := strings.ReplaceAll(files, "facebook_admin = 'foo'", "facebook_app_id = 'bar'")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p1/index.html",
		`<meta property="fb:app_id" content="bar">`,
		`! fb:admins`,
	)

	// Issue 14433
	b.AssertFileContent("public/index.html",
		`<meta property="og:title" content="Home">`,
		`<meta property="og:locale" content="en_US">`,
		`<meta property="og:type" content="website">`,
		`! false`,
	)

	// No description: no meta
	f = strings.ReplaceAll(files, `description = "m <em>n</em> and **o** can't."`, "")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p5/index.html", `! og:description`)
}

func TestPaginationTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
pagination.pagerSize = 1
-- layouts/home.html --
{{ .Paginate (where site.RegularPages "Section" "s1") }}PARTIAL
-- content/s1/p01.md --
---
title: p01
---
-- content/s1/p02.md --
---
title: p02
---
-- content/s1/p03.md --
---
title: p03
---
-- content/s1/p04.md --
---
title: p04
---
-- content/s1/p05.md --
---
title: p05
---
-- content/s1/p06.md --
---
title: p06
---
-- content/s1/p07.md --
---
title: p07
---
-- content/s1/p08.md --
---
title: p08
---
-- content/s1/p09.md --
---
title: p09
---
-- content/s1/p10.md --
---
title: p10
---
`

	test := func(variant string, expectedOutput string) {
		b := hugolib.Test(t, strings.ReplaceAll(files, "PARTIAL", variant))
		b.AssertFileContent("public/index.html", expectedOutput)
	}

	expectedOutputDefaultFormat := "Pager 1\n    <ul class=\"pagination pagination-default\">\n      <li class=\"page-item disabled\">\n        <a aria-disabled=\"true\" aria-label=\"First\" class=\"page-link\" role=\"button\" tabindex=\"-1\"><span aria-hidden=\"true\">&laquo;&laquo;</span></a>\n      </li>\n      <li class=\"page-item disabled\">\n        <a aria-disabled=\"true\" aria-label=\"Previous\" class=\"page-link\" role=\"button\" tabindex=\"-1\"><span aria-hidden=\"true\">&laquo;</span></a>\n      </li>\n      <li class=\"page-item active\">\n        <a aria-current=\"page\" aria-label=\"Page 1\" class=\"page-link\" role=\"button\">1</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Page 2\" class=\"page-link\" role=\"button\">2</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/3/\" aria-label=\"Page 3\" class=\"page-link\" role=\"button\">3</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/4/\" aria-label=\"Page 4\" class=\"page-link\" role=\"button\">4</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/5/\" aria-label=\"Page 5\" class=\"page-link\" role=\"button\">5</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Next\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;</span></a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/10/\" aria-label=\"Last\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;&raquo;</span></a>\n      </li>\n    </ul>"
	expectedOutputTerseFormat := "Pager 1\n    <ul class=\"pagination pagination-terse\">\n      <li class=\"page-item active\">\n        <a aria-current=\"page\" aria-label=\"Page 1\" class=\"page-link\" role=\"button\">1</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Page 2\" class=\"page-link\" role=\"button\">2</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/3/\" aria-label=\"Page 3\" class=\"page-link\" role=\"button\">3</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Next\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;</span></a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/10/\" aria-label=\"Last\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;&raquo;</span></a>\n      </li>\n    </ul>"

	test(`{{ partial "pagination.html" . }}`, expectedOutputDefaultFormat)
	test(`{{ partial "pagination.html" (dict "page" .) }}`, expectedOutputDefaultFormat)
	test(`{{ partial "pagination.html" (dict "page" . "format" "default") }}`, expectedOutputDefaultFormat)
	test(`{{ partial "pagination.html" (dict "page" . "format" "terse") }}`, expectedOutputTerseFormat)

	// Default format, last page: active First and Prev, disabled Next and Last
	f := strings.ReplaceAll(files, "PARTIAL", `{{ partial "pagination.html" . }}`)
	b := hugolib.Test(t, f)
	b.AssertFileContent("public/page/10/index.html",
		`<a href="/" aria-label="First" class="page-link" role="button">`,
		`<a href="/page/9/" aria-label="Previous" class="page-link" role="button">`,
		`aria-disabled="true" aria-label="Next"`,
		`aria-disabled="true" aria-label="Last"`,
	)

	// Terse format, last page: First and Prev shown, Next and Last hidden
	f = strings.ReplaceAll(files, "PARTIAL", `{{ partial "pagination.html" (dict "page" . "format" "terse") }}`)
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/page/10/index.html",
		`<a href="/" aria-label="First" class="page-link" role="button">`,
		`<a href="/page/9/" aria-label="Previous" class="page-link" role="button">`,
		`! aria-label="Next"`,
		`! aria-label="Last"`,
	)

	// TotalPages <= 1: no output
	f = `
-- hugo.toml --
-- layouts/home.html --
{{ .Paginate (where site.RegularPages "Section" "s1") }}{{ partial "pagination.html" . }}
-- content/s1/p1.md --
---
title: p1
---
`
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", `! class="pagination"`)

	// Unsupported format: build error
	f = strings.ReplaceAll(files, "PARTIAL", `{{ partial "pagination.html" (dict "page" . "format" "fancy") }}`)
	b, err := hugolib.TestE(t, f)
	b.Assert(err, qt.IsNotNil)

	// Map without page key: build error
	f = strings.ReplaceAll(files, "PARTIAL", `{{ partial "pagination.html" (dict "format" "default") }}`)
	b, err = hugolib.TestE(t, f)
	b.Assert(err, qt.IsNotNil)
}

// Issue 12432
func TestSchemaTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
capitalizeListTitles = false
disableKinds = ['rss','sitemap']
[markup.goldmark.renderer]
unsafe = true
[params]
description = "m <em>n</em> and **o** can't."
[taxonomies]
tag = 'tags'
-- layouts/list.html --
{{ partial "schema.html" . }}
-- layouts/single.html --
{{ partial "schema.html" . }}
-- content/s1/p1.md --
---
title: p1
date: 2024-04-24T08:00:00-07:00
lastmod: 2024-04-24T11:00:00-07:00
images: [a.jpg,b.jpg]
tags: [t1,t2]
---
a <em>b</em> and **c** can't.
-- content/s1/p2.md --
---
title: p2
---
d <em>e</em> and **f** can't.
<!--more-->
-- content/s1/p3.md --
---
title: p3
summary: g <em>h</em> and **i** can't.
---
-- content/s1/p4.md --
---
title: p4
description: j <em>k</em> and **l** can't.
---
-- content/s1/p5.md --
---
title: p5
---
-- content/s1/p6.md --
---
title: p6
keywords: [k1,k2]
---
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/p1/index.html", `
		<meta itemprop="name" content="p1">
		<meta itemprop="description" content="a b and c can’t.">
		<meta itemprop="datePublished" content="2024-04-24T08:00:00-07:00">
		<meta itemprop="dateModified" content="2024-04-24T11:00:00-07:00">
		<meta itemprop="wordCount" content="5">
		<meta itemprop="image" content="/a.jpg">
		<meta itemprop="image" content="/b.jpg">
		<meta itemprop="keywords" content="t1,t2">
  		`,
	)

	// Description from content excerpt
	b.AssertFileContent("public/s1/p2/index.html",
		`<meta itemprop="description" content="d e and f can’t.">`,
	)

	// Description from front matter summary
	b.AssertFileContent("public/s1/p3/index.html",
		`<meta itemprop="description" content="g h and i can’t.">`,
	)

	// The markdown is intentionally not rendered to HTML.
	b.AssertFileContent("public/s1/p4/index.html",
		`<meta itemprop="description" content="j k and **l** can&#39;t.">`,
	)

	// The markdown is intentionally not rendered to HTML.
	b.AssertFileContent("public/s1/p5/index.html",
		`<meta itemprop="description" content="m n and **o** can&#39;t.">`,
	)

	// Front matter keywords
	b.AssertFileContent("public/s1/p6/index.html",
		`<meta itemprop="keywords" content="k1,k2">`,
	)

	// No keywords: no meta
	b.AssertFileContent("public/s1/p5/index.html", `! itemprop="keywords"`)

	// Keywords taxonomy
	f := strings.ReplaceAll(files, "tag = 'tags'", "tag = 'tags'\nkeyword = 'keywords'")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p6/index.html", `<meta itemprop="keywords" content="k1,k2">`)

	// Any taxonomy fallback
	f = strings.ReplaceAll(files, "tag = 'tags'", "series = 'series'")
	f = strings.ReplaceAll(f, "tags: [t1,t2]", "series: [s1]")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p1/index.html", `<meta itemprop="keywords" content="s1">`)

	// No description: no meta
	f = strings.ReplaceAll(files, `description = "m <em>n</em> and **o** can't."`, "")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p5/index.html", `! itemprop="description"`)

	// name falls back to site.Title when page title is empty
	f = `
-- hugo.toml --
title = 'My Site'
disableKinds = ['rss','sitemap']
-- layouts/single.html --
{{ partial "schema.html" . }}
-- content/p1.md --
---
title: ''
---
`
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/p1/index.html", `<meta itemprop="name" content="My Site">`)
}

// Issue 12433
func TestTwitterCardsTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
capitalizeListTitles = false
disableKinds = ['rss','sitemap','taxonomy','term']
[markup.goldmark.renderer]
unsafe = true
[params]
description = "m <em>n</em> and **o** can't."
[params.social]
twitter = 'foo'
-- layouts/list.html --
{{ partial "twitter_cards.html" . }}
-- layouts/single.html --
{{ partial "twitter_cards.html" . }}
-- content/s1/p1.md --
---
title: p1
images: [a.jpg,b.jpg]
---
a <em>b</em> and **c** can't.
-- content/s1/p2.md --
---
title: p2
---
d <em>e</em> and **f** can't.
<!--more-->
-- content/s1/p3.md --
---
title: p3
summary: g <em>h</em> and **i** can't.
---
-- content/s1/p4.md --
---
title: p4
description: j <em>k</em> and **l** can't.
---
-- content/s1/p5.md --
---
title: p5
---
-- content/s1/p6.md --
---
title: ''
---
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/p1/index.html", `
		<meta name="twitter:card" content="summary_large_image">
		<meta name="twitter:image" content="/a.jpg">
		<meta name="twitter:title" content="p1">
		<meta name="twitter:description" content="a b and c can’t.">
		<meta name="twitter:site" content="@foo">
		`,
	)

	// Description from content excerpt
	b.AssertFileContent("public/s1/p2/index.html",
		`<meta name="twitter:card" content="summary">`,
		`<meta name="twitter:description" content="d e and f can’t.">`,
	)

	// Description from front matter summary
	b.AssertFileContent("public/s1/p3/index.html",
		`<meta name="twitter:description" content="g h and i can’t.">`,
	)

	// The markdown is intentionally not rendered to HTML.
	b.AssertFileContent("public/s1/p4/index.html",
		`<meta name="twitter:description" content="j k and **l** can&#39;t.">`,
	)

	// The markdown is intentionally not rendered to HTML.
	b.AssertFileContent("public/s1/p5/index.html",
		`<meta name="twitter:description" content="m n and **o** can&#39;t.">`,
	)

	// No title: no meta
	b.AssertFileContent("public/s1/p6/index.html", `! twitter:title`)

	// site.Title fallback when page title is empty
	f := strings.ReplaceAll(files, "capitalizeListTitles = false", "title = 'My Site'\ncapitalizeListTitles = false")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p6/index.html", `<meta name="twitter:title" content="My Site">`)

	// site.Params.title fallback when page title and site title are empty
	f = strings.ReplaceAll(files, "[params]\n", "[params]\ntitle = 'My Params Site'\n")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p6/index.html", `<meta name="twitter:title" content="My Params Site">`)

	// No description: no meta
	f = strings.ReplaceAll(files, `description = "m <em>n</em> and **o** can't."`, "")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p6/index.html", `! twitter:description`)

	// @-prefixed twitter handle: no double @
	f = strings.ReplaceAll(files, "twitter = 'foo'", "twitter = '@foo'")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p1/index.html", `<meta name="twitter:site" content="@foo">`)

	// site.Params.social not a map: no twitter:site
	f = strings.ReplaceAll(files, "[params.social]\ntwitter = 'foo'", "social = 'foo'")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p1/index.html", `! twitter:site`)

	// twitter handle empty: no twitter:site
	f = strings.ReplaceAll(files, "twitter = 'foo'", "twitter = ''")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/s1/p1/index.html", `! twitter:site`)
}

func TestGetPageImagesTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
[params]
images = ["a.jpg", "b.jpg"]
-- assets/c.jpg --
Yy5qcGc=
-- layouts/single.html --
{{ .Title}}
{{ $images := partial "_funcs/get-page-images.html" . }}
{{ range $images }}
 	{{ .Permalink }}
{{ end }}
-- content/p1/index.md --
---
title: p1
---
-- content/p1/d-featured.jpg --
ZC1mZWF0dXJlZC5qcGc=
-- content/p2/index.md --
---
title: p2
---
-- content/p2/e-cover.jpg --
ZS1jb3Zlci5qcGc=
-- content/p3/index.md --
---
title: p3
---
-- content/p3/f-thumbnail.jpg --
Zi10aHVtYm5haWwuanBn
-- content/p4.md --
---
title: p4
images: ["g.jpg"]
---
-- content/p5.md --
---
title: p5
images: ["https://example.com/h.png"]
---
-- content/p6/index.md --
---
title: p6
images: ["i.jpg"]
---
-- content/p6/i.jpg --
aS5qcGc=
-- content/p7.md --
---
title: p7
images: ["c.jpg"]
---
-- content/p8.md --
---
title: p8
---
`

	b := hugolib.Test(t, files)

	// Featured image resource (*feature* glob match)
	b.AssertFileContent("public/p1/index.html", `https://example.org/p1/d-featured.jpg`)
	b.AssertFileContent("public/p1/d-featured.jpg", "d-featured.jpg")

	// Cover image resource (*cover* glob match)
	b.AssertFileContent("public/p2/index.html", `https://example.org/p2/e-cover.jpg`)
	b.AssertFileContent("public/p2/e-cover.jpg", "e-cover.jpg")

	// Thumbnail image resource (*thumbnail* glob match)
	b.AssertFileContent("public/p3/index.html", `https://example.org/p3/f-thumbnail.jpg`)
	b.AssertFileContent("public/p3/f-thumbnail.jpg", "f-thumbnail.jpg")

	// Explicit images param: non-resource relative path (absURL fallback)
	b.AssertFileContent("public/p4/index.html", `https://example.org/g.jpg`)

	// Explicit images param: external URL
	b.AssertFileContent("public/p5/index.html", `https://example.com/h.png`)

	// Explicit images param: page resource
	b.AssertFileContent("public/p6/index.html", `https://example.org/p6/i.jpg`)
	b.AssertFileContent("public/p6/i.jpg", "i.jpg")

	// Explicit images param: global resource
	b.AssertFileContent("public/p7/index.html", `https://example.org/c.jpg`)
	b.AssertFileContent("public/c.jpg", "c.jpg")

	// Site params fallback: only the first site image is used
	b.AssertFileContent("public/p8/index.html", `https://example.org/a.jpg`, `! b.jpg`)

	// Site params fallback with no site images: no output
	f := strings.ReplaceAll(files, `images = ["a.jpg", "b.jpg"]`, `images = []`)
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/p8/index.html", `! https://`)
}
