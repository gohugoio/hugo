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

// Issue 9599
func TestReadDirWorkDir(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
theme = "mytheme"
-- myproject.txt --
Hello project!
-- themes/mytheme/mytheme.txt --
Hello theme!
-- layouts/index.html --
{{ $entries := (readDir ".") }}
START:|{{ range $entry := $entries }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:


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

// Issue 9620
func TestReadFileNotExists(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- layouts/index.html --
{{ $fi := (readFile "doesnotexist") }}
{{ if $fi }}Failed{{ else }}OK{{ end }}


  `

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
OK
`)
}

func TestReadDirMountsVirtualDirectorySizeIsZero(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
[module]
[[module.imports]]
path = "module1"
[[module.imports.mounts]]
source = "assets"
target = "assets/virtual"
[[module.imports.mounts]]
source = "assets/file.json"
target = "assets/testing.json"
-- themes/module1/assets/file.json --
{}
-- layouts/index.html --
{{ $entries := readDir "assets" false }}
START:|{{ range $entry := $entries }}{{ $entry.Name }}={{ $entry.Size }}|{{ end }}:END:
{{ $entries = readDir "assets/virtual" false }}
START:|{{ range $entry := $entries }}{{ $entry.Name }}={{ $entry.Size }}|{{ end }}:END:
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
START:|file.json=2|virtual=0|:END:
`)
}

func TestReadDirMountsTopDirectory(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
[module]
	[[module.imports]]
		path = "module1"
		[[module.imports.mounts]]
			source = "assets"
			target = "assets"

-- myproject.txt --
Hello project!
-- themes/module1/assets/file.json --
{}
-- layouts/index.html --
{{ $entries := (readDir "" false) }}
START:|{{ range $entry := $entries }}{{ $entry.Name }}|{{ end }}:END:
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
START:|assets|hugo.toml|layouts|myproject.txt|themes|:END:
`)
}

func TestReadDirMergeContents(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
[module]
[[module.imports]]
path = "module1"
[[module.imports]]
path = "module2"

-- myproject.txt --
Hello project!
-- themes/module1/assets/file1.json --
{}
-- themes/module2/assets/files/raw.txt --
Nothing here
-- themes/module2/assets/file2.json --
{}
-- layouts/index.html --
{{ $entries := (readDir "assets" false) }}
START:|{{ range $entry := $entries }}{{ $entry.Name }}|{{ end }}:END:
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
START:|file1.json|file2.json|files|:END:
`)
}
