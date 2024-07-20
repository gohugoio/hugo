// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/modules/npm"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/testmodBuilder/mods"
)

func TestHugoModulesVariants(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}

	tomlConfig := `
baseURL="https://example.org"
workingDir = %q

[module]
[[module.imports]]
path="github.com/gohugoio/hugoTestModule2"
%s
`

	createConfig := func(workingDir, moduleOpts string) string {
		return fmt.Sprintf(tomlConfig, workingDir, moduleOpts)
	}

	newTestBuilder := func(t testing.TB, moduleOpts string) *sitesBuilder {
		b := newTestSitesBuilder(t)
		tempDir := t.TempDir()
		workingDir := filepath.Join(tempDir, "myhugosite")
		b.Assert(os.MkdirAll(workingDir, 0o777), qt.IsNil)
		cfg := config.New()
		cfg.Set("workingDir", workingDir)
		cfg.Set("publishDir", "public")
		b.Fs = hugofs.NewDefault(cfg)
		b.WithWorkingDir(workingDir).WithConfigFile("toml", createConfig(workingDir, moduleOpts))
		b.WithTemplates(
			"index.html", `
Param from module: {{ site.Params.Hugo }}|
{{ $js := resources.Get "jslibs/alpinejs/alpine.js" }}
JS imported in module: {{ with $js }}{{ .RelPermalink }}{{ end }}|
`,
			"_default/single.html", `{{ .Content }}`)
		b.WithContent("p1.md", `---
title: "Page"
---

[A link](https://bep.is)

`)
		b.WithSourceFile("go.mod", `
module github.com/gohugoio/tests/testHugoModules


`)

		b.WithSourceFile("go.sum", `
github.com/gohugoio/hugoTestModule2 v0.0.0-20200131160637-9657d7697877 h1:WLM2bQCKIWo04T6NsIWsX/Vtirhf0TnpY66xyqGlgVY=
github.com/gohugoio/hugoTestModule2 v0.0.0-20200131160637-9657d7697877/go.mod h1:CBFZS3khIAXKxReMwq0le8sEl/D8hcXmixlOHVv+Gd0=
`)

		return b
	}

	t.Run("Target in subfolder", func(t *testing.T) {
		b := newTestBuilder(t, "ignoreImports=true")
		b.Build(BuildCfg{})

		b.AssertFileContent("public/p1/index.html", `<p>Page|https://bep.is|Title: |Text: A link|END</p>`)
	})

	t.Run("Ignore config", func(t *testing.T) {
		b := newTestBuilder(t, "ignoreConfig=true")
		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
Param from module: |
JS imported in module: |
`)
	})

	t.Run("Ignore imports", func(t *testing.T) {
		b := newTestBuilder(t, "ignoreImports=true")
		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
Param from module: Rocks|
JS imported in module: |
`)
	})

	t.Run("Create package.json", func(t *testing.T) {
		b := newTestBuilder(t, "")

		b.WithSourceFile("package.json", `{
		"name": "mypack",
		"version": "1.2.3",
        "scripts": {
          "client": "wait-on http://localhost:1313 && open http://localhost:1313",
          "start": "run-p client server",
		  "test": "echo 'hoge' > hoge"
		},
          "dependencies": {
        	"nonon": "error"
        	}
}`)

		b.WithSourceFile("package.hugo.json", `{
		"name": "mypack",
		"version": "1.2.3",
        "scripts": {
          "client": "wait-on http://localhost:1313 && open http://localhost:1313",
          "start": "run-p client server",
		  "test": "echo 'hoge' > hoge"
		},
          "dependencies": {
        	"foo": "1.2.3"
        	},
        "devDependencies": {
                "postcss-cli": "7.8.0",
                "tailwindcss": "1.8.0"

        }
}`)

		b.Build(BuildCfg{})
		b.Assert(npm.Pack(b.H.BaseFs.ProjectSourceFs, b.H.BaseFs.AssetsWithDuplicatesPreserved.Fs), qt.IsNil)

		b.AssertFileContentFn("package.json", func(s string) bool {
			return s == `{
  "comments": {
    "dependencies": {
      "foo": "project",
      "react-dom": "github.com/gohugoio/hugoTestModule2"
    },
    "devDependencies": {
      "@babel/cli": "github.com/gohugoio/hugoTestModule2",
      "@babel/core": "github.com/gohugoio/hugoTestModule2",
      "@babel/preset-env": "github.com/gohugoio/hugoTestModule2",
      "postcss-cli": "project",
      "tailwindcss": "project"
    }
  },
  "dependencies": {
    "foo": "1.2.3",
    "react-dom": "^16.13.1"
  },
  "devDependencies": {
    "@babel/cli": "7.8.4",
    "@babel/core": "7.9.0",
    "@babel/preset-env": "7.9.5",
    "postcss-cli": "7.8.0",
    "tailwindcss": "1.8.0"
  },
  "name": "mypack",
  "scripts": {
    "client": "wait-on http://localhost:1313 && open http://localhost:1313",
    "start": "run-p client server",
    "test": "echo 'hoge' > hoge"
  },
  "version": "1.2.3"
}
`
		})
	})

	t.Run("Create package.json, no default", func(t *testing.T) {
		b := newTestBuilder(t, "")

		const origPackageJSON = `{
		"name": "mypack",
		"version": "1.2.3",
        "scripts": {
          "client": "wait-on http://localhost:1313 && open http://localhost:1313",
          "start": "run-p client server",
		  "test": "echo 'hoge' > hoge"
		},
          "dependencies": {
           "moo": "1.2.3"
        	}
}`

		b.WithSourceFile("package.json", origPackageJSON)

		b.Build(BuildCfg{})
		b.Assert(npm.Pack(b.H.BaseFs.ProjectSourceFs, b.H.BaseFs.AssetsWithDuplicatesPreserved.Fs), qt.IsNil)

		b.AssertFileContentFn("package.json", func(s string) bool {
			return s == `{
  "comments": {
    "dependencies": {
      "moo": "project",
      "react-dom": "github.com/gohugoio/hugoTestModule2"
    },
    "devDependencies": {
      "@babel/cli": "github.com/gohugoio/hugoTestModule2",
      "@babel/core": "github.com/gohugoio/hugoTestModule2",
      "@babel/preset-env": "github.com/gohugoio/hugoTestModule2",
      "postcss-cli": "github.com/gohugoio/hugoTestModule2",
      "tailwindcss": "github.com/gohugoio/hugoTestModule2"
    }
  },
  "dependencies": {
    "moo": "1.2.3",
    "react-dom": "^16.13.1"
  },
  "devDependencies": {
    "@babel/cli": "7.8.4",
    "@babel/core": "7.9.0",
    "@babel/preset-env": "7.9.5",
    "postcss-cli": "7.1.0",
    "tailwindcss": "1.2.0"
  },
  "name": "mypack",
  "scripts": {
    "client": "wait-on http://localhost:1313 && open http://localhost:1313",
    "start": "run-p client server",
    "test": "echo 'hoge' > hoge"
  },
  "version": "1.2.3"
}
`
		})

		// https://github.com/gohugoio/hugo/issues/7690
		b.AssertFileContent("package.hugo.json", origPackageJSON)
	})

	t.Run("Create package.json, no default, no package.json", func(t *testing.T) {
		b := newTestBuilder(t, "")

		b.Build(BuildCfg{})
		b.Assert(npm.Pack(b.H.BaseFs.ProjectSourceFs, b.H.BaseFs.AssetsWithDuplicatesPreserved.Fs), qt.IsNil)

		b.AssertFileContentFn("package.json", func(s string) bool {
			return s == `{
  "comments": {
    "dependencies": {
      "react-dom": "github.com/gohugoio/hugoTestModule2"
    },
    "devDependencies": {
      "@babel/cli": "github.com/gohugoio/hugoTestModule2",
      "@babel/core": "github.com/gohugoio/hugoTestModule2",
      "@babel/preset-env": "github.com/gohugoio/hugoTestModule2",
      "postcss-cli": "github.com/gohugoio/hugoTestModule2",
      "tailwindcss": "github.com/gohugoio/hugoTestModule2"
    }
  },
  "dependencies": {
    "react-dom": "^16.13.1"
  },
  "devDependencies": {
    "@babel/cli": "7.8.4",
    "@babel/core": "7.9.0",
    "@babel/preset-env": "7.9.5",
    "postcss-cli": "7.1.0",
    "tailwindcss": "1.2.0"
  },
  "name": "myhugosite",
  "version": "0.1.0"
}
`
		})
	})
}

// TODO(bep) this fails when testmodBuilder is also building ...
func TestHugoModulesMatrix(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip (relative) long running modules test when running locally")
	}
	t.Parallel()

	if !htesting.IsCI() || hugo.GoMinorVersion() < 12 {
		// https://github.com/golang/go/issues/26794
		// There were some concurrent issues with Go modules in < Go 12.
		t.Skip("skip this on local host and for Go <= 1.11 due to a bug in Go's stdlib")
	}

	if testing.Short() {
		t.Skip()
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	gooss := []string{"linux", "darwin", "windows"}
	goos := gooss[rnd.Intn(len(gooss))]
	ignoreVendor := rnd.Intn(2) == 0
	testmods := mods.CreateModules(goos).Collect()
	rnd.Shuffle(len(testmods), func(i, j int) { testmods[i], testmods[j] = testmods[j], testmods[i] })

	for _, m := range testmods[:2] {
		c := qt.New(t)

		workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-modules-test")
		c.Assert(err, qt.IsNil)
		defer clean()

		v := config.New()
		v.Set("workingDir", workingDir)
		v.Set("publishDir", "public")

		configTemplate := `
baseURL = "https://example.com"
title = "My Modular Site"
workingDir = %q
theme = %q
ignoreVendorPaths = %q

`

		ignoreVendorPaths := ""
		if ignoreVendor {
			ignoreVendorPaths = "github.com/**"
		}
		config := fmt.Sprintf(configTemplate, workingDir, m.Path(), ignoreVendorPaths)

		b := newTestSitesBuilder(t)

		// Need to use OS fs for this.
		b.Fs = hugofs.NewDefault(v)

		b.WithWorkingDir(workingDir).WithConfigFile("toml", config)
		b.WithContent("page.md", `
---
title: "Foo"
---
`)
		b.WithTemplates("home.html", `

{{ $mod := .Site.Data.modinfo.module }}
Mod Name: {{ $mod.name }}
Mod Version: {{ $mod.version }}
----
{{ range $k, $v := .Site.Data.modinfo }}
- {{ $k }}: {{ range $kk, $vv := $v }}{{ $kk }}: {{ $vv }}|{{ end -}}
{{ end }}


`)
		b.WithSourceFile("go.mod", `
module github.com/gohugoio/tests/testHugoModules


`)

		b.Build(BuildCfg{})

		// Verify that go.mod is autopopulated with all the modules in config.toml.
		b.AssertFileContent("go.mod", m.Path())

		b.AssertFileContent("public/index.html",
			"Mod Name: "+m.Name(),
			"Mod Version: v1.4.0")

		b.AssertFileContent("public/index.html", createChildModMatchers(m, ignoreVendor, m.Vendor)...)

	}
}

func createChildModMatchers(m *mods.Md, ignoreVendor, vendored bool) []string {
	// Child dependencies are one behind.
	expectMinorVersion := 3

	if !ignoreVendor && vendored {
		// Vendored modules are stuck at v1.1.0.
		expectMinorVersion = 1
	}

	expectVersion := fmt.Sprintf("v1.%d.0", expectMinorVersion)

	var matchers []string

	for _, mm := range m.Children {
		matchers = append(
			matchers,
			fmt.Sprintf("%s: name: %s|version: %s", mm.Name(), mm.Name(), expectVersion))
		matchers = append(matchers, createChildModMatchers(mm, ignoreVendor, vendored || mm.Vendor)...)
	}
	return matchers
}

func TestModulesWithContent(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).WithWorkingDir("/site").WithConfigFile("toml", `
baseURL="https://example.org"

workingDir="/site"

defaultContentLanguage = "en"

[module]
[[module.imports]]
path="a"
[[module.imports.mounts]]
source="myacontent"
target="content/blog"
lang="en"
[[module.imports]]
path="b"
[[module.imports.mounts]]
source="mybcontent"
target="content/blog"
lang="nn"
[[module.imports]]
path="c"
[[module.imports]]
path="d"

[languages]

[languages.en]
title = "Title in English"
languageName = "English"
weight = 1
[languages.nn]
languageName = "Nynorsk"
weight = 2
title = "Tittel på nynorsk"
[languages.nb]
languageName = "Bokmål"
weight = 3
title = "Tittel på bokmål"
[languages.fr]
languageName = "French"
weight = 4
title = "French Title"


`)

	b.WithTemplatesAdded("index.html", `
{{ range .Site.RegularPages }}
|{{ .Title }}|{{ .RelPermalink }}|{{ .Plain }}
{{ end }}
{{ $data := .Site.Data }}
Data Common: {{ $data.common.value }}
Data C: {{ $data.c.value }}
Data D: {{ $data.d.value }}
All Data: {{ $data }}

i18n hello1: {{ i18n "hello1" . }}
i18n theme: {{ i18n "theme" . }}
i18n theme2: {{ i18n "theme2" . }}
`)

	content := func(id string) string {
		return fmt.Sprintf(`---
title: Title %s
---
Content %s

`, id, id)
	}

	i18nContent := func(id, value string) string {
		return fmt.Sprintf(`
[%s]
other = %q
`, id, value)
	}

	// Content files
	b.WithSourceFile("themes/a/myacontent/page.md", content("theme-a-en"))
	b.WithSourceFile("themes/b/mybcontent/page.md", content("theme-b-nn"))
	b.WithSourceFile("themes/c/content/blog/c.md", content("theme-c-nn"))

	// Data files
	b.WithSourceFile("data/common.toml", `value="Project"`)
	b.WithSourceFile("themes/c/data/common.toml", `value="Theme C"`)
	b.WithSourceFile("themes/c/data/c.toml", `value="Hugo Rocks!"`)
	b.WithSourceFile("themes/d/data/c.toml", `value="Hugo Rodcks!"`)
	b.WithSourceFile("themes/d/data/d.toml", `value="Hugo Rodks!"`)

	// i18n files
	b.WithSourceFile("i18n/en.toml", i18nContent("hello1", "Project"))
	b.WithSourceFile("themes/c/i18n/en.toml", `
[hello1]
other="Theme C Hello"
[theme]
other="Theme C"
`)
	b.WithSourceFile("themes/d/i18n/en.toml", i18nContent("theme", "Theme D"))
	b.WithSourceFile("themes/d/i18n/en.toml", i18nContent("theme2", "Theme2 D"))

	// Static files
	b.WithSourceFile("themes/c/static/hello.txt", `Hugo Rocks!"`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "|Title theme-a-en|/blog/page/|Content theme-a-en")
	b.AssertFileContent("public/nn/index.html", "|Title theme-b-nn|/nn/blog/page/|Content theme-b-nn")

	// Data
	b.AssertFileContent("public/index.html",
		"Data Common: Project",
		"Data C: Hugo Rocks!",
		"Data D: Hugo Rodks!",
	)

	// i18n
	b.AssertFileContent("public/index.html",
		"i18n hello1: Project",
		"i18n theme: Theme C",
		"i18n theme2: Theme2 D",
	)
}

func TestModulesIgnoreConfig(t *testing.T) {
	b := newTestSitesBuilder(t).WithWorkingDir("/site").WithConfigFile("toml", `
baseURL="https://example.org"

workingDir="/site"

[module]
[[module.imports]]
path="a"
ignoreConfig=true

`)

	b.WithSourceFile("themes/a/config.toml", `
[params]
a = "Should Be Ignored!"
`)

	b.WithTemplatesAdded("index.html", `Params: {{ .Site.Params }}`)

	b.Build(BuildCfg{})

	b.AssertFileContentFn("public/index.html", func(s string) bool {
		return !strings.Contains(s, "Ignored")
	})
}

func TestModulesDisabled(t *testing.T) {
	b := newTestSitesBuilder(t).WithWorkingDir("/site").WithConfigFile("toml", `
baseURL="https://example.org"

workingDir="/site"

[module]
[[module.imports]]
path="a"
[[module.imports]]
path="b"
disable=true


`)

	b.WithSourceFile("themes/a/config.toml", `
[params]
a = "A param"
`)

	b.WithSourceFile("themes/b/config.toml", `
[params]
b = "B param"
`)

	b.WithTemplatesAdded("index.html", `Params: {{ .Site.Params }}`)

	b.Build(BuildCfg{})

	b.AssertFileContentFn("public/index.html", func(s string) bool {
		return strings.Contains(s, "A param") && !strings.Contains(s, "B param")
	})
}

func TestModulesIncompatible(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).WithWorkingDir("/site").WithConfigFile("toml", `
baseURL="https://example.org"

workingDir="/site"

[module]
[[module.imports]]
path="ok"
[[module.imports]]
path="incompat1"
[[module.imports]]
path="incompat2"
[[module.imports]]
path="incompat3"

`)

	b.WithSourceFile("themes/ok/data/ok.toml", `title = "OK"`)

	b.WithSourceFile("themes/incompat1/config.toml", `

[module]
[module.hugoVersion]
min = "0.33.2"
max = "0.45.0"

`)

	// Old setup.
	b.WithSourceFile("themes/incompat2/theme.toml", `
min_version = "5.0.0"

`)

	// Issue 6162
	b.WithSourceFile("themes/incompat3/theme.toml", `
min_version = 0.55.0

`)

	logger := loggers.NewDefault()
	b.WithLogger(logger)

	b.Build(BuildCfg{})

	c := qt.New(t)

	c.Assert(logger.LoggCount(logg.LevelWarn), qt.Equals, 3)
}

func TestMountsProject(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL="https://example.org"

[module]
[[module.mounts]]
source="mycontent"
target="content"
-- layouts/_default/single.html --
Permalink: {{ .Permalink }}|
-- mycontent/mypage.md --
---
title: "My Page"
---
`
	b := Test(t, files)

	b.AssertFileContent("public/mypage/index.html", "Permalink: https://example.org/mypage/|")
}

// https://github.com/gohugoio/hugo/issues/6684
func TestMountsContentFile(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]
disableLiveReload = true
[module]
[[module.mounts]]
source = "README.md"
target = "content/_index.md"	
-- README.md --
# Hello World
-- layouts/index.html --
Home: {{ .Title }}|{{ .Content }}|
`
	b := Test(t, files)
	b.AssertFileContent("public/index.html", "Home: |<h1 id=\"hello-world\">Hello World</h1>\n|")
}

// https://github.com/gohugoio/hugo/issues/6299
func TestSiteWithGoModButNoModules(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
-- go.mod --

`

	b := Test(t, files, TestOptWithConfig(func(cfg *IntegrationTestConfig) {
		cfg.WorkingDir = tempDir
	}))

	b.Build()
}

// https://github.com/gohugoio/hugo/issues/6622
func TestModuleAbsMount(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	// We need to use the OS fs for this.
	workDir, clean1, err := htesting.CreateTempDir(hugofs.Os, "hugo-project")
	c.Assert(err, qt.IsNil)
	absContentDir, clean2, err := htesting.CreateTempDir(hugofs.Os, "hugo-content")
	c.Assert(err, qt.IsNil)

	cfg := config.New()
	cfg.Set("workingDir", workDir)
	cfg.Set("publishDir", "public")
	fs := hugofs.NewFromOld(hugofs.Os, cfg)

	config := fmt.Sprintf(`
workingDir=%q

[module]
  [[module.mounts]]
    source = %q
    target = "content"

`, workDir, absContentDir)

	defer clean1()
	defer clean2()

	b := newTestSitesBuilder(t)
	b.Fs = fs

	contentFilename := filepath.Join(absContentDir, "p1.md")
	afero.WriteFile(hugofs.Os, contentFilename, []byte(`
---
title: Abs
---

Content.
`), 0o777)

	b.WithWorkingDir(workDir).WithConfigFile("toml", config)
	b.WithContent("dummy.md", "")

	b.WithTemplatesAdded("index.html", `
{{ $p1 := site.GetPage "p1" }}
P1: {{ $p1.Title }}|{{ $p1.RelPermalink }}|Filename: {{ $p1.File.Filename }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "P1: Abs|/p1/", "Filename: "+contentFilename)
}

// Issue 9426
func TestMountSameSource(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'Hugo GitHub Issue #9426'

disableKinds = ['RSS','sitemap','taxonomy','term']

[[module.mounts]]
source = "content"
target = "content"

[[module.mounts]]
source = "extra-content"
target = "content/resources-a"

[[module.mounts]]
source = "extra-content"
target = "content/resources-b"
-- layouts/_default/single.html --
Single
-- content/p1.md --
-- extra-content/_index.md --
-- extra-content/subdir/_index.md --
-- extra-content/subdir/about.md --
"
`
	b := Test(t, files)

	b.AssertFileContent("public/resources-a/subdir/about/index.html", "Single")
	b.AssertFileContent("public/resources-b/subdir/about/index.html", "Single")
}

func TestMountData(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = 'https://example.org/'
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]

[[module.mounts]]
source = "data"
target = "data"

[[module.mounts]]
source = "extra-data"
target = "data/extra"
-- extra-data/test.yaml --
message: Hugo Rocks
-- layouts/index.html --
{{ site.Data.extra.test.message }}
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Hugo Rocks")
}
