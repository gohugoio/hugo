// Copyright 2018 The Hugo Authors. All rights reserved.
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

package transform_test

import (
	"testing"

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl/transform"

	qt "github.com/frankban/quicktest"
)

func TestRemarshal(t *testing.T) {
	t.Parallel()

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t},
	).Build()

	ns := transform.New(b.H.Deps)
	c := qt.New(t)

	c.Run("Roundtrip variants", func(c *qt.C) {
		tomlExample := `title = 'Test Metadata'
		
[[resources]]
  src = '**image-4.png'
  title = 'The Fourth Image!'
  [resources.params]
    byline = 'picasso'

[[resources]]
  name = 'my-cool-image-:counter'
  src = '**.png'
  title = 'TOML: The Image #:counter'
  [resources.params]
    byline = 'bep'
`

		yamlExample := `resources:
- params:
    byline: picasso
  src: '**image-4.png'
  title: The Fourth Image!
- name: my-cool-image-:counter
  params:
    byline: bep
  src: '**.png'
  title: 'TOML: The Image #:counter'
title: Test Metadata
`

		jsonExample := `{
   "resources": [
      {
         "params": {
            "byline": "picasso"
         },
         "src": "**image-4.png",
         "title": "The Fourth Image!"
      },
      {
         "name": "my-cool-image-:counter",
         "params": {
            "byline": "bep"
         },
         "src": "**.png",
         "title": "TOML: The Image #:counter"
      }
   ],
   "title": "Test Metadata"
}
`
		xmlExample := `<root>
		  <resources>
			<params>
			  <byline>picasso</byline>
			</params>
			<src>**image-4.png</src>
			<title>The Fourth Image!</title>
		  </resources>
		  <resources>
			<name>my-cool-image-:counter</name>
			<params>
			  <byline>bep</byline>
			</params>
			<src>**.png</src>
			<title>TOML: The Image #:counter</title>
		  </resources>
		  <title>Test Metadata</title>
		</root>
		`

		variants := []struct {
			format string
			data   string
		}{
			{"yaml", yamlExample},
			{"json", jsonExample},
			{"toml", tomlExample},
			{"TOML", tomlExample},
			{"Toml", tomlExample},
			{" TOML ", tomlExample},
			{"XML", xmlExample},
		}

		for _, v1 := range variants {
			for _, v2 := range variants {
				// Both from and to may be the same here, but that is fine.
				fromTo := qt.Commentf("%s => %s", v2.format, v1.format)

				converted, err := ns.Remarshal(v1.format, v2.data)
				c.Assert(err, qt.IsNil, fromTo)
				diff := htesting.DiffStrings(v1.data, converted)
				if len(diff) > 0 {
					t.Errorf("[%s] Expected \n%v\ngot\n%v\ndiff:\n%v", fromTo, v1.data, converted, diff)
				}

			}
		}
	})

	c.Run("Comments", func(c *qt.C) {
		input := `
Hugo = "Rules"
		
# It really does!

[m]
# A comment
a = "b"

`

		expected := `Hugo = 'Rules'
[m]
a = 'b'
`

		for _, format := range []string{"json", "yaml", "toml"} {
			fromTo := qt.Commentf("%s => %s", "toml", format)

			converted := input
			var err error
			// Do a round-trip conversion
			for _, toFormat := range []string{format, "toml"} {
				converted, err = ns.Remarshal(toFormat, converted)
				c.Assert(err, qt.IsNil, fromTo)
			}

			diff := htesting.DiffStrings(expected, converted)
			if len(diff) > 0 {
				t.Fatalf("[%s] Expected \n%v\ngot\n>>%v\ndiff:\n%v\n", fromTo, expected, converted, diff)
			}
		}
	})

	// Issue 8850
	c.Run("TOML Indent", func(c *qt.C) {
		input := `

[params]
[params.variables]
a = "b"

`

		converted, err := ns.Remarshal("toml", input)
		c.Assert(err, qt.IsNil)
		c.Assert(converted, qt.Equals, "[params]\n  [params.variables]\n    a = 'b'\n")
	})

	c.Run("Map input", func(c *qt.C) {
		input := map[string]any{
			"hello": "world",
		}

		output, err := ns.Remarshal("toml", input)
		c.Assert(err, qt.IsNil)
		c.Assert(output, qt.Equals, "hello = 'world'\n")
	})

	c.Run("Error", func(c *qt.C) {
		_, err := ns.Remarshal("asdf", "asdf")
		c.Assert(err, qt.Not(qt.IsNil))

		_, err = ns.Remarshal("json", "asdf")
		c.Assert(err, qt.Not(qt.IsNil))
	})
}

func TestRemarshaBillionLaughs(t *testing.T) {
	t.Parallel()

	yamlBillionLaughs := `
a: &a [_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]
b: &b [*a, *a, *a, *a, *a, *a, *a, *a, *a, *a]
c: &c [*b, *b, *b, *b, *b, *b, *b, *b, *b, *b]
d: &d [*c, *c, *c, *c, *c, *c, *c, *c, *c, *c]
e: &e [*d, *d, *d, *d, *d, *d, *d, *d, *d, *d]
f: &f [*e, *e, *e, *e, *e, *e, *e, *e, *e, *e]
g: &g [*f, *f, *f, *f, *f, *f, *f, *f, *f, *f]
h: &h [*g, *g, *g, *g, *g, *g, *g, *g, *g, *g]
i: &i [*h, *h, *h, *h, *h, *h, *h, *h, *h, *h]
`

	yamlMillionLaughs := `
a: &a [_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]
b: &b [*a, *a, *a, *a, *a, *a, *a, *a, *a, *a]
c: &c [*b, *b, *b, *b, *b, *b, *b, *b, *b, *b]
d: &d [*c, *c, *c, *c, *c, *c, *c, *c, *c, *c]
e: &e [*d, *d, *d, *d, *d, *d, *d, *d, *d, *d]
f: &f [*e, *e, *e, *e, *e, *e, *e, *e, *e, *e]
`

	yamlTenThousandLaughs := `
a: &a [_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]
b: &b [*a, *a, *a, *a, *a, *a, *a, *a, *a, *a]
c: &c [*b, *b, *b, *b, *b, *b, *b, *b, *b, *b]
d: &d [*c, *c, *c, *c, *c, *c, *c, *c, *c, *c]

`

	yamlThousandLaughs := `
a: &a [_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]
b: &b [*a, *a, *a, *a, *a, *a, *a, *a, *a, *a]
c: &c [*b, *b, *b, *b, *b, *b, *b, *b, *b, *b]

`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t},
	).Build()

	ns := transform.New(b.H.Deps)

	for _, test := range []struct {
		name string
		data string
	}{
		{"10k", yamlTenThousandLaughs},
		{"1M", yamlMillionLaughs},
		{"1B", yamlBillionLaughs},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)
			_, err := ns.Remarshal("json", test.data)
			c.Assert(err, qt.Not(qt.IsNil))
		})
	}

	// Thousand laughs should be ok.
	// It produces about 29KB of JSON,
	// which is still a large output for such a large input,
	// but there may be use cases for this.
	_, err := ns.Remarshal("json", yamlThousandLaughs)
	c := qt.New(t)
	c.Assert(err, qt.IsNil)
}
