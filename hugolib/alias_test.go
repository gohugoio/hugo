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
	"testing"

	"github.com/spf13/hugo/deps"
	"github.com/stretchr/testify/require"
)

const pageWithAlias = `---
title: Has Alias
aliases: ["foo/bar/"]
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.
`

const basicTemplate = "<html><body>{{.Content}}</body></html>"
const aliasTemplate = "<html><body>ALIASTEMPLATE</body></html>"

func TestAlias(t *testing.T) {
	t.Parallel()

	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg}
	)

	writeSource(t, fs, filepath.Join("content", "page.md"), pageWithAlias)
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), basicTemplate)

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	// the real page
	th.assertFileContent(t, fs, filepath.Join("public", "page", "index.html"), false, "For some moments the old man")
	// the alias redirector
	th.assertFileContent(t, fs, filepath.Join("public", "foo", "bar", "index.html"), false, "<meta http-equiv=\"refresh\" content=\"0; ")
}

func TestAliasTemplate(t *testing.T) {
	t.Parallel()

	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg}
	)

	writeSource(t, fs, filepath.Join("content", "page.md"), pageWithAlias)
	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), basicTemplate)
	writeSource(t, fs, filepath.Join("layouts", "alias.html"), aliasTemplate)

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, sites.Build(BuildCfg{}))

	// the real page
	th.assertFileContent(t, fs, filepath.Join("public", "page", "index.html"), false, "For some moments the old man")
	// the alias redirector
	th.assertFileContent(t, fs, filepath.Join("public", "foo", "bar", "index.html"), false, "ALIASTEMPLATE")
}
