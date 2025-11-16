// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"bytes"
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
)

func TestSiteStats(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	files := `
-- hugo.toml --
baseURL = "http://example.com/blog"

defaultContentLanguage = "nn"

[pagination]
pagerSize = 1

[languages]
[languages.nn]
languageName = "Nynorsk"
weight = 1
title = "Hugo på norsk"

[languages.en]
languageName = "English"
weight = 2
title = "Hugo in English"
-- layouts/single.html --
Single|{{ .Title }}|{{ .Content }}
{{ $img1 := resources.Get "myimage1.png" }}
{{ $img1 = $img1.Fit "100x100" }}
-- layouts/list.html --
List|{{ .Title }}|Pages: {{ .Paginator.TotalPages }}|{{ .Content }}
-- layouts/terms.html --
Terms List|{{ .Title }}|{{ .Content }}

`

	pageTemplate := `---
title: "T%d"
tags:
%s
categories:
%s
aliases: [/Ali%d]
---
# Doc
`

	for i := range 2 {
		for j := range 2 {
			pageID := i + j + 1
			files += fmt.Sprintf("\n-- content/p%d.md --\n", pageID)
			files += fmt.Sprintf(pageTemplate, pageID, fmt.Sprintf("- tag%d", j), fmt.Sprintf("- category%d", j), pageID)
		}
	}

	for i := range 5 {
		files += fmt.Sprintf("\n-- assets/myimage%d.png --\niVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", i+1)
	}

	b := Test(t, files)
	h := b.H

	stats := []*helpers.ProcessingStats{
		h.Sites[0].PathSpec.ProcessingStats,
		h.Sites[1].PathSpec.ProcessingStats,
	}

	var buff bytes.Buffer

	helpers.ProcessingStatsTable(&buff, stats...)
	s := buff.String()

	c.Assert(s, qt.Contains, "Pages            │ 19 │  7")
	c.Assert(s, qt.Contains, "Processed images │  1 │")
}

func TestSiteLastmod(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
-- content/_index.md --
---
date: 2023-01-01
---
-- content/posts/_index.md --
---
date: 2023-02-01
---
-- content/posts/post-1.md --
---
date: 2023-03-01
---
-- content/posts/post-2.md --
---
date: 2023-04-01
---
-- layouts/index.html --
site.Lastmod: {{ .Site.Lastmod.Format "2006-01-02" }}
home.Lastmod: {{ site.Home.Lastmod.Format "2006-01-02" }}

`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "site.Lastmod: 2023-04-01\nhome.Lastmod: 2023-01-01")
}
