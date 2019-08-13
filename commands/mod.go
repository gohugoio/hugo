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
	"os"

	"github.com/gohugoio/hugo/modules"
	"github.com/spf13/cobra"
)

var _ cmder = (*modCmd)(nil)

type modCmd struct {
	*baseBuilderCmd
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

Run "go help get" for more information. All flags available for "go get" is also relevant here.
` + commonUsage,
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(false, func(c *modules.Client) error {
					// We currently just pass on the flags we get to Go and
					// need to do the flag handling manually.
					if len(args) == 1 && args[0] == "-h" {
						return cmd.Help()
					}
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
		&cobra.Command{
			Use:   "tidy",
			Short: "Remove unused entries in go.mod and go.sum.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.withModsClient(true, func(c *modules.Client) error {
					return c.Tidy()
				})
			},
		},
		&cobra.Command{
			Use:   "clean",
			Short: "Delete the entire Hugo Module cache.",
			Long: `Delete the entire Hugo Module cache.

Note that after you run this command, all of your dependencies will be re-downloaded next time you run "hugo".

Also note that if you configure a positive maxAge for the "modules" file cache, it will also be cleaned as part of "hugo --gc".
 
`,
			RunE: func(cmd *cobra.Command, args []string) error {
				com, err := c.initConfig(true)
				if err != nil {
					return err
				}

				_, err = com.hugo.FileCaches.ModulesCache().Prune(true)
				return err

			},
		},
	)

	c.baseBuilderCmd = b.newBuilderCmd(cmd)

	return c

}

func (c *modCmd) withModsClient(failOnMissingConfig bool, f func(*modules.Client) error) error {
	com, err := c.initConfig(failOnMissingConfig)
	if err != nil {
		return err
	}

	return f(com.hugo.ModulesClient)
}

func (c *modCmd) initConfig(failOnNoConfig bool) (*commandeer, error) {
	com, err := initializeConfig(failOnNoConfig, false, &c.hugoBuilderCommon, c, nil)
	if err != nil {
		return nil, err
	}
	return com, nil
}
