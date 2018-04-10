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
	"path/filepath"

	"github.com/spf13/viper"

	"errors"

	"github.com/gohugoio/hugo/create"

	"fmt"

	"strings"

	"github.com/gohugoio/hugo/parser"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*newSiteCmd)(nil)

type newSiteCmd struct {
	configFormat string

	*baseCmd
}

func newNewSiteCmd() *newSiteCmd {
	ccmd := &newSiteCmd{}

	cmd := &cobra.Command{
		Use:   "site [path]",
		Short: "Create a new site (skeleton)",
		Long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use ` + "`hugo new [contentPath]`" + ` to create new content.`,
		RunE: ccmd.newSite,
	}

	cmd.Flags().StringVarP(&ccmd.configFormat, "format", "f", "toml", "config & frontmatter format")
	cmd.Flags().Bool("force", false, "init inside non-empty directory")

	ccmd.baseCmd = newBaseCmd(cmd)

	return ccmd

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
			return errors.New(basepath + " already exists and is not empty")

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
			return fmt.Errorf("Failed to create dir: %s", err)
		}
	}

	createConfig(fs, basepath, n.configFormat)

	// Create a defaul archetype file.
	helpers.SafeWriteToDisk(filepath.Join(archeTypePath, "default.md"),
		strings.NewReader(create.ArchetypeTemplateTemplate), fs.Source)

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

	return n.doNewSite(hugofs.NewDefault(viper.New()), createpath, forceNew)
}

func createConfig(fs *hugofs.Fs, inpath string, kind string) (err error) {
	in := map[string]string{
		"baseURL":      "http://example.org/",
		"title":        "My New Hugo Site",
		"languageCode": "en-us",
	}
	kind = parser.FormatSanitize(kind)

	var buf bytes.Buffer
	err = parser.InterfaceToConfig(in, parser.FormatToLeadRune(kind), &buf)
	if err != nil {
		return err
	}

	return helpers.WriteToDisk(filepath.Join(inpath, "config."+kind), &buf, fs.Source)
}
