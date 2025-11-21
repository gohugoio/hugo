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
	"sync"

	radix "github.com/gohugoio/go-radix"
)

// Tree is a non thread safe radix tree that holds T.
type Tree[T any] interface {
	TreeCommon[T]
	WalkPrefix(s string, f func(s string, v T) (bool, error)) error
	WalkPath(s string, f func(s string, v T) (bool, error)) error
	All() iter.Seq2[string, T]
}

// TreeThreadSafe is a thread safe radix tree that holds T.
type TreeThreadSafe[T any] interface {
	TreeCommon[T]
	WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error
	WalkPath(lockType LockType, s string, f func(s string, v T) (bool, error)) error
	All(lockType LockType) iter.Seq2[string, T]
}

type TreeCommon[T any] interface {
	Get(s string) T
	LongestPrefix(s string) (string, T)
	Insert(s string, v T) T
}

func NewSimpleTree[T any]() *SimpleTree[T] {
	return &SimpleTree[T]{tree: radix.New[T]()}
}

// SimpleTree is a radix tree that holds T.
// This tree is not thread safe.
type SimpleTree[T any] struct {
	tree *radix.Tree[T]
	zero T
}

func (tree *SimpleTree[T]) Get(s string) T {
	v, _ := tree.tree.Get(s)
	return v
}

func (tree *SimpleTree[T]) LongestPrefix(s string) (string, T) {
	s, v, _ := tree.tree.LongestPrefix(s)
	return s, v
}

func (tree *SimpleTree[T]) Insert(s string, v T) T {
	tree.tree.Insert(s, v)
	return v
}

func (tree *SimpleTree[T]) Walk(f func(s string, v T) (bool, error)) error {
	var walkFn radix.WalkFn[T] = func(s string, v T) (radix.WalkFlag, T, error) {
		var b bool
		b, err := f(s, v)
		if b || err != nil {
			return radix.WalkStop, tree.zero, err
		}
		return radix.WalkContinue, tree.zero, nil
	}
	return tree.tree.Walk(walkFn)
}

func (tree *SimpleTree[T]) WalkPrefix(s string, f func(s string, v T) (bool, error)) error {
	var walkFn radix.WalkFn[T] = func(s string, v T) (radix.WalkFlag, T, error) {
		b, err := f(s, v)
		if b || err != nil {
			return radix.WalkStop, tree.zero, err
		}
		return radix.WalkContinue, tree.zero, nil
	}
	return tree.tree.WalkPrefix(s, walkFn)
}

func (tree *SimpleTree[T]) WalkPath(s string, f func(s string, v T) (bool, error)) error {
	var err error
	var walkFn radix.WalkFn[T] = func(s string, v T) (radix.WalkFlag, T, error) {
		var b bool
		b, err = f(s, v)
		if b || err != nil {
			return radix.WalkStop, tree.zero, err
		}
		return radix.WalkContinue, tree.zero, nil
	}
	tree.tree.WalkPath(s, walkFn)
	return err
}

func (tree *SimpleTree[T]) All() iter.Seq2[string, T] {
	return func(yield func(s string, v T) bool) {
		var walkFn radix.WalkFn[T] = func(s string, v T) (radix.WalkFlag, T, error) {
			if !yield(s, v) {
				return radix.WalkStop, tree.zero, nil
			}
			return radix.WalkContinue, tree.zero, nil
		}
		tree.tree.Walk(walkFn)
	}
}

// NewSimpleThreadSafeTree creates a new SimpleTree.
func NewSimpleThreadSafeTree[T any]() *SimpleThreadSafeTree[T] {
	return &SimpleThreadSafeTree[T]{tree: radix.New[T](), mu: new(sync.RWMutex)}
}

// SimpleThreadSafeTree is a thread safe radix tree that holds T.
type SimpleThreadSafeTree[T any] struct {
	mu   *sync.RWMutex
	tree *radix.Tree[T]
	zero T
}

func (tree *SimpleThreadSafeTree[T]) Get(s string) T {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	v, _ := tree.tree.Get(s)
	return v
}

func (tree *SimpleThreadSafeTree[T]) LongestPrefix(s string) (string, T) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	s, v, _ := tree.tree.LongestPrefix(s)
	return s, v
}

func (tree *SimpleThreadSafeTree[T]) Insert(s string, v T) T {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.tree.Insert(s, v)
	return v
}

func (tree *SimpleThreadSafeTree[T]) Lock(lockType LockType) {
	switch lockType {
	case LockTypeRead:
		tree.mu.RLock()
	case LockTypeWrite:
		tree.mu.Lock()
	}
}

func (tree *SimpleThreadSafeTree[T]) Unlock(lockType LockType) {
	switch lockType {
	case LockTypeRead:
		tree.mu.RUnlock()
	case LockTypeWrite:
		tree.mu.Unlock()
	}
}

func (tree *SimpleThreadSafeTree[T]) WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	tree.Lock(lockType)
	defer tree.Unlock(lockType)
	var walkFn radix.WalkFn[T] = func(s string, v T) (radix.WalkFlag, T, error) {
		var b bool
		b, err := f(s, v)
		if b || err != nil {
			return radix.WalkStop, tree.zero, err
		}
		return radix.WalkContinue, tree.zero, nil
	}
	return tree.tree.WalkPrefix(s, walkFn)
}

func (tree *SimpleThreadSafeTree[T]) WalkPath(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	tree.Lock(lockType)
	defer tree.Unlock(lockType)
	var err error
	var walkFn radix.WalkFn[T] = func(s string, v T) (radix.WalkFlag, T, error) {
		var b bool
		b, err = f(s, v)
		if b || err != nil {
			return radix.WalkStop, tree.zero, err
		}
		return radix.WalkContinue, tree.zero, nil
	}
	tree.tree.WalkPath(s, walkFn)

	return err
}

func (tree *SimpleThreadSafeTree[T]) All(lockType LockType) iter.Seq2[string, T] {
	return func(yield func(s string, v T) bool) {
		tree.Lock(lockType)
		defer tree.Unlock(lockType)
		var walkFn radix.WalkFn[T] = func(s string, v T) (radix.WalkFlag, T, error) {
			if !yield(s, v) {
				return radix.WalkStop, tree.zero, nil
			}
			return radix.WalkContinue, tree.zero, nil
		}
		tree.tree.Walk(walkFn)
	}
}
