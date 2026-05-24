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
	"os/exec"
	"runtime"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
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

	filesTemplate := strings.ReplaceAll(`
-- hugo.toml --
baseURL = "https://example.org"
disableKinds = ["home", "section", "taxonomy", "term", "rss", "sitemap"]
[security.exec]
allow = [RST_EXEC_ALLOW]
[markup.rst.highlight]
classNaming = 'CLASS_NAMING'
-- layouts/page.html --
{{ .Content }}
-- content/p.rst --
---
title: p
---

.. code:: go

   if true {}
`, "RST_EXEC_ALLOW", execAllow)

	for _, test := range []struct {
		name             string
		classNaming      string
		requiresPygments bool
	}{
		{"None", "none", false},
		{"Short", "short", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.requiresPygments && !supportsPygments() {
				qt.Assert(t, htesting.IsGitHubAction(), qt.Equals, false, qt.Commentf("Pygments not installed"))
				t.Skip("Pygments not installed")
			}

			files := strings.ReplaceAll(filesTemplate, "CLASS_NAMING", test.classNaming)
			b := Test(t, files)
			content := b.FileContent("public/p/index.html")
			b.Assert(content, qt.Not(qt.Contains), `Pygments package not found`)

			if test.classNaming == "none" {
				b.Assert(content, qt.Contains, `if true {}`)
				b.Assert(content, qt.Not(qt.Contains), `<span class="keyword">if</span>`)
				b.Assert(content, qt.Not(qt.Contains), `<span class="k">if</span>`)
			} else {
				b.Assert(content, qt.Contains, `<span class="k">if</span>`)
				b.Assert(content, qt.Not(qt.Contains), `<span class="keyword">if</span>`)
			}
		})
	}
}

func supportsPygments() bool {
	for _, python := range []string{"python3", "python"} {
		if err := exec.Command(python, "-c", "import pygments").Run(); err == nil {
			return true
		}
	}
	return false
}
