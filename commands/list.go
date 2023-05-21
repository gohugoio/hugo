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
	"context"
	"encoding/csv"
	"time"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

// newListCommand creates a new list command and its subcommands.
func newListCommand() *listCommand {

	list := func(cd *simplecobra.Commandeer, r *rootCommand, createRecord func(page.Page) []string, opts ...any) error {
		bcfg := hugolib.BuildCfg{SkipRender: true}
		cfg := config.New()
		for i := 0; i < len(opts); i += 2 {
			cfg.Set(opts[i].(string), opts[i+1])
		}
		h, err := r.Build(cd, bcfg, cfg)
		if err != nil {
			return err
		}

		writer := csv.NewWriter(r.Out)
		defer writer.Flush()

		for _, p := range h.Pages() {
			if record := createRecord(p); record != nil {
				if err := writer.Write(record); err != nil {
					return err
				}
				if err != nil {
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
				short: "List all drafts",
				long:  `List all of the drafts in your content directory.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					createRecord := func(p page.Page) []string {
						if !p.Draft() || p.File().IsZero() {
							return nil
						}
						return []string{
							p.File().Path(),
							p.PublishDate().Format(time.RFC3339)}

					}
					return list(cd, r, createRecord,
						"buildDrafts", true,
						"buildFuture", true,
						"buildExpired", true,
					)
				},
			},
			&simpleCommand{
				name:  "future",
				short: "List all posts dated in the future",
				long:  `List all of the posts in your content directory which will be posted in the future.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					createRecord := func(p page.Page) []string {
						if !resource.IsFuture(p) || p.File().IsZero() {
							return nil
						}
						return []string{
							p.File().Path(),
							p.PublishDate().Format(time.RFC3339),
						}

					}
					return list(cd, r, createRecord,
						"buildFuture", true,
						"buildDrafts", true,
					)
				},
			},
			&simpleCommand{
				name:  "expired",
				short: "List all posts already expired",
				long:  `List all of the posts in your content directory which has already expired.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					createRecord := func(p page.Page) []string {
						if !resource.IsExpired(p) || p.File().IsZero() {
							return nil
						}
						return []string{
							p.File().Path(),
							p.PublishDate().Format(time.RFC3339),
						}

					}
					return list(cd, r, createRecord,
						"buildExpired", true,
						"buildDrafts", true,
					)
				},
			},
			&simpleCommand{
				name:  "all",
				short: "List all posts",
				long:  `List all of the posts in your content directory, include drafts, future and expired pages.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					createRecord := func(p page.Page) []string {
						if p.File().IsZero() {
							return nil
						}
						return []string{
							p.File().Path(),
							p.PublishDate().Format(time.RFC3339),
						}

					}
					return list(cd, r, createRecord, "buildDrafts", true, "buildFuture", true, "buildExpired", true)
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

	return nil
}

func (c *listCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	return nil
}
