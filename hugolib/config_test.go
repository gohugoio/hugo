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
	"github.com/gohugoio/hugo/config/allconfig"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/spf13/afero"
)

func TestLoadConfigLanguageParamsOverrideIssue10620(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "setion"]
title = "Base Title"
staticDir = "mystatic"
[params]
[params.comments]
color = "blue"
title = "Default Comments Title"
[languages]
[languages.en]
title = "English Title"
[languages.en.params.comments]
title = "English Comments Title"



`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	enSite := b.H.Sites[0]
	b.Assert(enSite.Title(), qt.Equals, "English Title")
	b.Assert(enSite.Home().Title(), qt.Equals, "English Title")
	b.Assert(enSite.Params(), qt.DeepEquals, maps.Params{
		"comments": maps.Params{
			"color": "blue",
			"title": "English Comments Title",
		},
	},
	)

}

func TestLoadConfig(t *testing.T) {

	t.Run("2 languages", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "setion"]
title = "Base Title"
staticDir = "mystatic"
[params]
p1 = "p1base"
p2 = "p2base"
[languages]
[languages.en]
title = "English Title"
[languages.en.params]
myparam = "enParamValue"
p1 = "p1en"
weight = 1
[languages.sv]
title = "Svensk Title"
staticDir = "mysvstatic"
weight = 2
[languages.sv.params]
myparam = "svParamValue"


`
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()

		enSite := b.H.Sites[0]
		svSite := b.H.Sites[1]
		b.Assert(enSite.Title(), qt.Equals, "English Title")
		b.Assert(enSite.Home().Title(), qt.Equals, "English Title")
		b.Assert(enSite.Params()["myparam"], qt.Equals, "enParamValue")
		b.Assert(enSite.Params()["p1"], qt.Equals, "p1en")
		b.Assert(enSite.Params()["p2"], qt.Equals, "p2base")
		b.Assert(svSite.Params()["p1"], qt.Equals, "p1base")
		b.Assert(enSite.conf.StaticDir[0], qt.Equals, "mystatic")

		b.Assert(svSite.Title(), qt.Equals, "Svensk Title")
		b.Assert(svSite.Home().Title(), qt.Equals, "Svensk Title")
		b.Assert(svSite.Params()["myparam"], qt.Equals, "svParamValue")
		b.Assert(svSite.conf.StaticDir[0], qt.Equals, "mysvstatic")

	})

	t.Run("disable default language", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "setion"]
title = "Base Title"
defaultContentLanguage = "sv"
disableLanguages = ["sv"]
[languages.en]
weight = 1
[languages.sv]
weight = 2
`
		b, err := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).BuildE()

		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, "cannot disable default content language")

	})

	t.Run("no internal config from outside", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
baseURL = "https://example.com"
[internal]
running = true
`
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()

		b.Assert(b.H.Conf.Running(), qt.Equals, false)

	})

	t.Run("env overrides", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "setion"]
title = "Base Title"
[params]
p1 = "p1base"
p2 = "p2base"
[params.pm2]
pm21 = "pm21base"
pm22 = "pm22base"
-- layouts/index.html --
p1: {{ .Site.Params.p1 }}
p2: {{ .Site.Params.p2 }}
pm21: {{ .Site.Params.pm2.pm21 }}
pm22: {{ .Site.Params.pm2.pm22 }}
pm31: {{ .Site.Params.pm3.pm31 }}



`
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
				Environ:     []string{"HUGO_PARAMS_P2=p2env", "HUGO_PARAMS_PM2_PM21=pm21env", "HUGO_PARAMS_PM3_PM31=pm31env"},
			},
		).Build()

		b.AssertFileContent("public/index.html", "p1: p1base\np2: p2env\npm21: pm21env\npm22: pm22base\npm31: pm31env")

	})

}

func TestLoadConfigThemeLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- /hugo.toml --
baseURL = "https://example.com"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
theme = "mytheme"
[languages]
[languages.en]
title = "English Title"
weight = 1
[languages.sv]
weight = 2
-- themes/mytheme/hugo.toml --
[params]
p1 = "p1base"
[languages]
[languages.en]
title = "English Title Theme"
[languages.en.params]
p2 = "p2en"
[languages.en.params.sub]
sub1 = "sub1en"
[languages.sv]
title = "Svensk Title Theme"
-- layouts/index.html --
title: {{ .Title }}|
p1: {{ .Site.Params.p1 }}|
p2: {{ .Site.Params.p2 }}|
sub: {{ .Site.Params.sub }}|
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/en/index.html", `
title: English Title|
p1: p1base
p2: p2en
sub: map[sub1:sub1en]
`)

}

func TestLoadMultiConfig(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContentBase := `
	Paginate = 32
	PaginatePath = "side"
	`
	configContentSub := `
	PaginatePath = "top"
	`
	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "base.toml", configContentBase)

	writeToFs(t, mm, "override.toml", configContentSub)

	all, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{Fs: mm, Filename: "base.toml,override.toml"})
	c.Assert(err, qt.IsNil)
	cfg := all.Base

	c.Assert(cfg.PaginatePath, qt.Equals, "top")
	c.Assert(cfg.Paginate, qt.Equals, 32)

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

		got := b.Configs.Base

		b.Assert(got.Params, qt.DeepEquals, maps.Params{
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

		c.Assert(got.BaseURL, qt.Equals, "https://example.com/")
	})

	c.Run("Merge shallow", func(c *qt.C) {
		b := buildForStrategy(c, fmt.Sprintf("_merge=%q", "shallow"))

		got := b.Configs.Base.Params

		// Shallow merge, only add new keys to params.
		b.Assert(got, qt.DeepEquals, maps.Params{
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

		got := b.Configs.Base.Params

		b.Assert(got, qt.DeepEquals, maps.Params{
			"p1": "p1 theme",
		})
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

			got := b.Configs.Base

			if mergeStrategy == "none" {
				b.Assert(got.Sitemap, qt.DeepEquals, config.SitemapConfig{ChangeFreq: "", Priority: -1, Filename: "sitemap.xml"})

				b.AssertFileContent("public/sitemap.xml", "schemas/sitemap")
			} else {
				b.Assert(got.Sitemap, qt.DeepEquals, config.SitemapConfig{ChangeFreq: "monthly", Priority: -1, Filename: "sitemap.xml"})
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

	got := b.Configs.Base.Params

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

	c.Assert(b.H.Sites[0].Config().Privacy.YouTube.PrivacyEnhanced, qt.Equals, true)
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

	modulesClient := b.H.Configs.ModulesClient
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

func TestInvalidDefaultMarkdownHandler(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup]
defaultMarkdownHandler = 'blackfriday'
-- content/_index.md --
## Foo
-- layouts/index.html --
{{ .Content }}

`

	b, err := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "Configured defaultMarkdownHandler \"blackfriday\" not found. Did you mean to use goldmark? Blackfriday was removed in Hugo v0.100.0.")

}

// Issue 8979
func TestHugoConfig(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
theme = "mytheme"
[params]
rootparam = "rootvalue"
-- config/_default/hugo.toml --
[params]
rootconfigparam = "rootconfigvalue"
-- themes/mytheme/config/_default/hugo.toml --
[params]
themeconfigdirparam = "themeconfigdirvalue"
-- themes/mytheme/hugo.toml --
[params]
themeparam = "themevalue"
-- layouts/index.html --
rootparam: {{ site.Params.rootparam }}
rootconfigparam: {{ site.Params.rootconfigparam }}
themeparam: {{ site.Params.themeparam }}
themeconfigdirparam: {{ site.Params.themeconfigdirparam }}


`

	for _, configName := range []string{"hugo.toml", "config.toml"} {
		configName := configName
		t.Run(configName, func(t *testing.T) {
			t.Parallel()

			files := strings.ReplaceAll(filesTemplate, "hugo.toml", configName)

			b, err := NewIntegrationTestBuilder(
				IntegrationTestConfig{
					T:           t,
					TxtarString: files,
				},
			).BuildE()

			b.Assert(err, qt.IsNil)
			b.AssertFileContent("public/index.html",
				"rootparam: rootvalue",
				"rootconfigparam: rootconfigvalue",
				"themeparam: themevalue",
				"themeconfigdirparam: themeconfigdirvalue",
			)

		})
	}

}

func TestConfigOutputFormatDefinedInTheme(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
theme = "mytheme"
[outputFormats]
[outputFormats.myotherformat]
baseName = 'myotherindex'
mediaType = 'text/html'
[outputs]
  home = ['myformat']
-- themes/mytheme/hugo.toml --
[outputFormats]
[outputFormats.myformat]
baseName = 'myindex'
mediaType = 'text/html'
-- layouts/index.html --
Home.



`

	b, err := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.IsNil)
	b.AssertFileContent("public/myindex.html", "Home.")

}

func TestConfigParamSetOnLanguageLevel(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT"]
[languages]
[languages.en]
title = "English Title"
thisIsAParam = "thisIsAParamValue"
[languages.en.params]
myparam = "enParamValue"
[languages.sv]
title = "Svensk Title"
[languages.sv.params]
myparam = "svParamValue"
-- layouts/index.html --
MyParam: {{ site.Params.myparam }}
ThisIsAParam: {{ site.Params.thisIsAParam }}


`

	b, err := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.IsNil)
	b.AssertFileContent("public/index.html", `		
MyParam: enParamValue
ThisIsAParam: thisIsAParamValue	
`)
}

func TestReproCommentsIn10947(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT"]
[languages]
[languages.en]
languageCode = "en-US"
title = "English Title"
[languages.en.params]
myparam = "enParamValue"
[languages.sv]
title = "Svensk Title"
[languages.sv.params]
myparam = "svParamValue"
-- content/mysection/_index.en.md --
---
title: "My English Section"
---
-- content/mysection/_index.sv.md --
---
title: "My Swedish Section"
---
-- layouts/index.html --
LanguageCode: {{ eq site.LanguageCode site.Language.LanguageCode }}|{{ site.Language.LanguageCode }}|
{{ range $i, $e := (slice site .Site) }}
{{ $i }}|AllPages: {{ len .AllPages }}|Sections: {{ if .Sections }}true{{ end }}| Author: {{ .Authors }}|BuildDrafts: {{ .BuildDrafts }}|IsMultiLingual: {{ .IsMultiLingual }}|Param: {{ .Language.Params.myparam }}|Language string: {{ .Language }}|Languages: {{ .Languages }}
{{ end }}



`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	{
		b.Assert(b.H.Log.LogCounters().WarnCounter.Count(), qt.Equals, uint64(2))
	}
	b.AssertFileContent("public/index.html", `
AllPages: 4|
Sections: true|
Param: enParamValue	
Param: enParamValue	
IsMultiLingual: true
LanguageCode: true|en-US|
`)

	b.AssertFileContent("public/sv/index.html", `
Param: svParamValue
LanguageCode: true|sv|

`)

}

func TestConfigEmptyMainSections(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.yml --
params:
  mainSections:
-- content/mysection/_index.md --
-- content/mysection/mycontent.md --
-- layouts/index.html --
mainSections: {{ site.Params.mainSections }}

`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
mainSections: []
`)

}

func TestConfigHugoWorkingDir(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/index.html --
WorkingDir: {{ hugo.WorkingDir }}|

`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			WorkingDir:  "myworkingdir",
		},
	).Build()

	b.AssertFileContent("public/index.html", `
WorkingDir: myworkingdir|
`)

}

func TestConfigMergeLanguageDeepEmptyLefSide(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[params]
p1 = "p1base"
[languages.en]
languageCode = 'en-US'
languageName = 'English'
weight = 1
[languages.en.markup.goldmark.extensions.typographer]
leftDoubleQuote = '&ldquo;'   # default &ldquo;
rightDoubleQuote = '&rdquo;'  # default &rdquo;

[languages.de]
languageCode = 'de-DE'
languageName = 'Deutsch'
weight = 2
[languages.de.params]
p1 = "p1de"
[languages.de.markup.goldmark.extensions.typographer]
leftDoubleQuote = '&laquo;'   # default &ldquo;
rightDoubleQuote = '&raquo;'  # default &rdquo;
-- layouts/index.html --
{{ .Content }}
p1: {{ site.Params.p1 }}|
-- content/_index.en.md --
---
title: "English Title"
---
A "quote" in English.
-- content/_index.de.md --
---
title: "Deutsch Title"
---
Ein "Zitat" auf Deutsch.



`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", "p1: p1base", "<p>A &ldquo;quote&rdquo; in English.</p>")
	b.AssertFileContent("public/de/index.html", "p1: p1de", "<p>Ein &laquo;Zitat&raquo; auf Deutsch.</p>")

}

func TestConfigLegacyValues(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
# taxonomyTerm was renamed to term in Hugo 0.60.0.
disableKinds = ["taxonomyTerm"]

-- layouts/index.html --
Home

`

	b, err := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.IsNil)
	b.AssertFileContent("public/index.html", `		
Home
`)

	conf := b.H.Configs.Base
	b.Assert(conf.IsKindEnabled("term"), qt.Equals, false)
}

// Issue #11000
func TestConfigEmptyTOMLString(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[mediaTypes]
[mediaTypes."text/htaccess"]
suffixes = ["htaccess"]
[outputFormats]
[outputFormats.htaccess]
mediaType = "text/htaccess"
baseName = ""
isPlainText = false
notAlternative = true
-- content/_index.md --
---
outputs: ["html", "htaccess"]
---
-- layouts/index.html --
HTML.
-- layouts/_default/list.htaccess --
HTACCESS.


	
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/.htaccess", "HTACCESS")

}

func TestConfigLanguageCodeTopLevel(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
languageCode = "en-US"
-- layouts/index.html --
LanguageCode: {{ .Site.LanguageCode }}|{{ site.Language.LanguageCode }}|

	
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", "LanguageCode: en-US|en-US|")

}

func TestConfigMiscPanics(t *testing.T) {
	t.Parallel()

	// Issue 11047,
	t.Run("empty params", func(t *testing.T) {

		files := `
-- hugo.yaml --
params:
-- layouts/index.html --
Foo: {{ site.Params.foo }}|
	
		
	`
		b := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).Build()

		b.AssertFileContent("public/index.html", "Foo: |")
	})

	// Issue 11046
	t.Run("invalid language setup", func(t *testing.T) {

		files := `
-- hugo.toml --
baseURL = "https://example.org"
languageCode = "en-us"
title = "Blog of me"
defaultContentLanguage = "en"

[languages]
	[en]
	lang = "en"
	languageName = "English"
	weight = 1
-- layouts/index.html --
Foo: {{ site.Params.foo }}|
	
		
	`
		b, err := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).BuildE()

		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, "no languages")
	})

	// Issue 11044
	t.Run("invalid defaultContentLanguage", func(t *testing.T) {

		files := `
-- hugo.toml --
baseURL = "https://example.org"
defaultContentLanguage = "sv"

[languages]
[languages.en]
languageCode = "en"
languageName = "English"
weight = 1

	
		
	`
		b, err := NewIntegrationTestBuilder(
			IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		).BuildE()

		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, "defaultContentLanguage does not match any language definition")
	})

}

// Issue #11040
func TestConfigModuleDefaultMountsInConfig(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
contentDir = "mycontent"
-- layouts/index.html --
Home.

	
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.Assert(b.H.Configs.Base.Module.Mounts, qt.HasLen, 7)
	b.Assert(b.H.Configs.LanguageConfigSlice[0].Module.Mounts, qt.HasLen, 7)

}
