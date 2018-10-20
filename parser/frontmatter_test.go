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

package parser

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestInterfaceToConfig(t *testing.T) {
	cases := []struct {
		input interface{}
		mark  byte
		want  []byte
		isErr bool
	}{
		// TOML
		{map[string]interface{}{}, TOMLLead[0], nil, false},
		{
			map[string]interface{}{"title": "test 1"},
			TOMLLead[0],
			[]byte("title = \"test 1\"\n"),
			false,
		},

		// YAML
		{map[string]interface{}{}, YAMLLead[0], []byte("{}\n"), false},
		{
			map[string]interface{}{"title": "test 1"},
			YAMLLead[0],
			[]byte("title: test 1\n"),
			false,
		},

		// JSON
		{map[string]interface{}{}, JSONLead[0], []byte("{}\n"), false},
		{
			map[string]interface{}{"title": "test 1"},
			JSONLead[0],
			[]byte("{\n   \"title\": \"test 1\"\n}\n"),
			false,
		},

		// Errors
		{nil, TOMLLead[0], nil, true},
		{map[string]interface{}{}, '$', nil, true},
	}

	for i, c := range cases {
		var buf bytes.Buffer

		err := InterfaceToConfig(c.input, rune(c.mark), &buf)
		if err != nil {
			if c.isErr {
				continue
			}
			t.Fatalf("[%d] unexpected error value: %v", i, err)
		}

		if !reflect.DeepEqual(buf.Bytes(), c.want) {
			t.Errorf("[%d] not equal:\nwant %q,\n got %q", i, c.want, buf.Bytes())
		}
	}
}

func TestInterfaceToFrontMatter(t *testing.T) {
	cases := []struct {
		input interface{}
		mark  rune
		want  []byte
		isErr bool
	}{
		// TOML
		{map[string]interface{}{}, '+', []byte("+++\n\n+++\n"), false},
		{
			map[string]interface{}{"title": "test 1"},
			'+',
			[]byte("+++\ntitle = \"test 1\"\n\n+++\n"),
			false,
		},

		// YAML
		{map[string]interface{}{}, '-', []byte("---\n{}\n---\n"), false}, //
		{
			map[string]interface{}{"title": "test 1"},
			'-',
			[]byte("---\ntitle: test 1\n---\n"),
			false,
		},

		// JSON
		{map[string]interface{}{}, '{', []byte("{}\n"), false},
		{
			map[string]interface{}{"title": "test 1"},
			'{',
			[]byte("{\n   \"title\": \"test 1\"\n}\n"),
			false,
		},

		// Errors
		{nil, '+', nil, true},
		{map[string]interface{}{}, '$', nil, true},
	}

	for i, c := range cases {
		var buf bytes.Buffer
		err := InterfaceToFrontMatter(c.input, c.mark, &buf)
		if err != nil {
			if c.isErr {
				continue
			}
			t.Fatalf("[%d] unexpected error value: %v", i, err)
		}

		if !reflect.DeepEqual(buf.Bytes(), c.want) {
			t.Errorf("[%d] not equal:\nwant %q,\n got %q", i, c.want, buf.Bytes())
		}
	}
}

func TestFormatToLeadRune(t *testing.T) {
	for i, this := range []struct {
		kind   string
		expect rune
	}{
		{"yaml", '-'},
		{"yml", '-'},
		{"toml", '+'},
		{"tml", '+'},
		{"json", '{'},
		{"js", '{'},
		{"org", '#'},
		{"unknown", '+'},
	} {
		result := FormatToLeadRune(this.kind)

		if result != this.expect {
			t.Errorf("[%d] got %q but expected %q", i, result, this.expect)
		}
	}
}

func TestRemoveTOMLIdentifier(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"a = 1", "a = 1"},
		{"a = 1\r\n", "a = 1\r\n"},
		{"+++\r\na = 1\r\n+++\r\n", "a = 1\r\n"},
		{"+++\na = 1\n+++\n", "a = 1\n"},
		{"+++\nb = \"+++ oops +++\"\n+++\n", "b = \"+++ oops +++\"\n"},
		{"+++\nc = \"\"\"+++\noops\n+++\n\"\"\"\"\n+++\n", "c = \"\"\"+++\noops\n+++\n\"\"\"\"\n"},
		{"+++\nd = 1\n+++", "d = 1\n"},
	}

	for i, c := range cases {
		res := removeTOMLIdentifier([]byte(c.input))
		if string(res) != c.want {
			t.Errorf("[%d] given %q\nwant: %q\n got: %q", i, c.input, c.want, res)
		}
	}
}

func BenchmarkFrontmatterTags(b *testing.B) {

	for _, frontmatter := range []string{"JSON", "YAML", "YAML2", "TOML"} {
		for i := 1; i < 60; i += 20 {
			doBenchmarkFrontmatter(b, frontmatter, i)
		}
	}
}

func doBenchmarkFrontmatter(b *testing.B, fileformat string, numTags int) {
	yamlTemplate := `---
name: "Tags"
tags:
%s
---
`

	yaml2Template := `---
name: "Tags"
tags: %s
---
`
	tomlTemplate := `+++
name = "Tags"
tags = %s
+++
`

	jsonTemplate := `{
	"name": "Tags",
	"tags": [
		%s
	]
}`
	name := fmt.Sprintf("%s:%d", fileformat, numTags)
	b.Run(name, func(b *testing.B) {
		tags := make([]string, numTags)
		var (
			tagsStr             string
			frontmatterTemplate string
		)
		for i := 0; i < numTags; i++ {
			tags[i] = fmt.Sprintf("Hugo %d", i+1)
		}
		if fileformat == "TOML" {
			frontmatterTemplate = tomlTemplate
			tagsStr = strings.Replace(fmt.Sprintf("%q", tags), " ", ", ", -1)
		} else if fileformat == "JSON" {
			frontmatterTemplate = jsonTemplate
			tagsStr = strings.Replace(fmt.Sprintf("%q", tags), " ", ", ", -1)
		} else if fileformat == "YAML2" {
			frontmatterTemplate = yaml2Template
			tagsStr = strings.Replace(fmt.Sprintf("%q", tags), " ", ", ", -1)
		} else {
			frontmatterTemplate = yamlTemplate
			for _, tag := range tags {
				tagsStr += "\n- " + tag
			}
		}

		frontmatter := fmt.Sprintf(frontmatterTemplate, tagsStr)

		p := page{frontmatter: []byte(frontmatter)}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			meta, err := p.Metadata()
			if err != nil {
				b.Fatal(err)
			}
			if meta == nil {
				b.Fatal("Meta is nil")
			}
		}
	})
}
