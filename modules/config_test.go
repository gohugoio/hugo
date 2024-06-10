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
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/config"

	qt "github.com/frankban/quicktest"
)

func TestConfigHugoVersionIsValid(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		in     HugoVersion
		expect bool
	}{
		{HugoVersion{Min: "0.33.0"}, true},
		{HugoVersion{Min: "0.56.0-DEV"}, true},
		{HugoVersion{Min: "0.33.0", Max: "0.55.0"}, false},
		{HugoVersion{Min: "0.33.0", Max: "0.199.0"}, true},
	} {
		c.Assert(test.in.IsValid(), qt.Equals, test.expect, qt.Commentf("%#v", test.in))
	}
}

func TestDecodeConfig(t *testing.T) {
	c := qt.New(t)

	c.Run("Basic", func(c *qt.C) {
		tempDir := c.TempDir()
		tomlConfig := `
workingDir = %q
[module]
workspace = "hugo.work"
[module.hugoVersion]
min = "0.54.2"
max = "0.199.0"
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

		hugoWorkFilename := filepath.Join(tempDir, "hugo.work")
		f, _ := os.Create(hugoWorkFilename)
		f.Close()
		cfg, err := config.FromConfigString(fmt.Sprintf(tomlConfig, tempDir), "toml")
		c.Assert(err, qt.IsNil)

		mcfg, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)

		v056 := hugo.VersionString("0.56.0")

		hv := mcfg.HugoVersion

		c.Assert(v056.Compare(hv.Min), qt.Equals, -1)
		c.Assert(v056.Compare(hv.Max), qt.Equals, 1)
		c.Assert(hv.Extended, qt.Equals, true)

		if hugo.IsExtended {
			c.Assert(hv.IsValid(), qt.Equals, true)
		}

		c.Assert(mcfg.Workspace, qt.Equals, hugoWorkFilename)

		c.Assert(len(mcfg.Mounts), qt.Equals, 1)
		c.Assert(len(mcfg.Imports), qt.Equals, 1)
		imp := mcfg.Imports[0]
		imp.Path = "github.com/bep/mycomponent"
		c.Assert(imp.Mounts[1].Source, qt.Equals, "src/markdown/blog")
		c.Assert(imp.Mounts[1].Target, qt.Equals, "content/blog")
		c.Assert(imp.Mounts[1].Lang, qt.Equals, "en")
	})

	c.Run("Replacements", func(c *qt.C) {
		for _, tomlConfig := range []string{`
[module]
replacements="a->b,github.com/bep/mycomponent->c"
[[module.imports]]
path="github.com/bep/mycomponent"
`, `
[module]
replacements=["a->b","github.com/bep/mycomponent->c"]
[[module.imports]]
path="github.com/bep/mycomponent"
`} {

			cfg, err := config.FromConfigString(tomlConfig, "toml")
			c.Assert(err, qt.IsNil)

			mcfg, err := DecodeConfig(cfg)
			c.Assert(err, qt.IsNil)
			c.Assert(mcfg.Replacements, qt.DeepEquals, []string{"a->b", "github.com/bep/mycomponent->c"})
			c.Assert(mcfg.replacementsMap, qt.DeepEquals, map[string]string{
				"a":                          "b",
				"github.com/bep/mycomponent": "c",
			})

			c.Assert(mcfg.Imports[0].Path, qt.Equals, "c")

		}
	})
}

func TestDecodeConfigBothOldAndNewProvided(t *testing.T) {
	c := qt.New(t)
	tomlConfig := `

theme = ["b", "c"]

[module]
[[module.imports]]
path="a"

`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	modCfg, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(len(modCfg.Imports), qt.Equals, 3)
	c.Assert(modCfg.Imports[0].Path, qt.Equals, "a")
}

// Test old style theme import.
func TestDecodeConfigTheme(t *testing.T) {
	c := qt.New(t)
	tomlConfig := `

theme = ["a", "b"]
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	mcfg, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)

	c.Assert(len(mcfg.Imports), qt.Equals, 2)
	c.Assert(mcfg.Imports[0].Path, qt.Equals, "a")
	c.Assert(mcfg.Imports[1].Path, qt.Equals, "b")
}
