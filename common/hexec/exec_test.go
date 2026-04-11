// Copyright 2026 The Hugo Authors. All rights reserved.
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

package hexec

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config/security"
)

func TestNodePermissionArgs(t *testing.T) {
	c := qt.New(t)

	// Use t.TempDir() so paths are absolute on any OS (avoids Windows volume assumptions).
	base := t.TempDir()
	site := filepath.Join(base, "site")
	tmp := filepath.Join(base, "tmp")
	cacheDir := filepath.Join(base, "home", "user", ".cache", "hugo_cache", "modules")

	c.Run("Default config tailwindcss", func(c *qt.C) {
		e := &Exec{
			sc:         security.DefaultConfig,
			workingDir: site,
		}
		args := e.nodePermissionArgs("tailwindcss", "")
		c.Assert(args, qt.DeepEquals, []string{
			"--permission",
			"--allow-fs-read=" + site,
			"--allow-addons",
			"--allow-worker",
			"--disable-warning=SecurityWarning",
		})
	})

	c.Run("Default config postcss", func(c *qt.C) {
		e := &Exec{
			sc:         security.DefaultConfig,
			workingDir: site,
		}
		args := e.nodePermissionArgs("postcss", "")
		c.Assert(args, qt.DeepEquals, []string{
			"--permission",
			"--allow-fs-read=" + site,
		})
	})

	c.Run("Multiple paths", func(c *qt.C) {
		cfg := security.DefaultConfig
		cfg.Node.Permissions.AllowRead = []string{".", tmp}
		cfg.Node.Permissions.AllowWrite = []string{"."}
		e := &Exec{
			sc:         cfg,
			workingDir: site,
		}
		args := e.nodePermissionArgs("tailwindcss", "")
		c.Assert(args, qt.DeepEquals, []string{
			"--permission",
			"--allow-fs-read=" + site,
			"--allow-fs-read=" + tmp,
			"--allow-fs-write=" + site,
			"--allow-addons",
			"--allow-worker",
			"--disable-warning=SecurityWarning",
		})
	})

	c.Run("Wildcard", func(c *qt.C) {
		cfg := security.DefaultConfig
		cfg.Node.Permissions.AllowRead = []string{"*"}
		cfg.Node.Permissions.AllowWrite = []string{"*"}
		e := &Exec{
			sc:         cfg,
			workingDir: site,
		}
		args := e.nodePermissionArgs("tailwindcss", "")
		c.Assert(args, qt.DeepEquals, []string{
			"--permission",
			"--allow-fs-read=*",
			"--allow-fs-write=*",
			"--allow-addons",
			"--allow-worker",
			"--disable-warning=SecurityWarning",
		})
	})

	c.Run("Disabled", func(c *qt.C) {
		cfg := security.DefaultConfig
		cfg.Node.Permissions.Disable = true
		e := &Exec{
			sc:         cfg,
			workingDir: site,
		}
		args := e.nodePermissionArgs("tailwindcss", "")
		c.Assert(args, qt.IsNil)
	})

	c.Run("No fs flags", func(c *qt.C) {
		cfg := security.DefaultConfig
		cfg.Node.Permissions.AllowRead = nil
		cfg.Node.Permissions.AllowAddons = nil
		cfg.Node.Permissions.AllowWorker = nil
		e := &Exec{
			sc:         cfg,
			workingDir: site,
		}
		args := e.nodePermissionArgs("postcss", "")
		c.Assert(args, qt.DeepEquals, []string{"--permission"})
	})

	c.Run("Read only", func(c *qt.C) {
		cfg := security.DefaultConfig
		cfg.Node.Permissions.AllowRead = []string{"."}
		cfg.Node.Permissions.AllowWrite = nil
		e := &Exec{
			sc:         cfg,
			workingDir: site,
		}
		args := e.nodePermissionArgs("postcss", "")
		c.Assert(args, qt.DeepEquals, []string{
			"--permission",
			"--allow-fs-read=" + site,
		})
	})

	c.Run("With additional read paths", func(c *qt.C) {
		e := &Exec{
			sc:            security.DefaultConfig,
			workingDir:    site,
			nodeReadPaths: []string{cacheDir},
		}
		args := e.nodePermissionArgs("postcss", "")
		c.Assert(args, qt.DeepEquals, []string{
			"--permission",
			"--allow-fs-read=" + site,
			"--allow-fs-read=" + cacheDir,
		})
	})

	c.Run("Global install script path", func(c *qt.C) {
		e := &Exec{
			sc:         security.DefaultConfig,
			workingDir: site,
		}
		globalNM := filepath.Join(base, "nvm", "lib", "node_modules")
		script := filepath.Join(globalNM, "postcss-cli", "bin", "postcss")
		args := e.nodePermissionArgs("postcss", script)
		c.Assert(args, qt.DeepEquals, []string{
			"--permission",
			"--allow-fs-read=" + site,
			"--allow-fs-read=" + globalNM,
		})
	})
}

func TestNodeScriptReadPath(t *testing.T) {
	c := qt.New(t)

	base := t.TempDir()
	nm := filepath.Join(base, "node_modules")
	globalNM := filepath.Join(base, "nvm", "lib", "node_modules")

	c.Assert(nodeScriptReadPath(""), qt.Equals, "")
	c.Assert(nodeScriptReadPath(filepath.Join(nm, "postcss-cli", "index.js")), qt.Equals, nm)
	c.Assert(nodeScriptReadPath(filepath.Join(nm, "@babel", "cli", "bin", "babel.js")), qt.Equals, nm)
	c.Assert(nodeScriptReadPath(filepath.Join(globalNM, "postcss-cli", "bin", "postcss")), qt.Equals, globalNM)

	loose := filepath.Join(base, "tools", "script.js")
	c.Assert(nodeScriptReadPath(loose), qt.Equals, filepath.Dir(loose))
}

func TestResolveNodeBin(t *testing.T) {
	c := qt.New(t)

	// Create a fake node_modules structure.
	dir := t.TempDir()
	nodeModules := filepath.Join(dir, "node_modules")
	binDir := filepath.Join(nodeModules, ".bin")

	// Create target JS files.
	postcssJS := filepath.Join(nodeModules, "postcss-cli", "index.js")
	babelJS := filepath.Join(nodeModules, "@babel", "cli", "bin", "babel.js")
	mkdirAndWrite(t, postcssJS, "#!/usr/bin/env node\nconsole.log('postcss');\n")
	mkdirAndWrite(t, babelJS, "#!/usr/bin/env node\nconsole.log('babel');\n")
	os.MkdirAll(binDir, 0o755)

	c.Run("Symlink to JS file", func(c *qt.C) {
		if runtime.GOOS == "windows" {
			c.Skip("Symlinks may require elevated privileges on Windows")
		}
		link := filepath.Join(binDir, "postcss-link")
		os.Remove(link)
		c.Assert(os.Symlink(postcssJS, link), qt.IsNil)

		resolved := resolveNodeBin(link)
		t.Logf("Symlink: link=%q, resolved=%q", link, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, postcssJS), qt.IsTrue)
	})

	c.Run("Symlink to scoped package", func(c *qt.C) {
		if runtime.GOOS == "windows" {
			c.Skip("Symlinks may require elevated privileges on Windows")
		}
		link := filepath.Join(binDir, "babel-link")
		os.Remove(link)
		c.Assert(os.Symlink(babelJS, link), qt.IsNil)

		resolved := resolveNodeBin(link)
		t.Logf("Scoped symlink: link=%q, resolved=%q", link, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, babelJS), qt.IsTrue)
	})

	c.Run("Shell wrapper", func(c *qt.C) {
		wrapper := filepath.Join(binDir, "postcss-sh")
		content := "#!/bin/sh\n" +
			`basedir=$(dirname "$(echo "$0" | sed -e 's,\\,/,g')")` + "\n" +
			`exec node  "$basedir/../postcss-cli/index.js" "$@"` + "\n"
		mkdirAndWrite(t, wrapper, content)

		resolved := resolveNodeBin(wrapper)
		t.Logf("Shell wrapper: wrapper=%q, resolved=%q", wrapper, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, postcssJS), qt.IsTrue)
	})

	c.Run("Cmd wrapper", func(c *qt.C) {
		wrapper := filepath.Join(binDir, "postcss.cmd")
		content := "@ECHO off\r\n" +
			"SETLOCAL\r\n" +
			`endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "%_prog%"  "%dp0%\..\postcss-cli\index.js" %*` + "\r\n"
		mkdirAndWrite(t, wrapper, content)

		resolved := resolveNodeBin(wrapper)
		t.Logf("Cmd wrapper: wrapper=%q, resolved=%q", wrapper, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, postcssJS), qt.IsTrue)
	})

	c.Run("Cmd wrapper scoped package", func(c *qt.C) {
		wrapper := filepath.Join(binDir, "babel.cmd")
		content := "@ECHO off\r\n" +
			`endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "%_prog%"  "%dp0%\..\@babel\cli\bin\babel.js" %*` + "\r\n"
		mkdirAndWrite(t, wrapper, content)

		resolved := resolveNodeBin(wrapper)
		t.Logf("Cmd wrapper (scoped): wrapper=%q, resolved=%q", wrapper, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, babelJS), qt.IsTrue)
	})

	c.Run("Node script with shebang", func(c *qt.C) {
		script := filepath.Join(binDir, "node-global")
		mkdirAndWrite(t, script, "#!/usr/bin/env node\nconsole.log('global');\n")

		resolved := resolveNodeBin(script)
		t.Logf("Node script: path=%q, resolved=%q", script, resolved)
		c.Assert(resolved, qt.Equals, script)
	})

	c.Run("Native binary", func(c *qt.C) {
		native := filepath.Join(binDir, "native-tool")
		mkdirAndWrite(t, native, "\x7fELF\x00\x00\x00")

		resolved := resolveNodeBin(native)
		t.Logf("Native binary: path=%q, resolved=%q", native, resolved)
		c.Assert(resolved, qt.Equals, "")
	})

	c.Run("Nonexistent file", func(c *qt.C) {
		c.Assert(resolveNodeBin("/nonexistent/path"), qt.Equals, "")
	})

	c.Run("Wrapper with missing target", func(c *qt.C) {
		wrapper := filepath.Join(binDir, "missing-target")
		mkdirAndWrite(t, wrapper, "#!/bin/sh\nexec node \"$basedir/../no-such-pkg/index.js\" \"$@\"\n")

		resolved := resolveNodeBin(wrapper)
		t.Logf("Missing target: wrapper=%q, resolved=%q", wrapper, resolved)
		c.Assert(resolved, qt.Equals, "")
	})
}

func TestExtractNodeEntryPointRegex(t *testing.T) {
	c := qt.New(t)

	cases := []struct {
		name    string
		content string
		want    string // expected capture group (with original separators)
	}{
		{"shell postcss", `"$basedir/../postcss-cli/index.js"`, "../postcss-cli/index.js"},
		{"shell babel", `"$basedir/../@babel/cli/bin/babel.js"`, "../@babel/cli/bin/babel.js"},
		{"shell tailwind mjs", `"$basedir/../@tailwindcss/cli/dist/index.mjs"`, "../@tailwindcss/cli/dist/index.mjs"},
		{"cmd postcss", `"%dp0%\..\postcss-cli\index.js"`, `..\postcss-cli\index.js`},
		{"cmd babel", `"%dp0%\..\@babel\cli\bin\babel.js"`, `..\@babel\cli\bin\babel.js`},
		{"cmd postcss no ext", `"%dp0%\..\postcss-cli\bin\postcss"`, `..\postcss-cli\bin\postcss`},
		{"shell postcss no ext", `"$basedir/../postcss-cli/bin/postcss"`, "../postcss-cli/bin/postcss"},
		{"cmd postcss global", `"%dp0%\node_modules\postcss-cli\index.js"`, `node_modules\postcss-cli\index.js`},
		{"cmd babel global", `"%dp0%\node_modules\@babel\cli\bin\babel.js"`, `node_modules\@babel\cli\bin\babel.js`},
		{"shell postcss global", `"$basedir/node_modules/postcss-cli/index.js"`, "node_modules/postcss-cli/index.js"},
	}

	for _, tc := range cases {
		c.Run(tc.name, func(c *qt.C) {
			m := nodeEntryPointRe.FindStringSubmatch(tc.content)
			t.Logf("regex match for %q: %v", tc.name, m)
			c.Assert(m, qt.Not(qt.IsNil))
			c.Assert(m[1], qt.Equals, tc.want)
		})
	}

	c.Run("No match", func(c *qt.C) {
		for _, s := range []string{"@ECHO off", "#!/bin/bash\necho hello", "\x7fELF"} {
			c.Assert(nodeEntryPointRe.FindStringSubmatch(s), qt.IsNil)
		}
	})
}

// TestResolveNodeBinWindows tests wrapper resolution on all platforms
// by simulating Windows-style wrapper files. On Windows CI, this also
// tests the native .cmd resolution path.
func TestResolveNodeBinWindows(t *testing.T) {
	c := qt.New(t)

	dir := t.TempDir()
	nodeModules := filepath.Join(dir, "node_modules")
	binDir := filepath.Join(nodeModules, ".bin")

	// Create target JS file.
	targetJS := filepath.Join(nodeModules, "postcss-cli", "index.js")
	mkdirAndWrite(t, targetJS, "#!/usr/bin/env node\nconsole.log('postcss');\n")
	os.MkdirAll(binDir, 0o755)

	// Simulate what npm creates on Windows: a .cmd wrapper and a shell script.
	cmdWrapper := filepath.Join(binDir, "postcss.cmd")
	cmdContent := "@ECHO off\r\n" +
		"SETLOCAL\r\n" +
		"CALL :find_dp0\r\n" +
		`endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "%_prog%"  "%dp0%\..\postcss-cli\index.js" %*` + "\r\n"
	mkdirAndWrite(t, cmdWrapper, cmdContent)

	shWrapper := filepath.Join(binDir, "postcss")
	shContent := "#!/bin/sh\n" +
		`exec node  "$basedir/../postcss-cli/index.js" "$@"` + "\n"
	mkdirAndWrite(t, shWrapper, shContent)

	t.Logf("GOOS=%s", runtime.GOOS)
	t.Logf("cmd wrapper: %s", cmdWrapper)
	t.Logf("sh wrapper: %s", shWrapper)
	t.Logf("target JS: %s", targetJS)

	c.Run("cmd wrapper resolves to JS", func(c *qt.C) {
		resolved := resolveNodeBin(cmdWrapper)
		t.Logf("resolveNodeBin(%q) = %q", cmdWrapper, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, targetJS), qt.IsTrue)
	})

	c.Run("sh wrapper resolves to JS", func(c *qt.C) {
		resolved := resolveNodeBin(shWrapper)
		t.Logf("resolveNodeBin(%q) = %q", shWrapper, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, targetJS), qt.IsTrue)
	})

	// Simulate postcss-cli 7, whose wrappers point to an extensionless Node
	// shebang script (bin/postcss) rather than a .js file.
	targetNoExt := filepath.Join(nodeModules, "postcss-cli", "bin", "postcss")
	mkdirAndWrite(t, targetNoExt, "#!/usr/bin/env node\nrequire('../');\n")

	cmdNoExt := filepath.Join(binDir, "postcssne.cmd")
	mkdirAndWrite(t, cmdNoExt, "@ECHO off\r\n"+
		`"%dp0%\..\postcss-cli\bin\postcss" %*`+"\r\n")

	shNoExt := filepath.Join(binDir, "postcssne")
	mkdirAndWrite(t, shNoExt, "#!/bin/sh\n"+
		`exec node "$basedir/../postcss-cli/bin/postcss" "$@"`+"\n")

	c.Run("cmd wrapper resolves to extensionless script", func(c *qt.C) {
		resolved := resolveNodeBin(cmdNoExt)
		t.Logf("resolveNodeBin(%q) = %q", cmdNoExt, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, targetNoExt), qt.IsTrue)
	})

	c.Run("sh wrapper resolves to extensionless script", func(c *qt.C) {
		resolved := resolveNodeBin(shNoExt)
		t.Logf("resolveNodeBin(%q) = %q", shNoExt, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, targetNoExt), qt.IsTrue)
	})

	// Simulate `npm install -g` on Windows: the wrapper sits at the npm
	// global prefix and references node_modules as a child (no `..`).
	globalDir := filepath.Join(dir, "global")
	globalTarget := filepath.Join(globalDir, "node_modules", "postcss-cli", "index.js")
	mkdirAndWrite(t, globalTarget, "#!/usr/bin/env node\nconsole.log('postcss');\n")

	globalCmd := filepath.Join(globalDir, "postcss.cmd")
	mkdirAndWrite(t, globalCmd, "@ECHO off\r\n"+
		`endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "%_prog%"  "%dp0%\node_modules\postcss-cli\index.js" %*`+"\r\n")

	globalSh := filepath.Join(globalDir, "postcss")
	mkdirAndWrite(t, globalSh, "#!/bin/sh\n"+
		`exec node  "$basedir/node_modules/postcss-cli/index.js" "$@"`+"\n")

	c.Run("global cmd wrapper resolves to JS", func(c *qt.C) {
		resolved := resolveNodeBin(globalCmd)
		t.Logf("resolveNodeBin(%q) = %q", globalCmd, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, globalTarget), qt.IsTrue)
	})

	c.Run("global sh wrapper resolves to JS", func(c *qt.C) {
		resolved := resolveNodeBin(globalSh)
		t.Logf("resolveNodeBin(%q) = %q", globalSh, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, globalTarget), qt.IsTrue)
	})

	globalScopedTarget := filepath.Join(globalDir, "node_modules", "@babel", "cli", "bin", "babel.js")
	mkdirAndWrite(t, globalScopedTarget, "#!/usr/bin/env node\nconsole.log('babel');\n")
	globalScopedCmd := filepath.Join(globalDir, "babel.cmd")
	mkdirAndWrite(t, globalScopedCmd, "@ECHO off\r\n"+
		`endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "%_prog%"  "%dp0%\node_modules\@babel\cli\bin\babel.js" %*`+"\r\n")

	c.Run("global cmd wrapper scoped package", func(c *qt.C) {
		resolved := resolveNodeBin(globalScopedCmd)
		t.Logf("resolveNodeBin(%q) = %q", globalScopedCmd, resolved)
		c.Assert(resolved, qt.Not(qt.Equals), "")
		c.Assert(sameFile(t, resolved, globalScopedTarget), qt.IsTrue)
	})
}

func sameFile(t *testing.T, a, b string) bool {
	t.Helper()
	infoA, errA := os.Stat(a)
	infoB, errB := os.Stat(b)
	if errA != nil || errB != nil {
		t.Logf("sameFile: stat errors: a=%v, b=%v", errA, errB)
		return false
	}
	return os.SameFile(infoA, infoB)
}

func mkdirAndWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}
