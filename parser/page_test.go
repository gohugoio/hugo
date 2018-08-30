package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPage(t *testing.T) {
	cases := []struct {
		raw string

		content     string
		frontmatter string
		renderable  bool
		metadata    map[string]interface{}
	}{
		{
			testPageLeader + jsonPageFrontMatter + "\n" + testPageTrailer + jsonPageContent,
			jsonPageContent,
			jsonPageFrontMatter,
			true,
			map[string]interface{}{
				"title": "JSON Test 1",
				"social": []interface{}{
					[]interface{}{"a", "#"},
					[]interface{}{"b", "#"},
				},
			},
		},
		{
			testPageLeader + tomlPageFrontMatter + testPageTrailer + tomlPageContent,
			tomlPageContent,
			tomlPageFrontMatter,
			true,
			map[string]interface{}{
				"title": "TOML Test 1",
				"social": []interface{}{
					[]interface{}{"a", "#"},
					[]interface{}{"b", "#"},
				},
			},
		},
		{
			testPageLeader + yamlPageFrontMatter + testPageTrailer + yamlPageContent,
			yamlPageContent,
			yamlPageFrontMatter,
			true,
			map[string]interface{}{
				"title": "YAML Test 1",
				"social": []interface{}{
					[]interface{}{"a", "#"},
					[]interface{}{"b", "#"},
				},
			},
		},
		{
			testPageLeader + orgPageFrontMatter + orgPageContent,
			orgPageContent,
			orgPageFrontMatter,
			true,
			map[string]interface{}{
				"TITLE":      "Org Test 1",
				"categories": []string{"a", "b"},
			},
		},
	}

	for i, c := range cases {
		p := pageMust(ReadFrom(strings.NewReader(c.raw)))
		meta, err := p.Metadata()

		mesg := fmt.Sprintf("[%d]", i)

		require.Nil(t, err, mesg)
		assert.Equal(t, c.content, string(p.Content()), mesg+" content")
		assert.Equal(t, c.frontmatter, string(p.FrontMatter()), mesg+" frontmatter")
		assert.Equal(t, c.renderable, p.IsRenderable(), mesg+" renderable")
		assert.Equal(t, c.metadata, meta, mesg+" metadata")
	}
}

var (
	testWhitespace  = "\t\t\n\n"
	testPageLeader  = "\ufeff" + testWhitespace + "<!--[metadata]>\n"
	testPageTrailer = "\n<![end-metadata]-->\n"

	jsonPageContent     = "# JSON Test\n"
	jsonPageFrontMatter = `{
	"title": "JSON Test 1",
	"social": [
		["a", "#"],
		["b", "#"]
	]
}`

	tomlPageContent     = "# TOML Test\n"
	tomlPageFrontMatter = `+++
title = "TOML Test 1"
social = [
	["a", "#"],
	["b", "#"],
]
+++
`

	yamlPageContent     = "# YAML Test\n"
	yamlPageFrontMatter = `---
title: YAML Test 1
social:
  - - "a"
    - "#"
  - - "b"
    - "#"
---
`

	orgPageContent     = "* Org Test\n"
	orgPageFrontMatter = `#+TITLE: Org Test 1
#+categories: a b
`

	pageHTMLComment = `<!--
	This is a sample comment.
-->
`
)
