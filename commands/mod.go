// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gohugoio/hugo/modules"
	"github.com/spf13/cobra"
)

var _ cmder = (*modCmd)(nil)

type modCmd struct {
	*baseBuilderCmd
}

func (c *modCmd) newVerifyCmd() *cobra.Command {
	var clean bool

	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify dependencies.",
		Long: `Verify checks that the dependencies of the current module, which are stored in a local downloaded source cache, have not been modified since being downloaded.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.withModsClient(true, func(c *modules.Client) error {
				return c.Verify(clean)
			})
		},
	}

	verifyCmd.Flags().BoolVarP(&clean, "clean", "", false, "delete module cache for dependencies that fail verification")

	return verifyCmd
}

var moduleNotFoundRe = regexp.MustCompile("module.*not found")

func (c *modCmd) newCleanCmd() *cobra.Command {
	var pattern string
	var all bool
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Delete the Hugo Module cache for the current project.",
		Long: `Delete the Hugo Module cache for the current project.

Note that after you run this command, all of your dependencies will be re-downloaded next time you run "hugo".

Also note that if you configure a positive maxAge for the "modules" file cache, it will also be cleaned as part of "hugo --gc".
 
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if all {
				com, err := c.initConfig(false)

				if err != nil && !moduleNotFoundRe.MatchString(err.Error()) {
					return err
				}

				_, err = com.hugo().FileCaches.ModulesCache().Prune(true)
				return err
			}
			return c.withModsClient(true, func(c *modules.Client) error {
				return c.Clean(pattern)
			})
		},
	}

	cmd.Flags().StringVarP(&pattern, "pattern", "", "", `pattern matching module paths to clean (all if not set), e.g. "**hugo*"`)
	cmd.Flags().BoolVarP(&all, "all", "", false, "clean entire module cache")

	return cmd
}

func (b *commandsBuilder) newModCmd() *modCmd {

	c := &modCmd{}

	const commonUsage = `
Note that Hugo will always start out by resolving the components defined in the site
configuration, provided by a _vendor directory (if no --ignoreVendor flag provided),
Go Modules, or a folder inside the themes directory, in that order.

See https://gohugo.io/hugo-modules/ for more information.

`

	cmd := &cobra.Command{
		Use:   "mod",
		Short: "Various Hugo Modules helpers.",
		Long: `Various helpers to help manage the modules in your project's dependency graph.

Most operations here requires a Go version installed on your system (>= Go 1.12) and the relevant VCS client (typically Git).
This is not needed if you only operate on modules inside /themes or if you have vendored them via "hugo mod vendor".

` + commonUsage,

		RunE: nil,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:                "get",
			DisableFlagParsing: true,
			Short:              "Resolves dependencies in your current Hugo Project.",
			Long: `
Resolves dependencies in your current Hugo Project.

Some examples:

Install the latest version possible for a given module:

    hugo mod get github.com/gohugoio/testshortcodes
    
Install a specific version:

    hugo mod get github.com/gohugoio/testshortcodes@v0.3.0

Install the latest versions of all module dependencies:

    hugo mod get -u
    hugo mod get -u ./... (recursive)

Run "go help get" for more information. All flags available for "go get" is also relevant here.
` + commonUsage,
			RunE: func(cmd *cobra.Command, args []string) error {
				// We currently just pass on the flags we get to Go and
				// need to do the flag handling manually.
				if len(args) == 1 && args[0] == "-h" {
					return cmd.Help()
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

					// Sanity check. We do recursive walking and want to avoid
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
							fmt.Println("Update module in", dir)
							c.source = dir
							err := c.withModsClient(false, func(c *modules.Client) error {
								if len(args) == 1 && args[0] == "-h" {
									return cmd.Help()
								}
								return c.Get(args...)
							})
							if err != nil {
								return err
							}

						}

						return nil
					})

					return nil
				}

				return c.withModsClient(false, func(c *modules.Client) error {
					return c.Get(args...)
				})
			},
		},
		&cobra.Command{
			Use:   "graph",
			Short: "Print a module dependency graph.",
			Long: `Print a module dependency graph with information about module status (disabled, vendored).
Note that for vendored modules, that is the version listed and not the one from go.mod.
`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(true, func(c *modules.Client) error {
					return c.Graph(os.Stdout)
				})
			},
		},
		&cobra.Command{
			Use:   "init",
			Short: "Initialize this project as a Hugo Module.",
			Long: `Initialize this project as a Hugo Module.
It will try to guess the module path, but you may help by passing it as an argument, e.g:

    hugo mod init github.com/gohugoio/testshortcodes

Note that Hugo Modules supports multi-module projects, so you can initialize a Hugo Module
inside a subfolder on GitHub, as one example.
`,
			RunE: func(cmd *cobra.Command, args []string) error {
				var path string
				if len(args) >= 1 {
					path = args[0]
				}
				return c.withModsClient(false, func(c *modules.Client) error {
					return c.Init(path)
				})
			},
		},
		&cobra.Command{
			Use:   "vendor",
			Short: "Vendor all module dependencies into the _vendor directory.",
			Long: `Vendor all module dependencies into the _vendor directory.

If a module is vendored, that is where Hugo will look for it's dependencies.
`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(true, func(c *modules.Client) error {
					return c.Vendor()
				})
			},
		},
		c.newVerifyCmd(),
		&cobra.Command{
			Use:   "tidy",
			Short: "Remove unused entries in go.mod and go.sum.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(true, func(c *modules.Client) error {
					return c.Tidy()
				})
			},
		},
		c.newCleanCmd(),
	)

	c.baseBuilderCmd = b.newBuilderCmd(cmd)

	return c

}

func (c *modCmd) withModsClient(failOnMissingConfig bool, f func(*modules.Client) error) error {
	com, err := c.initConfig(failOnMissingConfig)
	if err != nil {
		return err
	}

	return f(com.hugo().ModulesClient)
}

func (c *modCmd) initConfig(failOnNoConfig bool) (*commandeer, error) {
	com, err := initializeConfig(failOnNoConfig, false, &c.hugoBuilderCommon, c, nil)
	if err != nil {
		return nil, err
	}
	return com, nil
}
