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
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/modules/npm"
	"github.com/spf13/cobra"
)

const commonUsageMod = `
Note that Hugo will always start out by resolving the components defined in the site
configuration, provided by a _vendor directory (if no --ignoreVendorPaths flag provided),
Go Modules, or a folder inside the themes directory, in that order.

See https://gohugo.io/hugo-modules/ for more information.

`

// buildConfigCommands creates a new config command and its subcommands.
func newModCommands() *modCommands {
	var (
		clean   bool
		pattern string
		all     bool
	)

	npmCommand := &simpleCommand{
		name:  "npm",
		short: "Various npm helpers.",
		long:  `Various npm (Node package manager) helpers.`,
		commands: []simplecobra.Commander{
			&simpleCommand{
				name:  "pack",
				short: "Experimental: Prepares and writes a composite package.json file for your project.",
				long: `Prepares and writes a composite package.json file for your project.

On first run it creates a "package.hugo.json" in the project root if not already there. This file will be used as a template file
with the base dependency set. 

This set will be merged with all "package.hugo.json" files found in the dependency tree, picking the version closest to the project.

This command is marked as 'Experimental'. We think it's a great idea, so it's not likely to be
removed from Hugo, but we need to test this out in "real life" to get a feel of it,
so this may/will change in future versions of Hugo.
`,
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					applyLocalFlagsBuildConfig(cmd, r)
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					h, err := r.Hugo(flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					return npm.Pack(h.BaseFs.ProjectSourceFs, h.BaseFs.AssetsWithDuplicatesPreserved.Fs)
				},
			},
		},
	}

	return &modCommands{
		commands: []simplecobra.Commander{
			&simpleCommand{
				name:  "init",
				short: "Initialize this project as a Hugo Module.",
				long: `Initialize this project as a Hugo Module.
	It will try to guess the module path, but you may help by passing it as an argument, e.g:
	
		hugo mod init github.com/gohugoio/testshortcodes
	
	Note that Hugo Modules supports multi-module projects, so you can initialize a Hugo Module
	inside a subfolder on GitHub, as one example.
	`,
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					applyLocalFlagsBuildConfig(cmd, r)
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					h, err := r.getOrCreateHugo(flagsToCfg(cd, nil), true)
					if err != nil {
						return err
					}
					var initPath string
					if len(args) >= 1 {
						initPath = args[0]
					}
					c := h.Configs.ModulesClient
					if err := c.Init(initPath); err != nil {
						return err
					}
					return nil
				},
			},
			&simpleCommand{
				name:  "verify",
				short: "Verify dependencies.",
				long:  `Verify checks that the dependencies of the current module, which are stored in a local downloaded source cache, have not been modified since being downloaded.`,
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					applyLocalFlagsBuildConfig(cmd, r)
					cmd.Flags().BoolVarP(&clean, "clean", "", false, "delete module cache for dependencies that fail verification")
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					conf, err := r.ConfigFromProvider(configKey{counter: r.configVersionID.Load()}, flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					client := conf.configs.ModulesClient
					return client.Verify(clean)
				},
			},
			&simpleCommand{
				name:  "graph",
				short: "Print a module dependency graph.",
				long: `Print a module dependency graph with information about module status (disabled, vendored).
Note that for vendored modules, that is the version listed and not the one from go.mod.
`,
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					applyLocalFlagsBuildConfig(cmd, r)
					cmd.Flags().BoolVarP(&clean, "clean", "", false, "delete module cache for dependencies that fail verification")
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					conf, err := r.ConfigFromProvider(configKey{counter: r.configVersionID.Load()}, flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					client := conf.configs.ModulesClient
					return client.Graph(os.Stdout)
				},
			},
			&simpleCommand{
				name:  "clean",
				short: "Delete the Hugo Module cache for the current project.",
				long:  `Delete the Hugo Module cache for the current project.`,
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					applyLocalFlagsBuildConfig(cmd, r)
					cmd.Flags().StringVarP(&pattern, "pattern", "", "", `pattern matching module paths to clean (all if not set), e.g. "**hugo*"`)
					_ = cmd.RegisterFlagCompletionFunc("pattern", cobra.NoFileCompletions)
					cmd.Flags().BoolVarP(&all, "all", "", false, "clean entire module cache")
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					h, err := r.Hugo(flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					if all {
						modCache := h.ResourceSpec.FileCaches.ModulesCache()
						count, err := modCache.Prune(true)
						r.Printf("Deleted %d files from module cache.", count)
						return err
					}

					return h.Configs.ModulesClient.Clean(pattern)
				},
			},
			&simpleCommand{
				name:  "tidy",
				short: "Remove unused entries in go.mod and go.sum.",
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					applyLocalFlagsBuildConfig(cmd, r)
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					h, err := r.Hugo(flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					return h.Configs.ModulesClient.Tidy()
				},
			},
			&simpleCommand{
				name:  "vendor",
				short: "Vendor all module dependencies into the _vendor directory.",
				long: `Vendor all module dependencies into the _vendor directory.
	If a module is vendored, that is where Hugo will look for it's dependencies.
	`,
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					applyLocalFlagsBuildConfig(cmd, r)
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					h, err := r.Hugo(flagsToCfg(cd, nil))
					if err != nil {
						return err
					}
					return h.Configs.ModulesClient.Vendor()
				},
			},

			&simpleCommand{
				name:  "get",
				short: "Resolves dependencies in your current Hugo Project.",
				long: `
Resolves dependencies in your current Hugo Project.

Some examples:

Install the latest version possible for a given module:

    hugo mod get github.com/gohugoio/testshortcodes
    
Install a specific version:

    hugo mod get github.com/gohugoio/testshortcodes@v0.3.0

Install the latest versions of all direct module dependencies:

    hugo mod get
    hugo mod get ./... (recursive)

Install the latest versions of all module dependencies (direct and indirect):

    hugo mod get -u
    hugo mod get -u ./... (recursive)

Run "go help get" for more information. All flags available for "go get" is also relevant here.
` + commonUsageMod,
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.DisableFlagParsing = true
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					// We currently just pass on the flags we get to Go and
					// need to do the flag handling manually.
					if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
						return errHelp
					}

					var lastArg string
					if len(args) != 0 {
						lastArg = args[len(args)-1]
					}

					if lastArg == "./..." {
						args = args[:len(args)-1]
						// Do a recursive update.
						dirname, err := os.Getwd()
						if err != nil {
							return err
						}

						// Sanity chesimplecobra. We do recursive walking and want to avoid
						// accidents.
						if len(dirname) < 5 {
							return errors.New("must not be run from the file system root")
						}

						filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
							if info.IsDir() {
								return nil
							}
							if info.Name() == "go.mod" {
								// Found a module.
								dir := filepath.Dir(path)

								cfg := config.New()
								cfg.Set("workingDir", dir)
								conf, err := r.ConfigFromProvider(configKey{counter: r.configVersionID.Add(1)}, flagsToCfg(cd, cfg))
								if err != nil {
									return err
								}
								r.Println("Update module in", conf.configs.Base.WorkingDir)
								client := conf.configs.ModulesClient
								return client.Get(args...)

							}
							return nil
						})
						return nil
					} else {
						conf, err := r.ConfigFromProvider(configKey{counter: r.configVersionID.Load()}, flagsToCfg(cd, nil))
						if err != nil {
							return err
						}
						client := conf.configs.ModulesClient
						return client.Get(args...)
					}
				},
			},
			npmCommand,
		},
	}
}

type modCommands struct {
	r *rootCommand

	commands []simplecobra.Commander
}

func (c *modCommands) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *modCommands) Name() string {
	return "mod"
}

func (c *modCommands) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	_, err := c.r.ConfigFromProvider(configKey{counter: c.r.configVersionID.Load()}, nil)
	if err != nil {
		return err
	}
	// config := conf.configs.Base

	return nil
}

func (c *modCommands) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Various Hugo Modules helpers."
	cmd.Long = `Various helpers to help manage the modules in your project's dependency graph.
Most operations here requires a Go version installed on your system (>= Go 1.12) and the relevant VCS client (typically Git).
This is not needed if you only operate on modules inside /themes or if you have vendored them via "hugo mod vendor".

` + commonUsageMod
	cmd.RunE = nil
	return nil
}

func (c *modCommands) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.r = cd.Root.Command.(*rootCommand)
	return nil
}
