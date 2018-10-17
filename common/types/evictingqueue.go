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

// EvictingStringQueue is a queue which automatically evicts elements from the head of
// the queue when attempting to add new elements onto the queue and it is full.
// This queue orders elements LIFO (last-in-first-out). It throws away duplicates.
// Note: This queue currently does not contain any remove (poll etc.) methods.
type EvictingStringQueue struct {
	size int
	vals []string
	set  map[string]bool
	mu   sync.Mutex
}

// NewEvictingStringQueue creates a new queue with the given size.
func NewEvictingStringQueue(size int) *EvictingStringQueue {
	return &EvictingStringQueue{size: size, set: make(map[string]bool)}
}

// Add adds a new string to the tail of the queue if it's not already there.
func (q *EvictingStringQueue) Add(v string) {
	q.mu.Lock()
	if q.set[v] {
		q.mu.Unlock()
		return
	}

	if len(q.set) == q.size {
		// Full
		delete(q.set, q.vals[0])
		q.vals = append(q.vals[:0], q.vals[1:]...)
	}
	q.set[v] = true
	q.vals = append(q.vals, v)
	q.mu.Unlock()
}

// Contains returns whether the queue contains v.
func (q *EvictingStringQueue) Contains(v string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.set[v]
}

// Peek looks at the last element added to the queue.
func (q *EvictingStringQueue) Peek() string {
	q.mu.Lock()
	l := len(q.vals)
	if l == 0 {
		q.mu.Unlock()
		return ""
	}
	elem := q.vals[l-1]
	q.mu.Unlock()
	return elem
}

// PeekAll looks at all the elements in the queue, with the newest first.
func (q *EvictingStringQueue) PeekAll() []string {
	q.mu.Lock()
	vals := make([]string, len(q.vals))
	copy(vals, q.vals)
	q.mu.Unlock()
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}
	return vals
}

// PeekAllSet returns PeekAll as a set.
func (q *EvictingStringQueue) PeekAllSet() map[string]bool {
	all := q.PeekAll()
	set := make(map[string]bool)
	for _, v := range all {
		set[v] = true
	}

	return set
}
