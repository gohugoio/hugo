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
