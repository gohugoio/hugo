// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"context"

	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/releaser"
	"github.com/spf13/cobra"
)

// Note: This is a command only meant for internal use and must be run
// via "go run -tags release main.go release" on the actual code base that is in the release.
func newReleaseCommand() simplecobra.Commander {
	var (
		step     int
		skipPush bool
		try      bool
	)

	return &simpleCommand{
		name:  "release",
		short: "Release a new version of Hugo.",
		run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
			rel, err := releaser.New(skipPush, try, step)
			if err != nil {
				return err
			}

			return rel.Run()
		},
		withc: func(cmd *cobra.Command, r *rootCommand) {
			cmd.Hidden = true
			cmd.ValidArgsFunction = cobra.NoFileCompletions
			cmd.PersistentFlags().BoolVarP(&skipPush, "skip-push", "", false, "skip pushing to remote")
			cmd.PersistentFlags().BoolVarP(&try, "try", "", false, "no changes")
			cmd.PersistentFlags().IntVarP(&step, "step", "", 0, "step to run (1: set new version 2: prepare next dev version)")
			_ = cmd.RegisterFlagCompletionFunc("step", cobra.FixedCompletions([]string{"1", "2"}, cobra.ShellCompDirectiveNoFileComp))
		},
	}
}
