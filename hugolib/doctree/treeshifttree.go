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

package doctree

var _ Tree[string] = (*TreeShiftTree[string])(nil)

type TreeShiftTree[T any] struct {
	// This tree is shiftable in one dimension.
	d int

	// The value of the current dimension.
	v int

	// Will be of length equal to the length of the dimension.
	trees []*SimpleTree[T]
}

func NewTreeShiftTree[T any](d, length int) *TreeShiftTree[T] {
	if length <= 0 {
		panic("length must be > 0")
	}
	trees := make([]*SimpleTree[T], length)
	for i := 0; i < length; i++ {
		trees[i] = NewSimpleTree[T]()
	}
	return &TreeShiftTree[T]{d: d, trees: trees}
}

func (t TreeShiftTree[T]) Shape(d, v int) *TreeShiftTree[T] {
	if d != t.d {
		panic("dimension mismatch")
	}
	if v >= len(t.trees) {
		panic("value out of range")
	}
	t.v = v
	return &t
}

func (t *TreeShiftTree[T]) Get(s string) T {
	return t.trees[t.v].Get(s)
}

func (t *TreeShiftTree[T]) LongestPrefix(s string) (string, T) {
	return t.trees[t.v].LongestPrefix(s)
}

func (t *TreeShiftTree[T]) Insert(s string, v T) T {
	return t.trees[t.v].Insert(s, v)
}

func (t *TreeShiftTree[T]) WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.trees[t.v].WalkPrefix(lockType, s, f)
}

func (t *TreeShiftTree[T]) Delete(key string) {
	for _, tt := range t.trees {
		tt.tree.Delete(key)
	}
}

func (t *TreeShiftTree[T]) DeletePrefix(prefix string) int {
	var count int
	for _, tt := range t.trees {
		count += tt.tree.DeletePrefix(prefix)
	}
	return count
}

func (t *TreeShiftTree[T]) Lock(writable bool) (commit func()) {
	if writable {
		for _, tt := range t.trees {
			tt.mu.Lock()
		}
		return func() {
			for _, tt := range t.trees {
				tt.mu.Unlock()
			}
		}
	}

	for _, tt := range t.trees {
		tt.mu.RLock()
	}
	return func() {
		for _, tt := range t.trees {
			tt.mu.RUnlock()
		}
	}
}
