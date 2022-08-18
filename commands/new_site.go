// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/create"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/parser"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*newSiteCmd)(nil)

type newSiteCmd struct {
	configFormat string

	*baseBuilderCmd
}

func (b *commandsBuilder) newNewSiteCmd() *newSiteCmd {
	cc := &newSiteCmd{}

	cmd := &cobra.Command{
		Use:   "site [path]",
		Short: "Create a new site (skeleton)",
		Long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use ` + "`hugo new [contentPath]`" + ` to create new content.`,
		RunE: cc.newSite,
	}

	cmd.Flags().StringVarP(&cc.configFormat, "format", "f", "toml", "config file format")
	cmd.Flags().Bool("force", false, "init inside non-empty directory")

	cc.baseBuilderCmd = b.newBuilderBasicCmd(cmd)

	return cc
}

func (n *newSiteCmd) doNewSite(fs *hugofs.Fs, basepath string, force bool) error {
	archeTypePath := filepath.Join(basepath, "archetypes")
	dirs := []string{
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		archeTypePath,
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
		filepath.Join(basepath, "themes"),
	}

	if exists, _ := helpers.Exists(basepath, fs.Source); exists {
		if isDir, _ := helpers.IsDir(basepath, fs.Source); !isDir {
			return errors.New(basepath + " already exists but not a directory")
		}

		isEmpty, _ := helpers.IsEmpty(basepath, fs.Source)

		switch {
		case !isEmpty && !force:
			return errors.New(basepath + " already exists and is not empty. See --force.")

		case !isEmpty && force:
			all := append(dirs, filepath.Join(basepath, "config."+n.configFormat))
			for _, path := range all {
				if exists, _ := helpers.Exists(path, fs.Source); exists {
					return errors.New(path + " already exists")
				}
			}
		}
	}

	for _, dir := range dirs {
		if err := fs.Source.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("Failed to create dir: %w", err)
		}
	}

	createConfig(fs, basepath, n.configFormat)

	// Create a default archetype file.
	helpers.SafeWriteToDisk(filepath.Join(archeTypePath, "default.md"),
		strings.NewReader(create.DefaultArchetypeTemplateTemplate), fs.Source)

	jww.FEEDBACK.Printf("Congratulations! Your new Hugo site is created in %s.\n\n", basepath)
	jww.FEEDBACK.Println(nextStepsText())

	return nil
}

// newSite creates a new Hugo site and initializes a structured Hugo directory.
func (n *newSiteCmd) newSite(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return newUserError("path needs to be provided")
	}

	createpath, err := filepath.Abs(filepath.Clean(args[0]))
	if err != nil {
		return newUserError(err)
	}

	forceNew, _ := cmd.Flags().GetBool("force")
	cfg := config.New()
	cfg.Set("workingDir", createpath)
	cfg.Set("publishDir", "public")
	return n.doNewSite(hugofs.NewDefault(cfg), createpath, forceNew)
}

func createConfig(fs *hugofs.Fs, inpath string, kind string) (err error) {
	in := map[string]string{
		"baseURL":      "http://example.org/",
		"title":        "My New Hugo Site",
		"languageCode": "en-us",
	}

	var buf bytes.Buffer
	err = parser.InterfaceToConfig(in, metadecoders.FormatFromString(kind), &buf)
	if err != nil {
		return err
	}

	return helpers.WriteToDisk(filepath.Join(inpath, "config."+kind), &buf, fs.Source)
}

func nextStepsText() string {
	var nextStepsText bytes.Buffer

	nextStepsText.WriteString(`Just a few more steps and you're ready to go:

1. Download a theme into the same-named folder.
   Choose a theme from https://themes.gohugo.io/ or
   create your own with the "hugo new theme <THEMENAME>" command.
2. Perhaps you want to add some content. You can add single files
   with "hugo new `)

	nextStepsText.WriteString(filepath.Join("<SECTIONNAME>", "<FILENAME>.<FORMAT>"))

	nextStepsText.WriteString(`".
3. Start the built-in live server via "hugo server".

Visit https://gohugo.io/ for quickstart guide and full documentation.`)

	return nextStepsText.String()
}
