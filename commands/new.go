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
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/create"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*newCmd)(nil)

type newCmd struct {
	contentEditor string
	contentType   string

	*baseBuilderCmd
}

func (b *commandsBuilder) newNewCmd() *newCmd {
	cmd := &cobra.Command{
		Use:   "new [path]",
		Short: "Create new content for your site",
		Long: `Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.

You can also specify the kind with ` + "`-k KIND`" + `.

If archetypes are provided in your theme or site, they will be used.

Ensure you run this within the root directory of your site.`,
	}

	cc := &newCmd{baseBuilderCmd: b.newBuilderCmd(cmd)}

	cmd.Flags().StringVarP(&cc.contentType, "kind", "k", "", "content type to create")
	cmd.Flags().StringVar(&cc.contentEditor, "editor", "", "edit new content with this editor, if provided")

	cmd.AddCommand(b.newNewSiteCmd().getCommand())
	cmd.AddCommand(b.newNewThemeCmd().getCommand())

	cmd.RunE = cc.newContent

	return cc
}

func (n *newCmd) newContent(cmd *cobra.Command, args []string) error {
	cfgInit := func(c *commandeer) error {
		if cmd.Flags().Changed("editor") {
			c.Set("newContentEditor", n.contentEditor)
		}
		return nil
	}

	c, err := initializeConfig(true, false, &n.hugoBuilderCommon, n, cfgInit)

	if err != nil {
		return err
	}

	if len(args) < 1 {
		return newUserError("path needs to be provided")
	}

	createPath := args[0]

	var kind string

	createPath, kind = newContentPathSection(c.hugo(), createPath)

	if n.contentType != "" {
		kind = n.contentType
	}

	return create.NewContent(c.hugo(), kind, createPath)
}

func mkdir(x ...string) {
	p := filepath.Join(x...)

	err := os.MkdirAll(p, 0777) // before umask
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}

func touchFile(fs afero.Fs, x ...string) {
	inpath := filepath.Join(x...)
	mkdir(filepath.Dir(inpath))
	err := helpers.WriteToDisk(inpath, bytes.NewReader([]byte{}), fs)
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}

func newContentPathSection(h *hugolib.HugoSites, path string) (string, string) {
	// Forward slashes is used in all examples. Convert if needed.
	// Issue #1133
	createpath := filepath.FromSlash(path)

	if h != nil {
		for _, dir := range h.BaseFs.Content.Dirs {
			createpath = strings.TrimPrefix(createpath, dir.Meta().Filename())
		}
	}

	var section string
	// assume the first directory is the section (kind)
	if strings.Contains(createpath[1:], helpers.FilePathSeparator) {
		parts := strings.Split(strings.TrimPrefix(createpath, helpers.FilePathSeparator), helpers.FilePathSeparator)
		if len(parts) > 0 {
			section = parts[0]
		}

	}

	return createpath, section
}
