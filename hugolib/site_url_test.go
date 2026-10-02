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
	"strings"
	"testing"
)

func TestSectionsEntries(t *testing.T) {
	files := `
-- hugo.toml --
-- content/withfile/_index.md --
-- content/withoutfile/p1.md --
-- layouts/list.html --
SectionsEntries: {{ .SectionsEntries }}


`

	b := Test(t, files)

	b.AssertFileContent("public/withfile/index.html", "SectionsEntries: [withfile]")
	b.AssertFileContent("public/withoutfile/index.html", "SectionsEntries: [withoutfile]")
}

// See #14625.
func TestDefaultBaseURL(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
CONFIG
-- content/_index.md --
-- layouts/list.html --
BaseURL: {{ .Site.BaseURL }}
`
	files := strings.ReplaceAll(filesTemplate, "CONFIG", `baseURL = "https://example.com/"`)
	b := Test(t, files)
	b.AssertFileContent("public/index.html", "BaseURL: https://example.com/")

	files = strings.ReplaceAll(filesTemplate, "CONFIG", "")
	b = Test(t, files)
	b.AssertFileContent("public/index.html", "BaseURL: https://example.org/")

	files = strings.ReplaceAll(filesTemplate, "CONFIG", `baseURL = "/"`)
	b = Test(t, files)
	b.AssertFileContent("public/index.html", "BaseURL: /")
}
