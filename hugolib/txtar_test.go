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
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

// See issue 10363.
func TestTxtarInProjectRoot(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.txtar --
== hugo.toml ==
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
theme = "mytheme"
== config/_default/params.toml ==
foo = "foo from txtar"
== content/_index.md ==
---
title: "Home"
---
== content/p1.md ==
---
title: "P1"
---
P1 content.
== content/p2.md ==
---
title: "P2 from txtar"
---
== layouts/home.html ==
Home|{{ .Title }}|{{ site.Params.foo }}|{{ range .Pages }}{{ .Title }},{{ end }}|
== layouts/page.html ==
Page|{{ .Title }}|{{ .Content }}|{{ partial "p.html" . }}
== layouts/_partials/p.html ==
Partial from txtar
== static/txtar.txt ==
Static from txtar
== data/mydata.toml ==
v = "data from txtar"
-- content/p2.md --
---
title: "P2 from disk"
---
-- themes/mytheme/layouts/_partials/p.html --
Partial from theme
-- themes/mytheme/layouts/section.html --
Section from theme
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Home|Home|foo from txtar|P1,P2 from disk,|")
	b.AssertFileContent("public/p1/index.html", "Page|P1|<p>P1 content.</p>", "Partial from txtar")
	// Files on disk take precedence over files in the archive.
	b.AssertFileContent("public/p2/index.html", "Page|P2 from disk|")
	static, err := afero.ReadFile(b.H.BaseFs.StaticFs(""), "txtar.txt")
	b.Assert(err, qt.IsNil)
	b.Assert(string(static), qt.Equals, "Static from txtar\n")
}

// See issue 10363.
func TestTxtarInProjectRootRebuild(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term", "rss", "sitemap", "home"]
-- hugo.txtar --
== content/p1.md ==
---
title: "P1"
---
P1 content.
== content/p2.md ==
---
title: "P2"
---
== layouts/page.html ==
Page|{{ .Title }}|{{ .Content }}|{{ partial "p.html" . }}
== layouts/_partials/p.html ==
Partial v1
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/p1/index.html", "Page|P1|<p>P1 content.</p>", "Partial v1")
	b.AssertFileContent("public/p2/index.html", "Page|P2|")

	b.EditFileReplaceFunc("hugo.txtar", func(s string) string {
		s = strings.ReplaceAll(s, "P1 content.", "P1 content edited.")
		s = strings.ReplaceAll(s, "Partial v1", "Partial v2")
		s = strings.ReplaceAll(s, "-- content/p2.md --", "-- content/p3.md --")
		return s
	}).Build()

	b.AssertFileContent("public/p1/index.html", "Page|P1|<p>P1 content edited.</p>", "Partial v2")
	b.AssertFileContent("public/p3/index.html", "Page|P2|")
	b.Assert(b.H.Sites[0].RegularPages().Len(), qt.Equals, 2)
	b.Assert(b.H.GetContentPage(filepath.FromSlash("/content/p2.md")), qt.IsNil)
}
