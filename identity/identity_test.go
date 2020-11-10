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

package identity_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/identity/identitytesting"
)

func TestIdentityManager(t *testing.T) {
	c := qt.New(t)

	newM := func() identity.Manager {
		m1 := identity.NewManager(testIdentity{"base", "root"})
		m2 := identity.NewManager(identity.Anonymous)
		m3 := identity.NewManager(testIdentity{"base3", "id3"})
		m1.AddIdentity(
			testIdentity{"base", "id1"},
			testIdentity{"base2", "id2"},
			m2,
			m3,
		)

		m2.AddIdentity(testIdentity{"base4", "id4"})

		return m1
	}

	c.Run("Contains", func(c *qt.C) {
		im := newM()
		c.Assert(im.Contains(testIdentity{"base", "root"}), qt.IsTrue)
		c.Assert(im.Contains(testIdentity{"base", "id1"}), qt.IsTrue)
		c.Assert(im.Contains(testIdentity{"base3", "id3"}), qt.IsTrue)
		c.Assert(im.Contains(testIdentity{"base", "notfound"}), qt.IsFalse)

		im.Reset()
		c.Assert(im.Contains(testIdentity{"base", "root"}), qt.IsTrue)
		c.Assert(im.Contains(testIdentity{"base", "id1"}), qt.IsFalse)
	})

	c.Run("ContainsProbably", func(c *qt.C) {
		im := newM()
		c.Assert(im.ContainsProbably(testIdentity{"base", "id1"}), qt.IsTrue)
		c.Assert(im.ContainsProbably(testIdentity{"base", "notfound"}), qt.IsTrue)
		c.Assert(im.ContainsProbably(testIdentity{"base2", "notfound"}), qt.IsTrue)
		c.Assert(im.ContainsProbably(testIdentity{"base3", "notfound"}), qt.IsTrue)
		c.Assert(im.ContainsProbably(testIdentity{"base4", "notfound"}), qt.IsTrue)
		c.Assert(im.ContainsProbably(testIdentity{"base5", "notfound"}), qt.IsFalse)

		im.Reset()
		c.Assert(im.Contains(testIdentity{"base", "root"}), qt.IsTrue)
		c.Assert(im.ContainsProbably(testIdentity{"base", "notfound"}), qt.IsTrue)
	})

	c.Run("Anonymous", func(c *qt.C) {
		im := newM()
		im.AddIdentity(identity.Anonymous)
		c.Assert(im.Contains(identity.Anonymous), qt.IsFalse)
		c.Assert(im.ContainsProbably(identity.Anonymous), qt.IsFalse)
		c.Assert(identity.IsNotDependent(identity.Anonymous, identity.Anonymous), qt.IsTrue)
	})

	c.Run("GenghisKhan", func(c *qt.C) {
		im := newM()
		c.Assert(im.Contains(identity.GenghisKhan), qt.IsFalse)
		c.Assert(im.ContainsProbably(identity.GenghisKhan), qt.IsTrue)
		c.Assert(identity.IsNotDependent(identity.GenghisKhan, identity.GenghisKhan), qt.IsTrue)
	})
}

func BenchmarkIdentityManager(b *testing.B) {
	createIds := func(num int) []identity.Identity {
		ids := make([]identity.Identity, num)
		for i := 0; i < num; i++ {
			name := fmt.Sprintf("id%d", i)
			ids[i] = &testIdentity{base: name, name: name}
		}
		return ids
	}

	b.Run("identity.NewManager", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := identity.NewManager(identity.Anonymous)
			if m == nil {
				b.Fatal("manager is nil")
			}
		}
	})

	b.Run("Add unique", func(b *testing.B) {
		ids := createIds(b.N)
		im := identity.NewManager(identity.Anonymous)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			im.AddIdentity(ids[i])
		}

		b.StopTimer()
	})

	b.Run("Add duplicates", func(b *testing.B) {
		id := &testIdentity{base: "a", name: "b"}
		im := identity.NewManager(identity.Anonymous)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			im.AddIdentity(id)
		}

		b.StopTimer()
	})

	b.Run("Nop StringIdentity const", func(b *testing.B) {
		const id = identity.StringIdentity("test")
		for i := 0; i < b.N; i++ {
			identity.NopManager.AddIdentity(id)
		}
	})

	b.Run("Nop StringIdentity const other package", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			identity.NopManager.AddIdentity(identitytesting.TestIdentity)
		}
	})

	b.Run("Nop StringIdentity var", func(b *testing.B) {
		id := identity.StringIdentity("test")
		for i := 0; i < b.N; i++ {
			identity.NopManager.AddIdentity(id)
		}
	})

	b.Run("Nop pointer identity", func(b *testing.B) {
		id := &testIdentity{base: "a", name: "b"}
		for i := 0; i < b.N; i++ {
			identity.NopManager.AddIdentity(id)
		}
	})

	b.Run("Nop Anonymous", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			identity.NopManager.AddIdentity(identity.Anonymous)
		}
	})

	runContainsBenchmark := func(b *testing.B, im identity.Manager, fn func(id identity.Identity) bool, shouldFind bool) {
		if shouldFind {
			ids := createIds(b.N)

			for i := 0; i < b.N; i++ {
				im.AddIdentity(ids[i])
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				found := fn(ids[i])
				if !found {
					b.Fatal("id not found")
				}
			}
		} else {
			noMatchQuery := &testIdentity{base: "notfound", name: "notfound"}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				found := fn(noMatchQuery)
				if found {
					b.Fatal("id  found")
				}
			}
		}
	}

	b.Run("Contains", func(b *testing.B) {
		im := identity.NewManager(identity.Anonymous)
		runContainsBenchmark(b, im, im.Contains, true)
	})

	b.Run("ContainsNotFound", func(b *testing.B) {
		im := identity.NewManager(identity.Anonymous)
		runContainsBenchmark(b, im, im.Contains, false)
	})

	b.Run("ContainsProbably", func(b *testing.B) {
		im := identity.NewManager(identity.Anonymous)
		runContainsBenchmark(b, im, im.ContainsProbably, true)
	})

	b.Run("ContainsProbablyNotFound", func(b *testing.B) {
		im := identity.NewManager(identity.Anonymous)
		runContainsBenchmark(b, im, im.ContainsProbably, false)
	})
}

type testIdentity struct {
	base string
	name string
}

func (id testIdentity) IdentifierBase() interface{} {
	return id.base
}

func (id testIdentity) Name() string {
	return id.name
}

type testIdentityManager struct {
	testIdentity
	identity.Manager
}
