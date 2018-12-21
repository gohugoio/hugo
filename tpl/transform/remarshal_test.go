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

package transform

import (
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestRemarshal(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))
	assert := require.New(t)

	tomlExample := `title = "Test Metadata"
		
[[resources]]
  src = "**image-4.png"
  title = "The Fourth Image!"
  [resources.params]
    byline = "picasso"

[[resources]]
  name = "my-cool-image-:counter"
  src = "**.png"
  title = "TOML: The Image #:counter"
  [resources.params]
    byline = "bep"
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
	}

	for _, v1 := range variants {
		for _, v2 := range variants {
			// Both from and to may be the same here, but that is fine.
			fromTo := fmt.Sprintf("%s => %s", v2.format, v1.format)

			converted, err := ns.Remarshal(v1.format, v2.data)
			assert.NoError(err, fromTo)
			diff := helpers.DiffStrings(v1.data, converted)
			if len(diff) > 0 {
				t.Errorf("[%s] Expected \n%v\ngot\n%v\ndiff:\n%v", fromTo, v1.data, converted, diff)
			}

		}
	}

}

func TestRemarshalComments(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))

	assert := require.New(t)

	input := `
Hugo = "Rules"
		
# It really does!

[m]
# A comment
a = "b"

`

	expected := `
Hugo = "Rules"
		
[m]
  a = "b"
`

	for _, format := range []string{"json", "yaml", "toml"} {
		fromTo := fmt.Sprintf("%s => %s", "toml", format)

		converted := input
		var err error
		// Do a round-trip conversion
		for _, toFormat := range []string{format, "toml"} {
			converted, err = ns.Remarshal(toFormat, converted)
			assert.NoError(err, fromTo)
		}

		diff := helpers.DiffStrings(expected, converted)
		if len(diff) > 0 {
			t.Fatalf("[%s] Expected \n%v\ngot\n%v\ndiff:\n%v\n", fromTo, expected, converted, diff)
		}
	}
}

func TestTestRemarshalError(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))
	assert := require.New(t)

	_, err := ns.Remarshal("asdf", "asdf")
	assert.Error(err)

	_, err = ns.Remarshal("json", "asdf")
	assert.Error(err)

}
