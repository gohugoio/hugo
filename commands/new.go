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

package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/create"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	configFormat  string
	contentEditor string
	contentType   string
)

func init() {
	newSiteCmd.Flags().StringVarP(&configFormat, "format", "f", "toml", "config & frontmatter format")
	newSiteCmd.Flags().Bool("force", false, "init inside non-empty directory")
	newCmd.Flags().StringVarP(&contentType, "kind", "k", "", "content type to create")
	newCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "filesystem path to read files relative from")
	newCmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	newCmd.Flags().StringVar(&contentEditor, "editor", "", "edit new content with this editor, if provided")

	newCmd.AddCommand(newSiteCmd)
	newCmd.AddCommand(newThemeCmd)

}

var newCmd = &cobra.Command{
	Use:   "new [path]",
	Short: "Create new content for your site",
	Long: `Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.

You can also specify the kind with ` + "`-k KIND`" + `.

If archetypes are provided in your theme or site, they will be used.`,

	RunE: NewContent,
}

var newSiteCmd = &cobra.Command{
	Use:   "site [path]",
	Short: "Create a new site (skeleton)",
	Long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use ` + "`hugo new [contentPath]`" + ` to create new content.`,
	RunE: NewSite,
}

var newThemeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "Create a new theme",
	Long: `Create a new theme (skeleton) called [name] in the current directory.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
as you see fit.`,
	RunE: NewTheme,
}

// NewContent adds new content to a Hugo site.
func NewContent(cmd *cobra.Command, args []string) error {
	cfg, err := InitializeConfig()

	if err != nil {
		return err
	}

	c, err := newCommandeer(cfg)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("editor") {
		c.Set("newContentEditor", contentEditor)
	}

	if len(args) < 1 {
		return newUserError("path needs to be provided")
	}

	createPath := args[0]

	var kind string

	createPath, kind = newContentPathSection(createPath)

	if contentType != "" {
		kind = contentType
	}

	ps, err := helpers.NewPathSpec(cfg.Fs, cfg.Cfg)
	if err != nil {
		return err
	}

	// If a site isn't in use in the archetype template, we can skip the build.
	siteFactory := func(filename string, siteUsed bool) (*hugolib.Site, error) {
		if !siteUsed {
			return hugolib.NewSite(*cfg)
		}
		var s *hugolib.Site
		if err := c.initSites(); err != nil {
			return nil, err
		}

		if err := Hugo.Build(hugolib.BuildCfg{SkipRender: true, PrintStats: false}); err != nil {
			return nil, err
		}

		s = Hugo.Sites[0]

		if len(Hugo.Sites) > 1 {
			// Find the best match.
			for _, ss := range Hugo.Sites {
				if strings.Contains(createPath, "."+ss.Language.Lang) {
					s = ss
					break
				}
			}
		}
		return s, nil
	}

	return create.NewContent(ps, siteFactory, kind, createPath)
}

func doNewSite(fs *hugofs.Fs, basepath string, force bool) error {
	archeTypePath := filepath.Join(basepath, "archetypes")
	dirs := []string{
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		archeTypePath,
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
		filepath.Join(basepath, "themes"),
	}

	if exists, _ := helpers.Exists(basepath, fs.Source); exists {
		if isDir, _ := helpers.IsDir(basepath, fs.Source); !isDir {
			return errors.New(basepath + " already exists but not a directory")
		}

		isEmpty, _ := helpers.IsEmpty(basepath, fs.Source)

		switch {
		case !isEmpty && !force:
			return errors.New(basepath + " already exists and is not empty")

		case !isEmpty && force:
			all := append(dirs, filepath.Join(basepath, "config."+configFormat))
			for _, path := range all {
				if exists, _ := helpers.Exists(path, fs.Source); exists {
					return errors.New(path + " already exists")
				}
			}
		}
	}

	for _, dir := range dirs {
		if err := fs.Source.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("Failed to create dir: %s", err)
		}
	}

	createConfig(fs, basepath, configFormat)

	// Create a defaul archetype file.
	helpers.SafeWriteToDisk(filepath.Join(archeTypePath, "default.md"),
		strings.NewReader(create.ArchetypeTemplateTemplate), fs.Source)

	jww.FEEDBACK.Printf("Congratulations! Your new Hugo site is created in %s.\n\n", basepath)
	jww.FEEDBACK.Println(nextStepsText())

	return nil
}

func nextStepsText() string {
	var nextStepsText bytes.Buffer

	nextStepsText.WriteString(`Just a few more steps and you're ready to go:

1. Download a theme into the same-named folder.
   Choose a theme from https://themes.gohugo.io/, or
   create your own with the "hugo new theme <THEMENAME>" command.
2. Perhaps you want to add some content. You can add single files
   with "hugo new `)

	nextStepsText.WriteString(filepath.Join("<SECTIONNAME>", "<FILENAME>.<FORMAT>"))

	nextStepsText.WriteString(`".
3. Start the built-in live server via "hugo server".

Visit https://gohugo.io/ for quickstart guide and full documentation.`)

	return nextStepsText.String()
}

// NewSite creates a new Hugo site and initializes a structured Hugo directory.
func NewSite(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return newUserError("path needs to be provided")
	}

	createpath, err := filepath.Abs(filepath.Clean(args[0]))
	if err != nil {
		return newUserError(err)
	}

	forceNew, _ := cmd.Flags().GetBool("force")

	return doNewSite(hugofs.NewDefault(viper.New()), createpath, forceNew)
}

// NewTheme creates a new Hugo theme.
func NewTheme(cmd *cobra.Command, args []string) error {
	cfg, err := InitializeConfig()

	if err != nil {
		return err
	}

	if len(args) < 1 {
		return newUserError("theme name needs to be provided")
	}

	c, err := newCommandeer(cfg)
	if err != nil {
		return err
	}

	createpath := c.PathSpec().AbsPathify(filepath.Join(c.Cfg.GetString("themesDir"), args[0]))
	jww.INFO.Println("creating theme at", createpath)

	if x, _ := helpers.Exists(createpath, cfg.Fs.Source); x {
		return newUserError(createpath, "already exists")
	}

	mkdir(createpath, "layouts", "_default")
	mkdir(createpath, "layouts", "partials")

	touchFile(cfg.Fs.Source, createpath, "layouts", "index.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "404.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "_default", "list.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "_default", "single.html")

	touchFile(cfg.Fs.Source, createpath, "layouts", "partials", "header.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "partials", "footer.html")

	mkdir(createpath, "archetypes")

	archDefault := []byte("+++\n+++\n")

	err = helpers.WriteToDisk(filepath.Join(createpath, "archetypes", "default.md"), bytes.NewReader(archDefault), cfg.Fs.Source)
	if err != nil {
		return err
	}

	mkdir(createpath, "static", "js")
	mkdir(createpath, "static", "css")

	by := []byte(`The MIT License (MIT)

Copyright (c) ` + time.Now().Format("2006") + ` YOUR_NAME_HERE

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`)

	err = helpers.WriteToDisk(filepath.Join(createpath, "LICENSE.md"), bytes.NewReader(by), cfg.Fs.Source)
	if err != nil {
		return err
	}

	createThemeMD(cfg.Fs, createpath)

	return nil
}

func mkdir(x ...string) {
	p := filepath.Join(x...)

	err := os.MkdirAll(p, 0777) // before umask
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}

func touchFile(fs afero.Fs, x ...string) {
	inpath := filepath.Join(x...)
	mkdir(filepath.Dir(inpath))
	err := helpers.WriteToDisk(inpath, bytes.NewReader([]byte{}), fs)
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}

func createThemeMD(fs *hugofs.Fs, inpath string) (err error) {

	by := []byte(`# theme.toml template for a Hugo theme
# See https://github.com/gohugoio/hugoThemes#themetoml for an example

name = "` + strings.Title(helpers.MakeTitle(filepath.Base(inpath))) + `"
license = "MIT"
licenselink = "https://github.com/yourname/yourtheme/blob/master/LICENSE.md"
description = ""
homepage = "http://example.com/"
tags = []
features = []
min_version = "0.31.1"

[author]
  name = ""
  homepage = ""

# If porting an existing theme
[original]
  name = ""
  homepage = ""
  repo = ""
`)

	err = helpers.WriteToDisk(filepath.Join(inpath, "theme.toml"), bytes.NewReader(by), fs.Source)
	if err != nil {
		return
	}

	return nil
}

func newContentPathSection(path string) (string, string) {
	// Forward slashes is used in all examples. Convert if needed.
	// Issue #1133
	createpath := filepath.FromSlash(path)
	var section string
	// assume the first directory is the section (kind)
	if strings.Contains(createpath[1:], helpers.FilePathSeparator) {
		section = helpers.GuessSection(createpath)
	}

	return createpath, section
}

func createConfig(fs *hugofs.Fs, inpath string, kind string) (err error) {
	in := map[string]string{
		"baseURL":      "http://example.org/",
		"title":        "My New Hugo Site",
		"languageCode": "en-us",
	}
	kind = parser.FormatSanitize(kind)

	var buf bytes.Buffer
	err = parser.InterfaceToConfig(in, parser.FormatToLeadRune(kind), &buf)
	if err != nil {
		return err
	}

	return helpers.WriteToDisk(filepath.Join(inpath, "config."+kind), &buf, fs.Source)
}
