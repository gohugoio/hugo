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
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*listCmd)(nil)

type listCmd struct {
	*baseBuilderCmd
}

func (lc *listCmd) buildSites(config map[string]any) (*hugolib.HugoSites, error) {
	cfgInit := func(c *commandeer) error {
		for key, value := range config {
			c.Set(key, value)
		}
		return nil
	}

	c, err := initializeConfig(true, true, false, &lc.hugoBuilderCommon, lc, cfgInit)
	if err != nil {
		return nil, err
	}

	sites, err := hugolib.NewHugoSites(*c.DepsCfg)
	if err != nil {
		return nil, newSystemError("Error creating sites", err)
	}

	if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
		return nil, newSystemError("Error Processing Source Content", err)
	}

	return sites, nil
}

func (b *commandsBuilder) newListCmd() *listCmd {
	cc := &listCmd{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listing out various types of content",
		Long: `Listing out various types of content.

List requires a subcommand, e.g. ` + "`hugo list drafts`.",
		RunE: nil,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "drafts",
			Short: "List all drafts",
			Long:  `List all of the drafts in your content directory.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				sites, err := cc.buildSites(map[string]any{"buildDrafts": true})
				if err != nil {
					return newSystemError("Error building sites", err)
				}

				for _, p := range sites.Pages() {
					if p.Draft() {
						jww.FEEDBACK.Println(strings.TrimPrefix(p.File().Filename(), sites.WorkingDir+string(os.PathSeparator)))
					}
				}

				return nil
			},
		},
		&cobra.Command{
			Use:   "future",
			Short: "List all posts dated in the future",
			Long:  `List all of the posts in your content directory which will be posted in the future.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				sites, err := cc.buildSites(map[string]any{"buildFuture": true})
				if err != nil {
					return newSystemError("Error building sites", err)
				}

				if err != nil {
					return newSystemError("Error building sites", err)
				}

				writer := csv.NewWriter(os.Stdout)
				defer writer.Flush()

				for _, p := range sites.Pages() {
					if resource.IsFuture(p) {
						err := writer.Write([]string{
							strings.TrimPrefix(p.File().Filename(), sites.WorkingDir+string(os.PathSeparator)),
							p.PublishDate().Format(time.RFC3339),
						})
						if err != nil {
							return newSystemError("Error writing future posts to stdout", err)
						}
					}
				}

				return nil
			},
		},
		&cobra.Command{
			Use:   "expired",
			Short: "List all posts already expired",
			Long:  `List all of the posts in your content directory which has already expired.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				sites, err := cc.buildSites(map[string]any{"buildExpired": true})
				if err != nil {
					return newSystemError("Error building sites", err)
				}

				if err != nil {
					return newSystemError("Error building sites", err)
				}

				writer := csv.NewWriter(os.Stdout)
				defer writer.Flush()

				for _, p := range sites.Pages() {
					if resource.IsExpired(p) {
						err := writer.Write([]string{
							strings.TrimPrefix(p.File().Filename(), sites.WorkingDir+string(os.PathSeparator)),
							p.ExpiryDate().Format(time.RFC3339),
						})
						if err != nil {
							return newSystemError("Error writing expired posts to stdout", err)
						}
					}
				}

				return nil
			},
		},
		&cobra.Command{
			Use:   "all",
			Short: "List all posts",
			Long:  `List all of the posts in your content directory, include drafts, future and expired pages.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				sites, err := cc.buildSites(map[string]any{
					"buildExpired": true,
					"buildDrafts":  true,
					"buildFuture":  true,
				})
				if err != nil {
					return newSystemError("Error building sites", err)
				}

				writer := csv.NewWriter(os.Stdout)
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
				})
				for _, p := range sites.Pages() {
					if !p.IsPage() {
						continue
					}
					err := writer.Write([]string{
						strings.TrimPrefix(p.File().Filename(), sites.WorkingDir+string(os.PathSeparator)),
						p.Slug(),
						p.Title(),
						p.Date().Format(time.RFC3339),
						p.ExpiryDate().Format(time.RFC3339),
						p.PublishDate().Format(time.RFC3339),
						strconv.FormatBool(p.Draft()),
						p.Permalink(),
					})
					if err != nil {
						return newSystemError("Error writing posts to stdout", err)
					}
				}

				return nil
			},
		},
	)

	cc.baseBuilderCmd = b.newBuilderBasicCmd(cmd)

	return cc
}
