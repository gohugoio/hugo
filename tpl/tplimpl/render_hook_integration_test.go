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

package tplimpl_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestEmbeddedLinkRenderHook(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term']
[markup.goldmark.renderHooks.link]
enableDefault = true
-- layouts/_default/list.html --
{{ .Content }}
-- layouts/_default/single.html --
{{ .Content }}
-- assets/a.txt --
irrelevant
-- content/_index.md --
---
title: home
---
-- content/s1/_index.md --
---
title: s1
---
-- content/s1/p1.md --
---
title: s1/p1
---
-- content/s1/p2/index.md --
---
title: s1/p2
---
[500](a.txt) // global resource
[510](b.txt) // page resource
[520](./b.txt) // page resource
-- content/s1/p2/b.txt --
irrelevant
-- content/s1/p3.md --
---
title: s1/p3
---
// Remote
[10](https://a.org)

// fragment
[100](/#foo)
[110](#foo)
[120](p1#foo)
[130](p1/#foo)

// section page
[200](s1)
[210](/s1)
[220](../s1)
[230](s1/)
[240](/s1/)
[250](../s1/)

// regular page
[300](p1)
[310](/s1/p1)
[320](../s1/p1)
[330](p1/)
[340](/s1/p1/)
[350](../s1/p1/)

// leaf bundle
[400](p2)
[410](/s1/p2)
[420](../s1/p2)
[430](p2/)
[440](/s1/p2/)
[450](../s1/p2/)
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/p3/index.html",
		`<a href="https://a.org">10</a>`,
		`<a href="/#foo">100</a>`,
		`<a href="/s1/p3/#foo">110</a>`,
		`<a href="/s1/p1/#foo">120</a>`,
		`<a href="/s1/p1/#foo">130</a>`,

		`<a href="/s1/">200</a>`,
		`<a href="/s1/">210</a>`,
		`<a href="/s1/">220</a>`,
		`<a href="/s1/">230</a>`,
		`<a href="/s1/">240</a>`,
		`<a href="/s1/">250</a>`,

		`<a href="/s1/p1/">300</a>`,
		`<a href="/s1/p1/">310</a>`,
		`<a href="/s1/p1/">320</a>`,
		`<a href="/s1/p1/">330</a>`,
		`<a href="/s1/p1/">340</a>`,
		`<a href="/s1/p1/">350</a>`,

		`<a href="/s1/p2/">400</a>`,
		`<a href="/s1/p2/">410</a>`,
		`<a href="/s1/p2/">420</a>`,
		`<a href="/s1/p2/">430</a>`,
		`<a href="/s1/p2/">440</a>`,
		`<a href="/s1/p2/">450</a>`,
	)

	b.AssertFileContent("public/s1/p2/index.html",
		`<a href="/a.txt">500</a>`,
		`<a href="/s1/p2/b.txt">510</a>`,
		`<a href="/s1/p2/b.txt">520</a>`,
	)
}

// Issue 12203
// Issue 12468
// Issue 12514
func TestEmbeddedImageRenderHook(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'https://example.org/dir/'
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
[markup.goldmark.parser]
wrapStandAloneImageWithinParagraph = false
[markup.goldmark.parser.attribute]
block = false
[markup.goldmark.renderHooks.image]
enableDefault = true
-- content/p1/index.md --
![alt1](./pixel.png)

![alt2](pixel.png?a=b&c=d#fragment)
{.foo #bar}
-- content/p1/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html",
		`<img alt="alt1" src="/dir/p1/pixel.png">`,
		`<img alt="alt2" src="/dir/p1/pixel.png?a=b&c=d#fragment">`,
	)

	files = strings.Replace(files, "block = false", "block = true", -1)

	b = hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html",
		`<img alt="alt1" src="/dir/p1/pixel.png">`,
		`<img alt="alt2" class="foo" id="bar" src="/dir/p1/pixel.png?a=b&c=d#fragment">`,
	)
}
