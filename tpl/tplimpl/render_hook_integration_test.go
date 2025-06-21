// Copyright 2025 The Hugo Authors. All rights reserved.
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

// empty
[]()
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

		`<a href=""></a>`,
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
[markup.goldmark.extensions.typographer]
disable = true
[markup.goldmark.parser]
wrapStandAloneImageWithinParagraph = false
[markup.goldmark.parser.attribute]
block = false
[markup.goldmark.renderHooks.image]
enableDefault = true
-- content/p1/index.md --
![]()

![alt1](./pixel.png)

![alt2-&<>'](pixel.png "&<>'")

![alt3](pixel.png?a=b&c=d#fragment)
{.foo #bar}

![alt4](pixel.png)
{id="\"><script>alert()</script>"}
-- content/p1/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html",
		`<img src="" alt="">`,
		`<img src="/dir/p1/pixel.png" alt="alt1">`,
		`<img src="/dir/p1/pixel.png" alt="alt2-&amp;&lt;&gt;&#39;" title="&amp;&lt;&gt;&#39;">`,
		`<img src="/dir/p1/pixel.png?a=b&amp;c=d#fragment" alt="alt3">`,
		`<img src="/dir/p1/pixel.png" alt="alt4">`,
	)

	files = strings.Replace(files, "block = false", "block = true", -1)

	b = hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html",
		`<img src="" alt="">`,
		`<img src="/dir/p1/pixel.png" alt="alt1">`,
		`<img src="/dir/p1/pixel.png" alt="alt2-&amp;&lt;&gt;&#39;" title="&amp;&lt;&gt;&#39;">`,
		`<img src="/dir/p1/pixel.png?a=b&amp;c=d#fragment" alt="alt3" class="foo" id="bar">`,
		`<img src="/dir/p1/pixel.png" alt="alt4" id="&#34;&gt;&lt;script&gt;alert()&lt;/script&gt;">`,
	)
}

// Issue 13535
func TestEmbeddedLinkAndImageRenderHookConfig(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']

[markup.goldmark]
duplicateResourceFiles = false

[markup.goldmark.renderHooks.image]
#KEY_VALUE

[markup.goldmark.renderHooks.link]
#KEY_VALUE

#LANGUAGES
-- content/s1/p1/index.md --
---
title: p1
---
[p2](p2)

[a](a.txt)

![b](b.jpg)
-- content/s1/p1/a.txt --
-- content/s1/p1/b.jpg --
-- content/s1/p2.md --
---
title: p2
---
-- layouts/all.html --
{{ .Content }}
`

	const customHooks string = `
-- layouts/s1/_markup/render-link.html --
custom link render hook: {{ .Text }}|{{ .Destination }}
-- layouts/s1/_markup/render-image.html --
custom image render hook: {{ .Text }}|{{ .Destination }}
`

	const languages string = `
[languages.en]
[languages.fr]
`

	const (
		fileToCheck         = "public/s1/p1/index.html"
		wantCustom   string = "<p>custom link render hook: p2|p2</p>\n<p>custom link render hook: a|a.txt</p>\n<p>custom image render hook: b|b.jpg</p>"
		wantEmbedded string = "<p><a href=\"/s1/p2/\">p2</a></p>\n<p><a href=\"/s1/p1/a.txt\">a</a></p>\n<p><img src=\"/s1/p1/b.jpg\" alt=\"b\"></p>"
		wantGoldmark string = "<p><a href=\"p2\">p2</a></p>\n<p><a href=\"a.txt\">a</a></p>\n<p><img src=\"b.jpg\" alt=\"b\"></p>"
	)

	tests := []struct {
		id             string // the test id
		isMultilingual bool   // whether the site is multilingual single-host
		hasCustomHooks bool   // whether the site has custom link and image render hooks
		keyValuePair   string // the enableDefault (deprecated in v0.148.0) or useEmbedded key-value pair
		want           string // the expected content of public/s1/p1/index.html
	}{
		{"01", false, false, "", wantGoldmark},                         // monolingual
		{"02", false, false, "enableDefault = false", wantGoldmark},    // monolingual, enableDefault = false
		{"03", false, false, "enableDefault = true", wantEmbedded},     // monolingual, enableDefault = true
		{"04", false, false, "useEmbedded = 'always'", wantEmbedded},   // monolingual, useEmbedded = 'always'
		{"05", false, false, "useEmbedded = 'auto'", wantGoldmark},     // monolingual, useEmbedded = 'auto'
		{"06", false, false, "useEmbedded = 'fallback'", wantEmbedded}, // monolingual, useEmbedded = 'fallback'
		{"07", false, false, "useEmbedded = 'never'", wantGoldmark},    // monolingual, useEmbedded = 'never'
		{"08", false, true, "", wantCustom},                            // monolingual, with custom hooks
		{"09", false, true, "enableDefault = false", wantCustom},       // monolingual, with custom hooks, enableDefault = false
		{"10", false, true, "enableDefault = true", wantCustom},        // monolingual, with custom hooks, enableDefault = true
		{"11", false, true, "useEmbedded = 'always'", wantEmbedded},    // monolingual, with custom hooks, useEmbedded = 'always'
		{"12", false, true, "useEmbedded = 'auto'", wantCustom},        // monolingual, with custom hooks, useEmbedded = 'auto'
		{"13", false, true, "useEmbedded = 'fallback'", wantCustom},    // monolingual, with custom hooks, useEmbedded = 'fallback'
		{"14", false, true, "useEmbedded = 'never'", wantCustom},       // monolingual, with custom hooks, useEmbedded = 'never'
		{"15", true, false, "", wantEmbedded},                          // multilingual
		{"16", true, false, "enableDefault = false", wantGoldmark},     // multilingual, enableDefault = false
		{"17", true, false, "enableDefault = true", wantEmbedded},      // multilingual, enableDefault = true
		{"18", true, false, "useEmbedded = 'always'", wantEmbedded},    // multilingual, useEmbedded = 'always'
		{"19", true, false, "useEmbedded = 'auto'", wantEmbedded},      // multilingual, useEmbedded = 'auto'
		{"20", true, false, "useEmbedded = 'fallback'", wantEmbedded},  // multilingual, useEmbedded = 'fallback'
		{"21", true, false, "useEmbedded = 'never'", wantGoldmark},     // multilingual, useEmbedded = 'never'
		{"22", true, true, "", wantCustom},                             // multilingual, with custom hooks
		{"23", true, true, "enableDefault = false", wantCustom},        // multilingual, with custom hooks, enableDefault = false
		{"24", true, true, "enableDefault = true", wantCustom},         // multilingual, with custom hooks, enableDefault = true
		{"25", true, true, "useEmbedded = 'always'", wantEmbedded},     // multilingual, with custom hooks, useEmbedded = 'always'
		{"26", true, true, "useEmbedded = 'auto'", wantCustom},         // multilingual, with custom hooks, useEmbedded = 'auto'
		{"27", true, true, "useEmbedded = 'fallback'", wantCustom},     // multilingual, with custom hooks, useEmbedded = 'fallback'
		{"28", true, true, "useEmbedded = 'never'", wantCustom},        // multilingual, with custom hooks, useEmbedded = 'never'
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			f := files
			if tt.isMultilingual {
				f = strings.ReplaceAll(f, "#LANGUAGES", languages)
			}
			if tt.hasCustomHooks {
				f = f + customHooks
			}
			f = strings.ReplaceAll(f, "#KEY_VALUE", tt.keyValuePair)

			b := hugolib.Test(t, f)
			b.AssertFileContent(fileToCheck, tt.want)
		})
	}
}
