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

package tplimpl

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestShortcodesTemplate(t *testing.T) {

	t.Run("isShortcode", func(t *testing.T) {
		c := qt.New(t)
		c.Assert(isShortcode("shortcodes/figures.html"), qt.Equals, true)
		c.Assert(isShortcode("_internal/shortcodes/figures.html"), qt.Equals, true)
		c.Assert(isShortcode("shortcodes\\figures.html"), qt.Equals, false)
		c.Assert(isShortcode("myshortcodes"), qt.Equals, false)

	})

	t.Run("variantsFromName", func(t *testing.T) {
		c := qt.New(t)
		c.Assert(templateVariants("figure.html"), qt.DeepEquals, []string{"", "html", "html"})
		c.Assert(templateVariants("figure.no.html"), qt.DeepEquals, []string{"no", "no", "html"})
		c.Assert(templateVariants("figure.no.amp.html"), qt.DeepEquals, []string{"no", "amp", "html"})
		c.Assert(templateVariants("figure.amp.html"), qt.DeepEquals, []string{"amp", "amp", "html"})

		name, variants := templateNameAndVariants("figure.html")
		c.Assert(name, qt.Equals, "figure")
		c.Assert(variants, qt.DeepEquals, []string{"", "html", "html"})

	})

	t.Run("compareVariants", func(t *testing.T) {
		c := qt.New(t)
		var s *shortcodeTemplates

		tests := []struct {
			name     string
			name1    string
			name2    string
			expected int
		}{
			{"Same suffix", "figure.html", "figure.html", 6},
			{"Same suffix and output format", "figure.html.html", "figure.html.html", 6},
			{"Same suffix, output format and language", "figure.no.html.html", "figure.no.html.html", 6},
			{"No suffix", "figure", "figure", 6},
			{"Different output format", "figure.amp.html", "figure.html.html", -1},
			{"One with output format, one without", "figure.amp.html", "figure.html", -1},
		}

		for _, test := range tests {
			w := s.compareVariants(templateVariants(test.name1), templateVariants(test.name2))
			c.Assert(w, qt.Equals, test.expected)
		}

	})

	t.Run("indexOf", func(t *testing.T) {
		c := qt.New(t)

		s := &shortcodeTemplates{
			variants: []shortcodeVariant{
				{variants: []string{"a", "b", "c"}},
				{variants: []string{"a", "b", "d"}},
			},
		}

		c.Assert(s.indexOf([]string{"a", "b", "c"}), qt.Equals, 0)
		c.Assert(s.indexOf([]string{"a", "b", "d"}), qt.Equals, 1)
		c.Assert(s.indexOf([]string{"a", "b", "x"}), qt.Equals, -1)

	})

	t.Run("Name", func(t *testing.T) {
		c := qt.New(t)

		c.Assert(templateBaseName(templateShortcode, "shortcodes/foo.html"), qt.Equals, "foo.html")
		c.Assert(templateBaseName(templateShortcode, "_internal/shortcodes/foo.html"), qt.Equals, "foo.html")
		c.Assert(templateBaseName(templateShortcode, "shortcodes/test/foo.html"), qt.Equals, "test/foo.html")

		c.Assert(true, qt.Equals, true)

	})
}
