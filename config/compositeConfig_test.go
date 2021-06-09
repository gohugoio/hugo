// Copyright 2021 The Hugo Authors. All rights reserved.
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

package config

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCompositeConfig(t *testing.T) {
	c := qt.New(t)

	c.Run("Set and get", func(c *qt.C) {
		base, layer := New(), New()
		cfg := NewCompositeConfig(base, layer)

		layer.Set("a1", "av")
		base.Set("b1", "bv")
		cfg.Set("c1", "cv")

		c.Assert(cfg.Get("a1"), qt.Equals, "av")
		c.Assert(cfg.Get("b1"), qt.Equals, "bv")
		c.Assert(cfg.Get("c1"), qt.Equals, "cv")
		c.Assert(cfg.IsSet("c1"), qt.IsTrue)
		c.Assert(layer.IsSet("c1"), qt.IsTrue)
		c.Assert(base.IsSet("c1"), qt.IsFalse)
	})
}
