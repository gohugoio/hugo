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

func TestGitInfoFromGitModuleWithVersionQuery(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
enableGitInfo = true

[module]
[[module.imports]]
path = "github.com/bep/hugo-testing-git-versions"
version = "v3.0.0"
[[module.imports.mounts]]
source = "content"
target = "content"
[[module.imports]]
path = "github.com/bep/hugo-testing-git-versions"
version = "v3.0.1"
[[module.imports.mounts]]
source = "content"
target = "content/301"
-- layouts/page.html --
Title: {{ .Title }}|
GitInfo: {{ with .GitInfo }}Hash: {{ .Hash }}|Subject: {{ .Subject }}|AuthorName: {{ .AuthorName }}{{ end }}|
Content: {{ .Content }}|
-- layouts/_default/list.html --
List: {{ .Title }}
`

	b := Test(t, files, TestOptOsFs())

	b.AssertFileContent("public/docs/functions/mean/index.html",
		"Hash: 9769b63d4a4abdd406e333be5fb8b5d48737d3a9|",
		"AuthorName: Bjørn Erik Pedersen|",
	)

	b.AssertFileContent("public/docs/functions/standard-deviation/index.html",
		"Hash: 2d92492a7f1ec4968529ee12cf62ed652eb45950|Subject: v3.0.0 edits|",
		"AuthorName: Bjørn Erik Pedersen|",
	)

	b.AssertFileContent("public/301/docs/functions/standard-deviation/index.html",
		"Hash: 3e0f3930f1ec9a29a7442da5f1bfc0b7e58f167a|Subject: v3.0.1 edit|",
		"AuthorName: Bjørn Erik Pedersen|",
	)
}

func TestGitInfoFromGitModuleWithGoMod(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
enableGitInfo = true

[module]
[[module.imports]]
path = "github.com/bep/hugo-mod-testing-content"
[[module.imports.mounts]]
source = "content"
target = "content"
-- layouts/page.html --
Title: {{ .Title }}|
GitInfo: {{ with .GitInfo }}Hash: {{ .Hash }}|Subject: {{ .Subject }}|AuthorName: {{ .AuthorName }}{{ end }}|
Content: {{ .Content }}|
-- layouts/_default/list.html --
List: {{ .Title }}
-- go.mod --
module hugotest
`

	b := Test(t, files, TestOptOsFs())

	b.AssertFileContent("public/docs/functions/kurtosis/index.html",
		"Hash: 668663b54d0937df05185d144765d13c3ffda489|",
		"AuthorName: Bjørn Erik Pedersen|",
	)
}
