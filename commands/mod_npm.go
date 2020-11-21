// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/modules/npm"
	"github.com/spf13/cobra"
)

func newModNPMCmd(c *modCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "npm",
		Short: "Various npm helpers.",
		Long:  `Various npm (Node package manager) helpers.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.withHugo(func(h *hugolib.HugoSites) error {
				return nil
			})
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "pack",
		Short: "Experimental: Prepares and writes a composite package.json file for your project.",
		Long: `Prepares and writes a composite package.json file for your project.

On first run it creates a "package.hugo.json" in the project root if not already there. This file will be used as a template file
with the base dependency set. 

This set will be merged with all "package.hugo.json" files found in the dependency tree, picking the version closest to the project.

This command is marked as 'Experimental'. We think it's a great idea, so it's not likely to be
removed from Hugo, but we need to test this out in "real life" to get a feel of it,
so this may/will change in future versions of Hugo.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.withHugo(func(h *hugolib.HugoSites) error {
				return npm.Pack(h.BaseFs.SourceFs, h.BaseFs.Assets.Dirs)
			})
		},
	})

	return cmd
}
