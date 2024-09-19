// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"strings"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/create"
	"github.com/gohugoio/hugo/create/skeletons"
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
						return newUserError("path needs to be provided")
					}
					h, err := r.Hugo(flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					return create.NewContent(h, contentType, args[0], force)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
						if len(args) != 0 {
							return []string{}, cobra.ShellCompDirectiveNoFileComp
						}
						return []string{}, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveFilterDirs
					}
					cmd.Flags().StringVarP(&contentType, "kind", "k", "", "content type to create")
					cmd.Flags().String("editor", "", "edit new content with this editor, if provided")
					_ = cmd.RegisterFlagCompletionFunc("editor", cobra.NoFileCompletions)
					cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite file if it already exists")
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
						return newUserError("path needs to be provided")
					}
					createpath, err := filepath.Abs(filepath.Clean(args[0]))
					if err != nil {
						return err
					}

					cfg := config.New()
					cfg.Set("workingDir", createpath)
					cfg.Set("publishDir", "public")

					conf, err := r.ConfigFromProvider(configKey{counter: r.configVersionID.Load()}, flagsToCfg(cd, cfg))
					if err != nil {
						return err
					}
					sourceFs := conf.fs.Source

					err = skeletons.CreateSite(createpath, sourceFs, force, format)
					if err != nil {
						return err
					}

					r.Printf("Congratulations! Your new Hugo site was created in %s.\n\n", createpath)
					r.Println(c.newSiteNextStepsText(createpath, format))

					return nil
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
						if len(args) != 0 {
							return []string{}, cobra.ShellCompDirectiveNoFileComp
						}
						return []string{}, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveFilterDirs
					}
					cmd.Flags().BoolVarP(&force, "force", "f", false, "init inside non-empty directory")
					cmd.Flags().StringVar(&format, "format", "toml", "preferred file format (toml, yaml or json)")
					_ = cmd.RegisterFlagCompletionFunc("format", cobra.FixedCompletions([]string{"toml", "yaml", "json"}, cobra.ShellCompDirectiveNoFileComp))
				},
			},
			&simpleCommand{
				name:  "theme",
				use:   "theme [name]",
				short: "Create a new theme (skeleton)",
				long: `Create a new theme (skeleton) called [name] in ./themes.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
according to your needs.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					if len(args) < 1 {
						return newUserError("theme name needs to be provided")
					}
					cfg := config.New()
					cfg.Set("publishDir", "public")

					conf, err := r.ConfigFromProvider(configKey{counter: r.configVersionID.Load()}, flagsToCfg(cd, cfg))
					if err != nil {
						return err
					}
					sourceFs := conf.fs.Source
					createpath := paths.AbsPathify(conf.configs.Base.WorkingDir, filepath.Join(conf.configs.Base.ThemesDir, args[0]))
					r.Println("Creating new theme in", createpath)

					err = skeletons.CreateTheme(createpath, sourceFs)
					if err != nil {
						return err
					}

					return nil
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
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

	cmd.RunE = nil
	return nil
}

func (c *newCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.rootCmd = cd.Root.Command.(*rootCommand)
	return nil
}

func (c *newCommand) newSiteNextStepsText(path string, format string) string {
	format = strings.ToLower(format)
	var nextStepsText bytes.Buffer

	nextStepsText.WriteString(`Just a few more steps...

1. Change the current directory to ` + path + `.
2. Create or install a theme:
   - Create a new theme with the command "hugo new theme <THEMENAME>"
   - Or, install a theme from https://themes.gohugo.io/
3. Edit hugo.` + format + `, setting the "theme" property to the theme name.
4. Create new content with the command "hugo new content `)

	nextStepsText.WriteString(filepath.Join("<SECTIONNAME>", "<FILENAME>.<FORMAT>"))

	nextStepsText.WriteString(`".
5. Start the embedded web server with the command "hugo server --buildDrafts".

See documentation at https://gohugo.io/.`)

	return nextStepsText.String()
}
