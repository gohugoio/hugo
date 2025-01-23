// Copyright 2017-present The Hugo Authors. All rights reserved.
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

// Package types contains types shared between packages in Hugo.
package types

import (
	"sync"
)

// EvictingQueue is a queue which automatically evicts elements from the head of
// the queue when attempting to add new elements onto the queue and it is full.
// This queue orders elements LIFO (last-in-first-out). It throws away duplicates.
type EvictingQueue[T comparable] struct {
	size int
	vals []T
	set  map[T]bool
	mu   sync.Mutex
	zero T
}

// NewEvictingQueue creates a new queue with the given size.
func NewEvictingQueue[T comparable](size int) *EvictingQueue[T] {
	return &EvictingQueue[T]{size: size, set: make(map[T]bool)}
}

// Add adds a new string to the tail of the queue if it's not already there.
func (q *EvictingQueue[T]) Add(v T) *EvictingQueue[T] {
	q.mu.Lock()
	if q.set[v] {
		q.mu.Unlock()
		return q
	}

	if len(q.set) == q.size {
		// Full
		delete(q.set, q.vals[0])
		q.vals = append(q.vals[:0], q.vals[1:]...)
	}
	q.set[v] = true
	q.vals = append(q.vals, v)
	q.mu.Unlock()

	return q
}

func (q *EvictingQueue[T]) Len() int {
	if q == nil {
		return 0
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.vals)
}

// Contains returns whether the queue contains v.
func (q *EvictingQueue[T]) Contains(v T) bool {
	if q == nil {
		return false
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.set[v]
}

// Peek looks at the last element added to the queue.
func (q *EvictingQueue[T]) Peek() T {
	q.mu.Lock()
	l := len(q.vals)
	if l == 0 {
		q.mu.Unlock()
		return q.zero
	}
	elem := q.vals[l-1]
	q.mu.Unlock()
	return elem
}

// PeekAll looks at all the elements in the queue, with the newest first.
func (q *EvictingQueue[T]) PeekAll() []T {
	if q == nil {
		return nil
	}
	q.mu.Lock()
	vals := make([]T, len(q.vals))
	copy(vals, q.vals)
	q.mu.Unlock()
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}
	return vals
}

// PeekAllSet returns PeekAll as a set.
func (q *EvictingQueue[T]) PeekAllSet() map[T]bool {
	all := q.PeekAll()
	set := make(map[T]bool)
	for _, v := range all {
		set[v] = true
	}

	return set
}
