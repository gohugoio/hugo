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

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/modules/npm"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/common/hugo"

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
		b.Assert(os.MkdirAll(workingDir, 0777), qt.IsNil)
		cfg := config.NewWithTestDefaults()
		cfg.Set("workingDir", workingDir)
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
		b.Assert(npm.Pack(b.H.BaseFs.SourceFs, b.H.BaseFs.Assets.Dirs), qt.IsNil)

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
		b.Assert(npm.Pack(b.H.BaseFs.SourceFs, b.H.BaseFs.Assets.Dirs), qt.IsNil)

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
		b.Assert(npm.Pack(b.H.BaseFs.SourceFs, b.H.BaseFs.Assets.Dirs), qt.IsNil)

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

		v := config.NewWithTestDefaults()
		v.Set("workingDir", workingDir)

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

	logger := loggers.NewWarningLogger()
	b.WithLogger(logger)

	b.Build(BuildCfg{})

	c := qt.New(t)

	c.Assert(logger.LogCounters().WarnCounter.Count(), qt.Equals, uint64(3))
}

func TestModulesSymlinks(t *testing.T) {
	skipSymlink(t)

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)
	workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-mod-sym")
	c.Assert(err, qt.IsNil)

	// We need to use the OS fs for this.
	cfg := config.NewWithTestDefaults()
	cfg.Set("workingDir", workingDir)
	fs := hugofs.NewFrom(hugofs.Os, cfg)

	defer clean()

	const homeTemplate = `
Data: {{ .Site.Data }}
`

	createDirsAndFiles := func(baseDir string) {
		for _, dir := range files.ComponentFolders {
			realDir := filepath.Join(baseDir, dir, "real")
			c.Assert(os.MkdirAll(realDir, 0777), qt.IsNil)
			c.Assert(afero.WriteFile(fs.Source, filepath.Join(realDir, "data.toml"), []byte("[hello]\nother = \"hello\""), 0777), qt.IsNil)
		}

		c.Assert(afero.WriteFile(fs.Source, filepath.Join(baseDir, "layouts", "index.html"), []byte(homeTemplate), 0777), qt.IsNil)
	}

	// Create project dirs and files.
	createDirsAndFiles(workingDir)
	// Create one module inside the default themes folder.
	themeDir := filepath.Join(workingDir, "themes", "mymod")
	createDirsAndFiles(themeDir)

	createSymlinks := func(baseDir, id string) {
		for _, dir := range files.ComponentFolders {
			// Issue #9119: private use language tags cannot exceed 8 characters.
			if dir != "i18n" {
				c.Assert(os.Chdir(filepath.Join(baseDir, dir)), qt.IsNil)
				c.Assert(os.Symlink("real", fmt.Sprintf("realsym%s", id)), qt.IsNil)
				c.Assert(os.Chdir(filepath.Join(baseDir, dir, "real")), qt.IsNil)
				c.Assert(os.Symlink("data.toml", fmt.Sprintf(filepath.FromSlash("datasym%s.toml"), id)), qt.IsNil)
			}
		}
	}

	createSymlinks(workingDir, "project")
	createSymlinks(themeDir, "mod")

	config := `
baseURL = "https://example.com"
theme="mymod"
defaultContentLanguage="nn"
defaultContentLanguageInSubDir=true

[languages]
[languages.nn]
weight = 1
[languages.en]
weight = 2


`

	b := newTestSitesBuilder(t).WithNothingAdded().WithWorkingDir(workingDir)
	b.WithLogger(loggers.NewErrorLogger())
	b.Fs = fs

	b.WithConfigFile("toml", config)
	c.Assert(os.Chdir(workingDir), qt.IsNil)

	b.Build(BuildCfg{})

	b.AssertFileContentFn(filepath.Join("public", "en", "index.html"), func(s string) bool {
		// Symbolic links only followed in project. There should be WARNING logs.
		return !strings.Contains(s, "symmod") && strings.Contains(s, "symproject")
	})

	bfs := b.H.BaseFs

	for i, componentFs := range []afero.Fs{
		bfs.Static[""].Fs,
		bfs.Archetypes.Fs,
		bfs.Content.Fs,
		bfs.Data.Fs,
		bfs.Assets.Fs,
		bfs.I18n.Fs,
	} {

		if i != 0 {
			continue
		}

		for j, id := range []string{"mod", "project"} {

			statCheck := func(fs afero.Fs, filename string, isDir bool) {
				shouldFail := j == 0
				if !shouldFail && i == 0 {
					// Static dirs only supports symlinks for files
					shouldFail = isDir
				}

				_, err := fs.Stat(filepath.FromSlash(filename))
				if err != nil {
					if i > 0 && strings.HasSuffix(filename, "toml") && strings.Contains(err.Error(), "files not supported") {
						// OK
						return
					}
				}

				if shouldFail {
					c.Assert(err, qt.Not(qt.IsNil))
					c.Assert(err, qt.Equals, hugofs.ErrPermissionSymlink)
				} else {
					c.Assert(err, qt.IsNil)
				}
			}

			c.Logf("Checking %d:%d %q", i, j, id)

			statCheck(componentFs, fmt.Sprintf("realsym%s", id), true)
			statCheck(componentFs, fmt.Sprintf("real/datasym%s.toml", id), false)

		}
	}
}

func TestMountsProject(t *testing.T) {
	t.Parallel()

	config := `

baseURL="https://example.org"

[module]
[[module.mounts]]
source="mycontent"
target="content"

`
	b := newTestSitesBuilder(t).
		WithConfigFile("toml", config).
		WithSourceFile(filepath.Join("mycontent", "mypage.md"), `
---
title: "My Page"
---

`)

	b.Build(BuildCfg{})

	// helpers.PrintFs(b.H.Fs.Source, "public", os.Stdout)

	b.AssertFileContent("public/mypage/index.html", "Permalink: https://example.org/mypage/")
}

// https://github.com/gohugoio/hugo/issues/6684
func TestMountsContentFile(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-modules-content-file")
	c.Assert(err, qt.IsNil)
	defer clean()

	configTemplate := `
baseURL = "https://example.com"
title = "My Modular Site"
workingDir = %q

[module]
  [[module.mounts]]
    source = "README.md"
    target = "content/_index.md"
  [[module.mounts]]
    source = "mycontent"
    target = "content/blog"

`

	tomlConfig := fmt.Sprintf(configTemplate, workingDir)

	b := newTestSitesBuilder(t).Running()

	cfg := config.NewWithTestDefaults()
	cfg.Set("workingDir", workingDir)

	b.Fs = hugofs.NewDefault(cfg)

	b.WithWorkingDir(workingDir).WithConfigFile("toml", tomlConfig)
	b.WithTemplatesAdded("index.html", `
{{ .Title }}
{{ .Content }}

{{ $readme := .Site.GetPage "/README.md" }}
{{ with $readme }}README: {{ .Title }}|Filename: {{ path.Join .File.Filename }}|Path: {{ path.Join .File.Path }}|FilePath: {{ path.Join .File.FileInfo.Meta.PathFile }}|{{ end }}


{{ $mypage := .Site.GetPage "/blog/mypage.md" }}
{{ with $mypage }}MYPAGE: {{ .Title }}|Path: {{ path.Join .File.Path }}|FilePath: {{ path.Join .File.FileInfo.Meta.PathFile }}|{{ end }}
{{ $mybundle := .Site.GetPage "/blog/mybundle" }}
{{ with $mybundle }}MYBUNDLE: {{ .Title }}|Path: {{ path.Join .File.Path }}|FilePath: {{ path.Join .File.FileInfo.Meta.PathFile }}|{{ end }}


`, "_default/_markup/render-link.html", `
{{ $link := .Destination }}
{{ $isRemote := strings.HasPrefix $link "http" }}
{{- if not $isRemote -}}
{{ $url := urls.Parse .Destination }}
{{ $fragment := "" }}
{{- with $url.Fragment }}{{ $fragment = printf "#%s" . }}{{ end -}}
{{- with .Page.GetPage $url.Path }}{{ $link = printf "%s%s" .Permalink $fragment }}{{ end }}{{ end -}}
<a href="{{ $link | safeURL }}"{{ with .Title}} title="{{ . }}"{{ end }}{{ if $isRemote }} target="_blank"{{ end }}>{{ .Text | safeHTML }}</a>
`)

	os.Mkdir(filepath.Join(workingDir, "mycontent"), 0777)
	os.Mkdir(filepath.Join(workingDir, "mycontent", "mybundle"), 0777)

	b.WithSourceFile("README.md", `---
title: "Readme Title"
---

Readme Content.
`,
		filepath.Join("mycontent", "mypage.md"), `
---
title: "My Page"
---


* [Relative Link From Page](mybundle)
* [Relative Link From Page, filename](mybundle/index.md)
* [Link using original path](/mycontent/mybundle/index.md)


`, filepath.Join("mycontent", "mybundle", "index.md"), `
---
title: "My Bundle"
---

* [Dot Relative Link From Bundle](../mypage.md)
* [Link using original path](/mycontent/mypage.md)
* [Link to Home](/)
* [Link to Home, README.md](/README.md)
* [Link to Home, _index.md](/_index.md)

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
README: Readme Title
/README.md|Path: _index.md|FilePath: README.md
Readme Content.
MYPAGE: My Page|Path: blog/mypage.md|FilePath: mycontent/mypage.md|
MYBUNDLE: My Bundle|Path: blog/mybundle/index.md|FilePath: mycontent/mybundle/index.md|
`)
	b.AssertFileContent("public/blog/mypage/index.html", `
<a href="https://example.com/blog/mybundle/">Relative Link From Page</a>
<a href="https://example.com/blog/mybundle/">Relative Link From Page, filename</a>
<a href="https://example.com/blog/mybundle/">Link using original path</a>

`)
	b.AssertFileContent("public/blog/mybundle/index.html", `
<a href="https://example.com/blog/mypage/">Dot Relative Link From Bundle</a>
<a href="https://example.com/blog/mypage/">Link using original path</a>
<a href="https://example.com/">Link to Home</a>
<a href="https://example.com/">Link to Home, README.md</a>
<a href="https://example.com/">Link to Home, _index.md</a>
`)

	b.EditFiles("README.md", `---
title: "Readme Edit"
---
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
Readme Edit
`)
}

func TestMountsPaths(t *testing.T) {
	c := qt.New(t)

	type test struct {
		b          *sitesBuilder
		clean      func()
		workingDir string
	}

	prepare := func(c *qt.C, mounts string) test {
		workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-mounts-paths")
		c.Assert(err, qt.IsNil)

		configTemplate := `
baseURL = "https://example.com"
title = "My Modular Site"
workingDir = %q

%s

`
		tomlConfig := fmt.Sprintf(configTemplate, workingDir, mounts)
		tomlConfig = strings.ReplaceAll(tomlConfig, "WORKING_DIR", workingDir)

		b := newTestSitesBuilder(c).Running()

		cfg := config.NewWithTestDefaults()
		cfg.Set("workingDir", workingDir)
		b.Fs = hugofs.NewDefault(cfg)

		os.MkdirAll(filepath.Join(workingDir, "content", "blog"), 0777)

		b.WithWorkingDir(workingDir).WithConfigFile("toml", tomlConfig)

		return test{
			b:          b,
			clean:      clean,
			workingDir: workingDir,
		}
	}

	c.Run("Default", func(c *qt.C) {
		mounts := ``

		test := prepare(c, mounts)
		b := test.b
		defer test.clean()

		b.WithContent("blog/p1.md", `---
title: P1
---`)

		b.Build(BuildCfg{})

		p := b.GetPage("blog/p1.md")
		f := p.File().FileInfo().Meta()
		b.Assert(filepath.ToSlash(f.Path), qt.Equals, "blog/p1.md")
		b.Assert(filepath.ToSlash(f.PathFile()), qt.Equals, "content/blog/p1.md")

		b.Assert(b.H.BaseFs.Layouts.Path(filepath.Join(test.workingDir, "layouts", "_default", "single.html")), qt.Equals, filepath.FromSlash("_default/single.html"))
	})

	c.Run("Mounts", func(c *qt.C) {
		absDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-mounts-paths-abs")
		c.Assert(err, qt.IsNil)
		defer clean()

		mounts := `[module]
  [[module.mounts]]
    source = "README.md"
    target = "content/_index.md"
  [[module.mounts]]
    source = "mycontent"
    target = "content/blog"
   [[module.mounts]]
    source = "subdir/mypartials"
    target = "layouts/partials"
   [[module.mounts]]
    source = %q
    target = "layouts/shortcodes"
`
		mounts = fmt.Sprintf(mounts, filepath.Join(absDir, "/abs/myshortcodes"))

		test := prepare(c, mounts)
		b := test.b
		defer test.clean()

		subContentDir := filepath.Join(test.workingDir, "mycontent", "sub")
		os.MkdirAll(subContentDir, 0777)
		myPartialsDir := filepath.Join(test.workingDir, "subdir", "mypartials")
		os.MkdirAll(myPartialsDir, 0777)

		absShortcodesDir := filepath.Join(absDir, "abs", "myshortcodes")
		os.MkdirAll(absShortcodesDir, 0777)

		b.WithSourceFile("README.md", "---\ntitle: Readme\n---")
		b.WithSourceFile("mycontent/sub/p1.md", "---\ntitle: P1\n---")

		b.WithSourceFile(filepath.Join(absShortcodesDir, "myshort.html"), "MYSHORT")
		b.WithSourceFile(filepath.Join(myPartialsDir, "mypartial.html"), "MYPARTIAL")

		b.Build(BuildCfg{})

		p1_1 := b.GetPage("/blog/sub/p1.md")
		p1_2 := b.GetPage("/mycontent/sub/p1.md")
		b.Assert(p1_1, qt.Not(qt.IsNil))
		b.Assert(p1_2, qt.Equals, p1_1)

		f := p1_1.File().FileInfo().Meta()
		b.Assert(filepath.ToSlash(f.Path), qt.Equals, "blog/sub/p1.md")
		b.Assert(filepath.ToSlash(f.PathFile()), qt.Equals, "mycontent/sub/p1.md")
		b.Assert(b.H.BaseFs.Layouts.Path(filepath.Join(myPartialsDir, "mypartial.html")), qt.Equals, filepath.FromSlash("partials/mypartial.html"))
		b.Assert(b.H.BaseFs.Layouts.Path(filepath.Join(absShortcodesDir, "myshort.html")), qt.Equals, filepath.FromSlash("shortcodes/myshort.html"))
		b.Assert(b.H.BaseFs.Content.Path(filepath.Join(subContentDir, "p1.md")), qt.Equals, filepath.FromSlash("blog/sub/p1.md"))
		b.Assert(b.H.BaseFs.Content.Path(filepath.Join(test.workingDir, "README.md")), qt.Equals, filepath.FromSlash("_index.md"))
	})
}

// https://github.com/gohugoio/hugo/issues/6299
func TestSiteWithGoModButNoModules(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	// We need to use the OS fs for this.
	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-no-mod")
	c.Assert(err, qt.IsNil)

	cfg := config.NewWithTestDefaults()
	cfg.Set("workingDir", workDir)
	fs := hugofs.NewFrom(hugofs.Os, cfg)

	defer clean()

	b := newTestSitesBuilder(t)
	b.Fs = fs

	b.WithWorkingDir(workDir).WithViper(cfg)

	b.WithSourceFile("go.mod", "")
	b.Build(BuildCfg{})
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

	cfg := config.NewWithTestDefaults()
	cfg.Set("workingDir", workDir)
	fs := hugofs.NewFrom(hugofs.Os, cfg)

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
`), 0777)

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
	config := `baseURL = 'https://example.org/'
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
`
	b := newTestSitesBuilder(t).WithConfigFile("toml", config)

	b.WithContent("p1.md", "")

	b.WithSourceFile(
		"extra-content/_index.md", "",
		"extra-content/subdir/_index.md", "",
		"extra-content/subdir/about.md", "",
	)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/resources-a/subdir/about/index.html", "Single")
	b.AssertFileContent("public/resources-b/subdir/about/index.html", "Single")
}
