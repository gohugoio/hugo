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
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue 4926
// Issue 8232
// Issue 12342
func TestHashSignInPermalink(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['section','rss','sitemap','taxonomy']
[permalinks]
s1 = '/:section/:slug'
-- layouts/_default/list.html --
{{ range site.Pages }}{{ .RelPermalink }}|{{ end }}
-- layouts/_default/single.html --
{{ .Title }}
-- content/s1/p1.md --
---
title: p#1
tags: test#tag#
---
-- content/s2/p#2.md --
---
title: p#2
---
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/s1/p#1/index.html", true)
	b.AssertFileExists("public/s2/p#2/index.html", true)
	b.AssertFileExists("public/tags/test#tag#/index.html", true)

	b.AssertFileContentExact("public/index.html", "/|/s1/p%231/|/s2/p%232/|/tags/test%23tag%23/|")
}

// Issues: 13829, 4428, 7497.
func TestMiscPathIssues(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
uglyURLs = false

[outputFormats.print]
isPlainText = true
mediaType = 'text/plain'
path = 'print'

[outputs]
home = ['html','print']
page = ['html','print']
section = ['html','print']
taxonomy = ['html','print']
term = ['html','print']

[taxonomies]
tag = 'tags'
-- content/_index.md --
---
title: home
---
-- content/s1/_index.md --
---
title: s1
---
-- content/s1/p1.md --
---
title: p1
tags: ['red']
---
-- content/tags/_index.md --
---
title: tags
---
-- content/tags/red/_index.md --
---
title: red
---
`

	templates := []string{
		"layouts/home.html",
		"layouts/home.print.txt",
		"layouts/page.html",
		"layouts/page.print.txt",
		"layouts/section.html",
		"layouts/section.print.txt",
		"layouts/taxonomy.html",
		"layouts/taxonomy.print.txt",
		"layouts/term.html",
		"layouts/term.print.txt",
	}

	const code string = "TITLE: {{ .Title }} | AOFRP: {{ range .AlternativeOutputFormats }}{{ .RelPermalink }}{{ end }} | TEMPLATE: {{ templates.Current.Name }}"

	for _, template := range templates {
		filesTemplate += "-- " + template + " --\n" + code + "\n"
	}

	files := filesTemplate

	b := hugolib.Test(t, files)

	// uglyURLs: false, outputFormat: html
	b.AssertFileContent("public/index.html", "TITLE: home | AOFRP: /print/index.txt | TEMPLATE: home.html")
	b.AssertFileContent("public/s1/index.html", "TITLE: s1 | AOFRP: /print/s1/index.txt | TEMPLATE: section.html")
	b.AssertFileContent("public/s1/p1/index.html", "TITLE: p1 | AOFRP: /print/s1/p1/index.txt | TEMPLATE: page.html")
	b.AssertFileContent("public/tags/index.html", "TITLE: tags | AOFRP: /print/tags/index.txt | TEMPLATE: taxonomy.html")
	b.AssertFileContent("public/tags/red/index.html", "TITLE: red | AOFRP: /print/tags/red/index.txt | TEMPLATE: term.html")

	// uglyURLs: false, outputFormat: print
	b.AssertFileContent("public/print/index.txt", "TITLE: home | AOFRP: / | TEMPLATE: home.print.txt")
	b.AssertFileContent("public/print/s1/index.txt", "TITLE: s1 | AOFRP: /s1/ | TEMPLATE: section.print.txt")
	b.AssertFileContent("public/print/s1/p1/index.txt", "TITLE: p1 | AOFRP: /s1/p1/ | TEMPLATE: page.print.txt")
	b.AssertFileContent("public/print/tags/index.txt", "TITLE: tags | AOFRP: /tags/ | TEMPLATE: taxonomy.print.txt")
	b.AssertFileContent("public/print/tags/red/index.txt", "TITLE: red | AOFRP: /tags/red/ | TEMPLATE: term.print.txt")

	files = strings.ReplaceAll(filesTemplate, "uglyURLs = false", "uglyURLs = true")
	b = hugolib.Test(t, files)

	// uglyURLs: true, outputFormat: html
	b.AssertFileContent("public/index.html", "TITLE: home | AOFRP: /print/index.txt | TEMPLATE: home.html")
	b.AssertFileContent("public/s1/index.html", "TITLE: s1 | AOFRP: /print/s1/index.txt | TEMPLATE: section.html")
	b.AssertFileContent("public/s1/p1.html", "TITLE: p1 | AOFRP: /print/s1/p1.txt | TEMPLATE: page.html")
	b.AssertFileContent("public/tags/index.html", "TITLE: tags | AOFRP: /print/tags/index.txt | TEMPLATE: taxonomy.html")
	b.AssertFileContent("public/tags/red.html", "TITLE: red | AOFRP: /print/tags/red.txt | TEMPLATE: term.html")

	// uglyURLs: true, outputFormat: print
	b.AssertFileContent("public/print/index.txt", "TITLE: home | AOFRP: /index.html | TEMPLATE: home.print.txt")
	b.AssertFileContent("public/print/s1/index.txt", "TITLE: s1 | AOFRP: /s1/index.html | TEMPLATE: section.print.txt")
	b.AssertFileContent("public/print/s1/p1.txt", "TITLE: p1 | AOFRP: /s1/p1.html | TEMPLATE: page.print.txt")
	b.AssertFileContent("public/print/tags/index.txt", "TITLE: tags | AOFRP: /tags/index.html | TEMPLATE: taxonomy.print.txt")
	b.AssertFileContent("public/print/tags/red.txt", "TITLE: red | AOFRP: /tags/red.html | TEMPLATE: term.print.txt")
}
