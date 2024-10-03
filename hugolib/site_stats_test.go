// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"io"
	"testing"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
)

func TestSiteStats(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	siteConfig := `
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

	b := newTestSitesBuilder(t).WithConfigFile("toml", siteConfig)

	b.WithTemplates(
		"_default/single.html", "Single|{{ .Title }}|{{ .Content }}",
		"_default/list.html", `List|{{ .Title }}|Pages: {{ .Paginator.TotalPages }}|{{ .Content }}`,
		"_default/terms.html", "Terms List|{{ .Title }}|{{ .Content }}",
	)

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			pageID := i + j + 1
			b.WithContent(fmt.Sprintf("content/sect/p%d.md", pageID),
				fmt.Sprintf(pageTemplate, pageID, fmt.Sprintf("- tag%d", j), fmt.Sprintf("- category%d", j), pageID))
		}
	}

	for i := 0; i < 5; i++ {
		b.WithContent(fmt.Sprintf("assets/image%d.png", i+1), "image")
	}

	b.Build(BuildCfg{})
	h := b.H

	stats := []*helpers.ProcessingStats{
		h.Sites[0].PathSpec.ProcessingStats,
		h.Sites[1].PathSpec.ProcessingStats,
	}

	stats[0].Table(io.Discard)
	stats[1].Table(io.Discard)

	var buff bytes.Buffer

	helpers.ProcessingStatsTable(&buff, stats...)

	c.Assert(buff.String(), qt.Contains, "Pages            | 21 |  7")
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
