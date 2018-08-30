// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"fmt"
	"time"

	src "github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/hugolib"

	"path/filepath"

	"github.com/gohugoio/hugo/parser"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	_ cmder = (*convertCmd)(nil)
)

type convertCmd struct {
	hugoBuilderCommon

	outputDir string
	unsafe    bool

	*baseCmd
}

func newConvertCmd() *convertCmd {
	cc := &convertCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "convert",
		Short: "Convert your content to different formats",
		Long: `Convert your content (e.g. front matter) to different formats.

See convert's subcommands toJSON, toTOML and toYAML for more information.`,
		RunE: nil,
	})

	cc.cmd.AddCommand(
		&cobra.Command{
			Use:   "toJSON",
			Short: "Convert front matter to JSON",
			Long: `toJSON converts all front matter in the content directory
to use JSON for the front matter.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return cc.convertContents(rune([]byte(parser.JSONLead)[0]))
			},
		},
		&cobra.Command{
			Use:   "toTOML",
			Short: "Convert front matter to TOML",
			Long: `toTOML converts all front matter in the content directory
to use TOML for the front matter.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return cc.convertContents(rune([]byte(parser.TOMLLead)[0]))
			},
		},
		&cobra.Command{
			Use:   "toYAML",
			Short: "Convert front matter to YAML",
			Long: `toYAML converts all front matter in the content directory
to use YAML for the front matter.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return cc.convertContents(rune([]byte(parser.YAMLLead)[0]))
			},
		},
	)

	cc.cmd.PersistentFlags().StringVarP(&cc.outputDir, "output", "o", "", "filesystem path to write files to")
	cc.cmd.PersistentFlags().StringVarP(&cc.source, "source", "s", "", "filesystem path to read files relative from")
	cc.cmd.PersistentFlags().BoolVar(&cc.unsafe, "unsafe", false, "enable less safe operations, please backup first")
	cc.cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})

	return cc
}

func (cc *convertCmd) convertContents(mark rune) error {
	if cc.outputDir == "" && !cc.unsafe {
		return newUserError("Unsafe operation not allowed, use --unsafe or set a different output path")
	}

	c, err := initializeConfig(true, false, &cc.hugoBuilderCommon, cc, nil)
	if err != nil {
		return err
	}

	h, err := hugolib.NewHugoSites(*c.DepsCfg)
	if err != nil {
		return err
	}

	if err := h.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
		return err
	}

	site := h.Sites[0]

	site.Log.FEEDBACK.Println("processing", len(site.AllPages), "content files")
	for _, p := range site.AllPages {
		if err := cc.convertAndSavePage(p, site, mark); err != nil {
			return err
		}
	}
	return nil
}

func (cc *convertCmd) convertAndSavePage(p *hugolib.Page, site *hugolib.Site, mark rune) error {
	// The resources are not in .Site.AllPages.
	for _, r := range p.Resources.ByType("page") {
		if err := cc.convertAndSavePage(r.(*hugolib.Page), site, mark); err != nil {
			return err
		}
	}

	if p.Filename() == "" {
		// No content file.
		return nil
	}

	site.Log.INFO.Println("Attempting to convert", p.LogicalName())
	newPage, err := site.NewPage(p.LogicalName())
	if err != nil {
		return err
	}

	f, _ := p.File.(src.ReadableFile)
	file, err := f.Open()
	if err != nil {
		site.Log.ERROR.Println("Error reading file:", p.Path())
		file.Close()
		return nil
	}

	psr, err := parser.ReadFrom(file)
	if err != nil {
		site.Log.ERROR.Println("Error processing file:", p.Path())
		file.Close()
		return err
	}

	file.Close()

	metadata, err := psr.Metadata()
	if err != nil {
		site.Log.ERROR.Println("Error processing file:", p.Path())
		return err
	}

	// better handling of dates in formats that don't have support for them
	if mark == parser.FormatToLeadRune("json") || mark == parser.FormatToLeadRune("yaml") || mark == parser.FormatToLeadRune("toml") {
		newMetadata := cast.ToStringMap(metadata)
		for k, v := range newMetadata {
			switch vv := v.(type) {
			case time.Time:
				newMetadata[k] = vv.Format(time.RFC3339)
			}
		}
		metadata = newMetadata
	}

	newPage.SetSourceContent(psr.Content())
	if err = newPage.SetSourceMetaData(metadata, mark); err != nil {
		site.Log.ERROR.Printf("Failed to set source metadata for file %q: %s. For more info see For more info see https://github.com/gohugoio/hugo/issues/2458", newPage.FullFilePath(), err)
		return nil
	}

	newFilename := p.Filename()
	if cc.outputDir != "" {
		newFilename = filepath.Join(cc.outputDir, p.Dir(), newPage.LogicalName())
	}

	if err = newPage.SaveSourceAs(newFilename); err != nil {
		return fmt.Errorf("Failed to save file %q: %s", newFilename, err)
	}

	return nil
}
