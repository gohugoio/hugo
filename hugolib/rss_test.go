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
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/deps"
)

func TestRSSOutput(t *testing.T) {
	t.Parallel()

	rssLimit := len(weightedSources) - 1

	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub/")
	cfg.Set("title", "RSSTest")
	cfg.Set("rssLimit", rssLimit)
	th, configs := newTestHelperFromProvider(cfg, fs, t)

	rssURI := "index.xml"

	for _, src := range weightedSources {
		writeSource(t, fs, filepath.Join("content", "sect", src[0]), src[1])
	}

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

	// Home RSS
	th.assertFileContent(filepath.Join("public", rssURI), "<?xml", "rss version", "RSSTest")
	// Section RSS
	th.assertFileContent(filepath.Join("public", "sect", rssURI), "<?xml", "rss version", "Sects on RSSTest")
	// Taxonomy RSS
	th.assertFileContent(filepath.Join("public", "categories", "hugo", rssURI), "<?xml", "rss version", "Hugo on RSSTest")

	// RSS Item Limit
	content := readWorkingDir(t, fs, filepath.Join("public", rssURI))
	c := strings.Count(content, "<item>")
	if c != rssLimit {
		t.Errorf("incorrect RSS item count: expected %d, got %d", rssLimit, c)
	}

	// Encoded summary
	th.assertFileContent(filepath.Join("public", rssURI), "<?xml", "description", "A &lt;em&gt;custom&lt;/em&gt; summary")
}

// Before Hugo 0.49 we set the pseudo page kind RSS on the page when output to RSS.
// This had some unintended side effects, esp. when the only output format for that page
// was RSS.
// For the page kinds that can have multiple output formats, the Kind should be one of the
// standard home, page etc.
// This test has this single purpose: Check that the Kind is that of the source page.
// See https://github.com/gohugoio/hugo/issues/5138
func TestRSSKind(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded("index.rss.xml", `RSS Kind: {{ .Kind }}`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.xml", "RSS Kind: home")
}

func TestRSSCanonifyURLs(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded("index.rss.xml", `<rss>{{ range .Pages }}<item>{{ .Content | html }}</item>{{ end }}</rss>`)
	b.WithContent("page.md", `---
Title: My Page
---

Figure:

{{< figure src="/images/sunset.jpg" title="Sunset" >}}



`)
	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.xml", "img src=&#34;http://example.com/images/sunset.jpg")
}
