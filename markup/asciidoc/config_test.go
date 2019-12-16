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

package asciidoc

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"

	qt "github.com/frankban/quicktest"
)

func TestDefaultConfig(t *testing.T) {
	c := qt.New(t)
	cfg, err := config.FromConfigString("", "toml")
	c.Assert(err, qt.IsNil)

	acfg, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)

	c.Assert(acfg.CurrentContent, qt.Equals, false)
	c.Assert(strings.Join(acfg.Args, " "), qt.Equals, "--no-header-footer --safe --trace")
}

func TestDecodeConfig(t *testing.T) {
	c := qt.New(t)
	tomlConfig := `
[asciidoctor]
args = ["--no-header-footer", "-r", "asciidoctor-html5s", "-b", "html5s", "-r", "asciidoctor-diagram"]
currentContent = true
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	acfg, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)

	c.Assert(acfg.CurrentContent, qt.Equals, true)
	c.Assert(strings.Join(acfg.Args, " "), qt.Equals, "--no-header-footer -r asciidoctor-html5s -b html5s -r asciidoctor-diagram")
}
