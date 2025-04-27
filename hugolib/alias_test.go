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

package hugolib

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
)

const pageWithAlias = `---
title: Has Alias
aliases: ["/foo/bar/", "rel"]
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.
`

const pageWithAliasMultipleOutputs = `---
title: Has Alias for HTML and AMP
aliases: ["/foo/bar/"]
outputs: ["HTML", "AMP", "JSON"]
---
For some moments the old man did not reply. He stood with bowed head, buried in deep thought. But at last he spoke.
`

const (
	basicTemplate = "<html><body>{{.Content}}</body></html>"
	aliasTemplate = "<html><body>ALIASTEMPLATE</body></html>"
)

func TestAlias(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	tests := []struct {
		fileSuffix string
		urlPrefix  string
		urlSuffix  string
		settings   map[string]any
	}{
		{"/index.html", "http://example.com", "/", map[string]any{"baseURL": "http://example.com"}},
		{"/index.html", "http://example.com/some/path", "/", map[string]any{"baseURL": "http://example.com/some/path"}},
		{"/index.html", "http://example.com", "/", map[string]any{"baseURL": "http://example.com", "canonifyURLs": true}},
		{"/index.html", "../..", "/", map[string]any{"relativeURLs": true}},
		{".html", "", ".html", map[string]any{"uglyURLs": true}},
	}

	for _, test := range tests {
		b := newTestSitesBuilder(t)
		b.WithSimpleConfigFileAndSettings(test.settings).WithContent("blog/page.md", pageWithAlias)
		b.CreateSites().Build(BuildCfg{})

		c.Assert(len(b.H.Sites), qt.Equals, 1)
		c.Assert(len(b.H.Sites[0].RegularPages()), qt.Equals, 1)

		// the real page
		b.AssertFileContent("public/blog/page"+test.fileSuffix, "For some moments the old man")
		// the alias redirectors
		b.AssertFileContent("public/foo/bar"+test.fileSuffix, "<meta http-equiv=\"refresh\" content=\"0; url="+test.urlPrefix+"/blog/page"+test.urlSuffix+"\">")
		b.AssertFileContent("public/blog/rel"+test.fileSuffix, "<meta http-equiv=\"refresh\" content=\"0; url="+test.urlPrefix+"/blog/page"+test.urlSuffix+"\">")
	}
}

func TestAliasMultipleOutputFormats(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithContent("blog/page.md", pageWithAliasMultipleOutputs)

	b.WithTemplates(
		"_default/single.html", basicTemplate,
		"_default/single.amp.html", basicTemplate,
		"_default/single.json", basicTemplate)

	b.CreateSites().Build(BuildCfg{})

	b.H.Sites[0].m.debugPrint("", 999, os.Stdout)

	// the real pages
	b.AssertFileContent("public/blog/page/index.html", "For some moments the old man")
	b.AssertFileContent("public/amp/blog/page/index.html", "For some moments the old man")
	b.AssertFileContent("public/blog/page/index.json", "For some moments the old man")

	// the alias redirectors
	b.AssertFileContent("public/foo/bar/index.html", "<meta http-equiv=\"refresh\" content=\"0; ")
	b.AssertFileContent("public/amp/foo/bar/index.html", "<meta http-equiv=\"refresh\" content=\"0; ")
	c.Assert(b.CheckExists("public/foo/bar/index.json"), qt.Equals, false)
}

func TestAliasTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com"
-- layouts/single.html --
Single.
-- layouts/home.html --
Home.
-- layouts/alias.html --
ALIASTEMPLATE
-- content/page.md --
---
title: "Page"
aliases: ["/foo/bar/"]
---
`

	b := Test(t, files)

	// the real page
	b.AssertFileContent("public/page/index.html", "Single.")
	// the alias redirector
	b.AssertFileContent("public/foo/bar/index.html", "ALIASTEMPLATE")
}

func TestTargetPathHTMLRedirectAlias(t *testing.T) {
	h := newAliasHandler(nil, loggers.NewDefault(), false)

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
			t.Errorf("Expected: %q, got: %q", test.expected, path)
		}
	}
}
