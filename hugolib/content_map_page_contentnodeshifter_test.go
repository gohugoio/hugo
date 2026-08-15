// Copyright 2025 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
)

type testContentNodeForSite struct {
	vec sitesmatrix.Vector
}

func (n *testContentNodeForSite) Path() string                   { return "/test" }
func (n *testContentNodeForSite) nodeCategorySingle()            {}
func (n *testContentNodeForSite) siteVector() sitesmatrix.Vector { return n.vec }
func (n *testContentNodeForSite) forEeachContentNode(f func(sitesmatrix.Vector, contentNode) bool) bool {
	return f(n.vec, n)
}

// See issue 15207.
func TestContentNodeShifterDeleteMultipleFromNodes(t *testing.T) {
	c := qt.New(t)

	s := &contentNodeShifter{}
	vec := sitesmatrix.Vector{2, 0, 0}
	other := sitesmatrix.Vector{1, 0, 0}

	newNodes := func() contentNodes {
		return contentNodes{
			&testContentNodeForSite{vec: other},
			&testContentNodeForSite{vec: vec},
			&testContentNodeForSite{vec: other},
			&testContentNodeForSite{vec: other},
			&testContentNodeForSite{vec: vec},
		}
	}

	updated, deleted, wasDeleted, isEmpty := s.Delete(newNodes(), vec)
	c.Assert(wasDeleted, qt.IsTrue)
	c.Assert(isEmpty, qt.IsFalse)
	c.Assert(deleted.(contentNodes), qt.HasLen, 2)
	remaining := updated.(contentNodes)
	c.Assert(remaining, qt.HasLen, 3)
	for _, n := range remaining {
		c.Assert(n.(contentNodeForSite).siteVector(), qt.Equals, other)
	}

	// Delete all.
	updated, deleted, wasDeleted, isEmpty = s.Delete(contentNodes{
		&testContentNodeForSite{vec: vec},
		&testContentNodeForSite{vec: vec},
	}, vec)
	c.Assert(wasDeleted, qt.IsTrue)
	c.Assert(isEmpty, qt.IsTrue)
	c.Assert(deleted.(contentNodes), qt.HasLen, 2)
	c.Assert(updated, qt.IsNil)

	// Delete none.
	nodes := contentNodes{&testContentNodeForSite{vec: other}}
	updated, deleted, wasDeleted, isEmpty = s.Delete(nodes, vec)
	c.Assert(wasDeleted, qt.IsFalse)
	c.Assert(isEmpty, qt.IsFalse)
	c.Assert(deleted, qt.IsNil)
	c.Assert(updated.(contentNodes), qt.HasLen, 1)
}
