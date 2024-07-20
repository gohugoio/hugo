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

package markup_config

import (
	"testing"

	"github.com/gohugoio/hugo/config"

	qt "github.com/frankban/quicktest"
)

func TestConfig(t *testing.T) {
	c := qt.New(t)

	c.Run("Decode", func(c *qt.C) {
		c.Parallel()
		v := config.New()

		v.Set("markup", map[string]any{
			"goldmark": map[string]any{
				"renderer": map[string]any{
					"unsafe": true,
				},
			},
			"asciidocext": map[string]any{
				"workingFolderCurrent": true,
				"safeMode":             "save",
				"extensions":           []string{"asciidoctor-html5s"},
			},
		})

		conf, err := Decode(v)

		c.Assert(err, qt.IsNil)
		c.Assert(conf.Goldmark.Renderer.Unsafe, qt.Equals, true)
		c.Assert(conf.Goldmark.Parser.Attribute.Title, qt.Equals, true)
		c.Assert(conf.Goldmark.Parser.Attribute.Block, qt.Equals, false)

		c.Assert(conf.AsciidocExt.WorkingFolderCurrent, qt.Equals, true)
		c.Assert(conf.AsciidocExt.Extensions[0], qt.Equals, "asciidoctor-html5s")
	})

	c.Run("Decode legacy typographer", func(c *qt.C) {
		c.Parallel()
		v := config.New()

		// typographer was changed from a bool to a struct in 0.112.0.
		v.Set("markup", map[string]any{
			"goldmark": map[string]any{
				"extensions": map[string]any{
					"typographer": false,
				},
			},
		})

		conf, err := Decode(v)

		c.Assert(err, qt.IsNil)
		c.Assert(conf.Goldmark.Extensions.Typographer.Disable, qt.Equals, true)

		v.Set("markup", map[string]any{
			"goldmark": map[string]any{
				"extensions": map[string]any{
					"typographer": true,
				},
			},
		})

		conf, err = Decode(v)

		c.Assert(err, qt.IsNil)
		c.Assert(conf.Goldmark.Extensions.Typographer.Disable, qt.Equals, false)
		c.Assert(conf.Goldmark.Extensions.Typographer.Ellipsis, qt.Equals, "&hellip;")
	})
}
