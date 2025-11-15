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

func TestMap(t *testing.T) {
	c := qt.New(t)

	m := NewMap[string, int]()

	m.Set("b", 42)
	v, found := m.Lookup("b")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 42)
	v = m.Get("b")
	c.Assert(v, qt.Equals, 42)
	v, found = m.Lookup("c")
	c.Assert(found, qt.Equals, false)
	c.Assert(v, qt.Equals, 0)
	v = m.Get("c")
	c.Assert(v, qt.Equals, 0)
	v, err := m.GetOrCreate("d", func() (int, error) {
		return 100, nil
	})
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, 100)
	v, found = m.Lookup("d")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 100)

	v, err = m.GetOrCreate("d", func() (int, error) {
		return 200, nil
	})
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, 100)

	wasSet := m.SetIfAbsent("e", 300)
	c.Assert(wasSet, qt.Equals, true)
	v, found = m.Lookup("e")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 300)

	wasSet = m.SetIfAbsent("e", 400)
	c.Assert(wasSet, qt.Equals, false)
	v, found = m.Lookup("e")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 300)

	m.WithWriteLock(func(m map[string]int) {
		m["f"] = 500
	})
	v, found = m.Lookup("f")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 500)
}
