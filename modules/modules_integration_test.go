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

func TestModuleImportWithVersion(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
[[module.imports]]
path    = "github.com/bep/hugo-mod-misc/dummy-content"
version = "v0.2.0"
[[module.imports]]
path    = "github.com/bep/hugo-mod-misc/dummy-content"
version = "v0.1.0"
[[module.imports.mounts]]
source = "content"
target = "content/v1"
-- layouts/all.html --
Title: {{ .Title }}|Summary: {{ .Summary }}|
Deps: {{ range hugo.Deps}}{{ printf "%s@%s" .Path .Version }}|{{ end }}$


`

	b := hugolib.Test(t, files, hugolib.TestOptWithOSFs()).Build()

	b.AssertFileContent("public/index.html", "Deps: project@|github.com/bep/hugo-mod-misc/dummy-content@v0.2.0|github.com/bep/hugo-mod-misc/dummy-content@v0.1.0|$")

	b.AssertFileContent("public/blog/music/autumn-leaves/index.html", "Autumn Leaves is a popular jazz standard") // v0.2.0
	b.AssertFileContent("public/v1/blog/music/autumn-leaves/index.html", "Lorem markdownum, placidi peremptis")   // v0.1.0
}
