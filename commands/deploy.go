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

//go:build withdeploy
// +build withdeploy

package commands

import (
	"context"

	"github.com/gohugoio/hugo/deploy"

	"github.com/bep/simplecobra"
	"github.com/spf13/cobra"
)

func newDeployCommand() simplecobra.Commander {
	return &simpleCommand{
		name:  "deploy",
		short: "Deploy your site to a cloud provider",
		long: `Deploy your site to a cloud provider

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
			applyDeployFlags(cmd, r)
		},
	}
}
