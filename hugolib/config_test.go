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
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/media"
	"github.com/google/go-cmp/cmp"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/spf13/afero"
)

func TestLoadConfig(t *testing.T) {

	c := qt.New(t)

	loadConfig := func(c *qt.C, configContent string, fromDir bool) config.Provider {
		mm := afero.NewMemMapFs()
		filename := "config.toml"
		descriptor := ConfigSourceDescriptor{Fs: mm}
		if fromDir {
			filename = filepath.Join("config", "_default", filename)
			descriptor.AbsConfigDir = "config"
		}
		writeToFs(t, mm, filename, configContent)
		cfg, _, err := LoadConfig(descriptor)
		c.Assert(err, qt.IsNil)
		return cfg
	}

	c.Run("Basic", func(c *qt.C) {
		c.Parallel()
		// Add a random config variable for testing.
		// side = page in Norwegian.
		cfg := loadConfig(c, `PaginatePath = "side"`, false)
		c.Assert(cfg.GetString("paginatePath"), qt.Equals, "side")
	})

	// Issue #8763
	for _, fromDir := range []bool{false, true} {
		testName := "Taxonomy overrides"
		if fromDir {
			testName += " from dir"
		}
		c.Run(testName, func(c *qt.C) {
			c.Parallel()
			cfg := loadConfig(c, `[taxonomies]
appellation = "appellations"
vigneron = "vignerons"`, fromDir)

			c.Assert(cfg.Get("taxonomies"), qt.DeepEquals, maps.Params{
				"appellation": "appellations",
				"vigneron":    "vignerons",
			})
		})
	}
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

func TestLoadConfigFromThemes(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	mainConfigTemplate := `
theme = "test-theme"
baseURL = "https://example.com/"

[frontmatter]
date = ["date","publishDate"]

[params]
MERGE_PARAMS
p1 = "p1 main"
[params.b]
b1 = "b1 main"
[params.b.c]
bc1 = "bc1 main"

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
disableKinds = ["taxonomy", "term"]

# Can not be set in theme.
[frontmatter]
expiryDate = ["date"]

[params]
p1 = "p1 theme"
p2 = "p2 theme"
[params.b]
b1 = "b1 theme"
b2 = "b2 theme"
[params.b.c]
bc1 = "bc1 theme"
bc2 = "bc2 theme"
[params.b.c.d]
bcd1 = "bcd1 theme"

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

	buildForConfig := func(t testing.TB, mainConfig, themeConfig string) *sitesBuilder {
		b := newTestSitesBuilder(t)
		b.WithConfigFile("toml", mainConfig).WithThemeConfigFile("toml", themeConfig)
		return b.Build(BuildCfg{})
	}

	buildForStrategy := func(t testing.TB, s string) *sitesBuilder {
		mainConfig := strings.ReplaceAll(mainConfigTemplate, "MERGE_PARAMS", s)
		return buildForConfig(t, mainConfig, themeConfig)
	}

	c.Run("Merge default", func(c *qt.C) {
		b := buildForStrategy(c, "")

		got := b.Cfg.Get("").(maps.Params)

		// Issue #8866
		b.Assert(b.Cfg.Get("disableKinds"), qt.IsNil)

		b.Assert(got["params"], qt.DeepEquals, maps.Params{
			"b": maps.Params{
				"b1": "b1 main",
				"c": maps.Params{
					"bc1": "bc1 main",
					"bc2": "bc2 theme",
					"d":   maps.Params{"bcd1": string("bcd1 theme")},
				},
				"b2": "b2 theme",
			},
			"p2": "p2 theme",
			"p1": "p1 main",
		})

		b.Assert(got["mediatypes"], qt.DeepEquals, maps.Params{
			"text/m2": maps.Params{
				"suffixes": []interface{}{
					"m2theme",
				},
			},
			"text/m1": maps.Params{
				"suffixes": []interface{}{
					"m1main",
				},
			},
		})

		var eq = qt.CmpEquals(
			cmp.Comparer(func(m1, m2 media.Type) bool {
				if m1.SubType != m2.SubType {
					return false
				}
				return m1.FirstSuffix == m2.FirstSuffix
			}),
		)

		mediaTypes := b.H.Sites[0].mediaTypesConfig
		m1, _ := mediaTypes.GetByType("text/m1")
		m2, _ := mediaTypes.GetByType("text/m2")

		b.Assert(got["outputformats"], eq, maps.Params{
			"o1": maps.Params{
				"mediatype": m1,
				"basename":  "o1main",
			},
			"o2": maps.Params{
				"basename":  "o2theme",
				"mediatype": m2,
			},
		})

		b.Assert(got["languages"], qt.DeepEquals, maps.Params{
			"en": maps.Params{
				"languagename": "English",
				"params": maps.Params{
					"pl2": "p2-en-theme",
					"pl1": "p1-en-main",
				},
				"menus": maps.Params{
					"main": []interface{}{
						map[string]interface{}{
							"name": "menu-lang-en-main",
						},
					},
					"theme": []interface{}{
						map[string]interface{}{
							"name": "menu-lang-en-theme",
						},
					},
				},
			},
			"nb": maps.Params{
				"languagename": "Norsk",
				"params": maps.Params{
					"top": "top-nb-theme",
					"pl1": "p1-nb-main",
					"pl2": "p2-nb-theme",
				},
				"menus": maps.Params{
					"main": []interface{}{
						map[string]interface{}{
							"name": "menu-lang-nb-main",
						},
					},
					"theme": []interface{}{
						map[string]interface{}{
							"name": "menu-lang-nb-theme",
						},
					},
					"top": []interface{}{
						map[string]interface{}{
							"name": "menu-lang-nb-top",
						},
					},
				},
			},
		})

		c.Assert(got["baseurl"], qt.Equals, "https://example.com/")
	})

	c.Run("Merge shallow", func(c *qt.C) {
		b := buildForStrategy(c, fmt.Sprintf("_merge=%q", "shallow"))

		got := b.Cfg.Get("").(maps.Params)

		// Shallow merge, only add new keys to params.
		b.Assert(got["params"], qt.DeepEquals, maps.Params{
			"p1": "p1 main",
			"b": maps.Params{
				"b1": "b1 main",
				"c": maps.Params{
					"bc1": "bc1 main",
				},
			},
			"p2": "p2 theme",
		})
	})

	c.Run("Merge no params in project", func(c *qt.C) {
		b := buildForConfig(
			c,
			"baseURL=\"https://example.org\"\ntheme = \"test-theme\"\n",
			"[params]\np1 = \"p1 theme\"\n",
		)

		got := b.Cfg.Get("").(maps.Params)

		b.Assert(got["params"], qt.DeepEquals, maps.Params{
			"p1": "p1 theme",
		})
	})

	c.Run("Merge language no menus or params in project", func(c *qt.C) {
		b := buildForConfig(
			c,
			`
theme = "test-theme"
baseURL = "https://example.com/"

[languages]
[languages.en]
languageName = "English"

`,
			`
[languages]
[languages.en]
languageName = "EnglishTheme"

[languages.en.params]
p1="themep1"

[[languages.en.menus.main]]
name   = "menu-theme"
`,
		)

		got := b.Cfg.Get("").(maps.Params)

		b.Assert(got["languages"], qt.DeepEquals,
			maps.Params{
				"en": maps.Params{
					"languagename": "English",
					"menus": maps.Params{
						"main": []interface{}{
							map[string]interface{}{
								"name": "menu-theme",
							},
						},
					},
					"params": maps.Params{
						"p1": "themep1",
					},
				},
			},
		)
	})

	// Issue #8724
	for _, mergeStrategy := range []string{"none", "shallow"} {
		c.Run(fmt.Sprintf("Merge with sitemap config in theme, mergestrategy %s", mergeStrategy), func(c *qt.C) {

			smapConfigTempl := `[sitemap]
  changefreq = %q
  filename = "sitemap.xml"
  priority = 0.5`

			b := buildForConfig(
				c,
				fmt.Sprintf("_merge=%q\nbaseURL=\"https://example.org\"\ntheme = \"test-theme\"\n", mergeStrategy),
				"baseURL=\"http://example.com\"\n"+fmt.Sprintf(smapConfigTempl, "monthly"),
			)

			got := b.Cfg.Get("").(maps.Params)

			if mergeStrategy == "none" {
				b.Assert(got["sitemap"], qt.DeepEquals, maps.Params{
					"priority": int(-1),
					"filename": "sitemap.xml",
				})

				b.AssertFileContent("public/sitemap.xml", "schemas/sitemap")
			} else {
				b.Assert(got["sitemap"], qt.DeepEquals, maps.Params{
					"priority":   int(-1),
					"filename":   "sitemap.xml",
					"changefreq": "monthly",
				})

				b.AssertFileContent("public/sitemap.xml", "<changefreq>monthly</changefreq>")
			}

		})
	}

}

func TestLoadConfigFromThemeDir(t *testing.T) {
	t.Parallel()

	mainConfig := `
theme = "test-theme"

[params]
m1 = "mv1"	
`

	themeConfig := `
[params]
t1 = "tv1"	
t2 = "tv2"
`

	themeConfigDir := filepath.Join("themes", "test-theme", "config")
	themeConfigDirDefault := filepath.Join(themeConfigDir, "_default")
	themeConfigDirProduction := filepath.Join(themeConfigDir, "production")

	projectConfigDir := "config"

	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", mainConfig).WithThemeConfigFile("toml", themeConfig)
	b.Assert(b.Fs.Source.MkdirAll(themeConfigDirDefault, 0777), qt.IsNil)
	b.Assert(b.Fs.Source.MkdirAll(themeConfigDirProduction, 0777), qt.IsNil)
	b.Assert(b.Fs.Source.MkdirAll(projectConfigDir, 0777), qt.IsNil)

	b.WithSourceFile(filepath.Join(projectConfigDir, "config.toml"), `[params]
m2 = "mv2"
`)
	b.WithSourceFile(filepath.Join(themeConfigDirDefault, "config.toml"), `[params]
t2 = "tv2d"
t3 = "tv3d"
`)

	b.WithSourceFile(filepath.Join(themeConfigDirProduction, "config.toml"), `[params]
t3 = "tv3p"
`)

	b.Build(BuildCfg{})

	got := b.Cfg.Get("params").(maps.Params)

	b.Assert(got, qt.DeepEquals, maps.Params{
		"t3": "tv3p",
		"m1": "mv1",
		"t1": "tv1",
		"t2": "tv2d",
	})

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

theme = "mytheme"
environment = "production"
enableGitInfo = true
intSlice = [5,7,9]
floatSlice = [3.14, 5.19]
stringSlice = ["a", "b"]

[outputFormats]
[outputFormats.ofbase]
mediaType = "text/plain"

[params]
paramWithNoEnvOverride="nooverride"
[params.api_config]
api_key="default_key"
another_key="default another_key"

[imaging]
anchor = "smart"
quality = 75 
`

	newB := func(t testing.TB) *sitesBuilder {
		b := newTestSitesBuilder(t).WithConfigFile("toml", baseConfig)

		b.WithSourceFile("themes/mytheme/config.toml", `

[outputFormats]
[outputFormats.oftheme]
mediaType = "text/plain"
[outputFormats.ofbase]
mediaType = "application/xml"

[params]
[params.mytheme_section]
theme_param="themevalue"
theme_param_nooverride="nooverride"
[params.mytheme_section2]
theme_param="themevalue2"

`)

		return b
	}

	c.Run("Variations", func(c *qt.C) {

		b := newB(c)

		b.WithEnviron(
			"HUGO_ENVIRONMENT", "test",
			"HUGO_NEW", "new", // key not in config.toml
			"HUGO_ENABLEGITINFO", "false",
			"HUGO_IMAGING_ANCHOR", "top",
			"HUGO_IMAGING_RESAMPLEFILTER", "CatmullRom",
			"HUGO_STRINGSLICE", `["c", "d"]`,
			"HUGO_INTSLICE", `[5, 8, 9]`,
			"HUGO_FLOATSLICE", `[5.32]`,
			// Issue #7829
			"HUGOxPARAMSxAPI_CONFIGxAPI_KEY", "new_key",
			// Delimiters are case sensitive.
			"HUGOxPARAMSxAPI_CONFIGXANOTHER_KEY", "another_key",
			// Issue #8346
			"HUGOxPARAMSxMYTHEME_SECTIONxTHEME_PARAM", "themevalue_changed",
			"HUGOxPARAMSxMYTHEME_SECTION2xTHEME_PARAM", "themevalue2_changed",
			"HUGO_PARAMS_EMPTY", ``,
			"HUGO_PARAMS_HTML", `<a target="_blank" />`,
			// Issue #8618
			"HUGO_SERVICES_GOOGLEANALYTICS_ID", `gaid`,
			"HUGO_PARAMS_A_B_C", "abc",
		)

		b.Build(BuildCfg{})

		cfg := b.H.Cfg
		s := b.H.Sites[0]
		scfg := s.siteConfigConfig.Services

		c.Assert(cfg.Get("environment"), qt.Equals, "test")
		c.Assert(cfg.GetBool("enablegitinfo"), qt.Equals, false)
		c.Assert(cfg.Get("new"), qt.Equals, "new")
		c.Assert(cfg.Get("imaging.anchor"), qt.Equals, "top")
		c.Assert(cfg.Get("imaging.quality"), qt.Equals, int64(75))
		c.Assert(cfg.Get("imaging.resamplefilter"), qt.Equals, "CatmullRom")
		c.Assert(cfg.Get("stringSlice"), qt.DeepEquals, []interface{}{"c", "d"})
		c.Assert(cfg.Get("floatSlice"), qt.DeepEquals, []interface{}{5.32})
		c.Assert(cfg.Get("intSlice"), qt.DeepEquals, []interface{}{5, 8, 9})
		c.Assert(cfg.Get("params.api_config.api_key"), qt.Equals, "new_key")
		c.Assert(cfg.Get("params.api_config.another_key"), qt.Equals, "default another_key")
		c.Assert(cfg.Get("params.mytheme_section.theme_param"), qt.Equals, "themevalue_changed")
		c.Assert(cfg.Get("params.mytheme_section.theme_param_nooverride"), qt.Equals, "nooverride")
		c.Assert(cfg.Get("params.mytheme_section2.theme_param"), qt.Equals, "themevalue2_changed")
		c.Assert(cfg.Get("params.empty"), qt.Equals, ``)
		c.Assert(cfg.Get("params.html"), qt.Equals, `<a target="_blank" />`)

		params := cfg.Get("params").(maps.Params)
		c.Assert(params["paramwithnoenvoverride"], qt.Equals, "nooverride")
		c.Assert(cfg.Get("params.paramwithnoenvoverride"), qt.Equals, "nooverride")
		c.Assert(scfg.GoogleAnalytics.ID, qt.Equals, "gaid")
		c.Assert(cfg.Get("params.a.b"), qt.DeepEquals, maps.Params{
			"c": "abc",
		})

		ofBase, _ := s.outputFormatsConfig.GetByName("ofbase")
		ofTheme, _ := s.outputFormatsConfig.GetByName("oftheme")

		c.Assert(ofBase.MediaType, qt.Equals, media.TextType)
		c.Assert(ofTheme.MediaType, qt.Equals, media.TextType)

	})

	// Issue #8709
	c.Run("Set in string", func(c *qt.C) {
		b := newB(c)

		b.WithEnviron(
			"HUGO_ENABLEGITINFO", "false",
			// imaging.anchor is a string, and it's not possible
			// to set a child attribute.
			"HUGO_IMAGING_ANCHOR_FOO", "top",
		)

		b.Build(BuildCfg{})

		cfg := b.H.Cfg
		c.Assert(cfg.Get("imaging.anchor"), qt.Equals, "smart")

	})

}
