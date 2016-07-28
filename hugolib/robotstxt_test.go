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
	"bytes"
	"testing"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
)

const robotTxtTemplate = `User-agent: Googlebot
  {{ range .Data.Pages }}
	Disallow: {{.RelPermalink}}
	{{ end }}
`

func TestRobotsTXTOutput(t *testing.T) {
	testCommonResetState()

	hugofs.InitMemFs()

	viper.Set("baseurl", "http://auth/bub/")
	viper.Set("enableRobotsTXT", true)

	s := &Site{
		Source:   &source.InMemorySource{ByteSource: weightedSources},
		Language: newDefaultLanguage(),
	}

	if err := buildAndRenderSite(s, "robots.txt", robotTxtTemplate); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	robotsFile, err := hugofs.Destination().Open("public/robots.txt")

	if err != nil {
		t.Fatalf("Unable to locate: robots.txt")
	}

	robots := helpers.ReaderToBytes(robotsFile)
	if !bytes.HasPrefix(robots, []byte("User-agent: Googlebot")) {
		t.Errorf("Robots file should start with 'User-agentL Googlebot'. %s", robots)
	}
}
