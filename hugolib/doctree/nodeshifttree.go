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
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	radix "github.com/armon/go-radix"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/resources/resource"
)

type (
	Config[T any] struct {
		// Shifter handles tree transformations.
		Shifter        Shifter[T]
		TransformerRaw Transformer[T]
	}

	// Shifter handles tree transformations.
	Shifter[T any] interface {
		// ForEeachInDimension will call the given function for each value in the given dimension d.
		// This is typicall used to e.g. walk all language translations of a given node.
		// If the function returns false, the walk will stop.
		ForEeachInDimension(n T, vec sitesmatrix.Vector, d int, f func(T) bool)

		// ForEeachInAllDimensions will call the given function for each value in all dimensions.
		// If the function returns false, the walk will stop.
		ForEeachInAllDimensions(n T, f func(T) bool)

		// Insert inserts new into the tree into the dimension it provides.
		// It may replace old.
		// It returns the updated and existing T
		// and a bool indicating if an existing record is updated.
		Insert(old, new T) (T, T, bool)

		// Delete deletes T from the given dimension and returns the deleted T and whether the dimension was deleted and if  it's empty after the delete.
		Delete(v T, dimension sitesmatrix.Vector) (T, bool, bool)

		// DeleteFunc deletes nodes in v from the tree where the given function returns true.
		// It returns true if it's empty after the delete.
		DeleteFunc(v T, f func(n T) bool) bool

		// Shift shifts v into the given dimension,
		// if fallback is true, it will fall back a fallback match if found.
		// It returns the shifted T and a bool indicating if the shift was successful.
		Shift(v T, dimension sitesmatrix.Vector, fallback bool) (T, bool)
	}

	Transformer[T any] interface {
		// Append appends vs to t and returns the updated or replaced T and a bool indicating if T was replaced.
		// Note that t may be the zero value and should be ignored.
		Append(t T, ts ...T) (T, bool)
	}
)

// NodeShiftTree is the root of a tree that can be shaped using the Shape method.
// Note that multiplied shapes of the same tree is meant to be used concurrently,
// so use the applicable locking when needed.
type NodeShiftTree[T any] struct {
	tree *radix.Tree

	// [language, version, role].
	siteVector     sitesmatrix.Vector
	shifter        Shifter[T]
	transformerRaw Transformer[T]

	mu *sync.RWMutex
}

func New[T any](cfg Config[T]) *NodeShiftTree[T] {
	if cfg.Shifter == nil {
		panic("Shifter is required")
	}

	return &NodeShiftTree[T]{
		mu: &sync.RWMutex{},

		shifter:        cfg.Shifter,
		transformerRaw: cfg.TransformerRaw,
		tree:           radix.New(),
	}
}

// SiteVector returns the site vector of the current dimension.
func (r *NodeShiftTree[T]) SiteVector() sitesmatrix.Vector {
	return r.siteVector
}

// Delete deletes the node at the given key in the current dimension.
func (r *NodeShiftTree[T]) Delete(key string) (T, bool) {
	return r.delete(key)
}

// DeleteFuncRaw deletes nodes in the tree at the given key where the given function returns true.
// It will delete nodes in all dimensions.
func (r *NodeShiftTree[T]) DeleteFuncRaw(key string, f func(T) bool) (T, int) {
	var count int
	var lastDeleted T

	v, ok := r.tree.Get(key)
	if !ok {
		return lastDeleted, count
	}

	isEmpty := r.shifter.DeleteFunc(v.(T), func(n T) bool {
		if f(n) {
			count++
			lastDeleted = n
			return true
		}
		return false
	})

	if isEmpty {
		r.tree.Delete(key)
	}

	return lastDeleted, count
}

// DeleteRaw deletes the node at the given key without any shifting or transformation.
func (r *NodeShiftTree[T]) DeleteRaw(key string) {
	r.tree.Delete(key)
}

// DeletePrefix deletes all nodes with the given prefix in the current dimension.
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
		deleted, wasDeleted, isEmpty = r.shifter.Delete(v.(T), r.siteVector)
		if isEmpty {
			r.tree.Delete(key)
		}
	}
	return deleted, wasDeleted
}

func (t *NodeShiftTree[T]) DeletePrefixRaw(prefix string) int {
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

// Insert inserts v into the tree at the given key and the
// dimension defined by the value.
// It returns the updated and existing T and a bool indicating if an existing record is updated.
func (r *NodeShiftTree[T]) Insert(s string, v T) (T, T, bool) {
	s = mustValidateKey(cleanKey(s))
	var (
		updated  bool
		existing T
	)
	if vv, ok := r.tree.Get(s); ok {
		v, existing, updated = r.shifter.Insert(vv.(T), v)
	}
	r.insert(s, v)
	return v, existing, updated
}

func (r *NodeShiftTree[T]) insert(s string, v any) (any, bool) {
	if v == nil {
		panic("nil value")
	}
	n, updated := r.tree.Insert(s, v)

	return n, updated
}

// InsertRaw inserts v into the tree at the given key without any shifting or transformation.
func (r *NodeShiftTree[T]) InsertRaw(s string, v any) (any, bool) {
	return r.insert(s, v)
}

func (r *NodeShiftTree[T]) InsertRawWithLock(s string, v any) (any, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.insert(s, v)
}

// It returns the updated and existing T and a bool indicating if an existing record is updated.
func (r *NodeShiftTree[T]) InsertWithLock(s string, v T) (T, T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Insert(s, v)
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
// Note that NodeShiftTree is not thread-safe outside of this transaction construct.
func (t *NodeShiftTree[T]) Lock(writable bool) (commit func()) {
	if writable {
		t.mu.Lock()
	} else {
		t.mu.RLock()
	}

	if writable {
		return t.mu.Unlock
	}
	return t.mu.RUnlock
}

// LongestPrefix finds the longest prefix of s that exists in the tree that also matches the predicate (if set).
// Set exact to true to only match exact in the current dimension (e.g. language).
func (r *NodeShiftTree[T]) LongestPrefix(s string, fallback bool, predicate func(v T) bool) (string, T) {
	for {
		longestPrefix, v, found := r.tree.LongestPrefix(s)

		if found {
			if v, ok := r.shift(v.(T), fallback); ok && (predicate == nil || predicate(v)) {
				return longestPrefix, v
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

// LongestPrefiValueRaw returns the longest prefix and value considering all dimensions
func (r *NodeShiftTree[T]) LongestPrefiValueRaw(s string) (string, T) {
	s, v, found := r.tree.LongestPrefix(s)
	if !found {
		var t T
		return s, t
	}
	return s, v.(T)
}

// LongestPrefixRaw returns the longest prefix considering all dimensions.
func (r *NodeShiftTree[T]) LongestPrefixRaw(s string) (string, bool) {
	s, _, found := r.tree.LongestPrefix(s)
	return s, found
}

// GetRaw returns the raw value at the given key without any shifting or transformation.
func (r *NodeShiftTree[T]) GetRaw(s string) (T, bool) {
	v, ok := r.tree.Get(s)
	if !ok {
		var t T
		return t, false
	}
	return v.(T), true
}

// AppendRaw appends ts to the node at the given key without any shifting or transformation.
func (r *NodeShiftTree[T]) AppendRaw(s string, ts ...T) (T, bool) {
	if r.transformerRaw == nil {
		panic("transformerRaw is required")
	}
	n, found := r.GetRaw(s)
	n2, replaced := r.transformerRaw.Append(n, ts...)
	if replaced || !found {
		r.insert(s, n2)
	}
	return n2, replaced || !found
}

// WalkPrefixRaw walks all nodes with the given prefix in the tree without any shifting or transformation.
func (r *NodeShiftTree[T]) WalkPrefixRaw(prefix string, walker func(key string, value T) bool) {
	walker2 := func(key string, value any) bool {
		return walker(key, value.(T))
	}
	r.tree.WalkPrefix(prefix, walker2)
}

// Shape returns a new NodeShiftTree shaped to the given dimension.
func (t *NodeShiftTree[T]) Shape(v sitesmatrix.Vector) *NodeShiftTree[T] {
	x := t.clone()
	x.siteVector = v
	return x
}

func (t *NodeShiftTree[T]) String() string {
	return fmt.Sprintf("Root{%v}", t.siteVector)
}

func (r *NodeShiftTree[T]) Get(s string) T {
	t, _ := r.get(s)
	return t
}

func (r *NodeShiftTree[T]) ForEeachInDimension(s string, dims sitesmatrix.Vector, d int, f func(T) bool) {
	s = cleanKey(s)
	v, ok := r.tree.Get(s)
	if !ok {
		return
	}
	r.shifter.ForEeachInDimension(v.(T), dims, d, f)
}

func (r *NodeShiftTree[T]) ForEeachInAllDimensions(s string, f func(T) bool) {
	s = cleanKey(s)
	v, ok := r.tree.Get(s)
	if !ok {
		return
	}
	r.shifter.ForEeachInAllDimensions(v.(T), f)
}

type WalkFunc[T any] func(string, T) (bool, error)

//go:generate stringer -type=NodeTransformState
type NodeTransformState int

const (
	NodeTransformStateNone      NodeTransformState = iota
	NodeTransformStateUpdated                      // Node is updated in place.
	NodeTransformStateReplaced                     // Node is replaced and needs to be re-inserted into the tree.
	NodeTransformStateDeleted                      // Node is deleted and should be removed from the tree.
	NodeTransformStateSkip                         // Skip this node, but continue the walk.
	NodeTransformStateTerminate                    // Terminate the walk.
)

type NodeShiftTreeWalker[T any] struct {
	// The tree to walk.
	Tree *NodeShiftTree[T]

	// Transform will be called for each node in the main tree.
	// v2 will replace v1 in the tree.
	// The first bool indicates if the value was replaced and needs to be re-inserted into the tree.
	// the second bool indicates if the walk should skip this node.
	// the third bool indicates if the walk should terminate.
	Transform func(s string, v1 T) (v2 T, state NodeTransformState, err error)

	// When set, will add inserts to WalkContext.HooksPost1 to be performed after the walk.
	TransformDelayInsert bool

	// Handle will be called for each node in the main tree.
	// If the callback returns true, the walk will stop.
	// The callback can optionally return a callback for the nested tree.
	Handle func(s string, v T) (terminate bool, err error)

	// Optional prefix filter.
	Prefix string

	// IncludeFilter is an optional filter that can be used to filter nodes.
	// If it returns false, the node will be skipped.
	// Note that v is the shifted value from the tree.
	IncludeFilter func(s string, v T) bool

	// IncludeRawFilter is an optional filter that can be used to filter nodes.
	// If it returns false, the node will be skipped.
	// Note that v is the raw value from the tree.
	IncludeRawFilter func(s string, v T) bool

	// Enable read or write locking if needed.
	LockType LockType

	// When set, no dimension shifting will be performed.
	NoShift bool

	// When set, will try to fall back to alternative match,
	// typically a shared resource common for all languages.
	Fallback bool

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

// Extend returns a new NodeShiftTreeWalker with the same configuration
// and the same WalkContext as the original.
// Any local state is reset.
func (r *NodeShiftTreeWalker[T]) Extend() *NodeShiftTreeWalker[T] {
	return &NodeShiftTreeWalker[T]{
		Tree:                 r.Tree,
		Transform:            r.Transform,
		TransformDelayInsert: r.TransformDelayInsert,
		Handle:               r.Handle,
		Prefix:               r.Prefix,
		IncludeFilter:        r.IncludeFilter,
		IncludeRawFilter:     r.IncludeRawFilter,
		LockType:             r.LockType,
		NoShift:              r.NoShift,
		Fallback:             r.Fallback,
		Debug:                r.Debug,
		WalkContext:          r.WalkContext,
	}
}

// WithPrefix returns a new NodeShiftTreeWalker with the given prefix.
func (r *NodeShiftTreeWalker[T]) WithPrefix(prefix string) *NodeShiftTreeWalker[T] {
	r2 := r.Extend()
	r2.Prefix = prefix
	return r2
}

// SkipPrefix adds a prefix to be skipped in the walk.
func (r *NodeShiftTreeWalker[T]) SkipPrefix(prefix ...string) {
	r.skipPrefixes = append(r.skipPrefixes, prefix...)
}

// ShouldSkip returns whether the given key should be skipped in the walk.
func (r *NodeShiftTreeWalker[T]) ShouldSkip(s string, v T) bool {
	for _, prefix := range r.skipPrefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	if r.IncludeRawFilter != nil {
		if !r.IncludeRawFilter(s, v) {
			return true
		}
	}
	return false
}

func (r *NodeShiftTreeWalker[T]) Walk(ctx context.Context) error {
	if r.Tree == nil {
		panic("Tree is required")
	}

	var deletes []string

	handleT := func(s string, t T) (ns NodeTransformState, err error) {
		if r.IncludeFilter != nil && !r.IncludeFilter(s, t) {
			return
		}
		if r.Transform != nil {
			if !r.NoShift {
				panic("Transform must be performed with NoShift=true")
			}
			if ns, err = func() (ns NodeTransformState, err error) {
				t, ns, err = r.Transform(s, t)
				if ns >= NodeTransformStateSkip || err != nil {
					return
				}

				switch ns {
				case NodeTransformStateReplaced:
					// Delay insert until after the walk to
					// avoid locking issues.
					if r.TransformDelayInsert {
						r.WalkContext.HooksPost1().Push(
							func() error {
								r.Tree.InsertRaw(s, t)
								return nil
							},
						)
					} else {
						r.Tree.InsertRaw(s, t)
					}
				case NodeTransformStateDeleted:
					// Delay delete until after the walk.
					deletes = append(deletes, s)
					ns = NodeTransformStateSkip
				}
				return
			}(); ns >= NodeTransformStateSkip || err != nil {
				return
			}
		}

		if r.Handle != nil {
			var terminate bool
			terminate, err = r.Handle(s, t)
			if terminate || err != nil {
				return
			}
		}

		return
	}

	return func() error {
		if r.LockType > LockTypeNone {
			unlock := r.Tree.Lock(r.LockType == LockTypeWrite)
			defer unlock()
		}

		r.resetLocalState()

		main := r.Tree

		var err error

		handleV := func(s string, v any) (terminate bool) {
			// Context cancellation check.
			if ctx != nil && ctx.Err() != nil {
				err = ctx.Err()
				return true
			}
			if r.ShouldSkip(s, v.(T)) {
				return false
			}
			var t T
			if r.NoShift {
				t = v.(T)
			} else {
				var ok bool
				t, ok = r.toT(r.Tree, v)
				if !ok {
					return false
				}
			}
			var ns NodeTransformState
			ns, err = handleT(s, t)
			if ns == NodeTransformStateTerminate || err != nil {
				return true
			}
			return false
		}

		if r.Prefix != "" {
			main.tree.WalkPrefix(r.Prefix, handleV)
		} else {
			main.tree.Walk(handleV)
		}

		// This is currently only performed with no shift.
		for _, s := range deletes {
			main.tree.Delete(s)
		}

		return err
	}()
}

func (r *NodeShiftTreeWalker[T]) resetLocalState() {
	r.skipPrefixes = nil
}

func (r *NodeShiftTreeWalker[T]) toT(tree *NodeShiftTree[T], v any) (T, bool) {
	return tree.shift(v.(T), r.Fallback)
}

func (r *NodeShiftTree[T]) Has(s string) bool {
	_, ok := r.get(s)
	return ok
}

func (t NodeShiftTree[T]) clone() *NodeShiftTree[T] {
	return &t
}

func (r *NodeShiftTree[T]) shift(t T, fallback bool) (T, bool) {
	return r.shifter.Shift(t, r.siteVector, fallback)
}

func (r *NodeShiftTree[T]) get(s string) (T, bool) {
	s = cleanKey(s)
	v, ok := r.tree.Get(s)

	if !ok {
		var t T
		return t, false
	}
	if v, ok := r.shift(v.(T), false); ok {
		return v, true
	}

	var t T
	return t, false
}
