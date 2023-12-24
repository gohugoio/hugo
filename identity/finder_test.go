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

// Package provides ways to identify values in Hugo. Used for dependency tracking etc.
package identity_test

import (
	"testing"

	"github.com/gohugoio/hugo/identity"
)

func BenchmarkFinder(b *testing.B) {
	m1 := identity.NewManager("")
	m2 := identity.NewManager("")
	m3 := identity.NewManager("")
	m1.AddIdentity(
		testIdentity{"base", "id1", "", "pe1"},
		testIdentity{"base2", "id2", "eq1", ""},
		m2,
		m3,
	)

	b4 := testIdentity{"base4", "id4", "", ""}
	b5 := testIdentity{"base5", "id5", "", ""}

	m2.AddIdentity(b4)

	f := identity.NewFinder(identity.FinderConfig{})

	b.Run("Find one", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := f.Contains(b4, m1, -1)
			if r == 0 {
				b.Fatal("not found")
			}
		}
	})

	b.Run("Find none", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := f.Contains(b5, m1, -1)
			if r > 0 {
				b.Fatal("found")
			}
		}
	})
}
