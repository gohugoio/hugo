// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"github.com/spf13/cobra"
)

var _ cmder = (*deployCmd)(nil)

// deployCmd supports deploying sites to Cloud providers.
type deployCmd struct {
	*baseBuilderCmd
}

// TODO: In addition to the "deploy" command, consider adding a "--deploy"
// flag for the default command; this would build the site and then deploy it.
// It's not obvious how to do this; would all of the deploy-specific flags
// have to exist at the top level as well?

// TODO:  The output files change every time "hugo" is executed, it looks
// like because of map order randomization. This means that you can
// run "hugo && hugo deploy" again and again and upload new stuff every time. Is
// this intended?

func (b *commandsBuilder) newDeployCmd() *deployCmd {
	cc := &deployCmd{}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your site to a Cloud provider.",
		Long: `Deploy your site to a Cloud provider.

See https://gohugo.io/hosting-and-deployment/hugo-deploy/ for detailed
documentation.
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			cfgInit := func(c *commandeer) error {
				return nil
			}
			comm, err := initializeConfig(true, false, &cc.hugoBuilderCommon, cc, cfgInit)
			if err != nil {
				return err
			}
			deployer, err := deploy.New(comm.Cfg, comm.hugo().PathSpec.PublishFs)
			if err != nil {
				return err
			}
			return deployer.Deploy(context.Background())
		},
	}

	cmd.Flags().String("target", "", "target deployment from deployments section in config file; defaults to the first one")
	cmd.Flags().Bool("confirm", false, "ask for confirmation before making changes to the target")
	cmd.Flags().Bool("dryRun", false, "dry run")
	cmd.Flags().Bool("force", false, "force upload of all files")
	cmd.Flags().Bool("invalidateCDN", true, "invalidate the CDN cache listed in the deployment target")
	cmd.Flags().Int("maxDeletes", 256, "maximum # of files to delete, or -1 to disable")

	cc.baseBuilderCmd = b.newBuilderBasicCmd(cmd)

	return cc
}
