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

// Package highlight provides code highlighting.
package highlight

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
)

func TestConfig(t *testing.T) {
	c := qt.New(t)

	c.Run("applyLegacyConfig", func(c *qt.C) {
		v := config.New()
		v.Set("pygmentsStyle", "hugo")
		v.Set("pygmentsUseClasses", false)
		v.Set("pygmentsCodeFences", false)
		v.Set("pygmentsOptions", "linenos=inline")

		cfg := DefaultConfig
		err := ApplyLegacyConfig(v, &cfg)
		c.Assert(err, qt.IsNil)
		c.Assert(cfg.Style, qt.Equals, "hugo")
		c.Assert(cfg.NoClasses, qt.Equals, true)
		c.Assert(cfg.CodeFences, qt.Equals, false)
		c.Assert(cfg.LineNos, qt.Equals, true)
		c.Assert(cfg.LineNumbersInTable, qt.Equals, false)
	})

	c.Run("parseOptions", func(c *qt.C) {
		cfg := DefaultConfig
		opts := "noclasses=true,linenos=inline,linenostart=32,hl_lines=3-8 10-20"
		err := applyOptionsFromString(opts, &cfg)

		c.Assert(err, qt.IsNil)
		c.Assert(cfg.NoClasses, qt.Equals, true)
		c.Assert(cfg.LineNos, qt.Equals, true)
		c.Assert(cfg.LineNumbersInTable, qt.Equals, false)
		c.Assert(cfg.LineNoStart, qt.Equals, 32)
		c.Assert(cfg.Hl_Lines, qt.Equals, "3-8 10-20")
	})

	c.Run("applyOptionsFromMap", func(c *qt.C) {
		cfg := DefaultConfig
		err := applyOptionsFromMap(map[string]any{
			"noclasses":   true,
			"lineNos":     "inline", // mixed case key, should work after normalization
			"linenostart": 32,
			"hl_lines":    "3-8 10-20",
		}, &cfg)

		c.Assert(err, qt.IsNil)
		c.Assert(cfg.NoClasses, qt.Equals, true)
		c.Assert(cfg.LineNos, qt.Equals, true)
		c.Assert(cfg.LineNumbersInTable, qt.Equals, false)
		c.Assert(cfg.LineNoStart, qt.Equals, 32)
		c.Assert(cfg.Hl_Lines, qt.Equals, "3-8 10-20")
	})
}
