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

	"github.com/gohugoio/hugo/common/hugo"
)

var _ cmder = (*versionCmd)(nil)

type versionCmd struct {
	json bool
	*baseCmd
}

func newVersionCmd() *versionCmd {
	v := &versionCmd{}

	v.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if v.json {
				printHugoVersionJson()
				return nil
			}
			printHugoVersion()

			return nil
		},
	})

	v.cmd.Flags().BoolVar(&v.json, "json", false, "print JSON format")

	return v
}

func printHugoVersion() {
	jww.FEEDBACK.Println(hugo.BuildVersionString())
}

func printHugoVersionJson() {
	jww.FEEDBACK.Println(hugo.BuildVersionJSON())
}
