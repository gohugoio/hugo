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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/create"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
)

var siteType string
var configFormat string
var contentType string
var contentFormat string
var contentFrontMatter string

func init() {
	//newSiteCmd.Flags().StringVarP(&siteType, "type", "t", "blog", "What type of site to new")
	newSiteCmd.Flags().StringVarP(&configFormat, "format", "f", "yaml", "Config file format")
	newCmd.Flags().StringVarP(&contentType, "kind", "k", "", "Content type to create")
	newCmd.AddCommand(newSiteCmd)
}

var newCmd = &cobra.Command{
	Use:   "new [path]",
	Short: "Create new content for your site",
	Long: `Create will create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.
You can also specify the kind with -k KIND
If archetypes are provided in your theme or site, they will be used.
`,
	Run: NewContent,
}

func NewContent(cmd *cobra.Command, args []string) {
	InitializeConfig()

	if len(args) < 1 {
		jww.FATAL.Fatalln("path needs to be provided")
	}

	createpath := args[0]

	var kind string

	// assume the first directory is the section (kind)
	if strings.Contains(createpath[1:], "/") {
		kind = helpers.GuessSection(createpath)
	}

	if contentType != "" {
		kind = contentType
	}

	err := create.NewContent(kind, createpath)
	if err != nil {
		jww.ERROR.Println(err)
	}
}

var newSiteCmd = &cobra.Command{
	Use:   "site [type]",
	Short: "Create a new site of [type]",
	Long:  `Create a new site as a (blog, project, etc)`,
	Run:   NewSite,
}

func NewSite(cmd *cobra.Command, args []string) {
	InitializeConfig()

	fmt.Println("new site called")
	fmt.Println(args)
}
