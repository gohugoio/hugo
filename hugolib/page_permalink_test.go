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
	"fmt"
	"html/template"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gohugoio/hugo/deps"
)

func TestPermalink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		file         string
		base         template.URL
		slug         string
		url          string
		uglyURLs     bool
		canonifyURLs bool
		expectedAbs  string
		expectedRel  string
	}{
		{"x/y/z/boofar.md", "", "", "", false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "", "", false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		// Issue #1174
		{"x/y/z/boofar.md", "http://gopher.com/", "", "", false, true, "http://gopher.com/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://gopher.com/", "", "", true, true, "http://gopher.com/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "boofar", "", false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", false, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", false, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "", "", true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "", "", true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "boofar", "", true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", true, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", true, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, false, "http://barnew/boo/x/y/z/booslug.html", "/boo/x/y/z/booslug.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, true, "http://barnew/boo/x/y/z/booslug/", "/x/y/z/booslug/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, false, "http://barnew/boo/x/y/z/booslug/", "/boo/x/y/z/booslug/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, true, "http://barnew/boo/x/y/z/booslug.html", "/x/y/z/booslug.html"},
		{"x/y/z/boofar.md", "http://barnew/boo", "booslug", "", true, true, "http://barnew/boo/x/y/z/booslug.html", "/x/y/z/booslug.html"},
		// Issue #4666
		{"x/y/z/boo-makeindex.md", "http://barnew/boo", "", "", true, true, "http://barnew/boo/x/y/z/boo-makeindex.html", "/x/y/z/boo-makeindex.html"},

		// test URL overrides
		{"x/y/z/boofar.md", "", "", "/z/y/q/", false, false, "/z/y/q/", "/z/y/q/"},
	}

	for i, test := range tests {

		cfg, fs := newTestCfg()

		cfg.Set("uglyURLs", test.uglyURLs)
		cfg.Set("canonifyURLs", test.canonifyURLs)
		cfg.Set("baseURL", test.base)

		pageContent := fmt.Sprintf(`---
title: Page
slug: %q
url: %q
---
Content
`, test.slug, test.url)

		writeSource(t, fs, filepath.Join("content", filepath.FromSlash(test.file)), pageContent)

		s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})
		require.Len(t, s.RegularPages, 1)

		p := s.RegularPages[0]

		u := p.Permalink()

		expected := test.expectedAbs
		if u != expected {
			t.Fatalf("[%d] Expected abs url: %s, got: %s", i, expected, u)
		}

		u = p.RelPermalink()

		expected = test.expectedRel
		if u != expected {
			t.Errorf("[%d] Expected rel url: %s, got: %s", i, expected, u)
		}
	}
}
