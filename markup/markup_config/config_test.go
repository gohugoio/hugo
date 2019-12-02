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

	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
)

func TestConfig(t *testing.T) {
	c := qt.New(t)

	c.Run("Decode", func(c *qt.C) {
		c.Parallel()
		v := viper.New()

		v.Set("markup", map[string]interface{}{
			"goldmark": map[string]interface{}{
				"renderer": map[string]interface{}{
					"unsafe": true,
				},
			},
		})

		conf, err := Decode(v)

		c.Assert(err, qt.IsNil)
		c.Assert(conf.Goldmark.Renderer.Unsafe, qt.Equals, true)
		c.Assert(conf.BlackFriday.Fractions, qt.Equals, true)

	})

	c.Run("legacy", func(c *qt.C) {
		c.Parallel()
		v := viper.New()

		v.Set("blackfriday", map[string]interface{}{
			"angledQuotes": true,
		})

		v.Set("footnoteAnchorPrefix", "myprefix")
		v.Set("footnoteReturnLinkContents", "myreturn")
		v.Set("pygmentsStyle", "hugo")
		v.Set("pygmentsCodefencesGuessSyntax", true)
		conf, err := Decode(v)

		c.Assert(err, qt.IsNil)
		c.Assert(conf.BlackFriday.AngledQuotes, qt.Equals, true)
		c.Assert(conf.BlackFriday.FootnoteAnchorPrefix, qt.Equals, "myprefix")
		c.Assert(conf.BlackFriday.FootnoteReturnLinkContents, qt.Equals, "myreturn")
		c.Assert(conf.Highlight.Style, qt.Equals, "hugo")
		c.Assert(conf.Highlight.CodeFences, qt.Equals, true)
		c.Assert(conf.Highlight.GuessSyntax, qt.Equals, true)
	})

}
