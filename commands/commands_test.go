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

	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {

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
		commands []string
		flags    []string
	}{
		{[]string{"check", "ulimit"}, nil},
		{[]string{"env"}, nil},
		{[]string{"version"}, nil},
		// no args = hugo build
		{nil, []string{sourceFlag}},
		// TODO(bep) cli refactor remove the HugoSites global and enable the below
		//{nil, []string{sourceFlag, "--renderToMemory"}},
		{[]string{"benchmark"}, []string{sourceFlag, "-n=1"}},
		{[]string{"convert", "toTOML"}, []string{sourceFlag, "-o=" + filepath.Join(dirOut, "toml")}},
		{[]string{"convert", "toYAML"}, []string{sourceFlag, "-o=" + filepath.Join(dirOut, "yaml")}},
		{[]string{"convert", "toJSON"}, []string{sourceFlag, "-o=" + filepath.Join(dirOut, "json")}},
		{[]string{"gen", "autocomplete"}, []string{"--completionfile=" + filepath.Join(dirOut, "autocomplete.txt")}},
		{[]string{"gen", "chromastyles"}, []string{"--style=manni"}},
		{[]string{"gen", "doc"}, []string{"--dir=" + filepath.Join(dirOut, "doc")}},
		{[]string{"gen", "man"}, []string{"--dir=" + filepath.Join(dirOut, "man")}},
		{[]string{"list", "drafts"}, []string{sourceFlag}},
		{[]string{"list", "expired"}, []string{sourceFlag}},
		{[]string{"list", "future"}, []string{sourceFlag}},
		{[]string{"new", "new-page.md"}, []string{sourceFlag}},
		{[]string{"new", "site", filepath.Join(dirOut, "new-site")}, nil},
		// TODO(bep) cli refactor fix https://github.com/gohugoio/hugo/issues/4450
		//{[]string{"new", "theme", filepath.Join(dirOut, "new-theme")}, nil},
	}

	for _, test := range tests {

		hugoCmd := newHugoCompleteCmd()
		test.flags = append(test.flags, "--quiet")
		hugoCmd.SetArgs(append(test.commands, test.flags...))

		// TODO(bep) capture output and add some simple asserts

		assert.NoError(hugoCmd.Execute(), fmt.Sprintf("%v", test.commands))
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
