// Copyright 2024 The Hugo Authors. All rights reserved.
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

package httpcache

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGlobMatcher(t *testing.T) {
	c := qt.New(t)

	g := GlobMatcher{
		Includes: []string{"**/*.jpg", "**.png", "**/bar/**"},
		Excludes: []string{"**/foo.jpg", "**.css"},
	}

	p, err := g.CompilePredicate()
	c.Assert(err, qt.IsNil)

	c.Assert(p("foo.jpg"), qt.IsFalse)
	c.Assert(p("foo.png"), qt.IsTrue)
	c.Assert(p("foo/bar.jpg"), qt.IsTrue)
	c.Assert(p("foo/bar.png"), qt.IsTrue)
	c.Assert(p("foo/bar/foo.jpg"), qt.IsFalse)
	c.Assert(p("foo/bar/foo.css"), qt.IsFalse)
	c.Assert(p("foo.css"), qt.IsFalse)
	c.Assert(p("foo/bar/foo.css"), qt.IsFalse)
	c.Assert(p("foo/bar/foo.xml"), qt.IsTrue)
}
