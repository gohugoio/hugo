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

package doctree

import (
	"fmt"
	"testing"
)

func BenchmarkSimpleThreadSafeTree(b *testing.B) {
	newTestTree := func() TreeThreadSafe[int] {
		t := NewSimpleThreadSafeTree[int]()
		for i := 0; i < 1000; i++ {
			t.Insert(fmt.Sprintf("key%d", i), i)
		}
		return t
	}

	b.Run("Get", func(b *testing.B) {
		t := newTestTree()
		b.ResetTimer()
		for b.Loop() {
			t.Get("key500")
		}
	})

	b.Run("Insert", func(b *testing.B) {
		t := newTestTree()
		b.ResetTimer()
		for b.Loop() {
			t.Insert("key500", 501)
		}
	})

	b.Run("LongestPrefix", func(b *testing.B) {
		t := newTestTree()
		b.ResetTimer()
		for b.Loop() {
			t.LongestPrefix("key500extra")
		}
	})
}
