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

	"github.com/gohugoio/hugo/deps"

	"github.com/gohugoio/hugo/hugolib"

	"fmt"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/create"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewContent(t *testing.T) {
	assert := require.New(t)

	cases := []struct {
		kind     string
		path     string
		expected []string
	}{
		{"post", "post/sample-1.md", []string{`title = "Post Arch title"`, `test = "test1"`, "date = \"2015-01-12T19:20:04-07:00\""}},
		{"post", "post/org-1.org", []string{`#+title: ORG-1`}},
		{"emptydate", "post/sample-ed.md", []string{`title = "Empty Date Arch title"`, `test = "test1"`}},
		{"stump", "stump/sample-2.md", []string{`title: "Sample 2"`}},      // no archetype file
		{"", "sample-3.md", []string{`title: "Sample 3"`}},                 // no archetype
		{"product", "product/sample-4.md", []string{`title = "SAMPLE-4"`}}, // empty archetype front matter
		{"lang", "post/lang-1.md", []string{`Site Lang: en|Name: Lang 1|i18n: Hugo Rocks!`}},
		{"lang", "post/lang-2.en.md", []string{`Site Lang: en|Name: Lang 2|i18n: Hugo Rocks!`}},
		{"lang", "post/lang-3.nn.md", []string{`Site Lang: nn|Name: Lang 3|i18n: Hugo Rokkar!`}},
		{"lang", "content_nn/post/lang-4.md", []string{`Site Lang: nn|Name: Lang 4|i18n: Hugo Rokkar!`}},
		{"lang", "content_nn/post/lang-5.en.md", []string{`Site Lang: en|Name: Lang 5|i18n: Hugo Rocks!`}},
		{"lang", "post/my-bundle/index.md", []string{`Site Lang: en|Name: My Bundle|i18n: Hugo Rocks!`}},
		{"lang", "post/my-bundle/index.en.md", []string{`Site Lang: en|Name: My Bundle|i18n: Hugo Rocks!`}},
		{"lang", "post/my-bundle/index.nn.md", []string{`Site Lang: nn|Name: My Bundle|i18n: Hugo Rokkar!`}},
		{"shortcodes", "shortcodes/go.md", []string{
			`title = "GO"`,
			"{{< myshortcode >}}",
			"{{% myshortcode %}}",
			"{{</* comment */>}}\n{{%/* comment */%}}"}}, // shortcodes
	}

	for i, c := range cases {
		cfg, fs := newTestCfg(assert)
		assert.NoError(initFs(fs))
		h, err := hugolib.NewHugoSites(deps.DepsCfg{Cfg: cfg, Fs: fs})
		assert.NoError(err)

		assert.NoError(create.NewContent(h, c.kind, c.path))

		fname := filepath.FromSlash(c.path)
		if !strings.HasPrefix(fname, "content") {
			fname = filepath.Join("content", fname)
		}
		content := readFileFromFs(t, fs.Source, fname)
		for _, v := range c.expected {
			found := strings.Contains(content, v)
			if !found {
				t.Fatalf("[%d] %q missing from output:\n%q", i, v, content)
			}
		}
	}
}

func TestNewContentFromDir(t *testing.T) {
	assert := require.New(t)
	cfg, fs := newTestCfg(assert)
	assert.NoError(initFs(fs))

	archetypeDir := filepath.Join("archetypes", "my-bundle")
	assert.NoError(fs.Source.Mkdir(archetypeDir, 0755))

	archetypeThemeDir := filepath.Join("themes", "mytheme", "archetypes", "my-theme-bundle")
	assert.NoError(fs.Source.Mkdir(archetypeThemeDir, 0755))

	contentFile := `
File: %s
Site Lang: {{ .Site.Language.Lang  }} 	
Name: {{ replace .Name "-" " " | title }}
i18n: {{ T "hugo" }}
`

	assert.NoError(afero.WriteFile(fs.Source, filepath.Join(archetypeDir, "index.md"), []byte(fmt.Sprintf(contentFile, "index.md")), 0755))
	assert.NoError(afero.WriteFile(fs.Source, filepath.Join(archetypeDir, "index.nn.md"), []byte(fmt.Sprintf(contentFile, "index.nn.md")), 0755))

	assert.NoError(afero.WriteFile(fs.Source, filepath.Join(archetypeDir, "pages", "bio.md"), []byte(fmt.Sprintf(contentFile, "bio.md")), 0755))
	assert.NoError(afero.WriteFile(fs.Source, filepath.Join(archetypeDir, "resources", "hugo1.json"), []byte(`hugo1: {{ printf "no template handling in here" }}`), 0755))
	assert.NoError(afero.WriteFile(fs.Source, filepath.Join(archetypeDir, "resources", "hugo2.xml"), []byte(`hugo2: {{ printf "no template handling in here" }}`), 0755))

	assert.NoError(afero.WriteFile(fs.Source, filepath.Join(archetypeThemeDir, "index.md"), []byte(fmt.Sprintf(contentFile, "index.md")), 0755))
	assert.NoError(afero.WriteFile(fs.Source, filepath.Join(archetypeThemeDir, "resources", "hugo1.json"), []byte(`hugo1: {{ printf "no template handling in here" }}`), 0755))

	h, err := hugolib.NewHugoSites(deps.DepsCfg{Cfg: cfg, Fs: fs})
	assert.NoError(err)
	assert.Equal(2, len(h.Sites))

	assert.NoError(create.NewContent(h, "my-bundle", "post/my-post"))

	assertContains(assert, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-post/resources/hugo1.json")), `hugo1: {{ printf "no template handling in here" }}`)
	assertContains(assert, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-post/resources/hugo2.xml")), `hugo2: {{ printf "no template handling in here" }}`)

	// Content files should get the correct site context.
	// TODO(bep) archetype check i18n
	assertContains(assert, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-post/index.md")), `File: index.md`, `Site Lang: en`, `Name: My Post`, `i18n: Hugo Rocks!`)
	assertContains(assert, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-post/index.nn.md")), `File: index.nn.md`, `Site Lang: nn`, `Name: My Post`, `i18n: Hugo Rokkar!`)

	assertContains(assert, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-post/pages/bio.md")), `File: bio.md`, `Site Lang: en`, `Name: My Post`)

	assert.NoError(create.NewContent(h, "my-theme-bundle", "post/my-theme-post"))
	assertContains(assert, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-theme-post/index.md")), `File: index.md`, `Site Lang: en`, `Name: My Theme Post`, `i18n: Hugo Rocks!`)
	assertContains(assert, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-theme-post/resources/hugo1.json")), `hugo1: {{ printf "no template handling in here" }}`)

}

func initFs(fs *hugofs.Fs) error {
	perm := os.FileMode(0755)
	var err error

	// create directories
	dirs := []string{
		"archetypes",
		"content",
		filepath.Join("themes", "sample", "archetypes"),
	}
	for _, dir := range dirs {
		err = fs.Source.Mkdir(dir, perm)
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
			path:    filepath.Join("archetypes", "post.org"),
			content: "#+title: {{ .BaseFileName  | upper }}",
		},
		{
			path: filepath.Join("archetypes", "product.md"),
			content: `+++
title = "{{ .BaseFileName  | upper }}"
+++`,
		},
		{
			path:    filepath.Join("archetypes", "emptydate.md"),
			content: "+++\ndate =\"\"\ntitle = \"Empty Date Arch title\"\ntest = \"test1\"\n+++\n",
		},
		{
			path:    filepath.Join("archetypes", "lang.md"),
			content: `Site Lang: {{ .Site.Language.Lang  }}|Name: {{ replace .Name "-" " " | title }}|i18n: {{ T "hugo" }}`,
		},
		// #3623x
		{
			path: filepath.Join("archetypes", "shortcodes.md"),
			content: `+++
title = "{{ .BaseFileName  | upper }}"
+++

{{< myshortcode >}}

Some text.

{{% myshortcode %}}
{{</* comment */>}}
{{%/* comment */%}}


`,
		},
	} {
		f, err := fs.Source.Create(v.path)
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

func assertContains(assert *require.Assertions, v interface{}, matches ...string) {
	for _, m := range matches {
		assert.Contains(v, m)
	}
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

func newTestCfg(assert *require.Assertions) (*viper.Viper, *hugofs.Fs) {

	cfg := `

theme = "mytheme"
	
[languages]
[languages.en]
weight = 1
languageName = "English"
[languages.nn]
weight = 2
languageName = "Nynorsk"
contentDir = "content_nn"

`

	mm := afero.NewMemMapFs()

	assert.NoError(afero.WriteFile(mm, filepath.Join("i18n", "en.toml"), []byte(`[hugo]
other = "Hugo Rocks!"`), 0755))
	assert.NoError(afero.WriteFile(mm, filepath.Join("i18n", "nn.toml"), []byte(`[hugo]
other = "Hugo Rokkar!"`), 0755))

	assert.NoError(afero.WriteFile(mm, "config.toml", []byte(cfg), 0755))

	v, _, err := hugolib.LoadConfig(hugolib.ConfigSourceDescriptor{Fs: mm, Filename: "config.toml"})
	assert.NoError(err)

	return v, hugofs.NewFrom(mm, v)

}
