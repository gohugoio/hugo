// Copyright 2016-present The Hugo Authors. All rights reserved.
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

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContent := `
	PaginatePath = "side"
	`

	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "hugo.toml", configContent)

	cfg, _, err := LoadConfig(ConfigSourceDescriptor{Fs: mm, Filename: "hugo.toml"})
	require.NoError(t, err)

	assert.Equal("side", cfg.GetString("paginatePath"))
	// default
	assert.Equal("layouts", cfg.GetString("layoutDir"))
	// no themes
	assert.False(cfg.IsSet("allThemes"))
}

func TestLoadMultiConfig(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContentBase := `
	DontChange = "same"
	PaginatePath = "side"
	`
	configContentSub := `
	PaginatePath = "top"
	`
	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "base.toml", configContentBase)

	writeToFs(t, mm, "override.toml", configContentSub)

	cfg, _, err := LoadConfig(ConfigSourceDescriptor{Fs: mm, Filename: "base.toml,override.toml"})
	require.NoError(t, err)

	assert.Equal("top", cfg.GetString("paginatePath"))
	assert.Equal("same", cfg.GetString("DontChange"))
}

func TestLoadConfigFromTheme(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	mainConfigBasic := `
theme = "test-theme"
baseURL = "https://example.com/"

`
	mainConfig := `
theme = "test-theme"
baseURL = "https://example.com/"

[frontmatter]
date = ["date","publishDate"]

[params]
p1 = "p1 main"
p2 = "p2 main"
top = "top"

[mediaTypes]
[mediaTypes."text/m1"]
suffixes = ["m1main"]

[outputFormats.o1]
mediaType = "text/m1"
baseName = "o1main"

[languages]
[languages.en]
languageName = "English"
[languages.en.params]
pl1 = "p1-en-main"
[languages.nb]
languageName = "Norsk"
[languages.nb.params]
pl1 = "p1-nb-main"

[[menu.main]]
name = "menu-main-main"

[[menu.top]]
name = "menu-top-main"

`

	themeConfig := `
baseURL = "http://bep.is/"

# Can not be set in theme.
[frontmatter]
expiryDate = ["date"]

[params]
p1 = "p1 theme"
p2 = "p2 theme"
p3 = "p3 theme"

[mediaTypes]
[mediaTypes."text/m1"]
suffixes = ["m1theme"]
[mediaTypes."text/m2"]
suffixes = ["m2theme"]

[outputFormats.o1]
mediaType = "text/m1"
baseName = "o1theme"
[outputFormats.o2]
mediaType = "text/m2"
baseName = "o2theme"

[languages]
[languages.en]
languageName = "English2"
[languages.en.params]
pl1 = "p1-en-theme"
pl2 = "p2-en-theme"
[[languages.en.menu.main]]
name   = "menu-lang-en-main"
[[languages.en.menu.theme]]
name   = "menu-lang-en-theme"
[languages.nb]
languageName = "Norsk2"
[languages.nb.params]
pl1 = "p1-nb-theme"
pl2 = "p2-nb-theme"
top = "top-nb-theme"
[[languages.nb.menu.main]]
name   = "menu-lang-nb-main"
[[languages.nb.menu.theme]]
name   = "menu-lang-nb-theme"
[[languages.nb.menu.top]]
name   = "menu-lang-nb-top"

[[menu.main]]
name = "menu-main-theme"

[[menu.thememenu]]
name = "menu-theme"

`

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", mainConfig).WithThemeConfigFile("toml", themeConfig)
	b.CreateSites().Build(BuildCfg{})

	got := b.Cfg.(*viper.Viper).AllSettings()

	b.AssertObject(`
map[string]interface {}{
  "p1": "p1 main",
  "p2": "p2 main",
  "p3": "p3 theme",
  "test-theme": map[string]interface {}{
    "p1": "p1 theme",
    "p2": "p2 theme",
    "p3": "p3 theme",
  },
  "top": "top",
}`, got["params"])

	b.AssertObject(`
map[string]interface {}{
  "date": []interface {}{
    "date",
    "publishDate",
  },
}`, got["frontmatter"])

	b.AssertObject(`
map[string]interface {}{
  "text/m1": map[string]interface {}{
    "suffixes": []interface {}{
      "m1main",
    },
  },
  "text/m2": map[string]interface {}{
    "suffixes": []interface {}{
      "m2theme",
    },
  },
}`, got["mediatypes"])

	b.AssertObject(`
map[string]interface {}{
  "o1": map[string]interface {}{
    "basename": "o1main",
    "mediatype": Type{
      MainType: "text",
      SubType: "m1",
      Delimiter: ".",
      Suffixes: []string{
        "m1main",
      },
    },
  },
  "o2": map[string]interface {}{
    "basename": "o2theme",
    "mediatype": Type{
      MainType: "text",
      SubType: "m2",
      Delimiter: ".",
      Suffixes: []string{
        "m2theme",
      },
    },
  },
}`, got["outputformats"])

	b.AssertObject(`map[string]interface {}{
  "en": map[string]interface {}{
    "languagename": "English",
    "menus": map[string]interface {}{
      "theme": []map[string]interface {}{
        map[string]interface {}{
          "name": "menu-lang-en-theme",
        },
      },
    },
    "params": map[string]interface {}{
      "pl1": "p1-en-main",
      "pl2": "p2-en-theme",
      "test-theme": map[string]interface {}{
        "pl1": "p1-en-theme",
        "pl2": "p2-en-theme",
      },
    },
  },
  "nb": map[string]interface {}{
    "languagename": "Norsk",
    "menus": map[string]interface {}{
      "theme": []map[string]interface {}{
        map[string]interface {}{
          "name": "menu-lang-nb-theme",
        },
      },
    },
    "params": map[string]interface {}{
      "pl1": "p1-nb-main",
      "pl2": "p2-nb-theme",
      "test-theme": map[string]interface {}{
        "pl1": "p1-nb-theme",
        "pl2": "p2-nb-theme",
        "top": "top-nb-theme",
      },
    },
  },
}
`, got["languages"])

	b.AssertObject(`
map[string]interface {}{
  "main": []map[string]interface {}{
    map[string]interface {}{
      "name": "menu-main-main",
    },
  },
  "thememenu": []map[string]interface {}{
    map[string]interface {}{
      "name": "menu-theme",
    },
  },
  "top": []map[string]interface {}{
    map[string]interface {}{
      "name": "menu-top-main",
    },
  },
}
`, got["menus"])

	assert.Equal("https://example.com/", got["baseurl"])

	if true {
		return
	}
	// Test variants with only values from theme
	b = newTestSitesBuilder(t)
	b.WithConfigFile("toml", mainConfigBasic).WithThemeConfigFile("toml", themeConfig)
	b.CreateSites().Build(BuildCfg{})

	got = b.Cfg.(*viper.Viper).AllSettings()

	b.AssertObject(`map[string]interface {}{
  "p1": "p1 theme",
  "p2": "p2 theme",
  "p3": "p3 theme",
  "test-theme": map[string]interface {}{
    "p1": "p1 theme",
    "p2": "p2 theme",
    "p3": "p3 theme",
  },
}`, got["params"])

	assert.Nil(got["languages"])
	b.AssertObject(`
map[string]interface {}{
  "text/m1": map[string]interface {}{
    "suffix": "m1theme",
  },
  "text/m2": map[string]interface {}{
    "suffix": "m2theme",
  },
}`, got["mediatypes"])

	b.AssertObject(`
map[string]interface {}{
  "o1": map[string]interface {}{
    "basename": "o1theme",
    "mediatype": Type{
      MainType: "text",
      SubType: "m1",
      Suffix: "m1theme",
      Delimiter: ".",
    },
  },
  "o2": map[string]interface {}{
    "basename": "o2theme",
    "mediatype": Type{
      MainType: "text",
      SubType: "m2",
      Suffix: "m2theme",
      Delimiter: ".",
    },
  },
}`, got["outputformats"])
	b.AssertObject(`
map[string]interface {}{
  "main": []interface {}{
    map[string]interface {}{
      "name": "menu-main-theme",
    },
  },
  "thememenu": []interface {}{
    map[string]interface {}{
      "name": "menu-theme",
    },
  },
}`, got["menu"])

}

func TestPrivacyConfig(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	tomlConfig := `

someOtherValue = "foo"

[privacy]
[privacy.youtube]
privacyEnhanced = true
`

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", tomlConfig)
	b.Build(BuildCfg{SkipRender: true})

	assert.True(b.H.Sites[0].Info.Config().Privacy.YouTube.PrivacyEnhanced)

}
