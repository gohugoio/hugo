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
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/spf13/cobra"
)

func newConvertCommand() *convertCommand {
	var c *convertCommand
	c = &convertCommand{
		commands: []simplecobra.Commander{
			&simpleCommand{
				name:  "toJSON",
				short: "Convert front matter to JSON",
				long: `toJSON converts all front matter in the content directory
to use JSON for the front matter.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					return c.convertContents(metadecoders.JSON)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
			&simpleCommand{
				name:  "toTOML",
				short: "Convert front matter to TOML",
				long: `toTOML converts all front matter in the content directory
to use TOML for the front matter.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					return c.convertContents(metadecoders.TOML)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
			&simpleCommand{
				name:  "toYAML",
				short: "Convert front matter to YAML",
				long: `toYAML converts all front matter in the content directory
to use YAML for the front matter.`,
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					return c.convertContents(metadecoders.YAML)
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
				},
			},
		},
	}
	return c
}

type convertCommand struct {
	// Flags.
	outputDir string
	unsafe    bool

	// Deps.
	r *rootCommand
	h *hugolib.HugoSites

	// Commands.
	commands []simplecobra.Commander
}

func (c *convertCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *convertCommand) Name() string {
	return "convert"
}

func (c *convertCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	return nil
}

func (c *convertCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Convert your content to different formats"
	cmd.Long = `Convert your content (e.g. front matter) to different formats.

See convert's subcommands toJSON, toTOML and toYAML for more information.`

	cmd.PersistentFlags().StringVarP(&c.outputDir, "output", "o", "", "filesystem path to write files to")
	_ = cmd.MarkFlagDirname("output")
	cmd.PersistentFlags().BoolVar(&c.unsafe, "unsafe", false, "enable less safe operations, please backup first")

	cmd.RunE = nil
	return nil
}

func (c *convertCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.r = cd.Root.Command.(*rootCommand)
	cfg := config.New()
	cfg.Set("buildDrafts", true)
	h, err := c.r.Hugo(flagsToCfg(cd, cfg))
	if err != nil {
		return err
	}
	c.h = h
	return nil
}

func (c *convertCommand) convertAndSavePage(p page.Page, site *hugolib.Site, targetFormat metadecoders.Format) error {
	// The resources are not in .Site.AllPages.
	for _, r := range p.Resources().ByType("page") {
		if err := c.convertAndSavePage(r.(page.Page), site, targetFormat); err != nil {
			return err
		}
	}

	if p.File() == nil {
		// No content file.
		return nil
	}

	errMsg := fmt.Errorf("error processing file %q", p.File().Path())

	site.Log.Infoln("attempting to convert", p.File().Filename())

	f := p.File()
	file, err := f.FileInfo().Meta().Open()
	if err != nil {
		site.Log.Errorln(errMsg)
		file.Close()
		return nil
	}

	pf, err := pageparser.ParseFrontMatterAndContent(file)
	if err != nil {
		site.Log.Errorln(errMsg)
		file.Close()
		return err
	}

	file.Close()

	// better handling of dates in formats that don't have support for them
	if pf.FrontMatterFormat == metadecoders.JSON || pf.FrontMatterFormat == metadecoders.YAML || pf.FrontMatterFormat == metadecoders.TOML {
		for k, v := range pf.FrontMatter {
			switch vv := v.(type) {
			case time.Time:
				pf.FrontMatter[k] = vv.Format(time.RFC3339)
			}
		}
	}

	var newContent bytes.Buffer
	err = parser.InterfaceToFrontMatter(pf.FrontMatter, targetFormat, &newContent)
	if err != nil {
		site.Log.Errorln(errMsg)
		return err
	}

	newContent.Write(pf.Content)

	newFilename := p.File().Filename()

	if c.outputDir != "" {
		contentDir := strings.TrimSuffix(newFilename, p.File().Path())
		contentDir = filepath.Base(contentDir)

		newFilename = filepath.Join(c.outputDir, contentDir, p.File().Path())
	}

	fs := hugofs.Os
	if err := helpers.WriteToDisk(newFilename, &newContent, fs); err != nil {
		return fmt.Errorf("failed to save file %q:: %w", newFilename, err)
	}

	return nil
}

func (c *convertCommand) convertContents(format metadecoders.Format) error {
	if c.outputDir == "" && !c.unsafe {
		return newUserError("Unsafe operation not allowed, use --unsafe or set a different output path")
	}

	if err := c.h.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
		return err
	}

	site := c.h.Sites[0]

	var pagesBackedByFile page.Pages
	for _, p := range site.AllPages() {
		if p.File() == nil {
			continue
		}
		pagesBackedByFile = append(pagesBackedByFile, p)
	}

	site.Log.Println("processing", len(pagesBackedByFile), "content files")
	for _, p := range site.AllPages() {
		if err := c.convertAndSavePage(p, site, format); err != nil {
			return err
		}
	}
	return nil
}
