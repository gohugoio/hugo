// Copyright 2021 The Hugo Authors. All rights reserved.
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

func TestDateFormatMultilingual(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `
baseURL = "https://example.org"

defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true

[languages]
[languages.en]
weight=10
[languages.nn]
weight=20
	
`)

	pageWithDate := `---
title: Page
date: 2021-07-18
---	
`

	b.WithContent(
		"_index.en.md", pageWithDate,
		"_index.nn.md", pageWithDate,
	)

	b.WithTemplatesAdded("index.html", `
Date: {{ .Date | time.Format ":date_long" }}
	`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/en/index.html", `Date: July 18, 2021`)
	b.AssertFileContent("public/nn/index.html", `Date: 18. juli 2021`)

}
