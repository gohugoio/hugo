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
	"testing"
)

// Issue 14921
func TestResourceMetaCounterPlaceholder(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
-- content/p1/index.md --
---
title: p1
resources:
  - src: "*specs.pdf"
    title: "Specification #:counter"
  - src: "**.pdf"
    name: "pdf-file-:counter.pdf"
---
-- content/p1/checklist.pdf --
-- content/p1/guide.pdf --
-- content/p1/other_specs.pdf --
-- content/p1/photo_specs.pdf --
-- layouts/page.html --
<ul>
	{{- range seq 1 4 }}
		{{- with $.Resources.Get (printf "pdf-file-%d.pdf" .) }}
			<li><a href="{{ .RelPermalink }}">{{ .Title }}</a></li>
		{{- end }}
	{{- end }}
</ul>
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		`<a href="/p1/checklist.pdf">checklist.pdf</a>`,
		`<a href="/p1/guide.pdf">guide.pdf</a>`,
		`<a href="/p1/other_specs.pdf">Specification #1</a>`,
		`<a href="/p1/photo_specs.pdf">Specification #2</a>`,
	)
}
