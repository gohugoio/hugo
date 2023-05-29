// Copyright 2023 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"

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
			afs.Mkdir(base, 0755)
			content := []byte(fmt.Sprintf("content:%s:%s", theme, dir))
			afero.WriteFile(afs, filenameTheme, content, 0755)
			afero.WriteFile(afs, filenameOverlap, content, 0755)
		}
		// Write some files to the root of the theme
		base := filepath.Join(workingDir, "themes", theme)
		afero.WriteFile(afs, filepath.Join(base, fmt.Sprintf("theme-root-%s.txt", theme)), []byte(fmt.Sprintf("content:%s", theme)), 0755)
		afero.WriteFile(afs, filepath.Join(base, "file-theme-root.txt"), []byte(fmt.Sprintf("content:%s", theme)), 0755)
	}

	afero.WriteFile(afs, filepath.Join(workingDir, "file-root.txt"), []byte("content-project"), 0755)

	afero.WriteFile(afs, filepath.Join(workingDir, "themes", "btheme", "config.toml"), []byte(`
theme = ["atheme"]
`), 0755)

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

	c.Assert(bfs.IsData(filepath.Join(workingDir, "mydata", "file1.txt")), qt.Equals, true)
	c.Assert(bfs.IsI18n(filepath.Join(workingDir, "myi18n", "file1.txt")), qt.Equals, true)
	c.Assert(bfs.IsLayout(filepath.Join(workingDir, "mylayouts", "file1.txt")), qt.Equals, true)
	c.Assert(bfs.IsStatic(filepath.Join(workingDir, "mystatic", "file1.txt")), qt.Equals, true)
	c.Assert(bfs.IsAsset(filepath.Join(workingDir, "myassets", "file1.txt")), qt.Equals, true)

	contentFilename := filepath.Join(workingDir, "mycontent", "file1.txt")
	c.Assert(bfs.IsContent(contentFilename), qt.Equals, true)
	rel := bfs.RelContentDir(contentFilename)
	c.Assert(rel, qt.Equals, "file1.txt")

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

	afs := hugofs.Os

	defer func() {
		os.RemoveAll(root)
		os.RemoveAll(themesDir)
	}()

	c.Assert(afs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf1"), 0755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf2"), 0755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf2"), 0755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf3"), 0755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(root, "resources"), 0755), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(themesDir, "mytheme", "resources"), 0755), qt.IsNil)

	c.Assert(afs.MkdirAll(filepath.Join(root, "myassets", "js", "f2"), 0755), qt.IsNil)

	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "scss", "sf1", "a1.scss")), []byte("content"), 0755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "scss", "sf2", "a3.scss")), []byte("content"), 0755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "scss", "a2.scss")), []byte("content"), 0755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf2", "a3.scss")), []byte("content"), 0755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf3", "a4.scss")), []byte("content"), 0755)

	afero.WriteFile(afs, filepath.Join(filepath.Join(themesDir, "mytheme", "resources", "t1.txt")), []byte("content"), 0755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "resources", "p1.txt")), []byte("content"), 0755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "resources", "p2.txt")), []byte("content"), 0755)

	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "js", "f2", "a1.js")), []byte("content"), 0755)
	afero.WriteFile(afs, filepath.Join(filepath.Join(root, "myassets", "js", "a2.js")), []byte("content"), 0755)

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

	afero.WriteFile(afs, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir2, "f2.txt"), []byte("Hugo Themes Rocks in t2!"), 0755)

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

func TestStaticFsMultiHost(t *testing.T) {
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

	afero.WriteFile(afs, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0755)
	afero.WriteFile(afs, filepath.Join(workDir, "static_no", "f1.txt"), []byte("Hugo Rocks in Norway!"), 0755)

	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0755)
	afero.WriteFile(afs, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0755)

	conf := testconfig.GetTestConfig(afs, v)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())

	fmt.Println("IS", conf.IsMultihost())

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
	c := qt.New(t)
	v := config.New()
	afs := afero.NewMemMapFs()
	workDir := "mywork"
	v.Set("workingDir", workDir)

	c.Assert(afs.MkdirAll(filepath.Join(workDir, "dist", "d1"), 0777), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(workDir, "static", "d2"), 0777), qt.IsNil)
	c.Assert(afs.MkdirAll(filepath.Join(workDir, "dust", "d2"), 0777), qt.IsNil)

	moduleCfg := map[string]any{
		"mounts": []any{
			map[string]any{
				"source": "dist",
				"target": "static/mydist",
			},
			map[string]any{
				"source": "dust",
				"target": "static/foo/bar",
			},
			map[string]any{
				"source": "static",
				"target": "static",
			},
		},
	}

	v.Set("module", moduleCfg)

	conf := testconfig.GetTestConfig(afs, v)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())

	p, err := paths.New(fs, conf)
	c.Assert(err, qt.IsNil)
	bfs, err := filesystems.NewBase(p, nil)
	c.Assert(err, qt.IsNil)

	sfs := bfs.Static[""]
	c.Assert(sfs, qt.Not(qt.IsNil))

	makeRel := func(s string) string {
		r, _ := sfs.MakePathRelative(s)
		return r
	}

	c.Assert(makeRel(filepath.Join(workDir, "dist", "d1", "foo.txt")), qt.Equals, filepath.FromSlash("mydist/d1/foo.txt"))
	c.Assert(makeRel(filepath.Join(workDir, "static", "d2", "foo.txt")), qt.Equals, filepath.FromSlash("d2/foo.txt"))
	c.Assert(makeRel(filepath.Join(workDir, "dust", "d3", "foo.txt")), qt.Equals, filepath.FromSlash("foo/bar/d3/foo.txt"))
}

func checkFileCount(fs afero.Fs, dirname string, c *qt.C, expected int) {
	c.Helper()
	count, _, err := countFilesAndGetFilenames(fs, dirname)
	c.Assert(err, qt.IsNil)
	c.Assert(count, qt.Equals, expected)
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

	wf := func(path string, info hugofs.FileMetaInfo, err error) error {
		if err != nil {
			return err
		}
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
	fs.Mkdir(val, 0755)
	for i := 0; i < num; i++ {
		filename := filepath.Join(workingDir, val, fmt.Sprintf("f%d.txt", i+1))
		afero.WriteFile(fs, filename, []byte(fmt.Sprintf("content:%s:%d", key, i+1)), 0755)
	}
}
