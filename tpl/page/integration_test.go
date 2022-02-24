// Copyright 2022 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestPageGlobal(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- content/_index.md --
---
title: "Home"
---
{{< shortcode >}}

## Heading

[I'm an inline-style link](https://www.google.com)

![alt text](https://github.com/adam-p/markdown-here/raw/master/src/common/images/icon48.png "Logo Title Text 1")

$$$bash
echo "hello";
$$$

-- layouts/_default/_markup/render-heading.html --
{{ if page.IsHome }}
Heading OK.
{{ end }}
-- layouts/_default/_markup/render-image.html --
{{ if page.IsHome }}
Image OK.
{{ end }}
-- layouts/_default/_markup/render-link.html --
{{ if page.IsHome }}
Link OK.
{{ end }}
-- layouts/_default/_markup/render-codeblock.html --
{{ if page.IsHome }}
Codeblock OK.
{{ end }}
-- layouts/index.html --
{{ if eq page . }}
Page OK.
{{ end }}
{{ .Content }}
partial: {{ partials.Include "foo.html" . }}
-- layouts/partials/foo.html --
{{ if page.IsHome }}
Partial OK.
{{ end }}
-- layouts/shortcodes/shortcode.html --
{{ if page.IsHome }}
Shortcode OK.
{{ end }}
  `

	// Fenced code blocks.
	files = strings.ReplaceAll(files, "$$$", "```")

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
Heading OK.
Image OK.
Link OK.
Codeblock OK.
Page OK.
Partial OK.
Shortcode OK.
`)
}
