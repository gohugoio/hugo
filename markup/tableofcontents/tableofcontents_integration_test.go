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

package tableofcontents_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue #10776
func TestHeadingsLevel(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ range .Fragments.HeadingsMap }}
	{{ printf "%s|%d|%s" .ID .Level .Title }}
{{ end }}
-- content/_index.md --
## Heading L2
### Heading L3
##### Heading L5
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html",
		"heading-l2|2|Heading L2",
		"heading-l3|3|Heading L3",
		"heading-l5|5|Heading L5",
	)
}

// Issue #13107
func TestToHTMLArgTypes(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','section','rss','sitemap','taxonomy','term']
-- layouts/_default/single.html --
{{ .Fragments.ToHTML .Params.toc.startLevel .Params.toc.endLevel false }}
-- content/json.md --
{
  "title": "json",
  "params": {
    "toc": {
      "startLevel": 2,
      "endLevel": 4
    }
  }
}
CONTENT
-- content/toml.md --
+++
title = 'toml'
[params.toc]
startLevel = 2
endLevel = 4
+++
CONTENT
-- content/yaml.md --
---
title: yaml
params:
  toc:
    startLevel: 2
    endLevel: 4
---
CONTENT
`

	content := `
# Level One
## Level Two
### Level Three
#### Level Four
##### Level Five
###### Level Six
	`

	want := `
<nav id="TableOfContents">
  <ul>
    <li><a href="#level-two">Level Two</a>
      <ul>
        <li><a href="#level-three">Level Three</a>
          <ul>
            <li><a href="#level-four">Level Four</a></li>
          </ul>
        </li>
      </ul>
    </li>
  </ul>
</nav>
`

	files = strings.ReplaceAll(files, "CONTENT", content)

	b := hugolib.Test(t, files)
	b.AssertFileContentEquals("public/json/index.html", strings.TrimSpace(want))
	b.AssertFileContentEquals("public/toml/index.html", strings.TrimSpace(want))
	b.AssertFileContentEquals("public/yaml/index.html", strings.TrimSpace(want))

	files = strings.ReplaceAll(files, `2`, `"x"`)

	b, _ = hugolib.TestE(t, files)
	b.AssertLogMatches(`error calling ToHTML: startLevel: unable to cast "x" of type string`)
}
