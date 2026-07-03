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

package rst_test

import (
	"runtime"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/rst"
)

// --syntax-highlight, default = long.
func TestRSTSyntaxHighlight(t *testing.T) {
	if !rst.Supports() {
		t.Skip("rst not installed")
	}
	if !rst.SupportsPygments() {
		t.Skip("Pygments not installed")
	}

	execAllow := `'^rst2html(\.py)?$'`
	if runtime.GOOS == "windows" {
		execAllow = `'^python(\.exe)?$'`
	}

	filesTemplate := strings.ReplaceAll(`
-- hugo.toml --
baseURL = "https://example.org"
disableKinds = ["home", "section", "taxonomy", "term", "rss", "sitemap"]

markup.rst.syntaxHighlight = 'SYNTAX_HIGHLIGHT'

[security.exec]
allow = [RST_EXEC_ALLOW]

-- layouts/page.html --
{{ .Content }}
-- content/p.rst --
---
title: p
---

.. code:: go

   if true {}
`, "RST_EXEC_ALLOW", execAllow)

	files := strings.ReplaceAll(filesTemplate, "SYNTAX_HIGHLIGHT", "none")
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p/index.html", `
! Pygments			
`)

	b.AssertFileContent("public/p/index.html", `
if true {}
! if</span>
`)

	files = strings.ReplaceAll(filesTemplate, "SYNTAX_HIGHLIGHT", "short")
	b = hugolib.Test(t, files)

	b.AssertFileContent("public/p/index.html", `
<span class="k">if</span>
! <span class="keyword">if</span>
`)

	files = strings.ReplaceAll(filesTemplate, "SYNTAX_HIGHLIGHT", "long")
	b = hugolib.Test(t, files)

	b.AssertFileContent("public/p/index.html", `
! <span class="k">if</span>
<span class="keyword">if</span>
`)

	// Default long.
	files = strings.ReplaceAll(filesTemplate, "markup.rst.syntaxHighlight = 'SYNTAX_HIGHLIGHT'", "")
	b = hugolib.Test(t, files)

	b.AssertFileContent("public/p/index.html", `
! <span class="k">if</span>
<span class="keyword">if</span>
`)

	// Invalid.
	files = strings.ReplaceAll(filesTemplate, "SYNTAX_HIGHLIGHT", "foo")
	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "rst: invalid value for syntaxHighlight: \"foo\"")
}
