// Copyright 2015 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gohugoio/hugo/deps"
	"github.com/stretchr/testify/require"
)

const pageWithAlias = `---
title: Has Alias
aliases: ["foo/bar/"]
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.
`

const pageWithAliasMultipleOutputs = `---
title: Has Alias for HTML and AMP
aliases: ["foo/bar/"]
outputs: ["HTML", "AMP", "JSON"]
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.
`

const basicTemplate = "<html><body>{{.Content}}</body></html>"
const aliasTemplate = "<html><body>ALIASTEMPLATE</body></html>"

func TestAlias(t *testing.T) {
	t.Parallel()

	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	writeSource(t, fs, filepath.Join("content", "page.md"), pageWithAlias)
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), basicTemplate)

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	// the real page
	th.assertFileContent(filepath.Join("public", "page", "index.html"), "For some moments the old man")
	// the alias redirector
	th.assertFileContent(filepath.Join("public", "foo", "bar", "index.html"), "<meta http-equiv=\"refresh\" content=\"0; ")
}

func TestAliasMultipleOutputFormats(t *testing.T) {
	t.Parallel()

	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	writeSource(t, fs, filepath.Join("content", "page.md"), pageWithAliasMultipleOutputs)
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), basicTemplate)
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.amp.html"), basicTemplate)
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.json"), basicTemplate)

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	// the real pages
	th.assertFileContent(filepath.Join("public", "page", "index.html"), "For some moments the old man")
	th.assertFileContent(filepath.Join("public", "amp", "page", "index.html"), "For some moments the old man")
	th.assertFileContent(filepath.Join("public", "page", "index.json"), "For some moments the old man")

	// the alias redirectors
	th.assertFileContent(filepath.Join("public", "foo", "bar", "index.html"), "<meta http-equiv=\"refresh\" content=\"0; ")
	th.assertFileContent(filepath.Join("public", "foo", "bar", "amp", "index.html"), "<meta http-equiv=\"refresh\" content=\"0; ")
	require.False(t, destinationExists(th.Fs, filepath.Join("public", "foo", "bar", "index.json")))
}

func TestAliasTemplate(t *testing.T) {
	t.Parallel()

	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	writeSource(t, fs, filepath.Join("content", "page.md"), pageWithAlias)
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), basicTemplate)
	writeSource(t, fs, filepath.Join("layouts", "alias.html"), aliasTemplate)

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, sites.Build(BuildCfg{}))

	// the real page
	th.assertFileContent(filepath.Join("public", "page", "index.html"), "For some moments the old man")
	// the alias redirector
	th.assertFileContent(filepath.Join("public", "foo", "bar", "index.html"), "ALIASTEMPLATE")
}

func TestTargetPathHTMLRedirectAlias(t *testing.T) {
	h := newAliasHandler(nil, newErrorLogger(), false)

	errIsNilForThisOS := runtime.GOOS != "windows"

	tests := []struct {
		value    string
		expected string
		errIsNil bool
	}{
		{"", "", false},
		{"s", filepath.FromSlash("s/index.html"), true},
		{"/", "", false},
		{"alias 1", filepath.FromSlash("alias 1/index.html"), true},
		{"alias 2/", filepath.FromSlash("alias 2/index.html"), true},
		{"alias 3.html", "alias 3.html", true},
		{"alias4.html", "alias4.html", true},
		{"/alias 5.html", "alias 5.html", true},
		{"/трям.html", "трям.html", true},
		{"../../../../tmp/passwd", "", false},
		{"/foo/../../../../tmp/passwd", filepath.FromSlash("tmp/passwd/index.html"), true},
		{"foo/../../../../tmp/passwd", "", false},
		{"C:\\Windows", filepath.FromSlash("C:\\Windows/index.html"), errIsNilForThisOS},
		{"/trailing-space /", filepath.FromSlash("trailing-space /index.html"), errIsNilForThisOS},
		{"/trailing-period./", filepath.FromSlash("trailing-period./index.html"), errIsNilForThisOS},
		{"/tab\tseparated/", filepath.FromSlash("tab\tseparated/index.html"), errIsNilForThisOS},
		{"/chrome/?p=help&ctx=keyboard#topic=3227046", filepath.FromSlash("chrome/?p=help&ctx=keyboard#topic=3227046/index.html"), errIsNilForThisOS},
		{"/LPT1/Printer/", filepath.FromSlash("LPT1/Printer/index.html"), errIsNilForThisOS},
	}

	for _, test := range tests {
		path, err := h.targetPathAlias(test.value)
		if (err == nil) != test.errIsNil {
			t.Errorf("Expected err == nil => %t, got: %t. err: %s", test.errIsNil, err == nil, err)
			continue
		}
		if err == nil && path != test.expected {
			t.Errorf("Expected: \"%s\", got: \"%s\"", test.expected, path)
		}
	}
}
