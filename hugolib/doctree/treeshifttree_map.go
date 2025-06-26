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

import (
	"iter"

	"github.com/gohugoio/hugo/hugolib/sitematrix"
)

var _ TreeThreadSafe[string] = (*TreeShiftTreeMap[string])(nil)

type TreeShiftTreeMap[T comparable] struct {
	// v is the index of the current dimension.
	v sitematrix.Vector

	// The zero value of T.
	zero T

	dimensions map[sitematrix.Vector]*SimpleThreadSafeTree[T]
}

func NewTreeShiftTreeMap[T comparable](v sitematrix.Vector) *TreeShiftTreeMap[T] {
	dimensions := make(map[sitematrix.Vector]*SimpleThreadSafeTree[T])
	for i := 0; i < v[0]; i++ {
		for j := 0; j < v[1]; j++ {
			for k := 0; k < v[2]; k++ {
				vec := sitematrix.Vector{i, j, k}
				dimensions[vec] = NewSimpleThreadSafeTree[T]()
			}
		}
	}
	return &TreeShiftTreeMap[T]{dimensions: dimensions}
}

func (t TreeShiftTreeMap[T]) Shape(v sitematrix.Vector) *TreeShiftTreeMap[T] {
	t.v = v
	return &t
}

func (t *TreeShiftTreeMap[T]) tree() *SimpleThreadSafeTree[T] {
	return t.dimensions[t.v]
}

func (t *TreeShiftTreeMap[T]) Get(s string) T {
	return t.tree().Get(s)
}

func (t *TreeShiftTreeMap[T]) DeleteAllFunc(s string, f func(s string, v T) bool) {
	for _, tt := range t.dimensions {
		if v := tt.Get(s); v != t.zero {
			if f(s, v) {
				// Delete.
				tt.tree.Delete(s)
			}
		}
	}
}

func (t *TreeShiftTreeMap[T]) Trees() iter.Seq[*SimpleThreadSafeTree[T]] {
	return func(yield func(v *SimpleThreadSafeTree[T]) bool) {
		for _, tree := range t.dimensions {
			if !yield(tree) {
				return
			}
		}
	}
}

func (t *TreeShiftTreeMap[T]) LongestPrefix(s string) (string, T) {
	return t.tree().LongestPrefix(s)
}

func (t *TreeShiftTreeMap[T]) Insert(s string, v T) T {
	return t.tree().Insert(s, v)
}

func (t *TreeShiftTreeMap[T]) Lock(lockType LockType) func() {
	return t.tree().Lock(lockType)
}

func (t *TreeShiftTreeMap[T]) WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.tree().WalkPrefix(lockType, s, f)
}

func (t *TreeShiftTreeMap[T]) WalkPrefixRaw(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	for tt := range t.Trees() {
		if err := tt.WalkPrefix(lockType, s, f); err != nil {
			return err
		}
	}
	return nil
}

func (t *TreeShiftTreeMap[T]) WalkPath(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.tree().WalkPath(lockType, s, f)
}

func (t *TreeShiftTreeMap[T]) All(lockType LockType) iter.Seq2[string, T] {
	return t.tree().All(lockType)
}

func (t *TreeShiftTreeMap[T]) LenRaw() int {
	var count int
	for tt := range t.Trees() {
		count += tt.tree.Len()
	}
	return count
}

func (t *TreeShiftTreeMap[T]) Delete(key string) {
	for tt := range t.Trees() {
		tt.tree.Delete(key)
	}
}

func (t *TreeShiftTreeMap[T]) DeletePrefix(prefix string) int {
	var count int
	for tt := range t.Trees() {
		count += tt.tree.DeletePrefix(prefix)
	}
	return count
}
