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
	"flag"
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
	NumLangs     int
	RootSections int
	Render       bool
	Shortcodes   bool
	NumTags      int
	TagsPerPage  int
}

func (s siteBuildingBenchmarkConfig) String() string {
	// Make it comma separated with no spaces, so it is both Bash and regexp friendly.
	// To make it a short as possible, we only shows bools when enabled and ints when >= 0 (RootSections > 1)
	sep := ","
	id := s.Frontmatter + sep
	id += fmt.Sprintf("num_langs=%d%s", s.NumLangs, sep)

	if s.RootSections > 1 {
		id += fmt.Sprintf("num_root_sections=%d%s", s.RootSections, sep)
	}
	id += fmt.Sprintf("num_pages=%d%s", s.NumPages, sep)

	if s.NumTags > 0 {
		id += fmt.Sprintf("num_tags=%d%s", s.NumTags, sep)
	}

	if s.TagsPerPage > 0 {
		id += fmt.Sprintf("tags_per_page=%d%s", s.TagsPerPage, sep)
	}

	if s.Shortcodes {
		id += "shortcodes" + sep
	}

	if s.Render {
		id += "render" + sep
	}

	return strings.TrimSuffix(id, sep)
}

var someLangs = []string{"en", "fr", "nn"}

func BenchmarkSiteBuilding(b *testing.B) {
	var (
		// The below represents the full matrix of benchmarks. Big!
		allFrontmatters    = []string{"YAML", "TOML"}
		allNumRootSections = []int{1, 5}
		allNumTags         = []int{0, 1, 10, 20, 50, 100, 500, 1000, 5000}
		allTagsPerPage     = []int{0, 1, 5, 20, 50, 80}
		allNumPages        = []int{1, 10, 100, 500, 1000, 5000, 10000}
		allDoRender        = []bool{false, true}
		allDoShortCodes    = []bool{false, true}
		allNumLangs        = []int{1, 3}
	)

	var runDefault bool

	visitor := func(a *flag.Flag) {
		if a.Name == "test.bench" && len(a.Value.String()) < 40 {
			// The full suite is too big, so fall back to some smaller default if no
			// restriction is set.
			runDefault = true
		}
	}

	flag.Visit(visitor)

	if runDefault {
		allFrontmatters = allFrontmatters[1:]
		allNumRootSections = allNumRootSections[0:2]
		allNumTags = allNumTags[0:2]
		allTagsPerPage = allTagsPerPage[2:3]
		allNumPages = allNumPages[2:5]
		allDoRender = allDoRender[1:2]
		allDoShortCodes = allDoShortCodes[1:2]
	}

	var conf siteBuildingBenchmarkConfig
	for _, numLangs := range allNumLangs {
		conf.NumLangs = numLangs
		for _, frontmatter := range allFrontmatters {
			conf.Frontmatter = frontmatter
			for _, rootSections := range allNumRootSections {
				conf.RootSections = rootSections
				for _, numTags := range allNumTags {
					conf.NumTags = numTags
					for _, tagsPerPage := range allTagsPerPage {
						conf.TagsPerPage = tagsPerPage
						for _, numPages := range allNumPages {
							conf.NumPages = numPages
							for _, render := range allDoRender {
								conf.Render = render
								for _, shortcodes := range allDoShortCodes {
									conf.Shortcodes = shortcodes
									doBenchMarkSiteBuilding(conf, b)
								}
							}
						}
					}
				}
			}
		}
	}
}

func doBenchMarkSiteBuilding(conf siteBuildingBenchmarkConfig, b *testing.B) {
	b.Run(conf.String(), func(b *testing.B) {
		b.StopTimer()
		sites := createHugoBenchmarkSites(b, b.N, conf)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			h := sites[0]

			err := h.Build(BuildCfg{SkipRender: !conf.Render})
			if err != nil {
				b.Fatal(err)
			}

			// Try to help the GC
			sites[0] = nil
			sites = sites[1:]
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
content starts at 4-columns in :smile:.

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
enableEmoji = true

[outputs]
home = [ "HTML" ]
section = [ "HTML" ]
taxonomy = [ "HTML" ]
taxonomyTerm = [ "HTML" ]
page = [ "HTML" ]

[languages]
%s

[Taxonomies]
tag = "tags"
category = "categories"
`

	langConfigTemplate := `
[languages.%s]
languageName = "Lang %s"
weight = %d
`

	langConfig := ""

	for i := 0; i < cfg.NumLangs; i++ {
		langCode := someLangs[i]
		langConfig += fmt.Sprintf(langConfigTemplate, langCode, langCode, i+1)
	}

	siteConfig = fmt.Sprintf(siteConfig, langConfig)

	numTags := cfg.NumTags

	if cfg.TagsPerPage > numTags {
		numTags = cfg.TagsPerPage
	}

	var (
		contentPagesContent [3]string
		tags                = make([]string, numTags)
		pageTemplate        string
	)

	for i := 0; i < numTags; i++ {
		tags[i] = fmt.Sprintf("Hugo %d", i+1)
	}

	var tagsStr string

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
			"layouts/_default/single.html", `Single HTML|{{ .Title }}|{{ .Content }}|{{ partial "myPartial" . }}`,
			"layouts/_default/list.html", `List HTML|{{ .Title }}|{{ .Content }}|GetPage: {{ with .Site.GetPage "page" "sect3/page3.md" }}{{ .Title }}{{ end }}`,
			"layouts/partials/myPartial.html", `Partial: {{ "Hello **world**!" | markdownify }}`,
			"layouts/shortcodes/myShortcode.html", `<p>MyShortcode</p>`)

		fs := th.Fs

		pagesPerSection := cfg.NumPages / cfg.RootSections / cfg.NumLangs
		for li := 0; li < cfg.NumLangs; li++ {
			fileLangCodeID := ""
			if li > 0 {
				fileLangCodeID = "." + someLangs[li] + "."
			}

			for i := 0; i < cfg.RootSections; i++ {
				for j := 0; j < pagesPerSection; j++ {
					var tagsSlice []string

					if numTags > 0 {
						tagsStart := rand.Intn(numTags) - cfg.TagsPerPage
						if tagsStart < 0 {
							tagsStart = 0
						}
						tagsSlice = tags[tagsStart : tagsStart+cfg.TagsPerPage]
					}

					if cfg.Frontmatter == "TOML" {
						pageTemplate = pageTemplateTOML
						tagsStr = "[]"
						if cfg.TagsPerPage > 0 {
							tagsStr = strings.Replace(fmt.Sprintf("%q", tagsSlice), " ", ", ", -1)
						}
					} else {
						// YAML
						pageTemplate = pageTemplateYAML
						for _, tag := range tagsSlice {
							tagsStr += "\n- " + tag
						}
					}

					content := fmt.Sprintf(pageTemplate, fmt.Sprintf("Title%d_%d", i, j), tagsStr, contentPagesContent[rand.Intn(3)])

					contentFilename := fmt.Sprintf("page%d%s.md", j, fileLangCodeID)

					writeSource(b, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), contentFilename), content)
				}

				content := fmt.Sprintf(pageTemplate, fmt.Sprintf("Section %d", i), "[]", contentPagesContent[rand.Intn(3)])
				indexContentFilename := fmt.Sprintf("_index%s.md", fileLangCodeID)
				writeSource(b, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), indexContentFilename), content)
			}
		}

		sites[i] = h
	}

	return sites
}
