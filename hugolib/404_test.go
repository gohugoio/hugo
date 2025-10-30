// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"fmt"
	"testing"
)

func Test404(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "taxonomy", "term"]
baseURL = "http://example.com/"	
-- layouts/all.html --
All. {{ .Kind }}. {{ .Title }}|Lastmod: {{ .Lastmod.Format "2006-01-02" }}|
-- layouts/404.html --
{{ $home := site.Home }}
404: 
Parent: {{ .Parent.Kind }}|{{ .Parent.Path }}|
IsAncestor: {{ .IsAncestor $home }}/{{ $home.IsAncestor . }}
IsDescendant: {{ .IsDescendant $home }}/{{ $home.IsDescendant . }}
CurrentSection: {{ .CurrentSection.Kind }}|
FirstSection: {{ .FirstSection.Kind }}|
InSection: {{ .InSection $home.Section }}|{{ $home.InSection . }}
Sections: {{ len .Sections }}|
Page: {{ .Page.RelPermalink }}|
Data: {{ len .Data }}|
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", "All. home. |")

	// Note: We currently have only 1 404 page. One might think that we should have
	// multiple, to follow the Custom Output scheme, but I don't see how that would work
	// right now.
	b.AssertFileContent("public/404.html", `
  404:
Parent: home
IsAncestor: false/true
IsDescendant: true/false
CurrentSection: home|
FirstSection: home|
InSection: false|true
Sections: 0|
Page: /404.html|
Data: 1|
        
`)
}

func Test404WithBase(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplates("404.html", `{{ define "main" }}
Page not found
{{ end }}`,
		"baseof.html", `Base: {{ block "main" . }}{{ end }}`).WithContent("page.md", ``)

	b.Build(BuildCfg{})

	// Note: We currently have only 1 404 page. One might think that we should have
	// multiple, to follow the Custom Output scheme, but I don't see how that would work
	// right now.
	b.AssertFileContent("public/404.html", `
Base:
Page not found`)
}

func Test404EditTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
disableLiveReload = true
[internal]
fastRenderMode = true
-- layouts/_default/baseof.html --
Base: {{ block "main" . }}{{ end }}
-- layouts/404.html --
{{ define "main" }}
Not found.
{{ end }}
	
	`

	b := TestRunning(t, files)

	b.AssertFileContent("public/404.html", `Not found.`)

	b.EditFiles("layouts/404.html", `Not found. Updated.`).Build()

	fmt.Println("Rebuilding")
	b.BuildPartial("/does-not-exist")

	b.AssertFileContent("public/404.html", `Not found. Updated.`)
}
