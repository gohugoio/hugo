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
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/cobra"
)

// newListCommand creates a new list command and its subcommands.
func newListCommand() *listCommand {
	createRecord := func(workingDir string, p page.Page) []string {
		return []string{
			filepath.ToSlash(strings.TrimPrefix(p.File().Filename(), workingDir+string(os.PathSeparator))),
			p.Slug(),
			p.Title(),
			p.Date().Format(time.RFC3339),
			p.ExpiryDate().Format(time.RFC3339),
			p.PublishDate().Format(time.RFC3339),
			strconv.FormatBool(p.Draft()),
			p.Permalink(),
			p.Kind(),
			p.Section(),
		}
	}

	list := func(cd *simplecobra.Commandeer, r *rootCommand, shouldInclude func(page.Page) bool, opts ...any) error {
		bcfg := hugolib.BuildCfg{SkipRender: true}
		cfg := flagsToCfg(cd, nil)
		for i := 0; i < len(opts); i += 2 {
			cfg.Set(opts[i].(string), opts[i+1])
		}
		h, err := r.Build(cd, bcfg, cfg)
		if err != nil {
			return err
		}

		writer := csv.NewWriter(r.Out)
		defer writer.Flush()

		writer.Write([]string{
			"path",
			"slug",
			"title",
			"date",
			"expiryDate",
			"publishDate",
			"draft",
			"permalink",
			"kind",
			"section",
		})

		for _, p := range h.Pages() {
			if shouldInclude(p) {
				record := createRecord(h.Conf.BaseConfig().WorkingDir, p)
				if err := writer.Write(record); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return &listCommand{
		commands: []simplecobra.Commander{
			&simpleCommand{
				name:  "drafts",
				short: "List draft content",
				long:  `List draft content.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					shouldInclude := func(p page.Page) bool {
						if !p.Draft() || p.File() == nil {
							return false
						}
						return true
					}
					return list(cd, r, shouldInclude,
						"buildDrafts", true,
						"buildFuture", true,
						"buildExpired", true,
					)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
			&simpleCommand{
				name:  "future",
				short: "List future content",
				long:  `List content with a future publication date.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					shouldInclude := func(p page.Page) bool {
						if !resource.IsFuture(p) || p.File() == nil {
							return false
						}
						return true
					}
					return list(cd, r, shouldInclude,
						"buildFuture", true,
						"buildDrafts", true,
					)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
			&simpleCommand{
				name:  "expired",
				short: "List expired content",
				long:  `List content with a past expiration date.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					shouldInclude := func(p page.Page) bool {
						if !resource.IsExpired(p) || p.File() == nil {
							return false
						}
						return true
					}
					return list(cd, r, shouldInclude,
						"buildExpired", true,
						"buildDrafts", true,
					)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
			&simpleCommand{
				name:  "all",
				short: "List all content",
				long:  `List all content including draft, future, and expired.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					shouldInclude := func(p page.Page) bool {
						return p.File() != nil
					}
					return list(cd, r, shouldInclude, "buildDrafts", true, "buildFuture", true, "buildExpired", true)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
			&simpleCommand{
				name:  "published",
				short: "List published content",
				long:  `List content that is not draft, future, or expired.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					shouldInclude := func(p page.Page) bool {
						return !p.Draft() && !resource.IsFuture(p) && !resource.IsExpired(p) && p.File() != nil
					}
					return list(cd, r, shouldInclude)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
		},
	}
}

type listCommand struct {
	commands []simplecobra.Commander
}

func (c *listCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *listCommand) Name() string {
	return "list"
}

func (c *listCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	// Do nothing.
	return nil
}

func (c *listCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Listing out various types of content"
	cmd.Long = `Listing out various types of content.

List requires a subcommand, e.g. hugo list drafts`

	cmd.RunE = nil
	return nil
}

func (c *listCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	return nil
}
