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
	// The zero value of T.
	zero T

	// The current dimension.
	d int

	// The value of the current dimension.
	v int

	// Will be of length equal to the length of the dimension.
	trees []*SimpleThreadSafeTree[T]

	dimensions [][]*SimpleThreadSafeTree[T]
}

func NewTreeShiftTree[T comparable](numDimensions int, lengths []int) *TreeShiftTree[T] {
	if numDimensions <= 0 {
		panic("dimensions must be > 0")
	}
	if len(lengths) != numDimensions {
		panic("lengths must match dimensions")
	}
	dimensions := make([][]*SimpleThreadSafeTree[T], numDimensions)

	for d := 0; d < numDimensions; d++ {
		length := lengths[d]
		trees := make([]*SimpleThreadSafeTree[T], length)
		for i := 0; i < length; i++ {
			trees[i] = NewSimpleThreadSafeTree[T]()
		}
		dimensions[d] = trees
	}

	return &TreeShiftTree[T]{dimensions: dimensions}
}

func (t TreeShiftTree[T]) Shape(d, v int) *TreeShiftTree[T] {
	t.d = d
	t.v = v
	return &t
}

func (t *TreeShiftTree[T]) tree() *SimpleThreadSafeTree[T] {
	return t.dimensions[t.d][t.v]
}

func (t *TreeShiftTree[T]) Get(s string) T {
	return t.tree().Get(s)
}

func (t *TreeShiftTree[T]) forEeach(f func(tt *SimpleThreadSafeTree[T]) error) error {
	for _, dd := range t.dimensions {
		for _, tt := range dd {
			if err := f(tt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TreeShiftTree[T]) DeleteAllFunc(s string, f func(s string, v T) bool) {
	_ = t.forEeach(func(tt *SimpleThreadSafeTree[T]) error {
		if v := tt.Get(s); v != t.zero {
			if f(s, v) {
				// Delete.
				tt.tree.Delete(s)
			}
		}
		return nil
	})
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
	return t.forEeach(func(tt *SimpleThreadSafeTree[T]) error {
		return tt.WalkPrefix(lockType, s, f)
	})
}

func (t *TreeShiftTree[T]) WalkPath(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	return t.trees[t.v].WalkPath(lockType, s, f)
}

func (t *TreeShiftTree[T]) All(lockType LockType) iter.Seq2[string, T] {
	return t.trees[t.v].All(lockType)
}

func (t *TreeShiftTree[T]) LenRaw() int {
	var count int
	_ = t.forEeach(func(tt *SimpleThreadSafeTree[T]) error {
		count += tt.tree.Len()
		return nil
	})
	return count
}

func (t *TreeShiftTree[T]) Delete(key string) {
	_ = t.forEeach(func(tt *SimpleThreadSafeTree[T]) error {
		tt.tree.Delete(key)
		return nil
	})
}

func (t *TreeShiftTree[T]) DeletePrefix(prefix string) int {
	var count int
	_ = t.forEeach(func(tt *SimpleThreadSafeTree[T]) error {
		count += tt.tree.DeletePrefix(prefix)
		return nil
	})
	return count
}
