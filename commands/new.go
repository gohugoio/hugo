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
	hugoBuilderCommon
	contentEditor string
	contentType   string

	*baseCmd
}

func newNewCmd() *newCmd {
	cc := &newCmd{}
	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "new [path]",
		Short: "Create new content for your site",
		Long: `Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.

You can also specify the kind with ` + "`-k KIND`" + `.

If archetypes are provided in your theme or site, they will be used.`,

		RunE: cc.newContent,
	})

	cc.cmd.Flags().StringVarP(&cc.contentType, "kind", "k", "", "content type to create")
	cc.cmd.PersistentFlags().StringVarP(&cc.source, "source", "s", "", "filesystem path to read files relative from")
	cc.cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	cc.cmd.Flags().StringVar(&cc.contentEditor, "editor", "", "edit new content with this editor, if provided")

	cc.cmd.AddCommand(newNewSiteCmd().getCommand())
	cc.cmd.AddCommand(newNewThemeCmd().getCommand())

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

	createPath, kind = newContentPathSection(createPath)

	if n.contentType != "" {
		kind = n.contentType
	}

	cfg := c.DepsCfg

	ps, err := helpers.NewPathSpec(cfg.Fs, cfg.Cfg)
	if err != nil {
		return err
	}

	// If a site isn't in use in the archetype template, we can skip the build.
	siteFactory := func(filename string, siteUsed bool) (*hugolib.Site, error) {
		if !siteUsed {
			return hugolib.NewSite(*cfg)
		}
		var s *hugolib.Site

		if err := c.hugo.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
			return nil, err
		}

		s = c.hugo.Sites[0]

		if len(c.hugo.Sites) > 1 {
			// Find the best match.
			for _, ss := range c.hugo.Sites {
				if strings.Contains(createPath, "."+ss.Language.Lang) {
					s = ss
					break
				}
			}
		}
		return s, nil
	}

	return create.NewContent(ps, siteFactory, kind, createPath)
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

func newContentPathSection(path string) (string, string) {
	// Forward slashes is used in all examples. Convert if needed.
	// Issue #1133
	createpath := filepath.FromSlash(path)
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
