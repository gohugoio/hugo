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

package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/create"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newNewCommand() *newCommand {
	var (
		force       bool
		contentType string
		format      string
	)

	var c *newCommand
	c = &newCommand{
		commands: []simplecobra.Commander{
			&simpleCommand{
				name:  "content",
				use:   "content [path]",
				short: "Create new content for your site",
				long: `Create a new content file and automatically set the date and title.
		It will guess which kind of file to create based on the path provided.
		
		You can also specify the kind with ` + "`-k KIND`" + `.
		
		If archetypes are provided in your theme or site, they will be used.
		
		Ensure you run this within the root directory of your site.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					if len(args) < 1 {
						return errors.New("path needs to be provided")
					}
					h, err := r.Hugo(flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					return create.NewContent(h, contentType, args[0], force)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.Flags().StringVarP(&contentType, "kind", "k", "", "content type to create")
					cmd.Flags().String("editor", "", "edit new content with this editor, if provided")
					cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite file if it already exists")
					cmd.Flags().StringVar(&format, "format", "toml", "preferred file format (toml, yaml or json)")
					applyLocalFlagsBuildConfig(cmd, r)

				},
			},
			&simpleCommand{
				name:  "site",
				use:   "site [path]",
				short: "Create a new site (skeleton)",
				long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use ` + "`hugo new [contentPath]`" + ` to create new content.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					if len(args) < 1 {
						return errors.New("path needs to be provided")
					}
					createpath, err := filepath.Abs(filepath.Clean(args[0]))
					if err != nil {
						return err
					}

					cfg := config.New()
					cfg.Set("workingDir", createpath)
					cfg.Set("publishDir", "public")

					conf, err := r.ConfigFromProvider(r.configVersionID.Load(), flagsToCfg(cd, cfg))
					if err != nil {
						return err
					}
					sourceFs := conf.fs.Source

					archeTypePath := filepath.Join(createpath, "archetypes")
					dirs := []string{
						archeTypePath,
						filepath.Join(createpath, "assets"),
						filepath.Join(createpath, "content"),
						filepath.Join(createpath, "data"),
						filepath.Join(createpath, "layouts"),
						filepath.Join(createpath, "static"),
						filepath.Join(createpath, "themes"),
					}

					if exists, _ := helpers.Exists(createpath, sourceFs); exists {
						if isDir, _ := helpers.IsDir(createpath, sourceFs); !isDir {
							return errors.New(createpath + " already exists but not a directory")
						}

						isEmpty, _ := helpers.IsEmpty(createpath, sourceFs)

						switch {
						case !isEmpty && !force:
							return errors.New(createpath + " already exists and is not empty. See --force.")

						case !isEmpty && force:
							all := append(dirs, filepath.Join(createpath, "hugo."+format))
							for _, path := range all {
								if exists, _ := helpers.Exists(path, sourceFs); exists {
									return errors.New(path + " already exists")
								}
							}
						}
					}

					for _, dir := range dirs {
						if err := sourceFs.MkdirAll(dir, 0777); err != nil {
							return fmt.Errorf("failed to create dir: %w", err)
						}
					}

					c.newSiteCreateConfig(sourceFs, createpath, format)

					// Create a default archetype file.
					helpers.SafeWriteToDisk(filepath.Join(archeTypePath, "default.md"),
						strings.NewReader(create.DefaultArchetypeTemplateTemplate), sourceFs)

					r.Printf("Congratulations! Your new Hugo site is created in %s.\n\n", createpath)
					r.Println(c.newSiteNextStepsText())

					return nil
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.Flags().BoolVarP(&force, "force", "f", false, "init inside non-empty directory")
				},
			},
			&simpleCommand{
				name:  "theme",
				use:   "theme [path]",
				short: "Create a new site (skeleton)",
				long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use ` + "`hugo new [contentPath]`" + ` to create new content.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					h, err := r.Hugo(flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					ps := h.PathSpec
					sourceFs := ps.Fs.Source
					themesDir := h.Configs.LoadingInfo.BaseConfig.ThemesDir
					createpath := ps.AbsPathify(filepath.Join(themesDir, args[0]))
					r.Println("Creating theme at", createpath)

					if x, _ := helpers.Exists(createpath, sourceFs); x {
						return errors.New(createpath + " already exists")
					}

					for _, filename := range []string{
						"index.html",
						"404.html",
						"_default/list.html",
						"_default/single.html",
						"partials/head.html",
						"partials/header.html",
						"partials/footer.html",
					} {
						touchFile(sourceFs, filepath.Join(createpath, "layouts", filename))
					}

					baseofDefault := []byte(`<!DOCTYPE html>
<html>
    {{- partial "head.html" . -}}
    <body>
        {{- partial "header.html" . -}}
        <div id="content">
        {{- block "main" . }}{{- end }}
        </div>
        {{- partial "footer.html" . -}}
    </body>
</html>
`)

					err = helpers.WriteToDisk(filepath.Join(createpath, "layouts", "_default", "baseof.html"), bytes.NewReader(baseofDefault), sourceFs)
					if err != nil {
						return err
					}

					mkdir(createpath, "archetypes")

					archDefault := []byte("+++\n+++\n")

					err = helpers.WriteToDisk(filepath.Join(createpath, "archetypes", "default.md"), bytes.NewReader(archDefault), sourceFs)
					if err != nil {
						return err
					}

					mkdir(createpath, "static", "js")
					mkdir(createpath, "static", "css")

					by := []byte(`The MIT License (MIT)

Copyright (c) ` + htime.Now().Format("2006") + ` YOUR_NAME_HERE

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

					err = helpers.WriteToDisk(filepath.Join(createpath, "LICENSE"), bytes.NewReader(by), sourceFs)
					if err != nil {
						return err
					}

					c.createThemeMD(ps.Fs.Source, createpath)

					return nil
				},
			},
		},
	}

	return c

}

type newCommand struct {
	rootCmd *rootCommand

	commands []simplecobra.Commander
}

func (c *newCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *newCommand) Name() string {
	return "new"
}

func (c *newCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	return nil
}

func (c *newCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Create new content for your site"
	cmd.Long = `Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.

You can also specify the kind with ` + "`-k KIND`" + `.

If archetypes are provided in your theme or site, they will be used.

Ensure you run this within the root directory of your site.`
	return nil
}

func (c *newCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.rootCmd = cd.Root.Command.(*rootCommand)
	return nil
}

func (c *newCommand) newSiteCreateConfig(fs afero.Fs, inpath string, kind string) (err error) {
	in := map[string]string{
		"baseURL":      "http://example.org/",
		"title":        "My New Hugo Site",
		"languageCode": "en-us",
	}

	var buf bytes.Buffer
	err = parser.InterfaceToConfig(in, metadecoders.FormatFromString(kind), &buf)
	if err != nil {
		return err
	}

	return helpers.WriteToDisk(filepath.Join(inpath, "hugo."+kind), &buf, fs)
}

func (c *newCommand) newSiteNextStepsText() string {
	var nextStepsText bytes.Buffer

	nextStepsText.WriteString(`Just a few more steps and you're ready to go:

1. Download a theme into the same-named folder.
   Choose a theme from https://themes.gohugo.io/ or
   create your own with the "hugo new theme <THEMENAME>" command.
2. Perhaps you want to add some content. You can add single files
   with "hugo new `)

	nextStepsText.WriteString(filepath.Join("<SECTIONNAME>", "<FILENAME>.<FORMAT>"))

	nextStepsText.WriteString(`".
3. Start the built-in live server via "hugo server".

Visit https://gohugo.io/ for quickstart guide and full documentation.`)

	return nextStepsText.String()
}

func (c *newCommand) createThemeMD(fs afero.Fs, inpath string) (err error) {

	by := []byte(`# theme.toml template for a Hugo theme
# See https://github.com/gohugoio/hugoThemes#themetoml for an example

name = "` + strings.Title(helpers.MakeTitle(filepath.Base(inpath))) + `"
license = "MIT"
licenselink = "https://github.com/yourname/yourtheme/blob/master/LICENSE"
description = ""
homepage = "http://example.com/"
tags = []
features = []
min_version = "0.112.5"

[author]
  name = ""
  homepage = ""

# If porting an existing theme
[original]
  name = ""
  homepage = ""
  repo = ""
`)

	err = helpers.WriteToDisk(filepath.Join(inpath, "theme.toml"), bytes.NewReader(by), fs)
	if err != nil {
		return
	}

	err = helpers.WriteToDisk(filepath.Join(inpath, "hugo.toml"), strings.NewReader("# Theme config.\n"), fs)
	if err != nil {
		return
	}

	return nil
}
