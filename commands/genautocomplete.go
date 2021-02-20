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
	"io"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*genautocompleteCmd)(nil)

type genautocompleteCmd struct {
	autocompleteTarget string

	// bash, zsh, fish or powershell
	autocompleteType string

	*baseCmd
}

func newGenautocompleteCmd() *genautocompleteCmd {
	cc := &genautocompleteCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "autocomplete",
		Short: "Generate shell autocompletion script for Hugo",
		Long: `Generates a shell autocompletion script for Hugo.

The script is written to the console (stdout).

To write to file, add the ` + "`--completionfile=/path/to/file`" + ` flag.

Add ` + "`--type={bash, zsh, fish or powershell}`" + ` flag to set alternative
shell type.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion or /path/to/file`,

		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var target io.Writer

			if cc.autocompleteTarget == "" {
				target = os.Stdout
			} else {
				target, _ = os.OpenFile(cc.autocompleteTarget, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			}

			switch cc.autocompleteType {
			case "bash":
				err = cmd.Root().GenBashCompletion(target)
			case "zsh":
				err = cmd.Root().GenZshCompletion(target)
			case "fish":
				err = cmd.Root().GenFishCompletion(target, true)
			case "powershell":
				err = cmd.Root().GenPowerShellCompletion(target)
			default:
				return newUserError("Unsupported completion type")
			}

			if err != nil {
				return err
			}

			if cc.autocompleteTarget != "" {
				jww.FEEDBACK.Println(cc.autocompleteType+" completion file for Hugo saved to", cc.autocompleteTarget)
			}
			return nil
		},
	})

	cc.cmd.PersistentFlags().StringVarP(&cc.autocompleteTarget, "completionfile", "f", "", "autocompletion file, defaults to stdout")
	cc.cmd.PersistentFlags().StringVarP(&cc.autocompleteType, "type", "t", "bash", "autocompletion type (bash, zsh, fish, or powershell)")

	return cc
}
