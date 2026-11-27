// Copyright 2026 The Hugo Authors. All rights reserved.
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

package strings_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// See issue 15668.
func TestTruncatePlainTextOutputFormat(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
[outputFormats.llms]
mediaType = "text/plain"
baseName = "llms"
isPlainText = true
[outputs]
home = ["html", "llms"]
-- content/p.md --
---
title: P
description: 'Voltage < 10 kV & "rated" > 5 A'
---
-- layouts/home.llms.txt --
{{ range site.RegularPages }}
description:        {{ .Description }}
plainify:           {{ .Description | plainify }}
plainify|truncate:  {{ .Description | plainify | truncate 200 }}
truncate:           {{ .Description | truncate 200 }}
truncate short:     {{ .Description | truncate 13 }}
{{ end }}
-- layouts/home.html --
{{ range site.RegularPages }}
truncate:           {{ .Description | truncate 200 }}
truncate short:     {{ .Description | truncate 13 }}
truncate html:      {{ "<b>a</b> b c" | safeHTML | truncate 3 }}
truncate html ellipsis: {{ truncate 3 ("<i>…</i>" | safeHTML) "a < b c" }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/llms.txt",
		`description:        Voltage < 10 kV & "rated" > 5 A`,
		`plainify:           Voltage < 10 kV & "rated" > 5 A`,
		`plainify|truncate:  Voltage < 10 kV & "rated" > 5 A`,
		`truncate:           Voltage < 10 kV & "rated" > 5 A`,
		`truncate short:     Voltage < 10 …`,
	)
	b.AssertFileContent("public/index.html",
		`truncate:           Voltage &lt; 10 kV &amp; &#34;rated&#34; &gt; 5 A`,
		`truncate short:     Voltage &lt; 10 …`,
		`truncate html:      <b>a</b> b …`,
		`truncate html ellipsis: a &lt;<i>…</i>`,
	)
}
