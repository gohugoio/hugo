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

package filesystems

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/langs"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/paths"
	"github.com/gohugoio/hugo/modules"
	"github.com/spf13/viper"
)

func initConfig(fs afero.Fs, cfg config.Provider) error {
	if _, err := langs.LoadLanguageSettings(cfg, nil); err != nil {
		return err
	}

	modConfig, err := modules.DecodeConfig(cfg)
	if err != nil {
		return err
	}

	workingDir := cfg.GetString("workingDir")
	themesDir := cfg.GetString("themesDir")
	if !filepath.IsAbs(themesDir) {
		themesDir = filepath.Join(workingDir, themesDir)
	}
	modulesClient := modules.NewClient(modules.ClientConfig{
		Fs:           fs,
		WorkingDir:   workingDir,
		ThemesDir:    themesDir,
		ModuleConfig: modConfig,
		IgnoreVendor: true,
	})

	moduleConfig, err := modulesClient.Collect()
	if err != nil {
		return err
	}

	if err := modules.ApplyProjectConfigDefaults(cfg, moduleConfig.ActiveModules[0]); err != nil {
		return err
	}

	cfg.Set("allModules", moduleConfig.ActiveModules)

	return nil
}

func TestNewBaseFs(t *testing.T) {
	c := qt.New(t)
	v := viper.New()

	fs := hugofs.NewMem(v)

	themes := []string{"btheme", "atheme"}

	workingDir := filepath.FromSlash("/my/work")
	v.Set("workingDir", workingDir)
	v.Set("contentDir", "content")
	v.Set("themesDir", "themes")
	v.Set("defaultContentLanguage", "en")
	v.Set("theme", themes[:1])

	// Write some data to the themes
	for _, theme := range themes {
		for _, dir := range []string{"i18n", "data", "archetypes", "layouts"} {
			base := filepath.Join(workingDir, "themes", theme, dir)
			filenameTheme := filepath.Join(base, fmt.Sprintf("theme-file-%s.txt", theme))
			filenameOverlap := filepath.Join(base, "f3.txt")
			fs.Source.Mkdir(base, 0755)
			content := []byte(fmt.Sprintf("content:%s:%s", theme, dir))
			afero.WriteFile(fs.Source, filenameTheme, content, 0755)
			afero.WriteFile(fs.Source, filenameOverlap, content, 0755)
		}
		// Write some files to the root of the theme
		base := filepath.Join(workingDir, "themes", theme)
		afero.WriteFile(fs.Source, filepath.Join(base, fmt.Sprintf("theme-root-%s.txt", theme)), []byte(fmt.Sprintf("content:%s", theme)), 0755)
		afero.WriteFile(fs.Source, filepath.Join(base, "file-theme-root.txt"), []byte(fmt.Sprintf("content:%s", theme)), 0755)
	}

	afero.WriteFile(fs.Source, filepath.Join(workingDir, "file-root.txt"), []byte("content-project"), 0755)

	afero.WriteFile(fs.Source, filepath.Join(workingDir, "themes", "btheme", "config.toml"), []byte(`
theme = ["atheme"]
`), 0755)

	setConfigAndWriteSomeFilesTo(fs.Source, v, "contentDir", "mycontent", 3)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "i18nDir", "myi18n", 4)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "layoutDir", "mylayouts", 5)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "staticDir", "mystatic", 6)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "dataDir", "mydata", 7)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "archetypeDir", "myarchetypes", 8)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "assetDir", "myassets", 9)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "resourceDir", "myrsesource", 10)

	v.Set("publishDir", "public")
	c.Assert(initConfig(fs.Source, v), qt.IsNil)

	p, err := paths.New(fs, v)
	c.Assert(err, qt.IsNil)

	bfs, err := NewBase(p, nil)
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
	checkFileCount(bfs.Work, "", c, 82)

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

func createConfig() *viper.Viper {
	v := viper.New()
	v.Set("contentDir", "mycontent")
	v.Set("i18nDir", "myi18n")
	v.Set("staticDir", "mystatic")
	v.Set("dataDir", "mydata")
	v.Set("layoutDir", "mylayouts")
	v.Set("archetypeDir", "myarchetypes")
	v.Set("assetDir", "myassets")
	v.Set("resourceDir", "resources")
	v.Set("publishDir", "public")
	v.Set("defaultContentLanguage", "en")

	return v
}

func TestNewBaseFsEmpty(t *testing.T) {
	c := qt.New(t)
	v := createConfig()
	fs := hugofs.NewMem(v)
	c.Assert(initConfig(fs.Source, v), qt.IsNil)

	p, err := paths.New(fs, v)
	c.Assert(err, qt.IsNil)
	bfs, err := NewBase(p, nil)
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
	v := createConfig()
	fs := hugofs.NewDefault(v)
	sfs := fs.Source

	root, err := afero.TempDir(sfs, "", "realdir")
	c.Assert(err, qt.IsNil)
	themesDir, err := afero.TempDir(sfs, "", "themesDir")
	c.Assert(err, qt.IsNil)
	defer func() {
		os.RemoveAll(root)
		os.RemoveAll(themesDir)
	}()

	v.Set("workingDir", root)
	v.Set("themesDir", themesDir)
	v.Set("theme", "mytheme")

	c.Assert(sfs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf1"), 0755), qt.IsNil)
	c.Assert(sfs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf2"), 0755), qt.IsNil)
	c.Assert(sfs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf2"), 0755), qt.IsNil)
	c.Assert(sfs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf3"), 0755), qt.IsNil)
	c.Assert(sfs.MkdirAll(filepath.Join(root, "resources"), 0755), qt.IsNil)
	c.Assert(sfs.MkdirAll(filepath.Join(themesDir, "mytheme", "resources"), 0755), qt.IsNil)

	c.Assert(sfs.MkdirAll(filepath.Join(root, "myassets", "js", "f2"), 0755), qt.IsNil)

	afero.WriteFile(sfs, filepath.Join(filepath.Join(root, "myassets", "scss", "sf1", "a1.scss")), []byte("content"), 0755)
	afero.WriteFile(sfs, filepath.Join(filepath.Join(root, "myassets", "scss", "sf2", "a3.scss")), []byte("content"), 0755)
	afero.WriteFile(sfs, filepath.Join(filepath.Join(root, "myassets", "scss", "a2.scss")), []byte("content"), 0755)
	afero.WriteFile(sfs, filepath.Join(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf2", "a3.scss")), []byte("content"), 0755)
	afero.WriteFile(sfs, filepath.Join(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf3", "a4.scss")), []byte("content"), 0755)

	afero.WriteFile(sfs, filepath.Join(filepath.Join(themesDir, "mytheme", "resources", "t1.txt")), []byte("content"), 0755)
	afero.WriteFile(sfs, filepath.Join(filepath.Join(root, "resources", "p1.txt")), []byte("content"), 0755)
	afero.WriteFile(sfs, filepath.Join(filepath.Join(root, "resources", "p2.txt")), []byte("content"), 0755)

	afero.WriteFile(sfs, filepath.Join(filepath.Join(root, "myassets", "js", "f2", "a1.js")), []byte("content"), 0755)
	afero.WriteFile(sfs, filepath.Join(filepath.Join(root, "myassets", "js", "a2.js")), []byte("content"), 0755)

	c.Assert(initConfig(fs.Source, v), qt.IsNil)

	p, err := paths.New(fs, v)
	c.Assert(err, qt.IsNil)
	bfs, err := NewBase(p, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(bfs, qt.Not(qt.IsNil))

	checkFileCount(bfs.Assets.Fs, "", c, 6)

	realDirs := bfs.Assets.RealDirs("scss")
	c.Assert(len(realDirs), qt.Equals, 2)
	c.Assert(realDirs[0], qt.Equals, filepath.Join(root, "myassets/scss"))
	c.Assert(realDirs[len(realDirs)-1], qt.Equals, filepath.Join(themesDir, "mytheme/assets/scss"))

	c.Assert(bfs.theBigFs, qt.Not(qt.IsNil))

}

func TestStaticFs(t *testing.T) {
	c := qt.New(t)
	v := createConfig()
	workDir := "mywork"
	v.Set("workingDir", workDir)
	v.Set("themesDir", "themes")
	v.Set("theme", []string{"t1", "t2"})

	fs := hugofs.NewMem(v)

	themeStaticDir := filepath.Join(workDir, "themes", "t1", "static")
	themeStaticDir2 := filepath.Join(workDir, "themes", "t2", "static")

	afero.WriteFile(fs.Source, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir2, "f2.txt"), []byte("Hugo Themes Rocks in t2!"), 0755)

	c.Assert(initConfig(fs.Source, v), qt.IsNil)

	p, err := paths.New(fs, v)
	c.Assert(err, qt.IsNil)
	bfs, err := NewBase(p, nil)
	c.Assert(err, qt.IsNil)

	sfs := bfs.StaticFs("en")
	checkFileContent(sfs, "f1.txt", c, "Hugo Rocks!")
	checkFileContent(sfs, "f2.txt", c, "Hugo Themes Still Rocks!")

}

func TestStaticFsMultiHost(t *testing.T) {
	c := qt.New(t)
	v := createConfig()
	workDir := "mywork"
	v.Set("workingDir", workDir)
	v.Set("themesDir", "themes")
	v.Set("theme", "t1")
	v.Set("defaultContentLanguage", "en")

	langConfig := map[string]interface{}{
		"no": map[string]interface{}{
			"staticDir": "static_no",
			"baseURL":   "https://example.org/no/",
		},
		"en": map[string]interface{}{
			"baseURL": "https://example.org/en/",
		},
	}

	v.Set("languages", langConfig)

	fs := hugofs.NewMem(v)

	themeStaticDir := filepath.Join(workDir, "themes", "t1", "static")

	afero.WriteFile(fs.Source, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(workDir, "static_no", "f1.txt"), []byte("Hugo Rocks in Norway!"), 0755)

	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0755)

	c.Assert(initConfig(fs.Source, v), qt.IsNil)

	p, err := paths.New(fs, v)
	c.Assert(err, qt.IsNil)
	bfs, err := NewBase(p, nil)
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
	v := createConfig()
	fs := hugofs.NewMem(v)
	workDir := "mywork"
	v.Set("workingDir", workDir)

	c.Assert(fs.Source.MkdirAll(filepath.Join(workDir, "dist", "d1"), 0777), qt.IsNil)
	c.Assert(fs.Source.MkdirAll(filepath.Join(workDir, "static", "d2"), 0777), qt.IsNil)
	c.Assert(fs.Source.MkdirAll(filepath.Join(workDir, "dust", "d2"), 0777), qt.IsNil)

	moduleCfg := map[string]interface{}{
		"mounts": []interface{}{
			map[string]interface{}{
				"source": "dist",
				"target": "static/mydist",
			},
			map[string]interface{}{
				"source": "dust",
				"target": "static/foo/bar",
			},
			map[string]interface{}{
				"source": "static",
				"target": "static",
			},
		},
	}

	v.Set("module", moduleCfg)

	c.Assert(initConfig(fs.Source, v), qt.IsNil)

	p, err := paths.New(fs, v)
	c.Assert(err, qt.IsNil)
	bfs, err := NewBase(p, nil)
	c.Assert(err, qt.IsNil)

	sfs := bfs.Static[""]
	c.Assert(sfs, qt.Not(qt.IsNil))

	c.Assert(sfs.MakePathRelative(filepath.Join(workDir, "dist", "d1", "foo.txt")), qt.Equals, filepath.FromSlash("mydist/d1/foo.txt"))
	c.Assert(sfs.MakePathRelative(filepath.Join(workDir, "static", "d2", "foo.txt")), qt.Equals, filepath.FromSlash("d2/foo.txt"))
	c.Assert(sfs.MakePathRelative(filepath.Join(workDir, "dust", "d3", "foo.txt")), qt.Equals, filepath.FromSlash("foo/bar/d3/foo.txt"))

}

func checkFileCount(fs afero.Fs, dirname string, c *qt.C, expected int) {
	count, _, err := countFileaAndGetFilenames(fs, dirname)
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

func countFileaAndGetFilenames(fs afero.Fs, dirname string) (int, []string, error) {
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

func setConfigAndWriteSomeFilesTo(fs afero.Fs, v *viper.Viper, key, val string, num int) {
	workingDir := v.GetString("workingDir")
	v.Set(key, val)
	fs.Mkdir(val, 0755)
	for i := 0; i < num; i++ {
		filename := filepath.Join(workingDir, val, fmt.Sprintf("f%d.txt", i+1))
		afero.WriteFile(fs, filename, []byte(fmt.Sprintf("content:%s:%d", key, i+1)), 0755)
	}
}
