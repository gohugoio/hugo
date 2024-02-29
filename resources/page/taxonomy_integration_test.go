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
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestTaxonomiesGetAndCount(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap']
[taxonomies]
author = 'authors'
-- layouts/_default/home.html --
John Smith count: {{ site.Taxonomies.authors.Count "John Smith" }}
Robert Jones count: {{ (site.Taxonomies.authors.Get "Robert Jones").Pages.Len }}
-- layouts/_default/single.html --
{{ .Title }}|
-- layouts/_default/list.html --
{{ .Title }}|
-- content/p1.md --
---
title: p1
authors: [John Smith,Robert Jones]
---
-- content/p2.md --
---
title: p2
authors: [John Smith]
---
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"John Smith count: 2",
		"Robert Jones count: 1",
	)
}
