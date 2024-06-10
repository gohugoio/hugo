// Copyright 2024 The Hugo Authors. All rights reserved.
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

package filesystems_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/hugolib"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/hugolib/paths"
)

func TestNewBaseFs(t *testing.T) {
	c := qt.New(t)
	v := config.New()

	themes := []string{"btheme", "atheme"}

	workingDir := filepath.FromSlash("/my/work")
	v.Set("workingDir", workingDir)
	v.Set("contentDir", "content")
	v.Set("themesDir", "themes")
	v.Set("defaultContentLanguage", "en")
	v.Set("theme", themes[:1])
	v.Set("publishDir", "public")

	afs := afero.NewMemMapFs()

	// Write some data to the themes
	for _, theme := range themes {
		for _, dir := range []string{"i18n", "data", "archetypes", "layouts"} {
			base := filepath.Join(workingDir, "themes", theme, dir)
			filenameTheme := filepath.Join(base, fmt.Sprintf("theme-file-%s.txt", theme))
			filenameOverlap := filepath.Join(base, "f3.txt")
			afs.Mkdir(base, 0o755)
			content := []byte(fmt.Sprintf("content:%s:%s", theme, dir))
			afero.WriteFile(afs, filenameTheme, content, 0o755)
			afero.WriteFile(afs, filenameOverlap, content, 0o755)
		}
		// Write some files to the root of the theme
		base := filepath.Join(workingDir, "themes", theme)
		afero.WriteFile(afs, filepath.Join(base, fmt.Sprintf("theme-root-%s.txt", theme)), []byte(fmt.Sprintf("content:%s", theme)), 0o755)
		afero.WriteFile(afs, filepath.Join(base, "file-theme-root.txt"), []byte(fmt.Sprintf("content:%s", theme)), 0o755)
	}

	afero.WriteFile(afs, filepath.Join(workingDir, "file-root.txt"), []byte("content-project"), 0o755)

	afero.WriteFile(afs, filepath.Join(workingDir, "themes", "btheme", "config.toml"), []byte(`
theme = ["atheme"]
`), 0o755)

	setConfigAndWriteSomeFilesTo(afs, v, "contentDir", "mycontent", 3)
	setConfigAndWriteSomeFilesTo(afs, v, "i18nDir", "myi18n", 4)
	setConfigAndWriteSomeFilesTo(afs, v, "layoutDir", "mylayouts", 5)
	setConfigAndWriteSomeFilesTo(afs, v, "staticDir", "mystatic", 6)
	setConfigAndWriteSomeFilesTo(afs, v, "dataDir", "mydata", 7)
	setConfigAndWriteSomeFilesTo(afs, v, "archetypeDir", "myarchetypes", 8)
	setConfigAndWriteSomeFilesTo(afs, v, "assetDir", "myassets", 9)
	setConfigAndWriteSomeFilesTo(afs, v, "resourceDir", "myrsesource", 10)

	conf := testconfig.GetTestConfig(afs, v)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())

	p, err := paths.New(fs, conf)
	c.Assert(err, qt.IsNil)

	bfs, err := filesystems.NewBase(p, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(bfs, qt.Not(qt.IsNil))

	root, err := bfs.I18n.Fs.Open("")
	c.Assert(err, qt.IsNil)
	dirnames, err := root.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.DeepEquals, []string{"f1.txt", "f2.txt", "f3.txt", "f4.txt", "f3.txt", "theme-file-btheme.txt", "f3.txt", "theme-file-atheme.txt"})

	root, err = bfs.Data.Fs.Open("")
	c.Assert(err, qt.IsNil)
	dirnames, err = root.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.DeepEquals, []string{"f1.txt", "f2.txt", "f3.txt", "f4.txt", "f5.txt", "f6.txt", "f7.txt", "f3.txt", "theme-file-btheme.txt", "f3.txt", "theme-file-atheme.txt"})

	checkFileCount(bfs.Layouts.Fs, "", c, 7)

	checkFileCount(bfs.Content.Fs, "", c, 3)
	checkFileCount(bfs.I18n.Fs, "", c, 8) // 4 + 4 themes

	checkFileCount(bfs.Static[""].Fs, "", c, 6)
	checkFileCount(bfs.Data.Fs, "", c, 11)       // 7 + 4 themes
	checkFileCount(bfs.Archetypes.Fs, "", c, 10) // 8 + 2 themes
	checkFileCount(bfs.Assets.Fs, "", c, 9)
	checkFileCount(bfs.Work, "", c, 90)

	c.Assert(bfs.IsStatic(filepath.Join(workingDir, "mystatic", "file1.txt")), qt.Equals, true)

	contentFilename := filepath.Join(workingDir, "mycontent", "file1.txt")
	c.Assert(bfs.IsContent(contentFilename), qt.Equals, true)
	// Check Work fs vs theme
	checkFileContent(bfs.Work, "file-root.txt", c, "content-project")
	checkFileContent(bfs.Work, "theme-root-atheme.txt", c, "content:atheme")

	// https://github.com/gohugoio/hugo/issues/5318
	// Check both project and theme.
	for _, fs := range []afero.Fs{bfs.Archetypes.Fs, bfs.Layouts.Fs} {
		for _, filename := range []string{"/f1.txt", "/theme-file-atheme.txt"} {
			filename = filepath.FromSlash(filename)
			f, err := fs.Open(filename)
			c.Assert(err, qt.IsNil)
			f.Close()
		}
	}
}

func TestNewBaseFsEmpty(t *testing.T) {
	c := qt.New(t)
	afs := afero.NewMemMapFs()
	conf := testconfig.GetTestConfig(afs, nil)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())
	p, err := paths.New(fs, conf)
	c.Assert(err, qt.IsNil)
	bfs, err := filesystems.NewBase(p, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(bfs, qt.Not(qt.IsNil))
	c.Assert(bfs.Archetypes.Fs, qt.Not(qt.IsNil))
	c.Assert(bfs.Layouts.Fs, qt.Not(qt.IsNil))
	c.Assert(bfs.Data.Fs, qt.Not(qt.IsNil))
	c.Assert(bfs.I18n.Fs, qt.Not(qt.IsNil))
	c.Assert(bfs.Work, qt.Not(qt.IsNil))
	c.Assert(bfs.Content.Fs, qt.Not(qt.IsNil))
	c.Assert(bfs.Static, qt.Not(qt.IsNil))
}

func TestRealDirs(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	root, themesDir := t.TempDir(), t.TempDir()
	v.Set("workingDir", root)
	v.Set("themesDir", themesDir)
	v.Set("assetDir", "myassets")
	v.Set("theme", "mytheme")

	afs := &hugofs.OpenFilesFs{Fs: hugofs.Os}

	c.Assert(afs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf1"), 0o755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf2"), 0o755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf2"), 0o755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf3"), 0o755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(root, "resources"), 0o755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(themesDir, "mytheme", "resources"), 0o755), qt.IsNil)

	c.Assert(afs.MkdirAll(filepath.Join(root, "myassets", "js", "f2"), 0o755), qt.IsNil)

	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "scss", "sf1", "a1.scss")), []byte("content"), 0o755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "scss", "sf2", "a3.scss")), []byte("content"), 0o755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "scss", "a2.scss")), []byte("content"), 0o755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf2", "a3.scss")), []byte("content"), 0o755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf3", "a4.scss")), []byte("content"), 0o755)

	afero.WriteFile(afs, filepath.Join(filepath.Join(themesDir, "mytheme", "resources", "t1.txt")), []byte("content"), 0o755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "resources", "p1.txt")), []byte("content"), 0o755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "resources", "p2.txt")), []byte("content"), 0o755)

	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "js", "f2", "a1.js")), []byte("content"), 0o755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "js", "a2.js")), []byte("content"), 0o755)

	conf := testconfig.GetTestConfig(afs, v)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())
	p, err := paths.New(fs, conf)
	c.Assert(err, qt.IsNil)
	bfs, err := filesystems.NewBase(p, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(bfs, qt.Not(qt.IsNil))

	checkFileCount(bfs.Assets.Fs, "", c, 6)

	realDirs := bfs.Assets.RealDirs("scss")
	c.Assert(len(realDirs), qt.Equals, 2)
	c.Assert(realDirs[0], qt.Equals, filepath.Join(root, "myassets/scss"))
	c.Assert(realDirs[len(realDirs)-1], qt.Equals, filepath.Join(themesDir, "mytheme/assets/scss"))

	realDirs = bfs.Assets.RealDirs("foo")
	c.Assert(len(realDirs), qt.Equals, 0)

	c.Assert(afs.OpenFiles(), qt.HasLen, 0)
}

func TestWatchFilenames(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
theme = "t1"
[[module.mounts]]
source = 'content'
target = 'content'
[[module.mounts]]
source = 'content2'
target = 'content/c2'
[[module.mounts]]
source = "hugo_stats.json"
target = "assets/watching/hugo_stats.json"
-- hugo_stats.json --
Some stats.
-- content/foo.md --
foo
-- content2/bar.md --
-- themes/t1/layouts/_default/single.html --
{{ .Content }}
-- themes/t1/static/f1.txt --
`
	b := hugolib.Test(t, files)
	bfs := b.H.BaseFs
	watchFilenames := bfs.WatchFilenames()
	//   []string{"/hugo_stats.json", "/content", "/content2", "/themes/t1/layouts", "/themes/t1/layouts/_default", "/themes/t1/static"}
	b.Assert(watchFilenames, qt.HasLen, 6)
}

func TestNoSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on Windows")
	}
	files := `
-- hugo.toml --
theme = "t1"
-- content/a/foo.md --
foo
-- static/a/f1.txt --
F1 text
-- themes/t1/layouts/_default/single.html --
{{ .Content }}
-- themes/t1/static/a/f1.txt --
`
	tmpDir := t.TempDir()

	wd, _ := os.Getwd()

	for _, component := range []string{"content", "static"} {
		aDir := filepath.Join(tmpDir, component, "a")
		bDir := filepath.Join(tmpDir, component, "b")
		os.MkdirAll(aDir, 0o755)
		os.MkdirAll(bDir, 0o755)
		os.Chdir(bDir)
		os.Symlink("../a", "c")
	}

	os.Chdir(wd)

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			WorkingDir:  tmpDir,
		},
	).Build()

	bfs := b.H.BaseFs
	watchFilenames := bfs.WatchFilenames()
	b.Assert(watchFilenames, qt.HasLen, 10)
}

func TestStaticFs(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	workDir := "mywork"
	v.Set("workingDir", workDir)
	v.Set("themesDir", "themes")
	v.Set("staticDir", "mystatic")
	v.Set("theme", []string{"t1", "t2"})

	afs := afero.NewMemMapFs()

	themeStaticDir := filepath.Join(workDir, "themes", "t1", "static")
	themeStaticDir2 := filepath.Join(workDir, "themes", "t2", "static")

	afero.WriteFile(afs, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0o755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0o755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0o755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir2, "f2.txt"), []byte("Hugo Themes Rocks in t2!"), 0o755)

	conf := testconfig.GetTestConfig(afs, v)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())
	p, err := paths.New(fs, conf)

	c.Assert(err, qt.IsNil)
	bfs, err := filesystems.NewBase(p, nil)
	c.Assert(err, qt.IsNil)

	sfs := bfs.StaticFs("en")

	checkFileContent(sfs, "f1.txt", c, "Hugo Rocks!")
	checkFileContent(sfs, "f2.txt", c, "Hugo Themes Still Rocks!")
}

func TestStaticFsMultihost(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	workDir := "mywork"
	v.Set("workingDir", workDir)
	v.Set("themesDir", "themes")
	v.Set("staticDir", "mystatic")
	v.Set("theme", "t1")
	v.Set("defaultContentLanguage", "en")

	langConfig := map[string]any{
		"no": map[string]any{
			"staticDir": "static_no",
			"baseURL":   "https://example.org/no/",
		},
		"en": map[string]any{
			"baseURL": "https://example.org/en/",
		},
	}

	v.Set("languages", langConfig)

	afs := afero.NewMemMapFs()

	themeStaticDir := filepath.Join(workDir, "themes", "t1", "static")

	afero.WriteFile(afs, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0o755)
	afero.WriteFile(afs, filepath.Join(workDir, "static_no", "f1.txt"), []byte("Hugo Rocks in Norway!"), 0o755)

	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0o755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0o755)

	conf := testconfig.GetTestConfig(afs, v)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())

	p, err := paths.New(fs, conf)
	c.Assert(err, qt.IsNil)
	bfs, err := filesystems.NewBase(p, nil)
	c.Assert(err, qt.IsNil)
	enFs := bfs.StaticFs("en")
	checkFileContent(enFs, "f1.txt", c, "Hugo Rocks!")
	checkFileContent(enFs, "f2.txt", c, "Hugo Themes Still Rocks!")

	noFs := bfs.StaticFs("no")
	checkFileContent(noFs, "f1.txt", c, "Hugo Rocks in Norway!")
	checkFileContent(noFs, "f2.txt", c, "Hugo Themes Still Rocks!")
}

func TestMakePathRelative(t *testing.T) {
	files := `
-- hugo.toml --
[[module.mounts]]
source = "bar.txt"
target = "assets/foo/baz.txt"
[[module.imports]]
path = "t1"
[[module.imports.mounts]]
source = "src"
target = "assets/foo/bar"
-- bar.txt --
Bar.
-- themes/t1/src/main.js --
Main.
`
	b := hugolib.Test(t, files)

	rel, found := b.H.BaseFs.Assets.MakePathRelative(filepath.FromSlash("/themes/t1/src/main.js"), true)
	b.Assert(found, qt.Equals, true)
	b.Assert(rel, qt.Equals, filepath.FromSlash("foo/bar/main.js"))

	rel, found = b.H.BaseFs.Assets.MakePathRelative(filepath.FromSlash("/bar.txt"), true)
	b.Assert(found, qt.Equals, true)
	b.Assert(rel, qt.Equals, filepath.FromSlash("foo/baz.txt"))
}

func TestAbsProjectContentDir(t *testing.T) {
	tempDir := t.TempDir()

	files := `
-- hugo.toml --
[[module.mounts]]
source = "content"
target = "content"
-- content/foo.md --
---
title: "Foo"
---
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			WorkingDir:  tempDir,
			TxtarString: files,
		},
	).Build()

	abs1 := filepath.Join(tempDir, "content", "foo.md")
	rel, abs2, err := b.H.BaseFs.AbsProjectContentDir("foo.md")
	b.Assert(err, qt.IsNil)
	b.Assert(abs2, qt.Equals, abs1)
	b.Assert(rel, qt.Equals, filepath.FromSlash("foo.md"))
	rel2, abs3, err := b.H.BaseFs.AbsProjectContentDir(abs1)
	b.Assert(err, qt.IsNil)
	b.Assert(abs3, qt.Equals, abs1)
	b.Assert(rel2, qt.Equals, rel)
}

func TestContentReverseLookup(t *testing.T) {
	files := `
-- README.md --
---
title: README
---
-- blog/b1.md --
---
title: b1
---
-- docs/d1.md --
---
title: d1
---
-- hugo.toml --
baseURL = "https://example.com/"
[module]
[[module.mounts]]
source = "layouts"
target = "layouts"
[[module.mounts]]
source = "README.md"
target = "content/_index.md"
[[module.mounts]]
source = "blog"
target = "content/posts"
[[module.mounts]]
source = "docs"
target = "content/mydocs"
-- layouts/index.html --
Home.

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Home.")

	stat := func(path string) hugofs.FileMetaInfo {
		ps, err := b.H.BaseFs.Content.ReverseLookup(filepath.FromSlash(path), true)
		b.Assert(err, qt.IsNil)
		b.Assert(ps, qt.HasLen, 1)
		first := ps[0]
		fi, err := b.H.BaseFs.Content.Fs.Stat(filepath.FromSlash(first.Path))
		b.Assert(err, qt.IsNil)
		b.Assert(fi, qt.Not(qt.IsNil))
		return fi.(hugofs.FileMetaInfo)
	}

	sfs := b.H.Fs.Source

	_, err := sfs.Stat("blog/b1.md")
	b.Assert(err, qt.Not(qt.IsNil))

	_ = stat("blog/b1.md")
}

func TestReverseLookupShouldOnlyConsiderFilesInCurrentComponent(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com/"
[module]
[[module.mounts]]
source = "files/layouts"
target = "layouts"
[[module.mounts]]
source = "files/layouts/assets"
target = "assets"
-- files/layouts/l1.txt --
l1
-- files/layouts/assets/l2.txt --
l2
`
	b := hugolib.Test(t, files)

	assetsFs := b.H.Assets

	for _, checkExists := range []bool{false, true} {
		cps, err := assetsFs.ReverseLookup(filepath.FromSlash("files/layouts/assets/l2.txt"), checkExists)
		b.Assert(err, qt.IsNil)
		b.Assert(cps, qt.HasLen, 1)
		cps, err = assetsFs.ReverseLookup(filepath.FromSlash("files/layouts/l2.txt"), checkExists)
		b.Assert(err, qt.IsNil)
		b.Assert(cps, qt.HasLen, 0)
	}
}

func TestAssetsIssue12175(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com/"
[module]
[[module.mounts]]
source = "node_modules/@foo/core/assets"
target = "assets"
[[module.mounts]]
source = "assets"
target = "assets"
-- node_modules/@foo/core/assets/js/app.js --
JS.
-- node_modules/@foo/core/assets/scss/app.scss --
body { color: red; }
-- assets/scss/app.scss --
body { color: blue; }
-- layouts/index.html --
Home.
SCSS: {{ with resources.Get "scss/app.scss" }}{{ .RelPermalink }}|{{ .Content }}{{ end }}|
# Note that the pattern below will match 2 resources, which doesn't make much sense,
# but is how the current (and also < v0.123.0) merge logic works, and for most practical purposes, it doesn't matter.
SCSS Match: {{ with resources.Match "**.scss" }}{{ . | len }}|{{ range .}}{{ .RelPermalink }}|{{ end }}{{ end }}|

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
SCSS: /scss/app.scss|body { color: blue; }|
SCSS Match: 2|
`)
}

func TestStaticComposite(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[module]
[[module.mounts]]
source = "myfiles/f1.txt"
target = "static/files/f1.txt"
[[module.mounts]]
source = "f3.txt"
target = "static/f3.txt"
[[module.mounts]]
source = "static"
target = "static"
-- static/files/f2.txt --
f2
-- myfiles/f1.txt --
f1
-- f3.txt --
f3
-- layouts/home.html --
Home.

`
	b := hugolib.Test(t, files)

	b.AssertFs(b.H.BaseFs.StaticFs(""), `
. true
f3.txt false
files true
files/f1.txt false
files/f2.txt false
`)
}

func TestMountIssue12141(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[module]
[[module.mounts]]
source = "myfiles"
target = "static"
[[module.mounts]]
source = "myfiles/f1.txt"
target = "static/f2.txt"
-- myfiles/f1.txt --
f1
`
	b := hugolib.Test(t, files)
	fs := b.H.BaseFs.StaticFs("")

	b.AssertFs(fs, `
. true
f1.txt false
f2.txt false
`)
}

func checkFileCount(fs afero.Fs, dirname string, c *qt.C, expected int) {
	c.Helper()
	count, names, err := countFilesAndGetFilenames(fs, dirname)
	namesComment := qt.Commentf("filenames: %v", names)
	c.Assert(err, qt.IsNil, namesComment)
	c.Assert(count, qt.Equals, expected, namesComment)
}

func checkFileContent(fs afero.Fs, filename string, c *qt.C, expected ...string) {
	b, err := afero.ReadFile(fs, filename)
	c.Assert(err, qt.IsNil)

	content := string(b)

	for _, e := range expected {
		c.Assert(content, qt.Contains, e)
	}
}

func countFilesAndGetFilenames(fs afero.Fs, dirname string) (int, []string, error) {
	if fs == nil {
		return 0, nil, errors.New("no fs")
	}

	counter := 0
	var filenames []string

	wf := func(path string, info hugofs.FileMetaInfo) error {
		if !info.IsDir() {
			counter++
		}

		if info.Name() != "." {
			name := info.Name()
			name = strings.Replace(name, filepath.FromSlash("/my/work"), "WORK_DIR", 1)
			filenames = append(filenames, name)
		}

		return nil
	}

	w := hugofs.NewWalkway(hugofs.WalkwayConfig{Fs: fs, Root: dirname, WalkFn: wf})

	if err := w.Walk(); err != nil {
		return -1, nil, err
	}

	return counter, filenames, nil
}

func setConfigAndWriteSomeFilesTo(fs afero.Fs, v config.Provider, key, val string, num int) {
	workingDir := v.GetString("workingDir")
	v.Set(key, val)
	fs.Mkdir(val, 0o755)
	for i := 0; i < num; i++ {
		filename := filepath.Join(workingDir, val, fmt.Sprintf("f%d.txt", i+1))
		afero.WriteFile(fs, filename, []byte(fmt.Sprintf("content:%s:%d", key, i+1)), 0o755)
	}
}
