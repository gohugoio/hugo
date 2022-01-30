// Copyright 2020 The Hugo Authors. All rights reserved.
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

package page

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPageMatcher(t *testing.T) {
	c := qt.New(t)

	p1, p2, p3 := &testPage{path: "/p1", kind: "section", lang: "en"}, &testPage{path: "p2", kind: "page", lang: "no"}, &testPage{path: "p3", kind: "page", lang: "en"}

	c.Run("Matches", func(c *qt.C) {
		m := PageMatcher{Kind: "section"}

		c.Assert(m.Matches(p1), qt.Equals, true)
		c.Assert(m.Matches(p2), qt.Equals, false)

		m = PageMatcher{Kind: "page"}
		c.Assert(m.Matches(p1), qt.Equals, false)
		c.Assert(m.Matches(p2), qt.Equals, true)
		c.Assert(m.Matches(p3), qt.Equals, true)

		m = PageMatcher{Kind: "page", Path: "/p2"}
		c.Assert(m.Matches(p1), qt.Equals, false)
		c.Assert(m.Matches(p2), qt.Equals, true)
		c.Assert(m.Matches(p3), qt.Equals, false)

		m = PageMatcher{Path: "/p*"}
		c.Assert(m.Matches(p1), qt.Equals, true)
		c.Assert(m.Matches(p2), qt.Equals, true)
		c.Assert(m.Matches(p3), qt.Equals, true)

		m = PageMatcher{Lang: "en"}
		c.Assert(m.Matches(p1), qt.Equals, true)
		c.Assert(m.Matches(p2), qt.Equals, false)
		c.Assert(m.Matches(p3), qt.Equals, true)
	})

	c.Run("Decode", func(c *qt.C) {
		var v PageMatcher
		c.Assert(DecodePageMatcher(map[string]interface{}{"kind": "foo"}, &v), qt.Not(qt.IsNil))
		c.Assert(DecodePageMatcher(map[string]interface{}{"kind": "{foo,bar}"}, &v), qt.Not(qt.IsNil))
		c.Assert(DecodePageMatcher(map[string]interface{}{"kind": "taxonomy"}, &v), qt.IsNil)
		c.Assert(DecodePageMatcher(map[string]interface{}{"kind": "{taxonomy,foo}"}, &v), qt.IsNil)
		c.Assert(DecodePageMatcher(map[string]interface{}{"kind": "{taxonomy,term}"}, &v), qt.IsNil)
		c.Assert(DecodePageMatcher(map[string]interface{}{"kind": "*"}, &v), qt.IsNil)
		c.Assert(DecodePageMatcher(map[string]interface{}{"kind": "home", "path": filepath.FromSlash("/a/b/**")}, &v), qt.IsNil)
		c.Assert(v, qt.Equals, PageMatcher{Kind: "home", Path: "/a/b/**"})
	})
}
