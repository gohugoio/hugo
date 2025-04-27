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

import "iter"

var _ TreeThreadSafe[string] = (*TreeShiftTree[string])(nil)

type TreeShiftTree[T comparable] struct {
	// This tree is shiftable in one dimension.
	d int

	// The value of the current dimension.
	v int

	// The zero value of T.
	zero T

	dimensions [][]*SimpleThreadSafeTree[T]
}

func NewTreeShiftTree[T comparable](dimensionsLengths []int) *TreeShiftTree[T] {
	dimensions := make([][]*SimpleThreadSafeTree[T], len(dimensionsLengths))
	for d := 0; d < len(dimensionsLengths); d++ {
		length := dimensionsLengths[d]
		trees := make([]*SimpleThreadSafeTree[T], length)
		for i := 0; i < length; i++ {
			trees[i] = NewSimpleThreadSafeTree[T]()
		}
		dimensions[d] = trees
	}

	return &TreeShiftTree[T]{dimensions: dimensions}
}

func (t TreeShiftTree[T]) Shape(d, v int) *TreeShiftTree[T] {
	if v < 0 || v >= len(t.dimensions[d]) {
		panic("dimension value out of range")
	}
	t.v = v
	t.d = d
	return &t
}

func (t *TreeShiftTree[T]) tree() *SimpleThreadSafeTree[T] {
	return t.dimensions[t.d][t.v]
}

func (t *TreeShiftTree[T]) Get(s string) T {
	return t.tree().Get(s)
}

func (t *TreeShiftTree[T]) DeleteAllFunc(s string, f func(s string, v T) bool) {
	for tt := range t.Trees() {
		if v := tt.Get(s); v != t.zero {
			if f(s, v) {
				// Delete.
				tt.tree.Delete(s)
			}
		}
	}
}

func (t *TreeShiftTree[T]) Trees() iter.Seq[*SimpleThreadSafeTree[T]] {
	return func(yield func(v *SimpleThreadSafeTree[T]) bool) {
		for _, dd := range t.dimensions {
			for _, tt := range dd {
				if !yield(tt) {
					return
				}
			}
		}
	}
}

func (t *TreeShiftTree[T]) LongestPrefix(s string) (string, T) {
	return t.tree().LongestPrefix(s)
}

func (t *TreeShiftTree[T]) Insert(s string, v T) T {
	return t.tree().Insert(s, v)
}

func (t *TreeShiftTree[T]) Lock(lockType LockType) func() {
	return t.tree().Lock(lockType)
}

func (t *TreeShiftTree[T]) WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.tree().WalkPrefix(lockType, s, f)
}

func (t *TreeShiftTree[T]) WalkPrefixRaw(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	for tt := range t.Trees() {
		if err := tt.WalkPrefix(lockType, s, f); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeShiftTree[T]) WalkPath(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.tree().WalkPath(lockType, s, f)
}

func (t *TreeShiftTree[T]) All(lockType LockType) iter.Seq2[string, T] {
	return t.tree().All(lockType)
}

func (t *TreeShiftTree[T]) LenRaw() int {
	var count int
	for tt := range t.Trees() {
		count += tt.tree.Len()
	}
	return count
}

func (t *TreeShiftTree[T]) Delete(key string) {
	for tt := range t.Trees() {
		tt.tree.Delete(key)
	}
}

func (t *TreeShiftTree[T]) DeletePrefix(prefix string) int {
	var count int
	for tt := range t.Trees() {
		count += tt.tree.DeletePrefix(prefix)
	}
	return count
}
