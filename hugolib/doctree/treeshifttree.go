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

type TreeShiftTree[T comparable] struct {
	// This tree is shiftable in one dimension.
	d int

	// The value of the current dimension.
	v int

	// The zero value of T.
	zero T

	// Will be of length equal to the length of the dimension.
	trees []*SimpleTree[T]
}

func NewTreeShiftTree[T comparable](d, length int) *TreeShiftTree[T] {
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

func (t *TreeShiftTree[T]) DeleteAllFunc(s string, f func(s string, v T) bool) {
	for _, tt := range t.trees {
		if v := tt.Get(s); v != t.zero {
			if f(s, v) {
				// Delete.
				tt.tree.Delete(s)
			}
		}
	}
}

func (t *TreeShiftTree[T]) LongestPrefix(s string) (string, T) {
	return t.trees[t.v].LongestPrefix(s)
}

func (t *TreeShiftTree[T]) Insert(s string, v T) T {
	return t.trees[t.v].Insert(s, v)
}

func (t *TreeShiftTree[T]) Lock(lockType LockType) func() {
	return t.trees[t.v].Lock(lockType)
}

func (t *TreeShiftTree[T]) WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.trees[t.v].WalkPrefix(lockType, s, f)
}

func (t *TreeShiftTree[T]) WalkPrefixRaw(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	for _, tt := range t.trees {
		if err := tt.WalkPrefix(lockType, s, f); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeShiftTree[T]) LenRaw() int {
	var count int
	for _, tt := range t.trees {
		count += tt.tree.Len()
	}
	return count
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
