// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
)

var genmandir string
var genmanCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages for the Hugo CLI",
	Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		header := &doc.GenManHeader{
			Section: "1",
			Manual:  "Hugo Manual",
			Source:  fmt.Sprintf("Hugo %s", helpers.CurrentHugoVersion),
		}
		if !strings.HasSuffix(genmandir, helpers.FilePathSeparator) {
			genmandir += helpers.FilePathSeparator
		}
		if found, _ := helpers.Exists(genmandir, hugofs.Os); !found {
			jww.FEEDBACK.Println("Directory", genmandir, "does not exist, creating...")
			if err := hugofs.Os.MkdirAll(genmandir, 0777); err != nil {
				return err
			}
		}
		cmd.Root().DisableAutoGenTag = true

		jww.FEEDBACK.Println("Generating Hugo man pages in", genmandir, "...")
		doc.GenManTree(cmd.Root(), header, genmandir)

		jww.FEEDBACK.Println("Done.")

		return nil
	},
}

func init() {
	genmanCmd.PersistentFlags().StringVar(&genmandir, "dir", "man/", "the directory to write the man pages.")

	// For bash-completion
	genmanCmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})
}
