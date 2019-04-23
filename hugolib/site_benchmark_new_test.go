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
	"testing"
)

// TODO(bep) eventually remove the old (too complicated setup).
func BenchmarkSiteNew(b *testing.B) {
	// TODO(bep) create some common and stable data set

	const pageContent = `---
title: "My Page"
---

My page content.

`

	config := `
baseURL = "https://example.com"

`

	benchmarks := []struct {
		name   string
		create func(i int) *sitesBuilder
		check  func(s *sitesBuilder)
	}{
		{"Bundle with image", func(i int) *sitesBuilder {
			sb := newTestSitesBuilder(b).WithConfigFile("toml", config)
			sb.WithContent("content/blog/mybundle/index.md", pageContent)
			sb.WithSunset("content/blog/mybundle/sunset1.jpg")

			return sb
		},
			func(s *sitesBuilder) {
				s.AssertFileContent("public/blog/mybundle/index.html", "/blog/mybundle/sunset1.jpg")
				s.CheckExists("public/blog/mybundle/sunset1.jpg")

			},
		},
		{"Bundle with JSON file", func(i int) *sitesBuilder {
			sb := newTestSitesBuilder(b).WithConfigFile("toml", config)
			sb.WithContent("content/blog/mybundle/index.md", pageContent)
			sb.WithContent("content/blog/mybundle/mydata.json", `{ "hello": "world" }`)

			return sb
		},
			func(s *sitesBuilder) {
				s.AssertFileContent("public/blog/mybundle/index.html", "Resources: application/json: /blog/mybundle/mydata.json")
				s.CheckExists("public/blog/mybundle/mydata.json")

			},
		},
		{"Multiple languages", func(i int) *sitesBuilder {
			sb := newTestSitesBuilder(b).WithConfigFile("toml", `
baseURL = "https://example.com"

[languages]
[languages.en]
weight=1
[languages.fr]
weight=2
			
`)

			return sb
		},
			func(s *sitesBuilder) {

			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			sites := make([]*sitesBuilder, b.N)
			for i := 0; i < b.N; i++ {
				sites[i] = bm.create(i)
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
