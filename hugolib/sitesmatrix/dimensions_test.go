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

package sitesmatrix

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDimensionsCompare(t *testing.T) {
	c := qt.New(t)

	c.Assert(Vector{1, 2, 3}.Compare(Vector{1, 2, 8}), qt.Equals, -1)
	c.Assert(Vector{1, 2, 3}.Compare(Vector{1, 2, 3}), qt.Equals, 0)
	c.Assert(Vector{1, 2, 3}.Compare(Vector{1, 2, 0}), qt.Equals, 1)
	c.Assert(Vector{1, 2, 3}.Compare(Vector{1, 0, 3}), qt.Equals, 1)
	c.Assert(Vector{1, 2, 3}.Compare(Vector{0, 3, 2}), qt.Equals, 1)
	c.Assert(Vector{1, 2, 3}.Compare(Vector{0, 0, 0}), qt.Equals, 1)
	c.Assert(Vector{0, 0, 0}.Compare(Vector{1, 2, 3}), qt.Equals, -1)
	c.Assert(Vector{0, 0, 0}.Compare(Vector{0, 0, 0}), qt.Equals, 0)
	c.Assert(Vector{0, 0, 0}.Compare(Vector{1, 0, 0}), qt.Equals, -1)
	c.Assert(Vector{0, 0, 0}.Compare(Vector{0, 1, 0}), qt.Equals, -1)
	c.Assert(Vector{0, 0, 0}.Compare(Vector{0, 0, 1}), qt.Equals, -1)
}

func TestDimensionsDistance(t *testing.T) {
	c := qt.New(t)

	c.Assert(Vector{1, 2, 3}.Distance(Vector{1, 2, 8}), qt.Equals, -5)
	c.Assert(Vector{1, 2, 3}.Distance(Vector{1, 2, 3}), qt.Equals, 0)
	c.Assert(Vector{1, 2, 3}.Distance(Vector{1, 2, 0}), qt.Equals, 3)
	c.Assert(Vector{1, 2, 3}.Distance(Vector{1, 0, 3}), qt.Equals, 2)
	c.Assert(Vector{1, 2, 3}.Distance(Vector{0, 3, 2}), qt.Equals, 1)
}

func BenchmarkCompare(b *testing.B) {
	b.Run("Equal", func(b *testing.B) {
		v1 := Vector{1, 2, 3}
		v2 := Vector{1, 2, 3}
		for b.Loop() {
			v1.Compare(v2)
		}
	})

	b.Run("First different", func(b *testing.B) {
		v1 := Vector{1, 2, 3}
		v2 := Vector{2, 2, 3}
		for b.Loop() {
			v1.Compare(v2)
		}
	})

	b.Run("Last different", func(b *testing.B) {
		v1 := Vector{1, 2, 3}
		v2 := Vector{1, 2, 4}
		for b.Loop() {
			v1.Compare(v2)
		}
	})
}
