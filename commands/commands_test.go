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

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
)

func TestExecute(t *testing.T) {

	c := qt.New(t)

	createSite := func(c *qt.C) (string, func()) {
		dir, clean, err := createSimpleTestSite(t, testSiteConfig{})
		c.Assert(err, qt.IsNil)
		return dir, clean
	}

	c.Run("hugo", func(c *qt.C) {
		dir, clean := createSite(c)
		defer clean()
		resp := Execute([]string{"-s=" + dir})
		c.Assert(resp.Err, qt.IsNil)
		result := resp.Result
		c.Assert(len(result.Sites) == 1, qt.Equals, true)
		c.Assert(len(result.Sites[0].RegularPages()) == 1, qt.Equals, true)
		c.Assert(result.Sites[0].Info.Params()["myparam"], qt.Equals, "paramproduction")
	})

	c.Run("hugo, set environment", func(c *qt.C) {
		dir, clean := createSite(c)
		defer clean()
		resp := Execute([]string{"-s=" + dir, "-e=staging"})
		c.Assert(resp.Err, qt.IsNil)
		result := resp.Result
		c.Assert(result.Sites[0].Info.Params()["myparam"], qt.Equals, "paramstaging")
	})

	c.Run("convert toJSON", func(c *qt.C) {
		dir, clean := createSite(c)
		output := filepath.Join(dir, "myjson")
		defer clean()
		resp := Execute([]string{"convert", "toJSON", "-s=" + dir, "-e=staging", "-o=" + output})
		c.Assert(resp.Err, qt.IsNil)
		converted := readFileFrom(c, filepath.Join(output, "content", "p1.md"))
		c.Assert(converted, qt.Equals, "{\n   \"title\": \"P1\",\n   \"weight\": 1\n}\n\nContent\n\n", qt.Commentf(converted))
	})

	c.Run("config, set environment", func(c *qt.C) {
		dir, clean := createSite(c)
		defer clean()
		out, err := captureStdout(func() error {
			resp := Execute([]string{"config", "-s=" + dir, "-e=staging"})
			return resp.Err
		})
		c.Assert(err, qt.IsNil)
		c.Assert(out, qt.Contains, "params = map[myparam:paramstaging]", qt.Commentf(out))
	})

	c.Run("deploy, environment set", func(c *qt.C) {
		dir, clean := createSite(c)
		defer clean()
		resp := Execute([]string{"deploy", "-s=" + dir, "-e=staging", "--target=mydeployment", "--dryRun"})
		c.Assert(resp.Err, qt.Not(qt.IsNil))
		c.Assert(resp.Err.Error(), qt.Contains, `no provider registered for "hugocloud"`)
	})

	c.Run("list", func(c *qt.C) {
		dir, clean := createSite(c)
		defer clean()
		out, err := captureStdout(func() error {
			resp := Execute([]string{"list", "all", "-s=" + dir, "-e=staging"})
			return resp.Err
		})
		c.Assert(err, qt.IsNil)
		c.Assert(out, qt.Contains, "p1.md")
	})

	c.Run("new theme", func(c *qt.C) {
		dir, clean := createSite(c)
		defer clean()
		themesDir := filepath.Join(dir, "mythemes")
		resp := Execute([]string{"new", "theme", "mytheme", "-s=" + dir, "-e=staging", "--themesDir=" + themesDir})
		c.Assert(resp.Err, qt.IsNil)
		themeTOML := readFileFrom(c, filepath.Join(themesDir, "mytheme", "theme.toml"))
		c.Assert(themeTOML, qt.Contains, "name = \"Mytheme\"")
	})

	c.Run("new site", func(c *qt.C) {
		dir, clean := createSite(c)
		defer clean()
		siteDir := filepath.Join(dir, "mysite")
		resp := Execute([]string{"new", "site", siteDir, "-e=staging"})
		c.Assert(resp.Err, qt.IsNil)
		config := readFileFrom(c, filepath.Join(siteDir, "config.toml"))
		c.Assert(config, qt.Contains, "baseURL = \"http://example.org/\"")
		checkNewSiteInited(c, siteDir)
	})

}

func checkNewSiteInited(c *qt.C, basepath string) {
	paths := []string{
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		filepath.Join(basepath, "archetypes"),
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
		filepath.Join(basepath, "config.toml"),
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		c.Assert(err, qt.IsNil)
	}
}

func readFileFrom(c *qt.C, filename string) string {
	c.Helper()
	filename = filepath.Clean(filename)
	b, err := afero.ReadFile(hugofs.Os, filename)
	c.Assert(err, qt.IsNil)
	return string(b)
}

func TestCommandsPersistentFlags(t *testing.T) {
	c := qt.New(t)

	noOpRunE := func(cmd *cobra.Command, args []string) error {
		return nil
	}

	tests := []struct {
		args  []string
		check func(command []cmder)
	}{{[]string{"server",
		"--config=myconfig.toml",
		"--configDir=myconfigdir",
		"--contentDir=mycontent",
		"--disableKinds=page,home",
		"--environment=testing",
		"--configDir=myconfigdir",
		"--layoutDir=mylayouts",
		"--theme=mytheme",
		"--gc",
		"--themesDir=mythemes",
		"--cleanDestinationDir",
		"--navigateToChanged",
		"--disableLiveReload",
		"--noHTTPCache",
		"--i18n-warnings",
		"--destination=/tmp/mydestination",
		"-b=https://example.com/b/",
		"--port=1366",
		"--renderToDisk",
		"--source=mysource",
		"--path-warnings",
	}, func(commands []cmder) {
		var sc *serverCmd
		for _, command := range commands {
			if b, ok := command.(commandsBuilderGetter); ok {
				v := b.getCommandsBuilder().hugoBuilderCommon
				c.Assert(v.cfgFile, qt.Equals, "myconfig.toml")
				c.Assert(v.cfgDir, qt.Equals, "myconfigdir")
				c.Assert(v.source, qt.Equals, "mysource")
				c.Assert(v.baseURL, qt.Equals, "https://example.com/b/")
			}

			if srvCmd, ok := command.(*serverCmd); ok {
				sc = srvCmd
			}
		}

		c.Assert(sc, qt.Not(qt.IsNil))
		c.Assert(sc.navigateToChanged, qt.Equals, true)
		c.Assert(sc.disableLiveReload, qt.Equals, true)
		c.Assert(sc.noHTTPCache, qt.Equals, true)
		c.Assert(sc.renderToDisk, qt.Equals, true)
		c.Assert(sc.serverPort, qt.Equals, 1366)
		c.Assert(sc.environment, qt.Equals, "testing")

		cfg := viper.New()
		sc.flagsToConfig(cfg)
		c.Assert(cfg.GetString("publishDir"), qt.Equals, "/tmp/mydestination")
		c.Assert(cfg.GetString("contentDir"), qt.Equals, "mycontent")
		c.Assert(cfg.GetString("layoutDir"), qt.Equals, "mylayouts")
		c.Assert(cfg.GetStringSlice("theme"), qt.DeepEquals, []string{"mytheme"})
		c.Assert(cfg.GetString("themesDir"), qt.Equals, "mythemes")
		c.Assert(cfg.GetString("baseURL"), qt.Equals, "https://example.com/b/")

		c.Assert(cfg.Get("disableKinds"), qt.DeepEquals, []string{"page", "home"})

		c.Assert(cfg.GetBool("gc"), qt.Equals, true)

		// The flag is named path-warnings
		c.Assert(cfg.GetBool("logPathWarnings"), qt.Equals, true)

		// The flag is named i18n-warnings
		c.Assert(cfg.GetBool("logI18nWarnings"), qt.Equals, true)

	}}}

	for _, test := range tests {
		b := newCommandsBuilder()
		root := b.addAll().build()

		for _, c := range b.commands {
			if c.getCommand() == nil {
				continue
			}
			// We are only intereseted in the flag handling here.
			c.getCommand().RunE = noOpRunE
		}
		rootCmd := root.getCommand()
		rootCmd.SetArgs(test.args)
		c.Assert(rootCmd.Execute(), qt.IsNil)
		test.check(b.commands)
	}

}

func TestCommandsExecute(t *testing.T) {

	c := qt.New(t)

	dir, clean, err := createSimpleTestSite(t, testSiteConfig{})
	c.Assert(err, qt.IsNil)

	dirOut, clean2, err := htesting.CreateTempDir(hugofs.Os, "hugo-cli-out")
	c.Assert(err, qt.IsNil)

	defer clean()
	defer clean2()

	sourceFlag := fmt.Sprintf("-s=%s", dir)

	tests := []struct {
		commands           []string
		flags              []string
		expectErrToContain string
	}{
		// TODO(bep) permission issue on my OSX? "operation not permitted" {[]string{"check", "ulimit"}, nil, false},
		{[]string{"env"}, nil, ""},
		{[]string{"version"}, nil, ""},
		// no args = hugo build
		{nil, []string{sourceFlag}, ""},
		{nil, []string{sourceFlag, "--renderToMemory"}, ""},
		{[]string{"config"}, []string{sourceFlag}, ""},
		{[]string{"convert", "toTOML"}, []string{sourceFlag, "-o=" + filepath.Join(dirOut, "toml")}, ""},
		{[]string{"convert", "toYAML"}, []string{sourceFlag, "-o=" + filepath.Join(dirOut, "yaml")}, ""},
		{[]string{"convert", "toJSON"}, []string{sourceFlag, "-o=" + filepath.Join(dirOut, "json")}, ""},
		{[]string{"gen", "autocomplete"}, []string{"--completionfile=" + filepath.Join(dirOut, "autocomplete.txt")}, ""},
		{[]string{"gen", "chromastyles"}, []string{"--style=manni"}, ""},
		{[]string{"gen", "doc"}, []string{"--dir=" + filepath.Join(dirOut, "doc")}, ""},
		{[]string{"gen", "man"}, []string{"--dir=" + filepath.Join(dirOut, "man")}, ""},
		{[]string{"list", "drafts"}, []string{sourceFlag}, ""},
		{[]string{"list", "expired"}, []string{sourceFlag}, ""},
		{[]string{"list", "future"}, []string{sourceFlag}, ""},
		{[]string{"new", "new-page.md"}, []string{sourceFlag}, ""},
		{[]string{"new", "site", filepath.Join(dirOut, "new-site")}, nil, ""},
		{[]string{"unknowncommand"}, nil, "unknown command"},
		// TODO(bep) cli refactor fix https://github.com/gohugoio/hugo/issues/4450
		//{[]string{"new", "theme", filepath.Join(dirOut, "new-theme")}, nil,false},
	}

	for _, test := range tests {
		b := newCommandsBuilder().addAll().build()
		hugoCmd := b.getCommand()
		test.flags = append(test.flags, "--quiet")
		hugoCmd.SetArgs(append(test.commands, test.flags...))

		// TODO(bep) capture output and add some simple asserts
		// TODO(bep) misspelled subcommands does not return an error. We should investigate this
		// but before that, check for "Error: unknown command".

		_, err := hugoCmd.ExecuteC()
		if test.expectErrToContain != "" {
			c.Assert(err, qt.Not(qt.IsNil))
			c.Assert(err.Error(), qt.Contains, test.expectErrToContain)
		} else {
			c.Assert(err, qt.IsNil)
		}

		// Assert that we have not left any development debug artifacts in
		// the code.
		if b.c != nil {
			_, ok := b.c.destinationFs.(types.DevMarker)
			c.Assert(ok, qt.Equals, false)
		}

	}

}

type testSiteConfig struct {
	configTOML string
	contentDir string
}

func createSimpleTestSite(t *testing.T, cfg testSiteConfig) (string, func(), error) {
	d, clean, e := htesting.CreateTempDir(hugofs.Os, "hugo-cli")
	if e != nil {
		return "", nil, e
	}

	cfgStr := `

baseURL = "https://example.org"
title = "Hugo Commands"


`

	contentDir := "content"

	if cfg.configTOML != "" {
		cfgStr = cfg.configTOML
	}
	if cfg.contentDir != "" {
		contentDir = cfg.contentDir
	}

	os.MkdirAll(filepath.Join(d, "public"), 0777)

	// Just the basic. These are for CLI tests, not site testing.
	writeFile(t, filepath.Join(d, "config.toml"), cfgStr)
	writeFile(t, filepath.Join(d, "config", "staging", "params.toml"), `myparam="paramstaging"`)
	writeFile(t, filepath.Join(d, "config", "staging", "deployment.toml"), `
[[targets]]
name = "mydeployment"
URL = "hugocloud://hugotestbucket"
`)

	writeFile(t, filepath.Join(d, "config", "testing", "params.toml"), `myparam="paramtesting"`)
	writeFile(t, filepath.Join(d, "config", "production", "params.toml"), `myparam="paramproduction"`)

	writeFile(t, filepath.Join(d, contentDir, "p1.md"), `
---
title: "P1"
weight: 1
---

Content

`)

	writeFile(t, filepath.Join(d, "layouts", "_default", "single.html"), `

Single: {{ .Title }}

`)

	writeFile(t, filepath.Join(d, "layouts", "_default", "list.html"), `

List: {{ .Title }}
Environment: {{ hugo.Environment }}

`)

	return d, clean, nil

}

func writeFile(t *testing.T, filename, content string) {
	must(t, os.MkdirAll(filepath.Dir(filename), os.FileMode(0755)))
	must(t, ioutil.WriteFile(filename, []byte(content), os.FileMode(0755)))
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
