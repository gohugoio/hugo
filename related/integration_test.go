// Copyright 2023 The Hugo Authors. All rights reserved.
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

package related_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func BenchmarkRelatedSite(b *testing.B) {
	files := `
-- config.toml --
baseURL = "http://example.com/"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT"]
[related]
  includeNewer = false
  threshold = 80
  toLower = false
[[related.indices]]
  name = 'keywords'
  weight = 70	
-- layouts/_default/single.html --
Len related: {{ site.RegularPages.Related . | len }}
`

	createContent := func(n int) string {
		base := `---
title: "Page %d"
keywords: ['k%d']
---
`

		for i := 0; i < 32; i++ {
			base += fmt.Sprintf("\n## Title %d", rand.Intn(100))
		}

		return fmt.Sprintf(base, n, rand.Intn(32))

	}

	for i := 1; i < 100; i++ {
		files += fmt.Sprintf("\n-- content/posts/p%d.md --\n"+createContent(i+1), i+1)
	}

	cfg := hugolib.IntegrationTestConfig{
		T:           b,
		TxtarString: files,
	}
	builders := make([]*hugolib.IntegrationTestBuilder, b.N)

	for i := range builders {
		builders[i] = hugolib.NewIntegrationTestBuilder(cfg)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builders[i].Build()
	}
}
