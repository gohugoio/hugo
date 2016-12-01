// Copyright 2016 The Hugo Authors. All rights reserved.
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

package create_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/create"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/viper"
)

func TestNewContent(t *testing.T) {
	initViper()

	err := initFs()
	if err != nil {
		t.Fatalf("initialization error: %s", err)
	}

	cases := []struct {
		kind     string
		path     string
		expected []string
	}{
		{"post", "post/sample-1.md", []string{`title = "Post Arch title"`, `test = "test1"`, "date = \"2015-01-12T19:20:04-07:00\""}},
		{"stump", "stump/sample-2.md", []string{`title = "sample 2"`}},     // no archetype file
		{"", "sample-3.md", []string{`title = "sample 3"`}},                // no archetype
		{"product", "product/sample-4.md", []string{`title = "sample 4"`}}, // empty archetype front matter
	}

	for i, c := range cases {
		err = create.NewContent(hugofs.Source(), c.kind, c.path)
		if err != nil {
			t.Errorf("[%d] NewContent: %s", i, err)
		}

		fname := filepath.Join("content", filepath.FromSlash(c.path))
		content := readFileFromFs(t, hugofs.Source(), fname)

		for i, v := range c.expected {
			found := strings.Contains(content, v)
			if !found {
				t.Errorf("[%d] %q missing from output:\n%q", i, v, content)
			}
		}
	}
}

func TestNewContentInitCaps(t *testing.T) {
	initViper()
	viper.Set("coerceTitleFormat", "initCaps")

	err := initFs()
	if err != nil {
		t.Fatalf("initialization error: %s", err)
	}

	cases := []struct {
		kind     string
		path     string
		expected []string
	}{
		{"post", "post/sample-one-1.md", []string{`title = "Post Arch title"`, `test = "test1"`, "date = \"2015-01-12T19:20:04-07:00\""}},
		{"stump", "stump/sample-tWO-2.md", []string{`title = "Sample TWO 2"`}},       // no archetype file
		{"", "sample-three-3.md", []string{`title = "Sample Three 3"`}},              // no archetype
		{"product", "product/sample-four-4.md", []string{`title = "Sample Four 4"`}}, // empty archetype front matter
	}

	for i, c := range cases {
		err = create.NewContent(hugofs.Source(), c.kind, c.path)
		if err != nil {
			t.Errorf("[%d] NewContent: %s", i, err)
		}

		fname := filepath.Join("content", filepath.FromSlash(c.path))
		content := readFileFromFs(t, hugofs.Source(), fname)

		for i, v := range c.expected {
			found := strings.Contains(content, v)
			if !found {
				t.Errorf("[%d] %q missing from output:\n%q", i, v, content)
			}
		}
	}
}

func initViper() {
	viper.Reset()
	viper.Set("metaDataFormat", "toml")
	viper.Set("archetypeDir", "archetypes")
	viper.Set("contentDir", "content")
	viper.Set("themesDir", "themes")
	viper.Set("theme", "sample")
}

func initFs() error {
	hugofs.InitMemFs()
	perm := os.FileMode(0755)
	var err error

	// create directories
	dirs := []string{
		"archetypes",
		"content",
		filepath.Join("themes", "sample", "archetypes"),
	}
	for _, dir := range dirs {
		err = hugofs.Source().Mkdir(dir, perm)
		if err != nil {
			return err
		}
	}

	// create files
	for _, v := range []struct {
		path    string
		content string
	}{
		{
			path:    filepath.Join("archetypes", "post.md"),
			content: "+++\ndate = \"2015-01-12T19:20:04-07:00\"\ntitle = \"Post Arch title\"\ntest = \"test1\"\n+++\n",
		},
		{
			path:    filepath.Join("archetypes", "product.md"),
			content: "+++\n+++\n",
		},
	} {
		f, err := hugofs.Source().Create(v.path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write([]byte(v.content))
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO(bep) extract common testing package with this and some others
func readFileFromFs(t *testing.T, fs afero.Fs, filename string) string {
	filename = filepath.FromSlash(filename)
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		// Print some debug info
		root := strings.Split(filename, helpers.FilePathSeparator)[0]
		afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() {
				fmt.Println("    ", path)
			}

			return nil
		})
		t.Fatalf("Failed to read file: %s", err)
	}
	return string(b)
}
