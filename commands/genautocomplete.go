// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*genautocompleteCmd)(nil)

type genautocompleteCmd struct {
	autocompleteTarget string

	// bash for now (zsh and others will come)
	autocompleteType string

	*baseCmd
}

func newGenautocompleteCmd() *genautocompleteCmd {
	cc := &genautocompleteCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "autocomplete",
		Short: "Generate shell autocompletion script for Hugo",
		Long: `Generates a shell autocompletion script for Hugo.

NOTE: The current version supports Bash only.
      This should work for *nix systems with Bash installed.

By default, the file is written directly to /etc/bash_completion.d
for convenience, and the command may need superuser rights, e.g.:

	$ sudo hugo gen autocomplete

Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
file-path and name.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if cc.autocompleteType != "bash" {
				return newUserError("Only Bash is supported for now")
			}

			err := cmd.Root().GenBashCompletionFile(cc.autocompleteTarget)

			if err != nil {
				return err
			}

			jww.FEEDBACK.Println("Bash completion file for Hugo saved to", cc.autocompleteTarget)

			return nil
		},
	})

	cc.cmd.PersistentFlags().StringVarP(&cc.autocompleteTarget, "completionfile", "", "/etc/bash_completion.d/hugo.sh", "autocompletion file")
	cc.cmd.PersistentFlags().StringVarP(&cc.autocompleteType, "type", "", "bash", "autocompletion type (currently only bash supported)")

	// For bash-completion
	cc.cmd.PersistentFlags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})

	return cc
}
