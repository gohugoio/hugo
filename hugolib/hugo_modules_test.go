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

package hugolib

import (
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugoTestHelpers/testmodBuilder/mods"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestHugoModules(t *testing.T) {
	t.Parallel()

	testmods := mods.CreateModules().Collect()

	for _, m := range testmods {
		if len(m.Paths()) == 0 {
			continue
		}

		t.Run(m.Name(), func(t *testing.T) {
			assert := require.New(t)

			v := viper.New()

			workingDir, clean, err := createTempDir("hugo-modules-test")
			assert.NoError(err)
			defer clean()

			v.Set("workingDir", workingDir)

			configTemplate := `
baseURL = "https://example.com"
title = "My Modular Site"
workingDir = %q
theme = %q

`

			config := fmt.Sprintf(configTemplate, workingDir, m.Path())

			b := newTestSitesBuilder(t)

			// Need to use OS fs for this.
			b.Fs = hugofs.NewDefault(v)

			b.WithWorkingDir(workingDir).WithConfigFile("toml", config)
			b.WithContent("page.md", `
---
title: "Foo"
---
`)
			b.WithTemplates("home.html", `

{{ $mod := .Site.Data.modinfo.module }}
Mod Name: {{ $mod.name }}
Mod Version: {{ $mod.version }}
----
{{ range $k, $v := .Site.Data.modinfo }}
- {{ $k }}: {{ range $kk, $vv := $v }}{{ $kk }}: {{ $vv }}|{{ end -}}
{{ end }}


`)
			b.WithSourceFile("go.mod", `
module github.com/gohugoio/tests/testHugoModules


`)

			b.Build(BuildCfg{})

			// Verify that go.mod is autopopulated with all the modules in config.toml.
			b.AssertFileContent("go.mod", m.Path())

			b.AssertFileContent("public/index.html",
				"Mod Name: "+m.Name(),
				"Mod Version: v1.4.0")

			b.AssertFileContent("public/index.html", createModMatchers(m, m.Vendor)...)

		})
	}

}

func createModMatchers(m *mods.Md, vendored bool) []string {
	// Child depdendencies are one behind.
	expectMinorVersion := 3

	if vendored {
		// Vendored modules are stuck at v1.1.0.
		expectMinorVersion = 1
	}

	expectVersion := fmt.Sprintf("v1.%d.0", expectMinorVersion)

	var matchers []string
	for _, mm := range m.Children {
		matchers = append(
			matchers,
			fmt.Sprintf("%s: name: %s|version: %s", mm.Name(), mm.Name(), expectVersion))

		matchers = append(matchers, createModMatchers(mm, vendored || mm.Vendor)...)
	}

	return matchers
}
