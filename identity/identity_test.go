// Copyright 2019 The Hugo Authors. All rights reserved.
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

package identity

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestIdentityManager(t *testing.T) {
	c := qt.New(t)

	id1 := testIdentity{name: "id1"}

	newM := func() Manager {
		m := NewManager(id1)
		m.Add(testIdentity{name: "id2"})

		return m
	}

	c.Run("Search", func(c *qt.C) {
		im := newM()
		c.Assert(im.Search(id1).GetIdentity(), qt.Equals, id1)
		c.Assert(im.Search(testIdentity{name: "notfound"}), qt.Equals, nil)
	})

	c.Run("IsNotDependent", func(c *qt.C) {
		im := newM()
		c.Assert(im.IsNotDependent(testIdentity{name: "notfound"}), qt.Equals, true)
		c.Assert(im.IsNotDependent(testIdentity{name: "id1"}), qt.Equals, false)
		c.Assert(im.IsNotDependent(testIdentity{name: "id2"}), qt.Equals, false)
	})
}

func TestKeyValueIdentity(t *testing.T) {
	c := qt.New(t)

	c.Assert(KeyValueIdentity{Key: "a", Value: "b"}.Eq(KeyValueIdentity{Key: "a", Value: "b"}), qt.Equals, true)
	c.Assert(KeyValueIdentity{Key: "a", Value: "b"}.Eq(KeyValueIdentity{Key: "a", Value: "c"}), qt.Equals, false)
	c.Assert(KeyValueIdentity{Key: "a", Value: "b"}.Eq(KeyValueIdentity{Key: "c", Value: "b"}), qt.Equals, false)

	c.Assert(KeyValueIdentity{Key: "a", Value: "b"}.ProbablyEq(KeyValueIdentity{Key: "a", Value: "b"}), qt.Equals, true)
	c.Assert(KeyValueIdentity{Key: "a", Value: "b"}.ProbablyEq(KeyValueIdentity{Key: "a", Value: "c"}), qt.Equals, false)
	c.Assert(KeyValueIdentity{Key: "a", Value: "b"}.ProbablyEq(KeyValueIdentity{Key: "c", Value: "b"}), qt.Equals, false)
}

func BenchmarkIdentityManager(b *testing.B) {

	createIds := func(num int) []Provider {
		ids := make([]Provider, num)
		for i := 0; i < num; i++ {
			ids[i] = testIdentity{name: fmt.Sprintf("id%d", i)}
		}
		return ids
	}

	b.Run("Add", func(b *testing.B) {
		c := qt.New(b)
		b.StopTimer()
		ids := createIds(b.N)
		im := NewManager(testIdentity{"first"})
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			im.Add(ids[i])
		}

		b.StopTimer()
		c.Assert(im.GetIdentities(), qt.HasLen, b.N+1)
	})

	b.Run("Search", func(b *testing.B) {
		c := qt.New(b)
		b.StopTimer()
		ids := createIds(b.N)
		im := NewManager(testIdentity{"first"})

		for i := 0; i < b.N; i++ {
			im.Add(ids[i])
		}

		b.StartTimer()

		for i := 0; i < b.N; i++ {
			name := "id" + strconv.Itoa(rand.Intn(b.N))
			id := im.Search(testIdentity{name: name})
			c.Assert(id.GetIdentity().Name(), qt.Equals, name)
		}
	})
}

type testIdentity struct {
	name string
}

func (id testIdentity) IsNotDependent(other Provider) bool {
	return id != other.GetIdentity()
}

func (id testIdentity) GetIdentity() Identity {
	return id
}

func (id testIdentity) Base() interface{} {
	return id
}

func (id testIdentity) Name() string {
	return id.name
}

func (id testIdentity) Eq(other interface{}) bool {
	return false
}
func (id testIdentity) ProbablyEq(other interface{}) bool {
	return false
}
