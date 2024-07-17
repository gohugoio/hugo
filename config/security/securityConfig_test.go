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

package security

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
)

func TestDecodeConfigFromTOML(t *testing.T) {
	c := qt.New(t)

	c.Run("Slice whitelist", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `


someOtherValue = "bar"

[security]
enableInlineShortcodes=true
[security.exec]
allow=["a", "b"]
osEnv=["a", "b", "c"]
[security.funcs]
getEnv=["a", "b"]

`

		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)

		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc, qt.Not(qt.IsNil))
		c.Assert(pc.EnableInlineShortcodes, qt.IsTrue)
		c.Assert(pc.Exec.Allow.Accept("a"), qt.IsTrue)
		c.Assert(pc.Exec.Allow.Accept("d"), qt.IsFalse)
		c.Assert(pc.Exec.OsEnv.Accept("a"), qt.IsTrue)
		c.Assert(pc.Exec.OsEnv.Accept("e"), qt.IsFalse)
		c.Assert(pc.Funcs.Getenv.Accept("a"), qt.IsTrue)
		c.Assert(pc.Funcs.Getenv.Accept("c"), qt.IsFalse)
	})

	c.Run("String whitelist", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `


someOtherValue = "bar"

[security]
[security.exec]
allow="a"
osEnv="b"

`

		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)

		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc, qt.Not(qt.IsNil))
		c.Assert(pc.Exec.Allow.Accept("a"), qt.IsTrue)
		c.Assert(pc.Exec.Allow.Accept("d"), qt.IsFalse)
		c.Assert(pc.Exec.OsEnv.Accept("b"), qt.IsTrue)
		c.Assert(pc.Exec.OsEnv.Accept("e"), qt.IsFalse)
	})

	c.Run("Default exec.osEnv", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `


someOtherValue = "bar"

[security]
[security.exec]
allow="a"

`

		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)

		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc, qt.Not(qt.IsNil))
		c.Assert(pc.Exec.Allow.Accept("a"), qt.IsTrue)
		c.Assert(pc.Exec.OsEnv.Accept("PATH"), qt.IsTrue)
		c.Assert(pc.Exec.OsEnv.Accept("e"), qt.IsFalse)
	})

	c.Run("Enable inline shortcodes, legacy", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `


someOtherValue = "bar"
enableInlineShortcodes=true

[security]
[security.exec]
allow="a"
osEnv="b"

`

		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)

		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc.EnableInlineShortcodes, qt.IsTrue)
	})
}

func TestToTOML(t *testing.T) {
	c := qt.New(t)

	got := DefaultConfig.ToTOML()

	c.Assert(got, qt.Equals,
		"[security]\n  enableInlineShortcodes = false\n\n  [security.exec]\n    allow = ['^(dart-)?sass(-embedded)?$', '^go$', '^git$', '^npx$', '^postcss$', '^tailwindcss$']\n    osEnv = ['(?i)^((HTTPS?|NO)_PROXY|PATH(EXT)?|APPDATA|TE?MP|TERM|GO\\w+|(XDG_CONFIG_)?HOME|USERPROFILE|SSH_AUTH_SOCK|DISPLAY|LANG|SYSTEMDRIVE)$']\n\n  [security.funcs]\n    getenv = ['^HUGO_', '^CI$']\n\n  [security.http]\n    methods = ['(?i)GET|POST']\n    urls = ['.*']",
	)
}

func TestDecodeConfigDefault(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	pc, err := DecodeConfig(config.New())
	c.Assert(err, qt.IsNil)
	c.Assert(pc, qt.Not(qt.IsNil))
	c.Assert(pc.Exec.Allow.Accept("a"), qt.IsFalse)
	c.Assert(pc.Exec.Allow.Accept("npx"), qt.IsTrue)
	c.Assert(pc.Exec.Allow.Accept("Npx"), qt.IsFalse)

	c.Assert(pc.HTTP.URLs.Accept("https://example.org"), qt.IsTrue)
	c.Assert(pc.HTTP.Methods.Accept("POST"), qt.IsTrue)
	c.Assert(pc.HTTP.Methods.Accept("GET"), qt.IsTrue)
	c.Assert(pc.HTTP.Methods.Accept("get"), qt.IsTrue)
	c.Assert(pc.HTTP.Methods.Accept("DELETE"), qt.IsFalse)
	c.Assert(pc.HTTP.MediaTypes.Accept("application/msword"), qt.IsFalse)

	c.Assert(pc.Exec.OsEnv.Accept("PATH"), qt.IsTrue)
	c.Assert(pc.Exec.OsEnv.Accept("GOROOT"), qt.IsTrue)
	c.Assert(pc.Exec.OsEnv.Accept("HOME"), qt.IsTrue)
	c.Assert(pc.Exec.OsEnv.Accept("SSH_AUTH_SOCK"), qt.IsTrue)
	c.Assert(pc.Exec.OsEnv.Accept("a"), qt.IsFalse)
	c.Assert(pc.Exec.OsEnv.Accept("e"), qt.IsFalse)
	c.Assert(pc.Exec.OsEnv.Accept("MYSECRET"), qt.IsFalse)
}
