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

import (
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

type siteBenchmarkTestcase struct {
	name   string
	create func(t testing.TB) *sitesBuilder
	check  func(s *sitesBuilder)
}

func getBenchmarkSiteNewTestCases() []siteBenchmarkTestcase {

	const markdownSnippets = `

## Links


This is [an example](http://example.com/ "Title") inline link.

[This link](http://example.net/) has no title attribute.

This is [Relative](/all-is-relative).

See my [About](/about/) page for details. 
`

	pageContentForMarkdown := func(size int, markdown string) string {
		return `---
title: "My Page"
---

My page content.

` + strings.Repeat(markdown, size)
	}

	pageContent := func(size int) string {
		return pageContentForMarkdown(size, markdownSnippets)
	}

	config := `
baseURL = "https://example.com"
`

	benchmarks := []siteBenchmarkTestcase{
		{"Bundle with image", func(b testing.TB) *sitesBuilder {
			sb := newTestSitesBuilder(b).WithConfigFile("toml", config)
			sb.WithContent("content/blog/mybundle/index.md", pageContent(1))
			sb.WithSunset("content/blog/mybundle/sunset1.jpg")

			return sb
		},
			func(s *sitesBuilder) {
				s.AssertFileContent("public/blog/mybundle/index.html", "/blog/mybundle/sunset1.jpg")
				s.CheckExists("public/blog/mybundle/sunset1.jpg")

			},
		},
		{"Bundle with JSON file", func(b testing.TB) *sitesBuilder {
			sb := newTestSitesBuilder(b).WithConfigFile("toml", config)
			sb.WithContent("content/blog/mybundle/index.md", pageContent(1))
			sb.WithContent("content/blog/mybundle/mydata.json", `{ "hello": "world" }`)

			return sb
		},
			func(s *sitesBuilder) {
				s.AssertFileContent("public/blog/mybundle/index.html", "Resources: application/json: /blog/mybundle/mydata.json")
				s.CheckExists("public/blog/mybundle/mydata.json")

			},
		},
		{"Tags and categories", func(b testing.TB) *sitesBuilder {
			sb := newTestSitesBuilder(b).WithConfigFile("toml", `
title = "Tags and Cats"
baseURL = "https://example.com"

`)

			const pageTemplate = `
---
title: "Some tags and cats"
categories: ["caGR", "cbGR"]
tags: ["taGR", "tbGR"]
---

Some content.
			
`
			for i := 1; i <= 100; i++ {
				content := strings.Replace(pageTemplate, "GR", strconv.Itoa(i/3), -1)
				sb.WithContent(fmt.Sprintf("content/page%d.md", i), content)
			}

			return sb
		},
			func(s *sitesBuilder) {
				s.AssertFileContent("public/page3/index.html", "/page3/|Permalink: https://example.com/page3/")
				s.AssertFileContent("public/tags/ta3/index.html", "|ta3|")
			},
		},
		{"Canonify URLs", func(b testing.TB) *sitesBuilder {
			sb := newTestSitesBuilder(b).WithConfigFile("toml", `
title = "Canon"
baseURL = "https://example.com"
canonifyURLs = true

`)
			for i := 1; i <= 100; i++ {
				sb.WithContent(fmt.Sprintf("content/page%d.md", i), pageContent(i))
			}

			return sb
		},
			func(s *sitesBuilder) {
				s.AssertFileContent("public/page8/index.html", "https://example.com/about/")
			},
		},

		{"Deep content tree", func(b testing.TB) *sitesBuilder {

			sb := newTestSitesBuilder(b).WithConfigFile("toml", `
baseURL = "https://example.com"

[languages]
[languages.en]
weight=1
contentDir="content/en"
[languages.fr]
weight=2
contentDir="content/fr"
[languages.no]
weight=3
contentDir="content/no"
[languages.sv]
weight=4
contentDir="content/sv"
			
`)

			createContent := func(dir, name string) {
				sb.WithContent(filepath.Join("content", dir, name), pageContent(1))
			}

			createBundledFiles := func(dir string) {
				sb.WithContent(filepath.Join("content", dir, "data.json"), `{ "hello": "world" }`)
				for i := 1; i <= 3; i++ {
					sb.WithContent(filepath.Join("content", dir, fmt.Sprintf("page%d.md", i)), pageContent(1))
				}
			}

			for _, lang := range []string{"en", "fr", "no", "sv"} {
				for level := 1; level <= 5; level++ {
					sectionDir := path.Join(lang, strings.Repeat("section/", level))
					createContent(sectionDir, "_index.md")
					createBundledFiles(sectionDir)
					for i := 1; i <= 3; i++ {
						leafBundleDir := path.Join(sectionDir, fmt.Sprintf("bundle%d", i))
						createContent(leafBundleDir, "index.md")
						createBundledFiles(path.Join(leafBundleDir, "assets1"))
						createBundledFiles(path.Join(leafBundleDir, "assets1", "assets2"))
					}
				}
			}

			return sb
		},
			func(s *sitesBuilder) {
				s.CheckExists("public/blog/mybundle/index.html")
				s.Assert(len(s.H.Sites), qt.Equals, 4)
				s.Assert(len(s.H.Sites[0].RegularPages()), qt.Equals, len(s.H.Sites[1].RegularPages()))
				s.Assert(len(s.H.Sites[0].RegularPages()), qt.Equals, 30)

			},
		},
		{"Many HTML templates", func(b testing.TB) *sitesBuilder {

			pageTemplateTemplate := `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{ if not .IsPage }}{{ .Title }}{{ else }}{{ printf "Site: %s" site.Title }}{{ end }}</title>
    <style>
     body {
       margin: 3rem;
     }
    </style>
  </head>
  <body>
    <div class="page">{{ .Content }}</div>
    <ul>
    {{ with .Pages }}
    {{ range . }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }} {{ if not .IsNode }} (Page){{ end }}</a></li>
    {{ end }}
    {{ end }}
    </ul>
  </body>
</html>
`

			sb := newTestSitesBuilder(b).WithConfigFile("toml", `
baseURL = "https://example.com"

[languages]
[languages.en]
weight=1
contentDir="content/en"
[languages.fr]
weight=2
contentDir="content/fr"
[languages.no]
weight=3
contentDir="content/no"
[languages.sv]
weight=4
contentDir="content/sv"
			
`)

			createContent := func(dir, name string) {
				sb.WithContent(filepath.Join("content", dir, name), pageContent(1))
			}

			for _, lang := range []string{"en", "fr", "no", "sv"} {
				sb.WithTemplatesAdded(fmt.Sprintf("_default/single.%s.html", lang), pageTemplateTemplate)
				sb.WithTemplatesAdded(fmt.Sprintf("_default/list.%s.html", lang), pageTemplateTemplate)

				for level := 1; level <= 5; level++ {
					sectionDir := path.Join(lang, strings.Repeat("section/", level))
					createContent(sectionDir, "_index.md")
					for i := 1; i <= 3; i++ {
						leafBundleDir := path.Join(sectionDir, fmt.Sprintf("bundle%d", i))
						createContent(leafBundleDir, "index.md")
					}
				}
			}

			return sb
		},
			func(s *sitesBuilder) {
				s.CheckExists("public/blog/mybundle/index.html")
				s.Assert(len(s.H.Sites), qt.Equals, 4)
				s.Assert(len(s.H.Sites[0].RegularPages()), qt.Equals, len(s.H.Sites[1].RegularPages()))
				s.Assert(len(s.H.Sites[0].RegularPages()), qt.Equals, 15)

			},
		},
		{"Page collections", func(b testing.TB) *sitesBuilder {

			pageTemplateTemplate := `
{{ if .IsNode }}
{{ len .Paginator.Pages }}
{{ end }}
{{ len .Sections }}
{{ len .Pages }}
{{ len .RegularPages }}
{{ len site.RegularPages }}
{{ len site.Pages }}

`

			sb := newTestSitesBuilder(b).WithConfigFile("toml", `
baseURL = "https://example.com"

[languages]
[languages.en]
weight=1
contentDir="content/en"
[languages.fr]
weight=2
contentDir="content/fr"
[languages.no]
weight=3
contentDir="content/no"
[languages.sv]
weight=4
contentDir="content/sv"
			
`)

			sb.WithTemplatesAdded("index.html", pageTemplateTemplate)
			sb.WithTemplatesAdded("_default/single.html", pageTemplateTemplate)
			sb.WithTemplatesAdded("_default/list.html", pageTemplateTemplate)

			createContent := func(dir, name string) {
				sb.WithContent(filepath.Join("content", dir, name), pageContent(1))
			}

			createBundledFiles := func(dir string) {
				sb.WithContent(filepath.Join("content", dir, "data.json"), `{ "hello": "world" }`)
				for i := 1; i <= 3; i++ {
					sb.WithContent(filepath.Join("content", dir, fmt.Sprintf("page%d.md", i)), pageContent(1))
				}
			}

			r := rand.New(rand.NewSource(99))

			for _, lang := range []string{"en", "fr", "no", "sv"} {
				for level := 1; level <= r.Intn(5)+1; level++ {
					sectionDir := path.Join(lang, strings.Repeat("section/", level))
					createContent(sectionDir, "_index.md")
					createBundledFiles(sectionDir)
					for i := 1; i <= r.Intn(20)+1; i++ {
						leafBundleDir := path.Join(sectionDir, fmt.Sprintf("bundle%d", i))
						createContent(leafBundleDir, "index.md")
						createBundledFiles(path.Join(leafBundleDir, "assets1"))
						createBundledFiles(path.Join(leafBundleDir, "assets1", "assets2"))
					}
				}
			}

			return sb
		},
			func(s *sitesBuilder) {
				s.CheckExists("public/blog/mybundle/index.html")
				s.Assert(len(s.H.Sites), qt.Equals, 4)
				s.Assert(len(s.H.Sites[0].RegularPages()), qt.Equals, 26)

			},
		},
	}

	return benchmarks

}

// Run the benchmarks below as tests. Mostly useful when adding new benchmark
// variants.
func TestBenchmarkSiteNew(b *testing.T) {
	benchmarks := getBenchmarkSiteNewTestCases()
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.T) {
			s := bm.create(b)

			err := s.BuildE(BuildCfg{})
			if err != nil {
				b.Fatal(err)
			}
			bm.check(s)

		})
	}
}

// TODO(bep) eventually remove the old (too complicated setup).
func BenchmarkSiteNew(b *testing.B) {
	benchmarks := getBenchmarkSiteNewTestCases()

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			sites := make([]*sitesBuilder, b.N)
			for i := 0; i < b.N; i++ {
				sites[i] = bm.create(b)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := sites[i]
				err := s.BuildE(BuildCfg{})
				if err != nil {
					b.Fatal(err)
				}
				bm.check(s)
			}
		})
	}
}
