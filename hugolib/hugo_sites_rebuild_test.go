// Copyright 2019 The Hugo Authors. All rights reserved.
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

func TestSitesRebuild(t *testing.T) {

	configFile := `
baseURL = "https://example.com"
title = "Rebuild this"
contentDir = "content"


`

	contentFilename := "content/blog/page1.md"

	b := newTestSitesBuilder(t).WithConfigFile("toml", configFile)

	// To simulate https://github.com/gohugoio/hugo/issues/5838, the home page
	// needs a content page.
	b.WithContent("content/_index.md", `---
title: Home, Sweet Home!
---

`)

	b.WithContent(contentFilename, `
---
title: "Page 1"
summary: "Initial summary"
paginate: 3
---

Content.

`)

	b.WithTemplatesAdded("index.html", `
{{ range (.Paginate .Site.RegularPages).Pages }}
* Page: {{ .Title }}|Summary: {{ .Summary }}|Content: {{ .Content }}
{{ end }}
`)

	b.Running().Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "* Page: Page 1|Summary: Initial summary|Content: <p>Content.</p>")

	b.EditFiles(contentFilename, `
---
title: "Page 1 edit"
summary: "Edited summary"
---

Edited content.

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "* Page: Page 1 edit|Summary: Edited summary|Content: <p>Edited content.</p>")

}
