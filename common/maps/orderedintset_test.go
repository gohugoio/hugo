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

package maps

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestOrderedIntSet(t *testing.T) {
	c := qt.New(t)

	m := NewOrderedIntSet(2, 1, 3, 7)

	c.Assert(m.Len(), qt.Equals, 4)
	c.Assert(m.Has(1), qt.Equals, true)
	c.Assert(m.Has(4), qt.Equals, false)
	c.Assert(m.String(), qt.Equals, "[2 1 3 7]")
	m.Set(4)
	c.Assert(m.Len(), qt.Equals, 5)
	c.Assert(m.Has(4), qt.Equals, true)
	c.Assert(m.Next(0), qt.Equals, 1)
	c.Assert(m.Next(1), qt.Equals, 1)
	c.Assert(m.Next(2), qt.Equals, 2)
	c.Assert(m.Next(3), qt.Equals, 3)
	c.Assert(m.Next(4), qt.Equals, 4)
	c.Assert(m.Next(7), qt.Equals, 7)
	c.Assert(m.Next(8), qt.Equals, -1)
	c.Assert(m.String(), qt.Equals, "[2 1 3 7 4]")

	var nilset *OrderedIntSet
	c.Assert(nilset.Len(), qt.Equals, 0)
	c.Assert(nilset.Has(1), qt.Equals, false)
	c.Assert(nilset.String(), qt.Equals, "[]")

	var collected []int
	m.ForEachKey(func(key int) bool {
		collected = append(collected, key)
		return true
	})
	c.Assert(collected, qt.DeepEquals, []int{2, 1, 3, 7, 4})
}

func BenchmarkOrderedIntSet(b *testing.B) {
	smallSet := NewOrderedIntSet()
	for i := range 8 {
		smallSet.Set(i)
	}
	mediumSet := NewOrderedIntSet()
	for i := range 64 {
		mediumSet.Set(i)
	}
	largeSet := NewOrderedIntSet()
	for i := range 1024 {
		largeSet.Set(i)
	}

	b.Run("New", func(b *testing.B) {
		for b.Loop() {
			NewOrderedIntSet(1, 2, 3, 4, 5, 6, 7, 8)
		}
	})

	b.Run("Has small", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			smallSet.Has(i % 32)
		}
	})

	b.Run("Has medium", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			mediumSet.Has(i % 32)
		}
	})

	b.Run("Next", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			mediumSet.Next(i % 32)
		}
	})

	b.Run("ForEachKey small", func(b *testing.B) {
		for b.Loop() {
			smallSet.ForEachKey(func(key int) bool {
				return true
			})
		}
	})

	b.Run("ForEachKey medium", func(b *testing.B) {
		for b.Loop() {
			mediumSet.ForEachKey(func(key int) bool {
				return true
			})
		}
	})

	b.Run("ForEachKey large", func(b *testing.B) {
		for b.Loop() {
			largeSet.ForEachKey(func(key int) bool {
				return true
			})
		}
	})
}
