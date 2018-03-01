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
	}

	afero.WriteFile(fs.Source, filepath.Join(workingDir, "themes", "btheme", "config.toml"), []byte(`
theme = ["atheme"]
`), 0755)

	setConfigAndWriteSomeFilesTo(fs.Source, v, "contentDir", "mycontent", 3)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "i18nDir", "myi18n", 4)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "layoutDir", "mylayouts", 5)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "staticDir", "mystatic", 6)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "dataDir", "mydata", 7)
	setConfigAndWriteSomeFilesTo(fs.Source, v, "archetypeDir", "myarchetypes", 8)

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

	checkFileCount(bfs.ContentFs, "", assert, 3)
	checkFileCount(bfs.I18n.Fs, "", assert, 6) // 4 + 2 themes
	checkFileCount(bfs.Layouts.Fs, "", assert, 5)
	checkFileCount(bfs.Static[""].Fs, "", assert, 6)
	checkFileCount(bfs.Data.Fs, "", assert, 9) // 7 + 2 themes
	checkFileCount(bfs.Archetypes.Fs, "", assert, 8)

	assert.Equal([]string{filepath.FromSlash("/my/work/mydata"), filepath.FromSlash("/my/work/themes/btheme/data"), filepath.FromSlash("/my/work/themes/atheme/data")}, bfs.Data.Dirnames)

	assert.True(bfs.IsData(filepath.Join(workingDir, "mydata", "file1.txt")))
	assert.True(bfs.IsI18n(filepath.Join(workingDir, "myi18n", "file1.txt")))
	assert.True(bfs.IsLayout(filepath.Join(workingDir, "mylayouts", "file1.txt")))
	assert.True(bfs.IsStatic(filepath.Join(workingDir, "mystatic", "file1.txt")))
	contentFilename := filepath.Join(workingDir, "mycontent", "file1.txt")
	assert.True(bfs.IsContent(contentFilename))
	rel, _ := bfs.RelContentDir(contentFilename)
	assert.Equal("file1.txt", rel)

}

func TestNewBaseFsEmpty(t *testing.T) {
	assert := require.New(t)
	v := viper.New()
	v.Set("contentDir", "mycontent")
	v.Set("i18nDir", "myi18n")
	v.Set("staticDir", "mystatic")
	v.Set("dataDir", "mydata")
	v.Set("layoutDir", "mylayouts")
	v.Set("archetypeDir", "myarchetypes")

	fs := hugofs.NewMem(v)
	p, err := paths.New(fs, v)
	bfs, err := NewBase(p)
	assert.NoError(err)
	assert.NotNil(bfs)
	assert.Equal(hugofs.NoOpFs, bfs.Archetypes.Fs)
	assert.Equal(hugofs.NoOpFs, bfs.Layouts.Fs)
	assert.Equal(hugofs.NoOpFs, bfs.Data.Fs)
	assert.Equal(hugofs.NoOpFs, bfs.I18n.Fs)
	assert.NotNil(hugofs.NoOpFs, bfs.ContentFs)
	assert.NotNil(hugofs.NoOpFs, bfs.Static)
}

func checkFileCount(fs afero.Fs, dirname string, assert *require.Assertions, expected int) {
	count, _, err := countFileaAndGetDirs(fs, dirname)
	assert.NoError(err)
	assert.Equal(expected, count)
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
