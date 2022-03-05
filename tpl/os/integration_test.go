// Copyright 2022 The Hugo Authors. All rights reserved.
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

package os_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue 9599
func TestReadDir(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
theme = "mytheme"
-- myproject.txt --
Hello project!
-- themes/mytheme/mytheme.txt --
Hello theme!
-- themes/mytheme/archetypes/mytheme-default.md --
draft: true
-- archetypes/default.md --
draft: true
-- content/content-example.md --
draft: true
-- layouts/404.html --
<html></html>
-- data/data.csv --
1,2,3
-- assets/asset.jpg --
YWJjMTIzIT8kKiYoKSctPUB+
-- i18n/en-NZ.yaml --
summarized: "summarised"
-- layouts/index.html --
{{ $archetypeentries := (readDir "." "archetypes") }}
START archteypes:|{{ range $entry := $archetypeentries }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
{{ $content := (readDir "." "content") }}
START content:|{{ range $entry := $content }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
{{ $layouts := (readDir "." "layouts") }}
START layouts:|{{ range $entry := $layouts }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
{{ $data := (readDir "." "data") }}
START data:|{{ range $entry := $data }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
{{ $assets := (readDir "." "assets") }}
START assets:|{{ range $entry := $assets }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
{{ $i18n := (readDir "." "i18n") }}
START i18n:|{{ range $entry := $i18n }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
{{ $entries := (readDir ".") }}
START work:|{{ range $entry := $entries }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/index.html",
		"START archteypes:|default.md|mytheme-default.md|:END:",
		"START content:|content-example.md|:END:",
		"START layouts:|404.html|index.html|:END:",
		"START data:|data.csv|:END:",
		"START assets:|asset.jpg|:END:",
		"START i18n:|en-NZ.yaml|:END:",
		"START work:|config.toml|myproject.txt|:END:")
}
