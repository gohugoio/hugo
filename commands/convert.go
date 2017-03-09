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
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
)

var outputDir string
var unsafe bool

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert your content to different formats",
	Long: `Convert your content (e.g. front matter) to different formats.

See convert's subcommands toJSON, toTOML and toYAML for more information.`,
	RunE: nil,
}

var toJSONCmd = &cobra.Command{
	Use:   "toJSON",
	Short: "Convert front matter to JSON",
	Long: `toJSON converts all front matter in the content directory
to use JSON for the front matter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return convertContents(rune([]byte(parser.JSONLead)[0]))
	},
}

var toTOMLCmd = &cobra.Command{
	Use:   "toTOML",
	Short: "Convert front matter to TOML",
	Long: `toTOML converts all front matter in the content directory
to use TOML for the front matter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return convertContents(rune([]byte(parser.TOMLLead)[0]))
	},
}

var toYAMLCmd = &cobra.Command{
	Use:   "toYAML",
	Short: "Convert front matter to YAML",
	Long: `toYAML converts all front matter in the content directory
to use YAML for the front matter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return convertContents(rune([]byte(parser.YAMLLead)[0]))
	},
}

func init() {
	convertCmd.AddCommand(toJSONCmd)
	convertCmd.AddCommand(toTOMLCmd)
	convertCmd.AddCommand(toYAMLCmd)
	convertCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "filesystem path to write files to")
	convertCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "filesystem path to read files relative from")
	convertCmd.PersistentFlags().BoolVar(&unsafe, "unsafe", false, "enable less safe operations, please backup first")
	convertCmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
}

func convertContents(mark rune) error {
	cfg, err := InitializeConfig()
	if err != nil {
		return err
	}

	h, err := hugolib.NewHugoSites(*cfg)
	if err != nil {
		return err
	}

	site := h.Sites[0]

	if err = site.Initialise(); err != nil {
		return err
	}

	if site.Source == nil {
		panic("site.Source not set")
	}
	if len(site.Source.Files()) < 1 {
		return errors.New("No source files found")
	}

	contentDir := site.PathSpec.AbsPathify(site.Cfg.GetString("contentDir"))
	site.Log.FEEDBACK.Println("processing", len(site.Source.Files()), "content files")
	for _, file := range site.Source.Files() {
		site.Log.INFO.Println("Attempting to convert", file.LogicalName())
		page, err := site.NewPage(file.LogicalName())
		if err != nil {
			return err
		}

		psr, err := parser.ReadFrom(file.Contents)
		if err != nil {
			site.Log.ERROR.Println("Error processing file:", file.Path())
			return err
		}
		metadata, err := psr.Metadata()
		if err != nil {
			site.Log.ERROR.Println("Error processing file:", file.Path())
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

		page.SetDir(filepath.Join(contentDir, file.Dir()))
		page.SetSourceContent(psr.Content())
		if err = page.SetSourceMetaData(metadata, mark); err != nil {
			site.Log.ERROR.Printf("Failed to set source metadata for file %q: %s. For more info see For more info see https://github.com/spf13/hugo/issues/2458", page.FullFilePath(), err)
			continue
		}

		if outputDir != "" {
			if err = page.SaveSourceAs(filepath.Join(outputDir, page.FullFilePath())); err != nil {
				return fmt.Errorf("Failed to save file %q: %s", page.FullFilePath(), err)
			}
		} else {
			if unsafe {
				if err = page.SaveSource(); err != nil {
					return fmt.Errorf("Failed to save file %q: %s", page.FullFilePath(), err)
				}
			} else {
				site.Log.FEEDBACK.Println("Unsafe operation not allowed, use --unsafe or set a different output path")
			}
		}
	}
	return nil
}
