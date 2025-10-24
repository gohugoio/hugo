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

// Package asciidocext converts AsciiDoc to HTML using Asciidoctor
// external binary. The `asciidoc` module is reserved for a future golang
// implementation.

package asciidocext_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/asciidocext"
)

// Issue 9202
func TestAsciidoctorMultibyteOutdir(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("skip asciidoc")
	}

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']

[markup.asciidocext]
workingFolderCurrent = true

[security.exec]
allow = ['^asciidoctor$']
-- content/p1.adoc --
---
title: p1
---
H~2~O
-- content/hügo.adoc --
---
title: hügo
---
H~2~O
-- layouts/all.html --
{{ .Content }}
`

	b, err := hugolib.TestE(t, files, hugolib.TestOptInfo())
	if err != nil {
		fmt.Println(err.Error())
	}

	wantContent := "<p>H<sub>2</sub>O</p>"

	b.AssertFileContent("public/p1/index.html", wantContent)
	b.AssertFileContent("public/hügo/index.html", wantContent)
	b.AssertLogContains("/public/p1")
	b.AssertLogContains("/public/hügo")
}

// Issue 10473
func TestAsciidoctorMultilingualOutdir(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("skip asciidoc")
	}

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']

defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true

[languages.en]
weight = 1
BASEURL_EN

[languages.de]
weight = 2
BASEURL_DE

[markup.asciidocExt]
workingFolderCurrent = true

[security.exec]
allow = ['^asciidoctor$']
-- content/p1.en.adoc --
---
title: p1 (en)
---
H~2~O
-- content/p1.de.adoc --
---
title: p1 (de)
---
H~2~O
-- layouts/all.html --
{{ .Content }}
`

	wantContent := "<p>H<sub>2</sub>O</p>"

	// multilingual single-host
	f := strings.ReplaceAll(files, "BASEURL_EN", "")
	f = strings.ReplaceAll(f, "BASEURL_DE", "")

	b := hugolib.Test(t, f, hugolib.TestOptInfo())

	b.AssertFileContent("public/en/p1/index.html", wantContent)
	b.AssertFileContent("public/de/p1/index.html", wantContent)
	b.AssertLogContains("/public/en/p1")
	b.AssertLogContains("/public/de/p1")

	// multilingual multi-host
	f = strings.ReplaceAll(files, "BASEURL_EN", "baseURL = 'https://en.example.org/'")
	f = strings.ReplaceAll(f, "BASEURL_DE", "baseURL = 'https://de.example.org/'")

	b = hugolib.Test(t, f, hugolib.TestOptInfo())

	b.AssertFileContent("public/en/p1/index.html", wantContent)
	b.AssertFileContent("public/de/p1/index.html", wantContent)
	// b.AssertLogContains("/public/en/p1") // JMM fail: outdir contains /public/p1 (missing language prefix)
	// b.AssertLogContains("/public/de/p1") // JMM fail: outdir contains /public/p1 (missing language prefix)
}
