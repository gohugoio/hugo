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
		"[security]\n  enableInlineShortcodes = false\n\n  [security.exec]\n    allow = ['^(dart-)?sass(-embedded)?$', '^go$', '^git$', '^node$', '^postcss$', '^tailwindcss$']\n    osEnv = ['(?i)^((HTTPS?|NO)_PROXY|PATH(EXT)?|APPDATA|TE?MP|TERM|GO\\w+|(XDG_CONFIG_)?HOME|USERPROFILE|SSH_AUTH_SOCK|DISPLAY|LANG|SYSTEMDRIVE|PROGRAMDATA)$']\n\n  [security.funcs]\n    getenv = ['^HUGO_', '^CI$']\n\n  [security.http]\n    methods = ['(?i)GET|POST']\n    urls = ['(?i)^https?://[a-z]', '! (?i)localhost', '! @']\n\n  [security.node]\n    [security.node.permissions]\n      allowAddons = ['tailwindcss']\n      allowChildProcess = ['tailwindcss']\n      allowRead = ['.']\n      allowWorker = ['tailwindcss']\n      allowWrite = []\n      disable = false",
	)
}

func TestDecodeConfigDefault(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	pc, err := DecodeConfig(config.New())
	c.Assert(err, qt.IsNil)
	c.Assert(pc, qt.Not(qt.IsNil))
	c.Assert(pc.Exec.Allow.Accept("a"), qt.IsFalse)
	c.Assert(pc.Exec.Allow.Accept("node"), qt.IsTrue)
	c.Assert(pc.Exec.Allow.Accept("npx"), qt.IsFalse)

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

	c.Assert(pc.Node.Permissions.IsEnabled(), qt.IsTrue)
	c.Assert(pc.Node.Permissions.AllowRead, qt.DeepEquals, []string{"."})
	c.Assert(pc.Node.Permissions.AllowWrite, qt.DeepEquals, []string{})
	c.Assert(pc.Node.Permissions.AllowChildProcess, qt.DeepEquals, []string{"tailwindcss"})
}

func TestCheckAllowedHTTPURLHardenedDefaultsIssue14792(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	c.Run("Public URLs allowed by default", func(c *qt.C) {
		c.Parallel()
		pc, err := DecodeConfig(config.New())
		c.Assert(err, qt.IsNil)
		for _, u := range []string{
			"https://example.org/",
			"https://example.org:8443/foo",
			"https://sub.example.org/path",
		} {
			c.Assert(pc.CheckAllowedHTTPURL(u), qt.IsNil, qt.Commentf(u))
		}
	})

	c.Run("Private/loopback URLs denied by default", func(c *qt.C) {
		c.Parallel()
		pc, err := DecodeConfig(config.New())
		c.Assert(err, qt.IsNil)
		for _, u := range []string{
			"http://localhost/",
			"http://LOCALHOST:8080/",
			"http://foo.localhost/",
			"http://127.0.0.1/",
			"http://127.1.2.3:8080/x",
			"http://user:pass@127.0.0.1/", // userinfo must not sneak past the deny.
			"http://10.0.0.1/",
			"http://172.16.0.1/",
			"http://192.168.1.1/",
			"http://169.254.169.254/latest/meta-data/", // AWS/GCP metadata.
			"http://0.0.0.0/",
			"http://[::1]/",
			"http://[fe80::1]/",
			"http://[fc00::1]/",
			// Public IP literals are blocked as collateral; users can override.
			"http://93.184.216.34/",
			"https://[2001:db8::1]/",
		} {
			err := pc.CheckAllowedHTTPURL(u)
			c.Assert(err, qt.IsNotNil, qt.Commentf(u))
			c.Assert(err, qt.ErrorMatches, `(?s).*is not whitelisted in policy "security\.http\.urls".*`, qt.Commentf(u))
		}
	})

	c.Run("Explicit user config bypasses hardening", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `
[security.http]
urls = ['http://127\.0\.0\.1.*', 'http://localhost.*']
`
		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)
		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc.CheckAllowedHTTPURL("http://127.0.0.1:8080/foo"), qt.IsNil)
		c.Assert(pc.CheckAllowedHTTPURL("http://localhost:1313/"), qt.IsNil)
	})

	c.Run("User can deny with the ! prefix", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `
[security.http]
urls = ['.*', '! ^https?://evil\.example\.com']
`
		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)
		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc.CheckAllowedHTTPURL("https://good.example.com/"), qt.IsNil)
		err = pc.CheckAllowedHTTPURL("https://evil.example.com/x")
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*is not whitelisted in policy "security\.http\.urls".*`)
	})
}

func TestDecodeConfigNodePermissions(t *testing.T) {
	c := qt.New(t)

	c.Run("Custom paths", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `
[security.node.permissions]
allowRead = ["/tmp", "."]
allowWrite = ["."]
allowChildProcess = ["tailwindcss"]
`
		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)
		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc.Node.Permissions.IsEnabled(), qt.IsTrue)
		c.Assert(pc.Node.Permissions.AllowRead, qt.DeepEquals, []string{"/tmp", "."})
		c.Assert(pc.Node.Permissions.AllowWrite, qt.DeepEquals, []string{"."})
		c.Assert(pc.Node.Permissions.AllowChildProcess, qt.DeepEquals, []string{"tailwindcss"})
	})

	c.Run("Disabled", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `
[security.node.permissions]
disable = true
`
		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)
		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc.Node.Permissions.IsEnabled(), qt.IsFalse)
	})

	c.Run("Wildcard", func(c *qt.C) {
		c.Parallel()
		tomlConfig := `
[security.node.permissions]
allowRead = ["*"]
allowWrite = ["*"]
`
		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)
		pc, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(pc.Node.Permissions.IsEnabled(), qt.IsTrue)
		c.Assert(pc.Node.Permissions.AllowRead, qt.DeepEquals, []string{"*"})
	})
}
