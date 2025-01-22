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

package maps

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestOrdered(t *testing.T) {
	c := qt.New(t)

	m := NewOrdered[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	c.Assert(m.Keys(), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(m.Values(), qt.DeepEquals, []int{1, 2, 3})

	v, found := m.Get("b")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 2)

	m.Set("b", 22)
	c.Assert(m.Keys(), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(m.Values(), qt.DeepEquals, []int{1, 22, 3})

	m.Delete("b")

	c.Assert(m.Keys(), qt.DeepEquals, []string{"a", "c"})
	c.Assert(m.Values(), qt.DeepEquals, []int{1, 3})
}

func TestOrderedHash(t *testing.T) {
	c := qt.New(t)

	m := NewOrdered[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	h1, err := m.Hash()
	c.Assert(err, qt.IsNil)

	m.Set("d", 4)

	h2, err := m.Hash()
	c.Assert(err, qt.IsNil)

	c.Assert(h1, qt.Not(qt.Equals), h2)

	m = NewOrdered[string, int]()
	m.Set("b", 2)
	m.Set("a", 1)
	m.Set("c", 3)

	h3, err := m.Hash()
	c.Assert(err, qt.IsNil)
	// Order does not matter.
	c.Assert(h1, qt.Equals, h3)
}

func TestOrderedNil(t *testing.T) {
	c := qt.New(t)

	var m *Ordered[string, int]

	m.Set("a", 1)
	c.Assert(m.Keys(), qt.IsNil)
	c.Assert(m.Values(), qt.IsNil)
	v, found := m.Get("a")
	c.Assert(found, qt.Equals, false)
	c.Assert(v, qt.Equals, 0)
	m.Delete("a")
	var b bool
	m.Range(func(k string, v int) bool {
		b = true
		return true
	})
	c.Assert(b, qt.Equals, false)
	c.Assert(m.Len(), qt.Equals, 0)
	c.Assert(m.Clone(), qt.IsNil)
	h, err := m.Hash()
	c.Assert(err, qt.IsNil)
	c.Assert(h, qt.Equals, uint64(0))
}
