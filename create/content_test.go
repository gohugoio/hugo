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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/config/testconfig"

	"github.com/gohugoio/hugo/deps"

	"github.com/gohugoio/hugo/hugolib"

	"github.com/gohugoio/hugo/hugofs"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/create"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
)

// TODO(bep) clean this up. Export the test site builder in Hugolib or something.
func TestNewContentFromFile(t *testing.T) {
	cases := []struct {
		name     string
		kind     string
		path     string
		expected any
	}{
		{"Post", "post", "post/sample-1.md", []string{`title = "Post Arch title"`, `test = "test1"`, "date = \"2015-01-12T19:20:04-07:00\""}},
		{"Post org-mode", "post", "post/org-1.org", []string{`#+title: ORG-1`}},
		{"Post, unknown content filetype", "post", "post/sample-1.pdoc", false},
		{"Empty date", "emptydate", "post/sample-ed.md", []string{`title = "Empty Date Arch title"`, `test = "test1"`}},
		{"Archetype file not found", "stump", "stump/sample-2.md", []string{`title: "Sample 2"`}}, // no archetype file
		{"No archetype", "", "sample-3.md", []string{`title: "Sample 3"`}},                        // no archetype
		{"Empty archetype", "product", "product/sample-4.md", []string{`title = "SAMPLE-4"`}},     // empty archetype front matter
		{"Filenames", "filenames", "content/mypage/index.md", []string{"title = \"INDEX\"\n+++\n\n\nContentBaseName: mypage"}},
		{"Branch Name", "name", "content/tags/tag-a/_index.md", []string{"+++\ntitle = 'Tag A'\n+++"}},

		{"Lang 1", "lang", "post/lang-1.md", []string{`Site Lang: en|Name: Lang 1|i18n: Hugo Rocks!`}},
		{"Lang 2", "lang", "post/lang-2.en.md", []string{`Site Lang: en|Name: Lang 2|i18n: Hugo Rocks!`}},
		{"Lang nn file", "lang", "content/post/lang-3.nn.md", []string{`Site Lang: nn|Name: Lang 3|i18n: Hugo Rokkar!`}},
		{"Lang nn dir", "lang", "content_nn/post/lang-4.md", []string{`Site Lang: nn|Name: Lang 4|i18n: Hugo Rokkar!`}},
		{"Lang en in nn dir", "lang", "content_nn/post/lang-5.en.md", []string{`Site Lang: en|Name: Lang 5|i18n: Hugo Rocks!`}},
		{"Lang en default", "lang", "post/my-bundle/index.md", []string{`Site Lang: en|Name: My Bundle|i18n: Hugo Rocks!`}},
		{"Lang en file", "lang", "post/my-bundle/index.en.md", []string{`Site Lang: en|Name: My Bundle|i18n: Hugo Rocks!`}},
		{"Lang nn bundle", "lang", "content/post/my-bundle/index.nn.md", []string{`Site Lang: nn|Name: My Bundle|i18n: Hugo Rokkar!`}},
		{"Site", "site", "content/mypage/index.md", []string{"RegularPages .Site: 10", "RegularPages site: 10"}},
		{"Shortcodes", "shortcodes", "shortcodes/go.md", []string{
			`title = "GO"`,
			"{{< myshortcode >}}",
			"{{% myshortcode %}}",
			"{{</* comment */>}}\n{{%/* comment */%}}",
		}}, // shortcodes
	}

	c := qt.New(t)

	for i, cas := range cases {
		cas := cas

		c.Run(cas.name, func(c *qt.C) {
			c.Parallel()

			mm := afero.NewMemMapFs()
			c.Assert(initFs(mm), qt.IsNil)
			cfg, fs := newTestCfg(c, mm)
			conf := testconfig.GetTestConfigs(fs.Source, cfg)
			h, err := hugolib.NewHugoSites(deps.DepsCfg{Configs: conf, Fs: fs})
			c.Assert(err, qt.IsNil)
			err = create.NewContent(h, cas.kind, cas.path, false)

			if b, ok := cas.expected.(bool); ok && !b {
				if !b {
					c.Assert(err, qt.Not(qt.IsNil))
				}
				return
			}

			c.Assert(err, qt.IsNil)

			fname := filepath.FromSlash(cas.path)
			if !strings.HasPrefix(fname, "content") {
				fname = filepath.Join("content", fname)
			}

			content := readFileFromFs(c, fs.Source, fname)

			for _, v := range cas.expected.([]string) {
				found := strings.Contains(content, v)
				if !found {
					c.Fatalf("[%d] %q missing from output:\n%q", i, v, content)
				}
			}
		})

	}
}

func TestNewContentFromDirSiteFunction(t *testing.T) {
	mm := afero.NewMemMapFs()
	c := qt.New(t)

	archetypeDir := filepath.Join("archetypes", "my-bundle")
	defaultArchetypeDir := filepath.Join("archetypes", "default")
	c.Assert(mm.MkdirAll(archetypeDir, 0o755), qt.IsNil)
	c.Assert(mm.MkdirAll(defaultArchetypeDir, 0o755), qt.IsNil)

	contentFile := `
File: %s
site RegularPages: {{ len site.RegularPages  }} 	

`

	c.Assert(afero.WriteFile(mm, filepath.Join(archetypeDir, "index.md"), []byte(fmt.Sprintf(contentFile, "index.md")), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(mm, filepath.Join(defaultArchetypeDir, "index.md"), []byte("default archetype index.md"), 0o755), qt.IsNil)

	c.Assert(initFs(mm), qt.IsNil)
	cfg, fs := newTestCfg(c, mm)

	conf := testconfig.GetTestConfigs(fs.Source, cfg)
	h, err := hugolib.NewHugoSites(deps.DepsCfg{Configs: conf, Fs: fs})
	c.Assert(err, qt.IsNil)
	c.Assert(len(h.Sites), qt.Equals, 2)

	c.Assert(create.NewContent(h, "my-bundle", "post/my-post", false), qt.IsNil)
	cContains(c, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-post/index.md")), `site RegularPages: 10`)

	// Default bundle archetype
	c.Assert(create.NewContent(h, "", "post/my-post2", false), qt.IsNil)
	cContains(c, readFileFromFs(t, fs.Source, filepath.Join("content", "post/my-post2/index.md")), `default archetype index.md`)

	// Regular file with bundle kind.
	c.Assert(create.NewContent(h, "my-bundle", "post/foo.md", false), qt.IsNil)
	cContains(c, readFileFromFs(t, fs.Source, filepath.Join("content", "post/foo.md")), `draft: true`)

	// Regular files should fall back to the default archetype (we have no regular file archetype).
	c.Assert(create.NewContent(h, "my-bundle", "mypage.md", false), qt.IsNil)
	cContains(c, readFileFromFs(t, fs.Source, filepath.Join("content", "mypage.md")), `draft: true`)
}

func initFs(fs afero.Fs) error {
	perm := os.FileMode(0o755)
	var err error

	// create directories
	dirs := []string{
		"archetypes",
		"content",
		filepath.Join("themes", "sample", "archetypes"),
	}
	for _, dir := range dirs {
		err = fs.Mkdir(dir, perm)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	// create some dummy content
	for i := 1; i <= 10; i++ {
		filename := filepath.Join("content", fmt.Sprintf("page%d.md", i))
		afero.WriteFile(fs, filename, []byte(`---
title: Test
---
`), 0o666)
	}

	// create archetype files
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
			path: filepath.Join("archetypes", "name.md"),
			content: `+++
title = '{{ replace .Name "-" " " | title }}'
+++`,
		},
		{
			path: filepath.Join("archetypes", "product.md"),
			content: `+++
title = "{{ .BaseFileName  | upper }}"
+++`,
		},
		{
			path: filepath.Join("archetypes", "filenames.md"),
			content: `...
title = "{{ .BaseFileName  | upper }}"
+++


ContentBaseName: {{ .File.ContentBaseName }}

`,
		},
		{
			path: filepath.Join("archetypes", "site.md"),
			content: `...
title = "{{ .BaseFileName  | upper }}"
+++

Len RegularPages .Site: {{ len .Site.RegularPages }}
Len RegularPages site: {{ len site.RegularPages }}


`,
		},
		{
			path:    filepath.Join("archetypes", "emptydate.md"),
			content: "+++\ndate =\"\"\ntitle = \"Empty Date Arch title\"\ntest = \"test1\"\n+++\n",
		},
		{
			path:    filepath.Join("archetypes", "lang.md"),
			content: `Site Lang: {{ site.Language.Lang  }}|Name: {{ replace .Name "-" " " | title }}|i18n: {{ T "hugo" }}`,
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
		f, err := fs.Create(v.path)
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

func cContains(c *qt.C, v any, matches ...string) {
	for _, m := range matches {
		c.Assert(v, qt.Contains, m)
	}
}

// TODO(bep) extract common testing package with this and some others
func readFileFromFs(t testing.TB, fs afero.Fs, filename string) string {
	t.Helper()
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

func newTestCfg(c *qt.C, mm afero.Fs) (config.Provider, *hugofs.Fs) {
	cfg := `

theme = "mytheme"
[languages]
[languages.en]
weight = 1
languageName = "English"
[languages.nn]
weight = 2
languageName = "Nynorsk"

[module]
[[module.mounts]]
  source = 'archetypes'
  target = 'archetypes'
[[module.mounts]]
  source = 'content'
  target = 'content'
  lang = 'en'
[[module.mounts]]
  source = 'content_nn'
  target = 'content'
  lang = 'nn'
`
	if mm == nil {
		mm = afero.NewMemMapFs()
	}

	mm.MkdirAll(filepath.FromSlash("content_nn"), 0o777)

	mm.MkdirAll(filepath.FromSlash("themes/mytheme"), 0o777)

	c.Assert(afero.WriteFile(mm, filepath.Join("i18n", "en.toml"), []byte(`[hugo]
other = "Hugo Rocks!"`), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(mm, filepath.Join("i18n", "nn.toml"), []byte(`[hugo]
other = "Hugo Rokkar!"`), 0o755), qt.IsNil)

	c.Assert(afero.WriteFile(mm, "config.toml", []byte(cfg), 0o755), qt.IsNil)

	res, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{Fs: mm, Filename: "config.toml"})
	c.Assert(err, qt.IsNil)

	return res.LoadingInfo.Cfg, hugofs.NewFrom(mm, res.LoadingInfo.BaseConfig)
}
