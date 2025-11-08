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
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestI18n(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name     string
		langCode string
	}{
		{
			name:     "pt-br lowercase",
			langCode: "pt-br",
		},
		{
			name:     "pt-br uppercase",
			langCode: "PT-BR",
		},
	}

	for _, tc := range testCases {
		tc := tc
		c.Run(tc.name, func(c *qt.C) {
			files := fmt.Sprintf(`
-- hugo.toml --
baseURL = "https://example.com"
defaultContentLanguage = "%s"

[languages]
[languages.%s]
weight = 1
-- i18n/%s.toml --
hello.one = "Hello"
-- layouts/index.html --
Hello: {{ i18n "hello" 1 }}
-- content/p1.md --
`, tc.langCode, tc.langCode, tc.langCode)

			b := Test(c, files)
			b.AssertFileContent("public/index.html", "Hello: Hello")
		})
	}
}

func TestLanguageBugs(t *testing.T) {
	c := qt.New(t)

	// Issue #8672
	c.Run("Config with language, menu in root only", func(c *qt.C) {
		files := `
-- hugo.toml --
theme = "test-theme"
[[menus.foo]]
name = "foo-a"
[languages.en]
-- themes/test-theme/hugo.toml --
[languages.en]
`
		b := Test(c, files)

		menus := b.H.Sites[0].Menus()
		c.Assert(menus, qt.HasLen, 1)
	})
}

func TestLanguageNumberFormatting(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org"

defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true

[languages]
[languages.en]
timeZone="UTC"
weight=10
[languages.nn]
weight=20
-- layouts/index.html --

FormatNumber: {{ 512.5032 | lang.FormatNumber 2 }}
FormatPercent: {{ 512.5032 | lang.FormatPercent 2 }}
FormatCurrency: {{ 512.5032 | lang.FormatCurrency 2 "USD" }}
FormatAccounting: {{ 512.5032 | lang.FormatAccounting 2 "NOK" }}
FormatNumberCustom: {{ lang.FormatNumberCustom 2 12345.6789 }}
-- content/p1.md --
`

	b := Test(t, files)

	b.AssertFileContent("public/en/index.html", `
FormatNumber: 512.50
FormatPercent: 512.50%
FormatCurrency: $512.50
FormatAccounting: NOK512.50
FormatNumberCustom: 12,345.68
`,
	)

	b.AssertFileContent("public/nn/index.html", `
FormatNumber: 512,50
FormatPercent: 512,50 %
FormatCurrency: 512,50 USD
FormatAccounting: 512,50 kr
FormatNumberCustom: 12,345.68

`)
}

// Issue 11993.
func TestI18nDotFile(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
-- i18n/.keep --
-- data/.keep --
`
	Test(t, files)
}
