// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

type siteBuildingBenchmarkConfig struct {
	Frontmatter  string
	NumPages     int
	RootSections int
	Render       bool
	Shortcodes   bool
	TagsPerPage  int
}

func (s siteBuildingBenchmarkConfig) String() string {
	// Make it comma separated with no spaces, so it is both Bash and regexp friendly.
	return fmt.Sprintf("frontmatter=%s,num_root_sections=%d,num_pages=%d,tags_per_page=%d,shortcodes=%t,render=%t", s.Frontmatter, s.RootSections, s.NumPages, s.TagsPerPage, s.Shortcodes, s.Render)
}

func BenchmarkSiteBuilding(b *testing.B) {
	var conf siteBuildingBenchmarkConfig
	for _, frontmatter := range []string{"YAML", "TOML"} {
		conf.Frontmatter = frontmatter
		for _, rootSections := range []int{1, 5} {
			conf.RootSections = rootSections
			for _, tagsPerPage := range []int{0, 1, 5, 20} {
				conf.TagsPerPage = tagsPerPage
				for _, numPages := range []int{10, 100, 500, 1000, 5000, 10000} {
					conf.NumPages = numPages
					for _, render := range []bool{false, true} {
						conf.Render = render
						for _, shortcodes := range []bool{false, true} {
							conf.Shortcodes = shortcodes
							doBenchMarkSiteBuilding(conf, b)
						}
					}
				}
			}
		}
	}
}

func doBenchMarkSiteBuilding(conf siteBuildingBenchmarkConfig, b *testing.B) {
	b.Run(conf.String(), func(b *testing.B) {
		sites := createHugoBenchmarkSites(b, b.N, conf)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := sites[0]

			err := h.Build(BuildCfg{SkipRender: !conf.Render})
			if err != nil {
				b.Fatal(err)
			}

			// Try to help the GC
			sites[0] = nil
			sites = sites[1:len(sites)]
		}
	})
}

func createHugoBenchmarkSites(b *testing.B, count int, cfg siteBuildingBenchmarkConfig) []*HugoSites {
	someMarkdown := `
An h1 header
============

Paragraphs are separated by a blank line.

2nd paragraph. *Italic* and **bold**. Itemized lists
look like:

  * this one
  * that one
  * the other one

Note that --- not considering the asterisk --- the actual text
content starts at 4-columns in.

> Block quotes are
> written like so.
>
> They can span multiple paragraphs,
> if you like.

Use 3 dashes for an em-dash. Use 2 dashes for ranges (ex., "it's all
in chapters 12--14"). Three dots ... will be converted to an ellipsis.
Unicode is supported. ☺
`

	someMarkdownWithShortCode := someMarkdown + `

{{< myShortcode >}}

`

	pageTemplateTOML := `+++
title = "%s"
tags = %s
+++
%s

`

	pageTemplateYAML := `---
title: "%s"
tags:
%s
---
%s

`

	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 10
defaultContentLanguage = "en"

[Taxonomies]
tag = "tags"
category = "categories"
`
	var (
		contentPagesContent [3]string
		tags                = make([]string, cfg.TagsPerPage)
		pageTemplate        string
	)

	tagOffset := rand.Intn(10)

	for i := 0; i < len(tags); i++ {
		tags[i] = fmt.Sprintf("Hugo %d", i+tagOffset)
	}

	var tagsStr string

	if cfg.Frontmatter == "TOML" {
		pageTemplate = pageTemplateTOML
		tagsStr = "[]"
		if cfg.TagsPerPage > 0 {
			tagsStr = strings.Replace(fmt.Sprintf("%q", tags[0:cfg.TagsPerPage]), " ", ", ", -1)
		}
	} else {
		// YAML
		pageTemplate = pageTemplateYAML
		for _, tag := range tags {
			tagsStr += "\n- " + tag
		}
	}

	if cfg.Shortcodes {
		contentPagesContent = [3]string{
			someMarkdownWithShortCode,
			strings.Repeat(someMarkdownWithShortCode, 2),
			strings.Repeat(someMarkdownWithShortCode, 3),
		}
	} else {
		contentPagesContent = [3]string{
			someMarkdown,
			strings.Repeat(someMarkdown, 2),
			strings.Repeat(someMarkdown, 3),
		}
	}

	sites := make([]*HugoSites, count)
	for i := 0; i < count; i++ {
		// Maybe consider reusing the Source fs
		mf := afero.NewMemMapFs()
		th, h := newTestSitesFromConfig(b, mf, siteConfig,
			"layouts/_default/single.html", `Single HTML|{{ .Title }}|{{ .Content }}`,
			"layouts/_default/list.html", `List HTML|{{ .Title }}|{{ .Content }}`,
			"layouts/shortcodes/myShortcode.html", `<p>MyShortcode</p>`)

		fs := th.Fs

		pagesPerSection := cfg.NumPages / cfg.RootSections

		for i := 0; i < cfg.RootSections; i++ {
			for j := 0; j < pagesPerSection; j++ {
				content := fmt.Sprintf(pageTemplate, fmt.Sprintf("Title%d_%d", i, j), tagsStr, contentPagesContent[rand.Intn(3)])

				writeSource(b, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", j)), content)
			}
		}

		sites[i] = h
	}

	return sites
}
