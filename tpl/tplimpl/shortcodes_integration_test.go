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

	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/gohugoio/hugo/hugolib"
)

func TestCommentShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ .Content }}
-- content/_index.md --
---
title: home
---
a{{< comment >}}b{{< /comment >}}c
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())
	b.AssertFileContent("public/index.html", "<p>ac</p>")
	b.AssertLogContains(`WARN  The "comment" shortcode was deprecated in v0.143.0 and will be removed in a future release. Please use HTML comments instead.`)
}

func TestDetailsShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ .Content }}
-- content/_index.md --
---
title: home
---
{{< details >}}
A: An _emphasized_ word.
{{< /details >}}

{{< details
  class="my-class"
  name="my-name"
  open=true
  summary="A **bold** word"
  title="my-title"
>}}
B: An _emphasized_ word.
{{< /details >}}

{{< details open=false >}}
C: An _emphasized_ word.
{{< /details >}}

{{< details open="false" >}}
D: An _emphasized_ word.
{{< /details >}}

{{< details open=0 >}}
E: An _emphasized_ word.
{{< /details >}}
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"<details>\n  <summary>Details</summary>\n  <p>A: An <em>emphasized</em> word.</p>\n</details>",
		"<details class=\"my-class\" name=\"my-name\" open title=\"my-title\">\n  <summary>A <strong>bold</strong> word</summary>\n  <p>B: An <em>emphasized</em> word.</p>\n</details>",
		"<details>\n  <summary>Details</summary>\n  <p>C: An <em>emphasized</em> word.</p>\n</details>",
		"<details>\n  <summary>Details</summary>\n  <p>D: An <em>emphasized</em> word.</p>\n</details>",
		"<details>\n  <summary>Details</summary>\n  <p>D: An <em>emphasized</em> word.</p>\n</details>",
	)
}

func TestGistShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ .Content }}
-- content/_index.md --
---
title: home
---
{{< gist jmooring 23932424365401ffa5e9d9810102a477 >}}
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())
	b.AssertFileContent("public/index.html", `<script src="https://gist.github.com/jmooring/23932424365401ffa5e9d9810102a477.js"></script>`)
	b.AssertLogContains(`WARN  The "gist" shortcode was deprecated in v0.143.0 and will be removed in a future release. See https://gohugo.io/shortcodes/gist for instructions to create a replacement.`)
}

func TestInstagramShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
privacy.instagram.simple = false
-- content/_index.md --
---
title: home
---
{{< instagram CxOWiQNP2MO >}}
-- layouts/index.html --
Hash: {{ .Content | hash.XxHash }}
Content: {{ .Content }}
`

	// Regular mode
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "a7937c49665872d3")

	// Simple mode
	files = strings.ReplaceAll(files, "privacy.instagram.simple = false", "privacy.instagram.simple = true")
	b = hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "2c1dce3881be0513")
}

func TestQRShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ .Content }}
-- content/_index.md --
---
title: home
---
{{< qr
	text="https://gohugo.io"
	level="high"
	scale=4
	targetDir="codes"
	alt="QR code linking to https://gohugo.io"
	class="my-class"
	id="my-id"
	title="My Title"
/>}}

{{< qr >}}
https://gohugo.io"
{{< /qr >}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		`<img src="/codes/qr_be5d263c2671bcbd.png" width="148" height="148" alt="QR code linking to https://gohugo.io" class="my-class" id="my-id" title="My Title">`,
		`<img src="/qr_472aab57ec7a6e3d.png" width="132" height="132">`,
	)
}

func TestVimeoShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
privacy.vimeo.simple = false
-- content/_index.md --
---
title: home
---
{{< vimeo 55073825 >}}
-- layouts/index.html --
Hash: {{ .Content | hash.XxHash }}
Content: {{ .Content }}
`

	// Regular mode
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "d5b2a079cc37d0ed")

	// Simple mode
	files = strings.ReplaceAll(files, "privacy.vimeo.simple = false", "privacy.vimeo.simple = true")
	b = hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "73b8767ce8bdf694")

	// Simple mode with non-existent id
	files = strings.ReplaceAll(files, "{{< vimeo 55073825 >}}", "{{< vimeo __id_does_not_exist__ >}}")
	b = hugolib.Test(t, files, hugolib.TestOptWarn())
	b.AssertLogContains(`WARN  The "vimeo" shortcode was unable to retrieve the remote data.`)
}

// Issue 13214
// We deprecated the twitter, tweet (alias of twitter), and twitter_simple
// shortcodes in v0.141.0, replacing them with x and x_simple.
func TestXShortcodes(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
#CONFIG
-- content/p1.md --
---
title: p1
---
{{< x user="SanDiegoZoo" id="1453110110599868418" >}}
-- content/p2.md --
---
title: p2
---
{{< twitter user="SanDiegoZoo" id="1453110110599868418" >}}
-- content/p3.md --
---
title: p3
---
{{< tweet user="SanDiegoZoo" id="1453110110599868418" >}}
-- content/p4.md --
---
title: p4
---
{{< x_simple user="SanDiegoZoo" id="1453110110599868418" >}}
-- content/p5.md --
---
title: p5
---
{{< twitter_simple user="SanDiegoZoo" id="1453110110599868418" >}}
-- layouts/_default/single.html --
{{ .Content | strings.TrimSpace | safeHTML }}
--
`

	b := hugolib.Test(t, files)

	// Test x, twitter, and tweet shortcodes
	want := `<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Owl bet you&#39;ll lose this staring contest 🦉 <a href="https://t.co/eJh4f2zncC">pic.twitter.com/eJh4f2zncC</a></p>&mdash; San Diego Zoo Wildlife Alliance (@sandiegozoo) <a href="https://twitter.com/sandiegozoo/status/1453110110599868418?ref_src=twsrc%5Etfw">October 26, 2021</a></blockquote>
	<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`
	b.AssertFileContent("public/p1/index.html", want)

	htmlFiles := []string{
		b.FileContent("public/p1/index.html"),
		b.FileContent("public/p2/index.html"),
		b.FileContent("public/p3/index.html"),
	}

	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	// Test x_simple and twitter_simple shortcodes
	wantSimple := "<style type=\"text/css\">\n      .twitter-tweet {\n        font:\n          14px/1.45 -apple-system,\n          BlinkMacSystemFont,\n          \"Segoe UI\",\n          Roboto,\n          Oxygen-Sans,\n          Ubuntu,\n          Cantarell,\n          \"Helvetica Neue\",\n          sans-serif;\n        border-left: 4px solid #2b7bb9;\n        padding-left: 1.5em;\n        color: #555;\n      }\n      .twitter-tweet a {\n        color: #2b7bb9;\n        text-decoration: none;\n      }\n      blockquote.twitter-tweet a:hover,\n      blockquote.twitter-tweet a:focus {\n        text-decoration: underline;\n      }\n    </style><blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Owl bet you&#39;ll lose this staring contest 🦉 <a href=\"https://t.co/eJh4f2zncC\">pic.twitter.com/eJh4f2zncC</a></p>&mdash; San Diego Zoo Wildlife Alliance (@sandiegozoo) <a href=\"https://twitter.com/sandiegozoo/status/1453110110599868418?ref_src=twsrc%5Etfw\">October 26, 2021</a></blockquote>\n--"
	b.AssertFileContent("public/p4/index.html", wantSimple)

	htmlFiles = []string{
		b.FileContent("public/p4/index.html"),
		b.FileContent("public/p5/index.html"),
	}
	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	filesOriginal := files

	// Test privacy.twitter.simple
	files = strings.ReplaceAll(filesOriginal, "#CONFIG", "privacy.twitter.simple=true")
	b = hugolib.Test(t, files)
	htmlFiles = []string{
		b.FileContent("public/p2/index.html"),
		b.FileContent("public/p3/index.html"),
		b.FileContent("public/p5/index.html"),
	}
	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	// Test privacy.x.simple
	files = strings.ReplaceAll(filesOriginal, "#CONFIG", "privacy.x.simple=true")
	b = hugolib.Test(t, files)
	htmlFiles = []string{
		b.FileContent("public/p1/index.html"),
		b.FileContent("public/p4/index.html"),
		b.FileContent("public/p4/index.html"),
	}
	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	htmlFiles = []string{
		b.FileContent("public/p2/index.html"),
		b.FileContent("public/p3/index.html"),
	}
	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	// Test privacy.twitter.disable
	files = strings.ReplaceAll(filesOriginal, "#CONFIG", "privacy.twitter.disable = true")
	b = hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html", "")
	htmlFiles = []string{
		b.FileContent("public/p1/index.html"),
		b.FileContent("public/p2/index.html"),
		b.FileContent("public/p3/index.html"),
		b.FileContent("public/p4/index.html"),
		b.FileContent("public/p4/index.html"),
	}
	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	// Test privacy.x.disable
	files = strings.ReplaceAll(filesOriginal, "#CONFIG", "privacy.x.disable = true")
	b = hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html", "")
	htmlFiles = []string{
		b.FileContent("public/p1/index.html"),
		b.FileContent("public/p4/index.html"),
	}
	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	htmlFiles = []string{
		b.FileContent("public/p2/index.html"),
		b.FileContent("public/p3/index.html"),
	}
	b.Assert(htmlFiles, hqt.IsAllElementsEqual)

	// Test warnings
	files = `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- content/_index.md --
---
title: home
---
{{< x user="__user_does_not_exist__" id="__id_does_not_exist__" >}}
{{< x_simple user="__user_does_not_exist__" id="__id_does_not_exist__" >}}
{{< twitter user="__user_does_not_exist__" id="__id_does_not_exist__" >}}
{{< twitter_simple user="__user_does_not_exist__" id="__id_does_not_exist__" >}}
-- layouts/index.html --
{{ .Content }}
`

	b = hugolib.Test(t, files, hugolib.TestOptWarn())
	b.AssertLogContains(
		`WARN  The "x" shortcode was unable to retrieve the remote data.`,
		`WARN  The "x_simple" shortcode was unable to retrieve the remote data.`,
		`WARN  The "twitter", "tweet", and "twitter_simple" shortcodes were deprecated in v0.142.0 and will be removed in a future release.`,
		`WARN  The "twitter" shortcode was unable to retrieve the remote data.`,
		`WARN  The "twitter_simple" shortcode was unable to retrieve the remote data.`,
	)
}
