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
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContent := `
	PaginatePath = "side"
	`

	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "hugo.toml", configContent)

	cfg, _, err := LoadConfig(ConfigSourceDescriptor{Fs: mm, Filename: "hugo.toml"})
	c.Assert(err, qt.IsNil)

	c.Assert(cfg.GetString("paginatePath"), qt.Equals, "side")

}

func TestLoadMultiConfig(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

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
	c.Assert(err, qt.IsNil)

	c.Assert(cfg.GetString("paginatePath"), qt.Equals, "top")
	c.Assert(cfg.GetString("DontChange"), qt.Equals, "same")
}

func TestLoadConfigFromTheme(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

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

	c.Assert(got["baseurl"], qt.Equals, "https://example.com/")

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

	c.Assert(got["languages"], qt.IsNil)
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

	c := qt.New(t)

	tomlConfig := `

someOtherValue = "foo"

[privacy]
[privacy.youtube]
privacyEnhanced = true
`

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", tomlConfig)
	b.Build(BuildCfg{SkipRender: true})

	c.Assert(b.H.Sites[0].Info.Config().Privacy.YouTube.PrivacyEnhanced, qt.Equals, true)

}

func TestLoadConfigModules(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	// https://github.com/gohugoio/hugoThemes#themetoml

	const (
		// Before Hugo 0.56 each theme/component could have its own theme.toml
		// with some settings, mostly used on the Hugo themes site.
		// To preserve combability we read these files into the new "modules"
		// section in config.toml.
		o1t = `
name = "Component o1"
license = "MIT"
min_version = 0.38
`
		// This is the component's config.toml, using the old theme syntax.
		o1c = `
theme = ["n2"]
`

		n1 = `
title = "Component n1"

[module]
description = "Component n1 description"
[module.hugoVersion]
min = "0.40.0"
max = "0.50.0"
extended = true
[[module.imports]]
path="o1"
[[module.imports]]
path="n3"


`

		n2 = `
title = "Component n2"
`

		n3 = `
title = "Component n3"
`

		n4 = `
title = "Component n4"
`
	)

	b := newTestSitesBuilder(t)

	writeThemeFiles := func(name, configTOML, themeTOML string) {
		b.WithSourceFile(filepath.Join("themes", name, "data", "module.toml"), fmt.Sprintf("name=%q", name))
		if configTOML != "" {
			b.WithSourceFile(filepath.Join("themes", name, "config.toml"), configTOML)
		}
		if themeTOML != "" {
			b.WithSourceFile(filepath.Join("themes", name, "theme.toml"), themeTOML)
		}
	}

	writeThemeFiles("n1", n1, "")
	writeThemeFiles("n2", n2, "")
	writeThemeFiles("n3", n3, "")
	writeThemeFiles("n4", n4, "")
	writeThemeFiles("o1", o1c, o1t)

	b.WithConfigFile("toml", `
[module]
[[module.imports]]
path="n1"
[[module.imports]]
path="n4"

`)

	b.Build(BuildCfg{})

	modulesClient := b.H.Paths.ModulesClient
	var graphb bytes.Buffer
	modulesClient.Graph(&graphb)

	expected := `project n1
n1 o1
o1 n2
n1 n3
project n4
`

	c.Assert(graphb.String(), qt.Equals, expected)

}

func TestLoadConfigWithOsEnvOverrides(t *testing.T) {

	c := qt.New(t)

	baseConfig := `

environment = "production"
enableGitInfo = true
intSlice = [5,7,9]
floatSlice = [3.14, 5.19]
stringSlice = ["a", "b"]

[imaging]
anchor = "smart"
quality = 75 
resamplefilter = "CatmullRom"
`

	b := newTestSitesBuilder(t).WithConfigFile("toml", baseConfig)

	b.WithEnviron(
		"HUGO_ENVIRONMENT", "test",
		"HUGO_NEW", "new", // key not in config.toml
		"HUGO_ENABLEGITINFO", "false",
		"HUGO_IMAGING_ANCHOR", "top",
		"HUGO_STRINGSLICE", `["c", "d"]`,
		"HUGO_INTSLICE", `[5, 8, 9]`,
		"HUGO_FLOATSLICE", `[5.32]`,
	)

	b.Build(BuildCfg{})

	cfg := b.H.Cfg

	c.Assert(cfg.Get("environment"), qt.Equals, "test")
	c.Assert(cfg.GetBool("enablegitinfo"), qt.Equals, false)
	c.Assert(cfg.Get("new"), qt.Equals, "new")
	c.Assert(cfg.Get("imaging.anchor"), qt.Equals, "top")
	c.Assert(cfg.Get("imaging.quality"), qt.Equals, int64(75))
	c.Assert(cfg.Get("stringSlice"), qt.DeepEquals, []interface{}{"c", "d"})
	c.Assert(cfg.Get("floatSlice"), qt.DeepEquals, []interface{}{5.32})
	c.Assert(cfg.Get("intSlice"), qt.DeepEquals, []interface{}{5, 8, 9})

}
