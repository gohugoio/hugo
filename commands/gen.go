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
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"gopkg.in/yaml.v2"
)

func newGenCommand() *genCommand {
	var (
		// Flags.
		gendocdir string
		genmandir string

		// Chroma flags.
		style                  string
		highlightStyle         string
		lineNumbersInlineStyle string
		lineNumbersTableStyle  string
	)

	newChromaStyles := func() simplecobra.Commander {
		return &simpleCommand{
			name:  "chromastyles",
			short: "Generate CSS stylesheet for the Chroma code highlighter",
			long: `Generate CSS stylesheet for the Chroma code highlighter for a given style. This stylesheet is needed if markup.highlight.noClasses is disabled in config.

See https://xyproto.github.io/splash/docs/all.html for a preview of the available styles`,

			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				builder := styles.Get(style).Builder()
				if highlightStyle != "" {
					builder.Add(chroma.LineHighlight, highlightStyle)
				}
				if lineNumbersInlineStyle != "" {
					builder.Add(chroma.LineNumbers, lineNumbersInlineStyle)
				}
				if lineNumbersTableStyle != "" {
					builder.Add(chroma.LineNumbersTable, lineNumbersTableStyle)
				}
				style, err := builder.Build()
				if err != nil {
					return err
				}
				formatter := html.New(html.WithAllClasses(true))
				formatter.WriteCSS(os.Stdout, style)
				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVar(&style, "style", "friendly", "highlighter style (see https://xyproto.github.io/splash/docs/)")
				_ = cmd.RegisterFlagCompletionFunc("style", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&highlightStyle, "highlightStyle", "", `foreground and background colors for highlighted lines, e.g. --highlightStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("highlightStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&lineNumbersInlineStyle, "lineNumbersInlineStyle", "", `foreground and background colors for inline line numbers, e.g. --lineNumbersInlineStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("lineNumbersInlineStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&lineNumbersTableStyle, "lineNumbersTableStyle", "", `foreground and background colors for table line numbers, e.g. --lineNumbersTableStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("lineNumbersTableStyle", cobra.NoFileCompletions)
			},
		}
	}

	newMan := func() simplecobra.Commander {
		return &simpleCommand{
			name:  "man",
			short: "Generate man pages for the Hugo CLI",
			long: `This command automatically generates up-to-date man pages of Hugo's
	command-line interface.  By default, it creates the man page files
	in the "man" directory under the current directory.`,

			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				header := &doc.GenManHeader{
					Section: "1",
					Manual:  "Hugo Manual",
					Source:  fmt.Sprintf("Hugo %s", hugo.CurrentVersion),
				}
				if !strings.HasSuffix(genmandir, helpers.FilePathSeparator) {
					genmandir += helpers.FilePathSeparator
				}
				if found, _ := helpers.Exists(genmandir, hugofs.Os); !found {
					r.Println("Directory", genmandir, "does not exist, creating...")
					if err := hugofs.Os.MkdirAll(genmandir, 0o777); err != nil {
						return err
					}
				}
				cd.CobraCommand.Root().DisableAutoGenTag = true

				r.Println("Generating Hugo man pages in", genmandir, "...")
				doc.GenManTree(cd.CobraCommand.Root(), header, genmandir)

				r.Println("Done.")

				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVar(&genmandir, "dir", "man/", "the directory to write the man pages.")
				_ = cmd.MarkFlagDirname("dir")
			},
		}
	}

	newGen := func() simplecobra.Commander {
		const gendocFrontmatterTemplate = `---
title: "%s"
slug: %s
url: %s
---
`

		return &simpleCommand{
			name:  "doc",
			short: "Generate Markdown documentation for the Hugo CLI.",
			long: `Generate Markdown documentation for the Hugo CLI.
			This command is, mostly, used to create up-to-date documentation
	of Hugo's command-line interface for https://gohugo.io/.

	It creates one Markdown file per command with front matter suitable
	for rendering in Hugo.`,
			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				cd.CobraCommand.VisitParents(func(c *cobra.Command) {
					// Disable the "Auto generated by spf13/cobra on DATE"
					// as it creates a lot of diffs.
					c.DisableAutoGenTag = true
				})
				if !strings.HasSuffix(gendocdir, helpers.FilePathSeparator) {
					gendocdir += helpers.FilePathSeparator
				}
				if found, _ := helpers.Exists(gendocdir, hugofs.Os); !found {
					r.Println("Directory", gendocdir, "does not exist, creating...")
					if err := hugofs.Os.MkdirAll(gendocdir, 0o777); err != nil {
						return err
					}
				}
				prepender := func(filename string) string {
					name := filepath.Base(filename)
					base := strings.TrimSuffix(name, path.Ext(name))
					url := "/commands/" + strings.ToLower(base) + "/"
					return fmt.Sprintf(gendocFrontmatterTemplate, strings.Replace(base, "_", " ", -1), base, url)
				}

				linkHandler := func(name string) string {
					base := strings.TrimSuffix(name, path.Ext(name))
					return "/commands/" + strings.ToLower(base) + "/"
				}
				r.Println("Generating Hugo command-line documentation in", gendocdir, "...")
				doc.GenMarkdownTreeCustom(cd.CobraCommand.Root(), gendocdir, prepender, linkHandler)
				r.Println("Done.")

				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVar(&gendocdir, "dir", "/tmp/hugodoc/", "the directory to write the doc.")
				_ = cmd.MarkFlagDirname("dir")
			},
		}
	}

	var docsHelperTarget string

	newDocsHelper := func() simplecobra.Commander {
		return &simpleCommand{
			name:  "docshelper",
			short: "Generate some data files for the Hugo docs.",

			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				r.Println("Generate docs data to", docsHelperTarget)

				var buf bytes.Buffer
				jsonEnc := json.NewEncoder(&buf)

				configProvider := func() docshelper.DocProvider {
					conf := hugolib.DefaultConfig()
					conf.CacheDir = "" // The default value does not make sense in the docs.
					defaultConfig := parser.NullBoolJSONMarshaller{Wrapped: parser.LowerCaseCamelJSONMarshaller{Value: conf}}
					return docshelper.DocProvider{"config": defaultConfig}
				}

				docshelper.AddDocProviderFunc(configProvider)
				if err := jsonEnc.Encode(docshelper.GetDocProvider()); err != nil {
					return err
				}

				// Decode the JSON to a map[string]interface{} and then unmarshal it again to the correct format.
				var m map[string]interface{}
				if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
					return err
				}

				targetFile := filepath.Join(docsHelperTarget, "docs.yaml")

				f, err := os.Create(targetFile)
				if err != nil {
					return err
				}
				defer f.Close()
				yamlEnc := yaml.NewEncoder(f)
				if err := yamlEnc.Encode(m); err != nil {
					return err
				}

				r.Println("Done!")
				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.Hidden = true
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVarP(&docsHelperTarget, "dir", "", "docs/data", "data dir")
			},
		}
	}

	return &genCommand{
		commands: []simplecobra.Commander{
			newChromaStyles(),
			newGen(),
			newMan(),
			newDocsHelper(),
		},
	}
}

type genCommand struct {
	rootCmd *rootCommand

	commands []simplecobra.Commander
}

func (c *genCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *genCommand) Name() string {
	return "gen"
}

func (c *genCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	return nil
}

func (c *genCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "A collection of several useful generators."

	cmd.RunE = nil
	return nil
}

func (c *genCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.rootCmd = cd.Root.Command.(*rootCommand)
	return nil
}
