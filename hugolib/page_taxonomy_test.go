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
	"strings"
	"testing"
)

var PAGE_YAML_WITH_TAXONOMIES_A = `---
tags: ['a', 'B', 'c']
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var PAGE_YAML_WITH_TAXONOMIES_B = `---
tags:
 - "a"
 - "B"
 - "c"
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var PAGE_YAML_WITH_TAXONOMIES_C = `---
tags: 'E'
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var PAGE_JSON_WITH_TAXONOMIES = `{
  "categories": "D",
  "tags": [
    "a",
    "b",
    "c"
  ]
}
JSON Front Matter with tags and categories`

var PAGE_TOML_WITH_TAXONOMIES = `+++
tags = [ "a", "B", "c" ]
categories = "d"
+++
TOML Front Matter with tags and categories`

func TestParseTaxonomies(t *testing.T) {
	for _, test := range []string{PAGE_TOML_WITH_TAXONOMIES,
		PAGE_JSON_WITH_TAXONOMIES,
		PAGE_YAML_WITH_TAXONOMIES_A,
		PAGE_YAML_WITH_TAXONOMIES_B,
		PAGE_YAML_WITH_TAXONOMIES_C,
	} {

		p, _ := NewPage("page/with/taxonomy")
		_, err := p.ReadFrom(strings.NewReader(test))
		if err != nil {
			t.Fatalf("Failed parsing %q: %s", test, err)
		}

		param := p.GetParam("tags")

		if params, ok := param.([]string); ok {
			expected := []string{"a", "b", "c"}
			if !compareStringSlice(params, expected) {
				t.Errorf("Expected %s: got: %s", expected, params)
			}
		} else if params, ok := param.(string); ok {
			expected := "e"
			if params != expected {
				t.Errorf("Expected %s: got: %s", expected, params)
			}
		}

		param = p.GetParam("categories")
		singleparam := param.(string)

		if singleparam != "d" {
			t.Fatalf("Expected: d, got: %s", singleparam)
		}
	}
}

func compareStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}

	return true
}
