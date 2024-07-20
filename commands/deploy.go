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

//go:build !nodeploy
// +build !nodeploy

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

	"github.com/gohugoio/hugo/deploy"
	"github.com/gohugoio/hugo/deploy/deployconfig"

	"github.com/bep/simplecobra"
	"github.com/spf13/cobra"
)

func newDeployCommand() simplecobra.Commander {
	return &simpleCommand{
		name:  "deploy",
		short: "Deploy your site to a Cloud provider.",
		long: `Deploy your site to a Cloud provider.

See https://gohugo.io/hosting-and-deployment/hugo-deploy/ for detailed
documentation.
`,
		run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
			h, err := r.Hugo(flagsToCfgWithAdditionalConfigBase(cd, nil, "deployment"))
			if err != nil {
				return err
			}
			deployer, err := deploy.New(h.Configs.GetFirstLanguageConfig(), h.Log, h.PathSpec.PublishFs)
			if err != nil {
				return err
			}
			return deployer.Deploy(ctx)
		},
		withc: func(cmd *cobra.Command, r *rootCommand) {
			cmd.ValidArgsFunction = cobra.NoFileCompletions
			cmd.Flags().String("target", "", "target deployment from deployments section in config file; defaults to the first one")
			_ = cmd.RegisterFlagCompletionFunc("target", cobra.NoFileCompletions)
			cmd.Flags().Bool("confirm", false, "ask for confirmation before making changes to the target")
			cmd.Flags().Bool("dryRun", false, "dry run")
			cmd.Flags().Bool("force", false, "force upload of all files")
			cmd.Flags().Bool("invalidateCDN", deployconfig.DefaultConfig.InvalidateCDN, "invalidate the CDN cache listed in the deployment target")
			cmd.Flags().Int("maxDeletes", deployconfig.DefaultConfig.MaxDeletes, "maximum # of files to delete, or -1 to disable")
			_ = cmd.RegisterFlagCompletionFunc("maxDeletes", cobra.NoFileCompletions)
			cmd.Flags().Int("workers", deployconfig.DefaultConfig.Workers, "number of workers to transfer files. defaults to 10")
			_ = cmd.RegisterFlagCompletionFunc("workers", cobra.NoFileCompletions)
		},
	}
}
