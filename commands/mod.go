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

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/modules"
	"github.com/spf13/cobra"
)

var _ cmder = (*modCmd)(nil)

type modCmd struct {
	hugoBuilderCommon
	*baseCmd
}

func newModCmd() *modCmd {
	c := &modCmd{}

	c.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "mod",
		Short: "Various Hugo Modules helpers.",
		RunE:  nil,
	})

	c.cmd.AddCommand(
		&cobra.Command{
			// go get [-d] [-m] [-u] [-v] [-insecure] [build flags] [packages]
			Use:   "get",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) >= 1 {
					return c.getModulesHandler(nil).Get(args[0])
				}

				// Collect any modules defined in config.toml
				_, err := c.initConfig()
				return err

			},
		},
		&cobra.Command{
			Use:   "graph",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.getModulesHandler(nil).Graph()
			},
		},
		&cobra.Command{
			Use:   "init",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				var path string
				if len(args) >= 1 {
					path = args[0]
				}
				return c.getModulesHandler(nil).Init(path)
			},
		},
		&cobra.Command{
			Use:   "vendor",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				com, err := c.initConfig()
				if err != nil {
					return err
				}
				return c.getModulesHandler(com.Cfg).Vendor()
			},
		},
		&cobra.Command{
			Use:   "tidy",
			Short: "TODO(bep)",
			RunE: func(cmd *cobra.Command, args []string) error {
				com, err := c.initConfig()
				if err != nil {
					return err
				}
				return c.getModulesHandler(com.Cfg).Tidy()
			},
		},
	)

	c.handleCommonBuilderFlags(c.cmd)

	return c

}

func (c *modCmd) initConfig() (*commandeer, error) {
	com, err := initializeConfig(true, false, &c.hugoBuilderCommon, c, nil)
	if err != nil {
		return nil, err
	}
	return com, nil
}

func (c *modCmd) getModulesHandler(cfg config.Provider) *modules.Handler {
	var (
		workingDir string
		themesDir  string
		themes     []string
	)

	if c.source != "" {
		workingDir = c.source
	} else {
		workingDir, _ = os.Getwd()
	}

	if cfg != nil {
		// TODO(bep) mod remember this if we change
		themesDir = cfg.GetString("themesDir")
		themes = cfg.GetStringSlice("theme")
	}

	fs := hugofs.Os
	return modules.New(fs, workingDir, themesDir, themes)
}
