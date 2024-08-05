// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	qt "github.com/frankban/quicktest"
)

func TestI18n(t *testing.T) {
	c := qt.New(t)

	// https://github.com/gohugoio/hugo/issues/7804
	c.Run("pt-br should be case insensitive", func(c *qt.C) {
		b := newTestSitesBuilder(c)
		langCode := func() string {
			c := "pt-br"
			if htesting.RandBool() {
				c = strings.ToUpper(c)
			}
			return c
		}

		b.WithConfigFile(`toml`, fmt.Sprintf(`
baseURL = "https://example.com"
defaultContentLanguage = "%s"

[languages]
[languages.%s]
weight = 1
`, langCode(), langCode()))

		b.WithI18n(fmt.Sprintf("i18n/%s.toml", langCode()), `hello.one = "Hello"`)
		b.WithTemplates("index.html", `Hello: {{ i18n "hello" 1 }}`)
		b.WithContent("p1.md", "")
		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", "Hello: Hello")
	})
}

func TestLanguageBugs(t *testing.T) {
	c := qt.New(t)

	// Issue #8672
	c.Run("Config with language, menu in root only", func(c *qt.C) {
		b := newTestSitesBuilder(c)
		b.WithConfigFile("toml", `
theme = "test-theme"
[[menus.foo]]
name = "foo-a"
[languages.en]

`,
		)

		b.WithThemeConfigFile("toml", `[languages.en]`)

		b.Build(BuildCfg{})

		menus := b.H.Sites[0].Menus()
		c.Assert(menus, qt.HasLen, 1)
	})
}

func TestLanguageNumberFormatting(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `
baseURL = "https://example.org"

defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true

[languages]
[languages.en]
timeZone="UTC"
weight=10
[languages.nn]
weight=20
	
`)

	b.WithTemplates("index.html", `

FormatNumber: {{ 512.5032 | lang.FormatNumber 2 }}
FormatPercent: {{ 512.5032 | lang.FormatPercent 2 }}
FormatCurrency: {{ 512.5032 | lang.FormatCurrency 2 "USD" }}
FormatAccounting: {{ 512.5032 | lang.FormatAccounting 2 "NOK" }}
FormatNumberCustom: {{ lang.FormatNumberCustom 2 12345.6789 }}



	
`)
	b.WithContent("p1.md", "")

	b.Build(BuildCfg{})

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
-- hugo.toml --{}
baseURL = "https://example.com"
-- i18n/.keep --
-- data/.keep --
`
	Test(t, files)
}
