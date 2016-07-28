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
	"html/template"
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
)

func TestPermalink(t *testing.T) {
	testCommonResetState()

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
		{"x/y/z/boofar.md", "http://barnew/boo/", "boofar", "", true, false, "http://barnew/boo/x/y/z/boofar.html", "/boo/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "boofar", "", false, true, "http://barnew/boo/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "boofar", "", false, false, "http://barnew/boo/x/y/z/boofar/", "/boo/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "boofar", "", true, true, "http://barnew/boo/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo", "boofar", "", true, true, "http://barnew/boo/x/y/z/boofar.html", "/x/y/z/boofar.html"},

		// test URL overrides
		{"x/y/z/boofar.md", "", "", "/z/y/q/", false, false, "/z/y/q/", "/z/y/q/"},
	}

	viper.Set("DefaultExtension", "html")

	for i, test := range tests {
		viper.Set("uglyurls", test.uglyURLs)
		viper.Set("canonifyurls", test.canonifyURLs)
		p := &Page{
			Node: Node{
				URLPath: URLPath{
					Section: "z",
					URL:     test.url,
				},
				Site: &SiteInfo{
					BaseURL: test.base,
				},
			},
			Source: Source{File: *source.NewFile(filepath.FromSlash(test.file))},
		}

		if test.slug != "" {
			p.update(map[string]interface{}{
				"slug": test.slug,
			})
		}

		u, err := p.Permalink()
		if err != nil {
			t.Errorf("Test %d: Unable to process permalink: %s", i, err)
		}

		expected := test.expectedAbs
		if u != expected {
			t.Errorf("Test %d: Expected abs url: %s, got: %s", i, expected, u)
		}

		u, err = p.RelPermalink()
		if err != nil {
			t.Errorf("Test %d: Unable to process permalink: %s", i, err)
		}

		expected = test.expectedRel
		if u != expected {
			t.Errorf("Test %d: Expected rel url: %s, got: %s", i, expected, u)
		}
	}
}
