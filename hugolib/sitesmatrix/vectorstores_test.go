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

package sitesmatrix_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
)

func newTestDims() *sitesmatrix.ConfiguredDimensions {
	return sitesmatrix.NewTestingDimensions([]string{"en", "no"}, []string{"v1", "v2", "v3"}, []string{"admin", "editor", "viewer", "guest"})
}

func TestIntSets(t *testing.T) {
	c := qt.New(t)
	testDims := newTestDims()

	sets := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets2 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 4),
	).Build()

	sets3 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(2, 1),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(2, 1, 4),
	).Build()

	sets4 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(3, 4, 5),
		maps.NewOrderedIntSet(3, 4, 5),
		maps.NewOrderedIntSet(3, 4, 5),
	).Build()

	c.Assert(hashing.HashStringHex(sets), qt.Equals, "790f6004619934ff")
	c.Assert(hashing.HashStringHex(sets2), qt.Equals, "99abddc51cd22f24")
	c.Assert(hashing.HashStringHex(sets3), qt.Equals, "99abddc51cd22f24")

	c.Assert(sets.HasVector(sitesmatrix.Vector{1, 2, 3}), qt.Equals, true)
	c.Assert(sets.HasVector(sitesmatrix.Vector{3, 2, 3}), qt.Equals, false)
	c.Assert(sets.HasAnyVector(sets2), qt.Equals, true)
	c.Assert(sets.HasAnyVector(sets3), qt.Equals, true)
	c.Assert(sets.HasAnyVector(sets4), qt.Equals, false)

	c.Assert(sets.VectorSample(), qt.Equals, sitesmatrix.Vector{1, 1, 1})
	c.Assert(sets.EqualsVector(sets), qt.Equals, true)
	c.Assert(sets.EqualsVector(
		sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
			maps.NewOrderedIntSet(1, 2),
			maps.NewOrderedIntSet(1, 2, 3),
			maps.NewOrderedIntSet(1, 2, 3),
		).Build(),
	), qt.Equals, true)
	c.Assert(sets.EqualsVector(
		sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
			maps.NewOrderedIntSet(1, 2, 3),
			maps.NewOrderedIntSet(1, 2, 3, 4),
			maps.NewOrderedIntSet(1, 2, 3, 4),
		).Build(),
	), qt.Equals, false)

	c.Assert(sets.EqualsVector(
		sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
			maps.NewOrderedIntSet(1, 2),
			maps.NewOrderedIntSet(1, 2, 3),
			maps.NewOrderedIntSet(2, 3, 4),
		).Build(),
	), qt.Equals, false)

	allCount := 0
	seen := make(map[sitesmatrix.Vector]bool)
	ok := sets.ForEachVector(func(v sitesmatrix.Vector) bool {
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

	type values struct {
		v0 []int
		v1 []int
		v2 []int
	}

	type test struct {
		name     string
		left     values
		input    []values
		expected []sitesmatrix.Vector
	}

	runOne := func(c *qt.C, test test) {
		c.Helper()

		testDims := newTestDims()

		self := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
			maps.NewOrderedIntSet(test.left.v0...),
			maps.NewOrderedIntSet(test.left.v1...),
			maps.NewOrderedIntSet(test.left.v2...),
		).Build()

		var input []sitesmatrix.VectorProvider
		for _, v := range test.input {
			input = append(input, sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
				maps.NewOrderedIntSet(v.v0...),
				maps.NewOrderedIntSet(v.v1...),
				maps.NewOrderedIntSet(v.v2...),
			).Build())
		}

		result := self.Complement(input...)

		vectors := result.Vectors()
		c.Assert(vectors, qt.DeepEquals, test.expected)

		c.Assert(hashing.HashStringHex(result), qt.Not(qt.Equals), hashing.HashStringHex(self))
	}

	for _, test := range []test{
		{
			name: "Test 3",
			left: values{
				v0: []int{1, 3, 5},
				v1: []int{1},
				v2: []int{1},
			},
			input: []values{
				{v0: []int{2}, v1: []int{1}, v2: []int{1}},
				{v0: []int{2, 3}, v1: []int{1}, v2: []int{1}},
			},
			expected: []sitesmatrix.Vector{
				{1, 1, 1},
				{5, 1, 1},
			},
		},
		{
			name: "Test 1",
			left: values{
				v0: []int{1},
				v1: []int{1},
				v2: []int{1},
			},
			input: []values{
				{v0: []int{2}, v1: []int{1}, v2: []int{1}},
			},
			expected: []sitesmatrix.Vector{
				{1, 1, 1},
			},
		},
		{
			name: "Same values",
			left: values{
				v0: []int{1},
				v1: []int{1},
				v2: []int{1},
			},
			input: []values{
				{v0: []int{1}, v1: []int{1}, v2: []int{1}},
			},
			expected: nil,
		},
		{
			name: "Test 2",
			left: values{
				v0: []int{1, 3, 5},
				v1: []int{1},
				v2: []int{1},
			},
			input: []values{
				{v0: []int{2}, v1: []int{1}, v2: []int{1}},
			},
			expected: []sitesmatrix.Vector{
				{1, 1, 1},
				{3, 1, 1},
				{5, 1, 1},
			},
		},
		{
			name: "Many",
			left: values{
				v0: []int{1, 3, 5, 6, 7},
				v1: []int{1, 3, 5, 6, 7},
				v2: []int{1, 3, 5, 6, 7},
			},
			input: []values{
				{
					v0: []int{1, 3, 5, 6, 7},
					v1: []int{1, 3, 5, 6},
					v2: []int{1, 3, 5, 6, 7},
				},
			},
			expected: []sitesmatrix.Vector{
				{1, 7, 1},
				{1, 7, 3},
				{1, 7, 5},
				{1, 7, 6},
				{1, 7, 7},
				{3, 7, 1},
				{3, 7, 3},
				{3, 7, 5},
				{3, 7, 6},
				{3, 7, 7},
				{5, 7, 1},
				{5, 7, 3},
				{5, 7, 5},
				{5, 7, 6},
				{5, 7, 7},
				{6, 7, 1},
				{6, 7, 3},
				{6, 7, 5},
				{6, 7, 6},
				{6, 7, 7},
				{7, 7, 1},
				{7, 7, 3},
				{7, 7, 5},
				{7, 7, 6},
				{7, 7, 7},
			},
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			runOne(c, test)
		})
	}
}

func TestIntSetsComplementOfComplement(t *testing.T) {
	c := qt.New(t)

	testDims := newTestDims()

	sets1 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1),
		maps.NewOrderedIntSet(1),
		maps.NewOrderedIntSet(1),
	).Build()

	sets2 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1),
		maps.NewOrderedIntSet(1, 3),
	).Build()

	c1 := sets2.Complement(sets1)

	vectors := c1.Vectors()
	//	c.Assert(len(vectors), qt.Equals, 3)
	c.Assert(vectors, qt.DeepEquals, []sitesmatrix.Vector{
		{1, 1, 3},
		{2, 1, 1},
		{2, 1, 3},
	})

	c.Assert(hashing.HashStringHex(c1), qt.Not(qt.Equals), hashing.HashStringHex(sets1))
	c.Assert(hashing.HashStringHex(c1), qt.Not(qt.Equals), hashing.HashStringHex(sets2))

	c2 := sets1.Complement(c1)
	c.Assert(c2.Vectors(), qt.DeepEquals, []sitesmatrix.Vector{
		{1, 1, 1},
	})
}

func BenchmarkIntSetsComplement(b *testing.B) {
	testDims := newTestDims()

	sets1 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets1Copy := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets2 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4, 6),
	).Build()

	setsLanguage1 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1),
		maps.NewOrderedIntSet(1),
		maps.NewOrderedIntSet(1),
	).Build()

	b.ResetTimer()
	b.Run("sub set", func(b *testing.B) {
		for b.Loop() {
			_ = sets2.Complement(sets1)
		}
	})

	b.Run("super set", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.Complement(sets2)
		}
	})

	b.Run("self", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.Complement(sets1)
		}
	})

	b.Run("self multiple", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.Complement(sets1, sets1, sets1, sets1, sets1, sets1, sets1)
		}
	})

	b.Run("same", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.Complement(sets1Copy)
		}
	})

	b.Run("one overlapping language", func(b *testing.B) {
		for b.Loop() {
			_ = setsLanguage1.Complement(sets1)
		}
	})
}

func BenchmarkSets(b *testing.B) {
	testDims := newTestDims()

	sets1 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets1Copy := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2),
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3),
	).Build()

	sets2 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(1, 2, 3),
		maps.NewOrderedIntSet(1, 2, 3, 4),
		maps.NewOrderedIntSet(1, 2, 3, 4),
	).Build()

	v1 := sitesmatrix.Vector{1, 2, 3}
	v2 := sitesmatrix.Vector{3, 2, 3}

	b.ResetTimer()
	b.Run("Build", func(b *testing.B) {
		for b.Loop() {
			_ = sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
				maps.NewOrderedIntSet(1, 2),
				maps.NewOrderedIntSet(1, 2, 3),
				maps.NewOrderedIntSet(1, 2, 3),
			).Build()
		}
	})
	b.Run("HasVector", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.HasVector(v1)
			_ = sets1.HasVector(v2)
		}
	})

	b.Run("HasAnyVector(Sets)", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.HasAnyVector(sets2)
		}
	})
	b.Run("HasAnyVector(Vector)", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.HasAnyVector(v1)
		}
	})

	b.Run("HasAnyVector(&Vector)", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.HasAnyVector(&v1)
		}
	})

	b.Run("FirstVector", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.VectorSample()
		}
	})

	b.Run("LenVectors", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.LenVectors()
		}
	})

	b.Run("ForEachVector", func(b *testing.B) {
		for b.Loop() {
			allCount := 0
			ok := sets1.ForEachVector(func(v sitesmatrix.Vector) bool {
				allCount++
				_ = v
				return true
			})

			if !ok {
				b.Fatal("Expected ForEachVector to return true")
			}

			if allCount != 18 {
				b.Fatalf("Expected 18 combinations, got %d", allCount)
			}
		}
	})

	b.Run("EqualsVector pointer equal", func(b *testing.B) {
		for b.Loop() {
			if !sets1.EqualsVector(sets1) {
				b.Fatal("Expected sets1 to equal itself")
			}
		}
	})

	b.Run("EqualsVector equal copy", func(b *testing.B) {
		for b.Loop() {
			if !sets1.EqualsVector(sets1Copy) {
				b.Fatal("Expected sets1 to equal its copy")
			}
		}
	})

	b.Run("EqualsVector different sets", func(b *testing.B) {
		for b.Loop() {
			if sets1.EqualsVector(sets2) {
				b.Fatal("Expected sets1 to not equal sets2")
			}
		}
	})

	b.Run("Distance", func(b *testing.B) {
		for b.Loop() {
			_ = sets1.VectorSample().Distance(v1)
			_ = sets1.VectorSample().Distance(v2)
		}
	})
}

func TestVectorStoreMap(t *testing.T) {
	c := qt.New(t)
	testDims := newTestDims()

	c.Run("Complement", func(c *qt.C) {
		v1 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
			maps.NewOrderedIntSet(1, 3),
			maps.NewOrderedIntSet(1),
			maps.NewOrderedIntSet(1, 3),
		).Build()

		m := v1.ToVectorStoreMap()

		v2 := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
			maps.NewOrderedIntSet(1, 2),
			maps.NewOrderedIntSet(1),
			maps.NewOrderedIntSet(1, 3),
		).Build()

		complement := m.Complement(v2)

		c.Assert(complement.Vectors(), qt.DeepEquals, []sitesmatrix.Vector{
			{3, 1, 1},
			{3, 1, 3},
		})
	})
}

func BenchmarkHasAnyVectorSingle(b *testing.B) {
	testDims := newTestDims()
	set := sitesmatrix.NewIntSetsBuilder(testDims).WithSets(
		maps.NewOrderedIntSet(0),
		maps.NewOrderedIntSet(0),
		maps.NewOrderedIntSet(0),
	).Build()

	v := sitesmatrix.Vector{0, 0, 0}

	for b.Loop() {
		_ = set.HasAnyVector(v)
	}
}
