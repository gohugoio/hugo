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

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/testmodBuilder/mods"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// TODO(bep) this fails when testmodBuilder is also building ...
func TestHugoModules(t *testing.T) {
	t.Parallel()

	if hugo.GoMinorVersion() < 12 {
		// https://github.com/golang/go/issues/26794
		// There were some concurrent issues with Go modules in < Go 12.
		t.Skip("skip this for Go <= 1.11 due to a bug in Go's stdlib")
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
		assert := require.New(t)

		v := viper.New()

		workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-modules-test")
		assert.NoError(err)
		defer clean()

		configTemplate := `
baseURL = "https://example.com"
title = "My Modular Site"
workingDir = %q
theme = %q
ignoreVendor = %t

`

		config := fmt.Sprintf(configTemplate, workingDir, m.Path(), ignoreVendor)

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
	// Child depdendencies are one behind.
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

	logger := loggers.NewWarningLogger()
	b.WithLogger(logger)

	b.Build(BuildCfg{})

	assert := require.New(t)

	assert.Equal(uint64(2), logger.WarnCounter.Count())

}

func TestModulesSymlinks(t *testing.T) {
	skipSymlink(t)

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	assert := require.New(t)
	// We need to use the OS fs for this.
	cfg := viper.New()
	fs := hugofs.NewFrom(hugofs.Os, cfg)

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-mod-sym")
	assert.NoError(err)

	defer clean()

	const homeTemplate = `
Data: {{ .Site.Data }}
`

	createDirsAndFiles := func(baseDir string) {
		for _, dir := range files.ComponentFolders {
			realDir := filepath.Join(baseDir, dir, "real")
			assert.NoError(os.MkdirAll(realDir, 0777))
			assert.NoError(afero.WriteFile(fs.Source, filepath.Join(realDir, "data.toml"), []byte("[hello]\nother = \"hello\""), 0777))
		}

		assert.NoError(afero.WriteFile(fs.Source, filepath.Join(baseDir, "layouts", "index.html"), []byte(homeTemplate), 0777))
	}

	// Create project dirs and files.
	createDirsAndFiles(workDir)
	// Create one module inside the default themes folder.
	themeDir := filepath.Join(workDir, "themes", "mymod")
	createDirsAndFiles(themeDir)

	createSymlinks := func(baseDir, id string) {
		for _, dir := range files.ComponentFolders {
			assert.NoError(os.Chdir(filepath.Join(baseDir, dir)))
			assert.NoError(os.Symlink("real", fmt.Sprintf("realsym%s", id)))
			assert.NoError(os.Chdir(filepath.Join(baseDir, dir, "real")))
			assert.NoError(os.Symlink("data.toml", fmt.Sprintf(filepath.FromSlash("datasym%s.toml"), id)))
		}
	}

	createSymlinks(workDir, "project")
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

	b := newTestSitesBuilder(t).WithNothingAdded().WithWorkingDir(workDir)
	b.WithLogger(loggers.NewErrorLogger())
	b.Fs = fs

	b.WithConfigFile("toml", config)
	assert.NoError(os.Chdir(workDir))

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
		bfs.I18n.Fs} {

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
					assert.Error(err)
					assert.Equal(hugofs.ErrPermissionSymlink, err, filename)
				} else {
					assert.NoError(err, filename)
				}
			}

			statCheck(componentFs, fmt.Sprintf("realsym%s", id), true)
			statCheck(componentFs, fmt.Sprintf("real/datasym%s.toml", id), false)

		}
	}
}

func TestMountsProject(t *testing.T) {

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

	//helpers.PrintFs(b.H.Fs.Source, "public", os.Stdout)

	b.AssertFileContent("public/mypage/index.html", "Permalink: https://example.org/mypage/")
}
