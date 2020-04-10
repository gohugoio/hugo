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
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/gohugoio/hugo/parser/pageparser"

	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/hugolib"

	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	_ cmder = (*convertCmd)(nil)
)

type convertCmd struct {
	outputDir string
	unsafe    bool

	*baseBuilderCmd
}

func (b *commandsBuilder) newConvertCmd() *convertCmd {
	cc := &convertCmd{}

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert your content to different formats",
		Long: `Convert your content (e.g. front matter) to different formats.

See convert's subcommands toJSON, toTOML and toYAML for more information.`,
		RunE: nil,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "toJSON",
			Short: "Convert front matter to JSON",
			Long: `toJSON converts all front matter in the content directory
to use JSON for the front matter.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return cc.convertContents(metadecoders.JSON)
			},
		},
		&cobra.Command{
			Use:   "toTOML",
			Short: "Convert front matter to TOML",
			Long: `toTOML converts all front matter in the content directory
to use TOML for the front matter.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return cc.convertContents(metadecoders.TOML)
			},
		},
		&cobra.Command{
			Use:   "toYAML",
			Short: "Convert front matter to YAML",
			Long: `toYAML converts all front matter in the content directory
to use YAML for the front matter.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return cc.convertContents(metadecoders.YAML)
			},
		},
	)

	cmd.PersistentFlags().StringVarP(&cc.outputDir, "output", "o", "", "filesystem path to write files to")
	cmd.PersistentFlags().BoolVar(&cc.unsafe, "unsafe", false, "enable less safe operations, please backup first")

	cc.baseBuilderCmd = b.newBuilderBasicCmd(cmd)

	return cc
}

func (cc *convertCmd) convertContents(format metadecoders.Format) error {
	if cc.outputDir == "" && !cc.unsafe {
		return newUserError("Unsafe operation not allowed, use --unsafe or set a different output path")
	}

	c, err := initializeConfig(true, false, &cc.hugoBuilderCommon, cc, nil)
	if err != nil {
		return err
	}

	c.Cfg.Set("buildDrafts", true)

	h, err := hugolib.NewHugoSites(*c.DepsCfg)
	if err != nil {
		return err
	}

	if err := h.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
		return err
	}

	site := h.Sites[0]

	site.Log.FEEDBACK.Println("processing", len(site.AllPages()), "content files")
	for _, p := range site.AllPages() {
		if err := cc.convertAndSavePage(p, site, format); err != nil {
			return err
		}
	}
	return nil
}

func (cc *convertCmd) convertAndSavePage(p page.Page, site *hugolib.Site, targetFormat metadecoders.Format) error {
	// The resources are not in .Site.AllPages.
	for _, r := range p.Resources().ByType("page") {
		if err := cc.convertAndSavePage(r.(page.Page), site, targetFormat); err != nil {
			return err
		}
	}

	if p.File().IsZero() {
		// No content file.
		return nil
	}

	errMsg := fmt.Errorf("Error processing file %q", p.Path())

	site.Log.INFO.Println("Attempting to convert", p.File().Filename())

	f := p.File()
	file, err := f.FileInfo().Meta().Open()
	if err != nil {
		site.Log.ERROR.Println(errMsg)
		file.Close()
		return nil
	}

	pf, err := pageparser.ParseFrontMatterAndContent(file)
	if err != nil {
		site.Log.ERROR.Println(errMsg)
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
		site.Log.ERROR.Println(errMsg)
		return err
	}

	newContent.Write(pf.Content)

	newFilename := p.File().Filename()

	if cc.outputDir != "" {
		contentDir := strings.TrimSuffix(newFilename, p.Path())
		contentDir = filepath.Base(contentDir)

		newFilename = filepath.Join(cc.outputDir, contentDir, p.Path())
	}

	fs := hugofs.Os
	if err := helpers.WriteToDisk(newFilename, &newContent, fs); err != nil {
		return errors.Wrapf(err, "Failed to save file %q:", newFilename)
	}

	return nil
}

type parsedFile struct {
	frontMatterFormat metadecoders.Format
	frontMatterSource []byte
	frontMatter       map[string]interface{}

	// Everything after Front Matter
	content []byte
}
