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
	"reflect"
	"strings"
	"testing"
)

var pageYamlWithTaxonomiesA = `---
tags: ['a', 'B', 'c']
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var pageYamlWithTaxonomiesB = `---
tags:
 - "a"
 - "B"
 - "c"
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var pageYamlWithTaxonomiesC = `---
tags: 'E'
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var pageJSONWithTaxonomies = `{
  "categories": "D",
  "tags": [
    "a",
    "b",
    "c"
  ]
}
JSON Front Matter with tags and categories`

var pageTomlWithTaxonomies = `+++
tags = [ "a", "B", "c" ]
categories = "d"
+++
TOML Front Matter with tags and categories`

func TestParseTaxonomies(t *testing.T) {
	t.Parallel()
	for _, test := range []string{pageTomlWithTaxonomies,
		pageJSONWithTaxonomies,
		pageYamlWithTaxonomiesA,
		pageYamlWithTaxonomiesB,
		pageYamlWithTaxonomiesC,
	} {

		s := newTestSite(t)
		p, _ := s.NewPage("page/with/taxonomy")
		_, err := p.ReadFrom(strings.NewReader(test))
		if err != nil {
			t.Fatalf("Failed parsing %q: %s", test, err)
		}

		param := p.getParamToLower("tags")

		if params, ok := param.([]string); ok {
			expected := []string{"a", "b", "c"}
			if !reflect.DeepEqual(params, expected) {
				t.Errorf("Expected %s: got: %s", expected, params)
			}
		} else if params, ok := param.(string); ok {
			expected := "e"
			if params != expected {
				t.Errorf("Expected %s: got: %s", expected, params)
			}
		}

		param = p.getParamToLower("categories")
		singleparam := param.(string)

		if singleparam != "d" {
			t.Fatalf("Expected: d, got: %s", singleparam)
		}
	}
}
