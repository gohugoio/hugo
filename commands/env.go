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
	"runtime"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var _ cmder = (*envCmd)(nil)

type envCmd struct {
	*baseCmd
}

func newEnvCmd() *envCmd {
	return &envCmd{
		baseCmd: newBaseCmd(&cobra.Command{
			Use:   "env",
			Short: "Print Hugo version and environment info",
			Long: `Print Hugo version and environment info. This is useful in Hugo bug reports.

If you add the -v flag, you will get a full dependency list.
`,
			RunE: func(cmd *cobra.Command, args []string) error {
				printHugoVersion()
				jww.FEEDBACK.Printf("GOOS=%q\n", runtime.GOOS)
				jww.FEEDBACK.Printf("GOARCH=%q\n", runtime.GOARCH)
				jww.FEEDBACK.Printf("GOVERSION=%q\n", runtime.Version())

				isVerbose, _ := cmd.Flags().GetBool("verbose")

				if isVerbose {
					deps := hugo.GetDependencyList()
					for _, dep := range deps {
						jww.FEEDBACK.Printf("%s\n", dep)
					}
				}

				return nil
			},
		}),
	}

}
