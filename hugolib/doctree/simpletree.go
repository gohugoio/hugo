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

	radix "github.com/armon/go-radix"
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
	return &SimpleTree[T]{tree: radix.New()}
}

// SimpleTree is a radix tree that holds T.
// This tree is not thread safe.
type SimpleTree[T any] struct {
	tree *radix.Tree
	zero T
}

func (tree *SimpleTree[T]) Get(s string) T {
	if v, ok := tree.tree.Get(s); ok {
		return v.(T)
	}
	return tree.zero
}

func (tree *SimpleTree[T]) LongestPrefix(s string) (string, T) {
	if s, v, ok := tree.tree.LongestPrefix(s); ok {
		return s, v.(T)
	}
	return "", tree.zero
}

func (tree *SimpleTree[T]) Insert(s string, v T) T {
	tree.tree.Insert(s, v)
	return v
}

func (tree *SimpleTree[T]) Walk(f func(s string, v T) (bool, error)) error {
	var err error
	tree.tree.Walk(func(s string, v any) bool {
		var b bool
		b, err = f(s, v.(T))
		if err != nil {
			return true
		}
		return b
	})
	return err
}

func (tree *SimpleTree[T]) WalkPrefix(s string, f func(s string, v T) (bool, error)) error {
	var err error
	tree.tree.WalkPrefix(s, func(s string, v any) bool {
		var b bool
		b, err = f(s, v.(T))
		if err != nil {
			return true
		}
		return b
	})

	return err
}

func (tree *SimpleTree[T]) WalkPath(s string, f func(s string, v T) (bool, error)) error {
	var err error
	tree.tree.WalkPath(s, func(s string, v any) bool {
		var b bool
		b, err = f(s, v.(T))
		if err != nil {
			return true
		}
		return b
	})
	return err
}

func (tree *SimpleTree[T]) All() iter.Seq2[string, T] {
	return func(yield func(s string, v T) bool) {
		tree.tree.Walk(func(s string, v any) bool {
			return !yield(s, v.(T))
		})
	}
}

// NewSimpleThreadSafeTree creates a new SimpleTree.
func NewSimpleThreadSafeTree[T any]() *SimpleThreadSafeTree[T] {
	return &SimpleThreadSafeTree[T]{tree: radix.New(), mu: new(sync.RWMutex)}
}

// SimpleThreadSafeTree is a thread safe radix tree that holds T.
type SimpleThreadSafeTree[T any] struct {
	mu     *sync.RWMutex
	noLock bool
	tree   *radix.Tree
	zero   T
}

var noopFunc = func() {}

func (tree *SimpleThreadSafeTree[T]) readLock() func() {
	if tree.noLock {
		return noopFunc
	}
	tree.mu.RLock()
	return tree.mu.RUnlock
}

func (tree *SimpleThreadSafeTree[T]) writeLock() func() {
	if tree.noLock {
		return noopFunc
	}
	tree.mu.Lock()
	return tree.mu.Unlock
}

func (tree *SimpleThreadSafeTree[T]) Get(s string) T {
	unlock := tree.readLock()
	defer unlock()

	if v, ok := tree.tree.Get(s); ok {
		return v.(T)
	}
	return tree.zero
}

func (tree *SimpleThreadSafeTree[T]) LongestPrefix(s string) (string, T) {
	unlock := tree.readLock()
	defer unlock()

	if s, v, ok := tree.tree.LongestPrefix(s); ok {
		return s, v.(T)
	}
	return "", tree.zero
}

func (tree *SimpleThreadSafeTree[T]) Insert(s string, v T) T {
	unlock := tree.writeLock()
	defer unlock()

	tree.tree.Insert(s, v)
	return v
}

func (tree *SimpleThreadSafeTree[T]) Lock(lockType LockType) func() {
	switch lockType {
	case LockTypeNone:
		return noopFunc
	case LockTypeRead:
		tree.mu.RLock()
		return tree.mu.RUnlock
	case LockTypeWrite:
		tree.mu.Lock()
		return tree.mu.Unlock
	}
	return noopFunc
}

func (tree SimpleThreadSafeTree[T]) LockTree(lockType LockType) (TreeThreadSafe[T], func()) {
	unlock := tree.Lock(lockType)
	tree.noLock = true
	return &tree, unlock // create a copy of tree with the noLock flag set to true.
}

func (tree *SimpleThreadSafeTree[T]) WalkPrefix(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	commit := tree.Lock(lockType)
	defer commit()
	var err error
	tree.tree.WalkPrefix(s, func(s string, v any) bool {
		var b bool
		b, err = f(s, v.(T))
		if err != nil {
			return true
		}
		return b
	})

	return err
}

func (tree *SimpleThreadSafeTree[T]) WalkPath(lockType LockType, s string, f func(s string, v T) (bool, error)) error {
	commit := tree.Lock(lockType)
	defer commit()
	var err error
	tree.tree.WalkPath(s, func(s string, v any) bool {
		var b bool
		b, err = f(s, v.(T))
		if err != nil {
			return true
		}
		return b
	})

	return err
}

func (tree *SimpleThreadSafeTree[T]) All(lockType LockType) iter.Seq2[string, T] {
	commit := tree.Lock(lockType)
	defer commit()
	return func(yield func(s string, v T) bool) {
		tree.tree.Walk(func(s string, v any) bool {
			return !yield(s, v.(T))
		})
	}
}

// iter.Seq[*TemplWithBaseApplied]
