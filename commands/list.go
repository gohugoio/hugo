// Copyright © 2014 Steve Francia <spf@spf13.com>.
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
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/viper"
)

func init() {
	listCmd.AddCommand(listDraftsCmd)
	listCmd.AddCommand(listFutureCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing out various types of content",
	Long:  `Listing out various types of content. List requires a subcommand, eg. hugo list drafts`,
	Run:   nil,
}

var listDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "List all drafts",
	Long:  `List all of the drafts in your content directory`,
	Run: func(cmd *cobra.Command, args []string) {

		InitializeConfig()
		viper.Set("BuildDrafts", true)

		site := &hugolib.Site{}

		if err := site.Process(); err != nil {
			fmt.Println("Error Processing Source Content", err)
		}

		for _, p := range site.Pages {
			if p.IsDraft() {
				fmt.Println(filepath.Join(p.File.Dir(), p.File.LogicalName()))
			}

		}

	},
}

var listFutureCmd = &cobra.Command{
	Use:   "future",
	Short: "List all posts dated in the future",
	Long:  `List all of the posts in your content directory who will be posted in the future`,
	Run: func(cmd *cobra.Command, args []string) {

		InitializeConfig()
		viper.Set("BuildFuture", true)

		site := &hugolib.Site{}

		if err := site.Process(); err != nil {
			fmt.Println("Error Processing Source Content", err)
		}

		for _, p := range site.Pages {
			if p.IsFuture() {
				fmt.Println(filepath.Join(p.File.Dir(), p.File.LogicalName()))
			}

		}

	},
}
