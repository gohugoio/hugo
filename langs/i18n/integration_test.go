// Copyright 2022 The Hugo Authors. All rights reserved.
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

package i18n_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestI18nFromTheme(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[module]
[[module.imports]]
path = "mytheme"
-- i18n/en.toml --
[l1]
other = 'l1main'
[l2]
other = 'l2main'
-- themes/mytheme/i18n/en.toml --
[l1]
other = 'l1theme'
[l2]
other = 'l2theme'
[l3]
other = 'l3theme'
-- layouts/index.html --
l1: {{ i18n "l1"  }}|l2: {{ i18n "l2"  }}|l3: {{ i18n "l3"  }}

`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
l1: l1main|l2: l2main|l3: l3theme
	`)
}

func TestHasLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = "https://example.org"
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
weight=10
[languages.nn]
weight=20
-- i18n/en.toml --
key1.other = "en key1"
key2.other = "en key2"

-- i18n/nn.toml --
key1.other = "nn key1"
key3.other = "nn key2"
-- layouts/index.html --
key1: {{ lang.HasTranslation "key1" }}|
key2: {{ lang.HasTranslation "key2" }}|
key3: {{ lang.HasTranslation "key3" }}|

  `

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/en/index.html", "key1: true|\nkey2: true|\nkey3: false|")
	b.AssertFileContent("public/nn/index.html", "key1: true|\nkey2: false|\nkey3: true|")
}
