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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/viper"
)

func TestTemplateLookupOrder(t *testing.T) {
	t.Parallel()
	var (
		fs  *hugofs.Fs
		cfg *viper.Viper
		th  testHelper
	)

	// Variants base templates:
	//   1. <current-path>/<template-name>-baseof.<suffix>, e.g. list-baseof.<suffix>.
	//   2. <current-path>/baseof.<suffix>
	//   3. _default/<template-name>-baseof.<suffix>, e.g. list-baseof.<suffix>.
	//   4. _default/baseof.<suffix>
	for i, this := range []struct {
		setup  func(t *testing.T)
		assert func(t *testing.T)
	}{
		{
			// Variant 1
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "section", "sect1-baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "section", "sect1.html"), `{{define "main"}}sect{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base: sect")
			},
		},
		{
			// Variant 2
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "index.html"), `{{define "main"}}index{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "index.html"), "Base: index")
			},
		},
		{
			// Variant 3
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "_default", "list-baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"), `{{define "main"}}list{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base: list")
			},
		},
		{
			// Variant 4
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "_default", "baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"), `{{define "main"}}list{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base: list")
			},
		},
		{
			// Variant 1, theme,  use project's base
			func(t *testing.T) {
				cfg.Set("theme", "mytheme")
				writeSource(t, fs, filepath.Join("layouts", "section", "sect1-baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "section", "sect-baseof.html"), `Base Theme: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "section", "sect1.html"), `{{define "main"}}sect{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base: sect")
			},
		},
		{
			// Variant 1, theme,  use theme's base
			func(t *testing.T) {
				cfg.Set("theme", "mytheme")
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "section", "sect1-baseof.html"), `Base Theme: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "section", "sect1.html"), `{{define "main"}}sect{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base Theme: sect")
			},
		},
		{
			// Variant 4, theme, use project's base
			func(t *testing.T) {
				cfg.Set("theme", "mytheme")
				writeSource(t, fs, filepath.Join("layouts", "_default", "baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "_default", "baseof.html"), `Base Theme: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "_default", "list.html"), `{{define "main"}}list{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base: list")
			},
		},
		{
			// Variant 4, theme, use themes's base
			func(t *testing.T) {
				cfg.Set("theme", "mytheme")
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "_default", "baseof.html"), `Base Theme: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "_default", "list.html"), `{{define "main"}}list{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base Theme: list")
			},
		},
		{
			// Test section list and single template selection.
			// Issue #3116
			func(t *testing.T) {
				cfg.Set("theme", "mytheme")

				writeSource(t, fs, filepath.Join("layouts", "_default", "baseof.html"), `Base: {{block "main" .}}block{{end}}`)

				// Both single and list template in /SECTION/
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "sect1", "list.html"), `sect list`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "_default", "list.html"), `default list`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "sect1", "single.html"), `sect single`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "_default", "single.html"), `default single`)

				// sect2 with list template in /section
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "section", "sect2.html"), `sect2 list`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "sect list")
				th.assertFileContent(filepath.Join("public", "sect1", "page1", "index.html"), "sect single")
				th.assertFileContent(filepath.Join("public", "sect2", "index.html"), "sect2 list")
			},
		},
		{
			// Test section list and single template selection with base template.
			// Issue #2995
			func(t *testing.T) {

				writeSource(t, fs, filepath.Join("layouts", "_default", "baseof.html"), `Base Default: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "sect1", "baseof.html"), `Base Sect1: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "section", "sect2-baseof.html"), `Base Sect2: {{block "main" .}}block{{end}}`)

				// Both single and list + base template in /SECTION/
				writeSource(t, fs, filepath.Join("layouts", "sect1", "list.html"), `{{define "main"}}sect1 list{{ end }}`)
				writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"), `{{define "main"}}default list{{ end }}`)
				writeSource(t, fs, filepath.Join("layouts", "sect1", "single.html"), `{{define "main"}}sect single{{ end }}`)
				writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `{{define "main"}}default single{{ end }}`)

				// sect2 with list template in /section
				writeSource(t, fs, filepath.Join("layouts", "section", "sect2.html"), `{{define "main"}}sect2 list{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Base Sect1", "sect1 list")
				th.assertFileContent(filepath.Join("public", "sect1", "page1", "index.html"), "Base Sect1", "sect single")
				th.assertFileContent(filepath.Join("public", "sect2", "index.html"), "Base Sect2", "sect2 list")

				// Note that this will get the default base template and not the one in /sect2 -- because there are no
				// single template defined in /sect2.
				th.assertFileContent(filepath.Join("public", "sect2", "page2", "index.html"), "Base Default", "default single")
			},
		},
	} {

		if i != 9 {
			continue
		}

		cfg, fs = newTestCfg()
		th = testHelper{cfg, fs, t}

		for i := 1; i <= 3; i++ {
			writeSource(t, fs, filepath.Join("content", fmt.Sprintf("sect%d", i), fmt.Sprintf("page%d.md", i)), `---
title: Template test
---
Some content
`)
		}

		this.setup(t)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
		t.Log("Template Lookup test", i)
		this.assert(t)

	}
}
