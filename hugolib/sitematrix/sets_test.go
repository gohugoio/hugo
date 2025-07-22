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

package sitematrix_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/hugolib/sitematrix"
)

func TestIntSets(t *testing.T) {
	c := qt.New(t)

	sets := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets2 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 4),
	).Build()

	c.Assert(hashing.HashStringHex(sets), qt.Equals, "fb6349a290fa5920")
	c.Assert(hashing.HashStringHex(sets2), qt.Equals, "5218dc8d5d3b944")

	c.Assert(sets.HasVector(sitematrix.Vector{1, 2, 3}), qt.Equals, true)
	c.Assert(sets.HasVector(sitematrix.Vector{3, 2, 3}), qt.Equals, false)
	c.Assert(sets.FirstVector(), qt.Equals, sitematrix.Vector{1, 1, 1})
	c.Assert(sets.EqualsVector(sets), qt.Equals, true)
	c.Assert(sets.EqualsVector(
		sitematrix.NewIntSetsBuilder(0).WithSets(
			maps.NewOrderedIntSet(1, 2),
			maps.NewOrderedIntSet(1, 2, 3),
			maps.NewOrderedIntSet(1, 2, 3),
		).Build(),
	), qt.Equals, true)
	c.Assert(sets.EqualsVector(
		sitematrix.NewIntSetsBuilder(0).WithSets(
			maps.NewOrderedIntSet(1, 2, 3),
			maps.NewOrderedIntSet(1, 2, 3, 4),
			maps.NewOrderedIntSet(1, 2, 3, 4),
		).Build(),
	), qt.Equals, false)

	c.Assert(sets.EqualsVector(
		sitematrix.NewIntSetsBuilder(0).WithSets(
			maps.NewOrderedIntSet(1, 2),
			maps.NewOrderedIntSet(1, 2, 3),
			maps.NewOrderedIntSet(2, 3, 4),
		).Build(),
	), qt.Equals, false)

	allCount := 0
	seen := make(map[sitematrix.Vector]bool)
	ok := sets.ForEeachVector(func(v sitematrix.Vector) bool {
		c.Assert(seen[v], qt.IsFalse)
		seen[v] = true
		allCount++
		return true
	})

	c.Assert(ok, qt.IsTrue)

	// 2 languages * 3 versions * 3 roles = 18 combinations.
	c.Assert(allCount, qt.Equals, 18)
}

func TestIntSetsComplement(t *testing.T) {
	c := qt.New(t)

	sets1 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets1Copy := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets2 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4, 6),
	).Build()

	c1 := sets2.Complement(sets1)
	k1, k2, k3 := c1.KeysSorted()
	c.Assert(k1, qt.DeepEquals, []int{4})
	c.Assert(k2, qt.DeepEquals, []int{4})
	c.Assert(k3, qt.DeepEquals, []int{4, 6})
	c.Assert(hashing.HashStringHex(c1), qt.Not(qt.Equals), hashing.HashStringHex(sets1))
	c.Assert(hashing.HashStringHex(c1), qt.Not(qt.Equals), hashing.HashStringHex(sets2))
	c2 := sets1.Complement(sets1Copy)
	c.Assert(c2, qt.IsNil)
}

func BenchmarkIntSetsComplement(b *testing.B) {
	sets1 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets1Copy := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	setsLanguage1 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1),
		nil,
		nil,
	).Build()

	sets2 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4, 6),
	).Build()

	b.ResetTimer()
	b.Run("two different sets, some overlap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sets2.Complement(sets1)
		}
	})

	b.Run("self", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sets1.Complement(sets1)
		}
	})

	b.Run("same", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sets1.Complement(sets1Copy)
		}
	})

	b.Run("one overlapping language", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = setsLanguage1.Complement(sets1)
		}
	})
}

func BenchmarkSets(b *testing.B) {
	sets1 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets1Copy := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets2 := sitematrix.NewIntSetsBuilder(0).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4),
	).Build()

	v1 := sitematrix.Vector{1, 2, 3}
	v2 := sitematrix.Vector{3, 2, 3}

	b.ResetTimer()
	b.Run("HasVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sets1.HasVector(v1)
			_ = sets1.HasVector(v2)
		}
	})

	b.Run("FirstVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sets1.FirstVector()
		}
	})

	b.Run("LenVectors", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sets1.LenVectors()
		}
	})

	b.Run("ForEeachVector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			allCount := 0
			ok := sets1.ForEeachVector(func(v sitematrix.Vector) bool {
				allCount++
				_ = v
				return true
			})

			if !ok {
				b.Fatal("Expected ForEeachVector to return true")
			}

			if allCount != 18 {
				b.Fatalf("Expected 18 combinations, got %d", allCount)
			}
		}
	})

	b.Run("EqualsVector pointer equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !sets1.EqualsVector(sets1) {
				b.Fatal("Expected sets1 to equal itself")
			}
		}
	})

	b.Run("EqualsVector equal copy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !sets1.EqualsVector(sets1Copy) {
				b.Fatal("Expected sets1 to equal its copy")
			}
		}
	})

	b.Run("EqualsVector different sets", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if sets1.EqualsVector(sets2) {
				b.Fatal("Expected sets1 to not equal sets2")
			}
		}
	})

	b.Run("Distance", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sets1.FirstVector().Distance(v1)
			_ = sets1.FirstVector().Distance(v2)
		}
	})
}
