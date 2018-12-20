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

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {

	assert := require.New(t)

	dir, err := createSimpleTestSite(t)
	assert.NoError(err)

	defer func() {
		os.RemoveAll(dir)
	}()

	resp := Execute([]string{"-s=" + dir})
	assert.NoError(resp.Err)
	result := resp.Result
	assert.True(len(result.Sites) == 1)
	assert.True(len(result.Sites[0].RegularPages) == 1)
}

func TestCommandsPersistentFlags(t *testing.T) {
	assert := require.New(t)

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
	}, func(commands []cmder) {
		var sc *serverCmd
		for _, command := range commands {
			if b, ok := command.(commandsBuilderGetter); ok {
				v := b.getCommandsBuilder().hugoBuilderCommon
				assert.Equal("myconfig.toml", v.cfgFile)
				assert.Equal("myconfigdir", v.cfgDir)
				assert.Equal("mysource", v.source)
				assert.Equal("https://example.com/b/", v.baseURL)
			}

			if srvCmd, ok := command.(*serverCmd); ok {
				sc = srvCmd
			}
		}

		assert.NotNil(sc)
		assert.True(sc.navigateToChanged)
		assert.True(sc.disableLiveReload)
		assert.True(sc.noHTTPCache)
		assert.True(sc.renderToDisk)
		assert.Equal(1366, sc.serverPort)
		assert.Equal("testing", sc.environment)

		cfg := viper.New()
		sc.flagsToConfig(cfg)
		assert.Equal("/tmp/mydestination", cfg.GetString("publishDir"))
		assert.Equal("mycontent", cfg.GetString("contentDir"))
		assert.Equal("mylayouts", cfg.GetString("layoutDir"))
		assert.Equal("mytheme", cfg.GetString("theme"))
		assert.Equal("mythemes", cfg.GetString("themesDir"))
		assert.Equal("https://example.com/b/", cfg.GetString("baseURL"))

		assert.Equal([]string{"page", "home"}, cfg.Get("disableKinds"))

		assert.True(cfg.GetBool("gc"))

		// The flag is named i18n-warnings
		assert.True(cfg.GetBool("logI18nWarnings"))

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
		assert.NoError(rootCmd.Execute())
		test.check(b.commands)
	}

}

func TestCommandsExecute(t *testing.T) {

	assert := require.New(t)

	dir, err := createSimpleTestSite(t)
	assert.NoError(err)

	dirOut, err := ioutil.TempDir("", "hugo-cli-out")
	assert.NoError(err)

	defer func() {
		os.RemoveAll(dir)
		os.RemoveAll(dirOut)
	}()

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

		hugoCmd := newCommandsBuilder().addAll().build().getCommand()
		test.flags = append(test.flags, "--quiet")
		hugoCmd.SetArgs(append(test.commands, test.flags...))

		// TODO(bep) capture output and add some simple asserts
		// TODO(bep) misspelled subcommands does not return an error. We should investigate this
		// but before that, check for "Error: unknown command".

		_, err := hugoCmd.ExecuteC()
		if test.expectErrToContain != "" {
			assert.Error(err, fmt.Sprintf("%v", test.commands))
			assert.Contains(err.Error(), test.expectErrToContain)
		} else {
			assert.NoError(err, fmt.Sprintf("%v", test.commands))
		}

	}

}

func createSimpleTestSite(t *testing.T) (string, error) {
	d, e := ioutil.TempDir("", "hugo-cli")
	if e != nil {
		return "", e
	}

	// Just the basic. These are for CLI tests, not site testing.
	writeFile(t, filepath.Join(d, "config.toml"), `

baseURL = "https://example.org"
title = "Hugo Commands"

`)

	writeFile(t, filepath.Join(d, "content", "p1.md"), `
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

	return d, nil

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
