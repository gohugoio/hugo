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

package modules

import (
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/config"

	"github.com/stretchr/testify/require"
)

func TestConfigHugoVersionIsValid(t *testing.T) {
	assert := require.New(t)

	for i, test := range []struct {
		in     HugoVersion
		expect bool
	}{
		{HugoVersion{Min: "0.33.0"}, true},
		{HugoVersion{Min: "0.56.0-DEV"}, true},
		{HugoVersion{Min: "0.33.0", Max: "0.55.0"}, false},
		{HugoVersion{Min: "0.33.0", Max: "0.99.0"}, true},
	} {
		assert.Equal(test.expect, test.in.IsValid(), fmt.Sprintf("test %d", i))
	}
}

func TestDecodeConfig(t *testing.T) {
	assert := require.New(t)
	tomlConfig := `
[module]

[module.hugoVersion]
min = "0.54.2"
max = "0.99.0"
extended = true

[[module.mounts]]
source="src/project/blog"
target="content/blog"
lang="en"
[[module.imports]]
path="github.com/bep/mycomponent"
[[module.imports.mounts]]
source="scss"
target="assets/bootstrap/scss"
[[module.imports.mounts]]
source="src/markdown/blog"
target="content/blog"
lang="en"
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	assert.NoError(err)

	mcfg, err := DecodeConfig(cfg)
	assert.NoError(err)

	v056 := hugo.VersionString("0.56.0")

	hv := mcfg.HugoVersion

	assert.Equal(-1, v056.Compare(hv.Min))
	assert.Equal(1, v056.Compare(hv.Max))
	assert.True(hv.Extended)

	if hugo.IsExtended {
		assert.True(hv.IsValid())
	}

	assert.Len(mcfg.Mounts, 1)
	assert.Len(mcfg.Imports, 1)
	imp := mcfg.Imports[0]
	imp.Path = "github.com/bep/mycomponent"
	assert.Equal("src/markdown/blog", imp.Mounts[1].Source)
	assert.Equal("content/blog", imp.Mounts[1].Target)
	assert.Equal("en", imp.Mounts[1].Lang)

}

func TestDecodeConfigBothOldAndNewProvided(t *testing.T) {
	assert := require.New(t)
	tomlConfig := `

theme = ["b", "c"]

[module]
[[module.imports]]
path="a"

`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	assert.NoError(err)

	modCfg, err := DecodeConfig(cfg)
	assert.NoError(err)
	assert.Len(modCfg.Imports, 3)
	assert.Equal("a", modCfg.Imports[0].Path)

}

// Test old style theme import.
func TestDecodeConfigTheme(t *testing.T) {
	assert := require.New(t)
	tomlConfig := `

theme = ["a", "b"]
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	assert.NoError(err)

	mcfg, err := DecodeConfig(cfg)
	assert.NoError(err)

	assert.Len(mcfg.Imports, 2)
	assert.Equal("a", mcfg.Imports[0].Path)
	assert.Equal("b", mcfg.Imports[1].Path)
}
