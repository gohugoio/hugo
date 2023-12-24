// Copyright 2021 The Hugo Authors. All rights reserved.
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

package openapi3_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	files := `
-- assets/api/myapi.yaml --
openapi: 3.0.0
info:
  title: Sample API
  description: Optional multiline or single-line description in [CommonMark](http://commonmark.org/help/) or HTML.
  version: 0.1.9
servers:
  - url: http://api.example.com/v1
    description: Optional server description, e.g. Main (production) server
  - url: http://staging-api.example.com
    description: Optional server description, e.g. Internal staging server for testing
paths:
  /users:
    get:
      summary: Returns a list of users.
      description: Optional extended description in CommonMark or HTML.
      responses:
        '200':    # status code
          description: A JSON array of user names
          content:
            application/json:
              schema: 
                type: array
                items: 
                  type: string
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
{{ $api := resources.Get "api/myapi.yaml" | openapi3.Unmarshal }}
API: {{ $api.Info.Title | safeHTML }}
  `

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			Running:     true,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", `API: Sample API`)

	b.
		EditFileReplaceFunc("assets/api/myapi.yaml", func(s string) string { return strings.ReplaceAll(s, "Sample API", "Hugo API") }).
		Build()

	b.AssertFileContent("public/index.html", `API: Hugo API`)
}
