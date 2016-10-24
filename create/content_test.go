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
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/create"
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
		kind          string
		path          string
		resultStrings []string
	}{
		{"post", "post/sample-1.md", []string{`title = "sample 1"`, `test = "test1"`}},
		{"stump", "stump/sample-2.md", []string{`title = "sample 2"`}},     // no archetype file
		{"", "sample-3.md", []string{`title = "sample 3"`}},                // no archetype
		{"product", "product/sample-4.md", []string{`title = "sample 4"`}}, // empty archetype front matter
	}

	for i, c := range cases {
		err = create.NewContent(hugofs.Source(), c.kind, c.path)
		if err != nil {
			t.Errorf("[%d] NewContent: %s", i, err)
		}

		fname := filepath.Join(os.TempDir(), "content", filepath.FromSlash(c.path))
		_, err = hugofs.Source().Stat(fname)
		if err != nil {
			t.Errorf("[%d] Stat: %s", i, err)
		}

		for _, v := range c.resultStrings {
			found, err := afero.FileContainsBytes(hugofs.Source(), fname, []byte(v))
			if err != nil {
				t.Errorf("[%d] FileContainsBytes: %s", i, err)
			}
			if !found {
				t.Errorf("content missing from output: %q", v)
			}
		}
	}
}

func initViper() {
	viper.Reset()
	viper.Set("metaDataFormat", "toml")
	viper.Set("archetypeDir", filepath.Join(os.TempDir(), "archetypes"))
	viper.Set("contentDir", filepath.Join(os.TempDir(), "content"))
	viper.Set("themesDir", filepath.Join(os.TempDir(), "themes"))
	viper.Set("theme", "sample")
}

func initFs() error {
	hugofs.SetSource(new(afero.MemMapFs))
	perm := os.FileMode(0755)
	var err error

	// create directories
	dirs := []string{
		"archetypes",
		"content",
		filepath.Join("themes", "sample", "archetypes"),
	}
	for _, dir := range dirs {
		dir = filepath.Join(os.TempDir(), dir)
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
			path:    filepath.Join(os.TempDir(), "archetypes", "post.md"),
			content: "+++\ndate = \"2015-01-12T19:20:04-07:00\"\ntitle = \"post arch\"\ntest = \"test1\"\n+++\n",
		},
		{
			path:    filepath.Join(os.TempDir(), "archetypes", "product.md"),
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
