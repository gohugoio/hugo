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
