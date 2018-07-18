// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"testing"

	"github.com/gohugoio/hugo/langs"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/paths"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewBaseFs(t *testing.T) {
	assert := require.New(t)
	v := viper.New()

	fs := hugofs.NewMem(v)

	themes := []string{"btheme", "atheme"}

	workingDir := filepath.FromSlash("/my/work")
	v.Set("workingDir", workingDir)
	v.Set("themesDir", "themes")
	v.Set("theme", themes[:1])

	// Write some data to the themes
	for _, theme := range themes {
		for _, dir := range []string{"i18n", "data"} {
			base := filepath.Join(workingDir, "themes", theme, dir)
			fs.Source.Mkdir(base, 0755)
			afero.WriteFile(fs.Source, filepath.Join(base, fmt.Sprintf("theme-file-%s-%s.txt", theme, dir)), []byte(fmt.Sprintf("content:%s:%s", theme, dir)), 0755)
		}
		// Write some files to the root of the theme
		base := filepath.Join(workingDir, "themes", theme)
		afero.WriteFile(fs.Source, filepath.Join(base, fmt.Sprintf("theme-root-%s.txt", theme)), []byte(fmt.Sprintf("content:%s", theme)), 0755)
		afero.WriteFile(fs.Source, filepath.Join(base, "file-root.txt"), []byte(fmt.Sprintf("content:%s", theme)), 0755)
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

	p, err := paths.New(fs, v)
	assert.NoError(err)

	bfs, err := NewBase(p)
	assert.NoError(err)
	assert.NotNil(bfs)

	root, err := bfs.I18n.Fs.Open("")
	assert.NoError(err)
	dirnames, err := root.Readdirnames(-1)
	assert.NoError(err)
	assert.Equal([]string{projectVirtualFolder, "btheme", "atheme"}, dirnames)
	ff, err := bfs.I18n.Fs.Open("myi18n")
	assert.NoError(err)
	_, err = ff.Readdirnames(-1)
	assert.NoError(err)

	root, err = bfs.Data.Fs.Open("")
	assert.NoError(err)
	dirnames, err = root.Readdirnames(-1)
	assert.NoError(err)
	assert.Equal([]string{projectVirtualFolder, "btheme", "atheme"}, dirnames)
	ff, err = bfs.I18n.Fs.Open("mydata")
	assert.NoError(err)
	_, err = ff.Readdirnames(-1)
	assert.NoError(err)

	checkFileCount(bfs.Content.Fs, "", assert, 3)
	checkFileCount(bfs.I18n.Fs, "", assert, 6) // 4 + 2 themes
	checkFileCount(bfs.Layouts.Fs, "", assert, 5)
	checkFileCount(bfs.Static[""].Fs, "", assert, 6)
	checkFileCount(bfs.Data.Fs, "", assert, 9) // 7 + 2 themes
	checkFileCount(bfs.Archetypes.Fs, "", assert, 8)
	checkFileCount(bfs.Assets.Fs, "", assert, 9)
	checkFileCount(bfs.Resources.Fs, "", assert, 10)
	checkFileCount(bfs.Work.Fs, "", assert, 69)

	assert.Equal([]string{filepath.FromSlash("/my/work/mydata"), filepath.FromSlash("/my/work/themes/btheme/data"), filepath.FromSlash("/my/work/themes/atheme/data")}, bfs.Data.Dirnames)

	assert.True(bfs.IsData(filepath.Join(workingDir, "mydata", "file1.txt")))
	assert.True(bfs.IsI18n(filepath.Join(workingDir, "myi18n", "file1.txt")))
	assert.True(bfs.IsLayout(filepath.Join(workingDir, "mylayouts", "file1.txt")))
	assert.True(bfs.IsStatic(filepath.Join(workingDir, "mystatic", "file1.txt")))
	assert.True(bfs.IsAsset(filepath.Join(workingDir, "myassets", "file1.txt")))

	contentFilename := filepath.Join(workingDir, "mycontent", "file1.txt")
	assert.True(bfs.IsContent(contentFilename))
	rel := bfs.RelContentDir(contentFilename)
	assert.Equal("file1.txt", rel)

	// Check Work fs vs theme
	checkFileContent(bfs.Work.Fs, "file-root.txt", assert, "content-project")
	checkFileContent(bfs.Work.Fs, "theme-root-atheme.txt", assert, "content:atheme")

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

	return v
}

func TestNewBaseFsEmpty(t *testing.T) {
	assert := require.New(t)
	v := createConfig()
	fs := hugofs.NewMem(v)
	p, err := paths.New(fs, v)
	assert.NoError(err)
	bfs, err := NewBase(p)
	assert.NoError(err)
	assert.NotNil(bfs)
	assert.Equal(hugofs.NoOpFs, bfs.Archetypes.Fs)
	assert.Equal(hugofs.NoOpFs, bfs.Layouts.Fs)
	assert.Equal(hugofs.NoOpFs, bfs.Data.Fs)
	assert.Equal(hugofs.NoOpFs, bfs.Assets.Fs)
	assert.Equal(hugofs.NoOpFs, bfs.I18n.Fs)
	assert.NotNil(bfs.Work.Fs)
	assert.NotNil(bfs.Content.Fs)
	assert.NotNil(bfs.Static)
}

func TestRealDirs(t *testing.T) {
	assert := require.New(t)
	v := createConfig()
	fs := hugofs.NewDefault(v)
	sfs := fs.Source

	root, err := afero.TempDir(sfs, "", "realdir")
	assert.NoError(err)
	themesDir, err := afero.TempDir(sfs, "", "themesDir")
	assert.NoError(err)
	defer func() {
		os.RemoveAll(root)
		os.RemoveAll(themesDir)
	}()

	v.Set("workingDir", root)
	v.Set("themesDir", themesDir)
	v.Set("theme", "mytheme")

	assert.NoError(sfs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf1"), 0755))
	assert.NoError(sfs.MkdirAll(filepath.Join(root, "myassets", "scss", "sf2"), 0755))
	assert.NoError(sfs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf2"), 0755))
	assert.NoError(sfs.MkdirAll(filepath.Join(themesDir, "mytheme", "assets", "scss", "sf3"), 0755))
	assert.NoError(sfs.MkdirAll(filepath.Join(root, "resources"), 0755))
	assert.NoError(sfs.MkdirAll(filepath.Join(themesDir, "mytheme", "resources"), 0755))

	assert.NoError(sfs.MkdirAll(filepath.Join(root, "myassets", "js", "f2"), 0755))

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

	p, err := paths.New(fs, v)
	assert.NoError(err)
	bfs, err := NewBase(p)
	assert.NoError(err)
	assert.NotNil(bfs)
	checkFileCount(bfs.Assets.Fs, "", assert, 6)

	realDirs := bfs.Assets.RealDirs("scss")
	assert.Equal(2, len(realDirs))
	assert.Equal(filepath.Join(root, "myassets/scss"), realDirs[0])
	assert.Equal(filepath.Join(themesDir, "mytheme/assets/scss"), realDirs[len(realDirs)-1])

	checkFileCount(bfs.Resources.Fs, "", assert, 3)

}

func TestStaticFs(t *testing.T) {
	assert := require.New(t)
	v := createConfig()
	workDir := "mywork"
	v.Set("workingDir", workDir)
	v.Set("themesDir", "themes")
	v.Set("theme", "t1")

	fs := hugofs.NewMem(v)

	themeStaticDir := filepath.Join(workDir, "themes", "t1", "static")

	afero.WriteFile(fs.Source, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0755)

	p, err := paths.New(fs, v)
	assert.NoError(err)
	bfs, err := NewBase(p)
	sfs := bfs.StaticFs("en")
	checkFileContent(sfs, "f1.txt", assert, "Hugo Rocks!")
	checkFileContent(sfs, "f2.txt", assert, "Hugo Themes Still Rocks!")

}

func TestStaticFsMultiHost(t *testing.T) {
	assert := require.New(t)
	v := createConfig()
	workDir := "mywork"
	v.Set("workingDir", workDir)
	v.Set("themesDir", "themes")
	v.Set("theme", "t1")
	v.Set("multihost", true)

	vn := viper.New()
	vn.Set("staticDir", "nn_static")

	en := langs.NewLanguage("en", v)
	no := langs.NewLanguage("no", v)
	no.Set("staticDir", "static_no")

	languages := langs.Languages{
		en,
		no,
	}

	v.Set("languagesSorted", languages)

	fs := hugofs.NewMem(v)

	themeStaticDir := filepath.Join(workDir, "themes", "t1", "static")

	afero.WriteFile(fs.Source, filepath.Join(workDir, "mystatic", "f1.txt"), []byte("Hugo Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(workDir, "static_no", "f1.txt"), []byte("Hugo Rocks in Norway!"), 0755)

	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f1.txt"), []byte("Hugo Themes Rocks!"), 0755)
	afero.WriteFile(fs.Source, filepath.Join(themeStaticDir, "f2.txt"), []byte("Hugo Themes Still Rocks!"), 0755)

	p, err := paths.New(fs, v)
	assert.NoError(err)
	bfs, err := NewBase(p)
	enFs := bfs.StaticFs("en")
	checkFileContent(enFs, "f1.txt", assert, "Hugo Rocks!")
	checkFileContent(enFs, "f2.txt", assert, "Hugo Themes Still Rocks!")

	noFs := bfs.StaticFs("no")
	checkFileContent(noFs, "f1.txt", assert, "Hugo Rocks in Norway!")
	checkFileContent(noFs, "f2.txt", assert, "Hugo Themes Still Rocks!")
}

func checkFileCount(fs afero.Fs, dirname string, assert *require.Assertions, expected int) {
	count, _, err := countFileaAndGetDirs(fs, dirname)
	assert.NoError(err)
	assert.Equal(expected, count)
}

func checkFileContent(fs afero.Fs, filename string, assert *require.Assertions, expected ...string) {

	b, err := afero.ReadFile(fs, filename)
	assert.NoError(err)

	content := string(b)

	for _, e := range expected {
		assert.Contains(content, e)
	}
}

func countFileaAndGetDirs(fs afero.Fs, dirname string) (int, []string, error) {
	if fs == nil {
		return 0, nil, errors.New("no fs")
	}

	counter := 0
	var dirs []string

	afero.Walk(fs, dirname, func(path string, info os.FileInfo, err error) error {
		if info != nil {
			if !info.IsDir() {
				counter++
			} else if info.Name() != "." {
				dirs = append(dirs, filepath.Join(path, info.Name()))
			}
		}

		return nil
	})

	return counter, dirs, nil
}

func setConfigAndWriteSomeFilesTo(fs afero.Fs, v *viper.Viper, key, val string, num int) {
	workingDir := v.GetString("workingDir")
	v.Set(key, val)
	fs.Mkdir(val, 0755)
	for i := 0; i < num; i++ {
		afero.WriteFile(fs, filepath.Join(workingDir, val, fmt.Sprintf("file%d.txt", i+1)), []byte(fmt.Sprintf("content:%s:%d", key, i+1)), 0755)
	}
}
