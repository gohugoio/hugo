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
	"iter"

	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
)

var _ TreeThreadSafe[string] = (*TreeShiftTreeSlice[string])(nil)

type TreeShiftTreeSlice[T comparable] struct {
	// v points to a specific tree in the slice.
	v sitesmatrix.Vector

	// The zero value of T.
	zero T

	// trees is a 3D slice that holds all the trees.
	// Note that we have tested a version backed by a map, which is as fast to use, but is twice as epxensive/slow to create.
	trees [][][]*SimpleThreadSafeTree[T]
}

func NewTreeShiftTree[T comparable](v sitesmatrix.Vector) *TreeShiftTreeSlice[T] {
	trees := make([][][]*SimpleThreadSafeTree[T], v[0])
	for i := 0; i < v[0]; i++ {
		trees[i] = make([][]*SimpleThreadSafeTree[T], v[1])
		for j := 0; j < v[1]; j++ {
			trees[i][j] = make([]*SimpleThreadSafeTree[T], v[2])
			for k := 0; k < v[2]; k++ {
				trees[i][j][k] = NewSimpleThreadSafeTree[T]()
			}
		}
	}
	return &TreeShiftTreeSlice[T]{trees: trees}
}

func (t TreeShiftTreeSlice[T]) Shape(v sitesmatrix.Vector) *TreeShiftTreeSlice[T] {
	t.v = v
	return &t
}

func (t *TreeShiftTreeSlice[T]) tree() *SimpleThreadSafeTree[T] {
	return t.trees[t.v[0]][t.v[1]][t.v[2]]
}

func (t *TreeShiftTreeSlice[T]) Get(s string) T {
	return t.tree().Get(s)
}

func (t *TreeShiftTreeSlice[T]) DeleteAllFunc(s string, f func(s string, v T) bool) {
	for tt := range t.Trees() {
		if v := tt.Get(s); v != t.zero {
			if f(s, v) {
				// Delete.
				tt.tree.Delete(s)
			}
		}
	}
}

func (t *TreeShiftTreeSlice[T]) Trees() iter.Seq[*SimpleThreadSafeTree[T]] {
	return func(yield func(v *SimpleThreadSafeTree[T]) bool) {
		for _, l1 := range t.trees {
			for _, l2 := range l1 {
				for _, l3 := range l2 {
					if !yield(l3) {
						return
					}
				}
			}
		}
	}
}

func (t *TreeShiftTreeSlice[T]) LongestPrefix(s string) (string, T) {
	return t.tree().LongestPrefix(s)
}

func (t *TreeShiftTreeSlice[T]) Insert(s string, v T) T {
	return t.tree().Insert(s, v)
}

func (t *TreeShiftTreeSlice[T]) Lock(lockType LockType) func() {
	return t.tree().Lock(lockType)
}

func (t *TreeShiftTreeSlice[T]) WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.tree().WalkPrefix(lockType, s, f)
}

func (t *TreeShiftTreeSlice[T]) WalkPrefixRaw(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	for tt := range t.Trees() {
		if err := tt.WalkPrefix(lockType, s, f); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeShiftTreeSlice[T]) WalkPath(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.tree().WalkPath(lockType, s, f)
}

func (t *TreeShiftTreeSlice[T]) All(lockType LockType) iter.Seq2[string, T] {
	return t.tree().All(lockType)
}

func (t *TreeShiftTreeSlice[T]) LenRaw() int {
	var count int
	for tt := range t.Trees() {
		count += tt.tree.Len()
	}
	return count
}

func (t *TreeShiftTreeSlice[T]) Delete(key string) {
	for tt := range t.Trees() {
		tt.tree.Delete(key)
	}
}

func (t *TreeShiftTreeSlice[T]) DeletePrefix(prefix string) int {
	var count int
	for tt := range t.Trees() {
		count += tt.tree.DeletePrefix(prefix)
	}
	return count
}
