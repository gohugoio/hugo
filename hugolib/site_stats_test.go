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
	"io/ioutil"
	"testing"

	"github.com/gohugoio/hugo/helpers"

	qt "github.com/frankban/quicktest"
)

func TestSiteStats(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 1
defaultContentLanguage = "nn"

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
		h.Sites[1].PathSpec.ProcessingStats}

	stats[0].Table(ioutil.Discard)
	stats[1].Table(ioutil.Discard)

	var buff bytes.Buffer

	helpers.ProcessingStatsTable(&buff, stats...)

	c.Assert(buff.String(), qt.Contains, "Pages            | 19 |  6")

}
