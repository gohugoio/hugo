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

package hugolib

import (
	"testing"

	"github.com/spf13/viper"
)

const robotTxtTemplate = `User-agent: Googlebot
  {{ range .Data.Pages }}
	Disallow: {{.RelPermalink}}
	{{ end }}
`

func TestRobotsTXTOutput(t *testing.T) {
	t.Parallel()

	cfg := viper.New()
	cfg.Set("baseURL", "http://auth/bub/")
	cfg.Set("enableRobotsTXT", true)

	b := newTestSitesBuilder(t).WithViper(cfg)
	b.WithTemplatesAdded("layouts/robots.txt", robotTxtTemplate)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/robots.txt", "User-agent: Googlebot")

}
