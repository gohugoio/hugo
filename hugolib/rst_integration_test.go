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

package hugolib

import (
	"runtime"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/markup/rst"
)

func TestRSTHighlightClassNamingConfigIssue5349(t *testing.T) {
	if !rst.Supports() {
		t.Skip("rst not installed")
	}

	execAllow := `'^rst2html(\.py)?$'`
	if runtime.GOOS == "windows" {
		execAllow = `'^python(\.exe)?$'`
	}

	files := strings.ReplaceAll(`
-- hugo.toml --
baseURL = "https://example.org"
disableKinds = ["home", "section", "taxonomy", "term", "rss", "sitemap"]
[security.exec]
allow = [RST_EXEC_ALLOW]
[markup.rst.highlight]
classNaming = 'none'
-- layouts/page.html --
{{ .Content }}
-- content/p.rst --
---
title: p
---

.. code:: go

   if true {}
`, "RST_EXEC_ALLOW", execAllow)

	b := Test(t, files)
	content := b.FileContent("public/p/index.html")
	b.Assert(content, qt.Contains, `if true {}`)
	b.Assert(content, qt.Not(qt.Contains), `Pygments package not found`)
	b.Assert(content, qt.Not(qt.Contains), `<span class="keyword">if</span>`)
	b.Assert(content, qt.Not(qt.Contains), `<span class="k">if</span>`)
}
