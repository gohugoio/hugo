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
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/htesting"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigDir(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	configContent := `
baseURL = "https://example.org"
paginagePath = "pag_root"

[languages.en]
weight = 0
languageName = "English"

[languages.no]
weight = 10
languageName = "FOO"

[params]
p1 = "p1_base"

`

	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "hugo.toml", configContent)

	fb := htesting.NewTestdataBuilder(mm, "config/_default", t)

	fb.Add("config.toml", `paginatePath = "pag_default"`)

	fb.Add("params.yaml", `
p2: "p2params_default"
p3: "p3params_default"
p4: "p4params_default"
`)
	fb.Add("menus.toml", `
[[docs]]
name = "About Hugo"
weight = 1
[[docs]]
name = "Home"
weight = 2
	`)

	fb.Add("menus.no.toml", `
	[[docs]]
	name = "Om Hugo"
	weight = 1
	`)

	fb.Add("params.no.toml",
		`
p3 = "p3params_no_default"
p4 = "p4params_no_default"`,
	)
	fb.Add("languages.no.toml", `languageName = "Norsk_no_default"`)

	fb.Build()

	fb = fb.WithWorkingDir("config/production")

	fb.Add("config.toml", `paginatePath = "pag_production"`)

	fb.Add("params.no.toml", `
p2 = "p2params_no_production"
p3 = "p3params_no_production"
`)

	fb.Build()

	fb = fb.WithWorkingDir("config/development")

	// This is set in all the config.toml variants above, but this will win.
	fb.Add("config.toml", `paginatePath = "pag_development"`)

	fb.Add("params.no.toml", `p3 = "p3params_no_development"`)
	fb.Add("params.toml", `p3 = "p3params_development"`)

	fb.Build()

	cfg, _, err := LoadConfig(ConfigSourceDescriptor{Fs: mm, Environment: "development", Filename: "hugo.toml", AbsConfigDir: "config"})
	assert.NoError(err)

	assert.Equal("pag_development", cfg.GetString("paginatePath")) // /config/development/config.toml

	assert.Equal(10, cfg.GetInt("languages.no.weight"))                          //  /config.toml
	assert.Equal("Norsk_no_default", cfg.GetString("languages.no.languageName")) // /config/_default/languages.no.toml

	assert.Equal("p1_base", cfg.GetString("params.p1"))
	assert.Equal("p2params_default", cfg.GetString("params.p2")) // Is in both _default and production
	assert.Equal("p3params_development", cfg.GetString("params.p3"))
	assert.Equal("p3params_no_development", cfg.GetString("languages.no.params.p3"))

	assert.Equal(2, len(cfg.Get("menus.docs").(([]map[string]interface{}))))
	noMenus := cfg.Get("languages.no.menus.docs")
	assert.NotNil(noMenus)
	assert.Equal(1, len(noMenus.(([]map[string]interface{}))))

}

func TestLoadConfigDirError(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	configContent := `
baseURL = "https://example.org"

`

	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "hugo.toml", configContent)

	fb := htesting.NewTestdataBuilder(mm, "config/development", t)

	fb.Add("config.toml", `invalid & syntax`).Build()

	_, _, err := LoadConfig(ConfigSourceDescriptor{Fs: mm, Environment: "development", Filename: "hugo.toml", AbsConfigDir: "config"})
	assert.Error(err)

	fe := herrors.UnwrapErrorWithFileContext(err)
	assert.NotNil(fe)
	assert.Equal(filepath.FromSlash("config/development/config.toml"), fe.Position().Filename)

}
