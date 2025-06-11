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

package sitematrix

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/maps"
)

func TestIntSetsVectorProvider(t *testing.T) {
	c := qt.New(t)

	sets := &IntSets{
		Languages: maps.NewOrderedIntSet(1, 2),
		Versions:  maps.NewOrderedIntSet(1, 2, 3),
		Roles:     maps.NewOrderedIntSet(1, 2, 3),
	}

	c.Assert(sets.HasVector(Vector{1, 2, 3}), qt.Equals, true)
	c.Assert(sets.HasVector(Vector{3, 2, 3}), qt.Equals, false)
	c.Assert(sets.FirstVector(), qt.Equals, Vector{1, 1, 1})
	c.Assert(sets.EqualsVector(sets), qt.Equals, true)
	c.Assert(sets.EqualsVector(&IntSets{
		Languages: maps.NewOrderedIntSet(1, 2),
		Versions:  maps.NewOrderedIntSet(1, 2, 3),
		Roles:     maps.NewOrderedIntSet(1, 2, 3),
	}), qt.Equals, true)
	c.Assert(sets.EqualsVector(&IntSets{
		Languages: maps.NewOrderedIntSet(1, 2, 3),
		Versions:  maps.NewOrderedIntSet(1, 2, 3, 4),
		Roles:     maps.NewOrderedIntSet(1, 2, 3, 4),
	}), qt.Equals, false)

	c.Assert(sets.EqualsVector(&IntSets{
		Languages: maps.NewOrderedIntSet(1, 2),
		Versions:  maps.NewOrderedIntSet(1, 2, 3),
		Roles:     maps.NewOrderedIntSet(2, 3, 4),
	}), qt.Equals, false)

	alllCount := 0
	seen := make(map[Vector]bool)
	ok := sets.ForEeachVector(func(v Vector) bool {
		c.Assert(seen[v], qt.IsFalse)
		seen[v] = true
		alllCount++
		return true
	})

	c.Assert(ok, qt.IsTrue)

	// 2 languages * 3 versions * 3 roles = 18 combinations.
	c.Assert(alllCount, qt.Equals, 18)
}

func BenchmarkSets(b *testing.B) {
	sets1 := &IntSets{
		Languages: maps.NewOrderedIntSet(1, 2),
		Versions:  maps.NewOrderedIntSet(1, 2, 3),
		Roles:     maps.NewOrderedIntSet(1, 2, 3),
	}

	sets1Copy := &IntSets{
		Languages: maps.NewOrderedIntSet(1, 2),
		Versions:  maps.NewOrderedIntSet(1, 2, 3),
		Roles:     maps.NewOrderedIntSet(1, 2, 3),
	}

	sets2 := &IntSets{
		Languages: maps.NewOrderedIntSet(1, 2, 3),
		Versions:  maps.NewOrderedIntSet(1, 2, 3, 4),
		Roles:     maps.NewOrderedIntSet(1, 2, 3, 4),
	}

	v1 := Vector{1, 2, 3}
	v2 := Vector{3, 2, 3}

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
			ok := sets1.ForEeachVector(func(v Vector) bool {
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

	b.Run("SetFrom", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newSets := &IntSets{}
			newSets.SetFrom(sets1)
			_ = newSets
		}
	})
}
