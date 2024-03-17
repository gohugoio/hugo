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
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	radix "github.com/armon/go-radix"
	"github.com/gohugoio/hugo/resources/resource"
)

type (
	Config[T any] struct {
		// Shifter handles tree transformations.
		Shifter Shifter[T]
	}

	// Shifter handles tree transformations.
	Shifter[T any] interface {
		// ForEeachInDimension will call the given function for each value in the given dimension.
		// If the function returns true, the walk will stop.
		ForEeachInDimension(n T, d int, f func(T) bool)

		// Insert inserts new into the tree into the dimension it provides.
		// It may replace old.
		// It returns the updated and existing T
		// and a bool indicating if an existing record is updated.
		Insert(old, new T) (T, T, bool)

		// Insert inserts new into the given dimension.
		// It may replace old.
		// It returns the updated and existing T
		// and a bool indicating if an existing record is updated.
		InsertInto(old, new T, dimension Dimension) (T, T, bool)

		// Delete deletes T from the given dimension and returns the deleted T and whether the dimension was deleted and if  it's empty after the delete.
		Delete(v T, dimension Dimension) (T, bool, bool)

		// Shift shifts T into the given dimension
		// and returns the shifted T and a bool indicating if the shift was successful and
		// how accurate a match T is according to its dimensions.
		Shift(v T, dimension Dimension, exact bool) (T, bool, DimensionFlag)
	}
)

// NodeShiftTree is the root of a tree that can be shaped using the Shape method.
// Note that multipled shapes of the same tree is meant to be used concurrently,
// so use the applicable locking when needed.
type NodeShiftTree[T any] struct {
	tree *radix.Tree

	// E.g. [language, role].
	dims    Dimension
	shifter Shifter[T]

	mu *sync.RWMutex
}

func New[T any](cfg Config[T]) *NodeShiftTree[T] {
	if cfg.Shifter == nil {
		panic("Shifter is required")
	}

	return &NodeShiftTree[T]{
		mu:      &sync.RWMutex{},
		shifter: cfg.Shifter,
		tree:    radix.New(),
	}
}

func (r *NodeShiftTree[T]) Delete(key string) (T, bool) {
	return r.delete(key)
}

func (r *NodeShiftTree[T]) DeleteRaw(key string) {
	r.delete(key)
}

func (r *NodeShiftTree[T]) DeleteAll(key string) {
	r.tree.WalkPrefix(key, func(key string, value any) bool {
		v, ok := r.tree.Delete(key)
		if ok {
			resource.MarkStale(v)
		}
		return false
	})
}

func (r *NodeShiftTree[T]) DeletePrefix(prefix string) int {
	count := 0
	var keys []string
	r.tree.WalkPrefix(prefix, func(key string, value any) bool {
		keys = append(keys, key)
		return false
	})
	for _, key := range keys {
		if _, ok := r.delete(key); ok {
			count++
		}
	}
	return count
}

func (r *NodeShiftTree[T]) delete(key string) (T, bool) {
	var wasDeleted bool
	var deleted T
	if v, ok := r.tree.Get(key); ok {
		var isEmpty bool
		deleted, wasDeleted, isEmpty = r.shifter.Delete(v.(T), r.dims)
		if isEmpty {
			r.tree.Delete(key)
		}
	}
	return deleted, wasDeleted
}

func (t *NodeShiftTree[T]) DeletePrefixAll(prefix string) int {
	count := 0

	t.tree.WalkPrefix(prefix, func(key string, value any) bool {
		if v, ok := t.tree.Delete(key); ok {
			resource.MarkStale(v)
			count++
		}
		return false
	})

	return count
}

// Increment the value of dimension d by 1.
func (t *NodeShiftTree[T]) Increment(d int) *NodeShiftTree[T] {
	return t.Shape(d, t.dims[d]+1)
}

func (r *NodeShiftTree[T]) InsertIntoCurrentDimension(s string, v T) (T, T, bool) {
	s = mustValidateKey(cleanKey(s))
	var (
		updated  bool
		existing T
	)
	if vv, ok := r.tree.Get(s); ok {
		v, existing, updated = r.shifter.InsertInto(vv.(T), v, r.dims)
	}
	r.tree.Insert(s, v)
	return v, existing, updated
}

// InsertIntoValuesDimension inserts v into the tree at the given key and the
// dimension defined by the value.
// It returns the updated and existing T and a bool indicating if an existing record is updated.
func (r *NodeShiftTree[T]) InsertIntoValuesDimension(s string, v T) (T, T, bool) {
	s = mustValidateKey(cleanKey(s))
	var (
		updated  bool
		existing T
	)
	if vv, ok := r.tree.Get(s); ok {
		v, existing, updated = r.shifter.Insert(vv.(T), v)
	}
	r.tree.Insert(s, v)
	return v, existing, updated
}

func (r *NodeShiftTree[T]) InsertRawWithLock(s string, v any) (any, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tree.Insert(s, v)
}

// It returns the updated and existing T and a bool indicating if an existing record is updated.
func (r *NodeShiftTree[T]) InsertIntoValuesDimensionWithLock(s string, v T) (T, T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.InsertIntoValuesDimension(s, v)
}

func (t *NodeShiftTree[T]) Len() int {
	return t.tree.Len()
}

func (t *NodeShiftTree[T]) CanLock() bool {
	ok := t.mu.TryLock()
	if ok {
		t.mu.Unlock()
	}
	return ok
}

// Lock locks the data store for read or read/write access until commit is invoked.
// Note that Root is not thread-safe outside of this transaction construct.
func (t *NodeShiftTree[T]) Lock(writable bool) (commit func()) {
	if writable {
		t.mu.Lock()
	} else {
		t.mu.RLock()
	}
	return func() {
		if writable {
			t.mu.Unlock()
		} else {
			t.mu.RUnlock()
		}
	}
}

// LongestPrefix finds the longest prefix of s that exists in the tree that also matches the predicate (if set).
// Set exact to true to only match exact in the current dimension (e.g. language).
func (r *NodeShiftTree[T]) LongestPrefix(s string, exact bool, predicate func(v T) bool) (string, T) {
	for {
		longestPrefix, v, found := r.tree.LongestPrefix(s)

		if found {
			if t, ok, _ := r.shift(v.(T), exact); ok && (predicate == nil || predicate(t)) {
				return longestPrefix, t
			}
		}

		if s == "" || s == "/" {
			var t T
			return "", t
		}

		// Walk up to find a node in the correct dimension.
		s = path.Dir(s)

	}
}

// LongestPrefixAll returns the longest prefix considering all tree dimensions.
func (r *NodeShiftTree[T]) LongestPrefixAll(s string) (string, bool) {
	s, _, found := r.tree.LongestPrefix(s)
	return s, found
}

func (r *NodeShiftTree[T]) GetRaw(s string) (T, bool) {
	v, ok := r.tree.Get(s)
	if !ok {
		var t T
		return t, false
	}
	return v.(T), true
}

func (r *NodeShiftTree[T]) WalkPrefixRaw(prefix string, walker func(key string, value T) bool) {
	walker2 := func(key string, value any) bool {
		return walker(key, value.(T))
	}
	r.tree.WalkPrefix(prefix, walker2)
}

// Shape the tree for dimension d to value v.
func (t *NodeShiftTree[T]) Shape(d, v int) *NodeShiftTree[T] {
	x := t.clone()
	x.dims[d] = v
	return x
}

func (t *NodeShiftTree[T]) String() string {
	return fmt.Sprintf("Root{%v}", t.dims)
}

func (r *NodeShiftTree[T]) Get(s string) T {
	t, _ := r.get(s)
	return t
}

func (r *NodeShiftTree[T]) ForEeachInDimension(s string, d int, f func(T) bool) {
	s = cleanKey(s)
	v, ok := r.tree.Get(s)
	if !ok {
		return
	}
	r.shifter.ForEeachInDimension(v.(T), d, f)
}

type WalkFunc[T any] func(string, T) (bool, error)

type NodeShiftTreeWalker[T any] struct {
	// The tree to walk.
	Tree *NodeShiftTree[T]

	// Handle will be called for each node in the main tree.
	// If the callback returns true, the walk will stop.
	// The callback can optionally return a callback for the nested tree.
	Handle func(s string, v T, exact DimensionFlag) (terminate bool, err error)

	// Optional prefix filter.
	Prefix string

	// Enable read or write locking if needed.
	LockType LockType

	// When set, no dimension shifting will be performed.
	NoShift bool

	// Don't fall back to alternative dimensions (e.g. language).
	Exact bool

	// Used in development only.
	Debug bool

	// Optional context.
	// Note that this is copied to the nested walkers using Extend.
	// This means that walkers can pass data (down) and events (up) to
	// the related walkers.
	WalkContext *WalkContext[T]

	// Local state.
	// This is scoped to the current walker and not copied to the nested walkers.
	skipPrefixes []string
}

// Extend returns a new NodeShiftTreeWalker with the same configuration as the
// and the same WalkContext as the original.
// Any local state is reset.
func (r NodeShiftTreeWalker[T]) Extend() *NodeShiftTreeWalker[T] {
	r.resetLocalState()
	return &r
}

// SkipPrefix adds a prefix to be skipped in the walk.
func (r *NodeShiftTreeWalker[T]) SkipPrefix(prefix ...string) {
	r.skipPrefixes = append(r.skipPrefixes, prefix...)
}

// ShouldSkip returns whether the given key should be skipped in the walk.
func (r *NodeShiftTreeWalker[T]) ShouldSkip(s string) bool {
	for _, prefix := range r.skipPrefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func (r *NodeShiftTreeWalker[T]) Walk(ctx context.Context) error {
	if r.Tree == nil {
		panic("Tree is required")
	}
	r.resetLocalState()

	if r.LockType > LockTypeNone {
		commit1 := r.Tree.Lock(r.LockType == LockTypeWrite)
		defer commit1()
	}

	main := r.Tree

	var err error
	fnMain := func(s string, v interface{}) bool {
		if r.ShouldSkip(s) {
			return false
		}

		t, ok, exact := r.toT(r.Tree, v)
		if !ok {
			return false
		}

		var terminate bool
		terminate, err = r.Handle(s, t, exact)
		if terminate || err != nil {
			return true
		}
		return false
	}

	if r.Prefix != "" {
		main.tree.WalkPrefix(r.Prefix, fnMain)
	} else {
		main.tree.Walk(fnMain)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *NodeShiftTreeWalker[T]) resetLocalState() {
	r.skipPrefixes = nil
}

func (r *NodeShiftTreeWalker[T]) toT(tree *NodeShiftTree[T], v any) (t T, ok bool, exact DimensionFlag) {
	if r.NoShift {
		t = v.(T)
		ok = true
	} else {
		t, ok, exact = tree.shift(v.(T), r.Exact)
	}
	return
}

func (r *NodeShiftTree[T]) Has(s string) bool {
	_, ok := r.get(s)
	return ok
}

func (t NodeShiftTree[T]) clone() *NodeShiftTree[T] {
	return &t
}

func (r *NodeShiftTree[T]) shift(t T, exact bool) (T, bool, DimensionFlag) {
	return r.shifter.Shift(t, r.dims, exact)
}

func (r *NodeShiftTree[T]) get(s string) (T, bool) {
	s = cleanKey(s)
	v, ok := r.tree.Get(s)
	if !ok {
		var t T
		return t, false
	}
	t, ok, _ := r.shift(v.(T), true)
	return t, ok
}

type WalkConfig[T any] struct {
	// Optional prefix filter.
	Prefix string

	// Callback will be called for each node in the tree.
	// If the callback returns true, the walk will stop.
	Callback func(ctx *WalkContext[T], s string, t T) (bool, error)

	// Enable read or write locking if needed.
	LockType LockType

	// When set, no dimension shifting will be performed.
	NoShift bool

	// Exact will only match exact in the current dimension (e.g. language),
	// and will not look for alternatives.
	Exact bool
}
