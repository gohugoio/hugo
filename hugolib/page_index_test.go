package hugolib

import (
	"strings"
	"testing"
)

var PAGE_YAML_WITH_INDEXES_A = `---
tags: ['a', 'b', 'c']
categories: 'd'
---
YAML frontmatter with tags and categories index.`

var PAGE_YAML_WITH_INDEXES_B = `---
tags: 
 - "a"
 - "b"
 - "c"
categories: 'd'
---
YAML frontmatter with tags and categories index.`

var PAGE_JSON_WITH_INDEXES = `{
  "categories": "d",
  "tags": [
    "a", 
    "b", 
    "c"
  ]
}
JSON Front Matter with tags and categories`

var PAGE_TOML_WITH_INDEXES = `+++
tags = [ "a", "b", "c" ]
categories = "d"
+++
TOML Front Matter with tags and categories`

func TestParseIndexes(t *testing.T) {
	for _, test := range []string{PAGE_TOML_WITH_INDEXES,
		PAGE_JSON_WITH_INDEXES,
		PAGE_YAML_WITH_INDEXES_A,
		PAGE_YAML_WITH_INDEXES_B,
	} {
		p, err := ReadFrom(strings.NewReader(test), "page/with/index")
		if err != nil {
			t.Fatalf("Failed parsing page: %s", err)
		}

		param := p.GetParam("tags")
		params := param.([]string)

		expected := []string{"a", "b", "c"}
		if !compareStringSlice(params, expected) {
			t.Errorf("Expected %s: got: %s", expected, params)
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
