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

package modules_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestModulesFiltersFiles(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[module]
[[module.mounts]]
source = "content"
target = "content"
files = ["! **b.md", "**.md"] # The order matters here.
-- content/a.md --
+++
title = "A"
+++
Content A
-- content/b.md --
+++
title = "B"
+++
Content B
-- content/c.md --
+++
title = "C"
+++
Content C
-- layouts/all.html --
All. {{ .Title }}|

`

	b := hugolib.Test(t, files, hugolib.TestOptInfo())
	b.AssertLogContains("! deprecated")
	b.AssertFileContent("public/a/index.html", "All. A|")
	b.AssertFileContent("public/c/index.html", "All. C|")
	b.AssertFileExists("public/b/index.html", false)
}

// File filter format <= 0.152.0.
func TestModulesFiltersFilesLegacy(t *testing.T) {
	// This cannot be parallel.

	files := `
-- hugo.toml --
[module]
[[module.mounts]]
source = "content"
target = "content"
includefiles = ["**{a,c}.md"] #includes was evaluated before excludes <= 0.152.0
excludefiles = ["**b.md"]
-- content/a.md --
+++
title = "A"
+++
Content A
-- content/b.md --
+++
title = "B"
+++
Content B
-- content/c.md --
+++
title = "C"
+++
Content C
-- layouts/all.html --
All. {{ .Title }}|

`

	b := hugolib.Test(t, files, hugolib.TestOptInfo())

	b.AssertFileContent("public/a/index.html", "All. A|")
	b.AssertFileContent("public/c/index.html", "All. C|")
	b.AssertFileExists("public/b/index.html", false)
}
