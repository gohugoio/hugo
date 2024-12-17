// Copyright 2024 The Hugo Authors. All rights reserved.
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

package os_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)


func TestReadDirMountDir2(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
theme = "mytheme"
[module]
	[[module.imports]]
		path = "module1"
		[[module.imports.mounts]]
			source = "assets"
			target = "assets"
	[[module.imports]]
		path = "mytheme"
		[[module.imports.mounts]]
			source = "private"
			target = "assets"
		[[module.imports.mounts]]
			source = "private/test"
			target = "assets"

-- myproject.txt --
Hello project!
-- themes/module1/hugo.toml --
-- themes/module1/go.mod --
module github.com/rymut/hugo-issues-mre/hugo-os/hugo-os-module2

go 1.23.2
-- themes/module1/assets/file.json --
{}
-- themes/mytheme/mytheme.txt --
Hello theme!
-- themes/mytheme/data/not.my_content.md --
test
-- themes/mytheme/layouts/partials/mypartial.html --
test
-- themes/mytheme/private/test/should_not_mount.txt --
Empty file
-- layouts/index.html --
{{ $entries := (readDir "" false) }}
START:|{{ range $entry := $entries }}{{ $entry.Name }}|{{ end }}:END:
-- files/layouts/l1.txt --
l1
-- files/layouts/assets/l2.txt --
l2
	`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
START:|config.toml|myproject.txt|:END:
`)
}
