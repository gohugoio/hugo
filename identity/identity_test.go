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

package identity_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/identity/identitytesting"
)

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
			m := identity.NewManager("")
			if m == nil {
				b.Fatal("manager is nil")
			}
		}
	})

	b.Run("Add unique", func(b *testing.B) {
		ids := createIds(b.N)
		im := identity.NewManager("")

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			im.AddIdentity(ids[i])
		}

		b.StopTimer()
	})

	b.Run("Add duplicates", func(b *testing.B) {
		id := &testIdentity{base: "a", name: "b"}
		im := identity.NewManager("")

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
}

func BenchmarkIsNotDependent(b *testing.B) {
	runBench := func(b *testing.B, id1, id2 identity.Identity) {
		for i := 0; i < b.N; i++ {
			isNotDependent(id1, id2)
		}
	}

	newNestedManager := func(depth, count int) identity.Manager {
		m1 := identity.NewManager("")
		for i := 0; i < depth; i++ {
			m2 := identity.NewManager("")
			m1.AddIdentity(m2)
			for j := 0; j < count; j++ {
				id := fmt.Sprintf("id%d", j)
				m2.AddIdentity(&testIdentity{id, id, "", ""})
			}
			m1 = m2
		}
		return m1
	}

	type depthCount struct {
		depth int
		count int
	}

	for _, dc := range []depthCount{{10, 5}} {
		b.Run(fmt.Sprintf("Nested not found %d %d", dc.depth, dc.count), func(b *testing.B) {
			im := newNestedManager(dc.depth, dc.count)
			id1 := identity.StringIdentity("idnotfound")
			b.ResetTimer()
			runBench(b, im, id1)
		})
	}
}

func TestIdentityManager(t *testing.T) {
	c := qt.New(t)

	newNestedManager := func() identity.Manager {
		m1 := identity.NewManager("")
		m2 := identity.NewManager("")
		m3 := identity.NewManager("")
		m1.AddIdentity(
			testIdentity{"base", "id1", "", "pe1"},
			testIdentity{"base2", "id2", "eq1", ""},
			m2,
			m3,
		)

		m2.AddIdentity(testIdentity{"base4", "id4", "", ""})

		return m1
	}

	c.Run("Anonymous", func(c *qt.C) {
		im := newNestedManager()
		c.Assert(im.GetIdentity(), qt.Equals, identity.Anonymous)
		im.AddIdentity(identity.Anonymous)
		c.Assert(isNotDependent(identity.Anonymous, identity.Anonymous), qt.IsTrue)
	})

	c.Run("GenghisKhan", func(c *qt.C) {
		c.Assert(isNotDependent(identity.GenghisKhan, identity.GenghisKhan), qt.IsTrue)
	})
}

type testIdentity struct {
	base string
	name string

	idEq         string
	idProbablyEq string
}

func (id testIdentity) Eq(other any) bool {
	ot, ok := other.(testIdentity)
	if !ok {
		return false
	}
	if ot.idEq == "" || id.idEq == "" {
		return false
	}
	return ot.idEq == id.idEq
}

func (id testIdentity) IdentifierBase() string {
	return id.base
}

func (id testIdentity) Name() string {
	return id.name
}

func (id testIdentity) ProbablyEq(other any) bool {
	ot, ok := other.(testIdentity)
	if !ok {
		return false
	}
	if ot.idProbablyEq == "" || id.idProbablyEq == "" {
		return false
	}
	return ot.idProbablyEq == id.idProbablyEq
}

func isNotDependent(a, b identity.Identity) bool {
	f := identity.NewFinder(identity.FinderConfig{})
	r := f.Contains(b, a, -1)
	return r == 0
}
