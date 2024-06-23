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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/hugofs"

	qt "github.com/frankban/quicktest"
)

func TestClient(t *testing.T) {
	modName := "hugo-modules-basic-test"
	modPath := "github.com/gohugoio/tests/" + modName
	defaultImport := "modh2_2"
	expect := `github.com/gohugoio/tests/hugo-modules-basic-test github.com/gohugoio/hugoTestModules1_darwin/modh2_2@v1.4.0
github.com/gohugoio/hugoTestModules1_darwin/modh2_2@v1.4.0 github.com/gohugoio/hugoTestModules1_darwin/modh2_2_1v@v1.3.0
github.com/gohugoio/hugoTestModules1_darwin/modh2_2@v1.4.0 github.com/gohugoio/hugoTestModules1_darwin/modh2_2_2@v1.3.0
`

	c := qt.New(t)
	var clientID uint64 // we increment this to get each test in its own directory.

	newClient := func(c *qt.C, withConfig func(cfg *ClientConfig), imp string) (*Client, func()) {
		atomic.AddUint64(&clientID, uint64(1))
		workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, fmt.Sprintf("%s-%d", modName, clientID))
		c.Assert(err, qt.IsNil)
		themesDir := filepath.Join(workingDir, "themes")
		err = os.Mkdir(themesDir, 0o777)
		c.Assert(err, qt.IsNil)
		publishDir := filepath.Join(workingDir, "public")
		err = os.Mkdir(publishDir, 0o777)
		c.Assert(err, qt.IsNil)

		ccfg := ClientConfig{
			Fs:         hugofs.Os,
			CacheDir:   filepath.Join(workingDir, "modcache"),
			WorkingDir: workingDir,
			ThemesDir:  themesDir,
			PublishDir: publishDir,
			Exec:       hexec.New(security.DefaultConfig, ""),
		}

		withConfig(&ccfg)
		ccfg.ModuleConfig.Imports = []Import{{Path: "github.com/gohugoio/hugoTestModules1_darwin/" + imp}}
		client := NewClient(ccfg)

		return client, clean
	}

	c.Run("All", func(c *qt.C) {
		client, clean := newClient(c, func(cfg *ClientConfig) {
			cfg.ModuleConfig = DefaultModuleConfig
		}, defaultImport)
		defer clean()

		// Test Init
		c.Assert(client.Init(modPath), qt.IsNil)

		// Test Collect
		mc, err := client.Collect()
		c.Assert(err, qt.IsNil)
		c.Assert(len(mc.AllModules), qt.Equals, 4)
		for _, m := range mc.AllModules {
			c.Assert(m, qt.Not(qt.IsNil))
		}

		// Test Graph
		var graphb bytes.Buffer
		c.Assert(client.Graph(&graphb), qt.IsNil)

		c.Assert(graphb.String(), qt.Equals, expect)

		// Test Vendor
		c.Assert(client.Vendor(), qt.IsNil)
		graphb.Reset()
		c.Assert(client.Graph(&graphb), qt.IsNil)

		expectVendored := `project github.com/gohugoio/hugoTestModules1_darwin/modh2_2@v1.4.0+vendor
project github.com/gohugoio/hugoTestModules1_darwin/modh2_2_1v@v1.3.0+vendor
project github.com/gohugoio/hugoTestModules1_darwin/modh2_2_2@v1.3.0+vendor
`

		c.Assert(graphb.String(), qt.Equals, expectVendored)

		// Test Tidy
		c.Assert(client.Tidy(), qt.IsNil)
	})

	c.Run("IgnoreVendor", func(c *qt.C) {
		client, clean := newClient(
			c, func(cfg *ClientConfig) {
				cfg.ModuleConfig = DefaultModuleConfig
				cfg.IgnoreVendor = globAll
			}, defaultImport)
		defer clean()

		c.Assert(client.Init(modPath), qt.IsNil)
		_, err := client.Collect()
		c.Assert(err, qt.IsNil)
		c.Assert(client.Vendor(), qt.IsNil)

		var graphb bytes.Buffer
		c.Assert(client.Graph(&graphb), qt.IsNil)
		c.Assert(graphb.String(), qt.Equals, expect)
	})

	c.Run("NoVendor", func(c *qt.C) {
		mcfg := DefaultModuleConfig
		mcfg.NoVendor = "**"
		client, clean := newClient(
			c, func(cfg *ClientConfig) {
				cfg.ModuleConfig = mcfg
			}, defaultImport)
		defer clean()

		c.Assert(client.Init(modPath), qt.IsNil)
		_, err := client.Collect()
		c.Assert(err, qt.IsNil)
		c.Assert(client.Vendor(), qt.IsNil)

		var graphb bytes.Buffer
		c.Assert(client.Graph(&graphb), qt.IsNil)
		c.Assert(graphb.String(), qt.Equals, expect)
	})

	c.Run("VendorClosest", func(c *qt.C) {
		mcfg := DefaultModuleConfig
		mcfg.VendorClosest = true

		client, clean := newClient(
			c, func(cfg *ClientConfig) {
				cfg.ModuleConfig = mcfg
				s := "github.com/gohugoio/hugoTestModules1_darwin/modh1_1v"
				g, _ := glob.GetGlob(s)
				cfg.IgnoreVendor = g
			}, "modh1v")
		defer clean()

		c.Assert(client.Init(modPath), qt.IsNil)
		_, err := client.Collect()
		c.Assert(err, qt.IsNil)
		c.Assert(client.Vendor(), qt.IsNil)

		var graphb bytes.Buffer
		c.Assert(client.Graph(&graphb), qt.IsNil)

		c.Assert(graphb.String(), qt.Contains, "github.com/gohugoio/hugoTestModules1_darwin/modh1_1v@v1.3.0 github.com/gohugoio/hugoTestModules1_darwin/modh1_1_1v@v1.1.0+vendor")
	})

	// https://github.com/gohugoio/hugo/issues/7908
	c.Run("createThemeDirname", func(c *qt.C) {
		mcfg := DefaultModuleConfig
		client, clean := newClient(
			c, func(cfg *ClientConfig) {
				cfg.ModuleConfig = mcfg
			}, defaultImport)
		defer clean()

		dirname, err := client.createThemeDirname("foo", false)
		c.Assert(err, qt.IsNil)
		c.Assert(dirname, qt.Equals, filepath.Join(client.ccfg.ThemesDir, "foo"))

		dirname, err = client.createThemeDirname("../../foo", true)
		c.Assert(err, qt.IsNil)
		c.Assert(dirname, qt.Equals, filepath.Join(client.ccfg.ThemesDir, "../../foo"))

		_, err = client.createThemeDirname("../../foo", false)
		c.Assert(err, qt.Not(qt.IsNil))

		absDir := filepath.Join(client.ccfg.WorkingDir, "..", "..")
		dirname, err = client.createThemeDirname(absDir, true)
		c.Assert(err, qt.IsNil)
		c.Assert(dirname, qt.Equals, absDir)
		dirname, err = client.createThemeDirname(absDir, false)
		fmt.Println(dirname)
		c.Assert(err, qt.Not(qt.IsNil))
	})
}

var globAll, _ = glob.GetGlob("**")

func TestGetModlineSplitter(t *testing.T) {
	c := qt.New(t)

	gomodSplitter := getModlineSplitter(true)

	c.Assert(gomodSplitter("\tgithub.com/BurntSushi/toml v0.3.1"), qt.DeepEquals, []string{"github.com/BurntSushi/toml", "v0.3.1"})
	c.Assert(gomodSplitter("\tgithub.com/cpuguy83/go-md2man v1.0.8 // indirect"), qt.DeepEquals, []string{"github.com/cpuguy83/go-md2man", "v1.0.8"})
	c.Assert(gomodSplitter("require ("), qt.IsNil)

	gosumSplitter := getModlineSplitter(false)
	c.Assert(gosumSplitter("github.com/BurntSushi/toml v0.3.1"), qt.DeepEquals, []string{"github.com/BurntSushi/toml", "v0.3.1"})
}
