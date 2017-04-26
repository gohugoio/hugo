// +build release

// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"github.com/spf13/hugo/releaser"
)

func init() {
	HugoCmd.AddCommand(createReleaser().cmd)
}

type releaseCommandeer struct {
	cmd *cobra.Command

	// Will be zero for main releases.
	patchLevel int

	skipPublish bool

	step int
}

func createReleaser() *releaseCommandeer {
	// Note: This is a command only meant for internal use and must be run
	// via "go run -tags release main.go release" on the actual code base that is in the release.
	r := &releaseCommandeer{
		cmd: &cobra.Command{
			Use:    "release",
			Short:  "Release a new version of Hugo.",
			Hidden: true,
		},
	}

	r.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return r.release()
	}

	r.cmd.PersistentFlags().IntVarP(&r.patchLevel, "patch", "p", 0, "patch level, defaults to 0 for main releases")
	r.cmd.PersistentFlags().IntVarP(&r.step, "step", "s", -1, "release step, defaults to -1 for all steps.")
	r.cmd.PersistentFlags().BoolVarP(&r.skipPublish, "skip-publish", "", false, "skip all publishing pipes of the release")

	return r
}

func (r *releaseCommandeer) release() error {
	return releaser.New(r.patchLevel, r.step, r.skipPublish).Run()
}
