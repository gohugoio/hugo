// Copyright 2017 The Hugo Authors. All rights reserved.
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

package source

import (
	"testing"

	"github.com/gohugoio/hugo/helpers"

	"fmt"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var logger = jww.NewNotepad(jww.LevelInfo, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)

func TestStaticDirs(t *testing.T) {
	assert := require.New(t)

	tests := []struct {
		setup    func(cfg config.Provider, fs *hugofs.Fs) config.Provider
		expected []string
	}{

		{func(cfg config.Provider, fs *hugofs.Fs) config.Provider {
			cfg.Set("staticDir", "s1")
			return cfg
		}, []string{"s1"}},
		{func(cfg config.Provider, fs *hugofs.Fs) config.Provider {
			cfg.Set("staticDir", []string{"s2", "s1", "s2"})
			return cfg
		}, []string{"s1", "s2"}},
		{func(cfg config.Provider, fs *hugofs.Fs) config.Provider {
			cfg.Set("theme", "mytheme")
			cfg.Set("themesDir", "themes")
			cfg.Set("staticDir", []string{"s1", "s2"})
			return cfg
		}, []string{filepath.FromSlash("themes/mytheme/static"), "s1", "s2"}},
		{func(cfg config.Provider, fs *hugofs.Fs) config.Provider {
			cfg.Set("staticDir", "s1")

			l1 := helpers.NewLanguage("en", cfg)
			l1.Set("staticDir", []string{"l1s1", "l1s2"})
			return l1

		}, []string{"l1s1", "l1s2"}},
		{func(cfg config.Provider, fs *hugofs.Fs) config.Provider {
			cfg.Set("staticDir", "s1")

			l1 := helpers.NewLanguage("en", cfg)
			l1.Set("staticDir2", []string{"l1s1", "l1s2"})
			return l1

		}, []string{"s1", "l1s1", "l1s2"}},
		{func(cfg config.Provider, fs *hugofs.Fs) config.Provider {
			cfg.Set("staticDir", []string{"s1", "s2"})

			l1 := helpers.NewLanguage("en", cfg)
			l1.Set("staticDir2", []string{"l1s1", "l1s2"})
			return l1

		}, []string{"s1", "s2", "l1s1", "l1s2"}},
		{func(cfg config.Provider, fs *hugofs.Fs) config.Provider {
			cfg.Set("staticDir", "s1")

			l1 := helpers.NewLanguage("en", cfg)
			l1.Set("staticDir2", []string{"l1s1", "l1s2"})
			l2 := helpers.NewLanguage("nn", cfg)
			l2.Set("staticDir3", []string{"l2s1", "l2s2"})
			l2.Set("staticDir", []string{"l2"})

			cfg.Set("languagesSorted", helpers.Languages{l1, l2})
			return cfg

		}, []string{"s1", "l1s1", "l1s2", "l2", "l2s1", "l2s2"}},
	}

	for i, test := range tests {
		msg := fmt.Sprintf("Test %d", i)
		v := viper.New()
		fs := hugofs.NewMem(v)
		cfg := test.setup(v, fs)
		cfg.Set("workingDir", filepath.FromSlash("/work"))
		_, isLanguage := cfg.(*helpers.Language)
		if !isLanguage && !cfg.IsSet("languagesSorted") {
			cfg.Set("languagesSorted", helpers.Languages{helpers.NewDefaultLanguage(cfg)})
		}
		dirs, err := NewDirs(fs, cfg, logger)
		assert.NoError(err)
		assert.Equal(test.expected, dirs.staticDirs, msg)
		assert.Len(dirs.AbsStaticDirs, len(dirs.staticDirs))

		for i, d := range dirs.staticDirs {
			abs := dirs.AbsStaticDirs[i]
			assert.Equal(filepath.Join("/work", d)+helpers.FilePathSeparator, abs)
			assert.True(dirs.IsStatic(filepath.Join(abs, "logo.png")))
			rel := dirs.MakeStaticPathRelative(filepath.Join(abs, "logo.png"))
			assert.Equal("logo.png", rel)
		}

		assert.False(dirs.IsStatic(filepath.FromSlash("/some/other/dir/logo.png")))

	}

}

func TestStaticDirsFs(t *testing.T) {
	assert := require.New(t)
	v := viper.New()
	fs := hugofs.NewMem(v)
	v.Set("workingDir", filepath.FromSlash("/work"))
	v.Set("theme", "mytheme")
	v.Set("themesDir", "themes")
	v.Set("staticDir", []string{"s1", "s2"})
	v.Set("languagesSorted", helpers.Languages{helpers.NewDefaultLanguage(v)})

	writeToFs(t, fs.Source, "/work/s1/f1.txt", "s1-f1")
	writeToFs(t, fs.Source, "/work/s2/f2.txt", "s2-f2")
	writeToFs(t, fs.Source, "/work/s1/f2.txt", "s1-f2")
	writeToFs(t, fs.Source, "/work/themes/mytheme/static/f1.txt", "theme-f1")
	writeToFs(t, fs.Source, "/work/themes/mytheme/static/f3.txt", "theme-f3")

	dirs, err := NewDirs(fs, v, logger)
	assert.NoError(err)

	sfs, err := dirs.CreateStaticFs()
	assert.NoError(err)

	assert.Equal("s1-f1", readFileFromFs(t, sfs, "f1.txt"))
	assert.Equal("s2-f2", readFileFromFs(t, sfs, "f2.txt"))
	assert.Equal("theme-f3", readFileFromFs(t, sfs, "f3.txt"))

}

func TestRemoveDuplicatesKeepRight(t *testing.T) {
	in := []string{"a", "b", "c", "a"}
	out := removeDuplicatesKeepRight(in)

	require.Equal(t, []string{"b", "c", "a"}, out)
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}

func readFileFromFs(t testing.TB, fs afero.Fs, filename string) string {
	filename = filepath.FromSlash(filename)
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		afero.Walk(fs, "", func(path string, info os.FileInfo, err error) error {
			fmt.Println("    ", path, " ", info)
			return nil
		})
		t.Fatalf("Failed to read file: %s", err)
	}
	return string(b)
}
