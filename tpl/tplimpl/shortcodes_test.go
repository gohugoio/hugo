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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortcodesTemplate(t *testing.T) {

	t.Run("isShortcode", func(t *testing.T) {
		assert := require.New(t)
		assert.True(isShortcode("shortcodes/figures.html"))
		assert.True(isShortcode("_internal/shortcodes/figures.html"))
		assert.False(isShortcode("shortcodes\\figures.html"))
		assert.False(isShortcode("myshortcodes"))

	})

	t.Run("variantsFromName", func(t *testing.T) {
		assert := require.New(t)
		assert.Equal([]string{"", "html", "html"}, templateVariants("figure.html"))
		assert.Equal([]string{"no", "no", "html"}, templateVariants("figure.no.html"))
		assert.Equal([]string{"no", "amp", "html"}, templateVariants("figure.no.amp.html"))
		assert.Equal([]string{"amp", "amp", "html"}, templateVariants("figure.amp.html"))

		name, variants := templateNameAndVariants("figure.html")
		assert.Equal("figure", name)
		assert.Equal([]string{"", "html", "html"}, variants)

	})

	t.Run("compareVariants", func(t *testing.T) {
		assert := require.New(t)
		var s *shortcodeTemplates

		tests := []struct {
			name     string
			name1    string
			name2    string
			expected int
		}{
			{"Same suffix", "figure.html", "figure.html", 3},
			{"Same suffix and output format", "figure.html.html", "figure.html.html", 3},
			{"Same suffix, output format and language", "figure.no.html.html", "figure.no.html.html", 3},
			{"No suffix", "figure", "figure", 3},
			{"Different output format", "figure.amp.html", "figure.html.html", -1},
			{"One with output format, one without", "figure.amp.html", "figure.html", -1},
		}

		for i, test := range tests {
			w := s.compareVariants(templateVariants(test.name1), templateVariants(test.name2))
			assert.Equal(test.expected, w, fmt.Sprintf("[%d] %s", i, test.name))
		}

	})

	t.Run("indexOf", func(t *testing.T) {
		assert := require.New(t)

		s := &shortcodeTemplates{
			variants: []shortcodeVariant{
				{variants: []string{"a", "b", "c"}},
				{variants: []string{"a", "b", "d"}},
			},
		}

		assert.Equal(0, s.indexOf([]string{"a", "b", "c"}))
		assert.Equal(1, s.indexOf([]string{"a", "b", "d"}))
		assert.Equal(-1, s.indexOf([]string{"a", "b", "x"}))

	})

	t.Run("Template", func(t *testing.T) {
		assert := require.New(t)

		assert.True(true)

	})
}
