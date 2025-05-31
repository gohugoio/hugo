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

	m := NewOrderedIntSet(1, 2, 3)

	c.Assert(m.Len(), qt.Equals, 3)
	c.Assert(m.Has(1), qt.Equals, true)
	c.Assert(m.Has(4), qt.Equals, false)
	c.Assert(m.String(), qt.Equals, "[1 2 3]")
	m.Set(4)
	c.Assert(m.Len(), qt.Equals, 4)
	c.Assert(m.Has(4), qt.Equals, true)
	c.Assert(m.String(), qt.Equals, "[1 2 3 4]")

	var nilset *OrderedIntSet
	c.Assert(nilset.Len(), qt.Equals, 0)
	c.Assert(nilset.Has(1), qt.Equals, false)
	c.Assert(nilset.String(), qt.Equals, "[]")
}

func BenchmarkOrderedIntSetHasInSmallSet(b *testing.B) {
	m := NewOrderedIntSet()
	for i := range 8 {
		m.Set(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Has(i % 32)
	}
}

func BenchmarkOrderedIntSetNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewOrderedIntSet(1, 2, 3, 4, 5, 6, 7, 8)
	}
}
