// Copyright © 2014-2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
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
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/create"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/parser"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var siteType string
var configFormat string
var contentType string
var contentFormat string
var contentFrontMatter string

func init() {
	newSiteCmd.Flags().StringVarP(&configFormat, "format", "f", "toml", "config & frontmatter format")
	newCmd.Flags().StringVarP(&configFormat, "format", "f", "toml", "frontmatter format")
	newCmd.Flags().StringVarP(&contentType, "kind", "k", "", "Content type to create")
	newCmd.AddCommand(newSiteCmd)
	newCmd.AddCommand(newThemeCmd)
}

var newCmd = &cobra.Command{
	Use:   "new [path]",
	Short: "Create new content for your site",
	Long: `Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.
You can also specify the kind with -k KIND
If archetypes are provided in your theme or site, they will be used.
`,
	Run: NewContent,
}

var newSiteCmd = &cobra.Command{
	Use:   "site [path]",
	Short: "Create a new site (skeleton)",
	Long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use 'hugo new [contentPath]' to create new content.
	`,
	Run: NewSite,
}

var newThemeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "Create a new theme",
	Long: `Create a new theme (skeleton) called [name] in the current directory.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
as you see fit.
	`,
	Run: NewTheme,
}

// NewContent adds new content to a Hugo site.
func NewContent(cmd *cobra.Command, args []string) {
	InitializeConfig()

	if cmd.Flags().Lookup("format").Changed {
		viper.Set("MetaDataFormat", configFormat)
	}

	if len(args) < 1 {
		cmd.Usage()
		jww.FATAL.Fatalln("path needs to be provided")
	}

	createpath := args[0]

	var kind string

	createpath, kind = newContentPathSection(createpath)

	if contentType != "" {
		kind = contentType
	}

	err := create.NewContent(kind, createpath)
	if err != nil {
		jww.ERROR.Println(err)
	}
}

// NewSite creates a new hugo site and initializes a structured Hugo directory.
func NewSite(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		jww.FATAL.Fatalln("path needs to be provided")
	}

	createpath, err := filepath.Abs(filepath.Clean(args[0]))
	if err != nil {
		cmd.Usage()
		jww.FATAL.Fatalln(err)
	}

	if x, _ := helpers.Exists(createpath, hugofs.SourceFs); x {
		y, _ := helpers.IsDir(createpath, hugofs.SourceFs)
		if z, _ := helpers.IsEmpty(createpath, hugofs.SourceFs); y && z {
			jww.INFO.Println(createpath, "already exists and is empty")
		} else {
			jww.FATAL.Fatalln(createpath, "already exists and is not empty")
		}
	}

	mkdir(createpath, "layouts")
	mkdir(createpath, "content")
	mkdir(createpath, "archetypes")
	mkdir(createpath, "static")
	mkdir(createpath, "data")

	createConfig(createpath, configFormat)
}

// NewTheme creates a new Hugo theme.
func NewTheme(cmd *cobra.Command, args []string) {
	InitializeConfig()

	if len(args) < 1 {
		cmd.Usage()
		jww.FATAL.Fatalln("theme name needs to be provided")
	}

	createpath := helpers.AbsPathify(filepath.Join("themes", args[0]))
	jww.INFO.Println("creating theme at", createpath)

	if x, _ := helpers.Exists(createpath, hugofs.SourceFs); x {
		jww.FATAL.Fatalln(createpath, "already exists")
	}

	mkdir(createpath, "layouts", "_default")
	mkdir(createpath, "layouts", "partials")

	touchFile(createpath, "layouts", "index.html")
	touchFile(createpath, "layouts", "_default", "list.html")
	touchFile(createpath, "layouts", "_default", "single.html")

	touchFile(createpath, "layouts", "partials", "header.html")
	touchFile(createpath, "layouts", "partials", "footer.html")

	mkdir(createpath, "archetypes")
	touchFile(createpath, "archetypes", "default.md")

	mkdir(createpath, "static", "js")
	mkdir(createpath, "static", "css")

	by := []byte(`The MIT License (MIT)

Copyright (c) ` + time.Now().Format("2006") + ` YOUR_NAME_HERE

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

	err := helpers.WriteToDisk(filepath.Join(createpath, "LICENSE.md"), bytes.NewReader(by), hugofs.SourceFs)
	if err != nil {
		jww.FATAL.Fatalln(err)
	}

	createThemeMD(createpath)
}

func mkdir(x ...string) {
	p := filepath.Join(x...)

	err := os.MkdirAll(p, 0777) // before umask
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}

func touchFile(x ...string) {
	inpath := filepath.Join(x...)
	mkdir(filepath.Dir(inpath))
	err := helpers.WriteToDisk(inpath, bytes.NewReader([]byte{}), hugofs.SourceFs)
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}

func createThemeMD(inpath string) (err error) {

	by := []byte(`# theme.toml template for a Hugo theme
# See https://github.com/spf13/hugoThemes#themetoml for an example

name = "` + strings.Title(helpers.MakeTitle(filepath.Base(inpath))) + `"
license = "MIT"
licenselink = "https://github.com/yourname/yourtheme/blob/master/LICENSE.md"
description = ""
homepage = "http://siteforthistheme.com/"
tags = ["", ""]
features = ["", ""]
min_version = 0.13

[author]
  name = ""
  homepage = ""

# If porting an existing theme
[original]
  name = ""
  homepage = ""
  repo = ""
`)

	err = helpers.WriteToDisk(filepath.Join(inpath, "theme.toml"), bytes.NewReader(by), hugofs.SourceFs)
	if err != nil {
		return
	}

	return nil
}

func newContentPathSection(path string) (string, string) {
	// Forward slashes is used in all examples. Convert if needed.
	// Issue #1133
	createpath := strings.Replace(path, "/", helpers.FilePathSeparator, -1)
	var section string
	// assume the first directory is the section (kind)
	if strings.Contains(createpath[1:], helpers.FilePathSeparator) {
		section = helpers.GuessSection(createpath)
	}

	return createpath, section
}

func createConfig(inpath string, kind string) (err error) {
	in := map[string]string{
		"baseurl":      "http://replace-this-with-your-hugo-site.com/",
		"title":        "My New Hugo Site",
		"languageCode": "en-us",
	}
	kind = parser.FormatSanitize(kind)

	by, err := parser.InterfaceToConfig(in, parser.FormatToLeadRune(kind))
	if err != nil {
		return err
	}

	err = helpers.WriteToDisk(filepath.Join(inpath, "config."+kind), bytes.NewReader(by), hugofs.SourceFs)
	if err != nil {
		return
	}

	return nil
}
