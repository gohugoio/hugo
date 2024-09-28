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
	"fmt"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/kinds"
)

func TestUglyURLsPerSection(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	const dt = `---
title: Do not go gentle into that good night
---

Wild men who caught and sang the sun in flight,
And learn, too late, they grieved it on its way,
Do not go gentle into that good night.

`

	cfg, fs := newTestCfg()

	cfg.Set("uglyURLs", map[string]bool{
		"sect2": true,
	})
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)

	writeSource(t, fs, filepath.Join("content", "sect1", "p1.md"), dt)
	writeSource(t, fs, filepath.Join("content", "sect2", "p2.md"), dt)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{SkipRender: true})

	c.Assert(len(s.RegularPages()), qt.Equals, 2)

	notUgly := s.getPageOldVersion(kinds.KindPage, "sect1/p1.md")
	c.Assert(notUgly, qt.Not(qt.IsNil))
	c.Assert(notUgly.Section(), qt.Equals, "sect1")
	c.Assert(notUgly.RelPermalink(), qt.Equals, "/sect1/p1/")

	ugly := s.getPageOldVersion(kinds.KindPage, "sect2/p2.md")
	c.Assert(ugly, qt.Not(qt.IsNil))
	c.Assert(ugly.Section(), qt.Equals, "sect2")
	c.Assert(ugly.RelPermalink(), qt.Equals, "/sect2/p2.html")
}

func TestSectionWithURLInFrontMatter(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	const st = `---
title: Do not go gentle into that good night
url: %s
---

Wild men who caught and sang the sun in flight,
And learn, too late, they grieved it on its way,
Do not go gentle into that good night.

`

	const pt = `---
title: Wild men who caught and sang the sun in flight
---

Wild men who caught and sang the sun in flight,
And learn, too late, they grieved it on its way,
Do not go gentle into that good night.

`

	cfg, fs := newTestCfg()
	cfg.Set("pagination.pagerSize", 1)
	th, configs := newTestHelperFromProvider(cfg, fs, t)

	writeSource(t, fs, filepath.Join("content", "sect1", "_index.md"), fmt.Sprintf(st, "/ss1/"))
	writeSource(t, fs, filepath.Join("content", "sect2", "_index.md"), fmt.Sprintf(st, "/ss2/"))

	for i := 0; i < 5; i++ {
		writeSource(t, fs, filepath.Join("content", "sect1", fmt.Sprintf("p%d.md", i+1)), pt)
		writeSource(t, fs, filepath.Join("content", "sect2", fmt.Sprintf("p%d.md", i+1)), pt)
	}

	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), "<html><body>{{.Content}}</body></html>")
	writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"),
		"<html><body>P{{.Paginator.PageNumber}}|URL: {{.Paginator.URL}}|{{ if .Paginator.HasNext }}Next: {{.Paginator.Next.URL }}{{ end }}</body></html>")

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

	c.Assert(len(s.RegularPages()), qt.Equals, 10)

	sect1 := s.getPageOldVersion(kinds.KindSection, "sect1")
	c.Assert(sect1, qt.Not(qt.IsNil))
	c.Assert(sect1.RelPermalink(), qt.Equals, "/ss1/")
	th.assertFileContent(filepath.Join("public", "ss1", "index.html"), "P1|URL: /ss1/|Next: /ss1/page/2/")
	th.assertFileContent(filepath.Join("public", "ss1", "page", "2", "index.html"), "P2|URL: /ss1/page/2/|Next: /ss1/page/3/")
}

func TestSectionsEntries(t *testing.T) {
	files := `
-- hugo.toml --
-- content/withfile/_index.md --
-- content/withoutfile/p1.md --
-- layouts/_default/list.html --
SectionsEntries: {{ .SectionsEntries }}


`

	b := Test(t, files)

	b.AssertFileContent("public/withfile/index.html", "SectionsEntries: [withfile]")
	b.AssertFileContent("public/withoutfile/index.html", "SectionsEntries: [withoutfile]")
}
