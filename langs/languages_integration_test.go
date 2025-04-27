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

package langs_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestLanguagesContentSimple(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 2
[languages.nn]
weight = 1
-- content/_index.md --
---
title: "Home"
---
Welcome to the home page.
-- content/_index.nn.md --
---
title: "Heim"
---
Welkomen heim!
-- layouts/all.html --
title: {{ .Title }}|
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/index.html", `title: Home|`)
	b.AssertFileContent("public/nn/index.html", `title: Heim|`)
}
