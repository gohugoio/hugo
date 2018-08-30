// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
)

func TestThemesGraph(t *testing.T) {
	t.Parallel()

	const (
		themeStandalone = `
title = "Theme Standalone"
[params]
v1 = "v1s"
v2 = "v2s"
`
		themeCyclic = `
title = "Theme Cyclic"
theme = "theme3"
[params]
v1 = "v1c"
v2 = "v2c"
`
		theme1 = `
title = "Theme #1"
theme = "themeStandalone"
[params]
v2 = "v21"
`

		theme2 = `
title = "Theme #2"
theme = "theme1"
[params]
v1 = "v12"
`

		theme3 = `
title = "Theme #3"
theme = ["theme2", "themeStandalone", "themeCyclic"]
[params]
v1 = "v13"
v2 = "v24"
`

		theme4 = `
title = "Theme #4"
theme = "theme3"
[params]
v1 = "v14"
v2 = "v24"
`

		site1 = `
			theme = "theme4"
			
			[params]
			v1 = "site"
`
		site2 = `
			theme = ["theme2", "themeStandalone"]
`
	)

	var (
		testConfigs = []struct {
			siteConfig string

			// The name of theme somewhere in the middle to write custom key/files.
			offset string

			check func(b *sitesBuilder)
		}{
			{site1, "theme3", func(b *sitesBuilder) {

				// site1: theme4 theme3 theme2 theme1 themeStandalone themeCyclic

				// Check data
				// theme3 should win the offset competition
				b.AssertFileContent("public/index.html", "theme1o::[offset][v]theme3", "theme4o::[offset][v]theme3", "themeStandaloneo::[offset][v]theme3")
				b.AssertFileContent("public/index.html", "nproject::[inner][other]project|[project][other]project|[theme][other]theme4|[theme1][other]theme1")
				b.AssertFileContent("public/index.html", "ntheme::[inner][other]theme4|[theme][other]theme4|[theme1][other]theme1|[theme2][other]theme2|[theme3][other]theme3")
				b.AssertFileContent("public/index.html", "theme1::[inner][other]project|[project][other]project|[theme][other]theme1|[theme1][other]theme1|")
				b.AssertFileContent("public/index.html", "theme4::[inner][other]project|[project][other]project|[theme][other]theme4|[theme4][other]theme4|")

				// Check layouts
				b.AssertFileContent("public/index.html", "partial ntheme: theme4", "partial theme2o: theme3")

				// Check i18n
				b.AssertFileContent("public/index.html", "i18n: project theme4")

				// Check static files
				// TODO(bep) static files not currently part of the build b.AssertFileContent("public/nproject.txt", "TODO")

				// Check site params
				b.AssertFileContent("public/index.html", "v1::site", "v2::v24")
			}},
			{site2, "", func(b *sitesBuilder) {

				// site2: theme2 theme1 themeStandalone
				b.AssertFileContent("public/index.html", "nproject::[inner][other]project|[project][other]project|[theme][other]theme2|[theme1][other]theme1|[theme2][other]theme2|[themeStandalone][other]themeStandalone|")
				b.AssertFileContent("public/index.html", "ntheme::[inner][other]theme2|[theme][other]theme2|[theme1][other]theme1|[theme2][other]theme2|[themeStandalone][other]themeStandalone|")
				b.AssertFileContent("public/index.html", "i18n: project theme2")
				b.AssertFileContent("public/index.html", "partial ntheme: theme2")

				// Params only set in themes
				b.AssertFileContent("public/index.html", "v1::v12", "v2::v21")

			}},
		}

		themeConfigs = []struct {
			name   string
			config string
		}{
			{"themeStandalone", themeStandalone},
			{"themeCyclic", themeCyclic},
			{"theme1", theme1},
			{"theme2", theme2},
			{"theme3", theme3},
			{"theme4", theme4},
		}
	)

	for i, testConfig := range testConfigs {
		t.Log(fmt.Sprintf("Test %d", i))
		b := newTestSitesBuilder(t).WithLogger(loggers.NewErrorLogger())
		b.WithConfigFile("toml", testConfig.siteConfig)

		for _, tc := range themeConfigs {
			var variationsNameBase = []string{"nproject", "ntheme", tc.name}

			themeRoot := filepath.Join("themes", tc.name)
			b.WithSourceFile(filepath.Join(themeRoot, "config.toml"), tc.config)

			b.WithSourceFile(filepath.Join("layouts", "partials", "m.html"), `{{- range $k, $v := . }}{{ $k }}::{{ template "printv" $v }}
{{ end }}	
{{ define "printv" }}
{{- $tp := printf "%T" . -}}
{{- if (strings.HasSuffix $tp "map[string]interface {}") -}}
{{- range $k, $v := . }}[{{ $k }}]{{ template "printv" $v }}{{ end -}}
{{- else -}}
{{- . }}|
{{- end -}}
{{ end }}
`)

			for _, nameVariaton := range variationsNameBase {
				roots := []string{"", themeRoot}

				for _, root := range roots {
					name := tc.name
					if root == "" {
						name = "project"
					}

					if nameVariaton == "ntheme" && name == "project" {
						continue
					}

					// static
					b.WithSourceFile(filepath.Join(root, "static", nameVariaton+".txt"), name)

					// layouts
					if i == 1 {
						b.WithSourceFile(filepath.Join(root, "layouts", "partials", "theme2o.html"), "Not Set")
					}
					b.WithSourceFile(filepath.Join(root, "layouts", "partials", nameVariaton+".html"), name)
					if root != "" && testConfig.offset == tc.name {
						for _, tc2 := range themeConfigs {
							b.WithSourceFile(filepath.Join(root, "layouts", "partials", tc2.name+"o.html"), name)
						}
					}

					// i18n + data

					var dataContent string
					if root == "" {
						dataContent = fmt.Sprintf(`
[%s]
other = %q

[inner]
other = %q

`, name, name, name)
					} else {
						dataContent = fmt.Sprintf(`
[%s]
other = %q

[inner]
other = %q

[theme]
other = %q

`, name, name, name, name)
					}

					b.WithSourceFile(filepath.Join(root, "data", nameVariaton+".toml"), dataContent)
					b.WithSourceFile(filepath.Join(root, "i18n", "en.toml"), dataContent)

					// If an offset is set, duplicate a data key with a winner in the middle.
					if root != "" && testConfig.offset == tc.name {
						for _, tc2 := range themeConfigs {
							dataContent := fmt.Sprintf(`
[offset]
v = %q
`, tc.name)
							b.WithSourceFile(filepath.Join(root, "data", tc2.name+"o.toml"), dataContent)
						}
					}
				}

			}

		}

		for _, themeConfig := range themeConfigs {
			b.WithSourceFile(filepath.Join("themes", "config.toml"), themeConfig.config)
		}

		b.WithContent(filepath.Join("content", "page.md"), `---
title: "Page"
---

`)

		homeTpl := `
data: {{ partial "m" .Site.Data }}
i18n: {{ i18n "inner" }} {{ i18n "theme" }}
partial ntheme: {{ partial "ntheme" . }}
partial theme2o: {{ partial "theme2o" . }}
params: {{ partial "m" .Site.Params }} 
		
`

		b.WithTemplates(filepath.Join("layouts", "home.html"), homeTpl)

		b.Build(BuildCfg{})

		var _ = os.Stdout

		//	printFs(b.H.Deps.BaseFs.LayoutsFs, "", os.Stdout)
		testConfig.check(b)

	}

}
