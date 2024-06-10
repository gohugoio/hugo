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

package hugo_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestIsMultilingualAndIsMultihost(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
defaultContentLanguageInSubdir = true
[languages.de]
baseURL = 'https://de.example.org/'
[languages.en]
baseURL = 'https://en.example.org/'
-- content/_index.md --
---
title: home
---
-- layouts/index.html --
multilingual={{ hugo.IsMultilingual }}
multihost={{ hugo.IsMultihost }}
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/de/index.html",
		"multilingual=true",
		"multihost=true",
	)
	b.AssertFileContent("public/en/index.html",
		"multilingual=true",
		"multihost=true",
	)

	files = strings.ReplaceAll(files, "baseURL = 'https://de.example.org/'", "")
	files = strings.ReplaceAll(files, "baseURL = 'https://en.example.org/'", "")

	b = hugolib.Test(t, files)

	b.AssertFileContent("public/de/index.html",
		"multilingual=true",
		"multihost=false",
	)
	b.AssertFileContent("public/en/index.html",
		"multilingual=true",
		"multihost=false",
	)

	files = strings.ReplaceAll(files, "[languages.de]", "")
	files = strings.ReplaceAll(files, "[languages.en]", "")

	b = hugolib.Test(t, files)

	b.AssertFileContent("public/en/index.html",
		"multilingual=false",
		"multihost=false",
	)
}
