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

package hugolib

import "testing"

func TestFrontMatterParamsInItsOwnSection(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
-- content/_index.md --
+++
title = "Home"
[[cascade]]
background = 'yosemite.jpg'
[cascade.params]
a = "home-a"
b = "home-b"
[cascade._target]
kind = 'page'
+++
-- content/p1.md --
---
title: "P1"
summary: "frontmatter.summary"
params:
   a: "p1-a"
   summary: "params.summary"
---	
-- layouts/_default/single.html --
Params: {{ range $k, $v := .Params }}{{ $k }}: {{ $v }}|{{ end }}$
Summary: {{ .Summary }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"Params: a: p1-a|b: home-b|background: yosemite.jpg|draft: false|iscjklanguage: false|summary: params.summary|title: P1|$",
		"Summary: frontmatter.summary|",
	)
}
