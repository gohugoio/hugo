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
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/hugofs"

	"github.com/spf13/viper"
)

func TestBaseGoTemplate(t *testing.T) {
	t.Parallel()
	var (
		fs  *hugofs.Fs
		cfg *viper.Viper
		th  testHelper
	)

	// Variants:
	//   1. <current-path>/<template-name>-baseof.<suffix>, e.g. list-baseof.<suffix>.
	//   2. <current-path>/baseof.<suffix>
	//   3. _default/<template-name>-baseof.<suffix>, e.g. list-baseof.<suffix>.
	//   4. _default/baseof.<suffix>
	for _, this := range []struct {
		setup  func(t *testing.T)
		assert func(t *testing.T)
	}{
		{
			// Variant 1
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "section", "sect-baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "section", "sect.html"), `{{define "main"}}sect{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(t, fs, filepath.Join("public", "sect", "index.html"), false, "Base: sect")
			},
		},
		{
			// Variant 2
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "index.html"), `{{define "main"}}index{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(t, fs, filepath.Join("public", "index.html"), false, "Base: index")
			},
		},
		{
			// Variant 3
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "_default", "list-baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"), `{{define "main"}}list{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(t, fs, filepath.Join("public", "sect", "index.html"), false, "Base: list")
			},
		},
		{
			// Variant 4
			func(t *testing.T) {
				writeSource(t, fs, filepath.Join("layouts", "_default", "baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"), `{{define "main"}}list{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(t, fs, filepath.Join("public", "sect", "index.html"), false, "Base: list")
			},
		},
		{
			// Variant 1, theme,  use project's base
			func(t *testing.T) {
				cfg.Set("theme", "mytheme")
				writeSource(t, fs, filepath.Join("layouts", "section", "sect-baseof.html"), `Base: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "section", "sect-baseof.html"), `Base Theme: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "section", "sect.html"), `{{define "main"}}sect{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(t, fs, filepath.Join("public", "sect", "index.html"), false, "Base: sect")
			},
		},
		{
			// Variant 1, theme,  use theme's base
			func(t *testing.T) {
				cfg.Set("theme", "mytheme")
				writeSource(t, fs, filepath.Join("themes", "mytheme", "layouts", "section", "sect-baseof.html"), `Base Theme: {{block "main" .}}block{{end}}`)
				writeSource(t, fs, filepath.Join("layouts", "section", "sect.html"), `{{define "main"}}sect{{ end }}`)

			},
			func(t *testing.T) {
				th.assertFileContent(t, fs, filepath.Join("public", "sect", "index.html"), false, "Base Theme: sect")
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
				th.assertFileContent(t, fs, filepath.Join("public", "sect", "index.html"), false, "Base: list")
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
				th.assertFileContent(t, fs, filepath.Join("public", "sect", "index.html"), false, "Base Theme: list")
			},
		},
	} {

		cfg, fs = newTestCfg()
		th = testHelper{cfg}

		writeSource(t, fs, filepath.Join("content", "sect", "page.md"), `---
title: Template test
---
Some content
`)
		this.setup(t)

		buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

		this.assert(t)

	}
}
