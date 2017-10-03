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
		file              string
		base              template.URL
		slug              string
		url               string
		uglyURLs          bool
		canonifyURLs      bool
		trimTrailingSlash bool
		expectedAbs       string
		expectedRel       string
	}{
		// canonifyURLs=false, trimTrailingSlash=false, uglyURLs=false
		{"x/y/z/boofar.md", "", "", "", false, false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", false, false, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", false, false, false, "http://barnew/boo/x/y/z/boofar/", "/boo/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "boofar", "", false, false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "", "/z/y/q/", false, false, false, "/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", false, false, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, false, false, "http://barnew/boo/x/y/z/booslug/", "/boo/x/y/z/booslug/"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q/", false, false, false, "http://barnew/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q/", false, false, false, "http://barnew/boo/z/y/q/", "/boo/z/y/q/"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q/", false, false, false, "/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q/", false, false, false, "http://barnew/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q/", false, false, false, "http://barnew/boo/z/y/q/", "/boo/z/y/q/"},

		// canonifyURLs=true, trimTrailingSlash=false, uglyURLs=false
		{"x/y/z/boofar.md", "", "", "", false, true, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", false, true, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", false, true, false, "http://barnew/boo/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "boofar", "", false, true, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "", "/z/y/q/", false, true, false, "/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", false, true, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, true, false, "http://barnew/boo/x/y/z/booslug/", "/x/y/z/booslug/"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q/", false, true, false, "http://barnew/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q/", false, true, false, "http://barnew/boo/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q/", false, true, false, "/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q/", false, true, false, "http://barnew/z/y/q/", "/z/y/q/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q/", false, true, false, "http://barnew/boo/z/y/q/", "/z/y/q/"},

		// canonifyURLs=false, trimTrailingSlash=true, uglyURLs=false
		{"x/y/z/boofar.md", "", "", "", false, false, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", false, false, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", false, false, true, "http://barnew/boo/x/y/z/boofar", "/boo/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "boofar", "", false, false, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "", "/z/y/q.html", false, false, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", false, false, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, false, true, "http://barnew/boo/x/y/z/booslug", "/boo/x/y/z/booslug"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q.html", false, false, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q.html", false, false, true, "http://barnew/boo/z/y/q", "/boo/z/y/q"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q.html", false, false, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q.html", false, false, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q.html", false, false, true, "http://barnew/boo/z/y/q", "/boo/z/y/q"},

		// canonifyURLs=false, trimTrailingSlash=false, uglyURLs=true
		{"x/y/z/boofar.md", "", "", "", true, false, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", true, false, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", true, false, false, "http://barnew/boo/x/y/z/boofar.html", "/boo/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "boofar", "", true, false, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "", "/z/y/q.html", true, false, false, "/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", true, false, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, false, false, "http://barnew/boo/x/y/z/booslug.html", "/boo/x/y/z/booslug.html"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q.html", true, false, false, "http://barnew/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q.html", true, false, false, "http://barnew/boo/z/y/q.html", "/boo/z/y/q.html"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q.html", true, false, false, "/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q.html", true, false, false, "http://barnew/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q.html", true, false, false, "http://barnew/boo/z/y/q.html", "/boo/z/y/q.html"},

		// canonifyURLs=true, trimTrailingSlash=true, uglyURLs=false
		{"x/y/z/boofar.md", "", "", "", false, true, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", false, true, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", false, true, true, "http://barnew/boo/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "boofar", "", false, true, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "", "/z/y/q.html", false, true, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", false, true, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, true, true, "http://barnew/boo/x/y/z/booslug", "/x/y/z/booslug"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q.html", false, true, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q.html", false, true, true, "http://barnew/boo/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q.html", false, true, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q.html", false, true, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q.html", false, true, true, "http://barnew/boo/z/y/q", "/z/y/q"},

		// canonifyURLs=true, trimTrailingSlash=false, uglyURLs=true
		{"x/y/z/boofar.md", "", "", "", true, true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", true, true, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", true, true, false, "http://barnew/boo/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "boofar", "", true, true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "", "/z/y/q.html", true, true, false, "/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", true, true, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, true, false, "http://barnew/boo/x/y/z/booslug.html", "/x/y/z/booslug.html"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q.html", true, true, false, "http://barnew/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q.html", true, true, false, "http://barnew/boo/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q.html", true, true, false, "/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q.html", true, true, false, "http://barnew/z/y/q.html", "/z/y/q.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q.html", true, true, false, "http://barnew/boo/z/y/q.html", "/z/y/q.html"},

		// canonifyURLs=false, trimTrailingSlash=true, uglyURLs=true
		{"x/y/z/boofar.md", "", "", "", true, false, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", true, false, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", true, false, true, "http://barnew/boo/x/y/z/boofar", "/boo/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "boofar", "", true, false, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "", "/z/y/q.html", true, false, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", true, false, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, false, true, "http://barnew/boo/x/y/z/booslug", "/boo/x/y/z/booslug"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q.html", true, false, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q.html", true, false, true, "http://barnew/boo/z/y/q", "/boo/z/y/q"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q.html", true, false, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q.html", true, false, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q.html", true, false, true, "http://barnew/boo/z/y/q", "/boo/z/y/q"},

		// canonifyURLs=true, trimTrailingSlash=true, uglyURLs=true
		{"x/y/z/boofar.md", "", "", "", true, true, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", true, true, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "", true, true, true, "http://barnew/boo/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "boofar", "", true, true, true, "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "", "", "/z/y/q.html", true, true, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", true, true, true, "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, true, true, "http://barnew/boo/x/y/z/booslug", "/x/y/z/booslug"},
		{"x/y/z/boofar.md", "http://barnew/", "", "/z/y/q.html", true, true, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "", "/z/y/q.html", true, true, true, "http://barnew/boo/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "", "boofar", "/z/y/q.html", true, true, true, "/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "/z/y/q.html", true, true, true, "http://barnew/z/y/q", "/z/y/q"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "/z/y/q.html", true, true, true, "http://barnew/boo/z/y/q", "/z/y/q"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("uglyURLs=%t,canonifyURLs=%t,trimTrailingSlash=%t", test.uglyURLs, test.canonifyURLs, test.trimTrailingSlash), func(t *testing.T) {
			cfg, fs := newTestCfg()

			cfg.Set("uglyURLs", test.uglyURLs)
			cfg.Set("canonifyURLs", test.canonifyURLs)
			cfg.Set("trimTrailingSlash", test.trimTrailingSlash)
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
		})
	}
}
