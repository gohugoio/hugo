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
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*newThemeCmd)(nil)

type newThemeCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newNewThemeCmd() *newThemeCmd {
	cc := &newThemeCmd{}

	cmd := &cobra.Command{
		Use:   "theme [name]",
		Short: "Create a new theme",
		Long: `Create a new theme (skeleton) called [name] in ./themes.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
as you see fit.`,
		RunE: cc.newTheme,
	}

	cc.baseBuilderCmd = b.newBuilderBasicCmd(cmd)

	return cc
}

// newTheme creates a new Hugo theme template
func (n *newThemeCmd) newTheme(cmd *cobra.Command, args []string) error {
	c, err := initializeConfig(false, false, false, &n.hugoBuilderCommon, n, nil)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return newUserError("theme name needs to be provided")
	}

	createpath := c.hugo().PathSpec.AbsPathify(filepath.Join(c.Cfg.GetString("themesDir"), args[0]))
	jww.FEEDBACK.Println("Creating theme at", createpath)

	cfg := c.DepsCfg

	if x, _ := helpers.Exists(createpath, cfg.Fs.Source); x {
		return errors.New(createpath + " already exists")
	}

	mkdir(createpath, "layouts", "_default")
	mkdir(createpath, "layouts", "partials")

	touchFile(cfg.Fs.Source, createpath, "layouts", "index.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "404.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "_default", "list.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "_default", "single.html")

	baseofDefault := []byte(`<!DOCTYPE html>
<html>
    {{- partial "head.html" . -}}
    <body>
        {{- partial "header.html" . -}}
        <div id="content">
        {{- block "main" . }}{{- end }}
        </div>
        {{- partial "footer.html" . -}}
    </body>
</html>
`)
	err = helpers.WriteToDisk(filepath.Join(createpath, "layouts", "_default", "baseof.html"), bytes.NewReader(baseofDefault), cfg.Fs.Source)
	if err != nil {
		return err
	}

	touchFile(cfg.Fs.Source, createpath, "layouts", "partials", "head.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "partials", "header.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "partials", "footer.html")

	mkdir(createpath, "archetypes")

	archDefault := []byte("+++\n+++\n")

	err = helpers.WriteToDisk(filepath.Join(createpath, "archetypes", "default.md"), bytes.NewReader(archDefault), cfg.Fs.Source)
	if err != nil {
		return err
	}

	mkdir(createpath, "static", "js")
	mkdir(createpath, "static", "css")

	by := []byte(`The MIT License (MIT)

Copyright (c) ` + htime.Now().Format("2006") + ` YOUR_NAME_HERE

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`)

	err = helpers.WriteToDisk(filepath.Join(createpath, "LICENSE"), bytes.NewReader(by), cfg.Fs.Source)
	if err != nil {
		return err
	}

	n.createThemeMD(cfg.Fs, createpath)

	return nil
}

func (n *newThemeCmd) createThemeMD(fs *hugofs.Fs, inpath string) (err error) {
	by := []byte(`# theme.toml template for a Hugo theme
# See https://github.com/gohugoio/hugoThemes#themetoml for an example

name = "` + strings.Title(helpers.MakeTitle(filepath.Base(inpath))) + `"
license = "MIT"
licenselink = "https://github.com/yourname/yourtheme/blob/master/LICENSE"
description = ""
homepage = "http://example.com/"
tags = []
features = []
min_version = "0.41.0"

[author]
  name = ""
  homepage = ""

# If porting an existing theme
[original]
  name = ""
  homepage = ""
  repo = ""
`)

	err = helpers.WriteToDisk(filepath.Join(inpath, "theme.toml"), bytes.NewReader(by), fs.Source)
	if err != nil {
		return
	}

	return nil
}
