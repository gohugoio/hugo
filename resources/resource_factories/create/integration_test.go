// Copyright 2023 The Hugo Authors. All rights reserved.
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

package create_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestGetResourceHead(t *testing.T) {

	files := `
-- config.toml --
[security]
  [security.http]
    methods = ['(?i)GET|POST|HEAD']
    urls = ['.*gohugo\.io.*']

-- layouts/index.html --
{{ $url := "https://gohugo.io/img/hugo.png" }}
{{ $opts := dict "method" "head" }}
{{ with resources.GetRemote $url $opts }}
  {{ with .Err }}
    {{ errorf "Unable to get remote resource: %s" . }}
  {{ else }}
    Head Content: {{ .Content }}.
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource: %s" $url }}
{{ end }}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)

	b.Build()

	b.AssertFileContent("public/index.html", "Head Content: .")

}
