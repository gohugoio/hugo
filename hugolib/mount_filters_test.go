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

package hugolib

import (
	"testing"
)

func TestMountFilters(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
[module]
[[module.mounts]]
source = 'content'
target = 'content'
excludeFiles = "/a/c/**"
[[module.mounts]]
source = 'static'
target = 'static'
[[module.mounts]]
source = 'layouts'
target = 'layouts'
excludeFiles = "/**/foo.html"
[[module.mounts]]
source = 'data'
target = 'data'
includeFiles = "/mydata/**"
[[module.mounts]]
source = 'assets'
target = 'assets'
excludeFiles = ["/**exclude.*", "/moooo.*"]
[[module.mounts]]
source = 'i18n'
target = 'i18n'
[[module.mounts]]
source = 'archetypes'
target = 'archetypes'
-- layouts/_default/single.html --
Single page.
-- content/a/b/p1.md --
---
title: Include
---
-- content/a/c/p2.md --
---
title: Exclude
---
-- data/mydata/b.toml --
b1='bval'
-- data/nodata/c.toml --
c1='bval'
-- layouts/partials/foo.html --
foo
-- assets/exclude.txt --
foo
-- assets/js/exclude.js --
foo
-- assets/js/include.js --
foo
-- layouts/index.html --
Data: {{ site.Data }}:END

Template: {{ templates.Exists "partials/foo.html" }}:END
Resource1: {{ resources.Get "js/include.js" }}:END
Resource2: {{ resources.Get "js/exclude.js" }}:END
Resource3: {{ resources.Get "exclude.txt" }}:END
Resources: {{ resources.Match "**.js" }}
`
	b := Test(t, files)

	b.AssertFileExists("public/a/b/p1/index.html", true)
	b.AssertFileExists("public/a/c/p2/index.html", false)

	b.AssertFileContent("public/index.html", `
Data: map[mydata:map[b:map[b1:bval]]]:END
Template: false
Resource1: /js/include.js:END
Resource2: :END
Resource3: :END
Resources: [/js/include.js]
`)
}
