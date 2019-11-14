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
	"testing"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/hugofs"

	qt "github.com/frankban/quicktest"
)

func TestClient(t *testing.T) {
	if hugo.GoMinorVersion() < 12 {
		// https://github.com/golang/go/issues/26794
		// There were some concurrent issues with Go modules in < Go 12.
		t.Skip("skip this for Go <= 1.11 due to a bug in Go's stdlib")
	}

	t.Parallel()

	modName := "hugo-modules-basic-test"
	modPath := "github.com/gohugoio/tests/" + modName
	modConfig := DefaultModuleConfig
	modConfig.Imports = []Import{Import{Path: "github.com/gohugoio/hugoTestModules1_darwin/modh2_2"}}

	c := qt.New(t)

	workingDir, clean, err := htesting.CreateTempDir(hugofs.Os, modName)
	c.Assert(err, qt.IsNil)
	defer clean()

	client := NewClient(ClientConfig{
		Fs:           hugofs.Os,
		WorkingDir:   workingDir,
		ModuleConfig: modConfig,
	})

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

	expect := `github.com/gohugoio/tests/hugo-modules-basic-test github.com/gohugoio/hugoTestModules1_darwin/modh2_2@v1.4.0
github.com/gohugoio/hugoTestModules1_darwin/modh2_2@v1.4.0 github.com/gohugoio/hugoTestModules1_darwin/modh2_2_1v@v1.3.0
github.com/gohugoio/hugoTestModules1_darwin/modh2_2@v1.4.0 github.com/gohugoio/hugoTestModules1_darwin/modh2_2_2@v1.3.0
`

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

	// Test the ignoreVendor setting
	clientIgnoreVendor := NewClient(ClientConfig{
		Fs:           hugofs.Os,
		WorkingDir:   workingDir,
		ModuleConfig: modConfig,
		IgnoreVendor: true,
	})

	graphb.Reset()
	c.Assert(clientIgnoreVendor.Graph(&graphb), qt.IsNil)
	c.Assert(graphb.String(), qt.Equals, expect)

	// Test Tidy
	c.Assert(client.Tidy(), qt.IsNil)

}

func TestGetModlineSplitter(t *testing.T) {

	c := qt.New(t)

	gomodSplitter := getModlineSplitter(true)

	c.Assert(gomodSplitter("\tgithub.com/BurntSushi/toml v0.3.1"), qt.DeepEquals, []string{"github.com/BurntSushi/toml", "v0.3.1"})
	c.Assert(gomodSplitter("\tgithub.com/cpuguy83/go-md2man v1.0.8 // indirect"), qt.DeepEquals, []string{"github.com/cpuguy83/go-md2man", "v1.0.8"})
	c.Assert(gomodSplitter("require ("), qt.IsNil)

	gosumSplitter := getModlineSplitter(false)
	c.Assert(gosumSplitter("github.com/BurntSushi/toml v0.3.1"), qt.DeepEquals, []string{"github.com/BurntSushi/toml", "v0.3.1"})

}
