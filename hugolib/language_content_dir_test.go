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

func TestLanguageContentRoot(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 10
contentDir = "content/en"
[languages.nn]
weight = 20
contentDir = "content/nn"
-- content/en/_index.md --
---
title: "Home"
---
-- content/nn/_index.md --
---
title: "Heim"
---
-- content/en/myfiles/file1.txt --
file 1 en
-- content/en/myfiles/file2.txt --
file 2 en
-- content/nn/myfiles/file1.txt --
file 1 nn
-- layouts/index.html --
Title: {{ .Title }}|
Len Resources: {{ len .Resources }}|
{{ range $i, $e := .Resources }}
{{ $i }}|{{ .RelPermalink }}|{{ .Content }}|
{{ end }}

`
	b := Test(t, files)
	b.AssertFileContent("public/en/index.html", "Home", "0|/en/myfiles/file1.txt|file 1 en|\n\n1|/en/myfiles/file2.txt|file 2 en|")
	b.AssertFileContent("public/nn/index.html", "Heim", "0|/nn/myfiles/file1.txt|file 1 nn|\n\n1|/en/myfiles/file2.txt|file 2 en|")
}
