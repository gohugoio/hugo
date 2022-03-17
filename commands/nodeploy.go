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

//go:build nodeploy
// +build nodeploy

package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var _ cmder = (*deployCmd)(nil)

// deployCmd supports deploying sites to Cloud providers.
type deployCmd struct {
	*baseBuilderCmd
}

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
			return errors.New("build without HUGO_BUILD_TAGS=nodeploy to use this command")
		},
	}

	cc.baseBuilderCmd = b.newBuilderBasicCmd(cmd)

	return cc
}
