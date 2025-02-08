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

//go:build !withdeploy

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
	"errors"

	"github.com/bep/simplecobra"
	"github.com/spf13/cobra"
)

func newDeployCommand() simplecobra.Commander {
	return &simpleCommand{
		name: "deploy",
		run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
			return errors.New("deploy not supported in this version of Hugo; install a release with 'withdeploy' in the archive filename or build yourself with the 'withdeploy' build tag. Also see https://github.com/gohugoio/hugo/pull/12995")
		},
		withc: func(cmd *cobra.Command, r *rootCommand) {
			applyDeployFlags(cmd, r)
			cmd.Hidden = true
		},
	}
}
